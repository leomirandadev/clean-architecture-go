package models

import "time"

type ResetPassword struct {
	ID        string    `db:"id" json:"id"`
	Token     string    `db:"token" json:"token"`
	UserID    string    `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type CreateResetPassword struct {
	ID     string `db:"id" json:"id"`
	Token  string `db:"token" json:"token"`
	UserID string `db:"user_id" json:"user_id"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required"`
}

type ChangePassword struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}
