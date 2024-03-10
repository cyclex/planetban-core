package main

import (
	"context"
	"time"

	"github.com/cyclex/planet-ban/domain"
	"github.com/robfig/cron"
)

var processing, processingAuth, processingRedeem, processingCalculateQuota, processingSendReply, processingRecurringQuota bool

func InitCron(orderUcase domain.OrdersUcase, chatUcase domain.ChatUcase, cmsUcase domain.CmsUcase, timeout time.Duration) {

	c := context.Background()
	_, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	cr := cron.New()

	cr.AddFunc("* * * * * *", func() {
		if !processing {
			SendOrders(&processing, orderUcase, chatUcase, c)
		}
	})

	cr.AddFunc("@weekly", func() {
		if !processingAuth {
			RefreshToken(&processingAuth, chatUcase, c)
		}
	})

	cr.Start()

}
