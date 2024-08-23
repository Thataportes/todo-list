# TODO List

Este é um projeto de linha de comando desenvolvido em Go para gerenciar uma lista de tarefas (TODO list). O projeto utiliza MySQL como banco de dados e Docker para facilitar o desenvolvimento e a execução da aplicação.

## Objetivo

Desenvolver uma aplicação CLI para gerenciar uma lista de tarefas, permitindo criar, listar, editar, marcar como concluído e excluir tarefas. 

## Requisitos

- **Go** (versão 1.20 ou superior)
- **Docker** e **Docker Compose**

## Estrutura do Projeto

- **main.go**: Arquivo principal que contém a lógica da aplicação.
- **cli.go**: Configuração dos comandos CLI.
- **tasks.go**: Definições das estruturas e interações com o banco de dados.
- **handler.go**: Manipuladores para as operações CRUD.
- **Dockerfile**: Arquivo para construir a imagem Docker da aplicação.
- **docker-compose.yml**: Arquivo para orquestrar os contêineres Docker.
- **.env**: Arquivo de configuração com variáveis de ambiente.
- **README.md**: Documentação do projeto.

## Instalação e Execução


### Clonar o Repositório

```bash
git clone https://github.com/Thataportes/todo-list.git
cd todo-list
