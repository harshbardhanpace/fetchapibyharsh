package v2

import (
	"encoding/json"
	"net/http"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var theFetchFundsProviderV2 models.FetchFundsProvider

func InitFundsProviderV2(funds models.FetchFundsProvider) {
	defer models.HandlePanic()
	theFetchFundsProviderV2 = funds
}

// Funds
// @Tags space funds V2
// @Description Fetch Funds V2
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param clientId query string true "clientId Query Parameter" dataType(string)
// @Param type query string true "type Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.FetchFundsResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/funds/view/fetchFunds [GET]
func FetchFunds(c *gin.Context) {
	var reqParams models.FetchFundsRequest
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	clientID := c.Query("clientId")
	type1 := c.Query("type")
	if clientID == "" || type1 == "" {
		loggerconfig.Error("FetchFundsV2 (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.ClientID = clientID
	reqParams.Type = type1

	if reqParams.ClientID == "" {
		loggerconfig.Error("FetchFunds V2(controller), Empty Client Id requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientID)
		apihelpers.ErrorMessage(c, constants.InvalidClient)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FetchFunds V2(controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientID)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("FetchFunds V2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("FetchFunds V2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("FetchFunds V2(controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
	code, resp := theFetchFundsProviderV2.FetchFunds(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}

// CancelPayout
// @Tags space funds V2
// @Description Cancel Payout V2
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.CancelPayoutReq true "funds"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/funds/view/cancelPayout [PUT]
func CancelPayout(c *gin.Context) {
	var reqParams models.CancelPayoutReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	errr := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if errr != nil {
		loggerconfig.Error("CancelPayout V2(controller), error decoding body, error:", errr)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if reqParams.ClientID == "" {
		loggerconfig.Error("CancelPayout V2(controller), Empty Client Id requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientID)
		apihelpers.ErrorMessage(c, constants.InvalidClient)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("CancelPayout V2(controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientID)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("CancelPayout V2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("CancelPayout V2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("CancelPayout V2(controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
	code, resp := theFetchFundsProviderV2.CancelPayout(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}

// ClientTransactions
// @Tags space funds V2
// @Description ClientTransactions V2
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param clientId query string true "clientId Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.ClientTransactionsResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/funds/view/clientTransactions [GET]
func ClientTransactions(c *gin.Context) {
	var reqParams models.ClientTransactionsRequest
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	clientID := c.Query("clientId")
	if clientID == "" {
		loggerconfig.Error("ClientTransactions V2(controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.ClientID = clientID

	if reqParams.ClientID == "" {
		loggerconfig.Error("ClientTransactions V2(controller), Empty Client Id requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientID)
		apihelpers.ErrorMessage(c, constants.InvalidClient)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("ClientTransactions V2(controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientID)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("ClientTransactions V2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("ClientTransactions V2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("ClientTransactions V2(controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
	code, resp := theFetchFundsProviderV2.ClientTransactions(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}
