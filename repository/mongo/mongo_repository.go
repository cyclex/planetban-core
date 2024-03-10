package mongo

import (
	"context"
	"time"

	"github.com/cyclex/planet-ban/domain/model"
	"github.com/cyclex/planet-ban/domain/repository"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepo struct {
	DB            *mongo.Database
	c             context.Context
	channel       string
	expiredInHour time.Duration
}

func NewmongoRepository(c context.Context, db *mongo.Database, channel string, expired time.Duration) repository.QueueRepository {
	return &mongoRepo{
		DB:            db,
		c:             c,
		channel:       channel,
		expiredInHour: expired,
	}
}

func (self *mongoRepo) GetQueue(collection string) (queue []model.QueueChat, err error) {

	condition := QueueBasedOnState(1)

	dt, err := self.DB.Collection(collection).Find(
		self.c,
		condition,
	)

	if err != nil {
		err = errors.Wrap(err, "[mongo.GetQueue] Find")
		return nil, err
	}

	for dt.Next(self.c) {
		var queueHolder model.QueueChat

		err := dt.Decode(&queueHolder)
		if err != nil {
			err = errors.Wrap(err, "[mongo.GetQueue] Decode")
			return nil, err
		}

		queue = append(queue, queueHolder)
	}

	defer dt.Close(self.c)

	return queue, nil

}

func (self *mongoRepo) UpdateQueue(collection string, id primitive.ObjectID) (err error) {

	id, err = primitive.ObjectIDFromHex(id.Hex())
	if err != nil {
		return err
	}

	col := self.DB.Collection(collection)

	updatedData := AckQueue(1 * time.Hour)

	_, err = col.UpdateOne(
		self.c,
		bson.M{"_id": id},
		bson.D{
			{"$set", updatedData},
		})

	if err != nil {
		err = errors.Wrap(err, "[mongo.UpdateQueue] UpdateOne")
	}

	return

}

func AckQueue(exp time.Duration) bson.D {

	data := bson.D{
		{"expired_at", time.Now().Local().Add(exp)},
		{"state", 2},
	}

	return data
}

func QueueBasedOnState(stateID int) bson.D {

	data := bson.D{
		{"state", stateID},
	}

	return data
}

func (self *mongoRepo) CreateQueue(collection string, data model.QueueChat) (err error) {

	data.ExpiredAt = time.Now().Local().Add(6 * time.Hour)

	_, err = self.DB.Collection(collection).InsertOne(self.c, data)
	if err != nil {
		err = errors.Wrap(err, "[mongo.CreateQueue] InsertOne")
	}

	return

}
