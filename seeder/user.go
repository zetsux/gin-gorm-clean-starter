package seeder

import (
	"errors"

	"github.com/zetsux/gin-gorm-template-clean/common"
	"github.com/zetsux/gin-gorm-template-clean/entity"
	"gorm.io/gorm"
)

func UserSeeder(db *gorm.DB) error {
	var dummyUsers = []entity.User{
		{
			Name:     "Admin",
			Email:    "admin@gmail.com",
			Password: "admin1",
			Role:     common.ENUM_ROLE_ADMIN,
		},
		{
			Name:     "User",
			Email:    "user@gmail.com",
			Password: "user1",
			Role:     common.ENUM_ROLE_USER,
		},
	}

	hasTable := db.Migrator().HasTable(&entity.User{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.User{}); err != nil {
			return err
		}
	}

	for _, data := range dummyUsers {
		var user entity.User
		err := db.Where(&entity.User{Email: data.Email}).First(&user).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		isData := db.Find(&user, "email = ?", data.Email).RowsAffected
		if isData == 0 {
			if err := db.Create(&data).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
