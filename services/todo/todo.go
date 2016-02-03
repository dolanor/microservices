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
				helper.GenResponse(c, http.StatusUnauthorized, "todo.tmpl", gin.H{"title": "TODO", "data": nil})
				return
			}
		case error:
			switch err {
			case api.ErrUnauthorized:
				helper.GenResponse(c, http.StatusUnauthorized, "todo.tmpl", gin.H{"title": "TODO", "data": nil})
				return
			case api.ErrConnectingEndpoint:
				helper.GenResponse(c, http.StatusServiceUnavailable, "todo.tmpl", gin.H{"title": "TODO", "data": nil})
				return
			case api.ErrDataNotFound:
				helper.GenResponse(c, http.StatusNotFound, "todo.tmpl", gin.H{"title": "TODO", "data": nil})
				return
			default:
				helper.GenResponse(c, http.StatusInternalServerError, "todo.tmpl", gin.H{"title": "TODO", "data": nil})
				return
			}
		default:
			helper.GenResponse(c, http.StatusInternalServerError, "todo.tmpl", gin.H{"title": "TODO", "data": nil})
			return
		}
	}

	var todo api.Todo
	err = json.Unmarshal(data, &todo)
	if err != nil {
		helper.GenResponse(c, http.StatusBadRequest, "todo.tmpl", gin.H{"title": "TODO", "data": nil})
		return
	}
	// If we're here, we can get these informations already without errors
	token, _ := helper.GetTokenFromContext(c)
	username, _ := helper.GetUsernameFromToken(token)
	helper.GenResponse(c, http.StatusOK, "todo.tmpl", gin.H{"title": "TODO", "data": todo, "username": username})
}

func main() {
	r := gin.Default()
	store := sessions.NewCookieStore([]byte(helper.Cookiesecret))

	r.Use(sessions.Sessions("todo_session", store))

	r.LoadHTMLGlob("templates/*")

	r.GET("/todo", displayTodo)

	r.Run(":8200")
}
