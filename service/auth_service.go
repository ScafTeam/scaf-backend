package service

import (
	"backend/database"
	"backend/model"
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
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
		// user := res.Result()
		scaf_user := map[string]interface{}{
			"email":    req.Email,
			"avatar":   "",
			"bio":      "",
			"nickname": "",
			"projects": []string{},
		}

		_, err := database.Client.
			Doc("users/"+req.Email).
			Set(context.Background(), scaf_user)

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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BadRequst",
			"message": res.ErrorMessage(),
		})
		return
	}
}

func GetUserData(c *gin.Context) {
	email := c.Param("user_email")

	doc, err := database.Client.Doc("users/" + email).Get(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	var user model.ScafUser
	doc.DataTo(&user)

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"data":   user,
	})
}

func UpdateUserData(c *gin.Context) {
	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BadRequst",
			"message": err.Error(),
		})
		return
	}

	email := c.Param("user_email")

	// if field is empty, don't update
	updateData := processUpdateData(req)

	_, err := database.Client.
		Doc("users/"+email).
		Update(context.Background(), updateData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Update user data success",
	})
}

func UpdateUserPassword(c *gin.Context) {
	var req model.UpdateUserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BadRequst",
			"message": err.Error(),
		})
		return
	}

	email := c.Param("user_email")

	res := auth.SignInWithEmailAndPassword(email, req.OldPassword)
	if !res.Status() {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": res.ErrorMessage(),
		})
		return
	}

	user := res.Result()

	res2 := auth.UpdatePassword(user, req.NewPassword)
	if res2.Status() {
		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": "Update password success",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BadRequst",
			"message": res2.ErrorMessage(),
		})
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

func addUserProjects(c *gin.Context, email, projectName string) bool {
	_, err := database.Client.
		Doc("users/"+email).
		Update(context.Background(), []firestore.Update{
			{
				Path:  "projects",
				Value: firestore.ArrayUnion(projectName),
			},
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		c.Abort()
		return false
	}

	return true
}

func removeUserProjects(c *gin.Context, email, projectName string) bool {
	_, err := database.Client.
		Doc("users/"+email).
		Update(context.Background(), []firestore.Update{
			{
				Path:  "projects",
				Value: firestore.ArrayRemove(projectName),
			},
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		c.Abort()
		return false
	}

	return true
}
