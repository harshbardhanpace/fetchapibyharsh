package models

//PlaceOrderRequest place order request
type OrderReq struct {
	ClientID      string      `json:"clientId" binding:"required" example:"abc123"`
	OrderType     string      `json:"orderType" binding:"required"`
	RequestPacket interface{} `json:"requestPacket" binding:"required"`
}

type PlaceOrderRequest struct {
	Exchange          string  `json:"exchange" enums:"NSE,BSE,NFO,CDS,MCX,BFO" example:"NSE" validate:"oneof=NSE BSE NFO CDS MCX BFO"`
	InstrumentToken   string  `json:"instrumentToken" example:"Represents the unique id of instrument."`
	ClientID          string  `json:"clientId" binding:"required" example:"abc123"`
	OrderType         string  `json:"orderType" enums:"LIMIT,MARKET,SL,SLM" example:"LIMIT" validate:"oneof=LIMIT MARKET SL SLM"`
	Price             float64 `json:"price" validate:"gte=0"`
	Quantity          int     `json:"quantity" validate:"gt=0"`
	DisclosedQuantity int     `json:"disclosedQuantity"`
	Validity          string  `json:"validity" enums:"DAY,IOC" example:"DAY" validate:"oneof=DAY IOC"`
	Product           string  `json:"product" enums:"CNC,MIS,NRML,MTF" example:"CNC" validate:"oneof=CNC MIS NRML MTF"`
	OrderSide         string  `json:"orderSide" enums:"BUY,SELL" example:"BUY" validate:"oneof=BUY SELL"`
	TriggerPrice      float64 `json:"triggerPrice"`
	ExecutionType     string  `json:"executionType" example:"REGULAR" validate:"oneof=REGULAR AMO"`
	NoOfLegs          int     `json:"noOfLegs"`
	Device            string  `json:"device"`
}

type PlaceOrderMTFRequest struct {
	Exchange          string  `json:"exchange" enums:"NSE,BSE,NFO,CDS,MCX,BFO" example:"NSE" validate:"oneof=NSE BSE NFO CDS MCX BFO"`
	InstrumentToken   string  `json:"instrumentToken" example:"Represents the unique id of instrument."`
	ClientID          string  `json:"clientId" binding:"required" example:"abc123"`
	OrderType         string  `json:"orderType" enums:"LIMIT,MARKET,SL,SLM" example:"LIMIT" validate:"oneof=LIMIT MARKET SL SLM"`
	Price             string  `json:"price"`
	Quantity          int     `json:"quantity" validate:"gt=0"`
	DisclosedQuantity int     `json:"disclosedQuantity"`
	Validity          string  `json:"validity" enums:"DAY,IOC" example:"DAY" validate:"oneof=DAY IOC"`
	Product           string  `json:"product" enums:"MTF" example:"MTF" validate:"oneof=MTF"`
	NoOfLegs          int     `json:"noOfLegs"`
	OrderSide         string  `json:"orderSide" enums:"BUY,SELL" example:"BUY" validate:"oneof=BUY SELL"`
	Device            string  `json:"device"`
	UserOrderId       int     `json:"userOrderId"`
	TriggerPrice      float64 `json:"triggerPrice"`
}

type MTFEPledgeRequest struct {
	Depository  string           `json:"depository"`
	ClientID    string           `json:"clientId"`
	Exchange    string           `json:"exchange"`
	BoID        string           `json:"boId"`
	Segment     string           `json:"segment"`
	RequestType string           `json:"requestType"`
	IsinDetails []MTFIsinDetails `json:"isinDetails"`
	Order       MTFOrder         `json:"order"`
}

type MTFIsinDetails struct {
	IsinName string `json:"isinName"`
	Isin     string `json:"isin"`
	Quantity string `json:"quantity"`
	Price    string `json:"price"`
}

type MTFOrder struct {
	Price             string `json:"price"`
	Device            string `json:"device"`
	Product           string `json:"product"`
	Exchange          string `json:"exchange"`
	Quantity          int    `json:"quantity"`
	Validity          string `json:"validity"`
	ClientID          string `json:"clientId"`
	OrderSide         string `json:"orderSide"`
	OrderType         string `json:"orderType"`
	UserOrderID       string `json:"userOrderId"`
	InstrumentToken   string `json:"instrumentToken"`
	DisclosedQuantity int    `json:"disclosedQuantity"`
}

type MTFPledgeListResponse struct {
	List       []MTDPledgeList `json:"list"`
	TotalCount int             `json:"totalCount"`
}

type MTDPledgeList struct {
	ClientID             string  `json:"clientId"`
	Isin                 string  `json:"isin"`
	TotalPledgeQuantity  int     `json:"totalPledgeQuantity"`
	CtdQuantity          int     `json:"ctdQuantity"`
	Symbol               string  `json:"symbol"`
	AvgPrice             float64 `json:"avgPrice"`
	MarginMultiplier     int     `json:"marginMultiplier"`
	CtdMarginValue       int     `json:"ctdMarginValue"`
	Token                int     `json:"token"`
	Exchange             string  `json:"exchange"`
	CreatedAt            string  `json:"createdAt"`
	UpdatedAt            string  `json:"updatedAt"`
	EdisApprovedQuantity int     `json:"edisApprovedQuantity"`
	ObligationQuantity   int     `json:"obligationQuantity"`
	UsedQuantity         int     `json:"usedQuantity"`
	LoginID              string  `json:"loginId"`
	MarginValue          float64 `json:"marginValue"`
	TotalInvestedAmount  float64 `json:"totalInvestedAmount"`
	BrokerAmount         float64 `json:"brokerAmount"`
}

//PlaceOrderResponse place order response
type PlaceOrderResponse struct {
	OmsOrderID  string `json:"omsOrderId" example:"adsad123"`
	UserOrderID int    `json:"userOrderId" example:"21441523"`
}

//ModifyOrderRequest modify order request
type ModifyOrderRequest struct {
	Exchange              string  `json:"exchange" enums:"NSE,BSE,NFO,CDS,MCX,BFO" example:"NSE" validate:"oneof=NSE BSE NFO CDS MCX BFO"`
	InstrumentToken       string  `json:"instrumentToken" example:"Represents the unique id of instrument."`
	ClientID              string  `json:"clientId" binding:"required" example:"abc123"`
	OrderType             string  `json:"orderType" enums:"LIMIT,MARKET,SL,SLM" example:"LIMIT" validate:"oneof=LIMIT MARKET SL SLM"`
	Price                 float64 `json:"price" validate:"gte=0"`
	Quantity              int     `json:"quantity" validate:"gt=0"`
	DisclosedQuantity     int     `json:"disclosedQuantity"`
	Validity              string  `json:"validity" enums:"DAY,IOC" example:"DAY" validate:"oneof=DAY IOC"`
	Product               string  `json:"product" enums:"CNC,MIS,NRML,MTF" example:"CNC" validate:"oneof=CNC MIS NRML MTF"`
	OmsOrderID            string  `json:"omsOrderId" example:"adsad123"`
	TriggerPrice          float64 `json:"triggerPrice"`
	ExecutionType         string  `json:"executionType" example:"REGULAR" validate:"oneof=REGULAR"`
	ExchangeOrderID       string  `json:"exchangeOrderID"`
	FilledQuantity        int     `json:"filledQuantity"`
	RemainingQuantity     int     `json:"remainingQuantity"`
	LastActivityReference int64   `json:"lastActivityReference"`
}

//CancelOrderReqeust cancel order request
type CancelOrderRequest struct {
	ClientID      string `json:"clientId"`
	OmsOrderId    string `json:"omsOrderId"`
	ExecutionType string `json:"executionType" example:"REGULAR" validate:"oneof=REGULAR AMO"`
}

//ModifyOrCancelOrderResponse modify or cancel order response
type ModifyOrCancelOrderResponse struct {
	OmsOrderID string `json:"omsOrderId"`
}

//ModifyAMORequest modify AMO request
type ModifyAMORequest struct {
	Exchange              string  `json:"exchange" enums:"NSE,BSE,NFO,CDS,MCX,BFO" example:"NSE" validate:"oneof=NSE BSE NFO CDS MCX BFO"`
	InstrumentToken       string  `json:"instrumentToken" example:"Represents the unique id of instrument."`
	ClientID              string  `json:"clientId" binding:"required" example:"abc123"`
	OrderType             string  `json:"orderType" enums:"LIMIT,MARKET,SL,SLM" example:"LIMIT" validate:"oneof=LIMIT MARKET SL SLM"`
	Price                 float64 `json:"price" validate:"gte=0"`
	Quantity              int     `json:"quantity" validate:"gt=0"`
	DisclosedQuantity     int     `json:"disclosedQuantity"`
	Validity              string  `json:"validity" enums:"DAY,IOC" example:"DAY" validate:"oneof=DAY IOC"`
	Product               string  `json:"product" enums:"CNC,MIS,NRML" example:"CNC" validate:"oneof=CNC MIS NRML"`
	OmsOrderID            string  `json:"omsOrderId" example:"adsad123"`
	ExchangeOrderID       string  `json:"exchangeOrderId"`
	FilledQuantity        int     `json:"filledQuantity"`
	RemainingQuantity     int     `json:"remainingQuantity"`
	LastActivityReference int     `json:"lastActivityReference"`
	TriggerPrice          float64 `json:"triggerPrice"`
	ExecutionType         string  `json:"executionType" example:"AMO" validate:"oneof=AMO"`
}

//AMOOrderResponse AMO response
type AMOOrderResponse struct {
	OmsOrderID string `json:"omsOrderId"`
}

type PendingOrderRequest struct {
	Type     string `json:"type" example:"pending" validate:"oneof=pending"`
	ClientID string `json:"clientId" binding:"required" example:"abc123"`
}

//FetchingPendingOrder fetching pending order response
type PendingOrderResponse struct {
	Orders []PendingOrderResponseOrders `json:"orders"`
}

type PendingOrderResponseOrders struct {
	TradingSymbol              string  `json:"tradingSymbol"`
	AverageTradePrice          float64 `json:"averageTradePrice"`
	Exchange                   string  `json:"exchange"`
	ProCli                     string  `json:"proCli"`
	MarketProtectionPercentage int     `json:"marketProtectionPercentage"`
	OrderEntryTime             int     `json:"orderEntryTime"`
	Mode                       string  `json:"mode"`
	OmsOrderID                 string  `json:"omsOrderId"`
	TrailingStopLoss           float64 `json:"trailingStopLoss"`
	Deposit                    int     `json:"deposit"`
	SquareOffValue             float64 `json:"squareOffValue"`
	DisclosedQuantity          int     `json:"disclosedQuantity"`
	StopLossValue              float64 `json:"stopLossValue"`
	Price                      float64 `json:"price"`
	OrderTag                   string  `json:"orderTag"`
	Device                     string  `json:"device"`
	RemainingQuantity          int     `json:"remainingQuantity"`
	LastActivityReference      int64   `json:"lastActivityReference"`
	AveragePrice               float64 `json:"averagePrice"`
	SquareOff                  bool    `json:"squareOff"`
	OrderStatusInfo            string  `json:"orderStatusInfo"`
	Quantity                   int     `json:"quantity"`
	ExecutionType              string  `json:"executionType"`
	ClientID                   string  `json:"clientId"`
	ExchangeTime               int     `json:"exchangeTime"`
	OrderSide                  string  `json:"orderSide"`
	LoginID                    string  `json:"loginId"`
	Validity                   string  `json:"validity"`
	InstrumentToken            int     `json:"instrumentToken"`
	Product                    string  `json:"product"`
	TriggerPrice               float64 `json:"triggerPrice"`
	Segment                    string  `json:"segment"`
	TradePrice                 float64 `json:"tradePrice"`
	OrderType                  string  `json:"orderType"`
	//ContractDescription        struct {
	//} `json:"contractDescription"`
	RejectionCode     int            `json:"rejectionCode"`
	LegOrderIndicator string         `json:"legOrderIndicator"`
	ExchangeOrderID   string         `json:"exchangeOrderId"`
	OrderStatus       string         `json:"orderStatus"`
	FilledQuantity    int            `json:"filledQuantity"`
	TargetPriceType   string         `json:"targetPriceType"`
	IsTrailing        bool           `json:"isTrailing"`
	UserOrderID       string         `json:"userOrderId"`
	LotSize           int            `json:"lotSize"`
	Series            string         `json:"series"`
	NnfID             int64          `json:"nnfId"`
	RejectionReason   string         `json:"rejectionReason"`
	AdditionalInfo    AdditionalInfo `json:"additionalInfo"`
}

type CompletedOrderRequest struct {
	Type     string `json:"type" example:"completed" validate:"oneof=completed"`
	ClientID string `json:"clientId" binding:"required" example:"abc123"`
}

type CompletedOrderResponse struct {
	Orders []CompletedOrderResponseOrders `json:"orders"`
}

type CompletedOrderResponseOrders struct {
	TradingSymbol              string  `json:"tradingSymbol"`
	AverageTradePrice          float64 `json:"averageTradePrice"`
	Exchange                   string  `json:"exchange"`
	ProCli                     string  `json:"proCli"`
	MarketProtectionPercentage int     `json:"marketProtectionPercentage"`
	OrderEntryTime             int     `json:"orderEntryTime"`
	Mode                       string  `json:"mode"`
	OmsOrderID                 string  `json:"omsOrderId"`
	TrailingStopLoss           float64 `json:"trailingStopLoss"`
	Deposit                    int     `json:"deposit"`
	SquareOffValue             float64 `json:"squareOffValue"`
	DisclosedQuantity          int     `json:"disclosedQuantity"`
	StopLossValue              float64 `json:"stopLossValue"`
	Price                      float64 `json:"price"`
	OrderTag                   string  `json:"orderTag"`
	Device                     string  `json:"device"`
	RemainingQuantity          int     `json:"remainingQuantity"`
	LastActivityReference      int     `json:"lastActivityReference"`
	AveragePrice               float64 `json:"averagePrice"`
	SquareOff                  bool    `json:"squareOff"`
	OrderStatusInfo            string  `json:"orderStatusInfo"`
	Quantity                   int     `json:"quantity"`
	ExecutionType              string  `json:"executionType"`
	ClientID                   string  `json:"clientId"`
	ExchangeTime               int     `json:"exchangeTime"`
	OrderSide                  string  `json:"orderSide"`
	LoginID                    string  `json:"loginId"`
	Validity                   string  `json:"validity"`
	InstrumentToken            int     `json:"instrumentToken"`
	Product                    string  `json:"product"`
	TriggerPrice               float64 `json:"triggerPrice"`
	Segment                    string  `json:"segment"`
	TradePrice                 float64 `json:"tradePrice"`
	OrderType                  string  `json:"orderType"`
	//ContractDescription   struct {
	//} `json:"contractDescription"`
	RejectionCode     int            `json:"rejectionCode"`
	LegOrderIndicator string         `json:"legOrderIndicator"`
	ExchangeOrderID   string         `json:"exchangeOrderId"`
	OrderStatus       string         `json:"orderStatus"`
	FilledQuantity    int            `json:"filledQuantity"`
	TargetPriceType   string         `json:"targetPriceType"`
	IsTrailing        bool           `json:"isTrailing"`
	UserOrderID       string         `json:"userOrderId"`
	LotSize           int            `json:"lotSize"`
	Series            string         `json:"series"`
	NnfID             int64          `json:"nnfId"`
	RejectionReason   string         `json:"rejectionReason"`
	AdditionalInfo    AdditionalInfo `json:"additionalInfo"`
}

//TradeBookRequest
type TradeBookRequest struct {
	ClientID string `json:"clientId" binding:"required" example:"abc123"`
}

//TradeBookResponse
type TradeBookResponse struct {
	Trades []TradeBookResponseData `json:"trades" mask:"struct"`
}

type TradeBookResponseData struct {
	BookType              string         `json:"bookType"`
	BrokerID              string         `json:"brokerId"`
	ClientID              string         `json:"clientId"`
	DisclosedVol          int            `json:"disclosedVol"`
	DisclosedVolRemaining int            `json:"disclosedVolRemaining"`
	Exchange              string         `json:"exchange"`
	ExchangeOrderID       string         `json:"exchangeOrderId"`
	ExchangeTime          int            `json:"exchangeTime"`
	FillNumber            string         `json:"fillNumber"`
	FilledQuantity        int            `json:"filledQuantity"`
	GoodTillDate          int            `json:"goodTillDate"`
	InstrumentToken       int            `json:"instrumentToken"`
	LoginID               string         `json:"loginId"`
	OmsOrderID            string         `json:"omsOrderId"`
	OrderEntryTime        int            `json:"orderEntryTime"`
	OrderPrice            float64        `json:"orderPrice"`
	OrderSide             string         `json:"orderSide"`
	OrderType             string         `json:"orderType"`
	OriginalVol           int            `json:"originalVol"`
	Pan                   string         `json:"pan" mask:"id"`
	ProCli                int            `json:"proCli"`
	Product               string         `json:"product"`
	RemainingQuantity     int            `json:"remainingQuantity"`
	TradeNumber           string         `json:"tradeNumber"`
	TradePrice            float64        `json:"tradePrice"`
	TradeQuantity         int            `json:"tradeQuantity"`
	TradeTime             int            `json:"tradeTime"`
	TradingSymbol         string         `json:"tradingSymbol"`
	TriggerPrice          float64        `json:"triggerPrice"`
	VLoginID              string         `json:"vLoginId"`
	VolFilledToday        int            `json:"volFilledToday"`
	AdditionalInfo        AdditionalInfo `json:"additionalInfo"`
}

//OrderHistoryRequest
type OrderHistoryRequest struct {
	OmsOrderID string `json:"omsOrderId" binding:"required" example:"20220920-4"`
	ClientID   string `json:"clientId" binding:"required" example:"abc123"`
}

//OrderHistoryResponse
type OrderHistoryResponse struct {
	OrderHistory []OrderHistoryResponseData `json:"orderHistory"`
}

type OrderHistoryResponseData struct {
	AvgPrice          float64 `json:"avgPrice"`
	ClientID          string  `json:"clientId"`
	ClientOrderID     string  `json:"clientOrderId"`
	CreatedAt         int     `json:"createdAt"`
	DisclosedQuantity int     `json:"disclosedQuantity"`
	Exchange          string  `json:"exchange"`
	ExchangeOrderID   string  `json:"exchangeOrderId"`
	ExchangeTime      int     `json:"exchangeTime"`
	FillQuantity      int     `json:"fillQuantity"`
	LastModified      int64   `json:"lastModified"`
	LoginID           string  `json:"loginId"`
	ModifiedAt        int     `json:"modifiedAt"`
	OrderID           string  `json:"orderId"`
	OrderMode         string  `json:"orderMode"`
	OrderSide         string  `json:"orderSide"`
	OrderType         string  `json:"orderType"`
	Price             float64 `json:"price"`
	Product           string  `json:"product"`
	Quantity          int     `json:"quantity"`
	RejectReason      string  `json:"rejectReason"`
	RemainingQuantity int     `json:"remainingQuantity"`
	Segment           string  `json:"segment"`
	Status            string  `json:"status"`
	Symbol            string  `json:"symbol"`
	Token             int     `json:"token"`
	TriggerPrice      float64 `json:"triggerPrice"`
	UnderlyingToken   int     `json:"underlyingToken"`
	Validity          string  `json:"validity"`
}

type CreateGTTOrderRequest struct {
	ActionType                 string  `json:"actionType" example:"single_order" validate:"oneof=single_order oco_order oco"`
	ExpiryTime                 string  `json:"expiryTime"`
	ClientID                   string  `json:"clientId"`
	DisclosedQuantity          int     `json:"disclosedQuantity"`
	Exchange                   string  `json:"exchange" enums:"NSE,BSE,NFO,CDS,MCX,BFO" example:"NSE" validate:"oneof=NSE BSE NFO CDS MCX BFO"`
	InstrumentToken            string  `json:"instrumentToken"`
	MarketProtectionPercentage int     `json:"marketProtectionPercentage"`
	OrderSide                  string  `json:"orderSide"`
	OrderType                  string  `json:"orderType"`
	Price                      float64 `json:"price" validate:"gte=0"`
	Product                    string  `json:"product"`
	Quantity                   int     `json:"quantity" validate:"gt=0"`
	SlOrderPrice               float64 `json:"slOrderPrice"`
	SlOrderQuantity            int     `json:"slOrderQuantity"`
	SlTriggerPrice             float64 `json:"slTriggerPrice"`
	TriggerPrice               float64 `json:"triggerPrice"`
	UserOrderID                int     `json:"userOrderId"`
}

type GTTOrderResponse struct {
	ID string `json:"id"`
}

type ModifyGTTOrderRequest struct {
	ExpiryTime string                `json:"expiryTime"`
	ActionType string                `json:"actionType"`
	ID         string                `json:"id"`
	Order      ModifyGTTOrderDetails `json:"order"`
}

type ModifyGTTOrderDetails struct {
	ClientID string `json:"clientId"`
	// Device                     string  `json:"device" example:"WEB, ANDROID" validate:"oneof=WEB ANDROID"`
	DisclosedQuantity          int     `json:"disclosedQuantity"`
	Exchange                   string  `json:"exchange" enums:"NSE,BSE,NFO,CDS,MCX,BFO" example:"NSE" validate:"oneof=NSE BSE NFO CDS MCX BFO"`
	InstrumentToken            string  `json:"instrumentToken"`
	MarketProtectionPercentage int     `json:"marketProtectionPercentage"`
	OrderSide                  string  `json:"orderSide"`
	OrderType                  string  `json:"orderType"`
	Price                      float64 `json:"price" validate:"gte=0"`
	Product                    string  `json:"product"`
	Quantity                   int     `json:"quantity" validate:"gt=0"`
	SlOrderPrice               int     `json:"slOrderPrice"`
	SlOrderQuantity            int     `json:"slOrderQuantity"`
	SlTriggerPrice             int     `json:"slTriggerPrice"`
	TriggerPrice               float64 `json:"triggerPrice"`
	UserOrderID                int     `json:"userOrderId"`
}

type CancelGTTOrderRequest struct {
	ClientId string `json:"clientId"`
	Id       string `json:"id"`
}

type FetchGTTOrderRequest struct {
	ClientId string `json:"clientId"`
}

type FetchGTTOrderResponse struct {
	FetchGTTOrderData []FetchGTTOrderResponseData `json:"fetchGTTOrderData"`
}

type FetchGTTOrderResponseData struct {
	ActionType string `json:"actionType"`
	ClientID   string `json:"clientId"`
	CreatedAt  string `json:"createdAt"`
	ExpiryTime string `json:"expiryTime"`
	ID         string `json:"id"`
	LoginID    string `json:"loginId"`
	Order      struct {
		DisclosedQty     int     `json:"disclosedQty"`
		Exchange         string  `json:"exchange"`
		ExecutionType    string  `json:"executionType"`
		Mode             string  `json:"mode"`
		OrderSide        string  `json:"orderSide"`
		OrderType        string  `json:"orderType"`
		Price            float64 `json:"price"`
		ProCli           string  `json:"proCli"`
		ProdType         string  `json:"prodType"`
		Quantity         int     `json:"quantity"`
		Segment          string  `json:"segment"`
		SlOrderPrice     float64 `json:"slOrderPrice"`
		SlOrderQuantity  int     `json:"slOrderQuantity"`
		SlTriggerPrice   float64 `json:"slTriggerPrice"`
		SquareOffPrice   float64 `json:"squareOffPrice"`
		Token            int     `json:"token"`
		TradingSymbol    string  `json:"tradingSymbol"`
		TrailingStopLoss int     `json:"trailingStopLoss"`
		TriggerPrice     float64 `json:"triggerPrice"`
		Validity         string  `json:"validity"`
		VendorCode       string  `json:"vendorCode"`
	} `json:"order"`
	RejectCode   int    `json:"rejectCode"`
	RejectReason string `json:"rejectReason"`
	Status       string `json:"status"`
	Type         string `json:"type"`
	UpdatedAt    string `json:"updatedAt"`
}

type MarginCalculationRequest struct {
	Data []MarginCalculationRequestData `json:"data"`
}

type MarginCalculationRequestData struct {
	Segment    string `json:"segment"`
	Series     string `json:"series"`
	Exchange   string `json:"exchange"`
	Side       string `json:"side"`
	Mode       string `json:"mode"`
	Symbol     string `json:"symbol"`
	Underlying string `json:"underlying"`
	Token      string `json:"token"`
	Quantity   string `json:"quantity"` // can't add validation here because quantity is in string
	Price      string `json:"price"`    // can't add validation here because price is in string
	Product    string `json:"product"`
}

type MarginResultData struct {
	CombinedMargin         CombinedMarginData           `json:"combinedMargin"`
	IndividualMarginValues []IndividualMarginValuesData `json:"individualMarginValues"`
}

type CombinedMarginData struct {
	DeliveryMargin     float64   `json:"deliveryMargin"`
	Span               float64   `json:"span"`
	SomtierMargin      int       `json:"somtierMargin"`
	AdditionalMargin   float64   `json:"additionalMargin"`
	SpanSpreadMargin   float64   `json:"spanSpreadMargin"`
	VarMargin          float64   `json:"varMargin"`
	ExposureMargin     float64   `json:"exposureMargin"`
	PremiumMargin      float64   `json:"premiumMargin"`
	PremiumBenefit     float64   `json:"premiumBenefit"`
	ExtremeLossMargin  float64   `json:"extremeLossMargin"`
	MaxSpan            int       `json:"maxSpan"`
	NetSpan            int       `json:"netSpan"`
	NetSpanArray       []float64 `json:"netSpanArray"`
	CompositeDelta     float64   `json:"compositeDelta"`
	FutureBuyQuantity  int       `json:"futureBuyQuantity"`
	FutureSellQuantity int       `json:"futureSellQuantity"`
	OptionSellQuantity int       `json:"optionSellQuantity"`
	OptionBuyQuantity  int       `json:"optionBuyQuantity"`
	UnderlyingToken    int       `json:"underlyingToken"`
	SomRate            int       `json:"somRate"`
	SpreadRate         int       `json:"spreadRate"`
}

type IndividualMarginValuesData struct {
	DeliveryMargin     float64   `json:"deliveryMargin"`
	Span               float64   `json:"span"`
	SomtierMargin      int       `json:"somtierMargin"`
	AdditionalMargin   float64   `json:"additionalMargin"`
	SpanSpreadMargin   float64   `json:"spanSpreadMargin"`
	VarMargin          float64   `json:"varMargin"`
	ExposureMargin     float64   `json:"exposureMargin"`
	PremiumMargin      float64   `json:"premiumMargin"`
	PremiumBenefit     float64   `json:"premiumBenefit"`
	ExtremeLossMargin  float64   `json:"extremeLossMargin"`
	MaxSpan            int       `json:"maxSpan"`
	NetSpan            int       `json:"netSpan"`
	NetSpanArray       []float64 `json:"netSpanArray"`
	CompositeDelta     float64   `json:"compositeDelta"`
	FutureBuyQuantity  int       `json:"futureBuyQuantity"`
	FutureSellQuantity int       `json:"futureSellQuantity"`
	OptionSellQuantity int       `json:"optionSellQuantity"`
	OptionBuyQuantity  int       `json:"optionBuyQuantity"`
	UnderlyingToken    int       `json:"underlyingToken"`
	SomRate            int       `json:"somRate"`
	SpreadRate         int       `json:"spreadRate"`
}

type LastTradedPriceRequest struct {
	Exchange string `json:"exchange"`
	Segment  string `json:"segment"`
	Token    string `json:"token"`
}

type LastTradedPriceResponse struct {
	Price float64 `json:"price"`
}

type IcebergOrderReq struct {
	Exchange                   string  `json:"exchange" enums:"NSE,BSE,NFO,MCX,BFO" example:"NSE" validate:"oneof=NSE BSE NFO MCX"`
	InstrumentToken            string  `json:"instrumentToken"`
	ClientID                   string  `json:"clientId" binding:"required" example:"abc123"`
	OrderType                  string  `json:"orderType" enums:"LIMIT,MARKET,SL,SLM" example:"LIMIT" validate:"oneof=LIMIT MARKET SL SLM"`
	Price                      string  `json:"price"`
	Quantity                   string  `json:"quantity" validate:"gt=0"`
	DisclosedQuantity          int     `json:"disclosedQuantity"`
	Validity                   string  `json:"validity" enums:"DAY,IOC" example:"DAY" validate:"oneof=DAY IOC"`
	Product                    string  `json:"product" enums:"CNC,MIS,NRML" example:"CNC" validate:"oneof=CNC MIS NRML MTF"`
	NoOfLegs                   string  `json:"noOfLegs" validate:"gt=0"`
	GttPrice                   string  `json:"gttPrice"`
	OrderSide                  string  `json:"orderSide" enums:"BUY,SELL" example:"BUY" validate:"oneof=BUY SELL"`
	Device                     string  `json:"device"`
	UserOrderId                int     `json:"userOrderId"`
	TriggerPrice               float64 `json:"triggerPrice"`
	ExecutionType              string  `json:"executionType"`
	MarketProtectionPercentage float64 `json:"marketProtectionPercentage"`
}

type ModifyIcebergOrderReq struct {
	Exchange                   string  `json:"exchange" enums:"NSE,BSE,NFO,MCX,BFO" example:"NSE" validate:"oneof=NSE BSE NFO MCX"`
	InstrumentToken            int     `json:"instrumentToken"`
	ClientID                   string  `json:"clientId" binding:"required" example:"abc123"`
	OrderType                  string  `json:"orderType" enums:"LIMIT,MARKET,SL,SLM" example:"LIMIT" validate:"oneof=LIMIT MARKET SL SLM"`
	Price                      string  `json:"price"`
	Quantity                   int     `json:"quantity" validate:"gt=0"`
	DisclosedQuantity          int     `json:"disclosedQuantity"`
	Validity                   string  `json:"validity" enums:"DAY,IOC" example:"DAY" validate:"oneof=DAY IOC"`
	Product                    string  `json:"product" enums:"CNC,MIS,NRML" example:"CNC" validate:"oneof=CNC MIS NRML MTF"`
	GttPrice                   string  `json:"gttPrice"`
	OmsOrderId                 string  `json:"omsOrderId"`
	ExchangeOrderId            string  `json:"exchangeOrderId"`
	FilledQuantity             int     `json:"filledQuantity"`
	RemainingQuantity          int     `json:"remainingQuantity"`
	LastActivityReference      int64   `json:"lastActivityReference"`
	TriggerPrice               float64 `json:"triggerPrice"`
	ExecutionType              string  `json:"executionType"`
	MarketProtectionPercentage float64 `json:"marketProtectionPercentage"`
}

type CancelIcebergOrderReq struct {
	ClientId      string `json:"clientId"`
	OmsOrderID    string `json:"omsOrderId"`
	ExecutionType string `json:"execution_type" example:"REGULAR" validate:"oneof=REGULAR AMO"`
}

type IcebergOrderResponse struct {
	ID string `json:"id"`
}

type IcebergCanelOrderResponse struct {
	OMSOrderID string `json:"omsOrderID"`
}

type MTFCTDReq struct {
	MtfCtdValues []MtfCtdValues `json:"mtfCtdValues"`
	UserType     string         `json:"userType"`
	LoginID      string         `json:"loginId"`
	ClientID     string         `json:"clientId"`
}

type MtfCtdValues struct {
	Isin        string `json:"isin"`
	CtdQuantity int    `json:"ctdQuantity"`
}

type MTFCTDResDataRes struct {
	SuccessCount int `json:"successCount"`
	FailureCount int `json:"failureCount"`
	TotalCount   int `json:"totalCount"`
}

type CreateGttOCORequest struct {
	ExpiryTime                 string  `json:"expiryTime"`
	ClientID                   string  `json:"clientId"`
	DisclosedQuantity          int     `json:"disclosedQuantity"`
	Exchange                   string  `json:"exchange"`
	InstrumentToken            string  `json:"instrumentToken"`
	MarketProtectionPercentage int     `json:"marketProtectionPercentage"`
	OrderSide                  string  `json:"orderSide"`
	OrderType                  string  `json:"orderType"`
	Price                      float64 `json:"price" validate:"gte=0"`
	Product                    string  `json:"product"`
	Quantity                   int     `json:"quantity" validate:"gt=0"`
	SlOrderPrice               float64 `json:"slOrderPrice"`
	SlOrderQuantity            int     `json:"slOrderQuantity"`
	SlTriggerPrice             float64 `json:"slTriggerPrice"`
	TriggerPrice               float64 `json:"triggerPrice"`
	UserOrderID                int     `json:"userOrderId"`
}
