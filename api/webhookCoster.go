package api

type Metadata struct {
	PhoneNumberID      string `json:"phone_number_id"`
	DisplayPhoneNumber string `json:"display_phone_number"`
}

type Origin struct {
	Type string `json:"type"`
}

type Conversation struct {
	ID                  string      `json:"id"`
	Origin              Origin      `json:"origin"`
	ExpirationTimestamp interface{} `json:"expiration_timestamp"`
}

type Status struct {
	ID           string       `json:"id"`
	Errors       interface{}  `json:"errors"`
	Status       string       `json:"status"`
	Pricing      Pricing      `json:"pricing"`
	Timestamp    string       `json:"timestamp"`
	Conversation Conversation `json:"conversation"`
	RecipientID  string       `json:"recipient_id"`
}

type Pricing struct {
	Billable     bool   `json:"billable"`
	Category     string `json:"category"`
	PricingModel string `json:"pricing_model"`
}

type MessageValue struct {
	Metadata         Metadata `json:"metadata"`
	Statuses         []Status `json:"statuses"`
	MessagingProduct string   `json:"messaging_product"`
}

type Change struct {
	Field string       `json:"field"`
	Value MessageValue `json:"value"`
}

type Entry struct {
	ID      string   `json:"id"`
	Time    int      `json:"time"`
	Changes []Change `json:"changes"`
}

type Reroute struct {
	DateHit        string `json:"date_hit"`
	MessageID      string `json:"message_id"`
	Description    string `json:"description"`
	DateReceived   string `json:"date_received"`
	DeliveryStatus int    `json:"delivery_status"`
}

type webhookCoster struct {
	XID     string    `json:"xid"`
	Object  string    `json:"object"`
	Entry   []Entry   `json:"entry"`
	Reroute []Reroute `json:"reroute"`
}
