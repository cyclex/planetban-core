package domain

import (
	"context"
	"time"

	"github.com/cyclex/planet-ban/api"
	"github.com/cyclex/planet-ban/domain/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrdersUcase interface {
	GetQueueChat(c context.Context, collection string) (data []model.QueueChat, err error)
	UpdateQueueChat(c context.Context, collection string, id primitive.ObjectID) (err error)
	CreateQueueChat(c context.Context, collection string, msg api.ResSendMessage) (err error)
}

type CronChatbot struct {
	ID              primitive.ObjectID
	Err             error
	TrxChatBotMsgID string
	ServerTime      time.Time
	Status          string
}
