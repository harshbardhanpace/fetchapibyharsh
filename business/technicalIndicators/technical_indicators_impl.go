package technicalindicators

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"strconv"
	"sync"

	apihelpers "space/apiHelpers"
	"space/business/tradelab"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
)

type TIObj struct {
	tradeLabURL string
}

func InitTechnicalIndicators() TIObj {
	defer models.HandlePanic()

	tiObj := TIObj{
		tradeLabURL: constants.TLURL,
	}
	return tiObj
}

var FetchChartData = func(url string, reqH models.ReqHeader) ([]models.TLChartCandleData, bool) {
	return FetchChartDataActual(url, reqH)
}

func FetchChartDataActual(url string, reqH models.ReqHeader) ([]models.TLChartCandleData, bool) {
	candleData := make([]models.TLChartCandleData, 0)

	payload := new(bytes.Buffer)
	//call api
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	if err != nil {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " FetchChartData call api error =", err, " requestId:", reqH.RequestId)
		return candleData, false
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := tradelab.TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == tradelab.TLERROR {
		loggerconfig.Error("FetchChartData res error =", tlErrorRes.Message, " requestId:", reqH.RequestId)
		return candleData, false
	}

	tlChartDataBenchmark := tradelab.TradeLabChartDataResponse{}
	json.Unmarshal([]byte(string(body)), &tlChartDataBenchmark)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("FetchChartData tl status not ok =", tlChartDataBenchmark.Status, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)
		return candleData, false
	}
	loggerconfig.Info("FetchChartData TradelabResponse:=", helpers.LogStructAsJSON(tlChartDataBenchmark))

	start := 0
	for i := 0; i < len(tlChartDataBenchmark.Data.Candles); i++ {
		var candleDataEntry models.TLChartCandleData
		write := true
		for j := 0; j < len(tlChartDataBenchmark.Data.Candles[i]); j++ {
			if j == 0 {
				candleDataEntry.Timestamp = tlChartDataBenchmark.Data.Candles[i][j].(string)
				if i > 0 {
					if tlChartDataBenchmark.Data.Candles[i][j].(string) == tlChartDataBenchmark.Data.Candles[start][j].(string) {
						write = false
					} else {
						start++
					}
				}
			}
			switch j {
			case 1:
				candleDataEntry.Open = tlChartDataBenchmark.Data.Candles[i][j].(float64)
			case 2:
				candleDataEntry.High = tlChartDataBenchmark.Data.Candles[i][j].(float64)
			case 3:
				candleDataEntry.Low = tlChartDataBenchmark.Data.Candles[i][j].(float64)
			case 4:
				candleDataEntry.Close = tlChartDataBenchmark.Data.Candles[i][j].(float64)
			case 5:
				candleDataEntry.Volume = tlChartDataBenchmark.Data.Candles[i][j].(float64)
			}

		}
		if write {
			candleData = append(candleData, candleDataEntry)
		}
	}
	return candleData, true
}

func (obj TIObj) TechnicalIndicatorsValues(req models.TechnicalIndicatorsValuesReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes
	starttime := helpers.GetCurrentTimeInIST().Unix() - constants.EIGHTYDAYS
	endtime := helpers.GetCurrentTimeInIST().Unix()

	urlChartData := obj.tradeLabURL + tradelab.Charts + "?exchange=" + req.Exchange + "&token=" + req.Token + "&candletype=" + constants.DayWise + "&starttime=" + strconv.Itoa(int(starttime)) + "&endtime=" + strconv.Itoa(int(endtime)) + "&data_duration=" + constants.DataDuration

	chartData, status := FetchChartData(urlChartData, reqH)
	if !status {
		loggerconfig.Error("TechnicalIndicatorsValues FetchChartData failed, clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		apiRes.Message = constants.ErrorCodeMap[constants.TLChartDataFetchFailed]
		apiRes.Status = false
		return http.StatusInternalServerError, apiRes
	}
	var response models.TechnicalIndicatorsValuesRes

	var tiWG sync.WaitGroup
	tiWG.Add(7)
	tiCH := make(chan apihelpers.ReturnValue, 7)

	go func(tiWG *sync.WaitGroup, ch chan<- apihelpers.ReturnValue) {
		defer tiWG.Done()
		sma20, err := SMA(chartData[len(chartData)-constants.TwentyDays:])
		//sma20, err := SMA(chartData) during unit test use this
		if err != nil {
			loggerconfig.Error("TechnicalIndicatorsValues SMA(20) failed, error:", err, "  clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
			statusCode, res := apihelpers.SendInternalServerError()
			var returnValueSendInternalServerError apihelpers.ReturnValue
			returnValueSendInternalServerError.StatusCode = statusCode
			returnValueSendInternalServerError.ApiRes = res
			ch <- returnValueSendInternalServerError
			return
		}
		response.SMA = sma20
	}(&tiWG, tiCH)

	go func(tiWG *sync.WaitGroup, ch chan<- apihelpers.ReturnValue) {
		defer tiWG.Done()
		ema20, err := EMA(chartData, constants.TwentyDays)
		if err != nil {
			loggerconfig.Error("TechnicalIndicatorsValues EMA(20) failed, error:", err, "  clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
			statusCode, res := apihelpers.SendInternalServerError()
			var returnValueSendInternalServerError apihelpers.ReturnValue
			returnValueSendInternalServerError.StatusCode = statusCode
			returnValueSendInternalServerError.ApiRes = res
			ch <- returnValueSendInternalServerError
			return
		}
		response.EMA = ema20
	}(&tiWG, tiCH)

	go func(tiWG *sync.WaitGroup, ch chan<- apihelpers.ReturnValue) {
		defer tiWG.Done()
		rsi, err := RSI(chartData)
		if err != nil {
			loggerconfig.Error("TechnicalIndicatorsValues RSI(14) failed, error:", err, "  clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
			statusCode, res := apihelpers.SendInternalServerError()
			var returnValueSendInternalServerError apihelpers.ReturnValue
			returnValueSendInternalServerError.StatusCode = statusCode
			returnValueSendInternalServerError.ApiRes = res
			ch <- returnValueSendInternalServerError
			return
		}
		response.RSI = rsi
	}(&tiWG, tiCH)

	go func(tiWG *sync.WaitGroup, ch chan<- apihelpers.ReturnValue) {
		defer tiWG.Done()
		macd, err := MACD(chartData)
		if err != nil {
			loggerconfig.Error("TechnicalIndicatorsValues MACD(12, 26) failed, error:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
			statusCode, res := apihelpers.SendInternalServerError()
			var returnValueSendInternalServerError apihelpers.ReturnValue
			returnValueSendInternalServerError.StatusCode = statusCode
			returnValueSendInternalServerError.ApiRes = res
			ch <- returnValueSendInternalServerError
			return
		}
		response.MACD = macd
	}(&tiWG, tiCH)

	go func(tiWG *sync.WaitGroup, ch chan<- apihelpers.ReturnValue) {
		defer tiWG.Done()
		macdSignal, err := MACDSignal(chartData)
		if err != nil {
			loggerconfig.Error("TechnicalIndicatorsValues MACDSignal failed, error:", err, "  clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
			statusCode, res := apihelpers.SendInternalServerError()
			var returnValueSendInternalServerError apihelpers.ReturnValue
			returnValueSendInternalServerError.StatusCode = statusCode
			returnValueSendInternalServerError.ApiRes = res
			ch <- returnValueSendInternalServerError
			return
		}
		response.MACDSignal = macdSignal
	}(&tiWG, tiCH)

	go func(tiWG *sync.WaitGroup, ch chan<- apihelpers.ReturnValue) {
		defer tiWG.Done()
		awesomeOscillator, err := AwesomeOscillator(chartData)
		if err != nil {
			loggerconfig.Error("TechnicalIndicatorsValues AwesomeOscillator failed, error:", err, "  clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
			statusCode, res := apihelpers.SendInternalServerError()
			var returnValueSendInternalServerError apihelpers.ReturnValue
			returnValueSendInternalServerError.StatusCode = statusCode
			returnValueSendInternalServerError.ApiRes = res
			ch <- returnValueSendInternalServerError
			return
		}
		response.AwesomeOscillator = awesomeOscillator
	}(&tiWG, tiCH)

	go func(tiWG *sync.WaitGroup, ch chan<- apihelpers.ReturnValue) {
		defer tiWG.Done()
		cci, err := CCI(chartData, constants.TwentyDays)
		if err != nil {
			loggerconfig.Error("TechnicalIndicatorsValues CCI failed, error:", err, "  clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
			statusCode, res := apihelpers.SendInternalServerError()
			var returnValueSendInternalServerError apihelpers.ReturnValue
			returnValueSendInternalServerError.StatusCode = statusCode
			returnValueSendInternalServerError.ApiRes = res
			ch <- returnValueSendInternalServerError
			return
		}
		response.CCI = cci
	}(&tiWG, tiCH)

	//var lastValue models.TLChartCandleData during unit test use this
	lastValue := chartData[len(chartData)-1]
	r1, r2, r3, s1, s2, s3 := SupportAndResistance(lastValue)
	response.R1 = r1
	response.R2 = r2
	response.R3 = r3
	response.S1 = s1
	response.S2 = s2
	response.S3 = s3

	apiRes.Data = response
	apiRes.Message = "SUCCESS"
	apiRes.ErrorCode = ""
	apiRes.Status = true
	return http.StatusOK, apiRes
}

var SMA = func(chartdata []models.TLChartCandleData) (float64, error) {
	return SMAActual(chartdata)
}

func SMAActual(chartdata []models.TLChartCandleData) (float64, error) {
	var sma float64

	if len(chartdata) == 0 {
		return sma, errors.New(constants.ChartDataSizeError)
	}

	for i := 0; i < len(chartdata); i++ {
		sma += chartdata[i].Close
	}

	return sma / float64(len(chartdata)), nil
}

var EMA = func(chartdata []models.TLChartCandleData, period int) (float64, error) {
	return EMAActual(chartdata, period)
}

func EMAActual(chartdata []models.TLChartCandleData, period int) (float64, error) {
	var ema float64

	if len(chartdata) < period {
		return ema, errors.New(constants.ChartDataSizeError)
	}

	smoothingFactor := 2 / (float64(period) + 1)

	sma, err := SMA(chartdata[:period])
	if err != nil {
		return ema, err
	}

	ema = sma
	for i := period; i < len(chartdata); i++ {
		ema = (chartdata[i].Close-ema)*smoothingFactor + ema
	}

	return ema, nil
}

var RSI = func(chartdata []models.TLChartCandleData) (float64, error) {
	return RSIActual(chartdata)
}

func RSIActual(chartdata []models.TLChartCandleData) (float64, error) {
	var rsi float64

	if len(chartdata) < constants.FIFTEEN {
		return rsi, errors.New(constants.ChartDataSizeError)
	}

	j := 0
	priceChanges := make([]float64, constants.FOURTEEN)
	for i := len(chartdata) - 1; i > 0 && j < constants.FOURTEEN; i-- {
		priceChanges[j] = chartdata[i].Close - chartdata[i-1].Close
		j++
	}

	var sumGain, sumLoss float64
	for _, change := range priceChanges {
		if change > 0 {
			sumGain += change
		} else {
			sumLoss -= change
		}
	}
	avgGain := sumGain / float64(constants.FOURTEEN)
	avgLoss := sumLoss / float64(constants.FOURTEEN)

	if avgLoss == 0 {
		return 100, nil
	}

	relativeStrength := avgGain / avgLoss
	rsi = 100 - (100 / (1 + relativeStrength))

	return rsi, nil
}

var MACD = func(chartdata []models.TLChartCandleData) (float64, error) {
	return MACDActual(chartdata)
}

func MACDActual(chartdata []models.TLChartCandleData) (float64, error) {
	var macdLine float64
	if len(chartdata) < constants.TwentySixDays {
		return macdLine, errors.New(constants.ChartDataSizeError)
	}

	ema12, err := EMA(chartdata, constants.TwelveDays)
	if err != nil {
		return macdLine, err
	}

	ema26, err := EMA(chartdata, constants.TwentySixDays)
	if err != nil {
		return macdLine, err
	}
	macdLine = ema12 - ema26

	return macdLine, nil
}

var MACDSignal = func(chartdata []models.TLChartCandleData) (float64, error) {
	return MACDSignalActual(chartdata)
}

func MACDSignalActual(chartdata []models.TLChartCandleData) (float64, error) {
	var macdSignal float64
	var err error

	if len(chartdata) < constants.NineDays {
		return macdSignal, errors.New(constants.ChartDataSizeError)
	}
	macds := make([]float64, constants.NineDays)
	for i := 0; i < constants.NineDays; i++ {
		macds[i], err = MACD(chartdata[len(chartdata)-constants.TwentySixDays-i : len(chartdata)-i])
		if err != nil {
			return macdSignal, err
		}
	}
	smoothingFactor := 2 / (float64(constants.NineDays) + 1)
	macdSignal = macds[constants.NineDays-1]
	for i := constants.NineDays - 2; i >= 0; i-- {
		macdSignal = (macds[i] * smoothingFactor) + (macdSignal * (1 - smoothingFactor))
	}

	return macdSignal, nil
}

var AwesomeOscillator = func(chartdata []models.TLChartCandleData) (float64, error) {
	return AwesomeOscillatorActual(chartdata)
}

func AwesomeOscillatorActual(chartdata []models.TLChartCandleData) (float64, error) {
	var awesomeOscillator, sma34, sma5 float64

	if len(chartdata) < constants.ThirtyFourDays {
		return awesomeOscillator, errors.New(constants.ChartDataSizeError)
	}

	medianPriceSlice := make([]float64, constants.ThirtyFourDays)

	j := 0
	for i := len(chartdata) - 1; i >= 0 && j < constants.ThirtyFourDays; i-- {
		medianPriceSlice[j] = (chartdata[i].High + chartdata[i].Low) / 2
		sma34 += medianPriceSlice[j]
		if j < 5 {
			sma5 += medianPriceSlice[j]
		}
		j++
	}
	sma5 = sma5 / constants.FIVE
	sma34 = sma34 / constants.ThirtyFourDays
	awesomeOscillator = sma5 - sma34

	return awesomeOscillator, nil
}

var CCI = func(chartdata []models.TLChartCandleData, period int) (float64, error) {
	return CCIActual(chartdata, period)
}

func CCIActual(chartdata []models.TLChartCandleData, period int) (float64, error) {
	var cci, typicalPrice, typicalPriceSum float64
	if len(chartdata) < period {
		return cci, errors.New(constants.ChartDataSizeError)
	}

	typicalPriceSlice := make([]float64, 0)
	for i := len(chartdata) - period; i < len(chartdata); i++ {
		typicalPrice = (chartdata[i].Close + chartdata[i].High + chartdata[i].Low) / 3
		typicalPriceSlice = append(typicalPriceSlice, typicalPrice)
		typicalPriceSum += typicalPrice
	}

	sma, err := SMA(chartdata[len(chartdata)-constants.TwentyDays:])
	if err != nil {
		return cci, err
	}

	meanDeviation := 0.0
	for i := 0; i < period; i++ {
		meanDeviation += math.Abs(typicalPriceSlice[i] - sma)
	}
	meanDeviation /= float64(period)

	cci = (typicalPrice - typicalPriceSum/float64(period)) / (constants.CCIFactor * meanDeviation)

	return cci, nil
}

var SupportAndResistance = func(datapoint models.TLChartCandleData) (float64, float64, float64, float64, float64, float64) {
	return SupportAndResistanceActual(datapoint)
}

func SupportAndResistanceActual(datapoint models.TLChartCandleData) (float64, float64, float64, float64, float64, float64) {
	pivot := (datapoint.High + datapoint.Low + datapoint.Close) / 3

	r1 := 2*pivot - datapoint.Low
	r2 := pivot + (datapoint.High - datapoint.Low)
	r3 := datapoint.High - 2*(pivot-datapoint.Low)

	s1 := 2*pivot - datapoint.High
	s2 := pivot - (datapoint.High - datapoint.Low)
	s3 := datapoint.Low - (2 * (datapoint.High - datapoint.Low))

	return r1, r2, r3, s1, s2, s3
}
