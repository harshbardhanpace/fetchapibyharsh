package models

type GetSMAReq struct {
	TLChartData TLChartDataReq `json:"tlChartData"`
	SMAType     int            `json:"smaType" example:"5"`
}

type GetEMAReq struct {
	TLChartData TLChartDataReq `json:"tlChartData"`
	EMAType     int            `json:"emaType" example:"5"`
}

type GetHullMAReq struct {
	TLChartData TLChartDataReq `json:"tlChartData"`
	HullMAType  int            `json:"hullMAType" example:"9"`
}

type GetVWMAReq struct {
	TLChartData TLChartDataReq `json:"tlChartData"`
	VWMAType    int            `json:"vwmaType" example:"20"`
}

type GetRSIReq struct {
	TLChartData TLChartDataReq `json:"tlChartData"`
	RSIType     int            `json:"rsiType" example:"14"`
}

type GetCCIReq struct {
	TLChartData TLChartDataReq `json:"tlChartData"`
	CCIType     int            `json:"cciType" example:"20"`
}

type GetMACDReq struct {
	TLChartData  TLChartDataReq `json:"tlChartData"`
	FastPeriod   int            `json:"fastPeriod" example:"12"`
	SlowPeriod   int            `json:"slowPeriod" example:"26"`
	SignalPeriod int            `json:"signalPeriod" example:"8"`
}

type GetStochasticReq struct {
	TLChartData TLChartDataReq `json:"tlChartData"`
	KPeriod     int            `json:"kPeriod" example:"14"`
	DPeriod     int            `json:"dPeriod" example:"3"`
	Smooth      int            `json:"smooth" example:"3"`
}

type GetIchimokuBaseLineReq struct {
	TLChartData TLChartDataReq `json:"tlChartData"`
}

type GetADXReq struct {
	TLChartData TLChartDataReq `json:"tlChartData"`
	Period      int            `json:"period" example:"14"`
}

type GetAwesomeOscillatorReq struct {
	TLChartData TLChartDataReq `json:"tlChartData"`
}

type GetMomentumReq struct {
	TLChartData TLChartDataReq `json:"tlChartData"`
	Period      int            `json:"period" example:"10"`
}

type GetStochRSIFastReq struct {
	TLChartData TLChartDataReq `json:"tlChartData"`
	SmoothK     int            `json:"smoothK" example:"3"`
	SmoothD     int            `json:"smoothD" example:"3"`
	RsiPeriod   int            `json:"rsiPeriod" example:"14"`
	StochPeriod int            `json:"stochPeriod" example:"14"`
}

type GetWilliamsRangeReq struct {
	TLChartData TLChartDataReq `json:"tlChartData"`
	Period      int            `json:"period" example:"14"`
}

type GetUltimateOscillatorReq struct {
	TLChartData TLChartDataReq `json:"tlChartData"`
	Period1     int            `json:"period1" example:"7"`
	Period2     int            `json:"period2" example:"14"`
	Period3     int            `json:"period3" example:"28"`
}

type TechnicalIndicatorsRes struct {
	TimestampUnix string  `json:"timestampUnix"`
	Value         float64 `json:"value"`
	Signal        string  `json:"signal"`
}

type TechnicalIndicatorsResFull struct {
	Type string                   `json:"type"`
	Data []TechnicalIndicatorsRes `json:"data"`
}

type TLChartDataReq struct {
	Exchange     string `json:"exchange" example:"NSE"`
	Token        string `json:"token" example:"11536"`
	CandleType   string `json:"candleType" example:"1"`
	StartTime    string `json:"startTime" example:"1668623400"`
	EndTime      string `json:"endTime" example:"1670484551"`
	DataDuration string `json:"dataDuration" example:"1"`
}

type AllTechnicalIndicatorsRes struct {
	Entries []TechnicalIndicatorsResFull `json:"entries"`
}

type GetAllTechnicalIndicatorsReq struct {
	Exchange     string `json:"exchange" example:"NSE"`
	Token        string `json:"token" example:"2885"`
	TimeUnit     string `json:"timeUnit" example:"MONTH, WEEK, DAY, HOUR, MINUTE" validate:"oneof=DAY HOUR MINUTE WEEK MONTH"`
	TimeInterval int    `json:"timeInterval" example:1, 5, 10, 30`
}

type GetAllTechnicalIndicatorsRes struct {
	Entries []GetAllTechnicalIndicatorsResInternal `json:"entries"`
	Pivots  []PivotsValues                         `json:"pivots"`
}

type GetAllTechnicalIndicatorsResInternal struct {
	Type          string  `json:"type"`
	TimestampUnix string  `json:"timestampUnix"`
	Value         float64 `json:"value"`
	Signal        string  `json:"signal"`
}

type PivotPoints struct {
	P  float64 `json:"p"`  // Pivot Point
	R1 float64 `json:"r1"` // Resistance 1
	R2 float64 `json:"r2"` // Resistance 2
	R3 float64 `json:"r3"` // Resistance 3
	S1 float64 `json:"s1"` // Support 1
	S2 float64 `json:"s2"` // Support 2
	S3 float64 `json:"s3"` // Support 3
}

type PivotsValues struct {
	Type          string      `json:"type"`
	TimestampUnix string      `json:"timestampUnix"`
	Points        PivotPoints `json:"points"`
}
