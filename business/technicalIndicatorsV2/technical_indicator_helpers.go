package technicalindicatorsV2

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
	"time"

	"github.com/cinar/indicator"
)

type TradeLabErrorRes struct {
	Data struct {
	} `json:"data"`
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
	Status    string `json:"status"`
}

type TradeLabChartDataResponse struct {
	Status string `json:"status"`
	Data   struct {
		Candles [][]interface{} `json:"candles"`
	} `json:"data"`
}

func GetChartData(req models.ChartDataReq, reqH models.ReqHeader) (error, models.ChartDataResponse) {
	var chartDataResponse models.ChartDataResponse
	url := constants.TLURL + constants.Charts + "?exchange=" + req.Exchange + "&token=" + req.Token + "&candletype=" + req.CandleType + "&starttime=" + req.StartTime + "&endtime=" + req.EndTime + "&data_duration=" + req.DataDuration

	//make payload
	payload := new(bytes.Buffer)

	//call api
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ChartData", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("ChartData call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return err, chartDataResponse
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == "error" {
		loggerconfig.Error("ChartData res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return fmt.Errorf("Tradelab error, error parsing %s", tlErrorRes.Message), chartDataResponse
	}

	tlChartData := TradeLabChartDataResponse{}
	json.Unmarshal([]byte(string(body)), &tlChartData)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("ChartData tl status not ok =", tlChartData.Status, "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return fmt.Errorf("Tradelab error status not ok %d", res.StatusCode), chartDataResponse
	}

	chartDataResponse.Data = tlChartData.Data
	//loggerconfig.Info("chartDataRes tl resp=", helpers.LogStructAsJSON(chartDataResponse), "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
	return nil, chartDataResponse
}

// parseClosePrices extracts close prices from the candle data.
func parseClosePrices(data models.ChartDataResponse) ([]float64, error) {
	var closePrices []float64

	for _, candle := range data.Data.Candles {
		if len(candle) < 5 {
			return nil, fmt.Errorf("invalid candle data: %v", candle)
		}
		close, ok := candle[4].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid close price type: %v", candle[4])
		}
		closePrices = append(closePrices, close)
	}
	return closePrices, nil
}

func CalculateSMA(data models.ChartDataResponse, period int) ([]float64, error) {
	closePrices, err := parseClosePrices(data)
	if err != nil {
		return nil, err
	}
	if len(closePrices) < period {
		return nil, errors.New("not enough data points")
	}

	result := indicator.Sma(period, closePrices)

	return result, nil
}

func CalculateEMA(data models.ChartDataResponse, period int) ([]float64, error) {
	closePrices, err := parseClosePrices(data)
	if err != nil {
		return nil, err
	}
	if len(closePrices) < period {
		return nil, errors.New("not enough data points")
	}
	result := indicator.Ema(period, closePrices)

	return result, nil
}

// parseHighPrices extracts high prices from the candle data
func parseHighPrices(data models.ChartDataResponse) ([]float64, error) {
	var highPrices []float64
	for _, candle := range data.Data.Candles {
		if len(candle) < 3 {
			return nil, fmt.Errorf("invalid candle data: %v", candle)
		}
		high, ok := candle[2].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid high price type: %v", candle[2])
		}
		highPrices = append(highPrices, high)
	}
	return highPrices, nil
}

// parseLowPrices extracts low prices from the candle data
func parseLowPrices(data models.ChartDataResponse) ([]float64, error) {
	var lowPrices []float64
	for _, candle := range data.Data.Candles {
		if len(candle) < 4 {
			return nil, fmt.Errorf("invalid candle data: %v", candle)
		}
		low, ok := candle[3].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid low price type: %v", candle[3])
		}
		lowPrices = append(lowPrices, low)
	}
	return lowPrices, nil
}

// parseVolumes extracts volume data from the candle data
func parseVolumes(data models.ChartDataResponse) ([]int, error) {
	var volumes []int
	for _, candle := range data.Data.Candles {
		if len(candle) < 6 {
			return nil, fmt.Errorf("invalid candle data: %v", candle)
		}
		volumeFloat64, ok := candle[5].(float64) // Handle float64 volumes
		if !ok {
			return nil, fmt.Errorf("invalid volume type: %v", candle[5])
		}
		volumes = append(volumes, int(volumeFloat64)) // Convert to int
	}
	return volumes, nil
}

func calculateWMA(period int, values []float64) []float64 {
	var wma []float64
	denominator := (period * (period + 1)) / 2

	for i := 0; i < len(values); i++ {
		if i < period-1 {
			wma = append(wma, 0)
			continue
		}

		sum := 0.0
		weight := period
		for j := 0; j < period; j++ {
			sum += values[i-j] * float64(weight)
			weight--
		}
		wma = append(wma, sum/float64(denominator))
	}

	return wma
}

// CalculateHullMA calculates Hull Moving Average
func CalculateHullMA(data models.ChartDataResponse, period int) ([]float64, error) {
	closePrices, err := parseClosePrices(data)
	if err != nil {
		return nil, err
	}

	// Calculate WMA with period and period/2
	wma1 := calculateWMA(period, closePrices)
	wma2 := calculateWMA(period/2, closePrices)

	// Calculate 2*WMA(n/2) - WMA(n)
	var diffSeries []float64
	for i := 0; i < len(wma1); i++ {
		if i < len(wma2) {
			diff := 2*wma2[i] - wma1[i]
			diffSeries = append(diffSeries, diff)
		}
	}

	// Calculate final Hull MA
	sqrtPeriod := int(math.Sqrt(float64(period)))
	hma := calculateWMA(sqrtPeriod, diffSeries)

	return hma, nil
}

// CalculateVWMA calculates Volume Weighted Moving Average
func CalculateVWMA(data models.ChartDataResponse, period int) ([]float64, error) {
	closePrices, err := parseClosePrices(data)
	if err != nil {
		return nil, err
	}

	volumes, err := parseVolumes(data)
	if err != nil {
		return nil, err
	}

	var vwma []float64
	for i := 0; i < len(closePrices); i++ {
		if i < period-1 {
			vwma = append(vwma, 0)
			continue
		}

		sumCV := 0.0
		sumV := 0.0
		for j := 0; j < period; j++ {
			idx := i - j
			sumCV += closePrices[idx] * float64(volumes[idx])
			sumV += float64(volumes[idx])
		}
		if sumV > 0 {
			vwma = append(vwma, sumCV/sumV)
		} else {
			vwma = append(vwma, closePrices[i])
		}
	}

	return vwma, nil
}

// CalculateRSI calculates Relative Strength Index
func CalculateRSI(data models.ChartDataResponse, period int) ([]float64, error) {
	closePrices, err := parseClosePrices(data)
	if err != nil {
		return nil, err
	}

	// Create price changes for RSI calculation
	var changes []float64
	for i := 1; i < len(closePrices); i++ {
		change := closePrices[i] - closePrices[i-1]
		changes = append(changes, change)
	}

	// indicator.Rsi only takes price changes as input
	_, rsi := indicator.Rsi(changes)

	// Pad the beginning with zeros to match input length
	result := make([]float64, len(closePrices))
	copy(result[period:], rsi)

	return result, nil
}

// CalculateCCI calculates Commodity Channel Index
func CalculateCCI(data models.ChartDataResponse, period int) ([]float64, error) {
	highPrices, err := parseHighPrices(data)
	if err != nil {
		return nil, err
	}

	lowPrices, err := parseLowPrices(data)
	if err != nil {
		return nil, err
	}

	closePrices, err := parseClosePrices(data)
	if err != nil {
		return nil, err
	}

	return indicator.CommunityChannelIndex(period, highPrices, lowPrices, closePrices), nil
}

// CalculateMACD calculates Moving Average Convergence Divergence
func CalculateMACD(data models.ChartDataResponse, fastPeriod, slowPeriod, signalPeriod int) ([]float64, []float64, []float64, error) {
	closePrices, err := parseClosePrices(data)
	if err != nil {
		return nil, nil, nil, err
	}

	// Get MACD and signal line from indicator
	macd, signal := indicator.Macd(closePrices)

	// Calculate histogram (MACD - Signal)
	histogram := make([]float64, len(macd))
	for i := 0; i < len(macd); i++ {
		if i < len(signal) {
			histogram[i] = macd[i] - signal[i]
		}
	}

	return macd, signal, histogram, nil
}

// CalculateStochastic calculates Stochastic Oscillator
func CalculateStochastic(data models.ChartDataResponse, kPeriod, dPeriod, smooth int) ([]float64, []float64, error) {
	highPrices, err := parseHighPrices(data)
	if err != nil {
		return nil, nil, err
	}

	lowPrices, err := parseLowPrices(data)
	if err != nil {
		return nil, nil, err
	}

	closePrices, err := parseClosePrices(data)
	if err != nil {
		return nil, nil, err
	}

	if len(closePrices) < kPeriod {
		return nil, nil, fmt.Errorf("not enough data points")
	}

	k, d := indicator.StochasticOscillator(highPrices, lowPrices, closePrices)

	return k, d, nil
}

// CalculateIchimokuBaseLine calculates Ichimoku Base Line (Kijun-sen)
func CalculateIchimokuBaseLine(data models.ChartDataResponse) ([]float64, error) {
	highPrices, err := parseHighPrices(data)
	if err != nil {
		return nil, err
	}

	lowPrices, err := parseLowPrices(data)
	if err != nil {
		return nil, err
	}

	period := 26
	var baseline []float64

	for i := 0; i < len(highPrices); i++ {
		if i < period-1 {
			baseline = append(baseline, 0) // Not enough data for calculation
			continue
		}

		highest := math.SmallestNonzeroFloat64 // Initialize with very small value
		lowest := math.MaxFloat64              // Initialize with very large value

		for j := 0; j < period; j++ {
			idx := i - j
			highest = math.Max(highest, highPrices[idx])
			lowest = math.Min(lowest, lowPrices[idx])
		}

		baseline = append(baseline, (highest+lowest)/2)
	}

	return baseline, nil
}

func CalculateADX(data models.ChartDataResponse, period int) ([]float64, error) {
	highPrices, err := parseHighPrices(data)
	if err != nil {
		return nil, err
	}

	lowPrices, err := parseLowPrices(data)
	if err != nil {
		return nil, err
	}

	closePrices, err := parseClosePrices(data)
	if err != nil {
		return nil, err
	}

	if len(highPrices) <= period {
		return nil, fmt.Errorf("not enough data points to calculate ADX")
	}

	var tr, plusDM, minusDM []float64

	for i := 1; i < len(highPrices); i++ {
		highDiff := highPrices[i] - highPrices[i-1]
		lowDiff := lowPrices[i-1] - lowPrices[i]

		// True Range calculation
		trueRange := math.Max(highPrices[i]-lowPrices[i],
			math.Max(math.Abs(highPrices[i]-closePrices[i-1]), math.Abs(lowPrices[i]-closePrices[i-1])))
		tr = append(tr, trueRange)

		// Directional Movement
		if highDiff > lowDiff && highDiff > 0 {
			plusDM = append(plusDM, highDiff)
		} else {
			plusDM = append(plusDM, 0)
		}

		if lowDiff > highDiff && lowDiff > 0 {
			minusDM = append(minusDM, lowDiff)
		} else {
			minusDM = append(minusDM, 0)
		}
	}

	// Smoothing using a simple moving average
	smoothedTR := smooth(tr, period)
	smoothedPlusDM := smooth(plusDM, period)
	smoothedMinusDM := smooth(minusDM, period)

	var plusDI, minusDI, dx, adx []float64

	for i := 0; i < len(smoothedTR); i++ {
		if smoothedTR[i] == 0 {
			plusDI = append(plusDI, 0)
			minusDI = append(minusDI, 0)
			dx = append(dx, 0)
			continue
		}

		plusDIVal := (smoothedPlusDM[i] / smoothedTR[i]) * 100
		minusDIVal := (smoothedMinusDM[i] / smoothedTR[i]) * 100

		plusDI = append(plusDI, plusDIVal)
		minusDI = append(minusDI, minusDIVal)

		dxVal := (math.Abs(plusDIVal-minusDIVal) / (plusDIVal + minusDIVal)) * 100
		dx = append(dx, dxVal)
	}

	// Smooth DX to get ADX
	adx = smooth(dx, period)

	return adx, nil
}

// Helper function to smooth values (Simple Moving Average)
func smooth(values []float64, period int) []float64 {
	var smoothed []float64
	for i := 0; i <= len(values)-period; i++ {
		sum := 0.0
		for j := i; j < i+period; j++ {
			sum += values[j]
		}
		smoothed = append(smoothed, sum/float64(period))
	}
	return smoothed
}

// CalculateAwesomeOscillator calculates Awesome Oscillator
func CalculateAwesomeOscillator(data models.ChartDataResponse) ([]float64, error) {
	highPrices, err := parseHighPrices(data)
	if err != nil {
		return nil, err
	}

	lowPrices, err := parseLowPrices(data)
	if err != nil {
		return nil, err
	}

	var medianPrices []float64
	for i := 0; i < len(highPrices); i++ {
		medianPrices = append(medianPrices, (highPrices[i]+lowPrices[i])/2)
	}

	sma5 := indicator.Sma(5, medianPrices)
	sma34 := indicator.Sma(34, medianPrices)

	var ao []float64
	for i := 0; i < len(sma5); i++ {
		if i < len(sma34) {
			ao = append(ao, sma5[i]-sma34[i])
		}
	}

	return ao, nil
}

// CalculateMomentum calculates Momentum
func CalculateMomentum(data models.ChartDataResponse, period int) ([]float64, error) {
	closePrices, err := parseClosePrices(data)
	if err != nil {
		return nil, err
	}

	var momentum []float64
	for i := 0; i < len(closePrices); i++ {
		if i < period {
			momentum = append(momentum, 0)
			continue
		}
		momentum = append(momentum, closePrices[i]-closePrices[i-period])
	}

	return momentum, nil
}

func CalculateStochRSIFast(data models.ChartDataResponse, smoothK, smoothD, rsiPeriod, stochPeriod int) ([]float64, []float64, error) {
	closePrices, err := parseClosePrices(data)
	if err != nil {
		return nil, nil, err
	}

	if len(closePrices) < rsiPeriod+stochPeriod {
		return nil, nil, fmt.Errorf("not enough data points")
	}

	// Compute RSI first
	_, rsiValues := indicator.Rsi(closePrices)

	// Pre-allocate result arrays
	rawK := []float64{}

	// We need at least stochPeriod RSI values
	for i := stochPeriod; i <= len(rsiValues); i++ {
		// Get the window of RSI values for this calculation
		window := rsiValues[i-stochPeriod : i]

		// Find highest and lowest in the window
		highest := window[0]
		lowest := window[0]
		for _, val := range window {
			if val > highest {
				highest = val
			}
			if val < lowest {
				lowest = val
			}
		}

		// Calculate %K value
		var kVal float64
		if highest-lowest > 0 {
			kVal = 100 * (rsiValues[i-1] - lowest) / (highest - lowest)
		} else {
			kVal = 50 // Neutral value when no movement
		}
		rawK = append(rawK, kVal)
	}

	// Apply smoothing to %K
	smoothedK := indicator.Sma(smoothK, rawK)

	// Calculate %D (SMA of %K)
	d := indicator.Sma(smoothD, smoothedK)

	return smoothedK, d, nil
}

// CalculateWilliamsR calculates Williams %R
func CalculateWilliamsR(data models.ChartDataResponse, period int) ([]float64, error) {
	highPrices, err := parseHighPrices(data)
	if err != nil {
		return nil, err
	}

	lowPrices, err := parseLowPrices(data)
	if err != nil {
		return nil, err
	}

	closePrices, err := parseClosePrices(data)
	if err != nil {
		return nil, err
	}

	// Calculate Highest High and Lowest Low for the given period
	highestHighs := make([]float64, len(highPrices))
	lowestLows := make([]float64, len(lowPrices))

	for i := period - 1; i < len(highPrices); i++ {
		highSlice := highPrices[i-period+1 : i+1]
		lowSlice := lowPrices[i-period+1 : i+1]

		highestHigh := highSlice[0]
		lowestLow := lowSlice[0]

		for _, h := range highSlice {
			highestHigh = math.Max(highestHigh, h)
		}
		for _, l := range lowSlice {
			lowestLow = math.Min(lowestLow, l)
		}
		highestHighs[i] = highestHigh
		lowestLows[i] = lowestLow
	}

	wr := make([]float64, len(closePrices))
	for i := period - 1; i < len(closePrices); i++ {
		wr[i] = (highestHighs[i] - closePrices[i]) / (highestHighs[i] - lowestLows[i]) * -100
	}

	return wr[period-1:], nil // Return the WR values starting from the 'period' index
}

func CalculateUltimateOscillator(data models.ChartDataResponse, period1, period2, period3 int) ([]float64, error) {
	highPrices, err := parseHighPrices(data)
	if err != nil {
		return nil, err
	}

	lowPrices, err := parseLowPrices(data)
	if err != nil {
		return nil, err
	}

	closePrices, err := parseClosePrices(data)
	if err != nil {
		return nil, err
	}

	// Calculate buying pressure (BP) and true range (TR)
	var bp, tr []float64
	for i := 0; i < len(closePrices); i++ {
		if i == 0 {
			bp = append(bp, 0)
			tr = append(tr, highPrices[i]-lowPrices[i])
			continue
		}

		// Buying pressure = Close - Minimum(Low or Prior Close)
		minVal := math.Min(lowPrices[i], closePrices[i-1])
		bp = append(bp, closePrices[i]-minVal)

		// True Range = Maximum(High-Low, |High-Prior Close|, |Low-Prior Close|)
		trueRange := math.Max(highPrices[i]-lowPrices[i],
			math.Max(math.Abs(highPrices[i]-closePrices[i-1]),
				math.Abs(lowPrices[i]-closePrices[i-1])))
		tr = append(tr, trueRange)
	}

	// Calculate averages for different periods
	var ultimate []float64
	for i := 0; i < len(closePrices); i++ {
		if i < period3-1 {
			ultimate = append(ultimate, 0)
			continue
		}

		var avg1, avg2, avg3 float64
		var sum1BP, sum1TR, sum2BP, sum2TR, sum3BP, sum3TR float64

		// Calculate sums for each period
		for j := 0; j < period1; j++ {
			idx := i - j
			sum1BP += bp[idx]
			sum1TR += tr[idx]
		}
		for j := 0; j < period2; j++ {
			idx := i - j
			sum2BP += bp[idx]
			sum2TR += tr[idx]
		}
		for j := 0; j < period3; j++ {
			idx := i - j
			sum3BP += bp[idx]
			sum3TR += tr[idx]
		}

		// Calculate averages
		if sum1TR != 0 {
			avg1 = sum1BP / sum1TR
		}
		if sum2TR != 0 {
			avg2 = sum2BP / sum2TR
		}
		if sum3TR != 0 {
			avg3 = sum3BP / sum3TR
		}

		// Calculate Ultimate Oscillator
		// UO = 100 × [(4 × Average1) + (2 × Average2) + Average3] / (4 + 2 + 1)
		uo := 100 * ((4 * avg1) + (2 * avg2) + avg3) / 7
		ultimate = append(ultimate, uo)
	}

	return ultimate, nil
}

func formatIndicatorResults(data models.ChartDataResponse, indicatorValues []float64, signalFunction func(currentPrice, indicatorValue float64) string) ([]models.TechnicalIndicatorsRes, error) {
	result := make([]models.TechnicalIndicatorsRes, len(indicatorValues))
	closePrices, err := parseClosePrices(data)
	if err != nil {
		return nil, err
	}

	for i, value := range indicatorValues {
		var timestamp string
		if len(data.Data.Candles) > i {
			ts, ok := data.Data.Candles[i][0].(string)
			if ok {
				timestamp = ts
			}
		}

		if math.IsNaN(value) {
			value = 0.0
		}

		var signal string
		if i >= len(closePrices) || math.IsNaN(value) {
			signal = "neutral"
		} else {
			currentPrice := closePrices[i]
			signal = signalFunction(currentPrice, value)
		}

		result[i] = models.TechnicalIndicatorsRes{
			TimestampUnix: timestamp,
			Value:         value,
			Signal:        signal,
		}
	}

	return result, nil
}

func formatIndicatorResultsDiff(data models.ChartDataResponse, indicatorValues []float64, signalFunction func(previousValue, currentValue float64) string) ([]models.TechnicalIndicatorsRes, error) {
	result := make([]models.TechnicalIndicatorsRes, len(indicatorValues))

	for i, value := range indicatorValues {
		var timestamp string
		if len(data.Data.Candles) > i {
			ts, ok := data.Data.Candles[i][0].(string)
			if ok {
				timestamp = ts
			}
		}

		if math.IsNaN(value) {
			value = 0.0
		}

		// Ensure there's a previous value to compare
		var signal string
		if i == 0 || math.IsNaN(value) {
			signal = "neutral"
		} else {
			previousValue := indicatorValues[i-1] // Fix: Use previous value
			signal = signalFunction(previousValue, value)
		}

		result[i] = models.TechnicalIndicatorsRes{
			TimestampUnix: timestamp,
			Value:         value,
			Signal:        signal,
		}
	}

	return result, nil
}

func calculateMovingAverageSignal(currentPrice, maValue float64) string {
	if maValue == 0 {
		return "neutral"
	}

	if currentPrice-maValue > 0 {
		return "buy"
	} else if currentPrice-maValue < 0 {
		return "sell"
	}
	return "neutral"
}

func calculateRsiSignal(_, rsiValue float64) string {
	if rsiValue > 70 {
		return "sell" // Overbought
	} else if rsiValue < 30 {
		return "buy" // Oversold
	}
	return "neutral"
}

func calculateMacdSignal(_, macdValue float64) string {
	if macdValue > 0 {
		return "buy"
	} else if macdValue < 0 {
		return "sell"
	}
	return "neutral"
}

func calculateHullMASignal(currentPrice, hullValue float64) string {
	if hullValue == 0 {
		return "neutral"
	}

	if currentPrice-hullValue > 0 {
		return "buy"
	} else if currentPrice-hullValue < 0 {
		return "sell"
	}
	return "neutral"
}

// Volume Weighted Moving Average Signal
func calculateVWMASignal(currentPrice, vwmaValue float64) string {
	if vwmaValue == 0 {
		return "neutral"
	}

	if currentPrice-vwmaValue > 0 {
		return "buy"
	} else if currentPrice-vwmaValue < 0 {
		return "sell"
	}
	return "neutral"
}

// Commodity Channel Index Signal
func calculateCCISignal(_, cciValue float64) string {
	if cciValue > 100 {
		return "buy" // Strong upward trend
	} else if cciValue < -100 {
		return "sell" // Strong downward trend
	}
	return "neutral" // Ranging market
}

// Stochastic %K Signal
func calculateStochasticKSignal(_, stochValue float64) string {
	if stochValue > 80 {
		return "sell" // Overbought condition
	} else if stochValue < 20 {
		return "buy" // Oversold condition
	}
	return "neutral"
}

// Ichimoku Base Line (Kijun-sen) Signal
// This is simplified - a complete Ichimoku system would consider
// multiple lines and cloud positions
func calculateIchimokuBaseLineSignal(currentPrice, baseLineValue float64) string {
	if baseLineValue == 0 {
		return "neutral"
	}

	diff := currentPrice - baseLineValue

	if diff > 0 {
		return "buy" // Price above base line indicates bullish sentiment
	} else if diff < 0 {
		return "sell" // Price below base line indicates bearish sentiment
	}
	return "neutral"
}

// Average Directional Index Signal
func calculateADXSignal(_, adxValue float64) string {
	// ADX measures trend strength, not direction
	// Usually combined with +DI and -DI for direction
	// This is a simplified version
	if adxValue > 25 {
		return "trend" // Strong trend (direction not specified by ADX alone)
	}
	return "neutral" // Weak trend or ranging market
}

// Awesome Oscillator Signal
func calculateAwesomeOscillatorSignalWithCrossover(prevAO, currentAO float64) string {
	if prevAO < 0 && currentAO > 0 {
		return "buy" // Bullish crossover
	} else if prevAO > 0 && currentAO < 0 {
		return "sell" // Bearish crossover
	}
	return "neutral"
}

// Stochastic RSI Signal
func calculateStochRSISignal(_, stochRSIValue float64) string {
	if stochRSIValue > 80 {
		return "sell" // Overbought
	} else if stochRSIValue < 20 {
		return "buy" // Oversold
	}
	return "neutral"
}

// Williams %R Signal
func calculateWilliamsRSignal(_, williamsR float64) string {
	// Williams %R typically ranges from -100 to 0
	if williamsR >= -20 { // Changed from > to >= to include exactly -20
		return "sell" // Overbought condition (near 0)
	} else if williamsR < -80 {
		return "buy" // Oversold condition (near -100)
	}
	return "neutral"
}

// Ultimate Oscillator Signal
func calculateUltimateOscillatorSignal(_, uoValue float64) string {
	if uoValue > 70 {
		return "sell" // Overbought
	} else if uoValue < 30 {
		return "buy" // Oversold
	}
	return "neutral"
}

func calculateMomentumSignal(previousMomentum, currentMomentum float64) string {
	// Bullish crossover (Momentum moves from negative to positive)
	if previousMomentum <= 0 && currentMomentum > 0 {
		return "buy"
	}

	// Bearish crossover (Momentum moves from positive to negative)
	if previousMomentum >= 0 && currentMomentum < 0 {
		return "sell"
	}

	// **Fix: Upward momentum, even if negative, is bullish**
	if currentMomentum > previousMomentum {
		return "buy"
	}

	// **Fix: Downward momentum, even if positive, is bearish**
	if currentMomentum < previousMomentum {
		return "sell"
	}

	return "neutral" // No strong signal detected
}

// parseOHLC extracts high, low, close prices from the last candle data
func parseOHLC(data models.ChartDataResponse) (float64, float64, float64, float64, error) {
	if len(data.Data.Candles) == 0 {
		return 0, 0, 0, 0, fmt.Errorf("no candle data available")
	}

	lastCandle := data.Data.Candles[len(data.Data.Candles)-1]
	if len(lastCandle) < 5 {
		return 0, 0, 0, 0, fmt.Errorf("invalid candle data: %v", lastCandle)
	}

	open, ok := lastCandle[1].(float64)
	if !ok {
		return 0, 0, 0, 0, fmt.Errorf("invalid open price type: %v", lastCandle[1])
	}

	high, ok := lastCandle[2].(float64)
	if !ok {
		return 0, 0, 0, 0, fmt.Errorf("invalid high price type: %v", lastCandle[2])
	}

	low, ok := lastCandle[3].(float64)
	if !ok {
		return 0, 0, 0, 0, fmt.Errorf("invalid low price type: %v", lastCandle[3])
	}

	close, ok := lastCandle[4].(float64)
	if !ok {
		return 0, 0, 0, 0, fmt.Errorf("invalid close price type: %v", lastCandle[4])
	}

	return open, high, low, close, nil
}

// CalculateClassicPivots calculates Classic Pivot Points
func CalculateClassicPivots(data models.ChartDataResponse) (models.PivotPoints, error) {
	var pivots models.PivotPoints

	_, high, low, close, err := parseOHLC(data)
	if err != nil {
		return pivots, err
	}

	//Calculate Classic Pivot Point
	pivots.P = (high + low + close) / 3

	//Calculate Support and Resistance levels
	pivots.R1 = (2 * pivots.P) - low
	pivots.S1 = (2 * pivots.P) - high

	pivots.R2 = pivots.P + (high - low)
	pivots.S2 = pivots.P - (high - low)

	pivots.R3 = high + 2*(pivots.P-low)
	pivots.S3 = low - 2*(high-pivots.P)

	return pivots, nil
}

// CalculateFibonacciPivots calculates Fibonacci Pivot Points
func CalculateFibonacciPivots(data models.ChartDataResponse) (models.PivotPoints, error) {
	var pivots models.PivotPoints

	_, high, low, close, err := parseOHLC(data)
	if err != nil {
		return pivots, err
	}

	//Calculate Fibonacci Pivot Point
	pivots.P = (high + low + close) / 3

	//Fibonacci ratios
	pivots.R1 = pivots.P + 0.382*(high-low)
	pivots.R2 = pivots.P + 0.618*(high-low)
	pivots.R3 = pivots.P + 1.000*(high-low)

	pivots.S1 = pivots.P - 0.382*(high-low)
	pivots.S2 = pivots.P - 0.618*(high-low)
	pivots.S3 = pivots.P - 1.000*(high-low)

	return pivots, nil
}

// CalculateWoodiePivots calculates Woodie Pivot Points
func CalculateWoodiePivots(data models.ChartDataResponse) (models.PivotPoints, error) {
	var pivots models.PivotPoints

	_, high, low, close, err := parseOHLC(data)
	if err != nil {
		return pivots, err
	}

	//For Woodie, we need tomorrow's open price, but since we don't have it, we'll use the previous close as a substitute
	//Woodie's Pivot = (H + L + 2*C) / 4
	pivots.P = (high + low + 2*close) / 4

	//Calculate Support and Resistance levels
	pivots.R1 = (2 * pivots.P) - low
	pivots.S1 = (2 * pivots.P) - high

	pivots.R2 = pivots.P + (high - low)
	pivots.S2 = pivots.P - (high - low)

	//Woodie's R3 and S3 are less common but can be calculated as:
	pivots.R3 = high + 2*(pivots.P-low)
	pivots.S3 = low - 2*(high-pivots.P)

	return pivots, nil
}

// CalculateCamarillaPivots calculates Camarilla Pivot Points
func CalculateCamarillaPivots(data models.ChartDataResponse) (models.PivotPoints, error) {
	var pivots models.PivotPoints

	_, high, low, close, err := parseOHLC(data)
	if err != nil {
		return pivots, err
	}

	//Calculate Camarilla Pivot Point (same as classic)
	pivots.P = (high + low + close) / 3

	//Calculate Camarilla Support and Resistance levels
	range_ := high - low

	pivots.R1 = close + range_*1.1/12.0
	pivots.R2 = close + range_*1.1/6.0
	pivots.R3 = close + range_*1.1/4.0

	pivots.S1 = close - range_*1.1/12.0
	pivots.S2 = close - range_*1.1/6.0
	pivots.S3 = close - range_*1.1/4.0

	return pivots, nil
}

// ParsePreviousDayOHLC extracts high, low, close from previous day's data
func ParsePreviousDayOHLC(data models.ChartDataResponse) (float64, float64, float64, error) {
	if len(data.Data.Candles) < 1 {
		return 0, 0, 0, fmt.Errorf("not enough candle data available")
	}

	previousDayCandle := data.Data.Candles[len(data.Data.Candles)-1]
	if len(previousDayCandle) < 5 {
		return 0, 0, 0, fmt.Errorf("invalid candle data: %v", previousDayCandle)
	}

	high, ok := previousDayCandle[2].(float64)
	if !ok {
		return 0, 0, 0, fmt.Errorf("invalid high price type: %v", previousDayCandle[2])
	}

	low, ok := previousDayCandle[3].(float64)
	if !ok {
		return 0, 0, 0, fmt.Errorf("invalid low price type: %v", previousDayCandle[3])
	}

	close, ok := previousDayCandle[4].(float64)
	if !ok {
		return 0, 0, 0, fmt.Errorf("invalid close price type: %v", previousDayCandle[4])
	}

	return high, low, close, nil
}

func CalculateAllPivotPoints(data models.ChartDataResponse) (map[string]models.PivotPoints, error) {
	result := make(map[string]models.PivotPoints)

	classic, err := CalculateClassicPivots(data)
	if err != nil {
		return nil, err
	}
	result["classic"] = classic

	fibonacci, err := CalculateFibonacciPivots(data)
	if err != nil {
		return nil, err
	}
	result["fibonacci"] = fibonacci

	woodie, err := CalculateWoodiePivots(data)
	if err != nil {
		return nil, err
	}
	result["woodie"] = woodie

	camarilla, err := CalculateCamarillaPivots(data)
	if err != nil {
		return nil, err
	}
	result["camarilla"] = camarilla

	return result, nil
}
