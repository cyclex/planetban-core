package repository

import (
	"github.com/cyclex/planet-ban/domain/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QueueRepository interface {
	GetQueue(collection string) (queue []model.QueueChat, err error)
	UpdateQueue(collection string, id primitive.ObjectID) (err error)
	CreateQueue(collection string, data model.QueueChat) (err error)
}
