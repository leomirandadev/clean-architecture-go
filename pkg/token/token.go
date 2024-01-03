package token

import "time"

type TokenHash interface {
	Encrypt(data any) (string, error)
	Decrypt(bearerToken string) (bool, map[string]any, error)
}

type Options struct {
	Key        string        `mapstructure:"KEY"`
	Expiration time.Duration `mapstructure:"EXPIRATION"`
}
