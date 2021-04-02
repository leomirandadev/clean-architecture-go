package configs

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func ConnectGorm() *gorm.DB {

	var DB_CONNECTION string = os.Getenv("DB_CONNECTION")

	db, err := gorm.Open("mysql", DB_CONNECTION)

	if err != nil {
		panic("failed to connect database")
	}

	return db

}
