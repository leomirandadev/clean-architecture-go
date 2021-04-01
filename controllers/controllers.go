package controllers

import (
	"github.com/leomirandadev/clean-architecture-go/controllers/user"
	"github.com/leomirandadev/clean-architecture-go/services"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
	"github.com/leomirandadev/clean-architecture-go/utils/token"
)

// Container modelo para exportação dos serviços instanciados
type Container struct {
	User user.UserController
}

// Options struct de opções para a criação de uma instancia dos serviços
type Options struct {
	Srv   *services.Container
	Log   logger.Logger
	Token token.TokenHash
}

// New cria uma nova instancia dos serviços
func New(opts Options) *Container {

	return &Container{
		User: user.New(opts.Srv, opts.Log, opts.Token),
	}
}
