package taskapp

import (
	"TODO-list/business/domain/taskbus"
	"database/sql"
	"encoding/json"
	"time"
)

// AssignedTo is a custom type for serializing the 'assigned_to' field.
type AssignedTo struct {
	ID    int  `json:"id"`
	Valid bool `json:"-"`
}

// MarshalJSON customizes the JSON encoding of the 'AssignedTo' field.
func (a AssignedTo) MarshalJSON() ([]byte, error) {
	if !a.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(struct {
		ID int `json:"id"`
	}{
		ID: a.ID,
	})
}

// CreatedBy is a custom type for the 'created_by' field.
type CreatedBy struct {
	ID    int  `json:"id"`
	Valid bool `json:"-"`
}

// MarshalJSON customizes the JSON encoding of the 'CreatedBy' field.
func (c CreatedBy) MarshalJSON() ([]byte, error) {
	if !c.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(struct {
		ID int `json:"id"`
	}{
		ID: c.ID,
	})
}

// NewTask represents a new task to be created.
type NewTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedBy   int    `json:"created_by"`
	AssignedTo  *int   `json:"assigned_to"`
}

// Decode implements the decoder interface.
func (nt *NewTask) Decode(data []byte) error {
	return json.Unmarshal(data, &nt)
}

// toBusNewTask converts a NewTask from the application layer to the business layer.
func toBusNewTask(nt NewTask) taskbus.NewTask {
	assignedTo := sql.NullInt32{Int32: 1}
	if nt.AssignedTo != nil {
		assignedTo = sql.NullInt32{Int32: int32(*nt.AssignedTo), Valid: true}
	}
	return taskbus.NewTask{
		Title:       nt.Title,
		Description: nt.Description,
		CreatedBy:   nt.CreatedBy,
		AssignedTo:  assignedTo,
	}
}

// Task represents a task in the system.
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	FinishedAt  time.Time `json:"finished_at"`
	CreatedBy   int       `json:"created_by"`
	AssignedTo  int       `json:"assigned_to"`
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
		FinishedAt:  taskBus.FinishedAt.Time,
		CreatedBy:   taskBus.CreatedBy,
		AssignedTo:  int(taskBus.AssignedTo.Int32),
	}
}

// Tasks represents a list of tasks in the application layer.
type Tasks []Task

// Encode implements the web.Encoder interface for the Tasks type.
func (ts Tasks) Encode() ([]byte, string, error) {
	data, err := json.Marshal(ts)
	return data, "application/json", err
}

// toAppTasks converts a slice of business layer tasks to application layer tasks.
func toAppTasks(tasksBus []taskbus.Task) Tasks {
	tasksApp := make(Tasks, len(tasksBus))
	for i, taskBus := range tasksBus {
		tasksApp[i] = toAppTask(taskBus)
	}
	return tasksApp
}

// UpdateTask represents a task with updates to be applied.
type UpdateTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	AssignedTo  *int   `json:"assigned_to"`
}

// Decode implements the decoder interface.
func (ut *UpdateTask) Decode(data []byte) error {
	return json.Unmarshal(data, &ut)
}

// toBusUpdateTask converts an UpdateTask from the application layer to the business layer.
func toBusUpdateTask(ut UpdateTask) taskbus.UpdateTask {
	assignedTo := sql.NullInt32{Valid: false}
	if ut.AssignedTo != nil {
		assignedTo = sql.NullInt32{Int32: int32(*ut.AssignedTo), Valid: true}
	}

	return taskbus.UpdateTask{
		Title:       ut.Title,
		Description: ut.Description,
		AssignedTo:  assignedTo,
	}
}
