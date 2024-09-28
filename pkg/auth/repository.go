package auth

import (
	"github.com/ouz/gobackend/errors"
	entity "github.com/ouz/gobackend/pkg/entities"
	"gorm.io/gorm"
)

type Repository interface {
	SaveTokenPairs(tokens *[]entity.Token) error
	RevokeAllOldTokens(userID string) error
	FindClientBySecret(secret string) (*entity.Client, error)
	FindTokensByUserID(userID string) (*[]entity.Token, error)
	IsTokenRevoked(token string) (bool, error)
	CreateCredentials(credential *entity.Credential) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) SaveTokenPairs(tokens *[]entity.Token) error {
	if err := r.db.Create(tokens).Error; err != nil {
		return errors.AuthError("Failed to save token pairs", err)
	}
	return nil
}

func (r *repository) RevokeAllOldTokens(userID string) error {
	if err := r.db.Model(&entity.Token{}).
		Where("user_id = ? AND revoked = false", userID).
		Update("revoked", true).Error; err != nil {
		return errors.AuthError("Failed to revoke old tokens", err)
	}
	return nil
}

func (r *repository) FindClientBySecret(secret string) (*entity.Client, error) {
	var client entity.Client
	if err := r.db.Where("client_secret = ?", secret).First(&client).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundError("Client not found", err)
		}
		return nil, errors.InternalError("Failed to find client", err)
	}
	return &client, nil
}

func (r *repository) FindTokensByUserID(userID string) (*[]entity.Token, error) {
	var tokens []entity.Token
	if err := r.db.Preload("Client").Preload("User").
		Where("user_id = ?", userID).
		Find(&tokens).Error; err != nil {
		return nil, errors.InternalError("Failed to find tokens", err)
	}
	return &tokens, nil
}

func (r *repository) CreateCredentials(credential *entity.Credential) error {
	if err := r.db.Create(credential).Error; err != nil {
		return errors.InternalError("Failed to create credentials", err)
	}
	return nil
}

func (r *repository) IsTokenRevoked(token string) (bool, error) {
	var revoked bool
	if err := r.db.Model(&entity.Token{}).
		Where("token = ?", token).
		Select("revoked").
		First(&revoked).Error; err != nil {
		return false, errors.InternalError("Failed to check if token is revoked", err)
	}
	return revoked, nil
}
