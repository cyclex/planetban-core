package usecase

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/cyclex/planet-ban/api"
	"github.com/cyclex/planet-ban/domain"
	"github.com/cyclex/planet-ban/domain/model"
	"github.com/cyclex/planet-ban/domain/repository"
	"github.com/cyclex/planet-ban/pkg"
	"github.com/cyclex/planet-ban/pkg/httprequest"
	"github.com/jinzhu/gorm"
)

type chatUcase struct {
	m                  repository.ModelRepository
	q                  repository.QueueRepository
	contextTimeout     time.Duration
	urlSendMsg         string
	urlMedia           string
	namespace          string
	parameterNamespace string
	wabaAccountNumber  string
}

func NewChatUcase(m repository.ModelRepository, urlSendMsg, urlMedia, nameSpace, parameterNamespace, wabaAccountNumber string) domain.ChatUcase {

	return &chatUcase{
		m:                  m,
		urlSendMsg:         urlSendMsg,
		urlMedia:           urlMedia,
		namespace:          nameSpace,
		parameterNamespace: parameterNamespace,
		wabaAccountNumber:  wabaAccountNumber,
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
	url := self.urlSendMsg + "v2/messages"

	payload = api.ReqMessageCoster{
		XID:  "{{UNIQUE-ID-FROM-CLIENT}}",
		To:   waID,
		Type: "template",
		Template: api.TemplateCoster{
			Name: "{{TEMPLATE-NAME}}",
			Language: api.TemplateLanguage{
				Policy: "deterministic",
				Code:   "{{TEMPLATE-LANGUAGE}}",
			},
			Components: []api.Component{
				{
					Type: "header",
					Parameters: []api.Parameter{
						{
							Type: types,
							Text: "{{PARAM-HEADER-TEXT}}",
						},
					},
				},
			},
		},
	}

	// TODO get token from redis
	tokenChatbot, _ := self.m.FindToken()
	res, statusCode, err = httprequest.PostJson(url, payload, self.contextTimeout, "Bearer "+tokenChatbot.AccessToken)
	if err != nil {
		err = errors.Wrap(err, "[usecase.ChatToUser]")
	}
	return
}

func (self *chatUcase) ChatToUserV1(waID, chat, types, media string) (res []byte, statusCode int, err error) {

	var payload interface{}
	url := self.urlSendMsg + "v2/messages"

	if types == "text" {
		payload = api.ReqSendMessageText{
			RecipientType: "individual",
			To:            waID,
			Type:          types,
			Text: api.Text{
				Body: chat,
			},
		}
	} else if types == "image" {
		payload = api.ReqSendMessageImage{
			RecipientType: "individual",
			To:            waID,
			Type:          types,
			Image: api.Image{
				Link:    media,
				Caption: chat,
			},
		}
	} else if types == "broadcast" {
		var (
			dataParams []api.Parameters
			dataComp   []api.Components
		)

		tmpParam := api.Parameters{
			Type: "text",
			Text: chat,
		}
		for i := 0; i < 1; i++ {
			dataParams = append(dataParams, tmpParam)
		}

		tmpComp := api.Components{
			Type:       "body",
			Parameters: dataParams,
		}

		for i := 0; i < 1; i++ {
			dataComp = append(dataComp, tmpComp)
		}

		payload = api.ReqSendBroadcast{
			To:   waID,
			Type: "template",
			// Hsm: api.Hsm{
			// 	Namespace:   "26a3fc53_044a_4a7b_bf37_dd407f271c48",
			// 	ElementName: "info_hadiah",
			// 	Language: api.Language{
			// 		Code:   "id",
			// 		Policy: "deterministic",
			// 	},
			// 	LocalParam: dataParams,
			// },
			Template: api.Template{
				Namespace: self.namespace,
				Name:      self.parameterNamespace,
				Language: api.Language{
					Policy: "deterministic",
					Code:   "id",
				},
				Components: dataComp,
			},
		}

		// winLog.Infof("WhatsAppID:%s , payload:%+v", waID, payload)
	}

	// TODO get token from redis
	tokenChatbot, _ := self.m.FindToken()
	res, statusCode, err = httprequest.PostJson(url, payload, self.contextTimeout, "Bearer "+tokenChatbot.AccessToken)
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
		return
	}

	if statusCode != http.StatusOK {
		err = errors.Wrap(err, "[usecase.IncomingMessages]")
		return
	}

	var resChatBot api.ResponseChatbotCoster
	err = json.Unmarshal(res, &resChatBot)
	if err != nil {
		err = errors.Wrap(err, "[usecase.IncomingMessages] Unmarshal")
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

func (self *chatUcase) FindToken() (data model.Token, err error) {

	data, err = self.m.FindToken()
	if err != nil {
		err = errors.Wrap(err, "[usecase.FindToken]")
	}
	return
}

func (self *chatUcase) SetToken(updated map[string]interface{}) (err error) {

	err = self.m.SetToken(updated)
	if err != nil {
		err = errors.Wrap(err, "[usecase.SetToken]")
	}
	return
}

func (self *chatUcase) RefreshToken() (res []byte, statusCode int, err error) {

	url := self.urlSendMsg + "/wa/users/login"
	credential := viper.GetString("chatbot.username") + ":" + viper.GetString("chatbot.password")
	token := base64.StdEncoding.EncodeToString([]byte(credential))

	res, statusCode, err = httprequest.PostJson(url, nil, self.contextTimeout, "Basic "+token)
	if err != nil {
		err = errors.Wrap(err, "[usecase.RefreshToken]")
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
