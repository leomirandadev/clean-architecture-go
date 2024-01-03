package users

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/leomirandadev/clean-architecture-go/internal/handlers/middlewares"
	"github.com/leomirandadev/clean-architecture-go/internal/models"
	"github.com/leomirandadev/clean-architecture-go/internal/services"
	"github.com/leomirandadev/clean-architecture-go/pkg/customerr"
	"github.com/leomirandadev/clean-architecture-go/pkg/httprouter"
	"github.com/leomirandadev/clean-architecture-go/pkg/token"
	"github.com/leomirandadev/clean-architecture-go/pkg/validator"
)

func Init(
	mid *middlewares.Middleware,
	router httprouter.Router,
	token token.TokenHash,
	srv *services.Container,
	appDeepLinkBase string,
) {
	ctr := controllers{
		srv,
		token,
		appDeepLinkBase,
	}

	router.POST("/v1/users", mid.Auth.Public(ctr.Create))
	router.GET("/v1/users/{id}", mid.Auth.Admin(ctr.GetByID))
	router.GET("/v1/users/me", mid.Auth.Private(ctr.GetMe))
	router.POST("/v1/users/auth", mid.Auth.Public(ctr.Auth))
	router.POST("/v1/users/password/reset", mid.Auth.Public(ctr.ResetPassword))
	router.PUT("/v1/users/password", mid.Auth.Public(ctr.ChangePassword))
	router.GET("/v1/users/auth/google/sigin", mid.Auth.Public(ctr.GoogleSignin))
	router.GET("/v1/users/auth/google/callback", mid.Auth.Public(ctr.GoogleCallback))
}

type controllers struct {
	srv             *services.Container
	token           token.TokenHash
	appDeepLinkBase string
}

// user swagger document
// @Description Create one user
// @Tags users
// @Param user body models.UserRequest true "create new user"
// @Accept json
// @Produce json
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} customerr.Error
// @Failure 401 {object} customerr.Error
// @Failure 500 {object} customerr.Error
// @Router /v1/users [post]
func (ctr controllers) Create(c httprouter.Context) error {
	ctx := c.Context()

	var newUser models.UserRequest
	if err := c.Decode(&newUser); err != nil {
		return customerr.WithStatus(http.StatusBadRequest, "decode error", err.Error())
	}

	if err := validator.Validate(newUser); err != nil {
		slog.WarnContext(ctx, "validator.validate", "err", err)
		return customerr.WithStatus(http.StatusBadRequest, "invalid payload", err.Error())
	}

	user, err := ctr.srv.User.Create(ctx, newUser)
	if err != nil {
		slog.WarnContext(ctx, "srv.user.create", "err", err)
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

// user swagger document
// @Description Authenticate user
// @Tags users
// @Param user body models.UserAuth true "add user"
// @Accept json
// @Produce json
// @Success 200 {object} models.AuthToken
// @Failure 400 {object} customerr.Error
// @Failure 401 {object} customerr.Error
// @Failure 500 {object} customerr.Error
// @Router /v1/users/auth [post]
func (ctr controllers) Auth(c httprouter.Context) error {
	ctx := c.Context()

	var userLogin models.UserAuth
	if err := c.Decode(&userLogin); err != nil {
		return customerr.WithStatus(http.StatusBadRequest, "decode error", err.Error())
	}

	userFound, err := ctr.srv.User.GetUserByLogin(ctx, userLogin)
	if err != nil {
		slog.WarnContext(ctx, "srv.user.get_user_by_login", "err", err)
		return err
	}

	token, err := ctr.generateToken(userFound)
	if err != nil {
		slog.WarnContext(ctx, "token.encrypt", "err", err)
		return customerr.WithStatus(http.StatusInternalServerError, "create token fails", err)
	}

	return c.JSON(http.StatusOK, models.AuthToken{Token: token})
}

// user swagger document
// @Description This endpoint generate and send the token to user to allow reset password
// @Tags users
// @Param user body models.ResetPasswordRequest true "user details"
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} customerr.Error
// @Failure 500 {object} customerr.Error
// @Router /v1/users/password/reset [post]
func (ctr controllers) ResetPassword(c httprouter.Context) error {
	ctx := c.Context()

	var req models.ResetPasswordRequest
	if err := c.Decode(&req); err != nil {
		return customerr.WithStatus(http.StatusBadRequest, "decode error", err.Error())
	}

	if err := validator.Validate(req); err != nil {
		slog.WarnContext(ctx, "validator.validate", "err", err)
		return customerr.WithStatus(http.StatusBadRequest, "invalid payload", err.Error())
	}

	if err := ctr.srv.User.ResetPassword(ctx, req); err != nil {
		slog.WarnContext(ctx, "srv.user.reset_password", "err", err)
		return err
	}

	return c.JSON(http.StatusNoContent, nil)
}

// user swagger document
// @Description Change password passing the reset_password token
// @Tags users
// @Param user body models.ChangePassword true "new password and token"
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} customerr.Error
// @Failure 500 {object} customerr.Error
// @Router /v1/users/password [put]
func (ctr controllers) ChangePassword(c httprouter.Context) error {
	ctx := c.Context()

	var req models.ChangePassword
	if err := c.Decode(&req); err != nil {
		return customerr.WithStatus(http.StatusBadRequest, "decode error", err.Error())
	}

	if err := validator.Validate(req); err != nil {
		slog.WarnContext(ctx, "validator.validate", "err", err)
		return customerr.WithStatus(http.StatusBadRequest, "invalid payload", err.Error())
	}

	if err := ctr.srv.User.ChangePassword(ctx, req); err != nil {
		slog.WarnContext(ctx, "srv.user.reset_password", "err", err)
		return err
	}

	return c.JSON(http.StatusNoContent, nil)
}

// user swagger document
// @Description Get one user
// @Tags users
// @Param id path string true "User ID"
// @Accept json
// @Produce json
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} customerr.Error
// @Failure 401 {object} customerr.Error
// @Failure 500 {object} customerr.Error
// @Security BearerAuth
// @Router /v1/users/{id} [get]
func (ctr controllers) GetByID(c httprouter.Context) error {
	ctx := c.Context()

	userID := c.GetPathParam("id")
	if userID == "" {
		slog.InfoContext(ctx, "you need to provide an id")
		return customerr.WithStatus(http.StatusBadRequest, "you need to provide an id", map[string]string{"user_id": userID})
	}

	user, err := ctr.srv.User.GetByID(ctx, userID)
	if err != nil {
		slog.WarnContext(ctx, "srv.user.get_by_id", "err", err)
		return err
	}

	return c.JSON(http.StatusOK, user)
}

// user swagger document
// @Description Get me
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} customerr.Error
// @Failure 401 {object} customerr.Error
// @Failure 500 {object} customerr.Error
// @Security BearerAuth
// @Router /v1/users/me [get]
func (ctr controllers) GetMe(c httprouter.Context) error {
	ctx := c.Context()

	userID := c.GetFromHeader("payload.id")
	if userID == "" {
		slog.InfoContext(ctx, "user id not found on session")
		return customerr.WithStatus(http.StatusBadRequest, "user id not found on session", map[string]string{"user_id": userID})
	}

	user, err := ctr.srv.User.GetByID(ctx, userID)
	if err != nil {
		slog.WarnContext(ctx, "srv.user.me", "err", err)
		return err
	}

	return c.JSON(http.StatusOK, user)
}

// user swagger document
// @Description Google SSO redirect
// @Tags users
// @Success 307
// @Router /v1/users/auth/google/sigin [get]
func (ctr controllers) GoogleSignin(c httprouter.Context) error {
	var (
		ctx = c.Context()
		url = ctr.srv.User.GoogleSignin(ctx)
		w   = c.GetResponseWriter()
		r   = c.GetRequestReader()
	)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	return nil
}

// user swagger document
// @Description Google SSO callback
// @Tags users
// @Param user body models.GoogleCallbackReq true "google sso callback request"
// @Param state query string true "state"
// @Param code query string true "code"
// @Accept json
// @Produce json
// @Success 200 {object} models.AuthToken
// @Failure 400 {object} customerr.Error
// @Failure 401 {object} customerr.Error
// @Failure 500 {object} customerr.Error
// @Router /v1/users/auth/google/callback [get]
func (ctr controllers) GoogleCallback(c httprouter.Context) error {
	ctx := c.Context()

	req := models.GoogleCallbackReq{
		State: c.GetQueryParam("state"),
		Code:  c.GetQueryParam("code"),
	}

	if err := validator.Validate(req); err != nil {
		slog.WarnContext(ctx, "validator.validate", "err", err)
		return customerr.WithStatus(http.StatusBadRequest, "invalid payload", err.Error())
	}

	userFound, err := ctr.srv.User.GoogleCallback(ctx, req)
	if err != nil {
		slog.Error("error to get user data", "input", req)
		return customerr.WithStatus(http.StatusInternalServerError, "error to get user data", req)
	}

	token, err := ctr.generateToken(userFound)
	if err != nil {
		slog.WarnContext(ctx, "token.encrypt", "err", err)
		return customerr.WithStatus(http.StatusInternalServerError, "create token fails", err)
	}

	var (
		w   = c.GetResponseWriter()
		r   = c.GetRequestReader()
		url = fmt.Sprintf("%s/sso-callback?token=%s", ctr.appDeepLinkBase, token)
	)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	return nil
}

func (ctr controllers) generateToken(user *models.UserResponse) (string, error) {
	inputToEncrypt := map[string]string{"id": user.ID, "role": user.Role}
	token, err := ctr.token.Encrypt(inputToEncrypt)
	if err != nil {
		slog.Warn("token.encrypt", "err", err)
		return "", customerr.WithStatus(http.StatusInternalServerError, err.Error(), inputToEncrypt)
	}

	return token, nil
}
