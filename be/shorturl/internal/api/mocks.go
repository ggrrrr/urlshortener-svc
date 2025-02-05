package api

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/ggrrrr/urlshortener-svc/be/shorturl/models"
)

type MockApp struct {
	mock.Mock
}

// Create implements app.
func (m *MockApp) Create(ctx context.Context, request models.CreateShortURL) (string, error) {
	args := m.Called(request)
	return args.Get(0).(string), args.Error(1)
}

// Delete implements app.
func (m *MockApp) Delete(ctx context.Context, request models.DeleteShortURL) error {
	args := m.Called(request.Key)
	return args.Error(0)
}

// GetLongURL implements app.
func (m *MockApp) GetLongURL(ctx context.Context, key string) (url string, err error) {
	args := m.Called(key)
	return args.Get(0).(string), args.Error(1)
}

// ListForOwner implements app.
func (m *MockApp) ListForOwner(ctx context.Context) ([]*models.ShortURLRecord, error) {
	args := m.Called("key")
	return args.Get(0).([]*models.ShortURLRecord), args.Error(1)
}

// Update implements app.
func (m *MockApp) Update(ctx context.Context, request models.UpdateShortURL) error {
	args := m.Called(request.Key)
	return args.Error(0)
}

var _ (appShortURL) = (*MockApp)(nil)
