package main

import (
    "net/http"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

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