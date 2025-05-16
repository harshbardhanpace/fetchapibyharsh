package models

//Fetch Option Chain Request
type FetchOptionChainRequest struct {
	Token    int     `json:"token"`
	Num      int     `json:"num"`
	Price    float64 `json:"price" validate:"gte=0"`
	Exchange string  `json:"exchange" example:"NSE,BSE,NFO,BFO,MCX,CDS"  validate:"oneof=NSE BSE MCX NFO BFO CDS"`
}

//Fetch Option Chain Response
type FetchOptionChainResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Result []FetchOptionChainResponseData `json:"result"`
}

type FetchOptionChainResponseData struct {
	ExpiryDate string                                `json:"expiryDate"`
	Strikes    []FetchOptionChainResponseDataStrikes `json:"strikes"`
}

type FetchOptionChainResponseDataStrikes struct {
	StrikePrice float64 `json:"strikePrice"`
	CallOption  struct {
		Token         string  `json:"token"`
		Exchange      string  `json:"exchange"`
		Company       string  `json:"company"`
		Symbol        string  `json:"symbol"`
		TradingSymbol string  `json:"tradingSymbol"`
		DisplayName   string  `json:"displayName"`
		StrikePrice   float64 `json:"strikePrice"`
		ExpiryRaw     string  `json:"expiryRaw"`
		ClosePrice    string  `json:"closePrice"`
	} `json:"callOption"`
	PutOption struct {
		Token         string  `json:"token"`
		Exchange      string  `json:"exchange"`
		Company       string  `json:"company"`
		Symbol        string  `json:"symbol"`
		TradingSymbol string  `json:"tradingSymbol"`
		DisplayName   string  `json:"displayName"`
		StrikePrice   float64 `json:"strikePrice"`
		ExpiryRaw     string  `json:"expiryRaw"`
		ClosePrice    string  `json:"closePrice"`
	} `json:"putOption"`
}

type FetchFuturesChainReq struct {
	Token string `json:"token"`
}

type FetchFuturesChainRes struct {
	Result []FuturesChain `json:"result"`
}

type FuturesChain struct {
	ExpiryDate string `json:"expiryDate"`
	Strikes    Strike `json:"result"`
}

type Strike struct {
	Token         string `json:"token"`
	Exchange      string `json:"exchange"`
	Company       string `json:"company"`
	Symbol        string `json:"symbol"`
	TradingSymbol string `json:"tradingSymbol"`
	DisplayName   string `json:"displayName"`
	ExpiryRaw     string `json:"ExpiryRaw"`
	ClosePrice    string `json:"closePrice"`
}
