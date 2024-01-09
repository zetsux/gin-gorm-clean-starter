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
	MessageUserRegisterSuccess = "User register successful"
	MessageUserRegisterFailed  = "Failed to process user register request"

	MessageUserLoginSuccess    = "User login successful"
	MessageUserLoginFailed     = "Failed to process user login request"
	MessageUserWrongCredential = "Entered credentials invalid"

	MessageUsersFetchSuccess = "Users fetched successfully"
	MessageUsersFetchFailed  = "Failed to fetch users"
	MessageUserFetchSuccess  = "User fetched successfully"
	MessageUserFetchFailed   = "Failed to fetch user"

	MessageUserUpdateSuccess = "User update successful"
	MessageUserUpdateFailed  = "Failed to process user update request"

	MessageUserDeleteSuccess = "User delete successful"
	MessageUserDeleteFailed  = "Failed to process user delete request"

	MessageUserPictureUpdateSuccess = "User picture update successful"
	MessageUserPictureUpdateFailed  = "Failed to process user picture update request"

	MessageUserPictureDeleteSuccess = "User picture delete successful"
	MessageUserPictureDeleteFailed  = "Failed to process user picture delete request"
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
