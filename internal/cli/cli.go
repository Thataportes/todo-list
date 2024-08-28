package cli

import (
	"TODO-list/internal/service"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Representa o CLI para tarefas.
type taskCLI struct {
	service *service.TaskService
}

// Cria uma nova instância de taskCLI.
func NewTaskCLI(service *service.TaskService) *taskCLI {
	return &taskCLI{service: service}
}

// Executa o CLI.
func (cli *taskCLI) Run() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tasks <command> [arguments]")
		return
	}
	command := os.Args[1]

	switch command {
	case "list":
		cli.listTask()

	case "create":
		if len(os.Args) < 4 {
			fmt.Println("Usage: tasks create <title> <description>")
			return
		}
		title := os.Args[2]
		description := os.Args[3]
		cli.createTask(title, description)

	case "update":
		if len(os.Args) < 5 {
			fmt.Println("Usage: tasks update <id> <title> <description>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID:", os.Args[2])
			return
		}
		title := os.Args[3]
		description := os.Args[4]
		cli.updateTask(id, title, description)

	case "status":
		if len(os.Args) < 4 {
			fmt.Println("Usage: tasks status <id> <status>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID:", os.Args[2])
			return
		}
		status, err := strconv.ParseBool(os.Args[3])
		if err != nil {
			fmt.Println("Invalid status value:", os.Args[3])
			return
		}
		cli.statusTask(id, status)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: tasks delete <id>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID:", os.Args[2])
			return
		}
		cli.deleteTask(id)

	case "search":
		if len(os.Args) < 3 {
			fmt.Println("Usage: tasks search <title>")
			return
		}
		title := os.Args[2]
		cli.searchTasks(title)

	case "simulate":
		if len(os.Args) < 3 {
			fmt.Println("Usage: tasks simulate <task_id> <task_id> ...")
			return
		}
		taskIDsStr := os.Args[2:]
		cli.simulateReading(taskIDsStr)

	default:
		fmt.Println("Unknown command:", command)
	}
}

func (cli *taskCLI) listTask() {
	tasks, err := cli.service.GetTasks()
	if err != nil {
		fmt.Println("Error fetching tasks:", err)
		return
	}

	fmt.Println("Tasks:")
	for _, task := range tasks {
		fmt.Printf("ID: %d, Title: %s, Description: %s, Status: %t\n",
			task.ID, task.Title, task.Description, task.Status)
	}
}

func (cli *taskCLI) createTask(title, description string) {
	task := &service.Task{
		Title:       title,
		Description: description,
		Status:      false,
	}
	if err := cli.service.CreateTask(task); err != nil {
		fmt.Println("Error creating task:", err)
		return
	}
	fmt.Printf("Task created with ID: %d\n", task.ID)
}

func (cli *taskCLI) updateTask(id int, title, description string) {
	task := &service.Task{
		ID:          id,
		Title:       title,
		Description: description,
	}
	if err := cli.service.UpdateTask(task); err != nil {
		fmt.Println("Error updating task:", err)
		return
	}
	fmt.Println("Task updated successfully.")
}

func (cli *taskCLI) statusTask(id int, status bool) {
	if err := cli.service.StatusTask(id, status); err != nil {
		fmt.Println("Error updating task status:", err)
		return
	}
	fmt.Println("Task status updated successfully.")
}

func (cli *taskCLI) deleteTask(id int) {
	if err := cli.service.DeleteTask(id); err != nil {
		fmt.Println("Error deleting task:", err)
		return
	}
	fmt.Println("Task deleted successfully.")
}

// Busca e exibe tarefas com base no nome fornecido.
func (cli *taskCLI) searchTasks(title string) {
	tasks, err := cli.service.SearchTasksByName(title)
	if err != nil {
		fmt.Println("Error searching tasks:", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	fmt.Printf("%d tasks found \n", len(tasks))
	for _, task := range tasks {
		fmt.Printf("ID: %d, Title: %s, Description: %s, Status: %t\n",
			task.ID, task.Title, task.Description, task.Status)
	}
}

// Simula a leitura de tarefas com base nos IDs fornecidos.
func (cli *taskCLI) simulateReading(taskIDsStr []string) {
	var taskIDs []int
	for _, idStr := range taskIDsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Invalid task ID:", idStr)
			continue
		}
		taskIDs = append(taskIDs, id)
	}
	// Chama o serviço para simular a leitura de múltiplas tarefas.
	responses := cli.service.SimulateMultipleReadings(taskIDs, 2*time.Second)

	// Exibe os resultados
	for _, response := range responses {
		fmt.Println(response)
	}
}
