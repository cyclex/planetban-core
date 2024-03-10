package domain

import (
	"github.com/cyclex/planet-ban/api"
	"github.com/cyclex/planet-ban/domain/model"
)

type ChatUcase interface {
	IncomingMessages(payload api.Message) (trxChatBotID string, err error)
	ReplyMessages(waID, incoming string) (outgoing string, err error)
	ChatToUser(waID, chat, types, media string) (res []byte, statusCode int, err error)

	FindToken() (data model.Token, err error)
	SetToken(updated map[string]interface{}) (err error)
	RefreshToken() (res []byte, statusCode int, err error)
}
