package main

import (
	"fmt"
	"os"
)

type dbConfig struct {
	user     string
	password string
	host     string
	port     string
	name     string
}

func loadDBConfig() (dbConfig, error) {
	config := dbConfig{
		user:     os.Getenv("NUTMEG_DB_USER"),
		password: os.Getenv("NUTMEG_DB_PASSWORD"),
		host:     os.Getenv("NUTMEG_DB_HOST"),
		port:     getEnvOrDefault("NUTMEG_DB_PORT", "5432"),
		name:     os.Getenv("NUTMEG_DB_NAME"),
	}

	if config.user == "" {
		return dbConfig{}, fmt.Errorf("NUTMEG_DB_USERが設定されていません")
	}
	if config.password == "" {
		return dbConfig{}, fmt.Errorf("NUTMEG_DB_PASSWORDが設定されていません")
	}
	if config.host == "" {
		return dbConfig{}, fmt.Errorf("NUTMEG_DB_HOSTが設定されていません")
	}
	if config.name == "" {
		return dbConfig{}, fmt.Errorf("NUTMEG_DB_NAMEが設定されていません")
	}

	return config, nil
}

func getEnvOrDefault(name, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}

	return value
}
