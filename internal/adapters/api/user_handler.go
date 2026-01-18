package api

import (
	"net/http"

	"github.com/ouz/goboilerplate/internal/adapters/api/util"
	authDto "github.com/ouz/goboilerplate/internal/application/auth/dto"
	userDto "github.com/ouz/goboilerplate/internal/application/user/dto"
	"github.com/ouz/goboilerplate/internal/domain/user"
	"github.com/ouz/goboilerplate/pkg/log"
	resp "github.com/ouz/goboilerplate/pkg/response"
)

const (
	emailConfirmationTemplatePath = "internal/adapters/api/template/email_confirmation_response.html"
	notFoundTemplatePath          = "internal/adapters/api/template/not_found.html"
)

type UserHandler struct {
	logger      *log.Logger
	userService user.UserService
}

func NewUserHandler(logger *log.Logger, userService user.UserService) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
	}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var request authDto.UserRegisterRequest
	if err := resp.DecodeAndValidate(r, &request); err != nil {
		resp.Error(w, err)
		return
	}

	if err := h.userService.Register(r.Context(), request); err != nil {
		h.logger.Error("Failed to register user", "error", err, "email", request.Email)
		resp.Error(w, err)
		return
	}

	resp.JSON(w, http.StatusCreated, nil)
}

func (h *UserHandler) RegisterAnonymousUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.RegisterAnonymousUser(r.Context())
	if err != nil {
		h.logger.Error("Failed to register anonymous user", "error", err)
		resp.Error(w, err)
		return
	}

	response := authDto.AnonymousUserResponse{
		Email: user.Email,
	}

	resp.JSON(w, http.StatusCreated, response)
}

func (h *UserHandler) ConfirmUser(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	confirmation := vars.Get("key")
	if confirmation == "" {
		returnNotFound(w, r)
		return
	}

	if err := h.userService.ConfirmUser(r.Context(), confirmation); err != nil {
		returnNotFound(w, r)
		return
	}

	http.ServeFile(w, r, emailConfirmationTemplatePath)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := util.GetAuthenticatedUser(r)
	if err != nil {
		resp.Error(w, err)
		return
	}

	resp.JSON(w, http.StatusOK, userDto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Anonymous: user.Anonymous,
	})
}

func returnNotFound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, notFoundTemplatePath)
}
