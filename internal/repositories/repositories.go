package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/leomirandadev/clean-architecture-go/internal/repositories/reset_passwords"
	"github.com/leomirandadev/clean-architecture-go/internal/repositories/users"
)

type DB struct {
	User          users.IRepository
	ResetPassword reset_passwords.IRepository
}

type Options struct {
	WriterSqlx *sqlx.DB
	ReaderSqlx *sqlx.DB
}

func New(opts Options) *DB {
	return &DB{
		User:          users.NewSqlx(opts.WriterSqlx, opts.ReaderSqlx),
		ResetPassword: reset_passwords.NewSqlx(opts.WriterSqlx, opts.ReaderSqlx),
	}
}
