package services

import (
	"github.com/leomirandadev/clean-architecture-go/repositories"
	"github.com/leomirandadev/clean-architecture-go/services/user"
	"github.com/leomirandadev/clean-architecture-go/utils/cache"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
)

// Container modelo para exportação dos serviços instanciados
type Container struct {
	User user.UserService
}

// Options struct de opções para a criação de uma instancia dos serviços
type Options struct {
	Repo  *repositories.Container
	Log   logger.Logger
	Cache cache.Cache
}

// New cria uma nova instancia dos serviços
func New(opts Options) *Container {

	return &Container{
		User: user.New(opts.Repo, opts.Log),
	}
}
