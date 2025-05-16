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

var technicalIndicatorsProviderV2 models.TechnicalIndicatorsV2Provider

func InitTechnicalIndicatorsV2(provider models.TechnicalIndicatorsV2Provider) {
	defer models.HandlePanic()
	technicalIndicatorsProviderV2 = provider
}

// GetSMA
// @Tags space Technical Indicators V2
// @Description GetSMA - get simple moving averages, by input the type SMA5, SMA10, etc
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GetSMAReq true "GetSMA"
// @Success 200 {object} apihelpers.APIRes{data=models.TechnicalIndicatorsResFull}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/technicalIndicators/getSMA [POST]
func GetSMA(c *gin.Context) {
	var reqParams models.GetSMAReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetSMA (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetSMA (controller), Empty Device Type clientID: ", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetSMA (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetSMA (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := technicalIndicatorsProviderV2.GetSMA(reqParams, requestH)

	logDetail := "clientId: " + requestH.ClientId + " function: GetSMA requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetEMA
// @Tags space Technical Indicators V2
// @Description GetEMA - get simple moving averages, by input the type EMA5, EMA10, etc
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GetEMAReq true "GetEMA"
// @Success 200 {object} apihelpers.APIRes{data=models.TechnicalIndicatorsResFull}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/technicalIndicators/getEMA [POST]
func GetEMA(c *gin.Context) {
	var reqParams models.GetEMAReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetEMA (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetEMA (controller), Empty Device Type clientID: ", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetEMA (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetEMA (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := technicalIndicatorsProviderV2.GetEMA(reqParams, requestH)

	logDetail := "clientId: " + requestH.ClientId + " function: GetEMA requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetHullMA
// @Tags space Technical Indicators V2
// @Description GetHullMA - get hull moving averages
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GetHullMAReq true "GetHullMA"
// @Success 200 {object} apihelpers.APIRes{data=models.TechnicalIndicatorsResFull}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/technicalIndicators/getHullMA [POST]
func GetHullMA(c *gin.Context) {
	var reqParams models.GetHullMAReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetHullMA (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetHullMA (controller), Empty Device Type clientID: ", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetHullMA (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetHullMA (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := technicalIndicatorsProviderV2.GetHullMA(reqParams, requestH)

	logDetail := "clientId: " + requestH.ClientId + " function: GetHullMA requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetVWMA
// @Tags space Technical Indicators V2
// @Description GetVWMA - get volume weighted moving averages
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GetVWMAReq true "GetVWMA"
// @Success 200 {object} apihelpers.APIRes{data=models.TechnicalIndicatorsResFull}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/technicalIndicators/getVWMA [POST]
func GetVWMA(c *gin.Context) {
	var reqParams models.GetVWMAReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetVWMA (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetVWMA (controller), Empty Device Type clientID: ", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetVWMA (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetVWMA (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := technicalIndicatorsProviderV2.GetVWMA(reqParams, requestH)

	logDetail := "clientId: " + requestH.ClientId + " function: GetVWMA requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetRSI
// @Tags space Technical Indicators V2
// @Description Get RSI
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GetRSIReq true "GetRSIReq"
// @Success 200 {object} apihelpers.APIRes{data=models.TechnicalIndicatorsResFull}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/technicalIndicators/getRSI [POST]
func GetRSI(c *gin.Context) {
	var reqParams models.GetRSIReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetRSI (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetRSI (controller), Empty Device Type clientID: ", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetRSI (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetRSI (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := technicalIndicatorsProviderV2.GetRSI(reqParams, requestH)

	logDetail := "clientId: " + requestH.ClientId + " function: GetRSI requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetCCI
// @Tags space Technical Indicators V2
// @Description Get CCI
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GetCCIReq true "GetCCIReq"
// @Success 200 {object} apihelpers.APIRes{data=models.TechnicalIndicatorsResFull}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/technicalIndicators/getCCI [POST]
func GetCCI(c *gin.Context) {
	var reqParams models.GetCCIReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetCCI (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetCCI (controller), Empty Device Type clientID: ", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetCCI (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetCCI (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := technicalIndicatorsProviderV2.GetCCI(reqParams, requestH)

	logDetail := "clientId: " + requestH.ClientId + " function: GetCCI requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetMACD
// @Tags space Technical Indicators V2
// @Description Get MACD
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GetMACDReq true "GetMACDReq"
// @Success 200 {object} apihelpers.APIRes{data=models.TechnicalIndicatorsResFull}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/technicalIndicators/getMACD [POST]
func GetMACD(c *gin.Context) {
	var reqParams models.GetMACDReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetMACD (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetMACD (controller), Empty Device Type clientID: ", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetMACD (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetMACD (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := technicalIndicatorsProviderV2.GetMACD(reqParams, requestH)

	logDetail := "clientId: " + requestH.ClientId + " function: GetMACD requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetStochastic
// @Tags space Technical Indicators V2
// @Description Get Stochastic
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GetStochasticReq true "GetStochasticReq"
// @Success 200 {object} apihelpers.APIRes{data=models.TechnicalIndicatorsResFull}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/technicalIndicators/getStochastic [POST]
func GetStochastic(c *gin.Context) {
	var reqParams models.GetStochasticReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetStochastic (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetStochastic (controller), Empty Device Type clientID: ", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetStochastic (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetStochastic (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := technicalIndicatorsProviderV2.GetStochastic(reqParams, requestH)

	logDetail := "clientId: " + requestH.ClientId + " function: GetStochastic requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetIchimokuBaseLine
// @Tags space Technical Indicators V2
// @Description Get Ichimoku Base Line
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GetIchimokuBaseLineReq true "GetIchimokuBaseLineReq"
// @Success 200 {object} apihelpers.APIRes{data=models.TechnicalIndicatorsResFull}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/technicalIndicators/getIchimokuBaseLine [POST]
func GetIchimokuBaseLine(c *gin.Context) {
	var reqParams models.GetIchimokuBaseLineReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetIchimokuBaseLine (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetIchimokuBaseLine (controller), Empty Device Type clientID: ", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetIchimokuBaseLine (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetIchimokuBaseLine (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := technicalIndicatorsProviderV2.GetIchimokuBaseLine(reqParams, requestH)

	logDetail := "clientId: " + requestH.ClientId + " function: GetIchimokuBaseLine requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetADX
// @Tags space Technical Indicators V2
// @Description Get ADX
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GetADXReq true "GetADXReq"
// @Success 200 {object} apihelpers.APIRes{data=models.TechnicalIndicatorsResFull}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/technicalIndicators/getADX [POST]
func GetADX(c *gin.Context) {
	var reqParams models.GetADXReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetADX (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetADX (controller), Empty Device Type clientID: ", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetADX (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetADX (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := technicalIndicatorsProviderV2.GetADX(reqParams, requestH)

	logDetail := "clientId: " + requestH.ClientId + " function: GetADX requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetAwesomeOscillator
// @Tags space Technical Indicators V2
// @Description Get Awesome Oscillator
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GetAwesomeOscillatorReq true "GetAwesomeOscillatorReq"
// @Success 200 {object} apihelpers.APIRes{data=models.TechnicalIndicatorsResFull}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/technicalIndicators/getAwesomeOscillator [POST]
func GetAwesomeOscillator(c *gin.Context) {
	var reqParams models.GetAwesomeOscillatorReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetAwesomeOscillator (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetAwesomeOscillator (controller), Empty Device Type clientID: ", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetAwesomeOscillator (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetAwesomeOscillator (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := technicalIndicatorsProviderV2.GetAwesomeOscillator(reqParams, requestH)

	logDetail := "clientId: " + requestH.ClientId + " function: GetAwesomeOscillator requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetMomentum
// @Tags space Technical Indicators V2
// @Description Get Momentum
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GetMomentumReq true "GetMomentumReq"
// @Success 200 {object} apihelpers.APIRes{data=models.TechnicalIndicatorsResFull}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/technicalIndicators/getMomentum [POST]
func GetMomentum(c *gin.Context) {
	var reqParams models.GetMomentumReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetMomentum (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetMomentum (controller), Empty Device Type clientID: ", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetMomentum (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetMomentum (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := technicalIndicatorsProviderV2.GetMomentum(reqParams, requestH)

	logDetail := "clientId: " + requestH.ClientId + " function: GetMomentum requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetStochRSIFast
// @Tags space Technical Indicators V2
// @Description Get Stoch RSI Fast
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GetStochRSIFastReq true "GetStochRSIFastReq"
// @Success 200 {object} apihelpers.APIRes{data=models.TechnicalIndicatorsResFull}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/technicalIndicators/getStochRSIFast [POST]
func GetStochRSIFast(c *gin.Context) {
	var reqParams models.GetStochRSIFastReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetStochRSIFast (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetStochRSIFast (controller), Empty Device Type clientID: ", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetStochRSIFast (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetStochRSIFast (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := technicalIndicatorsProviderV2.GetStochRSIFast(reqParams, requestH)

	logDetail := "clientId: " + requestH.ClientId + " function: GetStochRSIFast requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetWilliamsRange
// @Tags space Technical Indicators V2
// @Description Get Williams Range
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GetWilliamsRangeReq true "GetWilliamsRangeReq"
// @Success 200 {object} apihelpers.APIRes{data=models.TechnicalIndicatorsResFull}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/technicalIndicators/getWilliamsRange [POST]
func GetWilliamsRange(c *gin.Context) {
	var reqParams models.GetWilliamsRangeReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetWilliamsRange (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetWilliamsRange (controller), Empty Device Type clientID: ", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetWilliamsRange (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetWilliamsRange (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := technicalIndicatorsProviderV2.GetWilliamsRange(reqParams, requestH)

	logDetail := "clientId: " + requestH.ClientId + " function: GetWilliamsRange requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetUltimateOscillator
// @Tags space Technical Indicators V2
// @Description Get Ultimate Oscillator
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GetUltimateOscillatorReq true "GetUltimateOscillatorReq"
// @Success 200 {object} apihelpers.APIRes{data=models.TechnicalIndicatorsResFull}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/technicalIndicators/getUltimateOscillator [POST]
func GetUltimateOscillator(c *gin.Context) {
	var reqParams models.GetUltimateOscillatorReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetUltimateOscillator (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetUltimateOscillator (controller), Empty Device Type clientID: ", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetUltimateOscillator (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetUltimateOscillator (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := technicalIndicatorsProviderV2.GetUltimateOscillator(reqParams, requestH)

	logDetail := "clientId: " + requestH.ClientId + " function: GetUltimateOscillator requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetAllTechnicalIndicators
// @Tags space Technical Indicators V2
// @Description Get All Technical Indicators
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GetAllTechnicalIndicatorsReq true "GetAllTechnicalIndicatorsReq"
// @Success 200 {object} apihelpers.APIRes{data=models.GetAllTechnicalIndicatorsRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/technicalIndicators/getAllTechnicalIndicators [POST]
func GetAllTechnicalIndicators(c *gin.Context) {
	var reqParams models.GetAllTechnicalIndicatorsReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetAllTechnicalIndicators (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetAllTechnicalIndicators (controller), Empty Device Type clientID: ", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetAllTechnicalIndicators (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetAllTechnicalIndicators (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := technicalIndicatorsProviderV2.GetAllTechnicalIndicators(reqParams, requestH)

	logDetail := "clientId: " + requestH.ClientId + " function: GetAllTechnicalIndicators requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
