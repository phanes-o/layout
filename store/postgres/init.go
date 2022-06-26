package postgres

import (
	"gorm.io/gorm"
	"phanes/config"
)

var db *gorm.DB

func Init() {
	db = config.DB
}
