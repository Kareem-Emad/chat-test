package main

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockDatastore struct {
	mock.Mock
}

func (m *MockDatastore) CheckUserExists(username string) (bool, error) {
	args := m.Called(username)
	return args.Bool(0), args.Error(1)
}

func (m *MockDatastore) CreateUser(username, hashedPassword string) error {
	args := m.Called(username, hashedPassword)
	return args.Error(0)
}

func (m *MockDatastore) Login(username string, password string) (bool, error) {
	args := m.Called(username, password)
	return args.Bool(0), args.Error(1)
}

func (m *MockDatastore) SendMessage(sender, recipient, content string, timestamp time.Time) error {
	args := m.Called(sender, recipient, content, timestamp)
	return args.Error(0)
}

func (m *MockDatastore) GetMessages(chat string, timestamp time.Time) ([]Message, error) {
	args := m.Called(chat, timestamp)
	return args.Get(0).([]Message), args.Error(1)
}
func (m *MockDatastore) Close() {
	m.Called()
}

func (m *MockDatastore) Connect(host, keyspace string) error {
	args := m.Called(host, keyspace)
	return args.Error(0)
}

type MockCache struct {
	mock.Mock
}

func (m *MockCache) GetCachedMessages(chat string) (string, error) {
	args := m.Called(chat)
	return args.String(0), args.Error(1)
}

func (m *MockCache) CacheMessages(chat string, messages []Message) error {
	args := m.Called(chat, messages)
	return args.Error(0)
}

func (m *MockCache) Connect(host, port string) error {
	args := m.Called(host, port)
	return args.Error(0)
}
