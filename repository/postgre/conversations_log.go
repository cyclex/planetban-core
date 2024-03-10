package postgre

import (
	"time"

	"github.com/cyclex/planet-ban/domain/model"
	"github.com/pkg/errors"
)

func (self *postgreRepo) CreateConversationsLog(new model.ConversationsLog) (err error) {

	new.CreatedAt = time.Now().Local()
	err = self.DB.Create(&new).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.CreateConversationsLog]")
	}

	return
}
