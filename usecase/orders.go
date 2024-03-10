package usecase

import (
	"context"
	"time"

	"github.com/cyclex/planet-ban/api"
	"github.com/cyclex/planet-ban/domain"
	"github.com/cyclex/planet-ban/domain/model"
	"github.com/cyclex/planet-ban/domain/repository"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ordersUcase struct {
	q              repository.QueueRepository
	contextTimeout time.Duration
}

func NewOrdersUcase(q repository.QueueRepository, timeout time.Duration) domain.OrdersUcase {

	return &ordersUcase{
		q:              q,
		contextTimeout: timeout,
	}
}

func (self *ordersUcase) GetQueueChat(c context.Context, collection string) (data []model.QueueChat, err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	data, err = self.q.GetQueue(collection)
	if err != nil {
		err = errors.Wrap(err, "[usecase.GetQueueChat]")
	}
	return
}

func (self *ordersUcase) UpdateQueueChat(c context.Context, collection string, id primitive.ObjectID) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	err = self.q.UpdateQueue(collection, id)
	if err != nil {
		err = errors.Wrap(err, "[usecase.UpdateQueueChat]")
	}
	return

}

func (self *ordersUcase) CreateQueueChat(c context.Context, collection string, msg api.ResSendMessage) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	id := primitive.NewObjectID()

	data := model.QueueChat{
		ID:       id,
		TrxId:    uuid.New().String(),
		State:    1,
		Messages: msg,
	}

	err = self.q.CreateQueue(collection, data)
	if err != nil {
		err = errors.Wrap(err, "[usecase.CreateQueueChat]")
	}
	return

}
