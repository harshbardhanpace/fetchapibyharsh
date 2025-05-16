package models

type FetchBondDataByIsinReq struct {
	Isin string `json:"isin" validate:"required,len=12,alphanum"`
}

type FetchBondDataByIsinResponse struct {
	CouponRate       string       `json:"couponRate"`
	CouponFrequency  string       `json:"couponFrequency"`
	Yield            string       `json:"yield"`
	Taxation         string       `json:"taxation"`
	Tenure           string       `json:"tenure"`
	MaturityDate     string       `json:"maturityDate"`
	Series           string       `json:"series"`
	TypeOfInstrument string       `json:"typeOfInstrument"`
	IssueSize        string       `json:"issueSize"`
	ISIN             string       `json:"isin"`
	CreditRating     RatingPacket `json:"creditRating"`
	Seniority        string       `json:"seniority"`
	Security         string       `json:"security"`
	IssuerName       string       `json:"issuerName"`
	Ownership        string       `json:"ownership"`
	BusinessSector   string       `json:"businessSector"`
	Source           string       `json:"source"`
	NseToken         string       `json:"nseToken"`
	BseToken         string       `json:"bseToken"`
}



type RatingPacket struct {
	CRISIL           string `json:"crisil"`
	ICRA             string `json:"icra"`
	CARE             string `json:"care"`
	ACUITE           string `json:"acuite"`
	BWR              string `json:"bwr"`
	IND              string `json:"ind"`
	BRICKWORK        string `json:"brickwork"`
	GovernmentBacked bool   `json:"governmentBacked"`
}
