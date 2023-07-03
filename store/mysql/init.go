package mysql

import (
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	log "phanes/collector/logger"
)

var db *gorm.DB

func Init(connectAddr string) func() {
	var (
		err   error
		sqlDB *sql.DB
	)

	if db, err = gorm.Open(mysql.Open(connectAddr), &gorm.Config{}); err != nil {
		log.Fatal(err.Error())
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
