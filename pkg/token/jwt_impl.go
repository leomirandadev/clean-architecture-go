package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func NewJWT(opts Options) TokenHash {
	return &jwtToken{
		key: []byte(opts.Key),
		exp: opts.Expiration,
	}
}

type jwtToken struct {
	key []byte
	exp time.Duration
}

func (j jwtToken) Encrypt(data any) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["data"] = data
	claims["exp"] = time.Now().Add(j.exp).Unix()

	tokenString, err := token.SignedString(j.key)

	if err != nil {
		fmt.Printf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func (j jwtToken) Decrypt(bearerToken string) (bool, map[string]any, error) {

	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("there was an error")
		}
		return j.key, nil
	})

	var claims map[string]any

	if claimsMap, ok := token.Claims.(jwt.MapClaims); ok {
		if data, ok := claimsMap["data"]; ok {
			claims = data.(map[string]any)
		}
	}

	return token.Valid, claims, err
}
