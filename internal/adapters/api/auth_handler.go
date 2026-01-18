package api

import (
	"net/http"

	resp "github.com/ouz/goauthboilerplate/pkg/response"
	"github.com/ouz/goauthboilerplate/internal/adapters/api/util"
	authDto "github.com/ouz/goauthboilerplate/internal/application/auth/dto"
	authService "github.com/ouz/goauthboilerplate/internal/domain/auth"
	"github.com/ouz/goauthboilerplate/pkg/errors"
	"github.com/ouz/goauthboilerplate/pkg/log"
)

type AuthHandler struct {
	logger      *log.Logger
	authService authService.AuthService
}

func NewAuthHandler(logger *log.Logger, authService authService.AuthService) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		authService: authService,
	}
}

func (h *AuthHandler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	var request authDto.RefreshAccessTokenRequest
	if err := resp.DecodeAndValidate(r, &request); err != nil {
		resp.Error(w, err)
		return
	}

	if request.RefreshToken == "" {
		resp.Error(w, errors.BadRequestError("Refresh token is required"))
		return
	}

	tokens, err := h.authService.RefreshAccessToken(r.Context(), request.RefreshToken)
	if err != nil {
		h.logger.Error("Failed to refresh access token", "error", err)
		resp.Error(w, err)
		return
	}

	response := authDto.TokenResponse{
		AccessToken:  tokens.AccessToken.RawToken,
		RefreshToken: tokens.RefreshToken.RawToken,
	}

	resp.JSON(w, http.StatusOK, response)
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var request authDto.UserLoginRequest
	if err := resp.DecodeAndValidate(r, &request); err != nil {
		resp.Error(w, err)
		return
	}

	tokens, err := h.authService.Login(r.Context(), request.Email, request.Password)
	if err != nil {
		h.logger.Error("Failed to login user", "error", err, "email", request.Email)
		resp.Error(w, err)
		return
	}

	response := authDto.TokenResponse{
		AccessToken:  tokens.AccessToken.RawToken,
		RefreshToken: tokens.RefreshToken.RawToken,
	}

	resp.JSON(w, http.StatusOK, response)
}

func (h *AuthHandler) LoginAnonymousUser(w http.ResponseWriter, r *http.Request) {
	var request authDto.AnonymousUserLoginRequest
	if err := resp.DecodeAndValidate(r, &request); err != nil {
		resp.Error(w, err)
		return
	}

	tokens, err := h.authService.LoginAnonymous(r.Context(), request.Email)
	if err != nil {
		h.logger.Error("Failed to login anonymous user", "error", err, "email", request.Email)
		resp.Error(w, err)
		return
	}

	response := authDto.TokenResponse{
		AccessToken:  tokens.AccessToken.RawToken,
		RefreshToken: tokens.RefreshToken.RawToken,
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
		h.logger.Error("Failed to logout user", "error", err, "userID", user.ID)
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
		h.logger.Error("Failed to logout all sessions", "error", err, "userID", user.ID)
		resp.Error(w, err)
		return
	}

	resp.JSON(w, http.StatusOK, nil)
}
