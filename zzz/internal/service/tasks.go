package service

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Representa uma tarefa no sistema.
type Task struct {
	ID          int
	Title       string
	Description string
	Status      bool
}

// Lida com a lógica de negócios e persistência de tarefas.
type TaskService struct {
	db *sql.DB
}

// Cria uma nova instância de TaskService.
func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{db: db}
}

// Cria uma nova tarefa no banco de dados.
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

// Retorna todas as tarefas do banco de dados.
func (s *TaskService) GetTasks() ([]Task, error) {
	query := "SELECT * FROM tasks"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil

}

// Retorna uma tarefa pelo seu ID.
func (s *TaskService) GetTaskByID(id int) (*Task, error) {
	query := "SELECT id, title, description, status FROM tasks WHERE id = ?"
	row := s.db.QueryRow(query, id)

	var task Task
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Status)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// Atualiza as informações de uma tarefa no banco de dados.
func (s *TaskService) UpdateTask(task *Task) error {
	query := "UPDATE tasks SET title=?, description=?, status=? WHERE id=? "
	_, err := s.db.Exec(query, task.Title, task.ID)
	return err
}

// Mostra se uma tarefa foi concluida ou nao no banco de dados
func (s *TaskService) StatusTask(id int, status bool) error {
	query := "UPDATE tasks SET status=? WHERE id=? "
	_, err := s.db.Exec(query, status, id)
	return err
}

// Deleta uma tarefa do banco de dados.
func (s *TaskService) DeleteTask(id int) error {
	query := "DELETE FROM tasks WHERE id=?"
	_, err := s.db.Exec(query, id)
	return err
}

// Busca tarefas pelo nome (título) no banco de dados.
func (s *TaskService) SearchTasksByName(name string) ([]Task, error) {
	query := "SELECT id, title, description, status FROM tasks WHERE title LIKE ?"
	rows, err := s.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// Simula a leitura de uma tarefa com base em um tempo de leitura.
func (s *TaskService) SimulateReading(taskId int, duration time.Duration, results chan<- string) {
	task, err := s.GetTaskByID(taskId)
	if err != nil || task == nil {
		results <- fmt.Sprintf("Task %d not found", taskId)
		return
	}

	time.Sleep(duration)
	results <- fmt.Sprintf("Task %s read", task.Title)
}

// Simula a leitura de múltiplas tarefas simultaneamente.
func (s *TaskService) SimulateMultipleReadings(taskIDs []int, duration time.Duration) []string {
	results := make(chan string, len(taskIDs))

	for _, id := range taskIDs {
		go func(taskId int) {
			s.SimulateReading(taskId, duration, results)
		}(id)
	}

	var responses []string
	for range taskIDs {
		responses = append(responses, <-results) // pause
	}
	close(results)
	return responses
}
