
MIGRATE_COMMAND=/bin/migrate.linux-amd64 -database "mysql://root:pass@tcp(mysql:3306)/todos?query"  -path db/migrations
MYSQL_COMMAND=mysql mysql -h mysql -u root -p

build:
	docker-compose build app

up:
	docker-compose up -d

run:
	docker-compose run --rm app go run .

console:
	docker-compose run --rm app bash

db-create:
	docker-compose run --rm mysql ${MYSQL_COMMAND} -e "CREATE DATABASE todos"

db-migrate:
	docker-compose run --rm app ${MIGRATE_COMMAND} up

db-migrate-down:
	docker-compose run --rm app ${MIGRATE_COMMAND} down

mysql-client:
	docker-compose exec mysql ${MYSQL_COMMAND} -D todos

fmt:
	go fmt ./...
