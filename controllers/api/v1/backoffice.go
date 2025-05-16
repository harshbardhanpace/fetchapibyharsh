package v1

import (
	"encoding/json"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var theBackofficeProvider models.BackofficeProvider

func InitBackOfficeProvider(provider models.BackofficeProvider) {
	defer models.HandlePanic()
	theBackofficeProvider = provider
}

func GetBackOfficeProvider() models.BackofficeProvider {
	return theBackofficeProvider
}

// TradeConfirmationDateRange
// @Tags space Shilpi V1
// @Description Trade Confirmation Date Range
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.TradeConfirmationDateRangeReq true "Create Basket"
// @Success 200 {object} apihelpers.APIRes{data=models.TradeConfirmationDateRangeRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/shilpi/tradeConfirmationDateRange [POST]
func TradeConfirmationDateRange(c *gin.Context) {
	var reqParams models.TradeConfirmationDateRangeReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("TradeConfirmationDateRange (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("TradeConfirmationDateRange (controller), Empty Device Type requestId: ", requestH.RequestId, "userID: ", reqParams.UserID, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("TradeConfirmationDateRange (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "userID: ", reqParams.UserID, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("TradeConfirmationDateRange (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "userID: ", reqParams.UserID, " deviceId: ", requestH.DeviceId)
	code, resp := theBackofficeProvider.TradeConfirmationDateRange(reqParams, requestH)
	logDetail := "clientId: " + reqParams.UserID + " function: TradeConfirmationDateRange requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetBillDetailsCdsl
// @Tags space Shilpi V1
// @Description Get Bill Details Cdsl
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.GetBillDetailsCdslReq true "Create Basket"
// @Success 200 {object} apihelpers.APIRes{data=models.GetBillDetailsCdslRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/shilpi/getBillDetailsCdsl [POST]
func GetBillDetailsCdsl(c *gin.Context) {
	var reqParams models.GetBillDetailsCdslReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetBillDetailsCdsl (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetBillDetailsCdsl (controller), Empty Device Type requestId: ", requestH.RequestId, "userID: ", reqParams.UserID, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetBillDetailsCdsl (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "userID: ", reqParams.UserID, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetBillDetailsCdsl (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "userID: ", reqParams.UserID, " deviceId: ", requestH.DeviceId)
	code, resp := theBackofficeProvider.GetBillDetailsCdsl(reqParams, requestH)
	logDetail := "clientId: " + reqParams.UserID + " function: GetBillDetailsCdsl requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// LongTermShortTerm
// @Tags space Shilpi V1
// @Description Long Term Short Term
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.LongTermShortTermReq true "Create Basket"
// @Success 200 {object} apihelpers.APIRes{data=models.LongTermShortTermRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/shilpi/longTermShortTerm [POST]
func LongTermShortTerm(c *gin.Context) {
	var reqParams models.LongTermShortTermReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("LongTermShortTerm (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("LongTermShortTerm (controller), Empty Device Type requestId: ", requestH.RequestId, "userID: ", reqParams.UserID, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("LongTermShortTerm (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "userID: ", reqParams.UserID, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("LongTermShortTerm (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "userID: ", reqParams.UserID, " deviceId: ", requestH.DeviceId)
	code, resp := theBackofficeProvider.LongTermShortTerm(reqParams, requestH)
	logDetail := "clientId: " + reqParams.UserID + " function: LongTermShortTerm requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchProfile
// @Tags space Shilpi V1
// @Description Fetch Profile
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.FetchProfileReq true "Create Basket"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchProfileRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/shilpi/fetchProfile [POST]
func FetchProfile(c *gin.Context) {
	var reqParams models.FetchProfileReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FetchProfile (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchProfile (controller), Empty Device Type requestId: ", requestH.RequestId, "userID: ", reqParams.UserID, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FetchProfile (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "userID: ", reqParams.UserID, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("FetchProfile (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "userID: ", reqParams.UserID, " deviceId: ", requestH.DeviceId)
	code, resp := theBackofficeProvider.FetchProfile(reqParams, requestH)
	logDetail := "clientId: " + reqParams.UserID + " function: FetchProfile requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// TradeConfirmationOnDate
// @Tags space Shilpi V1
// @Description Trade Confirmation On Date
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.TradeConfirmationOnDateReq true "Trade Confirmation On Date"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchProfileRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/shilpi/tradeConfirmationOnDate [POST]
func TradeConfirmationOnDate(c *gin.Context) {
	var reqParams models.TradeConfirmationOnDateReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("TradeConfirmationOnDate (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("TradeConfirmationOnDate (controller), Empty Device Type requestId: ", requestH.RequestId, "userID: ", reqParams.UserId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("TradeConfirmationOnDate (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "userID: ", reqParams.UserId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("TradeConfirmationOnDate (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "userID: ", reqParams.UserId, " deviceId: ", requestH.DeviceId)
	code, resp := theBackofficeProvider.TradeConfirmationOnDate(reqParams, requestH)
	logDetail := "clientId: " + reqParams.UserId + " function: TradeConfirmationOnDate requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// OpenPositions
// @Tags space Shilpi V1
// @Description Open Positions
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.OpenPositionsReq true "Open Positions"
// @Success 200 {object} apihelpers.APIRes{data=models.OpenPositionsRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/shilpi/openPositions [POST]
func OpenPositions(c *gin.Context) {
	var reqParams models.OpenPositionsReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("OpenPositions (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("OpenPositions (controller), Empty Device Type requestId: ", requestH.RequestId, "userID: ", reqParams.UserId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("OpenPositions (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "userID: ", reqParams.UserId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("OpenPositions (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "userID: ", reqParams.UserId, " deviceId: ", requestH.DeviceId)
	code, resp := theBackofficeProvider.OpenPositions(reqParams, requestH)
	logDetail := "clientId: " + reqParams.UserId + " function: OpenPositions requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetHolding
// @Tags space Shilpi V1
// @Description Get Holding
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.GetHoldingReq true "Get Holding"
// @Success 200 {object} apihelpers.APIRes{data=models.GetHoldingRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/shilpi/getHolding [POST]
func GetHolding(c *gin.Context) {
	var reqParams models.GetHoldingReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetHolding (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetHolding (controller), Empty Device Type requestId: ", requestH.RequestId, "userID: ", reqParams.UserId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetHolding (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "userID: ", reqParams.UserId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetHolding (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "userID: ", reqParams.UserId, " deviceId: ", requestH.DeviceId)
	code, resp := theBackofficeProvider.GetHolding(reqParams, requestH)
	logDetail := "clientId: " + reqParams.UserId + " function: GetHolding requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetMarginOnDate
// @Tags space Shilpi V1
// @Description Get Margin On Date
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.GetMarginOnDateReq true "Get Margin On Date"
// @Success 200 {object} apihelpers.APIRes{data=models.GetMarginOnDateRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/shilpi/getMarginOnDate [POST]
func GetMarginOnDate(c *gin.Context) {
	var reqParams models.GetMarginOnDateReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetMarginOnDate (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetMarginOnDate (controller), Empty Device Type requestId: ", requestH.RequestId, "userID: ", reqParams.UserId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetMarginOnDate (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "userID: ", reqParams.UserId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetMarginOnDate (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "userID: ", reqParams.UserId, " deviceId: ", requestH.DeviceId)
	code, resp := theBackofficeProvider.GetMarginOnDate(reqParams, requestH)
	logDetail := "clientId: " + reqParams.UserId + " function: GetMarginOnDate requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FinancialLedgerBalanceOnDate
// @Tags space Shilpi V1
// @Description Financial Ledger Balance On Date
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.FinancialLedgerBalanceOnDateReq true "Financial Ledger Balance On Date"
// @Success 200 {object} apihelpers.APIRes{data=models.FinancialLedgerBalanceOnDateRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/shilpi/financialLedgerBalanceOnDate [POST]
func FinancialLedgerBalanceOnDate(c *gin.Context) {
	var reqParams models.FinancialLedgerBalanceOnDateReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FinancialLedgerBalanceOnDate (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FinancialLedgerBalanceOnDate (controller), Empty Device Type requestId: ", requestH.RequestId, "userID: ", reqParams.UserId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FinancialLedgerBalanceOnDate (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "userID: ", reqParams.UserId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("FinancialLedgerBalanceOnDate (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "userID: ", reqParams.UserId, " deviceId: ", requestH.DeviceId)
	code, resp := theBackofficeProvider.FinancialLedgerBalanceOnDate(reqParams, requestH)
	logDetail := "clientId: " + reqParams.UserId + " function: FinancialLedgerBalanceOnDate requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetFinancial
// @Tags space Shilpi V1
// @Description Get Financial
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.GetFinancialReq true "GetFinancial"
// @Success 200 {object} apihelpers.APIRes{data=models.GetFinancialRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/shilpi/getFinancial [POST]
func GetFinancial(c *gin.Context) {
	var reqParams models.GetFinancialReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetFinancial (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetFinancial (controller), Empty Device Type requestId: ", requestH.RequestId, "userID: ", reqParams.UserId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetFinancial (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "userID: ", reqParams.UserId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetFinancial (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "userID: ", reqParams.UserId, " deviceId: ", requestH.DeviceId)
	code, resp := theBackofficeProvider.GetFinancial(reqParams, requestH)
	logDetail := "clientId: " + reqParams.UserId + " function: GetFinancial requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
