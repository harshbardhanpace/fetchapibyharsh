package tradelab

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
)

type AlertsObj struct {
	tradeLabURL string
}

func InitAlertsProvider() AlertsObj {
	defer models.HandlePanic()

	alertsObj := AlertsObj{
		tradeLabURL: constants.TLURL,
	}

	return alertsObj
}

func (obj AlertsObj) CreateAlert(req models.CreateAlertsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + ALERTSURL

	var tlCreateAlertsReq TradelabCreateAlertReq
	tlCreateAlertsReq.Exchange = req.Exchange
	tlCreateAlertsReq.Condition = req.Condition
	tlCreateAlertsReq.Expiry = req.Expiry
	tlCreateAlertsReq.Frequency = req.Frequency
	tlCreateAlertsReq.InstrumentToken = req.InstrumentToken
	tlCreateAlertsReq.StateAfterExpiry = req.StateAfterExpiry
	tlCreateAlertsReq.UserMessage = req.UserMessage
	tlCreateAlertsReq.UserSetValues = req.UserSetValues
	tlCreateAlertsReq.WaitTime = req.WaitTime

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlCreateAlertsReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "CreateAlert", duration, req.ClientId, reqH.RequestId)
	if err != nil {
		// logObj.Printf("BOOrderRes call api error =", err)
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " CreateAlert call api error", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("CreateAlert res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlCreateAlertsRes := TradelabCreateAlertRes{}
	json.Unmarshal([]byte(string(body)), &tlCreateAlertsRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " CreateAlert tl status not ok =", tlCreateAlertsRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlCreateAlertsRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var createAlertsRes models.CreateAlertsRes
	createAlertsRes.AlertID = tlCreateAlertsRes.Data.AlertID

	loggerconfig.Info("CreateAlert tl resp=", helpers.LogStructAsJSON(createAlertsRes), " uccId:", req.ClientId, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = createAlertsRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj AlertsObj) EditAlerts(req models.EditAlertsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + ALERTSURL + "/" + req.AlertId

	var tlCreateAlertsReq TradelabCreateAlertReq
	tlCreateAlertsReq.Exchange = req.Exchange
	tlCreateAlertsReq.Condition = req.Condition
	tlCreateAlertsReq.Expiry = req.Expiry
	tlCreateAlertsReq.Frequency = req.Frequency
	tlCreateAlertsReq.InstrumentToken = req.InstrumentToken
	tlCreateAlertsReq.StateAfterExpiry = req.StateAfterExpiry
	tlCreateAlertsReq.UserMessage = req.UserMessage
	tlCreateAlertsReq.UserSetValues = req.UserSetValues
	tlCreateAlertsReq.WaitTime = req.WaitTime

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlCreateAlertsReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPut, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "EditAlerts", duration, req.ClientId, reqH.RequestId)
	if err != nil {
		// logObj.Printf("BOOrderRes call api error =", err)
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " EditAlerts call api error", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("EditAlerts res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlEditAlertsRes := TradelabDeleteAlertsRes{}
	json.Unmarshal([]byte(string(body)), &tlEditAlertsRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " EditAlerts tl status not ok =", tlEditAlertsRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientId, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlEditAlertsRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("EditAlerts tl resp=", helpers.LogStructAsJSON(tlEditAlertsRes), " uccId:", req.ClientId, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = tlEditAlertsRes.Message
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj AlertsObj) GetAlerts(req models.GetAlertsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + ALERTSURL
	//make payload
	payload := new(bytes.Buffer) // empty payload

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "GetAlerts", duration, req.ClientId, reqH.RequestId)
	if err != nil {
		// logObj.Printf("BOOrderRes call api error =", err)
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " GetAlerts call api error", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("GetAlerts res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlGetAlertsRes := TradelabGetAlertsRes{}
	json.Unmarshal([]byte(string(body)), &tlGetAlertsRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " GetAlerts tl status not ok =", tlGetAlertsRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientId, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlGetAlertsRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	var getAlertsRes models.AlersRes

	alerts := make([]models.Alerts, 0)
	for _, tlAlert := range tlGetAlertsRes.Data {
		var alert models.Alerts
		alert.ClientID = tlAlert.ClientID
		alert.ConditionType = tlAlert.ConditionType
		alert.Exchange = tlAlert.Exchange
		alert.Expiry = tlAlert.Expiry
		alert.Frequency = tlAlert.Frequency
		alert.ID = tlAlert.ID
		alert.StateAfterExpiry = tlAlert.StateAfterExpiry
		alert.Status = tlAlert.Status
		alert.Token = tlAlert.Token
		alert.TradingSymbol = tlAlert.TradingSymbol
		alert.UserMessage = tlAlert.UserMessage
		alert.UserSetValues = tlAlert.UserSetValues
		alert.WaitTime = tlAlert.WaitTime
		alerts = append(alerts, alert)
	}
	getAlertsRes.Data = alerts

	loggerconfig.Info("GetAlerts tl resp=", helpers.LogStructAsJSON(getAlertsRes), " uccId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = getAlertsRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}

func (obj AlertsObj) PauseAlerts(req models.PauseAlertsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + ALERTSURL + "/" + strconv.Itoa(int(req.AlertId))

	var tlPauseAlertsReq TradelabPauseAlertsReq
	tlPauseAlertsReq.Status = req.Status

	//make payload
	payload := new(bytes.Buffer) // empty payload
	json.NewEncoder(payload).Encode(tlPauseAlertsReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPut, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "PauseAlerts", duration, req.ClientId, reqH.RequestId)
	if err != nil {
		// logObj.Printf("BOOrderRes call api error =", err)
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " PauseAlerts call api error", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("PauseAlerts res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " PauseAlerts tl status not ok =", res.StatusCode, " uccId:", req.ClientId, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = "FAILURE"
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	loggerconfig.Info("PauseAlerts tl success uccId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj AlertsObj) DeleteAlerts(req models.DeleteAlertsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + ALERTSURL + "/" + strconv.Itoa(int(req.AlertId))

	//make payload
	payload := new(bytes.Buffer) // empty payload

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodDelete, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "DeleteAlerts", duration, req.ClientId, reqH.RequestId)
	if err != nil {
		// logObj.Printf("BOOrderRes call api error =", err)
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " DeleteAlerts call api error", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("DeleteAlerts res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlDeleteAlertsRes := TradelabDeleteAlertsRes{}
	json.Unmarshal([]byte(string(body)), &tlDeleteAlertsRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " DeleteAlerts tl status not ok =", tlDeleteAlertsRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientId, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlDeleteAlertsRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	loggerconfig.Info("DeleteAlerts tl resp=", helpers.LogStructAsJSON(tlDeleteAlertsRes), " uccId:", req.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}
