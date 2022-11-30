package router

import (
	"backend/service"
	"github.com/gin-gonic/gin"
	"log"
)

func AddProjectRouter(project_router *gin.RouterGroup) {
	log.Println(service.User.Email)
	if service.User.Email == "" {
		return
	}
	project_router.GET("/list", service.ListAllProjects)
}
