package v2

import (
	"strconv"
	"strings"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
)

var theGainerLoserProviderV2 models.GainerLoserProvider

func InitGainerLoserProviderV2(provider models.GainerLoserProvider) {
	defer models.HandlePanic()
	theGainerLoserProviderV2 = provider
}

// GainerLoser
// @Tags space Screeners V2
// @Description Get Top Gainer Loser V2
// @Param ClientId header string false "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param index query string true "index Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.TopGainerLoserResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/screeners/gainersloser [GET]
func GainerLoser(c *gin.Context) {
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	var reqParams models.GainersLosersMostActiveVolumeReq
	index := c.Query("index")
	if index == "" {
		loggerconfig.Error("GainerLoser V2 (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.Index = index

	if requestH.DeviceType == "" {
		loggerconfig.Error("GainerLoser V2(controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("GainerLoser V2(controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "clientId: ", requestH.ClientId)
	var code int
	var resp apihelpers.APIRes
	if strings.ToLower(reqParams.Index) == constants.NSE {
		code, resp = theGainerLoserProviderV2.GainerLoserNse(requestH)
	} else {
		code, resp = theGainerLoserProviderV2.GainerLoserNiftyFifty(reqParams, requestH)
	}

	logDetail := "clientId: " + requestH.ClientId + " function: GainerLoser V2 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// MostActiveVolume
// @Tags space Screeners V2
// @Description Get Most Active Volume V2
// @Param ClientId header string false "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param index query string true "index Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.MostActiveVolumeData}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/screeners/mostActiveVolume [GET]
func MostActiveVolume(c *gin.Context) {
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	var reqParams models.GainersLosersMostActiveVolumeReq
	index := c.Query("index")
	if index == "" {
		loggerconfig.Error("MostActiveVolume V2 (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.Index = index

	if requestH.DeviceType == "" {
		loggerconfig.Error("MostActiveVolume V2(controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}
	loggerconfig.Info("MostActiveVolume V2(controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "clientId: ", requestH.ClientId)
	var code int
	var resp apihelpers.APIRes
	if strings.ToLower(reqParams.Index) == constants.NSE {
		code, resp = theGainerLoserProviderV2.MostActiveVolumeNSE(requestH)
	} else {
		code, resp = theGainerLoserProviderV2.MostActiveVolumeDataNifty50(reqParams, requestH)
	}

	logDetail := "clientId: " + requestH.ClientId + " function:MostActiveVolume V2 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GainerLoserNiftyFifty
// @Tags space Screeners V2
// @Description Get Most Active Volume NSE V2
// @Param ClientId header string false "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param exchange query string true "exchange Query Parameter" dataType(string)
// @Param token query string true "token Query Parameter" dataType(string)
// @Param candleType query string true "candleType Query Parameter" dataType(string)
// @Param startTime query string true "startTime Query Parameter" dataType(string)
// @Param endTime query string true "endTime Query Parameter" dataType(string)
// @Param dataDuration query string true "dataDuration Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.ChartDataResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/screeners/chartData [GET]
func ChartData(c *gin.Context) {
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	var reqParams models.ChartDataReq
	exchange := c.Query("exchange")
	token := c.Query("token")
	candleType := c.Query("candleType")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	dataDuration := c.Query("dataDuration")
	if exchange == "" || token == "" || candleType == "" || startTime == "" || endTime == "" || dataDuration == "" {
		loggerconfig.Error("ChartData V2 (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.Exchange = exchange
	reqParams.Token = token
	reqParams.CandleType = candleType
	reqParams.StartTime = startTime
	reqParams.EndTime = endTime
	reqParams.DataDuration = dataDuration

	if requestH.DeviceType == "" {
		loggerconfig.Error("ChartData V2(controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}
	loggerconfig.Info("ChartData V2(controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "clientId: ", requestH.ClientId)
	code, resp := theGainerLoserProviderV2.ChartData(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: ChartData requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ReturnOnInvestment
// @Tags space Screeners V2
// @Description Return On Investment V2
// @Param ClientId header string false "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param days query string true "days Query Parameter" dataType(string)
// @Param index query string true "index Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.ReturnOnInvestmentRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/screeners/returnOnInvestment [GET]
func ReturnOnInvestment(c *gin.Context) {
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	var reqParams models.ReturnOnInvestmentReq
	index := c.Query("index")
	daysStr := c.Query("days")
	days, errDays := strconv.Atoi(daysStr)
	if index == "" || errDays != nil || days <= 0 {
		loggerconfig.Error("ReturnOnInvestment V2 (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.Index = index
	reqParams.Days = days

	if requestH.DeviceType == "" {
		loggerconfig.Error("ReturnOnInvestment V2(controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("ReturnOnInvestment V2(controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "clientId: ", requestH.ClientId)
	code, resp := theGainerLoserProviderV2.ReturnOnInvestment(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: ReturnOnInvestment requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
