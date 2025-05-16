package models

type BrokerChargesReq struct {
	ClientID   string  `json:"clientId" validate:"required"`
	Price      float64 `json:"price" validate:"gt=0"`    // price of stock
	Quantity   int     `json:"quantity" validate:"gt=0"` // quantity of stock
	Segment    string  `json:"segment" validate:"required"`
	SubSegment string  `json:"subSegment" validate:"required"`
	Process    string  `json:"process" validate:"oneof=BUY SELL"`     // buy or sell
	Exchange   string  `json:"exchange" validate:"oneof=NSE BSE MCX"` // nse bse mcx
	Agri       bool    `json:"agriType"`                              // stock type
	GroupInfo  int     `json:"groupInfo"`                             // price calculated differently
	Product    string  `json:"product"`
}

type BrokerChargesRes struct {
	Price              float64 `json:"price"`
	Brokerage          float64 `json:"brokarage"`
	TotalCharge        float64 `json:"totalCharge"`
	SttOrCtt           float64 `json:"sttOrCtt"`
	TransactionCharges float64 `json:"transactionCharges"`
	SebiCharges        float64 `json:"sebiCharges"`
	Gst                float64 `json:"gst"`
	StampCharges       float64 `json:"stampCharges"`
}

type CombineBrokerChargesReq struct {
	ClientID      string             `json:"clientId"`
	BrokerCharges []BrokerChargesReq `json:"brokerCharges" validate:"dive"`
}

type CombineBrokerChargesRes struct {
	BrokerCharges []BrokerChargesRes `json:"brokerCharges"`
}

type FundsPayoutReq struct {
	ClientID string `json:"clientId"`
}

type FundsPayoutRes struct {
	ClientID              string  `json:"clientId"`
	PayoutAmount          float64 `json:"payoutAmount"`
	OpeningBalance        float64 `json:"openingBalance"`
	MarginUsed            float64 `json:"marginUsed"`
	LossOnClosedPositions float64 `json:"lossOnClosedPositions"`
	ChargesOnTrades       float64 `json:"chargesOnTrades"`
	Payin                 float64 `json:"payin"`
	ExtraPayoutAmount     float64 `json:"extraPayoutAmount"`
	ProfitNLoss           float64 `json:"pNlOnClosedPosition"`
	CmPnl                 float64 `json:"cmPnl"`
	UserPnl               float64 `json:"userPnl"`
	EquityCreditSell      float64 `json:"equityCreditSell"`
}
