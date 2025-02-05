package roles

import (
	"context"
	"fmt"
	"strings"

	"github.com/ggrrrr/urlshortener-svc/be/common/application"
)

type DummyGenerator struct{}

// Verify implements TokenVerifier.
func (m *DummyGenerator) Verify(ctx context.Context, token string) (AuthenticatedUser, error) {
	s := strings.Split(token, "@mocker")
	if len(s) != 1 {
		return AuthenticatedUser{}, application.NewUnauthorized()
	}
	return AuthenticatedUser{
		Username: s[0],
	}, nil
}

// Generate implements TokenGenerator.
func (m *DummyGenerator) Generate(ctx context.Context, username string) (string, error) {
	return fmt.Sprintf("%s@mocker", username), nil
}

var _ (TokenGenerator) = (*DummyGenerator)(nil)
var _ (TokenVerifier) = (*DummyGenerator)(nil)

func NewDummyGenerator() *DummyGenerator {
	return &DummyGenerator{}
}
