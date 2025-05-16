package models

type CreateAlertsReq struct {
	ClientId         string    `json:"clientId"`
	Exchange         string    `json:"exchange"`
	InstrumentToken  string    `json:"instrumentToken"`
	WaitTime         string    `json:"waitTime"`
	Condition        string    `json:"condition"`
	UserSetValues    []float64 `json:"userSetValues"`
	Frequency        string    `json:"frequency"`
	Expiry           int       `json:"expiry"`
	StateAfterExpiry string    `json:"stateAfterExpiry"`
	UserMessage      string    `json:"userMessage"`
}

type CreateAlertsRes struct {
	AlertID int64 `json:"alertId"`
}

type EditAlertsReq struct {
	ClientId         string    `json:"clientId"`
	AlertId          string    `json:"alertId"`
	Exchange         string    `json:"exchange"`
	InstrumentToken  string    `json:"instrumentToken"`
	WaitTime         string    `json:"waitTime"`
	Condition        string    `json:"condition"`
	UserSetValues    []float64 `json:"userSetValues"`
	Frequency        string    `json:"frequency"`
	Expiry           int       `json:"expiry"`
	StateAfterExpiry string    `json:"stateAfterExpiry"`
	UserMessage      string    `json:"userMessage"`
}

type GetAlertsReq struct {
	ClientId string `json:"clientId"`
}

type Alerts struct {
	WaitTime         int       `json:"waitTime"`
	UserSetValues    []float64 `json:"userSetValues"`
	UserMessage      string    `json:"userMessage"`
	TradingSymbol    string    `json:"tradingSymbol"`
	Token            string    `json:"token"`
	Status           string    `json:"status"`
	StateAfterExpiry string    `json:"stateAfterExpiry"`
	ID               int       `json:"id"`
	Frequency        string    `json:"frequency"`
	Expiry           int       `json:"expiry"`
	Exchange         string    `json:"exchange"`
	ConditionType    string    `json:"conditionType"`
	ClientID         string    `json:"clientId"`
}

type AlersRes struct {
	Data []Alerts `json:"data"`
}

type PauseAlertsReq struct {
	ClientId string `json:"clientId"`
	AlertId  int64  `json:"alertId"`
	Status   string `json:"status"`
}

type DeleteAlertsReq struct {
	ClientId string `json:"clientId"`
	AlertId  int64  `json:"alertId"`
}
