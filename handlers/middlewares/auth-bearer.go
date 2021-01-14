package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/leomirandadev/clean-architecture-go/utils/token"
)

type middlewareJWT struct {
	token token.TokenHash
}

type Response struct {
	Message string `json:"message"`
}

func newAuthBearer(tokenHasher token.TokenHash) AuthMiddleware {
	return &middlewareJWT{
		token: tokenHasher,
	}
}

func (m *middlewareJWT) Public(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if err := m.VerifyRoles(r, false); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Response{Message: "Permissão negada"})
			return
		}

		next.ServeHTTP(w, r)
	})

}

func (m *middlewareJWT) Admin(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if err := m.VerifyRoles(r, true, "admin"); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Response{Message: "Permissão negada"})
			return
		}

		next.ServeHTTP(w, r)
	})

}

func (m *middlewareJWT) VerifyRoles(r *http.Request, logged bool, roles ...string) error {

	if !logged {
		return nil
	}

	if r.Header["Authorization"] == nil {
		return errors.New("WITHOUT_AUTHORIZATION")
	}

	bearerSplited := strings.Split(r.Header["Authorization"][0], " ")

	if len(bearerSplited) != 2 {
		return errors.New("INVALID_AUTHORIZATION")
	}

	_, err := m.token.Decrypt(bearerSplited[1])

	if err != nil {
		return errors.New("UNAUTHORIZED")
	}

	return nil
}
