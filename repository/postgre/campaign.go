package postgre

import (
	"time"

	"github.com/cyclex/planet-ban/domain/model"
	"github.com/pkg/errors"
)

func (self *postgreRepo) FindCampaignBy(cond map[string]interface{}) (data []model.Campaign, err error) {

	err = self.DB.Table("campaign").Where(cond).Find(&data).Error

	return
}

func (self *postgreRepo) SetCampaign(id int64, campaign model.Campaign) (err error) {

	err = self.DB.Table("campaign").Where("id = ?", id).Updates(campaign).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.SetCampaign] Updates")
	}
	return

}

func (self *postgreRepo) RemoveCampaign(id []int64) (err error) {

	err = self.DB.Delete(&model.Campaign{}, id).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.RemoveCampaign] Delete")
	}
	return

}

func (self *postgreRepo) CreateCampaign(new model.Campaign, data []model.Kol) (err error) {

	tx := self.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		err = errors.Wrap(err, "[postgre.CreateCampaign]")
		return err
	}

	new.CreatedAt = time.Now().Local()
	err = tx.Create(&new).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.CreateCampaign] Create_1")
		tx.Rollback()
		return
	}

	for _, v := range data {
		v.CampaignID = int64(new.ID)
		err = self.DB.Create(v).Error
		if err != nil {
			err = errors.Wrap(err, "[postgre.CreateCampaign] Create_2")
			tx.Rollback()
			return
		}
	}

	err = tx.Commit().Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.CreateCampaign] Commit")
	}
	return

}
