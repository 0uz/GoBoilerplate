package auth

import (
	"reflect"
	"testing"
	"time"

	"github.com/ouz/goauthboilerplate/internal/domain/auth"
)

func TestNewToken(t *testing.T) {
	type args struct {
		userID     string
		tokenType  auth.TokenType
		clientType string
		jwtSecret  string
		expiration time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    *auth.Token
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := auth.NewToken(tt.args.userID, tt.args.tokenType, tt.args.clientType, tt.args.jwtSecret, tt.args.expiration)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToken_Validate(t *testing.T) {
	type args struct {
		jwtSecret string
	}
	tests := []struct {
		name    string
		tr      *auth.Token
		args    args
		want    *auth.TokenClaims
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.Validate(tt.args.jwtSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("Token.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Token.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToken_Revoke(t *testing.T) {
	tests := []struct {
		name string
		tr   *auth.Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.Revoke()
		})
	}
}

func TestToken_IsExpired(t *testing.T) {
	tests := []struct {
		name string
		tr   *auth.Token
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.IsExpired(); got != tt.want {
				t.Errorf("Token.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToken_GetCacheKey(t *testing.T) {
	tests := []struct {
		name string
		tr   *auth.Token
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.GetCacheKey(); got != tt.want {
				t.Errorf("Token.GetCacheKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
