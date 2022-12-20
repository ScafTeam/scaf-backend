package service

import (
	"backend/database"
	"backend/middleware"
	"backend/model"
	"context"
	"net/http"

	"encoding/json"

	"log"
	"time"

	// "github.com/ScafTeam/firebase-go-client/auth"
	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid"
	"google.golang.org/api/iterator"
)

func ListAllProjects(c *gin.Context) {
	userEmail := c.Param("user_email")
	log.Println("get projects")
	iter := database.Client.Collection("projects").
		Where("members", "array-contains", userEmail).
		Documents(context.Background())
	projects := make([]model.Project, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Internal Server Error",
				"message": err.Error(),
			})
			return
		}
		log.Println(doc.Data())
		jsonStr, err := json.Marshal(doc.Data())
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Internal Server Error",
				"message": err.Error(),
			})
			return
		}
		var project model.Project
		if err := json.Unmarshal(jsonStr, &project); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Internal Server Error",
				"message": err.Error(),
			})
			return
		}
		projects = append(projects, project)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   "OK",
		"message":  "Get projects successfully",
		"projects": projects,
	})
}

func CreateProject(c *gin.Context) {
	log.Println("create project")
	claims, _ := middleware.AuthMiddleware.GetClaimsFromJWT(c)

	userEmail := claims[middleware.IdentityKey].(string)

	var req model.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	// check project name without invalid characters
	if !model.CheckProjectName(req.Name) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Project name is invalid",
		})
		return
	}

	// check if project name is unique
	check, err := database.CheckProjectNameUnique(userEmail, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}
	if !check {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Project name is not unique",
		})
		return
	}

	project_uuid := shortuuid.New()
	project_create_on := time.Now().Format(time.RFC850)
	project := map[string]interface{}{
		"id":       project_uuid,
		"name":     req.Name,
		"createOn": project_create_on,
		"author":   userEmail,
		"members":  []string{userEmail},
		"repos":    []model.Repo{},
		"devTools": req.DevTools,
		"devMode":  req.DevMode,
	}

	_, err = database.Client.
		Doc("projects/"+project_uuid).
		Set(context.Background(), project)
	if err != nil {
		log.Printf("An error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	_, err = database.Client.
		Doc("kanbans/"+project_uuid).
		Set(context.Background(), map[string]interface{}{
			"projectId": project_uuid,
			"workflows": []model.Workflow{
				{
					Name:  "Todo",
					Tasks: []model.Task{},
				},
				{
					Name:  "InProgress",
					Tasks: []model.Task{},
				},
				{
					Name:  "Done",
					Tasks: []model.Task{},
				},
			},
		})

	if err != nil {
		log.Printf("An error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "Created",
		"message": "project created",
		"id":      project_uuid,
	})
}

func GetProject(c *gin.Context) {
	project_author := c.Param("user_email")
	project_name := c.Param("project_name")

	log.Println("get project")

	res, err := database.GetProjectDetail(project_author, project_name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	if res == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found",
			"message": "Project not found",
		})
		return
	}

	var project model.Project
	res.DataTo(&project)

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Get project successfully",
		"project": project,
	})
}

func UpdateProject(c *gin.Context) {
	project_name := c.Param("project_name")
	project_author := c.Param("user_email")

	res, err := database.GetProjectDetail(project_author, project_name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	if res == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found",
			"message": "Project not found",
		})
		return
	}

	project_id := res.Ref.ID
	req := model.Project{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	log.Println(project_id)

	_, err = database.Client.
		Doc("projects/"+project_id).
		Update(context.Background(), []firestore.Update{
			{
				Path:  "name",
				Value: req.Name,
			},
			{
				Path:  "devTools",
				Value: req.DevTools,
			},
			{
				Path:  "devMode",
				Value: req.DevMode,
			},
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Update project successfully",
	})
}

func DeleteProject(c *gin.Context) {
	project_author := c.Param("user_email")
	project_name := c.Param("project_name")

	res, err := database.GetProjectDetail(project_author, project_name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	if res == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found",
			"message": "Project not found",
		})
		return
	}

	project_id := res.Ref.ID

	_, err = database.Client.
		Doc("projects/" + project_id).
		Delete(context.Background())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "project deleted",
	})
}

func AddMember(c *gin.Context) {
	var req model.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	project_author := c.Param("user_email")
	project_name := c.Param("project_name")

	res, err := database.GetProjectDetail(project_author, project_name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	if res == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found",
			"message": "Project not found",
		})
		return
	}

	project_id := res.Ref.ID

	_, err = database.Client.
		Doc("projects/"+project_id).
		Update(context.Background(), []firestore.Update{
			{
				Path:  "members",
				Value: firestore.ArrayUnion(req.Email),
			},
		})

	if err != nil {
		log.Printf("An error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "member added",
	})
}
