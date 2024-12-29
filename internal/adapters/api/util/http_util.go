package util

import (
	"net/http"

	"github.com/ouz/goauthboilerplate/internal/domain/user"
)

const AuthenticatedUserKey string = "auth_user"

func ExtractClientSecret(r *http.Request) string {
	return r.Header.Get("x-client-key")
}

func GetAuthenticatedUser(r *http.Request) *user.User {
	return r.Context().Value(AuthenticatedUserKey).(*user.User)
}
