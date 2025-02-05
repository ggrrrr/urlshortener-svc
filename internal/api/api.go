package api

import (
	"context"

	"github.com/go-chi/chi/v5"

	"github.com/ggrrrr/urlshortener-svc/models"
)

type (
	app interface {
		GetLongURL(ctx context.Context, key string) (url string, err error)
		Create(ctx context.Context, request models.CreateShortURL) (models.ShortURLRecord, error)
		Delete(ctx context.Context, request models.DeleteShortURL) error
		Update(ctx context.Context, request models.UpdateShortURL) error
		ListForOwner(ctx context.Context) ([]models.ShortURLRecord, error)
	}

	server struct {
		app app
	}
)

func CreateRouter(app app) *chi.Mux {
	s := server{app: app}
	router := chi.NewRouter()

	router.Post("/admin/v1", s.handleCreateShortURL)
	router.Delete("/admin/v1", s.handleDeleteShortURL)
	router.Put("/admin/v1", s.handleUpdateShortURL)
	router.Get("/admin/v1", s.handleListShortURL)

	router.Get("/{key}", s.handleForward)
	return router
}
