package domain

import (
	"context"

	"github.com/cyclex/planet-ban/api"
)

type CmsUcase interface {
	Login(c context.Context, req api.Login) (data map[string]interface{}, err error)
	CheckToken(c context.Context, req api.CheckToken) error
	Access(c context.Context, req api.Access) (data []map[string]interface{}, err error)
	Report(c context.Context, req api.Report, category string) (data map[string]interface{}, err error)
	CreateCampaign(c context.Context, req api.Campaign) (err error)
	SetCampaign(c context.Context, req api.Campaign) (err error)
	DeleteCampaign(c context.Context, deletedID []int64) (err error)
	CreateKol(c context.Context, req api.Kol) (err error)
	SetKol(c context.Context, req api.Kol) (err error)
	DeleteKol(c context.Context, deletedID int64) (err error)
}
