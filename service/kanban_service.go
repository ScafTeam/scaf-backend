package service

import (
	"backend/database"
	"backend/model"
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid"
)

func CreateKanban(c *gin.Context) {
	log.Println("create kanban")
	project_id := c.Param("project_id")

	kanban := model.Kanban{
		ProjectId: project_id,
		Workflows: []model.Workflow{
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
	}

	_, err := database.Client.
		Collection("kanbans").
		Doc(project_id).
		Set(context.Background(), kanban)

	if err != nil {
		log.Printf("An error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "create kanban",
	})
}

func ListKanban(c *gin.Context) {
	log.Println("list kanban")
	project_id := c.Param("project_id")

	dsnap, err := database.Client.
		Doc("kanbans/" + project_id).
		Get(context.Background())

	if err != nil {
		log.Printf("An error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})

		return
	}

	var kanban model.Kanban
	dsnap.DataTo(&kanban)

	c.JSON(http.StatusOK, gin.H{
		"message": "list kanban",
		"kanban":  kanban,
	})
}

func AddWorkFlow(c *gin.Context) {
	log.Println("Add workflow")
	project_id := c.Param("project_id")
	req := make(map[string]interface{})
	c.BindJSON(&req)

	workflow := model.Workflow{
		Name:  req["Name"].(string),
		Tasks: []model.Task{},
	}

	_, err := database.Client.
		Doc("kanbans/"+project_id).
		Update(context.Background(), []firestore.Update{
			{
				Path:  "Workflows",
				Value: firestore.ArrayUnion(workflow),
			},
		})

	if err != nil {
		log.Printf("An error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Add workflow success",
	})
}

func AddTask(c *gin.Context) {
	project_id := c.Param("project_id")
	log.Println("add task")
	req := make(map[string]interface{})
	c.BindJSON(&req)
	task_id := shortuuid.New()

	workflow_name := req["WorkflowName"].(string)

	// check workflow_name in firebase kanbans
	dsnap, err := database.Client.
		Doc("kanbans/" + project_id).
		Get(context.Background())

	if err != nil {
		log.Printf("An error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	var kanban model.Kanban
	dsnap.DataTo(&kanban)

	var workflow_index int
	for idx, wf := range kanban.Workflows {
		if wf.Name == workflow_name {
			workflow_index = idx
		}
	}

	if kanban.Workflows[workflow_index].Name == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "workflow not found",
		})
		return
	}

	task := model.Task{
		Id:          task_id,
		Name:        req["Name"].(string),
		Description: req["Description"].(string),
	}

	kanban.Workflows[workflow_index].Tasks = append(kanban.Workflows[workflow_index].Tasks, task)

	_, err = database.Client.
		Doc("kanbans/"+project_id).
		Set(context.Background(), kanban)

	if err != nil {
		log.Printf("An error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "add todo task",
	})
}

func DeleteTask(c *gin.Context) {
	project_id := c.Param("project_id")
	req := make(map[string]interface{})
	c.BindJSON(&req)
	task_id := req["TaskId"].(string)

	workflow_name := req["WorkflowName"].(string)

	dsnap, err := database.Client.
		Doc("kanbans/" + project_id).
		Get(context.Background())

	if err != nil {
		log.Printf("An get task array error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})

		return
	}

	var kanban model.Kanban
	dsnap.DataTo(&kanban)

	var workflow_index int
	for idx, wf := range kanban.Workflows {
		if wf.Name == workflow_name {
			workflow_index = idx
		}
	}

	if kanban.Workflows[workflow_index].Name == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "workflow not found",
		})
		return
	}

	var task_index int
	for idx, task := range kanban.Workflows[workflow_index].Tasks {
		if task.Id == task_id {
			task_index = idx
		}
	}

	if kanban.Workflows[workflow_index].Tasks[task_index].Id == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "task not found",
		})
		return
	}

	kanban.Workflows[workflow_index].Tasks = append(kanban.Workflows[workflow_index].Tasks[:task_index], kanban.Workflows[workflow_index].Tasks[task_index+1:]...)

	_, err = database.Client.
		Doc("kanbans/"+project_id).
		Set(context.Background(), kanban)

	if err != nil {
		log.Printf("An error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": `delete ${workflow_name} task id = ${task_id} success`,
	})
}
