package model

type Kanban struct {
	ProjectId string     `json:"projectId"`
	Workflows []Workflow `json:"workflows"`
}

type Workflow struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AddWorkFlowRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateWorkFlowRequest struct {
	Id   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type DeleteWorkFlowRequest struct {
	Id string `json:"id" binding:"required"`
}

type AddTaskRequest struct {
	WorkflowId  string `json:"workflowId" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateTaskRequest struct {
	Id          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type MoveTaskRequest struct {
	Id            string `json:"id" binding:"required"`
	NewWorkflowId string `json:"newWorkflowId" binding:"required"`
}

type DeleteTaskRequest struct {
	Id         string `json:"id" binding:"required"`
	WorkflowId string `json:"workflowId" binding:"required"`
}
