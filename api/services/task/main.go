package main

import (
	"TODO-list/app/domain/taskapp"
	"TODO-list/business/domain/taskbus"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {

	// Open connection to the MySQL database.
	// The DSN includes 'parseTime=true' to ensure MySQL DATETIME and TIMESTAMP fields are automatically.
	// parsed into Go's time.Time type.
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/todolist?parseTime=true")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	// Initialize the business layer with the database connection.
	taskBus := taskbus.NewBusiness(db)

	// Initialize the application layer with the business layer.
	app := taskapp.NewApp(taskBus)

	// Create a context for the operations.
	ctx := context.Background()

	// Check if a command is provided.
	if len(os.Args) < 2 {
		fmt.Println("Usage: task <command> [arguments]")
		return
	}
	command := os.Args[1]

	// Switch on the command provided.
	switch command {
	case "query":
		query(ctx, app)
	case "create":
		if len(os.Args) < 4 {
			fmt.Println("Usage: task create <title> <description>")
			return
		}
		title := os.Args[2]
		description := os.Args[3]
		create(ctx, app, title, description)
	case "update":
		if len(os.Args) < 5 {
			fmt.Println("Usage: task update <id> <title> <description>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID:", os.Args[2])
			return
		}
		title := os.Args[3]
		description := os.Args[4]
		update(ctx, app, id, title, description)
	case "finish":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task finish <id>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID:", os.Args[2])
			return
		}
		finish(ctx, app, id)
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task delete <id>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID:", os.Args[2])
			return
		}
		delete(ctx, app, id)
	default:
		fmt.Println("Unknown command:", command)
	}

}

// Query retrieves and displays all tasks from the application layer.
func query(ctx context.Context, app *taskapp.App) {
	tasks, err := app.Query(ctx)
	if err != nil {
		fmt.Println("Error fetching tasks:", err)
		return
	}

	for _, task := range tasks {
		fmt.Printf("%d: %s - %s (Created: %s, Finished: %v)\n",
			task.ID, task.Title, task.Description, task.CreatedAt, task.FinishedAt)
	}
}

// create adds a new task using the application layer.
func create(ctx context.Context, app *taskapp.App, title, description string) {
	_, err := app.Create(ctx, taskapp.NewTask{Title: title, Description: description})
	if err != nil {
		fmt.Println("Error creating task:", err)
		return
	}
	fmt.Println("Task created successfully.")
}

// update modifies an existing task using the application layer.
func update(ctx context.Context, app *taskapp.App, id int, title, description string) {
	_, err := app.Update(ctx, taskapp.UpdateTask{ID: id, Title: title, Description: description})
	if err != nil {
		fmt.Println("Error updated task:", err)
		return
	}
	fmt.Println("Task updated successfully.")
}

// finish marks a task as finished using the application layer.
func finish(ctx context.Context, app *taskapp.App, id int) {
	err := app.Finish(ctx, id)
	if err != nil {
		fmt.Println("Error finishing task:", err)
		return
	}
	fmt.Println("Task finished successfully.")
}

// delete removes a task using the application layer.
func delete(ctx context.Context, app *taskapp.App, id int) {
	err := app.Delete(ctx, id)
	if err != nil {
		fmt.Println("Error deleting task:", err)
		return
	}
	fmt.Println("Task deleted successfully.")
}
