package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/dolanor/microservices/api"
	"github.com/dolanor/microservices/services/helper"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func displayTodo(c *gin.Context) {
	data, err := helper.QueryDataService(c)
	if err != nil {
		switch e := err.(type) {
		case *jwt.ValidationError:
			if e.Errors&jwt.ValidationErrorExpired == jwt.ValidationErrorExpired {
				c.Redirect(http.StatusTemporaryRedirect, "/login")
				return
			}
		case error:
			switch err {
			case api.ErrUnauthorized:
				c.Redirect(http.StatusTemporaryRedirect, "/login")
				return
			case api.ErrConnectingEndpoint:
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			case api.ErrDataNotFound:
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

	var todo api.Todo
	err = json.Unmarshal(data, &todo)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "todo.tmpl", gin.H{"title": "TODO", "todo": todo})
}

func main() {
	r := gin.Default()
	store := sessions.NewCookieStore([]byte(helper.Cookiesecret))

	r.Use(sessions.Sessions("tokens", store))

	r.LoadHTMLGlob("templates/*")

	r.GET("/todo", displayTodo)

	r.Run(":8200")
}
