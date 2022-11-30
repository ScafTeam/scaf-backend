package router

import (
	"backend/service"
	"github.com/gin-gonic/gin"
)

func AddAuthRouter(auth_router *gin.RouterGroup) {
	auth_router.POST("/login", service.UserLogin)
	auth_router.POST("/register", service.UserRegister)
}
