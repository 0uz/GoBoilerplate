package util

import (
	"context"
	"net/http"

	"github.com/ouz/goboilerplate/internal/domain/auth"
	"github.com/ouz/goboilerplate/internal/domain/user"
	"github.com/ouz/goboilerplate/pkg/errors"
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
