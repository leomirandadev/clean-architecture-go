package services

import (
	"github.com/leomirandadev/clean-architecture-go/entities"
	repositories "github.com/leomirandadev/clean-architecture-go/repositories/user"
	"github.com/leomirandadev/clean-architecture-go/utils/hasher"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
)

type UserService interface {
	New(newUser entities.User) error
	GetByID(ID int64) ([]entities.User, error)
}

type services struct {
	repositories repositories.UserRepository
	log          logger.Logger
}

func NewUserService(userRepository repositories.UserRepository, log logger.Logger) UserService {
	return &services{repositories: userRepository, log: log}
}

func (srv *services) New(newUser entities.User) error {
	hasherBcrypt := hasher.NewBcryptHasher()
	passwordHashed, errHash := hasherBcrypt.Generate(newUser.Password)

	if errHash != nil {
		srv.log.Error("Ctrl.Create: ", "Error generate hash password: ", newUser)
		return errHash
	}

	newUser.Password = passwordHashed
	return srv.repositories.Create(newUser)
}

func (srv *services) GetByID(ID int64) ([]entities.User, error) {
	return srv.repositories.GetByID(ID)
}
