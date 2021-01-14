package token

type TokenHash interface {
	Encrypt(data interface{}) (Output, error)
	Decrypt(bearerToken string) (interface{}, error)
}
