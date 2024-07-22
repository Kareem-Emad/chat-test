package main

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
)

type Datastore interface {
	Connect(host, keyspace string) error
	Close()
	CheckUserExists(username string) (bool, error)
	CreateUser(username, hashedPassword string) error
	Login(username string, password string) (bool, error)
	SendMessage(sender, recipient, content string, timestamp time.Time) error
	GetMessages(chat string, timestamp time.Time) ([]Message, error)
}

type Session interface {
	Query(stmt string, values ...interface{}) *gocql.Query
	Close()
}
type CassandraDatastore struct {
	session Session
}

func (ds *CassandraDatastore) Connect(host, keyspace string) error {
	cluster := gocql.NewCluster(host)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum

	if ds.session != nil {
		return nil
	}

	retryInterval := 5 * time.Second
	timeout := 1 * time.Minute
	startTime := time.Now()

	for {
		var err error
		ds.session, err = cluster.CreateSession()
		if err == nil {
			fmt.Println("Connected to Cassandra")
			return nil
		}

		if time.Since(startTime) > timeout {
			fmt.Println("Failed to connect to Cassandra within timeout")
			return err
		}

		fmt.Println("Cassandra is unavailable - retrying in", retryInterval)
		time.Sleep(retryInterval)
	}
}

func (ds *CassandraDatastore) Close() {
	ds.session.Close()
}

func (ds *CassandraDatastore) CheckUserExists(username string) (bool, error) {
	var existingUsername string
	if err := ds.session.Query(`SELECT username FROM users WHERE username = ?`, username).Scan(&existingUsername); err != nil {
		if err == gocql.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (ds *CassandraDatastore) CreateUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return ds.session.Query(`INSERT INTO users (username, password) VALUES (?, ?)`, username, hashedPassword).Exec()
}

func (ds *CassandraDatastore) Login(username string, password string) (bool, error) {
	storedPassword, err := ds.getUserPassword(username)
	if err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); err != nil {
		return false, nil
	}

	return true, nil
}
func (ds *CassandraDatastore) getUserPassword(username string) (string, error) {
	var storedPassword string
	if err := ds.session.Query(`SELECT password FROM users WHERE username = ? LIMIT 1`, username).Scan(&storedPassword); err != nil {
		if err == gocql.ErrNotFound {
			return "", nil
		}
		return "", err
	}
	return storedPassword, nil
}

func (ds *CassandraDatastore) SendMessage(sender, recipient, content string, timestamp time.Time) error {
	chat := fmt.Sprintf("%s:%s", sender, recipient)
	if sender > recipient {
		chat = fmt.Sprintf("%s:%s", recipient, sender)
	}
	return ds.session.Query(`INSERT INTO messages (sender, recipient, timestamp, content, chat) VALUES (?, ?, ?, ?, ?)`,
		sender, recipient, timestamp, content, chat).Exec()
}

func (ds *CassandraDatastore) GetMessages(chat string, timestamp time.Time) ([]Message, error) {
	limit := 20
	var messages []Message
	iter := ds.session.Query(`SELECT sender, recipient, timestamp, content FROM messages WHERE chat = ? AND timestamp < ? ORDER BY timestamp DESC LIMIT ?`, chat, timestamp, limit).Iter()
	var msg Message
	for iter.Scan(&msg.Sender, &msg.Recipient, &msg.Timestamp, &msg.Content) {
		messages = append(messages, msg)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return messages, nil
}
