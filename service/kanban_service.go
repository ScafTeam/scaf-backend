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

func AddTodoTask(c *gin.Context) {
	log.Println("add task")
	req := make(map[string]interface{})
	c.BindJSON(&req)
	project_id := c.Param("project_id")
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
				Path:  "Todo",
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

func AddInProgressTask(c *gin.Context) {
	log.Println("add task")
	req := make(map[string]interface{})
	c.BindJSON(&req)
	project_id := c.Param("project_id")
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
				Path:  "InProgress",
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
		"message": "add in progress task",
	})
}

func AddDoneTask(c *gin.Context) {
	log.Println("add task")
	req := make(map[string]interface{})
	c.BindJSON(&req)
	project_id := c.Param("project_id")
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
				Path:  "Done",
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
		"message": "add done task",
	})
}

// func DeleteTodoTask(c *gin.Context) {
// 	log.Println("delete Todo task")
// 	req := make(map[string]interface{})
// 	c.BindJSON(&req)
// 	project_id := c.Param("project_id")
// 	task_id := req["TaskId"].(string)

// 	_, err := database.Client.
// 		Doc("kanbans/"+project_id).
// 		Update(context.Background(), []firestore.Update{
// 			{
// 				FieldPath: []string{"Todo"},
// 				Remove:    []interface{}{},
// 			},
// 		})
// 	if err != nil {
// 		log.Printf("An error has occurred: %s", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": err,
// 		})
// 	}
// 	c.JSON(http.StatusCreated, gin.H{
// 		"message": "delete todo task",
// 	})
// }

// func DeleteInProgressTask(c *gin.Context) {
// 	log.Println("delete InProgress task")
// 	req := make(map[string]interface{})
// 	c.BindJSON(&req)
// 	project_id := c.Param("project_id")
// 	task_id := req["TaskId"].(string)

// 	_, err := database.Client.
// 		Doc("kanbans/"+project_id).
// 		Update(context.Background(), []firestore.Update{
// 			{
// 				Path:  "InProgress",
// 				Value: firestore.ArrayRemove(task_id),
// 			},
// 		})
// 	if err != nil {
// 		log.Printf("An error has occurred: %s", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": err,
// 		})
// 	}
// 	c.JSON(http.StatusCreated, gin.H{
// 		"message": "delete inprogress task",
// 	})
// }

// func DeleteDoneTask(c *gin.Context) {
// 	log.Println("delete Done task")
// 	req := make(map[string]interface{})
// 	c.BindJSON(&req)
// 	project_id := c.Param("project_id")
// 	task_id := req["TaskId"].(string)

// 	_, err := database.Client.
// 		Doc("kanbans/"+project_id).
// 		Update(context.Background(), []firestore.Update{
// 			{
// 				Path:  "Done",
// 				Value: firestore.ArrayRemove(task_id),
// 			},
// 		})
// 	if err != nil {
// 		log.Printf("An error has occurred: %s", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": err,
// 		})
// 	}
// 	c.JSON(http.StatusCreated, gin.H{
// 		"message": "delete done task",
// 	})
// }
