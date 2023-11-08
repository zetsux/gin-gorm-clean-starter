package service

import (
	"context"
	"reflect"

	"github.com/zetsux/gin-gorm-template-clean/dto"
	"github.com/zetsux/gin-gorm-template-clean/entity"
	"github.com/zetsux/gin-gorm-template-clean/repository"
	"github.com/zetsux/gin-gorm-template-clean/utils"

	"github.com/jinzhu/copier"
)

type userService struct {
	userRepository repository.UserRepository
}

type UserService interface {
	VerifyLogin(ctx context.Context, email string, password string) bool
	CreateNewUser(ctx context.Context, userDTO dto.UserRegisterRequest) (dto.UserResponse, error)
	GetAllUsers(ctx context.Context) ([]dto.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (dto.UserResponse, error)
	UpdateSelfName(ctx context.Context, userDTO dto.UserNameUpdateRequest, id string) (dto.UserResponse, error)
	GetUserByID(ctx context.Context, id string) (dto.UserResponse, error)
	DeleteSelfUser(ctx context.Context, id string) error
}

func NewUserService(userR repository.UserRepository) UserService {
	return &userService{userRepository: userR}
}

func (userS *userService) VerifyLogin(ctx context.Context, email string, password string) bool {
	userCheck, err := userS.userRepository.GetUserByEmail(ctx, nil, email)
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

func (userS *userService) CreateNewUser(ctx context.Context, userDTO dto.UserRegisterRequest) (dto.UserResponse, error) {
	userCheck, err := userS.userRepository.GetUserByEmail(ctx, nil, userDTO.Email)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if !(reflect.DeepEqual(userCheck, entity.User{})) {
		return dto.UserResponse{}, dto.ErrEmailAlreadyExists
	}

	// Fill user role
	userDTO.Role = "user"

	// Copy UserDTO to empty newly created user var
	var user entity.User
	copier.Copy(&user, &userDTO)

	// create new user
	newUser, err := userS.userRepository.CreateNewUser(ctx, nil, user)
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

func (userS *userService) GetAllUsers(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := userS.userRepository.GetAllUsers(ctx, nil)
	if err != nil {
		return []dto.UserResponse{}, err
	}

	var userResponse []dto.UserResponse
	for _, user := range users {
		userResponse = append(userResponse, dto.UserResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		})
	}

	return userResponse, nil
}

func (userS *userService) GetUserByEmail(ctx context.Context, email string) (dto.UserResponse, error) {
	user, err := userS.userRepository.GetUserByEmail(ctx, nil, email)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (userS *userService) GetUserByID(ctx context.Context, id string) (dto.UserResponse, error) {
	user, err := userS.userRepository.GetUserByID(ctx, nil, id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (userS *userService) UpdateSelfName(ctx context.Context, userDTO dto.UserNameUpdateRequest, id string) (dto.UserResponse, error) {
	user, err := userS.userRepository.GetUserByID(ctx, nil, id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	user, err = userS.userRepository.UpdateNameUser(ctx, nil, userDTO.Name, user)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (userS *userService) DeleteSelfUser(ctx context.Context, id string) error {
	userCheck, err := userS.userRepository.GetUserByID(ctx, nil, id)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(userCheck, entity.User{}) {
		return dto.ErrUserNotFound
	}

	err = userS.userRepository.DeleteUserByID(ctx, nil, id)
	if err != nil {
		return err
	}
	return nil
}
