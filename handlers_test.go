package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	mockDatastore *MockDatastore
	mockCache     *MockCache
	router        *gin.Engine
)

func setupRouter() {
	gin.SetMode(gin.TestMode)
	router = gin.Default()

	serviceRegistry = &ServiceRegistry{
		datastore: mockDatastore,
		cache:     mockCache,
		config:    ConfigManager{},
	}

	router.POST("/register", register)
	router.POST("/login", login)
	router.POST("/send", AuthMiddleware(), sendMessage)
	router.GET("/messages", AuthMiddleware(), getMessages)
}

func TestRegister(t *testing.T) {
	mockDatastore = new(MockDatastore)
	mockCache = new(MockCache)

	setupRouter()

	mockDatastore.On("CheckUserExists", "testuser").Return(false, nil)
	mockDatastore.On("CreateUser", "testuser", mock.Anything).Return(nil)

	body := `{"username": "testuser", "password": "testpass"}`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"message": "User registered successfully"}`, rec.Body.String())

	mockDatastore.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	mockDatastore = new(MockDatastore)
	mockCache = new(MockCache)

	setupRouter()

	mockDatastore.On("Login", "testuser", "testpass").Return(true, nil)

	token, _ := generateJWT("testuser")

	body := `{"username": "testuser", "password": "testpass"}`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"token":`)
	assert.Contains(t, rec.Body.String(), token)

	mockDatastore.AssertExpectations(t)
}

func TestSendMessage(t *testing.T) {
	mockDatastore = new(MockDatastore)
	mockCache = new(MockCache)

	setupRouter()

	mockDatastore.On("CheckUserExists", "recipientuser").Return(true, nil)
	mockDatastore.On("SendMessage", "test", "recipientuser", "Hello", mock.Anything).Return(nil)

	body := `{"recipient": "recipientuser", "content": "Hello"}`
	req, _ := http.NewRequest("POST", "/send", strings.NewReader(body))
	token, _ := generateJWT("test")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"message": "Message sent successfully"}`, rec.Body.String())

	mockDatastore.AssertExpectations(t)
}

func TestGetMessages(t *testing.T) {
	mockDatastore = new(MockDatastore)
	mockCache = new(MockCache)

	setupRouter()

	messages := []Message{
		{Sender: "sender", Recipient: "recipient", Timestamp: time.Now(), Content: "Hello"},
	}

	mockCache.On("GetCachedMessages", "recipient:sender").Return("", gocql.ErrNotFound)
	mockDatastore.On("GetMessages", "recipient:sender", mock.Anything).Return(messages, nil)
	mockCache.On("CacheMessages", "recipient:sender", mock.Anything).Return(nil)

	req, _ := http.NewRequest("GET", "/messages?recipient=recipient", nil)
	token, _ := generateJWT("sender")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	mockDatastore.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
