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
	VerifyLogin(ctx context.Context, identifier string, password string) bool
	CreateNewUser(ctx context.Context, userDTO dto.UserRegisterRequest) (dto.UserResponse, error)
	GetAllUsers(ctx context.Context) ([]dto.UserResponse, error)
	GetUserByIdentifier(ctx context.Context, identifier string) (dto.UserResponse, error)
	GetUserByUsernameOrEmail(ctx context.Context, username string, email string) (dto.UserResponse, error)
	UpdateSelfName(ctx context.Context, userDTO dto.UserNameUpdateRequest, id string) (dto.UserResponse, error)
	GetUserByID(ctx context.Context, id string) (dto.UserResponse, error)
	DeleteSelfUser(ctx context.Context, id string) error
}

func NewUserService(userR repository.UserRepository) UserService {
	return &userService{userRepository: userR}
}

func (userS *userService) VerifyLogin(ctx context.Context, identifier string, password string) bool {
	userCheck, err := userS.userRepository.GetUserByIdentifier(ctx, nil, identifier, identifier)
	if err != nil {
		return false
	}
	passwordCheck, err := utils.PasswordCompare(userCheck.Password, []byte(password))
	if err != nil {
		return false
	}

	if (userCheck.Username == identifier || userCheck.Email == identifier) && passwordCheck {
		return true
	}
	return false
}

func (userS *userService) CreateNewUser(ctx context.Context, userDTO dto.UserRegisterRequest) (dto.UserResponse, error) {
	userCheck, err := userS.GetUserByUsernameOrEmail(ctx, userDTO.Username, userDTO.Email)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if !(reflect.DeepEqual(userCheck, entity.User{})) {
		if userCheck.Email == userDTO.Email {
			return dto.UserResponse{}, dto.ErrEmailAlreadyExists
		} else if userCheck.Username == userDTO.Username {
			return dto.UserResponse{}, dto.ErrUsernameAlreadyExists
		}
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
		ID:       newUser.ID.String(),
		Name:     newUser.Name,
		Username: newUser.Username,
		Email:    newUser.Email,
		Role:     newUser.Role,
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
			ID:       user.ID.String(),
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		})
	}

	return userResponse, nil
}

func (userS *userService) GetUserByIdentifier(ctx context.Context, identifier string) (dto.UserResponse, error) {
	user, err := userS.userRepository.GetUserByIdentifier(ctx, nil, identifier, identifier)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:       user.ID.String(),
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (userS *userService) GetUserByUsernameOrEmail(ctx context.Context, username string, email string) (dto.UserResponse, error) {
	user, err := userS.userRepository.GetUserByIdentifier(ctx, nil, username, email)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:       user.ID.String(),
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (userS *userService) GetUserByID(ctx context.Context, id string) (dto.UserResponse, error) {
	user, err := userS.userRepository.GetUserByID(ctx, nil, id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:       user.ID.String(),
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
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
		ID:       user.ID.String(),
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (userS *userService) DeleteSelfUser(ctx context.Context, id string) error {
	err := userS.userRepository.DeleteUserByID(ctx, nil, id)
	if err != nil {
		return err
	}
	return nil
}
