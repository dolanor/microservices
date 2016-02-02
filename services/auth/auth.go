package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dolanor/microservices/api"
	"github.com/dolanor/microservices/services/helper"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
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
		tokenString, err := token.SignedString([]byte(helper.SymmetricKey))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		// save the tokenstring in the cookiestore (maybe use localstorage?)
		session := sessions.Default(c)
		session.Set("token", tokenString)
		session.Save()

		c.Redirect(301, "/user/"+form.Username)
	} else {
		c.Redirect(301, "/login")
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
	return []byte(helper.SymmetricKey), nil
}

func displayProfile(c *gin.Context) {
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
				c.HTML(http.StatusNotFound, "profile.tmpl", gin.H{"title": "Profile", "profile": nil})
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

	var profile api.User
	err = json.Unmarshal(data, &profile)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "profile.tmpl", gin.H{"title": "Profile", "profile": profile})
}

func main() {
	r := gin.Default()
	store := sessions.NewCookieStore([]byte(helper.Cookiesecret))

	r.Use(sessions.Sessions("tokens", store))

	r.LoadHTMLGlob("templates/*")

	r.GET("/login", displayLogin)
	r.POST("/login", postLogin)
	r.GET("/user/:username", displayProfile)

	r.Run(":8100")
}
