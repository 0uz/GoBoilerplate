package postgres

import (
	"context"

	"github.com/ouz/goboilerplate/internal/adapters/repo/postgres"
	"github.com/ouz/goboilerplate/internal/domain/auth"
	"github.com/ouz/goboilerplate/pkg/errors"
	"gorm.io/gorm"
)

type authRepository struct {
	postgres.BaseRepository
}

func NewAuthRepository(db *gorm.DB) auth.AuthRepository {
	return &authRepository{BaseRepository: postgres.BaseRepository{
		DB: db,
	}}
}

func (r *authRepository) FindClientBySecret(ctx context.Context, secret string) (*auth.Client, error) {
	var client auth.Client
	if err := r.GetDB(ctx).Where("client_secret = ?", secret).First(&client).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundError("Client not found", err)
		}
		return nil, errors.InternalError("Failed to find client", err)
	}
	return &client, nil
}
