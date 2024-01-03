package services

import (
	"github.com/leomirandadev/clean-architecture-go/internal/repositories"
	"github.com/leomirandadev/clean-architecture-go/internal/services/users"
	"github.com/leomirandadev/clean-architecture-go/pkg/mail"
	"github.com/leomirandadev/clean-architecture-go/pkg/sso/google"
)

type Container struct {
	User users.IService
}

type Options struct {
	Repo      *repositories.DB
	Mailing   mail.MailSender
	GoogleSSO google.GoogleSSO
}

func New(opts Options) *Container {
	return &Container{
		User: users.New(opts.Repo, opts.Mailing, opts.GoogleSSO),
	}
}
