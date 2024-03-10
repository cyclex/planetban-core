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
	From    string `json:"from" validate:"required"`
	To      string `json:"to" validate:"required"`
	Offset  string `json:"offset"`
	Limit   string `json:"limit"`
	Column  string `json:"column"`
	Keyword string `json:"keyword"`
}

type Campaign struct {
	CampaignID     int64  `json:"campaignID"`
	Name           string `json:"name" validate:"required"`
	StartDate      int64  `json:"startDate" validate:"required"`
	EndDate        int64  `json:"endDate" validate:"required"`
	DiscProduct    int64  `json:"discProduct" validate:"required"`
	DiscProductBan int64  `json:"discProductBan" validate:"required"`
	Influencer     []Kol  `json:"kol" validate:"required"`
	File           string `json:"file"`
}

type Kol struct {
	FileName    string `json:"file"`
	CampaignID  int64  `json:"campaignID"`
	KolID       int64  `json:"kolID"`
	Name        string `json:"name"`
	Source      string `json:"source"`
	VoucherCode string `json:"voucherCode"`
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
