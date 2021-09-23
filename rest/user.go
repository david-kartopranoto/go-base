package rest

import (
	"net/http"
	"strconv"

	"github.com/david-kartopranoto/go-base/entity"
	"github.com/david-kartopranoto/go-base/usecase/user"
	"github.com/gin-gonic/gin"
)

const registerQueue = "user-register"

type UserBrokerProvider interface {
	Publish(queue string, body interface{}) error
}

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

func registerUserV2(service UserBrokerProvider) func(c *gin.Context) {
	return func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")
		username := c.PostForm("username")
		err := service.Publish(registerQueue, entity.User{Email: email,
			Password: password,
			Username: username})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

//MakeUserHandlers make url handlers
func MakeUserHandlers(r *gin.Engine, service user.UseCase, broker UserBrokerProvider) {
	r.GET("/v1/user/:id", getUser(service))
	r.GET("/v1/user/list", listUsers(service))
	r.POST("/v1/user/search", searchUser(service))
	r.POST("/v1/user/register", registerUser(service))
	r.POST("/v2/user/register", registerUserV2(broker))
}
