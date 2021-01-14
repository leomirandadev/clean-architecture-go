package main

import (
	"github.com/leomirandadev/clean-architecture-go/handlers"
	"github.com/leomirandadev/clean-architecture-go/utils/httpRouter"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
	"github.com/leomirandadev/clean-architecture-go/utils/token"
)

func main() {

	router := httpRouter.NewMuxRouter()
	log := logger.NewLogrusLog()
	tokenHasher := token.NewJWT()

	handlers.ConfigHandlers(router, log, tokenHasher)

	router.SERVE(":80")
}
