package usecase

import (
	"context"

	"time"

	"github.com/cyclex/planet-ban/api"
	"github.com/cyclex/planet-ban/domain"
	"github.com/cyclex/planet-ban/domain/model"
	"github.com/cyclex/planet-ban/domain/repository"
	"github.com/cyclex/planet-ban/pkg"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type cmsUcase struct {
	m              repository.ModelRepository
	contextTimeout time.Duration
}

func NewCmsUcase(m repository.ModelRepository, ctxTimeout time.Duration) domain.CmsUcase {

	return &cmsUcase{
		m:              m,
		contextTimeout: ctxTimeout,
	}
}

func (self *cmsUcase) Access(c context.Context, req api.Access) (data []map[string]interface{}, err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	res, err := self.m.Access(req.PrivilegeID)
	if err != nil {
		err = errors.Wrap(err, "[usecase.Access]")
		return
	}

	for _, v := range res {
		tmp := map[string]interface{}{
			"title": v["title"],
			"path":  v["path"],
			"pid":   v["pid"],
		}

		data = append(data, tmp)
	}

	return
}

func (self *cmsUcase) Login(c context.Context, req api.Login) (data map[string]interface{}, err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	res, err := self.m.Login(req.Username, req.Password)
	if err != nil {
		return
	}

	tokenCms := pkg.TokenGenerator(16)
	err = self.m.SetTokenLogin(res.ID, tokenCms)
	if err != nil {
		err = errors.Wrap(err, "[usecase.Login]")
	}

	data = map[string]interface{}{
		"username": res.Username,
		"user_id":  res.ID,
		"level":    res.Level,
		"token":    tokenCms,
	}

	return
}

func (self *cmsUcase) CheckToken(c context.Context, req api.CheckToken) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	err = self.m.CheckToken(req.Token)
	if err != nil {
		err = errors.Wrap(err, "[usecase.CheckToken]")
	}

	return
}

func (self *cmsUcase) Report(c context.Context, req api.Report, category string) (data map[string]interface{}, err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	switch category {
	case "campaign":
		data, err = self.m.ReportCampaign(req)
		break
	case "detail":
		data, err = self.m.ReportDetail(req)
		break
	case "detail_summary":
		data, err = self.m.ReportDetailSummary(req)
		break
	case "summary_aggregate":
		data, err = self.m.ReportSummaryAggregate(req)
		break
	case "user":
		data, err = self.m.ReportUser(req)
		break
	default:
		data, err = self.m.ReportCampaign(req)
		break
	}

	if err != nil {
		err = errors.Wrap(err, "[usecase.Report]")
	}
	return
}

func (self *cmsUcase) CreateCampaign(c context.Context, req api.Campaign) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	dataCampaign := model.Campaign{
		Name:               req.Name,
		Status:             true,
		StartDate:          req.StartDate,
		EndDate:            req.EndDate,
		DiscountProduct:    req.DiscProduct,
		DiscountProductBan: req.DiscProductBan,
		ProductName:        req.ProductName,
	}

	var dataKol []model.Kol
	for _, v := range req.Influencer {
		dataKol = append(dataKol, model.Kol{
			UID:         "",
			Name:        v.Name,
			Source:      v.Source,
			VoucherCode: v.VoucherCode,
		})
	}

	err = self.m.CreateCampaign(dataCampaign, dataKol)
	if err != nil {
		err = errors.Wrap(err, "[usecase.CreateCampaign]")
	}
	return
}

func (self *cmsUcase) SetCampaign(c context.Context, req api.Campaign) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	dataCampaign := model.Campaign{
		Name:               req.Name,
		StartDate:          req.StartDate,
		EndDate:            req.EndDate,
		DiscountProduct:    req.DiscProduct,
		DiscountProductBan: req.DiscProductBan,
		ProductName:        req.ProductName,
	}

	err = self.m.SetCampaign(req.CampaignID, dataCampaign)
	if err != nil {
		err = errors.Wrap(err, "[usecase.SetCampaign]")
	}

	return
}

func (self *cmsUcase) DeleteCampaign(c context.Context, deletedID []int64) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	err = self.m.RemoveCampaign(deletedID)
	if err != nil {
		err = errors.Wrap(err, "[usecase.DeleteCampaign]")
	}

	return
}

func (self *cmsUcase) CreateKol(c context.Context, req api.Kol) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	if req.Name == "" && req.FileName == "" {
		return errors.New("invalid request")
	}

	var (
		dataKol      []model.Kol
		skipFirstRow bool
	)

	if req.Name != "" {
		dataKol = append(dataKol, model.Kol{
			UID:         pkg.ShortUUID(uuid.NewString()),
			CampaignID:  req.CampaignID,
			Name:        req.Name,
			Source:      req.Source,
			VoucherCode: req.VoucherCode,
			AdsPlatform: req.AdsPlatform,
		})
	}

	if req.FileName != "" {
		skipFirstRow = true
		rows, err := pkg.ReadFromFile(req.FileName)
		if err != nil {
			err = errors.Wrap(err, "[usecase.CreateKol]")
			return err
		}

		for _, v := range rows {
			dataKol = append(dataKol, model.Kol{
				UID:         pkg.ShortUUID(uuid.NewString()),
				CampaignID:  req.CampaignID,
				Name:        v[2],
				Source:      v[1],
				VoucherCode: v[0],
				AdsPlatform: v[3],
			})
		}
	}

	return self.m.CreateBulkKol(dataKol, skipFirstRow)
}

func (self *cmsUcase) DeleteKol(c context.Context, deletedID int64) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	err = self.m.RemoveKol([]int64{deletedID})
	if err != nil {
		err = errors.Wrap(err, "[usecase.DeleteKol]")
	}

	return
}

func (self *cmsUcase) SetKol(c context.Context, req api.Kol) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	dataCampaign := model.Kol{
		Name:        req.Name,
		Source:      req.Source,
		VoucherCode: req.VoucherCode,
		AdsPlatform: req.AdsPlatform,
	}

	return self.m.SetKol(req.KolID, dataCampaign)
}

func (self *cmsUcase) CreateUser(c context.Context, req api.User) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	if req.Username == "" && req.Level == "" {
		return errors.New("invalid request")
	}

	var dataUser model.UserCMS

	dataUser = model.UserCMS{
		Username: req.Username,
		Level:    req.Level,
		Password: req.Password,
	}

	err = self.m.CreateUser(dataUser)
	if err != nil {
		err = errors.Wrap(err, "[usecase.CreateUser]")
	}

	return
}

func (self *cmsUcase) DeleteUser(c context.Context, deletedID int64) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	err = self.m.RemoveUser([]int64{deletedID})
	if err != nil {
		err = errors.Wrap(err, "[usecase.DeleteUser]")
	}

	return
}

func (self *cmsUcase) SetUser(c context.Context, req api.User) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	dataUser := model.UserCMS{
		Username: req.Username,
		Level:    req.Level,
		Password: req.Password,
	}

	err = self.m.SetUser(req.ID, dataUser)
	if err != nil {
		err = errors.Wrap(err, "[usecase.SetUser]")
	}

	return
}

func (self *cmsUcase) SetUserPassword(c context.Context, req api.User) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	dataUser := model.UserCMS{
		Password: req.Password,
	}

	err = self.m.SetUserPassword(req.Username, dataUser)
	if err != nil {
		err = errors.Wrap(err, "[usecase.SetUser]")
	}

	return
}
