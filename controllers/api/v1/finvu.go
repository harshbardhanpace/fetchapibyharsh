package v1

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

var theFinvuProvider models.FinvuProvider

func InitFinvuProvider(provider models.FinvuProvider) {
	defer models.HandlePanic()
	theFinvuProvider = provider
}

// FinvuConsentRequestPlus
// @Tags space Finvu
// @Description Finvu Consent Request Plus
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string true "P-ClientPublicIP Header"
// @Param request body models.CreateConsentRequestPlusReq true "CreateConsentRequestPlusReq"
// @Success 200 {object} apihelpers.APIRes{data=models.ConsentsRequestPlusRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/finvu/finvuConsentRequestPlus [POST]
func FinvuConsentRequestPlus(c *gin.Context) {
	var reqParams models.CreateConsentRequestPlusReq

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FinvuConsentRequestPlus  (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("FinvuConsentRequestPlus (controller), Empty Device Type requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FinvuConsentRequestPlus (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("EpledgeRequest (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("FinvuConsentRequestPlus (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("FinvuConsentRequestPlus (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theFinvuProvider.FinvuConsentRequestPlus(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientId + " function: FinvuConsentRequestPlus requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FinvuGetBankStatement
// @Tags space Finvu
// @Description Finvu Consent Request Plus
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string true "P-ClientPublicIP Header"
// @Param request body models.FinvuGetBankStatementReq true "FinvuGetBankStatementReq"
// @Success 200 {object} apihelpers.APIRes{data=models.ConsentsRequestPlusRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/finvu/finvuGetBankStatement [POST]
func FinvuGetBankStatement(c *gin.Context) {
	var reqParams models.FinvuGetBankStatementReq

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FinvuGetBankStatement  (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("FinvuGetBankStatement (controller), Empty Device Type requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FinvuGetBankStatement (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("EpledgeRequest (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("FinvuGetBankStatement (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("FinvuGetBankStatement (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theFinvuProvider.FinvuGetBankStatement(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientId + " function: FinvuGetBankStatement requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
