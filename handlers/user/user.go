package user

import (
	"github.com/leomirandadev/clean-architecture-go/controllers"
	repositories "github.com/leomirandadev/clean-architecture-go/repositories/user"
	"github.com/leomirandadev/clean-architecture-go/services"
	"github.com/leomirandadev/clean-architecture-go/utils/httpRouter"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
)

func ConfigHandlers(router httpRouter.Router, log logger.Logger) {

	var (
		userRepo       repositories.UserRepository = repositories.NewSqlxRepository(log)
		userService    services.UserService        = services.NewUserService(userRepo, log)
		userController controllers.UserController  = controllers.NewUserController(userService, log)
	)

	userRepo.Migrate()
	router.POST("/users", userController.Create)
	router.GET("/users/{id}", userController.GetByID)

}
