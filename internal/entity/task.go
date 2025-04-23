package entity

type TaskStatus string

const (
	// statusPending    TaskStatus = "pending"
	statusInProgress TaskStatus = "in progress"
	statusCompleted  TaskStatus = "completed"
)

type Task struct {
	Status      TaskStatus `json:"status"`
	Result      string     `json:"result,omitempty"`
	Description string     `json:"description"`
	ID          uint       `json:"id"`
	Duration    int64      `json:"duration_in_seconds,omitempty"`
}

func NewTask(id uint, description string) *Task {
	return &Task{
		ID:          id,
		Status:      statusInProgress,
		Description: description,
	}
}

func (t *Task) Complete(err error) {
	t.Status = statusCompleted

	if err != nil {
		t.Result = err.Error()
	} else {
		t.Result = "successfully"
	}
}
