package main

import (
	"context"

	"github.com/caarlos0/env/v11"

	"github.com/ggrrrr/urlshortener-svc/be/common/roles"
	"github.com/ggrrrr/urlshortener-svc/be/common/system"
	"github.com/ggrrrr/urlshortener-svc/be/common/web"
	loginAPI "github.com/ggrrrr/urlshortener-svc/be/login/api"
	loginApp "github.com/ggrrrr/urlshortener-svc/be/login/app"
	shortAPI "github.com/ggrrrr/urlshortener-svc/be/shorturl/internal/api"
	shortApp "github.com/ggrrrr/urlshortener-svc/be/shorturl/internal/app"
	"github.com/ggrrrr/urlshortener-svc/be/shorturl/internal/repo/pg"
)

func main() {
	webCfg := web.Config{}
	err := env.Parse(&webCfg)
	if err != nil {
		panic(err)
	}

	repoCfg := pg.Config{}
	err = env.Parse(&repoCfg)
	if err != nil {
		panic(err)
	}

	db, err := pg.Connect(repoCfg)
	if err != nil {
		panic(err)
	}

	err = pg.Up(db)
	if err != nil {
		panic(err)
	}

	apiRouter := shortAPI.CreateRouter(shortApp.New(pg.NewRepo(db)))

	httpListener, err := web.NewListener(webCfg, roles.NewDummyGenerator())
	if err != nil {
		panic(err)
	}

	loginRouter := loginAPI.CreateRouter(loginApp.New(&loginApp.DummyUserRepo{}, roles.NewDummyGenerator()))
	httpListener.MountAPI("/login", loginRouter)
	httpListener.MountAPI("/", apiRouter)

	system, err := system.New(system.WithHTTPListener(httpListener))
	if err != nil {
		panic(err)
	}

	system.Start(context.Background())
}
