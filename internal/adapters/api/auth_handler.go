package api

import (
	"net/http"

	resp "github.com/ouz/goauthboilerplate/internal/adapters/api/response"
	"github.com/ouz/goauthboilerplate/internal/adapters/api/util"
	authDto "github.com/ouz/goauthboilerplate/internal/application/auth/dto"
	authService "github.com/ouz/goauthboilerplate/internal/domain/auth"
	"github.com/ouz/goauthboilerplate/pkg/errors"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	logger      *logrus.Logger
	authService authService.AuthService
}

func NewAuthHandler(logger *logrus.Logger, authService authService.AuthService) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		authService: authService,
	}
}

func (h *AuthHandler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	var request authDto.RefreshAccessTokenRequest
	if err := util.DecodeAndValidate(r, &request); err != nil {
		resp.Error(w, err)
		return
	}

	if request.RefreshToken == "" {
		resp.Error(w, errors.BadRequestError("Refresh token is required"))
		return
	}

	tokens, err := h.authService.RefreshAccessToken(r.Context(), request.RefreshToken)
	if err != nil {
		h.logger.WithError(err).Error("Failed to refresh access token")
		resp.Error(w, err)
		return
	}

	response := authDto.TokenResponse{
		AccessToken:  tokens[0].Token,
		RefreshToken: tokens[1].Token,
	}

	resp.JSON(w, http.StatusOK, response)
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var request authDto.UserLoginRequest
	if err := util.DecodeAndValidate(r, &request); err != nil {
		resp.Error(w, err)
		return
	}

	tokens, err := h.authService.Login(r.Context(), request.Email, request.Password)
	if err != nil {
		h.logger.WithError(err).WithField("email", request.Email).Error("Failed to login user")
		resp.Error(w, err)
		return
	}

	response := authDto.TokenResponse{
		AccessToken:  tokens[0].Token,
		RefreshToken: tokens[1].Token,
	}

	resp.JSON(w, http.StatusOK, response)
}

func (h *AuthHandler) LoginAnonymousUser(w http.ResponseWriter, r *http.Request) {
	var request authDto.AnonymousUserLoginRequest
	if err := util.DecodeAndValidate(r, &request); err != nil {
		resp.Error(w, err)
		return
	}

	tokens, err := h.authService.LoginAnonymous(r.Context(), request.Email)
	if err != nil {
		h.logger.WithError(err).WithField("email", request.Email).Error("Failed to login anonymous user")
		resp.Error(w, err)
		return
	}

	response := authDto.TokenResponse{
		AccessToken:  tokens[0].Token,
		RefreshToken: tokens[1].Token,
	}

	resp.JSON(w, http.StatusOK, response)
}

func (h *AuthHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	user, err := util.GetAuthenticatedUser(r)
	if err != nil {
		resp.Error(w, err)
		return
	}

	if err := h.authService.Logout(r.Context(), user.ID); err != nil {
		h.logger.WithError(err).WithField("userID", user.ID).Error("Failed to logout user")
		resp.Error(w, err)
		return
	}

	resp.JSON(w, http.StatusOK, nil)
}

func (h *AuthHandler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	user, err := util.GetAuthenticatedUser(r)
	if err != nil {
		resp.Error(w, err)
		return
	}

	if err := h.authService.LogoutAll(r.Context(), user.ID); err != nil {
		h.logger.WithError(err).WithField("userID", user.ID).Error("Failed to logout all sessions")
		resp.Error(w, err)
		return
	}

	resp.JSON(w, http.StatusOK, nil)
}
