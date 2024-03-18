package http

import (
	"net/http"
	"strconv"

	"github.com/cyclex/planet-ban/api"
	_AppMW "github.com/cyclex/planet-ban/delivery/http/middleware"
	"github.com/cyclex/planet-ban/domain"
	"github.com/cyclex/planet-ban/pkg"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

var cmsLog *logrus.Logger

type CmsHandler struct {
	CmsGw domain.CmsUcase
}

func NewCmsHandler(e *echo.Echo, gw domain.CmsUcase, debug bool) {

	cmsLog = pkg.New("cms", debug)

	handler := &CmsHandler{
		CmsGw: gw,
	}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: pkg.New("middleware", false).Out,
	}))

	e.POST("/v1/login", handler.login, _AppMW.ReqLogin)
	e.GET("/v1/checkToken/:token", handler.checkToken)

	e.POST("/v1/report/:type", handler.report)

	e.POST("/v1/campaign", handler.createCampaign, _AppMW.ReqCampaign)
	e.DELETE("/v1/campaign/:id", handler.deleteCampaign)
	e.PUT("/v1/campaign/:id", handler.setCampaign)

	e.POST("/v1/kol", handler.createKol)
	e.DELETE("/v1/kol/:id", handler.deleteKol)
	e.PUT("/v1/kol/:id", handler.setKol)
}

func (self *CmsHandler) login(c echo.Context) error {

	var (
		request api.Login
		res     interface{}
		code    = http.StatusForbidden
		ctx     = c.Request().Context()
	)

	c.Bind(&request)
	data, err := self.CmsGw.Login(ctx, request)
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseSuccess{
			Status:  true,
			Message: "success",
			Data:    data,
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) checkToken(c echo.Context) error {

	var (
		request api.CheckToken
		res     interface{}
		code    = http.StatusForbidden
		ctx     = c.Request().Context()
	)

	request.Token = c.Param("token")
	err := self.CmsGw.CheckToken(ctx, request)
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseError{
			Status: true,
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) report(c echo.Context) error {

	var (
		request api.Report
		res     interface{}
		code    = http.StatusInternalServerError
		ctx     = c.Request().Context()
	)

	c.Bind(&request)
	data, err := self.CmsGw.Report(ctx, request, c.Param("type"))
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseSuccess{
			Status:  true,
			Message: "success",
			Data:    data,
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) createCampaign(c echo.Context) error {

	var (
		request api.Campaign
		res     interface{}
		code    = http.StatusInternalServerError
		ctx     = c.Request().Context()
	)

	c.Bind(&request)
	err := self.CmsGw.CreateCampaign(ctx, request)
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseSuccess{
			Status:  true,
			Message: "success",
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) setCampaign(c echo.Context) error {

	var (
		request api.Campaign
		res     interface{}
		ctx     = c.Request().Context()
		code    = http.StatusInternalServerError
	)

	c.Bind(&request)
	n, _ := strconv.Atoi(c.Param("id"))
	request.CampaignID = int64(n)

	err := self.CmsGw.SetCampaign(ctx, request)
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseSuccess{
			Status:  true,
			Message: "success",
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) deleteCampaign(c echo.Context) error {

	var (
		res  interface{}
		code = http.StatusInternalServerError
		ctx  = c.Request().Context()
	)

	n, _ := strconv.Atoi(c.Param("id"))
	x := int64(n)

	err := self.CmsGw.DeleteCampaign(ctx, []int64{x})
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseSuccess{
			Status:  true,
			Message: "success",
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) createKol(c echo.Context) error {

	var (
		request api.Kol
		res     interface{}
		code    = http.StatusInternalServerError
		ctx     = c.Request().Context()
	)

	c.Bind(&request)
	err := self.CmsGw.CreateKol(ctx, request)
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseSuccess{
			Status:  true,
			Message: "success",
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) deleteKol(c echo.Context) error {

	var (
		res  interface{}
		code = http.StatusInternalServerError
		ctx  = c.Request().Context()
	)

	n, _ := strconv.Atoi(c.Param("id"))
	x := int64(n)

	err := self.CmsGw.DeleteKol(ctx, x)
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseSuccess{
			Status:  true,
			Message: "success",
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) setKol(c echo.Context) error {

	var (
		request api.Kol
		res     interface{}
		ctx     = c.Request().Context()
		code    = http.StatusInternalServerError
	)

	c.Bind(&request)
	n, _ := strconv.Atoi(c.Param("id"))
	request.KolID = int64(n)

	err := self.CmsGw.SetKol(ctx, request)
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseSuccess{
			Status:  true,
			Message: "success",
		}
	}

	return c.JSON(code, res)
}
