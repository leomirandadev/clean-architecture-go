package repositories

import (
	"github.com/leomirandadev/clean-architecture-go/repositories/user"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
)

// Container modelo para exportação dos repositórios instanciados
type Container struct {
	User user.UserRepository
}

// Options struct de opções para a criação de uma instancia dos serviços
type Options struct {
	Log logger.Logger
}

// New cria uma nova instancia dos repositórios
func New(opts Options) *Container {
	return &Container{
		User: user.NewGormRepository(opts.Log),
	}
}
