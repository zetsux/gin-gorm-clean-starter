package database

import (
	"fmt"

	"github.com/zetsux/gin-gorm-clean-starter/core/entity"
	"github.com/zetsux/gin-gorm-clean-starter/database/seeder"
	"gorm.io/gorm"
)

func DBMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		entity.User{},
	)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	if err := DBSeed(db); err != nil {
		panic(err)
	}
}

func DBSeed(db *gorm.DB) error {
	if err := seeder.UserSeeder(db); err != nil {
		return err
	}

	return nil
}
