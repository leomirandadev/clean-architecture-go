package reset_passwords

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/leomirandadev/clean-architecture-go/internal/models"
	"github.com/leomirandadev/clean-architecture-go/pkg/tracer"
)

type IRepository interface {
	Create(ctx context.Context, req models.CreateResetPassword) (string, error)
	GetByToken(ctx context.Context, token string) (*models.ResetPassword, error)
}

func NewSqlx(writer, reader *sqlx.DB) IRepository {
	return &repoSqlx{writer: writer, reader: reader}
}

type repoSqlx struct {
	writer *sqlx.DB
	reader *sqlx.DB
}

func (repo repoSqlx) Create(ctx context.Context, req models.CreateResetPassword) (string, error) {
	ctx, tr := tracer.Span(ctx, "repositories.reset_passwords.create")
	defer tr.End()

	req.ID = uuid.New().String()

	_, err := repo.writer.NamedExecContext(ctx, `
		INSERT INTO reset_passwords (id, token, user_id) 
		VALUES (:id, :token, :user_id)
	`, req)

	if err != nil {
		slog.WarnContext(ctx, "writer.exec_context", "err", err)
		return "", err
	}

	return req.ID, nil
}

func (repo repoSqlx) GetByToken(ctx context.Context, token string) (*models.ResetPassword, error) {
	ctx, tr := tracer.Span(ctx, "repositories.reset_passwords.get_by_token")
	defer tr.End()

	resp := new(models.ResetPassword)

	err := repo.reader.GetContext(ctx, resp, `
		SELECT * FROM reset_passwords WHERE token = $1 LIMIT 1
	`, token)

	if err != nil {
		slog.WarnContext(ctx, "reader.get_context", "err", err)
		return nil, err
	}

	return resp, nil
}
