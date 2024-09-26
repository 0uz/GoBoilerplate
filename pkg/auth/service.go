package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ouz/gobackend/config"
	"github.com/ouz/gobackend/errors"
	entity "github.com/ouz/gobackend/pkg/entities"
	"github.com/ouz/gobackend/pkg/user"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	ValidateTokenAndGetUser(token string) (*entity.User, error)
	CreateToken(c *fiber.Ctx, user *entity.User) (*[]entity.Token, error)
	CreateCredentials(user *entity.User, password string) error
}

type service struct {
	repository  Repository
	userService user.Service
}

func NewService(repository Repository, userService user.Service) Service {
	return &service{
		repository:  repository,
		userService: userService,
	}
}

func (s *service) ValidateTokenAndGetUser(token string) (*entity.User, error) {
	t, err := jwt.ParseWithClaims(token, &entity.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.UnauthorizedError("Invalid token signing method", nil)
		}
		return []byte(config.Conf.JWT_SECRET), nil
	})

	if err != nil {
		return nil, errors.UnauthorizedError("Invalid token", err)
	}

	claims, ok := t.Claims.(*entity.TokenClaims)
	if !ok {
		return nil, errors.UnauthorizedError("Invalid token claims", nil)
	}

	user, err := s.userService.FindById(claims.ID)
	if err != nil {
		return nil, errors.NotFoundError("User not found", err)
	}

	return user, nil
}

func (s *service) CreateToken(c *fiber.Ctx, user *entity.User) (*[]entity.Token, error) {
	now := time.Now()

	accessToken, err := s.generateToken(user.ID, config.Conf.JWT_ACCESS_EXPIRATION)
	if err != nil {
		return nil, errors.AuthError("Failed to generate access token", err)
	}

	refreshToken, err := s.generateToken(user.ID, config.Conf.JWT_REFRESH_EXPIRATION)
	if err != nil {
		return nil, errors.InternalError("Failed to generate refresh token", err)
	}

	if err := s.repository.RevokeAllOldTokens(user.ID); err != nil {
		return nil, errors.InternalError("Failed to revoke old tokens", err)
	}

	client, err := s.repository.FindClientBySecret("12345")
	if err != nil {
		return nil, errors.InternalError("Failed to find client", err)
	}

	tokens := []entity.Token{
		s.createTokenEntity(accessToken, entity.ACCESS_TOKEN, user.ID, client.ClientType, c.IP(), now.Add(config.Conf.JWT_ACCESS_EXPIRATION)),
		s.createTokenEntity(refreshToken, entity.REFRESH_TOKEN, user.ID, client.ClientType, c.IP(), now.Add(config.Conf.JWT_REFRESH_EXPIRATION)),
	}

	if err := s.repository.SaveTokenPairs(&tokens); err != nil {
		return nil, errors.InternalError("Failed to save token pairs", err)
	}

	return &tokens, nil
}

func (s *service) CreateCredentials(user *entity.User, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	credential := &entity.Credential{
		UserID:         user.ID,
		Hash:           string(hashedPassword),
		CredentialType: entity.PASSWORD,
	}

	if err := s.repository.CreateCredentials(credential); err != nil {
		return err
	}

	return nil
}

func (s *service) generateToken(userID string, expiration time.Duration) (string, error) {
	claims := entity.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
		ID: userID,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.Conf.JWT_SECRET))
}

func (s *service) createTokenEntity(token string, tokenType entity.TokenType, userID, clientType, ipAddress string, expiresAt time.Time) entity.Token {
	return entity.Token{
		Token:      token,
		TokenType:  tokenType,
		Revoked:    false,
		IpAddress:  ipAddress,
		ClientType: clientType,
		UserID:     userID,
		ExpiresAt:  expiresAt,
	}
}
