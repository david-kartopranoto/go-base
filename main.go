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
	config, err := util.LoadConfig("./config")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	router := gin.Default()

	conn := initDB(config)
	userRepo := repository.NewUserSQL(conn)
	userService := user.NewService(userRepo)

	rest.MakeUserHandlers(router, userService)

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to %s", c.FullPath())
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
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
	log.Println("Database connection established")

	return conn
}
