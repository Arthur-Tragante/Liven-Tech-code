package middlewares_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/arthur-tragante/liven-code-test/middlewares"
)

func generateToken(secret string, userID uint, exp time.Time) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	secret := "testsecret"
	r.Use(middlewares.AuthMiddleware(secret))

	r.GET("/test", func(c *gin.Context) {
		userID, _ := c.Get("userID")
		c.JSON(http.StatusOK, gin.H{"userID": userID})
	})

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "No Authorization Header",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   map[string]interface{}{"error": "Authorization header required"},
		},
		{
			name:           "Invalid Authorization Header Format",
			token:          "InvalidToken",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   map[string]interface{}{"error": "Authorization header format must be Bearer {token}"},
		},
		{
			name:           "Invalid Token",
			token:          "Bearer InvalidToken",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   map[string]interface{}{"error": "Invalid token"},
		},
		{
			name:           "Valid Token",
			token:          "",
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"userID": float64(123)},
		},
	}

	// Generate a valid token for the "Valid Token" test case
	token, err := generateToken(secret, 123, time.Now().Add(time.Hour))
	require.NoError(t, err)
	tests[3].token = "Bearer " + token

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			require.NoError(t, err)

			// Check if expected body is a subset of the actual response body
			for key, value := range tt.expectedBody {
				assert.Equal(t, value, responseBody[key])
			}
		})
	}
}
