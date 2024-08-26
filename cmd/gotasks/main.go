package main

import (
	"TODO-list/internal/cli"
	"TODO-list/internal/service"
	"TODO-list/internal/web"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// conectar com o Mysql
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/todolist")
	if err != nil {
		log.Fatal(err.Error())
	}

	// verificar a conexao
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Conectou com o banco de dados")

	// Inicializando o serviço
	taskService := service.NewTaskService(db)
	// Inicializando os handlers
	taskHandlers := web.NewTaskHandlers(taskService)
	// Verifica se o CLI foi chamado com o comando "search" ou "simulate"
	if len(os.Args) > 1 && (os.Args[1] == "search" || os.Args[1] == "simulate") {
		taskCLI := cli.NewTaskCLI(taskService)
		taskCLI.Run()
		return
	}

	// Criando o roteador com o novo servidor
	router := http.NewServeMux()

	// Configurando as rotas RESTful
	router.HandleFunc("GET /tasks", taskHandlers.GetTasks)
	router.HandleFunc("POST /tasks", taskHandlers.CreateTask)
	router.HandleFunc("GET /tasks/{id}", taskHandlers.GetTaskByID)
	router.HandleFunc("PUT /tasks/{id}", taskHandlers.UptadeTask)
	router.HandleFunc("PATCH /tasks/{id}", taskHandlers.StatusTask)
	router.HandleFunc("DELETE /tasks/{id}", taskHandlers.DeleteTask)

	// Iniciando o servidor
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", router)
}
