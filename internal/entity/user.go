package entity

import (
	"github.com/zetsux/gin-gorm-template-clean/common/standard"
	"github.com/zetsux/gin-gorm-template-clean/common/utils"

	"gorm.io/gorm"
)

type User struct {
	standard.Model
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Role     string `json:"role" gorm:"not null"`
	Picture  string `json:"picture"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	var err error
	u.Password, err = utils.PasswordHash(u.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if u.Password != "" {
		var err error
		u.Password, err = utils.PasswordHash(u.Password)
		if err != nil {
			return err
		}
	}
	return nil
}
