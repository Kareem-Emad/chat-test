package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

var (
	serviceRegistry *ServiceRegistry
)

func main() {
	var err error
	serviceRegistry, err = InitServiceRegistery()
	if err != nil {
		log.Fatalf("Failed to initialize service registry: %v", err)
	}

	defer serviceRegistry.datastore.Close()

	router := gin.Default()

	router.POST("/register", register)
	router.POST("/login", login)

	protected := router.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.POST("/send", sendMessage)
		protected.GET("/messages", getMessages)
	}
	router.Run(":8080")
}
