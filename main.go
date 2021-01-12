package main

import (
	"github.com/leomirandadev/clean-architecture-go/handlers"
	"github.com/leomirandadev/clean-architecture-go/utils/httpRouter"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
)

func main() {

	router := httpRouter.NewMuxRouter()
	log := logger.NewLogrusLog()

	handlers.ConfigHandlers(router, log)

	router.SERVE(":80")
}
