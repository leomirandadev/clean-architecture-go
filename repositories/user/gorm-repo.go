package user

import (
	"context"
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/leomirandadev/clean-architecture-go/entities"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
)

type repoGorm struct {
	log    logger.Logger
	writer *gorm.DB
	reader *gorm.DB
}

func NewGormRepository(log logger.Logger, writer, reader *gorm.DB) UserRepository {
	return &repoGorm{log: log}
}

func (repo *repoGorm) Migrate() {
	defer repo.writer.Close()

	repo.writer.AutoMigrate(&entities.User{})
}

func (repo *repoGorm) Create(ctx context.Context, newUser entities.User) error {
	defer repo.writer.Close()

	repo.writer.Table("users").Create(&newUser)

	return repo.writer.Error
}

func (repo *repoGorm) GetUserByEmail(ctx context.Context, userLogin entities.UserAuth) ([]entities.User, error) {
	defer repo.reader.Close()

	var user []entities.User
	repo.reader.Table("users").Where("email = ?", userLogin.Email).Find(&user)

	if repo.reader.Error != nil {
		repo.log.Error("GormRepo.GetByID: ", "Error on get User ID: ", userLogin, repo.reader.Error)
		return nil, errors.New("Usuário não encontrado")
	}

	return user, nil
}

func (repo *repoGorm) GetByID(ctx context.Context, ID int64) ([]entities.UserResponse, error) {
	defer repo.reader.Close()

	var user []entities.UserResponse
	repo.reader.Table("users").Where("ID = ?", ID).Find(&user)

	if repo.reader.Error != nil {
		repo.log.Error("GormRepo.GetByID: ", "Error on get User ID: ", ID, repo.reader.Error)
		return nil, errors.New("Usuário não encontrado")
	}

	return user, nil
}
