package models

type TradeConfirmationDateRangeReq struct {
	UserID   string `json:"userID"`
	DFDateFr string `json:"dFDateFr"`
	DFDateTo string `json:"dFDateTo"`
}

type TradeConfirmationDateRangeRes struct {
	AllTradeConfirmationDateRange []TradeConfirmationDateRange `json:"allTradeConfirmationDateRange"`
}

type TradeConfirmationDateRange struct {
	SellAvgRate string `json:"sellAvgRate"`
	Scripname   string `json:"scripname"`
	BuyValue    string `json:"buyValue"`
	SellValue   string `json:"sellValue"`
	SellQty     string `json:"sellQty"`
	NetQty      string `json:"netQty"`
	Segment     string `json:"segment"`
	NetValue    string `json:"netValue"`
	Buyqty      string `json:"buyqty"`
	BuyAvgRate  string `json:"buyAvgRate"`
	TradeDate   string `json:"tradeDate"`
	NetAvgRate  string `json:"netAvgRate"`
}

type GetBillDetailsCdslReq struct {
	UserID   string `json:"userID"`
	DFDateFr string `json:"dFDateFr"`
	DFDateTo string `json:"dFDateTo"`
}

type GetBillDetailsCdslRes struct {
	GetBillDetailsCdsl []GetBillDetailsCdsl `json:"shilpiGetBillDetailsCdsl"`
}

type GetBillDetailsCdsl struct {
	Charges        string `json:"charges"`
	TrxDate        string `json:"trxDate"`
	Qty            string `json:"qty"`
	Gst            string `json:"gst"`
	InstrumentName string `json:"instrumentName"`
	Isincode       string `json:"isincode"`
	ChargesDetails string `json:"chargesDetails"`
	TotalCharges   string `json:"totalCharges"`
}

type LongTermShortTermReq struct {
	UserID string `json:"userID" validate:"required"`
}

type LongTermShortTermRes struct {
	LongTermShortTerm []LongTermShortTerm `json:"shilpiLongTermShortTerm"`
}

type LongTermShortTerm struct {
	Scripname string `json:"scripname"`
	BuyRate   string `json:"buyRate"`
	SellQty   string `json:"sellQty"`
	Jobbing   string `json:"jobbing"`
	SellRate  string `json:"sellRate"`
	ShortTerm string `json:"shortTerm"`
	Buyqty    string `json:"buyqty"`
	LongTerm  string `json:"longTerm"`
	Scripcd   string `json:"scripcd"`
	Isin      string `json:"isin"`
}

type FetchProfileReq struct {
	UserID string `json:"userID" validate:"required"`
}

type FetchProfileRes struct {
	FetchProfile []FetchProfile `json:"shilpiFetchProfile"`
}

type FetchProfile struct {
	Emailno          string `json:"emailno"`
	Phnos            string `json:"phnos" mask:"id"`
	Occupationcode   string `json:"occupationcode"`
	Poaflag          string `json:"poaflag"`
	Gender           string `json:"gender"`
	City             string `json:"city"`
	Annualincomedate string `json:"annualincomedate"`
	Panno            string `json:"panno" mask:"id"`
	Annualincomecode string `json:"annualincomecode"`
	Dpid             string `json:"dpid"`
	Mobileno         string `json:"mobileno" mask:"id"`
	Micrcode         string `json:"micrcode"`
	Branch           string `json:"branch"`
	Ckycregnno       string `json:"ckycregnno"`
	Bankcity         string `json:"bankcity"`
	Groupclient      string `json:"groupclient"`
	Accountstatus    string `json:"accountstatus"`
	Bankadd3         string `json:"bankadd3"`
	Bankadd1         string `json:"bankadd1"`
	Bankadd2         string `json:"bankadd2"`
	State            string `json:"state"`
	Add1             string `json:"add1"`
	Activeexchange   string `json:"activeexchange"`
	Add4             string `json:"add4"`
	Branchname       string `json:"branchname"`
	Enableotradeing  string `json:"enableotradeing"`
	Pincode          string `json:"pincode"`
	Introdate        string `json:"introdate"`
	Motherfullname   string `json:"motherfullname"`
	Bankacno         string `json:"bankacno" mask:"id"`
	Ifsccode         string `json:"ifsccode" mask:"id"`
	Bankactype       string `json:"bankactype"`
	Dpaccountno      string `json:"dpaccountno" mask:"id"`
	Fathername       string `json:"fathername"`
	Nationality      string `json:"nationality"`
	Dob              string `json:"dob"`
	Name             string `json:"name"`
	Bankname         string `json:"bankname"`
	Kraupload        string `json:"kraupload"`
	Polexposueperson string `json:"polexposueperson"`
	Fatca            string `json:"fatca"`
	Married          string `json:"married"`
}

type TradeConfirmationOnDateReq struct {
	UserId    string `json:"userId"`
	TradeDate string `json:"tradeDate"`
}

type TradeConfirmationOnDateRes struct {
	TradeConfirmationOnDate []TradeConfirmationOnDate `json:"tradeConfirmationOnDate"`
}

type TradeConfirmationOnDate struct {
	SellAvgRate string `json:"SellAvgRate"`
	Scripname   string `json:"Scripname"`
	BuyValue    string `json:"BuyValue"`
	SellValue   string `json:"SellValue"`
	SellQty     string `json:"SellQty"`
	NetQty      string `json:"NetQty"`
	Segment     string `json:"segment"`
	NetValue    string `json:"NetValue"`
	Buyqty      string `json:"Buyqty"`
	BuyAvgRate  string `json:"buyAvgRate"`
	NetAvgRate  string `json:"NetAvgRate"`
}

type OpenPositionsReq struct {
	UserId string `json:"userId"`
	DsFlag string `json:"dsFlag"`
}

type OpenPositionsRes struct {
	OpenPositions []OpenPositions `json:"shilpiOpenPositions"`
}

type OpenPositions struct {
	Segment     string `json:"segment"`
	Exchange    string `json:"exchange"`
	Contract    string `json:"contract"`
	BuyQty      string `json:"buyQty"`
	BuyValue    string `json:"buyValue"`
	SellQty     string `json:"sellQty"`
	SellValue   string `json:"sellValue"`
	NetQty      string `json:"netQty"`
	NetValue    string `json:"netValue"`
	ClosingRate string `json:"closingRate"`
	PAndL       string `json:"pAndL"`
	Exposure    string `json:"exposure"`
}

type GetHoldingReq struct {
	UserId string `json:"userId" validate:"required"`
}

type GetHoldingRes struct {
	GetHolding []GetHolding `json:"shilpiGetHolding" mask:"struct"`
}

type GetHolding struct {
	MobileNo   string `json:"mobileNo" mask:"id"`
	LoginId    string `json:"loginId"`
	IsinCode   string `json:"isinCode"`
	IsinName   string `json:"isinName"`
	Holding    string `json:"holding"`
	ClosePrice string `json:"closePrice"`
	Valuation  string `json:"valuation"`
	ScripCode  string `json:"scripCode"`
}

type GetMarginOnDateReq struct {
	UserId     string `json:"userId"`
	MarginDate string `json:"marginDate"`
}

type GetMarginOnDateRes struct {
	GetMarginOnDate []GetMarginOnDate `json:"getMarginOnDate"`
}

type GetMarginOnDate struct {
	TotalMargin     string `json:"totalMargin"`
	Exposure        string `json:"exposure"`
	PeakDeposit     string `json:"peakDeposit"`
	VarMargin       string `json:"varMargin"`
	Deposit         string `json:"deposit"`
	ShortMargin     string `json:"shortMargin"`
	TotalPeakmargin string `json:"totalPeakMargin"`
	PeakMarginShort string `json:"peakMarginShort"`
	Span            string `json:"span"`
}

type FinancialLedgerBalanceOnDateReq struct {
	UserId   string `json:"userId"`
	AsOnDate string `json:"asOnDate"`
}

type FinancialLedgerBalanceOnDateRes struct {
	FinancialLedgerBalanceOnDate []FinancialLedgerBalanceOnDate `json:"shilpiFinancialLedgerBalanceOnDate"`
}

type FinancialLedgerBalanceOnDate struct {
	FinBalance string `json:"finBalance"`
}

type GetFinancialReq struct {
	UserId      string `json:"userId"`
	NoOfEntries string `json:"noOfEntries"`
}

type GetFinancialRes struct {
	GetFinancial []GetFinancial `json:"shilpiGetFinancial"`
}

type GetFinancial struct {
	Date      string `json:"date"`
	Narr      string `json:"narr"`
	Segment   string `json:"segment"`
	Exchange  string `json:"exchange"`
	Debit     string `json:"debit"`
	Credit    string `json:"credit"`
	Payrefno  string `json:"payrefno"`
	ValueDate string `json:"valueDate"`
}

type UserDetailsData struct {
	ClientID  string `json:"clientId"`
	Name      string `json:"name"`
	EmailID   string `json:"emailId"`
	BoID      string `json:"boId"`
	PanNumber string `json:"panNumber" mask:"id"`
}

type ViewDPchargesRes struct {
	UserDetails   ProfileDataResp `json:"userDetails"`
	Charges       float64         `json:"charges"`
	GST           float64         `json:"gst"`
	TotalCharges  float64         `json:"totalCharges"`
	DPChargesList []DPcharges     `json:"DPChargesList"`
}

type DPcharges struct {
	InstrumentName string `json:"instrumentName"`
	ChargesDetails string `json:"chargesDetails"`
	Qty            string `json:"qty"`
	Charges        string `json:"charges"`
	Gst            string `json:"gst"`
	TotalCharges   string `json:"totalCharges"`
}

type DPChargesReq struct {
	UserID   string `json:"userID"`
	DFDateFr string `json:"dFDateFr"`
	DFDateTo string `json:"dFDateTo"`
}

type DownloadDPChargesRes struct {
	DownloadUrl string `json:"downloadUrl"`
}

type ScripWiseCostingRes struct {
	UserDetails       ProfileDataResp    `json:"userDetails" mask:"struct"`
	TotalBrokerage    float64            `json:"totalBrokerage"`
	TotalGST          float64            `json:"totalGST"`
	TotalSEBITax      float64            `json:"totalSEBITax"`
	TotalSTT          float64            `json:"totalSTT"`
	TotalTurnCharges  float64            `json:"totalTurnCharges"`
	TotalStampDuty    float64            `json:"totalStampDuty"`
	TotalOtherCharges float64            `json:"totalOtherCharges"`
	TotalCharges      float64            `json:"totalCharges"`
	ScripWiseCosting  []ScripWiseCosting `json:"scripWiseCosting"`
}

type ScripWiseCosting struct {
	NetValue       float64 `json:"netvalue"`
	Brokerage      float64 `json:"brokerage"`
	NetQty         float64 `json:"netqty"`
	OrderNo        string  `json:"orderno"`
	OtherCharges   float64 `json:"othercharges"`
	ScripCode      string  `json:"scripcode"`
	GST            float64 `json:"gst"`
	Stamp          float64 `json:"stamp"`
	Price          float64 `json:"price"`
	TrxDate        string  `json:"trxdate"`
	ISINCode       string  `json:"isincode"`
	SEBIFee        float64 `json:"sebifee"`
	NetAmt         float64 `json:"netamt"`
	SellValue      float64 `json:"sellvalue"`
	STT            float64 `json:"stt"`
	BuyValue       float64 `json:"buyvalue"`
	ExchClg        float64 `json:"exchclg"`
	TurnTax        float64 `json:"turntax"`
	Exchange       string  `json:"exchange"`
	BrokType       string  `json:"broktype"`
	ScripName      string  `json:"scrip name"`
	BuySellType    string  `json:"buySellType"`
	Quantity       float64 `json:"quantity"`
	UnixTimeFormat int64   `json:"unixTime"`
}

type TradebookReq struct {
	UserID   string `json:"userID"`
	DFDateFr string `json:"dFDateFr"`
	DFDateTo string `json:"dFDateTo"`
}

type DownloadTradebookRes struct {
	DownloadUrl string `json:"downloadUrl"`
}

type FinancialLedgerRes struct {
	UserDetails     ProfileDataResp       `json:"userDetails"`
	OpeningBalance  float64               `json:"openingbalance"`
	Inflow          float64               `json:"inflow"`
	Outflow         float64               `json:"outflow"`
	FundsReceived   float64               `json:"fundsReceived"`
	FundsWithdrawn  float64               `json:"fundsWithdrawn"`
	ClosingBalance  float64               `json:"closingBalance"`
	FinancialLedger []FinancialLedgerData `json:"financialLedger"`
}

type FinancialLedgerData struct {
	TransactionDate    string  `json:"TransactionDate"`
	TransactionDetails string  `json:"transactionDetails"`
	Segment            string  `json:"segment"`
	Exchange           string  `json:"exchange"`
	Debit              float64 `json:"debit"`
	Credit             float64 `json:"credit"`
	SettlementNumber   string  `json:"SettlementNumber"`
	NetBalance         float64 `json:"netBalance"`
	SettlementDate     string  `json:"SettlementDate"`
	UnixTimeFormat     int64   `json:"unixTime"`
}

type GetFinancialLedgerDataReq struct {
	UserID   string `json:"userID"`
	DFDateFr string `json:"dFDateFr"`
	DFDateTo string `json:"dFDateTo"`
}

type LedgerReq struct {
	UserID   string `json:"userID"`
	DFDateFr string `json:"dFDateFr"`
	DFDateTo string `json:"dFDateTo"`
}

type DownloadLedgerRes struct {
	DownloadUrl string `json:"downloadUrl"`
}

type OpenPositionReq struct {
	UserID string `json:"userID"`
	Dsflag string `json:"dsflag" example:"D,S"  validate:"oneof=D S"`
}

type OpenPositionRes struct {
	UserDetails         ProfileDataResp    `json:"userDetails"`
	EquityFutureMTM     float64            `json:"equityFutureMTM"`
	CurrencyFutureMTM   float64            `json:"currencyFutureMTM"`
	CommodityFutureMTM  float64            `json:"commodityFutureMTM"`
	EquityOptionMTM     float64            `json:"equityOptionMTM"`
	CurrencyOptionMTM   float64            `json:"currencyOptionMTM"`
	CommodityOptionMTM  float64            `json:"commodityOptionMTM"`
	EquityDerivative    []OpenPositionData `json:"equityDerivative"`
	CurrencyDerivative  []OpenPositionData `json:"currencyDerivative"`
	CommodityDerivative []OpenPositionData `json:"commodityDerivative"`
}

type OpenPositionData struct {
	ScripName              string  `json:"scripName"`
	InstrumentType         string  `json:"instrumentType"`
	OptionType             string  `json:"optionType"`
	BuySellType            string  `json:"buySellType"`
	StrikePrice            float64 `json:"strikePrice"`
	ExpiryDate             string  `json:"expiryDate"`
	Exchange               string  `json:"exchange"`
	OpenQuantity           float64 `json:"openQuantity"`
	AveragePrice           float64 `json:"averagePrice"`
	ClosingPrice           float64 `json:"closingPrice"`
	UnrealisedProfitOrLoss float64 `json:"unrealisedProfitOrLoss"`
}

type DownloadOpenPositionRes struct {
	DownloadUrl string `json:"downloadUrl"`
}

type GetFONetPositionDataReq struct {
	UserID   string `json:"userID"`
	DFDateFr string `json:"dFDateFr"`
	DFDateTo string `json:"dFDateTo"`
}

type FONetPositionSummaryData struct {
	DateRange     string  `json:"dateRange"`
	Charges       float64 `json:"charges"`
	RealisedPNL   float64 `json:"realisedPNL"`
	UnRealisedPNL float64 `json:"unRealisedPNL"`
	NetPNL        float64 `json:"netPNL"`
}

type FONetPositionChargesData struct {
	Brockerage                 float64 `json:"brockerage"`
	ExchangeTransactionCharges float64 `json:"exchangeTransactionCharges"`
	ClearingCharges            float64 `json:"clearingCharges"`
	IntegratedGST              float64 `json:"integratedGST"`
	SecuritiesTransactionTax   float64 `json:"securitiesTransactionTax"`
	SEBIFees                   float64 `json:"SEBIFees"`
	StampDuty                  float64 `json:"stampDuty"`
	TotalCharges               float64 `json:"totalCharges"`
}

type FONetPositionDetailsData struct {
	Symbol               string  `json:"contract"`
	InstrumentType       string  `json:"instrumentType"`
	OptionType           string  `json:"optionType"`
	StrikePrice          float64 `json:"strikePrice"`
	ExpiryDate           string  `json:"expiryDate"`
	Quantity             int     `json:"quantity"`
	BuyValue             float64 `json:"buyValue"`
	BuyPrice             float64 `json:"buyPrice"`
	SellValue            float64 `json:"sellValue"`
	SellPrice            float64 `json:"sellPrice"`
	RealizedPNL          float64 `json:"realizedPNL"`
	PreviousClosingPrice float64 `json:"previousClosingPrice"`
	OpenQuantity         int     `json:"openQuantity"`
	OpenValue            float64 `json:"openValue"`
	UnrealizedPNL        float64 `json:"unrealizedPNL"`
}

type FONetPositionRes struct {
	UserDetails          ProfileDataResp            `json:"userDetails"`
	Summary              FONetPositionSummaryData   `json:"summary"`
	ChargesDetails       FONetPositionChargesData   `json:"chargesDetails"`
	FONetPositionDetails []FONetPositionDetailsData `josn:"FONetPositionDetails"`
}

type FnoPnlReq struct {
	UserID   string `json:"userID"`
	DFDateFr string `json:"dFDateFr"`
	DFDateTo string `json:"dFDateTo"`
}

type DownloadFnoPnlRes struct {
	DownloadUrl string `json:"downloadUrl"`
}

type GetHoldingFinancialDataReq struct {
	UserID string `json:"userID"`
}

type GetHoldingFinancialData struct {
	Isin                        string  `json:"isin"`
	Instrument                  string  `json:"instrument"`
	PledgedQty                  float64 `json:"pledgedQty"`
	FreeQty                     float64 `json:"freeQty"`
	TotalQty                    float64 `json:"totalQty"`
	TotalPledgedValue           float64 `json:"totalPledgedValue"`
	HaircutPercentage           float64 `json:"haircutPercentage"`
	MarginAvailableAfterHaircut float64 `json:"marginAvailableAfterHaircut"`
	AvgBuyPrice                 float64 `json:"avgBuyPrice"`
	ClosingPrice                float64 `json:"closingPrice"`
	InvestmentValue             float64 `json:"investmentValue"`
	CurrentValue                float64 `json:"currentValue"`
	ContributionPercentage      float64 `json:"contributionPercentage"`
	UnrealizedProfitLoss        float64 `json:"unrealizedProfitLoss"`
	NetChange                   float64 `json:"netChange"`
}

type HoldingSummaryData struct {
	InvestedValue                float64 `json:"investedValue"`
	CurrentValue                 float64 `json:"currentValue"`
	UnrealisedPNL                float64 `json:"unrealisedPNL"`
	TotalPledgeValue             float64 `json:"totalPledgeValue"`
	TotalMarginValueAfterHaircut float64 `json:"totalMarginValueAfterHaircut"`
}

type GetHoldingFinancialDataRes struct {
	UserDetails          ProfileDataResp           `json:"userDetails"`
	HoldingSummary       HoldingSummaryData        `json:"holdingSummary"`
	HoldingFinancialData []GetHoldingFinancialData `json:"holdingFinancialData"`
}

type DownloadHoldingFinancialRes struct {
	DownloadUrl string `json:"downloadUrl"`
}

type CommodityTradebookReq struct {
	UserID   string `json:"userID"`
	DFDateFr string `json:"dFDateFr"`
	DFDateTo string `json:"dFDateTo"`
}

type CommodityTransactionData struct {
	Symbol           string  `json:"symbol"`
	ExpiryDate       string  `json:"expirydate"`
	OptionType       string  `json:"optiontype"`
	ClientCode       string  `json:"clientcode"`
	TradeNo          string  `json:"tradeno"`
	StampDuty        float64 `json:"stampduty"`
	GST              float64 `json:"gst"`
	TradeStatus      string  `json:"tradestatus"`
	InstrumentType   string  `json:"instrumenttype"`
	MultiplierActual string  `json:"multiplieractual"`
	TurnoverTax      float64 `json:"turnovertax"`
	StrikePrice      float64 `json:"strikeprice"`
	TradeTime        string  `json:"tradetime"`
	TradePrice       float64 `json:"tradeprice"`
	CLGTax           float64 `json:"clgtax"`
	MarketLot        string  `json:"mktlot"`
	Brokerage        float64 `json:"brokerage"`
	CTT              float64 `json:"ctt"`
	Multiplier       string  `json:"multiplier"`
	OrderTime        string  `json:"ordertime"`
	TradeDate        string  `json:"trade_date"`
	SEBITax          float64 `json:"sebitax"`
	Unit             float64 `json:"unit"`
	OrderNo          string  `json:"ORDERNO"`
	Exchange         string  `json:"exchange"`
	BuySellInd       string  `json:"buysellind"`
	TradeQty         float64 `json:"trade_qty"`
	BrokerRound      string  `json:"brokround"`
	Segment          string  `json:"segment"`
	UnixTimeFormat   int64   `json:"unixTime"`
}

type CommodityTransactionRes struct {
	UserDetails           ProfileDataResp            `json:"userDetails" mask:"struct"`
	TotalBrokerage        float64                    `json:"totalBrokerage"`
	TotalGST              float64                    `json:"totalGST"`
	TotalSEBITax          float64                    `json:"totalSEBITax"`
	TotalCTT              float64                    `json:"totalCTT"`
	TotalTurnCharges      float64                    `json:"totalTurnCharges"`
	TotalStampDuty        float64                    `json:"totalStampDuty"`
	TotalCharges          float64                    `json:"totalCharges"`
	CommodityTransactions []CommodityTransactionData `json:"commodityTransactions"`
}

type DownloadCommodityTradebookRes struct {
	DownloadUrl string `json:"downloadUrl"`
}

type FNOTradebookReq struct {
	UserID   string `json:"userID"`
	DFDateFr string `json:"dFDateFr"`
	DFDateTo string `json:"dFDateTo"`
}

type FNOTransactionRes struct {
	UserDetails          ProfileDataResp      `json:"userDetails" mask:"struct"`
	TotalBrokerage       float64              `json:"totalBrokerage"`
	TotalGST             float64              `json:"totalGST"`
	TotalSEBITax         float64              `json:"totalSEBITax"`
	TotalSTT             float64              `json:"totalSTT"`
	TotalTurnCharges     float64              `json:"totalTurnCharges"`
	TotalStampDuty       float64              `json:"totalStampDuty"`
	TotalClearingCharges float64              `json:"totalClearingCharges"`
	TotalCharges         float64              `json:"totalCharges"`
	FNOTransactions      []FNOTransactionData `json:"fnoTransactions"`
}

type FNOTransactionData struct {
	Symbol          string  `json:"symbol"`
	InstrumentType  string  `json:"instrumenttype"`
	ExpiryDate      string  `json:"expirydate"`
	OptionType      string  `json:"optiontype"`
	StrikePrice     float64 `json:"strikeprice"`
	TradeDate       string  `json:"trade_date"`
	BuySellInd      string  `json:"buysellind"`
	TradePrice      float64 `json:"tradeprice"`
	TradeQty        float64 `json:"trade_qty"`
	Brokerage       float64 `json:"brokerage"`
	GST             float64 `json:"gst"`
	STT             float64 `json:"stt"`
	SEBITax         float64 `json:"sebitax"`
	TurnoverTax     float64 `json:"turnovertax"`
	StampDuty       float64 `json:"stampduty"`
	IPFTax          float64 `json:"ipftax"`
	ClearingCharges float64 `json:"clgtax"`
	TradeTime       string  `json:"tradetime"`
	Segment         string  `json:"segment"`
	Exchange        string  `json:"exchange"`
	OrderNo         string  `json:"ORDERNO"`
	TradeNo         string  `json:"tradeno"`
	UnixTimeFormat  int64   `json:"unixTime"`
}

type DownloadFnoTradebookRes struct {
	DownloadUrl string `json:"downloadUrl"`
}
