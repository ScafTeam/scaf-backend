package model

type ScafUser struct {
	Email    string   `json:"email"`
	Projects []string `json:"projects"`
}

type UserRegisterRequest struct {
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"required"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserForgotPasswordRequest struct {
	Email string `json:"email" binding:"email"`
}
