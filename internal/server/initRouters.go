package server

import (
	"net/http"

	httpAdapter "github.com/EvansTrein/Workmate/internal/adapters/http"
)

// Grouping of used handlers
type ActiveHandlers struct {
	*httpAdapter.TaskHandler
}

// Grouping of used Middlewares
type ActiveMiddlewares struct{}

func InitRouters(router *http.ServeMux, handlers *ActiveHandlers, middlewares *ActiveMiddlewares) {
	router.Handle("POST /task", handlers.Task())
	router.Handle("GET /status", handlers.Status())
}
