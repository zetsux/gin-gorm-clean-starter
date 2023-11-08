package seeder

import (
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	if err := UserSeeder(db); err != nil {
		return err
	}

	return nil
}
