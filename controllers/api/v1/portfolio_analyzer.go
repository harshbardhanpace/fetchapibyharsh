package v1

import (
	"encoding/json"
	"net/http"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var thePortfolioAnalyzerProvider models.PortfolioAnalyzer

func InitPortfolioAnalyzerProvider(provider models.PortfolioAnalyzer) {
	defer models.HandlePanic()
	thePortfolioAnalyzerProvider = provider
}

// HoldingsWeightages
// @Tags space Portfolio Analyzer V1
// @Description Holdings Weightages
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.Holdings}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/holdingsWeightages [POST]
func HoldingsWeightages(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("HoldingsWeightages (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("HoldingsWeightages (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("HoldingsWeightages (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("HoldingsWeightages (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("HoldingsWeightages (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("HoldingsWeightages (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.HoldingsWeightages(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: HoldingsWeightages requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// PortfolioBeta
// @Tags space Portfolio Analyzer V1
// @Description Portfolio Beta
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.PortfolioBeta}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/portfolioBeta [POST]
func PortfolioBeta(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("PortfolioBeta (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("PortfolioBeta (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("PortfolioBeta (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("PortfolioBeta (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("PortfolioBeta (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("PortfolioBeta (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.PortfolioBeta(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: PortfolioBeta requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// PortfolioPE
// @Tags space Portfolio Analyzer V1
// @Description Portfolio PE
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.PortfolioPE}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/portfolioPE [POST]
func PortfolioPE(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("PortfolioPE (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("PortfolioPE (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("PortfolioPE (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("PortfolioPE (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("PortfolioPE (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("PortfolioPE (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.PortfolioPE(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: PortfolioPE requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// PortfolioDE
// @Tags space Portfolio Analyzer V1
// @Description Portfolio DE
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.PortfolioDE}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/portfolioDE [POST]
func PortfolioDE(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("PortfolioDE (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("PortfolioDE (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("PortfolioDE (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("PortfolioDE (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("PortfolioDE (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("PortfolioDE (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.PortfolioDE(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: PortfolioDE requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// HighPledgedPromoterHoldings
// @Tags space Portfolio Analyzer V1
// @Description HighPledged Promoter Holdings
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.HighPledgePromoterHoldingRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/highPledgedPromoterHoldings [POST]
func HighPledgedPromoterHoldings(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("HighPledgedPromoterHoldings (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("HighPledgedPromoterHoldings (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("HighPledgedPromoterHoldings (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("HighPledgedPromoterHoldings (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("HighPledgedPromoterHoldings (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("HighPledgedPromoterHoldings (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.HighPledgedPromoterHoldings(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: HighPledgedPromoterHoldings requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// AdditionalSurveillanceMeasure
// @Tags space Portfolio Analyzer V1
// @Description Additional Surveillance Measure
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.AdditionalSurveillanceMeasureRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/additionalSurveillanceMeasure [POST]
func AdditionalSurveillanceMeasure(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("AdditionalSurveillanceMeasure (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("AdditionalSurveillanceMeasure (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("AdditionalSurveillanceMeasure (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("AdditionalSurveillanceMeasure (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("AdditionalSurveillanceMeasure (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("AdditionalSurveillanceMeasure (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.AdditionalSurveillanceMeasure(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: AdditionalSurveillanceMeasure requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GradedSurveillanceMeasure
// @Tags space Portfolio Analyzer V1
// @Description Graded Surveillance Measure
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.GradedSurveillanceMeasureRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/gradedSurveillanceMeasure [POST]
func GradedSurveillanceMeasure(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GradedSurveillanceMeasure (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("GradedSurveillanceMeasure (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GradedSurveillanceMeasure (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("GradedSurveillanceMeasure (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("GradedSurveillanceMeasure (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("GradedSurveillanceMeasure (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.GradedSurveillanceMeasure(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: GradedSurveillanceMeasure requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// HighDefaultProbability
// @Tags space Portfolio Analyzer V1
// @Description High Default Probability
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.HighDefaultProbability}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/highDefaultProbability [POST]
func HighDefaultProbability(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("HighDefaultProbability (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("HighDefaultProbability (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("HighDefaultProbability (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("HighDefaultProbability (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("HighDefaultProbability (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("HighDefaultProbability (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.HighDefaultProbability(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: HighDefaultProbability requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// LowROE
// @Tags space Portfolio Analyzer V1
// @Description Low ROE
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.LowROERes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/lowROE [POST]
func LowROE(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("LowROE (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("LowROE (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("LowROE (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("LowROE (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("LowROE (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("LowROE (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.LowROE(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: LowROE requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// LowProfitGrowth
// @Tags space Portfolio Analyzer V1
// @Description Low Profit Growth
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.ProfitabilityGrowthRedFlagRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/lowProfitGrowth [POST]
func LowProfitGrowth(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("LowProfitGrowth (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("LowProfitGrowth (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("LowProfitGrowth (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("LowProfitGrowth (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("LowProfitGrowth (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("LowProfitGrowth (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.LowProfitGrowth(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: LowProfitGrowth requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// HoldingStockContribution
// @Tags space Portfolio Analyzer V1
// @Description Holding Stock Contribution
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.HoldingStockContributionRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/holdingStockContribution [POST]
func HoldingStockContribution(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("HoldingStockContribution (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("HoldingStockContribution (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("HoldingStockContribution (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("HoldingStockContribution (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("HoldingStockContribution (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("HoldingStockContribution (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.HoldingStockContribution(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: HoldingStockContribution requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// InvestmentSector
// @Tags space Portfolio Analyzer V1
// @Description Investment Sector
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.InvestmentSectorRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/investmentSector [POST]
func InvestmentSector(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("InvestmentSector (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("InvestmentSector (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("InvestmentSector (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("InvestmentSector (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("InvestmentSector (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("InvestmentSector (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.InvestmentSector(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: InvestmentSector requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DeclineInPromoterHolding
// @Tags space Portfolio Analyzer V1
// @Description Decline In Promoter Holding
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.DeclineInPromoterHoldingRedFlagRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/declineInPromoterHolding [POST]
func DeclineInPromoterHolding(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("DeclineInPromoterHolding (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("DeclineInPromoterHolding (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("DeclineInPromoterHolding (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("DeclineInPromoterHolding (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("DeclineInPromoterHolding (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("DeclineInPromoterHolding (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.DeclineInPromoterHolding(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: DeclineInPromoterHolding requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// InterestCoverageRatio
// @Tags space Portfolio Analyzer V1
// @Description Interest Coverage Ratio
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.InterestCoverageRatioRedFlagRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/interestCoverageRatio [POST]
func InterestCoverageRatio(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("InterestCoverageRatio (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("DeclineInPromoterHolding (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("InterestCoverageRatio (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("DeclineInPromoterHolding (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("InterestCoverageRatio (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("InterestCoverageRatio (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.InterestCoverageRatio(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: InterestCoverageRatio requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DeclineInRevenueAndProfit
// @Tags space Portfolio Analyzer V1
// @Description Decline In Revenue And Profit
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.DeclineInRevenueAndProfitRedFlagRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/declineInRevenueAndProfit [POST]
func DeclineInRevenueAndProfit(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("DeclineInRevenueAndProfit (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("DeclineInRevenueAndProfit (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("DeclineInRevenueAndProfit (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("DeclineInRevenueAndProfit (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("DeclineInRevenueAndProfit (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("DeclineInRevenueAndProfit (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.DeclineInRevenueAndProfit(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: DeclineInRevenueAndProfit requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// LowNetWorth
// @Tags space Portfolio Analyzer V1
// @Description Net Worth
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.LowNetWorthDataRedFlagRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/lowNetWorth [POST]
func LowNetWorth(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("LowNetWorth (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("LowNetWorth (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("LowNetWorth (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("LowNetWorth (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("LowNetWorth (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("LowNetWorth (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.LowNetWorth(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: LowNetWorth requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DeclineInRevenue
// @Tags space Portfolio Analyzer V1
// @Description Decline In Revenue
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.DeclineInRevenueRedFlagRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/declineInRevenue [POST]
func DeclineInRevenue(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("DeclineInRevenue (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("DeclineInRevenue (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("DeclineInRevenue (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("DeclineInRevenue (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("DeclineInRevenue (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("DeclineInRevenue (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.DeclineInRevenue(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: DeclineInRevenue requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// PromoterPledge
// @Tags space Portfolio Analyzer V1
// @Description Promoter Pledge
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.PromoterPledgeRedFlagRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/promoterPledge [POST]
func PromoterPledge(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("PromoterPledge (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("PromoterPledge (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("PromoterPledge (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("PromoterPledge (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("PromoterPledge (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("PromoterPledge (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.PromoterPledge(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: PromoterPledge requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// PennyStocks
// @Tags space Portfolio Analyzer V1
// @Description Penny Stocks
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.PennyStocksRedFlagRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/pennyStocks [POST]
func PennyStocks(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("PennyStocks (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("PennyStocks (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("PennyStocks (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("PennyStocks (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("PennyStocks (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("PennyStocks (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.PennyStocks(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: PennyStocks requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// StockReturn
// @Tags space Portfolio Analyzer V1
// @Description Stock Return
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.StockReturneRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/stockReturn [POST]
func StockReturn(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("StockReturn (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("StockReturn (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("StockReturn (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("PennyStocks (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("StockReturn (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("StockReturn (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.StockReturn(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: StockReturn requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// NiftyVsPortfolio
// @Tags space Portfolio Analyzer V1
// @Description Nifty Vs Portfolio
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.NiftyVsPortfolioReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.NiftyVsPortfolioRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/niftyVsPortfolio [POST]
func NiftyVsPortfolio(c *gin.Context) {
	var reqParams models.NiftyVsPortfolioReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("NiftyVsPortfolio (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("NiftyVsPortfolio (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("NiftyVsPortfolio (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("NiftyVsPortfolio (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("NiftyVsPortfolio (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("NiftyVsPortfolio (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.NiftyVsPortfolio(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: NiftyVsPortfolio requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ChangeInInstitutionalHolding
// @Tags space Portfolio Analyzer V1
// @Description Change In Institutional Holding
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.ChangeInInstitutionalHoldingRedFlagRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/changeInInstitutionalHolding [POST]
func ChangeInInstitutionalHolding(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ChangeInInstitutionalHolding (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("ChangeInInstitutionalHolding (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("ChangeInInstitutionalHolding (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("ChangeInInstitutionalHolding (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("ChangeInInstitutionalHolding (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("ChangeInInstitutionalHolding (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.ChangeInInstitutionalHolding(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: ChangeInInstitutionalHolding requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// RoeAndStockReturn
// @Tags space Portfolio Analyzer V1
// @Description Roe And Stock Return
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.RoeAndStockReturnHoldingRedFlagRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/roeAndStockReturn [POST]
func RoeAndStockReturn(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("RoeAndStockReturn (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("RoeAndStockReturn (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("RoeAndStockReturn (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("RoeAndStockReturn (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("RoeAndStockReturn (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("RoeAndStockReturn (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.RoeAndStockReturn(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: RoeAndStockReturn requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// IlliquidStocks
// @Tags space Portfolio Analyzer V1
// @Description Illiquid Stocks
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PortfolioAnalyzerReq true "PortfolioAnalyzer"
// @Success 200 {object} apihelpers.APIRes{data=models.IlliquidStocksResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/portfolioAnalyzer/illiquidStocks [POST]
func IlliquidStocks(c *gin.Context) {
	var reqParams models.PortfolioAnalyzerReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("IlliquidStocks (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("IlliquidStocks (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("IlliquidStocks (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("IlliquidStocks (controller) CheckAuthWithClient invalid authtoken", " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("IlliquidStocks (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", requestH.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("IlliquidStocks (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := thePortfolioAnalyzerProvider.IlliquidStocks(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: IlliquidStocks requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
