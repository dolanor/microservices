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

// Login binds the data from the HTML form to gin data binding system
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
			helper.GenResponse(c, http.StatusBadRequest, "login_form.tmpl", gin.H{"title": "Log in"})
		}

		// save the tokenstring in the cookiestore (maybe use localstorage?)
		session := sessions.Default(c)
		session.Set("token", tokenString)
		session.Save()

		helper.GenResponse(c, http.StatusOK, "login_form.tmpl", gin.H{"title": "Log in", "data": form.Username, "username": form.Username})
	} else {
		helper.GenResponse(c, http.StatusUnauthorized, "login_form.tmpl", gin.H{"title": "Log in"})
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
				helper.GenResponse(c, http.StatusUnauthorized, "profile.tmpl", gin.H{"title": "Profile", "data": nil})
				return
			}
		case error:
			switch err {
			case api.ErrUnauthorized:
				helper.GenResponse(c, http.StatusUnauthorized, "profile.tmpl", gin.H{"title": "Profile", "data": nil})
				return
			case api.ErrConnectingEndpoint:
				helper.GenResponse(c, http.StatusServiceUnavailable, "profile.tmpl", gin.H{"title": "Profile", "data": nil})
				return
			case api.ErrDataNotFound:
				helper.GenResponse(c, http.StatusNotFound, "profile.tmpl", gin.H{"title": "Profile", "data": nil})
				return
			default:
				helper.GenResponse(c, http.StatusInternalServerError, "profile.tmpl", gin.H{"title": "Profile", "data": nil})
				return
			}
		default:
			helper.GenResponse(c, http.StatusInternalServerError, "profile.tmpl", gin.H{"title": "Profile", "data": nil})
			return
		}
	}

	var profile api.User
	err = json.Unmarshal(data, &profile)
	if err != nil {
		helper.GenResponse(c, http.StatusBadRequest, "profile.tmpl", gin.H{"title": "Profile", "data": nil})
		return
	}

	// If we're here, we can get these informations already without errors
	token, _ := helper.GetTokenFromContext(c)
	username, _ := helper.GetUsernameFromToken(token)
	helper.GenResponse(c, http.StatusOK, "profile.tmpl", gin.H{"title": "Profile", "data": profile, "username": username})
}

func main() {
	r := gin.Default()
	store := sessions.NewCookieStore([]byte(helper.Cookiesecret))

	r.Use(sessions.Sessions("todo_session", store))

	r.LoadHTMLGlob("templates/*")

	r.GET("/login", displayLogin)
	r.POST("/login", postLogin)
	r.GET("/user/:username", displayProfile)

	r.Run(":8100")
}
