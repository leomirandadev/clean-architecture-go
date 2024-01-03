package main

import (
	"log"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/leomirandadev/clean-architecture-go/internal/handlers"
	"github.com/leomirandadev/clean-architecture-go/internal/repositories"
	"github.com/leomirandadev/clean-architecture-go/internal/services"
	"github.com/leomirandadev/clean-architecture-go/pkg/envs"
	"github.com/leomirandadev/clean-architecture-go/pkg/httprouter"
	"github.com/leomirandadev/clean-architecture-go/pkg/mail"
	"github.com/leomirandadev/clean-architecture-go/pkg/sso/google"
	"github.com/leomirandadev/clean-architecture-go/pkg/token"
	"github.com/leomirandadev/clean-architecture-go/pkg/tracer"
	"github.com/leomirandadev/clean-architecture-go/pkg/tracer/otel_jaeger"
)

type Config struct {
	Port            string              `mapstructure:"PORT"`
	Env             string              `mapstructure:"ENV"`
	Tracer          otel_jaeger.Options `mapstructure:"TRACER"`
	JWT             token.Options       `mapstructure:"JWT"`
	Mailing         mail.Options        `mapstructure:"MAILING"`
	AppDeepLinkBase string              `mapstructure:"APP_DEEP_LINK_BASE"`
	GoogleSSO       google.Options      `mapstructure:"GOOGLE_SSO"`

	Database struct {
		Reader string `mapstructure:"READER"`
		Writer string `mapstructure:"WRITER"`
	} `mapstructure:"DATABASE"`

	BasicAuth struct {
		User     string `mapstructure:"USER"`
		Password string `mapstructure:"PASSWORD"`
	} `mapstructure:"BASIC_AUTH"`
}

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" before paste the token
func main() {
	var cfg Config
	if err := envs.Load(".", &cfg); err != nil {
		log.Fatal("cfg not loaded", err)
	}

	tools := toolsInit(cfg)
	defer tools.tr.Close()

	repo := repositories.New(repositories.Options{
		ReaderSqlx: sqlx.MustConnect("pgx", cfg.Database.Reader),
		WriterSqlx: sqlx.MustConnect("pgx", cfg.Database.Writer),
	})

	srv := services.New(services.Options{
		Repo:      repo,
		Mailing:   tools.mailing,
		GoogleSSO: tools.googleSSO,
	})

	handlers.New(handlers.Options{
		Srv:               srv,
		Router:            tools.router,
		Token:             tools.tokenizer,
		BasicAuthUser:     cfg.BasicAuth.User,
		BasicAuthPassword: cfg.BasicAuth.Password,
		AppDeepLinkBase:   cfg.AppDeepLinkBase,
	})

	tools.router.SERVE(cfg.Port)
}

type tools struct {
	router    httprouter.Router
	tokenizer token.TokenHash
	tr        tracer.ITracer
	mailing   mail.MailSender
	googleSSO google.GoogleSSO
}

func toolsInit(cfg Config) tools {

	slog.SetDefault(slog.New(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}),
	))

	jaegerCollector := otel_jaeger.NewCollector(cfg.Tracer)

	return tools{
		router:    httprouter.NewChiRouter(),
		tokenizer: token.NewJWT(cfg.JWT),
		tr:        tracer.New(jaegerCollector),
		mailing:   mail.NewSMTP(cfg.Mailing),
		googleSSO: google.New(cfg.GoogleSSO),
	}
}
