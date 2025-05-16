package models

type TopGainerLoserResponse struct {
	Losers  []LosersGainers `json:"losers"`
	Gainers []LosersGainers `json:"gainers"`
}

type LosersGainers struct {
	TurnoverValue   float64 `json:"turnoverValue"`
	TradedQuantity  int     `json:"tradedQuantity"`
	Symbol          string  `json:"symbol"`
	NetPrice        float64 `json:"netPrice"`
	Ltp             float64 `json:"ltp"`
	LotSize         int     `json:"lotSize"`
	InstrumentToken int     `json:"instrumentToken"`
	Exchange        string  `json:"exchange"`
	CompanyName     string  `json:"companyName"`
	ClosePrice      float64 `json:"closePrice"`
}

type GainersLosersMostActiveVolumeReq struct {
	Index string `json:"index" example:"nifty_50"`
}

type MostActiveVolumeData struct {
	MostActiveVolume []MostActiveVolume `json:"mostActiveVolume"`
}

type MostActiveVolume struct {
	TurnoverValue     float64 `json:"turnoverValue"`
	TradedQuantity    int     `json:"tradedQuantity"`
	TotalSellQuantity int     `json:"totalSellQuantity"`
	TotalBuyQuantity  int     `json:"totalBuyQuantity"`
	Symbol            string  `json:"symbol"`
	PreviousPrice     float64 `json:"previousPrice"`
	NetPrice          float64 `json:"netPrice"`
	Ltp               float64 `json:"ltp"`
	LotSize           int     `json:"lotSize"`
	InstrumentToken   int     `json:"instrumentToken"`
	Exchange          string  `json:"exchange"`
	CompanyName       string  `json:"companyName"`
}

type ChartDataReq struct {
	Exchange     string `json:"exchange" example:"NSE"`
	Token        string `json:"token" example:"11536"`
	CandleType   string `json:"candleType" example:"1"`
	StartTime    string `json:"startTime" example:"1668623400"`
	EndTime      string `json:"endTime" example:"1670484551"`
	DataDuration string `json:"dataDuration" example:"1"`
}

type ChartDataResponse struct {
	Data CandleData `json:"data"`
}

type CandleData struct {
	Candles [][]interface{} `json:"candles"`
}

type ReturnOnInvestmentReq struct {
	Days  int    `json:"days" validate:"gte=0"`
	Index string `json:"index"`
}

type ReturnOnInvestmentRes struct {
	Roi []RoiData `json:"roi"`
}

type RoiData struct {
	Volume           int     `json:"volume"`
	ReturnPercent    float64 `json:"returnPercent"`
	PercentageChange float64 `json:"percentageChange"`
	Ltp              float64 `json:"ltp"`
	InstrumentToken  int     `json:"instrumentToken"`
	Exchange         string  `json:"exchange"`
	DaysChange       int     `json:"daysChange"`
	ClosePrice       float64 `json:"closePrice"`
	Change           float64 `json:"change"`
	TradingSymbol    string  `json:"tradingSymbol"`
}

type HistoricPerformaceReq struct {
	Exchange string `json:"exchange" example:"NSE"`
	Token    string `json:"token" example:"11536"`
	Period   string `json:"period" example:"1Y"`
}

type HistoricPerformaceRes struct {
	OpenPrice  float64 `json:"openPrice" example:"11536"`
	ClosePrice float64 `json:"closePrice" example:"115360"`
	ROI        string  `json:"roi" example:"10 %"`
}

type AllHistoricPerformaceReq struct {
	Exchange string `json:"exchange" example:"NSE"`
	Token    string `json:"token" example:"11536"`
}

type AllPerformanceRes struct {
	Period1D string `json:"1D"`
	Period1W string `json:"1W"`
	Period1M string `json:"1M"`
	Period6M string `json:"6M"`
	Period1Y string `json:"1Y"`
	Period3Y string `json:"3Y"`
	Period5Y string `json:"5Y"`
}
