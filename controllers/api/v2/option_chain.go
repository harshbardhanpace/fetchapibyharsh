package v2

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

var theOptionChainProviderV2 models.OptionChainProviderV2

func InitOptionChainProviderV2(provider models.OptionChainProviderV2) {
	defer models.HandlePanic()
	theOptionChainProviderV2 = provider
}

// FetchOptionChainV2
// @Tags space optionchain V2
// @Description Fetch Option Chain V2
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.FetchOptionChainV2Request true "optionchain"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchOptionChainResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/optionchain/fetchOptionChain [POST]
func FetchOptionChainV2(c *gin.Context) {
	var reqParams models.FetchOptionChainV2Request
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FetchOptionChainV2 (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchOptionChainV2 (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FetchOptionChainV2 (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("FetchOptionChainV2 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theOptionChainProviderV2.FetchOptionChainV2(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: FetchOptionChainV2 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchOptionChainByExpiryV2
// @Tags space optionchain V2
// @Description Fetch Option Chain By Expiry V2
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.FetchOptionChainByExpiryV2Request true "optionchain"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchOptionChainResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/optionchain/fetchOptionChainByExpiry [POST]
func FetchOptionChainByExpiryV2(c *gin.Context) {
	var reqParams models.FetchOptionChainByExpiryV2Request
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FetchOptionChainByExpiryV2 (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchOptionChainByExpiryV2 (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FetchOptionChainByExpiryV2 (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("FetchOptionChainByExpiryV2 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theOptionChainProviderV2.FetchOptionChainByExpiryV2(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: FetchOptionChainByExpiryV2 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
