package postgres

import (
	"context"
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	log "phanes/collector/logger"
	"phanes/errors"
)

type contextTxKey struct{}

var (
	db              *gorm.DB
	enabled         = false
	NotEnabledError = errors.New("postgres not enabled")
	ContextTxKey    = contextTxKey{}
)

func GetDB(ctx context.Context) (*gorm.DB, error) {
	if !enabled {
		return nil, NotEnabledError
	}
	tx, ok := ctx.Value(ContextTxKey).(*gorm.DB)
	if ok {
		return tx, nil
	}
	return db.WithContext(ctx), nil
}

func Init(enabled bool, connectAddr string) func() {
	if !enabled {
		return func() {}
	}
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
	enabled = true
	return func() {
		if err = sqlDB.Close(); err != nil {
			log.Error(err.Error())
		}
	}
}
