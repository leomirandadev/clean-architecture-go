package hasher

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type bcryptHasher struct{}

func NewBcryptHasher() Hasher {
	return &bcryptHasher{}
}

func (*bcryptHasher) Generate(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return "", errors.New("Erro de criptografia")
	}

	return string(bytes), nil
}

func (*bcryptHasher) Compare(password1 string, password2 string) error {

	return bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2))

}
