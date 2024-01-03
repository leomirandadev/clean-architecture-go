package users

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/leomirandadev/clean-architecture-go/internal/models"
	"github.com/leomirandadev/clean-architecture-go/internal/permissions"
	"github.com/leomirandadev/clean-architecture-go/internal/repositories"
	"github.com/leomirandadev/clean-architecture-go/pkg/customerr"
	"github.com/leomirandadev/clean-architecture-go/pkg/hasher"
	"github.com/leomirandadev/clean-architecture-go/pkg/mail"
	"github.com/leomirandadev/clean-architecture-go/pkg/sso/google"
	"github.com/leomirandadev/clean-architecture-go/pkg/tracer"
)

type IService interface {
	Create(ctx context.Context, newUser models.UserRequest) (*models.UserResponse, error)
	GetByID(ctx context.Context, ID string) (*models.UserResponse, error)
	GetUserByLogin(ctx context.Context, userLogin models.UserAuth) (*models.UserResponse, error)
	ResetPassword(ctx context.Context, req models.ResetPasswordRequest) error
	ChangePassword(ctx context.Context, req models.ChangePassword) error
	GoogleSignin(ctx context.Context) string
	GoogleCallback(ctx context.Context, req models.GoogleCallbackReq) (*models.UserResponse, error)
}

type serviceImpl struct {
	repos     *repositories.DB
	hashDoer  hasher.Hasher
	mailing   mail.MailSender
	googleSSO google.GoogleSSO
}

func New(repos *repositories.DB, mailing mail.MailSender, googleSSO google.GoogleSSO) IService {
	hashDoer := hasher.NewBcryptHasher()
	return &serviceImpl{repos, hashDoer, mailing, googleSSO}
}

func (srv serviceImpl) Create(ctx context.Context, req models.UserRequest) (*models.UserResponse, error) {
	ctx, tr := tracer.Span(ctx, "srv.users.create")
	defer tr.End()

	userFound, _ := srv.repos.User.GetUserByEmail(ctx, req.Email)
	if userFound != nil {
		slog.WarnContext(ctx, "this email is already in use")
		return nil, customerr.WithStatus(http.StatusBadRequest, "this email is already in use", req)
	}

	salt, err := srv.hashDoer.Salt()
	if err != nil {
		slog.WarnContext(ctx, "hasher.salt", "err", err)
		return nil, customerr.WithStatus(http.StatusInternalServerError, "encrypt error", err)
	}

	passwordHashed, err := srv.hashDoer.Generate(salt, req.Password)
	if err != nil {
		slog.WarnContext(ctx, "hasher.generate", "err", err)
		return nil, customerr.WithStatus(http.StatusInternalServerError, "encrypt error", err)
	}

	id, err := srv.repos.User.Create(ctx, models.User{
		NickName: req.NickName,
		Name:     req.Name,
		Email:    req.Email,
		Password: passwordHashed,
		Salt:     salt,
		Role:     permissions.Free.String(),
	})
	if err != nil {
		slog.WarnContext(ctx, "repositories.user.create", "err", err)
		return nil, customerr.WithStatus(http.StatusInternalServerError, "create user fails", err)
	}

	result, err := srv.repos.User.GetByID(ctx, id)
	if err != nil {
		slog.WarnContext(ctx, "repositories.user.get_by_id", "err", err)
		return nil, customerr.WithStatus(http.StatusInternalServerError, "error to recover user", err)
	}

	resp := result.ToResponse()

	return &resp, nil
}

func (srv serviceImpl) GetUserByLogin(ctx context.Context, req models.UserAuth) (*models.UserResponse, error) {
	ctx, tr := tracer.Span(ctx, "srv.users.get_user_by_login")
	defer tr.End()

	userFound, err := srv.repos.User.GetUserByEmail(ctx, req.Email)
	if err != nil {
		slog.WarnContext(ctx, "repositories.user.get_user_by_email", "err", err)
		return nil, customerr.WithStatus(http.StatusBadRequest, "user not found", err)
	}

	err = srv.hashDoer.Compare(userFound.Password, userFound.Salt, req.Password)
	if err != nil {
		slog.WarnContext(ctx, "hasher.compare", "err", err)
		return nil, customerr.WithStatus(http.StatusBadRequest, "email or password wrong", req)
	}

	resp := userFound.ToResponse()

	return &resp, nil
}

func (srv serviceImpl) GetByID(ctx context.Context, ID string) (*models.UserResponse, error) {
	ctx, tr := tracer.Span(ctx, "srv.users.get_by_id")
	defer tr.End()

	result, err := srv.repos.User.GetByID(ctx, ID)
	if err != nil {
		slog.WarnContext(ctx, "repositories.user.get_by_id", "err", err)
		return nil, customerr.WithStatus(http.StatusNotFound, "user not found", err)
	}

	resp := result.ToResponse()

	return &resp, nil
}

func (srv serviceImpl) ResetPassword(ctx context.Context, req models.ResetPasswordRequest) error {
	ctx, tr := tracer.Span(ctx, "srv.users.reset_password")
	defer tr.End()

	userFound, err := srv.repos.User.GetUserByEmail(ctx, req.Email)
	if err != nil {
		slog.WarnContext(ctx, "repositories.user.get_user_by_email", "err", err)
		return customerr.WithStatus(http.StatusBadRequest, "user not found", err)
	}

	createResetRow := models.CreateResetPassword{
		UserID: userFound.ID,
		Token:  randToken(TOKEN_RESET_PASSWORD_SIZE),
	}

	if _, err := srv.repos.ResetPassword.Create(ctx, createResetRow); err != nil {
		slog.WarnContext(ctx, "save token fails", "err", err)
		return customerr.WithStatus(http.StatusInternalServerError, "save token fails", err)
	}

	go srv.sendTokenByEmail(userFound.Email, createResetRow.Token)

	slog.Info("token generated", "token", createResetRow.Token)

	return nil
}

func (srv serviceImpl) sendTokenByEmail(email, token string) {
	_, tr := tracer.Span(context.Background(), "srv.users.send_token_by_email")
	defer tr.End()

	msg := fmt.Sprintf("Your token is: %s", token)
	emails := []string{email}

	if err := srv.mailing.Send(emails, "Reset Password", msg); err != nil {
		slog.Error("email was not sent", "err", err)
	}
}

func (srv serviceImpl) ChangePassword(ctx context.Context, req models.ChangePassword) error {
	ctx, tr := tracer.Span(ctx, "srv.users.change_password")
	defer tr.End()

	resetPasswordResp, err := srv.repos.ResetPassword.GetByToken(ctx, req.Token)
	if err != nil {
		slog.WarnContext(ctx, "invalid token", "err", err)
		return customerr.WithStatus(http.StatusBadRequest, "invalid token", err)
	}

	userFound, err := srv.repos.User.GetByID(ctx, resetPasswordResp.UserID)
	if err != nil {
		slog.WarnContext(ctx, "user not found", "err", err)
		return customerr.WithStatus(http.StatusBadRequest, "user not found", err)
	}

	salt, err := srv.hashDoer.Salt()
	if err != nil {
		slog.WarnContext(ctx, "hasher.salt", "err", err)
		return customerr.WithStatus(http.StatusInternalServerError, "encrypt error", err)
	}

	passwordHashed, err := srv.hashDoer.Generate(salt, req.NewPassword)
	if err != nil {
		slog.WarnContext(ctx, "hasher.generate", "err", err)
		return customerr.WithStatus(http.StatusInternalServerError, "encrypt error", err)
	}

	if err := srv.repos.User.ChangePassword(ctx, userFound.ID, passwordHashed, salt); err != nil {
		slog.ErrorContext(ctx, "update password fails", "err", err)
		return customerr.WithStatus(http.StatusInternalServerError, "update password fails", err)
	}

	return nil
}

func (srv serviceImpl) GoogleSignin(ctx context.Context) string {
	_, tr := tracer.Span(context.Background(), "srv.users.google_sign_in")
	defer tr.End()

	return srv.googleSSO.GenerateURL()
}

func (srv serviceImpl) GoogleCallback(ctx context.Context, req models.GoogleCallbackReq) (*models.UserResponse, error) {
	ctx, tr := tracer.Span(ctx, "srv.users.google_callback")
	defer tr.End()

	if !srv.googleSSO.IsRandomStrValid(req.State) {
		slog.Debug("state not valid")
		return nil, customerr.WithStatus(http.StatusBadRequest, "state not valid", nil)
	}

	data, err := srv.googleSSO.GetUserData(ctx, req.State, req.Code)
	if err != nil {
		slog.Debug("get user data fails", "err", err)
		return nil, customerr.WithStatus(http.StatusInternalServerError, "get google data fails", err)
	}

	userFound, err := srv.repos.User.GetUserByEmail(ctx, data.Email)
	if err == nil {
		resp := userFound.ToResponse()
		return &resp, nil
	}

	id, err := srv.repos.User.Create(ctx, models.User{
		NickName: data.Email,
		Name:     data.Name,
		Email:    data.Email,
		Password: "",
		Salt:     "",
		Role:     permissions.Free.String(),
	})
	if err != nil {
		slog.WarnContext(ctx, "repositories.user.create", "err", err)
		return nil, customerr.WithStatus(http.StatusInternalServerError, "create user fails", err)
	}

	result, err := srv.repos.User.GetByID(ctx, id)
	if err != nil {
		slog.WarnContext(ctx, "repositories.user.get_by_id", "err", err)
		return nil, customerr.WithStatus(http.StatusInternalServerError, "error to recover user", err)
	}

	resp := result.ToResponse()

	return &resp, nil
}

func randToken(size int) string {
	var str strings.Builder
	for i := 0; i < size; i++ {
		str.WriteString(strconv.Itoa(rand.Intn(9)))
	}

	return str.String()
}
