package httpAdapter

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/EvansTrein/Workmate/pkg/validate"
)

// Universal structure for sending responses
type BaseHandlerResponce struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Status  int    `json:"status"`
}

// Common handler for all other handlers
// is needed to remove code duplication, as a single place to send http responses and handle errors
// use -> embed in any handler and it gets logger and validation capabilities, read request body, send errors
type BaseHandler struct {
	Log *slog.Logger
}

type BaseHandlerDeps struct {
	*slog.Logger
}

func NewBaseHandler(deps *BaseHandlerDeps) *BaseHandler {
	return &BaseHandler{Log: deps.Logger}
}

// Records the data in the response and puts a response code in the response
func (h *BaseHandler) SendJsonResp(w http.ResponseWriter, status int, data any) {
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		h.Log.Error("failed to marshal JSON", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(jsonResponse); err != nil {
		h.Log.Error("!!ATTENTION!! failed to write JSON response", "error", err)
	}
}

// The common, basic handle, deals with reading the request body and validating it
func (h *BaseHandler) Handle(w http.ResponseWriter, r *http.Request, decodeFunc func(io.Reader) (any, error)) (any, error) {
	op := "BaseHandler: Handle func"
	log := h.Log.With(slog.String("operation", op))

	// Decoding the request body
	data, err := decodeFunc(r.Body)
	if err != nil {
		log.Warn("failed to decode request body", "error", err)
		h.SendJsonResp(w, 400, &BaseHandlerResponce{
			Status:  http.StatusBadRequest,
			Message: "failed to decode request body",
			Error:   err.Error(),
		})
		return nil, err
	}

	// Data validation
	if err := validate.IsValid(data); err != nil {
		log.Warn("request body data failed validation", "error", err)
		h.SendJsonResp(w, 400, &BaseHandlerResponce{
			Status:  http.StatusBadRequest,
			Message: "request body data failed validation",
			Error:   err.Error(),
		})
		return nil, err
	}

	log.Debug("data successfully validated", "data", data)
	return data, nil
}

// General function for submitting errors, this is where ALL errors from ALL locations are collected
// If you want the http response to be specific, you can add the error you want here
func (h *BaseHandler) HandleError(w http.ResponseWriter, err error) {
	op := "BaseHandler: HandleError func"
	log := h.Log.With(slog.String("operation", op))

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		log.Error("request processing exceeded the allowed time limit", "error", err)
		h.SendJsonResp(w, 504, &BaseHandlerResponce{
			Status:  http.StatusGatewayTimeout,
			Message: "request processing exceeded the allowed time limit",
			Error:   err.Error(),
		})
	default:
		log.Error("internal server error", "error", err)
		h.SendJsonResp(w, 500, &BaseHandlerResponce{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
			Error:   err.Error(),
		})
	}
}
