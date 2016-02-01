package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dolanor/microservices/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Login binds the data from an HTML form to gin data binding system
type Login struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func displayLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login_form.tmpl", gin.H{"title": "Log in"})
}

var (
	// Single secret to do simple token signing
	symmetricKey = "symmetrickey"
	// password for the cookie store
	cookiesecret = "cookiesecret"
	// url to the DB Accessor. Need to replace with service discovery.
	dataServiceURL = "http://localhost:8300"
)

func postLogin(c *gin.Context) {
	// Temporary storage
	credentials := map[string]string{
		"dolanor": "test",
		"tanguy":  "pass",
	}

	var form Login

	if c.Bind(&form) != nil {
		return
	}

	if password, ok := credentials[form.Username]; ok && password == form.Password {
		// generate the jwt
		token := jwt.New(jwt.GetSigningMethod("HS256"))
		token.Claims["name"] = form.Username
		token.Claims["exp"] = time.Now().Add(1 * time.Minute).Unix()
		token.Claims["auth"] = true

		// sign the token
		tokenString, err := token.SignedString([]byte(symmetricKey))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		// save the tokenstring in the cookiestore (maybe use localstorage?)
		session := sessions.Default(c)
		session.Set("token", tokenString)
		session.Save()

		c.Redirect(301, "/user/profile")
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

type authTokenData struct {
	username      string
	authenticated bool
}

func verifyToken(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected Signing method: %v", token.Header["alg"])
	}
	// return the key used for signing the token
	return []byte(symmetricKey), nil
}

func displayProfile(c *gin.Context) {
	session := sessions.Default(c)
	tokenString, ok := session.Get("token").(string)
	if !ok {
		c.AbortWithError(http.StatusUnauthorized, errors.New("Token not found in sessions cookiestore"))
		return
	}

	token, err := jwt.Parse(tokenString, verifyToken)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		c.Redirect(http.StatusTemporaryRedirect, "/user/login")
		return
	}

	if token.Claims["auth"].(bool) {
		// Connect to DB service and lookup profile info for token.Claims["name"]
		client := &http.Client{}
		r := c.Request
		u, err := url.ParseRequestURI(r.RequestURI)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		user := models.UserProfile{Username: token.Claims["name"].(string)}
		jsonuser, err := json.Marshal(user)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", dataServiceURL, u.Path), bytes.NewBuffer(jsonuser))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		req.Header.Set("content-type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer resp.Body.Close()

		// If the data services doesn't have any profile associated with that
		// user
		if resp.StatusCode == 404 {
			//c.AbortWithStatus(resp.StatusCode)
			c.HTML(http.StatusOK, "profile.tmpl", gin.H{"title": "Profile", "profile": nil})
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var profile models.UserProfile
		err = json.Unmarshal(body, &profile)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.HTML(http.StatusOK, "profile.tmpl", gin.H{"title": "Profile", "profile": profile})
	} else {
		c.Redirect(http.StatusTemporaryRedirect, "/user/login")
	}
}

func main() {
	r := gin.Default()
	store := sessions.NewCookieStore([]byte(cookiesecret))

	r.Use(sessions.Sessions("tokens", store))

	r.LoadHTMLGlob("templates/*")

	auth := r.Group("/user")
	auth.GET("/login", displayLogin)
	auth.POST("/login", postLogin)

	r.GET("/user/profile", displayProfile)

	r.Run(":8100")
}
