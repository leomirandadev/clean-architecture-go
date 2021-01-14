package user

import (
	"errors"

	"github.com/leomirandadev/clean-architecture-go/configs"
	"github.com/leomirandadev/clean-architecture-go/entities"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
)

type repoGorm struct {
	log logger.Logger
}

func NewGormRepository(log logger.Logger) UserRepository {
	return &repoGorm{log: log}
}

func (repo *repoGorm) Migrate() {
	db := configs.ConnectGorm()
	defer db.Close()

	db.AutoMigrate(&entities.User{})
}

func (repo *repoGorm) Create(newUser entities.User) error {
	db := configs.ConnectGorm()
	defer db.Close()

	db.Table("users").Create(&newUser)

	return db.Error
}

func (repo *repoGorm) GetUserByEmail(userLogin entities.UserAuth) ([]entities.User, error) {
	db := configs.ConnectGorm()
	defer db.Close()

	var user []entities.User
	db.Table("users").Where("email = ?", userLogin.Email).Find(&user)

	if db.Error != nil {
		repo.log.Error("GormRepo.GetByID: ", "Error on get User ID: ", userLogin, db.Error)
		return nil, errors.New("Usuário não encontrado")
	}

	return user, nil
}

func (repo *repoGorm) GetByID(ID int64) ([]entities.UserResponse, error) {
	db := configs.ConnectGorm()
	defer db.Close()

	var user []entities.UserResponse
	db.Table("users").Where("ID = ?", ID).Find(&user)

	if db.Error != nil {
		repo.log.Error("GormRepo.GetByID: ", "Error on get User ID: ", ID, db.Error)
		return nil, errors.New("Usuário não encontrado")
	}

	return user, nil
}
