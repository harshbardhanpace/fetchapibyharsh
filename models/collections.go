package models

type CreateCollectionsRequest struct {
	CollectionName         string               `json:"collectionName"`
	CollectionShortDesc    string               `json:"collectionShortDesc"`
	CollectionLongDesc     string               `json:"collectionLongDesc"`
	CollectionExchange     string               `json:"collectionExchange"`
	CollectionBackColor    string               `json:"collectionBackColor"`
	CollectionImage        string               `json:"collectionImage"`
	CollectionIllustration string               `json:"collectionIllustration"` // small image
	CollectionType         string               `json:"collectionType"`         // basic/pro etc
	CollectionMetaData     []CollectionDataInfo `json:"collectionMetaData"`
}

type CreateCollectionsResponse struct {
	CollectionId       string               `json:"collectionId"`
	CollectionName     string               `json:"collectionName"`
	CollectionType     string               `json:"collectionType"`
	CollectionMetaData []CollectionDataInfo `json:"collectionMetaData"`
}

type FetchCollectionsDetailsRequest struct {
	CollectionId string `json:"collectionId"`
}

type FetchCollectionsDetailsResponse struct {
	CollectionId           string               `json:"collectionId"`
	CollectionName         string               `json:"collectionName"`
	CollectionShortDesc    string               `json:"collectionShortDesc"`
	CollectionLongDesc     string               `json:"collectionLongDesc"`
	CollectionExchange     string               `json:"collectionExchange"`
	CollectionBackColor    string               `json:"collectionBackColor"`
	CollectionImage        string               `json:"collectionImage"`
	CollectionIllustration string               `json:"collectionIllustration"` // small image
	CollectionType         string               `json:"collectionType"`         // basic/pro etc
	CollectionMetaData     []CollectionDataInfo `json:"collectionMetaData"`
}

type ModifyCollectionsRequest struct {
	CollectionName         string               `json:"collectionName"`
	CollectionShortDesc    string               `json:"collectionShortDesc"`
	CollectionLongDesc     string               `json:"collectionLongDesc"`
	CollectionExchange     string               `json:"collectionExchange"`
	CollectionMetaData     []CollectionDataInfo `json:"collectionMetaData"`
	CollectionId           string               `json:"collectionId"`
	CollectionBackColor    string               `json:"collectionBackColor"`
	CollectionImage        string               `json:"collectionImage"`
	CollectionIllustration string               `json:"collectionIllustration"` // small image
	CollectionType         string               `json:"collectionType"`         // basic/pro etc
}

type ModifyCollectionsResponse struct {
	CollectionId       string               `json:"collectionId"`
	CollectionType     string               `json:"collectionType"` // basic/pro etc
	CollectionMetaData []CollectionDataInfo `json:"collectionMetaData"`
}

type DeleteCollectionsRequest struct {
	CollectionId string `json:"collectionId"`
}

type DeleteCollecionsResponse struct {
	CollectionMetaData []CollectionDataInfo `json:"collectionMetaData"`
}

type CollectionDataInfo struct {
	StockId       string `json:"stockId"`
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

type FetchAllCollectionsDetailsResponse struct {
	FetchAllCollectionsDetailsResponse []FetchCollectionsDetailsResponse `json:"fetchAllCollectionsDetailsResponse"`
}
