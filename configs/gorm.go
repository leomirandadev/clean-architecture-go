package configs

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func ConnectGorm() *gorm.DB {

	var host string = os.Getenv("DB_HOST_MYSQL")
	var user string = os.Getenv("DB_USER_MYSQL")
	var password string = os.Getenv("DB_PASSWORD_MYSQL")
	var dbname string = os.Getenv("DB_NAME_MYSQL")

	db, err := gorm.Open("mysql", user+":"+password+"@("+host+")/"+dbname+"?charset=utf8\u0026readTimeout=30s\u0026writeTimeout=30s&parseTime=true&loc=Local")

	if err != nil {
		panic("failed to connect database")
	}

	return db

}
