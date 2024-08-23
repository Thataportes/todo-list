package cli

import (
	"TODO-list/internal/service"
	"fmt"
	"os"
	"strconv"
	"time"
)

type taskCLI struct {
	service *service.TaskService
}

func NewTaskCLI(service *service.TaskService) *taskCLI {
	return &taskCLI{service: service}
}

func (cli *taskCLI) Run() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tasks <command> [arguments]")
		return
	}
	command := os.Args[1]

	switch command {
	case "search":
		if len(os.Args) < 3 {
			fmt.Println("Usage: tasks search <task title>")
			return
		}
		taskName := os.Args[2]
		cli.searchTasks(taskName)

	case "simulate":
		if len(os.Args) < 3 {
			fmt.Println("Usage: tasks simulate <task_id> <task_id> <task_id> ...")
			return
		}
		TaskIDs := os.Args[2:]
		cli.simulateReading(TaskIDs)
	}

}

func (cli *taskCLI) searchTasks(name string) {
	tasks, err := cli.service.SearchTasksByName(name)
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
		fmt.Printf("ID: %d, Title: %s, Description: %s",
			task.ID, task.Title, task.Description,
		)
	}
}

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
	responses := cli.service.SimulateMultipleReadings(taskIDs, 2*time.Second)

	for _, response := range responses {
		fmt.Println(response)
	}
}
