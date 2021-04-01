package users

import (
	"github.com/leomirandadev/clean-architecture-go/controllers"
	"github.com/leomirandadev/clean-architecture-go/handlers/middlewares"
	"github.com/leomirandadev/clean-architecture-go/utils/httpRouter"
)

func New(mid *middlewares.Middleware, router httpRouter.Router, Ctrl *controllers.Container) {

	router.POST("/users", mid.Auth.Admin(Ctrl.User.Create))
	router.GET("/users/{id}", mid.Auth.Admin(Ctrl.User.GetByID))
	router.POST("/users/auth", mid.Auth.Public(Ctrl.User.Auth))

}
