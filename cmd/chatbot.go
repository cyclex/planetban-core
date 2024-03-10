package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/cyclex/planet-ban/api"
	"github.com/cyclex/planet-ban/domain"
	"github.com/cyclex/planet-ban/domain/model"
)

var (
	incoming = "incoming"
	outgoing = "outgoing"
)

func SendOrders(processing *bool, orderUcase domain.OrdersUcase, chatUcase domain.ChatUcase, ctx context.Context) {

	*processing = true
	resChannel := make(chan domain.CronChatbot)

	defer func() {
		close(resChannel)
		*processing = false
	}()

	queue, err := orderUcase.GetQueueChat(ctx, incoming)
	if err != nil {
		appLog.Error(err)
		return
	}

	var wg sync.WaitGroup

	doSend := 0
	for _, row := range queue {
		wg.Add(1)
		go workerSendOrders(ctx, row, chatUcase, resChannel, &wg)
		doSend++
	}

	wg.Wait()

	for i := 0; i < doSend; i++ {
		res := <-resChannel
		appLog.Debugln(fmt.Sprintf("Receive data ::: TrxID:%s Response:%+v Err:%s", res.ID, res, res.Err))

		if res.Status == "ok" {
			err = orderUcase.UpdateQueueChat(ctx, incoming, res.ID)
			if err != nil {
				appLog.Error(err)
				return
			}
		}
	}

	*processing = false

	return
}

func workerSendOrders(ctx context.Context, row model.QueueChat, chatUcase domain.ChatUcase, resChannel chan domain.CronChatbot, wg *sync.WaitGroup) {

	var resTrx domain.CronChatbot

	appLog.Debugln(fmt.Sprintf("Send data ::: TrxID:%s", row.ID))

	handleReturn := func() {
		wg.Done()

		if r := recover(); r != nil {
			appLog.Infoln(fmt.Sprintf(" -- workerSendOrders Recovered from panic: %s", r))
		}

		resChannel <- resTrx
	}

	defer handleReturn()

	resTrx = domain.CronChatbot{
		ID:         row.ID,
		ServerTime: time.Now().Local(),
	}

	trxChatbotID, err := chatUcase.IncomingMessages(row.Messages.Messages[0])
	if err != nil {
		appLog.Error(err)
		resTrx.Err = err
		return
	}

	resTrx.Status = "ok"
	resTrx.TrxChatBotMsgID = trxChatbotID

	return

}

func RefreshToken(processing *bool, chatUcase domain.ChatUcase, ctx context.Context) {

	*processing = true

	defer func() {
		*processing = false
	}()

	res, statusCode, err := chatUcase.RefreshToken()
	authLog.Infoln(fmt.Sprintf("Res:%s statusCode:%d Err:%s", res, statusCode, err))

	var token = api.ResLogin{}
	err = json.Unmarshal(res, &token)
	if err != nil {
		authLog.Error(err)
		return
	}

	if statusCode == http.StatusOK {
		updated := map[string]interface{}{
			"access_token": token.Users[0].Token,
			"expired_at":   token.Users[0].ExpiresAfter,
			"updated_at":   time.Now().Local(),
		}

		err = chatUcase.SetToken(updated)
		if err != nil {
			authLog.Error(err)
			return
		}
	}

	*processing = false

}
