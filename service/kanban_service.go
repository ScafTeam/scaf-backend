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
		ProjectId:  project_id,
		Todo:       []model.Task{},
		InProgress: []model.Task{},
		Done:       []model.Task{},
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
	}

	var kanban model.Kanban
	dsnap.DataTo(&kanban)

	c.JSON(http.StatusOK, gin.H{
		"message": "list kanban",
		"kanban":  kanban,
	})
}

func AddTask(c *gin.Context) {
	list_type := c.Param("mode")
	if list_type != "Todo" && list_type != "InProgress" && list_type != "Done" {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not Found",
		})
		return
	}

	project_id := c.Param("project_id")
	log.Println("add task")
	req := make(map[string]interface{})
	c.BindJSON(&req)
	log.Println(list_type)
	task_id := shortuuid.New()

	task := model.Task{
		Id:          task_id,
		Name:        req["Name"].(string),
		Description: req["Description"].(string),
	}
	_, err := database.Client.
		Doc("kanbans/"+project_id).
		Update(context.Background(), []firestore.Update{
			{
				Path:  list_type,
				Value: firestore.ArrayUnion(task),
			},
		})
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
	list_type := c.Param("mode")
	if list_type != "Todo" && list_type != "InProgress" && list_type != "Done" {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not Found",
		})
		return
	}
	project_id := c.Param("project_id")
	log.Println(`delete ${list_type} task`)
	req := make(map[string]interface{})
	c.BindJSON(&req)
	task_id := req["TaskId"].(string)

	dsnap, err := database.Client.
		Doc("kanbans/" + project_id).
		Get(context.Background())
	if err != nil {
		log.Printf("An get task array error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}
	tasks := dsnap.Data()[list_type].([]interface{})
	new_tasks := []model.Task{}
	for _, i := range tasks {
		task := model.Task{
			Id:          i.(map[string]interface{})["Id"].(string),
			Name:        i.(map[string]interface{})["Name"].(string),
			Description: i.(map[string]interface{})["Description"].(string),
		}
		if task.Id == task_id {
			continue
		}
		new_tasks = append(new_tasks, task)
	}

	_, err = database.Client.
		Doc("kanbans/"+project_id).
		Update(context.Background(), []firestore.Update{
			{
				Path:  list_type,
				Value: new_tasks,
			},
		})
	if err != nil {
		log.Printf("An error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": `delete ${list_type} task`,
	})
}
