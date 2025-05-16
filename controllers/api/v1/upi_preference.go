package v1

import (
	"encoding/json"
	"net/http"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	masker "github.com/ggwhite/go-masker"
	"github.com/gin-gonic/gin"
)

var theUpiPreferenceProvider models.UpiPreferenceProvider

func InitUpiPreferenceProvider(provider models.UpiPreferenceProvider) {
	defer models.HandlePanic()
	theUpiPreferenceProvider = provider
	maskObj = masker.New()
}

// SetUpiPreference
// @Tags space upi preference
// @Description Set Upi Preference
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.SetUpiPreferenceReq true "Set Upi Preference"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/upi/setUpiPreference [POST]
func SetUpiPreference(c *gin.Context) {
	var reqParams models.SetUpiPreferenceReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("SetUpiPreference (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("SetUpiPreference (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	maskedReq, err := maskObj.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("In Controller SetUpiPreference Error in masking request err: ", err, " clientId: ", requestH.ClientId, " requestid = ", requestH.RequestId)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("SetUpiPreference (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("SetUpiPreference (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("SetUpiPreference (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId)
	code, resp := theUpiPreferenceProvider.SetUpiPreference(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: SetUpiPreference requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchUpiPreference
// @Tags space upi preference
// @Description Fetch Upi Preference
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param clientId query string true "clientId Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.FetchUpiPreferenceRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/upi/fetchUpiPreference [GET]
func FetchUpiPreference(c *gin.Context) {
	var reqParams models.FetchUpiPreferenceReq
	clientID := c.Query("clientId")
	if clientID == "" {
		loggerconfig.Error("FetchUpiPreference (controller), error parsing the query params in Get request")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.ClientId = clientID

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchUpiPreference (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("FetchUpiPreference (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("FetchUpiPreference (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	maskedReq, err := maskObj.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("In Controller FetchUpiPreference Error in masking request err: ", err, " clientId: ", requestH.ClientId, " requestid = ", requestH.RequestId)
		return
	}

	loggerconfig.Info("FetchUpiPreference (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId)
	code, resp := theUpiPreferenceProvider.FetchUpiPreference(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: FetchUpiPreference requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DeleteUpiPreference
// @Tags space upi preference
// @Description Delete Upi Preference
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.DeleteUpiPreferenceReq true "Delete Upi Preference"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/upi/deleteUpiPreference [DELETE]
func DeleteUpiPreference(c *gin.Context) {
	var reqParams models.DeleteUpiPreferenceReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("DeleteUpiPreference (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("DeleteUpiPreference (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("DeleteUpiPreference (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("DeleteUpiPreference (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	maskedReq, err := maskObj.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("In Controller DeleteUpiPreference Error in masking request err: ", err, " clientId: ", requestH.ClientId, " requestid = ", requestH.RequestId)
		return
	}

	loggerconfig.Info("DeleteUpiPreference (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId)
	code, resp := theUpiPreferenceProvider.DeleteUpiPreference(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientId + " function: DeleteUpiPreference requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
