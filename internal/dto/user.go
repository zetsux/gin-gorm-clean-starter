package dto

import (
	"errors"
	"mime/multipart"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNoPicture      = errors.New("user don't have any picture")
)

const (
	MESSAGE_USER_REGISTER_SUCCESS = "User register successful"
	MESSAGE_USER_REGISTER_FAILED  = "Failed to process user register request"

	MESSAGE_USER_LOGIN_SUCCESS    = "User login successful"
	MESSAGE_USER_LOGIN_FAILED     = "Failed to process user login request"
	MESSAGE_USER_WRONG_CREDENTIAL = "Entered credentials invalid"

	MESSAGE_USERS_FETCH_SUCCESS = "Users fetched successfully"
	MESSAGE_USERS_FETCH_FAILED  = "Failed to fetch users"
	MESSAGE_USER_FETCH_SUCCESS  = "User fetched successfully"
	MESSAGE_USER_FETCH_FAILED   = "Failed to fetch user"

	MESSAGE_USER_UPDATE_SUCCESS = "User update successful"
	MESSAGE_USER_UPDATE_FAILED  = "Failed to process user update request"

	MESSAGE_USER_DELETE_SUCCESS = "User delete successful"
	MESSAGE_USER_DELETE_FAILED  = "Failed to process user delete request"

	MESSAGE_USER_PICTURE_UPDATE_SUCCESS = "User picture update successful"
	MESSAGE_USER_PICTURE_UPDATE_FAILED  = "Failed to process user picture update request"

	MESSAGE_USER_PICTURE_DELETE_SUCCESS = "User picture delete successful"
	MESSAGE_USER_PICTURE_DELETE_FAILED  = "Failed to process user picture delete request"
)

type (
	UserRegisterRequest struct {
		Name     string `json:"name" form:"name" binding:"required"`
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
		Role     string `json:"role"`
	}

	UserResponse struct {
		ID      string `json:"id"`
		Name    string `json:"name,omitempty"`
		Email   string `json:"email,omitempty"`
		Role    string `json:"role,omitempty"`
		Picture string `json:"picture,omitempty"`
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

	UserChangePictureRequest struct {
		Picture *multipart.FileHeader `json:"picture" form:"picture"`
	}
)
