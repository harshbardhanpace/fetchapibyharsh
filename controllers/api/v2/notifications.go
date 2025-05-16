package v2

import (
	"net/http"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
)

var theNotificationsProviderV2 models.NotificationsProvider

func InitNotificationsProviderV2(provider models.NotificationsProvider) {
	defer models.HandlePanic()
	theNotificationsProviderV2 = provider
}

// FetchAdminMessages
// @Tags space Notifications V2
// @Description This API returns all the admin Messages V2
// @Param ClientId header string false "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param clientId query string true "clientId Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.FetchAdminMessageRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/notifications/adminMessages [GET]
func FetchAdminMessages(c *gin.Context) {
	var reqParams models.FetchAdminMessageRequest
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	clientID := c.Query("clientId")
	if clientID == "" {
		loggerconfig.Error("FetchAdminMessages (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.ClientId = clientID

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchAdminMessages V2(controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("FetchAdminMessages V2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("FetchAdminMessages V2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("FetchAdminMessages V2(controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := theNotificationsProviderV2.FetchAdminMessages(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: FetchAdminMessages requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
