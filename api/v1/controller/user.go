package controller

import (
	"net/http"
	"reflect"

	"github.com/zetsux/gin-gorm-template-clean/common"
	"github.com/zetsux/gin-gorm-template-clean/internal/dto"
	"github.com/zetsux/gin-gorm-template-clean/internal/service"

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
	ChangePicture(ctx *gin.Context)
	DeletePicture(ctx *gin.Context)
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
	var req common.GetsRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, common.CreateFailResponse(
			dto.MESSAGE_USERS_FETCH_FAILED,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	users, pageMeta, err := uc.userService.GetAllUsers(ctx, req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, common.CreateFailResponse(
			dto.MESSAGE_USERS_FETCH_FAILED,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	if reflect.DeepEqual(pageMeta, common.PaginationResponse{}) {
		ctx.JSON(http.StatusOK, common.CreateSuccessResponse(
			dto.MESSAGE_USERS_FETCH_SUCCESS,
			http.StatusOK, users,
		))
	} else {
		ctx.JSON(http.StatusOK, common.CreatePaginatedResponse(
			dto.MESSAGE_USERS_FETCH_SUCCESS,
			http.StatusOK, users, pageMeta,
		))
	}
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

func (uc *userController) ChangePicture(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)

	var userDTO dto.UserChangePictureRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, common.CreateFailResponse(
			dto.MESSAGE_USER_PICTURE_UPDATE_FAILED,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	res, err := uc.userService.ChangePicture(ctx, userDTO, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, common.CreateFailResponse(
			dto.MESSAGE_USER_PICTURE_UPDATE_FAILED,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, common.CreateSuccessResponse(
		dto.MESSAGE_USER_PICTURE_UPDATE_SUCCESS,
		http.StatusOK, res,
	))
}

func (uc *userController) DeletePicture(ctx *gin.Context) {
	id := ctx.Param("user_id")

	err := uc.userService.DeletePicture(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, common.CreateFailResponse(
			dto.MESSAGE_USER_PICTURE_DELETE_FAILED,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, common.CreateSuccessResponse(
		dto.MESSAGE_USER_PICTURE_DELETE_SUCCESS,
		http.StatusOK, nil,
	))
}
