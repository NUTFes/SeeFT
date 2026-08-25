package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"
	"time"

	_ "github.com/lib/pq"
)

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
	query.Set("sslmode", config.sslMode)
	databaseURL.RawQuery = query.Encode()

	return databaseURL.String()
}

func pingDB(db *sql.DB, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("PostgreSQLへのPingに失敗しました: %w", err)
	}

	log.Println("データベースへの接続に成功しました")
	return nil
}

func closeDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("データベース接続の終了に失敗しました: %v", err)
	}
}
