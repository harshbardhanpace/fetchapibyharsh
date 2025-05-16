package backoffice

type ShilpiTradeConfirmationDateRangeRes []struct {
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
	TradeDate   string `json:"TradeDate"`
	NetAvgRate  string `json:"NetAvgRate"`
}

type ShilpiGetBillDetailsCdsl []struct {
	Charges        string `json:"charges"`
	TrxDate        string `json:"TrxDate"`
	Qty            string `json:"qty"`
	Gst            string `json:"gst"`
	InstrumentName string `json:"instrument name"`
	Isincode       string `json:"isincode"`
	ChargesDetails string `json:"charges details"`
	TotalCharges   string `json:"total charges"`
}

type ShilpiLongTermShortTerm []struct {
	Scripname string `json:"Scripname"`
	BuyRate   string `json:"BuyRate"`
	SellQty   string `json:"SellQty"`
	Jobbing   string `json:"Jobbing"`
	SellRate  string `json:"SellRate"`
	ShortTerm string `json:"ShortTerm"`
	Buyqty    string `json:"Buyqty"`
	LongTerm  string `json:"LongTerm"`
	Scripcd   string `json:"scripcd"`
	Isin      string `json:"isin"`
}

type ShilpiFetchProfile []struct {
	Emailno          string `json:"emailno"`
	Phnos            string `json:"phnos"`
	Occupationcode   string `json:"occupationcode"`
	Poaflag          string `json:"poaflag"`
	Gender           string `json:"gender"`
	City             string `json:"city"`
	Annualincomedate string `json:"annualincomedate"`
	Panno            string `json:"panno"`
	Annualincomecode string `json:"annualincomecode"`
	Dpid             string `json:"dpid"`
	Mobileno         string `json:"mobileno"`
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
	Bankacno         string `json:"bankacno"`
	Ifsccode         string `json:"ifsccode"`
	Bankactype       string `json:"bankactype"`
	Dpaccountno      string `json:"dpaccountno"`
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

type ShilpiTradeConfirmationOnDate []struct {
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

type ShilpiOpenPositions []struct {
	Num1  string `json:"1"`
	Num2  string `json:"2"`
	Num3  string `json:"3"`
	Num4  string `json:"4"`
	Num5  string `json:"5"`
	Num6  string `json:"6"`
	Num7  string `json:"7"`
	Num8  string `json:"8"`
	Num9  string `json:"9"`
	Num10 string `json:"10"`
	Num11 string `json:"11"`
	Num12 string `json:"12"`
}

type ShilpiGetHolding []struct {
	Num1 string `json:"1"`
	Num2 string `json:"2"`
	Num3 string `json:"3"`
	Num4 string `json:"4"`
	Num5 string `json:"5"`
	Num6 string `json:"6"`
	Num7 string `json:"7"`
	Num8 string `json:"8"`
}

type ShilpiGetMarginOnDate []struct {
	TotalMargin     string `json:"total_margin"`
	Exposure        string `json:"exposure"`
	PeakDeposit     string `json:"peakdeposit"`
	VarMargin       string `json:"varmargin"`
	Deposit         string `json:"deposit"`
	ShortMargin     string `json:"shortmargin"`
	TotalPeakmargin string `json:"totalpeakmargin"`
	PeakMarginShort string `json:"peakmarginshort"`
	Span            string `json:"span"`
}

type ShilpiFinancialLedgerBalanceOnDate []struct {
	Num1 string `json:"1"`
}

type ShilpiGetFinancial []struct {
	Date      string `json:"date"`
	Narr      string `json:"narr"`
	Segment   string `json:"segment"`
	Exchange  string `json:"exchange"`
	Debit     string `json:"debit"`
	Credit    string `json:"credit"`
	Payrefno  string `json:"payrefno"`
	ValueDate string `json:"value_date"`
}

type ShilpiScripWiseCosting []struct {
	NetValue  string `json:"netvalue"`
	Brok      string `json:"brok"`
	NetQty    string `json:"netqty"`
	OrderNo   string `json:"orderno"`
	OthChrg2  string `json:"othchrg2"`
	OthChrg1  string `json:"othchrg1"`
	ScripCode string `json:"scripcode"`
	GST       string `json:"gst"`
	Stamp     string `json:"stamp"`
	AvgRate   string `json:"avgrate"`
	TrxDate   string `json:"trxdate"`
	ISINCode  string `json:"isincode"`
	SellQty   string `json:"sellqty"`
	SEBIFee   string `json:"sebifee"`
	NetAmt    string `json:"netamt"`
	SellValue string `json:"sellvalue"`
	STT       string `json:"stt"`
	BuyQty    string `json:"buyqty"`
	BuyValue  string `json:"buyvalue"`
	ExchClg   string `json:"exchclg"`
	TurnTax   string `json:"turntax"`
	Exchange  string `json:"exchange"`
	BrokType  string `json:"broktype"`
	ScripName string `json:"scrip name"`
}

type ShilpiGetFinancials []struct {
	Date      string `json:"date"`
	Narr      string `json:"narr"`
	Segment   string `json:"segment"`
	Exchange  string `json:"exchange"`
	Debit     string `json:"debit"`
	Credit    string `json:"credit"`
	PayRefNo  string `json:"payrefno"`
	RuningBal string `json:"runingbal"`
	ValueDate string `json:"value_date"`
}

type ShilpiGetOpenPosition []struct {
	Segment     string `json:"1"`
	Exchange    string `json:"2"`
	Contract    string `json:"3"`
	BuyQty      string `json:"4"`
	BuyValue    string `json:"5"`
	SellQty     string `json:"6"`
	SellValue   string `json:"7"`
	NetQty      string `json:"8"`
	NetValue    string `json:"9"`
	ClosingRate string `json:"10"`
	PnL         string `json:"11"`
	Exposure    string `json:"12"`
	ExpiryDate  string `json:"13"`
	StrikePrice string `json:"14"`
	OptionType  string `json:"15"`
}

type ShilpiFONetPosition []struct {
	ID            string `json:"1"`
	NetValue      string `json:"netvalue"`
	OptionType    string `json:"optiontype"`
	NetQty        int    `json:"netqty"`
	NetPL         string `json:"netpl"`
	BuyAvgRate    string `json:"buyavgrate"`
	AvgRate       string `json:"avgrate"`
	Contracts     string `json:"contracts"`
	BuyQty        int    `json:"buyqty"`
	SaleAvgRate   string `json:"saleavgrate"`
	BuyValue      string `json:"buyvalue"`
	ClosPrice     string `json:"closprice"`
	StrikePrice   string `json:"strikeprice"`
	SaleQty       int    `json:"saleqty"`
	SaleValue     string `json:"salevalue"`
	ExpDate       string `json:"expdate"`
	OtherCharges2 string `json:"othercharges2"`
	IPFChrg       string `json:"ipfchrg"`
	OtherCharges1 string `json:"othercharges1"`
	StampDuty     string `json:"stampduty"`
	ClgCharges    string `json:"clgcharges"`
	MinBrok       string `json:"minbrok"`
	SebiFee       string `json:"sebifee"`
	NetPay        string `json:"netpay"`
	STT           string `json:"stt"`
	TotalGST      string `json:"totalgst"`
	FutMTM        string `json:"futmtm"`
	ExchClg       string `json:"exchclg"`
	TurnTax       string `json:"turntax"`
	TotPremium    string `json:"totpremium"`
}

type ShilpiHoldingFinancial []struct {
	MTFCollateral string `json:"MTF Collateral"`
	ScripName     string `json:"Scrip Name"`
	ScripCode     string `json:"ScripCode"`
	DPStock       string `json:"DP Stock"`
	Holding       string `json:"Holding"`
	BuyPrice      string `json:"Buy Price"`
	PledgeStock   string `json:"Pledge Stock"`
	ClosingPrice  string `json:"ClosingPrice"`
	ISIN          string `json:"ISIN"`
	TotalHolding  string `json:"Total Holding"`
	Variance      string `json:"Variance"`
}

type ShilpiCommodityTransaction []struct {
	Symbol           string `json:"symbol"`
	ExpiryDate       string `json:"expirydate"`
	OptionType       string `json:"optiontype"`
	ClientCode       string `json:"clientcode"`
	TradeNo          string `json:"tradeno"`
	StampDuty        string `json:"stampduty"`
	GST              string `json:"gst"`
	TradeStatus      string `json:"tradestatus"`
	InstrumentType   string `json:"instrumenttype"`
	MultiplierActual string `json:"multiplieractual"`
	TurnoverTax      string `json:"turnovertax"`
	StrikePrice      string `json:"strikeprice"`
	TradeTime        string `json:"tradetime"`
	TradePrice       string `json:"tradeprice"`
	CLGTax           string `json:"clgtax"`
	MarketLot        string `json:"mktlot"`
	Brokerage        string `json:"brokerage"`
	CTT              string `json:"ctt"`
	Multiplier       string `json:"multiplier"`
	OrderTime        string `json:"ordertime"`
	TradeDate        string `json:"trade_date"`
	SEBITax          string `json:"sebitax"`
	Unit             string `json:"unit"`
	OrderNo          string `json:"ORDERNO"`
	Exchange         string `json:"exchange"`
	BuySellInd       int    `json:"buysellind"`
	TradeQty         string `json:"trade_qty"`
	BrokerRound      string `json:"brokround"`
}

type ShilpiFNOTransaction []struct {
	Symbol         string `json:"symbol"`
	ExpiryDate     string `json:"expirydate"`
	OptionType     string `json:"optiontype"`
	ClientCode     string `json:"clientcode"`
	TradeNo        string `json:"tradeno"`
	StampDuty      string `json:"stampduty"`
	GST            string `json:"gst"`
	TradeStatus    string `json:"tradestatus"`
	TrnLot         string `json:"trnlot"`
	InstrumentType string `json:"instrumenttype"`
	TurnoverTax    string `json:"turnovertax"`
	StrikePrice    string `json:"strikeprice"`
	TradeTime      string `json:"tradetime"`
	TradePrice     string `json:"tradeprice"`
	ClgTax         string `json:"clgtax"`
	Brokerage      string `json:"brokerage"`
	IPFTax         string `json:"ipftax"`
	OrderTime      string `json:"ordertime"`
	TradeDate      string `json:"trade_date"`
	SebiTax        string `json:"sebitax"`
	STT            string `json:"stt"`
	OrderNo        string `json:"ORDERNO"`
	Exchange       string `json:"exchange"`
	BuySellInd     int    `json:"buysellind"`
	TradeQty       string `json:"trade_qty"`
	BrokRound      string `json:"brokround"`
}
