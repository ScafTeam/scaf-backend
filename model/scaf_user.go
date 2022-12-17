package model

type ScafUser struct {
	Email    string   `json:"email"`
	Projects []string `json:"projects"`
}
