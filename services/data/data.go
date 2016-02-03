package main

import (
	"github.com/dolanor/microservices/api"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func getUserTodoList(c *gin.Context) {
	todos := map[string]*api.Todo{
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
	users := map[string]*api.User{
		"dolanor": {"dolanor", "Tanguy Herrmann", time.Date(1983, 01, 01, 0, 0, 0, 0, time.UTC)},
	}

	var username string
	c.BindJSON(&username)
	if user, ok := users[c.Param("username")]; ok {
		c.JSON(http.StatusOK, user)
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func main() {
	r := gin.Default()

	data := r.Group("/data")
	data.POST("/user/:username", getUserProfile)
	data.POST("/todo", getUserTodoList)

	r.Run(":8300")
}
