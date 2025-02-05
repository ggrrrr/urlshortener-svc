package app

import (
	"context"
	"errors"

	"github.com/ggrrrr/urlshortener-svc/common/roles"
	"github.com/ggrrrr/urlshortener-svc/internal/keygenerator"
	"github.com/ggrrrr/urlshortener-svc/internal/repo"
	"github.com/ggrrrr/urlshortener-svc/models"
)

type (
	Application struct {
		//slog
		//otel
		repo      repo.ShortURLRepo
		generator keygenerator.KeyGenerator
	}
)

func (a *Application) GetLongURL(ctx context.Context, key string) (url string, err error) {
	url, err = a.repo.GetByKey(ctx, key)
	return
}

func (a *Application) Create(ctx context.Context, user roles.AuthenticatedUser, request models.CreateShortURL) (string, error) {
	var err error

	var newKey string
	newKey, err = a.generator.Generate(request.LongURL)
	if err != nil {
		return "", err
	}

	err = a.repo.Create(ctx, repo.NewRecord{
		Owner:   user.Username,
		Key:     newKey,
		LongURL: request.LongURL,
	})

	if err != nil {
		if errors.Is(err, repo.ErrRecordExists) {
			//loop
		}
	}

	return newKey, nil
}

func (*Application) Delete(ctx context.Context, user roles.AuthenticatedUser, request models.DeleteShortURL) error {
	return nil
}

func (*Application) Update(ctx context.Context, user roles.AuthenticatedUser, request models.UpdateShortURL) error {
	return nil
}

func (*Application) ListForOwner(ctx context.Context, user roles.AuthenticatedUser) ([]models.ShortURLRecord, error) {
	return nil, nil
}
