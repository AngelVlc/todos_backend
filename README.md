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


**Profiling**

CPU
```
	f, err := os.Create("cpuprof.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
```

MEMORY
```
	f, err := os.Create("memprof.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	runtime.GC()    // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
```