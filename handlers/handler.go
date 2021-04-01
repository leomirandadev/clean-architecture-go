package handlers

import (
	"github.com/leomirandadev/clean-architecture-go/controllers"
	"github.com/leomirandadev/clean-architecture-go/handlers/middlewares"
	"github.com/leomirandadev/clean-architecture-go/handlers/users"
	"github.com/leomirandadev/clean-architecture-go/utils/httpRouter"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
	"github.com/leomirandadev/clean-architecture-go/utils/token"
)

type Options struct {
	Ctrl   *controllers.Container
	Router httpRouter.Router
	Log    logger.Logger
	Token  token.TokenHash
}

func New(opts Options) {

	mid := middlewares.New(opts.Token)
	users.New(mid, opts.Router, opts.Ctrl)

}
