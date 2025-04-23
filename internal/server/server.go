package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/EvansTrein/Workmate/config"
	httpAdapter "github.com/EvansTrein/Workmate/internal/adapters/http"
	"github.com/EvansTrein/Workmate/internal/adapters/repository"
	"github.com/EvansTrein/Workmate/internal/usecase"
	"github.com/EvansTrein/Workmate/pkg/middleware"
)

const (
	gracefulShutdownTimer = time.Second * 10
)

type HttpServer struct {
	conf   *config.HTTPServer
	log    *slog.Logger
	server *http.Server
	router *http.ServeMux
}

type HttpServerDeps struct {
	*config.HTTPServer
	*slog.Logger
}

func New(deps *HttpServerDeps) *HttpServer {
	router := http.NewServeMux()

	// Initialize repository
	repoTask := repository.NewTaskRepo()

	// Initialize use case
	taskUC := usecase.NewWorkUseCase(&usecase.WorkUseCaseDeps{
		Logger:    deps.Logger,
		ITaskRepo: repoTask,
	})

	// Initialize handler
	baseHandler := httpAdapter.NewBaseHandler(&httpAdapter.BaseHandlerDeps{
		Logger: deps.Logger,
	})

	taskHandler := httpAdapter.NewTaskHandler(&httpAdapter.TaskHandlerDeps{
		BaseHandler:  baseHandler,
		ITaskUseCase: taskUC,
	})

	// Initialize Routers
	activeHandlers := &ActiveHandlers{
		TaskHandler: taskHandler,
	}

	activeMiddlewares := &ActiveMiddlewares{}

	InitRouters(router, activeHandlers, activeMiddlewares)

	return &HttpServer{
		conf:   deps.HTTPServer,
		log:    deps.Logger,
		router: router,
	}
}

func (s *HttpServer) Start() error {
	log := s.log.With(slog.String("Address", s.conf.Address+":"+s.conf.Port))
	log.Debug("HTTP server: started creating")

	LoggerHTTP := middleware.NewMiddlewareLogging(&middleware.MiddlewareLoggingDeps{
		Logger: s.log,
	})

	s.server = &http.Server{
		Addr: s.conf.Address + ":" + s.conf.Port,
		Handler: middleware.ChainMiddleware(
			middleware.CORS,
			LoggerHTTP.HandlersLog(),
		)(s.router),
	}

	log.Info("HTTP server: successfully started")
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *HttpServer) Stop() error {
	s.log.Debug("HTTP server: stop started")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimer)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Error("Server shutdown failed", "error", err)
		return err
	}

	s.server = nil
	s.log.Info("HTTP server: stop successful")

	return nil
}
