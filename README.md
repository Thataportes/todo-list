# TODO List

Este é um projeto de linha de comando desenvolvido em Go para gerenciar uma lista de tarefas (TODO list). O projeto utiliza MySQL como banco de dados e Docker Compose para facilitar o desenvolvimento e a execução da aplicação.

## Objetivo

Desenvolver uma aplicação CLI para gerenciar uma lista de tarefas, permitindo criar, listar, editar, marcar como concluído e excluir tarefas. 

## Requisitos

** Versionamento:

O código deve ser versionado utilizando o Git.
O repositório está disponível no GitHub para clonagem.

** Banco de Dados:

Utiliza MySQL para armazenar as informações das tarefas.
Configurações de acesso ao banco de dados são feitas através de um arquivo .env.

- **Go** (versão 1.20 ou superior)
- **Docker** e **Docker Compose**

### Clonar o Repositório

```bash
git clone https://github.com/Thataportes/todo-list.git
cd todo-list


###  Estrutura do Projeto

- **main.go**: Arquivo principal que contém a lógica da aplicação.
- **cli.go**: Configuração dos comandos CLI.
- **tasks.go**: Definições das estruturas e interações com o banco de dados.
- **handler.go**: Manipuladores para as operações CRUD.
- **Dockerfile**: Arquivo para construir a imagem Docker da aplicação.
- **docker-compose.yml**: Arquivo para orquestrar os contêineres Docker.
- **.env**: Arquivo de configuração com variáveis de ambiente.
- **README.md**: Documentação do projeto.

## Instalação e Execução

Configuração do Ambiente

** Configure e inicie o container docker.
- ** docker-compose up **

** Compile a aplicação.
- ** go build -o task ./cmd/gotasks.main.go **

** Entre no banco de dados Mysql.
- ** Docker exec -it todolist bash 
- ** mysql -u root -p
- ** root **

### Comandos da Aplicação 
** Entrar no banco de dados "todolist" e nas tabelas
- ** USE todolist **
- ** SHOW TABLES; **

** Listar todas as tarefas.
- ** SELECT * FROM task; **

** Criar uma tarefa
- ** INSERT INTO task (title, description, status)
- ** VALUES ("Compras", "Ir as compras tal dia", "pending") **

** Editar um item 
- ** UPDATE task
- ** SET title="Compras", description="ir no sabado", status="completed"
- ** id=1; **

** Marcar tarefa como concluida
- ** UPDATE task
- ** SET status="completed"
- ** id=1; **

** Deletar uma tarefa
- ** DELETE FROM task WHERE id= 2; **





 