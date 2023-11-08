package dto

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
)

type (
	UserRegisterRequest struct {
		Name     string `json:"name" form:"name" binding:"required"`
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
		Role     string `json:"role"`
	}

	UserResponse struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	UserLoginRequest struct {
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	UserNameUpdateRequest struct {
		Name string `json:"name" binding:"required"`
	}

	UserUpdateRequest struct {
		ID       string `json:"id"`
		Name     string `json:"name" form:"name"`
		Email    string `json:"email" form:"email"`
		Role     string `json:"role" form:"role"`
		Password string `json:"password" form:"password"`
	}
)
