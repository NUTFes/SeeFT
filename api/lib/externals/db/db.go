package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"os"
)

type client struct {
	db *sql.DB
}

type Client interface {
	DB() *sql.DB
	CloseDB()
}

func ConnectMySQL() (client, error) {
	err := godotenv.Load("env/dev.env")
	if err != nil {
		fmt.Println(err)
	}
	dbUser := os.Getenv("NUTMEG_DB_USER")
	dbPassword := os.Getenv("NUTMEG_DB_PASSWORD")
	dbHost := os.Getenv("NUTMEG_DB_HOST")
	dbPort := os.Getenv("NUTMEG_DB_PORT")
	dbName := os.Getenv("NUTMEG_DB_NAME")
	
	// MySQLに接続する
	// データベース接続部分
	// dbconf := "seeft:password@tcp(nutfes-seeft-db:3306)/seeft_db?charset=utf8mb4&parseTime=true"
	// dbconf := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=true"
	// db, err := sql.Open("mysql", dbconf)

	dns := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"
	db, err := sql.Open("postgres", dns);

	if err != nil {
		return client{}, err
	}
	
	err = db.Ping()

	if err != nil {
		fmt.Println(err)
		fmt.Println("[Failed] Not Connect to PostgreSQL") // 失敗
		return client{}, err
	} else {
		fmt.Println("[Success] Connect to PostgreSQL") // 成功
		return client{db}, nil
	}
}

func (c client) CloseDB() {
	if c.db != nil {
		c.db.Close()
	}
}

func (c client) DB() *sql.DB {
	return c.db
}
