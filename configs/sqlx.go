package configs

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func ConnectSqlx() *sqlx.DB {

	var DB_CONNECTION string = os.Getenv("DB_CONNECTION")

	db, err := sqlx.Connect("mysql", DB_CONNECTION)

	if err != nil {
		panic("failed to connect database")
	}

	return db

}
