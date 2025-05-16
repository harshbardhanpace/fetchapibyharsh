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

var theAlertsProvider models.AlertsProvider

func InitAlertsProvider(provider models.AlertsProvider) {
	defer models.HandlePanic()
	theAlertsProvider = provider
}

// SetAlerts
// @Tags space Alerts API's
// @Description Set   Alerts
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.CreateAlertsReq true "Alerts"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/alerts/setAlerts [POST]
func SetAlerts(c *gin.Context) {
	var reqParams models.CreateAlertsReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("SetAlerts (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("SetAlerts (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientId, " deviceId: ", requestH.DeviceId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("SetAlerts (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientId, " deviceId: ", requestH.DeviceId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("SetAlerts (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("SetAlerts (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("SetAlerts (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientId, " requestId:", requestH.RequestId, " deviceId: ", requestH.DeviceId, " clientVersion: ", requestH.ClientVersion)

	code, resp := theAlertsProvider.CreateAlert(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientId + " function: CreateAlert requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// EditAlerts
// @Tags space Alerts API's
// @Description Edit Alerts
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.EditAlertsReq true "Alerts"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/alerts/editAlerts [POST]
func EditAlerts(c *gin.Context) {
	var reqParams models.EditAlertsReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("EditAlerts (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("EditAlerts (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientId, " deviceId: ", requestH.DeviceId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("EditAlerts (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientId, " deviceId: ", requestH.DeviceId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("EditAlerts (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("EditAlerts (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("EditAlerts (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientId, "requestId:", requestH.RequestId, " deviceId: ", requestH.DeviceId, " clientVersion: ", requestH.ClientVersion)

	code, resp := theAlertsProvider.EditAlerts(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientId + " function: EditAlerts requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// PauseAlerts
// @Tags space Alerts API's
// @Description PauseAlerts
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.PauseAlertsReq true "Alerts"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/alerts/pauseAlerts [POST]
func PauseAlerts(c *gin.Context) {
	var reqParams models.PauseAlertsReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("PauseAlerts (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("PauseAlerts (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientId, " deviceId: ", requestH.DeviceId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("PauseAlerts (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientId, " deviceId: ", requestH.DeviceId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("PauseAlerts (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("PauseAlerts (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("PauseAlerts (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientId, "requestId:", requestH.RequestId, " deviceId: ", requestH.DeviceId, " clientVersion: ", requestH.ClientVersion)

	code, resp := theAlertsProvider.PauseAlerts(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientId + " function: PauseAlerts requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DeleteAlerts
// @Tags space Alerts API's
// @Description Delete Alerts
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.DeleteAlertsReq true "Alerts"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/alerts/deleteAlerts [POST]
func DeleteAlerts(c *gin.Context) {
	var reqParams models.DeleteAlertsReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("DeleteAlerts (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("DeleteAlerts (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientId, " deviceId: ", requestH.DeviceId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("DeleteAlerts (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientId, " deviceId: ", requestH.DeviceId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("DeleteAlerts (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("DeleteAlerts (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("DeleteAlerts (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientId, "requestId:", requestH.RequestId, " deviceId: ", requestH.DeviceId, " clientVersion: ", requestH.ClientVersion)

	code, resp := theAlertsProvider.DeleteAlerts(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientId + " function: DeleteAlerts requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetAlerts
// @Tags space Alerts API's
// @Description Get Alerts
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GetAlertsReq true "Alerts"
// @Success 200 {object} apihelpers.APIRes{data=models.AlersRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/alerts/getAlerts [POST]
func GetAlerts(c *gin.Context) {

	var reqParams models.GetAlertsReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetAlerts (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetAlerts (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientId, " deviceId: ", requestH.DeviceId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetAlerts (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientId, " deviceId: ", requestH.DeviceId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("GetAlerts (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("GetAlerts (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("GetAlerts (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientId, "requestId:", requestH.RequestId, " deviceId: ", requestH.DeviceId, " clientVersion: ", requestH.ClientVersion)

	code, resp := theAlertsProvider.GetAlerts(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientId + " function: GetAlerts requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
