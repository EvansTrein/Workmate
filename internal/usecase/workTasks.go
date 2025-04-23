package usecase

import (
	"context"
	"errors"
	"log/slog"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/EvansTrein/Workmate/internal/dto"
	"github.com/EvansTrein/Workmate/internal/entity"
)

const (
	minDelay = time.Minute * 3
	maxDelay = time.Minute * 5
)

var total uint64 = 0

type ITaskRepo interface {
	Add(task *entity.Task)
	FindByID(id uint) (*entity.Task, bool)
	All() []*entity.Task
}

type WorkUseCase struct {
	log  *slog.Logger
	repo ITaskRepo
}

type WorkUseCaseDeps struct {
	*slog.Logger
	ITaskRepo
}

func NewWorkUseCase(deps *WorkUseCaseDeps) *WorkUseCase {
	return &WorkUseCase{
		log:  deps.Logger,
		repo: deps.ITaskRepo,
	}
}

func (uc *WorkUseCase) CreateTask(ctx context.Context, data *dto.TaskCreateRequest) (*dto.TaskCreateResponce, error) {
	op := "WorkUseCase: adding task"
	log := uc.log.With(slog.String("operation", op))
	log.Debug("CreateTask func call", "data", data)

	id := atomic.AddUint64(&total, 1)

	task := entity.NewTask(uint(id), data.Description)

	uc.repo.Add(task)

	go uc.execute(task)

	log.Info("task successfully created", "taks", task)
	return &dto.TaskCreateResponce{Task: task}, nil
}

func (uc *WorkUseCase) AllTasks(ctx context.Context) (*dto.AllTasksResponce, error) {
	op := "WorkUseCase: getting all tasks"
	log := uc.log.With(slog.String("operation", op))
	log.Debug("AllTasks func call")

	data := uc.repo.All()

	return &dto.AllTasksResponce{Tasks: data}, nil
}

func (uc *WorkUseCase) execute(task *entity.Task) {
	op := "WorkUseCase: executing task"
	log := uc.log.With(slog.String("operation", op), slog.Uint64("task_id", uint64(task.ID)))

	startTime := time.Now()

	delay := minDelay + time.Duration(rand.Intn(int(maxDelay-minDelay)))
	log.Debug("Task execution started", "delay", delay)
	time.Sleep(delay)

	foundTask, ok := uc.repo.FindByID(task.ID)
	if !ok {
		log.Error("Task not found in repository")
		return
	}

	randomResult := rand.Intn(2)

	if randomResult == 0 {
		foundTask.Complete(nil)
	} else {
		foundTask.Complete(errors.New("random error"))
	}

	foundTask.Duration = int64(time.Since(startTime).Seconds())
	log.Debug("Task execution completed", "task", task)
}
