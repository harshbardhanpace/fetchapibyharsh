package v1

import (
	"encoding/json"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
)

var theTechnicalIndicatorsProvider models.TechnicalIndicators

func InitTechnicalIndicatorsProvider(provider models.TechnicalIndicators) {
	defer models.HandlePanic()
	theTechnicalIndicatorsProvider = provider
}

// TechnicalIndicatorsValues
// @Tags space Technical Indicators V1
// @Description This API returns technical indicator values for things like sma, ema, macd, rsi, etc.,
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param request body models.TechnicalIndicatorsValuesReq true "TechnicalIndicatorsValues"
// @Success 200 {object} apihelpers.APIRes{data=models.TechnicalIndicatorsValuesRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/technicalIndicators/technicalIndicatorsValues [POST]
func TechnicalIndicatorsValues(c *gin.Context) {
	var reqParams models.TechnicalIndicatorsValuesReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("TechnicalIndicatorsValues (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("TechnicalIndicatorsValues (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("TechnicalIndicatorsValues (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := theTechnicalIndicatorsProvider.TechnicalIndicatorsValues(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: TechnicalIndicatorsValues requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
