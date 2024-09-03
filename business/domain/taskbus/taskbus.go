package taskbus

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// Business handles business logic and persistence of tasks.
type Business struct {
	db *sql.DB
}

// NewBusiness creates a new instance of Business.
func NewBusiness(db *sql.DB) *Business {
	return &Business{db: db}
}

// Create adds a new task to the database and returns the created task.
// TODO: precisamos setar o campo createdAt.
func (s *Business) Create(ctx context.Context, nt NewTask) (Task, error) {
	query := "INSERT INTO tasks (title, description, createdAt) VALUES (?, ?)"
	result, err := s.db.ExecContext(ctx, query, nt.Title, nt.Description, nt.CreatedAt)
	if err != nil {
		return Task{}, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return Task{}, err
	}

	t := Task{
		ID:          int(lastInsertID),
		Title:       nt.Title,
		Description: nt.Description,
	}

	return t, nil
}

func (s *Business) Query() ([]Task, error) {
	query := "SELECT id, title, description FROM tasks"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil

}

func (s *Business) QueryByID(id int) (Task, error) {
	query := "SELECT id, title, description FROM tasks WHERE id =?"
	row := s.db.QueryRow(query, id)

	var task Task
	err := row.Scan(&task.ID, &task.Title, &task.Description)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

// Update modifies task information in the database and returns the updated task.ções de uma tarefa no banco de dados e retorna a tarefa atualizada.
func (s *Business) Update(ctx context.Context, ut UpdateTask) (Task, error) {
	query := "UPDATE tasks SET title=?, description=? WHERE id= ? "
	_, err := s.db.ExecContext(ctx, query, ut.Title, ut.Description, ut.ID)
	if err != nil {
		return Task{}, err
	}

	updatedTask, err := s.QueryByID(ut.ID)
	if err != nil {
		return Task{}, err
	}

	return updatedTask, nil
}

func (s *Business) Delete(id int) error {
	query := "DELETE FROM tasks WHERE id=?"
	_, err := s.db.Exec(query, id)
	return err
}

func (s *Business) Finish(id int) error {
	query := "UPDATE tasks SET status = completed' WHERE id = ?"
	_, err := s.db.Exec(query, id)
	return err
}
