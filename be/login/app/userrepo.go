package app

import (
	"context"
	"fmt"
)

type (
	DummyUserRepo struct{}
)

// GetByUsername implements loginUserRepo.
func (d *DummyUserRepo) GetByUsername(ctx context.Context, username string) (*User, error) {
	if username == "admin" {
		return &User{
			Username: "admin",
			Password: "mypass",
		}, nil
	}
	if username == "err" {
		return nil, fmt.Errorf("dummy user repo")
	}
	return nil, nil
}

var _ (loginUserRepo) = (*DummyUserRepo)(nil)
