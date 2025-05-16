package v1

import (
	"encoding/json"
	"strings"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var theOptionChainProvider models.OptionChainProvider

func InitOptionChainProvider(provider models.OptionChainProvider) {
	defer models.HandlePanic()
	theOptionChainProvider = provider
}

// FetchOptionChain
// @Tags space optionchain V1
// @Description Fetch Option Chain
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.FetchOptionChainRequest true "optionchain"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchOptionChainResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/optionchain/fetchOptionChain [POST]
func FetchOptionChain(c *gin.Context) {
	var reqParams models.FetchOptionChainRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FetchOptionChain (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchOptionChain (controller), Empty Device Type clientID: ", requestH.ClientId, "requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}
	reqParams.Exchange = strings.ToUpper(reqParams.Exchange)
	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FetchOptionChain (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("FetchOptionChain (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, "requestId:", requestH.RequestId)
	code, resp := theOptionChainProvider.FetchOptionChain(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: FetchOptionChain requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchFuturesChain
// @Tags space optionchain V1
// @Description Fetch Futures Chain by token
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.FetchFuturesChainReq true "futuresChain"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchFuturesChainRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/optionchain/fetchFuturesChain [POST]
func FetchFuturesChain(c *gin.Context) {
	var reqParams models.FetchFuturesChainReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FetchFuturesChain (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchFuturesChain (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FetchFuturesChain (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("FetchFuturesChain (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theOptionChainProvider.FetchFuturesChain(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: FetchFuturesChain requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
