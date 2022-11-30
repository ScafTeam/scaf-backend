package service

import (
	"log"
	"net/http"

	"github.com/ScafTeam/firebase-go-client/auth"
	"github.com/gin-gonic/gin"
)

var User = auth.User{}

func UserLogin(c *gin.Context) {
	json := make(map[string]interface{})
	c.BindJSON(&json)
	// log.Println(json)
	res := auth.SignInWithEmailAndPassword(json["email"].(string), json["password"].(string))
	if res.Status() {
		User = res.Result()
		log.Println(User.Email + " is signed in")
		c.JSON(http.StatusOK, gin.H{
			"status":  "authorized",
			"message": "Sign in success",
		})
	} else {
		log.Println(res.ErrorMessage())
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "unauthorized",
			"message": res.ErrorMessage(),
		})
	}
}

func UserRegister(c *gin.Context) {
	json := make(map[string]interface{})
	c.BindJSON(&json)
	res := auth.SignUpWithEmailAndPassword(json["email"].(string), json["password"].(string))
	log.Println(res.Status())
	if res.Status() {
		user := res.Result()
		log.Println(user.Email + " is signed up")
		c.JSON(http.StatusOK, gin.H{
			"status":  "authorized",
			"message": "Sign in success",
		})
	} else {
		log.Println(res.ErrorMessage())
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "unauthorized",
			"message": res.ErrorMessage(),
		})
	}
}

func UserForgotPassword(c *gin.Context) {
	json := make(map[string]interface{})
	c.BindJSON(&json)
	res := auth.ForgotPassword(json["email"].(string))
	if res.Status() {
		log.Println("Email is sent")
		log.Println(res)
		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": "Password reset email sent",
		})
	} else {
	// EMAIL_NOT_FOUND 沒有此用戶
		log.Println(res.ErrorMessage())
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "eamil not found",
			"message": res.ErrorMessage(),
		})
	}
}
