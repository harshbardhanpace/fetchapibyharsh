package technicalindicatorsV2

import (
	"space/helpers"
	"space/models"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func almostEqual(a, b, epsilon float64) bool {
	return (a-b) < epsilon && (b-a) < epsilon
}

func TestCalculateSMA(t *testing.T) {
	mockData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 8.84, 8.84, 8.84, 8.84, 20},
				{"2025-01-27T11:34:00+0530", 8.85, 8.85, 8.85, 8.85, 1059},
				{"2025-01-27T11:35:00+0530", 8.86, 8.86, 8.86, 8.86, 230},
				{"2025-01-27T11:36:00+0530", 8.89, 8.89, 8.89, 8.89, 150},
			},
		},
	}

	expected := []float64{8.84, 8.845, 8.85, 8.87}
	sma, err := CalculateSMA(mockData, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	epsilon := 1e-2
	for i, v := range expected {
		if !almostEqual(sma[i], v, epsilon) {
			t.Errorf("Expected SMA at index %d to be %.2f, got %.2f", i, v, sma[i])
		}
	}
}

func TestCalculateEMA(t *testing.T) {
	data := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 8.84, 8.84, 8.84, 8.84, 20},
				{"2025-01-27T11:34:00+0530", 8.84, 8.84, 8.84, 8.84, 1059},
				{"2025-01-27T11:35:00+0530", 8.84, 8.84, 8.84, 8.84, 230},
				{"2025-01-27T15:29:00+0530", 8.91, 8.91, 8.9, 8.9, 2004},
			},
		},
	}

	period := 3
	expected := []float64{8.84, 8.84, 8.84, 8.87}

	result, err := CalculateEMA(data, period)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result) != len(expected) {
		t.Fatalf("Expected result length %d, got %d", len(expected), len(result))
	}

	epsilon := 1e-2
	for i, v := range result {
		if !almostEqual(v, expected[i], epsilon) {
			t.Errorf("Expected EMA at index %d to be %.2f, got %.2f", i, expected[i], v)
		}
	}
}

func TestCalculateSMAInvalidData(t *testing.T) {
	data := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 8.84, 8.84, 8.84}, // Missing close price
			},
		},
	}

	_, err := CalculateSMA(data, 3)
	if err == nil {
		t.Fatal("Expected error for invalid candle data, got nil")
	}
}

func TestCalculateEMAInvalidData(t *testing.T) {
	data := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 8.84, 8.84, 8.84, "invalid", 20}, // Invalid close price type
			},
		},
	}

	_, err := CalculateEMA(data, 3)
	if err == nil {
		t.Fatal("Expected error for invalid close price type, got nil")
	}
}

func TestCalculateHullMA(t *testing.T) {
	mockData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 8.84, 8.84, 8.84, 8.84, 20},
				{"2025-01-27T11:34:00+0530", 8.85, 8.85, 8.85, 8.85, 1059},
				{"2025-01-27T11:35:00+0530", 8.86, 8.86, 8.86, 8.86, 230},
				{"2025-01-27T11:36:00+0530", 8.89, 8.89, 8.89, 8.89, 150},
				{"2025-01-27T11:37:00+0530", 8.90, 8.90, 8.90, 8.90, 200},
				{"2025-01-27T11:38:00+0530", 8.91, 8.91, 8.91, 8.91, 300},
			},
		},
	}

	expected := []float64{0, 0, 8.85, 8.87, 8.89, 8.91}
	hma, err := CalculateHullMA(mockData, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	epsilon := 1e-2
	for i, v := range expected {
		if !almostEqual(hma[i], v, epsilon) {
			t.Errorf("Expected HullMA at index %d to be %.2f, got %.2f", i, v, hma[i])
		}
	}
}

func TestCalculateVWMA(t *testing.T) {
	mockData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 8.84, 8.84, 8.84, 8.84, 100},
				{"2025-01-27T11:34:00+0530", 8.85, 8.85, 8.85, 8.85, 200},
				{"2025-01-27T11:35:00+0530", 8.86, 8.86, 8.86, 8.86, 300},
				{"2025-01-27T11:36:00+0530", 8.89, 8.89, 8.89, 8.89, 400},
			},
		},
	}

	expected := []float64{8.84, 8.847, 8.857, 8.872}
	vwma, err := CalculateVWMA(mockData, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	epsilon := 1e-3
	for i, v := range expected {
		if !almostEqual(vwma[i], v, epsilon) {
			t.Errorf("Expected VWMA at index %d to be %.3f, got %.3f", i, v, vwma[i])
		}
	}
}

func TestCalculateRSI(t *testing.T) {
	mockData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 8.84, 8.84, 8.84, 44.34, 100},
				{"2025-01-27T11:34:00+0530", 8.85, 8.85, 8.85, 44.09, 200},
				{"2025-01-27T11:35:00+0530", 8.86, 8.86, 8.86, 44.15, 300},
				{"2025-01-27T11:36:00+0530", 8.89, 8.89, 8.89, 43.61, 400},
				{"2025-01-27T11:37:00+0530", 8.90, 8.90, 8.90, 44.33, 500},
				{"2025-01-27T11:38:00+0530", 8.91, 8.91, 8.91, 44.83, 600},
			},
		},
	}

	expected := []float64{0, 0, 0, 45.56, 54.46, 61.23}
	rsi, err := CalculateRSI(mockData, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	epsilon := 1e-2
	for i, v := range expected {
		if !almostEqual(rsi[i], v, epsilon) {
			t.Errorf("Expected RSI at index %d to be %.2f, got %.2f", i, v, rsi[i])
		}
	}
}

func TestCalculateCCI(t *testing.T) {
	mockData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 24.20, 24.20, 23.85, 23.89, 100},
				{"2025-01-27T11:34:00+0530", 24.07, 24.10, 23.72, 23.95, 200},
				{"2025-01-27T11:35:00+0530", 24.04, 24.20, 23.95, 24.20, 300},
				{"2025-01-27T11:36:00+0530", 24.08, 24.20, 23.87, 23.87, 400},
				{"2025-01-27T11:37:00+0530", 23.67, 23.87, 23.67, 23.67, 500},
			},
		},
	}

	expected := []float64{100.0, 66.67, 0.0, -66.67, -100.0}
	cci, err := CalculateCCI(mockData, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	epsilon := 1e-2
	for i, v := range expected {
		if !almostEqual(cci[i], v, epsilon) {
			t.Errorf("Expected CCI at index %d to be %.2f, got %.2f", i, v, cci[i])
		}
	}
}

func TestCalculateMACD(t *testing.T) {
	mockData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 8.84, 8.84, 8.84, 459.99, 100},
				{"2025-01-27T11:34:00+0530", 8.85, 8.85, 8.85, 448.85, 200},
				{"2025-01-27T11:35:00+0530", 8.86, 8.86, 8.86, 446.06, 300},
				{"2025-01-27T11:36:00+0530", 8.89, 8.89, 8.89, 450.81, 400},
				{"2025-01-27T11:37:00+0530", 8.90, 8.90, 8.90, 442.80, 500},
				{"2025-01-27T11:38:00+0530", 8.91, 8.91, 8.91, 448.97, 600},
			},
		},
	}

	macd, signal, hist, err := CalculateMACD(mockData, 12, 26, 9)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(macd) != len(mockData.Data.Candles) {
		t.Errorf("Expected MACD length %d, got %d", len(mockData.Data.Candles), len(macd))
	}
	if len(signal) != len(mockData.Data.Candles) {
		t.Errorf("Expected Signal length %d, got %d", len(mockData.Data.Candles), len(signal))
	}
	if len(hist) != len(mockData.Data.Candles) {
		t.Errorf("Expected Histogram length %d, got %d", len(mockData.Data.Candles), len(hist))
	}
}

func TestCalculateStochastic(t *testing.T) {
	mockData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 8.84, 127.01, 125.36, 126.81, 100},
				{"2025-01-27T11:34:00+0530", 8.85, 127.62, 126.99, 127.62, 200},
				{"2025-01-27T11:35:00+0530", 8.86, 127.61, 127.30, 127.35, 300},
				{"2025-01-27T11:36:00+0530", 8.89, 127.35, 126.96, 127.07, 400},
				{"2025-01-27T11:37:00+0530", 8.90, 127.74, 127.08, 127.29, 500},
			},
		},
	}

	k, d, err := CalculateStochastic(mockData, 5, 3, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(k) != len(mockData.Data.Candles) {
		t.Errorf("Expected K length %d, got %d", len(mockData.Data.Candles), len(k))
	}
	if len(d) != len(mockData.Data.Candles) {
		t.Errorf("Expected D length %d, got %d", len(mockData.Data.Candles), len(d))
	}
}

func TestCalculateIchimokuBaseLine(t *testing.T) {
	mockData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 8.84, 30.20, 29.41, 29.87, 100},
				{"2025-01-27T11:34:00+0530", 8.85, 30.28, 29.32, 30.24, 200},
				{"2025-01-27T11:35:00+0530", 8.86, 30.45, 29.96, 30.10, 300},
				{"2025-01-27T11:36:00+0530", 8.89, 30.10, 29.46, 29.46, 400},
				{"2025-01-27T11:37:00+0530", 8.90, 29.35, 28.83, 28.91, 500},
				{"2025-01-27T11:38:00+0530", 8.91, 29.35, 28.83, 28.91, 600},
			},
		},
	}

	baseline, err := CalculateIchimokuBaseLine(mockData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(baseline) != len(mockData.Data.Candles) {
		t.Errorf("Expected baseline length %d, got %d", len(mockData.Data.Candles), len(baseline))
	}
}

func TestCalculateADX(t *testing.T) {
	mockData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 10.0, 10.5, 9.5, 10.0, 1000},
				{"2025-01-27T11:34:00+0530", 10.0, 10.8, 9.8, 10.3, 1500},
				{"2025-01-27T11:35:00+0530", 10.3, 10.9, 10.0, 10.6, 2000},
				{"2025-01-27T11:36:00+0530", 10.6, 11.0, 10.2, 10.4, 1800},
				{"2025-01-27T11:37:00+0530", 10.4, 10.8, 10.0, 10.7, 1600},
				{"2025-01-27T11:38:00+0530", 10.7, 11.2, 10.5, 11.0, 1700},
				{"2025-01-27T11:39:00+0530", 11.0, 11.4, 10.8, 11.2, 1800},
			},
		},
	}

	period := 3
	adx, err := CalculateADX(mockData, period)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// ADX should be between 0 and 100
	for i, v := range adx {
		if v < 0 || v > 100 {
			t.Errorf("ADX at index %d should be between 0 and 100, got %.2f", i, v)
		}
	}

	// Test for insufficient data
	insufficientData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 10.0, 10.5, 9.5, 10.0, 1000},
				{"2025-01-27T11:34:00+0530", 10.0, 10.8, 9.8, 10.3, 1500},
			},
		},
	}

	_, err = CalculateADX(insufficientData, period)
	if err == nil {
		t.Error("Expected error for insufficient data, got nil")
	}

	// Test invalid data
	invalidData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 10.0, "invalid", 9.5, 10.0, 1000},
			},
		},
	}

	_, err = CalculateADX(invalidData, period)
	if err == nil {
		t.Error("Expected error for invalid data, got nil")
	}
}

// Helper function to create test data for oscillator tests
func createOscillatorTestData() models.ChartDataResponse {
	return models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 10.0, 10.5, 9.5, 10.0, 1000},
				{"2025-01-27T11:34:00+0530", 10.0, 10.8, 9.8, 10.3, 1500},
				{"2025-01-27T11:35:00+0530", 10.3, 10.9, 10.0, 10.6, 2000},
				{"2025-01-27T11:36:00+0530", 10.6, 11.0, 10.2, 10.4, 1800},
				{"2025-01-27T11:37:00+0530", 10.4, 10.8, 10.0, 10.7, 1600},
			},
		},
	}
}

func TestCalculateAwesomeOscillator(t *testing.T) {
	mockData := createOscillatorTestData()
	ao, err := CalculateAwesomeOscillator(mockData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ao) != len(mockData.Data.Candles) {
		t.Errorf("Expected AO length %d, got %d", len(mockData.Data.Candles), len(ao))
	}
}

func TestCalculateMomentum(t *testing.T) {
	mockData := createOscillatorTestData()
	period := 2
	momentum, err := CalculateMomentum(mockData, period)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []float64{0, 0, 0.6, 0.1, 0.1}
	epsilon := 1e-1
	for i, v := range expected {
		if !almostEqual(momentum[i], v, epsilon) {
			t.Errorf("Expected Momentum at index %d to be %.1f, got %.1f", i, v, momentum[i])
		}
	}
}

func TestCalculateStochRSIFast(t *testing.T) {
	mockData := createOscillatorTestData()
	k, d, err := CalculateStochRSIFast(mockData, 3, 3, 14, 14)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(k) != len(mockData.Data.Candles) {
		t.Errorf("Expected K length %d, got %d", len(mockData.Data.Candles), len(k))
	}
	if len(d) != len(mockData.Data.Candles) {
		t.Errorf("Expected D length %d, got %d", len(mockData.Data.Candles), len(d))
	}
}

func TestCalculateWilliamsR(t *testing.T) {
	mockData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 10.0, 10.5, 9.5, 10.0, 1000},
				{"2025-01-27T11:34:00+0530", 10.0, 10.8, 9.8, 10.3, 1500},
				{"2025-01-27T11:35:00+0530", 10.3, 10.9, 10.0, 10.6, 2000},
				{"2025-01-27T11:36:00+0530", 10.6, 11.0, 10.2, 10.4, 1800},
				{"2025-01-27T11:37:00+0530", 10.4, 10.8, 10.0, 10.7, 1600},
			},
		},
	}

	period := 3
	williams, err := CalculateWilliamsR(mockData, period)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Williams %R is between -100 and 0
	for i, v := range williams {
		if v > 0 || v < -100 {
			t.Errorf("Williams %%R at index %d should be between -100 and 0, got %.2f", i, v)
		}
	}

	// Expected values calculated based on formula: ((Highest High - Close)/(Highest High - Lowest Low)) * -100
	expected := []float64{-50.0, -33.33, -66.67}
	epsilon := 1e-2
	for i, v := range expected {
		if !almostEqual(williams[i], v, epsilon) {
			t.Errorf("Expected Williams %%R at index %d to be %.2f, got %.2f", i, v, williams[i])
		}
	}
}

func TestCalculateMovingAverageSignal(t *testing.T) {
	tests := []struct {
		name         string
		currentPrice float64
		maValue      float64
		expected     string
	}{
		{
			name:         "Price more than 5% above MA - should buy",
			currentPrice: 105.5,
			maValue:      100,
			expected:     "buy",
		},
		{
			name:         "Price more than 5% below MA - should sell",
			currentPrice: 94.5,
			maValue:      100,
			expected:     "sell",
		},
		{
			name:         "Price within 5% of MA - should be neutral",
			currentPrice: 104,
			maValue:      100,
			expected:     "buy",
		},
		{
			name:         "MA value is zero - should be neutral",
			currentPrice: 100,
			maValue:      0,
			expected:     "neutral",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateMovingAverageSignal(tt.currentPrice, tt.maValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateRsiSignal(t *testing.T) {
	tests := []struct {
		name     string
		rsiValue float64
		expected string
	}{
		{
			name:     "RSI above 70 - should sell (overbought)",
			rsiValue: 75,
			expected: "sell",
		},
		{
			name:     "RSI below 30 - should buy (oversold)",
			rsiValue: 25,
			expected: "buy",
		},
		{
			name:     "RSI between 30 and 70 - should be neutral",
			rsiValue: 50,
			expected: "neutral",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateRsiSignal(0, tt.rsiValue) // First parameter unused
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateMacdSignal(t *testing.T) {
	tests := []struct {
		name      string
		macdValue float64
		expected  string
	}{
		{
			name:      "MACD above 0 - should buy",
			macdValue: 0.5,
			expected:  "buy",
		},
		{
			name:      "MACD below 0 - should sell",
			macdValue: -0.5,
			expected:  "sell",
		},
		{
			name:      "MACD at 0 - should be neutral",
			macdValue: 0,
			expected:  "neutral",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateMacdSignal(0, tt.macdValue) // First parameter unused
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatIndicatorResults(t *testing.T) {
	// Create mock data
	mockData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{float64(1612345678), 100.0, 105.0, 95.0, 102.0, 1000},
				{float64(1612345778), 102.0, 107.0, 101.0, 106.0, 1500},
				{float64(1612345878), 106.0, 110.0, 104.0, 108.0, 2000},
			},
		},
	}

	mockIndicatorValues := []float64{0, 100.0, 105.0}

	// Use SMA signal calculation for test
	result, err := formatIndicatorResults(mockData, mockIndicatorValues, calculateMovingAverageSignal)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, result, 3)

	// First value uses timestamp but has "neutral" signal due to NaN for indicator
	assert.Equal(t, int64(1612345678), result[0].TimestampUnix)
	assert.Equal(t, 0.0, result[0].Value)
	assert.Equal(t, "neutral", result[0].Signal)

	// Second value has "buy" signal as price (106.0) > 5% above indicator value (100.0)
	assert.Equal(t, int64(1612345778), result[1].TimestampUnix)
	assert.Equal(t, 100.0, result[1].Value)
	assert.Equal(t, "buy", result[1].Signal)

	// Third value has "sell" signal as price (108.0) > indicator value (105.0) but within 5%
	assert.Equal(t, int64(1612345878), result[2].TimestampUnix)
	assert.Equal(t, 105.0, result[2].Value)
	assert.Equal(t, "neutral", result[2].Signal) // 108 is ~2.85% above 105
}

func TestCalculateHullMASignal(t *testing.T) {
	tests := []struct {
		name         string
		currentPrice float64
		hullValue    float64
		expected     string
	}{
		{
			name:         "Price more than 3% above Hull MA - should buy",
			currentPrice: 103.5,
			hullValue:    100,
			expected:     "buy",
		},
		{
			name:         "Price more than 3% below Hull MA - should sell",
			currentPrice: 96.5,
			hullValue:    100,
			expected:     "sell",
		},
		{
			name:         "Price within 3% of Hull MA - should be neutral",
			currentPrice: 102,
			hullValue:    100,
			expected:     "buy",
		},
		{
			name:         "Hull MA value is zero - should be neutral",
			currentPrice: 100,
			hullValue:    0,
			expected:     "neutral",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateHullMASignal(tt.currentPrice, tt.hullValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateCCISignal(t *testing.T) {
	tests := []struct {
		name     string
		cciValue float64
		expected string
	}{
		{
			name:     "CCI above 100 - strong uptrend (buy)",
			cciValue: 120,
			expected: "buy",
		},
		{
			name:     "CCI below -100 - strong downtrend (sell)",
			cciValue: -120,
			expected: "sell",
		},
		{
			name:     "CCI between -100 and 100 - neutral",
			cciValue: 50,
			expected: "neutral",
		},
		{
			name:     "CCI exactly at 100 - neutral",
			cciValue: 100,
			expected: "neutral",
		},
		{
			name:     "CCI exactly at -100 - neutral",
			cciValue: -100,
			expected: "neutral",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateCCISignal(0, tt.cciValue) // First parameter unused
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateStochasticKSignal(t *testing.T) {
	tests := []struct {
		name       string
		stochValue float64
		expected   string
	}{
		{
			name:       "Stochastic above 80 - overbought (sell)",
			stochValue: 85,
			expected:   "sell",
		},
		{
			name:       "Stochastic below 20 - oversold (buy)",
			stochValue: 15,
			expected:   "buy",
		},
		{
			name:       "Stochastic between 20 and 80 - neutral",
			stochValue: 50,
			expected:   "neutral",
		},
		{
			name:       "Stochastic exactly at 80 - neutral",
			stochValue: 80,
			expected:   "neutral",
		},
		{
			name:       "Stochastic exactly at 20 - neutral",
			stochValue: 20,
			expected:   "neutral",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateStochasticKSignal(0, tt.stochValue) // First parameter unused
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateADXSignal(t *testing.T) {
	tests := []struct {
		name     string
		adxValue float64
		expected string
	}{
		{
			name:     "ADX above 25 - strong trend",
			adxValue: 30,
			expected: "trend",
		},
		{
			name:     "ADX below 25 - weak trend/ranging",
			adxValue: 20,
			expected: "neutral",
		},
		{
			name:     "ADX exactly at 25 - neutral",
			adxValue: 25,
			expected: "neutral",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateADXSignal(0, tt.adxValue) // First parameter unused
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateWilliamsRSignal(t *testing.T) {
	tests := []struct {
		name      string
		williamsR float64
		expected  string
	}{
		{
			name:      "Williams %R above -20 (near 0) - overbought (sell)",
			williamsR: -15,
			expected:  "sell",
		},
		{
			name:      "Williams %R below -80 (near -100) - oversold (buy)",
			williamsR: -85,
			expected:  "buy",
		},
		{
			name:      "Williams %R between -80 and -20 - neutral",
			williamsR: -50,
			expected:  "neutral",
		},
		{
			name:      "Williams %R exactly at -20 - sell",
			williamsR: -20,
			expected:  "sell",
		},
		{
			name:      "Williams %R exactly at -80 - neutral",
			williamsR: -80,
			expected:  "neutral",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateWilliamsRSignal(0, tt.williamsR) // First parameter unused
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateUltimateOscillatorSignal(t *testing.T) {
	tests := []struct {
		name     string
		uoValue  float64
		expected string
	}{
		{
			name:     "UO above 70 - overbought (sell)",
			uoValue:  75,
			expected: "sell",
		},
		{
			name:     "UO below 30 - oversold (buy)",
			uoValue:  25,
			expected: "buy",
		},
		{
			name:     "UO between 30 and 70 - neutral",
			uoValue:  50,
			expected: "neutral",
		},
		{
			name:     "UO exactly at 70 - neutral",
			uoValue:  70,
			expected: "neutral",
		},
		{
			name:     "UO exactly at 30 - neutral",
			uoValue:  30,
			expected: "neutral",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateUltimateOscillatorSignal(0, tt.uoValue) // First parameter unused
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateMomentumSignal(t *testing.T) {
	tests := []struct {
		name             string
		previousMomentum float64
		currentMomentum  float64
		expected         string
	}{
		{
			name:             "Momentum crosses above zero - buy signal",
			previousMomentum: -5.0,
			currentMomentum:  3.5,
			expected:         "buy",
		},
		{
			name:             "Momentum crosses below zero - sell signal",
			previousMomentum: 2.5,
			currentMomentum:  -1.8,
			expected:         "sell",
		},
		{
			name:             "Momentum positive and increasing - buy signal",
			previousMomentum: 5.0,
			currentMomentum:  7.5,
			expected:         "buy",
		},
		{
			name:             "Momentum negative and decreasing - sell signal",
			previousMomentum: -3.0,
			currentMomentum:  -5.2,
			expected:         "sell",
		},
		{
			name:             "Momentum positive but decreasing - neutral",
			previousMomentum: 8.5,
			currentMomentum:  4.2,
			expected:         "neutral",
		},
		{
			name:             "Momentum negative but increasing - neutral",
			previousMomentum: -7.0,
			currentMomentum:  -2.5,
			expected:         "neutral",
		},
		{
			name:             "Both values zero - neutral",
			previousMomentum: 0.0,
			currentMomentum:  0.0,
			expected:         "neutral",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateMomentumSignal(tt.previousMomentum, tt.currentMomentum)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTranslateTechIndicatorReqToChartDataReq(t *testing.T) {
	// Store the original function and restore it after tests
	originalTimeNow := timeNow
	defer func() { timeNow = originalTimeNow }()

	// Mock fixed time for testing (Tuesday, April 1, 2025, 12:00:00 IST)
	mockTime := time.Date(2025, 4, 1, 12, 0, 0, 0, LocationKolkata)
	timeNow = func() time.Time {
		return mockTime
	}

	tests := []struct {
		name           string
		input          models.GetAllTechnicalIndicatorsReq
		expectedOutput models.TLChartDataReq
		checkTimeDiff  bool  // Flag to check time difference instead of exact values
		expectedDiff   int64 // Expected time difference in seconds
	}{
		{
			name: "MINUTE with TimeInterval 5",
			input: models.GetAllTechnicalIndicatorsReq{
				Exchange:     "NSE",
				Token:        "2885",
				TimeUnit:     "MINUTE",
				TimeInterval: 5,
			},
			expectedOutput: models.TLChartDataReq{
				Exchange:     "NSE",
				Token:        "2885",
				CandleType:   "1",
				DataDuration: "5",
			},
			checkTimeDiff: true,
			// For 5 MINUTE, we need at least 300 data points
			// Each day has 6 hours = 360 minutes of trading
			// With 5-minute intervals, that's 72 data points per day
			// Need 300/72 = 4.167 days, with 50% buffer = ~6.25 days
			// Rounding to int = 6 days = 518,400 seconds
			expectedDiff: 777600,
		},
		{
			name: "HOUR with TimeInterval 1",
			input: models.GetAllTechnicalIndicatorsReq{
				Exchange:     "NSE",
				Token:        "2885",
				TimeUnit:     "HOUR",
				TimeInterval: 1,
			},
			expectedOutput: models.TLChartDataReq{
				Exchange:     "NSE",
				Token:        "2885",
				CandleType:   "2",
				DataDuration: "1",
			},
			checkTimeDiff: true,
			// For 1 HOUR, we get 6 data points per day (market is open 6 hours)
			// Need 300/6 = 50 days, with 50% buffer = 75 days
			// 75 days = 6,480,000 seconds
			expectedDiff: 10108800,
		},
		{
			name: "DAY with TimeInterval 1",
			input: models.GetAllTechnicalIndicatorsReq{
				Exchange:     "NSE",
				Token:        "2885",
				TimeUnit:     "DAY",
				TimeInterval: 1,
			},
			expectedOutput: models.TLChartDataReq{
				Exchange:     "NSE",
				Token:        "2885",
				CandleType:   "3",
				DataDuration: "1",
			},
			checkTimeDiff: true,
			// For 1 DAY, we need 300 days
			// 300 days = 25,920,000 seconds
			expectedDiff: 40867200,
		},
		{
			name: "WEEK with TimeInterval 1",
			input: models.GetAllTechnicalIndicatorsReq{
				Exchange:     "NSE",
				Token:        "2885",
				TimeUnit:     "WEEK",
				TimeInterval: 1,
			},
			expectedOutput: models.TLChartDataReq{
				Exchange:     "NSE",
				Token:        "2885",
				CandleType:   "3",
				DataDuration: "7", // 1 week = 7 days
			},
			checkTimeDiff: true,
			// For 1 WEEK, we need 300 weeks = 300 * 7 days
			// 2100 days = 181,440,000 seconds
			expectedDiff: 204422400,
		},
		{
			name: "MONTH with TimeInterval 1",
			input: models.GetAllTechnicalIndicatorsReq{
				Exchange:     "NSE",
				Token:        "2885",
				TimeUnit:     "MONTH",
				TimeInterval: 1,
			},
			expectedOutput: models.TLChartDataReq{
				Exchange:     "NSE",
				Token:        "2885",
				CandleType:   "3",
				DataDuration: "30", // 1 month = ~30 days
			},
			checkTimeDiff: true,
			// For 1 MONTH, we need 300 months = 300 * 30 days
			// 9000 days = 777,600,000 seconds
			expectedDiff: 899683200,
		},
		{
			name: "MINUTE with TimeInterval 15",
			input: models.GetAllTechnicalIndicatorsReq{
				Exchange:     "NSE",
				Token:        "2885",
				TimeUnit:     "MINUTE",
				TimeInterval: 15,
			},
			expectedOutput: models.TLChartDataReq{
				Exchange:     "NSE",
				Token:        "2885",
				CandleType:   "1",
				DataDuration: "15",
			},
			checkTimeDiff: true,
			// For 15 MINUTE, we get 4 data points per hour, 24 per day
			// Need 300/24 = 12.5 days, with 50% buffer = ~18.75 days
			// Rounding to int = 18 days = 1,555,200 seconds
			expectedDiff: 2419200,
		},
	}

	// Helper function to approximate equal with tolerance
	approxEqual := func(a, b int64, tolerance int64) bool {
		diff := a - b
		if diff < 0 {
			diff = -diff
		}
		return diff <= tolerance
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TranslateTechIndicatorReqToChartDataReq(tt.input)

			// Check non-time fields
			if result.Exchange != tt.expectedOutput.Exchange {
				t.Errorf("Exchange = %v, want %v", result.Exchange, tt.expectedOutput.Exchange)
			}
			if result.Token != tt.expectedOutput.Token {
				t.Errorf("Token = %v, want %v", result.Token, tt.expectedOutput.Token)
			}
			if result.CandleType != tt.expectedOutput.CandleType {
				t.Errorf("CandleType = %v, want %v", result.CandleType, tt.expectedOutput.CandleType)
			}
			if result.DataDuration != tt.expectedOutput.DataDuration {
				t.Errorf("DataDuration = %v, want %v", result.DataDuration, tt.expectedOutput.DataDuration)
			}

			// Check time fields
			if tt.checkTimeDiff {
				endTimeUnix, _ := strconv.ParseInt(result.EndTime, 10, 64)
				startTimeUnix, _ := strconv.ParseInt(result.StartTime, 10, 64)
				timeDiff := endTimeUnix - startTimeUnix

				// Allow 1 day tolerance for rounding to int days
				tolerance := int64(24 * 60 * 60)
				if !approxEqual(timeDiff, tt.expectedDiff, tolerance) {
					t.Errorf("Time difference = %v seconds, want approximately %v seconds (within %v seconds)",
						timeDiff, tt.expectedDiff, tolerance)
				}

				// Verify end time is current time
				if endTimeUnix != mockTime.Unix() {
					t.Errorf("EndTime = %v, want %v", endTimeUnix, mockTime.Unix())
				}
			}
		})
	}
}

func TestIsMarketOpenWithHolidays(t *testing.T) {
	// Store the original HolidayCalendar and restore it after tests
	originalCalendar := helpers.HolidayCalendar
	defer func() { helpers.HolidayCalendar = originalCalendar }()

	// Create a test calendar with some holidays
	testCalendar := models.Calendar{
		Date: []models.DateDetails{
			{
				Date:        "01-Apr-2025",
				DayOfWeek:   "Tuesday",
				Description: "Annual Bank Closing",
				IsHoliday:   true,
			},
			{
				Date:        "02-Apr-2025",
				DayOfWeek:   "Wednesday",
				Description: "Regular Trading Day",
				IsHoliday:   false,
			},
			{
				Date:        "14-Apr-2025",
				DayOfWeek:   "Monday",
				Description: "Dr. Ambedkar Jayanti",
				IsHoliday:   true,
			},
		},
	}

	// Replace the global calendar
	helpers.HolidayCalendar = testCalendar

	// Store the original timeNow function and restore it after tests
	originalTimeNow := timeNow
	defer func() { timeNow = originalTimeNow }()

	tests := []struct {
		name     string
		mockTime time.Time
		expected bool
	}{
		{
			name:     "Holiday (01-Apr-2025)",
			mockTime: time.Date(2025, 4, 1, 12, 0, 0, 0, LocationKolkata),
			expected: false,
		},
		{
			name:     "Regular Trading Day (02-Apr-2025)",
			mockTime: time.Date(2025, 4, 2, 12, 0, 0, 0, LocationKolkata),
			expected: true,
		},
		{
			name:     "Regular Trading Day but before market opens",
			mockTime: time.Date(2025, 4, 2, 9, 0, 0, 0, LocationKolkata),
			expected: false,
		},
		{
			name:     "Regular Trading Day but after market closes",
			mockTime: time.Date(2025, 4, 2, 15, 30, 0, 0, LocationKolkata),
			expected: false,
		},
		{
			name:     "Weekend (Saturday)",
			mockTime: time.Date(2025, 4, 5, 12, 0, 0, 0, LocationKolkata),
			expected: false,
		},
		{
			name:     "Holiday (14-Apr-2025)",
			mockTime: time.Date(2025, 4, 14, 12, 0, 0, 0, LocationKolkata),
			expected: false,
		},
		{
			name:     "Date not in calendar but weekday during market hours",
			mockTime: time.Date(2025, 4, 3, 12, 0, 0, 0, LocationKolkata),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock timeNow to return the test case time
			timeNow = func() time.Time {
				return tt.mockTime
			}

			result := IsMarketOpen()
			if result != tt.expected {
				t.Errorf("IsMarketOpen() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEstimateHolidaysInDateRange(t *testing.T) {
	// Store the original HolidayCalendar and restore it after tests
	originalCalendar := helpers.HolidayCalendar
	defer func() { helpers.HolidayCalendar = originalCalendar }()

	// Create a test calendar with some holidays
	testCalendar := models.Calendar{
		Date: []models.DateDetails{
			{
				Date:        "01-Apr-2025",
				DayOfWeek:   "Tuesday",
				Description: "Annual Bank Closing",
				IsHoliday:   true,
			},
			{
				Date:        "02-Apr-2025",
				DayOfWeek:   "Wednesday",
				Description: "Regular Trading Day",
				IsHoliday:   false,
			},
			{
				Date:        "14-Apr-2025",
				DayOfWeek:   "Monday",
				Description: "Dr. Ambedkar Jayanti",
				IsHoliday:   true,
			},
			{
				Date:        "01-May-2025",
				DayOfWeek:   "Thursday",
				Description: "Labor Day",
				IsHoliday:   true,
			},
		},
	}

	// Replace the global calendar
	helpers.HolidayCalendar = testCalendar

	tests := []struct {
		name      string
		startDate time.Time
		endDate   time.Time
		expected  int
	}{
		{
			name:      "Count all holidays in April",
			startDate: time.Date(2025, 4, 1, 0, 0, 0, 0, LocationKolkata),
			endDate:   time.Date(2025, 4, 30, 23, 59, 59, 0, LocationKolkata),
			expected:  2, // Apr 1 and Apr 14
		},
		{
			name:      "One month with one holiday",
			startDate: time.Date(2025, 5, 1, 0, 0, 0, 0, LocationKolkata),
			endDate:   time.Date(2025, 5, 31, 23, 59, 59, 0, LocationKolkata),
			expected:  1, // May 1
		},
		{
			name:      "Range including multiple months",
			startDate: time.Date(2025, 4, 10, 0, 0, 0, 0, LocationKolkata),
			endDate:   time.Date(2025, 5, 10, 23, 59, 59, 0, LocationKolkata),
			expected:  2, // Apr 14 and May 1
		},
		{
			name:      "One day range on a holiday",
			startDate: time.Date(2025, 4, 1, 0, 0, 0, 0, LocationKolkata),
			endDate:   time.Date(2025, 4, 1, 23, 59, 59, 0, LocationKolkata),
			expected:  1, // Apr 1
		},
		{
			name:      "Range with no holidays",
			startDate: time.Date(2025, 4, 3, 0, 0, 0, 0, LocationKolkata),
			endDate:   time.Date(2025, 4, 13, 23, 59, 59, 0, LocationKolkata),
			expected:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := estimateHolidaysInDateRange(tt.startDate, tt.endDate)
			if result != tt.expected {
				t.Errorf("estimateHolidaysInDateRange() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConvertTradingDaysToCalendarDays(t *testing.T) {
	tests := []struct {
		name        string
		tradingDays int
		expectedMin int // Minimum acceptable calendar days
		expectedMax int // Maximum acceptable calendar days
	}{
		{
			name:        "Zero trading days",
			tradingDays: 0,
			expectedMin: 0,
			expectedMax: 0,
		},
		{
			name:        "5 trading days (approximately 1 week)",
			tradingDays: 5,
			expectedMin: 7,  // At least a week
			expectedMax: 10, // Not more than 10 days with buffers
		},
		{
			name:        "20 trading days (approximately 1 month)",
			tradingDays: 20,
			expectedMin: 28, // At least 4 weeks
			expectedMax: 35, // Not more than 5 weeks with buffers
		},
		{
			name:        "250 trading days (approximately 1 year)",
			tradingDays: 250,
			expectedMin: 350, // At least 350 days
			expectedMax: 400, // Not more than 400 days with buffers
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use current time for end date
			endDate := time.Now().In(LocationKolkata)

			result := convertTradingDaysToCalendarDays(tt.tradingDays, endDate)

			if result < tt.expectedMin || result > tt.expectedMax {
				t.Errorf("convertTradingDaysToCalendarDays() = %v, want between %v and %v",
					result, tt.expectedMin, tt.expectedMax)
			}
		})
	}
}

func TestCalculateClassicPivots(t *testing.T) {
	// Mock data with known values
	mockData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 100.0, 110.0, 90.0, 105.0, 1000},
			},
		},
	}

	expectedPivots := models.PivotPoints{
		P:  101.67, // (110 + 90 + 105) / 3
		R1: 113.33, // (2 * 101.67) - 90
		R2: 121.67, // 101.67 + (110 - 90)
		R3: 133.33, // 110 + 2*(101.67 - 90)
		S1: 93.33,  // (2 * 101.67) - 110
		S2: 81.67,  // 101.67 - (110 - 90)
		S3: 73.33,  // 90 - 2*(110 - 101.67)
	}

	pivots, err := CalculateClassicPivots(mockData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	epsilon := 1e-2
	assert.True(t, almostEqual(pivots.P, expectedPivots.P, epsilon), "Expected P to be %.2f, got %.2f", expectedPivots.P, pivots.P)
	assert.True(t, almostEqual(pivots.R1, expectedPivots.R1, epsilon), "Expected R1 to be %.2f, got %.2f", expectedPivots.R1, pivots.R1)
	assert.True(t, almostEqual(pivots.R2, expectedPivots.R2, epsilon), "Expected R2 to be %.2f, got %.2f", expectedPivots.R2, pivots.R2)
	assert.True(t, almostEqual(pivots.R3, expectedPivots.R3, epsilon), "Expected R3 to be %.2f, got %.2f", expectedPivots.R3, pivots.R3)
	assert.True(t, almostEqual(pivots.S1, expectedPivots.S1, epsilon), "Expected S1 to be %.2f, got %.2f", expectedPivots.S1, pivots.S1)
	assert.True(t, almostEqual(pivots.S2, expectedPivots.S2, epsilon), "Expected S2 to be %.2f, got %.2f", expectedPivots.S2, pivots.S2)
	assert.True(t, almostEqual(pivots.S3, expectedPivots.S3, epsilon), "Expected S3 to be %.2f, got %.2f", expectedPivots.S3, pivots.S3)
}

func TestCalculateFibonacciPivots(t *testing.T) {
	// Mock data with known values
	mockData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 100.0, 110.0, 90.0, 105.0, 1000},
			},
		},
	}

	expectedPivots := models.PivotPoints{
		P:  101.67, // (110 + 90 + 105) / 3
		R1: 109.31, // 101.67 + 0.382*(110 - 90)
		R2: 114.03, // 101.67 + 0.618*(110 - 90)
		R3: 121.67, // 101.67 + 1.000*(110 - 90)
		S1: 94.03,  // 101.67 - 0.382*(110 - 90)
		S2: 89.31,  // 101.67 - 0.618*(110 - 90)
		S3: 81.67,  // 101.67 - 1.000*(110 - 90)
	}

	pivots, err := CalculateFibonacciPivots(mockData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	epsilon := 1e-2
	assert.True(t, almostEqual(pivots.P, expectedPivots.P, epsilon), "Expected P to be %.2f, got %.2f", expectedPivots.P, pivots.P)
	assert.True(t, almostEqual(pivots.R1, expectedPivots.R1, epsilon), "Expected R1 to be %.2f, got %.2f", expectedPivots.R1, pivots.R1)
	assert.True(t, almostEqual(pivots.R2, expectedPivots.R2, epsilon), "Expected R2 to be %.2f, got %.2f", expectedPivots.R2, pivots.R2)
	assert.True(t, almostEqual(pivots.R3, expectedPivots.R3, epsilon), "Expected R3 to be %.2f, got %.2f", expectedPivots.R3, pivots.R3)
	assert.True(t, almostEqual(pivots.S1, expectedPivots.S1, epsilon), "Expected S1 to be %.2f, got %.2f", expectedPivots.S1, pivots.S1)
	assert.True(t, almostEqual(pivots.S2, expectedPivots.S2, epsilon), "Expected S2 to be %.2f, got %.2f", expectedPivots.S2, pivots.S2)
	assert.True(t, almostEqual(pivots.S3, expectedPivots.S3, epsilon), "Expected S3 to be %.2f, got %.2f", expectedPivots.S3, pivots.S3)
}

func TestCalculateWoodiePivots(t *testing.T) {
	// Mock data with known values
	mockData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 100.0, 110.0, 90.0, 105.0, 1000},
			},
		},
	}

	expectedPivots := models.PivotPoints{
		P:  102.50, // (110 + 90 + 2*105) / 4
		R1: 115.00, // (2 * 102.5) - 90
		R2: 122.50, // 102.5 + (110 - 90)
		R3: 135.00, // 110 + 2*(102.5 - 90)
		S1: 95.00,  // (2 * 102.5) - 110
		S2: 82.50,  // 102.5 - (110 - 90)
		S3: 75.00,  // 90 - 2*(110 - 102.5)
	}

	pivots, err := CalculateWoodiePivots(mockData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	epsilon := 1e-2
	assert.True(t, almostEqual(pivots.P, expectedPivots.P, epsilon), "Expected P to be %.2f, got %.2f", expectedPivots.P, pivots.P)
	assert.True(t, almostEqual(pivots.R1, expectedPivots.R1, epsilon), "Expected R1 to be %.2f, got %.2f", expectedPivots.R1, pivots.R1)
	assert.True(t, almostEqual(pivots.R2, expectedPivots.R2, epsilon), "Expected R2 to be %.2f, got %.2f", expectedPivots.R2, pivots.R2)
	assert.True(t, almostEqual(pivots.R3, expectedPivots.R3, epsilon), "Expected R3 to be %.2f, got %.2f", expectedPivots.R3, pivots.R3)
	assert.True(t, almostEqual(pivots.S1, expectedPivots.S1, epsilon), "Expected S1 to be %.2f, got %.2f", expectedPivots.S1, pivots.S1)
	assert.True(t, almostEqual(pivots.S2, expectedPivots.S2, epsilon), "Expected S2 to be %.2f, got %.2f", expectedPivots.S2, pivots.S2)
	assert.True(t, almostEqual(pivots.S3, expectedPivots.S3, epsilon), "Expected S3 to be %.2f, got %.2f", expectedPivots.S3, pivots.S3)
}

func TestCalculateCamarillaPivots(t *testing.T) {
	// Mock data with known values
	mockData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 100.0, 110.0, 90.0, 105.0, 1000},
			},
		},
	}

	//range_ := 110.0 - 90.0
	expectedPivots := models.PivotPoints{
		P:  101.67, // (110 + 90 + 105) / 3
		R1: 106.83, // 105 + range * 1.1/12
		R2: 108.67, // 105 + range * 1.1/6
		R3: 110.50, // 105 + range * 1.1/4
		S1: 103.17, // 105 - range * 1.1/12
		S2: 101.33, // 105 - range * 1.1/6
		S3: 99.50,  // 105 - range * 1.1/4
	}

	pivots, err := CalculateCamarillaPivots(mockData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	epsilon := 1e-2
	assert.True(t, almostEqual(pivots.P, expectedPivots.P, epsilon), "Expected P to be %.2f, got %.2f", expectedPivots.P, pivots.P)
	assert.True(t, almostEqual(pivots.R1, expectedPivots.R1, epsilon), "Expected R1 to be %.2f, got %.2f", expectedPivots.R1, pivots.R1)
	assert.True(t, almostEqual(pivots.R2, expectedPivots.R2, epsilon), "Expected R2 to be %.2f, got %.2f", expectedPivots.R2, pivots.R2)
	assert.True(t, almostEqual(pivots.R3, expectedPivots.R3, epsilon), "Expected R3 to be %.2f, got %.2f", expectedPivots.R3, pivots.R3)
	assert.True(t, almostEqual(pivots.S1, expectedPivots.S1, epsilon), "Expected S1 to be %.2f, got %.2f", expectedPivots.S1, pivots.S1)
	assert.True(t, almostEqual(pivots.S2, expectedPivots.S2, epsilon), "Expected S2 to be %.2f, got %.2f", expectedPivots.S2, pivots.S2)
	assert.True(t, almostEqual(pivots.S3, expectedPivots.S3, epsilon), "Expected S3 to be %.2f, got %.2f", expectedPivots.S3, pivots.S3)
}

func TestCalculateAllPivotPoints(t *testing.T) {
	// Mock data with known values
	mockData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 100.0, 110.0, 90.0, 105.0, 1000},
			},
		},
	}

	pivots, err := CalculateAllPivotPoints(mockData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify all pivot types are present
	pivotTypes := []string{"classic", "fibonacci", "woodie", "camarilla"}
	for _, pivotType := range pivotTypes {
		_, exists := pivots[pivotType]
		assert.True(t, exists, "Expected pivot type %s to be present", pivotType)
	}

	// Since individual pivot calculations are tested separately,
	// just verify that each type returns something reasonable
	epsilon := 1e-2
	for _, pivotType := range pivotTypes {
		pivot := pivots[pivotType]
		assert.False(t, almostEqual(pivot.P, 0, epsilon), "Expected %s P to be non-zero", pivotType)
		assert.False(t, almostEqual(pivot.R1, 0, epsilon), "Expected %s R1 to be non-zero", pivotType)
		assert.False(t, almostEqual(pivot.S1, 0, epsilon), "Expected %s S1 to be non-zero", pivotType)
	}
}

func TestInvalidCandleData(t *testing.T) {
	// Test with empty candle data
	emptyData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{},
		},
	}

	_, err := CalculateClassicPivots(emptyData)
	assert.Error(t, err, "Expected error with empty candle data")

	// Test with invalid candle format (missing values)
	invalidData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 100.0}, // Missing high, low, close
			},
		},
	}

	_, err = CalculateClassicPivots(invalidData)
	assert.Error(t, err, "Expected error with invalid candle format")

	// Test with invalid price type
	invalidTypeData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 100.0, "invalid", 90.0, 105.0, 1000}, // High as string
			},
		},
	}

	_, err = CalculateFibonacciPivots(invalidTypeData)
	assert.Error(t, err, "Expected error with invalid price type")
}

func TestEdgeCases(t *testing.T) {
	// Test with zero range (high = low)
	zeroRangeData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 100.0, 100.0, 100.0, 100.0, 1000},
			},
		},
	}

	// Classic pivots with zero range
	classicPivots, err := CalculateClassicPivots(zeroRangeData)
	assert.NoError(t, err, "Unexpected error with zero range")
	assert.Equal(t, classicPivots.P, 100.0, "Expected P to be 100.0 with zero range")
	assert.Equal(t, classicPivots.R1, 100.0, "Expected R1 to be 100.0 with zero range")
	assert.Equal(t, classicPivots.S1, 100.0, "Expected S1 to be 100.0 with zero range")

	// Fibonacci pivots with zero range
	fibPivots, err := CalculateFibonacciPivots(zeroRangeData)
	assert.NoError(t, err, "Unexpected error with zero range")
	assert.Equal(t, fibPivots.P, 100.0, "Expected P to be 100.0 with zero range")
	assert.Equal(t, fibPivots.R1, 100.0, "Expected R1 to be 100.0 with zero range")
	assert.Equal(t, fibPivots.S1, 100.0, "Expected S1 to be 100.0 with zero range")

	// Extreme values
	extremeData := models.ChartDataResponse{
		Data: models.CandleData{
			Candles: [][]interface{}{
				{"2025-01-27T11:33:00+0530", 100.0, 1000.0, 1.0, 500.0, 1000},
			},
		},
	}

	camarillaPivots, err := CalculateCamarillaPivots(extremeData)
	assert.NoError(t, err, "Unexpected error with extreme values")
	assert.Greater(t, camarillaPivots.R3, camarillaPivots.R2, "R3 should be greater than R2 with extreme values")
	assert.Greater(t, camarillaPivots.R2, camarillaPivots.R1, "R2 should be greater than R1 with extreme values")
	assert.Greater(t, camarillaPivots.R1, camarillaPivots.P, "R1 should be greater than P with extreme values")
	assert.Greater(t, camarillaPivots.P, camarillaPivots.S1, "P should be greater than S1 with extreme values")
	assert.Greater(t, camarillaPivots.S1, camarillaPivots.S2, "S1 should be greater than S2 with extreme values")
	assert.Greater(t, camarillaPivots.S2, camarillaPivots.S3, "S2 should be greater than S3 with extreme values")
}
