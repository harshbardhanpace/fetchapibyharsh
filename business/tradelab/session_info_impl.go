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

type SessionInfoObj struct {
	tradeLabURL string
}

func InitSessionInfoProvider() SessionInfoObj {
	defer models.HandlePanic()

	sessionInfoObj := SessionInfoObj{
		tradeLabURL: constants.TLURL,
	}

	return sessionInfoObj
}

func (obj SessionInfoObj) SessionInfo(req models.SessionInfoReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + SESSIONINFOURL

	//make payload
	payload := new(bytes.Buffer)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallApi(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "SessionInfo", duration, reqH.ClientId, reqH.RequestId)
	defer res.Body.Close()
	if err != nil {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " SessionInfoReq call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("SessionInfoRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlSessionInfoRes := TradeLabSessionInfoResponse{}
	json.Unmarshal([]byte(string(body)), &tlSessionInfoRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " sessionInfoRes tl status not ok =", tlSessionInfoRes.Error.Message, "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlSessionInfoRes.Error.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	// fill up controller response
	var sessionInfoResponse models.SessionInfoRes

	for i := 0; i < len(tlSessionInfoRes.Result); i++ {
		var sessionInfo models.SessionInfo
		sessionInfo.Exchange = tlSessionInfoRes.Result[i].Exchange
		sessionInfo.SessionName = tlSessionInfoRes.Result[i].SessionName
		sessionInfo.AMOStartTime = tlSessionInfoRes.Result[i].AMOStartTime
		sessionInfo.AMOEndTime = tlSessionInfoRes.Result[i].AMOEndTime
		sessionInfo.IsActive = tlSessionInfoRes.Result[i].IsActive
		sessionInfo.IsHoliday = tlSessionInfoRes.Result[i].IsHoliday
		sessionInfo.BufferStartTime = tlSessionInfoRes.Result[i].BufferStartTime
		sessionInfo.BufferEndTime = tlSessionInfoRes.Result[i].BufferEndTime
		sessionInfo.MarketCloseTime = tlSessionInfoRes.Result[i].MarketCloseTime
		sessionInfo.PostClosingStart = tlSessionInfoRes.Result[i].PostClosingStart
		sessionInfo.PostClosingEnd = tlSessionInfoRes.Result[i].PostClosingEnd

		sessionInfoResponse.Result = append(sessionInfoResponse.Result, sessionInfo)
	}

	loggerconfig.Info("sessionInfoRes tl resp=", helpers.LogStructAsJSON(sessionInfoResponse), "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = sessionInfoResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}
