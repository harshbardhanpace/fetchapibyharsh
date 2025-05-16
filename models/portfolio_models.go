package models

// Fetch Demat Holdings Request
type FetchDematHoldingsRequest struct {
	ClientID string `json:"clientId" binding:"required" example:"Client1"`
}

// Fetch Demat Holdings Response
type FetchDematHoldingsResponse struct {
	Holdings []FetchDematHoldingsResponseData `json:"holdings"`
}

type FetchDematHoldingsResponseData struct {
	BranchCode        string  `json:"branchCode"`
	BuyAvg            float64 `json:"buyAvg"`
	BuyAvgMtm         float64 `json:"buyAvgMtm"`
	ClientID          string  `json:"clientId"`
	Exchange          string  `json:"exchange"`
	FreeQuantity      int     `json:"freeQuantity"`
	InstrumentDetails struct {
		Exchange        int    `json:"exchange"`
		InstrumentName  string `json:"instrumentName"`
		InstrumentToken int    `json:"instrumentToken"`
		TradingSymbol   string `json:"tradingSymbol"`
	} `json:"instrumentDetails"`
	Isin                  string         `json:"isin"`
	Ltp                   float64        `json:"ltp"`
	PendingQuantity       int            `json:"pendingQuantity"`
	PledgeAllow           bool           `json:"pledgeAllow"`
	PledgeQuantity        int            `json:"pledgeQuantity"`
	PreviousClose         float64        `json:"previousClose"`
	Quantity              int            `json:"quantity"`
	Symbol                string         `json:"symbol"`
	T0Price               float64        `json:"t0Price"`
	T0Quantity            int            `json:"t0Quantity"`
	T1Price               float64        `json:"t1Price"`
	T1Quantity            int            `json:"t1Quantity"`
	T2Price               float64        `json:"t2Price"`
	T2Quantity            int            `json:"t2Quantity"`
	TodayPledgeQuantity   int            `json:"todayPledgeQuantity"`
	TodayUnpledgeQuantity int            `json:"todayUnpledgeQuantity"`
	Token                 int            `json:"token"`
	TradingSymbol         string         `json:"tradingSymbol"`
	TransactionType       string         `json:"transactionType"`
	UsedQuantity          int            `json:"usedQuantity"`
	ActualBuyAvg          float64        `json:"actualBuyAvg"`
	NetHoldingQty         int            `json:"netHoldingsQty"`
	PledgePercentage      float64        `json:"pledgePercentage"`
	AdditionalInfo        AdditionalInfo `json:"additionalInfo"`
}

// Convert Positions Request
type ConvertPositionsRequest struct {
	ClientID        string `json:"clientId" binding:"required" example:"Client1"`
	Exchange        string `json:"exchange" enums:"NSE,BSE,NFO,CDS,MCX,BFO" example:"NSE" validate:"oneof=NSE BSE NFO CDS MCX BFO"`
	InstrumentToken int    `json:"instrumentToken"`
	Product         string `json:"product" enums:"CNC,MIS,NRML" example:"CNC" validate:"oneof=CNC MIS NRML"`
	NewProduct      string `json:"newProduct" enums:"CNC,MIS,NRML" example:"CNC" validate:"oneof=CNC MIS NRML"`
	Quantity        int    `json:"quantity" validate:"gt=0"`
	Validity        string `json:"validity" enums:"DAY,IOC" example:"DAY" validate:"oneof=DAY IOC"`
	OrderSide       string `json:"orderSide" enums:"BUY,SELL" example:"BUY" validate:"oneof=BUY SELL"`
}

// Convert Positions Resposne
type ConvertPositionsResponse struct {
	Data struct {
	} `json:"data"`
}

// GetPositionRequest
type GetPositionRequest struct {
	ClientID string `json:"clientId"`
	Type     string `json:"type" example:"live,historical" validate:"oneof=live historical"`
}

// GetPositionResponse
type GetPositionResponse struct {
	Data []GetPositionResponseData `json:"data"`
}

type GetPositionResponseData struct {
	ActualCfBuyAmount      float64        `json:"actualCfBuyAmount"`
	ActualCfSellAmount     float64        `json:"actualCfSellAmount"`
	ActualAverageBuyPrice  float64        `json:"actualAverageBuyPrice"`
	ActualAverageSellPrice float64        `json:"actualAverageSellPrice"`
	AverageBuyPrice        float64        `json:"averageBuyPrice"`
	AveragePrice           float64        `json:"averagePrice"`
	AverageSellPrice       float64        `json:"averageSellPrice"`
	BuyAmount              float64        `json:"buyAmount"`
	BuyQuantity            int            `json:"buyQuantity"`
	CfBuyAmount            float64        `json:"cfBuyAmount"`
	CfBuyQuantity          int            `json:"cfBuyQuantity"`
	CfSellAmount           float64        `json:"cfSellAmount"`
	CfSellQuantity         int            `json:"cfSellQuantity"`
	ClientID               string         `json:"clientId"`
	ClosePrice             float64        `json:"closePrice"`
	Exchange               string         `json:"exchange"`
	InstrumentToken        int            `json:"instrumentToken"`
	Ltp                    float64        `json:"ltp"`
	Multiplier             float64        `json:"multiplier"`
	NetAmount              float64        `json:"netAmount"`
	NetQuantity            int            `json:"netQuantity"`
	OtherMargin            any            `json:"otherMargin"`
	PreviousClose          float64        `json:"previousClose"`
	ProCli                 string         `json:"proCli"`
	ProdType               int            `json:"prodType"`
	Product                string         `json:"product"`
	RealizedMtm            float64        `json:"realizedMtm"`
	Segment                string         `json:"segment"`
	SellAmount             float64        `json:"sellAmount"`
	SellQuantity           int            `json:"sellQuantity"`
	Symbol                 string         `json:"symbol"`
	Token                  int            `json:"token"`
	TotalPledgeCollateral  any            `json:"totalPledgeCollateral"`
	TradingSymbol          string         `json:"tradingSymbol"`
	VLoginID               string         `json:"vLoginId"`
	AdditionalInfo         AdditionalInfo `json:"additionalInfo"`
}

type AdditionalInfo struct {
	Expiry           int    `json:"expiry"`
	Isin             string `json:"isin"`
	Type             string `json:"type"` //option or future
	StrikePrice      string `json:"strikePrice"`
	IsWeekly         bool   `json:"isWeekly"`
	UnderlyingSymbol string `json:"underlyingSymbol"`
}
