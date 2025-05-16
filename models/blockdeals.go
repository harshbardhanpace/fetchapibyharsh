package models

type BlockDeal struct {
	Cocode     int     `json:"cocode"`
	DealType   string  `json:"dealtype"`
	ScripCode  string  `json:"scripcode"`
	Serial     int     `json:"serial"`
	Date1      string  `json:"date1"`
	ScripName  string  `json:"scripname"`
	ClientName string  `json:"clientname"`
	BuySell    string  `json:"buysell"`
	QtyShares  float64 `json:"qtyshares"`
	AvgPrice   float64 `json:"avgprice"`
	UnixTime   int64   `json:"unixtime"`
}
