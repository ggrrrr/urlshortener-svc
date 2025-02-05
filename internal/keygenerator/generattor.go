package keygenerator

type (
	KeyGenerator interface {
		Generate(longURL string) (string, error)
	}
)
