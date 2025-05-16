package models

type MongoAdmin struct {
	UserId   string `json:"userId" bson:"userId"`
	Password string `json:"password" bson:"password"`
}

type MongoPocketsMetaData struct {
	Qty   string `json:"qty"`
	Token string `json:"token"`
}

type MongoPockets struct {
	PocketName             string            `json:"pocketName" bson:"pocketName"`
	PocketShortDesc        string            `json:"pocketShortDesc" bson:"pocketShortDesc"`
	PocketLongDesc         string            `json:"pocketLongDesc" bson:"pocketLongDesc"`
	PocketExchange         string            `json:"pocketExchange" bson:"pocketExchange"`
	PocketImage            string            `json:"pocketImage"`
	PocketWebImage         string            `json:"pocketWebImage"`
	PrimaryBackgroundColor string            `json:"primaryBackgroundColor"`
	PrimarySecondaryColor  string            `json:"primarySecondaryColor"`
	PocketTokens           []PocketsMetaData `json:"pocketTokens" bson:"pocketTokens"`
	PocketId               string            `json:"pocketId" bson:"pocketId"`
}

type MongoCollections struct {
	CollectionName         string               `json:"collectionName" bson:"collectionName"`
	CollectionShortDesc    string               `json:"collectionShortDesc" bson:"collectionShortDesc"`
	CollectionLongDesc     string               `json:"collectionLongDesc" bson:"collectionLongDesc"`
	CollectionExchange     string               `json:"collectionExchange" bson:"collectionExchange"`
	CollectionMetaData     []CollectionDataInfo `json:"collectionMetaData" bson:"collectionMetaData"`
	CollectionId           string               `json:"collectionId" bson:"collectionId"`
	CollectionBackColor    string               `json:"collectionBackColor"`
	CollectionImage        string               `json:"collectionImage"`
	CollectionIllustration string               `json:"collectionIllustration"` // small image
	CollectionType         string               `json:"collectionType"`         // basic/pro etc
}

type MongoWatchLists struct {
	ClientId               string `json:"clientId" bson:"clientId"`
	WatchListName          string `json:"watchListName" bson:"watchListName"`
	WatchListNameLongDesc  string `json:"watchListLongDesc" bson:"watchListLongDesc"`
	WatchListNameShortDesc string `json:"watchListShortDesc" bson:"watchListShortDesc"`
	WatchListId            string `json:"watchListId" bson:"watchListId"`
}

type MongoStocksWatchLists struct {
	ClientId    string       `json:"clientId" bson:"clientId"`
	WatchListId string       `json:"watchListId" bson:"watchListId"`
	Stock       StockDetails `json:"stock" bson:"stock"`
}

type MongoNewWatchLists struct {
	ClientId   string         `json:"clientId" bson:"clientId"`
	WatchList1 []StockDetails `json:"wl1" bson:"wl1"`
	WatchList2 []StockDetails `json:"wl2" bson:"wl2"`
	WatchList3 []StockDetails `json:"wl3" bson:"wl3"`
	WatchList4 []StockDetails `json:"wl4" bson:"wl4"`
	WatchList5 []StockDetails `json:"wl5" bson:"wl5"`
}

type MongoPinsMetadata struct {
	ClientId      string         `json:"clientId"`
	PinsMetaDatas []PinsMetaData `json:"pinsDetails"`
}

type MongoEmailUserId struct {
	EmailId  string `json:"emailId" bson:"emailId"`
	ClientId string `json:"clientId" bson:"clientId"`
}

type ConsentDetailsStore struct {
	ClientId      string `json:"clientId" bson:"clientId"`
	ConsentHandle string `json:"consentHandle" bson:"consentHandle"`
	CustomerId    string `json:"customerId" bson:"customerId"`
	ConsentStatus string `json:"consentStatus" bson:"consentStatus"`
}

type MongoBankStatementDetails struct {
	UserId                  string `json:"userid"`
	BankStatementS3Location string `json:"bankStatementS3Location"`
	Verified                bool   `json:"verified"`
	Rejection               string `json:"Rejection"`
}

type MongoFetchFinancials struct {
	CoCode     int                  `json:"coCode" bson:"coCode"`
	Isin       string               `json:"isin" bson:"isin"`
	Nsesymbol  string               `json:"nsesymbol" bson:"nsesymbol"`
	Bsecode    string               `json:"bsecode" bson:"bsecode"`
	Financials FetchFinancialsV4Res `json:"financials" bson:"financials"`
}

type UpiPreference struct {
	ClientId string   `json:"clientId" bson:"clientId"`
	UpiIds   []string `json:"upiIds" bson:"upiIds"`
}

type MongoNewWatchListsV2 struct {
	ClientId   string           `json:"clientId" bson:"clientId"`
	WatchList1 []StockDetailsV2 `json:"wl1" bson:"wl1"`
	WatchList2 []StockDetailsV2 `json:"wl2" bson:"wl2"`
	WatchList3 []StockDetailsV2 `json:"wl3" bson:"wl3"`
	WatchList4 []StockDetailsV2 `json:"wl4" bson:"wl4"`
	WatchList5 []StockDetailsV2 `json:"wl5" bson:"wl5"`
}

type MongoAccountFreeze struct {
	ClientId     string `json:"clientId" bson:"clientId"`
	OtpStatus    bool   `json:"otp"`
	FreezeStatus bool   `json:"freezeStatus"`
}

type MongoPocketsV3 struct {
	PocketName             string                 `json:"pocketName" bson:"pocketName"`
	PocketShortDesc        string                 `json:"pocketShortDesc" bson:"pocketShortDesc"`
	PocketLongDesc         string                 `json:"pocketLongDesc" bson:"pocketLongDesc"`
	PocketExchange         string                 `json:"pocketExchange" bson:"pocketExchange"`
	PocketImage            string                 `json:"pocketImage"`
	PocketWebImage         string                 `json:"pocketWebImage"`
	PrimaryBackgroundColor string                 `json:"primaryBackgroundColor"`
	PrimarySecondaryColor  string                 `json:"primarySecondaryColor"`
	PocketId               string                 `json:"pocketId" bson:"pocketId"`
	PocketLatestVersion    int                    `json:"pocketLatestVersion" bson:"pocketLatestVersion"`
	PocketBenchMark        string                 `json:"pocketBenchMark" bson:"pocketBenchMark"`
	PocketCreateTimeUnix   int64                  `json:"pocketCreateTimeUnix" bson:"pocketCreateTimeUnix"`
	PocketVersionDetails   []PocketVersionDetails `json:"pocketVersionDetails" bson:"pocketVersionDetails"`
	Tag                    string                 `json:"tag" bson:"tag"`
}

type PocketVersionDetails struct {
	Version                 int               `json:"version" bson:"version"`
	SupportDays             int               `json:"supportDays" bson:"supportDays"`
	FactSheet               string            `json:"factSheet" bson:"factSheet"`
	PocketVersionCreateTime int64             `json:"pocketVersionCreateTine" bson:"pocketVersionCreateTime"`
	PocketVersionUpdateTime int64             `json:"pocketVersionUpdateTine" bson:"pocketVersionUpdateTime"`
	Tokens                  []PocketsMetaData `json:"tokens" bson:"tokens"`
	PocketStatus            string            `json:"pocketStatus" bson:"pocketStatus"`
}

type MongoLatestPocketDetails struct {
	PocketTokens           []PocketsMetaData `json:"pocketTokens" bson:"pocketTokens"`
	PocketName             string            `json:"pocketName" bson:"pocketName"`
	PocketShortDesc        string            `json:"pocketShortDesc" bson:"pocketShortDesc"`
	PocketLongDesc         string            `json:"pocketLongDesc" bson:"pocketLongDesc"`
	PocketExchange         string            `json:"pocketExchange" bson:"pocketExchange"`
	PocketImage            string            `json:"pocketImage"`
	PocketWebImage         string            `json:"pocketWebImage"`
	PrimaryBackgroundColor string            `json:"primaryBackgroundColor"`
	PrimarySecondaryColor  string            `json:"primarySecondaryColor"`
	PocketId               string            `json:"pocketId" bson:"pocketId"`
	PocketVersion          int               `json:"pocketVersion" bson:"pocketVersion"`
	PocketBenchMark        string            `json:"pocketBenchMark" bson:"pocketBenchMark"`
	PocketCreateTimeUnix   int64             `json:"pocketCreateTimeUnix" bson:"pocketCreateTimeUnix"`
}

type PocketTransactionStoreV3 struct {
	ClientId           string                   `json:"clientId" bson:"clientId"`
	AllPocketPurchases []PurchagedPocketDetails `json:"allPocketPurchases" bson:"allPocketPurchases"`
}

type PurchagedPocketDetails struct {
	PocketTransactionId string            `json:"pocketTransactionId" bson:"pocketTransactionId"`
	PocketId            string            `json:"pocketId" bson:"pocketId"`
	TransactionStatus   int               `json:"transactionStatus" bson:"transactionStatus"` // 0 - Bought, 1 - SOLD/Exited, 2 - for future needs
	LotSize             int               `json:"lotSize" bson:"lotSize"`
	Action              string            `json:"action" bson:"action"`
	OrderCompletedPrice float64           `json:"orderCompletedPrice" bson:"orderCompletedPrice"`
	OrderCompleted      []PocketsMetaData `json:"orderCompleted" bson:"orderCompleted"`
	OrderCancelled      []PocketsMetaData `json:"orderCancelled" bson:"orderCancelled"`
}
