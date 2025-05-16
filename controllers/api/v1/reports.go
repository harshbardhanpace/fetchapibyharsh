package v1

import (
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var theReportsProvider models.ReportsProvider

func InitReportsProvider(provider models.ReportsProvider) {
	defer models.HandlePanic()
	theReportsProvider = provider
}

// ViewDPCharges
// @Tags space Reports V1
// @Description DP Charges
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param dfdatefr query string true "dfdatefr Query Parameter" dataType(string)
// @Param dfdateto query string true "dfdateto Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.ViewDPchargesRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/viewDPCharges [GET]
func ViewDPCharges(c *gin.Context) {
	var reqParams models.DPChargesReq
	dfDateFr := c.Query("dfdatefr")
	dfDateTo := c.Query("dfdateto")

	if dfDateFr == "" || dfDateTo == "" {
		loggerconfig.Error("ViewDPCharges (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if !helpers.ValidateDateQueryParam(dfDateFr, constants.DDMMYYYY) || !helpers.ValidateDateQueryParam(dfDateTo, constants.DDMMYYYY) || !helpers.ValidateDateRange(dfDateFr, dfDateTo, constants.DDMMYYYY) {
		loggerconfig.Error("ViewDPCharges (controller), Invalid date format !")
		apihelpers.ErrorMessage(c, constants.InvalidDate)
		return
	}

	dfDateFr, err := helpers.ConvertDateFormat(dfDateFr, constants.DDMMYYYY, constants.ShilpiDateFormat)
	if err != nil {
		loggerconfig.Error("ViewDPCharges (controller), Invalid date format !, unable to convert inputDate to ShilpiDateFormat, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidDate)
		return
	}

	dfDateTo, err = helpers.ConvertDateFormat(dfDateTo, constants.DDMMYYYY, constants.ShilpiDateFormat)
	if err != nil {
		loggerconfig.Error("ViewDPCharges (controller), Invalid date format !, unable to convert inputDate to ShilpiDateFormat, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidDate)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId
	reqParams.DFDateFr = dfDateFr
	reqParams.DFDateTo = dfDateTo

	if requestH.DeviceType == "" {
		loggerconfig.Error("ViewDPCharges (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("ViewDPCharges (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("ViewDPCharges (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.ViewDPCharges(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: ViewDPCharges requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DownloadDPCharges
// @Tags space Reports V1
// @Description DP Charges
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param dfdatefr query string true "dfdatefr Query Parameter" dataType(string)
// @Param dfdateto query string true "dfdateto Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.DownloadDPChargesRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/downloadDPCharges [GET]
func DownloadDPCharges(c *gin.Context) {
	var reqParams models.DPChargesReq
	dfDateFr := c.Query("dfdatefr")
	dfDateTo := c.Query("dfdateto")

	if dfDateFr == "" || dfDateTo == "" {
		loggerconfig.Error("DownloadDPCharges (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if !helpers.ValidateDateQueryParam(dfDateFr, constants.DDMMYYYY) || !helpers.ValidateDateQueryParam(dfDateTo, constants.DDMMYYYY) || !helpers.ValidateDateRange(dfDateFr, dfDateTo, constants.DDMMYYYY) {
		loggerconfig.Error("DownloadDPCharges (controller), Invalid date format !")
		apihelpers.ErrorMessage(c, constants.InvalidDate)
		return
	}

	dfDateFr, err := helpers.ConvertDateFormat(dfDateFr, constants.DDMMYYYY, constants.ShilpiDateFormat)
	if err != nil {
		loggerconfig.Error("ViewDPCharges (controller), Invalid date format !, unable to convert inputDate dfDateFr to ShilpiDateFormat, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidDate)
		return
	}

	dfDateTo, err = helpers.ConvertDateFormat(dfDateTo, constants.DDMMYYYY, constants.ShilpiDateFormat)
	if err != nil {
		loggerconfig.Error("ViewDPCharges (controller), Invalid date format !, unable to convert inputDate dfDateTo to ShilpiDateFormat, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidDate)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId
	reqParams.DFDateFr = dfDateFr
	reqParams.DFDateTo = dfDateTo

	if requestH.DeviceType == "" {
		loggerconfig.Error("DownloadDPCharges (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("DownloadDPCharges (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("DownloadDPCharges (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.DownloadDPCharges(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: DownloadDPCharges requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// SendEmailDPCharges
// @Tags space Reports V1
// @Description Send Email DP Charges
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param dfdatefr query string true "dfdatefr Query Parameter" dataType(string)
// @Param dfdateto query string true "dfdateto Query Parameter" dataType(string)
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/sendEmailDPCharges [GET]
func SendEmailDPCharges(c *gin.Context) {
	var reqParams models.DPChargesReq
	dfDateFr := c.Query("dfdatefr")
	dfDateTo := c.Query("dfdateto")

	if dfDateFr == "" || dfDateTo == "" {
		loggerconfig.Error("SendEmailDPCharges (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if !helpers.ValidateDateQueryParam(dfDateFr, constants.DDMMYYYY) || !helpers.ValidateDateQueryParam(dfDateTo, constants.DDMMYYYY) || !helpers.ValidateDateRange(dfDateFr, dfDateTo, constants.DDMMYYYY) {
		loggerconfig.Error("SendEmailDPCharges (controller), Invalid date format !")
		apihelpers.ErrorMessage(c, constants.InvalidDate)
		return
	}

	dfDateFr, err := helpers.ConvertDateFormat(dfDateFr, constants.DDMMYYYY, constants.ShilpiDateFormat)
	if err != nil {
		loggerconfig.Error("SendEmailDPCharges (controller), Invalid date format !, unable to convert inputDate dfDateFr to ShilpiDateFormat, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidDate)
		return
	}

	dfDateTo, err = helpers.ConvertDateFormat(dfDateTo, constants.DDMMYYYY, constants.ShilpiDateFormat)
	if err != nil {
		loggerconfig.Error("SendEmailDPCharges (controller), Invalid date format !, unable to convert inputDate dfDateTo to ShilpiDateFormat, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidDate)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId
	reqParams.DFDateFr = dfDateFr
	reqParams.DFDateTo = dfDateTo

	if requestH.DeviceType == "" {
		loggerconfig.Error("SendEmailDPCharges (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("SendEmailDPCharges (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("SendEmailDPCharges (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.SendEmailDPCharges(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: SendEmailDPCharges requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ViewTradebook
// @Tags space Reports V1
// @Description View Tradebook
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param dfdatefr query string true "dfdatefr Query Parameter" dataType(string)
// @Param dfdateto query string true "dfdateto Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.ScripWiseCostingRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/viewTradebook [GET]
func ViewTradebook(c *gin.Context) {
	var reqParams models.TradebookReq

	dfDateFr := c.Query("dfdatefr")
	dfDateTo := c.Query("dfdateto")

	if dfDateFr == "" || dfDateTo == "" {
		loggerconfig.Error("ViewTradebook (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if !helpers.ValidateDateQueryParam(dfDateFr, constants.DDMMYYYY) || !helpers.ValidateDateQueryParam(dfDateTo, constants.DDMMYYYY) || !helpers.ValidateDateRange(dfDateFr, dfDateTo, constants.DDMMYYYY) {
		loggerconfig.Error("ViewTradebook (controller), Invalid date format !")
		apihelpers.ErrorMessage(c, constants.InvalidDate)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId
	reqParams.DFDateFr = dfDateFr
	reqParams.DFDateTo = dfDateTo

	if requestH.DeviceType == "" {
		loggerconfig.Error("ViewTradebook (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("ViewTradebook (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("ViewTradebook (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.ViewTradebook(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: ViewTradebook requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DownloadTradebook
// @Tags space Reports V1
// @Description Download Tradebook
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param dfdatefr query string true "dfdatefr Query Parameter" dataType(string)
// @Param dfdateto query string true "dfdateto Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.DownloadTradebookRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/downloadTradebook [GET]
func DownloadTradebook(c *gin.Context) {
	var reqParams models.TradebookReq

	dfDateFr := c.Query("dfdatefr")
	dfDateTo := c.Query("dfdateto")

	if dfDateFr == "" || dfDateTo == "" {
		loggerconfig.Error("DownloadTradebook (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if !helpers.ValidateDateQueryParam(dfDateFr, constants.DDMMYYYY) || !helpers.ValidateDateQueryParam(dfDateTo, constants.DDMMYYYY) || !helpers.ValidateDateRange(dfDateFr, dfDateTo, constants.DDMMYYYY) {
		loggerconfig.Error("DownloadTradebook (controller), Invalid date format !")
		apihelpers.ErrorMessage(c, constants.InvalidDate)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId
	reqParams.DFDateFr = dfDateFr
	reqParams.DFDateTo = dfDateTo

	if requestH.DeviceType == "" {
		loggerconfig.Error("DownloadTradebook (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("DownloadTradebook (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("DownloadTradebook (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.DownloadTradebook(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: DownloadTradebook requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ViewLedger
// @Tags space Reports V1
// @Description View Ledger
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param dfdatefr query string true "dfdatefr Query Parameter" dataType(string)
// @Param dfdateto query string true "dfdateto Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.FinancialLedgerRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/viewLedger [GET]
func ViewLedger(c *gin.Context) {
	var reqParams models.LedgerReq

	dfDateFr := c.Query("dfdatefr")
	dfDateTo := c.Query("dfdateto")

	if dfDateFr == "" || dfDateTo == "" {
		loggerconfig.Error("ViewLedger (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if !helpers.ValidateDateQueryParam(dfDateFr, constants.DDMMYYYY) || !helpers.ValidateDateQueryParam(dfDateTo, constants.DDMMYYYY) || !helpers.ValidateDateRange(dfDateFr, dfDateTo, constants.DDMMYYYY) {
		loggerconfig.Error("ViewLedger (controller), Invalid date format !")
		apihelpers.ErrorMessage(c, constants.InvalidDate)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId
	reqParams.DFDateFr = dfDateFr
	reqParams.DFDateTo = dfDateTo

	if requestH.DeviceType == "" {
		loggerconfig.Error("ViewLedger (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("ViewLedger (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("ViewLedger (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.ViewLedger(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: ViewLedger requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DownloadLedger
// @Tags space Reports V1
// @Description Download Ledger
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param dfdatefr query string true "dfdatefr Query Parameter" dataType(string)
// @Param dfdateto query string true "dfdateto Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.DownloadLedgerRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/downloadLedger [GET]
func DownloadLedger(c *gin.Context) {
	var reqParams models.LedgerReq

	dfDateFr := c.Query("dfdatefr")
	dfDateTo := c.Query("dfdateto")

	if dfDateFr == "" || dfDateTo == "" {
		loggerconfig.Error("DownloadLedger (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if !helpers.ValidateDateQueryParam(dfDateFr, constants.DDMMYYYY) || !helpers.ValidateDateQueryParam(dfDateTo, constants.DDMMYYYY) || !helpers.ValidateDateRange(dfDateFr, dfDateTo, constants.DDMMYYYY) {
		loggerconfig.Error("DownloadLedger (controller), Invalid date format !")
		apihelpers.ErrorMessage(c, constants.InvalidDate)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId
	reqParams.DFDateFr = dfDateFr
	reqParams.DFDateTo = dfDateTo

	if requestH.DeviceType == "" {
		loggerconfig.Error("DownloadLedger (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("DownloadLedger (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("DownloadLedger (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.DownloadLedger(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: DownloadLedger requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ViewOpenPosition
// @Tags space Reports V1
// @Description View Open Position
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param dsFlag query string true "dsFlag Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.OpenPositionRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/viewOpenPosition [GET]
func ViewOpenPosition(c *gin.Context) {
	var reqParams models.OpenPositionReq
	dsFlag := c.Query("dsFlag")

	if dsFlag == "" {
		loggerconfig.Error("ViewOpenPosition (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId
	reqParams.Dsflag = dsFlag

	if requestH.DeviceType == "" {
		loggerconfig.Error("ViewOpenPosition (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("ViewOpenPosition (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("ViewOpenPosition (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.ViewOpenPosition(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: ViewOpenPosition requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DownloadOpenPosition
// @Tags space Reports V1
// @Description Download Open Position
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param dsFlag query string true "dsFlag Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.DownloadOpenPositionRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/downloadOpenPosition [GET]
func DownloadOpenPosition(c *gin.Context) {
	var reqParams models.OpenPositionReq

	dsFlag := c.Query("dsFlag")

	if dsFlag == "" {
		loggerconfig.Error("DownloadOpenPosition (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId
	reqParams.Dsflag = dsFlag

	if requestH.DeviceType == "" {
		loggerconfig.Error("DownloadOpenPosition (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("DownloadOpenPosition (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("DownloadOpenPosition (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.DownloadOpenPosition(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: DownloadOpenPosition requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ViewFnoPnl
// @Tags space Reports V1
// @Description View Fno Pnl
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param dfdatefr query string true "dfdatefr Query Parameter" dataType(string)
// @Param dfdateto query string true "dfdateto Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.FONetPositionRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/viewFnoPnl [GET]
func ViewFnoPnl(c *gin.Context) {
	var reqParams models.FnoPnlReq

	dfDateFr := c.Query("dfdatefr")
	dfDateTo := c.Query("dfdateto")

	if dfDateFr == "" || dfDateTo == "" {
		loggerconfig.Error("ViewFnoPnl (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if !helpers.ValidateDateQueryParam(dfDateFr, constants.DDMMYYYY) || !helpers.ValidateDateQueryParam(dfDateTo, constants.DDMMYYYY) || !helpers.ValidateDateRange(dfDateFr, dfDateTo, constants.DDMMYYYY) {
		loggerconfig.Error("ViewFnoPnl (controller), Invalid date format !")
		apihelpers.ErrorMessage(c, constants.InvalidDate)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId
	reqParams.DFDateFr = dfDateFr
	reqParams.DFDateTo = dfDateTo

	if requestH.DeviceType == "" {
		loggerconfig.Error("ViewFnoPnl (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("ViewFnoPnl (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("ViewFnoPnl (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.ViewFnoPnl(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: ViewFnoPnl requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DownloadFnoPnl
// @Tags space Reports V1
// @Description Download Fno Pnl
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param dfdatefr query string true "dfdatefr Query Parameter" dataType(string)
// @Param dfdateto query string true "dfdateto Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.DownloadFnoPnlRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/downloadFnoPnl [GET]
func DownloadFnoPnl(c *gin.Context) {
	var reqParams models.FnoPnlReq

	dfDateFr := c.Query("dfdatefr")
	dfDateTo := c.Query("dfdateto")

	if dfDateFr == "" || dfDateTo == "" {
		loggerconfig.Error("DownloadFnoPnl (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if !helpers.ValidateDateQueryParam(dfDateFr, constants.DDMMYYYY) || !helpers.ValidateDateQueryParam(dfDateTo, constants.DDMMYYYY) || !helpers.ValidateDateRange(dfDateFr, dfDateTo, constants.DDMMYYYY) {
		loggerconfig.Error("DownloadFnoPnl (controller), Invalid date format !")
		apihelpers.ErrorMessage(c, constants.InvalidDate)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId
	reqParams.DFDateFr = dfDateFr
	reqParams.DFDateTo = dfDateTo

	if requestH.DeviceType == "" {
		loggerconfig.Error("DownloadFnoPnl (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("DownloadFnoPnl (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("DownloadFnoPnl (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.DownloadFnoPnl(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: DownloadFnoPnl requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ViewHoldingFinancial
// @Tags space Reports V1
// @Description View Holding Financial
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Success 200 {object} apihelpers.APIRes{data=models.GetHoldingFinancialDataRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/viewHoldingFinancial [GET]
func ViewHoldingFinancial(c *gin.Context) {
	var reqParams models.GetHoldingFinancialDataReq

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId

	if requestH.DeviceType == "" {
		loggerconfig.Error("ViewHoldingFinancial (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("ViewHoldingFinancial (controller), requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.ViewHoldingFinancial(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: ViewHoldingFinancial requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DownloadHoldingFinancial
// @Tags space Reports V1
// @Description Download Holding Financial
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Success 200 {object} apihelpers.APIRes{data=models.DownloadHoldingFinancialRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/downloadHoldingFinancial [GET]
func DownloadHoldingFinancial(c *gin.Context) {
	var reqParams models.GetHoldingFinancialDataReq

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId

	if requestH.DeviceType == "" {
		loggerconfig.Error("DownloadHoldingFinancial (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("DownloadHoldingFinancial (controller), requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.DownloadHoldingFinancial(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: PendingOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// SendEmailHoldingFinancial
// @Tags space Reports V1
// @Description Send Email Holding Financial
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/sendEmailHoldingFinancial [GET]
func SendEmailHoldingFinancial(c *gin.Context) {
	var reqParams models.GetHoldingFinancialDataReq

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId

	if requestH.DeviceType == "" {
		loggerconfig.Error("SendEmailHoldingFinancial (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("SendEmailHoldingFinancial (controller), requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.SendEmailHoldingFinancial(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: SendEmailHoldingFinancial requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// SendEmailLedger
// @Tags space Reports V1
// @Description Send Email Ledger
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param dfdatefr query string true "dfdatefr Query Parameter" dataType(string)
// @Param dfdateto query string true "dfdateto Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/sendEmailLedger [GET]
func SendEmailLedger(c *gin.Context) {
	var reqParams models.LedgerReq

	dfDateFr := c.Query("dfdatefr")
	dfDateTo := c.Query("dfdateto")

	if dfDateFr == "" || dfDateTo == "" {
		loggerconfig.Error("SendEmailLedger (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if !helpers.ValidateDateQueryParam(dfDateFr, constants.DDMMYYYY) || !helpers.ValidateDateQueryParam(dfDateTo, constants.DDMMYYYY) || !helpers.ValidateDateRange(dfDateFr, dfDateTo, constants.DDMMYYYY) {
		loggerconfig.Error("SendEmailLedger (controller), Invalid date format !")
		apihelpers.ErrorMessage(c, constants.InvalidDate)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId
	reqParams.DFDateFr = dfDateFr
	reqParams.DFDateTo = dfDateTo

	if requestH.DeviceType == "" {
		loggerconfig.Error("SendEmailLedger (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("SendEmailLedger (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("SendEmailLedger (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.SendEmailLedger(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: SendEmailLedger requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ViewCommodityTradebook
// @Tags space Reports V1
// @Description View Commodity Tradebook
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param dfdatefr query string true "dfdatefr Query Parameter" dataType(string)
// @Param dfdateto query string true "dfdateto Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.CommodityTransactionRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/viewCommodityTradebook [GET]
func ViewCommodityTradebook(c *gin.Context) {
	var reqParams models.CommodityTradebookReq

	dfDateFr := c.Query("dfdatefr")
	dfDateTo := c.Query("dfdateto")

	if dfDateFr == "" || dfDateTo == "" {
		loggerconfig.Error("ViewCommodityTradebook (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId
	reqParams.DFDateFr = dfDateFr
	reqParams.DFDateTo = dfDateTo

	if requestH.DeviceType == "" {
		loggerconfig.Error("ViewCommodityTradebook (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("ViewCommodityTradebook (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("ViewCommodityTradebook (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.ViewCommodityTradebook(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: ViewCommodityTradebook requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DownloadCommodityTradebook
// @Tags space Reports V1
// @Description Download Commodity Tradebook
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param dfdatefr query string true "dfdatefr Query Parameter" dataType(string)
// @Param dfdateto query string true "dfdateto Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.DownloadCommodityTradebookRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/downloadCommodityTradebook [GET]
func DownloadCommodityTradebook(c *gin.Context) {
	var reqParams models.CommodityTradebookReq
	dfDateFr := c.Query("dfdatefr")
	dfDateTo := c.Query("dfdateto")

	if dfDateFr == "" || dfDateTo == "" {
		loggerconfig.Error("DownloadCommodityTradebook (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId
	reqParams.DFDateFr = dfDateFr
	reqParams.DFDateTo = dfDateTo

	if requestH.DeviceType == "" {
		loggerconfig.Error("DownloadCommodityTradebook (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("DownloadCommodityTradebook (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("DownloadCommodityTradebook (controller), requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.DownloadCommodityTradebook(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: DownloadCommodityTradebook requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// SendEmailCommodityTradebook
// @Tags space Reports V1
// @Description Send Email Commodity Tradebook
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param dfdatefr query string true "dfdatefr Query Parameter" dataType(string)
// @Param dfdateto query string true "dfdateto Query Parameter" dataType(string)
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/sendEmailCommodityTradebook [GET]
func SendEmailCommodityTradebook(c *gin.Context) {
	var reqParams models.CommodityTradebookReq
	dfDateFr := c.Query("dfdatefr")
	dfDateTo := c.Query("dfdateto")

	if dfDateFr == "" || dfDateTo == "" {
		loggerconfig.Error("SendEmailCommodityTradebook (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId
	reqParams.DFDateFr = dfDateFr
	reqParams.DFDateTo = dfDateTo

	if requestH.DeviceType == "" {
		loggerconfig.Error("SendEmailCommodityTradebook (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("SendEmailCommodityTradebook (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("SendEmailCommodityTradebook (controller), requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.SendEmailCommodityTradebook(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: SendEmailCommodityTradebook requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ViewFnoTradebook
// @Tags space Reports V1
// @Description View Fno Tradebook
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param dfdatefr query string true "dfdatefr Query Parameter" dataType(string)
// @Param dfdateto query string true "dfdateto Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.FNOTransactionRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/viewFnoTradebook [GET]
func ViewFnoTradebook(c *gin.Context) {
	var reqParams models.FNOTradebookReq

	dfDateFr := c.Query("dfdatefr")
	dfDateTo := c.Query("dfdateto")

	if dfDateFr == "" || dfDateTo == "" {
		loggerconfig.Error("ViewFnoTradebook (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId
	reqParams.DFDateFr = dfDateFr
	reqParams.DFDateTo = dfDateTo

	if requestH.DeviceType == "" {
		loggerconfig.Error("ViewFnoTradebook (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("ViewFnoTradebook (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("ViewFnoTradebook (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.ViewFnoTradebook(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: ViewFnoTradebook requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DownloadFnoTradebook
// @Tags space Reports V1
// @Description Download Fno Tradebook
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param dfdatefr query string true "dfdatefr Query Parameter" dataType(string)
// @Param dfdateto query string true "dfdateto Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.DownloadFnoTradebookRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/downloadFnoTradebook [GET]
func DownloadFnoTradebook(c *gin.Context) {
	var reqParams models.FNOTradebookReq

	dfDateFr := c.Query("dfdatefr")
	dfDateTo := c.Query("dfdateto")

	if dfDateFr == "" || dfDateTo == "" {
		loggerconfig.Error("DownloadFnoTradebook (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId
	reqParams.DFDateFr = dfDateFr
	reqParams.DFDateTo = dfDateTo

	if requestH.DeviceType == "" {
		loggerconfig.Error("DownloadFnoTradebook (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("DownloadFnoTradebook (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("DownloadFnoTradebook (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.DownloadFnoTradebook(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: DownloadFnoTradebook requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// SendEmailFnoTradebook
// @Tags space Reports V1
// @Description Send Email Fno Tradebook
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param dfdatefr query string true "dfdatefr Query Parameter" dataType(string)
// @Param dfdateto query string true "dfdateto Query Parameter" dataType(string)
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/reports/sendEmailFnoTradebook [GET]
func SendEmailFnoTradebook(c *gin.Context) {
	var reqParams models.FNOTradebookReq

	dfDateFr := c.Query("dfdatefr")
	dfDateTo := c.Query("dfdateto")

	if dfDateFr == "" || dfDateTo == "" {
		loggerconfig.Error("SendEmailFnoTradebook (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	profileData, _ := c.Get("profileData")
	profileInfo, _ := (profileData).(models.ProfileDataResp)

	reqParams.UserID = requestH.ClientId
	reqParams.DFDateFr = dfDateFr
	reqParams.DFDateTo = dfDateTo

	if requestH.DeviceType == "" {
		loggerconfig.Error("SendEmailFnoTradebook (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("SendEmailFnoTradebook (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("SendEmailFnoTradebook (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theReportsProvider.SendEmailFnoTradebook(reqParams, requestH, profileInfo)
	logDetail := "clientId: " + requestH.ClientId + " function: SendEmailFnoTradebook requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
