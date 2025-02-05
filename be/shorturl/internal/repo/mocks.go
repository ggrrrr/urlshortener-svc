package repo

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

// Create implements ShortURLRepo.
func (m *MockRepo) Create(ctx context.Context, record NewRecord) error {
	args := m.Called(record)
	return args.Error(0)
}

// Delete implements ShortURLRepo.
func (m *MockRepo) Delete(ctx context.Context, key string) error {
	args := m.Called(key)
	return args.Error(0)
}

// GetByKey implements ShortURLRepo.
func (m *MockRepo) GetByKey(ctx context.Context, key string) (*URLRecord, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*URLRecord), args.Error(1)
}

// ListByOwner implements ShortURLRepo.
func (m *MockRepo) ListByOwner(ctx context.Context, owner string) ([]*URLRecord, error) {
	args := m.Called(owner)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*URLRecord), args.Error(1)
}

// Update implements ShortURLRepo.
func (m *MockRepo) Update(ctx context.Context, key string, longURL string) error {
	args := m.Called(key, longURL)
	return args.Error(0)
}

var _ (ShortURLRepo) = (*MockRepo)(nil)
