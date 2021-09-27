package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/david-kartopranoto/go-base/entity"
	"github.com/david-kartopranoto/go-base/repository"
	"github.com/david-kartopranoto/go-base/rest"
	"github.com/david-kartopranoto/go-base/usecase/user"
	"github.com/david-kartopranoto/go-base/util"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig("./config", "consumer")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	metricService, err := util.NewPrometheusService()
	if err != nil {
		log.Fatal(err.Error())
	}

	brokerService, err := util.NewRabbitMQService(config, metricService)
	if err != nil {
		log.Fatal(err.Error())
	}

	conn := initDB(config)
	userRepo := repository.NewUserSQL(conn)
	userService := user.NewService(userRepo)

	stopChan := make(chan bool)

	brokerService.Consume(entity.RegisterQueue, stopChan, userService.ConsumeRegister)

	router := gin.Default()

	rest.MakeMetricsHandlers(router, metricService)

	router.Run()
}

func initDB(config util.Config) *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Password,
		config.DB.Name,
		config.DB.SSLMode)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Println("Open", err)
	}
	err = conn.Ping()
	if err != nil {
		log.Println("Ping", err)
	}

	return conn
}
