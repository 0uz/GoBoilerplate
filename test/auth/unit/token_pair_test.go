package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ouz/goauthboilerplate/internal/config"
	"github.com/ouz/goauthboilerplate/internal/domain/auth"
)

func TestNewTokenPair(t *testing.T) {
	type args struct {
		userID     string
		clientType string
		jwtConfig  config.JWTConfig
	}
	tests := []struct {
		name    string
		args    args
		want    *auth.TokenPair
		wantErr bool
	}{
		{
			name: "Valid token pair",
			args: args{
				userID:     uuid.New().String(),
				clientType: "WEB",
				jwtConfig: config.JWTConfig{
					Secret:            "test-secret",
					AccessExpiration:  time.Hour,
					RefreshExpiration: time.Hour * 24,
				},
			},
			wantErr: false,
		},
		{
			name: "Empty JWT secret",
			args: args{
				userID:     uuid.New().String(),
				clientType: "WEB",
				jwtConfig: config.JWTConfig{
					Secret:            "",
					AccessExpiration:  time.Hour,
					RefreshExpiration: time.Hour * 24,
				},
			},
			wantErr: true,
		},
		{
			name: "Zero access expiration",
			args: args{
				userID:     uuid.New().String(),
				clientType: "WEB",
				jwtConfig: config.JWTConfig{
					Secret:            "test-secret",
					AccessExpiration:  0,
					RefreshExpiration: time.Hour * 24,
				},
			},
			wantErr: false, // Expiration validation is not implemented
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := auth.NewTokenPair(tt.args.userID, tt.args.clientType, tt.args.jwtConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTokenPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewTokenPair() = nil, want non-nil")
					return
				}
				if got.AccessToken.UserID != tt.args.userID {
					t.Errorf("NewTokenPair().AccessToken.UserID = %v, want %v", got.AccessToken.UserID, tt.args.userID)
				}
				if got.RefreshToken.UserID != tt.args.userID {
					t.Errorf("NewTokenPair().RefreshToken.UserID = %v, want %v", got.RefreshToken.UserID, tt.args.userID)
				}
				if got.AccessToken.TokenType != auth.ACCESS_TOKEN {
					t.Errorf("NewTokenPair().AccessToken.TokenType = %v, want %v", got.AccessToken.TokenType, auth.ACCESS_TOKEN)
				}
				if got.RefreshToken.TokenType != auth.REFRESH_TOKEN {
					t.Errorf("NewTokenPair().RefreshToken.TokenType = %v, want %v", got.RefreshToken.TokenType, auth.REFRESH_TOKEN)
				}
			}
		})
	}
}

func TestTokenPair_ToTokenSlice(t *testing.T) {
	userID := uuid.New().String()
	clientType := "WEB"
	jwtConfig := config.JWTConfig{
		Secret:            "test-secret",
		AccessExpiration:  time.Hour,
		RefreshExpiration: time.Hour * 24,
	}

	tokenPair, _ := auth.NewTokenPair(userID, clientType, jwtConfig)

	tests := []struct {
		name string
		tp   *auth.TokenPair
		want []auth.Token
	}{
		{
			name: "Convert to slice",
			tp:   tokenPair,
			want: []auth.Token{tokenPair.AccessToken, tokenPair.RefreshToken},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tp.ToTokenSlice()
			if len(got) != 2 {
				t.Errorf("TokenPair.ToTokenSlice() length = %v, want %v", len(got), 2)
				return
			}
			if got[0].TokenType != auth.ACCESS_TOKEN {
				t.Errorf("First token type = %v, want %v", got[0].TokenType, auth.ACCESS_TOKEN)
			}
			if got[1].TokenType != auth.REFRESH_TOKEN {
				t.Errorf("Second token type = %v, want %v", got[1].TokenType, auth.REFRESH_TOKEN)
			}
		})
	}
}
