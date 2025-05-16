package v3

import (
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var theOptionChainProviderV3 models.OptionChainProvider

func InitOptionChainProviderV3(provider models.OptionChainProvider) {
	defer models.HandlePanic()
	theOptionChainProviderV3 = provider
}

// FetchOptionChain
// @Tags space optionchain V3
// @Description Fetch Option Chain V3
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param exchange query string true "exchange Query Parameter" format(string)
// @Param token query string true "token Query Parameter" format(string)
// @Param num query string true "num Query Parameter" format(string)
// @Param price query string true "price Query Parameter" format(string)
// @Success 200 {object} apihelpers.APIRes{data=models.FetchOptionChainResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/optionchain/fetchOptionChain [GET]
func FetchOptionChain(c *gin.Context) {
	var reqParams models.FetchOptionChainRequest
	tokenStr := c.Query("token")
	numStr := c.Query("num")
	priceStr := c.Query("price")
	exchangeStr := c.Query("exchange")

	token, errToken := strconv.Atoi(tokenStr)
	num, errNum := strconv.Atoi(numStr)
	price, errPrice := strconv.ParseFloat(priceStr, 64)

	if errToken != nil || errNum != nil || errPrice != nil || exchangeStr == "" {
		loggerconfig.Error("FetchOptionChain V3 (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	reqParams.Token = token
	reqParams.Num = num
	reqParams.Price = price
	reqParams.Exchange = exchangeStr

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchOptionChain V3 (controller), Empty Device Type clientID: ", requestH.ClientId, "requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FetchOptionChain V3 (controller), Error validating struct: ", err, "clientID: ", requestH.ClientId, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("FetchOptionChain V3 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "clientID: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := theOptionChainProviderV3.FetchOptionChain(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: FetchOptionChain v3 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchFuturesChain
// @Tags space optionchain V3
// @Description Fetch Futures Chain V3
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param token query string true "token Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.FetchFuturesChainRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/optionchain/fetchFuturesChain [GET]
func FetchFuturesChain(c *gin.Context) {
	var reqParams models.FetchFuturesChainReq
	token := c.Query("token")
	if token == "" {
		loggerconfig.Error("GetAlertsV2 (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.Token = token

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchFuturesChain V3 (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FetchFuturesChain V3 (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("FetchFuturesChain V3 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId)
	code, resp := theOptionChainProviderV3.FetchFuturesChain(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: FetchFuturesChain v3 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
