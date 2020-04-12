
MIGRATE_COMMAND=./migrate.linux-amd64 -database "mysql://root:pass@tcp(mysql:3306)/todos?query"  -path db/migrations

build:
	docker-compose build app

up:
	docker-compose up -d

console:
	docker-compose run --rm app bash

db-create:
	docker-compose run --rm mysql mysql -h mysql -u root -p -e "CREATE DATABASE todos"

db-migrate:
	docker-compose run --rm app ${MIGRATE_COMMAND} up

db-migrate-down:
	docker-compose run --rm app ${MIGRATE_COMMAND} down

mysql:
	docker-compose run mysql mysql -h mysql -u root -D todos -p