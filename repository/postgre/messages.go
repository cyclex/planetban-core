package postgre

import (
	"github.com/cyclex/planet-ban/domain/model"
	"github.com/pkg/errors"
)

func (self *postgreRepo) FindMsgBy(cond map[string]interface{}) (data model.Messages, err error) {

	err = self.DB.Table("messages").Where(cond).First(&data).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.FindMsgBy]")
	}

	return
}
