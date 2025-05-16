package portfolioanalyzer

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	apihelpers "space/apiHelpers"
	"space/business/pockets"
	"space/business/tradelab"
	"space/constants"
	"space/db"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
)

type PAObj struct {
	tradeLabURL string
}

func InitPortfolioAnalyzer() PAObj {
	defer models.HandlePanic()

	paObj := PAObj{
		tradeLabURL: constants.TLURL,
	}
	return paObj
}

var FetchHoldingsData = func(url string, clientID string, reqH models.ReqHeader) ([]models.HoldingsData, error) {
	return FetchHoldingsDataActual(url, clientID, reqH)
}

func FetchHoldingsDataActual(url string, clientID string, reqH models.ReqHeader) ([]models.HoldingsData, error) {
	var response []models.HoldingsData

	payload := new(bytes.Buffer)
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	if err != nil {
		loggerconfig.Error("FetchHoldingsData call api error =", err, " uccId:", clientID, " requestId:", reqH.RequestId)
		return response, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := tradelab.TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == tradelab.TLERROR {
		loggerconfig.Error("FetchHoldingsData res error =", tlErrorRes.Message, " uccId:", clientID, " requestId:", reqH.RequestId)
		return response, err
	}

	tlFetchDematHoldingsResponse := tradelab.TradeLabFetchDematHoldingsResponse{}
	json.Unmarshal([]byte(string(body)), &tlFetchDematHoldingsResponse)
	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("FetchHoldingsData tl status not ok =", tlFetchDematHoldingsResponse.Message, " uccId:", clientID, " requestId:", reqH.RequestId)
		return response, err
	}

	loggerconfig.Info("FetchHoldingsData tl resp=", helpers.LogStructAsJSON(tlFetchDematHoldingsResponse), " uccId:", clientID, " requestId:", reqH.RequestId)

	for i := 0; i < len(tlFetchDematHoldingsResponse.Data.Holdings); i++ {
		var holding models.HoldingsData
		holding.Exchange = tlFetchDematHoldingsResponse.Data.Holdings[i].Exchange
		holding.Symbol = tlFetchDematHoldingsResponse.Data.Holdings[i].Symbol
		holding.Isin = tlFetchDematHoldingsResponse.Data.Holdings[i].Isin
		holding.LTP = tlFetchDematHoldingsResponse.Data.Holdings[i].Ltp
		holding.Quantity = tlFetchDematHoldingsResponse.Data.Holdings[i].Quantity
		holding.Token = strconv.Itoa(tlFetchDematHoldingsResponse.Data.Holdings[i].Token)
		response = append(response, holding)
	}

	return response, nil
}

func (obj PAObj) HoldingsWeightages(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes
	var response []models.Holdings

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("HoldingsWeightages call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var totalValue float64
	for i := 0; i < len(holdings); i++ {
		totalValue = totalValue + (holdings[i].LTP * float64(holdings[i].Quantity))
	}
	for i := 0; i < len(holdings); i++ {
		var ith models.Holdings
		ith.PercentageOfPortfolio = float64(holdings[i].Quantity) * holdings[i].LTP / totalValue * 100.0
		ith.ValueOfHolding = float64(holdings[i].Quantity) * holdings[i].LTP
		ith.Isin = holdings[i].Isin
		ith.Token = holdings[i].Token
		ith.TradingSymbol = holdings[i].Symbol
		ith.Exchange = holdings[i].Exchange
		ith.SectorCode, ith.SectorName, err = db.GetPgObj().FetchSector(holdings[i].Isin)
		if err != nil {
			loggerconfig.Error("HoldingsWeightages DB fetch error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		response = append(response, ith)
	}

	loggerconfig.Info("HoldingsWeightages Successful, response:", helpers.LogStructAsJSON(response), " uccId:", req.ClientId, " requestId:", reqH.RequestId)

	apiRes.Data = response
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func CoVariance(benchmark []models.TLCandleData, pricesIth []models.TLCandleData) float64 {
	n := float64(len(pricesIth) - 1)
	var retIth, avgIth, retBenchmark, avgBenchmark, covariance float64

	for i := len(pricesIth) - 1; i > 0; i-- {
		avgBenchmark = avgBenchmark + ((benchmark[i].Close/benchmark[i-1].Close)-1)*100
	}
	avgBenchmark = avgBenchmark / float64(len(benchmark))

	for i := len(pricesIth) - 1; i > 0; i-- {
		retBenchmark = retBenchmark + ((pricesIth[i].Close/pricesIth[i-1].Close)-1)*100
	}
	retBenchmark = retBenchmark / float64(len(pricesIth))

	for i := len(pricesIth) - 1; i > 1; i-- {
		avgIth = ((benchmark[i].Close / benchmark[i-1].Close) - 1) * 100
		retIth = ((pricesIth[i].Close / pricesIth[i-1].Close) - 1) * 100
		covariance = covariance + (avgIth-avgBenchmark)*(retIth-retBenchmark)/n
	}
	return covariance
}

func Variance(benchmark []models.TLCandleData) float64 {
	nBen := float64(len(benchmark) - 1)
	var avgBenchmark, variance float64

	for i := len(benchmark) - 1; i > 0; i-- {
		avgBenchmark = avgBenchmark + ((benchmark[i].Close/benchmark[i-1].Close)-1)*100
	}
	avgBenchmark = avgBenchmark / float64(nBen)

	for i := len(benchmark) - 1; i > 0; i-- {
		varianceCalculation := ((((benchmark[i].Close / benchmark[i-1].Close) - 1) * 100) - avgBenchmark)
		variance = variance + varianceCalculation*varianceCalculation
	}
	variance = variance / float64(nBen)

	return variance
}

var CalculateBeta = func(benchmark []models.TLCandleData, pricesIth []models.TLCandleData) float64 {
	return CalculateBetaActual(benchmark, pricesIth)
}

func CalculateBetaActual(benchmark []models.TLCandleData, pricesIth []models.TLCandleData) float64 {
	covariance := CoVariance(benchmark, pricesIth)
	variance := Variance(benchmark)

	beta := covariance / variance
	return beta
}

func (obj PAObj) PortfolioBeta(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes
	var response models.PortfolioBeta

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("PortfolioBeta call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	starttime := helpers.GetCurrentTimeInIST().Unix() - constants.ONEYEAR
	endtime := helpers.GetCurrentTimeInIST().Unix()

	urlBenchmark := obj.tradeLabURL + tradelab.Charts + "?exchange=" + constants.IndexExchange + "&token=" + constants.NiftyToken + "&candletype=" + constants.DayWise + "&starttime=" + strconv.Itoa(int(starttime)) + "&endtime=" + strconv.Itoa(int(endtime)) + "&data_duration=" + constants.DataDuration

	var status bool
	benchmarkData, status := pockets.CallTLforData(urlBenchmark, reqH)
	if !status {
		apiRes.Message = constants.ErrorCodeMap[constants.TLChartDataFetchFailed]
		apiRes.Status = false
		return http.StatusInternalServerError, apiRes
	}

	var totalValue float64
	for i := 0; i < len(holdings); i++ {
		totalValue = totalValue + (holdings[i].LTP * float64(holdings[i].Quantity))
	}
	response.PortfolioTotalValue = totalValue

	var chartData string
	for i := 0; i < len(holdings); i++ {
		chartData = obj.tradeLabURL + tradelab.Charts + "?exchange=" + holdings[i].Exchange + "&token=" + holdings[i].Token + "&candletype=" + constants.DayWise + "&starttime=" + strconv.Itoa(int(starttime)) + "&endtime=" + strconv.Itoa(int(endtime)) + "&data_duration=" + constants.DataDuration
		var PricesIth []models.TLCandleData
		var individualBeta models.IndividualBeta
		PricesIth, status := pockets.CallTLforData(chartData, reqH)
		if !status {
			apiRes.Message = constants.ErrorCodeMap[constants.TLChartDataFetchFailed]
			apiRes.Status = false
			return http.StatusInternalServerError, apiRes
		}
		betaIth := CalculateBeta(benchmarkData, PricesIth)
		response.PortfolioBeta = response.PortfolioBeta + betaIth*(float64(holdings[i].Quantity)*holdings[i].LTP/totalValue)
		individualBeta.Beta = betaIth
		individualBeta.Exchange = holdings[i].Exchange
		individualBeta.Isin = holdings[i].Isin
		individualBeta.Token = holdings[i].Token
		individualBeta.Symbol = holdings[i].Symbol
		response.IndividualBeta = append(response.IndividualBeta, individualBeta)
	}

	loggerconfig.Info("PortfolioBeta Successful, response:", helpers.LogStructAsJSON(response), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = response
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) PortfolioPE(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes
	var response models.PortfolioPE

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("PortfolioPE call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var totalValue float64
	var portfolioPE float64
	for i := 0; i < len(holdings); i++ {
		totalValue = totalValue + (holdings[i].LTP * float64(holdings[i].Quantity))
	}
	fetchIndividualPE := make([]models.IndividualPE, 0)
	for i := 0; i < len(holdings); i++ {
		var ith models.IndividualPE
		ith.Isin = holdings[i].Isin
		ith.Token = holdings[i].Token
		ith.TradingSymbol = holdings[i].Symbol
		ith.Pe, err = db.GetPgObj().FetchPE(holdings[i].Isin)
		if err != nil {
			loggerconfig.Error("PortfolioPE DB fetch error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		fetchIndividualPE = append(fetchIndividualPE, ith)
		portfolioPE = portfolioPE + (float64(holdings[i].Quantity)*holdings[i].LTP/totalValue)*ith.Pe
	}
	response.PortfolioPE = portfolioPE
	response.IndividualPE = fetchIndividualPE

	loggerconfig.Info("PortfolioPE Successful, response:", helpers.LogStructAsJSON(response), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = response
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) PortfolioDE(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes
	var response models.PortfolioDE

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("PortfolioDE call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var totalValue float64
	var portfolioDE float64
	for i := 0; i < len(holdings); i++ {
		totalValue = totalValue + (holdings[i].LTP * float64(holdings[i].Quantity))
	}
	fetchIndividualDE := make([]models.IndividualDE, 0)
	for i := 0; i < len(holdings); i++ {
		var ith models.IndividualDE
		ith.Isin = holdings[i].Isin
		ith.Token = holdings[i].Token
		ith.TradingSymbol = holdings[i].Symbol
		ith.De, err = db.GetPgObj().FetchDE(holdings[i].Isin)
		if err != nil {
			loggerconfig.Error("PortfolioDE DB fetch error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		fetchIndividualDE = append(fetchIndividualDE, ith)
		portfolioDE = portfolioDE + (float64(holdings[i].Quantity)*holdings[i].LTP/totalValue)*ith.De
	}
	response.PortfolioDE = portfolioDE
	response.IndividualDE = fetchIndividualDE

	loggerconfig.Info("PortfolioDE Successful, response:", helpers.LogStructAsJSON(response), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = response
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) HighPledgedPromoterHoldings(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("HighPledgedPromoterHoldings HoldingsWeightages call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	holdingMap := make(map[string]models.HoldingRedFlagData)
	var allIsin models.AllIsin
	for i := 0; i < len(holdings); i++ {
		allIsin.Isin = append(allIsin.Isin, holdings[i].Isin)
		var holdingHpph models.HoldingRedFlagData
		holdingHpph.Isin = holdings[i].Isin
		holdingHpph.TradingSymbol = holdings[i].Symbol
		holdingHpph.Token = holdings[i].Token
		holdingMap[holdings[i].Isin] = holdingHpph
	}

	commonPledgeHolding, err := db.GetPgObj().FetchHighPledgePromoterHoldingMatchData(allIsin)

	var highPledgePromoterHoldingRes models.HighPledgePromoterHoldingRes
	for i := 0; i < len(commonPledgeHolding.HighPledgePromoterHoldingAll); i++ {
		highPledgePromoterHoldingRes.Holdings = append(highPledgePromoterHoldingRes.Holdings, holdingMap[commonPledgeHolding.HighPledgePromoterHoldingAll[i].Isin])
	}

	loggerconfig.Info("HighPledgedPromoterHoldings Successful, response:", helpers.LogStructAsJSON(highPledgePromoterHoldingRes), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = highPledgePromoterHoldingRes
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) AdditionalSurveillanceMeasure(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("AdditionalSurveillanceMeasure HoldingsWeightages call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	holdingMap := make(map[string]models.HoldingRedFlagData)
	var allIsin models.AllIsin
	for i := 0; i < len(holdings); i++ {
		allIsin.Isin = append(allIsin.Isin, holdings[i].Isin)
		var holdingASM models.HoldingRedFlagData
		holdingASM.Isin = holdings[i].Isin
		holdingASM.TradingSymbol = holdings[i].Symbol
		holdingASM.Token = holdings[i].Token
		holdingMap[holdings[i].Isin] = holdingASM
	}

	commonAdditionalSurveillanceMeasure, err := db.GetPgObj().FetchAdditionalSurveillanceMeasureData(allIsin)

	var additionalSurveillanceMeasureRes models.AdditionalSurveillanceMeasureRes
	for i := 0; i < len(commonAdditionalSurveillanceMeasure.AdditionalSurveillanceMeasureAll); i++ {
		additionalSurveillanceMeasureRes.Holdings = append(additionalSurveillanceMeasureRes.Holdings, holdingMap[commonAdditionalSurveillanceMeasure.AdditionalSurveillanceMeasureAll[i].Isin])
	}

	loggerconfig.Info("AdditionalSurveillanceMeasure Successful, response:", helpers.LogStructAsJSON(additionalSurveillanceMeasureRes), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = additionalSurveillanceMeasureRes
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) GradedSurveillanceMeasure(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("GradedSurveillanceMeasure HoldingsWeightages call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	holdingMap := make(map[string]models.HoldingRedFlagData)
	var allIsin models.AllIsin
	for i := 0; i < len(holdings); i++ {
		allIsin.Isin = append(allIsin.Isin, holdings[i].Isin)
		var holdingASM models.HoldingRedFlagData
		holdingASM.Isin = holdings[i].Isin
		holdingASM.TradingSymbol = holdings[i].Symbol
		holdingASM.Token = holdings[i].Token
		holdingMap[holdings[i].Isin] = holdingASM
	}

	commonGradedSurveillanceMeasure, err := db.GetPgObj().FetchGradedSurveillanceMeasureData(allIsin)

	var gradedSurveillanceMeasureRes models.GradedSurveillanceMeasureRes
	for i := 0; i < len(commonGradedSurveillanceMeasure.AllGradedSurveillanceMeasureAll); i++ {
		gradedSurveillanceMeasureRes.Holdings = append(gradedSurveillanceMeasureRes.Holdings, holdingMap[commonGradedSurveillanceMeasure.AllGradedSurveillanceMeasureAll[i].Isin])
	}

	loggerconfig.Info("GradedSurveillanceMeasure Successful, response:", helpers.LogStructAsJSON(gradedSurveillanceMeasureRes), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = gradedSurveillanceMeasureRes
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) HighDefaultProbability(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	return http.StatusOK, apiRes
}

func (obj PAObj) LowROE(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("LowROE HoldingsWeightages call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	holdingMap := make(map[string]models.RoeHoldingRedFlagData)
	var allIsin models.AllIsin
	for i := 0; i < len(holdings); i++ {
		allIsin.Isin = append(allIsin.Isin, holdings[i].Isin)
		var holdingASM models.RoeHoldingRedFlagData
		holdingASM.Isin = holdings[i].Isin
		holdingASM.TradingSymbol = holdings[i].Symbol
		holdingASM.Token = holdings[i].Token
		holdingMap[holdings[i].Isin] = holdingASM
	}

	commonGradedRoeData, err := db.GetPgObj().FetchRoeData(allIsin)

	var lowROERes models.LowROERes
	for i := 0; i < len(commonGradedRoeData.LowRoeAll); i++ {
		roeHoldingRedFlagData, present := holdingMap[commonGradedRoeData.LowRoeAll[i].Isin]
		if present {
			roeHoldingRedFlagData.Roe = commonGradedRoeData.LowRoeAll[i].Roe
			if roeHoldingRedFlagData.Roe < 10 {
				lowROERes.Holdings = append(lowROERes.Holdings, roeHoldingRedFlagData)
			}
		}

	}
	loggerconfig.Info("LowROE Successful, response:", helpers.LogStructAsJSON(lowROERes), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = lowROERes
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) LowProfitGrowth(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("LowProfitGrowth HoldingsWeightages call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	holdingMap := make(map[string]models.ProfitabilityGrowthRedFlagData)
	var allIsin models.AllIsin
	for i := 0; i < len(holdings); i++ {
		allIsin.Isin = append(allIsin.Isin, holdings[i].Isin)
		var profitabilityGrowthRedFlagData models.ProfitabilityGrowthRedFlagData
		profitabilityGrowthRedFlagData.Isin = holdings[i].Isin
		profitabilityGrowthRedFlagData.TradingSymbol = holdings[i].Symbol
		profitabilityGrowthRedFlagData.Token = holdings[i].Token
		holdingMap[holdings[i].Isin] = profitabilityGrowthRedFlagData
	}

	commonLowProfitGrowthData, err := db.GetPgObj().FetchLowProfitGrowthData(allIsin)

	var pGRedFlagRes models.ProfitabilityGrowthRedFlagRes
	for i := 0; i < len(commonLowProfitGrowthData.ProfitabilityGrowthAll); i++ {
		dataRedFlag, present := holdingMap[commonLowProfitGrowthData.ProfitabilityGrowthAll[i].Isin]
		if present {
			yearZero := commonLowProfitGrowthData.ProfitabilityGrowthAll[i].YZero
			yearFour := commonLowProfitGrowthData.ProfitabilityGrowthAll[i].YFour

			netProfit := math.Pow((yearFour/yearZero), 1.0/5.0) - 1
			if netProfit < 5 {
				dataRedFlag.NetProfit = netProfit
				pGRedFlagRes.Holdings = append(pGRedFlagRes.Holdings, dataRedFlag)
			}
		}
	}

	loggerconfig.Info("LowProfitGrowth Successful, response:", helpers.LogStructAsJSON(pGRedFlagRes), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = pGRedFlagRes
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) HoldingStockContribution(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("HoldingStockContribution HoldingsWeightages call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var totalValuePortfolio float64
	for i := 0; i < len(holdings); i++ {
		totalValuePortfolio += holdings[i].LTP * float64(holdings[i].Quantity)
	}

	var holdingStockContributionRes models.HoldingStockContributionRes

	for i := 0; i < len(holdings); i++ {
		var holdingStockContributionData models.HoldingStockContributionData
		holdingStockContributionData.Isin = holdings[i].Isin
		holdingStockContributionData.TradingSymbol = holdings[i].Symbol
		holdingStockContributionData.Token = holdings[i].Token
		holdingStockContributionData.StockInvestedPrice = holdings[i].LTP * float64(holdings[i].Quantity)
		holdingStockContributionData.PercentageShare = (holdingStockContributionData.StockInvestedPrice / totalValuePortfolio) * 100
		holdingStockContributionRes.Holdings = append(holdingStockContributionRes.Holdings, holdingStockContributionData)
	}

	sort.Slice(holdingStockContributionRes.Holdings, func(i, j int) bool {
		return holdingStockContributionRes.Holdings[i].PercentageShare > holdingStockContributionRes.Holdings[j].PercentageShare
	})

	for i := 0; i < len(holdingStockContributionRes.Holdings) && i < 5; i++ {
		holdingStockContributionRes.TopFiveHoldingPrice += holdingStockContributionRes.Holdings[i].StockInvestedPrice
	}

	holdingStockContributionRes.TopFiveCombinedPercentageShare = (holdingStockContributionRes.TopFiveHoldingPrice / totalValuePortfolio) * 100

	loggerconfig.Info("HoldingStockContribution Successful, response:", helpers.LogStructAsJSON(holdingStockContributionRes), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = holdingStockContributionRes
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) InvestmentSector(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("InvestmentSector HoldingsWeightages call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var totalValuePortfolio float64
	investmentDataMap := make(map[string]models.InvestmentData)
	var allIsin models.AllIsin
	for i := 0; i < len(holdings); i++ {
		allIsin.Isin = append(allIsin.Isin, holdings[i].Isin)
		var investmentData models.InvestmentData
		investmentData.Isin = holdings[i].Isin
		investmentData.LTP = holdings[i].LTP
		investmentData.Quantity = holdings[i].Quantity
		investmentDataMap[holdings[i].Isin] = investmentData
		totalValuePortfolio += holdings[i].LTP * float64(holdings[i].Quantity)
	}

	companyInfo, err := db.GetPgObj().FetchCompanyMasterData(allIsin)
	if err != nil {
		loggerconfig.Error("InvestmentSector call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	investmentSectorDataMap := make(map[string]models.InvestmentSectorData)

	for i := 0; i < len(companyInfo.CompanyMasterAll); i++ {
		isin := companyInfo.CompanyMasterAll[i].Isin
		if isin != "" {
			_, contains := investmentSectorDataMap[isin]
			if contains {
				var updatedInvestedDetails models.InvestmentSectorData
				updatedInvestedDetails.SectorCode = companyInfo.CompanyMasterAll[i].Sectorcode
				updatedInvestedDetails.SectorName = companyInfo.CompanyMasterAll[i].Sectorname
				updatedInvestedDetails.InvestedValue += investmentDataMap[isin].LTP * float64(investmentDataMap[isin].Quantity)
				investmentSectorDataMap[isin] = updatedInvestedDetails
			} else {
				var newInvestedDetails models.InvestmentSectorData
				newInvestedDetails.SectorCode = companyInfo.CompanyMasterAll[i].Sectorcode
				newInvestedDetails.SectorName = companyInfo.CompanyMasterAll[i].Sectorname
				newInvestedDetails.InvestedValue = investmentDataMap[isin].LTP * float64(investmentDataMap[isin].Quantity)
				investmentSectorDataMap[isin] = newInvestedDetails
			}
		}
	}

	var investmentSectorRes models.InvestmentSectorRes
	for _, investmentSectorData := range investmentSectorDataMap {
		investmentSectorData.InvestedPercentage = (investmentSectorData.InvestedValue / totalValuePortfolio) * 100
		investmentSectorRes.InvestmentSector = append(investmentSectorRes.InvestmentSector, investmentSectorData)
	}

	sort.Slice(investmentSectorRes.InvestmentSector, func(i, j int) bool {
		return investmentSectorRes.InvestmentSector[i].InvestedPercentage > investmentSectorRes.InvestmentSector[j].InvestedPercentage
	})

	loggerconfig.Info("InvestmentSector Successful, response:", helpers.LogStructAsJSON(investmentSectorRes), " uccId:", req.ClientId, " requestId:", reqH.RequestId)

	apiRes.Data = investmentSectorRes
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) DeclineInPromoterHolding(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("DeclineInPromoterHolding HoldingsWeightages call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	holdingMap := make(map[string]models.DeclineInPromoterHoldingRedFlagData)
	var allIsin models.AllIsin
	for i := 0; i < len(holdings); i++ {
		allIsin.Isin = append(allIsin.Isin, holdings[i].Isin)
		var declineInPromoterHoldingRedFlagData models.DeclineInPromoterHoldingRedFlagData
		declineInPromoterHoldingRedFlagData.Isin = holdings[i].Isin
		declineInPromoterHoldingRedFlagData.TradingSymbol = holdings[i].Symbol
		declineInPromoterHoldingRedFlagData.Token = holdings[i].Token
		holdingMap[holdings[i].Isin] = declineInPromoterHoldingRedFlagData
	}

	commonDeclineInPromoterHoldingData, err := db.GetPgObj().FetchDeclineInPromoterHoldingData(allIsin)
	if err != nil {
		loggerconfig.Error("DeclineInPromoterHolding, error in fetching FetchDeclineInPromoterHoldingData ", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var declineInPromoterHoldingRes models.DeclineInPromoterHoldingRedFlagRes
	for i := 0; i < len(commonDeclineInPromoterHoldingData.DeclineInPromoterHolding); i++ {
		dataRedFlag, present := holdingMap[commonDeclineInPromoterHoldingData.DeclineInPromoterHolding[i].Isin]
		if present {
			currentQuarterTPPS := commonDeclineInPromoterHoldingData.DeclineInPromoterHolding[i].CurrentQuarterTPPS
			previousQuarterTPPS := commonDeclineInPromoterHoldingData.DeclineInPromoterHolding[i].PreviousQuarterTPPS

			netDeclineInPromoterHolding := previousQuarterTPPS - currentQuarterTPPS
			if netDeclineInPromoterHolding > constants.TWO {
				dataRedFlag.NetDecline = netDeclineInPromoterHolding
				declineInPromoterHoldingRes.Holdings = append(declineInPromoterHoldingRes.Holdings, dataRedFlag)
			}
		}
	}

	loggerconfig.Info("DeclineInPromoterHolding Successful, response:", helpers.LogStructAsJSON(declineInPromoterHoldingRes), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = declineInPromoterHoldingRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) InterestCoverageRatio(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("InterestCoverageRatio HoldingsWeightages call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	holdingMap := make(map[string]models.HoldingRedFlagData)
	var allIsin models.AllIsin
	for i := 0; i < len(holdings); i++ {
		allIsin.Isin = append(allIsin.Isin, holdings[i].Isin)
		var InterestCoverageRatioRedFlagData models.HoldingRedFlagData
		InterestCoverageRatioRedFlagData.Isin = holdings[i].Isin
		InterestCoverageRatioRedFlagData.TradingSymbol = holdings[i].Symbol
		InterestCoverageRatioRedFlagData.Token = holdings[i].Token
		holdingMap[holdings[i].Isin] = InterestCoverageRatioRedFlagData
	}

	commonInterestCoverageRatioData, err := db.GetPgObj().FetchInterestCoverageRatioData(allIsin)
	if err != nil {
		loggerconfig.Error("InterestCoverageRatio, error in fetching FetchInterestCoverageRatioData ", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var icrRedFlagRes models.InterestCoverageRatioRedFlagRes
	for i := 0; i < len(commonInterestCoverageRatioData.InterestCoverageRatioData); i++ {
		dataRedFlag, present := holdingMap[commonInterestCoverageRatioData.InterestCoverageRatioData[i].Isin]
		if present {
			financeCost := commonInterestCoverageRatioData.InterestCoverageRatioData[i].FinanceCost
			profitBeforeTax := commonInterestCoverageRatioData.InterestCoverageRatioData[i].ProfitBeforeTax
			interestCoverRatio := commonInterestCoverageRatioData.InterestCoverageRatioData[i].InterestCoverRatio

			icr := (financeCost + profitBeforeTax) / financeCost
			twiceOfIcr := icr * constants.TWO

			if twiceOfIcr < interestCoverRatio {
				icrRedFlagRes.Holdings = append(icrRedFlagRes.Holdings, dataRedFlag)
			}
		}
	}

	loggerconfig.Info("InterestCoverageRatio Successful, response:", helpers.LogStructAsJSON(icrRedFlagRes), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = icrRedFlagRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) DeclineInRevenueAndProfit(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("DeclineInRevenueAndProfit HoldingsWeightages call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	holdingMap := make(map[string]models.HoldingRedFlagData)
	var allIsin models.AllIsin
	for i := 0; i < len(holdings); i++ {
		allIsin.Isin = append(allIsin.Isin, holdings[i].Isin)
		var declineInRevenueAndProfit models.HoldingRedFlagData
		declineInRevenueAndProfit.Isin = holdings[i].Isin
		declineInRevenueAndProfit.TradingSymbol = holdings[i].Symbol
		declineInRevenueAndProfit.Token = holdings[i].Token
		holdingMap[holdings[i].Isin] = declineInRevenueAndProfit
	}

	commonDeclineInRevenueAndProfitData, err := db.GetPgObj().DeclineInRevenueAndProfitData(allIsin)
	if err != nil {
		loggerconfig.Error("DeclineInRevenueAndProfit, error in fetching DeclineInRevenueAndProfitData ", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var declineInRPRedFlagRes models.DeclineInRevenueAndProfitRedFlagRes
	for i := 0; i < len(commonDeclineInRevenueAndProfitData.Holding); i++ {
		dataRedFlag, present := holdingMap[commonDeclineInRevenueAndProfitData.Holding[i].Isin]
		if present {
			declineInRPRedFlagRes.Holdings = append(declineInRPRedFlagRes.Holdings, dataRedFlag)
		}
	}

	loggerconfig.Info("DeclineInRevenueAndProfit Successful, response:", helpers.LogStructAsJSON(declineInRPRedFlagRes), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = declineInRPRedFlagRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) LowNetWorth(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("LowNetWorth call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	holdingMap := make(map[string]models.NetWorthDataRedFlagData)
	var allIsin models.AllIsin
	for i := 0; i < len(holdings); i++ {
		allIsin.Isin = append(allIsin.Isin, holdings[i].Isin)
		var netWorthData models.NetWorthDataRedFlagData
		netWorthData.Isin = holdings[i].Isin
		netWorthData.TradingSymbol = holdings[i].Symbol
		netWorthData.Token = holdings[i].Token
		holdingMap[holdings[i].Isin] = netWorthData
	}

	commonLowNetWorthData, err := db.GetPgObj().LowNetWorthData(allIsin)
	if err != nil {
		loggerconfig.Error("LowNetWorth, error in fetching DeclineInRevenueAndProfitData ", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var lowNetWorthRedFlagRes models.LowNetWorthDataRedFlagRes
	for i := 0; i < len(commonLowNetWorthData.NetWorthData); i++ {
		dataRedFlag, present := holdingMap[commonLowNetWorthData.NetWorthData[i].Isin]
		if present {
			dataRedFlag.NetWorth = commonLowNetWorthData.NetWorthData[i].NetWorth
			lowNetWorthRedFlagRes.Holdings = append(lowNetWorthRedFlagRes.Holdings, dataRedFlag)
		}
	}

	loggerconfig.Info("LowNetWorth Successful, response:", helpers.LogStructAsJSON(lowNetWorthRedFlagRes), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = lowNetWorthRedFlagRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) DeclineInRevenue(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("DeclineInRevenue HoldingsWeightages call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	holdingMap := make(map[string]models.HoldingRedFlagData)
	var allIsin models.AllIsin
	for i := 0; i < len(holdings); i++ {
		allIsin.Isin = append(allIsin.Isin, holdings[i].Isin)
		var declineInRevenue models.HoldingRedFlagData
		declineInRevenue.Isin = holdings[i].Isin
		declineInRevenue.TradingSymbol = holdings[i].Symbol
		declineInRevenue.Token = holdings[i].Token
		holdingMap[holdings[i].Isin] = declineInRevenue
	}

	commonDeclineInRevenueData, err := db.GetPgObj().DeclineInRevenueData(allIsin)
	if err != nil {
		loggerconfig.Error("DeclineInRevenue, error in fetching DeclineInRevenueData ", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var declineInRevenueRedFlagRes models.DeclineInRevenueRedFlagRes
	for i := 0; i < len(commonDeclineInRevenueData.Holding); i++ {
		dataRedFlag, present := holdingMap[commonDeclineInRevenueData.Holding[i].Isin]
		if present {
			declineInRevenueRedFlagRes.Holdings = append(declineInRevenueRedFlagRes.Holdings, dataRedFlag)
		}
	}

	loggerconfig.Info("DeclineInRevenue Successful, response:", helpers.LogStructAsJSON(declineInRevenueRedFlagRes), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = declineInRevenueRedFlagRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) PromoterPledge(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("PromoterPledge HoldingsWeightages call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	holdingMap := make(map[string]models.HoldingRedFlagData)
	var allIsin models.AllIsin
	for i := 0; i < len(holdings); i++ {
		allIsin.Isin = append(allIsin.Isin, holdings[i].Isin)
		var promoterPledge models.HoldingRedFlagData
		promoterPledge.Isin = holdings[i].Isin
		promoterPledge.TradingSymbol = holdings[i].Symbol
		promoterPledge.Token = holdings[i].Token
		holdingMap[holdings[i].Isin] = promoterPledge
	}

	commonPromoterPledgeData, err := db.GetPgObj().PromoterPledgeData(allIsin)
	if err != nil {
		loggerconfig.Error("PromoterPledge, error in fetching PromoterPledgeData ", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var promoterPledgeRedFlagRes models.PromoterPledgeRedFlagRes
	for i := 0; i < len(commonPromoterPledgeData.Holding); i++ {
		dataRedFlag, present := holdingMap[commonPromoterPledgeData.Holding[i].Isin]
		if present {
			promoterPledgeRedFlagRes.Holdings = append(promoterPledgeRedFlagRes.Holdings, dataRedFlag)
		}
	}

	loggerconfig.Info("PromoterPledge Successful, response:", helpers.LogStructAsJSON(promoterPledgeRedFlagRes), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = promoterPledgeRedFlagRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) PennyStocks(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("PennyStocks HoldingsWeightages call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	holdingMap := make(map[string]models.PennyStocksHoldingRedFlagData)
	var allIsin models.AllIsin
	for i := 0; i < len(holdings); i++ {
		allIsin.Isin = append(allIsin.Isin, holdings[i].Isin)
		var pennyStocks models.PennyStocksHoldingRedFlagData
		pennyStocks.Isin = holdings[i].Isin
		pennyStocks.TradingSymbol = holdings[i].Symbol
		pennyStocks.Token = holdings[i].Token
		holdingMap[holdings[i].Isin] = pennyStocks
	}

	commonPennyStocksData, err := db.GetPgObj().PennyStocksData(allIsin)
	if err != nil {
		loggerconfig.Error("PromoterPledge, error in fetching PromoterPledgeData ", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var pennyStocksRedFlagRes models.PennyStocksRedFlagRes
	for i := 0; i < len(commonPennyStocksData.Holding); i++ {
		dataRedFlag, present := holdingMap[commonPennyStocksData.Holding[i].Isin]
		if present {
			dataRedFlag.MarketCap = commonPennyStocksData.Holding[i].MarketCap
			pennyStocksRedFlagRes.Holdings = append(pennyStocksRedFlagRes.Holdings, dataRedFlag)
		}
	}

	loggerconfig.Info("PromoterPledge Successful, response:", helpers.LogStructAsJSON(pennyStocksRedFlagRes), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = pennyStocksRedFlagRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) StockReturn(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("StockReturn call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	holdingMap := make(map[string]models.HoldingRedFlagData)
	var allIsin models.AllIsin
	for i := 0; i < len(holdings); i++ {
		allIsin.Isin = append(allIsin.Isin, holdings[i].Isin)
		var holdingData models.HoldingRedFlagData
		holdingData.Isin = holdings[i].Isin
		holdingData.TradingSymbol = holdings[i].Symbol
		holdingData.Token = holdings[i].Token
		holdingMap[holdings[i].Isin] = holdingData
	}

	var stockReturneRes models.StockReturneRes

	stockReturnData, err := db.GetPgObj().StockReturnData(allIsin)
	if err != nil {
		loggerconfig.Error("StockReturn, error in fetching StockReturnData ", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	for i := 0; i < len(stockReturnData.StockReturn); i++ {
		redFlagData, present := holdingMap[stockReturnData.StockReturn[i].Isin]
		if present {
			stockReturneRes.Holdings = append(stockReturneRes.Holdings, redFlagData)
		}
	}

	loggerconfig.Info("StockReturn Red FLag Successful, response:", helpers.LogStructAsJSON(stockReturneRes), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = stockReturneRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) NiftyVsPortfolio(req models.NiftyVsPortfolioReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	var portfolioAnalyzerReq models.PortfolioAnalyzerReq
	portfolioAnalyzerReq.ClientId = req.ClientId

	status, res := obj.PortfolioBeta(portfolioAnalyzerReq, reqH)
	if status != http.StatusOK {
		loggerconfig.Error("NiftyVsPortfolio in PortfolioBeta status != 200", status, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	portfolioBeta, ok := res.Data.(models.PortfolioBeta)
	if !ok {
		loggerconfig.Error("In NiftyVsPortfolio PortfolioBeta interface parsing error", ok, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var niftyVsPortfolioRes models.NiftyVsPortfolioRes
	niftyVsPortfolioRes.PortfolioMovement = req.MovementOfNiftyPercentage * portfolioBeta.PortfolioBeta
	niftyVsPortfolioRes.ExpectedGainAmountForNiftyMovement = portfolioBeta.PortfolioBeta * portfolioBeta.PortfolioTotalValue
	for i := 0; i < len(portfolioBeta.IndividualBeta); i++ {
		var individualMoment models.IndividualMovement
		individualMoment.Exchange = portfolioBeta.IndividualBeta[i].Exchange
		individualMoment.Isin = portfolioBeta.IndividualBeta[i].Isin
		individualMoment.Symbol = portfolioBeta.IndividualBeta[i].Symbol
		individualMoment.Token = portfolioBeta.IndividualBeta[i].Token
		individualMoment.Movement = req.MovementOfNiftyPercentage * portfolioBeta.IndividualBeta[i].Beta
		niftyVsPortfolioRes.IndividualMovement = append(niftyVsPortfolioRes.IndividualMovement, individualMoment)
	}

	loggerconfig.Info("NiftyVsPortfolio Successful, response:", helpers.LogStructAsJSON(niftyVsPortfolioRes), " uccId:", req.ClientId, " requestId:", reqH.RequestId)

	apiRes.Data = niftyVsPortfolioRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) ChangeInInstitutionalHolding(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("ChangeInInstitutionalHolding HoldingsWeightages call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	holdingMap := make(map[string]models.ChangeInInstitutionalHoldingRedFlagData)
	var allIsin models.AllIsin
	for i := 0; i < len(holdings); i++ {
		allIsin.Isin = append(allIsin.Isin, holdings[i].Isin)
		var changeInInstitutionalHoldingRedFlagData models.ChangeInInstitutionalHoldingRedFlagData
		changeInInstitutionalHoldingRedFlagData.Isin = holdings[i].Isin
		changeInInstitutionalHoldingRedFlagData.TradingSymbol = holdings[i].Symbol
		changeInInstitutionalHoldingRedFlagData.Token = holdings[i].Token
		holdingMap[holdings[i].Isin] = changeInInstitutionalHoldingRedFlagData
	}

	commonChangeInInstitutionalHoldingData, err := db.GetPgObj().FetchChangeInInstitutionalHoldingData(allIsin)
	if err != nil {
		loggerconfig.Error("ChangeInInstitutionalHolding, error in fetching FetchChangeInInstitutionalHoldingData ", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var changeInInstitutionalHoldingRes models.ChangeInInstitutionalHoldingRedFlagRes
	for i := 0; i < len(commonChangeInInstitutionalHoldingData.ChangeInInstitutionalHoldingDb); i++ {
		dataRedFlag, present := holdingMap[commonChangeInInstitutionalHoldingData.ChangeInInstitutionalHoldingDb[i].Isin]
		if present {
			dataRedFlag.DiffrenceInInstitutionalHolding = commonChangeInInstitutionalHoldingData.ChangeInInstitutionalHoldingDb[i].DiffrenceInInstitutionalHolding
			changeInInstitutionalHoldingRes.Holdings = append(changeInInstitutionalHoldingRes.Holdings, dataRedFlag)
		}
	}

	loggerconfig.Info("ChangeInInstitutionalHolding Successful, response:", helpers.LogStructAsJSON(changeInInstitutionalHoldingRes), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = changeInInstitutionalHoldingRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) RoeAndStockReturn(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("RoeAndStockReturn HoldingsWeightages call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	holdingMap := make(map[string]models.RoeAndStockReturnHoldingRedFlagData)
	var allIsin models.AllIsin
	for i := 0; i < len(holdings); i++ {
		allIsin.Isin = append(allIsin.Isin, holdings[i].Isin)
		var roeAndStockReturnHoldingRedFlagData models.RoeAndStockReturnHoldingRedFlagData
		roeAndStockReturnHoldingRedFlagData.Isin = holdings[i].Isin
		roeAndStockReturnHoldingRedFlagData.TradingSymbol = holdings[i].Symbol
		roeAndStockReturnHoldingRedFlagData.Token = holdings[i].Token
		holdingMap[holdings[i].Isin] = roeAndStockReturnHoldingRedFlagData
	}

	roeAndStockReturndata, err := db.GetPgObj().FetchRoeAndStockReturnData(allIsin)
	if err != nil {
		loggerconfig.Error("RoeAndStockReturn, error in fetching FetchRoeAndStockReturnData ", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var roeAndStockReturnHoldingRedFlagRes models.RoeAndStockReturnHoldingRedFlagRes
	for i := 0; i < len(roeAndStockReturndata.RoeAndStockReturndb); i++ {
		dataRedFlag, present := holdingMap[roeAndStockReturndata.RoeAndStockReturndb[i].Isin]
		if present {
			returnOver3Years := math.Cbrt(roeAndStockReturndata.RoeAndStockReturndb[i].Return3Yrs / roeAndStockReturndata.RoeAndStockReturndb[i].Ltp)
			if returnOver3Years < 0.12 {
				year0Roe := roeAndStockReturndata.RoeAndStockReturndb[i].Y0ProfitAES / roeAndStockReturndata.RoeAndStockReturndb[i].Y0TotalShareholderFund
				year1Roe := roeAndStockReturndata.RoeAndStockReturndb[i].Y1ProfitAES / roeAndStockReturndata.RoeAndStockReturndb[i].Y1TotalShareholderFund
				year2Roe := roeAndStockReturndata.RoeAndStockReturndb[i].Y2ProfitAES / roeAndStockReturndata.RoeAndStockReturndb[i].Y2TotalShareholderFund

				avg3yearRoe := (year0Roe + year1Roe + year2Roe) / 3

				if avg3yearRoe < 0.12 {
					dataRedFlag.Return3Yrs = returnOver3Years
					dataRedFlag.AvgRoe = avg3yearRoe
					roeAndStockReturnHoldingRedFlagRes.RoeAndStockReturnHolding = append(roeAndStockReturnHoldingRedFlagRes.RoeAndStockReturnHolding, dataRedFlag)
				}
			}
		}
	}

	loggerconfig.Info("RoeAndStockReturn Successful, response:", helpers.LogStructAsJSON(roeAndStockReturnHoldingRedFlagRes), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = roeAndStockReturnHoldingRedFlagRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PAObj) IlliquidStocks(req models.PortfolioAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + tradelab.FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientId)

	holdings, err := FetchHoldingsData(url, req.ClientId, reqH)
	if err != nil {
		loggerconfig.Error("RoeAndStockReturn HoldingsWeightages call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var illiquidStocksResponse models.IlliquidStocksResponse

	for i := 0; i < len(holdings); i++ {
		var chartDataReq models.ChartDataReq
		chartDataReq.Exchange = strings.ToUpper(holdings[i].Exchange)
		chartDataReq.Token = holdings[i].Token
		chartDataReq.CandleType = constants.CandleTypeThree
		chartDataReq.StartTime = strconv.Itoa(int(helpers.GetCurrentTimeInIST().Unix() - 30*constants.UnixOneDaySeconds))
		chartDataReq.EndTime = strconv.Itoa(int(helpers.GetCurrentTimeInIST().Unix()))
		totalTradedVolumeOneMonth, avgTradedVolumeOneDay, err := volumeTraded(chartDataReq, reqH)
		if err != nil {
			loggerconfig.Error("IlliquidStocks error in getting volumeTraded err: ", err, "clientId: ", req.ClientId, " requestId:", reqH.RequestId)
			continue
		}

		if totalTradedVolumeOneMonth < constants.IlliquidStocksMinTradingVol {
			var illiquidStocksHolding models.IlliquidStocksHolding
			illiquidStocksHolding.Exchange = strings.ToUpper(holdings[i].Exchange)
			illiquidStocksHolding.Token = holdings[i].Token
			illiquidStocksHolding.Isin = holdings[i].Isin
			illiquidStocksHolding.TradingSymbol = holdings[i].Symbol
			illiquidStocksHolding.AvgVolumeDay = avgTradedVolumeOneDay
			illiquidStocksHolding.AvgVolumeMonth = totalTradedVolumeOneMonth
			illiquidStocksResponse.Holding = append(illiquidStocksResponse.Holding, illiquidStocksHolding)
		}

	}

	loggerconfig.Info("IlliquidStocks Successful, response:", helpers.LogStructAsJSON(illiquidStocksResponse), " uccId:", req.ClientId, " requestId:", reqH.RequestId)

	apiRes.Data = illiquidStocksResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func volumeTraded(req models.ChartDataReq, reqH models.ReqHeader) (float64, float64, error) {

	gainerLoserObj := tradelab.InitGainerLoserProvider()
	status, res := tradelab.GainerLoserObj.ChartData(gainerLoserObj, req, reqH)
	if status != http.StatusOK {
		loggerconfig.Error("volumeTraded ChartData status != 200", status, " requestId:", reqH.RequestId)
		return 0.0, 0.0, errors.New("volumeTraded ChartData status != 200")
	}

	chartData, ok := res.Data.(models.ChartDataResponse)
	if !ok {
		loggerconfig.Error("volumeTraded ChartData interface parsing error", ok, " requestId:", reqH.RequestId)
		return 0.0, 0.0, errors.New("volumeTraded ChartData interface parsing error")
	}

	totalVolumeOneMonth := 0.0
	avgTradedVolumeOneDay := 0.0

	for i := 0; i < len(chartData.Data.Candles); i++ {

		currVol, ok := chartData.Data.Candles[i][4].(float64)
		if !ok {
			loggerconfig.Error("volumeTraded ChartData interface parsing error", ok, " requestId:", reqH.RequestId)
			return 0.0, 0.0, errors.New("error in converstion to float")
		}

		totalVolumeOneMonth += currVol
	}

	avgTradedVolumeOneDay = totalVolumeOneMonth / float64(len(chartData.Data.Candles))

	return totalVolumeOneMonth, avgTradedVolumeOneDay, nil

}
