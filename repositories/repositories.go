package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	"github.com/leomirandadev/clean-architecture-go/repositories/user"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
)

// Container modelo para exportação dos repositórios instanciados
type Container struct {
	User user.UserRepository
}

// Options struct de opções para a criação de uma instancia dos serviços
type Options struct {
	Log        logger.Logger
	WriterSqlx *sqlx.DB
	ReaderSqlx *sqlx.DB
	WriterGorm *gorm.DB
	ReaderGorm *gorm.DB
}

// New cria uma nova instancia dos repositórios
func New(opts Options) *Container {
	return &Container{
		User: user.NewSqlxRepository(opts.Log, opts.WriterSqlx, opts.ReaderSqlx),
	}
}
