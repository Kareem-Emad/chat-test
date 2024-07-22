package main

import "time"

type User struct {
	Username string
	Password string
}

type Message struct {
	Sender    string    `json:"sender"`
	Recipient string    `json:"recipient"`
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
}
