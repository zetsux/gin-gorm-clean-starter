package controller

import (
	"net/http"
	"reflect"

	"github.com/zetsux/gin-gorm-template-clean/common/base"
	"github.com/zetsux/gin-gorm-template-clean/common/constant"
	"github.com/zetsux/gin-gorm-template-clean/core/helper/dto"
	"github.com/zetsux/gin-gorm-template-clean/core/helper/messages"
	"github.com/zetsux/gin-gorm-template-clean/core/service"

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
	UpdateUserByID(ctx *gin.Context)
	DeleteSelfUser(ctx *gin.Context)
	DeleteUserByID(ctx *gin.Context)
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MessageUserRegisterFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	newUser, err := uc.userService.CreateNewUser(ctx, userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MessageUserRegisterFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusCreated, base.CreateSuccessResponse(
		messages.MessageUserRegisterSuccess,
		http.StatusCreated, newUser,
	))
}

func (uc *userController) Login(ctx *gin.Context) {
	var userDTO dto.UserLoginRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MessageUserLoginFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	res := uc.userService.VerifyLogin(ctx, userDTO.Email, userDTO.Password)
	if !res {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, base.CreateFailResponse(
			messages.MessageUserWrongCredential,
			"", http.StatusBadRequest,
		))
		return
	}

	user, err := uc.userService.GetUserByPrimaryKey(ctx, constant.DBAttrEmail, userDTO.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MessageUserLoginFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	token := uc.jwtService.GenerateToken(user.ID, user.Role)
	authResp := base.CreateAuthResponse(token, user.Role)
	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MessageUserLoginSuccess,
		http.StatusOK, authResp,
	))
}

func (uc *userController) GetAllUsers(ctx *gin.Context) {
	var req base.GetsRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MessageUsersFetchFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	users, pageMeta, err := uc.userService.GetAllUsers(ctx, req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MessageUsersFetchFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	if reflect.DeepEqual(pageMeta, base.PaginationResponse{}) {
		ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
			messages.MessageUsersFetchSuccess,
			http.StatusOK, users,
		))
	} else {
		ctx.JSON(http.StatusOK, base.CreatePaginatedResponse(
			messages.MessageUsersFetchSuccess,
			http.StatusOK, users, pageMeta,
		))
	}
}

func (uc *userController) GetMe(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)
	user, err := uc.userService.GetUserByPrimaryKey(ctx, constant.DBAttrID, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MessageUserFetchFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MessageUserFetchSuccess,
		http.StatusOK, user,
	))
}

func (uc *userController) UpdateSelfName(ctx *gin.Context) {
	var userDTO dto.UserNameUpdateRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MessageUserUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	id := ctx.MustGet("ID").(string)
	user, err := uc.userService.UpdateSelfName(ctx, userDTO, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MessageUserUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MessageUserUpdateSuccess,
		http.StatusOK, user,
	))
}

func (uc *userController) UpdateUserByID(ctx *gin.Context) {
	id := ctx.Param("user_id")

	var userDTO dto.UserUpdateRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MessageUserUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	user, err := uc.userService.UpdateUserByID(ctx, userDTO, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MessageUserUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MessageUserUpdateSuccess,
		http.StatusOK, user,
	))
}

func (uc *userController) DeleteSelfUser(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)
	err := uc.userService.DeleteUserByID(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MessageUserDeleteFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MessageUserDeleteSuccess,
		http.StatusOK, nil,
	))
}

func (uc *userController) DeleteUserByID(ctx *gin.Context) {
	id := ctx.Param("user_id")
	err := uc.userService.DeleteUserByID(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MessageUserDeleteFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MessageUserDeleteSuccess,
		http.StatusOK, nil,
	))
}

func (uc *userController) ChangePicture(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)

	var userDTO dto.UserChangePictureRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MessageUserPictureUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	res, err := uc.userService.ChangePicture(ctx, userDTO, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MessageUserPictureUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MessageUserPictureUpdateSuccess,
		http.StatusOK, res,
	))
}

func (uc *userController) DeletePicture(ctx *gin.Context) {
	id := ctx.Param("user_id")
	err := uc.userService.DeletePicture(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MessageUserPictureDeleteFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MessageUserPictureDeleteSuccess,
		http.StatusOK, nil,
	))
}
