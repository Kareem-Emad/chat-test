package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ----------------------------------------------
// --------------- USER  HANDLERS ---------------
// ----------------------------------------------
func register(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the username already exists
	exists, err := serviceRegistry.datastore.CheckUserExists(user.Username)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check username"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Insert the user into the users table
	if err := serviceRegistry.datastore.CreateUser(user.Username, user.Password); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func login(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validLogin, err := serviceRegistry.datastore.Login(user.Username, user.Password)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate login"})
		return
	}
	if !validLogin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := generateJWT(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

// ----------------------------------------------
// ------------ messaging  HANDLERS -------------
// ----------------------------------------------
func sendMessage(c *gin.Context) {
	var msg Message
	if err := c.BindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, _ := c.Get("username")

	sender := username.(string)
	// not the best way, as there could be message delays, or we could use a scheduler to save messages
	msg.Timestamp = time.Now()

	exists, err := serviceRegistry.datastore.CheckUserExists(msg.Recipient)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate recipient"})

		return
	}
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Recipient not found"})
		return
	}

	if err := serviceRegistry.datastore.SendMessage(sender, msg.Recipient, msg.Content, msg.Timestamp); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}

func getMessages(c *gin.Context) {
	username, _ := c.Get("username")
	user := username.(string)

	recipient := c.Query("recipient")
	timestampStr := c.Query("timestamp")
	firstPage := timestampStr == ""
	if firstPage {
		timestampStr = time.Now().Format(time.RFC3339)
	}

	timestamp, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp"})
		return
	}

	var messages []Message
	chat := fmt.Sprintf("%s:%s", user, recipient)
	if user > recipient {
		chat = fmt.Sprintf("%s:%s", recipient, user)
	}

	if firstPage {
		cachedMessages, err := serviceRegistry.cache.GetCachedMessages(chat)

		if err == nil {
			// Return cached messages
			c.JSON(http.StatusOK, gin.H{"messages": cachedMessages})
			return
		}
	}

	messages, err = serviceRegistry.datastore.GetMessages(chat, timestamp)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}

	if firstPage {
		// Store the first page in Redis with 1-day expiration
		serviceRegistry.cache.CacheMessages(chat, messages)
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
