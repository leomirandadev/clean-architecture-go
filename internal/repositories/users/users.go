package users

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/leomirandadev/clean-architecture-go/internal/models"
	"github.com/leomirandadev/clean-architecture-go/pkg/tracer"
)

type IRepository interface {
	Create(ctx context.Context, newUser models.User) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, ID string) (*models.User, error)
	ChangePassword(ctx context.Context, userID, newPassword, salt string) error
}

func NewSqlx(writer, reader *sqlx.DB) IRepository {
	return &repoSqlx{writer: writer, reader: reader}
}

type repoSqlx struct {
	writer *sqlx.DB
	reader *sqlx.DB
}

func (repo repoSqlx) Create(ctx context.Context, newUser models.User) (string, error) {
	ctx, tr := tracer.Span(ctx, "repositories.users.create")
	defer tr.End()

	newUser.ID = uuid.New().String()

	_, err := repo.writer.NamedExecContext(ctx, `
		INSERT INTO users (id, nick_name, name, email, password, photo_url, salt) 
		VALUES (:id, :nick_name, :name, :email, :password, :photo_url, :salt)
	`, newUser)

	if err != nil {
		slog.WarnContext(ctx, "writer.exec_context", "err", err)
		return "", err
	}

	return newUser.ID, nil
}

func (repo repoSqlx) ChangePassword(ctx context.Context, userID, newPassword, salt string) error {
	ctx, tr := tracer.Span(ctx, "repositories.users.change_password")
	defer tr.End()

	_, err := repo.writer.ExecContext(ctx, `
		UPDATE users SET password = $1, salt = $2 WHERE id = $3
	`, newPassword, salt, userID)

	if err != nil {
		slog.Warn("exec_context fails", "error", err)
		return err
	}

	return nil
}

func (repo repoSqlx) GetByID(ctx context.Context, ID string) (*models.User, error) {
	ctx, tr := tracer.Span(ctx, "repositories.users.get_by_id")
	defer tr.End()

	user := new(models.User)

	err := repo.reader.GetContext(ctx, user, `
		SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL LIMIT 1
	`, ID)

	if err != nil {
		slog.WarnContext(ctx, "reader.get_context", "err", err)
		return nil, err
	}

	return user, nil
}

func (repo repoSqlx) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, tr := tracer.Span(ctx, "repositories.users.get_user_by_email")
	defer tr.End()

	user := new(models.User)

	err := repo.reader.GetContext(ctx, user, `
		SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL LIMIT 1
	`, email)

	if err != nil {
		slog.WarnContext(ctx, "reader.get_context", "err", err)
		return nil, err
	}

	return user, nil
}
