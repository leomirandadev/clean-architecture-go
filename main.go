package main

import (
	"os"
	"strconv"
	"time"

	"github.com/leomirandadev/clean-architecture-go/configs"
	"github.com/leomirandadev/clean-architecture-go/controllers"
	"github.com/leomirandadev/clean-architecture-go/handlers"
	"github.com/leomirandadev/clean-architecture-go/repositories"
	"github.com/leomirandadev/clean-architecture-go/services"
	"github.com/leomirandadev/clean-architecture-go/utils/cache"
	"github.com/leomirandadev/clean-architecture-go/utils/httpRouter"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
	"github.com/leomirandadev/clean-architecture-go/utils/mail"
	"github.com/leomirandadev/clean-architecture-go/utils/token"
)

func main() {

	router, log, tokenHasher, cacheStore, mailSMTP := toolsInit()

	repo := repositories.New(repositories.Options{
		Log:        log,
		ReaderSqlx: configs.GetReaderSqlx(),
		WriterSqlx: configs.GetWriterSqlx(),
		ReaderGorm: configs.GetReaderGorm(),
		WriterGorm: configs.GetWriterGorm(),
	})

	srv := services.New(services.Options{
		Log:   log,
		Repo:  repo,
		Cache: cacheStore,
		Mail:  mailSMTP,
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

func toolsInit() (httpRouter.Router, logger.Logger, token.TokenHash, cache.Cache, mail.MailSender) {

	router := httpRouter.NewMuxRouter()
	log := logger.NewLogrusLog()
	tokenHasher := token.NewJWT()

	mailSMTP := mail.NewSMTP(mail.SmtpConfigs{
		UserSender:     "",
		PasswordSender: "",
		HostSender:     "",
		NameSender:     "",
		EmailSender:    "",
		Port:           "",
	})

	cacheStore := cache.NewMemcache(cache.Options{
		URL: os.Getenv("CACHE_URL"),
		Expiration: func() time.Duration {
			cacheExpiration, _ := strconv.ParseInt(os.Getenv("CACHE_EXP"), 10, 64)
			return time.Duration(cacheExpiration)
		}(),
	}, log)

	return router, log, tokenHasher, cacheStore, mailSMTP
}
