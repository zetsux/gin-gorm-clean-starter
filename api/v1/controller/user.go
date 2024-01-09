package controller

import (
	"net/http"
	"reflect"

	"github.com/zetsux/gin-gorm-template-clean/common/standard"
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, standard.CreateFailResponse(
			dto.MessageUserRegisterFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	newUser, err := uc.userService.CreateNewUser(ctx, userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, standard.CreateFailResponse(
			dto.MessageUserRegisterFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusCreated, standard.CreateSuccessResponse(
		dto.MessageUserRegisterSuccess,
		http.StatusCreated, newUser,
	))
}

func (uc *userController) Login(ctx *gin.Context) {
	var userDTO dto.UserLoginRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, standard.CreateFailResponse(
			dto.MessageUserLoginFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	res := uc.userService.VerifyLogin(ctx, userDTO.Email, userDTO.Password)
	if !res {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, standard.CreateFailResponse(
			dto.MessageUserWrongCredential,
			"", http.StatusBadRequest,
		))
		return
	}

	user, err := uc.userService.GetUserByPrimaryKey(ctx, standard.DBAttrEmail, userDTO.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, standard.CreateFailResponse(
			dto.MessageUserLoginFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	token := uc.jwtService.GenerateToken(user.ID, user.Role)
	authResp := standard.CreateAuthResponse(token, user.Role)
	ctx.JSON(http.StatusOK, standard.CreateSuccessResponse(
		dto.MessageUserLoginSuccess,
		http.StatusOK, authResp,
	))
}

func (uc *userController) GetAllUsers(ctx *gin.Context) {
	var req standard.GetsRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, standard.CreateFailResponse(
			dto.MessageUsersFetchFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	users, pageMeta, err := uc.userService.GetAllUsers(ctx, req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, standard.CreateFailResponse(
			dto.MessageUsersFetchFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	if reflect.DeepEqual(pageMeta, standard.PaginationResponse{}) {
		ctx.JSON(http.StatusOK, standard.CreateSuccessResponse(
			dto.MessageUsersFetchSuccess,
			http.StatusOK, users,
		))
	} else {
		ctx.JSON(http.StatusOK, standard.CreatePaginatedResponse(
			dto.MessageUsersFetchSuccess,
			http.StatusOK, users, pageMeta,
		))
	}
}

func (uc *userController) GetMe(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)
	user, err := uc.userService.GetUserByPrimaryKey(ctx, standard.DBAttrID, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, standard.CreateFailResponse(
			dto.MessageUserFetchFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, standard.CreateSuccessResponse(
		dto.MessageUserFetchSuccess,
		http.StatusOK, user,
	))
}

func (uc *userController) UpdateSelfName(ctx *gin.Context) {
	var userDTO dto.UserNameUpdateRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, standard.CreateFailResponse(
			dto.MessageUserUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	id := ctx.MustGet("ID").(string)
	user, err := uc.userService.UpdateSelfName(ctx, userDTO, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, standard.CreateFailResponse(
			dto.MessageUserUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, standard.CreateSuccessResponse(
		dto.MessageUserUpdateSuccess,
		http.StatusOK, user,
	))
}

func (uc *userController) UpdateUserByID(ctx *gin.Context) {
	id := ctx.Param("user_id")

	var userDTO dto.UserUpdateRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, standard.CreateFailResponse(
			dto.MessageUserUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	user, err := uc.userService.UpdateUserByID(ctx, userDTO, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, standard.CreateFailResponse(
			dto.MessageUserUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, standard.CreateSuccessResponse(
		dto.MessageUserUpdateSuccess,
		http.StatusOK, user,
	))
}

func (uc *userController) DeleteSelfUser(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)
	err := uc.userService.DeleteUserByID(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, standard.CreateFailResponse(
			dto.MessageUserDeleteFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, standard.CreateSuccessResponse(
		dto.MessageUserDeleteSuccess,
		http.StatusOK, nil,
	))
}

func (uc *userController) DeleteUserByID(ctx *gin.Context) {
	id := ctx.Param("user_id")
	err := uc.userService.DeleteUserByID(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, standard.CreateFailResponse(
			dto.MessageUserDeleteFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, standard.CreateSuccessResponse(
		dto.MessageUserDeleteSuccess,
		http.StatusOK, nil,
	))
}

func (uc *userController) ChangePicture(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)

	var userDTO dto.UserChangePictureRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, standard.CreateFailResponse(
			dto.MessageUserPictureUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	res, err := uc.userService.ChangePicture(ctx, userDTO, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, standard.CreateFailResponse(
			dto.MessageUserPictureUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, standard.CreateSuccessResponse(
		dto.MessageUserPictureUpdateSuccess,
		http.StatusOK, res,
	))
}

func (uc *userController) DeletePicture(ctx *gin.Context) {
	id := ctx.Param("user_id")
	err := uc.userService.DeletePicture(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, standard.CreateFailResponse(
			dto.MessageUserPictureDeleteFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, standard.CreateSuccessResponse(
		dto.MessageUserPictureDeleteSuccess,
		http.StatusOK, nil,
	))
}
