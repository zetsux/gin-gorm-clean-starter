package repository

import (
	"context"
	"errors"
	"math"

	"github.com/zetsux/gin-gorm-template-clean/common"
	"github.com/zetsux/gin-gorm-template-clean/internal/entity"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	// db transaction
	BeginTx(ctx context.Context) (*gorm.DB, error)
	CommitTx(ctx context.Context, tx *gorm.DB) error
	RollbackTx(ctx context.Context, tx *gorm.DB)

	// functional
	CreateNewUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
	GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, error)
	GetUserByID(ctx context.Context, tx *gorm.DB, id string) (entity.User, error)
	GetAllUsers(ctx context.Context, req common.GetsRequest, tx *gorm.DB) ([]entity.User, int64, int64, error)
	UpdateNameUser(ctx context.Context, tx *gorm.DB, name string, user entity.User) (entity.User, error)
	UpdateUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
	DeleteUserByID(ctx context.Context, tx *gorm.DB, id string) error
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := ur.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (ur *userRepository) CommitTx(ctx context.Context, tx *gorm.DB) error {
	err := tx.WithContext(ctx).Commit().Error
	if err == nil {
		return err
	}
	return nil
}

func (ur *userRepository) RollbackTx(ctx context.Context, tx *gorm.DB) {
	tx.WithContext(ctx).Debug().Rollback()
}

func (ur *userRepository) CreateNewUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	var err error
	if tx == nil {
		tx = ur.db.WithContext(ctx).Debug().Create(&user)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Create(&user).Error
	}

	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (ur *userRepository) GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, error) {
	var err error
	var user entity.User
	if tx == nil {
		tx = ur.db.WithContext(ctx).Debug().Where("email = $1", email).Take(&user)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Where("email = $1", email).Take(&user).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return user, err
	}
	return user, nil
}

func (ur *userRepository) GetUserByID(ctx context.Context, tx *gorm.DB, id string) (entity.User, error) {
	var err error
	var user entity.User
	if tx == nil {
		tx = ur.db.WithContext(ctx).Debug().Where("id = $1", id).Take(&user)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Where("id = $1", id).Take(&user).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return user, err
	}
	return user, nil
}

func (ur *userRepository) GetAllUsers(ctx context.Context, req common.GetsRequest, tx *gorm.DB) ([]entity.User, int64, int64, error) {
	var err error
	var users []entity.User
	var total int64

	if tx == nil {
		tx = ur.db
	}

	stmt := tx.WithContext(ctx).Debug()
	if req.Search != "" {
		searchQuery := "%" + req.Search + "%"
		err = tx.WithContext(ctx).Model(&entity.User{}).Where("name ILIKE ? OR email ILIKE ?", searchQuery, searchQuery).Count(&total).Error
		if err != nil {
			return nil, 0, 0, err
		}
		stmt = stmt.Where("name ILIKE ? OR email ILIKE ?", searchQuery, searchQuery)
	} else {
		err = tx.WithContext(ctx).Model(&entity.User{}).Count(&total).Error
		if err != nil {
			return nil, 0, 0, err
		}
	}

	if req.Sort != "" {
		stmt = stmt.Order(req.Sort)
	}

	lastPage := int64(math.Ceil(float64(total) / float64(req.Limit)))
	if req.Limit == 0 {
		err = stmt.Find(&users).Error
	} else {
		if req.Page <= 0 || int64(req.Page) > lastPage {
			return nil, 0, 0, common.ErrInvalidPage
		}
		err = stmt.Offset(((req.Page - 1) * req.Limit)).Limit(req.Limit).Find(&users).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return users, 0, 0, err
	}
	return users, lastPage, total, nil
}

func (ur *userRepository) UpdateNameUser(ctx context.Context, tx *gorm.DB, name string, user entity.User) (entity.User, error) {
	var err error
	userUpdate := user
	userUpdate.Name = name
	if tx == nil {
		tx = ur.db.WithContext(ctx).Debug().Save(&userUpdate)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Save(&userUpdate).Error
	}

	if err != nil {
		return userUpdate, err
	}
	return userUpdate, nil
}

func (ur *userRepository) UpdateUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	if tx == nil {
		tx = ur.db.WithContext(ctx).Debug().Updates(&user)
		if err := tx.Error; err != nil {
			return entity.User{}, err
		}
	} else {
		if err := ur.db.Updates(&user).Error; err != nil {
			return entity.User{}, err
		}
	}

	return user, nil
}

func (ur *userRepository) DeleteUserByID(ctx context.Context, tx *gorm.DB, id string) error {
	var err error
	if tx == nil {
		tx = ur.db.WithContext(ctx).Debug().Delete(&entity.User{}, &id)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Delete(&entity.User{}, &id).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return err
	}
	return nil
}
