package model

import "github.com/jinzhu/gorm"

type Campaign struct {
	gorm.Model
	Name               string `gorm:"name" json:"name"`
	Status             bool   `gorm:"status" json:"status"`
	StartDate          int64  `gorm:"start_date" json:"startDate"`
	EndDate            int64  `gorm:"end_date" json:"endDate"`
	DiscountProduct    string `gorm:"discount_product" json:"discountProduct"`
	DiscountProductBan string `gorm:"discount_product_ban" json:"discountProductBan"`
	ProductName        string `gorm:"product_name" json:"productName"`
	CreatedAt          int64  `gorm:"autoCreateTime"`
	UpdatedAt          int64  `gorm:"autoUpdateTime"`
}

type Kol struct {
	gorm.Model
	UID         string `gorm:"uid" json:"uid"`
	CampaignID  int64  `gorm:"campaign_id" json:"campaignID"`
	Name        string `gorm:"name" json:"name"`
	Source      string `gorm:"source" json:"source"`
	VoucherCode string `gorm:"voucher_code" json:"voucherCode"`
	Status      int64  `gorm:"status" json:"status"`
	AdsPlatform string `gorm:"add_platform" json:"ads_platform"`
	CreatedAt   int64  `gorm:"autoCreateTime"`
	UpdatedAt   int64  `gorm:"autoUpdateTime"`
}

type Participant struct {
	gorm.Model
	MSISDN       string `gorm:"msisdn" json:"msisdn"`
	CampaignID   int64  `gorm:"campaign_id" json:"campaignID"`
	KolID        int64  `gorm:"kol_id" json:"kolID"`
	Status       bool   `gorm:"status" json:"status"`
	CreatedAt    int64  `gorm:"autoCreateTime"`
	UpdatedAt    int64  `gorm:"autoUpdateTime"`
	KSource      string `gorm:"k_source" json:"k_source"`
	KAdsPlatform string `gorm:"k_add_platform" json:"k_ads_platform"`
	KName        string `gorm:"k_name" json:"k_name"`
}
