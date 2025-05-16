package v2

import (
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
	"strings"

	"github.com/gin-gonic/gin"
)

var pinsProviderV2 models.PinsProvider


func InitPinsProviderV2(provider models.PinsProvider) {
	defer models.HandlePanic()
	pinsProviderV2 = provider
}
// FetchPins
// @Tags space pins V2
// @Description Fetch Pins
// @Param ClientId header string false "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Success 200 {object} apihelpers.APIRes{data=models.PinsResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/pins/fetchPins [GET]
func FetchPinsV2(c *gin.Context) {
	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("FetchPinsV2 (controller), error decoding header, error:", err)
	}

	if reqH.DeviceType == "" {
		loggerconfig.Error("FetchPinsV2 (controller), Empty Device Type requestId: ", reqH.RequestId, "ClientID: ", reqH.ClientId, " deviceId: ", reqH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("FetchPinsV2 (controller), FetchPinsV2 requestId: ", reqH.RequestId, "ClientID: ", reqH.ClientId, " deviceId: ", reqH.DeviceId)

	reqH.ClientId = strings.ToUpper(reqH.ClientId)
	var fetchPinsReq models.PinsRequest
	if reqH.ClientType == constants.GUESTUSERTYPE{
		fetchPinsReq.ClientId = constants.GUESTUSERTYPE
	} else {
		fetchPinsReq.ClientId = reqH.ClientId
	}

	loggerconfig.Info("FetchPinsV2 (controller), reqParams:", helpers.LogStructAsJSON(fetchPinsReq), " uccId: ", fetchPinsReq.ClientId, "requestId:", reqH.RequestId)
	code, resp := pinsProviderV2.FetchPins(fetchPinsReq, reqH)

	logDetail := "clientId: " + reqH.ClientId + " function: FetchPinsV2 requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
