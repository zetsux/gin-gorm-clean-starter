package repository

import (
	"context"
	"log"

	"gorm.io/gorm"
)

type txRepository struct {
	db *gorm.DB
}

type TxRepository interface {
	// db
	DB() *gorm.DB

	// tx
	BeginTx(ctx context.Context, db *gorm.DB) (*gorm.DB, error)
	CommitOrRollbackTx(ctx context.Context, tx *gorm.DB, err error)
}

func NewTxRepository(db *gorm.DB) *txRepository {
	return &txRepository{db: db}
}

func (txr txRepository) DB() *gorm.DB {
	return txr.db
}

func (txr txRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := txr.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (txr txRepository) CommitOrRollbackTx(ctx context.Context, tx *gorm.DB, err error) {
	if err != nil {
		log.Println("Error occurred: ", err)
		tx.WithContext(ctx).Debug().Rollback()
		return
	}

	err = tx.WithContext(ctx).Commit().Error
	if err != nil {
		log.Println("Commit failed: ", err)
		return
	}
	log.Println("Committed successfully")
}
