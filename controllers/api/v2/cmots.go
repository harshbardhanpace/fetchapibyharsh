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

var theCmotsProviderV2 models.CMOTSProviderV2

func InitCmotsProviderV2(provider models.CMOTSProviderV2) {
	defer models.HandlePanic()
	theCmotsProviderV2 = provider
}

// StocksOnNewsV2
// @Tags space cmots V2
// @Description Stocks On News V2
// @Param ClientId header string false "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.StocksOnNewsV2Req true "cmots"
// @Success 200 {object} apihelpers.APIRes{data=[]models.StocksOnNewsResponseData}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/cmots/stocksOnNewsV2 [POST]
func StocksOnNewsV2(c *gin.Context) {
	var reqParams models.StocksOnNewsV2Req
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("StocksOnNewsV2 (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("StocksOnNewsV2 (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("StocksOnNewsV2 (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("StocksOnNewsV2 (controller), reqParams :", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theCmotsProviderV2.StocksOnNewsV2(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: StocksOnNewsV2 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchFinancialsV2
// @Tags space cmots V2
// @Description Fetch Financials V2
// @Param ClientId header string false "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.FetchFinancialsReq true "cmots"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchFinancialsV2Res}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/cmots/fetchFinancialsV2 [POST]
func FetchFinancialsV2(c *gin.Context) {
	var reqParams models.FetchFinancialsReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FetchFinancialsV2 (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchFinancialsV2 (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FetchFinancialsV2 (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("FetchFinancialsV2 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theCmotsProviderV2.FetchFinancialsV2(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: FetchFinancialsV2 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchPeersV2
// @Tags space cmots V2
// @Description Fetch Peers V2
// @Param ClientId header string false "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.FetchPeersV2Req true "cmots"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchPeersV2Res}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/cmots/fetchPeers [POST]
func FetchPeersV2(c *gin.Context) {
	var reqParams models.FetchPeersV2Req
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FetchPeersV2 (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchPeersV2 (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FetchPeersV2 (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("FetchPeersV2 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theCmotsProviderV2.FetchPeersV2(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: FetchPeersV2 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchSectorList
// @Tags space cmots V2
// @Description Fetch Sector List V2
// @Param ClientId header string false "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param sectorCode query string false "sectorCode"
// @Success 200 {object} apihelpers.APIRes{data=[]models.SectorList}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/cmots/fetchSectorListV2 [GET]
func FetchSectorListV2(c *gin.Context) {
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchSectorListV2 (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	sectorCode := c.Query("sectorCode")

	code, resp := theCmotsProviderV2.FetchSectorListV2(sectorCode, requestH)
	logDetail := "FetchSectorListV2 (controller) clientId: " + requestH.ClientId + " requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchSectorWiseCompany
// @Tags space cmots V2
// @Description Fetch Sector Wise Company V2
// @Param ClientId header string false "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.FetchSectorWiseCompanyReqV2 true "cmots"
// @Success 200 {object} apihelpers.APIRes{data=[]models.SectorWiseCompany}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/cmots/fetchSectorWiseCompanyV2 [POST]
func FetchSectorWiseCompanyV2(c *gin.Context) {
	var reqParams models.FetchSectorWiseCompanyReqV2
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FetchSectorWiseCompanyV2 (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchSectorWiseCompanyV2 (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FetchSectorWiseCompanyV2 (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("FetchSectorWiseCompanyV2 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theCmotsProviderV2.FetchSectorWiseCompanyV2(reqParams, requestH)
	logDetail := "FetchSectorWiseCompanyV2 (controller) clientId: " + requestH.ClientId + " requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
