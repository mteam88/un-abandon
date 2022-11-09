package app

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	// change working directory to project root
	DashboardSetup()
}

func TestDashboard(t *testing.T) {
	tests := []struct {
		name         string
		route        string
		auth         bool // if true, the request will be authenticated
		expectedCode int
	}{
		{
			name:         "dashboard unauthenticated",
			route:        "/dashboard",
			auth:         false,
			expectedCode: 302, // redirects to /install
		},
//		{
//			name:         "dashboard authenticated",
//			route:        "/dashboard",
//			auth:         true,
//			expectedCode: 200,
//		},
// TODO: Figure how to test this function when user is authenticated
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := App.Test(httptest.NewRequest("GET", tt.route, nil))
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}
