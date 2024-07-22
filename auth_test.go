package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(AuthMiddleware())

	// Sample handler to test middleware
	router.GET("/protected", func(c *gin.Context) {
		username := c.MustGet("username").(string)
		c.JSON(http.StatusOK, gin.H{"username": username})
	})

	t.Run("Missing Authorization Header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.JSONEq(t, `{"error":"Authorization header is missing"}`, resp.Body.String())
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalidtoken")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.JSONEq(t, `{"error":"Invalid or missing token"}`, resp.Body.String())
	})

	t.Run("Valid Token", func(t *testing.T) {
		token, _ := generateJWT("testuser")
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, `{"username":"testuser"}`, resp.Body.String())
	})
}

func TestGenerateJWT(t *testing.T) {
	token, err := generateJWT("testuser")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate the token
	username, err := extractUsernameFromToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", username)
}

func TestExtractUsernameFromToken(t *testing.T) {
	token, err := generateJWT("testuser")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	t.Run("Valid Token", func(t *testing.T) {
		username, err := extractUsernameFromToken(token)
		assert.NoError(t, err)
		assert.Equal(t, "testuser", username)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		_, err := extractUsernameFromToken("invalidtoken")
		assert.Error(t, err)
	})

	t.Run("Expired Token", func(t *testing.T) {
		expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": "testuser",
			"exp":      time.Now().Add(-time.Hour).Unix(),
		})
		tokenStr, _ := expiredToken.SignedString(jwtSecret)

		_, err := extractUsernameFromToken(tokenStr)
		assert.Error(t, err)
	})
}
