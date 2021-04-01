# TODOS BACKEND

## Migrations

https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

**create a migration**
```
migrate create -ext sql -dir db/migrations -seq migration_name
```

**run migrations**
```
migrate -database "mysql://root:pass@tcp(localhost)/todos" -path src/db/migrations up 
```


**force database version**
```
migrate -database "mysql://root:pass@tcp(localhost)/todos?query"  -path src/db/migrations force version
```
