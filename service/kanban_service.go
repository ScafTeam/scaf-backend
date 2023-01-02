package service

import (
	"backend/database"
	"backend/model"
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid"
	"google.golang.org/api/iterator"
)

func createKanban(c *gin.Context) []map[string]interface{} {
	todo_id := shortuuid.New()
	inProgress_id := shortuuid.New()
	done_id := shortuuid.New()

	todo := map[string]interface{}{
		"id":   todo_id,
		"name": "To Do",
	}

	inProgress := map[string]interface{}{
		"id":   inProgress_id,
		"name": "In Progress",
	}

	done := map[string]interface{}{
		"id":   done_id,
		"name": "Done",
	}

	return []map[string]interface{}{todo, inProgress, done}
}

func setKanbanToFirestore(c *gin.Context, project_id string) bool {
	_, err := database.Client.
		Doc("kanbans/"+project_id).
		Set(context.Background(), map[string]interface{}{
			"projectId": project_id,
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return false
	}

	workflows := createKanban(c)

	for _, workflow := range workflows {
		_, err := database.Client.
			Doc("kanbans/"+project_id).
			Collection("workflows").
			Doc(workflow["id"].(string)).
			Set(context.Background(), map[string]interface{}{
				"name": workflow["name"],
			})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Internal Server Error",
				"message": err.Error(),
			})
			return false
		}
	}

	return true
}

func ListAllKanbans(c *gin.Context) {
	res := getProjectDetail(c)

	project_id := res.Ref.ID

	var kanban model.Kanban
	kanban.ProjectId = project_id

	iter := database.Client.
		Doc("kanbans/" + project_id).
		Collection("workflows").
		Documents(context.Background())

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Internal Server Error",
				"message": err.Error(),
			})
			return
		}

		var workflow model.Workflow
		workflow.Id = doc.Ref.ID
		workflow.Name = doc.Data()["name"].(string)

		kanban.Workflows = append(kanban.Workflows, workflow)
	}

	iter = database.Client.
		Doc("kanbans/" + project_id).
		Collection("tasks").
		Documents(context.Background())

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Internal Server Error",
				"message": err.Error(),
			})
			return
		}

		var task model.Task
		doc.DataTo(&task)
		task.Id = doc.Ref.ID
		workflowId := doc.Data()["workflowId"].(string)

		for i, workflow := range kanban.Workflows {
			if workflow.Id == workflowId {
				kanban.Workflows[i].Tasks = append(kanban.Workflows[i].Tasks, task)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "List Kanban",
		"kanban":  kanban,
	})
}

func AddWorkFlow(c *gin.Context) {
	project_id := getProjectId(c)

	var req model.AddWorkFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	new_workflow_id := shortuuid.New()
	workflow := map[string]interface{}{
		"name": req.Name,
	}

	_, err := database.Client.
		Doc("kanbans/"+project_id).
		Collection("workflows").
		Doc(new_workflow_id).
		Set(context.Background(), workflow)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      new_workflow_id,
		"status":  "Created",
		"message": "Add workflow " + req.Name + " success",
	})
}

func UpdateWorkFlow(c *gin.Context) {
	var req model.UpdateWorkFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	project_id := getProjectId(c)

	_, err := database.Client.
		Doc("kanbans/"+project_id).
		Collection("workflows").
		Doc(req.Id).
		Set(context.Background(), map[string]interface{}{
			"name": req.Name,
		}, firestore.MergeAll)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Update " + req.Name + " Success",
	})
}

func DeleteWorkFlow(c *gin.Context) {
	var req model.DeleteWorkFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	project_id := getProjectId(c)

	_, err := database.Client.
		Doc("kanbans/" + project_id).
		Collection("workflows").
		Doc(req.Id).Delete(context.Background())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Delete " + req.Id + " Success",
	})
}

func AddTask(c *gin.Context) {
	project_id := getProjectId(c)
	task_id := shortuuid.New()

	var req model.AddTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	// check workflow_Id in firebase kanbans
	_, ok := findWorkFlowById(c, project_id, req.WorkflowId)

	if !ok {
		return
	}

	task := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"workflowId":  req.WorkflowId,
	}

	_, err := database.Client.
		Doc("kanbans/"+project_id).
		Collection("tasks").
		Doc(task_id).
		Set(context.Background(), task)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "Created",
		"message": "Add " + req.Name + " Task Success",
	})
}

func UpdateTask(c *gin.Context) {
	project_id := getProjectId(c)

	var req model.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	updateData := processUpdateData(req)

	_, err := database.Client.
		Doc("kanbans/"+project_id).
		Collection("tasks").
		Doc(req.Id).
		Update(context.Background(), updateData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "Created",
		"message": "Update " + req.Name + " Task Success",
	})
}

func MoveTask(c *gin.Context) {
	project_id := getProjectId(c)

	var req model.MoveTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	// check new_workflow_Id in firebase kanbans
	_, ok := findWorkFlowById(c, project_id, req.NewWorkflowId)

	if !ok {
		return
	}

	_, err := database.Client.
		Doc("kanbans/"+project_id).
		Collection("tasks").
		Doc(req.Id).
		Set(context.Background(), map[string]interface{}{
			"workflowId": req.NewWorkflowId,
		}, firestore.MergeAll)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Move Task Success",
	})
}

func DeleteTask(c *gin.Context) {
	project_id := getProjectId(c)

	req := model.DeleteTaskRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	_, err := database.Client.
		Doc("kanbans/" + project_id).
		Collection("tasks").
		Doc(req.Id).
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
		"message": "Delete Task Success",
	})
}

func findWorkFlowById(c *gin.Context, project_id, workflow_Id string) (model.Workflow, bool) {
	dsnap, err := database.Client.
		Doc("kanbans/" + project_id).
		Collection("workflows").
		Doc(workflow_Id).
		Get(context.Background())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return model.Workflow{}, false
	}

	var workflow model.Workflow
	dsnap.DataTo(&workflow)

	return workflow, true
}
