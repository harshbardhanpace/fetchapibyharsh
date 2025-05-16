package models

type PinsMetaData struct {
	PinIndex int          `json:"pinIndex"`
	PinId    string       `json:"PinId"`
	StockDet StockDetails `json:"stockDetails"`
}

type PinsRequest struct {
	ClientId string `json:"clientId"`
}

type PinsResponse struct {
	ClientId      string         `json:"clientId"`
	PinsMetaDatas []PinsMetaData `json:"pinsDetails"`
}

type UpdatePins struct {
	ClientId      string         `json:"clientId"`
	PinsMetaDatas []PinsMetaData `json:"pinsDetails"`
}

type DeletePins struct {
	ClientId string   `json:"clientId"`
	PinId    []string `json:"PinId"`
}

type AddPinReq struct {
	ClientId         string          `json:"clientId"`
	StockDetailsData []AddPinDataReq `json:"stockDetailsData"`
}

type AddPinDataReq struct {
	Exchange      string `json:"exchange"`
	Token         string `json:"token"`
	Expiry        string `json:"expiry"`
	Company       string `json:"company"`
	Symbol        string `json:"symbol"`
	Segment       string `json:"segment"`
	TradingSymbol string `json:"tradingSymbol"`
	DisplayName   string `json:"displayName"`
}
