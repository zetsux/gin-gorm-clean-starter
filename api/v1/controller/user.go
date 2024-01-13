package controller

import (
	"net/http"
	"reflect"

	"github.com/zetsux/gin-gorm-clean-starter/common/base"
	"github.com/zetsux/gin-gorm-clean-starter/common/constant"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/dto"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/messages"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"

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
			messages.MsgUserRegisterFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	newUser, err := uc.userService.CreateNewUser(ctx, userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserRegisterFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusCreated, base.CreateSuccessResponse(
		messages.MsgUserRegisterSuccess,
		http.StatusCreated, newUser,
	))
}

func (uc *userController) Login(ctx *gin.Context) {
	var userDTO dto.UserLoginRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserLoginFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	res := uc.userService.VerifyLogin(ctx, userDTO.Email, userDTO.Password)
	if !res {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, base.CreateFailResponse(
			messages.MsgUserWrongCredential,
			"", http.StatusBadRequest,
		))
		return
	}

	user, err := uc.userService.GetUserByPrimaryKey(ctx, constant.DBAttrEmail, userDTO.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserLoginFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	token := uc.jwtService.GenerateToken(user.ID, user.Role)
	authResp := base.CreateAuthResponse(token, user.Role)
	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserLoginSuccess,
		http.StatusOK, authResp,
	))
}

func (uc *userController) GetAllUsers(ctx *gin.Context) {
	var req base.GetsRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUsersFetchFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	users, pageMeta, err := uc.userService.GetAllUsers(ctx, req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUsersFetchFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	if reflect.DeepEqual(pageMeta, base.PaginationResponse{}) {
		ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
			messages.MsgUsersFetchSuccess,
			http.StatusOK, users,
		))
	} else {
		ctx.JSON(http.StatusOK, base.CreatePaginatedResponse(
			messages.MsgUsersFetchSuccess,
			http.StatusOK, users, pageMeta,
		))
	}
}

func (uc *userController) GetMe(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)
	user, err := uc.userService.GetUserByPrimaryKey(ctx, constant.DBAttrID, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserFetchFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserFetchSuccess,
		http.StatusOK, user,
	))
}

func (uc *userController) UpdateSelfName(ctx *gin.Context) {
	var userDTO dto.UserNameUpdateRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	id := ctx.MustGet("ID").(string)
	user, err := uc.userService.UpdateSelfName(ctx, userDTO, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserUpdateSuccess,
		http.StatusOK, user,
	))
}

func (uc *userController) UpdateUserByID(ctx *gin.Context) {
	id := ctx.Param("user_id")

	var userDTO dto.UserUpdateRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	user, err := uc.userService.UpdateUserByID(ctx, userDTO, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserUpdateSuccess,
		http.StatusOK, user,
	))
}

func (uc *userController) DeleteSelfUser(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)
	err := uc.userService.DeleteUserByID(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserDeleteFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserDeleteSuccess,
		http.StatusOK, nil,
	))
}

func (uc *userController) DeleteUserByID(ctx *gin.Context) {
	id := ctx.Param("user_id")
	err := uc.userService.DeleteUserByID(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserDeleteFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserDeleteSuccess,
		http.StatusOK, nil,
	))
}

func (uc *userController) ChangePicture(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)

	var userDTO dto.UserChangePictureRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserPictureUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	res, err := uc.userService.ChangePicture(ctx, userDTO, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserPictureUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserPictureUpdateSuccess,
		http.StatusOK, res,
	))
}

func (uc *userController) DeletePicture(ctx *gin.Context) {
	id := ctx.Param("user_id")
	err := uc.userService.DeletePicture(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserPictureDeleteFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserPictureDeleteSuccess,
		http.StatusOK, nil,
	))
}
