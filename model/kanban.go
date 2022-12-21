package model

type Kanban struct {
	ProjectId string              `json:"projectId"`
	Workflows map[string]Workflow `json:"workflows"`
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
	WorkflowName string `json:"name" binding:"required"`
}

type AddTaskRequest struct {
	WorkflowId  string `json:"workflowId" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateTaskRequest struct {
	TaskId      string `json:"id" binding:"required"`
	WorkflowId  string `json:"workflowId" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type MoveTaskRequest struct {
	TaskId        string `json:"id" binding:"required"`
	WorkflowId    string `json:"workflowId" binding:"required"`
	NewWorkflowId string `json:"newWorkflowId" binding:"required"`
}

type DeleteTaskRequest struct {
	TaskId     string `json:"id" binding:"required"`
	WorkflowId string `json:"workflowId" binding:"required"`
}
