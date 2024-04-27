package entity

import (
	"github.com/google/uuid"
	"github.com/zetsux/gin-gorm-clean-starter/common/base"
	"github.com/zetsux/gin-gorm-clean-starter/common/util"

	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name     string    `json:"name" gorm:"not null"`
	Email    string    `json:"email" gorm:"unique;not null"`
	Password string    `json:"password" gorm:"not null"`
	Role     string    `json:"role" gorm:"not null"`
	Picture  string    `json:"picture"`
	base.Model
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	var err error
	u.Password, err = util.PasswordHash(u.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if u.Password != "" {
		var err error
		u.Password, err = util.PasswordHash(u.Password)
		if err != nil {
			return err
		}
	}
	return nil
}
