package api

import (
	"context"

	"github.com/go-chi/chi/v5"

	"github.com/ggrrrr/urlshortener-svc/be/login/models"
)

type (
	appLogin interface {
		Login(ctx context.Context, req models.UserPasswordRequest) (string, error)
	}

	server struct {
		app appLogin
	}
)

func CreateRouter(app appLogin) *chi.Mux {
	s := server{app: app}
	router := chi.NewRouter()

	router.Post("/v1", s.handleLogin)

	return router
}
