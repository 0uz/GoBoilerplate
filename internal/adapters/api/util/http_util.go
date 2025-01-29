package util

import (
	"context"
	"net/http"

	"github.com/ouz/goauthboilerplate/internal/domain/auth"
	"github.com/ouz/goauthboilerplate/internal/domain/user"
)

const AuthenticatedUserKey string = "auth_user"
const ClientHeader string = "x-client-key"
const ClientKey string = "client"

func GetClient(ctx context.Context) *auth.Client {
	return ctx.Value(ClientKey).(*auth.Client)
}

func ExtractClientSecret(r *http.Request) string {
	return r.Header.Get(ClientHeader)
}

func GetAuthenticatedUser(r *http.Request) *user.User {
	return r.Context().Value(AuthenticatedUserKey).(*user.User)
}
