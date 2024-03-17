package domain

import (
	"github.com/cyclex/planet-ban/api"
)

type ChatUcase interface {
	GetWhatsappTemplateMessage(id string) (message string, err error)
	GetWabaAccountNumber() (msisdn string)

	IncomingMessages(payload api.Message) (trxChatBotID string, err error)
	ReplyMessages(waID, incoming string) (outgoing string, err error)
	ChatToUser(waID, chat, types, media string) (res []byte, statusCode int, err error)
}
