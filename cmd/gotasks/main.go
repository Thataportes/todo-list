package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Estrutura para Todo
type Todo struct {
	Id          int
	Title       string
	Description string
	Status      string
}

func main() {

	// conectar com o Mysql
	db, err := sql.Open("mysql", "root:1053@/mysql")
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
}
