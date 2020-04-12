# TODOS BACKEND

## Migrations

**create a migration**
```
migrate create -ext sql -dir db/migrations -seq migration_name
```

**force database version**
```
migrate -database "mysql://user:pass@tcp(host:port)/dbname?query"  -path db/migrations force version
```
