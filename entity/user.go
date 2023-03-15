package entity

import (
	"fp-rpl/common"
	"fp-rpl/utils"

	"gorm.io/gorm"
)

type User struct {
	common.Model
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
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
	if u.Password != "" && tx.Statement.Changed("Password") {
		var err error
		u.Password, err = utils.PasswordHash(u.Password)
		if err != nil {
			return err
		}
	}
	return nil
}
