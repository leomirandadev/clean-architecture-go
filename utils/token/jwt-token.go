package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("captainjacksparrowsayshi")

type jwtToken struct{}

func NewJWT() TokenHash {
	return &jwtToken{}
}

func (*jwtToken) Encrypt(data interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["data"] = data
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func (*jwtToken) Decrypt(bearerToken string) (interface{}, error) {

	return jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return mySigningKey, nil
	})
}
