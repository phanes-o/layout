package config

import (
	"database/sql"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() {
	var (
		err error
		db  *sql.DB
	)

	if DB, err = gorm.Open(postgres.Open(Conf.Postgres), &gorm.Config{}); err != nil {
		log.Fatal()
	}
	if db, err = DB.DB(); err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(500)
	go func() {
		<-ExitC
		_ = db.Close()
	}()

	if err = db.Ping(); err != nil {
		panic(err)
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
