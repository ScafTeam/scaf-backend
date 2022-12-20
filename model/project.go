package model

import "strings"

type Project struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Author   string   `json:"author"`
	CreateOn string   `json:"createOn"`
	Repos    []Repo   `json:"repos"`
	Members  []string `json:"members"`
	DevTools []string `json:"devTools"`
	DevMode  string   `json:"devMode"`
}

type Repo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type CreateProjectRequest struct {
	Name     string   `json:"name" binding:"required"`
	DevTools []string `json:"devTools"`
	DevMode  string   `json:"devMode" binding:"required"`
}

type AddMemberRequest struct {
	Email string `json:"email" binding:"email"`
}

func CheckProjectName(name string) bool {
	// check project name without invalid characters
	return !strings.ContainsAny(name, "/\\?%*:|\"<>")
}
