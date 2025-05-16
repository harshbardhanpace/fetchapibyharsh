package models

import "time"

type ReqHeader struct {
	Authorization   string `header:"Authorization"`
	OperatingSystem string `header:"P-Operating-System"`
	AppName         string `header:"P-Appname"`
	DeviceType      string `header:"P-DeviceType"`
	Platform        string `header:"P-Platform"`
	ClientPublicIP  string `header:"P-ClientPublicIP"`
	RequestId       string `header:"P-RequestId"`
	ClientId        string `header:"clientId"`
	ClientType      string `header:"P-ClientType"`
	DeviceId        string `header:"P-DeviceId"`
	FCMToken        string `header:"DeviceToken"`
	AdminRequestKey string `header:"P-AdminRequestKey"`
	ClientVersion   string `header:"P-ClientVersion"`
}

type Pong struct {
	DT time.Time `json:"time"`
}

type TokenHeaders struct {
	BlacklistKey string `json:"blacklist_key"`
	ClientID     string `json:"client_id"`
	ClientToken  string `json:"client_token"`
	Device       string `json:"device"`
	DeviceID     string `json:"device_id"`
	IP           string `json:"ip"`
	Exp          int64  `json:"exp"`
}
