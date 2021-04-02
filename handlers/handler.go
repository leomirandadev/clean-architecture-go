package handlers

import (
	"github.com/leomirandadev/clean-architecture-go/controllers"
	"github.com/leomirandadev/clean-architecture-go/handlers/middlewares"
	"github.com/leomirandadev/clean-architecture-go/handlers/users"
	"github.com/leomirandadev/clean-architecture-go/utils/httpRouter"
)

type Options struct {
	Ctrl   *controllers.Container
	Mid    *middlewares.Middleware
	Router httpRouter.Router
}

func New(opts Options) {

	users.New(opts.Mid, opts.Router, opts.Ctrl)

}
