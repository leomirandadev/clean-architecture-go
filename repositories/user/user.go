package user

import "github.com/leomirandadev/clean-architecture-go/entities"

type UserRepository interface {
	Create(newUser entities.User) error
	GetByID(ID int64) ([]entities.User, error)
	Migrate()
}
