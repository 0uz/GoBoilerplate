package util

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ouz/goauthboilerplate/internal/adapters/api/validator"
	"github.com/ouz/goauthboilerplate/internal/domain/auth"
	"github.com/ouz/goauthboilerplate/internal/domain/user"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

type ContextKey string

const AuthenticatedUserKey ContextKey = "auth_user"
const ClientHeader string = "x-client-key"
const ClientKey ContextKey = "client"

func GetClient(ctx context.Context) (auth.Client, error) {
	rawClient := ctx.Value(ClientKey)
	if rawClient == nil {
		return auth.Client{}, errors.InternalError("Failed to get client", nil)
	}
	client, ok := rawClient.(auth.Client)
	if !ok {
		return auth.Client{}, errors.InternalError("Failed to convert client", nil)
	}
	return client, nil
}

func ExtractClientSecret(r *http.Request) string {
	return r.Header.Get(ClientHeader)
}

func GetAuthenticatedUser(r *http.Request) (user.User, error) {
	rawUser := r.Context().Value(AuthenticatedUserKey)
	if rawUser == nil {
		return user.User{}, errors.UnauthorizedError("User not found", nil)
	}
	u, ok := rawUser.(user.User)
	if !ok {
		return user.User{}, errors.InternalError("Failed to convert user", nil)
	}
	return u, nil
}

func DecodeAndValidate(r *http.Request, request interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		return errors.BadRequestError("Invalid request body")
	}

	if err := validator.Validator.Struct(request); err != nil {
		return errors.BadRequestError(err.Error())
	}

	return nil
}
