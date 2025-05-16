package models

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

type NudgeAlertReq struct {
	Isin     string `json:"isin"`
	ClientId string `json:"clientId"`
}

type NudgeAlertRes struct {
	AsmPresent bool `json:"asmPresent"`
	GsmPresent bool `json:"gsmPresent"`
}

func HandlePanic() {
	if r := recover(); r != nil {
		logrus.Printf("Alert Severity:P0-Critical, PanicRecover RECOVERED from: %v\n", r)
		stackTrace := make([]byte, 1024)
		runtime.Stack(stackTrace, true)
		logrus.Printf("Alert Severity:P0-Critical, PanicRecover Stack trace: %s\n", string(stackTrace))
	}
}
