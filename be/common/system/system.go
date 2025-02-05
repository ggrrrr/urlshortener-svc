package system

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/ggrrrr/urlshortener-svc/be/common/web"
)

type (
	SystemOptionFunc func(s *System) error

	System struct {
		listener *web.Listener

		shutdownFunc []func(context.Context) error
		startupFunc  []func(context.Context) error
	}
)

func New(opts ...SystemOptionFunc) (*System, error) {

	system := &System{
		shutdownFunc: make([]func(context.Context) error, 0, 1),
		startupFunc:  make([]func(context.Context) error, 0, 1),
	}

	for _, optFn := range opts {
		err := optFn(system)
		if err != nil {
			return nil, err
		}
	}

	return system, nil

}

func (s *System) Start(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// }

	g.Go(func() error {
		<-ctx.Done()
		cancel()
		slog.Info("git kill signal")

		for _, shutdownFunc := range s.shutdownFunc {
			err := shutdownFunc(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})

	g.Go(func() error {
		for _, startFunc := range s.startupFunc {
			err := startFunc(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return g.Wait()
}

func WithHTTPListener(l *web.Listener) SystemOptionFunc {
	return func(s *System) error {
		s.listener = l
		s.startupFunc = append(s.startupFunc, func(ctx context.Context) error {
			return l.Start(ctx)
		})
		s.shutdownFunc = append(s.shutdownFunc, func(ctx context.Context) error {
			return l.Shutdown(ctx)
		})
		return nil
	}
}
