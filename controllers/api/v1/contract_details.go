package v1

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

var theContractDetailsProvider models.ContractDetailsProvider

func InitContractDetailsProvider(provider models.ContractDetailsProvider) {
	defer models.HandlePanic()
	theContractDetailsProvider = provider
}

// SearchScrip
// @Tags space contractdetails V1
// @Description Search Scrip
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.SearchScripRequest true "contractdetails"
// @Success 200 {object} apihelpers.APIRes{data=models.SearchScripResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/contractdetails/searchScrip [POST]
func SearchScrip(c *gin.Context) {
	var reqParams models.SearchScripRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("SearchScrip (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("SearchScrip (controller), Empty Device Type requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("SearchScrip (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("SearchScrip (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theContractDetailsProvider.SearchScrip(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: SearchScrip requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ScripInfo
// @Tags space contractdetails V1
// @Description Scrip Info
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.ScripInfoRequest true "contractdetails"
// @Success 200 {object} apihelpers.APIRes{data=models.ScripInfoResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/contractdetails/scripInfo [POST]
func ScripInfo(c *gin.Context) {
	var reqParams models.ScripInfoRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ScripInfo (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("ScripInfo (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("ScripInfo (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("ScripInfo (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " clientVersion:", requestH.ClientVersion)
	code, resp := theContractDetailsProvider.ScripInfo(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + "clientId: " + requestH.ClientId + " function: ScripInfo requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
