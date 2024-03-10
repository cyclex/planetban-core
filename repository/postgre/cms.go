package postgre

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"

	"github.com/cyclex/planet-ban/api"
	"github.com/cyclex/planet-ban/domain/model"
	"github.com/jinzhu/gorm"

	"github.com/pkg/errors"
)

var (
	LayoutDefault  = time.RFC3339Nano
	Loc, _         = time.LoadLocation("Asia/Jakarta")
	LayoutDateTime = "2006-01-02 15:04:05"
)

func (self *postgreRepo) Access(privilegeID string) (data []map[string]interface{}, err error) {

	type tmp struct {
		Title string `json:"title"`
		Path  string `json:"path"`
		Pid   string `json:"pid"`
	}

	var res []tmp

	err = self.DB.Table("v_access").Where("privilege_id = ?", privilegeID).Find(&res).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.Access]")
		return
	}

	for _, v := range res {
		x := map[string]interface{}{
			"title": v.Title,
			"path":  v.Path,
			"pid":   v.Pid,
		}

		data = append(data, x)
	}

	return
}

func (self *postgreRepo) Login(username, password string) (data model.UserCMS, err error) {

	password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
	cond := map[string]interface{}{
		"username": username,
		"password": password,
	}

	err = self.DB.Table("user_cms").Where(cond).First(&data).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			err = errors.Wrap(err, "[postgre.Login]")
			return
		}
	}

	return
}

func (self *postgreRepo) CheckToken(token string) (err error) {

	var data model.UserCMS

	cond := map[string]interface{}{
		"token": token,
	}

	err = self.DB.Table("user_cms").Where(cond).Find(&data).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.CheckToken]")
	}

	return
}

func (self *postgreRepo) SetTokenLogin(id uint, token string) (err error) {

	err = self.DB.Table("user_cms").Where(map[string]interface{}{"id": id}).Updates(map[string]interface{}{"token": token}).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.SetTokenLogin]")
	}
	return

}

func (self *postgreRepo) ReportCampaign(req api.Report) (data map[string]interface{}, err error) {

	type tmp struct {
		Rnum         string `json:"rnum"`
		CampaignName string `json:"campaignName"`
		CreatedAt    string `json:"createdAt"`
		StartDate    string `json:"startDate"`
		EndDate      string `json:"endDate"`
	}

	var (
		res   []tmp
		cond  map[string]interface{}
		datas []map[string]interface{}
		rows  int64
	)

	q := self.DB.Table("campaign").Select("*, row_number() OVER () as rnum").Where(cond)

	if req.Keyword != "" {
		column := fmt.Sprintf("%s ilike ?", "name")
		q = q.Where(column, "%"+req.Keyword+"%")
	}

	q.Count(&rows)
	err = q.Order("created_at desc").Find(&res).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			err = errors.Wrap(err, "[postgre.ReportCampaign]")
			return
		}
	}

	for _, v := range res {

		createdAt, _ := time.Parse(LayoutDefault, v.CreatedAt)
		startDate, _ := time.Parse(LayoutDefault, v.StartDate)
		endDate, _ := time.Parse(LayoutDefault, v.EndDate)

		x := map[string]interface{}{
			"rNum":         v.Rnum,
			"campaignName": v.CampaignName,
			"createdAt":    createdAt.In(Loc).Format(LayoutDateTime),
			"startDate":    startDate.In(Loc).Format(LayoutDateTime),
			"endDate":      endDate.In(Loc).Format(LayoutDateTime),
		}

		datas = append(datas, x)
	}

	data = map[string]interface{}{
		"rows": rows,
		"data": datas,
	}

	return
}

func (self *postgreRepo) ReportDetail(req api.Report) (data map[string]interface{}, err error) {

	type summary struct {
		VoucherCode  string `json:"voucherCode"`
		Source       string `json:"source"`
		KolName      string `json:"kolName"`
		URLLink      string `json:"urlLink"`
		CampaignName string `json:"campaignName"`
		Msisdn       string `json:"msisdn"`
		CreatedAt    string `json:"createdAt"`
	}

	var (
		sum   []summary
		datas []map[string]interface{}
		cond  map[string]interface{}
	)

	cond = map[string]interface{}{"campaign_id": req.Keyword}

	err = self.DB.Table("v_campaign").Select("*").Where(cond).Order("created_at desc").Find(&sum).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			err = errors.Wrap(err, "[postgre.ReportDetail]")
			return
		}
	}

	for _, v := range sum {
		x := map[string]interface{}{
			"voucherCode": v.VoucherCode,
			"source":      v.Source,
			"kolName":     v.KolName,
			"urlLink":     v.URLLink,
		}

		datas = append(datas, x)
	}

	data = map[string]interface{}{
		"rows": len(datas),
		"data": datas,
	}

	return
}

func (self *postgreRepo) ReportSummary(req api.Report) (data map[string]interface{}, err error) {

	type summary struct {
		VoucherCode   string `json:"voucherCode"`
		Source        string `json:"source"`
		KolName       string `json:"kolName"`
		CampaignName  string `json:"campaignName"`
		TotalReceived string `json:"totalReceived"`
	}
	var (
		sum   []summary
		datas []map[string]interface{}
	)

	q := self.DB.Table("v_campaign_summary").Select("*")
	if req.From != "" {
		q = q.Where("created_at BETWEEN ? AND ?", req.From, req.To)
	}

	limit, _ := strconv.Atoi(req.Limit)
	offset, _ := strconv.Atoi(req.Offset)
	err = q.Order("created_at desc").Limit(limit).Offset(offset).Find(&sum).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			err = errors.Wrap(err, "[postgre.ReportSummary]")
			return
		}
	}

	for _, v := range sum {
		x := map[string]interface{}{
			"voucherCode":   v.VoucherCode,
			"source":        v.Source,
			"kolName":       v.KolName,
			"campaigName":   v.CampaignName,
			"totalReceived": v.TotalReceived,
		}
		datas = append(datas, x)
	}

	data = map[string]interface{}{
		"rows": len(datas),
		"data": datas,
	}

	return
}
