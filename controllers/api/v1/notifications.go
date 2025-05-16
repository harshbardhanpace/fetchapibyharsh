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
)

var theNotificationsProvider models.NotificationsProvider

func InitNotificationsProvider(provider models.NotificationsProvider) {
	defer models.HandlePanic()
	theNotificationsProvider = provider
}

// FetchAdminMessages
// @Tags space Notifications V1
// @Description This API returns all the admin Messages by clientId
// @Param ClientId header string false "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.FetchAdminMessageRequest true "FetchAdminMessageRequest"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchAdminMessageRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/notifications/adminMessages [POST]
func FetchAdminMessages(c *gin.Context) {
	var reqParams models.FetchAdminMessageRequest
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FetchAdminMessages  (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchAdminMessages  (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("FetchAdminMessages (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("FetchAdminMessages (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("FetchAdminMessages (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := theNotificationsProvider.FetchAdminMessages(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientId + " function: FetchAdminMessages requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// NotificationUpdates
// @Tags space Notifications V1
// @Description Notification Updates
// @Param ClientId header string false "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Success 200 {object} apihelpers.APIRes{data=models.NotificationUpdatesRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/notifications/notificationUpdates [GET]
func NotificationUpdates(c *gin.Context) {
	var reqParams models.NotificationUpdatesReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	clientID := requestH.ClientId // c.Query("ClientId")
	if clientID == "" {
		loggerconfig.Error("NotificationUpdates (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	reqParams.ClientId = clientID
	if requestH.DeviceType == "" {
		loggerconfig.Error("NotificationUpdates  (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("NotificationUpdates (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("NotificationUpdates (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("NotificationUpdates (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := theNotificationsProvider.NotificationUpdates(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientId + " function: NotificationUpdates requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
