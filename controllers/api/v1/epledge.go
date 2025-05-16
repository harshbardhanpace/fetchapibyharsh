package v1

import (
	"encoding/json"
	"net/http"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var theEpledgeProvider models.EpledgeProvider

func InitEpledgeProvider(provider models.EpledgeProvider) {
	defer models.HandlePanic()
	theEpledgeProvider = provider
}

// EpledgeRequest
// @Tags space epledge V1
// @Description Epledge Request
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.EpledgeReq true "EpledgeRequest"
// @Success 200 {object} apihelpers.APIRes{data=models.EpledgeRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/epledge/epledgeRequest [POST]
func EpledgeRequest(c *gin.Context) {
	var reqParams models.EpledgeReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("EpledgeRequest (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")

	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("EpledgeRequest (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("EpledgeRequest (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("EpledgeRequest (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("EpledgeRequest (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}
	for _, isinDetail := range reqParams.IsinDetails {
		quantity, err := strconv.Atoi(isinDetail.Quantity)
		if err != nil || quantity <= 0 {
			loggerconfig.Error("EpledgeRequest (controller) Invalid quantity value: ", isinDetail.Quantity, " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
			apihelpers.SendErrorController(c, false, constants.InvalidRequest, http.StatusBadRequest)
			return
		}
		price, err := strconv.ParseFloat(isinDetail.Price, 64)
		if err != nil || price <= 0 {
			loggerconfig.Error("EpledgeRequest (controller) Invalid price value: ", isinDetail.Price, " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
			apihelpers.SendErrorController(c, false, constants.InvalidRequest, http.StatusBadRequest)
			return
		}
	}

	loggerconfig.Info("EpledgeRequest (controller), reqParams:", reqParams, "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theEpledgeProvider.SendEpledgeRequest(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: EpledgeRequest requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// UnpledgeRequest
// @Tags space epledge V1
// @Description UnpledgeRequest Request
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.UnpledgeReq true "UnpledgeRequest"
// @Success 200 {object} apihelpers.APIRes{data=[]models.UnpledgeRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/epledge/unpledge [POST]
func UnpledgeRequest(c *gin.Context) {
	var reqParams models.UnpledgeReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("UnpledgeRequest (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")

	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("UnpledgeRequest (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("UnpledgeRequest (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("UnpledgeRequest (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("UnpledgeRequest (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}
	for _, unpledge := range reqParams.UnpledgeList {
		if unpledge.Quantity <= 0 {
			loggerconfig.Error("EpledgeRequest (controller) Invalid quantity value: ", unpledge.Quantity, " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
			apihelpers.SendErrorController(c, false, constants.InvalidRequest, http.StatusBadRequest)
			return
		}
	}

	loggerconfig.Info("UnpledgeRequest (controller), reqParams:", reqParams, "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theEpledgeProvider.UnpledgeRequest(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: UnpledgeRequest requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// MTFEpledgeRequest
// @Tags space epledge V1
// @Description MTFEpledgeRequest Request
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.MTFEPledgeRequest true "MTFEpledgeRequest"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/epledge/mtfEpledgeRequest [POST]
func MTFEpledgeRequest(c *gin.Context) {
	var reqParams models.MTFEPledgeRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("MTFEpledgeRequest (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")

	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("MTFEpledgeRequest (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("MTFEpledgeRequest (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("MTFEpledgeRequest (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("MTFEpledgeRequest (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("MTFEpledgeRequest (controller), reqParams:", reqParams, "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theEpledgeProvider.MTFEpledgeRequest(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: MTFEpledgeRequest requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetPledgeList
// @Tags space epledge V1
// @Description Get Pledge List Request
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Success 200 {object} apihelpers.APIRes{data=[]models.MTFPledgeListResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/epledge/getPledgeList [GET]
func GetPledgeList(c *gin.Context) {
	cRH, _ := c.Get("reqH")

	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetPledgeList (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("GetPledgeList (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("GetPledgeList (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("GetPledgeList (controller), requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theEpledgeProvider.GetPledgeList(requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: GetPledgeList requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetCTDQuantityList
// @Tags space epledge V1
// @Description Get CTD Quantity List
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.MTFCTDDataReq true "MTFCTDDataReq"
// @Success 200 {object} apihelpers.APIRes{data=[]models.MTFCTDDataRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/epledge/getCTDQuantityList [POST]
func GetCTDQuantityList(c *gin.Context) {
	var reqParams models.MTFCTDDataReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetCTDQuantityList (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	cRH, _ := c.Get("reqH")

	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetCTDQuantityList (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetCTDQuantityList (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("GetCTDQuantityList (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("GetCTDQuantityList (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("GetCTDQuantityList (controller), requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theEpledgeProvider.GetCTDQuantityList(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: GetCTDQuantityList requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetPledgeTransactions
// @Tags space epledge V1
// @Description Get all pledge transactions against user
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param clientId query string true "ClientId"
// @Param page query string true "Page number"
// @Param startDate query string true "start date"
// @Param endDate query string true "end date"
// @Success 200 {object} apihelpers.APIRes{data=[]models.PledgeData}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/epledge/getPledgeTransactions [GET]
func GetPledgeTransactions(c *gin.Context) {

	cRH, _ := c.Get("reqH")

	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetPledgeTransactions (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("GetPledgeTransactions (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("GetPledgeTransactions (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	// read clientId
	clientId := c.Query("clientId")
	if clientId == "" {
		loggerconfig.Error("GetPledgeTransactions (controller) invalid clientId", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidClient)
		return
	}

	// read page
	pageInt, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		loggerconfig.Error("GetPledgeTransactions (controller), error converting page to int type, Invalid page received: ", pageInt, " requestId:", requestH.RequestId, "ClientID: ", clientId)
		pageInt = 1
	}

	// read startDate
	startDate := c.Query("startDate")
	if startDate == "" {
		loggerconfig.Error("GetPledgeTransactions (controller) invalid startDate", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDate)
		return
	} else {
		_, err := time.Parse(constants.YYYYMMDD, startDate)
		if err != nil {
			loggerconfig.Error("GetPledgeTransactions (controller) invalid startDate", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidDate)
			return
		}
	}

	// read endDate
	endDate := c.Query("endDate")
	if endDate == "" {
		loggerconfig.Error("GetPledgeTransactions (controller) invalid endDate", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDate)
		return
	} else {
		_, err := time.Parse(constants.YYYYMMDD, endDate)
		if err != nil {
			loggerconfig.Error("GetPledgeTransactions (controller) invalid endDate", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidDate)
			return
		}

	}

	fetchEpledgeTxnReq := models.FetchEpledgeTxnReq{
		ClientID:  clientId,
		Page:      pageInt,
		StartDate: startDate,
		EndDate:   endDate,
	}

	loggerconfig.Info("GetPledgeTransactions (controller), requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theEpledgeProvider.GetPledgeTransactions(fetchEpledgeTxnReq, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: GetPledgeTransactions requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}

// MTFCTD
// @Tags space epledge V1
// @Description MTF CTD
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.MTFCTDReq true "MTFCTDReq"
// @Success 200 {object} apihelpers.APIRes{data=[]models.MTFCTDResDataRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/epledge/mtfCtd [POST]
func MTFCTD(c *gin.Context) {
	var reqParams models.MTFCTDReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("MTFCTD (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")

	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("MTFCTD (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("MTFCTD (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("MTFCTD (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("MTFCTD (controller), requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theEpledgeProvider.MTFCTD(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: MTFCTD requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
