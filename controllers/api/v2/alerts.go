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

var theAlertsProviderV2 models.AlertsProvider

func InitAlertsProviderV2(provider models.AlertsProvider) {
	defer models.HandlePanic()
	theAlertsProviderV2 = provider
}

// EditAlerts
// @Tags space Alerts API's V2
// @Description Edit Alerts V2
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.EditAlertsReq true "Alerts"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/alerts/editAlerts [PUT]
func EditAlerts(c *gin.Context) {
	var reqParams models.EditAlertsReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("EditAlertsV2 (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("EditAlertsV2 (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("EditAlertsV2 (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("EditAlertsV2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("EditAlertsV2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("EditAlertsV2 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientId, "requestId:", requestH.RequestId)

	code, resp := theAlertsProviderV2.EditAlerts(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: EditAlerts requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// PauseAlerts
// @Tags space Alerts API's V2
// @Description PauseAlerts V2
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PauseAlertsReq true "Alerts"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/alerts/pauseAlerts [PUT]
func PauseAlerts(c *gin.Context) {
	var reqParams models.PauseAlertsReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("PauseAlertsV2 (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("PauseAlertsV2 (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("PauseAlertsV2 (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("PauseAlertsV2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("PauseAlertsV2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("PauseAlertsV2 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientId, "requestId:", requestH.RequestId)

	code, resp := theAlertsProviderV2.PauseAlerts(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: PauseAlerts requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DeleteAlerts
// @Tags space Alerts API's V2
// @Description Delete Alerts V2
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.DeleteAlertsReq true "Alerts"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/alerts/deleteAlerts [DELETE]
func DeleteAlerts(c *gin.Context) {
	var reqParams models.DeleteAlertsReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("DeleteAlertsV2 (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("DeleteAlertsV2 (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("DeleteAlertsV2 (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("DeleteAlertsV2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("DeleteAlertsV2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("DeleteAlertsV2 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientId, "requestId:", requestH.RequestId)

	code, resp := theAlertsProviderV2.DeleteAlerts(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: DeleteAlerts requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetAlerts
// @Tags space Alerts API's V2
// @Description Get Alerts V2
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param clientId query string true "clientId Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.AlersRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/alerts/getAlerts [GET]
func GetAlerts(c *gin.Context) {
	var reqParams models.GetAlertsReq
	clientID := c.Query("clientId")
	if clientID == "" {
		loggerconfig.Error("GetAlertsV2 (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.ClientId = clientID
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetAlertsV2 (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetAlertsV2 (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("GetAlertsV2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("GetAlertsV2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("GetAlertsV2 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientId, "requestId:", requestH.RequestId)

	code, resp := theAlertsProviderV2.GetAlerts(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: GetAlerts requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
