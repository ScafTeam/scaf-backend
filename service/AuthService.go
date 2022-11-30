package service

import (
	"github.com/ScafTeam/firebase-go-client/auth"
	"github.com/gin-gonic/gin"
	"log"
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
	} else {
		log.Println(res.ErrorMessage())
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
