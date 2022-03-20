package user

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/leomirandadev/clean-architecture-go/entities"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
)

type repoSqlx struct {
	log    logger.Logger
	writer *sqlx.DB
	reader *sqlx.DB
}

func NewSqlxRepository(log logger.Logger, writer, reader *sqlx.DB) UserRepository {
	return &repoSqlx{log: log, writer: writer, reader: reader}
}

func (repo *repoSqlx) Create(ctx context.Context, newUser entities.User) error {

	_, err := repo.writer.ExecContext(ctx, `
		INSERT INTO users (nick_name,name,email,password,role) VALUES (:nick_name,:name,:email,:password,:role)
	`, newUser)

	return err
}

func (repo *repoSqlx) GetByID(ctx context.Context, ID int64) (entities.UserResponse, error) {

	var user entities.UserResponse

	err := repo.reader.GetContext(ctx, &user, `
		SELECT 
			id,
			nick_name,
			name,
			email,
			role,
			created_at,
			updated_at
		FROM users 
		WHERE id = ?
		LIMIT 1
	`, ID)

	if err != nil {
		repo.log.Error("SqlxRepo.GetByID", "Error on get User ID: ", ID, err)
		return user, errors.New("Usuário não encontrado")
	}

	return user, nil

}

func (repo *repoSqlx) GetUserByEmail(ctx context.Context, userLogin entities.UserAuth) (entities.User, error) {
	var user entities.User

	err := repo.reader.GetContext(ctx, &user, `
		SELECT 
			id,
			nick_name,
			name,
			email,
			role,
			password,
			created_at,
			updated_at
		FROM users 
		WHERE email = ?
		LIMIT 1
	`, userLogin.Email)

	if err != nil {
		repo.log.Error("SqlxRepo.GetByID", "Error on get User ID: ", userLogin, err)
		return user, errors.New("Usuário não encontrado")
	}

	return user, nil
}
