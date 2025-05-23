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
	"github.com/leomirandadev/clean-architecture-go/pkg/healthcheck"
	"github.com/leomirandadev/clean-architecture-go/pkg/httprouter"
	"github.com/leomirandadev/clean-architecture-go/pkg/mail"
	"github.com/leomirandadev/clean-architecture-go/pkg/slogtint"
	"github.com/leomirandadev/clean-architecture-go/pkg/sso/google"
	"github.com/leomirandadev/clean-architecture-go/pkg/token"
	"github.com/leomirandadev/clean-architecture-go/pkg/tracer"
	"github.com/leomirandadev/clean-architecture-go/pkg/tracer/otel_jaeger"
	"github.com/leomirandadev/clean-architecture-go/pkg/validator"
)

type Config struct {
	Port               string              `mapstructure:"PORT" validate:"required"`
	Env                string              `mapstructure:"ENV" validate:"required"`
	Tracer             otel_jaeger.Options `mapstructure:"TRACER" validate:"required"`
	JWT                token.Options       `mapstructure:"JWT" validate:"required"`
	Mailing            mail.Options        `mapstructure:"MAILING" validate:"required"`
	SSOBaseURLCallback string              `mapstructure:"SSO_BASE_URL_CALLBACK"`
	GoogleSSO          google.Options      `mapstructure:"GOOGLE_SSO" validate:"required"`

	Database struct {
		Reader string `mapstructure:"READER" validate:"required"`
		Writer string `mapstructure:"WRITER" validate:"required"`
	} `mapstructure:"DATABASE" validate:"required"`

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

	if err := validator.Validate(cfg); err != nil {
		log.Fatal("missing required fields", err)
	}

	tools := toolsInit(cfg)
	defer tools.tr.Close()

	writeDBConn := sqlx.MustConnect("pgx", cfg.Database.Reader)
	readDBConn := writeDBConn

	repo := repositories.New(repositories.Options{
		ReaderSqlx: readDBConn,
		WriterSqlx: writeDBConn,
	})

	srv := services.New(services.Options{
		Repo:      repo,
		Mailing:   tools.mailing,
		GoogleSSO: tools.googleSSO,
	})

	tools.healthChecker.
		Register("postgres_writer", writeDBConn.PingContext).
		Register("postgres_reader", readDBConn.PingContext)

	handlers.New(handlers.Options{
		Srv:                srv,
		Router:             tools.router,
		Token:              tools.tokenizer,
		BasicAuthUser:      cfg.BasicAuth.User,
		BasicAuthPassword:  cfg.BasicAuth.Password,
		SSOBaseURLCallback: cfg.SSOBaseURLCallback,
		HealthChecker:      tools.healthChecker,
	})

	tools.router.Serve(cfg.Port)
}

type tools struct {
	router        httprouter.Router
	tokenizer     token.TokenHash
	tr            tracer.ITracer
	mailing       mail.MailSender
	googleSSO     google.GoogleSSO
	healthChecker healthcheck.Health
}

func toolsInit(cfg Config) tools {

	slog.SetDefault(slog.New(
		slogtint.NewHandler(os.Stderr, &slogtint.Options{
			AddSource: true,
			Level:     slog.LevelDebug,
		}),
	))

	jaegerCollector := otel_jaeger.NewCollector(cfg.Tracer)

	return tools{
		router:        httprouter.NewChiRouter(),
		tokenizer:     token.NewJWT(cfg.JWT),
		tr:            tracer.New(jaegerCollector),
		mailing:       mail.NewSMTP(cfg.Mailing),
		googleSSO:     google.New(cfg.GoogleSSO),
		healthChecker: healthcheck.New(),
	}
}
