package app

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	// change working directory to project root
	os.Chdir("../")
	Setup()
}

func TestIndex(t *testing.T) {
	tests := []struct {
		name         string
		route        string
		expectedCode int
	}{
		{
			name:         "index",
			route:        "/",
			expectedCode: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := App.Test(httptest.NewRequest("GET", tt.route, nil))
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}

func Test404(t *testing.T) {
	tests := []struct {
		name         string
		route        string
		expectedCode int
	}{
		{
			name:         "404",
			route:        "/404",
			expectedCode: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := App.Test(httptest.NewRequest("GET", tt.route, nil))
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}