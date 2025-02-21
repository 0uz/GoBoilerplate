package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
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
		{
			name: "Valid access token",
			args: args{
				userID:     uuid.New().String(),
				tokenType:  auth.ACCESS_TOKEN,
				clientType: "WEB",
				jwtSecret:  "test-secret",
				expiration: time.Hour,
			},
			wantErr: false,
		},
		{
			name: "Valid refresh token",
			args: args{
				userID:     uuid.New().String(),
				tokenType:  auth.REFRESH_TOKEN,
				clientType: "WEB",
				jwtSecret:  "test-secret",
				expiration: time.Hour * 24,
			},
			wantErr: false,
		},
		{
			name: "Empty user ID",
			args: args{
				userID:     "",
				tokenType:  auth.ACCESS_TOKEN,
				clientType: "WEB",
				jwtSecret:  "test-secret",
				expiration: time.Hour,
			},
			wantErr: false, // UserID validation is not implemented
		},
		{
			name: "Empty JWT secret",
			args: args{
				userID:     uuid.New().String(),
				tokenType:  auth.ACCESS_TOKEN,
				clientType: "WEB",
				jwtSecret:  "",
				expiration: time.Hour,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := auth.NewToken(tt.args.userID, tt.args.tokenType, tt.args.clientType, tt.args.jwtSecret, tt.args.expiration)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewToken() = nil, want non-nil")
					return
				}
				if got.UserID != tt.args.userID {
					t.Errorf("NewToken().UserID = %v, want %v", got.UserID, tt.args.userID)
				}
				if got.TokenType != tt.args.tokenType {
					t.Errorf("NewToken().TokenType = %v, want %v", got.TokenType, tt.args.tokenType)
				}
				if got.ClientType != tt.args.clientType {
					t.Errorf("NewToken().ClientType = %v, want %v", got.ClientType, tt.args.clientType)
				}
				if got.Token == "" {
					t.Error("NewToken().Token is empty")
				}
			}
		})
	}
}

func TestToken_Validate(t *testing.T) {
	validToken, _ := auth.NewToken(uuid.New().String(), auth.ACCESS_TOKEN, "WEB", "test-secret", time.Hour)
	expiredToken, _ := auth.NewToken(uuid.New().String(), auth.ACCESS_TOKEN, "WEB", "test-secret", -time.Hour)
	revokedToken, _ := auth.NewToken(uuid.New().String(), auth.ACCESS_TOKEN, "WEB", "test-secret", time.Hour)
	revokedToken.Revoke()

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
		{
			name: "Valid token",
			tr:   validToken,
			args: args{
				jwtSecret: "test-secret",
			},
			wantErr: false,
		},
		{
			name: "Expired token",
			tr:   expiredToken,
			args: args{
				jwtSecret: "test-secret",
			},
			wantErr: true,
		},
		{
			name: "Revoked token",
			tr:   revokedToken,
			args: args{
				jwtSecret: "test-secret",
			},
			wantErr: true,
		},
		{
			name: "Wrong secret",
			tr:   validToken,
			args: args{
				jwtSecret: "wrong-secret",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.Validate(tt.args.jwtSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("Token.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("Token.Validate() = nil, want non-nil")
					return
				}
				if got.UserId != tt.tr.UserID {
					t.Errorf("Token.Validate().UserId = %v, want %v", got.UserId, tt.tr.UserID)
				}
			}
		})
	}
}

func TestToken_Revoke(t *testing.T) {
	token, _ := auth.NewToken(uuid.New().String(), auth.ACCESS_TOKEN, "WEB", "test-secret", time.Hour)

	tests := []struct {
		name string
		tr   *auth.Token
	}{
		{
			name: "Revoke token",
			tr:   token,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.tr.Revoked {
				t.Error("Token is already revoked before test")
			}
			tt.tr.Revoke()
			if !tt.tr.Revoked {
				t.Error("Token.Revoke() did not revoke the token")
			}
		})
	}
}

func TestToken_IsExpired(t *testing.T) {
	validToken, _ := auth.NewToken(uuid.New().String(), auth.ACCESS_TOKEN, "WEB", "test-secret", time.Hour)
	expiredToken, _ := auth.NewToken(uuid.New().String(), auth.ACCESS_TOKEN, "WEB", "test-secret", -time.Hour)

	tests := []struct {
		name string
		tr   *auth.Token
		want bool
	}{
		{
			name: "Valid token",
			tr:   validToken,
			want: false,
		},
		{
			name: "Expired token",
			tr:   expiredToken,
			want: true,
		},
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
	userID := uuid.New().String()
	clientType := "WEB"
	accessToken, _ := auth.NewToken(userID, auth.ACCESS_TOKEN, clientType, "test-secret", time.Hour)
	refreshToken, _ := auth.NewToken(userID, auth.REFRESH_TOKEN, clientType, "test-secret", time.Hour)

	tests := []struct {
		name string
		tr   *auth.Token
		want string
	}{
		{
			name: "Access token cache key",
			tr:   accessToken,
			want: "uat:" + userID + ":" + clientType,
		},
		{
			name: "Refresh token cache key",
			tr:   refreshToken,
			want: "urt:" + userID + ":" + clientType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.GetCacheKey(); got != tt.want {
				t.Errorf("Token.GetCacheKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
