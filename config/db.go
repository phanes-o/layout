package config

import (
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	log "phanes/collector/logger"
)

func initDB() func() {
	var (
		err error
		db  *sql.DB
	)

	if DB, err = gorm.Open(postgres.Open(Conf.Postgres), &gorm.Config{}); err != nil {
		log.Fatal(err)
	}
	if db, err = DB.DB(); err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(500)

	if err = db.Ping(); err != nil {
		panic(err)
	}

	return func() {
		if err = db.Close(); err != nil {
			log.Error(err)
		}
	}

}

func AutoMigrate(models ...interface{}) {
	if len(models) > 0 {
		for _, model := range models {
			if !DB.Migrator().HasTable(model) {
				if err := DB.Migrator().CreateTable(model); err != nil {
					panic(err)
				}
			}
		}
	}
}
