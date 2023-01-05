package service

import (
	"backend/database"
	"strings"

	"net/http"
	"reflect"

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

func processUpdateData(req interface{}) []firestore.Update {
	update_data := []firestore.Update{}

	req_type := reflect.TypeOf(req)
	req_value := reflect.ValueOf(req)

	for i := 0; i < req_type.NumField(); i++ {
		field := req_type.Field(i)
		value := req_value.Field(i)

		// Check if value is empty, [] is empty
		if !value.IsZero() {
			update_data = append(update_data, firestore.Update{
				Path:  strings.ToLower(field.Name[:1]) + field.Name[1:],
				Value: value.String(),
			})
		}
	}

	return update_data
}

func _processData(data interface{}) map[string]interface{} {
	data_type := reflect.TypeOf(data)
	data_value := reflect.ValueOf(data)

	res := map[string]interface{}{}

	for i := 0; i < data_type.NumField(); i++ {
		field := data_type.Field(i)
		value := data_value.Field(i)

		res[strings.ToLower(field.Name[:1])+field.Name[1:]] = value.String()
	}

	return res
}
