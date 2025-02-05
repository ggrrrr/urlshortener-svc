package api

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/ggrrrr/urlshortener-svc/models"
)

type MockApp struct {
	mock.Mock
}

// Create implements app.
func (m *MockApp) Create(ctx context.Context, request models.CreateShortURL) (models.ShortURLRecord, error) {
	args := m.Called(request.LongURL)
	return args.Get(0).(models.ShortURLRecord), args.Error(1)
}

// Delete implements app.
func (m *MockApp) Delete(ctx context.Context, request models.DeleteShortURL) error {
	panic("Delete unimplemented")
}

// GetLongURL implements app.
func (m *MockApp) GetLongURL(ctx context.Context, key string) (url string, err error) {
	args := m.Called(key)
	return args.Get(0).(string), args.Error(1)
}

// ListForOwner implements app.
func (m *MockApp) ListForOwner(ctx context.Context) ([]models.ShortURLRecord, error) {
	panic("ListForOwner unimplemented")
}

// Update implements app.
func (m *MockApp) Update(ctx context.Context, request models.UpdateShortURL) error {
	panic("Update unimplemented")
}

var _ (app) = (*MockApp)(nil)
