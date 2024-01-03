package middlewares

import (
	"github.com/leomirandadev/clean-architecture-go/internal/handlers/middlewares/auth"
	"github.com/leomirandadev/clean-architecture-go/pkg/token"
)

type Middleware struct {
	Auth auth.AuthMiddleware
}

func New(tokenHasher token.TokenHash, basicAuthUser, basicAuthPassword string) *Middleware {
	return &Middleware{
		Auth: auth.NewBearer(tokenHasher, basicAuthUser, basicAuthPassword),
	}
}
