package hasher

import "errors"

type Hasher interface {
	Generate(salt, secret string) (string, error)
	Compare(secretHash, salt, secret string) error
	Salt() (string, error)
}

var ErrCript = errors.New("cryptography error")

const saltSize = 15
