package main

import "log"

type ServiceRegistry struct {
	datastore Datastore
	cache     Cache
	config    ConfigManager
}

func InitServiceRegistery() (*ServiceRegistry, error) {
	cfg := ConfigManager{}
	cfg.LoadConfig()

	ds := &CassandraDatastore{}
	if err := ds.Connect(cfg.DatabaseHost, cfg.DatabaseKeyspace); err != nil {
		log.Fatalf("Failed to connect to datastore: %v", err)
	}

	cache := &RedisCache{}
	if err := cache.Connect(cfg.RedisHost, cfg.RedisPort); err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		return nil, err
	}

	return &ServiceRegistry{
		datastore: ds,
		cache:     cache,
		config:    cfg,
	}, nil
}
