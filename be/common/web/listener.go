package web

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ggrrrr/urlshortener-svc/be/common/roles"
)

type (
	Config struct {
		Timeout    time.Duration `env:"HTTP_TTL" envDefault:"10s"`
		ListenAddr string        `env:"LISTEN_ADDR" envDefault:":8080"`
		CORSHosts  string        `env:"CORS_HOSTS"`
	}

	Listener struct {
		cfg        Config
		verifier   roles.TokenVerifier
		httpServer *http.Server
		mux        *chi.Mux
	}
)

func NewListener(cfg Config, verifier roles.TokenVerifier) (*Listener, error) {

	s := &Listener{
		cfg:      cfg,
		verifier: verifier,
	}

	mux := chi.NewRouter()
	mux.Use(middleware.Heartbeat("/ping"))
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(s.httpHandlerAuth)
	if cfg.CORSHosts != "" {
		mux.Use(s.httpHandlerCORS)
	}
	s.mux = mux

	webServer := &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: mux,
	}

	s.httpServer = webServer
	return s, nil
}

func (l *Listener) MountAPI(pattern string, api http.Handler) {
	l.mux.Mount(pattern, api)
}

func (l *Listener) Start(ctx context.Context) error {
	slog.Info("Start", slog.Any("ListenAddr", l.cfg.ListenAddr))
	return l.httpServer.ListenAndServe()
}

func (l *Listener) Shutdown(ctx context.Context) error {
	defer slog.Info("Shutdown", slog.Any("ListenAddr", l.cfg.ListenAddr))
	return l.httpServer.Shutdown(ctx)
}
