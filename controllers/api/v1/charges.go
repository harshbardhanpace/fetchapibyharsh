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

var theChargesProvider models.ChargesProvider

func InitChargesProvider(provider models.ChargesProvider) {
	defer models.HandlePanic()
	theChargesProvider = provider
}

// BrokerCharges
// @Tags space charges
// @Description Broker Charges
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.BrokerChargesReq true "BrokerCharges"
// @Success 200 {object} apihelpers.APIRes{data=models.BrokerChargesRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/charges/brokerCharges [POST]
func BrokerCharges(c *gin.Context) {
	var brokerChargesReq models.BrokerChargesReq

	errr := json.NewDecoder(c.Request.Body).Decode(&brokerChargesReq)
	if errr != nil {
		loggerconfig.Error("BrokerCharges (controller), error decoding body, error:", errr)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("BrokerCharges (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", brokerChargesReq.ClientID, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(brokerChargesReq)
	if err != nil {
		loggerconfig.Error("BrokerCharges (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", brokerChargesReq.ClientID, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if !checkValidBrokerageRequest(brokerChargesReq, requestH) {
		apihelpers.SendErrorController(c, false, constants.InvalidRequest, http.StatusBadRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(brokerChargesReq.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("BrokerCharges (controller) CheckAuthWithClient invalid authtoken", " clientId: ", brokerChargesReq.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("BrokerCharges (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", brokerChargesReq.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("BrokerCharges (controller), reqParams:", helpers.LogStructAsJSON(brokerChargesReq), " clientId: ", brokerChargesReq.ClientID, "requestId:", requestH.RequestId, " deviceId: ", requestH.DeviceId)

	//call service
	code, resp := theChargesProvider.BrokerCharges(brokerChargesReq, requestH)

	//return response using api helper
	logDetail := "clientId: " + brokerChargesReq.ClientID + " function: BrokerCharges requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// CombineBrokerCharges
// @Tags space charges
// @Description Combine Broker Charges - Brokerage Charges for multiple orders
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.CombineBrokerChargesReq true "BrokerCharges"
// @Success 200 {object} apihelpers.APIRes{data=models.CombineBrokerChargesRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/charges/combineBrokerCharges [POST]
func CombineBrokerCharges(c *gin.Context) {
	var combineBrokerChargesReq models.CombineBrokerChargesReq

	errr := json.NewDecoder(c.Request.Body).Decode(&combineBrokerChargesReq)
	if errr != nil {
		loggerconfig.Error("CombineBrokerCharges (controller), error decoding body, error:", errr)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("CombineBrokerCharges (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", combineBrokerChargesReq.ClientID, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(combineBrokerChargesReq)
	if err != nil {
		loggerconfig.Error("CombineBrokerCharges (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", combineBrokerChargesReq.ClientID, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	for i := 0; i < len(combineBrokerChargesReq.BrokerCharges); i++ {
		if !checkValidBrokerageRequest(combineBrokerChargesReq.BrokerCharges[i], requestH) {
			apihelpers.SendErrorController(c, false, constants.InvalidRequest, http.StatusBadRequest)
			return
		}
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(combineBrokerChargesReq.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("CombineBrokerCharges (controller) CheckAuthWithClient invalid authtoken", " clientId: ", combineBrokerChargesReq.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}

	if !matchStatus {
		loggerconfig.Error("CombineBrokerCharges (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", combineBrokerChargesReq.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("CombineBrokerCharges (controller), reqParams:", helpers.LogStructAsJSON(combineBrokerChargesReq), " clientId: ", combineBrokerChargesReq.ClientID, "requestId:", requestH.RequestId, " deviceId: ", requestH.DeviceId)

	//call service
	code, resp := theChargesProvider.CombineBrokerCharges(combineBrokerChargesReq, requestH)

	//return response using api helper
	logDetail := "clientId: " + combineBrokerChargesReq.ClientID + " function: CombineBrokerCharges requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

func checkValidBrokerageRequest(brokerChargesReq models.BrokerChargesReq, requestH models.ReqHeader) bool {
	if !(strings.ToLower(brokerChargesReq.Segment) == constants.EQUITY || strings.ToLower(brokerChargesReq.Segment) == constants.CURRENCY || strings.ToLower(brokerChargesReq.Segment) == constants.COMMODITY) {
		loggerconfig.Error("checkValidBrokerageRequest (controller), Error invalid segment provided: ", strings.ToLower(brokerChargesReq.Segment), " requestId: ", requestH.RequestId, "clientId: ", brokerChargesReq.ClientID, " deviceId: ", requestH.DeviceId)
		return false
	}

	if !(strings.ToLower(brokerChargesReq.SubSegment) == constants.DELIVERY || strings.ToLower(brokerChargesReq.SubSegment) == constants.INTRADAY || strings.ToLower(brokerChargesReq.SubSegment) == constants.FUTURES || strings.ToLower(brokerChargesReq.SubSegment) == constants.OPTIONS) {
		loggerconfig.Error("checkValidBrokerageRequest (controller), Error invalid subSegment provided: ", strings.ToLower(brokerChargesReq.SubSegment), " requestId: ", requestH.RequestId, "clientId: ", brokerChargesReq.ClientID, " deviceId: ", requestH.DeviceId)
		return false
	}

	return true
}

// FundsPayout
// @Tags space charges
// @Description Funds Payout
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.FundsPayoutReq true "FundsPayout"
// @Success 200 {object} apihelpers.APIRes{data=models.FundsPayoutRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/charges/fundsPayout [POST]
func FundsPayout(c *gin.Context) {
	var fundsPayoutReq models.FundsPayoutReq

	errr := json.NewDecoder(c.Request.Body).Decode(&fundsPayoutReq)
	if errr != nil {
		loggerconfig.Error("FundsPayout error decoding body, error:", errr)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FundsPayout Empty Device Type requestId: ", requestH.RequestId, " clientId: ", fundsPayoutReq.ClientID, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	if requestH.ClientType != constants.ADMIN {
		matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(fundsPayoutReq.ClientID, requestH.Authorization)
		if !tokenValidStatus {
			loggerconfig.Error("FundsPayout (controller) CheckAuthWithClient invalid authtoken", " clientId: ", fundsPayoutReq.ClientID, " requestId:", requestH.RequestId)
			apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
			return
		}
		if !matchStatus {
			loggerconfig.Error("FundsPayout (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", fundsPayoutReq.ClientID, " requestId:", requestH.RequestId)
			apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
			return
		}
	}

	loggerconfig.Info("FundsPayout, reqParams:", helpers.LogStructAsJSON(fundsPayoutReq), " clientId: ", fundsPayoutReq.ClientID, "requestId:", requestH.RequestId, " deviceId: ", requestH.DeviceId)

	//call service
	code, resp := theChargesProvider.FundsPayout(fundsPayoutReq, requestH)

	//return response using api helper
	logDetail := "clientId: " + fundsPayoutReq.ClientID + " function: FundsPayout requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
