package api

import (
	"context"

	"github.com/go-chi/chi/v5"

	"github.com/ggrrrr/urlshortener-svc/be/shorturl/models"
)

type (
	appShortURL interface {
		GetLongURL(ctx context.Context, key string) (url string, err error)
		Create(ctx context.Context, request models.CreateShortURL) (*models.Key, error)
		Delete(ctx context.Context, request models.DeleteShortURL) error
		Update(ctx context.Context, request models.UpdateShortURL) error
		ListForOwner(ctx context.Context) ([]*models.ShortURLRecord, error)
	}

	server struct {
		app appShortURL
	}
)

func CreateRouter(app appShortURL) *chi.Mux {
	s := server{app: app}
	router := chi.NewRouter()

	router.Post("/admin/v1", s.handleCreateShortURL)
	router.Delete("/admin/v1", s.handleDeleteShortURL)
	router.Put("/admin/v1", s.handleUpdateShortURL)
	router.Get("/admin/v1", s.handleListShortURL)

	router.Get("/{key}", s.handleForward)
	return router
}
