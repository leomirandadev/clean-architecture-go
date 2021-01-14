package users

import (
	"github.com/leomirandadev/clean-architecture-go/controllers"
	"github.com/leomirandadev/clean-architecture-go/handlers/middlewares"
	repositories "github.com/leomirandadev/clean-architecture-go/repositories/user"
	"github.com/leomirandadev/clean-architecture-go/services"
	"github.com/leomirandadev/clean-architecture-go/utils/httpRouter"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
	"github.com/leomirandadev/clean-architecture-go/utils/token"
)

func ConfigHandlers(router httpRouter.Router, log logger.Logger, mid *middlewares.Middleware, tokenHasher token.TokenHash) {

	var (
		userRepo       repositories.UserRepository = repositories.NewSqlxRepository(log)
		userService    services.UserService        = services.NewUserService(userRepo, log)
		userController controllers.UserController  = controllers.NewUserController(userService, log, tokenHasher)
	)

	userRepo.Migrate()
	router.POST("/users", mid.Auth.Admin(userController.Create))
	router.GET("/users/{id}", mid.Auth.Admin(userController.GetByID))
	router.POST("/users/auth", mid.Auth.Public(userController.Auth))

}
