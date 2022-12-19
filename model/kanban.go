package model

type Kanban struct {
	ProjectId string     `json:"ProjectId"`
	Workflows []Workflow `json:"Workflows"`
}

type Workflow struct {
	Name  string `json:"Name"`
	Tasks []Task `json:"Tasks"`
}

type Task struct {
	Id          string `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}
