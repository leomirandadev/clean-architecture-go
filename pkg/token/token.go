package token

type TokenHash interface {
	Encrypt(data any) (string, error)
	Decrypt(bearerToken string) (bool, map[string]any, error)
}
