package service

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/zetsux/gin-gorm-template-clean/common/standard"
	"github.com/zetsux/gin-gorm-template-clean/common/utils"
	"github.com/zetsux/gin-gorm-template-clean/internal/dto"
	"github.com/zetsux/gin-gorm-template-clean/internal/entity"
	"github.com/zetsux/gin-gorm-template-clean/internal/repository"

	"github.com/jinzhu/copier"
)

type userService struct {
	userRepository repository.UserRepository
}

type UserService interface {
	VerifyLogin(ctx context.Context, email string, password string) bool
	CreateNewUser(ctx context.Context, ud dto.UserRegisterRequest) (dto.UserResponse, error)
	GetAllUsers(ctx context.Context, req standard.GetsRequest) ([]dto.UserResponse, standard.PaginationResponse, error)
	GetUserByPrimaryKey(ctx context.Context, key string, value string) (dto.UserResponse, error)
	UpdateSelfName(ctx context.Context, ud dto.UserNameUpdateRequest, id string) (dto.UserResponse, error)
	UpdateUserByID(ctx context.Context, ud dto.UserUpdateRequest, id string) (dto.UserResponse, error)
	DeleteUserByID(ctx context.Context, id string) error
	ChangePicture(ctx context.Context, req dto.UserChangePictureRequest, userID string) (dto.UserResponse, error)
	DeletePicture(ctx context.Context, userID string) error
}

func NewUserService(userR repository.UserRepository) UserService {
	return &userService{userRepository: userR}
}

func (us *userService) VerifyLogin(ctx context.Context, email string, password string) bool {
	userCheck, err := us.userRepository.GetUserByPrimaryKey(ctx, nil, standard.DBAttrEmail, email)
	if err != nil {
		return false
	}
	passwordCheck, err := utils.PasswordCompare(userCheck.Password, []byte(password))
	if err != nil {
		return false
	}

	if userCheck.Email == email && passwordCheck {
		return true
	}
	return false
}

func (us *userService) CreateNewUser(ctx context.Context, ud dto.UserRegisterRequest) (dto.UserResponse, error) {
	userCheck, err := us.userRepository.GetUserByPrimaryKey(ctx, nil, standard.DBAttrEmail, ud.Email)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if !(reflect.DeepEqual(userCheck, entity.User{})) {
		return dto.UserResponse{}, dto.ErrEmailAlreadyExists
	}

	// Fill user role
	ud.Role = standard.EnumRoleUser

	// Copy UserDTO to empty newly created user var
	var user entity.User
	err = copier.Copy(&user, &ud)
	if err != nil {
		return dto.UserResponse{}, err
	}

	// create new user
	newUser, err := us.userRepository.CreateNewUser(ctx, nil, user)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:    newUser.ID.String(),
		Name:  newUser.Name,
		Email: newUser.Email,
		Role:  newUser.Role,
	}, nil
}

func (us *userService) GetAllUsers(ctx context.Context, req standard.GetsRequest) (
	userResp []dto.UserResponse, pageResp standard.PaginationResponse, err error) {
	if req.Limit < 0 {
		req.Limit = 0
	}

	if req.Page < 0 {
		req.Page = 0
	}

	if req.Sort != "" && req.Sort[0] == '-' {
		req.Sort = req.Sort[1:] + " DESC"
	}

	users, lastPage, total, err := us.userRepository.GetAllUsers(ctx, req, nil)
	if err != nil {
		return []dto.UserResponse{}, standard.PaginationResponse{}, err
	}

	for _, user := range users {
		userResp = append(userResp, dto.UserResponse{
			ID:      user.ID.String(),
			Name:    user.Name,
			Email:   user.Email,
			Role:    user.Role,
			Picture: user.Picture,
		})
	}

	if req.Limit == 0 {
		return userResp, standard.PaginationResponse{}, nil
	}

	pageResp = standard.PaginationResponse{
		Page:     int64(req.Page),
		Limit:    int64(req.Limit),
		LastPage: lastPage,
		Total:    total,
	}
	return userResp, pageResp, nil
}

func (us *userService) GetUserByPrimaryKey(ctx context.Context, key string, val string) (dto.UserResponse, error) {
	user, err := us.userRepository.GetUserByPrimaryKey(ctx, nil, key, val)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:      user.ID.String(),
		Name:    user.Name,
		Email:   user.Email,
		Role:    user.Role,
		Picture: user.Picture,
	}, nil
}

func (us *userService) UpdateSelfName(ctx context.Context,
	ud dto.UserNameUpdateRequest, id string) (dto.UserResponse, error) {
	user, err := us.userRepository.GetUserByPrimaryKey(ctx, nil, standard.DBAttrID, id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	user, err = us.userRepository.UpdateNameUser(ctx, nil, ud.Name, user)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:   user.ID.String(),
		Name: user.Name,
	}, nil
}

func (us *userService) UpdateUserByID(ctx context.Context,
	ud dto.UserUpdateRequest, id string) (dto.UserResponse, error) {
	user, err := us.userRepository.GetUserByPrimaryKey(ctx, nil, standard.DBAttrID, id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if reflect.DeepEqual(user, entity.User{}) {
		return dto.UserResponse{}, dto.ErrUserNotFound
	}

	if ud.Email != "" && ud.Email != user.Email {
		us, err := us.userRepository.GetUserByPrimaryKey(ctx, nil, standard.DBAttrEmail, ud.Email)
		if err != nil {
			return dto.UserResponse{}, err
		}

		if !(reflect.DeepEqual(us, entity.User{})) {
			return dto.UserResponse{}, dto.ErrEmailAlreadyExists
		}
	}

	userEdit := entity.User{
		Name:     ud.Name,
		Email:    ud.Email,
		Role:     ud.Role,
		Password: ud.Password,
	}
	userEdit.ID = user.ID

	edited, err := us.userRepository.UpdateUser(ctx, nil, userEdit)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if edited.Name == "" {
		edited.Name = user.Name
	}
	if edited.Email == "" {
		edited.Email = user.Email
	}
	if edited.Role == "" {
		edited.Role = user.Role
	}

	return dto.UserResponse{
		ID:      edited.ID.String(),
		Name:    edited.Name,
		Email:   edited.Email,
		Role:    edited.Role,
		Picture: user.Picture,
	}, nil
}

func (us *userService) DeleteUserByID(ctx context.Context, id string) error {
	userCheck, err := us.userRepository.GetUserByPrimaryKey(ctx, nil, standard.DBAttrID, id)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(userCheck, entity.User{}) {
		return dto.ErrUserNotFound
	}

	err = us.userRepository.DeleteUserByID(ctx, nil, id)
	if err != nil {
		return err
	}
	return nil
}

func (us *userService) ChangePicture(ctx context.Context,
	req dto.UserChangePictureRequest, userID string) (dto.UserResponse, error) {
	user, err := us.userRepository.GetUserByPrimaryKey(ctx, nil, standard.DBAttrID, userID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if reflect.DeepEqual(user, entity.User{}) {
		return dto.UserResponse{}, dto.ErrUserNotFound
	}

	if user.Picture != "" {
		if err := utils.DeleteFile(user.Picture); err != nil {
			return dto.UserResponse{}, err
		}
	}

	picID := uuid.New()
	picPath := fmt.Sprintf("user_picture/%v", picID)

	userEdit := entity.User{
		Picture: picPath,
	}
	userEdit.ID = user.ID

	if err := utils.UploadFile(req.Picture, picPath); err != nil {
		return dto.UserResponse{}, err
	}

	userUpdate, err := us.userRepository.UpdateUser(ctx, nil, userEdit)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:      userUpdate.ID.String(),
		Picture: userUpdate.Picture,
	}, nil
}

func (us *userService) DeletePicture(ctx context.Context, userID string) error {
	user, err := us.userRepository.GetUserByPrimaryKey(ctx, nil, standard.DBAttrID, userID)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(user, entity.User{}) {
		return dto.ErrUserNotFound
	}

	if user.Picture == "" {
		return dto.ErrUserNoPicture
	}

	if err := utils.DeleteFile(user.Picture); err != nil {
		return err
	}

	userEdit := entity.User{
		Picture: "",
	}
	userEdit.ID = user.ID

	_, err = us.userRepository.UpdateUser(ctx, nil, userEdit)
	if err != nil {
		return err
	}

	return nil
}
