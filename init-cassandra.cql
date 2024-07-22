CREATE KEYSPACE chat WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

USE chat;

CREATE TABLE users (
    username TEXT PRIMARY KEY,
    password TEXT
);

CREATE TABLE messages (
    sender TEXT,
    recipient TEXT,
    chat TEXT,
    timestamp TIMESTAMP,
    content TEXT,
    PRIMARY KEY (chat, timestamp)
)  WITH CLUSTERING ORDER BY (timestamp DESC);
