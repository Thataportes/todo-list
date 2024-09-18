package taskapp

import (
	"TODO-list/business/domain/taskbus"
	"encoding/json"
	"time"
)

// NewTask represents a new task to be created.
type NewTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Decode implements the decoder interface.
func (nt *NewTask) Decode(data []byte) error {
	return json.Unmarshal(data, &nt)
}

// toBusNewTask converts a NewTask from the application layer to the business layer.
func toBusNewTask(nt NewTask) taskbus.NewTask {
	return taskbus.NewTask{
		Title:       nt.Title,
		Description: nt.Description,
	}
}

// Task represents a task in the system.
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	FinishedAt  time.Time `json:"finished_at"`
}

// Encode implements the web.Encoder interface for the Task type.
func (t Task) Encode() ([]byte, string, error) {
	data, err := json.Marshal(t)
	return data, "application/json", err
}

// toAppTask converts a task from the business layer to the application layer.
func toAppTask(taskBus taskbus.Task) Task {
	return Task{
		ID:          taskBus.ID,
		Title:       taskBus.Title,
		Description: taskBus.Description,
		CreatedAt:   taskBus.CreatedAt,
	}
}

// Tasks represents a list of tasks in the application layer.
type Tasks []Task

// Encode implements the web.Encoder interface for the ToAppTasks type.
func (ts Tasks) Encode() ([]byte, string, error) {
	data, err := json.Marshal(ts)
	return data, "application/json", err
}

// toAppTasks converts a slice of business layer tasks to application layer tasks.
func toAppTasks(tasksBus []taskbus.Task) Tasks {
	tasksApp := make([]Task, len(tasksBus))
	for i, taskBus := range tasksBus {
		tasksApp[i] = toAppTask(taskBus)
	}
	return tasksApp
}

// UpdateTask represents a task with updates to be applied.
type UpdateTask struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Decode implements the decoder interface.
func (ut *UpdateTask) Decode(data []byte) error {
	return json.Unmarshal(data, &ut)
}

// toBusUpdateTask converts an UpdateTask from the application layer to the business layer.
func toBusUpdateTask(ut UpdateTask) taskbus.UpdateTask {
	return taskbus.UpdateTask{
		ID:          ut.ID,
		Title:       ut.Title,
		Description: ut.Description,
	}
}
