package model

import (
	"time"

	"github.com/cyclex/planet-ban/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QueueChat struct {
	ID        primitive.ObjectID `bson:"_id"`
	TrxId     string             `bson:"trx_id"`
	Messages  api.ResSendMessage `bson:"messages"`
	State     int                `bson:"state"`
	ExpiredAt time.Time          `bson:"expired_at"`
}
