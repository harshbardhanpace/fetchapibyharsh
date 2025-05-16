package v1

import (
	"net/http"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
)

var theSipProvider models.SipProvider

func InitSipProvider(provider models.SipProvider) {
	defer models.HandlePanic()
	theSipProvider = provider
}

// fetchStockSips
// @Tags space sip v1
// @description fetch stock sips
// @param P-DeviceType header string true "P-DeviceType Header"
// @param P-Platform header string false "P-Platform Header"
// @param P-DeviceId header string false "P-DeviceId Header"
// @param Authorization header string true "Authorization Header"
// @param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param clientId query string true "clientId" dataType(string)
// @success 200 {object} apihelpers.APIRes{data=tradelab.TLGetSipResponse}
// @failure 400 {object} apihelpers.APIRes
// @failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/sip/fetchStockSipOrder [GET]
func FetchStockSips(c *gin.Context) {
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchStockSips (controller), Empty Device Type requestId: ", requestH.RequestId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	clientId := c.Query("clientId")
	if clientId == "" {
		loggerconfig.Error("FetchStockSips (controller), Empty ClientId requestId: ", requestH.RequestId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(clientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("FetchStockSips (controller), CheckAuthWithClient invalid authtoken", " clientId: ", clientId, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("FetchStockSips CheckAuthWithClient difference in authtoken-clientId and clientId, clientId: ", clientId, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
	}

	loggerconfig.Info("FetchStockSips (controller), uccId: ", clientId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theSipProvider.GetStockSips(clientId, requestH)
	logDetail := "clientId: " + clientId + " function: FetchStockSips requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// PlaceSipOrder
// @Tags space sip v1
// @description place sip for baskets
// @param P-DeviceType header string true "P-DeviceType Header"
// @param P-Platform header string false "P-Platform Header"
// @param P-DeviceId header string false "P-DeviceId Header"
// @param Authorization header string true "Authorization Header"
// @param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @param P-ClientVersion header string false "P-ClientVersion Header"
// @Accept json
// @Produce json
// @Param request body models.PlaceSipRequest true "Place SIP Request"
// @success 200 {object} apihelpers.APIRes
// @failure 400 {object} apihelpers.APIRes
// @failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/sip/placeSipOrder [POST]
func PlaceSipOrder(c *gin.Context) {
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("PlaceSipOrder (controller), Empty Device Type requestId: ", requestH.RequestId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	var request models.PlaceSipRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		loggerconfig.Error("PlaceSipOrder (controller), Invalid request body. Error: ", err, " requestId: ", requestH.RequestId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if request.ClientID == "" {
		loggerconfig.Error("PlaceSipOrder (controller), Empty ClientId requestId: ", requestH.RequestId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(request.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("PlaceSipOrder (controller), CheckAuthWithClient invalid authtoken", " clientId: ", request.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("PlaceSipOrder (controller) CheckAuthWithClient difference in authtoken-clientId and clientId, clientId: ", request.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("PlaceSipOrder (controller), uccId: ", request.ClientID, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theSipProvider.PlaceSipOrder(request, requestH)
	logDetail := "clientId: " + request.ClientID + " function: PlaceSipOrder requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DeleteSip
// @Tags space sip v1
// @description delete sip
// @param P-DeviceType header string true "P-DeviceType Header"
// @param P-Platform header string false "P-Platform Header"
// @param P-DeviceId header string false "P-DeviceId Header"
// @param Authorization header string true "Authorization Header"
// @param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @param P-ClientVersion header string false "P-ClientVersion Header"
// @Param clientId path string true "Client ID"
// @Param sipId path string true "SIP ID"
// @success 200 {object} apihelpers.APIRes
// @failure 400 {object} apihelpers.APIRes
// @failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/sip/deleteSipOrder/{clientId}/{sipId} [DELETE]
func DeleteSipOrder(c *gin.Context) {
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("DeleteSipOrder (controller), Empty Device Type requestId: ", requestH.RequestId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	clientId := c.Param("clientId")
	sipId := c.Param("sipId")

	if clientId == "" || sipId == "" {
		loggerconfig.Error("DeleteSipOrder (controller), Empty ClientId or SipId. ClientId: ", clientId, " SipId: ", sipId, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(clientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("DeleteSipOrder (controller), CheckAuthWithClient invalid authtoken", " clientId: ", clientId, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("DeleteSipOrder (controller) CheckAuthWithClient difference in authtoken-clientId and clientId, clientId: ", clientId, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("DeleteSipOrder (controller), uccId: ", clientId, " sipId: ", sipId, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theSipProvider.DeleteSipOrder(clientId, sipId, requestH)
	logDetail := "clientId: " + clientId + " sipId: " + sipId + " function: DeleteSipOrder requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ModifySip
// @Tags space sip v1
// @description modify sip
// @param P-DeviceType header string true "P-DeviceType Header"
// @param P-Platform header string false "P-Platform Header"
// @param P-DeviceId header string false "P-DeviceId Header"
// @param Authorization header string true "Authorization Header"
// @param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @param P-ClientVersion header string false "P-ClientVersion Header"
// @Accept json
// @Produce json
// @Param request body models.ModifySipRequest true "Modify SIP Request"
// @success 200 {object} apihelpers.APIRes
// @failure 400 {object} apihelpers.APIRes
// @failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/sip/modifySipOrder [PUT]
func ModifySipOrder(c *gin.Context) {
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("ModifySipOrder (controller), Empty Device Type requestId: ", requestH.RequestId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	var request models.ModifySipRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		loggerconfig.Error("ModifySipOrder (controller), Invalid request body. Error: ", err, " requestId: ", requestH.RequestId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if request.ClientID == "" || request.ID == "" {
		loggerconfig.Error("ModifySipOrder (controller), Empty ClientId or SipId. ClientId: ", request.ClientID, " SipId: ", request.ID, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(request.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("ModifySipOrder (controller), CheckAuthWithClient invalid authtoken", " clientId: ", request.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("ModifySipOrder (controller) CheckAuthWithClient difference in authtoken-clientId and clientId, clientId: ", request.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("ModifySipOrder (controller), uccId: ", request.ClientID, " sipId: ", request.ID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theSipProvider.ModifySipOrder(request, requestH)
	logDetail := "clientId: " + request.ClientID + " sipId: " + request.ID + " function: ModifySipOrder requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// UpdateSipStatus
// @Tags space sip v1
// @description update sip status (Active/Paused)
// @param P-DeviceType header string true "P-DeviceType Header"
// @Param ClientId header string true "ClientId"
// @param P-Platform header string false "P-Platform Header"
// @param P-DeviceId header string false "P-DeviceId Header"
// @param Authorization header string true "Authorization Header"
// @param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @param P-ClientVersion header string false "P-ClientVersion Header"
// @Accept json
// @Produce json
// @Param request body models.UpdateSipStatusRequest true "Update SIP Status Request"
// @success 200 {object} apihelpers.APIRes
// @failure 400 {object} apihelpers.APIRes
// @failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/sip/updateSipStatus [PUT]
func UpdateSipStatus(c *gin.Context) {
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("UpdateSipStatus (controller), Empty Device Type requestId: ", requestH.RequestId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	var request models.UpdateSipStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		loggerconfig.Error("UpdateSipStatus (controller), Invalid request body. Error: ", err, " requestId: ", requestH.RequestId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("UpdateSipStatus (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("UpdateSipStatus (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("UpdateSipStatus (controller), uccId: ", requestH.ClientId, " sipId: ", request.ID, " action: ", request.Action, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theSipProvider.UpdateSipStatus(request, requestH)
	logDetail := "sipId: " + request.ID + " action: " + request.Action + " function: UpdateSipStatus requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
