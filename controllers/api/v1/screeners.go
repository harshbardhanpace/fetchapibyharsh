package v1

import (
	"encoding/json"
	"strings"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var theGainerLoserProvider models.GainerLoserProvider

func InitGainerLoserProvider(provider models.GainerLoserProvider) {
	defer models.HandlePanic()
	theGainerLoserProvider = provider
}

// GainerLoser
// @Tags space Screeners V1
// @Description Get Top Gainer Loser
// @Param ClientId header string false "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GainersLosersMostActiveVolumeReq true "Gainers Losers"
// @Success 200 {object} apihelpers.APIRes{data=models.TopGainerLoserResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/screeners/gainersloser [POST]
func GainerLoser(c *gin.Context) {
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	var reqParams models.GainersLosersMostActiveVolumeReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GainerLoser (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("GainerLoser (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}
	loggerconfig.Info("GainerLoser (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	var code int
	var resp apihelpers.APIRes
	if strings.ToLower(reqParams.Index) == constants.NSE {
		code, resp = theGainerLoserProvider.GainerLoserNse(requestH)
	} else {
		code, resp = theGainerLoserProvider.GainerLoserNiftyFifty(reqParams, requestH)
	}

	logDetail := "clientId: " + requestH.ClientId + " function: GainerLoser requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// MostActiveVolume
// @Tags space Screeners V1
// @Description Get Most Active Volume
// @Param ClientId header string false "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GainersLosersMostActiveVolumeReq true "Most Active Volume"
// @Success 200 {object} apihelpers.APIRes{data=models.MostActiveVolumeData}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/screeners/mostActiveVolume [POST]
func MostActiveVolume(c *gin.Context) {
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	var reqParams models.GainersLosersMostActiveVolumeReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("MostActiveVolume (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("MostActiveVolume (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}
	loggerconfig.Info("MostActiveVolume (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	var code int
	var resp apihelpers.APIRes
	if strings.ToLower(reqParams.Index) == constants.NSE {
		code, resp = theGainerLoserProvider.MostActiveVolumeNSE(requestH)
	} else {
		code, resp = theGainerLoserProvider.MostActiveVolumeDataNifty50(reqParams, requestH)
	}

	logDetail := "clientId: " + requestH.ClientId + " function: MostActiveVolume requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GainerLoserNiftyFifty
// @Tags space Screeners V1
// @Description Get Most Active Volume NSE
// @Param ClientId header string false "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.ChartDataReq true "Chart Data"
// @Success 200 {object} apihelpers.APIRes{data=models.ChartDataResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/screeners/chartData [POST]
func ChartData(c *gin.Context) {
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	var reqParams models.ChartDataReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ChartData (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("ChartData (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}
	loggerconfig.Info("ChartData (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theGainerLoserProvider.ChartData(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: GainerLoserNiftyFifty requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ReturnOnInvestment
// @Tags space Screeners V1
// @Description Return On Investment
// @Param ClientId header string false "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.ReturnOnInvestmentReq true "Return On Investment"
// @Success 200 {object} apihelpers.APIRes{data=models.ReturnOnInvestmentRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/screeners/returnOnInvestment [POST]
func ReturnOnInvestment(c *gin.Context) {
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	var reqParams models.ReturnOnInvestmentReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ReturnOnInvestment (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("ReturnOnInvestment (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("ReturnOnInvestment (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("ReturnOnInvestment (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theGainerLoserProvider.ReturnOnInvestment(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: ReturnOnInvestment requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchHistoricPerformance
// @Tags space Screeners V1
// @Description Fetch historic performace stockwise. 1D, 1W, 1M, 6M, 1Y, 3Y, 5Y
// @Param ClientId header string false "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.HistoricPerformaceReq true "Historic-Performace-Req"
// @Success 200 {object} apihelpers.APIRes{data=models.ChartDataResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/screeners/fetchHistoricPerformance [POST]
func FetchHistoricPerformance(c *gin.Context) {
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	var reqParams models.HistoricPerformaceReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FetchHistoricPerformance (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	//validate Period(1D, 1W, 1M, 6M, 1Y, 3Y, 5Y)
	if reqParams.Period != "1D" && reqParams.Period != "1W" && reqParams.Period != "1M" && reqParams.Period != "6M" && reqParams.Period != "1Y" && reqParams.Period != "3Y" && reqParams.Period != "5Y" {
		loggerconfig.Error("FetchHistoricPerformance (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchHistoricPerformance (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}
	loggerconfig.Info("FetchHistoricPerformance (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theGainerLoserProvider.FetchHistoricPerformance(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: FetchHistoricPerformance requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchAllHistoricPerformance
// @Tags space Screeners V1
// @Description Fetch All historic performace stockwise.
// @Param ClientId header string false "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.AllHistoricPerformaceReq true "All-Historic-Performace-Req"
// @Success 200 {object} apihelpers.APIRes{data=models.ChartDataResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/screeners/fetchHistoricPerformance/all [POST]
func FetchAllHistoricPerformance(c *gin.Context) {
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	var reqParams models.AllHistoricPerformaceReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FetchALLHistoricPerformance (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchALLHistoricPerformance (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("FetchALLHistoricPerformance (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theGainerLoserProvider.FetchAllHistoricPerformance(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: FetchALLHistoricPerformance requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
