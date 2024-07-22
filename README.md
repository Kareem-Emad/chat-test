# Chat System

This is a simple chat system that supports messaging between two people. It is designed to allow users to send and receive messages in a one-to-one chat environment. This README file provides details on how to set up, run, and understand the project.

## Getting Started

To get started with the chat system, follow these steps:

### Running the Project

1. **Build and Start the Application**

   Use Docker Compose to build and run the application:

   ```sh
   docker-compose up --build
   ```
   This command will build the Docker images and start the containers defined in the docker-compose.yml file.

2. **Run Tests**
   To run the tests, use the following command:

   ```sh
   go test ./...
   ```

## Functionalities
The chat system supports the following functionalities:

- User Registration: Allows users to create new accounts.
- User Login: Authenticates users and generates JWT tokens.
- Send Message: Enables users to send messages to each other.
- Get Messages: Retrieves messages between two users, with support for pagination.

Note: This system currently supports one-to-one messaging only. Group chats or any form of multi-user messaging are not supported.

## Database Design
The database for this chat system uses Apache Cassandra. The design includes the following key components:

1. Keyspace

```
CREATE KEYSPACE chat WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
```
The keyspace chat is created with a simple replication strategy and a replication factor of 1, suitable for a development environment.

NOTE: this should be configured differently for a production setting to achieve max performance 


2. Tables

   - Users Table : The users table stores user credentials. The username field serves as the primary key.

   - Messages Table : The messages table stores chat messages. The primary key is a composite key consisting of chat and timestamp
     - The primary key is designed to optimize the retrieval of messages. By using chat as the partition key and timestamp as the clustering key:
     - Partition Key (chat): Ensures that all messages between a specific pair of users are stored together, which optimizes read and write operations for messages between those users.
     - Clustering Key (timestamp): Orders messages within each partition by timestamp in descending order, allowing efficient retrieval of recent messages.



## Code Structure
The project is organized to maintain clean separation of concerns and to facilitate dependency injection:

- main.go: Entry point of the application. Sets up the HTTP server and routes.

- handlers.go: Contains HTTP handler functions for user registration, login, sending messages, and retrieving messages.

- datastore.go: Defines the Datastore interface and its implementation for Cassandra. Handles interactions with the database.

- cache.go: Defines the Cache interface and its implementation for Redis. Handles caching of messages.

- service_registry.go: Manages the initialization of services and dependency injection. Ensures that all components are properly configured and wired together.


## API Documentation
API documentation is available in Swagger format. You can access the Swagger file here (swagger.yml). It provides detailed information about all the available endpoints, their parameters, and responses.

### Sample API Calls

```
curl --location 'http://localhost:8080/register' \
--header 'Content-Type: application/json' \
--data '{
      "username": "myuser",
      "password": "1234"
    }'

```
```
curl --location 'http://localhost:8080/login' \
--header 'Content-Type: application/json' \
--data '{
      "username": "myuser",
      "password": "1234"
    }'

```
```
curl --location 'http://localhost:8080/send' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <your_jwt_token>' \
--data '{
      "recipient": "recipientuser",
      "content": "Hello, how are you?"
    }'

```
```
curl --location 'http://localhost:8080/messages?recipient=recipientuser' \
--header 'Authorization: Bearer <your_jwt_token>'
```

 