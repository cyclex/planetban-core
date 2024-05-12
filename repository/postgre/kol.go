package postgre

import (
	"fmt"
	"time"

	"github.com/cyclex/planet-ban/domain/model"
	"github.com/pkg/errors"
)

func (self *postgreRepo) FindKolBy(cond map[string]interface{}) (data []model.Kol, err error) {

	err = self.DB.Where(cond).Find(&data).Error
	if err != nil {
		err = errors.New("Duplicate kol name")
	}

	return
}

func (self *postgreRepo) SetKol(id int64, kol model.Kol) (err error) {

	kol.UpdatedAt = time.Now().Local().Unix()
	err = self.DB.Where("id = ? and status = ?", id, 1).Updates(kol).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.SetKol]")
	}

	return

}

func (self *postgreRepo) RemoveKol(id []int64) (err error) {

	err = self.DB.Where("id = ? and status = ?", id, 1).Delete(&model.Kol{}).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.RemoveKol]")
	}

	return

}

func (self *postgreRepo) CreateKol(new model.Kol) (err error) {

	new.CreatedAt = time.Now().Local().Unix()
	new.Status = 1
	err = self.DB.Create(&new).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.CreateKol]")
	}

	return
}

func (self *postgreRepo) CreateBulkKol(rows []model.Kol, skipFirstRow bool) (err error) {

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

		if skipFirstRow {
			if x == 0 {
				x++
				continue
			}
		}
		x++

		row.CreatedAt = time.Now().Local().Unix()
		row.Status = 1
		if err := tx.Create(&row).Error; err != nil {
			tx.Rollback()

			if skipFirstRow {
				err = errors.New(fmt.Sprintf("Baris ke #%d => Duplicate kol name", x))
			} else {
				err = errors.New("Duplicate kol name")
			}
			return err
		}

	}

	err = tx.Commit().Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.CreateBulkKol]")
	}

	return
}
