package helpers

import (
	"space/constants"
	"space/loggerconfig"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// RecordAPILatency logs the details of an API call including its latency
func RecordAPILatency(url string, apiName string, duration time.Duration, clientId, requestID string) {
	latencyMs := duration.Milliseconds()
	logrus.WithFields(logrus.Fields{
		"url":       url,
		"apiName":   apiName,
		"latencyMS": latencyMs,
		"requestID": requestID,
	}).Info("API call completed")

	if latencyMs >= constants.LatencyThresholdHigh {
		loggerconfig.Error("Alert Severity:LatencyP0-Critical: " + apiName + " API is taking more than 5 seconds. clientId: " + clientId +
			" Request ID: " + requestID + ", Latency: " + strconv.FormatInt(latencyMs, 10) + "ms")
	} else if latencyMs >= constants.LatencyThresholdLow {
		loggerconfig.Error("Alert Severity:LatencyP1-High: " + apiName + " API is taking more than 1 second. clientId: " + clientId +
			"Request ID: " + requestID + ", Latency: " + strconv.FormatInt(latencyMs, 10) + "ms")
	}

}
