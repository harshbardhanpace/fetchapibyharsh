package models

type Instrument struct {
	InstrumentToken int    `json:"instrumentToken"`
	Exchange        string `json:"exchange"`
	Total           int    `json:"total"`
	Authorized      int    `json:"authorized"`
}

type EdisReq struct {
	ClientID    string       `json:"clientId"`
	Instruments []Instrument `json:"instruments"`
	RequestType string       `json:"requestType"`
}

type DataRes struct {
	Depository    string `json:"depository"`
	DpId          string `json:"dpId"`
	EncryptedDtls string `json:"encryptedDtls"`
	RequestId     string `json:"requestId"`
	Version       string `json:"version"`
}

type EdisRes struct {
	Data    DataRes `json:"data"`
	Html    bool    `json:"html"`
	Message string  `json:"message"`
	Status  string  `json:"status"`
}

type TpinReq struct {
	ClientID string `json:"clientId"`
	Boid     string `json:"boid"`
	Pan      string `json:"pan"`
	ReqFlag  string `json:"ReqFlag"`
	ReqTime  string `json:"ReqTime"`
}

type TpinRes struct {
	TlRespMessage string `json:"tlRespMessage"`
}
