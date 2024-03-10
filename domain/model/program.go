package model

import "github.com/jinzhu/gorm"

type Campaign struct {
	gorm.Model
	Name               string `gorm:"name" json:"name"`
	Status             bool   `gorm:"status" json:"status"`
	StartDate          int64  `gorm:"start_date" json:"startDate"`
	EndDate            int64  `gorm:"end_date" json:"endDate"`
	DiscountProduct    int64  `gorm:"discount_product" json:"discountProduct"`
	DiscountProductBan int64  `gorm:"discount_product_ban" json:"discountProductBan"`
}

type Kol struct {
	gorm.Model
	UID         string `gorm:"uid" json:"uid"`
	CampaignID  int64  `gorm:"campaign_id" json:"campaignID"`
	Name        string `gorm:"name" json:"name"`
	Source      string `gorm:"source" json:"source"`
	VoucherCode string `gorm:"voucher_code" json:"voucherCode"`
	URL         string `gorm:"url" json:"url"`
}

type Participant struct {
	gorm.Model
	MSISDN     string `gorm:"msisdn" json:"msisdn"`
	CampaignID int64  `gorm:"campaign_id" json:"campaignID"`
	KolID      int64  `gorm:"kol_id" json:"kolID"`
	Status     bool   `gorm:"status" json:"status"`
}

// type Report struct {
// 	gorm.Model
// 	ReportID     string `gorm:"report_id" json:"reportID"`
// 	RequestBy    string `gorm:"request_by" json:"requestBy"`
// 	StartDate    int64  `gorm:"start_date" json:"startDate"`
// 	EndDate      int64  `gorm:"end_date" json:"endDate"`
// 	Status       bool   `gorm:"status" json:"status"`
// 	DownloadFile string `gorm:"download_file" json:"downloadFile"`
// }
