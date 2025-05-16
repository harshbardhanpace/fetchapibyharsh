package models

type FetchBondDataReq struct {
	Isin string `json:"isin" validate:"required"`
}

type NseBondStoreDb struct {
	Data []NseBondStoreDbData `json:"data"`
}

type NseBondStoreDbData struct {
	Symbol       string `json:"symbol"`
	Series       string `json:"series"`
	BondType     string `json:"bondType"`
	Open         string `json:"bondOpen"`
	High         string `json:"high"`
	Low          string `json:"low"`
	LtP          string `json:"ltP"`
	Close        string `json:"bondClose"`
	Per          string `json:"per"`
	Qty          string `json:"qty"`
	TrdVal       string `json:"trdVal"`
	Coupr        string `json:"coupr"`
	CreditRating string `json:"creditRating"`
	RatingAgency string `json:"ratingAgency"`
	FaceValue    string `json:"faceValue"`
	NxtipDate    string `json:"nxtipDate"`
	MaturityDate string `json:"maturityDate"`
	BYield       string `json:"bYield"`
	Isin         string `json:"isin"`
	CompanyName  string `json:"companyName"`
	Industry     string `json:"industry"`
	IsFNOSec     bool   `json:"isFNOSec"`
	IsCASec      bool   `json:"isCASec"`
	IsSLBSec     bool   `json:"isSLBSec"`
	IsDebtSec    bool   `json:"isDebtSec"`
	IsSuspended  bool   `json:"isSuspended"`
	IsETFSec     bool   `json:"isETFSec"`
	IsDelisted   bool   `json:"isDelisted"`
}
