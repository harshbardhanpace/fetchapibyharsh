package technicalindicatorsV2

import (
	"fmt"
	"math"
	"net/http"
	apihelpers "space/apiHelpers"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
	"strconv"
	"strings"
	"time"
)

type TIV2Obj struct {
}

func InitTechnicalIndicatorsV2() TIV2Obj {
	defer models.HandlePanic()

	tiObj := TIV2Obj{}
	return tiObj
}

func (obj TIV2Obj) GetSMA(req models.GetSMAReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes
	smaRes, err := sma(req, reqH)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, GetSMA err:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var response models.TechnicalIndicatorsResFull
	response.Type = "SMA" + strconv.Itoa(req.SMAType)
	response.Data = smaRes

	apiRes.Data = response
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func sma(req models.GetSMAReq, reqH models.ReqHeader) ([]models.TechnicalIndicatorsRes, error) {
	var res []models.TechnicalIndicatorsRes

	var reqChartData models.ChartDataReq
	reqChartData.Exchange = req.TLChartData.Exchange
	reqChartData.Token = req.TLChartData.Token
	reqChartData.CandleType = req.TLChartData.CandleType
	reqChartData.StartTime = req.TLChartData.StartTime
	reqChartData.EndTime = req.TLChartData.EndTime
	reqChartData.DataDuration = req.TLChartData.DataDuration
	err, chartData := GetChartData(reqChartData, reqH)
	if err != nil {
		return res, err
	}

	rawData, err := CalculateSMA(chartData, req.SMAType)
	if err != nil {
		return res, err
	}

	res, err = formatIndicatorResults(chartData, rawData, calculateMovingAverageSignal)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (obj TIV2Obj) GetEMA(req models.GetEMAReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	emaRes, err := ema(req, reqH)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, GetEMA err:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var response models.TechnicalIndicatorsResFull
	response.Type = "EMA" + strconv.Itoa(req.EMAType)
	response.Data = emaRes

	apiRes.Data = response
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func ema(req models.GetEMAReq, reqH models.ReqHeader) ([]models.TechnicalIndicatorsRes, error) {
	var res []models.TechnicalIndicatorsRes

	var reqChartData models.ChartDataReq
	reqChartData.Exchange = req.TLChartData.Exchange
	reqChartData.Token = req.TLChartData.Token
	reqChartData.CandleType = req.TLChartData.CandleType
	reqChartData.StartTime = req.TLChartData.StartTime
	reqChartData.EndTime = req.TLChartData.EndTime
	reqChartData.DataDuration = req.TLChartData.DataDuration
	err, chartData := GetChartData(reqChartData, reqH)
	if err != nil {
		return res, err
	}

	rawData, err := CalculateEMA(chartData, req.EMAType)
	if err != nil {
		return res, err
	}

	res, err = formatIndicatorResults(chartData, rawData, calculateMovingAverageSignal)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (obj TIV2Obj) GetHullMA(req models.GetHullMAReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	hullmaRes, err := hullma(req, reqH)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, GetHullMA err:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var response models.TechnicalIndicatorsResFull
	response.Type = "HUllMA" + strconv.Itoa(req.HullMAType)
	response.Data = hullmaRes

	apiRes.Data = response
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func hullma(req models.GetHullMAReq, reqH models.ReqHeader) ([]models.TechnicalIndicatorsRes, error) {
	var res []models.TechnicalIndicatorsRes

	var reqChartData models.ChartDataReq
	reqChartData.Exchange = req.TLChartData.Exchange
	reqChartData.Token = req.TLChartData.Token
	reqChartData.CandleType = req.TLChartData.CandleType
	reqChartData.StartTime = req.TLChartData.StartTime
	reqChartData.EndTime = req.TLChartData.EndTime
	reqChartData.DataDuration = req.TLChartData.DataDuration
	err, chartData := GetChartData(reqChartData, reqH)
	if err != nil {
		return res, err
	}

	rawData, err := CalculateHullMA(chartData, req.HullMAType)
	if err != nil {
		return res, err
	}

	res, err = formatIndicatorResults(chartData, rawData, calculateHullMASignal)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (obj TIV2Obj) GetVWMA(req models.GetVWMAReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	vwmaRes, err := vwma(req, reqH)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, GetVWMA err:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var response models.TechnicalIndicatorsResFull
	response.Type = "VWMA" + strconv.Itoa(req.VWMAType)
	response.Data = vwmaRes

	apiRes.Data = response
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func vwma(req models.GetVWMAReq, reqH models.ReqHeader) ([]models.TechnicalIndicatorsRes, error) {
	var res []models.TechnicalIndicatorsRes

	var reqChartData models.ChartDataReq
	reqChartData.Exchange = req.TLChartData.Exchange
	reqChartData.Token = req.TLChartData.Token
	reqChartData.CandleType = req.TLChartData.CandleType
	reqChartData.StartTime = req.TLChartData.StartTime
	reqChartData.EndTime = req.TLChartData.EndTime
	reqChartData.DataDuration = req.TLChartData.DataDuration
	err, chartData := GetChartData(reqChartData, reqH)
	if err != nil {
		return res, err
	}

	rawData, err := CalculateVWMA(chartData, req.VWMAType)
	if err != nil {
		return res, err
	}

	res, err = formatIndicatorResults(chartData, rawData, calculateVWMASignal)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (obj TIV2Obj) GetRSI(req models.GetRSIReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	rsiRes, err := rsi(req, reqH)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, GetRSI err:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var response models.TechnicalIndicatorsResFull
	response.Type = "RSI" + strconv.Itoa(req.RSIType)
	response.Data = rsiRes

	apiRes.Data = response
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func rsi(req models.GetRSIReq, reqH models.ReqHeader) ([]models.TechnicalIndicatorsRes, error) {
	var res []models.TechnicalIndicatorsRes

	var reqChartData models.ChartDataReq
	reqChartData.Exchange = req.TLChartData.Exchange
	reqChartData.Token = req.TLChartData.Token
	reqChartData.CandleType = req.TLChartData.CandleType
	reqChartData.StartTime = req.TLChartData.StartTime
	reqChartData.EndTime = req.TLChartData.EndTime
	reqChartData.DataDuration = req.TLChartData.DataDuration
	err, chartData := GetChartData(reqChartData, reqH)
	if err != nil {
		return res, err
	}

	rawData, err := CalculateRSI(chartData, req.RSIType)
	if err != nil {
		return res, err
	}

	res, err = formatIndicatorResults(chartData, rawData, calculateRsiSignal)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (obj TIV2Obj) GetCCI(req models.GetCCIReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	cciRes, err := cci(req, reqH)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, GetCCI err:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var response models.TechnicalIndicatorsResFull
	response.Type = "CCI" + strconv.Itoa(req.CCIType)
	response.Data = cciRes

	apiRes.Data = response
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func cci(req models.GetCCIReq, reqH models.ReqHeader) ([]models.TechnicalIndicatorsRes, error) {
	var res []models.TechnicalIndicatorsRes

	var reqChartData models.ChartDataReq
	reqChartData.Exchange = req.TLChartData.Exchange
	reqChartData.Token = req.TLChartData.Token
	reqChartData.CandleType = req.TLChartData.CandleType
	reqChartData.StartTime = req.TLChartData.StartTime
	reqChartData.EndTime = req.TLChartData.EndTime
	reqChartData.DataDuration = req.TLChartData.DataDuration
	err, chartData := GetChartData(reqChartData, reqH)
	if err != nil {
		return res, err
	}

	rawData, err := CalculateCCI(chartData, req.CCIType)
	if err != nil {
		return res, err
	}

	res, err = formatIndicatorResults(chartData, rawData, calculateCCISignal)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (obj TIV2Obj) GetMACD(req models.GetMACDReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	macdRes, err := macd(req, reqH)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, GetMACD err:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var response models.TechnicalIndicatorsResFull
	response.Type = "MACD(" + strconv.Itoa(req.FastPeriod) + "," + strconv.Itoa(req.SlowPeriod) + "," + strconv.Itoa(req.SignalPeriod) + ")"
	response.Data = macdRes

	apiRes.Data = response
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func macd(req models.GetMACDReq, reqH models.ReqHeader) ([]models.TechnicalIndicatorsRes, error) {
	var res []models.TechnicalIndicatorsRes

	var reqChartData models.ChartDataReq
	reqChartData.Exchange = req.TLChartData.Exchange
	reqChartData.Token = req.TLChartData.Token
	reqChartData.CandleType = req.TLChartData.CandleType
	reqChartData.StartTime = req.TLChartData.StartTime
	reqChartData.EndTime = req.TLChartData.EndTime
	reqChartData.DataDuration = req.TLChartData.DataDuration
	err, chartData := GetChartData(reqChartData, reqH)
	if err != nil {
		return res, err
	}

	macd, _, _, err := CalculateMACD(chartData, req.FastPeriod, req.SlowPeriod, req.SignalPeriod)
	if err != nil {
		return res, err
	}

	res, err = formatIndicatorResults(chartData, macd, calculateMacdSignal)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (obj TIV2Obj) GetStochastic(req models.GetStochasticReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	stochRes, err := stochastic(req, reqH)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, GetStochastic err:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var response models.TechnicalIndicatorsResFull
	response.Type = "Stochastic(" + strconv.Itoa(req.KPeriod) + "," + strconv.Itoa(req.DPeriod) + "," + strconv.Itoa(req.Smooth) + ")"
	response.Data = stochRes

	apiRes.Data = response
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func stochastic(req models.GetStochasticReq, reqH models.ReqHeader) ([]models.TechnicalIndicatorsRes, error) {
	var res []models.TechnicalIndicatorsRes

	var reqChartData models.ChartDataReq
	reqChartData.Exchange = req.TLChartData.Exchange
	reqChartData.Token = req.TLChartData.Token
	reqChartData.CandleType = req.TLChartData.CandleType
	reqChartData.StartTime = req.TLChartData.StartTime
	reqChartData.EndTime = req.TLChartData.EndTime
	reqChartData.DataDuration = req.TLChartData.DataDuration
	err, chartData := GetChartData(reqChartData, reqH)
	if err != nil {
		return res, err
	}

	resultK, _, err := CalculateStochastic(chartData, req.KPeriod, req.DPeriod, req.Smooth)
	if err != nil {
		return res, err
	}

	res, err = formatIndicatorResults(chartData, resultK, calculateStochasticKSignal)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (obj TIV2Obj) GetIchimokuBaseLine(req models.GetIchimokuBaseLineReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	ibl, err := ichimokuBaseLine(req, reqH)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, GetIchimokuBaseLine err:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var response models.TechnicalIndicatorsResFull
	response.Type = "IchimokuBaseLine"
	response.Data = ibl

	apiRes.Data = response
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func ichimokuBaseLine(req models.GetIchimokuBaseLineReq, reqH models.ReqHeader) ([]models.TechnicalIndicatorsRes, error) {
	var res []models.TechnicalIndicatorsRes

	var reqChartData models.ChartDataReq
	reqChartData.Exchange = req.TLChartData.Exchange
	reqChartData.Token = req.TLChartData.Token
	reqChartData.CandleType = req.TLChartData.CandleType
	reqChartData.StartTime = req.TLChartData.StartTime
	reqChartData.EndTime = req.TLChartData.EndTime
	reqChartData.DataDuration = req.TLChartData.DataDuration
	err, chartData := GetChartData(reqChartData, reqH)
	if err != nil {
		return res, err
	}

	rawData, err := CalculateIchimokuBaseLine(chartData)
	if err != nil {
		return res, err
	}

	res, err = formatIndicatorResults(chartData, rawData, calculateIchimokuBaseLineSignal)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (obj TIV2Obj) GetADX(req models.GetADXReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	adx, err := adx(req, reqH)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, GetADX err:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var response models.TechnicalIndicatorsResFull
	response.Type = "ADX" + strconv.Itoa(req.Period)
	response.Data = adx

	apiRes.Data = response
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func adx(req models.GetADXReq, reqH models.ReqHeader) ([]models.TechnicalIndicatorsRes, error) {
	var res []models.TechnicalIndicatorsRes

	var reqChartData models.ChartDataReq
	reqChartData.Exchange = req.TLChartData.Exchange
	reqChartData.Token = req.TLChartData.Token
	reqChartData.CandleType = req.TLChartData.CandleType
	reqChartData.StartTime = req.TLChartData.StartTime
	reqChartData.EndTime = req.TLChartData.EndTime
	reqChartData.DataDuration = req.TLChartData.DataDuration
	err, chartData := GetChartData(reqChartData, reqH)
	if err != nil {
		return res, err
	}

	rawData, err := CalculateADX(chartData, req.Period)
	if err != nil {
		return res, err
	}

	res, err = formatIndicatorResults(chartData, rawData, calculateADXSignal)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (obj TIV2Obj) GetAwesomeOscillator(req models.GetAwesomeOscillatorReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	ao, err := awesomeOscillator(req, reqH)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, GetAwesomeOscillator err:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var response models.TechnicalIndicatorsResFull
	response.Type = "AwesomeOscillator"
	response.Data = ao

	apiRes.Data = response
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func awesomeOscillator(req models.GetAwesomeOscillatorReq, reqH models.ReqHeader) ([]models.TechnicalIndicatorsRes, error) {
	var res []models.TechnicalIndicatorsRes

	var reqChartData models.ChartDataReq
	reqChartData.Exchange = req.TLChartData.Exchange
	reqChartData.Token = req.TLChartData.Token
	reqChartData.CandleType = req.TLChartData.CandleType
	reqChartData.StartTime = req.TLChartData.StartTime
	reqChartData.EndTime = req.TLChartData.EndTime
	reqChartData.DataDuration = req.TLChartData.DataDuration
	err, chartData := GetChartData(reqChartData, reqH)
	if err != nil {
		return res, err
	}

	ao, err := CalculateAwesomeOscillator(chartData)
	if err != nil {
		return res, err
	}

	res, err = formatIndicatorResultsDiff(chartData, ao, calculateAwesomeOscillatorSignalWithCrossover)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (obj TIV2Obj) GetMomentum(req models.GetMomentumReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	momentumRes, err := momentum(req, reqH)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, GetMomentum err:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var response models.TechnicalIndicatorsResFull
	response.Type = "Momentum" + strconv.Itoa(req.Period)
	response.Data = momentumRes

	apiRes.Data = response
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func momentum(req models.GetMomentumReq, reqH models.ReqHeader) ([]models.TechnicalIndicatorsRes, error) {
	var res []models.TechnicalIndicatorsRes

	var reqChartData models.ChartDataReq
	reqChartData.Exchange = req.TLChartData.Exchange
	reqChartData.Token = req.TLChartData.Token
	reqChartData.CandleType = req.TLChartData.CandleType
	reqChartData.StartTime = req.TLChartData.StartTime
	reqChartData.EndTime = req.TLChartData.EndTime
	reqChartData.DataDuration = req.TLChartData.DataDuration
	err, chartData := GetChartData(reqChartData, reqH)
	if err != nil {
		return res, err
	}

	rawData, err := CalculateMomentum(chartData, req.Period)
	if err != nil {
		return res, err
	}

	res, err = formatIndicatorResultsDiff(chartData, rawData, calculateMomentumSignal)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (obj TIV2Obj) GetStochRSIFast(req models.GetStochRSIFastReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	stochRSI, err := stochasticGetStochRSIFast(req, reqH)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, GetStochRSIFast err:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var response models.TechnicalIndicatorsResFull
	response.Type = "StochRSI(" + strconv.Itoa(req.SmoothK) + "," + strconv.Itoa(req.SmoothD) + "," + strconv.Itoa(req.RsiPeriod) + "," + strconv.Itoa(req.StochPeriod) + ")"
	response.Data = stochRSI

	apiRes.Data = response
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func stochasticGetStochRSIFast(req models.GetStochRSIFastReq, reqH models.ReqHeader) ([]models.TechnicalIndicatorsRes, error) {
	var res []models.TechnicalIndicatorsRes

	var reqChartData models.ChartDataReq
	reqChartData.Exchange = req.TLChartData.Exchange
	reqChartData.Token = req.TLChartData.Token
	reqChartData.CandleType = req.TLChartData.CandleType
	reqChartData.StartTime = req.TLChartData.StartTime
	reqChartData.EndTime = req.TLChartData.EndTime
	reqChartData.DataDuration = req.TLChartData.DataDuration
	err, chartData := GetChartData(reqChartData, reqH)
	if err != nil {
		return res, err
	}

	smoothedK, _, err := CalculateStochRSIFast(chartData, req.SmoothK, req.SmoothD, req.RsiPeriod, req.StochPeriod)
	if err != nil {
		return res, err
	}

	res, err = formatIndicatorResults(chartData, smoothedK, calculateStochRSISignal)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (obj TIV2Obj) GetWilliamsRange(req models.GetWilliamsRangeReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	williamsRange, err := williamsRange(req, reqH)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, GetWilliamsRange err:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var response models.TechnicalIndicatorsResFull
	response.Type = "WilliamsRange" + strconv.Itoa(req.Period)
	response.Data = williamsRange

	apiRes.Data = response
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func williamsRange(req models.GetWilliamsRangeReq, reqH models.ReqHeader) ([]models.TechnicalIndicatorsRes, error) {
	var res []models.TechnicalIndicatorsRes

	var reqChartData models.ChartDataReq
	reqChartData.Exchange = req.TLChartData.Exchange
	reqChartData.Token = req.TLChartData.Token
	reqChartData.CandleType = req.TLChartData.CandleType
	reqChartData.StartTime = req.TLChartData.StartTime
	reqChartData.EndTime = req.TLChartData.EndTime
	reqChartData.DataDuration = req.TLChartData.DataDuration
	err, chartData := GetChartData(reqChartData, reqH)
	if err != nil {
		return res, err
	}

	rawData, err := CalculateWilliamsR(chartData, req.Period)
	if err != nil {
		return res, err
	}

	res, err = formatIndicatorResults(chartData, rawData, calculateWilliamsRSignal)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (obj TIV2Obj) GetUltimateOscillator(req models.GetUltimateOscillatorReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	uo, err := ultimateOscillator(req, reqH)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, GetUltimateOscillator err:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var response models.TechnicalIndicatorsResFull
	response.Type = "UltimateOscillator(" + strconv.Itoa(req.Period1) + "," + strconv.Itoa(req.Period2) + "," + strconv.Itoa(req.Period3) + ")"
	response.Data = uo

	apiRes.Data = response
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func ultimateOscillator(req models.GetUltimateOscillatorReq, reqH models.ReqHeader) ([]models.TechnicalIndicatorsRes, error) {
	var res []models.TechnicalIndicatorsRes

	var reqChartData models.ChartDataReq
	reqChartData.Exchange = req.TLChartData.Exchange
	reqChartData.Token = req.TLChartData.Token
	reqChartData.CandleType = req.TLChartData.CandleType
	reqChartData.StartTime = req.TLChartData.StartTime
	reqChartData.EndTime = req.TLChartData.EndTime
	reqChartData.DataDuration = req.TLChartData.DataDuration
	err, chartData := GetChartData(reqChartData, reqH)
	if err != nil {
		return res, err
	}

	rawData, err := CalculateUltimateOscillator(chartData, req.Period1, req.Period2, req.Period3)
	if err != nil {
		return res, err
	}

	res, err = formatIndicatorResults(chartData, rawData, calculateUltimateOscillatorSignal)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (obj TIV2Obj) GetAllTechnicalIndicators(req models.GetAllTechnicalIndicatorsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	reqPacket := TranslateTechIndicatorReqToChartDataReq(req)
	allTechnicals, pivots, err := getAllTechnicalIndicators(reqPacket, reqH)
	if err != nil {
		loggerconfig.Error("Alert Severity:P2-Mid, GetAllTechnicalIndicators err:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var respPacket models.GetAllTechnicalIndicatorsRes
	respPacket.Entries = make([]models.GetAllTechnicalIndicatorsResInternal, 0)

	for _, indicator := range allTechnicals.Entries {
		//check if the indicator has data, to avoid runtime crash
		if len(indicator.Data) > 0 {
			//taking the last data point from the indicator's data array
			lastDataPoint := indicator.Data[len(indicator.Data)-1]
			entry := models.GetAllTechnicalIndicatorsResInternal{
				Type:          indicator.Type,
				TimestampUnix: lastDataPoint.TimestampUnix,
				Value:         lastDataPoint.Value,
				Signal:        lastDataPoint.Signal,
			}
			respPacket.Entries = append(respPacket.Entries, entry)
		}
	}
	respPacket.Pivots = pivots

	apiRes.Data = respPacket
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func getAllTechnicalIndicators(req models.TLChartDataReq, reqH models.ReqHeader) (models.AllTechnicalIndicatorsRes, []models.PivotsValues, error) {
	var result models.AllTechnicalIndicatorsRes
	var pivotsResult []models.PivotsValues
	var individualEntries []models.TechnicalIndicatorsResFull

	reqChartData := models.ChartDataReq{
		Exchange:     req.Exchange,
		Token:        req.Token,
		CandleType:   req.CandleType,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		DataDuration: req.DataDuration,
	}

	err, chartData := GetChartData(reqChartData, reqH)
	if err != nil {
		return result, pivotsResult, err
	}

	smaPeriods := []int{10, 20, 30, 50, 100, 200}
	for _, period := range smaPeriods {
		entry, err := calculateAndFormatSMA(chartData, period)
		if err != nil {
			return result, pivotsResult, err
		}
		individualEntries = append(individualEntries, entry)
	}

	emaPeriods := []int{10, 20, 30, 50, 100, 200}
	for _, period := range emaPeriods {
		entry, err := calculateAndFormatEMA(chartData, period)
		if err != nil {
			return result, pivotsResult, err
		}
		individualEntries = append(individualEntries, entry)
	}

	rawDataHullMa, err := CalculateHullMA(chartData, 9)
	if err != nil {
		return result, pivotsResult, err
	}
	hullMA, err := formatIndicatorResults(chartData, rawDataHullMa, calculateHullMASignal)
	if err != nil {
		return result, pivotsResult, err
	}
	var hullMA9 models.TechnicalIndicatorsResFull
	hullMA9.Type = "HUllMA9"
	hullMA9.Data = hullMA
	individualEntries = append(individualEntries, hullMA9)

	rawDataVWMA, err := CalculateVWMA(chartData, 20)
	if err != nil {
		return result, pivotsResult, err
	}
	vwma, err := formatIndicatorResults(chartData, rawDataVWMA, calculateVWMASignal)
	if err != nil {
		return result, pivotsResult, err
	}
	var vwma20 models.TechnicalIndicatorsResFull
	vwma20.Type = "VWMA20"
	vwma20.Data = vwma
	individualEntries = append(individualEntries, vwma20)

	rawDataRSI, err := CalculateRSI(chartData, 14)
	if err != nil {
		return result, pivotsResult, err
	}
	rsi, err := formatIndicatorResults(chartData, rawDataRSI, calculateRsiSignal)
	if err != nil {
		return result, pivotsResult, err
	}
	var rsi14 models.TechnicalIndicatorsResFull
	rsi14.Type = "RSI14"
	rsi14.Data = rsi
	individualEntries = append(individualEntries, rsi14)

	rawDataCCI, err := CalculateCCI(chartData, 20)
	if err != nil {
		return result, pivotsResult, err
	}
	cci, err := formatIndicatorResults(chartData, rawDataCCI, calculateCCISignal)
	if err != nil {
		return result, pivotsResult, err
	}
	var cci20 models.TechnicalIndicatorsResFull
	cci20.Type = "CCI20"
	cci20.Data = cci
	individualEntries = append(individualEntries, cci20)

	macd, _, _, err := CalculateMACD(chartData, 12, 26, 8)
	if err != nil {
		return result, pivotsResult, err
	}
	macdVal, err := formatIndicatorResults(chartData, macd, calculateMacdSignal)
	if err != nil {
		return result, pivotsResult, err
	}
	var macd12268 models.TechnicalIndicatorsResFull
	macd12268.Type = "MACD(12,26,8)"
	macd12268.Data = macdVal
	individualEntries = append(individualEntries, macd12268)

	resultK, _, err := CalculateStochastic(chartData, 14, 3, 3)
	if err != nil {
		return result, pivotsResult, err
	}
	stochastic, err := formatIndicatorResults(chartData, resultK, calculateStochasticKSignal)
	if err != nil {
		return result, pivotsResult, err
	}
	var stochastic1433 models.TechnicalIndicatorsResFull
	stochastic1433.Type = "Stochastic(14,3,3)"
	stochastic1433.Data = stochastic
	individualEntries = append(individualEntries, stochastic1433)

	rawDataibl, err := CalculateIchimokuBaseLine(chartData)
	if err != nil {
		return result, pivotsResult, err
	}
	ibl, err := formatIndicatorResults(chartData, rawDataibl, calculateIchimokuBaseLineSignal)
	if err != nil {
		return result, pivotsResult, err
	}
	var ibl9265226 models.TechnicalIndicatorsResFull
	ibl9265226.Type = "IchimokuBaseLine"
	ibl9265226.Data = ibl
	individualEntries = append(individualEntries, ibl9265226)

	rawDataADX, err := CalculateADX(chartData, 14)
	if err != nil {
		return result, pivotsResult, err
	}
	adx, err := formatIndicatorResults(chartData, rawDataADX, calculateADXSignal)
	if err != nil {
		return result, pivotsResult, err
	}
	var adx14 models.TechnicalIndicatorsResFull
	adx14.Type = "ADX14"
	adx14.Data = adx
	individualEntries = append(individualEntries, adx14)

	ao, err := CalculateAwesomeOscillator(chartData)
	if err != nil {
		return result, pivotsResult, err
	}
	awesomeOscillator, err := formatIndicatorResultsDiff(chartData, ao, calculateAwesomeOscillatorSignalWithCrossover)
	if err != nil {
		return result, pivotsResult, err
	}
	var aores models.TechnicalIndicatorsResFull
	aores.Type = "AwesomeOscillator"
	aores.Data = awesomeOscillator
	individualEntries = append(individualEntries, aores)

	rawDataMomentum, err := CalculateMomentum(chartData, 10)
	if err != nil {
		return result, pivotsResult, err
	}
	resMomentum, err := formatIndicatorResultsDiff(chartData, rawDataMomentum, calculateMomentumSignal)
	if err != nil {
		return result, pivotsResult, err
	}
	var momentum10 models.TechnicalIndicatorsResFull
	momentum10.Type = "Momentum10"
	momentum10.Data = resMomentum
	individualEntries = append(individualEntries, momentum10)

	smoothedK, _, err := CalculateStochRSIFast(chartData, 3, 3, 14, 14)
	if err != nil {
		return result, pivotsResult, err
	}
	rsiFast, err := formatIndicatorResults(chartData, smoothedK, calculateStochRSISignal)
	if err != nil {
		return result, pivotsResult, err
	}
	var rsif models.TechnicalIndicatorsResFull
	rsif.Type = "StochRSI(3,3,14,14)"
	rsif.Data = rsiFast
	individualEntries = append(individualEntries, rsif)

	rawDataWR, err := CalculateWilliamsR(chartData, 14)
	if err != nil {
		return result, pivotsResult, err
	}
	wr, err := formatIndicatorResults(chartData, rawDataWR, calculateWilliamsRSignal)
	if err != nil {
		return result, pivotsResult, err
	}
	var wr14 models.TechnicalIndicatorsResFull
	wr14.Type = "WilliamsRange(14)"
	wr14.Data = wr
	individualEntries = append(individualEntries, wr14)

	rawDataUO, err := CalculateUltimateOscillator(chartData, 7, 14, 28)
	if err != nil {
		return result, pivotsResult, err
	}
	uo, err := formatIndicatorResults(chartData, rawDataUO, calculateUltimateOscillatorSignal)
	if err != nil {
		return result, pivotsResult, err
	}
	var uo713428 models.TechnicalIndicatorsResFull
	uo713428.Type = "UltimateOscillator(7,14,28)"
	uo713428.Data = uo
	individualEntries = append(individualEntries, uo713428)

	pivots, err := getAllPivots(chartData)
	if err != nil {
		return result, pivotsResult, err
	}
	pivotsResult = pivots
	result.Entries = individualEntries

	return result, pivotsResult, nil
}

func getAllPivots(chartData models.ChartDataResponse) ([]models.PivotsValues, error) {
	var response []models.PivotsValues

	if len(chartData.Data.Candles) < 1 {
		return response, fmt.Errorf("not enough candle data available")
	}

	lastCandle := chartData.Data.Candles[len(chartData.Data.Candles)-1]
	timestampInterface := lastCandle[0]
	timestamp, ok := timestampInterface.(string)
	if !ok {
		return response, fmt.Errorf("timestamp is not a string")
	}

	allPivotsResponse, err := CalculateAllPivotPoints(chartData)
	if err != nil {
		return response, err
	}

	for pivotType, points := range allPivotsResponse {
		pivotEntry := models.PivotsValues{
			Type:          pivotType,
			TimestampUnix: timestamp,
			Points:        points,
		}
		response = append(response, pivotEntry)
	}

	return response, nil
}

func calculateAndFormatSMA(chartData models.ChartDataResponse, period int) (models.TechnicalIndicatorsResFull, error) {
	rawData, err := CalculateSMA(chartData, period)
	if err != nil {
		return models.TechnicalIndicatorsResFull{}, err
	}

	res, err := formatIndicatorResults(chartData, rawData, calculateMovingAverageSignal)
	if err != nil {
		return models.TechnicalIndicatorsResFull{}, err
	}

	return models.TechnicalIndicatorsResFull{
		Type: fmt.Sprintf("SMA%d", period),
		Data: res,
	}, nil
}

func calculateAndFormatEMA(chartData models.ChartDataResponse, period int) (models.TechnicalIndicatorsResFull, error) {
	rawData, err := CalculateEMA(chartData, period)
	if err != nil {
		return models.TechnicalIndicatorsResFull{}, err
	}

	res, err := formatIndicatorResults(chartData, rawData, calculateMovingAverageSignal)
	if err != nil {
		return models.TechnicalIndicatorsResFull{}, err
	}

	return models.TechnicalIndicatorsResFull{
		Type: fmt.Sprintf("EMA%d", period),
		Data: res,
	}, nil
}

var (
	//Market hours in IST (09:15 to 15:15, which is 6 hours)
	MarketOpenHour  = 9
	MarketOpenMin   = 15
	MarketCloseHour = 15
	MarketCloseMin  = 15

	// Market open days (Monday to Friday)
	MarketOpenDays = 5

	// Hours market is open per day
	MarketHoursPerDay = 6

	// Minimum data points required
	MinDataPoints = 300

	// Add 50% buffer for safety
	DataPointBuffer = 1.5

	LocationKolkata = time.FixedZone("IST", 5*60*60+30*60)
)

var timeNow = time.Now

// Use this function instead of direct calls to time.Now()
func getNow() time.Time {
	return timeNow()
}

// accounts for holidays when calculating required duration
func TranslateTechIndicatorReqToChartDataReq(req models.GetAllTechnicalIndicatorsReq) models.TLChartDataReq {
	result := models.TLChartDataReq{
		Exchange: req.Exchange,
		Token:    req.Token,
	}

	switch req.TimeUnit {
	case "MINUTE":
		result.CandleType = "1" // 1 for minute
	case "HOUR":
		result.CandleType = "2" // 2 for hour
	case "DAY", "WEEK", "MONTH":
		result.CandleType = "3" // 3 for day
	}

	//Set DataDuration based on TimeUnit and TimeInterval
	switch req.TimeUnit {
	case "MINUTE", "HOUR", "DAY":
		result.DataDuration = strconv.Itoa(req.TimeInterval)
	case "WEEK":
		result.DataDuration = strconv.Itoa(req.TimeInterval * 7) // Week = 7 days
	case "MONTH":
		result.DataDuration = strconv.Itoa(req.TimeInterval * 30) // Month ~= 30 days
	}

	endTime := getNow().In(LocationKolkata)

	var requiredDuration time.Duration

	// Get estimated number of trading days needed
	tradingDaysNeeded := estimateTradingDays(req.TimeUnit, req.TimeInterval)

	// Convert trading days to calendar days considering holidays and weekends
	calendarDays := convertTradingDaysToCalendarDays(tradingDaysNeeded, endTime)

	// Set the duration
	requiredDuration = time.Duration(calendarDays) * 24 * time.Hour

	// Calculate start time by subtracting required duration from end time
	startTime := endTime.Add(-requiredDuration)

	// Set timestamps
	result.EndTime = strconv.FormatInt(endTime.Unix(), 10)
	result.StartTime = strconv.FormatInt(startTime.Unix(), 10)

	return result
}

// estimateTradingDays calculates how many trading days are needed to get the minimum required data points
func estimateTradingDays(timeUnit string, timeInterval int) int {
	switch timeUnit {
	case "MINUTE":
		// For minutes, calculate data points per day based on market hours
		dataPointsPerHour := 60 / timeInterval
		dataPointsPerDay := dataPointsPerHour * MarketHoursPerDay
		if dataPointsPerDay == 0 {
			dataPointsPerDay = 1 // Avoid division by zero
		}
		return int(math.Ceil(float64(MinDataPoints) / float64(dataPointsPerDay) * DataPointBuffer))

	case "HOUR":
		// For hours, we get MarketHoursPerDay/TimeInterval data points per day
		dataPointsPerDay := MarketHoursPerDay / timeInterval
		if dataPointsPerDay == 0 {
			dataPointsPerDay = 1 // Handle case where interval > market hours
		}
		return int(math.Ceil(float64(MinDataPoints) / float64(dataPointsPerDay) * DataPointBuffer))

	case "DAY":
		// For days, each trading day gives one data point
		return MinDataPoints * timeInterval

	case "WEEK":
		// For weeks, convert to trading days
		return MinDataPoints * timeInterval * 5 // Approximating 5 trading days per week

	case "MONTH":
		// For months, convert to trading days
		return MinDataPoints * timeInterval * 22 // Approximating 22 trading days per month
	}

	return MinDataPoints // Default fallback
}

func convertTradingDaysToCalendarDays(tradingDays int, endDate time.Time) int {
	if tradingDays <= 0 {
		return 0
	}

	calendarDays := int(float64(tradingDays) * 7.0 / 5.0)

	//Add buffer for holidays (approximately 10 holidays per year)
	//10 holidays / 365 days ≈ 0.027 holidays per day
	holidayBuffer := int(float64(calendarDays) * 0.027)

	//further buffer of 10% to be safe
	return calendarDays + holidayBuffer + int(float64(calendarDays)*0.1)
}

// For unit testing purposes
func estimateHolidaysInDateRange(startDate, endDate time.Time) int {
	count := 0

	startDateStr := startDate.Format("02-Jan-2006")
	endDateStr := endDate.Format("02-Jan-2006")

	for _, holidayDate := range helpers.HolidayCalendar.Date {
		if !holidayDate.IsHoliday {
			continue
		}

		parsedDate, err := time.Parse("02-Jan-2006", holidayDate.Date)
		if err != nil {
			continue
		}

		if (parsedDate.After(startDate) || parsedDate.Format("02-Jan-2006") == startDateStr) &&
			(parsedDate.Before(endDate) || parsedDate.Format("02-Jan-2006") == endDateStr) {
			count++
		}
	}

	return count
}

func IsMarketOpen() bool {
	now := getNow().In(LocationKolkata)
	currentDate := now.Format("02-Jan-2006")
	if isHoliday(currentDate) {
		return false
	}
	hour, min := now.Hour(), now.Minute()

	//Before market opens
	if hour < MarketOpenHour || (hour == MarketOpenHour && min < MarketOpenMin) {
		return false
	}

	//After market closes
	if hour > MarketCloseHour || (hour == MarketCloseHour && min > MarketCloseMin) {
		return false
	}

	return true
}

func isHoliday(date string) bool {
	normalizedDate := strings.ToLower(date)
	for _, holidayDate := range helpers.HolidayCalendar.Date {
		if strings.ToLower(holidayDate.Date) == normalizedDate {
			return holidayDate.IsHoliday
		}
	}

	parsedDate, err := time.Parse("02-Jan-2006", date)
	if err != nil {
		return false
	}

	weekday := parsedDate.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}
