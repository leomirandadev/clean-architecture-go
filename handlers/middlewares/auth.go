package middlewares

import (
	"net/http"
)

type AuthMiddleware interface {
	Public(next http.HandlerFunc) http.HandlerFunc
	Admin(next http.HandlerFunc) http.HandlerFunc
	VerifyRoles(r *http.Request, logged bool, roles ...string) error
}
