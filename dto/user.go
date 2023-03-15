package dto

type UserRegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type UserLoginRequest struct {
	UserIdentifier string `json:"user-identifier" binding:"required"`
	Password       string `json:"password" binding:"required"`
}

type UserNameUpdateRequest struct {
	Name string `json:"name" binding:"required"`
}
