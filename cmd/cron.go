package main

import (
	"context"
	"time"

	"github.com/cyclex/planet-ban/domain"
	"github.com/robfig/cron"
)

var processing, processingAuth, processingRedeem, processingCalculateQuota, processingSendReply, processingRecurringQuota bool

func InitCron(chatUcase domain.ChatUcase, timeout time.Duration) {

	c := context.Background()
	_, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	cr := cron.New()

	cr.AddFunc("@weekly", func() {
		if !processingAuth {
			RefreshToken(&processingAuth, chatUcase, c)
		}
	})

	cr.Start()

}
