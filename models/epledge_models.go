package models

import "time"

type Isin struct {
	IsinName string `json:"isinName"`
	Isin     string `json:"isin"`
	Quantity string `json:"quantity"`
	Price    string `json:"price"`
}

type EpledgeReq struct {
	Depository  string `json:"depository"`
	ClientID    string `json:"clientId"`
	Exchange    string `json:"exchange"`
	BoId        string `json:"boId"`
	Segment     string `json:"segment"`
	IsinDetails []Isin `json:"isinDetails"`
}
type Unpledge struct {
	Exchange string `json:"exchange"`
	Isin     string `json:"isin"`
	Quantity int64  `json:"quantity"`
}

type UnpledgeReq struct {
	ClientID     string     `json:"clientId"`
	Segment      string     `json:"segment"`
	UnpledgeList []Unpledge `json:"unpledgeList"`
}

type EpledgeDataRes struct {
	DpId       string `json:"dpid"`
	PledgedTls string `json:"pledgedtls"`
	ReqId      string `json:"reqid"`
	Version    string `json:"version"`
}

type EpledgeRes struct {
	Data    EpledgeDataRes `json:"data"`
	Message string         `json:"message"`
	Status  string         `json:"status"`
}
type ErrorRes struct {
	Code    string `json:"code"`
	Message int    `json:"message"`
}

type UnpledgeRes struct {
	Isin    string      `json:"isin"`
	Error   ErrorRes    `json:"error"`
	Result  string      `json:"result"`
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PledgeData struct {
	ID             int64     `json:"id"  db:"id"`
	ClientID       string    `json:"clientId"  db:"client_id"`
	SegmentID      string    `json:"segmentId"  db:"segment_id"`
	Timestamp      time.Time `json:"timestamp"  db:"timestamp"`
	ISIN           string    `json:"isin"  db:"isin"`
	Quantity       string    `json:"quantity"  db:"quantity"`
	Price          string    `json:"price" db:"price"`
	Exchange       string    `json:"exchange"  db:"exchange"`
	BOID           string    `json:"boId"  db:"bo_id"`
	Depository     string    `json:"depository"  db:"depository"`
	PledgeUnpledge string    `json:"pledgeUnpledge"  db:"pledge_unpledge"`
	DPID           string    `json:"dpId"  db:"dp_id"`
	PledgeTLS      string    `json:"pledgeTls"  db:"pledge_tls"`
	ReqID          string    `json:"reqId"  db:"req_id"`
	Version        string    `json:"version"  db:"version"`
	Status         string    `json:"status"  db:"status"`
}

type FetchEpledgeTxnReq struct {
	ClientID  string `json:"clientId"`
	Page      int    `json:"page"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type MTFCTDDataReq struct {
	PageSize  string `json:"pageSize" example:"100000"`
	PageNo    string `json:"pageNo" example:"1"`
	SortOrder string `json:"sortOrder" example:"ASC_UPDATED_AT"`
	FirstPage string `json:"firstPage" example:"true"`
	ClientId  string `json:"clientId" example:"AB0001"`
}

type MTFCTDDataRes struct {
	List       []CTDList `json:"list"`
	TotalCount int       `json:"totalCount"`
}

type CTDList struct {
	ClientID             string  `json:"clientId"`
	Isin                 string  `json:"isin"`
	TotalPledgeQuantity  int     `json:"totalPledgeQuantity"`
	CtdQuantity          int     `json:"ctdQuantity"`
	Symbol               string  `json:"symbol"`
	AvgPrice             float64 `json:"avgPrice"`
	MarginMultiplier     int     `json:"marginMultiplier"`
	CtdMarginValue       float64 `json:"ctdMarginValue"`
	Token                int     `json:"token"`
	Exchange             string  `json:"exchange"`
	CreatedAt            string  `json:"createdAt"`
	UpdatedAt            string  `json:"updatedAt"`
	EdisApprovedQuantity int     `json:"edisApprovedQuantity"`
	ObligationQuantity   int     `json:"obligationQuantity"`
	UsedQuantity         int     `json:"usedQuantity"`
	LoginID              string  `json:"loginId"`
	MarginValue          int     `json:"marginValue"`
	TotalInvestedAmount  int     `json:"totalInvestedAmount"`
	BrokerAmount         int     `json:"brokerAmount"`
}

type TLMtfPledgeListRes struct {
	List []MTFPledgeList `json:"list"`
}

type MTFPledgeList struct {
	ClientID            string    `json:"clientId"`
	Isin                string    `json:"isin"`
	PledgeQuantity      int       `json:"pledgeQuantity"`
	ToBePledgedQuantity int       `json:"toBePledgedQuantity"`
	Segment             string    `json:"segment"`
	Symbol              string    `json:"symbol"`
	MtfSettlementDate   string    `json:"mtfSettlementDate"`
	MtfSquareOffDate    string    `json:"mtfSquareOffDate"`
	NseToken            int       `json:"nseToken"`
	BseToken            int       `json:"bseToken"`
	AvgPrice            int       `json:"avgPrice"`
	MarginMultiplier    int       `json:"marginMultiplier"`
	MarginVarElm        int       `json:"marginVarElm"`
	MarginValue         float64   `json:"marginValue"`
	DaysTillSquareoff   int       `json:"daysTillSquareoff"`
	IsLastDayOfMtf      bool      `json:"isLastDayOfMtf"`
	IsCfObligation      bool      `json:"isCfObligation"`
	CreatedAt           time.Time `json:"CreatedAt"`
	UpdatedAt           time.Time `json:"UpdatedAt"`
}
