package repository

import (
	"github.com/cyclex/planet-ban/api"
	"github.com/cyclex/planet-ban/domain/model"
)

type ModelRepository interface {
	CreateConversationsLog(new model.ConversationsLog) (err error)

	FindMsgBy(cond map[string]interface{}) (data model.Messages, err error)

	FindKolBy(cond map[string]interface{}) (data []model.Kol, err error)
	SetKol(id int64, kol model.Kol) (err error)
	RemoveKol(id []int64) (err error)
	CreateBulkKol(new []model.Kol) (err error)

	FindCampaignBy(cond map[string]interface{}) (data []model.Campaign, err error)
	SetCampaign(id int64, campaign model.Campaign) (err error)
	RemoveCampaign(id []int64) (err error)
	CreateCampaign(campaign model.Campaign, kol []model.Kol) (err error)

	FindParticipant(cond map[string]interface{}) (data []model.Participant, err error)
	CreateParticipant(new model.Participant) (err error)

	FindToken() (data model.Token, err error)
	SetToken(updated map[string]interface{}) (err error)

	Login(username, password string) (data model.UserCMS, err error)
	SetTokenLogin(id uint, token string) error
	CheckToken(token string) error
	Access(userID string) (data []map[string]interface{}, err error)
	ReportCampaign(req api.Report) (data map[string]interface{}, err error)
	ReportDetail(req api.Report) (data map[string]interface{}, err error)
	ReportSummary(req api.Report) (data map[string]interface{}, err error)
}
