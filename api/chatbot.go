package api

type PayloadReply struct {
	WaID        string `json:"waID"`
	Chat        string `json:"chat"`
	SessionID   string `json:"sessionID"`
	MessageID   uint   `json:"msgID"`
	ScheduledAt int64  `json:"scheduledAt"`
}

type ResSendMessage struct {
	Messages []Message  `json:"messages"`
	Contacts []Contacts `json:"contacts"`
	Err      []ErrDesc  `json:"errors"`
}

type ErrDesc struct {
	Code    string `json:"code"`
	Title   string `json:"title"`
	Details string `json:"details"`
}

type Message struct {
	ID        string  `json:"id" bson:"id"`
	From      string  `json:"from" bson:"from"`
	Text      *Text   `json:"text,omitempty" bson:"text"`
	Timestamp string  `json:"timestamp" bson:"timestamp"`
	Type      string  `json:"type" bson:"type"`
	Images    *Images `json:"image,omitempty" bson:"image"`
}

type Text struct {
	Body string `json:"body" bson:"body"`
}

type Contacts struct {
	WaID    string `json:"wa_id"`
	Profile struct {
		Name string `json:"name"`
	} `json:"profile"`
}

type Image struct {
	Link    string `json:"link" bson:"link"`
	Caption string `json:"caption" bson:"caption"`
}

type Images struct {
	MimeType string `json:"mime_type" bson:"mime_type"`
	ID       string `json:"id" bson:"id"`
}

type ReqSendMessageText struct {
	RecipientType string `json:"recipient_type"`
	To            string `json:"to"`
	Text          Text   `json:"text"`
	Type          string `json:"type"`
}

type ReqSendMessageImage struct {
	RecipientType string `json:"recipient_type"`
	To            string `json:"to"`
	Type          string `json:"type"`
	Image         Image  `json:"image"`
}

type ReqSendBroadcast struct {
	To       string   `json:"to"`
	Type     string   `json:"type"`
	Template Template `json:"template"`
	// Hsm      Hsm      `json:"hsm"`
}

type Template struct {
	Namespace  string       `json:"namespace"`
	Name       string       `json:"name"`
	Language   Language     `json:"language"`
	Components []Components `json:"components"`
}

type Hsm struct {
	Namespace   string       `json:"namespace"`
	ElementName string       `json:"element_name"`
	Language    Language     `json:"language"`
	LocalParam  []LocalParam `json:"localizable_params"`
}

type LocalParam struct {
	Default string `json:"default"`
}

type Language struct {
	Code   string `json:"code"`
	Policy string `json:"policy"`
}

type Components struct {
	Type       string       `json:"type"`
	Parameters []Parameters `json:"parameters"`
}

type Parameters struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ResLogin struct {
	Users []User `json:"users"`
}

type User struct {
	Token        string `json:"token"`
	ExpiresAfter string `json:"expires_after"`
}
