package config

import "os"

type Config struct {
	Port    string
	Storage string
	DBUrl   string
}

func Load() Config {
	cfg := Config{
		Port:    os.Getenv("PORT"),
		Storage: os.Getenv("STORAGE"),
		DBUrl:   os.Getenv("DB_URL"),
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.Storage == "" {
		cfg.Storage = "memory"
	}
	return cfg
}
