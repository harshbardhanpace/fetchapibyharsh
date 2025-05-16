package tradelab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
)

type GainerLoserObj struct {
	tradeLabURL string
}

func InitGainerLoserProvider() GainerLoserObj {
	defer models.HandlePanic()

	gainerLoserObj := GainerLoserObj{
		tradeLabURL: constants.TLURL,
	}

	return gainerLoserObj
}

func (obj GainerLoserObj) GainerLoserNse(reqH models.ReqHeader) (int, apihelpers.APIRes) {

	url := obj.tradeLabURL + GainerLoserNse

	return GainerLoser(url, reqH)

}

func (obj GainerLoserObj) GainerLoserNiftyFifty(req models.GainersLosersMostActiveVolumeReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + GainerLoserNse + "?index=" + req.Index

	return GainerLoser(url, reqH)
}

func GainerLoser(url string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	//make payload
	payload := new(bytes.Buffer)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallApi(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "GainerLoserNiftyFifty", duration, reqH.ClientId, reqH.RequestId)
	defer res.Body.Close()
	if err != nil {
		loggerconfig.Error("gainerLoser call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("gainerLoser res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlGainerLoserRes := TradeLabTopGainerLoserResponse{}
	json.Unmarshal([]byte(string(body)), &tlGainerLoserRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("tlGainerLoserRes tl status not ok =", tlGainerLoserRes.Message, "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlGainerLoserRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	// fill up controller response
	var topGainerLoserResponse models.TopGainerLoserResponse

	allLosers := make([]models.LosersGainers, 0)

	for i := 0; i < len(tlGainerLoserRes.Data.Losers); i++ {
		var allLosersRes models.LosersGainers

		allLosersRes.TurnoverValue = tlGainerLoserRes.Data.Losers[i].TurnoverValue
		allLosersRes.TradedQuantity = tlGainerLoserRes.Data.Losers[i].TradedQuantity
		allLosersRes.Symbol = tlGainerLoserRes.Data.Losers[i].Symbol
		allLosersRes.NetPrice = tlGainerLoserRes.Data.Losers[i].NetPrice
		allLosersRes.Ltp = tlGainerLoserRes.Data.Losers[i].Ltp
		allLosersRes.LotSize = tlGainerLoserRes.Data.Losers[i].LotSize
		allLosersRes.InstrumentToken = tlGainerLoserRes.Data.Losers[i].InstrumentToken
		allLosersRes.Exchange = tlGainerLoserRes.Data.Losers[i].Exchange
		allLosersRes.CompanyName = tlGainerLoserRes.Data.Losers[i].CompanyName
		allLosersRes.ClosePrice = tlGainerLoserRes.Data.Losers[i].ClosePrice

		allLosers = append(allLosers, allLosersRes)
	}
	topGainerLoserResponse.Losers = allLosers

	allGainers := make([]models.LosersGainers, 0)

	for i := 0; i < len(tlGainerLoserRes.Data.Gainers); i++ {
		var allGainersRes models.LosersGainers

		allGainersRes.TurnoverValue = tlGainerLoserRes.Data.Gainers[i].TurnoverValue
		allGainersRes.TradedQuantity = tlGainerLoserRes.Data.Gainers[i].TradedQuantity
		allGainersRes.Symbol = tlGainerLoserRes.Data.Gainers[i].Symbol
		allGainersRes.NetPrice = tlGainerLoserRes.Data.Gainers[i].NetPrice
		allGainersRes.Ltp = tlGainerLoserRes.Data.Gainers[i].Ltp
		allGainersRes.LotSize = tlGainerLoserRes.Data.Gainers[i].LotSize
		allGainersRes.InstrumentToken = tlGainerLoserRes.Data.Gainers[i].InstrumentToken
		allGainersRes.Exchange = tlGainerLoserRes.Data.Gainers[i].Exchange
		allGainersRes.CompanyName = tlGainerLoserRes.Data.Gainers[i].CompanyName
		allGainersRes.ClosePrice = tlGainerLoserRes.Data.Gainers[i].ClosePrice

		allGainers = append(allGainers, allGainersRes)
	}
	topGainerLoserResponse.Gainers = allGainers

	loggerconfig.Info("gainerLoserRes tl resp=", helpers.LogStructAsJSON(topGainerLoserResponse), "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = topGainerLoserResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj GainerLoserObj) MostActiveVolumeNSE(reqH models.ReqHeader) (int, apihelpers.APIRes) {
	loggerconfig.Info("In MostActiveVolumeNSE, clientID: ", reqH.ClientId, "requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
	url := obj.tradeLabURL + Screeners + "/mostactivestocks"

	return MostActiveVolume(url, reqH)
}

func (obj GainerLoserObj) MostActiveVolumeDataNifty50(req models.GainersLosersMostActiveVolumeReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	loggerconfig.Info("In MostActiveVolumeDataNifty50, clientID: ", reqH.ClientId, "requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
	url := obj.tradeLabURL + Screeners + "/mostactivestocks" + "?index=" + req.Index

	return MostActiveVolume(url, reqH)
}

func (obj GainerLoserObj) ChartData(req models.ChartDataReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + Charts + "?exchange=" + req.Exchange + "&token=" + req.Token + "&candletype=" + req.CandleType + "&starttime=" + req.StartTime + "&endtime=" + req.EndTime + "&data_duration=" + req.DataDuration

	//make payload
	payload := new(bytes.Buffer)

	//call api
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ChartData", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("ChartData call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("ChartData res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlChartData := TradeLabChartDataResponse{}
	json.Unmarshal([]byte(string(body)), &tlChartData)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("ChartData tl status not ok =", tlChartData.Status, "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlChartData.Status
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	var chartDataResponse models.ChartDataResponse
	chartDataResponse.Data = tlChartData.Data
	loggerconfig.Info("chartDataRes tl resp=", helpers.LogStructAsJSON(chartDataResponse), "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	apiRes.Data = chartDataResponse
	return http.StatusOK, apiRes
}

func MostActiveVolume(url string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	//make payload
	payload := new(bytes.Buffer)

	//call api
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "MostActiveVolume", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("MostActiveVolume call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("MostActiveVolume res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlMostActiveVolumeNSE := TradeLabMostActiveVolumeResponse{}
	json.Unmarshal([]byte(string(body)), &tlMostActiveVolumeNSE)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("MostActiveVolume tl status not ok =", tlMostActiveVolumeNSE.Message, "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlMostActiveVolumeNSE.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var mostActiveVolumeData models.MostActiveVolumeData
	mostActiveVolumeDataEntries := make([]models.MostActiveVolume, 0)
	for i := 0; i < len(tlMostActiveVolumeNSE.Data.MostActiveVolume); i++ {
		var mostActiveVolumeDataInstance models.MostActiveVolume
		mostActiveVolumeDataInstance.TurnoverValue = tlMostActiveVolumeNSE.Data.MostActiveVolume[i].TurnoverValue
		mostActiveVolumeDataInstance.TradedQuantity = tlMostActiveVolumeNSE.Data.MostActiveVolume[i].TradedQuantity
		mostActiveVolumeDataInstance.TotalSellQuantity = tlMostActiveVolumeNSE.Data.MostActiveVolume[i].TotalSellQuantity
		mostActiveVolumeDataInstance.TotalBuyQuantity = tlMostActiveVolumeNSE.Data.MostActiveVolume[i].TotalBuyQuantity
		mostActiveVolumeDataInstance.Symbol = tlMostActiveVolumeNSE.Data.MostActiveVolume[i].Symbol
		mostActiveVolumeDataInstance.PreviousPrice = tlMostActiveVolumeNSE.Data.MostActiveVolume[i].PreviousPrice
		mostActiveVolumeDataInstance.NetPrice = tlMostActiveVolumeNSE.Data.MostActiveVolume[i].NetPrice
		mostActiveVolumeDataInstance.Ltp = tlMostActiveVolumeNSE.Data.MostActiveVolume[i].Ltp
		mostActiveVolumeDataInstance.LotSize = tlMostActiveVolumeNSE.Data.MostActiveVolume[i].LotSize
		mostActiveVolumeDataInstance.InstrumentToken = tlMostActiveVolumeNSE.Data.MostActiveVolume[i].InstrumentToken
		mostActiveVolumeDataInstance.Exchange = tlMostActiveVolumeNSE.Data.MostActiveVolume[i].Exchange
		mostActiveVolumeDataInstance.CompanyName = tlMostActiveVolumeNSE.Data.MostActiveVolume[i].CompanyName

		mostActiveVolumeDataEntries = append(mostActiveVolumeDataEntries, mostActiveVolumeDataInstance)
	}
	mostActiveVolumeData.MostActiveVolume = mostActiveVolumeDataEntries
	loggerconfig.Info("MostActiveVolume tl resp=", helpers.LogStructAsJSON(mostActiveVolumeData), "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	apiRes.Data = mostActiveVolumeData
	return http.StatusOK, apiRes
}

func (obj GainerLoserObj) ReturnOnInvestment(roiReq models.ReturnOnInvestmentReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	var apiRes apihelpers.APIRes

	var url string
	if roiReq.Index == "" || strings.ToLower(roiReq.Index) == constants.NSE {
		url = obj.tradeLabURL + RETURNONINVESTMENTURL + "?&days=" + strconv.Itoa(roiReq.Days)
	} else {
		url = obj.tradeLabURL + RETURNONINVESTMENTURL + "?index=" + roiReq.Index + "&days=" + strconv.Itoa(roiReq.Days)
	}

	//make payload
	payload := new(bytes.Buffer)

	//call api
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ReturnOnInvestment", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("ReturnOnInvestment call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("ReturnOnInvestment res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlReturnOnInvestment := TradeLabReturnOnInvestmentRes{}
	json.Unmarshal([]byte(string(body)), &tlReturnOnInvestment)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("MostActiveVolume tl status not ok =", tlReturnOnInvestment.Message, "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlReturnOnInvestment.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var returnOnInvestmentRes models.ReturnOnInvestmentRes
	returnOnInvestmentAllData := make([]models.RoiData, 0)
	for i := 0; i < len(tlReturnOnInvestment.Data.Roi); i++ {
		var returnOnInvestmentData models.RoiData
		returnOnInvestmentData.Volume = tlReturnOnInvestment.Data.Roi[i].Volume
		returnOnInvestmentData.ReturnPercent = tlReturnOnInvestment.Data.Roi[i].ReturnPercent
		returnOnInvestmentData.PercentageChange = tlReturnOnInvestment.Data.Roi[i].PercentageChange
		returnOnInvestmentData.Ltp = tlReturnOnInvestment.Data.Roi[i].Ltp
		returnOnInvestmentData.InstrumentToken = tlReturnOnInvestment.Data.Roi[i].InstrumentToken
		returnOnInvestmentData.Exchange = tlReturnOnInvestment.Data.Roi[i].Exchange
		returnOnInvestmentData.DaysChange = tlReturnOnInvestment.Data.Roi[i].DaysChange
		returnOnInvestmentData.ClosePrice = tlReturnOnInvestment.Data.Roi[i].ClosePrice
		returnOnInvestmentData.Change = tlReturnOnInvestment.Data.Roi[i].Change
		returnOnInvestmentData.TradingSymbol = tlReturnOnInvestment.Data.Roi[i].TradingSymbol
		returnOnInvestmentAllData = append(returnOnInvestmentAllData, returnOnInvestmentData)
	}
	returnOnInvestmentRes.Roi = returnOnInvestmentAllData

	loggerconfig.Info("ReturnOnInvestmentRes tl resp=%v", helpers.LogStructAsJSON(returnOnInvestmentRes), "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
	apiRes.Data = returnOnInvestmentRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj GainerLoserObj) FetchHistoricPerformance(req models.HistoricPerformaceReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	startTime, candleType, dataDuration, valid := getCandleConfigForPeriod(req.Period)
	if !valid {
		apiRes.Message = "Invalid Period"
		apiRes.Status = false
		return http.StatusBadRequest, apiRes
	}
	endTime := strconv.FormatInt(time.Now().Unix(), 10)

	url := obj.tradeLabURL + Charts + "?exchange=" + req.Exchange + "&token=" + req.Token + "&candletype=" + candleType + "&starttime=" + startTime + "&endtime=" + endTime + "&data_duration=" + dataDuration

	payload := new(bytes.Buffer)

	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ChartData", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("ChartData call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("ChartData res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlChartData := TradeLabChartDataResponse{}
	json.Unmarshal([]byte(string(body)), &tlChartData)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("ChartData tl status not ok =", tlChartData.Status, "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlChartData.Status
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	roiData := extractROIData(tlChartData.Data)

	apiRes.Data = roiData
	apiRes.Status = true
	apiRes.Message = "Success"

	return http.StatusOK, apiRes
}

func getCandleConfigForPeriod(period string) (string, string, string, bool) {
	now := time.Now()

	switch period {
	case constants.Period1D:
		return strconv.FormatInt(now.AddDate(0, 0, -1).Unix(), 10), "1", "1", true
	case constants.Period1W:
		return strconv.FormatInt(now.AddDate(0, 0, -7).Unix(), 10), "1", "15", true
	case constants.Period1M:
		return strconv.FormatInt(now.AddDate(0, -1, 0).Unix(), 10), "2", "1", true
	case constants.Period6M:
		return strconv.FormatInt(now.AddDate(0, -6, 0).Unix(), 10), "2", "4", true
	case constants.Period1Y:
		return strconv.FormatInt(now.AddDate(-1, 0, 0).Unix(), 10), "3", "1", true
	case constants.Period3Y:
		return strconv.FormatInt(now.AddDate(-3, 0, 0).Unix(), 10), "3", "1", true
	case constants.Period5Y:
		return strconv.FormatInt(now.AddDate(-5, 0, 0).Unix(), 10), "3", "1", true
	default:
		return "", "", "", false
	}
}

func extractROIData(chartData models.CandleData) models.HistoricPerformaceRes {
	var roiData models.HistoricPerformaceRes

	if len(chartData.Candles) == 0 {
		roiData.ROI = "0.00 %"
		return roiData
	}

	openPrice := getCandleValue(chartData.Candles[0], 1)
	closePrice := getCandleValue(chartData.Candles[len(chartData.Candles)-1], 4)

	roiData.OpenPrice = openPrice
	roiData.ClosePrice = closePrice

	if openPrice > 0 {
		roi := ((closePrice - openPrice) / openPrice) * 100
		roiData.ROI = fmt.Sprintf("%.2f %%", roi)
	} else {
		roiData.ROI = "0.00 %"
	}

	return roiData
}

func getCandleValue(candle []interface{}, index int) float64 {
	if len(candle) <= index {
		return 0
	}
	switch v := candle[index].(type) {
	case float64:
		return v
	case int:
		return float64(v)
	default:
		strVal := fmt.Sprintf("%v", v)
		val, _ := strconv.ParseFloat(strVal, 64)
		return val
	}
}

func (obj GainerLoserObj) FetchAllHistoricPerformance(req models.AllHistoricPerformaceReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes
	performance := models.AllPerformanceRes{}

	periodReq := models.HistoricPerformaceReq{Exchange: req.Exchange, Token: req.Token}

	for _, period := range constants.ChartDataPeriod {
		periodReq.Period = period
		statusCode, periodRes := obj.FetchHistoricPerformance(periodReq, reqH)

		resultValue := constants.NoDataPresent
		if statusCode == http.StatusOK && periodRes.Status {
			if periodData, ok := periodRes.Data.(models.HistoricPerformaceRes); ok {
				resultValue = periodData.ROI
			}
		} else {
			return statusCode, periodRes
		}

		// Assign value to the appropriate struct field
		switch period {
		case constants.Period1D:
			performance.Period1D = resultValue
		case constants.Period1W:
			performance.Period1W = resultValue
		case constants.Period1M:
			performance.Period1M = resultValue
		case constants.Period6M:
			performance.Period6M = resultValue
		case constants.Period1Y:
			performance.Period1Y = resultValue
		case constants.Period3Y:
			performance.Period3Y = resultValue
		case constants.Period5Y:
			performance.Period5Y = resultValue
		}
	}

	apiRes.Data = performance
	apiRes.Status = true
	apiRes.Message = "Success"

	return http.StatusOK, apiRes
}
