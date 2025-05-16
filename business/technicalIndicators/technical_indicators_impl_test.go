package technicalindicators

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/loggerconfig"
	"space/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSMA(t *testing.T) {
	// Test case 1: Empty chart data
	chartData := []models.TLChartCandleData{}
	expectedSMA := 0.0
	expectedError := errors.New(constants.ChartDataSizeError)

	sma, err := SMA(chartData)

	assert.Equal(t, expectedSMA, sma)
	assert.Equal(t, expectedError, err)

	// Test case 2: Non-empty chart data
	chartData = []models.TLChartCandleData{
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
	}
	expectedSMA = 15.0
	expectedError = nil

	sma, err = SMA(chartData)

	assert.Equal(t, expectedSMA, sma)
	assert.Equal(t, expectedError, err)
}

func TestEMA(t *testing.T) {
	// Test case 1: Insufficient chart data
	chartData := []models.TLChartCandleData{
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
	}
	period := 3
	expectedEMA := 0.0
	expectedError := errors.New(constants.ChartDataSizeError)

	ema, err := EMA(chartData, period)

	assert.Equal(t, expectedEMA, ema)
	assert.Equal(t, expectedError, err)

	// Test case 2: Sufficient chart data
	chartData = []models.TLChartCandleData{
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
	}
	period = 3
	expectedEMA = 15
	expectedError = nil

	ema, err = EMA(chartData, period)

	assert.Equal(t, expectedEMA, ema)
	assert.Equal(t, expectedError, err)
}

func TestRSI(t *testing.T) {
	// Test case 1: Insufficient chart data
	chartData := []models.TLChartCandleData{
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
	}
	expectedRSI := 0.0
	expectedError := errors.New(constants.ChartDataSizeError)

	rsi, err := RSI(chartData)

	assert.Equal(t, expectedRSI, rsi)
	assert.Equal(t, expectedError, err)

	// Test case 2: Sufficient chart data
	chartData = []models.TLChartCandleData{
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
	}
	expectedRSI = 55.55555555555556
	expectedError = nil

	rsi, err = RSI(chartData)

	assert.Equal(t, expectedRSI, rsi)
	assert.Equal(t, expectedError, err)
}

func TestMACD(t *testing.T) {
	// Test case 1: Insufficient chart data
	chartData := []models.TLChartCandleData{
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
	}
	expectedMACD := 0.0
	expectedError := errors.New(constants.ChartDataSizeError)

	macd, err := MACD(chartData)

	assert.Equal(t, expectedMACD, macd)
	assert.Equal(t, expectedError, err)

	// Test case 2: Sufficient chart data
	chartData = []models.TLChartCandleData{
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
	}
	expectedMACD = 0.36803500375868836
	expectedError = nil

	macd, err = MACD(chartData)

	assert.Equal(t, expectedMACD, macd)
	assert.Equal(t, expectedError, err)
}

func TestMACDSignal(t *testing.T) {
	// Test case 1: Insufficient chart data
	chartData := []models.TLChartCandleData{
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
	}
	expectedMACDSignal := 0.0
	expectedError := errors.New(constants.ChartDataSizeError)

	macdSignal, err := MACDSignal(chartData)

	assert.Equal(t, expectedMACDSignal, macdSignal)
	assert.Equal(t, expectedError, err)

	// Test case 2: Sufficient chart data
	chartData = []models.TLChartCandleData{
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
		{Timestamp: "2023-06-01", Close: 10.0},
		{Timestamp: "2023-06-02", Close: 15.0},
		{Timestamp: "2023-06-03", Close: 20.0},
	}
	expectedMACDSignal = -0.08661883408818222
	expectedError = nil

	macdSignal, err = MACDSignal(chartData)

	assert.Equal(t, expectedMACDSignal, macdSignal)
	assert.Equal(t, expectedError, err)
}

func TestAwesomeOscillator(t *testing.T) {
	// Test case 1: Insufficient chart data
	chartData := []models.TLChartCandleData{
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
	}
	expectedAwesomeOscillator := 0.0
	expectedError := errors.New(constants.ChartDataSizeError)

	awesomeOscillator, err := AwesomeOscillator(chartData)

	assert.Equal(t, expectedAwesomeOscillator, awesomeOscillator)
	assert.Equal(t, expectedError, err)

	// Test case 2: Sufficient chart data
	chartData = []models.TLChartCandleData{
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
	}
	expectedAwesomeOscillator = 0.7058823529411757
	expectedError = nil

	awesomeOscillator, err = AwesomeOscillator(chartData)

	assert.Equal(t, expectedAwesomeOscillator, awesomeOscillator)
	assert.Equal(t, expectedError, err)
}

func TestCCI(t *testing.T) {
	// Test case 1: Insufficient chart data
	chartData := []models.TLChartCandleData{
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
	}
	period := 5
	expectedCCI := 0.0
	expectedError := errors.New(constants.ChartDataSizeError)

	cci, err := CCI(chartData, period)

	assert.Equal(t, expectedCCI, cci)
	assert.Equal(t, expectedError, err)

	// Test case 2: Sufficient chart data
	chartData = []models.TLChartCandleData{
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
		{Timestamp: "2023-06-01", High: 10.0, Low: 5.0, Close: 7.5},
		{Timestamp: "2023-06-02", High: 15.0, Low: 8.0, Close: 11.5},
		{Timestamp: "2023-06-03", High: 20.0, Low: 12.0, Close: 16.0},
	}
	period = 20
	expectedCCI = 95.23809523809523
	expectedError = nil

	cci, err = CCI(chartData, period)

	assert.Equal(t, expectedCCI, cci)
	assert.Equal(t, expectedError, err)
}

func TestSupportAndResistance(t *testing.T) {
	// Test case 1
	datapoint := models.TLChartCandleData{
		High:  20.0,
		Low:   10.0,
		Close: 15.0,
	}
	expectedR1 := 20.0
	expectedR2 := 25.0
	expectedR3 := 10.0
	expectedS1 := 10.0
	expectedS2 := 5.0
	expectedS3 := -10.0

	r1, r2, r3, s1, s2, s3 := SupportAndResistance(datapoint)

	assert.Equal(t, expectedR1, r1)
	assert.Equal(t, expectedR2, r2)
	assert.Equal(t, expectedR3, r3)
	assert.Equal(t, expectedS1, s1)
	assert.Equal(t, expectedS2, s2)
	assert.Equal(t, expectedS3, s3)
}

func TestTechnicalIndicatorsObj_TechnicalIndicatorsValues(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.TechnicalIndicatorsValuesReq
		reqH models.ReqHeader
	}

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	field1 := fields{
		tradeLabURL: "http://test",
	}

	req1 := models.TechnicalIndicatorsValuesReq{
		Exchange: "NSE",
		Token:    "11536",
	}

	reqH1 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2f",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH1,
	}

	var response models.TechnicalIndicatorsValuesRes
	response.SMA = 0.89
	response.EMA = 0.89
	response.RSI = 0.89
	response.MACD = 0.89
	response.MACDSignal = 0.89
	response.CCI = 0.89
	response.AwesomeOscillator = 0.89
	response.R1 = 0.89
	response.R2 = 0.89
	response.R3 = 0.89
	response.S1 = 0.89
	response.S2 = 0.89
	response.S3 = 0.89

	res1 := apihelpers.APIRes{
		Data:    response,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock tl candles
	FetchChartData = func(url string, reqH models.ReqHeader) ([]models.TLChartCandleData, bool) {
		var candles []models.TLChartCandleData
		return candles, true
	}

	SMA = func(chartdata []models.TLChartCandleData) (float64, error) {
		return 0.89, nil
	}

	EMA = func(chartdata []models.TLChartCandleData, period int) (float64, error) {
		return 0.89, nil
	}

	RSI = func(chartdata []models.TLChartCandleData) (float64, error) {
		return 0.89, nil
	}

	MACD = func(chartdata []models.TLChartCandleData) (float64, error) {
		return 0.89, nil
	}

	MACDSignal = func(chartdata []models.TLChartCandleData) (float64, error) {
		return 0.89, nil
	}

	AwesomeOscillator = func(chartdata []models.TLChartCandleData) (float64, error) {
		return 0.89, nil
	}

	CCI = func(chartdata []models.TLChartCandleData, period int) (float64, error) {
		return 0.89, nil
	}

	SupportAndResistance = func(datapoint models.TLChartCandleData) (float64, float64, float64, float64, float64, float64) {
		return 0.89, 0.89, 0.89, 0.89, 0.89, 0.89
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", field1, arg1, http.StatusOK, res1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := TIObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.TechnicalIndicatorsValues(tt.args.req, tt.args.reqH)
			//got1 = res2
			if got != tt.want {
				t.Errorf("TIObj.TechnicalIndicatorsValues() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("TIObj.TechnicalIndicatorsValues() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end
}
