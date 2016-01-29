package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	//	jwt_gin "github.com/gin-gonic/contrib/jwt"
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

// Single secret to do simple signing
var symmetricKey = "symmetrickey"

func postLogin(c *gin.Context) {
	credentials := map[string]string{
		"dolanor": "test",
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
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Couldn't generate the jwt token"})
		}

		// save the tokenstring in the cookiestore (maybe use localstorage?)
		session := sessions.Default(c)
		session.Set("token", tokenString)
		session.Save()

		c.JSON(http.StatusOK, gin.H{"status": "You are logged in"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized access"})
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
		c.JSON(http.StatusExpectationFailed, gin.H{"status": "Couldn't get the token."})
		return
	}

	token, err := jwt.Parse(tokenString, verifyToken)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, gin.H{"status": "Token invalid: " + fmt.Sprintf("%v", err)})
		return
	}

	if token.Claims["auth"].(bool) {
		// Connect to DB service and lookup profile info for token.Claims["name"]
		c.JSON(http.StatusOK, gin.H{"status": "you will see great profile information"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized access"})
	}
}

// password for the cookie store
var cookiesecret = "cookiesecret"

func main() {
	r := gin.Default()
	store := sessions.NewCookieStore([]byte(cookiesecret))

	r.Use(sessions.Sessions("tokens", store))

	r.LoadHTMLGlob("templates/*")

	auth := r.Group("/user")
	auth.GET("/login", displayLogin)
	auth.POST("/login", postLogin)

	profile := r.Group("/user/profile")
	// Could use this if we do request from javascript with custom header
	// Authorize: Bearer <token>
	//profile.Use(jwt_gin.Auth(symmetricKey))
	profile.GET("/", displayProfile)

	r.Run(":8100")
}
