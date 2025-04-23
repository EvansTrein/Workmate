package httpAdapter

import (
	"context"
	"io"
	"net/http"

	"github.com/EvansTrein/Workmate/internal/dto"
	myErr "github.com/EvansTrein/Workmate/pkg/error"
	"github.com/EvansTrein/Workmate/pkg/utils"
)

type ITaskUseCase interface {
	CreateTask(ctx context.Context, data *dto.TaskCreateRequest) (*dto.TaskCreateResponce, error)
	AllTasks(ctx context.Context) (*dto.AllTasksResponce, error)
}

type TaskHandler struct {
	baseH *BaseHandler
	uc    ITaskUseCase
}

type TaskHandlerDeps struct {
	*BaseHandler
	ITaskUseCase
}

func NewTaskHandler(deps *TaskHandlerDeps) *TaskHandler {
	return &TaskHandler{
		baseH: deps.BaseHandler,
		uc:    deps.ITaskUseCase,
	}
}

func (h *TaskHandler) Task() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := h.baseH.Handle(w, r, func(body io.Reader) (any, error) {
			return utils.DecodeBody[dto.TaskCreateRequest](r.Body)
		})
		if err != nil {
			// the error in the response has already been written to h.baseH.Handle
			return
		}

		taskData, ok := data.(*dto.TaskCreateRequest)
		if !ok {
			h.baseH.HandleError(w, myErr.ErrTypeConversion)
			return
		}

		resp, err := h.uc.CreateTask(r.Context(), taskData)
		if err != nil {
			h.baseH.HandleError(w, err)
			return
		}

		h.baseH.SendJsonResp(w, 201, resp)
	}
}

func (h *TaskHandler) Status() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		resp, err := h.uc.AllTasks(r.Context())
		if err != nil {
			h.baseH.HandleError(w, err)
			return
		}

		h.baseH.SendJsonResp(w, 200, resp)
	}
}
