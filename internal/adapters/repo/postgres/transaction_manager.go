package postgres

import (
	"context"

	"gorm.io/gorm"
)

type contextKey string

const txKey contextKey = "pg-tx"

type TransactionManager interface {
	ExecuteInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type transactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &transactionManager{db: db}
}

func (tm *transactionManager) ExecuteInTransaction(ctx context.Context, operations func(ctx context.Context) error) error {
	tx := tm.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	txCtx := context.WithValue(ctx, txKey, tx)

	if err := operations(txCtx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
