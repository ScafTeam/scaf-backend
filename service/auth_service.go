package service

import (
	"backend/database"
	"backend/model"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ScafTeam/firebase-go-client/auth"
	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	context_json := make(map[string]interface{})
	c.BindJSON(&context_json)
	res := auth.SignUpWithEmailAndPassword(context_json["email"].(string), context_json["password"].(string))
	log.Println(res.Status())
	if res.Status() {
		user := res.Result()
		scaf_user := model.ScafUser{
			Email:    user.Email,
			Projects: []string{},
		}

		jsonStr, err := json.Marshal(scaf_user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err,
			})
			return
		}

		var mapData map[string]interface{}
		if err := json.Unmarshal(jsonStr, &mapData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err,
			})
		}
		_, err = database.Client.
			Doc("users/"+scaf_user.Email).
			Set(context.Background(), mapData)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "authorized",
			"message": "Sign up success",
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
