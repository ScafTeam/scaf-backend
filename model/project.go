package model

type Project struct {
	Id       string   `json:"Id"`
	Name     string   `json:"Name"`
	Author   string   `json:"Author"`
	CreateOn string   `json:"CreateOn"`
	Repos    []Repo   `json:"Repos"`
	Members  []string `json:"Members"`
	DevTools []string `json:"DevTools"`
	DevMode  string   `json:"DevMode"`
}

type Repo struct {
	Name string `json:"Name"`
	Url  string `json:"Url"`
}
