package configs

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func GetReaderSqlx() *sqlx.DB {
	var DB_CONNECTION string = os.Getenv("DB_CONNECTION")
	reader := sqlx.MustConnect("mysql", DB_CONNECTION)

	return reader
}

func GetWriterSqlx() *sqlx.DB {
	var DB_CONNECTION string = os.Getenv("DB_CONNECTION")
	writer := sqlx.MustConnect("mysql", DB_CONNECTION)

	return writer
}
