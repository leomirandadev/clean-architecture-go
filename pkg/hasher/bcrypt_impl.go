package hasher

import (
	"math/rand"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type bcryptHasher struct{}

func NewBcryptHasher() Hasher {
	return &bcryptHasher{}
}

const letterRunes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (bcryptHasher) Salt() (string, error) {
	var strBuilder strings.Builder
	sizeLetterRunes := len(letterRunes)

	for i := 0; i < saltSize; i++ {
		char := letterRunes[rand.Intn(sizeLetterRunes)]
		strBuilder.WriteByte(char)
	}

	return strBuilder.String(), nil
}

func (bcryptHasher) Generate(salt, secret string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(secret+salt), bcrypt.MinCost)
	if err != nil {
		return "", ErrCript
	}

	return string(bytes), nil
}

func (bcryptHasher) Compare(secretHash, salt, secret string) error {
	return bcrypt.CompareHashAndPassword([]byte(secretHash), []byte(secret+salt))
}
