package model

import "github.com/jinzhu/gorm"

type Conversations struct {
	gorm.Model
	SessionID string `gorm:"session_id" json:"sessionID"`
	WAID      string `gorm:"wa_id" json:"waID"`
	Status    int    `gorm:"status" json:"status"`
	State     int    `gorm:"state" json:"state"`
	MessageID int    `gorm:"message_id" json:"messageID"`
	ExpiredAt int64  `gorm:"expired_at" json:"expiredAt"`
}

type ConversationsLog struct {
	gorm.Model
	SessionID    string `gorm:"session_id" json:"sessionID"`
	MessageID    uint   `gorm:"message_id" json:"messageID"`
	WaPayload    string `gorm:"wa_payload" json:"waPayload"`
	Chat         string `gorm:"chat" json:"chat"`
	WAID         string `gorm:"wa_id" json:"waID"`
	ChatBotTrxID string `gorm:"chat_bot_trx_id" json:"chatbotTrxID"`
}
