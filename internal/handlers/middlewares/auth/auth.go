package auth

import (
	"github.com/leomirandadev/clean-architecture-go/pkg/httprouter"
)

type AuthMiddleware interface {
	Public(next httprouter.HandlerFunc) httprouter.HandlerFunc
	Private(next httprouter.HandlerFunc) httprouter.HandlerFunc
	Premium(next httprouter.HandlerFunc) httprouter.HandlerFunc
	Admin(next httprouter.HandlerFunc) httprouter.HandlerFunc
	BasicAuth(next httprouter.HandlerFunc) httprouter.HandlerFunc
}
