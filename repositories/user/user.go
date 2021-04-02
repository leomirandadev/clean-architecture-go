package user

import (
	"context"

	"github.com/leomirandadev/clean-architecture-go/entities"
)

type UserRepository interface {
	Create(ctx context.Context, newUser entities.User) error
	GetUserByEmail(ctx context.Context, userLogin entities.UserAuth) ([]entities.User, error)
	GetByID(ctx context.Context, ID int64) ([]entities.UserResponse, error)
	Migrate()
}
