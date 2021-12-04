
MIGRATE_COMMAND=/bin/migrate.linux-amd64 -database "mysql://root:pass@tcp(mysql:3306)/todos?query"  -path db/migrations
MYSQL_COMMAND=mysql mysql -h mysql -u root -p

build:
	docker-compose build app

up:
	docker-compose up -d app mysql

run:
	docker-compose run --rm -p 5001:5001 app go run cmd/api/main.go

console:
	docker-compose run --rm -p 5001:5001 app bash

db-create:
	docker-compose run --rm mysql ${MYSQL_COMMAND} -e "CREATE DATABASE todos"

db-migrate:
	docker-compose run --rm app ${MIGRATE_COMMAND} up

db-migrate-down:
	docker-compose run --rm app ${MIGRATE_COMMAND} down 1

db-version:
	docker-compose run --rm app ${MIGRATE_COMMAND} version

mysql-client:
	docker-compose exec mysql ${MYSQL_COMMAND}

test:
	docker-compose run --rm app go test --race ./...

test-e2e:
	docker-compose -f docker-compose-e2e.yml run --rm -e BASE_URL=http://app:5001 app-e2e
