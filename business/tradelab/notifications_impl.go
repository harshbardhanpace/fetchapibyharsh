package tradelab

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
	"strconv"
	"time"
)

type NotificationsObj struct {
	tradeLabURL string
}

func InitNotificationsObj() NotificationsObj {
	defer models.HandlePanic()

	notificationsObj := NotificationsObj{
		tradeLabURL: constants.TLURL,
	}

	return notificationsObj
}

func (obj NotificationsObj) FetchAdminMessages(req models.FetchAdminMessageRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + ADMINMESSAGE + "?type=" + constants.ADMINMESSAGE + "&client_id=" + url.QueryEscape(req.ClientId)

	//make payload
	payload := new(bytes.Buffer)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "FetchAdminMessages", duration, req.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " FetchAdminMsg call api error =", err, " clientId: ", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("FetchAdminMsg res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " clientId: ", req.ClientId, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlFetchAdminMessageResponse := TradelabFetchAdminMessageRes{}
	json.Unmarshal([]byte(string(body)), &tlFetchAdminMessageResponse)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " FetchAdminMsg tl status not ok StatusCode: ", res.StatusCode, " clientId: ", req.ClientId, " requestId:", reqH.RequestId)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	var adminMessagesRes models.FetchAdminMessageRes

	var adminMsgs []models.FetchAdminMessageUpdates

	for i := 0; i < len(tlFetchAdminMessageResponse.Data.Updates); i++ {
		var adminMessage models.FetchAdminMessageUpdates
		adminMessage.Type = tlFetchAdminMessageResponse.Data.Updates[i].Type
		adminMessage.UpdateEntryTime = tlFetchAdminMessageResponse.Data.Updates[i].UpdateEntryTime
		adminMessage.Message = tlFetchAdminMessageResponse.Data.Updates[i].Message
		adminMessage.Platform = tlFetchAdminMessageResponse.Data.Updates[i].Platform
		adminMessage.Title = tlFetchAdminMessageResponse.Data.Updates[i].Title
		adminMsgs = append(adminMsgs, adminMessage)
	}

	adminMessagesRes.Updates = adminMsgs

	loggerconfig.Info("FetchAdminMessages resp=", helpers.LogStructAsJSON(adminMessagesRes), " StatusCode: ", res.StatusCode, " clientId: ", req.ClientId, " requestId:", reqH.RequestId)

	apiRes.Data = adminMessagesRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}

func (obj NotificationsObj) NotificationUpdates(req models.NotificationUpdatesReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + ADMINMESSAGELATEST + "?type=" + constants.ALERT + "&client_id=" + url.QueryEscape(req.ClientId)

	//make payload
	payload := new(bytes.Buffer)

	//call api
	var apiRes apihelpers.APIRes
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	if err != nil {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " NotificationUpdates call api error =", err, " clientId: ", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("NotificationUpdates res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " clientId: ", req.ClientId, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlNotificationUpdates := TradelabNotificationUpdates{}
	json.Unmarshal((body), &tlNotificationUpdates)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " NotificationUpdates tl status not ok StatusCode: ", res.StatusCode, " clientId: ", req.ClientId, " requestId:", reqH.RequestId)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	var notificationUpdatesRes models.NotificationUpdatesRes

	var allnotificationUpdates []models.NotificationUpdatesResUpdate

	for i := 0; i < len(tlNotificationUpdates.Data.Updates); i++ {
		var notificationUpdates models.NotificationUpdatesResUpdate
		notificationUpdates.Type = tlNotificationUpdates.Data.Updates[i].Type
		notificationUpdates.UpdateEntryTime = tlNotificationUpdates.Data.Updates[i].UpdateEntryTime
		notificationUpdates.UpdateID = tlNotificationUpdates.Data.Updates[i].UpdateID
		notificationUpdates.AlertID = tlNotificationUpdates.Data.Updates[i].AlertID
		notificationUpdates.Condition = tlNotificationUpdates.Data.Updates[i].Condition
		notificationUpdates.Exchange = tlNotificationUpdates.Data.Updates[i].Exchange
		notificationUpdates.Expiry = tlNotificationUpdates.Data.Updates[i].Expiry
		notificationUpdates.Frequency = tlNotificationUpdates.Data.Updates[i].Frequency
		notificationUpdates.GeneratedAt = tlNotificationUpdates.Data.Updates[i].GeneratedAt
		notificationUpdates.InstrumentCode = tlNotificationUpdates.Data.Updates[i].InstrumentCode
		notificationUpdates.LotSize = tlNotificationUpdates.Data.Updates[i].LotSize
		notificationUpdates.NewValue = tlNotificationUpdates.Data.Updates[i].NewValue
		notificationUpdates.StateAfterExpiry = tlNotificationUpdates.Data.Updates[i].StateAfterExpiry
		notificationUpdates.TradingSymbol = tlNotificationUpdates.Data.Updates[i].TradingSymbol
		notificationUpdates.UserMessage = tlNotificationUpdates.Data.Updates[i].UserMessage
		notificationUpdates.UserSetValues = tlNotificationUpdates.Data.Updates[i].UserSetValues

		allnotificationUpdates = append(allnotificationUpdates, notificationUpdates)
	}

	notificationUpdatesRes.AllUpdate = allnotificationUpdates
	notificationUpdatesRes.Status = tlNotificationUpdates.Status

	loggerconfig.Info("NotificationUpdates resp=", helpers.LogStructAsJSON(notificationUpdatesRes), " StatusCode: ", res.StatusCode, " clientId: ", req.ClientId, " requestId:", reqH.RequestId)

	apiRes.Data = notificationUpdatesRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}
