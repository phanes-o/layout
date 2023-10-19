package store

import (
	"go.uber.org/zap"
	log "phanes/collector/logger"
	"phanes/config"
	"phanes/store/mysql"
	"phanes/store/postgres"
)

func Init() func() {
	var cancels = make([]func(), 0)

	if !config.Conf.Store.Enabled {
		return func() {}
	}

	if len(config.Conf.Store.DB) > 0 {
		for _, db := range config.Conf.Store.DB {
			if !db.Enabled {
				continue
			}
			switch db.Type {
			case "mysql":
				cancels = append(cancels, mysql.Init(true, db.Addr))
			case "postgres":
				cancels = append(cancels, postgres.Init(true, db.Addr))
			default:
				log.Error("unknown db type: ", zap.String("db_type", db.Type))
			}
		}
	}

	return func() {
		for _, cancel := range cancels {
			cancel()
		}
	}
}
