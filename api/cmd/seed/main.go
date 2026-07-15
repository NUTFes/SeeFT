package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const pingTimeout = 10 * time.Second

type dbConfig struct {
	user     string
	password string
	host     string
	port     string
	name     string
}

func main() {
	config := loadDBConfig()

	if err := validateDBConfig(config); err != nil {
		log.Fatalf("データベース接続設定が不正です: %v", err)
	}

	db, err := openDB(config)
	if err != nil {
		log.Fatalf("データベース接続の初期化に失敗しました: %v", err)
	}
	defer closeDB(db)

	if err := pingDB(db, pingTimeout); err != nil {
		log.Fatalf(
			"データベースへの接続に失敗しました: host=%s port=%s db=%s user=%s: %v",
			config.host,
			config.port,
			config.name,
			config.user,
			err,
		)
	}

	log.Printf(
		"データベースへの接続に成功しました: host=%s port=%s db=%s user=%s",
		config.host,
		config.port,
		config.name,
		config.user,
	)
}

func loadDBConfig() dbConfig {
	return dbConfig{
		user:     os.Getenv("NUTMEG_DB_USER"),
		password: os.Getenv("NUTMEG_DB_PASSWORD"),
		host:     os.Getenv("NUTMEG_DB_HOST"),
		port:     getEnvOrDefault("NUTMEG_DB_PORT", "5432"),
		name:     os.Getenv("NUTMEG_DB_NAME"),
	}
}

func validateDBConfig(config dbConfig) error {
	switch {
	case config.user == "":
		return fmt.Errorf("NUTMEG_DB_USERが設定されていません")
	case config.password == "":
		return fmt.Errorf("NUTMEG_DB_PASSWORDが設定されていません")
	case config.host == "":
		return fmt.Errorf("NUTMEG_DB_HOSTが設定されていません")
	case config.name == "":
		return fmt.Errorf("NUTMEG_DB_NAMEが設定されていません")
	default:
		return nil
	}
}

func openDB(config dbConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", buildDatabaseURL(config))
	if err != nil {
		return nil, fmt.Errorf("PostgreSQLドライバーの初期化に失敗しました: %w", err)
	}

	return db, nil
}

func buildDatabaseURL(config dbConfig) string {
	databaseURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(config.user, config.password),
		Host:   net.JoinHostPort(config.host, config.port),
		Path:   config.name,
	}

	query := databaseURL.Query()
	query.Set("sslmode", "disable")
	databaseURL.RawQuery = query.Encode()

	return databaseURL.String()
}

func pingDB(db *sql.DB, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("PostgreSQLへのPingに失敗しました: %w", err)
	}

	return nil
}

func closeDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("データベース接続の終了に失敗しました: %v", err)
	}
}

func getEnvOrDefault(name, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}

	return value
}
