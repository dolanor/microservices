package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dolanor/microservices/api"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
)

// GenResponse generate the body of the response in HTML or JSON given the Accept header sent by the client.
// Makes it easier to get JSON without changing the api path.
func GenResponse(c *gin.Context, httpStatusCode int, template string, ginH map[string]interface{}) {
	data := ginH["data"]

	if c.Request.Header.Get("Accept") == "application/json" {
		c.JSON(httpStatusCode, data)
		return
	}
	c.HTML(httpStatusCode, template, ginH)
}

// QueryDataService is a helper function to access another microservice endpoint.
// Currently, works for the DB Accessor. Might be extended for other services.
func QueryDataService(c *gin.Context) ([]byte, error) {
	session := sessions.Default(c)
	tokenString, ok := session.Get("token").(string)
	if !ok {
		return []byte{}, api.ErrTokenNotFound
	}

	token, err := jwt.Parse(tokenString, verifyToken)
	if err != nil {
		return []byte{}, err
	}

	if !token.Claims["auth"].(bool) {
		return []byte{}, api.ErrUnauthorized
	}

	// Connect to DB service and lookup profile info for token.Claims["name"]
	client := &http.Client{}
	r := c.Request
	u, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		return []byte{}, err
	}

	var username string
	if username, ok = token.Claims["name"].(string); !ok {
		return []byte{}, api.ErrTokenItemNotFound
	}

	if c.Param("username") != "" && c.Param("username") != username {
		return []byte{}, api.ErrUnauthorized
	}

	jsonusername, err := json.Marshal(username)
	if err != nil {
		return []byte{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", DataServiceURL, u.Path), bytes.NewBuffer(jsonusername))
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, api.ErrConnectingEndpoint
	}
	defer resp.Body.Close()

	// If the data services doesn't have any profile associated with that user
	if resp.StatusCode == 404 {
		return []byte{}, api.ErrDataNotFound
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}
