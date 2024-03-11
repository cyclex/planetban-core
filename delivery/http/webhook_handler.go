package http

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/cyclex/planet-ban/api"
	"github.com/cyclex/planet-ban/domain"
	"github.com/cyclex/planet-ban/pkg"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

var appLog, cmsLog *logrus.Logger

type OrderHandler struct {
	Ch domain.ChatUcase
}

func NewOrderHandler(e *echo.Echo, chatUcase domain.ChatUcase, debug bool) {

	appLog = pkg.New("app", debug)
	cmsLog = pkg.New("cms", debug)

	handler := &OrderHandler{
		Ch: chatUcase,
	}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: pkg.New("middleware", debug).Out,
	}))

	e.POST("/v1/webhooks/whatsapp", handler.webhooksWhatsapp)
	e.GET("go/:qp", handler.webhooksInfluencer)
}

func (self *OrderHandler) webhooksWhatsapp(c echo.Context) (err error) {

	var (
		request api.ResSendMessage
		code    = 400
	)

	defer func(code *int) {
		res := api.ResponseChatbot{
			Code:       *code,
			Message:    http.StatusText(*code),
			ServerTime: time.Now().Local().Unix(),
		}
		c.JSON(*code, res)
	}(&code)

	c.Bind(&request)

	if len(request.Messages) > 0 {

		if request.Messages[0].From == "" {
			return
		}

		code = 200
		_, err = self.Ch.IncomingMessages(request.Messages[0])
		if err != nil {
			appLog.Error(err)
		}
	}

	return
}

func (self *OrderHandler) webhooksInfluencer(c echo.Context) (err error) {

	id := c.Param("qp")
	if id == "" {
		return
	}

	message, err := self.Ch.GetWhatsappTemplateMessage(id)
	if err != nil {
		appLog.Error(err)
		return
	}

	url := fmt.Sprintf("https://wa.me/%s?text=%s", self.Ch.GetWabaAccountNumber(), url.QueryEscape(message))
	err = c.Redirect(http.StatusSeeOther, url)
	if err != nil {
		appLog.Error(err)
	}

	return
}
