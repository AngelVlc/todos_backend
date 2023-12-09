
MIGRATE_COMMAND=/bin/migrate.linux-amd64 -database "mysql://root:pass@tcp(mysql:3306)/todos?query"  -path db/migrations
MYSQL_COMMAND=mysql -h mysql -u root

build:
	docker-compose build app

up:
	docker-compose up -d app mysql

run:
	docker-compose run --rm -p 5001:5001 app go run ./...

console:
	docker-compose run --rm -p 5001:5001 app sh

db-create:
	docker-compose run --rm mysql ${MYSQL_COMMAND} -e "CREATE DATABASE todos" -p

db-migrate:
	docker-compose run --rm app ${MIGRATE_COMMAND} up

db-migrate-down:
	docker-compose run --rm app ${MIGRATE_COMMAND} down 1

db-version:
	docker-compose run --rm app ${MIGRATE_COMMAND} version

mysql-client:
	docker-compose exec mysql ${MYSQL_COMMAND} -D todos -p

test:
	docker-compose run --rm app go test --race ./...

test-e2e:
	docker-compose run --rm -e BASE_URL=http://app:5001 app go test --tags=e2e ./cmd/api
