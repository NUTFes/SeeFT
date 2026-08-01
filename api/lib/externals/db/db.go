package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	// "github.com/joho/godotenv"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type client struct {
	db *sql.DB
	gorm *gorm.DB
}

type Client interface {
	DB() *sql.DB
	GormDB() *gorm.DB
	CloseDB()
}

func ConnectMySQL() (client, error) {
	// err := godotenv.Load("env/dev.env")
	// if err != nil {
	// 	fmt.Println(err)
	// }
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

	// 接続数の上限を設定（無制限のままだと高負荷時にPostgresのmax_connectionsを超過し、
	// リクエストが500エラーで失敗する。ミニPC実測で400並列負荷時の500エラー率を
	// 37.8%→0%に改善することを確認済み）
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(30 * time.Minute)

	err = db.Ping()

	if err != nil {
		fmt.Println(err)
		fmt.Println("[Failed] Not Connect to PostgreSQL") // 失敗
		return client{}, err
	} else {
		fmt.Println("[Success] Connect to PostgreSQL") // 成功
		// initialize GORM using the existing *sql.DB connection (keep existing comments)
		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})
		if err != nil {
			return client{}, err
		}
		return client{db: db, gorm: gormDB}, nil
	}
}

func (c client) CloseDB() {
	if c.db != nil {
		_ = c.db.Close()
	}
}

func (c client) DB() *sql.DB {
	return c.db
}

func (c client) GormDB() *gorm.DB {
	return c.gorm
}
