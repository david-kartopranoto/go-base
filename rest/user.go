package rest

import (
	"net/http"
	"strconv"

	"github.com/david-kartopranoto/go-base/usecase/user"
	"github.com/gin-gonic/gin"
)

func getUser(service user.UseCase) func(c *gin.Context) {
	return func(c *gin.Context) {
		strID := c.Param("id")
		id, err := strconv.ParseInt(strID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		user, err := service.GetUser(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func listUsers(service user.UseCase) func(c *gin.Context) {
	return func(c *gin.Context) {
		users, err := service.ListUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func searchUser(service user.UseCase) func(c *gin.Context) {
	return func(c *gin.Context) {
		query := c.PostForm("query")
		users, err := service.SearchUsers(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func registerUser(service user.UseCase) func(c *gin.Context) {
	return func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")
		username := c.PostForm("username")
		newID, err := service.Register(email, password, username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": newID})
	}
}

//MakeUserHandlers make url handlers
func MakeUserHandlers(r *gin.Engine, service user.UseCase) {
	r.GET("/v1/user/:id", getUser(service))
	r.GET("/v1/user/list", listUsers(service))
	r.POST("/v1/user/search", searchUser(service))
	r.POST("/v1/user/register", registerUser(service))
}
