package middlewares

import "github.com/leomirandadev/clean-architecture-go/utils/token"

type Middleware struct {
	Auth AuthMiddleware
}

func New(tokenHasher token.TokenHash) *Middleware {
	return &Middleware{
		Auth: newAuthBearer(tokenHasher),
	}
}
