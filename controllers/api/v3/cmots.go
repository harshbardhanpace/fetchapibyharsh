package v3

import (
	"encoding/json"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var theCmotsProviderV3 models.CMOTSProviderV3

func InitCmotsProviderV3(provider models.CMOTSProviderV3) {
	defer models.HandlePanic()
	theCmotsProviderV3 = provider
}

// FetchFinancialsV3
// @Tags space cmots V3
// @Description Fetch Financials V3
// @Param ClientId header string false "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.FetchFinancialsReq true "cmots"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchFinancialsV3Res}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/cmots/fetchFinancials [POST]
func FetchFinancialsV3(c *gin.Context) {
	var reqParams models.FetchFinancialsReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FetchFinancialsV3 (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchFinancialsV3 (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FetchFinancialsV3 (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("FetchFinancialsV3 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theCmotsProviderV3.FetchFinancialsV3(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: FetchFinancialsV3 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
