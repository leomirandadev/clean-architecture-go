package hasher

import (
	"testing"
)

func TestGenerateAndCompare(t *testing.T) {
	h := NewBcryptHasher()

	firstPassword := "password"
	secondPassword := "password2"

	salt, err := h.Salt()
	if err != nil {
		t.Error(err)
	}

	password1Hash, err := h.Generate(salt, firstPassword)
	if err != nil {
		t.Error(err)
	}

	if err := h.Compare(password1Hash, salt, firstPassword); err != nil {
		t.Error(err)
	}

	if err := h.Compare(password1Hash, salt, secondPassword); err == nil {
		t.Error("password can't be equal")
	}
}
