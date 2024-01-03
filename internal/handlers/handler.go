package handlers

import (
	"github.com/leomirandadev/clean-architecture-go/internal/handlers/health"
	"github.com/leomirandadev/clean-architecture-go/internal/handlers/middlewares"
	"github.com/leomirandadev/clean-architecture-go/internal/handlers/swagger"
	"github.com/leomirandadev/clean-architecture-go/internal/handlers/users"
	"github.com/leomirandadev/clean-architecture-go/internal/services"
	"github.com/leomirandadev/clean-architecture-go/pkg/httprouter"
	"github.com/leomirandadev/clean-architecture-go/pkg/token"
)

type Options struct {
	Token             token.TokenHash
	Router            httprouter.Router
	Srv               *services.Container
	BasicAuthUser     string
	BasicAuthPassword string
	AppDeepLinkBase   string
}

func New(opts Options) {
	mid := middlewares.New(opts.Token, opts.BasicAuthUser, opts.BasicAuthPassword)

	users.Init(mid, opts.Router, opts.Token, opts.Srv, opts.AppDeepLinkBase)
	health.Init(mid, opts.Router)
	swagger.Init(mid, opts.Router)
}
