package auth

import (
	"reflect"
	"testing"

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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := auth.NewTokenPair(tt.args.userID, tt.args.clientType, tt.args.jwtConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTokenPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTokenPair() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenPair_ToTokenSlice(t *testing.T) {
	tests := []struct {
		name string
		tp   *auth.TokenPair
		want []auth.Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tp.ToTokenSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TokenPair.ToTokenSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
