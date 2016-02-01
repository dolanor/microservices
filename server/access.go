package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dolanor/microservices/errors"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
)

func QueryDataService(c *gin.Context) ([]byte, error) {
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

		req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", DataServiceURL, u.Path), bytes.NewBuffer(jsonusername))
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
