package token

import "time"

type TokenHash interface {
	Encrypt(data any) (string, error)
	Decrypt(bearerToken string) (bool, map[string]any, error)
}

type Options struct {
	Key        string        `mapstructure:"KEY" validate:"required"`
	Expiration time.Duration `mapstructure:"EXPIRATION" validate:"required"`
}
