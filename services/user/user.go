package user

import (
	"context"

	"github.com/leomirandadev/clean-architecture-go/entities"
	"github.com/leomirandadev/clean-architecture-go/repositories"
	"github.com/leomirandadev/clean-architecture-go/utils/hasher"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
)

type UserService interface {
	Create(ctx context.Context, newUser entities.User) error
	GetByID(ctx context.Context, ID int64) (entities.UserResponse, error)
	GetUserByLogin(ctx context.Context, userLogin entities.UserAuth) (entities.UserResponse, error)
}

type services struct {
	repositories *repositories.Container
	log          logger.Logger
}

func New(repo *repositories.Container, log logger.Logger) UserService {
	return &services{repositories: repo, log: log}
}

func (srv *services) Create(ctx context.Context, newUser entities.User) error {
	hasherBcrypt := hasher.NewBcryptHasher()
	passwordHashed, errHash := hasherBcrypt.Generate(newUser.Password)

	if errHash != nil {
		srv.log.Error("Srv.Create: ", "Error generate hash password: ", newUser)
		return errHash
	}

	newUser.Password = passwordHashed
	return srv.repositories.User.Create(ctx, newUser)
}

func (srv *services) GetUserByLogin(ctx context.Context, userLogin entities.UserAuth) (entities.UserResponse, error) {

	userFound, err := srv.repositories.User.GetUserByEmail(ctx, userLogin)

	if err != nil {
		srv.log.Error("Srv.Auth: ", "User not found", userFound)
		return entities.UserResponse{}, err
	}

	hasherBcrypt := hasher.NewBcryptHasher()
	err = hasherBcrypt.Compare(userFound.Password, userLogin.Password)
	if err != nil {
		return entities.UserResponse{}, err
	}

	return entities.UserResponse{
		ID:        userFound.ID,
		NickName:  userFound.NickName,
		Name:      userFound.Name,
		Email:     userFound.Email,
		CreatedAt: userFound.CreatedAt,
		UpdatedAt: userFound.UpdatedAt,
	}, nil
}

func (srv *services) GetByID(ctx context.Context, ID int64) (entities.UserResponse, error) {
	return srv.repositories.User.GetByID(ctx, ID)
}
