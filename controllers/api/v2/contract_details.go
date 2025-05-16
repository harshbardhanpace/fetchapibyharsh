package v2

import (
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var theContractDetailsProviderV2TL models.ContractDetailsProvider

func InitContractDetailsProviderV2TL(provider models.ContractDetailsProvider) {
	defer models.HandlePanic()
	theContractDetailsProviderV2TL = provider
}

// SearchScrip V2TL
// @Tags space contractdetails V2TL
// @Description Search Scrip V2TL
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param key query string true "key Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.SearchScripResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/contractdetails/searchScripTL [GET]
func SearchScrip(c *gin.Context) {
	var reqParams models.SearchScripRequest
	key := c.Query("key")
	if key == "" {
		loggerconfig.Error("SearchScrip V2TL (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.Key = key
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("SearchScrip V2TL (controller), Empty Device Type requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("SearchScrip V2TL (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("SearchScrip V2TL (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId)
	code, resp := theContractDetailsProviderV2TL.SearchScrip(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: SearchScrip requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ScripInfo V2TL
// @Tags space contractdetails V2TL
// @Description Scrip Info V2TL
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param exchange query string true "exchange Query Parameter" dataType(string)
// @Param info query string true "info Query Parameter" dataType(string)
// @Param token query string true "token Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.ScripInfoResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/contractdetails/scripInfoTL [GET]
func ScripInfo(c *gin.Context) {
	var reqParams models.ScripInfoRequest
	exchange := c.Query("exchange")
	info := c.Query("info")
	token := c.Query("token")
	if exchange == "" || info == "" || token == "" {
		loggerconfig.Error("ScripInfo V2TL (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.Exchange = exchange
	reqParams.Info = info
	reqParams.Token = token

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("ScripInfo V2TL (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("ScripInfo V2TL (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("ScripInfo V2TL (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId)
	code, resp := theContractDetailsProviderV2TL.ScripInfo(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: ScripInfo requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
