package handlers

import (
	"github.com/leomirandadev/clean-architecture-go/handlers/user"
	"github.com/leomirandadev/clean-architecture-go/utils/httpRouter"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
)

func ConfigHandlers(router httpRouter.Router, log logger.Logger) {
	user.ConfigHandlers(router, log)
}
