package v1

import (
	"encoding/json"
	"net/http"
	"strings"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var theIpoProvider models.IpoProvider

func InitIpoProvider(provider models.IpoProvider) {
	defer models.HandlePanic()
	theIpoProvider = provider
}

// GetAllIpo
// @Tags space ipo V1
// @Description Get All Ipo - It will display all ipo (upcoming, closed, current)
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param ClientId header string true "ClientId"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Success 200 {object} apihelpers.APIRes{data=models.GetAllIpoResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/tradeipo/getAllIpo [GET]
func GetAllIpo(c *gin.Context) {
	var reqParams models.GetAllIpoRequest

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetAllIpo (controller), Empty Device Type requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetAllIpo (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("GetAllIpo (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theIpoProvider.GetAllIpo(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: GetAllIpo requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// PlaceIpoOrder
// @Tags space ipo V1
// @Description Place Ipo Order - It will place current open ipos
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.PlaceIpoOrderRequest true "ipo"
// @Success 200 {object} apihelpers.APIRes{data=models.PlaceIpoOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/tradeipo/placeIpoOrder [POST]
func PlaceIpoOrder(c *gin.Context) {
	var reqParams models.PlaceIpoOrderRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("PlaceIpoOrder (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("PlaceIpoOrder (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("PlaceIpoOrder (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("PlaceIpoOrder (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("PlaceIpoOrder (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	maskedReq, err := maskObj.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("In Controller PlaceIpoOrder Error in masking request err: ", err, " clientId: ", requestH.ClientId, " requestid = ", requestH.RequestId)
		return
	}

	loggerconfig.Info("PlaceIpoOrder (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
	code, resp := theIpoProvider.PlaceIpoOrder(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: PlaceIpoOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchIpoOrder
// @Tags space ipo V1
// @Description Fetch Ipo Order - It will give applied IPO
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.FetchIpoOrderRequest true "ipo"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchIpoOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/tradeipo/fetchIpoOrder [POST]
func FetchIpoOrder(c *gin.Context) {
	var reqParams models.FetchIpoOrderRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FetchIpoOrder (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchIpoOrder (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FetchIpoOrder (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("FetchIpoOrder (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("FetchIpoOrder (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("FetchIpoOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
	code, resp := theIpoProvider.FetchIpoOrder(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: FetchIpoOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// CancelIpoOrder
// @Tags space ipo V1
// @Description Cancel Ipo Order - It will cancel placed IPO order
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.CancelIpoOrderRequest true "ipo"
// @Success 200 {object} apihelpers.APIRes{data=models.CancelIpoOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/tradeipo/cancelIpoOrder [POST]
func CancelIpoOrder(c *gin.Context) {
	var reqParams models.CancelIpoOrderRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("CancelIpoOrder (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("CancelIpoOrder (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("CancelIpoOrder (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("CancelIpoOrder (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("CancelIpoOrder (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	maskedReq, err := maskObj.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("In Controller CancelIpoOrder Error in masking request err: ", err, " clientId: ", requestH.ClientId, " requestid = ", requestH.RequestId)
		return
	}

	loggerconfig.Info("CancelIpoOrder (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
	code, resp := theIpoProvider.CancelIpoOrder(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: CancelIpoOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchIpoData
// @Tags space ipo V1
// @Description Fetch Ipo Data - It will give details of IPO from ipo name
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param request body models.FetchIpoDataRequest true "ipo"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchIpoDataResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/tradeipo/fetchIpoData [POST]
func FetchIpoData(c *gin.Context) {
	var reqParams models.FetchIpoDataRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FetchIpoData (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchIpoData (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FetchIpoData (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("FetchIpoData (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theIpoProvider.FetchIpoData(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: FetchIpoData requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchIpoGmpData
// @Tags space ipo V1
// @Description Fetch Ipo Gmp Data - It provides IPO Grey Market Premium
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.FetchIpoDataRequest true "ipo"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchIpoGmpDataResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/tradeipo/fetchIpoGmpData [POST]
func FetchIpoGmpData(c *gin.Context) {
	var reqParams models.FetchIpoDataRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FetchIpoGmpData (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchIpoGmpData (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FetchIpoGmpData (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("FetchIpoGmpData (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theIpoProvider.FetchIpoGmpData(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: FetchIpoGmpData requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchEIpo
// @Tags space ipo V1
// @Description Fetch EIpo Data - fetch ipo-current-issue, ipo-closed, ipo-upcoming
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param ClientId header string true "ClientId"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param request body models.FetchEipoReq true "ipo"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchNseIpoResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/tradeipo/fetchIpoDataNse [POST]
func FetchEIpo(c *gin.Context) {

	var reqParams models.FetchEipoReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("fetchIpoDataNse (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchEIpo (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	ipoSymbol := reqParams.IpoSymbol
	ipoStage := reqParams.IpoStage
	ipoStage = strings.ToUpper(ipoStage)

	// Ipo stage can only be UPCOMING, CLOSED or CURRENT
	if ipoStage != "" && ipoStage != constants.CLOSED && ipoStage != constants.CURRENT && ipoStage != constants.UPCOMING {
		loggerconfig.Error("FetchEIpo (controller), Invalid ipoStage: ", ipoStage)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	reqParams.IpoSymbol = ipoSymbol
	reqParams.IpoStage = ipoStage

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FetchEIpo (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("FetchEIpo (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theIpoProvider.FetchEIpo(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: FetchEIpo requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
