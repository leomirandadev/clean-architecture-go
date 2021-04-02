package middlewares

import (
	"github.com/leomirandadev/clean-architecture-go/handlers/middlewares/auth"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
	"github.com/leomirandadev/clean-architecture-go/utils/token"
)

type Middleware struct {
	Auth auth.AuthMiddleware
}

func New(tokenHasher token.TokenHash, log logger.Logger) *Middleware {
	return &Middleware{
		Auth: auth.NewBearer(tokenHasher, log),
	}
}
