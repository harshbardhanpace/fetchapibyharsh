package v1

import (
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var theSessionInfoProvider models.SessionInfoProvider

func InitSessionInfoProvider(provider models.SessionInfoProvider) {
	defer models.HandlePanic()
	theSessionInfoProvider = provider
}

// SessionInfo
// @Tags space Info V1
// @Description Session Info
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Success 200 {object} apihelpers.APIRes{data=models.SessionInfoRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/info/session [GET]
func SessionInfo(c *gin.Context) {
	var reqParams models.SessionInfoReq

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("SessionInfo (controller), Empty Device Type requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("SessionInfo (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("SessionInfo (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theSessionInfoProvider.SessionInfo(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: SessionInfo requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
