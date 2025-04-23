package app

import (
	"log/slog"

	"github.com/EvansTrein/Workmate/config"
	"github.com/EvansTrein/Workmate/internal/server"
)

type App struct {
	conf   *config.Config
	log    *slog.Logger
	server *server.HttpServer
}

type AppDeps struct {
	*config.Config
	*slog.Logger
}

func New(deps *AppDeps) *App {

	httpServer := server.New(&server.HttpServerDeps{
		HTTPServer: &deps.HTTPServer,
		Logger:     deps.Logger,
	})

	return &App{
		conf:   deps.Config,
		log:    deps.Logger,
		server: httpServer,
	}
}

func (a *App) MustStart() {
	a.log.Debug("app: started")

	a.log.Info("app: successfully started", "port", a.conf.HTTPServer.Port)
	if err := a.server.Start(); err != nil {
		panic(err)
	}
}

func (a *App) Stop() error {
	a.log.Debug("app: stop started")

	if err := a.server.Stop(); err != nil {
		a.log.Error("failed to stop HTTP server")
		return err
	}

	a.server = nil
	a.log.Info("app: stop successful")

	return nil
}
