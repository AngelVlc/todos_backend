package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	authDomain "github.com/AngelVlc/todos_backend/src/internal/api/auth/domain"
	sharedApp "github.com/AngelVlc/todos_backend/src/internal/api/shared/application"
	"github.com/honeybadger-io/honeybadger-go"
	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initHoneyBadger(cfg sharedApp.ConfigurationService) {
	if !cfg.InProduction() {
		log.Println("Init HoneyBadger skipped")
		return
	}

	configuration := honeybadger.Configuration{
		APIKey: cfg.GetHoneyBadgerApiKey(),
		Env:    cfg.GetEnvironment(),
		Sync:   true,
	}
	honeybadger.Configure(configuration)
}

func initNewRelic(cfg sharedApp.ConfigurationService) (*newrelic.Application, error) {
	if !cfg.InProduction() {
		log.Println("Init NewRelic skipped")
		return nil, nil
	}

	appName := fmt.Sprintf("todos_backend_%v", cfg.GetEnvironment())
	licenseKey := cfg.GetNewRelicLicenseKey()

	newRelicApp, err := newrelic.NewApplication(newrelic.ConfigAppName(appName), newrelic.ConfigLicense(licenseKey), newrelic.ConfigEnabled(true))
	if err != nil {
		return nil, err
	}

	newRelicApp.WaitForConnection(5 * time.Second)

	return newRelicApp, nil
}

func initDb(c sharedApp.ConfigurationService) (*gorm.DB, error) {
	sqlDb, err := sql.Open("nrmysql", c.GetDatasource())
	if err != nil {
		return nil, err
	}

	sqlDb.SetConnMaxLifetime(60 * time.Second)

	gormdb, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDb}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return gormdb, nil
}

func initDeleteExpiredTokensProcess(cfg sharedApp.ConfigurationService, authRepo authDomain.AuthRepository, newRelicApp *newrelic.Application) {
	duration := cfg.GetDeleteExpiredRefreshTokensIntervalDuration()
	ticker := time.NewTicker(duration)
	done := make(chan bool)
	log.Printf("Delete expired token process set every %v", duration)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				txn := newRelicApp.StartTransaction("deleteExpiredRefreshTokens")
				ctx := newrelic.NewContext(context.Background(), txn)
				authRepo.DeleteExpiredRefreshTokens(ctx, t)
				txn.End()
			}
		}
	}()
}
