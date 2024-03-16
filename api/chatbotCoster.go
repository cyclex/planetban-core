package api

type TemplateLanguage struct {
	Policy string `json:"policy"`
	Code   string `json:"code"`
}

type Parameter struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Component struct {
	Type       string      `json:"type"`
	Parameters []Parameter `json:"parameters"`
}

type TemplateCoster struct {
	Name       string           `json:"name"`
	Language   TemplateLanguage `json:"language"`
	Components []Component      `json:"components"`
}

type ReqMessageCoster struct {
	ID       interface{}    `json:"id,omitempty"`
	XID      string         `json:"xid"`
	To       string         `json:"to"`
	Type     string         `json:"type"`
	Template TemplateCoster `json:"template"`
}

type Meta struct {
	Author  string `json:"author"`
	Meta    string `json:"meta,omitempty"`
	Version string `json:"version,omitempty"`
}

type Contact struct {
	Input string `json:"input"`
	WAID  string `json:"wa_id"`
}

type ErrorDetail struct {
	Code   int    `json:"code"`
	Detail string `json:"detail"`
	Title  string `json:"title"`
}

type ResponseChatbotCoster struct {
	Error    ErrorDetail        `json:"error"`
	Meta     Meta               `json:"meta"`
	Contacts []Contact          `json:"contacts"`
	Messages []ReqMessageCoster `json:"messages"`
}
