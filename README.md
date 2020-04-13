# TODOS BACKEND

## Migrations

**create a migration**
```
migrate create -ext sql -dir db/migrations -seq migration_name
```

**force database version**
```
migrate -database "mysql://root:pass@tcp(localhost)/todos?query"  -path db/migrations force version
```
