package app

import (
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAuthenticateUser(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Figure how to test this function
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestCheckGHOauthToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name   string
		args   args
		wantOk bool
	}{
		{
			name:   "fail",
			args:   args{token: "obviously_invalid_token"},
			wantOk: false,
		},
		// TODO: test with valid token
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOk := CheckGHOauthToken(tt.args.token); gotOk != tt.wantOk {
				t.Errorf("CheckGHOauthToken() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
