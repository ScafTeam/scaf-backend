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

	var req model.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BadRequst",
			"message": err.Error(),
		})
		return
	}

	res := auth.SignUpWithEmailAndPassword(req.Email, req.Password)
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
				"status":  "Internal Server Error",
				"message": err.Error(),
			})
			return
		}

		var mapData map[string]interface{}
		if err := json.Unmarshal(jsonStr, &mapData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Internal Server Error",
				"message": err.Error(),
			})
			return
		}
		_, err = database.Client.
			Doc("users/"+scaf_user.Email).
			Set(context.Background(), mapData)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Internal Server Error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "Authorized",
			"message": "Sign up success",
		})
		return
	} else {
		log.Println(res.ErrorMessage())
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BadRequst",
			"message": res.ErrorMessage(),
		})
		return
	}
}

func UserForgotPassword(c *gin.Context) {

	var req model.UserForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BadRequst",
			"message": err.Error(),
		})
		return
	}

	res := auth.ForgotPassword(req.Email)
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
			"status":  "Eamil Not Found",
			"message": res.ErrorMessage(),
		})
	}
}
