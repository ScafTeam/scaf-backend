package service

import (
	"backend/database"

	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func getProjectDetail(c *gin.Context) *firestore.DocumentSnapshot {
	project_name := c.Param("project_name")
	user_email := c.Param("user_email")

	res, err := database.GetProjectDetail(user_email, project_name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		c.Abort()
		return nil
	}

	if res == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found",
			"message": "Project not found",
		})
		c.Abort()
		return nil
	}

	return res
}

func getProjectId(c *gin.Context) string {
	res := getProjectDetail(c)

	return res.Ref.ID
}
