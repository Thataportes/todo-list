package main

import (
	"TODO-list/internal/cli"
	"TODO-list/internal/service"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Conectar com o MySQL
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/todolist")
	if err != nil {
		log.Fatal(err.Error())
	}

	// Verificar a conexão
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Conectou com o banco de dados")

	// Inicializando o serviço
	taskService := service.NewTaskService(db)

	// Verificar se há argumentos CLI
	if len(os.Args) > 1 {
		taskCLI := cli.NewTaskCLI(taskService)
		taskCLI.Run()
		return
	}
	// Caso nao haja argumentos CLI
	fmt.Println("Usage: tasks <command> [arguments]")

}
