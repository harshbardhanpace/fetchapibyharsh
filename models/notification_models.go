package models

type FetchAdminMessageRequest struct {
	ClientId string `json:"clientId"`
}

type FetchAdminMessageRes struct {
	Updates []FetchAdminMessageUpdates `json:"updates"`
}

type FetchAdminMessageUpdates struct {
	Type            string `json:"type"`
	UpdateEntryTime string `json:"update_entry_time"`
	Message         string `json:"message"`
	Platform        string `json:"platform"`
	Title           string `json:"title"`
}

type NotificationUpdatesReq struct {
	ClientId string `json:"clientId"`
}

type NotificationUpdatesRes struct {
	AllUpdate []NotificationUpdatesResUpdate `json:"allUpdate"`
	Status    string                         `json:"status"`
}

type NotificationUpdatesResUpdate struct {
	Type             string    `json:"type"`
	UpdateEntryTime  string    `json:"updateEntryTime"`
	UpdateID         string    `json:"updateId"`
	AlertID          int       `json:"alertId"`
	Condition        string    `json:"condition"`
	Exchange         string    `json:"exchange"`
	Expiry           int       `json:"expiry"`
	Frequency        string    `json:"frequency"`
	GeneratedAt      int       `json:"generatedAt"`
	InstrumentCode   string    `json:"instrumentCode"`
	LotSize          string    `json:"lotSize"`
	NewValue         float64   `json:"newValue"`
	StateAfterExpiry string    `json:"stateAfterExpiry"`
	TradingSymbol    string    `json:"tradingSymbol"`
	UserMessage      string    `json:"userMessage"`
	UserSetValues    []float64 `json:"userSetValues"`
}
