package user

import (
	"github.com/ouz/gobackend/errors"
	entity "github.com/ouz/gobackend/pkg/entities"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindById(id string) (*entity.User, error)
	Create(user *entity.User) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]entity.User, error) {
	var users []entity.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, errors.InternalError("Failed to fetch users", err)
	}
	return users, nil
}

func (r *repository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Credentials").Where("username = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundError("User not found", err)
		}
		return nil, errors.InternalError("Failed to fetch user by email", err)
	}
	return &user, nil
}

func (r *repository) FindById(id string) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Credentials").Preload("Roles").Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundError("User not found", err)
		}
		return nil, errors.InternalError("Failed to fetch user by ID", err)
	}
	return &user, nil
}

func (r *repository) Create(user *entity.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return errors.InternalError("Failed to create user", err)
	}
	return nil
}
