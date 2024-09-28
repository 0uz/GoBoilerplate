package user

import (
	"github.com/ouz/gobackend/errors"
	entity "github.com/ouz/gobackend/pkg/entities"
)

type Service interface {
	FindAll() ([]entity.User, error)
	RegisterUser(request *entity.User) error
	RegisterAnonymousUser(request *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindById(id string) (*entity.User, error)
	FindByIdLazy(id string) (*entity.User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) FindAll() ([]entity.User, error) {
	return s.repository.FindAll()
}

func (s *service) FindByEmail(email string) (*entity.User, error) {
	return s.repository.FindByEmail(email)
}

func (s *service) FindById(id string) (*entity.User, error) {
	return s.repository.FindById(id)
}

func (s *service) FindByIdLazy(id string) (*entity.User, error) {
	return s.repository.FindByIdLazy(id)
}

func (s *service) RegisterUser(request *entity.User) error {
	if request == nil {
		return errors.ValidationError("Request is nil", nil)
	}

	existingUser, err := s.repository.FindByEmail(request.Username)

	if err != nil && !errors.IsNotFoundError(err) {
		return err
	}

	if existingUser != nil {
		return errors.ConflictError("User with this email already exists", nil)
	}

	if err := s.repository.Create(request); err != nil {
		return errors.InternalError("Failed to create user", err)
	}

	return nil
}

func (s *service) RegisterAnonymousUser(request *entity.User) error {
	if request == nil {
		return errors.ValidationError("Request is nil", nil)
	}

	if err := s.repository.Create(request); err != nil {
		return errors.InternalError("Failed to create user", err)
	}

	return nil
}
