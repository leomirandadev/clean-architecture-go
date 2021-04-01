package user

import (
	"github.com/leomirandadev/clean-architecture-go/entities"
	"github.com/leomirandadev/clean-architecture-go/repositories"
	"github.com/leomirandadev/clean-architecture-go/utils/hasher"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
)

type UserService interface {
	Create(newUser entities.User) error
	GetByID(ID int64) ([]entities.UserResponse, error)
	GetUserByLogin(userLogin entities.UserAuth) (entities.UserResponse, error)
}

type services struct {
	repositories *repositories.Container
	log          logger.Logger
}

func New(repo *repositories.Container, log logger.Logger) UserService {
	return &services{repositories: repo, log: log}
}

func (srv *services) Create(newUser entities.User) error {
	hasherBcrypt := hasher.NewBcryptHasher()
	passwordHashed, errHash := hasherBcrypt.Generate(newUser.Password)

	if errHash != nil {
		srv.log.Error("Srv.Create: ", "Error generate hash password: ", newUser)
		return errHash
	}

	newUser.Password = passwordHashed
	return srv.repositories.User.Create(newUser)
}

func (srv *services) GetUserByLogin(userLogin entities.UserAuth) (entities.UserResponse, error) {

	userFound, err := srv.repositories.User.GetUserByEmail(userLogin)

	if err != nil || len(userFound) <= 0 {
		srv.log.Error("Srv.Auth: ", "User not found", userFound)
		return entities.UserResponse{}, err
	}

	hasherBcrypt := hasher.NewBcryptHasher()
	err = hasherBcrypt.Compare(userFound[0].Password, userLogin.Password)
	if err != nil {
		return entities.UserResponse{}, err
	}

	return entities.UserResponse{
		ID:        userFound[0].ID,
		NickName:  userFound[0].NickName,
		Name:      userFound[0].Name,
		Email:     userFound[0].Email,
		CreatedAt: userFound[0].CreatedAt,
		UpdatedAt: userFound[0].UpdatedAt,
	}, nil
}

func (srv *services) GetByID(ID int64) ([]entities.UserResponse, error) {
	return srv.repositories.User.GetByID(ID)
}
