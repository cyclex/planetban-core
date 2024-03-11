package main

import (
	"context"
	"time"

	"github.com/cyclex/planet-ban/domain/model"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

func ConnectQueue(dsn string, c context.Context) (*mongo.Client, error) {

	timeOutRequest := viper.GetInt("context.timeout")
	_, cancel := context.WithTimeout(c, time.Duration(timeOutRequest)*time.Second)
	defer cancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dsn).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(c, opts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(c, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}
