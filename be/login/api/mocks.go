package api

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/ggrrrr/urlshortener-svc/be/login/models"
)

type MockApp struct {
	mock.Mock
}

// Login implements appLogin.
func (m *MockApp) Login(ctx context.Context, req models.UserPasswordRequest) (string, error) {
	args := m.Called(req.Username)
	return args.Get(0).(string), args.Error(1)
}

var _ (appLogin) = (*MockApp)(nil)
