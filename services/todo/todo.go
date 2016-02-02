package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/dolanor/microservices/errors"
	"github.com/dolanor/microservices/models"
	"github.com/dolanor/microservices/server"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func displayTodo(c *gin.Context) {
	data, err := server.QueryDataService(c)
	if err != nil {
		switch e := err.(type) {
		case *jwt.ValidationError:
			if e.Errors&jwt.ValidationErrorExpired == jwt.ValidationErrorExpired {
				c.Redirect(http.StatusTemporaryRedirect, "/user/login")
				return
			}
		case error:
			switch err {
			case errors.ErrUnauthorized:
				c.Redirect(http.StatusTemporaryRedirect, "/user/login")
				return
			case errors.ErrConnectingEndpoint:
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			case errors.ErrDataNotFound:
				c.HTML(http.StatusNotFound, "todo.tmpl", gin.H{"title": "TODO", "todo": nil})
				return
			default:
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		default:
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	var todo models.Todo
	err = json.Unmarshal(data, &todo)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "todo.tmpl", gin.H{"title": "TODO", "todo": todo})
}

func main() {
	r := gin.Default()
	store := sessions.NewCookieStore([]byte(server.Cookiesecret))

	r.Use(sessions.Sessions("tokens", store))

	r.LoadHTMLGlob("templates/*")

	r.GET("/todo", displayTodo)

	r.Run(":8200")
}
