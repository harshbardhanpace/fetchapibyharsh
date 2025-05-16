package models

type FetchOptionChainV2Request struct {
	Token int     `json:"token"`
	Num   int     `json:"num"`
	Price float64 `json:"price" validate:"gte=0"`
}

type FetchOptionChainByExpiryV2Request struct {
	Token  int     `json:"token"`
	Num    int     `json:"num"`
	Price  float64 `json:"price" validate:"gte=0"`
	Expiry string  `json:"expiry"`
}

type OptionData struct {
	Token         string
	Exchange      string
	Company       string
	Symbol        string
	TradingSymbol string
	DisplayName   string
	StrikePrice   string
	ExpiryRaw     string
	ClosePrice    string
}

type OptionDataByExpiry struct {
	ExpiryDate string       `json:"expiryDate"`
	Strikes    []OptionData `json:"strikes"`
}
