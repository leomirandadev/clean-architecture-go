package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/leomirandadev/clean-architecture-go/internal/permissions"
	"github.com/leomirandadev/clean-architecture-go/pkg/customerr"
	"github.com/leomirandadev/clean-architecture-go/pkg/httprouter"
	"github.com/leomirandadev/clean-architecture-go/pkg/token"
)

type middlewareJWT struct {
	token             token.TokenHash
	basicAuthUser     string
	basicAuthPassword string
}

func NewBearer(tokenHasher token.TokenHash, basicAuthUser, basicAuthPassword string) AuthMiddleware {
	return &middlewareJWT{
		token:             tokenHasher,
		basicAuthUser:     basicAuthUser,
		basicAuthPassword: basicAuthPassword,
	}
}

func (m *middlewareJWT) Public(next httprouter.HandlerFunc) httprouter.HandlerFunc {
	return func(httpCtx httprouter.Context) error {
		if err := m.verifyRoles(httpCtx.Context(), httpCtx.Headers(), false); err != nil {
			return customerr.WithStatus(http.StatusUnauthorized, "unauthorizated", err.Error())
		}

		return next(httpCtx)
	}

}

func (m *middlewareJWT) Private(next httprouter.HandlerFunc) httprouter.HandlerFunc {
	return func(httpCtx httprouter.Context) error {
		if err := m.verifyRoles(httpCtx.Context(), httpCtx.Headers(), true, permissions.Strings()...); err != nil {
			return customerr.WithStatus(http.StatusUnauthorized, "unauthorizated", err.Error())
		}

		return next(httpCtx)
	}
}

func (m *middlewareJWT) Admin(next httprouter.HandlerFunc) httprouter.HandlerFunc {
	return func(httpCtx httprouter.Context) error {
		if err := m.verifyRoles(httpCtx.Context(), httpCtx.Headers(), true, permissions.Admin.String()); err != nil {
			return customerr.WithStatus(http.StatusUnauthorized, "unauthorizated", err.Error())
		}

		return next(httpCtx)
	}
}

func (m *middlewareJWT) Premium(next httprouter.HandlerFunc) httprouter.HandlerFunc {
	return func(httpCtx httprouter.Context) error {
		if err := m.verifyRoles(httpCtx.Context(), httpCtx.Headers(), true, permissions.Premium.String()); err != nil {
			return customerr.WithStatus(http.StatusUnauthorized, "unauthorizated", err.Error())
		}

		return next(httpCtx)
	}
}

func (m *middlewareJWT) BasicAuth(next httprouter.HandlerFunc) httprouter.HandlerFunc {
	return func(httpCtx httprouter.Context) error {
		username, password, ok := httpCtx.GetRequestReader().BasicAuth()
		if ok {
			if username == m.basicAuthUser && password == m.basicAuthPassword {
				return next(httpCtx)
			}
		}

		httpCtx.GetResponseWriter().Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		return customerr.WithStatus(http.StatusUnauthorized, "unauthorizated", "")
	}
}

func (m *middlewareJWT) verifyRoles(ctx context.Context, header http.Header, logged bool, roles ...string) error {
	if !logged {
		return nil
	}

	if header["Authorization"] == nil {
		slog.InfoContext(ctx, "authorization nil")
		return errors.New("you must to provide a valid token")
	}

	bearerSplited := strings.Split(header["Authorization"][0], " ")
	if len(bearerSplited) != 2 {
		slog.InfoContext(ctx, "can't split bearer")
		return errors.New("error to get the token. Certify that you are providing the bearer prefix")
	}

	isValid, claims, err := m.token.Decrypt(bearerSplited[1])
	if err != nil {
		slog.WarnContext(ctx, "decrypt", "err", err)
		return err
	}

	if !isValid {
		slog.InfoContext(ctx, "TOKEN NOT VALID")
		return errors.New("this token is not valid anymore")
	}

	for _, role := range roles {
		if claims["role"] == role {
			m.InsertTokenFieldsOnPayload(claims, header)
			return nil
		}
	}

	slog.InfoContext(ctx, "role not found")
	return errors.New("your user needs more permissions to use this endpoint")
}

func (m *middlewareJWT) InsertTokenFieldsOnPayload(token map[string]any, header http.Header) {
	header.Add("payload.id", token["id"].(string))
	header.Add("payload.role", token["role"].(string))
}
