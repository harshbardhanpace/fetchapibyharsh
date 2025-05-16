package models

type CreateWatchListRequest struct {
	ClientId               string `json:"clientId"`
	WatchListName          string `json:"watchListName"`
	WatchListNameLongDesc  string `json:"watchListLongDesc"`
	WatchListNameShortDesc string `json:"watchListShortDesc"`
}

type CreateWatchListResponse struct {
	WatchListId string `json:"watchListId"`
}

type ModifyWatchListRequest struct {
	WatchListId            string `json:"clientId"`
	WatchListName          string `json:"watchListName"`
	WatchListNameLongDesc  string `json:"watchListLongDesc"`
	WatchListNameShortDesc string `json:"watchListShortDesc"`
}

type FetchWatchListsRequest struct {
	ClientId string `json:"clientId"`
}

type WatchListMetadata struct {
	WatchListId            string `json:"watchListId"`
	WatchListName          string `json:"watchListName"`
	WatchListNameLongDesc  string `json:"watchListLongDesc"`
	WatchListNameShortDesc string `json:"watchListShortDesc"`
}

type FetchWatchListResponse struct {
	WatchListIds []WatchListMetadata `json:"watchListIds"`
}

type FetchWatchListsDetailsRequest struct {
	ClientId    string `json:"clientId"`
	WatchListId string `json:"watchListId"`
}

type FetchWatchListsDetailsResponse struct {
	Stocks []StockDetails `json:"stocks"`
}

type StockDetails struct {
	StockId       string `json:"stockId"`
	Exchange      string `json:"exchange" example:"NSE,BSE,NFO,BFO,MCX,CDS" validate:"oneof=NSE BSE MCX NFO BFO CDS"`
	Token         string `json:"token"`
	Expiry        string `json:"expiry"`
	Company       string `json:"company"`
	Symbol        string `json:"symbol"`
	TradingSymbol string `json:"tradingSymbol"`
	DisplayName   string `json:"displayName"`
	Isin          string `json:"isin"`
	IsTradable    bool   `json:"isTradable"`
	Segment       string `json:"segment"`
	Execution     string `json:"execution"`
}

type AddStockToWatchListsRequest struct {
	ClientId    string       `json:"clientId"`
	WatchListId string       `json:"watchListId"`
	Stock       StockDetails `json:"stock"`
}

type DeleteStockToWatchListsRequest struct {
	ClientId    string       `json:"clientId"`
	WatchListId string       `json:"watchListId"`
	Stock       StockDetails `json:"stock"`
}

type DeleteWatchListRequest struct {
	ClientId    string `json:"clientId"`
	WatchListId string `json:"watchListId"`
}

type FetchWatchListV2Request struct {
	ClientId string `json:"clientId"`
}

type WatchListsDetails struct {
	WatchListId string         `json:"watchListId"`
	Stocks      []StockDetails `json:"stocks"`
}

type FetchWatchListV2Response struct {
	ClientId   string              `json:"clientId"`
	WatchLists []WatchListsDetails `json:"watchLists"`
}

type AddStockToWatchListV2Request struct {
	ClientId    string       `json:"clientId"`
	WatchListId []string     `json:"watchListId" validate:"required,min=1,dive,oneof=wl1 wl2 wl3 wl4 wl5"`
	Stock       StockDetails `json:"stock"`
}

type DeleteWatchListV2Request struct {
	ClientId    string   `json:"clientId"`
	WatchListId string   `json:"watchListId"`
	StockId     []string `json:"stockId"`
}

type ArrangeStocksWatchListV2Request struct {
	ClientId    string   `json:"clientId"`
	WatchListId string   `json:"watchListId"`
	StockIds    []string `json:"stockIds"`
}

type AddStockToWatchListConciseV2Request struct {
	ClientId    string      `json:"clientId"`
	WatchListId string      `json:"watchListId"`
	StockData   []StockData `json:"stockData"`
}

type StockData struct {
	Exchange        string `json:"exchange"`
	InstrumentToken string `json:"instrumentToken"`
}

type DeleteWatchListV2UpdatedRequest struct {
	ClientId    string      `json:"clientId"`
	WatchListId []WatchList `json:"watchListId"`
}

type WatchList struct {
	WatchlistId string   `json:"watchlistId"`
	StockId     []string `json:"stockId"`
}

type FetchWatchListV3Request struct {
	ClientId string `json:"clientId"`
}

type FetchWatchListV3Response struct {
	ClientId   string                `json:"clientId"`
	WatchLists []WatchListsDetailsV2 `json:"watchLists"`
}

type AddStockToWatchListV3Request struct {
	ClientId    string         `json:"clientId"`
	WatchListId []string       `json:"watchListId" validate:"required,min=1,dive,oneof=wl1 wl2 wl3 wl4 wl5"`
	Stock       StockDetailsV2 `json:"stock"`
}

type DeleteWatchListV3Request struct {
	ClientId    string   `json:"clientId"`
	WatchListId string   `json:"watchListId"`
	StockId     []string `json:"stockId"`
}

type StockDetailsWithExpiry struct {
	Expiry int64          `json:"expiry"`
	Stock  StockDetailsV2 `json:"stock"`
}

type StockDetailsV2 struct {
	StockId       string `json:"stockId"`
	Exchange      string `json:"exchange" example:"NSE,BSE,NFO,BFO,MCX,CDS" validate:"oneof=NSE BSE MCX NFO BFO CDS"`
	Token         string `json:"token"`
	Expiry        string `json:"expiry"`
	Company       string `json:"company"`
	Symbol        string `json:"symbol"`
	TradingSymbol string `json:"tradingSymbol"`
	DisplayName   string `json:"displayName"`
	Isin          string `json:"isin"`
	IsTradable    bool   `json:"isTradable"`
	Segment       string `json:"segment"`
	Execution     string `json:"execution"`
	IsinStockId   string `json:"isinStockId"`
}

type WatchListsDetailsV2 struct {
	WatchListId string           `json:"watchListId"`
	Stocks      []StockDetailsV2 `json:"stocks"`
}

type ArrangeStocksWatchListV3Request struct {
	ClientId    string   `json:"clientId"`
	WatchListId string   `json:"watchListId"`
	StockIds    []string `json:"stockIds"`
}

type DeleteWatchListV3UpdatedRequest struct {
	ClientId    string      `json:"clientId"`
	WatchListId []WatchList `json:"watchListId"`
}
