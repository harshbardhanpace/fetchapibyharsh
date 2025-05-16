package models

import "time"

// Get All Ipo Request
type GetAllIpoRequest struct {
}

// Get All Ipo Response
type GetAllIpoResponse struct {
	AllIpo struct {
		Data   []IpoState `json:"data"`
		Status string     `json:"status"`
	} `json:"allIpo"`
	OpenIpo     []IpoState `json:"openIpo"`
	UpcomingIpo []IpoState `json:"upcomingIpo"`
	ClosedIpo   []IpoState `json:"closedIpo"`
}

type IpoState struct {
	BiddingStartDate           string                        `json:"biddingStartDate"`
	Symbol                     string                        `json:"symbol"`
	MinBidQuantity             int                           `json:"minBidQuantity"`
	Registrar                  string                        `json:"registrar"`
	LotSize                    int                           `json:"lotSize"`
	T1ModEndDate               string                        `json:"t1ModEndDate"`
	DailyStartTime             string                        `json:"dailyStartTime"`
	T1ModStartTime             string                        `json:"t1ModStartTime"`
	BiddingEndDate             string                        `json:"biddingEndDate"`
	T1ModEndTime               string                        `json:"t1ModEndTime"`
	DailyEndTime               string                        `json:"dailyEndTime"`
	TickSize                   float64                       `json:"tickSize"`
	IssueType                  string                        `json:"issueType"`
	FaceValue                  float64                       `json:"faceValue"`
	MinPrice                   float64                       `json:"minPrice"`
	T1ModStartDate             string                        `json:"t1ModStartDate"`
	Name                       string                        `json:"name"`
	IssueSize                  int                           `json:"issueSize"`
	MaxPrice                   float64                       `json:"maxPrice"`
	CutOffPrice                float64                       `json:"cutOffPrice"`
	UnixBiddingEndDate         int                           `json:"unixBiddingEndDate"`
	UnixBiddingStartDate       int                           `json:"unixBiddingStartDate"`
	Isin                       string                        `json:"isin"`
	AllotmentDate              string                        `json:"allotmentDate"`
	ExchangeIssueType          string                        `json:"exchangeIssueType"`
	AllotmentBegins            string                        `json:"allotmentBegins"`
	RefundDate                 string                        `json:"refundDate"`
	ListingDate                string                        `json:"listingDate"`
	AboutCompany               string                        `json:"aboutCompany"`
	ParentCompany              string                        `json:"parentCompany"`
	FoundedYear                string                        `json:"foundedYear"`
	ProspectusFileURL          string                        `json:"prospectusFileUrl"`
	ManagingDirector           string                        `json:"managingDirector"`
	MaxLimit                   float64                       `json:"MaxLimit"`
	RetailDiscount             float64                       `json:"RetailDiscount"`
	NseExchangeListed          bool                          `json:"nseExchangeListed"`
	BseExchangeListed          bool                          `json:"bseExchangeListed"`
	AmoOrderEntryTime          string                        `json:"amoOrderEntryTime"`
	ApplicationRangeStart      int                           `json:"applicationRangeStart"`
	ApplicationRangeEnd        int                           `json:"applicationRangeEnd"`
	TotalApplicationRangeCount int                           `json:"totalApplicationRangeCount"`
	CategoryDetails            any                           `json:"categoryDetails"`
	SubCategorySettings        []IpoStateSubCategorySettings `json:"subCategorySettings"`
	IpoAllowed                 bool                          `json:"ipoAllowed"`
	BseAllowed                 bool                          `json:"bseAllowed"`
	NseAllowed                 bool                          `json:"nseAllowed"`
	SubType                    string                        `json:"subType"`
	EnablePio                  bool                          `json:"enablePio"`
	PioStartDate               time.Time                     `json:"pioStartDate"`
	PioEndDate                 time.Time                     `json:"pioEndDate"`
	PioEndTime                 time.Time                     `json:"pioEndTime"`
	PioStartTime               time.Time                     `json:"pioStartTime"`
	DematTransferDate          string                        `json:"dematTransferDate"`
	MandateEndDate             string                        `json:"mandateEndDate"`
	IsEmployeeCat              bool                          `json:"isEmployeeCat"`
	IsShareHolderCat           bool                          `json:"isShareHolderCat"`
}

type IpoStateSubCategorySettings struct {
	SubCatCode    string      `json:"subCatCode"`
	MinValue      interface{} `json:"minValue"`
	MaxUpiLimit   int         `json:"maxUpiLimit"`
	AllowCutOff   bool        `json:"allowCutOff"`
	AllowUpi      bool        `json:"allowUpi"`
	MaxValue      interface{} `json:"maxValue"`
	DiscountPrice interface{} `json:"discountPrice"`
	DiscountType  string      `json:"discountType"`
	MaxPrice      interface{} `json:"maxPrice"`
	CaCode        string      `json:"caCode"`
	Allowed       bool        `json:"allowed"`
	StartDate     string      `json:"startDate"`
	EndDate       string      `json:"endDate"`
	DisplayName   string      `json:"displayName"`
	MinLotSize    int         `json:"minLotSize"`
	StartTime     string      `json:"startTime"`
	EndTime       string      `json:"endTime"`
}

type PlaceIpoOrderRequest struct {
	ClientID      string                     `json:"clientId"`
	Symbol        string                     `json:"symbol"`
	UpiID         string                     `json:"upiId" mask:"id"`
	Bids          []PlaceIpoOrderRequestBids `json:"bids"`
	AllotmentMode string                     `json:"allotmentMode"`
	BankAccount   string                     `json:"bankAccount"`
	BankCode      string                     `json:"bankCode"`
	Broker        string                     `json:"broker"`
	CategoryCode  string                     `json:"categoryCode"`
	ClientBenID   string                     `json:"clientBenId"`
	ClientName    string                     `json:"clientName"`
	DpID          string                     `json:"dpId"`
	Ifsc          string                     `json:"ifsc" mask:"id"`
	LocationCode  string                     `json:"locationCode"`
	NonAsba       bool                       `json:"nonAsba"`
	Pan           string                     `json:"pan" mask:"id"`
	Category      string                     `json:"category"`
}

type PlaceIpoOrderRequestBids struct {
	ActivityType string `json:"activityType"`
	Quantity     int    `json:"quantity" validate:"gt=0"`
	AtCutOff     bool   `json:"atCutOff"`
	Price        int    `json:"price" validate:"gte=0"`
	Amount       int    `json:"amount"`
}

type PlaceIpoOrderResponse struct {
	Data interface{} `json:"data"`
}

type FetchIpoOrderRequest struct {
	ClientID string `json:"clientId"`
}

type FetchIpoOrderResponse struct {
	Data []FetchIpoOrderResponseData `json:"data" mask:"struct"`
}

type FetchIpoOrderResponseData struct {
	Symbol                          string                      `json:"symbol"`
	Reason                          string                      `json:"reason"`
	ApplicationNumber               string                      `json:"applicationNumber"`
	ClientName                      string                      `json:"clientName"`
	ChequeNumber                    string                      `json:"chequeNumber" mask:"id"`
	ReferenceNumber                 string                      `json:"referenceNumber"`
	DpVerStatusFlag                 string                      `json:"dpVerStatusFlag"`
	SubBrokerCode                   string                      `json:"subBrokerCode"`
	Depository                      string                      `json:"depository"`
	ReasonCode                      int                         `json:"reasonCode"`
	Pan                             string                      `json:"pan" mask:"id"`
	Ifsc                            string                      `json:"ifsc" mask:"id"`
	Timestamp                       string                      `json:"timestamp"`
	BankAccount                     string                      `json:"bankAccount" mask:"id"`
	BankCode                        string                      `json:"bankCode"`
	DpVerReason                     string                      `json:"dpVerReason"`
	DpID                            string                      `json:"dpId"`
	Upi                             string                      `json:"upi" mask:"id"`
	UpiAmtBlocked                   interface{}                 `json:"upiAmtBlocked"`
	Bids                            []FetchIpoOrderResponseBids `json:"bids"`
	AllotmentMode                   string                      `json:"allotmentMode"`
	DpVerFailCode                   string                      `json:"dpVerFailCode"`
	NonASBA                         bool                        `json:"nonASBA"`
	UpiFlag                         string                      `json:"upiFlag"`
	Category                        string                      `json:"category"`
	LocationCode                    string                      `json:"locationCode"`
	ClientBenID                     string                      `json:"clientBenId"`
	ClientID                        string                      `json:"clientId"`
	Status                          string                      `json:"status"`
	Mode                            string                      `json:"mode"`
	AllotmentStatus                 string                      `json:"allotmentStatus"`
	AllotmentDate                   string                      `json:"allotmentDate"`
	AllotmentUpdated                string                      `json:"allotmentUpdated"`
	AllotmentQuantity               int                         `json:"allotmentQuantity"`
	AllotmentPrice                  float64                     `json:"allotmentPrice"`
	CategoryCode                    string                      `json:"categoryCode"`
	CategoryDisplayName             string                      `json:"categoryDisplayName"`
	IsAmoOrder                      bool                        `json:"isAmoOrder"`
	PaymentMode                     string                      `json:"paymentMode"`
	AmtBlockTime                    string                      `json:"amtBlockTime"`
	Modify                          bool                        `json:"modify"`
	IsOrderModify                   bool                        `json:"isOrderModify"`
	IsBseIpo                        bool                        `json:"isBseIpo"`
	IsNseIpo                        bool                        `json:"isNseIpo"`
	UpiPaymentStatusMessage         string                      `json:"upiPaymentStatusMessage"`
	ExchangeUpdatedUpiBlockedAmount int                         `json:"exchangeUpdatedUpiBlockedAmount"`
	IsPioOrder                      bool                        `json:"isPioOrder"`
}

type FetchIpoOrderResponseBids struct {
	AtCutOff           bool    `json:"atCutOff"`
	Amount             int     `json:"amount"`
	Quantity           int     `json:"quantity"`
	BidReferenceNumber int64   `json:"bidReferenceNumber"`
	Series             string  `json:"series"`
	Price              float64 `json:"price"`
	ActivityType       string  `json:"activityType"`
	Status             string  `json:"status"`
}

type CancelIpoOrderRequest struct {
	ClientID          string                      `json:"clientId"`
	Symbol            string                      `json:"symbol"`
	ApplicationNumber string                      `json:"applicationNumber"`
	Bids              []CancelIpoOrderRequestBids `json:"bids"`
	UpiID             string                      `json:"upiId" mask:"id"`
	AllotmentMode     string                      `json:"allotmentMode"`
	BankAccount       string                      `json:"bankAccount" mask:"id"`
	BankCode          string                      `json:"bankCode"`
	Broker            string                      `json:"broker"`
	ClientBenID       string                      `json:"clientBenId"`
	ClientName        string                      `json:"clientName"`
	DpID              string                      `json:"dpId"`
	Ifsc              string                      `json:"ifsc" mask:"id"`
	LocationCode      string                      `json:"locationCode"`
	NonAsba           bool                        `json:"nonAsba"`
	Pan               string                      `json:"pan" mask:"id"`
}

type CancelIpoOrderRequestBids struct {
	Quantity           int     `json:"quantity" validate:"gt=0"`
	AtCutOff           bool    `json:"atCutOff"`
	Price              float64 `json:"price" validate:"gte=0"`
	Amount             int     `json:"amount"`
	BidReferenceNumber int64   `json:"bidReferenceNumber"`
	Series             string  `json:"series"`
	ActivityType       string  `json:"activityType"`
	Status             string  `json:"status"`
}

type CancelIpoOrderResponse struct {
	Data interface{} `json:"data"`
}

type FetchIpoDataRequest struct {
	Name string `json:"name" validate:"required"`
}

type FetchIpoDataResponse struct {
	Name            string   `json:"name" bson:"name"`
	Timetable       []string `json:"timetable" bson:""`
	Details         []string `json:"details" bson:"details"`
	Reservations    []string `json:"reservations" bson:"reservations"`
	LotsAvailable   []string `json:"lotsAvailable" bson:"lots_available"`
	PermotorHolding []string `json:"permotorHolding" bson:"permotor_holding"`
	CompanyFinance  []string `json:"companyFinance" bson:"company_finance"`
	Subscriptions   []string `json:"subscriptions" bson:"subscriptions"`
	Performance     []string `json:"performance" bson:"performance"`
	Documents       []string `json:"documents" bson:"documents"`
	IpoData         []string `json:"ipoData" bson:"ipo_data"`
	Overview        []string `json:"overview" bson:"overview"`
}

type FetchIpoGmpDataResponse struct {
	IpoName    string   `json:"ipoName" bson:"ipo_name"`
	GmpDetails []string `json:"gmpDetails" bson:"gmpDetails"`
}

type FetchEipoReq struct {
	IpoSymbol []string `json:"ipoSymbol" validate:"required,min=1,dive,required"`
	IpoStage  string   `json:"ipoStage"  validate:"required,alpha"`
}

type FetchNseIpoResponse struct {
	IpoSymbol  string      `json:"symbol"`
	IpoStage   string      `json:"stage"`
	Series     string      `json:"series"`
	DRHPLink   string      `json:"drhpLink"`
	BidDeatils interface{} `json:"bidDetails"`
	SubsTimes  string      `json:"subsTimes"` // No of times issue subscribed
}
