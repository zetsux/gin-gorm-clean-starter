package service

import (
	"context"
	"errors"
	"fp-rpl/dto"
	"fp-rpl/entity"
	"fp-rpl/repository"
	"fp-rpl/utils"
	"reflect"

	"github.com/jinzhu/copier"
)

type userService struct {
	userRepository repository.UserRepository
}

type UserService interface {
	VerifyLogin(ctx context.Context, identifier string, password string) bool
	CreateNewUser(ctx context.Context, userDTO dto.UserRegisterRequest) (entity.User, error)
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	GetUserByIdentifier(ctx context.Context, identifier string) (entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	UpdateSelfName(ctx context.Context, userDTO dto.UserNameUpdateRequest, id uint64) (entity.User, error)
	DeleteSelfUser(ctx context.Context, id uint64) error
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

func (userS *userService) CreateNewUser(ctx context.Context, userDTO dto.UserRegisterRequest) (entity.User, error) {
	// Copy UserDTO to empty newly created user var
	var user entity.User
	copier.Copy(&user, &userDTO)

	// Check for duplicate Username or Email
	userCheck, err := userS.userRepository.GetUserByIdentifier(ctx, nil, userDTO.Username, userDTO.Email)
	if err != nil {
		return entity.User{}, err
	}

	// Check if duplicate is found
	if !(reflect.DeepEqual(userCheck, entity.User{})) {
		if userCheck.Username == userDTO.Username {
			return entity.User{}, errors.New("username already exists")
		} else if userCheck.Email == userDTO.Email {
			return entity.User{}, errors.New("email already used")
		}
	}

	// create new user
	newUser, err := userS.userRepository.CreateNewUser(ctx, nil, user)
	if err != nil {
		return entity.User{}, err
	}
	return newUser, nil
}

func (userS *userService) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	users, err := userS.userRepository.GetAllUsers(ctx, nil)
	if err != nil {
		return []entity.User{}, err
	}
	return users, nil
}

func (userS *userService) GetUserByIdentifier(ctx context.Context, identifier string) (entity.User, error) {
	user, err := userS.userRepository.GetUserByIdentifier(ctx, nil, identifier, identifier)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (userS *userService) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	user, err := userS.userRepository.GetUserByIdentifier(ctx, nil, username, "")
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (userS *userService) UpdateSelfName(ctx context.Context, userDTO dto.UserNameUpdateRequest, id uint64) (entity.User, error) {
	user, err := userS.userRepository.GetUserByID(ctx, nil, id)
	if err != nil {
		return entity.User{}, err
	}

	user, err = userS.userRepository.UpdateNameUser(ctx, nil, userDTO.Name, user)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (userS *userService) DeleteSelfUser(ctx context.Context, id uint64) error {
	err := userS.userRepository.DeleteUserByID(ctx, nil, id)
	if err != nil {
		return err
	}
	return nil
}
