package models

type UserAdd struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required"`
}
