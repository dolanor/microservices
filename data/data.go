package main

import (
	"github.com/dolanor/microservices/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func getUserTodoList(c *gin.Context) {

}

func getUserProfile(c *gin.Context) {
	users := map[string]*models.UserProfile{
		"dolanor": {"dolanor", "Tanguy Herrmann", time.Date(1983, 01, 01, 0, 0, 0, 0, time.UTC)},
	}

	var profile models.UserProfile
	c.BindJSON(&profile)
	if user, ok := users[profile.Username]; ok {
		c.JSON(http.StatusOK, user)
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func main() {
	r := gin.Default()

	r.POST("/user/profile", getUserProfile)
	r.GET("/todo/", getUserTodoList)

	r.Run(":8300")
}
