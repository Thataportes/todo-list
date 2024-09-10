.PHONY: up down db-logs mysql task-dir

up:
	docker-compose -f zarf/docker-compose.yaml up -d --build

down:
	docker-compose -f zarf/docker-compose.yaml down -v

db-logs:
	docker-compose -f zarf/docker-compose.yaml logs dbmake
mysql:
	docker exec -it todolist bash -c "mysql -u root -p"
task-dir:
	cd api/services/task && zsh


