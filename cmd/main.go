package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/EvansTrein/Workmate/config"
	"github.com/EvansTrein/Workmate/internal/app"
	"github.com/EvansTrein/Workmate/pkg/logs"
)

func main() {
	var conf *config.Config
	var log *slog.Logger

	conf = config.MustLoad()
	log = logs.InitLog(conf.Env)

	application := app.New(&app.AppDeps{
		Config: conf,
		Logger: log,
	})

	go func() {
		application.MustStart()
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	if err := application.Stop(); err != nil {
		log.Error("an error occurred when stopping the application", "error", err)
		panic(err)
	}
}
