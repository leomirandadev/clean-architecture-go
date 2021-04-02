package main

import (
	"os"
	"strconv"
	"time"

	"github.com/leomirandadev/clean-architecture-go/controllers"
	"github.com/leomirandadev/clean-architecture-go/handlers"
	"github.com/leomirandadev/clean-architecture-go/repositories"
	"github.com/leomirandadev/clean-architecture-go/services"
	"github.com/leomirandadev/clean-architecture-go/utils/cache"
	"github.com/leomirandadev/clean-architecture-go/utils/httpRouter"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
	"github.com/leomirandadev/clean-architecture-go/utils/token"
)

func main() {

	router := httpRouter.NewMuxRouter()
	log := logger.NewLogrusLog()
	tokenHasher := token.NewJWT()

	cacheExpiration, _ := strconv.ParseInt(os.Getenv("CACHE_EXP"), 10, 64)

	cacheStore := cache.NewMemcache(cache.Options{
		URL:        os.Getenv("CACHE_URL"),
		Expiration: time.Duration(cacheExpiration),
	}, log)

	repo := repositories.New(repositories.Options{Log: log})

	srv := services.New(services.Options{
		Log:   log,
		Repo:  repo,
		Cache: cacheStore,
	})

	ctrl := controllers.New(controllers.Options{
		Log:   log,
		Srv:   srv,
		Token: tokenHasher,
	})

	handlers.New(handlers.Options{
		Ctrl:   ctrl,
		Router: router,
		Log:    log,
		Token:  tokenHasher,
	})

	router.SERVE(":8080")
}
