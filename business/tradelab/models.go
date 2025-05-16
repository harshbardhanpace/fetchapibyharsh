package tradelab

import "time"

type TradeLabLoginReq struct {
	LoginID  string `json:"login_id"`
	Password string `json:"password"`
	Device   string `json:"device"`
}

type TradeLabTwoFaQuestions struct {
	Question   string `json:"question"`
	QuestionID int    `json:"question_id"`
}

type TradeLabTwoFaDetails struct {
	Questions  []TradeLabTwoFaQuestions `json:"questions"`
	TwofaToken string                   `json:"twofa_token"`
	Type       string                   `json:"type"`
}

type TradeLabLoginRes struct {
	Data struct {
		Alert         string               `json:"alert"`
		AuthToken     string               `json:"auth_token"`
		LoginID       string               `json:"login_id"`
		ResetPassword bool                 `json:"reset_password"`
		ResetTwoFa    bool                 `json:"reset_two_fa"`
		Twofa         TradeLabTwoFaDetails `json:"twofa"`
		TwofaEnabled  bool                 `json:"twofa_enabled"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabSetPasswordReq struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

type TradeLabSetPasswordRes struct {
	Data struct {
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabForgetPasswordReq struct {
	LoginID string `json:"login_id"`
	Pan     string `json:"pan"`
}

type TradeLabForgetPasswordRes struct {
	Data struct {
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabForgetResetTwoFaRes struct {
	Data struct {
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabErrorRes struct {
	Data struct {
	} `json:"data"`
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
	Status    string `json:"status"`
}

type TradeLabTwoFaQuestion struct {
	QuestionID string `json:"question_id"`
	Answer     string `json:"answer"`
}

type TradeLabValidateTwoFaRequest struct {
	LoginID    string                  `json:"login_id"`
	Twofa      []TradeLabTwoFaQuestion `json:"twofa"`
	TwofaToken string                  `json:"twofa_token"`
	Type       string                  `json:"type"`
}

type TradeLabValidateTwoFaResponse struct {
	Data struct {
		AuthToken     string `json:"auth_token"`
		ResetPassword bool   `json:"reset_password"`
		ResetTwoFa    bool   `json:"reset_two_fa"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradelabSetTwoFaPinRequest struct {
	LoginID   string `json:"login_id"`
	Pin       string `json:"pin"`
	TwofaType string `json:"twofa_type"`
}

type TradelabSetTwoFaPinResponse struct {
	Data struct {
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

//Tradelab Order Models

type TradelabPlaceAMORequest struct {
	Exchange          string  `json:"exchange"`
	InstrumentToken   string  `json:"instrument_token"`
	ClientID          string  `json:"client_id"`
	OrderType         string  `json:"order_type"`
	Price             float64 `json:"price"`
	Quantity          int     `json:"quantity"`
	DisclosedQuantity int     `json:"disclosed_quantity"`
	Validity          string  `json:"validity"`
	Product           string  `json:"product"`
	OrderSide         string  `json:"order_side"`
	Device            string  `json:"device"`
	UserOrderID       int     `json:"user_order_id"`
	TriggerPrice      float64 `json:"trigger_price"`
	ExecutionType     string  `json:"execution_type"`
}

type TradelabModifyAMORequest struct {
	Exchange              string  `json:"exchange"`
	InstrumentToken       int     `json:"instrument_token"`
	ClientID              string  `json:"client_id"`
	OrderType             string  `json:"order_type"`
	Price                 float64 `json:"price"`
	Quantity              int     `json:"quantity"`
	DisclosedQuantity     int     `json:"disclosed_quantity"`
	Validity              string  `json:"validity"`
	Product               string  `json:"product"`
	OmsOrderID            string  `json:"oms_order_id"`
	ExchangeOrderID       string  `json:"exchange_order_id"`
	FilledQuantity        int     `json:"filled_quantity"`
	RemainingQuantity     int     `json:"remaining_quantity"`
	LastActivityReference int     `json:"last_activity_reference"`
	TriggerPrice          float64 `json:"trigger_price"`
	ExecutionType         string  `json:"execution_type"`
}

type TradelabDeleteAMORequest struct {
	ClientID      string `json:"client_id"`
	ExecutionType string `json:"execution_type"`
}

type TradelabAMOResponse struct {
	Data struct {
		OmsOrderID string `json:"oms_order_id"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradelabPlaceOrderRequest struct {
	Exchange          string  `json:"exchange"`
	InstrumentToken   string  `json:"instrument_token"`
	ClientID          string  `json:"client_id"`
	OrderType         string  `json:"order_type"`
	Amo               bool    `json:"amo"`
	Price             float64 `json:"price"`
	Quantity          int     `json:"quantity"`
	DisclosedQuantity int     `json:"disclosed_quantity"`
	Validity          string  `json:"validity"`
	Product           string  `json:"product"`
	OrderSide         string  `json:"order_side"`
	Device            string  `json:"device"`
	UserOrderID       int     `json:"user_order_id"`
	TriggerPrice      float64 `json:"trigger_price"`
	ExecutionType     string  `json:"execution_type"`
}

type TradelabPlaceOrderMTFRequest struct {
	Exchange          string  `json:"exchange"`
	InstrumentToken   string  `json:"instrument_token"`
	ClientID          string  `json:"client_id"`
	OrderType         string  `json:"order_type"`
	Price             string  `json:"price"`
	Quantity          int     `json:"quantity" validate:"gt=0"`
	DisclosedQuantity int     `json:"disclosedQuantity"`
	Validity          string  `json:"validity" enums:"DAY,IOC" example:"DAY" validate:"oneof=DAY IOC"`
	Product           string  `json:"product" enums:"MTF" example:"MTF" validate:"oneof=MTF"`
	NoOfLegs          int     `json:"no_of_legs"`
	OrderSide         string  `json:"order_side" enums:"BUY,SELL" example:"BUY" validate:"oneof=BUY SELL"`
	Device            string  `json:"device"`
	UserOrderId       int     `json:"user_order_id"`
	TriggerPrice      float64 `json:"trigger_price"`
}

type TradeLabMTFEPledgeRequest struct {
	Depository  string      `json:"depository"`
	ClientID    string      `json:"client_id"`
	Exchange    string      `json:"exchange"`
	BoID        string      `json:"bo_id"`
	Segment     string      `json:"segment"`
	RequestType string      `json:"request_type"`
	IsinDetails IsinDetails `json:"isin_details"`
	Order       struct {
		Price             string `json:"price"`
		Device            string `json:"device"`
		Product           string `json:"product"`
		Exchange          string `json:"exchange"`
		Quantity          int    `json:"quantity"`
		Validity          string `json:"validity"`
		ClientID          string `json:"client_id"`
		OrderSide         string `json:"order_side"`
		OrderType         string `json:"order_type"`
		UserOrderID       string `json:"user_order_id"`
		InstrumentToken   string `json:"instrument_token"`
		DisclosedQuantity int    `json:"disclosed_quantity"`
	} `json:"order"`
}

type IsinDetails struct {
	IsinName string `json:"isin_name"`
	Isin     string `json:"isin"`
	Quantity string `json:"quantity"`
	Price    string `json:"price"`
}

type TradeLabMTFPledgeListResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		List []struct {
			ClientID             string  `json:"client_id"`
			Isin                 string  `json:"isin"`
			TotalPledgeQuantity  int     `json:"total_pledge_quantity"`
			CtdQuantity          int     `json:"ctd_quantity"`
			Symbol               string  `json:"symbol"`
			AvgPrice             float64 `json:"avg_price"`
			MarginMultiplier     int     `json:"margin_multiplier"`
			CtdMarginValue       int     `json:"ctd_margin_value"`
			Token                int     `json:"token"`
			Exchange             string  `json:"exchange"`
			CreatedAt            string  `json:"created_at"`
			UpdatedAt            string  `json:"updated_at"`
			EdisApprovedQuantity int     `json:"edis_approved_quantity"`
			ObligationQuantity   int     `json:"obligation_quantity"`
			UsedQuantity         int     `json:"used_quantity"`
			LoginID              string  `json:"login_id"`
			MarginValue          float64 `json:"margin_value"`
			TotalInvestedAmount  float64 `json:"total_invested_amount"`
			BrokerAmount         float64 `json:"broker_amount"`
		} `json:"list"`
		TotalCount int `json:"total_count"`
	} `json:"data"`
}

type TradelabMTFResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	ErrorCode int    `json:"error_code"`
	Data      []any  `json:"data"`
}

type TLMTFCTDListRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		List       []MTFCTDList `json:"list"`
		TotalCount int          `json:"total_count"`
	} `json:"data"`
}

type TLMtfPledgeListRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		List []struct {
			ClientID            string    `json:"client_id"`
			Isin                string    `json:"isin"`
			PledgeQuantity      int       `json:"pledge_quantity"`
			ToBePledgedQuantity int       `json:"to_be_pledged_quantity"`
			Segment             string    `json:"segment"`
			Symbol              string    `json:"symbol"`
			MtfSettlementDate   string    `json:"mtf_settlement_date"`
			MtfSquareOffDate    string    `json:"mtf_square_off_date"`
			NseToken            int       `json:"nse_token"`
			BseToken            int       `json:"bse_token"`
			AvgPrice            int       `json:"avg_price"`
			MarginMultiplier    int       `json:"margin_multiplier"`
			MarginVarElm        int       `json:"margin_var_elm"`
			MarginValue         float64   `json:"margin_value"`
			DaysTillSquareoff   int       `json:"days_till_squareoff"`
			IsLastDayOfMtf      bool      `json:"is_last_day_of_mtf"`
			IsCfObligation      bool      `json:"is_cf_obligation"`
			CreatedAt           time.Time `json:"CreatedAt"`
			UpdatedAt           time.Time `json:"UpdatedAt"`
		} `json:"list"`
	} `json:"data"`
}

type MTFCTDList struct {
	ClientID             string  `json:"client_id"`
	Isin                 string  `json:"isin"`
	TotalPledgeQuantity  int     `json:"total_pledge_quantity"`
	CtdQuantity          int     `json:"ctd_quantity"`
	Symbol               string  `json:"symbol"`
	AvgPrice             float64 `json:"avg_price"`
	MarginMultiplier     int     `json:"margin_multiplier"`
	CtdMarginValue       float64 `json:"ctd_margin_value"`
	Token                int     `json:"token"`
	Exchange             string  `json:"exchange"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
	EdisApprovedQuantity int     `json:"edis_approved_quantity"`
	ObligationQuantity   int     `json:"obligation_quantity"`
	UsedQuantity         int     `json:"used_quantity"`
	LoginID              string  `json:"login_id"`
	MarginValue          int     `json:"margin_value"`
	TotalInvestedAmount  int     `json:"total_invested_amount"`
	BrokerAmount         int     `json:"broker_amount"`
}

type TradelabPlaceOrderResponse struct {
	Data struct {
		OmsOrderID  string `json:"oms_order_id"`
		UserOrderID int    `json:"user_order_id"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradelabModifyOrderRequest struct {
	Exchange              string  `json:"exchange"`
	InstrumentToken       int     `json:"instrument_token"`
	ClientID              string  `json:"client_id"`
	OrderType             string  `json:"order_type"`
	Price                 float64 `json:"price"`
	Quantity              int     `json:"quantity"`
	DisclosedQuantity     int     `json:"disclosed_quantity"`
	Validity              string  `json:"validity"`
	Product               string  `json:"product"`
	OmsOrderID            string  `json:"oms_order_id"`
	TriggerPrice          float64 `json:"trigger_price"`
	ExecutionType         string  `json:"execution_type"`
	ExchangeOrderID       string  `json:"exchange_order_id"`
	FilledQuantity        int     `json:"filled_quantity"`
	RemainingQuantity     int     `json:"remaining_quantity"`
	LastActivityReference int64   `json:"last_activity_reference"`
}

type TradelabCancelOrModifyResponse struct {
	Data struct {
		OmsOrderID string `json:"oms_order_id"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabPlaceBOOrderRequest struct {
	ClientID          string  `json:"client_id"`
	Device            string  `json:"device"`
	DisclosedQuantity int     `json:"disclosed_quantity"`
	Exchange          string  `json:"exchange"`
	ExecutionType     string  `json:"execution_type"`
	InstrumentToken   string  `json:"instrument_token"`
	IsTrailing        bool    `json:"is_trailing"`
	OrderSide         string  `json:"order_side"`
	OrderType         string  `json:"order_type"`
	Price             float64 `json:"price"`
	Product           string  `json:"product"`
	Quantity          int     `json:"quantity"`
	SquareOffValue    float64 `json:"square_off_value"`
	StopLossValue     float64 `json:"stop_loss_value"`
	TrailingStopLoss  string  `json:"trailing_stop_loss"`
	TriggerPrice      float64 `json:"trigger_price"`
	UserOrderID       int     `json:"user_order_id"`
	Validity          string  `json:"validity"`
}

type TradeLabModifyBOOrderRequest struct {
	Exchange              string  `json:"exchange"`
	InstrumentToken       string  `json:"instrument_token"`
	ClientID              string  `json:"client_id"`
	OrderType             string  `json:"order_type"`
	Price                 float64 `json:"price"`
	Quantity              int     `json:"quantity"`
	DisclosedQuantity     int     `json:"disclosed_quantity"`
	Validity              string  `json:"validity"`
	Product               string  `json:"product"`
	OmsOrderID            string  `json:"oms_order_id"`
	ExchangeOrderID       string  `json:"exchange_order_id"`
	FilledQuantity        int     `json:"filled_quantity"`
	RemainingQuantity     int     `json:"remaining_quantity"`
	LastActivityReference int64   `json:"last_activity_reference"`
	TriggerPrice          float64 `json:"trigger_price"`
	StopLossValue         float64 `json:"stop_loss_value"`
	SquareOffValue        float64 `json:"square_off_value"`
	TrailingStopLoss      float64 `json:"trailing_stop_loss"`
	IsTrailing            bool    `json:"is_trailing"`
	ExecutionType         string  `json:"execution_type"`
}

type TradeLabCancelBOOrderRequest struct {
	ClientID          string `json:"client_id"`
	ExchangeOrderID   string `json:"exchange_order_id"`
	ExecutionType     string `json:"execution_type"`
	LegOrderIndicator string `json:"leg_order_indicator"`
	OmsOrderID        string `json:"oms_order_id"`
	Status            string `json:"status"`
}

type TradeLabPlaceBOOrderResponse struct {
	Data struct {
		Data struct {
			BasketID string `json:"basket_id"`
			Message  string `json:"message"`
		} `json:"data"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabModOrExitBOOrderResponse struct {
	Data struct {
		BasketID string `json:"basket_id"`
		Message  string `json:"message"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabPlaceCOOrderRequest struct {
	Exchange          string  `json:"exchange"`
	InstrumentToken   string  `json:"instrument_token"`
	ClientID          string  `json:"client_id"`
	OrderType         string  `json:"order_type"`
	Price             float64 `json:"price"`
	Quantity          int     `json:"quantity"`
	DisclosedQuantity int     `json:"disclosed_quantity"`
	Validity          string  `json:"validity"`
	Product           string  `json:"product"`
	OrderSide         string  `json:"order_side"`
	Device            string  `json:"device"`
	UserOrderID       int     `json:"user_order_id"`
	ExecutionType     string  `json:"execution_type"`
	StopLossValue     float64 `json:"stop_loss_value"`
	TrailingStopLoss  float64 `json:"trailing_stop_loss"`
}

type TradeLabModifyCOOrderRequest struct {
	ClientID              string  `json:"client_id"`
	DisclosedQuantity     int     `json:"disclosed_quantity"`
	Exchange              string  `json:"exchange"`
	ExchangeOrderID       string  `json:"exchange_order_id"`
	ExecutionType         string  `json:"execution_type"`
	FilledQuantity        int     `json:"filled_quantity"`
	InstrumentToken       string  `json:"instrument_token"`
	LastActivityReference int64   `json:"last_activity_reference"`
	OmsOrderID            string  `json:"oms_order_id"`
	OrderType             string  `json:"order_type"`
	Price                 float64 `json:"price"`
	Product               string  `json:"product"`
	Quantity              int     `json:"quantity"`
	RemainingQuantity     int     `json:"remaining_quantity"`
	StopLossValue         float64 `json:"stop_loss_value"`
	TrailingStopLoss      float64 `json:"trailing_stop_loss"`
	Validity              string  `json:"validity"`
	LegOrderIndicator     string  `json:"leg_order_indicator"`
	TriggerPrice          float64 `json:"trigger_price"`
}

type TradeLabCancelCOOrderRequest struct {
	ClientID          string `json:"client_id"`
	ExchangeOrderID   string `json:"exchange_order_id"`
	ExecutionType     string `json:"execution_type"`
	LegOrderIndicator string `json:"leg_order_indicator"`
	OmsOrderID        string `json:"oms_order_id"`
}

type TradeLabPlaceCOOrderResponse struct {
	Data struct {
		Data struct {
			BasketID string `json:"basket_id"`
			Message  string `json:"message"`
		} `json:"data"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabModifyOrExitCOOrderResponse struct {
	Data struct {
		BasketID string `json:"basket_id"`
		Message  string `json:"message"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// TradelabPlaceSpreadOrderRequest tradelab place spread order request
type TradelabPlaceSpreadOrderRequest struct {
	Exchange          string  `json:"exchange"`
	InstrumentToken   string  `json:"instrument_token"`
	ClientID          string  `json:"client_id"`
	OrderType         string  `json:"order_type"`
	Price             float64 `json:"price"`
	Quantity          int     `json:"quantity"`
	DisclosedQuantity int     `json:"disclosed_quantity"`
	Validity          string  `json:"validity"`
	Product           string  `json:"product"`
	OrderSide         string  `json:"order_side"`
	Device            string  `json:"device"`
	UserOrderID       int     `json:"user_order_id"`
	ExecutionType     string  `json:"execution_type"`
}

// TradelabModifySpreadOrderRequest tradelab modify spread order response
type TradelabModifySpreadOrderRequest struct {
	ClientID          string  `json:"client_id"`
	DisclosedQuantity int     `json:"disclosed_quantity"`
	Exchange          string  `json:"exchange"`
	ExchangeOrderID   string  `json:"exchange_order_id"`
	ExecutionType     string  `json:"execution_type"`
	InstrumentToken   string  `json:"instrument_token"`
	IsTrailing        bool    `json:"is_trailing"`
	OmsOrderID        string  `json:"oms_order_id"`
	OrderType         string  `json:"order_type"`
	Price             float64 `json:"price"`
	ProdType          string  `json:"prod_type"`
	Product           string  `json:"product"`
	Quantity          int     `json:"quantity"`
	SquareOffValue    float64 `json:"square_off_value"`
	StopLossValue     float64 `json:"stop_loss_value"`
	TrailingStopLoss  float64 `json:"trailing_stop_loss"`
	TriggerPrice      float64 `json:"trigger_price"`
	Validity          string  `json:"validity"`
}

// TradelabCancelSpreadOrderRequest tradelab cancel spread order response
type TradelabCancelSpreadOrderRequest struct {
	ClientID          string `json:"client_id"`
	LegOrderIndicator string `json:"leg_order_indicator"`
	OmsOrderID        string `json:"oms_order_id"`
	Status            string `json:"status"`
	ExecutionType     string `json:"execution_type"`
	ExchangeOrderID   string `json:"exchange_order_id"`
}

// TradelabSpreadOrderResponse tradelab spread order response
type TradelabSpreadOrderResponse struct {
	Data struct {
		Data struct {
			BasketID string `json:"basket_id"`
			Message  string `json:"message"`
		} `json:"data"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradelabExitSpreadOrderResponse struct {
	Data struct {
		BasketID string `json:"basket_id"`
		Message  string `json:"message"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradelabPendingOrderResponse struct {
	Data    TradeLabPendingOrderResponseData `json:"data"`
	Message string                           `json:"message"`
	Status  string                           `json:"status"`
}

type TradeLabPendingOrderResponseData struct {
	Orders []TradelabPendingOrderResponseOrders `json:"orders"`
}

type TradelabPendingOrderResponseOrders struct {
	TradingSymbol              string  `json:"trading_symbol"`
	AverageTradePrice          float64 `json:"average_trade_price"`
	Exchange                   string  `json:"exchange"`
	ProCli                     string  `json:"pro_cli"`
	MarketProtectionPercentage int     `json:"market_protection_percentage"`
	OrderEntryTime             int     `json:"order_entry_time"`
	Mode                       string  `json:"mode"`
	OmsOrderID                 string  `json:"oms_order_id"`
	TrailingStopLoss           float64 `json:"trailing_stop_loss"`
	Deposit                    int     `json:"deposit"`
	SquareOffValue             float64 `json:"square_off_value"`
	DisclosedQuantity          int     `json:"disclosed_quantity"`
	StopLossValue              float64 `json:"stop_loss_value"`
	Price                      float64 `json:"price"`
	OrderTag                   string  `json:"order_tag"`
	Device                     string  `json:"device"`
	RemainingQuantity          int     `json:"remaining_quantity"`
	LastActivityReference      int64   `json:"last_activity_reference"`
	AveragePrice               float64 `json:"average_price"`
	SquareOff                  bool    `json:"square_off?"`
	OrderStatusInfo            string  `json:"order_status_info"`
	Quantity                   int     `json:"quantity"`
	ExecutionType              string  `json:"execution_type"`
	ClientID                   string  `json:"client_id"`
	ExchangeTime               int     `json:"exchange_time"`
	OrderSide                  string  `json:"order_side"`
	LoginID                    string  `json:"login_id"`
	Validity                   string  `json:"validity"`
	InstrumentToken            int     `json:"instrument_token"`
	Product                    string  `json:"product"`
	TriggerPrice               float64 `json:"trigger_price"`
	Segment                    string  `json:"segment"`
	TradePrice                 float64 `json:"trade_price"`
	OrderType                  string  `json:"order_type"`
	ContractDescription        struct {
	} `json:"contract_description"`
	RejectionCode     int    `json:"rejection_code"`
	LegOrderIndicator string `json:"leg_order_indicator"`
	ExchangeOrderID   string `json:"exchange_order_id"`
	OrderStatus       string `json:"order_status"`
	FilledQuantity    int    `json:"filled_quantity"`
	TargetPriceType   string `json:"target_price_type"`
	IsTrailing        bool   `json:"is_trailing"`
	UserOrderID       string `json:"user_order_id"`
	LotSize           int    `json:"lot_size"`
	Series            string `json:"series"`
	NnfID             int64  `json:"nnf_id"`
	RejectionReason   string `json:"rejection_reason"`
}

type TLOrderUpdatePacket struct {
	TriggerPrice               float64 `json:"trigger_price"`
	IsTrailing                 bool    `json:"is_trailing"`
	RejectionCode              string  `json:"rejection_code"`
	LastActivityReference      int64   `json:"last_activity_reference"`
	TransactionType            string  `json:"transaction_type"`
	NestRequestID              string  `json:"nest_request_id"`
	OrderEntryTime             int     `json:"order_entry_time"`
	Lotsize                    float64 `json:"lotsize"`
	Device                     string  `json:"device"`
	FilledQuantity             int     `json:"filled_quantity"`
	ExchangeOrderID            string  `json:"exchange_order_id"`
	Deposit                    float64 `json:"deposit"`
	InstrumentToken            int     `json:"instrument_token"`
	Product                    string  `json:"product"`
	Validity                   string  `json:"validity"`
	OrderStatus                string  `json:"order_status"`
	Price                      float64 `json:"price"`
	LotSize                    float64 `json:"lot_size"`
	ProCli                     string  `json:"pro_cli"`
	Nnfid                      int64   `json:"nnfid"`
	AveragePrice               float64 `json:"average_price"`
	ClientID                   string  `json:"client_id"`
	OrderSide                  string  `json:"order_side"`
	LoginID                    string  `json:"login_id"`
	SpreadToken                string  `json:"spread_token"`
	Exchange                   string  `json:"exchange"`
	LegOrderIndicator          string  `json:"leg_order_indicator"`
	ExecutionType              string  `json:"execution_type"`
	DisclosedQuantity          int     `json:"disclosed_quantity"`
	UserOrderID                string  `json:"user_order_id"`
	OrderType                  string  `json:"order_type"`
	AverageTradePrice          float64 `json:"average_trade_price"`
	BasketID                   string  `json:"basket_id"`
	TrailingStopLoss           float64 `json:"trailing_stop_loss"`
	RejectionReason            string  `json:"rejection_reason"`
	Quantity                   int     `json:"quantity"`
	Mode                       string  `json:"mode"`
	ExchangeTime               int     `json:"exchange_time"`
	Symbol                     string  `json:"symbol"`
	StopLossValue              float64 `json:"stop_loss_value"`
	Series                     string  `json:"series"`
	TradingSymbol              string  `json:"trading_symbol"`
	RemainingQuantity          int     `json:"remaining_quantity"`
	SquareOffValue             float64 `json:"square_off_value"`
	MarketProtectionPercentage float64 `json:"market_protection_percentage"`
	TradePrice                 float64 `json:"trade_price"`
	OmsOrderID                 string  `json:"oms_order_id"`
}

type TradelabCompletedOrderResponse struct {
	Data    TradelabCompletedOrderResponseData `json:"data"`
	Message string                             `json:"message"`
	Status  string                             `json:"status"`
}

type TradelabCompletedOrderResponseData struct {
	Orders []TradelabCompletedOrderResponseOrders `json:"orders"`
}

type TradelabCompletedOrderResponseOrders struct {
	TradingSymbol              string  `json:"trading_symbol"`
	AverageTradePrice          float64 `json:"average_trade_price"`
	Exchange                   string  `json:"exchange"`
	ProCli                     string  `json:"pro_cli"`
	MarketProtectionPercentage int     `json:"market_protection_percentage"`
	OrderEntryTime             int     `json:"order_entry_time"`
	Mode                       string  `json:"mode"`
	OmsOrderID                 string  `json:"oms_order_id"`
	TrailingStopLoss           float64 `json:"trailing_stop_loss"`
	Deposit                    int     `json:"deposit"`
	SquareOffValue             float64 `json:"square_off_value"`
	DisclosedQuantity          int     `json:"disclosed_quantity"`
	StopLossValue              float64 `json:"stop_loss_value"`
	Price                      float64 `json:"price"`
	OrderTag                   string  `json:"order_tag"`
	Device                     string  `json:"device"`
	RemainingQuantity          int     `json:"remaining_quantity"`
	LastActivityReference      int     `json:"last_activity_reference"`
	AveragePrice               float64 `json:"average_price"`
	SquareOff                  bool    `json:"square_off?"`
	OrderStatusInfo            string  `json:"order_status_info"`
	Quantity                   int     `json:"quantity"`
	ExecutionType              string  `json:"execution_type"`
	ClientID                   string  `json:"client_id"`
	ExchangeTime               int     `json:"exchange_time"`
	OrderSide                  string  `json:"order_side"`
	LoginID                    string  `json:"login_id"`
	Validity                   string  `json:"validity"`
	InstrumentToken            int     `json:"instrument_token"`
	Product                    string  `json:"product"`
	TriggerPrice               float64 `json:"trigger_price"`
	Segment                    string  `json:"segment"`
	TradePrice                 float64 `json:"trade_price"`
	OrderType                  string  `json:"order_type"`
	ContractDescription        struct {
	} `json:"contract_description"`
	RejectionCode     int    `json:"rejection_code"`
	LegOrderIndicator string `json:"leg_order_indicator"`
	ExchangeOrderID   string `json:"exchange_order_id"`
	OrderStatus       string `json:"order_status"`
	FilledQuantity    int    `json:"filled_quantity"`
	TargetPriceType   string `json:"target_price_type"`
	IsTrailing        bool   `json:"is_trailing"`
	UserOrderID       string `json:"user_order_id"`
	LotSize           int    `json:"lot_size"`
	Series            string `json:"series"`
	NnfID             int64  `json:"nnf_id"`
	RejectionReason   string `json:"rejection_reason"`
}

type TradeLabTradeBookResponse struct {
	Data    TradeLabTradeBookResponseData `json:"data"`
	Message string                        `json:"message"`
	Status  string                        `json:"status"`
}

type TradeLabTradeBookResponseData struct {
	Trades []TradeLabTradeBookResponseTrades `json:"trades"`
}

type TradeLabTradeBookResponseTrades struct {
	BookType              string  `json:"book_type"`
	BrokerID              string  `json:"broker_id"`
	ClientID              string  `json:"client_id"`
	DisclosedVol          int     `json:"disclosed_vol"`
	DisclosedVolRemaining int     `json:"disclosed_vol_remaining"`
	Exchange              string  `json:"exchange"`
	ExchangeOrderID       string  `json:"exchange_order_id"`
	ExchangeTime          int     `json:"exchange_time"`
	FillNumber            string  `json:"fill_number"`
	FilledQuantity        int     `json:"filled_quantity"`
	GoodTillDate          int     `json:"good_till_date"`
	InstrumentToken       int     `json:"instrument_token"`
	LoginID               string  `json:"login_id"`
	OmsOrderID            string  `json:"oms_order_id"`
	OrderEntryTime        int     `json:"order_entry_time"`
	OrderPrice            float64 `json:"order_price"`
	OrderSide             string  `json:"order_side"`
	OrderType             string  `json:"order_type"`
	OriginalVol           int     `json:"original_vol"`
	Pan                   string  `json:"pan"`
	ProCli                int     `json:"pro_cli"`
	Product               string  `json:"product"`
	RemainingQuantity     int     `json:"remaining_quantity"`
	TradeNumber           string  `json:"trade_number"`
	TradePrice            float64 `json:"trade_price"`
	TradeQuantity         int     `json:"trade_quantity"`
	TradeTime             int     `json:"trade_time"`
	TradingSymbol         string  `json:"trading_symbol"`
	TriggerPrice          float64 `json:"trigger_price"`
	VLoginID              string  `json:"v_login_id"`
	VolFilledToday        int     `json:"vol_filled_today"`
}

type TlTradeUpdate struct {
	VolumeFilledToday int     `json:"volume_filled_today"`
	VLoginID          string  `json:"v_login_id"`
	UserOrderID       string  `json:"user_order_id"`
	TransactionType   any     `json:"transaction_type"`
	TradingSymbol     string  `json:"trading_symbol"`
	TradePrice        float64 `json:"trade_price"`
	TradeID           string  `json:"trade_id"`
	Symbol            string  `json:"symbol"`
	Strike            string  `json:"strike"`
	Series            string  `json:"series"`
	Segment           string  `json:"segment"`
	RemainingQuantity int     `json:"remaining_quantity"`
	Product           string  `json:"product"`
	ProCli            string  `json:"pro_cli"`
	OrderSide         string  `json:"order_side"`
	OptionType        string  `json:"option_type"`
	OmsOrderID        string  `json:"oms_order_id"`
	LoginID           string  `json:"login_id"`
	InstrumentToken   int     `json:"instrument_token"`
	InstrumentName    string  `json:"instrument_name"`
	FilledQuantity    int     `json:"filled_quantity"`
	Expiry            string  `json:"expiry"`
	ExchangeTime      int     `json:"exchange_time"`
	ExchangeOrderID   string  `json:"exchange_order_id"`
	Exchange          string  `json:"exchange"`
	ClientID          string  `json:"client_id"`
}

type TradeLabOrderHistoryResponse struct {
	Data    []TradeLabOrderHistoryResponseData `json:"data"`
	Message string                             `json:"message"`
	Status  string                             `json:"status"`
}

type TradeLabOrderHistoryResponseData struct {
	AvgPrice          float64 `json:"avg_price"`
	ClientID          string  `json:"client_id"`
	ClientOrderID     string  `json:"client_order_id"`
	CreatedAt         int     `json:"created_at"`
	DisclosedQuantity int     `json:"disclosed_quantity"`
	Exchange          string  `json:"exchange"`
	ExchangeOrderID   string  `json:"exchange_order_id"`
	ExchangeTime      int     `json:"exchange_time"`
	FillQuantity      int     `json:"fill_quantity"`
	LastModified      int64   `json:"last_modified"`
	LoginID           string  `json:"login_id"`
	ModifiedAt        int     `json:"modified_at"`
	OrderID           string  `json:"order_id"`
	OrderMode         string  `json:"order_mode"`
	OrderSide         string  `json:"order_side"`
	OrderType         string  `json:"order_type"`
	Price             float64 `json:"price"`
	Product           string  `json:"product"`
	Quantity          int     `json:"quantity"`
	RejectReason      string  `json:"reject_reason"`
	RemainingQuantity int     `json:"remaining_quantity"`
	Segment           string  `json:"segment"`
	Status            string  `json:"status"`
	Symbol            string  `json:"symbol"`
	Token             int     `json:"token"`
	TriggerPrice      float64 `json:"trigger_price"`
	UnderlyingToken   int     `json:"underlying_token"`
	Validity          string  `json:"validity"`
}

// Fetch Demat Holdings Response
type TradeLabFetchDematHoldingsResponse struct {
	Data    TradelabFetchDematHoldingsHoldingsData `json:"data"`
	Message string                                 `json:"message"`
	Status  string                                 `json:"status"`
}

type TradelabFetchDematHoldingsHoldingsData struct {
	Holdings []TradelabFetchDematHoldingsHoldings `json:"holdings"`
}

type TradelabFetchDematHoldingsHoldings struct {
	BranchCode        string  `json:"branch_code"`
	BuyAvg            float64 `json:"buy_avg"`
	BuyAvgMtm         float64 `json:"buy_avg_mtm"`
	ClientID          string  `json:"client_id"`
	Exchange          string  `json:"exchange"`
	FreeQuantity      int     `json:"free_quantity"`
	InstrumentDetails struct {
		Exchange        int    `json:"exchange"`
		InstrumentName  string `json:"instrument_name"`
		InstrumentToken int    `json:"instrument_token"`
		TradingSymbol   string `json:"trading_symbol"`
	} `json:"instrument_details"`
	Isin                  string  `json:"isin"`
	Ltp                   float64 `json:"ltp"`
	PendingQuantity       int     `json:"pending_quantity"`
	PledgeAllow           bool    `json:"pledge_allow"`
	PledgeQuantity        int     `json:"pledge_quantity"`
	PreviousClose         float64 `json:"previous_close"`
	Quantity              int     `json:"quantity"`
	Symbol                string  `json:"symbol"`
	T0Price               float64 `json:"t0_price"`
	T0Quantity            int     `json:"t0_quantity"`
	T1Price               float64 `json:"t1_price"`
	T1Quantity            int     `json:"t1_quantity"`
	T2Price               float64 `json:"t2_price"`
	T2Quantity            int     `json:"t2_quantity"`
	TodayPledgeQuantity   int     `json:"today_pledge_quantity"`
	TodayUnpledgeQuantity int     `json:"today_unpledge_quantity"`
	Token                 int     `json:"token"`
	TradingSymbol         string  `json:"trading_symbol"`
	TransactionType       string  `json:"transaction_type"`
	UsedQuantity          int     `json:"used_quantity"`
	ActualBuyAvg          float64 `json:"actual_buy_avg"`
	NetHoldingQty         int     `json:"net_holding_quantity"`
	PledgePercentage      float64 `json:"pledge_percentage"`
}

// Convert Positions Request
type TradeLabConvertPositionsRequest struct {
	ClientID        string `json:"client_id"`
	Exchange        string `json:"exchange"`
	InstrumentToken int    `json:"instrument_token"`
	Product         string `json:"product"`
	NewProduct      string `json:"new_product"`
	Quantity        int    `json:"quantity"`
	Validity        string `json:"validity"`
	OrderSide       string `json:"order_side"`
}

// Convert Positions Resposne
type TradeLabConvertPositionsResponse struct {
	Data struct {
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabGetPositionResponse struct {
	Data    []TradeLabGetPositionResponseData `json:"data"`
	Message string                            `json:"message"`
	Status  string                            `json:"status"`
}
type TradeLabGetPositionResponseData struct {
	AverageSellPrice       float64 `json:"average_sell_price"`
	TotalPledgeCollateral  any     `json:"total_pledge_collateral"`
	Exchange               string  `json:"exchange"`
	ProCli                 string  `json:"pro_cli"`
	AveragePrice           float64 `json:"average_price"`
	CfBuyQuantity          int     `json:"cf_buy_quantity"`
	CfSellQuantity         int     `json:"cf_sell_quantity"`
	SellQuantity           int     `json:"sell_quantity"`
	ActualCfBuyAmount      float64 `json:"actual_cf_buy_amount"`
	ClientID               string  `json:"client_id"`
	ActualAverageBuyPrice  float64 `json:"actual_average_buy_price"`
	CfBuyAmount            float64 `json:"cf_buy_amount"`
	ActualAverageSellPrice float64 `json:"actual_average_sell_price"`
	ClosePrice             float64 `json:"close_price"`
	BuyQuantity            int     `json:"buy_quantity"`
	AverageBuyPrice        float64 `json:"average_buy_price"`
	SellAmount             float64 `json:"sell_amount"`
	OtherMargin            any     `json:"other_margin"`
	NetAmount              float64 `json:"net_amount"`
	RealizedMtm            float64 `json:"realized_mtm"`
	Multiplier             float64 `json:"multiplier"`
	PreviousClose          float64 `json:"previous_close"`
	Segment                string  `json:"segment"`
	Product                string  `json:"product"`
	BuyAmount              float64 `json:"buy_amount"`
	VLoginID               string  `json:"v_login_id"`
	CfSellAmount           float64 `json:"cf_sell_amount"`
	Token                  int     `json:"token"`
	TradingSymbol          string  `json:"trading_symbol"`
	NetQuantity            int     `json:"net_quantity"`
	Symbol                 string  `json:"symbol"`
	InstrumentToken        int     `json:"instrument_token"`
	Ltp                    float64 `json:"ltp"`
	ActualCfSellAmount     float64 `json:"actual_cf_sell_amount"`
	ProdType               int     `json:"prod_type"`
}

// Fetch Option Chain Response
type TradeLabFetchOptionChainResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Result []struct {
		ExpiryDate string `json:"expiry_date"`
		Strikes    []struct {
			StrikePrice float64 `json:"strike_price"`
			CallOption  struct {
				Token         string  `json:"token"`
				Exchange      string  `json:"exchange"`
				Company       string  `json:"company"`
				Symbol        string  `json:"symbol"`
				TradingSymbol string  `json:"trading_symbol"`
				DisplayName   string  `json:"display_name"`
				StrikePrice   float64 `json:"strike_price"`
				ExpiryRaw     string  `json:"ExpiryRaw"`
				ClosePrice    string  `json:"close_price"`
			} `json:"call_option"`
			PutOption struct {
				Token         string  `json:"token"`
				Exchange      string  `json:"exchange"`
				Company       string  `json:"company"`
				Symbol        string  `json:"symbol"`
				TradingSymbol string  `json:"trading_symbol"`
				DisplayName   string  `json:"display_name"`
				StrikePrice   float64 `json:"strike_price"`
				ExpiryRaw     string  `json:"ExpiryRaw"`
				ClosePrice    string  `json:"close_price"`
			} `json:"put_option"`
		} `json:"strikes"`
	} `json:"result"`
}

// Profile Response
type TradeLabProfileResponse struct {
	Data struct {
		Branch            string   `json:"branch"`
		BankBranchName    string   `json:"bank_branch_name"`
		OfficeAddr        string   `json:"office_addr"`
		DpID              []string `json:"dp_id"`
		City              string   `json:"city"`
		PermanentAddr     string   `json:"permanent_addr"`
		BankName          string   `json:"bank_name"`
		BankAccountNumber string   `json:"bank_account_number"`
		PanNumber         string   `json:"pan_number"`
		Role              struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"role"`
		EmailID                string   `json:"email_id"`
		BrokerID               string   `json:"broker_id"`
		ClientID               string   `json:"client_id"`
		BankState              string   `json:"bank_state"`
		AccountType            string   `json:"account_type"`
		Status                 string   `json:"status"`
		UserType               string   `json:"user_type"`
		LastPasswordChangeDate int      `json:"last_password_change_date"`
		BoID                   []string `json:"bo_id"`
		BasketEnabled          bool     `json:"basket_enabled"`
		TwofaEnabled           bool     `json:"twofa_enabled"`
		Name                   string   `json:"name"`
		Depository             string   `json:"depository"`
		ExchangeNnf            struct {
			Bse int `json:"BSE"`
			Mcx int `json:"MCX"`
			Nfo int `json:"NFO"`
			Nse int `json:"NSE"`
		} `json:"exchange_nnf"`
		PoaStatus           bool     `json:"poa_status"`
		BankCity            string   `json:"bank_city"`
		IfscCode            string   `json:"ifsc_code"`
		Dob                 string   `json:"dob"`
		ExchangesSubscribed []string `json:"exchanges_subscribed"`
		Sex                 string   `json:"sex"`
		PoaEnabled          bool     `json:"poa_enabled"`
		BackofficeLink      string   `json:"backoffice_link"`
		State               string   `json:"state"`
		PhoneNumber         string   `json:"phone_number"`
		ProductsEnabled     []string `json:"products_enabled"`
		ProfileURL          string   `json:"profile_url"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabAccountFreezeRes struct {
	Data struct {
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// Fetch Funds Response
type TradeLabFetchFundsResponse struct {
	Data struct {
		ClientID string     `json:"client_id"`
		Headers  []string   `json:"headers"`
		Values   [][]string `json:"values"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// Search Scrip Response
type TradeLabSearchScripResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Result []TradeLabSearchScripResponseResult `json:"result"`
}

type TradeLabSearchScripResponseResult struct {
	Token         string  `json:"token"`
	Exchange      string  `json:"exchange"`
	Execution     string  `json:"execution"`
	Company       string  `json:"company"`
	Symbol        string  `json:"symbol"`
	Isin          string  `json:"isin"`
	TradingSymbol string  `json:"trading_symbol"`
	DisplayName   string  `json:"display_name"`
	Score         float64 `json:"score"`
	ClosePrice    string  `json:"close_price"`
	IsTradable    bool    `json:"is_tradable"`
	Segment       string  `json:"segment"`
	Tag           string  `json:"tag"`
	Expiry        string  `json:"expiry"`
	Alternate     struct {
		Token         string  `json:"token"`
		Exchange      string  `json:"exchange"`
		Execution     string  `json:"execution"`
		Company       string  `json:"company"`
		Symbol        string  `json:"symbol"`
		TradingSymbol string  `json:"trading_symbol"`
		DisplayName   string  `json:"display_name"`
		Score         float64 `json:"score"`
		ClosePrice    string  `json:"close_price"`
		IsTradable    bool    `json:"is_tradable"`
		Segment       string  `json:"segment"`
		Tag           string  `json:"tag"`
		Expiry        string  `json:"expiry"`
	} `json:"alternate,omitempty"`
	//Alternate0 struct {
	//} `json:"alternate,omitempty"`
	//Alternate1 struct {
	//} `json:"alternate,omitempty"`
}

// ScripInfoResponse
type TradeLabScripInfoResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Result struct {
		BoardLotQuantity         int     `json:"board_lot_quantity"`
		ChangeInOi               int     `json:"change_in_oi"`
		Exchange                 int     `json:"exchange"`
		Expiry                   int     `json:"expiry"`
		HigherCircuitLimit       float64 `json:"higher_circuit_limit"`
		InstrumentName           string  `json:"instrument_name"`
		InstrumentToken          int     `json:"instrument_token"`
		Isin                     string  `json:"isin"`
		LowerCircuitLimit        float64 `json:"lower_circuit_limit"`
		Multiplier               int     `json:"multiplier"`
		OpenInterest             int     `json:"open_interest"`
		OptionType               string  `json:"option_type"`
		Precision                int     `json:"precision"`
		Series                   string  `json:"series"`
		Strike                   int     `json:"strike"`
		Symbol                   string  `json:"symbol"`
		TickSize                 float64 `json:"tick_size"`
		TradingSymbol            string  `json:"trading_symbol"`
		UnderlyingToken          int     `json:"underlying_token"`
		RawExpiry                int     `json:"raw_expiry"`
		Freeze                   int     `json:"freeze"`
		InstrumentType           string  `json:"instrument_type"`
		IssueRate                int     `json:"issue_rate"`
		IssueStartDate           string  `json:"issue_start_date"`
		ListDate                 string  `json:"list_date"`
		MaxOrderSize             int     `json:"max_order_size"`
		PriceNumerator           float64 `json:"price_numerator"`
		PriceDenominator         float64 `json:"price_denominator"`
		Comments                 string  `json:"comments"`
		CircuitRating            string  `json:"circuit_rating"`
		CompanyName              string  `json:"company_name"`
		DisplayName              string  `json:"display_name"`
		RawTickSize              int     `json:"raw_tick_size"`
		IsIndex                  bool    `json:"is_index"`
		Tradable                 bool    `json:"tradable"`
		MaxSingleQty             int     `json:"max_single_qty"`
		ExpiryString             string  `json:"expiry_string"`
		LocalUpdateTime          string  `json:"local_update_time"`
		MarketType               string  `json:"market_type"`
		PriceUnits               string  `json:"price_units"`
		TradingUnits             string  `json:"trading_units"`
		LastTradingDate          string  `json:"last_trading_date"`
		TenderPeriodEndDate      string  `json:"tender_period_end_date"`
		DeliveryStartDate        string  `json:"delivery_start_date"`
		PriceQuotation           float64 `json:"price_quotation"`
		GeneralDenominator       string  `json:"general_denominator"`
		TenderPeriodStartDate    string  `json:"tender_period_start_date"`
		DeliveryUnits            string  `json:"delivery_units"`
		DeliveryEndDate          string  `json:"delivery_end_date"`
		TradingUnitFactor        int     `json:"trading_unit_factor"`
		DeliveryUnitFactor       int     `json:"delivery_unit_factor"`
		BookClosureEndDate       string  `json:"book_closure_end_date"`
		BookClosureStartDate     string  `json:"book_closure_start_date"`
		NoDeliveryDateEnd        string  `json:"no_delivery_date_end"`
		NoDeliveryDateStart      string  `json:"no_delivery_date_start"`
		ReAdmissionDate          string  `json:"re_admission_date"`
		RecordDate               string  `json:"record_date"`
		Warning                  string  `json:"warning"`
		Dpr                      string  `json:"dpr"`
		TradeToTrade             bool    `json:"trade_to_trade"`
		SurveillanceIndicator    int     `json:"surveillance_indicator"`
		PartitionID              int     `json:"partition_id"`
		ProductID                int     `json:"product_id"`
		ProductCategory          string  `json:"product_category"`
		MonthIdentifier          int     `json:"month_identifier"`
		ClosePrice               string  `json:"close_price"`
		SpecialPreopen           int     `json:"special_preopen"`
		AlternateExchange        string  `json:"alternate_exchange"`
		AlternateToken           int     `json:"alternate_token"`
		Asm                      string  `json:"asm"`
		Gsm                      string  `json:"gsm"`
		Execution                string  `json:"execution"`
		Symbol2                  string  `json:"symbol2"`
		RawTenderPeriodStartDate string  `json:"raw_tender_period_start_date"`
		RawTenderPeriodEndDate   string  `json:"raw_tender_period_end_date"`
		YearlyHighPrice          string  `json:"yearly_high_price"`
		YearlyLowPrice           string  `json:"yearly_low_price"`
		IssueMaturityDate        int     `json:"issue_maturity_date"`
		Var                      string  `json:"var"`
		Exposure                 string  `json:"exposure"`
		Span                     []int   `json:"span"`
		HaveFutures              bool    `json:"have_futures"`
		HaveOptions              bool    `json:"have_options"`
		Tag                      string  `json:"tag"`
		ShortCode                string  `json:"short_code"`
		IsMisEligible            bool    `json:"is_mis_eligible"`
		IsMtfEligible            bool    `json:"is_mtf_eligible"`
		ExBonusDate              string  `json:"ex_bonus_date"`
		ExDate                   string  `json:"ex_date"`
		Exflag                   string  `json:"ex_flag"`
		ExRightDate              string  `json:"ex_right_date"`
		MtfMargin                float64 `json:"mtf_margin"`
	} `json:"result"`
}

// TradelabCreateGTTOrderRequest create gtt req
type TradeLabCreateGTTOrderRequest struct {
	ActionType string `json:"action_type"`
	ExpiryTime string `json:"expiry_time"`
	Order      struct {
		ClientID                   string  `json:"client_id"`
		Device                     string  `json:"device"`
		DisclosedQuantity          int     `json:"disclosed_quantity"`
		Exchange                   string  `json:"exchange"`
		InstrumentToken            string  `json:"instrument_token"`
		MarketProtectionPercentage int     `json:"market_protection_percentage"`
		OrderSide                  string  `json:"order_side"`
		OrderType                  string  `json:"order_type"`
		Price                      float64 `json:"price"`
		Product                    string  `json:"product"`
		Quantity                   int     `json:"quantity"`
		SlOrderPrice               float64 `json:"sl_order_price"`
		SlOrderQuantity            int     `json:"sl_order_quantity"`
		SlTriggerPrice             float64 `json:"sl_trigger_price"`
		TriggerPrice               float64 `json:"trigger_price"`
		UserOrderID                int     `json:"user_order_id"`
	} `json:"order"`
}

// TradelabGTTOrderResponse  gtt response
type TradelabGTTOrderResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// TradelabModifyGTTOrderRequest modify gtt order
type TradelabModifyGTTOrderRequest struct {
	ExpiryTime string `json:"expiry_time"`
	ActionType string `json:"action_type"`
	ID         string `json:"id"`
	Order      struct {
		ClientID                   string  `json:"client_id"`
		Device                     string  `json:"device"`
		DisclosedQuantity          int     `json:"disclosed_quantity"`
		Exchange                   string  `json:"exchange"`
		InstrumentToken            string  `json:"instrument_token"`
		MarketProtectionPercentage int     `json:"market_protection_percentage"`
		OrderSide                  string  `json:"order_side"`
		OrderType                  string  `json:"order_type"`
		Price                      float64 `json:"price"`
		Product                    string  `json:"product"`
		Quantity                   int     `json:"quantity"`
		SlOrderPrice               int     `json:"sl_order_price"`
		SlOrderQuantity            int     `json:"sl_order_quantity"`
		SlTriggerPrice             int     `json:"sl_trigger_price"`
		TriggerPrice               float64 `json:"trigger_price"`
		UserOrderID                int     `json:"user_order_id"`
	} `json:"order"`
}

// TradelabFetchGTTOrderResponse fetch gtt order response
type TradelabFetchGTTOrderResponse struct {
	Data    []TradelabFetchGTTOrderResponseData `json:"data"`
	Message string                              `json:"message"`
	Status  string                              `json:"status"`
}

type TradelabFetchGTTOrderResponseData struct {
	ActionType string `json:"action_type"`
	ClientID   string `json:"client_id"`
	CreatedAt  string `json:"created_at"`
	ExpiryTime string `json:"expiry_time"`
	ID         string `json:"id"`
	LoginID    string `json:"login_id"`
	Order      struct {
		DisclosedQty     int     `json:"disclosed_qty"`
		Exchange         string  `json:"exchange"`
		ExecutionType    string  `json:"execution_type"`
		InstrumentToken  int     `json:"instrument_token"`
		Mode             string  `json:"mode"`
		OrderSide        string  `json:"order_side"`
		OrderType        string  `json:"order_type"`
		Price            float64 `json:"price"`
		ProCli           string  `json:"pro_cli"`
		ProdType         string  `json:"prod_type"`
		Product          string  `json:"product"`
		Quantity         int     `json:"quantity"`
		Segment          string  `json:"segment"`
		SlOrderPrice     float64 `json:"sl_order_price"`
		SlOrderQuantity  int     `json:"sl_order_quantity"`
		SlTriggerPrice   float64 `json:"sl_trigger_price"`
		SquareOffPrice   float64 `json:"square_off_price"`
		Token            int     `json:"token"`
		TradingSymbol    string  `json:"trading_symbol"`
		TrailingStopLoss int     `json:"trailing_stop_loss"`
		TriggerPrice     float64 `json:"trigger_price"`
		Validity         string  `json:"validity"`
		VendorCode       string  `json:"vendor_code"`
	} `json:"order"`
	RejectCode   int    `json:"reject_code"`
	RejectReason string `json:"reject_reason"`
	Status       string `json:"status"`
	Type         string `json:"type"`
	UpdatedAt    string `json:"updated_at"`
}

// Get All Ipo Response
type TradeLabGetAllIpoResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		AllIpo struct {
			Data   []IpoState `json:"data"`
			Status string     `json:"status"`
		} `json:"all_ipo"`
		OpenIpo     []IpoState `json:"open_ipo"`
		UpcomingIpo []IpoState `json:"upcoming_ipo"`
		ClosedIpo   []IpoState `json:"closed_ipo"`
	} `json:"data"`
}

type IpoState struct {
	BiddingStartDate           string  `json:"biddingStartDate"`
	Symbol                     string  `json:"symbol"`
	MinBidQuantity             int     `json:"minBidQuantity"`
	Registrar                  string  `json:"registrar"`
	LotSize                    int     `json:"lotSize"`
	T1ModEndDate               string  `json:"t1ModEndDate"`
	DailyStartTime             string  `json:"dailyStartTime"`
	T1ModStartTime             string  `json:"t1ModStartTime"`
	BiddingEndDate             string  `json:"biddingEndDate"`
	T1ModEndTime               string  `json:"t1ModEndTime"`
	DailyEndTime               string  `json:"dailyEndTime"`
	TickSize                   float64 `json:"tickSize"`
	IssueType                  string  `json:"issueType"`
	FaceValue                  float64 `json:"faceValue"`
	MinPrice                   float64 `json:"minPrice"`
	T1ModStartDate             string  `json:"t1ModStartDate"`
	Name                       string  `json:"name"`
	IssueSize                  int     `json:"issueSize"`
	MaxPrice                   float64 `json:"maxPrice"`
	CutOffPrice                float64 `json:"cutOffPrice"`
	UnixBiddingEndDate         int     `json:"unixBiddingEndDate"`
	UnixBiddingStartDate       int     `json:"unixBiddingStartDate"`
	Isin                       string  `json:"isin"`
	AllotmentDate              string  `json:"allotmentDate"`
	ExchangeIssueType          string  `json:"exchange_issue_type"`
	AllotmentBegins            string  `json:"allotment_begins"`
	RefundDate                 string  `json:"refundDate"`
	ListingDate                string  `json:"listingDate"`
	AboutCompany               string  `json:"aboutCompany"`
	ParentCompany              string  `json:"parentCompany"`
	FoundedYear                string  `json:"foundedYear"`
	ProspectusFileURL          string  `json:"prospectusFileUrl"`
	ManagingDirector           string  `json:"managingDirector"`
	MaxLimit                   float64 `json:"MaxLimit"`
	RetailDiscount             float64 `json:"RetailDiscount"`
	NseExchangeListed          bool    `json:"nse_exchange_listed"`
	BseExchangeListed          bool    `json:"bse_exchange_listed"`
	AmoOrderEntryTime          string  `json:"amo_order_entry_time"`
	ApplicationRangeStart      int     `json:"application_range_start"`
	ApplicationRangeEnd        int     `json:"application_range_end"`
	TotalApplicationRangeCount int     `json:"total_application_range_count"`
	CategoryDetails            any     `json:"categoryDetails"`
	SubCategorySettings        []struct {
		SubCatCode    string `json:"subCatCode"`
		MinValue      any    `json:"minValue"`
		MaxUpiLimit   int    `json:"maxUpiLimit"`
		AllowCutOff   bool   `json:"allowCutOff"`
		AllowUpi      bool   `json:"allowUpi"`
		MaxValue      any    `json:"maxValue"`
		DiscountPrice any    `json:"discountPrice"`
		DiscountType  string `json:"discountType"`
		MaxPrice      any    `json:"maxPrice"`
		CaCode        string `json:"caCode"`
		Allowed       bool   `json:"allowed"`
		StartDate     string `json:"start_date"`
		EndDate       string `json:"end_date"`
		DisplayName   string `json:"displayName"`
		MinLotSize    int    `json:"min_lot_size"`
		StartTime     string `json:"startTime"`
		EndTime       string `json:"endTime"`
	} `json:"subCategorySettings"`
	IpoAllowed        bool      `json:"ipoAllowed"`
	BseAllowed        bool      `json:"bse_allowed"`
	NseAllowed        bool      `json:"nse_allowed"`
	SubType           string    `json:"subType"`
	EnablePio         bool      `json:"enable_pio"`
	PioStartDate      time.Time `json:"pio_start_date"`
	PioEndDate        time.Time `json:"pio_end_date"`
	PioEndTime        time.Time `json:"pio_end_time"`
	PioStartTime      time.Time `json:"pio_start_time"`
	DematTransferDate string    `json:"dematTransferDate"`
	MandateEndDate    string    `json:"mandateEndDate"`
	IsEmployeeCat     bool      `json:"is_employee_cat"`
	IsShareHolderCat  bool      `json:"is_share_holder_cat"`
}

type TradeLabPlaceIpoOrderRequest struct {
	ClientID      string                     `json:"client_id"`
	Symbol        string                     `json:"symbol"`
	UpiID         string                     `json:"upi_id"`
	Bids          []PlaceIpoOrderRequestBids `json:"bids"`
	AllotmentMode string                     `json:"allotment_mode"`
	BankAccount   string                     `json:"bank_account"`
	BankCode      string                     `json:"bank_code"`
	Broker        string                     `json:"broker"`
	CategoryCode  string                     `json:"category_code"`
	ClientBenID   string                     `json:"client_ben_id"`
	ClientName    string                     `json:"client_name"`
	DpID          string                     `json:"dp_id"`
	Ifsc          string                     `json:"ifsc"`
	LocationCode  string                     `json:"location_code"`
	NonAsba       bool                       `json:"non_asba"`
	Pan           string                     `json:"pan"`
	Category      string                     `json:"category"`
}

type PlaceIpoOrderRequestBids struct {
	ActivityType string `json:"activityType"`
	Quantity     int    `json:"quantity"`
	AtCutOff     bool   `json:"atCutOff"`
	Price        int    `json:"price"`
	Amount       int    `json:"amount"`
}

type TradeLabPlaceIpoOrderResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type TradeLabFetchIpoOrderRequest struct {
	ClientID string `json:"client_id"`
}

type TradeLabFetchIpoOrderResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []struct {
		Symbol            string      `json:"symbol"`
		Reason            string      `json:"reason"`
		ApplicationNumber string      `json:"applicationNumber"`
		ClientName        string      `json:"clientName"`
		ChequeNumber      string      `json:"chequeNumber"`
		ReferenceNumber   string      `json:"referenceNumber"`
		DpVerStatusFlag   string      `json:"dpVerStatusFlag"`
		SubBrokerCode     string      `json:"subBrokerCode"`
		Depository        string      `json:"depository"`
		ReasonCode        int         `json:"reasonCode"`
		Pan               string      `json:"pan"`
		Ifsc              string      `json:"ifsc"`
		Timestamp         string      `json:"timestamp"`
		BankAccount       string      `json:"bankAccount"`
		BankCode          string      `json:"bankCode"`
		DpVerReason       string      `json:"dpVerReason"`
		DpID              string      `json:"dpId"`
		Upi               string      `json:"upi"`
		UpiAmtBlocked     interface{} `json:"upiAmtBlocked"`
		Bids              []struct {
			AtCutOff           bool    `json:"atCutOff"`
			Amount             int     `json:"amount"`
			Quantity           int     `json:"quantity"`
			BidReferenceNumber int64   `json:"bidReferenceNumber"`
			Series             string  `json:"series"`
			Price              float64 `json:"price"`
			ActivityType       string  `json:"activityType"`
			Status             string  `json:"status"`
		} `json:"bids"`
		AllotmentMode                   string  `json:"allotmentMode"`
		DpVerFailCode                   string  `json:"dpVerFailCode"`
		NonASBA                         bool    `json:"nonASBA"`
		UpiFlag                         string  `json:"upiFlag"`
		Category                        string  `json:"category"`
		LocationCode                    string  `json:"locationCode"`
		ClientBenID                     string  `json:"clientBenId"`
		ClientID                        string  `json:"clientId"`
		Status                          string  `json:"status"`
		Mode                            string  `json:"mode"`
		Allotmentstatus                 string  `json:"allotmentstatus"`
		Allotmentdate                   string  `json:"allotmentdate"`
		Allotmentupdated                string  `json:"allotmentupdated"`
		Allotmentquantity               int     `json:"allotmentquantity"`
		Allotmentprice                  float64 `json:"allotmentprice"`
		CategoryCode                    string  `json:"category_code"`
		CategoryDisplayName             string  `json:"category_display_name"`
		IsAmoOrder                      bool    `json:"isAmoOrder"`
		PaymentMode                     string  `json:"paymentMode"`
		AmtBlockTime                    string  `json:"amtBlockTime"`
		Modify                          bool    `json:"modify"`
		IsOrderModify                   bool    `json:"is_order_modify"`
		IsBseIpo                        bool    `json:"is_bse_ipo"`
		IsNseIpo                        bool    `json:"is_nse_ipo"`
		UpiPaymentStatusMessage         string  `json:"upi_payment_status_message"`
		ExchangeUpdatedUpiBlockedAmount int     `json:"exchange_updated_upi_blocked_amount"`
		IsPioOrder                      bool    `json:"is_pio_order"`
	} `json:"data"`
}

type TradeLabCancelIpoOrderRequest struct {
	ClientID          string                      `json:"client_id"`
	Symbol            string                      `json:"symbol"`
	ApplicationNumber string                      `json:"applicationNumber"`
	Bids              []CancelIpoOrderRequestBids `json:"bids"`
	UpiID             string                      `json:"upi_id"`
	AllotmentMode     string                      `json:"allotment_mode"`
	BankAccount       string                      `json:"bank_account"`
	BankCode          string                      `json:"bank_code"`
	Broker            string                      `json:"broker"`
	ClientBenID       string                      `json:"client_ben_id"`
	ClientName        string                      `json:"client_name"`
	DpID              string                      `json:"dp_id"`
	Ifsc              string                      `json:"ifsc"`
	LocationCode      string                      `json:"location_code"`
	NonAsba           bool                        `json:"non_asba"`
	Pan               string                      `json:"pan"`
}

type CancelIpoOrderRequestBids struct {
	Quantity           int     `json:"quantity"`
	AtCutOff           bool    `json:"atCutOff"`
	Price              float64 `json:"price"`
	Amount             int     `json:"amount"`
	BidReferenceNumber int64   `json:"bidReferenceNumber"`
	Series             string  `json:"series"`
	ActivityType       string  `json:"activityType"`
	Status             string  `json:"status"`
}

type TradeLabCancelIpoOrderResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type MarginCalculationRequest []struct {
	Segment    string `json:"segment"`
	Series     string `json:"series"`
	Exchange   string `json:"exchange"`
	Side       string `json:"side"`
	Mode       string `json:"mode"`
	Symbol     string `json:"symbol"`
	Underlying string `json:"underlying"`
	Token      string `json:"token"`
	Quantity   string `json:"quantity" validate:"gt=0"`
	Price      string `json:"price" validate:"gte=0"`
	Product    string `json:"product"`
}

// type MarginCalculationResponse struct {
// 	Error struct {
// 		Code    int    `json:"code"`
// 		Message string `json:"message"`
// 	} `json:"error"`
// 	Result struct {
// 		CombinedMargin         CombinedMarginData           `json:"combined_margin"`
// 		IndividualMarginValues []IndividualMarginValuesData `json:"individual_margin_values"`
// 	} `json:"result"`
// }

type MarginCalculationResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Result struct {
		CombinedMargin struct {
			DeliveryMargin     float64   `json:"delivery_margin"`
			Span               float64   `json:"span"`
			SomtierMargin      int       `json:"somtier_margin"`
			AdditionalMargin   float64   `json:"additional_margin"`
			SpanSpreadMargin   float64   `json:"span_spread_margin"`
			VarMargin          float64   `json:"var_margin"`
			ExposureMargin     float64   `json:"exposure_margin"`
			PremiumMargin      float64   `json:"premium_margin"`
			PremiumBenefit     float64   `json:"premium_benefit"`
			ExtremeLossMargin  float64   `json:"extreme_loss_margin"`
			MaxSpan            int       `json:"max_span"`
			NetSpan            int       `json:"net_span"`
			NetSpanArray       []float64 `json:"net_span_array"`
			CompositeDelta     float64   `json:"composite_delta"`
			FutureBuyQuantity  int       `json:"future_buy_quantity"`
			FutureSellQuantity int       `json:"future_sell_quantity"`
			OptionSellQuantity int       `json:"option_sell_quantity"`
			OptionBuyQuantity  int       `json:"option_buy_quantity"`
			UnderlyingToken    int       `json:"underlying_token"`
			SomRate            int       `json:"som_rate"`
			SpreadRate         int       `json:"spread_rate"`
		} `json:"combined_margin"`
		IndividualMarginValues []struct {
			DeliveryMargin     float64   `json:"delivery_margin"`
			Span               float64   `json:"span"`
			SomtierMargin      int       `json:"somtier_margin"`
			AdditionalMargin   float64   `json:"additional_margin"`
			SpanSpreadMargin   float64   `json:"span_spread_margin"`
			VarMargin          float64   `json:"var_margin"`
			ExposureMargin     float64   `json:"exposure_margin"`
			PremiumMargin      float64   `json:"premium_margin"`
			PremiumBenefit     float64   `json:"premium_benefit"`
			ExtremeLossMargin  float64   `json:"extreme_loss_margin"`
			MaxSpan            int       `json:"max_span"`
			NetSpan            int       `json:"net_span"`
			NetSpanArray       []float64 `json:"net_span_array"`
			CompositeDelta     float64   `json:"composite_delta"`
			FutureBuyQuantity  int       `json:"future_buy_quantity"`
			FutureSellQuantity int       `json:"future_sell_quantity"`
			OptionSellQuantity int       `json:"option_sell_quantity"`
			OptionBuyQuantity  int       `json:"option_buy_quantity"`
			UnderlyingToken    int       `json:"underlying_token"`
			SomRate            int       `json:"som_rate"`
			SpreadRate         int       `json:"spread_rate"`
		} `json:"individual_margin_values"`
	} `json:"result"`
}

type CombinedMarginData struct {
	DeliveryMargin     int     `json:"delivery_margin"`
	Span               int     `json:"span"`
	SomtierMargin      int     `json:"somtier_margin"`
	AdditionalMargin   int     `json:"additional_margin"`
	SpanSpreadMargin   int     `json:"span_spread_margin"`
	VarMargin          float64 `json:"var_margin"`
	ExposureMargin     int     `json:"exposure_margin"`
	PremiumMargin      int     `json:"premium_margin"`
	PremiumBenefit     int     `json:"premium_benefit"`
	ExtremeLossMargin  int     `json:"extreme_loss_margin"`
	MaxSpan            int     `json:"max_span"`
	NetSpan            int     `json:"net_span"`
	NetSpanArray       []int   `json:"net_span_array"`
	CompositeDelta     int     `json:"composite_delta"`
	FutureBuyQuantity  int     `json:"future_buy_quantity"`
	FutureSellQuantity int     `json:"future_sell_quantity"`
	OptionSellQuantity int     `json:"option_sell_quantity"`
	OptionBuyQuantity  int     `json:"option_buy_quantity"`
	UnderlyingToken    int     `json:"underlying_token"`
	SomRate            int     `json:"som_rate"`
	SpreadRate         int     `json:"spread_rate"`
}

type IndividualMarginValuesData struct {
	DeliveryMargin     int     `json:"delivery_margin"`
	Span               int     `json:"span"`
	SomtierMargin      int     `json:"somtier_margin"`
	AdditionalMargin   int     `json:"additional_margin"`
	SpanSpreadMargin   int     `json:"span_spread_margin"`
	VarMargin          float64 `json:"var_margin"`
	ExposureMargin     int     `json:"exposure_margin"`
	PremiumMargin      int     `json:"premium_margin"`
	PremiumBenefit     int     `json:"premium_benefit"`
	ExtremeLossMargin  int     `json:"extreme_loss_margin"`
	MaxSpan            int     `json:"max_span"`
	NetSpan            int     `json:"net_span"`
	NetSpanArray       []int   `json:"net_span_array"`
	CompositeDelta     int     `json:"composite_delta"`
	FutureBuyQuantity  int     `json:"future_buy_quantity"`
	FutureSellQuantity int     `json:"future_sell_quantity"`
	OptionSellQuantity int     `json:"option_sell_quantity"`
	OptionBuyQuantity  int     `json:"option_buy_quantity"`
	UnderlyingToken    int     `json:"underlying_token"`
	SomRate            int     `json:"som_rate"`
	SpreadRate         int     `json:"spread_rate"`
}

type TradeLabTopGainerLoserResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Losers []struct {
			TurnoverValue   float64 `json:"turnover_value"`
			TradedQuantity  int     `json:"tradedQuantity"`
			Symbol          string  `json:"symbol"`
			NetPrice        float64 `json:"netPrice"`
			Ltp             float64 `json:"ltp"`
			LotSize         int     `json:"lot_size"`
			InstrumentToken int     `json:"instrument_token"`
			Exchange        string  `json:"exchange"`
			CompanyName     string  `json:"company_name"`
			ClosePrice      float64 `json:"close_price"`
		} `json:"losers"`
		Gainers []struct {
			TurnoverValue   float64 `json:"turnover_value"`
			TradedQuantity  int     `json:"tradedQuantity"`
			Symbol          string  `json:"symbol"`
			NetPrice        float64 `json:"netPrice"`
			Ltp             float64 `json:"ltp"`
			LotSize         int     `json:"lot_size"`
			InstrumentToken int     `json:"instrument_token"`
			Exchange        string  `json:"exchange"`
			CompanyName     string  `json:"company_name"`
			ClosePrice      float64 `json:"close_price"`
		} `json:"gainers"`
	} `json:"data"`
}

type TradeLabMostActiveVolumeResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		MostActiveVolume []struct {
			TurnoverValue     float64 `json:"turnover_value"`
			TradedQuantity    int     `json:"tradedQuantity"`
			TotalSellQuantity int     `json:"total_sell_quantity"`
			TotalBuyQuantity  int     `json:"total_buy_quantity"`
			Symbol            string  `json:"symbol"`
			PreviousPrice     float64 `json:"previousPrice"`
			NetPrice          float64 `json:"netPrice"`
			Ltp               float64 `json:"ltp"`
			LotSize           int     `json:"lot_size"`
			InstrumentToken   int     `json:"instrument_token"`
			Exchange          string  `json:"exchange"`
			CompanyName       string  `json:"company_name"`
		} `json:"most_active_volume"`
	} `json:"data"`
}

type TradeLabChartDataResponse struct {
	Status string `json:"status"`
	Data   struct {
		Candles [][]interface{} `json:"candles"`
	} `json:"data"`
}

type TradeLabLastTradedPrice struct {
	Data    float64 `json:"data"`
	Message string  `json:"message"`
	Status  string  `json:"status"`
}

type TradeLabMCXLastTradedPriceData struct {
	AskPrice          float64 `json:"ask_price"`
	AskQty            int     `json:"ask_qty"`
	AverageTradePrice float64 `json:"average_trade_price"`
	BidPrice          float64 `json:"bid_price"`
	BidQty            int     `json:"bid_qty"`
	ClosePrice        float64 `json:"close_price"`
	Exchange          string  `json:"exchange"`
	ExchangeTimestamp int64   `json:"exchange_timestamp"`
	HighPrice         float64 `json:"high_price"`
	InstrumentToken   int     `json:"instrument_token"`
	LastTradePrice    float64 `json:"last_trade_price"`
	LastTradeQty      int     `json:"last_trade_qty"`
	LastTradeTime     int64   `json:"last_trade_time"`
	LowPrice          float64 `json:"low_price"`
	OpenPrice         float64 `json:"open_price"`
	TotalBuyQty       int     `json:"total_buy_qty"`
	TotalSellQty      int     `json:"total_sell_qty"`
	TradeVolume       int     `json:"trade_volume"`
	YearlyHigh        float64 `json:"yearly_high"`
	YearlyLow         float64 `json:"yearly_low"`
}

type TradeLabMCXLastTradedPrice struct {
	Data    TradeLabMCXLastTradedPriceData `json:"data"`
	Message string                         `json:"message"`
	Status  string                         `json:"status"`
}

type TradeLabCreateBasketReq struct {
	LoginID     string `json:"login_id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	ProductType string `json:"product_type"`
	OrderType   string `json:"order_type"`
}

type TradeLabFetchBasketReq struct {
	LoginID string `json:"login_id"`
}

type TradeLabBasketRes struct {
	Data []struct {
		BasketID   string `json:"basket_id"`
		BasketType string `json:"basket_type"`
		IsExecuted bool   `json:"is_executed"`
		LoginID    string `json:"login_id"`
		Name       string `json:"name"`
		OrderType  string `json:"order_type"`
		Orders     []struct {
			OrderID   string `json:"order_id"`
			OrderInfo struct {
				TriggerPrice        float64 `json:"trigger_price"`
				UnderlyingToken     string  `json:"underlying_token"`
				Series              string  `json:"series"`
				UserOrderID         int     `json:"user_order_id"`
				Exchange            string  `json:"exchange"`
				SquareOff           bool    `json:"square_off"`
				Mode                string  `json:"mode"`
				RemainingQuantity   int     `json:"remaining_quantity"`
				AverageTradePrice   int     `json:"average_trade_price"`
				TradePrice          int     `json:"trade_price"`
				OrderTag            string  `json:"order_tag"`
				OrderStatusInfo     string  `json:"order_status_info"`
				OrderSide           string  `json:"order_side"`
				SquareOffValue      float64 `json:"square_off_value"`
				ContractDescription struct {
				} `json:"contract_description"`
				Segment                    string      `json:"segment"`
				ClientID                   string      `json:"client_id"`
				TradingSymbol              string      `json:"trading_symbol"`
				RejectionCode              int         `json:"rejection_code"`
				LotSize                    int         `json:"lot_size"`
				Quantity                   int         `json:"quantity"`
				LastActivityReference      int         `json:"last_activity_reference"`
				NnfID                      int         `json:"nnf_id"`
				ProCli                     string      `json:"pro_cli"`
				Price                      float64     `json:"price"`
				OrderType                  string      `json:"order_type"`
				Validity                   string      `json:"validity"`
				TargetPriceType            string      `json:"target_price_type"`
				InstrumentToken            int         `json:"instrument_token"`
				SlTriggerPrice             float64     `json:"sl_trigger_price"`
				IsTrailing                 bool        `json:"is_trailing"`
				SlOrderQuantity            int         `json:"sl_order_quantity"`
				OrderEntryTime             int         `json:"order_entry_time"`
				ExchangeTime               int         `json:"exchange_time"`
				LegOrderIndicator          interface{} `json:"leg_order_indicator"`
				TrailingStopLoss           float64     `json:"trailing_stop_loss"`
				LoginID                    interface{} `json:"login_id"`
				OmsOrderID                 string      `json:"oms_order_id"`
				MarketProtectionPercentage int         `json:"market_protection_percentage"`
				ExecutionType              string      `json:"execution_type"`
				DisclosedQuantity          int         `json:"disclosed_quantity"`
				RejectionReason            string      `json:"rejection_reason"`
				StopLossValue              float64     `json:"stop_loss_value"`
				Device                     interface{} `json:"device"`
				Product                    string      `json:"product"`
				SlOrderPrice               float64     `json:"sl_order_price"`
				FilledQuantity             int         `json:"filled_quantity"`
				ExchangeOrderID            string      `json:"exchange_order_id"`
				Deposit                    int         `json:"deposit"`
				AveragePrice               int         `json:"average_price"`
				SpreadToken                interface{} `json:"spread_token"`
				OrderStatus                interface{} `json:"order_status"`
			} `json:"order_info"`
		} `json:"orders"`
		ProductType string `json:"product_type"`
		SipEligible bool   `json:"sip_eligible"`
		SipEnabled  bool   `json:"sip_enabled"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabDeleteBasketReq struct {
	BasketID string `json:"basket_id"`
	Name     string `json:"name"`
	SipCount int    `json:"sip_count"`
}

type TradeLabAddBasketInstrumentReq struct {
	BasketID  string `json:"basket_id"`
	Name      string `json:"name"`
	OrderInfo struct {
		Exchange          string  `json:"exchange"`
		InstrumentToken   string  `json:"instrument_token"`
		ClientID          string  `json:"client_id"`
		OrderType         string  `json:"order_type"`
		Price             float64 `json:"price"`
		Quantity          int     `json:"quantity"`
		DisclosedQuantity int     `json:"disclosed_quantity"`
		Validity          string  `json:"validity"`
		Product           string  `json:"product"`
		TradingSymbol     string  `json:"trading_symbol"`
		OrderSide         string  `json:"order_side"`
		UserOrderID       int     `json:"user_order_id"`
		UnderlyingToken   string  `json:"underlying_token"`
		Series            string  `json:"series"`
		Device            string  `json:"device"`
		TriggerPrice      float64 `json:"trigger_price"`
		ExecutionType     string  `json:"execution_type"`
	} `json:"order_info"`
}

type TradeLabBasketInstrumentRes struct {
	Data struct {
		BasketID   string `json:"basket_id"`
		BasketType string `json:"basket_type"`
		IsExecuted bool   `json:"is_executed"`
		LoginID    string `json:"login_id"`
		Name       string `json:"name"`
		OrderType  string `json:"order_type"`
		Orders     []struct {
			OrderID   string `json:"order_id"`
			OrderInfo struct {
				RemainingQuantity          int         `json:"remaining_quantity"`
				LotSize                    int         `json:"lot_size"`
				Exchange                   string      `json:"exchange"`
				SlOrderPrice               float64     `json:"sl_order_price"`
				Mode                       string      `json:"mode"`
				OrderStatusInfo            string      `json:"order_status_info"`
				NnfID                      int         `json:"nnf_id"`
				Validity                   string      `json:"validity"`
				OrderEntryTime             int         `json:"order_entry_time"`
				SpreadToken                interface{} `json:"spread_token"`
				SquareOff                  bool        `json:"square_off"`
				ExchangeOrderID            string      `json:"exchange_order_id"`
				ClientID                   string      `json:"client_id"`
				Quantity                   int         `json:"quantity"`
				TargetPriceType            string      `json:"target_price_type"`
				AverageTradePrice          int         `json:"average_trade_price"`
				RejectionCode              int         `json:"rejection_code"`
				LegOrderIndicator          interface{} `json:"leg_order_indicator"`
				DisclosedQuantity          int         `json:"disclosed_quantity"`
				TriggerPrice               float64     `json:"trigger_price"`
				LastActivityReference      int         `json:"last_activity_reference"`
				Series                     string      `json:"series"`
				InstrumentToken            int         `json:"instrument_token"`
				Segment                    string      `json:"segment"`
				Price                      float64     `json:"price"`
				OrderType                  string      `json:"order_type"`
				RejectionReason            string      `json:"rejection_reason"`
				MarketProtectionPercentage int         `json:"market_protection_percentage"`
				FilledQuantity             int         `json:"filled_quantity"`
				UserOrderID                int         `json:"user_order_id"`
				TradingSymbol              string      `json:"trading_symbol"`
				SlOrderQuantity            int         `json:"sl_order_quantity"`
				TradePrice                 int         `json:"trade_price"`
				LoginID                    interface{} `json:"login_id"`
				AveragePrice               int         `json:"average_price"`
				Product                    string      `json:"product"`
				SlTriggerPrice             float64     `json:"sl_trigger_price"`
				OrderTag                   string      `json:"order_tag"`
				IsTrailing                 bool        `json:"is_trailing"`
				TrailingStopLoss           float64     `json:"trailing_stop_loss"`
				Deposit                    int         `json:"deposit"`
				OmsOrderID                 string      `json:"oms_order_id"`
				ExecutionType              string      `json:"execution_type"`
				UnderlyingToken            string      `json:"underlying_token"`
				ContractDescription        struct {
				} `json:"contract_description"`
				ExchangeTime   int         `json:"exchange_time"`
				Device         interface{} `json:"device"`
				OrderStatus    interface{} `json:"order_status"`
				SquareOffValue float64     `json:"square_off_value"`
				OrderSide      string      `json:"order_side"`
				StopLossValue  float64     `json:"stop_loss_value"`
				ProCli         string      `json:"pro_cli"`
			} `json:"order_info"`
		} `json:"orders"`
		ProductType string `json:"product_type"`
		SipEligible bool   `json:"sip_eligible"`
		SipEnabled  bool   `json:"sip_enabled"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabEditBasketInstrumentReq struct {
	BasketID  string `json:"basket_id"`
	Name      string `json:"name"`
	OrderID   string `json:"order_id"`
	OrderInfo struct {
		Exchange          string  `json:"exchange"`
		InstrumentToken   int     `json:"instrument_token"`
		ClientID          string  `json:"client_id"`
		OrderType         string  `json:"order_type"`
		Price             float64 `json:"price"`
		Quantity          int     `json:"quantity"`
		DisclosedQuantity int     `json:"disclosed_quantity"`
		Validity          string  `json:"validity"`
		Product           string  `json:"product"`
		TradingSymbol     string  `json:"trading_symbol"`
		OrderSide         string  `json:"order_side"`
		UserOrderID       int     `json:"user_order_id"`
		UnderlyingToken   string  `json:"underlying_token"`
		Series            string  `json:"series"`
		OmsOrderID        string  `json:"oms_order_id"`
		ExchangeOrderID   string  `json:"exchange_order_id"`
		TriggerPrice      float64 `json:"trigger_price"`
		ExecutionType     string  `json:"execution_type"`
	} `json:"order_info"`
}

type TradeLabDeleteBasketInstrumentReq struct {
	BasketID string `json:"basket_id"`
	OrderID  string `json:"order_id"`
	Name     string `json:"name"`
}

type TradeLabRenameBasketReq struct {
	BasketID string `json:"basket_id"`
	Name     string `json:"name"`
}

type TradeLabExecuteBasketReq struct {
	BasketID       string `json:"basket_id"`
	Name           string `json:"name"`
	ExecutionType  string `json:"execution_type"`
	SquareOff      bool   `json:"square_off"`
	ClientID       string `json:"client_id"`
	ExecutionState bool   `json:"execution_state"`
}

type TradeLabExecuteBasketRes struct {
	Data struct {
		Data struct {
			BasketID string `json:"basket_id"`
			Message  string `json:"message"`
		} `json:"data"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabCancelPayoutReq struct {
	Transactions []string `json:"transactions"`
	UserID       string   `json:"user_id"`
	Status       string   `json:"status"`
}

type TradeLabCancelPayoutRes struct {
	Data    string `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabAtomPayoutRequest struct {
	Amount        string `json:"amount"`
	ClientID      string `json:"clientId"`
	Ifsc          string `json:"ifsc"`
	AccountNumber string `json:"accountNumber"`
	BankName      string `json:"bank_name"`
}

type TradeLabAtomPayoutResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  string      `json:"status"`
}

type TradeLabClientTransactionsResponse struct {
	Data    []TradeLabClientTransactionsResponseData `json:"data"`
	Message string                                   `json:"message"`
	Status  string                                   `json:"status"`
}

type TradeLabClientTransactionsResponseData struct {
	AccountName                 string                                          `json:"account_name"`
	Amount                      string                                          `json:"amount"`
	BankName                    string                                          `json:"bank_name"`
	BankTransactionID           string                                          `json:"bank_transaction_id"`
	ClientID                    string                                          `json:"client_id"`
	CreatedAt                   string                                          `json:"created_at"`
	Ifsc                        string                                          `json:"ifsc"`
	MerchantTransactionID       string                                          `json:"merchant_transaction_id"`
	PaymentGatewayTransactionID string                                          `json:"payment_gateway_transaction_id"`
	PaymentGatewayUsername      string                                          `json:"payment_gateway_username"`
	PreviousBalance             int                                             `json:"previous_balance"`
	Status                      string                                          `json:"status"`
	StatusLifeCycle             []ClientTransactionsResponseDataStatusLifeCycle `json:"status_life_cycle"`
	TransactionID               string                                          `json:"transaction_id"`
	TransactionTimestamp        int                                             `json:"transaction_timestamp"`
	TransactionType             string                                          `json:"transaction_type"`
	UpdatedAt                   string                                          `json:"updated_at"`
	UpdatedBy                   string                                          `json:"updated_by"`
	UserID                      string                                          `json:"user_id"`
}

type ClientTransactionsResponseDataStatusLifeCycle struct {
	Status    string `json:"status"`
	UpdatedAt string `json:"updated_at"`
	UpdatedBy string `json:"updated_by"`
}

type TradelabCreateAlertReq struct {
	Exchange         string    `json:"exchange"`
	InstrumentToken  string    `json:"instrument_token"`
	WaitTime         string    `json:"wait_time"`
	Condition        string    `json:"condition"`
	UserSetValues    []float64 `json:"user_set_values"`
	Frequency        string    `json:"frequency"`
	Expiry           int       `json:"expiry"`
	StateAfterExpiry string    `json:"state_after_expiry"`
	UserMessage      string    `json:"user_message"`
}

type TradelabCreateAlertRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AlertID int64 `json:"alert_id"`
	} `json:"data"`
}

type TlAlerts struct {
	WaitTime         int       `json:"wait_time"`
	UserSetValues    []float64 `json:"user_set_values"`
	UserMessage      string    `json:"user_message"`
	TradingSymbol    string    `json:"trading_symbol"`
	Token            string    `json:"token"`
	Status           string    `json:"status"`
	StateAfterExpiry string    `json:"state_after_expiry"`
	ID               int       `json:"id"`
	Frequency        string    `json:"frequency"`
	Expiry           int       `json:"expiry"`
	Exchange         string    `json:"exchange"`
	ConditionType    string    `json:"condition_type"`
	ClientID         string    `json:"client_id"`
}

type TradelabGetAlertsRes struct {
	Status  string     `json:"status"`
	Message string     `json:"message"`
	Data    []TlAlerts `json:"data"`
}

type TradelabPauseAlertsReq struct {
	Status string `json:"status"`
}

type TradelabDeleteAlertsRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
	} `json:"data"`
}

type TradeLabUpdateBasketExecutionStateReq struct {
	BasketID       string `json:"basket_id"`
	Name           string `json:"name"`
	ClientID       string `json:"client_id"`
	ExecutionType  string `json:"execution_type"`
	SquareOff      bool   `json:"square_off"`
	ExecutionState bool   `json:"execution_state"`
}

type TradeLabUpdateBasketExecutionStateRes struct {
	BasketID       string `json:"basket_id"`
	Name           string `json:"name"`
	ClientID       string `json:"client_id"`
	ExecutionType  string `json:"execution_type"`
	SquareOff      bool   `json:"square_off"`
	ExecutionState bool   `json:"execution_state"`
}

// Session Info
type TradeLabSessionInfoResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Result []struct {
		Exchange         string `json:"exchange"`
		SessionName      string `json:"session_name"`
		AMOStartTime     string `json:"amo_start_time"`
		AMOEndTime       string `json:"amo_end_time"`
		IsActive         string `json:"is_active"`
		IsHoliday        string `json:"is_holiday"`
		BufferStartTime  string `json:"buffer_start_time"`
		BufferEndTime    string `json:"buffer_end_time"`
		MarketCloseTime  string `json:"market_close_time"`
		PostClosingStart string `json:"post_closing_start"`
		PostClosingEnd   string `json:"post_closing_end"`
	} `json:"result"`
}

type TradeLabLoginV2Req struct {
	ChannelID     string `json:"channel_id"`
	ChannelSecret string `json:"channel_secret"`
}

type TradelabLoginByEmailOtpReq struct {
	ChannelID string `json:"channel_id"`
}

type TradeLabLoginV2Response struct {
	Data struct {
		Alert          string `json:"alert"`
		AuthToken      string `json:"auth_token"`
		CheckPan       bool   `json:"check_pan"`
		LoginID        string `json:"login_id"`
		Name           string `json:"name"`
		ReferenceToken string `json:"reference_token"`
		ResetPassword  bool   `json:"reset_password"`
		ResetTwoFa     bool   `json:"reset_two_fa"`
		Twofa          struct {
			Questions []struct {
				Question   string `json:"question"`
				QuestionID int    `json:"question_id"`
			} `json:"questions"`
			TwofaToken string `json:"twofa_token"`
			Type       string `json:"type"`
		} `json:"twofa"`
		TwofaEnabled bool `json:"twofa_enabled"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabValidateTwofaV2Req struct {
	LoginID    string                  `json:"login_id"`
	Twofa      []TradeLabTwoFaQuestion `json:"twofa"`
	TwofaToken string                  `json:"twofa_token"`
	Type       string                  `json:"type"`
	DeviceType string                  `json:"device_type"`
}

type TradeLabValidateTwofaV2Res struct {
	Data struct {
		AuthToken     string `json:"auth_token"`
		ResetPassword bool   `json:"reset_password"`
		ResetTwoFa    bool   `json:"reset_two_fa"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabSetupTotpV2Req struct {
	ClientID string `json:"client_id"`
}

type TradeLabSetupTotpV2Res struct {
	Data    string `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabChooseTwofaV2Req struct {
	LoginID   string `json:"login_id"`
	TwofaType string `json:"twofa_type"`
	Totp      string `json:"totp"`
}

type TradeLabChooseTwofaV2Res struct {
	Data struct {
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabReturnOnInvestmentRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Roi []struct {
			Volume           int     `json:"volume"`
			ReturnPercent    float64 `json:"return_percent"`
			PercentageChange float64 `json:"percentage_change"`
			Ltp              float64 `json:"ltp"`
			InstrumentToken  int     `json:"instrument_token"`
			Exchange         string  `json:"exchange"`
			DaysChange       int     `json:"days_change"`
			ClosePrice       float64 `json:"close_price"`
			Change           float64 `json:"change"`
			TradingSymbol    string  `json:"trading_symbol"`
		} `json:"roi"`
	} `json:"data"`
}

type TradeLabAllBankAccountsUpdatedRes struct {
	Data struct {
		BankAccounts []struct {
			AccountType       string `json:"account_type"`
			BankAccountNumber string `json:"bank_account_number"`
			BankBranchName    string `json:"bank_branch_name"`
			BankID            string `json:"bank_id"`
			BankName          string `json:"bank_name"`
			City              string `json:"city"`
			Ifsc              string `json:"ifsc"`
			PanNumber         string `json:"pan_number"`
			State             string `json:"state"`
		} `json:"bank_accounts"`
		ClientId string `json:"client_id"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradelabFetchAdminMessageRes struct {
	Data struct {
		Updates []TradelabFetchAdminMessageDataUpdates `json:"updates"`
	} `json:"data"`
	Status string `json:"status"`
}

type TradelabFetchAdminMessageDataUpdates struct {
	Type            string `json:"type"`
	UpdateEntryTime string `json:"update_entry_time"`
	Message         string `json:"message"`
	Platform        string `json:"platform"`
	Title           string `json:"title"`
}

type TradelabNotificationUpdates struct {
	Data struct {
		Updates []struct {
			Type             string    `json:"type"`
			UpdateEntryTime  string    `json:"update_entry_time"`
			UpdateID         string    `json:"update_id"`
			AlertID          int       `json:"alert_id"`
			Condition        string    `json:"condition"`
			Exchange         string    `json:"exchange"`
			Expiry           int       `json:"expiry"`
			Frequency        string    `json:"frequency"`
			GeneratedAt      int       `json:"generated_at"`
			InstrumentCode   string    `json:"instrument_code"`
			LotSize          string    `json:"lot_size"`
			NewValue         float64   `json:"new_value"`
			StateAfterExpiry string    `json:"state_after_expiry"`
			TradingSymbol    string    `json:"trading_symbol"`
			UserMessage      string    `json:"user_message"`
			UserSetValues    []float64 `json:"user_set_values"`
		} `json:"updates"`
	} `json:"data"`
	Status string `json:"status"`
}

type TradeLabFuturesChain struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Result []struct {
		ExpiryDate string `json:"expiry_date"`
		Strikes    struct {
			Token         string `json:"token"`
			Exchange      string `json:"exchange"`
			Company       string `json:"company"`
			Symbol        string `json:"symbol"`
			TradingSymbol string `json:"trading_symbol"`
			DisplayName   string `json:"display_name"`
			ExpiryRaw     string `json:"ExpiryRaw"`
			ClosePrice    string `json:"close_price"`
		} `json:"strikes"`
	} `json:"result"`
}

type TradeLabForgetTotpV2Req struct {
	LoginID string `json:"login_id"`
	Pan     string `json:"pan"`
}

type TradeLabForgetTotpV2Res struct {
	Data    string `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabValidateLoginOtpV2Req struct {
	ReferenceToken string `json:"reference_token"`
	Otp            string `json:"otp"`
}

type TradeLabValidateLoginOtpV2Res struct {
	Data struct {
		Alert          string `json:"alert"`
		AuthToken      string `json:"auth_token"`
		CheckPan       bool   `json:"check_pan"`
		LoginID        string `json:"login_id"`
		Name           string `json:"name"`
		ReferenceToken string `json:"reference_token"`
		ResetPassword  bool   `json:"reset_password"`
		ResetTwoFa     bool   `json:"reset_two_fa"`
		Twofa          struct {
			Questions []struct {
				Question   string `json:"question"`
				QuestionID int    `json:"question_id"`
			} `json:"questions"`
			TwofaToken string `json:"twofa_token"`
			Type       string `json:"type"`
		} `json:"twofa"`
		TwofaEnabled bool `json:"twofa_enabled"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabSetupBiometricV2Req struct {
	ClientID    string `json:"client_id"`
	Fingerprint string `json:"fingerprint"`
}

type TradeLabSetupBiometricV2Res struct {
	Data struct {
		TwofaToken string `json:"twofa_token"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabDisableBiometricV2Res struct {
	Data struct {
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradeLabInstrument struct {
	InstrumentToken int    `json:"instrument_token"`
	Exchange        string `json:"exchange"`
	Total           int    `json:"total"`
	Authorized      int    `json:"authorized"`
}

type TradeLabEdisReq struct {
	ClientID    string               `json:"client_id"`
	Instruments []TradeLabInstrument `json:"instruments"`
	RequestType string               `json:"request_type"`
}

type TradeLabEdisDataRes struct {
	Depository    string `json:"depository"`
	DpId          string `json:"dp_id"`
	EncryptedDtls string `json:"encrypted_dtls"`
	RequestId     string `json:"request_id"`
	Version       string `json:"version"`
}

type TradeLabEdisRes struct {
	Data    TradeLabEdisDataRes `json:"data"`
	Html    bool                `json:"html"`
	Message string              `json:"message"`
	Status  string              `json:"status"`
}

type TradeLabIsin struct {
	IsinName string `json:"isin_name"`
	Isin     string `json:"isin"`
	Quantity string `json:"quantity"`
	Price    string `json:"price"`
}

type TradeLabEpledgeReq struct {
	Depository  string         `json:"depository"`
	ClientID    string         `json:"client_id"`
	Exchange    string         `json:"exchange"`
	BoId        string         `json:"bo_id"`
	Segment     string         `json:"segment"`
	IsinDetails []TradeLabIsin `json:"isin_details"`
}

type TradeLabEpledgeDataRes struct {
	DpId       string `json:"dpid"`
	PledgedTls string `json:"pledgedtls"`
	ReqId      string `json:"reqid"`
	Version    string `json:"version"`
}

type TradeLabEpledgeRes struct {
	Data    TradeLabEpledgeDataRes `json:"data"`
	Message string                 `json:"message"`
	Status  string                 `json:"status"`
}

type TradeLabTpinReq struct {
	Boid    string `json:"boid"`
	Pan     string `json:"pan"`
	ReqFlag string `json:"ReqFlag"`
	ReqTime string `json:"ReqTime"`
}

type TradeLabTpinRes struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type TradeLabUnPledgeReq struct {
	ClientID        string `json:"client_id"`
	Exchange        string `json:"exchange"`
	Segment         string `json:"segment"`
	Isin            string `json:"isin"`
	Quantity        int64  `json:"quantity"`
	TransactionType string `json:"transaction_type"`
}

type TlError struct {
	Code    string `json:"code"`
	Message int    `json:"message"`
}

type TradeLabUnpledgeRes struct {
	Error   TlError     `json:"error"`
	Result  string      `json:"result"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type TradeLabUnblockUserReq struct {
	LoginID string `json:"login_id"`
	Pan     string `json:"pan"`
}

type TradeLabCreateAppReq struct {
	AppName      string   `json:"app_name"`
	RedirectUris []string `json:"redirect_uris"`
	Scope        string   `json:"scope"`
	GrantTypes   []string `json:"grant_types"`
	Owner        string   `json:"owner"`
}

type TradeLabCreateAppRes struct {
	Data    TradeLabCreateAppData `json:"data"`
	Message string                `json:"message"`
	Status  string                `json:"status"`
}

type TradeLabCreateAppData struct {
	AppID              string   `json:"app_id"`
	AppName            string   `json:"app_name"`
	AppOwner           string   `json:"app_owner"`
	AppSecret          string   `json:"app_secret"`
	AppSecretExpiresAt int      `json:"app_secret_expires_at"`
	GrantTypes         []string `json:"grant_types"`
	RedirectUris       []string `json:"redirect_uris"`
	Scope              string   `json:"scope"`
}

type TradelabFetchAppResponse struct {
	AppID        string   `json:"app_id"`
	AppName      string   `json:"app_name"`
	RedirectURIs []string `json:"redirect_uris"`
	Scope        string   `json:"scope"`
}

type TradelabFetchAppDetailsResponse struct {
	Data []TradelabFetchAppResponse `json:"data"`
}

type TradelabDeleteResponse struct {
	Data    string `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TradelabOauthAccessTokenReq struct {
	GrantType   string `json:"grant_type"`
	Code        string `json:"code"`
	RedirectUri string `json:"redirect_uri"`
}

type TradelabOauthAccessTokenRes struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type TLGetSipResponse struct {
	Data    []TLGetSipResponseData `json:"data"`
	Message string                 `json:"message"`
	Status  string                 `json:"status"`
}

type TLGetSipResponseData struct {
	ActionType      string      `json:"action_type"`
	Baskets         []string    `json:"baskets"`
	ClientID        string      `json:"client_id"`
	Count           int         `json:"count"`
	CreatedAt       string      `json:"created_at"`
	Id              string      `json:"id"`
	LoginID         string      `json:"login_id"`
	Name            string      `json:"name"`
	NextExecutionAt string      `json:"next_execution_at"`
	Order           interface{} `json:"order"`
	RejectCode      int         `json:"reject_code"`
	RejectReason    string      `json:"reject_reason"`
	RemainingDays   int         `json:"remaining_days"`
	Schedules       Schedule    `json:"schedules"`
	Source          string      `json:"source"`
	Status          string      `json:"status"`
	Tags            []string    `json:"tags"`
	Type            string      `json:"type"`
	UpdatedAt       string      `json:"updated_at"`
}

type Schedule struct {
	Frequency  string      `json:"frequency"`
	TimeSquare []TimeEntry `json:"time_square"`
}

type TimeEntry struct {
	Day     string `json:"day"`
	Expiry  string `json:"expiry"`
	Time    string `json:"time"`
	Weekday string `json:"weekday"`
}

type TLUpdateSipOrderResponse struct {
	Data struct {
		Id string `json:"id"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TLUpdateSipStatusResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	Data UpdateSipStatusResData `json:"data"`
}

type UpdateSipStatusResData struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type TLIcebergOrderReq struct {
	Order struct {
		Exchange                   string `json:"exchange"`
		InstrumentToken            string `json:"instrument_token"`
		ClientId                   string `json:"client_id"`
		OrderType                  string `json:"order_type"`
		Price                      string `json:"price"`
		Quantity                   string `json:"quantity"`
		DisclosedQuantity          int    `json:"disclosed_quantity"`
		Validity                   string `json:"validity"`
		Product                    string `json:"product"`
		NoOfLegs                   string `json:"no_of_legs"`
		GttPrice                   string `json:"gtt_price"`
		OrderSide                  string `json:"order_side"`
		Device                     string `json:"device"`
		UserOrderId                int    `json:"user_order_id"`
		TriggerPrice               int    `json:"trigger_price"`
		ExecutionType              string `json:"execution_type"`
		MarketProtectionPercentage int    `json:"market_protection_percentage"`
	} `json:"order"`
}

type TLModifyIcebergOrderReq struct {
	Exchange                   string  `json:"exchange"`
	InstrumentToken            int     `json:"instrument_token"`
	ClientId                   string  `json:"client_id"`
	OrderType                  string  `json:"order_type"`
	Price                      string  `json:"price"`
	Quantity                   int     `json:"quantity"`
	DisclosedQuantity          int     `json:"disclosed_quantity"`
	Validity                   string  `json:"validity"`
	Product                    string  `json:"product"`
	GttPrice                   string  `json:"gtt_price"`
	OmsOrderId                 string  `json:"oms_order_id"`
	ExchangeOrderId            string  `json:"exchange_order_id"`
	FilledQuantity             int     `json:"filled_quantity"`
	RemainingQuantity          int     `json:"remaining_quantity"`
	LastActivityReference      int64   `json:"last_activity_reference"`
	TriggerPrice               float64 `json:"trigger_price"`
	ExecutionType              string  `json:"execution_type"`
	MarketProtectionPercentage float64 `json:"market_protection_percentage"`
}

type TLMTFCTDReq struct {
	MtfCtdValues []MtfCtdValues `json:"mtf_ctd_values"`
	UserType     string         `json:"user_type"`
	LoginID      string         `json:"login_id"`
	ClientID     string         `json:"client_id"`
}

type MtfCtdValues struct {
	Isin        string `json:"isin"`
	CtdQuantity int    `json:"ctd_quantity"`
}

type TLMTFCTDRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		SuccessCount int `json:"success_count"`
		FailureCount int `json:"failure_count"`
		TotalCount   int `json:"total_count"`
	} `json:"data"`
}
