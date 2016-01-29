package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Login struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func displayLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login_form.tmpl", gin.H{"title": "Log in"})
}

func postLogin(c *gin.Context) {
	credentials := map[string]string{
		"dolanor": "test",
	}

	var form Login

	if c.Bind(&form) != nil {
		return
	}

	if password, ok := credentials[form.Username]; ok && password == form.Password {
		c.JSON(http.StatusOK, gin.H{"status": "You are logged in"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized access"})
	}
}

func displayProfile(c *gin.Context) {
}

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	auth := r.Group("/user")
	auth.GET("/login", displayLogin)
	auth.POST("/login", postLogin)
	auth.GET("/", displayProfile)

	r.Run(":8100")
}
