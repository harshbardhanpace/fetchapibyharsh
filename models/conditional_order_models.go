package models

// PlaceBOOrderRequest place BO order request
type PlaceBOOrderRequest struct {
	ClientID          string  `json:"clientId" example:"client123"`
	DisclosedQuantity int     `json:"disclosedQuantity" example:"2"`
	Exchange          string  `json:"exchange" enums:"NSE,BSE,NFO,CDS,MCX,BFO" validate:"oneof=NSE BSE NFO CDS MCX BFO"`
	InstrumentToken   string  `json:"instrumentToken" example:"22"`
	IsTrailing        bool    `json:"isTrailing" enums:"TRUE,FALSE"`
	OrderSide         string  `json:"orderSide" enums:"BUY,SELL"`
	OrderType         string  `json:"orderType" enums:"LIMIT,MARKET,SL,SLM" validate:"oneof=LIMIT MARKET SL SLM"`
	Price             float64 `json:"price" example:"34.2"`
	Product           string  `json:"product" enums:"CNC,MIS,NRML" validate:"oneof=CNC MIS NRML"`
	Quantity          int     `json:"quantity" example:"2"`
	SquareOffValue    float64 `json:"squareOffValue" example:"45"`
	StopLossValue     float64 `json:"stopLossValue" example:"23"`
	TrailingStopLoss  string  `json:"trailingStopLoss" example:"22"`
	TriggerPrice      float64 `json:"triggerPrice" example:"44"`
	UserOrderID       int     `json:"userOrderId" example:"1027109"`
	Validity          string  `json:"validity" enums:"DAY,IOC" example:"DAY" validate:"oneof=DAY IOC"`
}

// ModifyBOOrderRequest modify BO order request
type ModifyBOOrderRequest struct {
	Exchange              string  `json:"exchange" enums:"NSE,BSE,NFO,CDS,MCX,BFO" validate:"oneof=NSE BSE NFO CDS MCX BFO"`
	InstrumentToken       string  `json:"instrumentToken" example:"22"`
	ClientID              string  `json:"clientId" example:"client123"`
	OrderType             string  `json:"orderType" enums:"LIMIT,MARKET,SL,SLM" validate:"oneof=LIMIT MARKET SL SLM"`
	Price                 float64 `json:"price" example:"34.2"`
	Quantity              int     `json:"quantity" example:"2"`
	DisclosedQuantity     int     `json:"disclosedQuantity" example:"2"`
	Validity              string  `json:"validity" enums:"DAY,IOC" example:"DAY" validate:"oneof=DAY IOC"`
	Product               string  `json:"product" enums:"CNC,MIS,NRML" validate:"oneof=CNC MIS NRML"`
	OmsOrderID            string  `json:"omsOrderId" example:"123445"`
	ExchangeOrderID       string  `json:"exchangeOrderId" example:"ionwdg123"`
	FilledQuantity        int     `json:"filledQuantity" example:"1"`
	RemainingQuantity     int     `json:"remainingQuantity" exmaple:"1"`
	LastActivityReference int64   `json:"lastActivityReference" example:"1325938440097498600"`
	TriggerPrice          float64 `json:"triggerPrice" example:"38"`
	StopLossValue         float64 `json:"stopLossValue" example:"39"`
	SquareOffValue        float64 `json:"squareOffValue" example:"40"`
	TrailingStopLoss      float64 `json:"trailingStopLoss" example:"41"`
	IsTrailing            bool    `json:"isTrailing" enums:"TRUE,FALSE"`
}

// ExitBOOrderRequest exits BO order request
type ExitBOOrderRequest struct {
	ClientID          string `json:"clientId" example:"client123"`
	ExchangeOrderID   string `json:"exchangeOrderId" example:"ionwdg123"`
	LegOrderIndicator string `json:"legOrderIndicator" enums:"Entry, Second or Third"`
	OmsOrderID        string `json:"omsOrderId" example:"123445"`
	Status            string `json:"status" example:"CONFIRMED"`
}

// BOOrderResponse BO order response
type BOOrderResponse struct {
	BasketID string `json:"basketId"`
	Message  string `json:"message"`
}

// PlaceCOOrderRequest place co order request
type PlaceCOOrderRequest struct {
	Exchange          string  `json:"exchange" enums:"NSE,BSE,NFO,CDS,MCX,BFO" validate:"oneof=NSE BSE NFO CDS MCX BFO"`
	InstrumentToken   string  `json:"instrumentToken" example:"22"`
	ClientID          string  `json:"clientId" example:"client123"`
	OrderType         string  `json:"orderType" enums:"LIMIT,MARKET,SL,SLM" validate:"oneof=LIMIT MARKET SL SLM"`
	Price             float64 `json:"price" example:"34.2"`
	Quantity          int     `json:"quantity" example:"2"`
	DisclosedQuantity int     `json:"disclosedQuantity" example:"2"`
	Validity          string  `json:"validity" enums:"DAY,IOC" validate:"oneof=DAY IOC"`
	Product           string  `json:"product" enums:"CNC,MIS,NRML" validate:"oneof=CNC MIS NRML"`
	OrderSide         string  `json:"orderSide" enums:"BUY,SELL"`
	UserOrderID       int     `json:"userOrderId" example:"91261928"`
	StopLossValue     float64 `json:"stopLossValue" example:"33"`
	TrailingStopLoss  float64 `json:"trailingStopLoss" example:"33"`
}

// ModifyCOOrderRequest place bo order request
type ModifyCOOrderRequest struct {
	ClientID              string  `json:"clientId" example:"client123"`
	DisclosedQuantity     int     `json:"disclosedQuantity" example:"2"`
	Exchange              string  `json:"exchange" enums:"NSE,BSE,NFO,CDS,MCX,BFO" validate:"oneof=NSE BSE NFO CDS MCX BFO"`
	ExchangeOrderID       string  `json:"exchangeOrderId" example:"ionwdg123"`
	FilledQuantity        int     `json:"filledQuantity" example:"1"`
	InstrumentToken       string  `json:"instrumentToken" example:"22"`
	LastActivityReference int64   `json:"lastActivityReference" example:"1325938440097498600"`
	OmsOrderID            string  `json:"omsOrderId" example:"123445"`
	OrderType             string  `json:"orderType" enums:"LIMIT,MARKET,SL,SLM" validate:"oneof=LIMIT MARKET SL SLM"`
	Price                 float64 `json:"price" example:"34.2"`
	Product               string  `json:"product" enums:"CNC,MIS,NRML" validate:"oneof=CNC MIS NRML"`
	Quantity              int     `json:"quantity" example:"2"`
	RemainingQuantity     int     `json:"remainingQuantity" exmaple:"1"`
	StopLossValue         float64 `json:"stopLossValue"`
	TrailingStopLoss      float64 `json:"trailingStopLoss" example:"33.2"`
	Validity              string  `json:"validity" enums:"DAY,IOC" validate:"oneof=DAY IOC"`
	LegOrderIndicator     string  `json:"legOrderIndicator"`
	TriggerPrice          float64 `json:"triggerPrice" example:"34.2"`
}

// ExitCOOrderRequest exit cover order request
type ExitCOOrderRequest struct {
	ClientID          string `json:"clientId" example:"client123"`
	ExchangeOrderID   string `json:"exchangeOrderId" example:"ionwdg123"`
	LegOrderIndicator string `json:"legOrderIndicator" enums:"Entry, Second or Third"`
	OmsOrderID        string `json:"omsOrderId" example:"123445"`
}

// COOrderResponse modify bo order response
type COOrderResponse struct {
	BasketID string `json:"basketId"`
	Message  string `json:"message"`
}

// PlaceSpreadOrderRequest place spread order request
type PlaceSpreadOrderRequest struct {
	Exchange          string  `json:"exchange" enums:"NSE,BSE,NFO,CDS,MCX,BFO" validate:"oneof=NSE BSE NFO CDS MCX BFO"`
	InstrumentToken   string  `json:"instrumentToken" example:"22"`
	ClientID          string  `json:"clientId" example:"RN2363"`
	OrderType         string  `json:"orderType" enums:"LIMIT,MARKET,SL,SLM" validate:"oneof=LIMIT MARKET SL SLM"`
	Price             float64 `json:"price" example:"34.2"`
	Quantity          int     `json:"quantity" example:"2"`
	DisclosedQuantity int     `json:"disclosedQuantity" example:"2"`
	Validity          string  `json:"validity" enums:"DAY,IOC" validate:"oneof=DAY IOC"`
	Product           string  `json:"product" enums:"CNC,MIS,NRML" validate:"oneof=CNC MIS NRML"`
	OrderSide         string  `json:"orderSide" enums:"BUY,SELL"`
	UserOrderID       int     `json:"userOrderId" example:"91261928"`
}

// ModifySpreadOrderRequest modify spread order request
type ModifySpreadOrderRequest struct {
	ClientID          string  `json:"clientId" example:"RN2363"`
	DisclosedQuantity int     `json:"disclosedQuantity"`
	Exchange          string  `json:"exchange" enums:"NSE,BSE,NFO,CDS,MCX,BFO" validate:"oneof=NSE BSE NFO CDS MCX BFO"`
	ExchangeOrderID   string  `json:"exchangeOrderId"`
	InstrumentToken   string  `json:"instrumentToken" example:"22"`
	IsTrailing        bool    `json:"isTrailing" enums:"TRUE,FALSE"`
	OmsOrderID        string  `json:"omsOrderId" example:"123445"`
	OrderType         string  `json:"orderType" enums:"LIMIT,MARKET,SL,SLM" validate:"oneof=LIMIT MARKET SL SLM"`
	Price             float64 `json:"price" example:"34.2"`
	ProdType          string  `json:"prodType"`
	Product           string  `json:"product" enums:"CNC,MIS,NRML"`
	Quantity          int     `json:"quantity" example:"2"`
	StopLossValue     float64 `json:"stopLossValue" example:"39"`
	SquareOffValue    float64 `json:"squareOffValue" example:"40"`
	TrailingStopLoss  float64 `json:"trailingStopLoss" example:"33"`
	TriggerPrice      float64 `json:"triggerPrice" example:"38"`
	Validity          string  `json:"validity" enums:"DAY,IOC" validate:"oneof=DAY IOC"`
}

// ExitSpreadOrderRequest exit spread order request
type ExitSpreadOrderRequest struct {
	ClientID          string `json:"clientId" example:"RN2363"`
	LegOrderIndicator string `json:"legOrderIndicator" enums:"Entry, Second or Third"`
	OmsOrderID        string `json:"omsOrderId" example:"123445"`
	Status            string `json:"status" example:"CONFIRMED"`
	ExchangeOrderID   string `json:"exchangeOrderId" example:"ionwdg123"`
}

// SpreadOrderResponse spread order response
type SpreadOrderResponse struct {
	BasketID string `json:"basketId" example:"20210531-30"`
	Message  string `json:"message"`
}
