package postgre

import (
	"time"

	"github.com/cyclex/planet-ban/domain/model"
	"github.com/pkg/errors"
)

func (self *postgreRepo) CreateParticipant(new model.Participant) (err error) {

	new.CreatedAt = time.Now().Local().Unix()
	err = self.DB.Create(&new).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.CreateParticipant]")
	}

	return
}

func (self *postgreRepo) FindParticipant(cond map[string]interface{}) (data []model.Participant, err error) {

	err = self.DB.Where(cond).Find(&data).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.FindParticipant]")
	}

	return
}
