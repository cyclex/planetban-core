package postgre

import (
	"time"

	"github.com/cyclex/planet-ban/domain/model"
	"github.com/pkg/errors"
)

func (self *postgreRepo) CreateParticipant(new model.Participant, kol model.Kol) (err error) {

	// new.CreatedAt = time.Now().Local().Unix()
	// err = self.DB.Create(&new).Error
	// if err != nil {
	// 	err = errors.Wrap(err, "[postgre.CreateParticipant]")
	// }

	// return

	tx := self.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		err = errors.Wrap(err, "[postgre.CreateParticipant]")
		return err
	}

	new.CreatedAt = time.Now().Local().Unix()
	err = tx.Create(&new).Error
	if err != nil {
		tx.Rollback()
		err = errors.Wrap(err, "[postgre.CreateParticipant]")
		return
	}

	kol.UpdatedAt = time.Now().Local().Unix()
	kol.Status = 2
	err = tx.Where("id = ?", kol.ID).Updates(kol).Error
	if err != nil {
		tx.Rollback()
		err = errors.Wrap(err, "[postgre.CreateParticipant]")
		return
	}

	err = tx.Commit().Error
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
