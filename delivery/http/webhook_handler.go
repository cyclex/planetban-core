package http

import (
	"net/http"
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
	Gw domain.OrdersUcase
	Ch domain.ChatUcase
}

func NewOrderHandler(e *echo.Echo, ordersUcase domain.OrdersUcase, chatUcase domain.ChatUcase, debug bool) {

	appLog = pkg.New("app", debug)
	cmsLog = pkg.New("cms", debug)

	handler := &OrderHandler{
		Gw: ordersUcase,
		Ch: chatUcase,
	}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: pkg.New("middleware", debug).Out,
	}))

	e.POST("/v1/webhooks/whatsapp", handler.webhooksWhatsapp)
}

func (self *OrderHandler) webhooksWhatsapp(c echo.Context) (err error) {

	var (
		request api.ResSendMessage
		code    = 200
		ctx     = c.Request().Context()
	)

	c.Bind(&request)
	if len(request.Messages) > 0 {

		if request.Messages[0].Type == "text" {
			err = self.Gw.CreateQueueChat(ctx, "incoming", request)
			if err != nil {
				code = http.StatusInternalServerError
				appLog.Error(err)
			}
		} else {
			_, err = self.Ch.IncomingMessages(request.Messages[0])
			if err != nil {
				appLog.Error(err)
			}
		}
	}

	res := api.ResponseChatbot{
		Code:       code,
		Message:    http.StatusText(code),
		ServerTime: time.Now().Local().Unix(),
	}

	return c.JSON(code, res)
}
