package postgres

import (
	"context"

	"gorm.io/gorm"
)

type BaseRepository struct {
	DB *gorm.DB
}

func (r *BaseRepository) GetDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}

	return r.DB.WithContext(ctx)
}
