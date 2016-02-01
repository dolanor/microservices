package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dolanor/microservices/errors"
	"github.com/dolanor/microservices/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	// Single secret to do simple token signing
	symmetricKey = "symmetrickey"
	// password for the cookie store
	cookiesecret = "cookiesecret"
	// url to the DB Accessor. Need to replace with service discovery.
	dataServiceURL = "http://localhost:8300"
)

func verifyToken(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected Signing method: %v", token.Header["alg"])
	}
	// return the key used for signing the token
	return []byte(symmetricKey), nil
}

func queryDataService(c *gin.Context) ([]byte, error) {
	session := sessions.Default(c)
	tokenString, ok := session.Get("token").(string)
	if !ok {
		return []byte{}, errors.ErrTokenNotFound
	}

	token, err := jwt.Parse(tokenString, verifyToken)
	if err != nil {
		return []byte{}, err
	}

	if token.Claims["auth"].(bool) {
		// Connect to DB service and lookup profile info for token.Claims["name"]
		client := &http.Client{}
		r := c.Request
		u, err := url.ParseRequestURI(r.RequestURI)
		if err != nil {
			return []byte{}, err
		}

		var username string
		if username, ok = token.Claims["name"].(string); !ok {
			return []byte{}, errors.ErrTokenItemNotFound
		}

		jsonusername, err := json.Marshal(username)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return []byte{}, err
		}

		req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", dataServiceURL, u.Path), bytes.NewBuffer(jsonusername))
		if err != nil {
			return []byte{}, err
		}

		req.Header.Set("content-type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			return []byte{}, errors.ErrConnectingEndpoint
		}
		defer resp.Body.Close()

		// If the data services doesn't have any profile associated with that user
		if resp.StatusCode == 404 {
			return []byte{}, errors.ErrDataNotFound
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []byte{}, err
		}
		return body, nil
	} else {
		return []byte{}, errors.ErrUnauthorized
	}
}

func displayTodo(c *gin.Context) {
	data, err := queryDataService(c)
	if err != nil {
		switch e := err.(type) {
		case jwt.ValidationError:
			if e.Errors&jwt.ValidationErrorExpired == jwt.ValidationErrorExpired {
				c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			}
		case error:
			switch err {
			case errors.ErrUnauthorized:
				c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			case errors.ErrConnectingEndpoint:
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			case errors.ErrDataNotFound:
				c.HTML(http.StatusOK, "todo.tmpl", gin.H{"title": "TODO", "todo": nil})
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
	store := sessions.NewCookieStore([]byte(cookiesecret))

	r.Use(sessions.Sessions("tokens", store))

	r.LoadHTMLGlob("templates/*")

	r.GET("/todo", displayTodo)

	r.Run(":8200")
}
