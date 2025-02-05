package app

import (
	"context"
	"log/slog"

	"github.com/ggrrrr/urlshortener-svc/be/common/application"
	"github.com/ggrrrr/urlshortener-svc/be/common/roles"
	"github.com/ggrrrr/urlshortener-svc/be/shorturl/internal/keygenerator"
	"github.com/ggrrrr/urlshortener-svc/be/shorturl/internal/repo"
	"github.com/ggrrrr/urlshortener-svc/be/shorturl/models"
)

type (
	Application struct {
		repo      repo.ShortURLRepo
		generator keygenerator.KeyGenerator
	}
)

func New(urlRepo repo.ShortURLRepo) *Application {
	return &Application{
		repo:      urlRepo,
		generator: &keygenerator.Generator{},
	}
}

func (a *Application) GetLongURL(ctx context.Context, key string) (url string, err error) {
	record, err := a.repo.GetByKey(ctx, key)
	if err != nil {
		return "", err
	}
	url = record.LongURL
	return
}

func (a *Application) Create(ctx context.Context, request models.CreateShortURL) (string, error) {
	var err error

	authUser, err := roles.ExtractUser(ctx)
	if err != nil {
		return "", err
	}

	newRecord := repo.NewRecord{
		Key:     "",
		Owner:   authUser.Username,
		LongURL: request.LongURL,
	}

	newKey, err := a.retryLoop(ctx, newRecord, 10)
	if err != nil {
		return "", application.NewSystemError("unable to create short url", err)
	}

	return newKey, nil
}

func (a *Application) retryLoop(ctx context.Context, newRecord repo.NewRecord, counter int) (string, error) {
	counter -= 1
	newRecord.Key = a.generator.Generate()

	slog.Info("retryLoop", slog.Any("key", newRecord.Key), slog.Int("counter", counter))

	err := a.repo.Create(ctx, newRecord)
	if err != nil {
		if counter > 0 {
			return a.retryLoop(ctx, newRecord, counter)
		}
		return "", err
	}
	return newRecord.Key, nil
}

func (a *Application) Delete(ctx context.Context, request models.DeleteShortURL) error {
	authUser, err := roles.ExtractUser(ctx)
	if err != nil {
		return err
	}

	record, err := a.repo.GetByKey(ctx, request.Key)
	if err != nil {
		return err
	}

	if record.Owner != authUser.Username {
		return application.NewForbidden("not allowed")
	}

	err = a.repo.Delete(ctx, request.Key)
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) Update(ctx context.Context, request models.UpdateShortURL) error {
	authUser, err := roles.ExtractUser(ctx)
	if err != nil {
		return err
	}

	record, err := a.repo.GetByKey(ctx, request.Key)
	if err != nil {
		return err
	}

	if record.Owner != authUser.Username {
		return application.NewForbidden("not allowed")
	}

	err = a.repo.Update(ctx, request.Key, request.LongURL)
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) ListForOwner(ctx context.Context) ([]*models.ShortURLRecord, error) {
	authUser, err := roles.ExtractUser(ctx)
	if err != nil {
		return nil, err
	}

	records, err := a.repo.ListByOwner(ctx, authUser.Username)
	if err != nil {
		return nil, err
	}

	out := make([]*models.ShortURLRecord, 0, len(records))
	for _, r := range records {
		out = append(out, &models.ShortURLRecord{
			Key:       r.Key,
			LongURL:   r.LongURL,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}

	return out, nil
}
