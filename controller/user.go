package controller

import (
	"net/http"

	"github.com/zetsux/gin-gorm-template-clean/common"
	"github.com/zetsux/gin-gorm-template-clean/dto"
	"github.com/zetsux/gin-gorm-template-clean/service"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userService service.UserService
	jwtService  service.JWTService
}

type UserController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	GetAllUsers(ctx *gin.Context)
	GetMe(ctx *gin.Context)
	UpdateSelfName(ctx *gin.Context)
	UpdateUserById(ctx *gin.Context)
	DeleteSelfUser(ctx *gin.Context)
	DeleteUserById(ctx *gin.Context)
}

func NewUserController(userS service.UserService, jwtS service.JWTService) UserController {
	return &userController{
		userService: userS,
		jwtService:  jwtS,
	}
}

func (uc *userController) Register(ctx *gin.Context) {
	var userDTO dto.UserRegisterRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, common.CreateFailResponse(
			dto.MESSAGE_USER_REGISTER_FAILED,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	newUser, err := uc.userService.CreateNewUser(ctx, userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, common.CreateFailResponse(
			dto.MESSAGE_USER_REGISTER_FAILED,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusCreated, common.CreateSuccessResponse(
		dto.MESSAGE_USER_REGISTER_SUCCESS,
		http.StatusCreated, newUser,
	))
}

func (uc *userController) Login(ctx *gin.Context) {
	var userDTO dto.UserLoginRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, common.CreateFailResponse(
			dto.MESSAGE_USER_LOGIN_FAILED,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	res := uc.userService.VerifyLogin(ctx.Request.Context(), userDTO.Email, userDTO.Password)
	if !res {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, common.CreateFailResponse(
			dto.MESSAGE_USER_WRONG_CREDENTIAL,
			"", http.StatusBadRequest,
		))
		return
	}

	user, err := uc.userService.GetUserByEmail(ctx.Request.Context(), userDTO.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, common.CreateFailResponse(
			dto.MESSAGE_USER_LOGIN_FAILED,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	token := uc.jwtService.GenerateToken(user.ID, user.Role)
	authResp := common.CreateAuthResponse(token, user.Role)
	ctx.JSON(http.StatusOK, common.CreateSuccessResponse(
		dto.MESSAGE_USER_LOGIN_SUCCESS,
		http.StatusOK, authResp,
	))
}

func (uc *userController) GetAllUsers(ctx *gin.Context) {
	users, err := uc.userService.GetAllUsers(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, common.CreateFailResponse(
			dto.MESSAGE_USERS_FETCH_FAILED,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, common.CreateSuccessResponse(
		dto.MESSAGE_USERS_FETCH_SUCCESS,
		http.StatusOK, users,
	))
}

func (uc *userController) GetMe(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)
	user, err := uc.userService.GetUserByID(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, common.CreateFailResponse(
			dto.MESSAGE_USER_FETCH_FAILED,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, common.CreateSuccessResponse(
		dto.MESSAGE_USER_FETCH_SUCCESS,
		http.StatusOK, user,
	))
}

func (uc *userController) UpdateSelfName(ctx *gin.Context) {
	var userDTO dto.UserNameUpdateRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, common.CreateFailResponse(
			dto.MESSAGE_USER_UPDATE_FAILED,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	id := ctx.MustGet("ID").(string)
	user, err := uc.userService.UpdateSelfName(ctx, userDTO, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, common.CreateFailResponse(
			dto.MESSAGE_USER_UPDATE_FAILED,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, common.CreateSuccessResponse(
		dto.MESSAGE_USER_UPDATE_SUCCESS,
		http.StatusOK, user,
	))
}

func (uc *userController) UpdateUserById(ctx *gin.Context) {
	id := ctx.Param("user_id")

	var userDTO dto.UserUpdateRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, common.CreateFailResponse(
			dto.MESSAGE_USER_UPDATE_FAILED,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	user, err := uc.userService.UpdateUserById(ctx, userDTO, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, common.CreateFailResponse(
			dto.MESSAGE_USER_UPDATE_FAILED,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, common.CreateSuccessResponse(
		dto.MESSAGE_USER_UPDATE_SUCCESS,
		http.StatusOK, user,
	))
}

func (uc *userController) DeleteSelfUser(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)
	err := uc.userService.DeleteUserById(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, common.CreateFailResponse(
			dto.MESSAGE_USER_DELETE_FAILED,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, common.CreateSuccessResponse(
		dto.MESSAGE_USER_DELETE_SUCCESS,
		http.StatusOK, nil,
	))
}

func (uc *userController) DeleteUserById(ctx *gin.Context) {
	id := ctx.Param("user_id")
	err := uc.userService.DeleteUserById(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, common.CreateFailResponse(
			dto.MESSAGE_USER_DELETE_FAILED,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, common.CreateSuccessResponse(
		dto.MESSAGE_USER_DELETE_SUCCESS,
		http.StatusOK, nil,
	))
}
