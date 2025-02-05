package app

import (
	"context"

	"github.com/ggrrrr/urlshortener-svc/be/common/application"
	"github.com/ggrrrr/urlshortener-svc/be/common/roles"
	"github.com/ggrrrr/urlshortener-svc/be/login/models"
)

type (
	User struct {
		Username string
		Password string
	}

	Application struct {
		repo           loginUserRepo
		tokenGenerator roles.TokenGenerator
	}

	loginUserRepo interface {
		GetByUsername(ctx context.Context, username string) (*User, error)
	}
)

func New(repo loginUserRepo, tokener roles.TokenGenerator) *Application {
	return &Application{
		repo:           repo,
		tokenGenerator: tokener,
	}
}

func (a Application) Login(ctx context.Context, req models.UserPasswordRequest) (string, error) {
	user, err := a.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return "", application.NewSystemError("unable to fetch user", err)
	}
	if user == nil {
		// User not found
		return "", application.NewUnauthorized()
	}

	// Quick easy verification of password
	// not good for production
	if user.Password != req.Password {
		// Password dont match
		return "", application.NewUnauthorized()
	}

	token, err := a.tokenGenerator.Generate(ctx, user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}
