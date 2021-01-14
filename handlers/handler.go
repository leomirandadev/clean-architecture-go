package handlers

import (
	"github.com/leomirandadev/clean-architecture-go/handlers/middlewares"
	"github.com/leomirandadev/clean-architecture-go/handlers/users"
	"github.com/leomirandadev/clean-architecture-go/utils/httpRouter"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
	"github.com/leomirandadev/clean-architecture-go/utils/token"
)

func ConfigHandlers(router httpRouter.Router, log logger.Logger, tokenHasher token.TokenHash) {

	mid := middlewares.New(tokenHasher)
	users.ConfigHandlers(router, log, mid, tokenHasher)
}
