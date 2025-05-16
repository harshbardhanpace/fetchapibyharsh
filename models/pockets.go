package models

import "time"

type AdminLoginRequest struct {
	UserId   string `json:"userId"`
	Password string `json:"password" mask:"id"`
}

type AdminLoginResponse struct {
	AuthToken string `json:"authToken"`
}

type PocketsMetaData struct {
	StockId       string `json:"stockId"`
	Qty           string `json:"qty"`
	Token         string `json:"token"`
	Exchange      string `json:"exchange"`
	Company       string `json:"company"`
	Symbol        string `json:"symbol"`
	TradingSymbol string `json:"tradingSymbol"`
	DisplayName   string `json:"displayName"`
	Isin          string `json:"isin"`
	IsTradable    bool   `json:"isTradable"`
	Segment       string `json:"segment"`
	StockImage    string `json:"stockImage"`
}

type CreatePocketsRequest struct {
	PocketName      string            `json:"pocketName"`
	PocketShortDesc string            `json:"pocketShortDesc"`
	PocketLongDesc  string            `json:"pocketLongDesc"`
	PocketExchange  string            `json:"pocketExchange"`
	PocketImage     string            `json:"pocketImage"`
	PocketTokens    []PocketsMetaData `json:"pocketTokens"`
}

type CreatePocketsResponse struct {
	PocketId     string            `json:"pocketId"`
	PocketName   string            `json:"pocketName"`
	PocketTokens []PocketsMetaData `json:"pocketTokens"`
}

type FetchPocketsDetailsRequest struct {
	PocketId string `json:"pocketId"`
}

type FetchPocketsDetailsResponse struct {
	PocketId               string            `json:"pocketId"`
	PocketName             string            `json:"pocketName"`
	PocketShortDesc        string            `json:"pocketShortDesc"`
	PocketLongDesc         string            `json:"pocketLongDesc"`
	PocketExchange         string            `json:"pocketExchange"`
	PocketImage            string            `json:"pocketImage"`
	PocketWebImage         string            `json:"pocketWebImage"`
	PrimaryBackgroundColor string            `json:"primaryBackgroundColor"`
	PrimarySecondaryColor  string            `json:"primarySecondaryColor"`
	PocketTokens           []PocketsMetaData `json:"pocketTokens"`
}

type ModifyPocketsRequest struct {
	PocketName      string            `json:"pocketName"`
	PocketShortDesc string            `json:"pocketShortDesc"`
	PocketLongDesc  string            `json:"pocketLongDesc"`
	PocketExchange  string            `json:"pocketExchange"`
	PocketImage     string            `json:"pocketImage"`
	PocketTokens    []PocketsMetaData `json:"pocketTokens"`
	PocketId        string            `json:"pocketId"`
}

type ModifyPocketsResponse struct {
	PocketId     string            `json:"pocketId"`
	PocketTokens []PocketsMetaData `json:"pocketTokens"`
}

type DeletePocketsRequest struct {
	PocketId string `json:"pocketId"`
}

type DeletePocketsResponse struct {
	PocketTokens []PocketsMetaData `json:"pocketTokens"`
}

type FetchAllPocketsDetailsResponse struct {
	FetchAllPocketsDetailsResponse []FetchPocketsDetailsResponse `json:"fetchAllPocketsDetailsResponse"`
}

type FetchPocketsRequest struct {
	ClientId string `json:"clientId"`
}

type FetchPocketsMetadata struct {
	Exchange   string `json:"pocketExchange"`
	BuyPrice   string `json:"pocketBuyPrice"`
	BuyQty     string `json:"pocketBuyQty"`
	Token      string `json:"pocketToken"`
	BoughtDate string `json:"pocketBoughtDate"`
}

type PocketsMeta struct {
	PocketID        string                 `json:"pocketId"`
	PocketName      string                 `json:"pocketName"`
	PocketShortDesc string                 `json:"pocketShortDesc"`
	PocketLongDesc  string                 `json:"pocketLongDesc"`
	Details         []FetchPocketsMetadata `json:"details"`
}

type FetchPocketsResponse struct {
	Pockets []PocketsMeta `json:"pockets"`
}

type ExecutePocketRequest struct {
	PocketId    string  `json:"pocketId"`
	ClientId    string  `json:"clientId"`
	PocketPrice float64 `json:"pocketPrice"`
}

type ExecutePocketResponse struct {
	BasketID string `json:"basketId"`
	Message  string `json:"message"`
}

type FetchPocketPortfolioRequest struct {
	ClientId string `json:"clientId"`
}

type ExitPocketRequest struct {
	PocketId string `json:"pocketId"`
}

type PocketsCalculationsReq struct {
	PocketId          string `json:"pocketId"`
	BenchmarkToken    string `json:"benchmarkToken"`
	BenchmarkExchange string `json:"benchmarkExchange"`
	TimeInterval      int    `json:"timeInterval"`
}

type PocketsCalculationsRes struct {
	Benchmark []CalculatedChartData `json:"benchmark"`
	Pocket    []CalculatedChartData `json:"pocket"`
}

type MultipleAndIndividualStocksCalculationsReq struct {
	Stocks            []StocksCalculationsStruct `json:"stocks"`
	BenchmarkToken    string                     `json:"benchmarkToken"`
	BenchmarkExchange string                     `json:"benchmarkExchange"`
	TimeInterval      int                        `json:"timeInterval"`
}

type StocksCalculationsStruct struct {
	StockToken    string `json:"stockToken"`
	StockExchange string `json:"stockExchange"`
	StockQuantity int    `json:"stockQuantity"`
}

type MultipleAndIndividualStocksCalculationsRes struct {
	MultipleOrIndividualStocks []CalculatedChartData `json:"multipleOrIndividualStocks"`
}

type CalculatedChartData struct {
	PercentageGain          float64 `json:"percentageGain"`
	PercentageGainBenchmark float64 `json:"percentageGainBenchmark"`
	Date                    string  `json:"date"`
}

type CalculationDataTemp struct {
	PercentageGain float64 `json:"percentageGain"`
	Date           string  `json:"date"`
}

type TLCandleData struct {
	Timestamp string
	Close     float64
}

type TotalPrice struct {
	TokenPrice float64
	Quantity   int
}

type PocketTransaction struct {
	ClientId          string    `json:"clientId"`
	PocketID          string    `json:"pocketID"`
	PocketName        string    `json:"pocketName"`
	TransactionDate   time.Time `json:"transactionDate"`
	TotalInvestment   float64   `json:"totalInvestment"`
	TransactionID     string    `json:"transactionID"`
	TransactionStatus int       `json:"transactionStatus"` // 0 - Bought, 1 - SOLD/Exited, 2 - for future needs
	BasketCounter     int       `json:"basketCounter"`
	BasketName        string    `json:"basketName"`
}

type PocketTransactionComplete struct {
	ClientId           string       `json:"clientId"`
	AllPocketPurchases []PocketInfo `json:"allPocketPurchases"`
	PocketCounter      int          `json:"pocketCounter"`
}

type FetchPocketTransactionReq struct {
	ClientId string `json:"clientId"`
}

type FetchPocketTransactionRes struct {
	AllPocketTransactionComplete []PocketTransactionComplete `json:"allPocketTransactionComplete"`
}

type PocketInfo struct {
	PocketID               string    `json:"pocketID"`
	PocketName             string    `json:"pocketName"`
	PocketImage            string    `json:"pocketImage"`
	PocketWebImage         string    `json:"pocketWebImage"`
	PrimaryBackgroundColor string    `json:"primaryBackgroundColor"`
	PrimarySecondaryColor  string    `json:"primarySecondaryColor"`
	TransactionDate        time.Time `json:"transactionDate"`
	PocketExecutedPrice    float64   `json:"pocketExecutedPrice"`
	TransactionID          string    `json:"transactionID"`
	TransactionStatus      int       `json:"transactionStatus"` // 0 - Bought, 1 - SOLD/Exited, 2 - for future needs
	// BasketCounter     int       `json:"basketCounter"`
	BasketName string `json:"basketName"`
	Qty        int    `json:"qty"`
}

type FetchPocketPortfolioResponse struct {
	ClientId         string   `json:"clientId"`
	PortfolioDetails []Pocket `json:"portfolioDetails"`
}

type Pocket struct {
	PocketID               string  `json:"pocketId"`
	PocketName             string  `json:"pocketName"`
	PocketImage            string  `json:"pocketImage"`
	PocketWebImage         string  `json:"pocketWebImage"`
	PrimaryBackgroundColor string  `json:"primaryBackgroundColor"`
	PrimarySecondaryColor  string  `json:"primarySecondaryColor"`
	AveragePrice           float64 `json:"averagePrice"`
	TotalBuyPockets        int     `json:"totalBuyPockets"`
	TotalSellPockets       int     `json:"totalSellPockets"`
	TotalInvestment        float64 `json:"totalInvestment"`
}

type StorePocketTransactionReq struct {
	PocketId    string  `json:"pocketId"`
	ClientId    string  `json:"clientId"`
	PocketPrice float64 `json:"pocketPrice"`
	OrderSide   string  `json:"orderSide"`
	Qty         int     `json:"qty"`
}

type TokenExchange struct {
	Exchange string `json:"exchange"`
	Token    string `json:"token"`
}

type ExecutePocketV2Request struct {
	PocketId string `json:"pocketId"`
	ClientId string `json:"clientId"`
	LotSize  int    `json:"lotSize"`
}

type PocketTransactionStoreV2 struct {
	ClientId           string            `json:"clientId"`
	AllPocketPurchages []PocketDetailsV2 `json:"pocketDetailsV2"`
}

type PocketDetailsV2 struct {
	PocketTransactionId string            `json:"pocketTransactionId"`
	PocketID            string            `json:"pocketID"`
	TransactionStatus   int               `json:"transactionStatus"` // 0 - Bought, 1 - SOLD/Exited, 2 - for future needs
	LotSize             int               `json:"lotSize"`
	OrderCompletedPrice float64           `json:"orderCompletedPrice"`
	OrderCompleted      []PocketsMetaData `json:"orderCompleted"`
	OrderCancelled      []PocketsMetaData `json:"orderCancelled"`
}

// pocket-V3 ...

type ExecutePocketV3Request struct {
	PocketId string `json:"pocketId" binding:"required"`
	ClientId string `json:"clientId" binding:"required"`
	LotSize  int    `json:"lotSize" binding:"required,gt=0"`
}

// Represents the user's collection of all pockets they own
type UserPocketHolding struct {
	ClientID string       `json:"clientId" bson:"clientId"`
	Pockets  []UserPocket `json:"pockets" bson:"pockets"`
}

// Represents each specific pocket a user holds
type UserPocket struct {
	PocketID string `json:"pocketId" bson:"pocketId"`
	Version  int    `json:"version" bson:"version"`
	LotSize  int    `json:"lotSize" bson:"lotSize"`
}

// Represents user's current holdings fetched from the demat account
type UserHolding struct {
	Token         int    `json:"token"`
	Quantity      int    `json:"quantity"`
	Symbol        string `json:"symbol"`
	TradingSymbol string `json:"tradingSymbol"`
	Company       string `json:"company"`
}

type RebalanceResponse struct {
	ClientId        string            `json:"clientId"`
	PocketId        string            `json:"pocketId"`
	BuyRequirements []PocketsMetaData `json:"buyRequirements"` // List of stocks to buy with required quantities
	Message         string            `json:"message"`
	Action          string            `json:"action"`    // eg: repair or rebalance
	OrderSide       string            `json:"orderSide"` // BUY or SELL
	PocketVersion   int               `json:"pocketVersion"`
}

type BuyReqStocks struct {
	OrderCompletedPrice float64           `json:"orderCompletedPrice"`
	OrderCompleted      []PocketsMetaData `json:"orderCompleted"`
	OrderCancelled      []PocketsMetaData `json:"orderCancelled"`
}

type FetchPocketPortfolioResponseV2 struct {
	ClientId         string                   `json:"clientId"`
	PortfolioDetails []PocketPortfolioDetails `json:"portfolioDetails"`
}
type PocketPortfolioDetails struct {
	PocketId               string            `json:"pocketId"`
	PocketName             string            `json:"pocketName"`
	PocketShortDesc        string            `json:"pocketShortDesc" bson:"pocketShortDesc"`
	PocketLongDesc         string            `json:"pocketLongDesc" bson:"pocketLongDesc"`
	PocketExchange         string            `json:"pocketExchange" bson:"pocketExchange"`
	PocketImage            string            `json:"pocketImage"`
	PocketWebImage         string            `json:"pocketWebImage"`
	PrimaryBackgroundColor string            `json:"primaryBackgroundColor"`
	PrimarySecondaryColor  string            `json:"primarySecondaryColor"`
	PocketVersion          int               `json:"pocketVersion" bson:"pocketVersion"`
	UsersPocketVersion     int               `json:"usersPocketVersion"`
	PocketBenchMark        string            `json:"pocketBenchMark" bson:"pocketBenchMark"`
	PocketCreateTimeUnix   int64             `json:"pocketCreateTimeUnix" bson:"pocketCreateTimeUnix"`
	AveragePrice           float64           `json:"averagePrice"`
	TotalBuyPockets        int               `json:"totalBuyPockets"`
	TotalSellPockets       int               `json:"totalSellPockets"`
	TotalInvestment        float64           `json:"totalInvestment"`
	LatestPocketTokens     []PocketsMetaData `json:"pocketTokens"`
}

type ProcessOrderReq struct {
	PocketId          string `json:"pocketId"`
	ClientId          string `json:"clientId"`
	LotSize           int    `json:"lotSize"`
	TransactionStatus int    `json:"transactionStatus"`
	Action            string `json:"orderSide"`
}

type CheckActionRequiredReq struct {
	ClientId     string `json:"clientId"`
	PocketId     string `json:"pocketId"`
	UserVersion  int    `json:"userVersion"`
	UsersLotSize int    `json:"userslotSize"`
}

type PocketOrdersRes struct {
	OrderCompleted []PocketsMetaData `json:"orderCompleted"`
	OrderCancelled []PocketsMetaData `json:"orderCancelled"`
}
