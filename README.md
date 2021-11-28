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

**Load test**

hey -m POST -d '{"username": "admin","password": "893210"}'  http://localhost:5001/auth/login


**Profiling**

CPU
```
	f, err := os.Create("cpuprof.pprof")
	if err != nil {
		log.Fatal(err)
	}
	runtime.SetCPUProfileRate(500)
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
```

go tool pprof cpuprof.pprof

MEMORY
```
	f, err := os.Create("memprof.pprof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	runtime.GC()    // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
```

```
	import "net/http/pprof"

	pprofSubRouter := router.PathPrefix("/debug/pprof").Subrouter()
	pprofSubRouter.Handle("/cmdline", http.HandlerFunc(pprof.Cmdline))
	pprofSubRouter.Handle("/profile", http.HandlerFunc(pprof.Profile))
	pprofSubRouter.Handle("/symbol", http.HandlerFunc(pprof.Symbol))
	pprofSubRouter.Handle("/heap", pprof.Handler("heap"))
	pprofSubRouter.Handle("/block", pprof.Handler("block"))
	pprofSubRouter.Handle("/goroutine", pprof.Handler("goroutine"))
	pprofSubRouter.Handle("/threadcreate", pprof.Handler("threadcreate"))
```

curl http://localhost:5001/debug/pprof/heap > heap.1.pprof

go tool pprof -http=:8080 -inuse_objects -base heap.0.pprof heap.1.pprof

```
profileFunc := profile.Start(profile.MemProfile, profile.MemProfileRate(1), profile.ProfilePath("."), profile.NoShutdownHook)


profileFunc.Stop()
log.Println("gracefully stopped")
```