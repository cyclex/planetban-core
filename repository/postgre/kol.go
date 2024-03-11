package postgre

import (
	"time"

	"github.com/cyclex/planet-ban/domain/model"
	"github.com/pkg/errors"
)

func (self *postgreRepo) FindKolBy(cond map[string]interface{}) (data []model.Kol, err error) {

	err = self.DB.Where(cond).Find(&data).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.FindKolBy]")
	}

	return
}

func (self *postgreRepo) SetKol(id int64, kol model.Kol) (err error) {

	err = self.DB.Table("kol").Where("id = ?", id).Updates(kol).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.SetKol]")
	}

	return

}

func (self *postgreRepo) RemoveKol(id []int64) (err error) {

	err = self.DB.Delete(&model.Kol{}, id).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.RemoveKol]")
	}

	return

}

func (self *postgreRepo) CreateKol(new model.Kol) (err error) {

	new.CreatedAt = time.Now().Local()
	err = self.DB.Create(&new).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.CreateKol]")
	}

	return
}

func (self *postgreRepo) CreateBulkKol(rows []model.Kol) (err error) {

	tx := self.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		err = errors.Wrap(err, "[postgre.CreateBulkKol]")
		return err
	}

	x := 0
	for _, row := range rows {

		if x == 0 {
			x++
			continue
		}
		x++

		if err := tx.Create(&row).Error; err != nil {
			tx.Rollback()
			err = errors.Wrapf(err, "[postgre.CreateBulkKol] Baris ke #%d => ", x)
			return err
		}

	}

	err = tx.Commit().Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.CreateBulkKol]")
	}

	return
}
