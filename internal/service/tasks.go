package service

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Estrutura de dados
type Task struct {
	ID          int
	Title       string
	Description string
	status      bool
}

type TaskService struct {
	db *sql.DB
}

func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{db: db}
}

// metodo que cria uma tarefa o bd
func (s *TaskService) CreateTask(task *Task) error {
	query := "INSERT INTO tasks (title, description) VALUES (?, ?)"
	result, err := s.db.Exec(query, task.Title, task.Description)
	if err != nil {
		return err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	task.ID = int(lastInsertID)
	return nil
}

func (s *TaskService) GetTasks() ([]Task, error) {
	query := "SELECT id , title, description, status FROM tasks"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil

}

func (s *TaskService) GetTaskByID(id int) (*Task, error) {
	query := "SELECT id, title, description, completed FROM tasks WHERE id = ?"
	row := s.db.QueryRow(query, id)

	var task Task
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.status)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *TaskService) UptadeTask(task *Task) error {
	query := "UPDATE tasks SET title=?, description=?, WHERE id=? "
	_, err := s.db.Exec(query, task.Title, task.ID)
	return err
}

func (s *TaskService) StatusTask(id int) error {
	query := "UPDATE tasks SET status = TRUE WHERE id=?"
	_, err := s.db.Exec(query, id)
	return err
}

func (s *TaskService) DeleteTask(id int) error {
	query := "DELETE FROM tasks WHERE id=?"
	_, err := s.db.Exec(query, id)
	return err
}

func (s *TaskService) SearchTasksByName(name string) ([]Task, error) {
	query := "SELECT id, title, description, status FROM tasks WHERE title like ?"
	rows, err := s.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *TaskService) SimulateReading(taskId int, duration time.Duration, results chan<- string) {
	task, err := s.GetTaskByID(taskId)
	if err != nil || task == nil {
		results <- fmt.Sprintf("Task %d not found", taskId)
		return
	}

	time.Sleep(duration)
	results <- fmt.Sprintf("Task %s read", task.Title)
}

func (s *TaskService) SimulateMultipleReadings(taskIDd []int, duration time.Duration) []string {
	results := make(chan string, len(taskIDd))

	for _, id := range taskIDd {
		go func(taskId int) {
			s.SimulateReading(taskId, duration, results)
		}(id)
	}

	var responses []string
	for range taskIDd {
		responses = append(responses, <-results) // pause
	}
	close(results)
	return responses
}
