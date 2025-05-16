package models

type TLChartCandleData struct {
	Timestamp string
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

type TechnicalIndicatorsValuesReq struct {
	Exchange string `json:"exchange"`
	Token    string `json:"token"`
}

type TechnicalIndicatorsValuesRes struct {
	SMA               float64 `json:"sma"`
	EMA               float64 `json:"ema"`
	RSI               float64 `json:"rsi"`
	MACD              float64 `json:"macd"`
	MACDSignal        float64 `json:"macdSignal"`
	CCI               float64 `json:"cci"`
	AwesomeOscillator float64 `json:"awesomeOscillator"`
	S1                float64 `json:"s1"`
	S2                float64 `json:"s2"`
	S3                float64 `json:"s3"`
	R1                float64 `json:"r1"`
	R2                float64 `json:"r2"`
	R3                float64 `json:"r3"`
}
