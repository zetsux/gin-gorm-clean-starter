package controller

import (
	"net/http"
	"reflect"

	"github.com/zetsux/gin-gorm-template-clean/common"
	"github.com/zetsux/gin-gorm-template-clean/dto"
	"github.com/zetsux/gin-gorm-template-clean/entity"
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
	GetUserByUsername(ctx *gin.Context)
	GetMe(ctx *gin.Context)
	UpdateSelfName(ctx *gin.Context)
	DeleteSelfUser(ctx *gin.Context)
}

func NewUserController(userS service.UserService, jwtS service.JWTService) UserController {
	return &userController{
		userService: userS,
		jwtService:  jwtS,
	}
}

func (userC *userController) Register(ctx *gin.Context) {
	var userDTO dto.UserRegisterRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		resp := common.CreateFailResponse("Failed to process user register request", err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	newUser, err := userC.userService.CreateNewUser(ctx, userDTO)
	if err != nil {
		resp := common.CreateFailResponse("Failed to process user register request", err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	resp := common.CreateSuccessResponse("Successfully registered user", http.StatusCreated, newUser)
	ctx.JSON(http.StatusCreated, resp)
}

func (userC *userController) Login(ctx *gin.Context) {
	var userDTO dto.UserLoginRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		resp := common.CreateFailResponse("Failed to process user login request", err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	res := userC.userService.VerifyLogin(ctx.Request.Context(), userDTO.UserIdentifier, userDTO.Password)
	if !res {
		response := common.CreateFailResponse("Entered credentials invalid", err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	user, err := userC.userService.GetUserByIdentifier(ctx.Request.Context(), userDTO.UserIdentifier)
	if err != nil {
		response := common.CreateFailResponse("Failed to process user login request", err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	token := userC.jwtService.GenerateToken(user.ID, user.Role)
	authResp := common.CreateAuthResponse(token, user.Role)
	resp := common.CreateSuccessResponse("User login successful", http.StatusOK, authResp)
	ctx.JSON(http.StatusOK, resp)
}

func (userC *userController) GetAllUsers(ctx *gin.Context) {
	users, err := userC.userService.GetAllUsers(ctx)
	if err != nil {
		resp := common.CreateFailResponse("Failed to fetch all users", err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	var resp common.Response
	if len(users) == 0 {
		resp = common.CreateSuccessResponse("No user found", http.StatusOK, users)
	} else {
		resp = common.CreateSuccessResponse("Successfully fetched all users", http.StatusOK, users)
	}
	ctx.JSON(http.StatusOK, resp)
}

func (userC *userController) GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")
	user, err := userC.userService.GetUserByIdentifier(ctx, username)
	if err != nil {
		resp := common.CreateFailResponse("Failed to fetch user", err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	var resp common.Response
	if reflect.DeepEqual(user, entity.User{}) {
		resp = common.CreateSuccessResponse("User not found", http.StatusOK, nil)
	} else {
		resp = common.CreateSuccessResponse("Successfully fetched user", http.StatusOK, user)
	}
	ctx.JSON(http.StatusOK, resp)
}

func (userC *userController) GetMe(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)
	user, err := userC.userService.GetUserByID(ctx, id)
	if err != nil {
		resp := common.CreateFailResponse("Failed to fetch user", err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	var resp common.Response
	if reflect.DeepEqual(user, entity.User{}) {
		resp = common.CreateSuccessResponse("User not found", http.StatusOK, nil)
	} else {
		resp = common.CreateSuccessResponse("Successfully fetched user", http.StatusOK, user)
	}
	ctx.JSON(http.StatusOK, resp)
}

func (userC *userController) UpdateSelfName(ctx *gin.Context) {
	var userDTO dto.UserNameUpdateRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		resp := common.CreateFailResponse("Failed to process user name update request", err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	id := ctx.MustGet("ID").(string)
	user, err := userC.userService.UpdateSelfName(ctx, userDTO, id)
	if err != nil {
		resp := common.CreateFailResponse("Failed to process user name update request", err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	var resp common.Response
	if reflect.DeepEqual(user, entity.User{}) {
		resp = common.CreateSuccessResponse("User not found", http.StatusOK, nil)
	} else {
		resp = common.CreateSuccessResponse("Successfully updated user", http.StatusOK, user)
	}
	ctx.JSON(http.StatusOK, resp)
}

func (userC *userController) DeleteSelfUser(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)
	err := userC.userService.DeleteSelfUser(ctx, id)
	if err != nil {
		resp := common.CreateFailResponse("Failed to delete user", err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	resp := common.CreateSuccessResponse("Successfully deleted user", http.StatusOK, nil)
	ctx.JSON(http.StatusOK, resp)
}
