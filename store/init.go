package store

import (
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
	log "phanes/collector/logger"
	"phanes/config"
	"phanes/model"
	"phanes/store/mysql"
	"phanes/store/postgres"
	"phanes/store/redis"
	"phanes/utils"
)

func Init() func() {
	var cancels = make([]func(), 0)
	var err error

	if len(config.Conf.DB) > 0 {
		for _, db := range config.Conf.DB {
			if !db.Enabled {
				continue
			}
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

	// init snowID
	if config.Conf.IdGen.Type == model.IdGenTypeSnow {
		utils.SnowGen, err = snowflake.NewNode(int64(config.Conf.IdGen.Node))
		if err != nil {
			panic(err)
		}
	}

	return func() {
		for _, cancel := range cancels {
			cancel()
		}
	}
}
