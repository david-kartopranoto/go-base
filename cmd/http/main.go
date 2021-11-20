package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/david-kartopranoto/go-base/repository"
	"github.com/david-kartopranoto/go-base/rest"
	"github.com/david-kartopranoto/go-base/usecase/user"
	"github.com/david-kartopranoto/go-base/util"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig("./config", "app")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	metricService, err := util.NewPrometheusService()
	if err != nil {
		log.Fatal(err.Error())
	}

	limiterService, err := util.NewLimiterService(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	brokerService, err := util.NewRabbitMQService(config, metricService)
	if err != nil {
		log.Fatal(err.Error())
	}

	authService, err := util.NewAuthService(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	router := gin.Default()

	conn := initDB(config)
	userRepo := repository.NewUserSQL(conn)
	userService := user.NewService(userRepo)

	router.Use(rest.HistogramMiddleware(metricService))
	router.Use(rest.LimiterMiddleware(limiterService))
	rest.MakeUserHandlers(router, userService, brokerService, authService)
	rest.MakeMetricsHandlers(router, metricService)
	rest.MakeAuthHandlers(router, authService)

	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"DB": conn.Stats(),
		})
	})

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
