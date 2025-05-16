package models

type SessionInfoReq struct {
}

type SessionInfo struct {
	Exchange         string `json:"exchange"`
	SessionName      string `json:"sessionName"`
	AMOStartTime     string `json:"amoStartTime"`
	AMOEndTime       string `json:"amoEndTime"`
	IsActive         string `json:"isActive"`
	IsHoliday        string `json:"isHoliday"`
	BufferStartTime  string `json:"bufferStartTime"`
	BufferEndTime    string `json:"bufferEndTime"`
	MarketCloseTime  string `json:"marketCloseTime"`
	PostClosingStart string `json:"postClosingStart"`
	PostClosingEnd   string `json:"postClosingEnd"`
}

type SessionInfoRes struct {
	Result []SessionInfo `json:"result"`
}
