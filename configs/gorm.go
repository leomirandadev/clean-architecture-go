package configs

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func GetReaderGorm() *gorm.DB {
	var DB_CONNECTION string = os.Getenv("DB_CONNECTION")
	reader, err := gorm.Open("mysql", DB_CONNECTION)

	if err != nil {
		panic("failed to connect database")
	}

	return reader
}

func GetWriterGorm() *gorm.DB {
	var DB_CONNECTION string = os.Getenv("DB_CONNECTION")
	writer, err := gorm.Open("mysql", DB_CONNECTION)

	if err != nil {
		panic("failed to connect database")
	}

	return writer
}
