package model

type ScafUser struct {
	Email    string   `json:"email"`
	Avatar   string   `json:"avatar"`
	Bio      string   `json:"bio"`
	Nickname string   `json:"nickname"`
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

type UpdateUserRequest struct {
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	Bio      string `json:"bio"`
}

type UpdateUserPasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}
