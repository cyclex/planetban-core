package api

type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"pass" validate:"required"`
}

type CheckToken struct {
	Token string `json:"token" validate:"required"`
}

type Access struct {
	UserID      string `json:"user_id"`
	PrivilegeID string `json:"privilege_id" validate:"required"`
}

type Report struct {
	From    int    `json:"from" validate:"required"`
	To      int    `json:"to" validate:"required"`
	Offset  int    `json:"offset"`
	Limit   int    `json:"limit"`
	Column  string `json:"column"`
	Keyword string `json:"keyword"`
}

type Campaign struct {
	CampaignID     int64  `json:"campaignID"`
	Name           string `json:"name" validate:"required"`
	StartDate      int64  `json:"startDate" validate:"required"`
	EndDate        int64  `json:"endDate" validate:"required"`
	DiscProduct    string `json:"discProduct" validate:"required"`
	DiscProductBan string `json:"discProductBan" validate:"required"`
	ProductName    string `json:"productName" validate:"required"`
	Influencer     []Kol  `json:"kol"`
	File           string `json:"file"`
}

type Kol struct {
	FileName    string `json:"file"`
	CampaignID  int64  `json:"campaignID" validate:"required"`
	KolID       int64  `json:"kolID"`
	Name        string `json:"name"`
	Source      string `json:"source" validate:"required"`
	VoucherCode string `json:"voucherCode" validate:"required"`
	AdsPlatform string `json:"adsPlatform" validate:"required"`
}

type ResponseError struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type ResponseSuccess struct {
	Status  bool                   `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type ResponseReport struct {
	Status  bool                     `json:"status"`
	Message string                   `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}
