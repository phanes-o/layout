package postgres

import (
	"context"
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	log "phanes/collector/logger"
)

var db *gorm.DB

type contextTxKey struct{}

var ContextTxKey = contextTxKey{}

func dbWithContext(ctx context.Context) *gorm.DB {
	return GetDB(ctx)
}

func GetDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(ContextTxKey).(*gorm.DB)
	if ok {
		return tx
	}
	return db.WithContext(ctx)
}

func Init(connectAddr string) func() {
	var (
		err   error
		sqlDB *sql.DB
	)

	if db, err = gorm.Open(postgres.Open(connectAddr), &gorm.Config{}); err != nil {
		log.Error(err.Error())
	}
	if sqlDB, err = db.DB(); err != nil {
		log.Fatal(err.Error())
	}
	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetMaxOpenConns(500)

	if err = sqlDB.Ping(); err != nil {
		panic(err)
	}

	for _, s := range migrates {
		s.Init()
	}
	return func() {
		if err = sqlDB.Close(); err != nil {
			log.Error(err.Error())
		}
	}
}
