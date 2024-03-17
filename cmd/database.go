package main

import (
	"github.com/cyclex/planet-ban/domain/model"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(driver, dsn string, debug bool) (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if debug {
		db = db.Debug()
	}

	db.AutoMigrate(
		&model.ConversationsLog{},
		&model.Messages{},
		&model.Campaign{}, &model.Kol{}, &model.Participant{},
		&model.AccessCMS{}, &model.Modul{}, &model.Privilege{}, &model.RuleModul{}, &model.Rule{}, &model.UserCMS{},
		&model.Token{},
	)

	return db, nil
}
