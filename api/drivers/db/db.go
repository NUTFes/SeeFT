package db

import (
	"fmt"
	"github.com/NUTFes/seeft-api/lib/config"
	"github.com/NUTFes/seeft-api/lib/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db  *gorm.DB
	err error
)

func Init() {
	config.NewConfig()
	c := config.C
	conn := fmt.Sprintf("%v:%v@tcp(%v)/%v?parseTime=true", c.Database.Root, c.Database.Password, c.Database.Host, c.Database.Dbname)
	db, err = gorm.Open("mysql", conn)

	if err != nil {
		panic(err)
	}

	autoMigration()
}

func GetDB() *gorm.DB {
	return db
}

func Close() {
	if err := db.Close(); err != nil {
		panic(err)
	}
}

func autoMigration() {
	db.AutoMigrate(&entity.User{})
}