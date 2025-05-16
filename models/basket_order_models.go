package models

type CreateBasketReq struct {
	LoginID     string `json:"loginId"`
	Name        string `json:"name" validate:"required,max=30"`
	Type        string `json:"type"`
	ProductType string `json:"productType"`
	OrderType   string `json:"orderType"`
}

type CreateBasketRes struct {
	CreateBasketResData BasketDataRes `json:"createBasketDataRes"`
}

type FetchBasketReq struct {
	LoginID string `json:"loginId"`
}

type BasketRes struct {
	Data []BasketDataRes `json:"data"`
}

type BasketDataRes struct {
	BasketID    string               `json:"basketId"`
	BasketType  string               `json:"basketType"`
	IsExecuted  bool                 `json:"isExecuted"`
	LoginID     string               `json:"loginId"`
	Name        string               `json:"name"`
	OrderType   string               `json:"orderType"`
	Orders      []BasketDataOrderRes `json:"orders"`
	ProductType string               `json:"productType"`
	SipEligible bool                 `json:"sipEligible"`
	SipEnabled  bool                 `json:"sipEnabled"`
}

type BasketDataOrderRes struct {
	OrderID   string `json:"orderId"`
	OrderInfo struct {
		TriggerPrice        float64 `json:"triggerPrice"`
		UnderlyingToken     string  `json:"underlyingToken"`
		Series              string  `json:"series"`
		UserOrderID         int     `json:"userOrderId"`
		Exchange            string  `json:"exchange"`
		SquareOff           bool    `json:"squareOff"`
		Mode                string  `json:"mode"`
		RemainingQuantity   int     `json:"remainingQuantity"`
		AverageTradePrice   int     `json:"averageTradePrice"`
		TradePrice          int     `json:"tradePrice"`
		OrderTag            string  `json:"orderTag"`
		OrderStatusInfo     string  `json:"orderStatusInfo"`
		OrderSide           string  `json:"orderSide"`
		SquareOffValue      float64 `json:"squareOffValue"`
		ContractDescription struct {
		} `json:"contractDescription"`
		Segment                    string      `json:"segment"`
		ClientID                   string      `json:"clientId"`
		TradingSymbol              string      `json:"tradingSymbol"`
		RejectionCode              int         `json:"rejectionCode"`
		LotSize                    int         `json:"lotSize"`
		Quantity                   int         `json:"quantity"`
		LastActivityReference      int         `json:"lastActivityReference"`
		NnfID                      int         `json:"nnfId"`
		ProCli                     string      `json:"proCli"`
		Price                      float64     `json:"price"`
		OrderType                  string      `json:"orderType"`
		Validity                   string      `json:"validity"`
		TargetPriceType            string      `json:"targetPriceType"`
		InstrumentToken            int         `json:"instrumentToken"`
		SlTriggerPrice             float64     `json:"slTriggerPrice"`
		IsTrailing                 bool        `json:"isTrailing"`
		SlOrderQuantity            int         `json:"slOrderQuantity"`
		OrderEntryTime             int         `json:"orderEntryTime"`
		ExchangeTime               int         `json:"exchangeTime"`
		LegOrderIndicator          interface{} `json:"legOrderIndicator"`
		TrailingStopLoss           float64     `json:"trailingStopLoss"`
		LoginID                    interface{} `json:"loginId"`
		OmsOrderID                 string      `json:"omsOrderId"`
		MarketProtectionPercentage int         `json:"marketProtectionPercentage"`
		ExecutionType              string      `json:"executionType"`
		DisclosedQuantity          int         `json:"disclosedQuantity"`
		RejectionReason            string      `json:"rejectionReason"`
		StopLossValue              float64     `json:"stopLossValue"`
		Device                     interface{} `json:"device"`
		Product                    string      `json:"product"`
		SlOrderPrice               float64     `json:"slOrderPrice"`
		FilledQuantity             int         `json:"filledQuantity"`
		ExchangeOrderID            string      `json:"exchangeOrderId"`
		Deposit                    int         `json:"deposit"`
		AveragePrice               int         `json:"averagePrice"`
		SpreadToken                interface{} `json:"spreadToken"`
		OrderStatus                interface{} `json:"orderStatus"`
	} `json:"orderInfo"`
}

type DeleteBasketReq struct {
	BasketID string `json:"basketId"`
	Name     string `json:"name"`
	SipCount int    `json:"sipCount"`
}

type AddBasketInstrumentReq struct {
	BasketID  string `json:"basketId"`
	Name      string `json:"name"`
	OrderInfo struct {
		Exchange          string  `json:"exchange"`
		InstrumentToken   int     `json:"instrumentToken"`
		ClientID          string  `json:"clientId"`
		OrderType         string  `json:"orderType"`
		Price             float64 `json:"price" validate:"gte=0"`
		Quantity          int     `json:"quantity" validate:"gt=0"`
		DisclosedQuantity int     `json:"disclosedQuantity"`
		Validity          string  `json:"validity"`
		Product           string  `json:"product"`
		TradingSymbol     string  `json:"tradingSymbol"`
		OrderSide         string  `json:"orderSide"`
		UserOrderID       int     `json:"userOrderId"`
		UnderlyingToken   string  `json:"underlyingToken"`
		Series            string  `json:"series"`
		TriggerPrice      float64 `json:"triggerPrice"`

		ExecutionType string `json:"executionType"`
	} `json:"orderInfo"`
}

type BasketInstrumentRes struct {
	Data BasketDataRes `json:"data"`
}

type EditBasketInstrumentReq struct {
	BasketID  string `json:"basketId"`
	Name      string `json:"name"`
	OrderID   string `json:"orderId"`
	OrderInfo struct {
		Exchange          string  `json:"exchange"`
		InstrumentToken   int     `json:"instrumentToken"`
		ClientID          string  `json:"clientId"`
		OrderType         string  `json:"orderType"`
		Price             float64 `json:"price" validate:"gte=0"`
		Quantity          int     `json:"quantity" validate:"gt=0"`
		DisclosedQuantity int     `json:"disclosedQuantity"`
		Validity          string  `json:"validity"`
		Product           string  `json:"product"`
		TradingSymbol     string  `json:"tradingSymbol"`
		OrderSide         string  `json:"orderSide"`
		UserOrderID       int     `json:"userOrderId"`
		UnderlyingToken   string  `json:"underlyingToken"`
		Series            string  `json:"series"`
		OmsOrderID        string  `json:"omsOrderId"`
		ExchangeOrderID   string  `json:"exchangeOrderId"`
		TriggerPrice      float64 `json:"triggerPrice"`
		ExecutionType     string  `json:"executionType"`
	} `json:"orderInfo"`
}

type DeleteBasketInstrumentReq struct {
	BasketID string `json:"basketId"`
	OrderID  string `json:"orderId"`
	Name     string `json:"name"`
}

type RenameBasketReq struct {
	BasketID string `json:"basketId" validate:"required"`
	Name     string `json:"name" validate:"required,max=30"`
}

type ExecuteBasketReq struct {
	BasketID       string `json:"basketId"`
	Name           string `json:"name"`
	ExecutionType  string `json:"executionType"`
	SquareOff      bool   `json:"squareOff"`
	ClientID       string `json:"clientId"`
	ExecutionState bool   `json:"executionState"`
}

type ExecuteBasketRes struct {
	Data struct {
		BasketID string `json:"basketId"`
		Message  string `json:"message"`
	} `json:"data"`
}

type UpdateBasketExecutionStateReq struct {
	BasketID       string `json:"basketId"`
	Name           string `json:"name"`
	ClientID       string `json:"clientId"`
	ExecutionType  string `json:"executionType"`
	SquareOff      bool   `json:"squareOff"`
	ExecutionState bool   `json:"executionState"`
}
