package repository

import (
	"sync"

	"github.com/EvansTrein/Workmate/internal/entity"
)

type TaskRepo struct {
	*sync.Map
}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{
		Map: &sync.Map{},
	}
}

func (r *TaskRepo) Add(task *entity.Task) {
	r.Store(task.ID, task)
}

func (r *TaskRepo) FindByID(id uint) (*entity.Task, bool) {
	value, exists := r.Load(id)
	if !exists {
		return nil, false
	}
	task, ok := value.(*entity.Task)
	if !ok {
		return nil, false
	}
	return task, true
}

func (r *TaskRepo) All() []*entity.Task {
	var result []*entity.Task

	r.Range(func(key, value any) bool {
		task := value.(*entity.Task)
		result = append(result, task)
		return true
	})

	return result
}
