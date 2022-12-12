package model

import (
	"github.com/ScafTeam/firebase-go-client/auth"
)

type ScafUser struct {
	auth.User
	Projects []Project `json:"projects"`
}
