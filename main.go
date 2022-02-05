package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Data struct {
	Message string
	Page    string
	Action  string
}

var database = map[string]string{}

func openLoginPage(c *gin.Context) {
	data := Data{
		Message: "Login Below",
		Page:    "Login Page",
		Action:  "Login",
	}
	c.HTML(http.StatusOK, "login.html", gin.H{
		"Message": data.Message,
		"Page":    data.Page,
		"Action":  data.Action,
	})
}

func openSignUpPage(c *gin.Context) {
	data := Data{
		Message: "Sign Up Below",
		Page:    "Sign Up Page",
		Action:  "Sign Up"}
	c.HTML(http.StatusOK, "login.html", gin.H{"Message": data.Message,
		"Action": data.Action,
		"Page":   data.Page,
	})
}

func signUp(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	data := Data{
		Message: "Sign Up Failed: Username already exists",
		Page:    "Failure",
		Action:  "Sign-up"}
	for k := range database {
		if k == username {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{"Message": data.Message,
				"Action": data.Action,
				"Page":   data.Page})
			return
		}
	}
	// get the hash of the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	data = Data{
		Message: "Sign Up Failed: Failed to generate password",
		Page:    "Failure",
		Action:  "Sign-up"}
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"Message": data.Message,
			"Action": data.Action,
			"Page":   data.Page})
		return
	}
	// store in db
	database[username] = string(passwordHash)
	data = Data{
		Message: fmt.Sprintf("Sign Up Successful, Welcome %v", username),
		Page:    "Success",
		Action:  "None",
	}
	c.HTML(http.StatusOK, "login.html", gin.H{"Message": data.Message,
		"Action": data.Action,
		"Page":   data.Page})
}

func loginUser(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	data := Data{
		Message: "Login Failed: Invalid Password",
		Page:    "Failure",
		Action:  "Login"}
	// validate user
	if err := bcrypt.CompareHashAndPassword([]byte(database[username]), []byte(password)); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"Message": data.Message,
			"Action": data.Action,
			"Page":   data.Page})
		return
	}
	data = Data{
		Message: "Login Successful",
		Page:    "Success",
		Action:  "None",
	}
	c.HTML(http.StatusOK, "login.html", gin.H{"Message": data.Message,
		"Action": data.Action,
		"Page":   data.Page})
}

func logoutUser(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "/login")
}

func main() {
	ginObj := gin.Default()
	ginObj.LoadHTMLFiles("login.html")
	ginObj.GET("/login", openLoginPage)
	ginObj.GET("/signup", openSignUpPage)
	ginObj.POST("/signup", signUp)
	ginObj.POST("/login", loginUser)
	ginObj.GET("/logout", logoutUser)
	ginObj.Run(":8000")
}
