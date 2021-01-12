package hasher

type Hasher interface {
	Generate(password string) (string, error)
	Compare(password1 string, password2 string) error
}
