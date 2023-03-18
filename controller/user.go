package controller

import (
	"fp-rpl/common"
	"fp-rpl/dto"
	"fp-rpl/entity"
	"fp-rpl/service"
	"net/http"
	"reflect"

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
		resp := common.CreateFailResponse("failed to process user register request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	// Check for duplicate Username or Email
	userCheck, err := userC.userService.GetUserByUsernameOrEmail(ctx, userDTO.Username, userDTO.Email)
	if err != nil {
		resp := common.CreateFailResponse("failed to process user register request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	// Check if duplicate is found
	if !(reflect.DeepEqual(userCheck, entity.User{})) {
		var resp common.Response
		if userCheck.Username == userDTO.Username && userCheck.Email == userDTO.Email {
			resp = common.CreateFailResponse("username and email are already used", http.StatusBadRequest)
		} else if userCheck.Username == userDTO.Username {
			resp = common.CreateFailResponse("username is already used", http.StatusBadRequest)
		} else if userCheck.Email == userDTO.Email {
			resp = common.CreateFailResponse("email is already used", http.StatusBadRequest)
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	newUser, err := userC.userService.CreateNewUser(ctx, userDTO)
	if err != nil {
		resp := common.CreateFailResponse("failed to process user register request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	resp := common.CreateSuccessResponse("successfully registered user", http.StatusCreated, newUser)
	ctx.JSON(http.StatusCreated, resp)
}

func (userC *userController) Login(ctx *gin.Context) {
	var userDTO dto.UserLoginRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		resp := common.CreateFailResponse("failed to process user login request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	res := userC.userService.VerifyLogin(ctx.Request.Context(), userDTO.UserIdentifier, userDTO.Password)
	if !res {
		response := common.CreateFailResponse("entered credentials invalid", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	user, err := userC.userService.GetUserByIdentifier(ctx.Request.Context(), userDTO.UserIdentifier)
	if err != nil {
		response := common.CreateFailResponse("failed to process user login request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	token := userC.jwtService.GenerateToken(user.ID, user.Role)
	authResp := common.CreateAuthResponse(token, user.Role)
	resp := common.CreateSuccessResponse("user login successful", http.StatusOK, authResp)
	ctx.JSON(http.StatusOK, resp)
}

func (userC *userController) GetAllUsers(ctx *gin.Context) {
	users, err := userC.userService.GetAllUsers(ctx)
	if err != nil {
		resp := common.CreateFailResponse("failed to fetch all users", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	var resp common.Response
	if len(users) == 0 {
		resp = common.CreateSuccessResponse("no user found", http.StatusOK, users)
	} else {
		resp = common.CreateSuccessResponse("successfully fetched all users", http.StatusOK, users)
	}
	ctx.JSON(http.StatusOK, resp)
}

func (userC *userController) GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")
	user, err := userC.userService.GetUserByIdentifier(ctx, username)
	if err != nil {
		resp := common.CreateFailResponse(err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	var resp common.Response
	if reflect.DeepEqual(user, entity.User{}) {
		resp = common.CreateSuccessResponse("user not found", http.StatusOK, nil)
	} else {
		resp = common.CreateSuccessResponse("successfully fetched user", http.StatusOK, user)
	}
	ctx.JSON(http.StatusOK, resp)
}

func (userC *userController) GetMe(ctx *gin.Context) {
	id := ctx.GetUint64("ID")
	user, err := userC.userService.GetUserByID(ctx, id)
	if err != nil {
		resp := common.CreateFailResponse("failed to fetch user", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	var resp common.Response
	if reflect.DeepEqual(user, entity.User{}) {
		resp = common.CreateSuccessResponse("user not found", http.StatusOK, nil)
	} else {
		resp = common.CreateSuccessResponse("successfully fetched user", http.StatusOK, user)
	}
	ctx.JSON(http.StatusOK, resp)
}

func (userC *userController) UpdateSelfName(ctx *gin.Context) {
	var userDTO dto.UserNameUpdateRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		resp := common.CreateFailResponse("failed to process user name update request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	id := ctx.GetUint64("ID")
	user, err := userC.userService.UpdateSelfName(ctx, userDTO, id)
	if err != nil {
		resp := common.CreateFailResponse(err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	var resp common.Response
	if reflect.DeepEqual(user, entity.User{}) {
		resp = common.CreateSuccessResponse("user not found", http.StatusOK, nil)
	} else {
		resp = common.CreateSuccessResponse("successfully updated user", http.StatusOK, user)
	}
	ctx.JSON(http.StatusOK, resp)
}

func (userC *userController) DeleteSelfUser(ctx *gin.Context) {
	id := ctx.GetUint64("ID")
	err := userC.userService.DeleteSelfUser(ctx, id)
	if err != nil {
		resp := common.CreateFailResponse(err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	resp := common.CreateSuccessResponse("successfully deleted user", http.StatusOK, nil)
	ctx.JSON(http.StatusOK, resp)
}
