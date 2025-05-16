package models

type SetUpiPreferenceReq struct {
	ClientId string `json:"clientId"`
	UpiId    string `json:"upiId" mask:"id"`
}

type FetchUpiPreferenceReq struct {
	ClientId string `json:"clientId"`
}

type FetchUpiPreferenceRes struct {
	ClientId string   `json:"clientId"`
	UpiIds   []string `json:"upiIds"`
}

type DeleteUpiPreferenceReq struct {
	ClientId string   `json:"clientId"`
	UpiIds   []string `json:"upiIds" mask:"id"`
}
