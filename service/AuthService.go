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
		log.Println(User.Email + " is login")
		c.JSON(http.StatusOK, gin.H{
			"status":  "active",
			"message": "login success",
		})
	} else {
		log.Println(res.ErrorMessage())
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "active",
			"message": "login success",
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
	} else {
		log.Println(res.ErrorMessage())
	}
}
