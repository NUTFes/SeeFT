package db

import (
  "fmt"
  // "database/sql"
  // "gorm.io/driver/mysql"
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
  "os"
)


type gormClient struct {
	db *gorm.DB
}

type GormClient interface {
	DB() *gorm.DB
}

func (c gormClient) DB() *gorm.DB {
	return c.db
}

func ConnectMySQLFromGorm() (GormClient, error) {
	dbUser := os.Getenv("NUTMEG_DB_USER")
	dbPassword := os.Getenv("NUTMEG_DB_PASSWORD")
	dbHost := os.Getenv("NUTMEG_DB_HOST")
	dbPort := os.Getenv("NUTMEG_DB_PORT")
	dbName := os.Getenv("NUTMEG_DB_NAME")
	dns := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"
	fmt.Println(dns)
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})

	if err != nil {
		fmt.Println("[Failed] Not Connect to PostgreSQL form Grom")
		return nil, err
	} else {
		fmt.Println("[Success] Connect to PostgreSQL form Grom")
		return gormClient{db}, nil
	}
}
