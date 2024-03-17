package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/cyclex/planet-ban/api"
	"github.com/cyclex/planet-ban/domain"
	"github.com/cyclex/planet-ban/domain/model"
	"github.com/cyclex/planet-ban/domain/repository"
	"github.com/cyclex/planet-ban/pkg"
	"github.com/cyclex/planet-ban/pkg/httprequest"
	"github.com/jinzhu/gorm"
)

type chatUcase struct {
	m                 repository.ModelRepository
	urlSendMsg        string
	wabaAccountNumber string
	DivisionID        string
	AccountID         string
	AccessToken       string
}

func NewChatUcase(m repository.ModelRepository, urlSendMsg, divisionID, accountID, accessToken, wabaAccountNumber string) domain.ChatUcase {

	return &chatUcase{
		m:                 m,
		urlSendMsg:        urlSendMsg,
		DivisionID:        divisionID,
		AccountID:         accountID,
		AccessToken:       accessToken,
		wabaAccountNumber: wabaAccountNumber,
	}
}

func (self *chatUcase) ReplyMessages(waID, incoming string) (outgoing string, err error) {

	fmt.Println(incoming)
	var cond = map[string]interface{}{}
	outgoing = "Maaf, kami tidak mengerti maksud anda. Silahkan menggunakan format chat yang sudah ditentukan"

	usernames := pkg.ExtractUsernames(incoming)
	campaign := pkg.ExtractSentencesAfterWord(incoming, "promo")
	fmt.Println(usernames, campaign)

	if len(usernames) == 0 || len(campaign) == 0 {
		return
	}

	cond = map[string]interface{}{
		"name": campaign[0],
	}
	dataCampaign, err := self.m.FindCampaignBy(cond)
	if err != nil {
		err = errors.Wrap(err, "[usecase.ReplyMessages]")
		return
	}

	if len(dataCampaign) == 0 {
		return
	}

	if !pkg.IsTimeInBetween(dataCampaign[0].StartDate, dataCampaign[0].EndDate) {
		outgoing = "Mohon maaf, program yang anda ikuti sudah berakhir, silahkan menggunakan kode yang sedang digunakan"
		return
	}

	cond = map[string]interface{}{
		"name":        usernames[0],
		"campaign_id": dataCampaign[0].ID,
	}
	dataKol, err := self.m.FindKolBy(cond)
	if err != nil {
		err = errors.Wrap(err, "[usecase.ReplyMessages]")
		return
	}

	if len(dataKol) == 0 {
		return
	}

	cond = map[string]interface{}{
		"msisdn":      waID,
		"campaign_id": dataCampaign[0].ID,
		"kol_id":      dataKol[0].ID,
	}
	dataParticipant, err := self.m.FindParticipant(cond)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			err = errors.Wrap(err, "[usecase.ReplyMessages]")
			return
		}
	}

	timeLeft := time.Unix(dataCampaign[0].EndDate, 0).Local().Format("2006-01-02")
	if len(dataParticipant) > 0 {
		outgoing = fmt.Sprintf("Halo Planeters. Segera gunakan kode voucher kamu yang berlaku sd %s di seluruh toko Planet", timeLeft)
		return
	}

	createCP := model.Participant{
		MSISDN:     waID,
		CampaignID: int64(dataCampaign[0].ID),
		KolID:      int64(dataKol[0].ID),
		Status:     true,
	}
	err = self.m.CreateParticipant(createCP)
	if err != nil {
		err = errors.Wrap(err, "[usecase.ReplyMessages]")
		return
	}

	outgoing = fmt.Sprintf("Halo Planeters! Berikut kode voucher kamu: %s Kode voucher bisa digunakan untuk mendapat Diskon %s pada Pembelian produk %s Voucher berlaku sd %s di seluruh toko Planet Ban.", dataKol[0].VoucherCode, dataCampaign[0].DiscountProduct, dataCampaign[0].ProductName, timeLeft)
	return
}

func (self *chatUcase) ChatToUser(waID, chat, types, media string) (res []byte, statusCode int, err error) {

	var payload interface{}
	url := self.urlSendMsg

	payload = api.ReqSendMessageText{
		XID:         uuid.NewString(),
		ChannelID:   "whatsapp-cloud",
		AccountID:   self.AccountID,
		DivisionID:  self.DivisionID,
		IsHelpdesk:  false,
		MessageType: "outbound",
		Data: api.Data{
			MessagingProduct: "whatsapp",
			RecipientType:    "individual",
			To:               waID,
			Type:             types,
			Text: api.Text{
				Body: chat,
			},
		},
	}

	res, statusCode, err = httprequest.PostJson(url, payload, 15*time.Second, self.AccessToken)
	if err != nil {
		err = errors.Wrap(err, "[usecase.ChatToUser]")
	}
	return
}

func (self *chatUcase) IncomingMessages(payload api.Message) (trxChatBotID string, err error) {

	jsonPayload, _ := json.Marshal(payload)
	waID := payload.From
	sessID := uuid.NewString()
	newCLog := model.ConversationsLog{
		SessionID: sessID,
		Chat:      strings.ToUpper(payload.Text.Body),
		WAID:      waID,
		WaPayload: string(jsonPayload),
	}
	err = self.m.CreateConversationsLog(newCLog)
	if err != nil {
		err = errors.Wrap(err, "[usecase.IncomingMessages]")
		return
	}

	outgoing, err := self.ReplyMessages(waID, strings.ToLower(newCLog.Chat))
	if err != nil {
		err = errors.Wrap(err, "[usecase.IncomingMessages]")
		log.Error(err.Error())
	}

	fmt.Println(outgoing)
	res, statusCode, err := self.ChatToUser(waID, outgoing, "text", "")
	if err != nil {
		err = errors.Wrap(err, "[usecase.IncomingMessages]")
		fmt.Println(err.Error())
		return
	}

	if statusCode != http.StatusOK {
		err = errors.New(http.StatusText(statusCode))
		err = errors.Wrap(err, "[usecase.IncomingMessages]")
		fmt.Println(err.Error())
		fmt.Printf("%s\n", res)
		return
	}

	newCLog = model.ConversationsLog{
		SessionID:    sessID,
		Chat:         outgoing,
		WAID:         waID,
		WaPayload:    string(res),
		ChatBotTrxID: trxChatBotID,
	}
	err = self.m.CreateConversationsLog(newCLog)
	if err != nil {
		err = errors.Wrap(err, "[usecase.IncomingMessages]")
	}

	return
}

func (self *chatUcase) GetWhatsappTemplateMessage(id string) (message string, err error) {

	cond := map[string]interface{}{
		"uid": id,
	}
	dataKol, err := self.m.FindKolBy(cond)
	if err != nil {
		err = errors.Wrap(err, "[usecase.GetWhatsappTemplateMessage]")
		return
	}

	if len(dataKol) == 0 {
		return
	}

	cond = map[string]interface{}{
		"id": dataKol[0].CampaignID,
	}
	dataCampaign, err := self.m.FindCampaignBy(cond)
	if err != nil {
		err = errors.Wrap(err, "[usecase.GetWhatsappTemplateMessage]")
		return
	}

	message = fmt.Sprintf("halo %s anda di %s", dataKol[0].Name, dataCampaign[0].Name)

	return
}

func (self *chatUcase) GetWabaAccountNumber() (msisdn string) {
	return self.wabaAccountNumber
}
