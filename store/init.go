package store

import (
	"go.uber.org/zap"
	log "phanes/collector/logger"
	"phanes/config"
	"phanes/store/mysql"
	"phanes/store/postgres"
	"phanes/store/redis"
)

func Init() func() {
	var cancels = make([]func(), 0)

	if len(config.Conf.DB) > 0 {
		for _, db := range config.Conf.DB {
			switch db.Type {
			case "mysql":
				cancels = append(cancels, mysql.Init(db.Addr))
			case "postgres":
				cancels = append(cancels, postgres.Init(db.Addr))
			case "redis":
				cancels = append(cancels, redis.Init(db.Addr, db.Pwd))
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
