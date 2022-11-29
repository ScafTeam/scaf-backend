package router

import (
	"backend/service"
	"github.com/gin-gonic/gin"
)

func AddAuthRouter(r *gin.RouterGroup) {
	auth_router := r.Group("/")
	auth_router.POST("/login", service.UserLogin)
	auth_router.POST("/register", service.UserRegister)
}
