package models

// Search Scrip Request
type SearchScripRequest struct {
	Key string `json:"key"`
}

type SearchScripAPIRequest struct {
	SearchText string `json:"searchText"`
	Exchange   string `json:"exchange"`
	Page       int    `json:"page"`
}

// Search Scrip Response
type SearchScripResponse struct {
	// Error struct {
	// 	Code    int    `json:"code"`
	// 	Message string `json:"message"`
	// } `json:"error"`
	Result []SearchScripResponseResult `json:"result"`
}

type ContractDetails struct {
	Exchange        string  `json:"exchange"`
	Token1          string  `json:"token1"`
	Token2          string  `json:"token2"`
	UnderlyingToken string  `json:"underlyingToken"`
	Symbol          string  `json:"symbol"`
	Series          string  `json:"series"`
	Instname        string  `json:"instName"`
	Name            string  `json:"name"`
	TradingSymbol   string  `json:"tradingSymbol"`
	LowDpr          float64 `json:"lowDpr"`
	HighDpr         float64 `json:"highDpr"`
	LotSize         int     `json:"lotSize"`
	TickSize        int     `json:"tickSize"`
	Multiplier      int     `json:"multiplier"`
	Precision       int     `json:"precision"`
	Strike          int     `json:"strike"`
	Isin            string  `json:"isin"`
	IsTradable      bool    `json:"isTradable"`
	Expiry1         string  `json:"expiry1"`
	Expiry2         string  `json:"expiry2"`
	OptionType      string  `json:"optionType"`
	AlternateToken  string  `json:"alternateToken"`
	IsMtfEligible   bool    `json:"isMtfEligible"`
}

type SearchScripResponseResult struct {
	Token         string  `json:"token"`
	Exchange      string  `json:"exchange"`
	Execution     string  `json:"execution"`
	Company       string  `json:"company"`
	Symbol        string  `json:"symbol"`
	Isin          string  `json:"isin"`
	TradingSymbol string  `json:"tradingSymbol"`
	DisplayName   string  `json:"displayName"`
	Score         float64 `json:"score"`
	ClosePrice    string  `json:"closePrice"`
	IsTradable    bool    `json:"isTradable"`
	Segment       string  `json:"segment"`
	Tag           string  `json:"tag"`
	Expiry        string  `json:"expiry"`
	Series        string  `json:"series"`
	Strike        int     `json:"strike"`
	IsMtfEligible bool    `json:"isMtfEligible"`
	Alternate     struct {
		Token         string  `json:"token"`
		Exchange      string  `json:"exchange"`
		Execution     string  `json:"execution"`
		Company       string  `json:"company"`
		Symbol        string  `json:"symbol"`
		TradingSymbol string  `json:"tradingSymbol"`
		DisplayName   string  `json:"displayName"`
		Score         float64 `json:"score"`
		ClosePrice    string  `json:"closePrice"`
		IsTradable    bool    `json:"isTradable"`
		Segment       string  `json:"segment"`
		Tag           string  `json:"tag"`
		Expiry        string  `json:"expiry"`
		Series        string  `json:"series"`
		Strike        int     `json:"strike"`
		IsMtfEligible bool    `json:"isMtfEligible"`
	}
}

// ScripInfoRequest
type ScripInfoRequest struct {
	Exchange string `json:"exchange" enums:"NSE,BSE,NFO,CDS,MCX,BFO" validate:"oneof=NSE BSE NFO CDS MCX BFO"`
	Info     string `json:"info"`
	Token    string `json:"token"`
}

// ScripInfoResponse
type ScripInfoResponse struct {
	// Error struct {
	// 	Code    int    `json:"code"`
	// 	Message string `json:"message"`
	// } `json:"error"`
	Result struct {
		BoardLotQuantity         int     `json:"boardLotQuantity"`
		ChangeInOi               int     `json:"changeInOi"`
		Exchange                 int     `json:"exchange"`
		Expiry                   int     `json:"expiry"`
		HigherCircuitLimit       float64 `json:"higherCircuitLimit"`
		InstrumentName           string  `json:"instrumentName"`
		InstrumentToken          int     `json:"instrumentToken"`
		Isin                     string  `json:"isin"`
		LowerCircuitLimit        float64 `json:"lowerCircuitLimit"`
		Multiplier               int     `json:"multiplier"`
		OpenInterest             int     `json:"openInterest"`
		OptionType               string  `json:"optionType"`
		Precision                int     `json:"precision"`
		Series                   string  `json:"series"`
		Strike                   int     `json:"strike"`
		Symbol                   string  `json:"symbol"`
		TickSize                 float64 `json:"tickSize"`
		TradingSymbol            string  `json:"tradingSymbol"`
		UnderlyingToken          int     `json:"underlyingToken"`
		RawExpiry                int     `json:"rawExpiry"`
		Freeze                   int     `json:"freeze"`
		InstrumentType           string  `json:"instrumentType"`
		IssueRate                int     `json:"issueRate"`
		IssueStartDate           string  `json:"issueStartDate"`
		ListDate                 string  `json:"listDate"`
		MaxOrderSize             int     `json:"maxOrderSize"`
		PriceNumerator           float64 `json:"priceNumerator"`
		PriceDenominator         float64 `json:"priceDenominator"`
		Comments                 string  `json:"comments"`
		CircuitRating            string  `json:"circuitRating"`
		CompanyName              string  `json:"companyName"`
		DisplayName              string  `json:"displayName"`
		RawTickSize              int     `json:"rawTickSize"`
		IsIndex                  bool    `json:"isIndex"`
		Tradable                 bool    `json:"tradable"`
		MaxSingleQty             int     `json:"maxSingleQty"`
		ExpiryString             string  `json:"expiryString"`
		LocalUpdateTime          string  `json:"localUpdateTime"`
		MarketType               string  `json:"marketType"`
		PriceUnits               string  `json:"priceUnits"`
		TradingUnits             string  `json:"tradingUnits"`
		LastTradingDate          string  `json:"lastTradingDate"`
		TenderPeriodEndDate      string  `json:"tenderPeriodEndDate"`
		DeliveryStartDate        string  `json:"deliveryStartDate"`
		PriceQuotation           float64 `json:"priceQuotation"`
		GeneralDenominator       string  `json:"generalDenominator"`
		TenderPeriodStartDate    string  `json:"tenderPeriodStartDate"`
		DeliveryUnits            string  `json:"deliveryUnits"`
		DeliveryEndDate          string  `json:"deliveryEndDate"`
		TradingUnitFactor        int     `json:"tradingUnitFactor"`
		DeliveryUnitFactor       int     `json:"deliveryUnitFactor"`
		BookClosureEndDate       string  `json:"bookClosureEndDate"`
		BookClosureStartDate     string  `json:"bookClosureStartDate"`
		NoDeliveryDateEnd        string  `json:"noDeliveryDateEnd"`
		NoDeliveryDateStart      string  `json:"noDeliveryDateStart"`
		ReAdmissionDate          string  `json:"reAdmissionDate"`
		RecordDate               string  `json:"recordDate"`
		Warning                  string  `json:"warning"`
		Dpr                      string  `json:"dpr"`
		TradeToTrade             bool    `json:"tradeToTrade"`
		SurveillanceIndicator    int     `json:"surveillanceIndicator"`
		PartitionID              int     `json:"partitionId"`
		ProductID                int     `json:"productId"`
		ProductCategory          string  `json:"productCategory"`
		MonthIdentifier          int     `json:"monthIdentifier"`
		ClosePrice               string  `json:"closePrice"`
		SpecialPreopen           int     `json:"specialPreopen"`
		AlternateExchange        string  `json:"alternateExchange"`
		AlternateToken           int     `json:"alternateToken"`
		Asm                      string  `json:"asm"`
		Gsm                      string  `json:"gsm"`
		Execution                string  `json:"execution"`
		Symbol2                  string  `json:"symbol2"`
		RawTenderPeriodStartDate string  `json:"rawTenderPeriodStartDate"`
		RawTenderPeriodEndDate   string  `json:"rawTenderPeriodEndDate"`
		YearlyHighPrice          string  `json:"yearlyHighPrice"`
		YearlyLowPrice           string  `json:"yearlyLowPrice"`
		IssueMaturityDate        int     `json:"issueMaturityDate"`
		Var                      string  `json:"var"`
		Exposure                 string  `json:"exposure"`
		Span                     []int   `json:"span"`
		HaveFutures              bool    `json:"haveFutures"`
		HaveOptions              bool    `json:"haveOptions"`
		Tag                      string  `json:"tag"`
		ShortCode                string  `json:"shortCode"`
		IsMtfEligible            bool    `json:"isMtfEligible"`
		IsMisEligible            bool    `json:"isMisEligible"`
		ExBonusDate              string  `json:"exBonusDate"`
		ExDate                   string  `json:"exDate"`
		Exflag                   string  `json:"exFlag"`
		ExRightDate              string  `json:"exRightDate"`
		MtfMargin                float64 `json:"mtf_margin"`
	} `json:"result"`
}
