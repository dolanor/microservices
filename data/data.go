package main

import (
	"github.com/dolanor/microservices/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func getUserTodoList(c *gin.Context) {
	todos := models.Todos{
		"dolanor": {"Make the bed", "Eat", "Change the world", "Sleep"},
		"tanguy":  {"Code", "Prepare food", "Read"},
	}

	var username string

	c.BindJSON(&username)
	if todo, ok := todos[username]; ok {
		c.JSON(http.StatusOK, todo)
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
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
	r.POST("/todo", getUserTodoList)

	r.Run(":8300")
}
