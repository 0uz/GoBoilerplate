package postgres

import (
	"context"

	"github.com/ouz/goauthboilerplate/internal/adapters/repo/postgres"
	"github.com/ouz/goauthboilerplate/internal/domain/auth"
	"github.com/ouz/goauthboilerplate/pkg/errors"
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

func (r *authRepository) IsTokenRevoked(ctx context.Context, jti, userId string) (bool, error) {
	var revoked bool
	if err := r.GetDB(ctx).Model(&TokenEntity{}).
		Where("id = ? AND user_id = ?", jti, userId).
		Select("revoked").
		First(&revoked).Error; err != nil {
		return false, errors.InternalError("Failed to check if token is revoked", err)
	}
	return revoked, nil
}

func (r *authRepository) RevokeAllOldTokens(ctx context.Context, userID string) error {
	if err := r.GetDB(ctx).Model(&TokenEntity{}).
		Where("user_id = ? AND revoked = false", userID).
		Update("revoked", true).Error; err != nil {
		return errors.AuthError("Failed to revoke old tokens", err)
	}
	return nil
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

func (r *authRepository) SaveTokenPairs(ctx context.Context, tokens *[]auth.Token) error {
	entities := make([]TokenEntity, len(*tokens))
	for i, token := range *tokens {
		entities[i] = FromDomain(&token)
	}
	if err := r.GetDB(ctx).Create(entities).Error; err != nil {
		return errors.AuthError("Failed to save token pairs", err)
	}
	return nil
}
