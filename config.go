package main

import "os"

type ConfigManager struct {
	DatabaseHost     string
	DatabaseKeyspace string
	RedisHost        string
	RedisPort        string
}

func (cm *ConfigManager) LoadConfig() {
	cm.DatabaseHost = "cassandra"
	cm.RedisHost = os.Getenv("REDIS_HOST")
	cm.RedisPort = "6379"
	cm.DatabaseKeyspace = os.Getenv("CASSANDRA_KEYSPACE")
}
