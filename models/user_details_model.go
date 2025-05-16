package models

type GetAllBankAccountsReq struct {
	UserId string `json:"userId"`
}

type MongoAllBankAccountDetails struct {
	UserId          string                    `json:"userid"`
	AllBankAccounts []MongoBankAccountDetails `json:"allBankAccounts" mask:"struct"`
}

type MongoBankAccountDetails struct {
	UserId              string                               `json:"userid"`
	Verified            bool                                 `json:"verified"`
	IsPrimary           bool                                 `json:"isPrimary"`
	Rejection           string                               `json:"Rejection"`
	BankAccountMetadata BankAccountDetails                   `json:"bankAccountMetadata" mask:"struct"`
	DigioBankMetadata   DigioBankVerificationSuccessResponse `json:"digioBankMetadata"`
}

type BankAccountDetails struct {
	UserId          string `json:"userid"`
	UserAccNumber   string `json:"beneficiary_account_no" example:"123110023204445" mask:"id"`
	UserIfscCode    string `json:"beneficiary_ifsc" example:"UTIB0000888" mask:"id"`
	UserFullName    string `json:"beneficiary_name" example:"Dummy Name"`
	UserBankName    string `json:"bank_name" example:"Axis Bank"`
	UserBankAddress string `json:"bank_address" example:"Okhla Branch"`
	SetPrimary      bool   `json:"setPrimary"`
}

type DigioBankVerificationSuccessResponse struct {
	Id                         string `json:"id"`
	Verified                   bool   `json:"verified"`
	Verified_at                string `json:"verified_at"`
	Beneficiary_name_with_bank string `json:"beneficiary_name_with_bank"`
	Fuzzy_match_score          int    `json:"fuzzy_match_score"` // Number returned only of account is verified and name match is expected as per request.
}

type GetAllBankAccountsUpdatedReq struct {
	ClientId string `json:"clientId"`
}

type GetAllBankAccountsUpdatedRes struct {
	BankAccounts []BankAccount `json:"bankAccounts" mask:"struct"`
	ClientId     string        `json:"clientId"`
}

type BankAccount struct {
	AccountType       string `json:"accountType"`
	BankAccountNumber string `json:"bankAccountNumber" mask:"id"`
	BankBranchName    string `json:"bankBranchName"`
	BankID            string `json:"bankId"`
	BankName          string `json:"bankName"`
	City              string `json:"city"`
	Ifsc              string `json:"ifsc" mask:"id"`
	PanNumber         string `json:"panNumber" mask:"id"`
	State             string `json:"state"`
}

type GetUserIdReq struct {
	Id     string `json:"id"`
	IdType string `json:"idType" example:"emailid,clientid"`
}

type GetUserIdRes struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
}

type UnblockUserReq struct {
	LoginID string `json:"loginId" binding:"required"`
	Pan     string `json:"pan" binding:"required"`
}

type UnblockUserRes struct {
	Data struct {
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type UserNotificationsReq struct {
	ClientId string `json:"clientId"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

type MongoNotificationStore struct {
	ClientId            string              `json:"clientId"`
	StoredAt            int64               `json:"storedAt"`
	TLOrderUpdatePacket TLOrderUpdatePacket `json:"tlOrderUpdatePacket"`
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

type MongoClientsDetails struct {
	KycUserId  string `json:"kycUserId"`
	ClientID   string `json:"clientId"`
	ClientName string `json:"clientName"`
	IntroDate  string `json:"introDate"`
	DOB        string `json:"dob"`
	MobileNo   string `json:"mobileNo"`
	PhNo       string `json:"phNo"`
	Email      string `json:"email"`
	PAN        string `json:"pan"`
	DPID       string `json:"dpId"`
	BOID       string `json:"boId"`
	ClientType string `json:"clientType"`
}

type UnblockUserV2Req struct {
	ClientID string `json:"ClientId"`
	EmailID  string `json:"emailId"`
	Pan      string `json:"pan" mask:"id"`
}

type GetClientStatusRes struct {
	UserStatus int    `json:"userStatus"`
	KycUserId  string `json:"kycUserId"`
}
