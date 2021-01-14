package token

type TokenHash interface {
	Encrypt(data interface{}) (string, error)
	Decrypt(bearerToken string) (interface{}, error)
}
