package postgre

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"

	"github.com/cyclex/planet-ban/api"
	"github.com/cyclex/planet-ban/domain/model"

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

	self.DB.Table("v_access").Where("privilege_id = ?", privilegeID).First(&res)

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

	err = self.DB.Where(cond).First(&data).Error
	if err != nil {
		err = errors.New("Invalid username or password. Please try again.")
	}

	return
}

func (self *postgreRepo) CheckToken(token string) (err error) {

	var data model.UserCMS

	cond := map[string]interface{}{
		"token": token,
	}

	err = self.DB.Where(cond).First(&data).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.CheckToken]")
	}

	return
}

func (self *postgreRepo) SetTokenLogin(id uint, token string) (err error) {

	err = self.DB.Model(&model.UserCMS{}).Where(map[string]interface{}{"id": id}).Updates(map[string]interface{}{"token": token}).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.SetTokenLogin]")
	}
	return

}

func (self *postgreRepo) ReportCampaign(req api.Report) (data map[string]interface{}, err error) {

	type tmp struct {
		Rnum      string `json:"rnum"`
		Id        string `json:"id"`
		Name      string `json:"name"`
		CreatedAt string `json:"created_at"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	var (
		res   []tmp
		cond  map[string]interface{}
		datas []map[string]interface{}
		rows  int64
	)

	q := self.DB.Model(&model.Campaign{}).Select("*, row_number() OVER () as rnum").Where(cond)

	if req.Keyword != "" {
		column := fmt.Sprintf("%s ilike ?", "name")
		q = q.Where(column, "%"+req.Keyword+"%")
	}

	q.Count(&rows)
	err = q.Order("created_at desc").Limit(req.Limit).Offset(req.Offset).Find(&res).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.ReportCampaign]")
		return
	}

	for _, v := range res {
		createdAt, _ := strconv.ParseInt(v.CreatedAt, 10, 64)
		startDate, _ := strconv.ParseInt(v.StartDate, 10, 64)
		endDate, _ := strconv.ParseInt(v.EndDate, 10, 64)

		x := map[string]interface{}{
			"rNum":         v.Rnum,
			"campaignName": v.Name,
			"createdAt":    time.Unix(createdAt, 0).Local().Format("2006-01-02 15:04:05"),
			"startDate":    time.Unix(startDate, 0).Local().Format("2006-01-02 15:04:05"),
			"endDate":      time.Unix(endDate, 0).Local().Format("2006-01-02 15:04:05"),
			"id":           v.Id,
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

	type tmp struct {
		Rnum        string `json:"rnum"`
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Source      string `json:"source"`
		AdsPlatform string `json:"ads_platform"`
		VoucherCode string `json:"voucher_code"`
		UID         string `json:"uid"`
	}

	var (
		res   []tmp
		datas []map[string]interface{}
		cond  map[string]interface{}
		rows  int64
	)

	cond = map[string]interface{}{"campaign_id": req.Keyword}

	q := self.DB.Model(&model.Kol{}).Select("*, row_number() OVER () as rnum").Where(cond)

	q.Count(&rows)
	err = q.Order("created_at desc").Limit(req.Limit).Offset(req.Offset).Find(&res).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.ReportDetail]")
		return
	}

	for _, v := range res {
		x := map[string]interface{}{
			"rNum":        v.Rnum,
			"voucherCode": v.VoucherCode,
			"source":      v.Source,
			"kolName":     v.Name,
			"adsPlatform": v.AdsPlatform,
			"urlLink":     fmt.Sprintf("%s/%s", self.UrlHost, v.UID),
			"id":          v.ID,
		}

		datas = append(datas, x)
	}

	data = map[string]interface{}{
		"rows": rows,
		"data": datas,
	}

	return
}

func (self *postgreRepo) ReportSummary(req api.Report) (data map[string]interface{}, err error) {

	type summary struct {
		Rnum         string `json:"rnum"`
		VoucherCode  string `json:"voucher_code"`
		Source       string `json:"source"`
		Name         string `json:"name"`
		CampaignName string `json:"campaign_name"`
		AdsPlatform  string `json:"ads_platform"`
		Msisdn       string `json:"msisdn"`
		CreatedAt    string `json:"created_at"`
	}
	var (
		sum   []summary
		datas []map[string]interface{}
		rows  int64
	)

	q := self.DB.Model(&model.Participant{}).Select("k.voucher_code,k.source,k.name, k.ads_platform,participants.msisdn, participants.created_at, c.name as campaign_name, row_number() OVER () as rnum").Joins("join kols k on participants.kol_id = k.id").Joins("join campaigns c on c.id = k.campaign_id")

	if req.From != 0 || req.To != 0 {
		q = q.Where("participants.created_at BETWEEN ? AND ?", req.From, req.To)
	}

	if req.Keyword != "" {
		column := fmt.Sprintf("%s ilike ?", "c.name")
		q = q.Where(column, "%"+req.Keyword+"%")
	}

	q.Count(&rows)
	err = q.Order("participants.created_at desc").Limit(req.Limit).Offset(req.Offset).Find(&sum).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.ReportSummary]")
		return
	}

	for _, v := range sum {
		createdAt, _ := strconv.ParseInt(v.CreatedAt, 10, 64)

		x := map[string]interface{}{
			"rNum":         v.Rnum,
			"voucherCode":  v.VoucherCode,
			"source":       v.Source,
			"kolName":      v.Name,
			"campaignName": v.CampaignName,
			"adsPlatform":  v.AdsPlatform,
			"msisdn":       v.Msisdn,
			"created_at":   time.Unix(createdAt, 0).Local().Format("2006-01-02 15:04:05"),
		}
		datas = append(datas, x)
	}

	data = map[string]interface{}{
		"rows": rows,
		"data": datas,
	}

	return
}

func (self *postgreRepo) ReportSummaryAggregate(req api.Report) (data map[string]interface{}, err error) {

	type summary struct {
		Rnum          string `json:"rnum"`
		VoucherCode   string `json:"voucher_code"`
		Source        string `json:"source"`
		Name          string `json:"name"`
		CampaignName  string `json:"campaign_name"`
		AdsPlatform   string `json:"ads_platform"`
		TotalReceived int64  `json:"total_received"`
	}
	var (
		sum   []summary
		datas []map[string]interface{}
		rows  int64
	)

	q := self.DB.Model(&model.Participant{}).Select("k.voucher_code,k.source,k.name, k.ads_platform, k.created_at, COUNT(k.name) AS total_received, c.name as campaign_name, row_number() OVER () as rnum").Joins("join kols k on participants.kol_id = k.id").Joins("join campaigns c on c.id = k.campaign_id")

	if req.From != 0 || req.To != 0 {
		q = q.Where("participants.created_at BETWEEN ? AND ?", req.From, req.To)
	}

	if req.Keyword != "" {
		column := fmt.Sprintf("%s ilike ?", "c.name")
		q = q.Where(column, "%"+req.Keyword+"%")
	}

	q.Count(&rows)
	q = q.Order("k.created_at desc")
	q = q.Group("k.voucher_code, k.source, k.name, k.ads_platform, k.created_at, c.name")
	err = q.Limit(req.Limit).Offset(req.Offset).Find(&sum).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.ReportSummary]")
		return
	}

	for _, v := range sum {
		x := map[string]interface{}{
			"rNum":          v.Rnum,
			"voucherCode":   v.VoucherCode,
			"source":        v.Source,
			"kolName":       v.Name,
			"campaignName":  v.CampaignName,
			"adsPlatform":   v.AdsPlatform,
			"totalReceived": v.TotalReceived,
		}
		datas = append(datas, x)
	}

	data = map[string]interface{}{
		"rows": rows,
		"data": datas,
	}

	return
}
