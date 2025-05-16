package tradelab

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
	"strconv"
	"time"
)

type SipServiceObj struct {
	tradeLabURL string
}

func InitSipSerivceProvider() SipServiceObj {
	defer models.HandlePanic()

	sipObj := SipServiceObj{
		tradeLabURL: constants.TLURL,
	}

	return sipObj
}

func (obj SipServiceObj) GetStockSips(clientId string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + SIPURL + "/" + clientId

	//make payload
	payload := new(bytes.Buffer)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "GetStockSips", duration, clientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " GetStockSips call api error =", err, " uccId:", clientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}

	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("GetStockSips TLres error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", clientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlGetSipResponse := TLGetSipResponse{}
	json.Unmarshal([]byte(string(body)), &tlGetSipResponse)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " GetStockSips tl status not ok =", tlGetSipResponse.Message, " uccId:", clientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlGetSipResponse.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("GetStockSips TLres =", tlGetSipResponse, " uccId:", clientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = tlGetSipResponse.Data
	apiRes.Message = tlGetSipResponse.Message
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj SipServiceObj) PlaceSipOrder(request models.PlaceSipRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + SIPURL

	requestJSON, err := json.Marshal(request)
	if err != nil {
		loggerconfig.Error("PlaceSipOrder marshal request error =", err, " clientId:", request.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	// Make payload
	payload := bytes.NewBuffer(requestJSON)

	// Call API
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "PlaceSipOrder", duration, request.ClientID, reqH.RequestId)

	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " PlaceSipOrder call api error =", err, " uccId:", request.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}

	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("PlaceSipOrder res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", request.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	placeSipResponse := TLUpdateSipOrderResponse{}
	json.Unmarshal([]byte(string(body)), &placeSipResponse)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " PlaceSipOrder tl status not ok =", placeSipResponse.Message, " uccId:", request.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = placeSipResponse.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("PlaceSipOrder TLres =", placeSipResponse, " uccId:", request.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
	apiRes.Data = placeSipResponse.Data
	apiRes.Message = placeSipResponse.Message
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj SipServiceObj) DeleteSipOrder(clientId, sipId string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + SIPURL + "/" + clientId + "/" + sipId

	// Make empty payload
	payload := new(bytes.Buffer)

	// Call API
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodDelete, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "DeleteSipOrder", duration, clientId, reqH.RequestId)

	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " DeleteSipOrder call api error =", err, " uccId:", clientId, " sipId:", sipId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}

	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("DeleteSipOrder res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", clientId, " sipId:", sipId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	sipResponse := TLUpdateSipOrderResponse{}
	json.Unmarshal([]byte(string(body)), &sipResponse)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " DeleteSipOrder tl status not ok =", sipResponse.Message, " uccId:", clientId, " sipId:", sipId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = sipResponse.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("DeleteSipOrder TLres =", sipResponse, " uccId:", clientId, " sipId:", sipId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
	apiRes.Data = sipResponse.Data
	apiRes.Message = sipResponse.Message
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj SipServiceObj) ModifySipOrder(request models.ModifySipRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + SIPURL

	requestJSON, err := json.Marshal(request)
	if err != nil {
		loggerconfig.Error("ModifySipOrder marshal request error =", err, " clientId:", request.ClientID, " sipId:", request.ID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	// Make payload
	payload := bytes.NewBuffer(requestJSON)

	// Call API
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPut, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ModifySipOrder", duration, request.ClientID, reqH.RequestId)

	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " ModifySipOrder call api error =", err, " uccId:", request.ClientID, " sipId:", request.ID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}

	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("ModifySipOrder res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", request.ClientID, " sipId:", request.ID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	modifySipResponse := TLUpdateSipOrderResponse{}
	json.Unmarshal([]byte(string(body)), &modifySipResponse)
	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " ModifySipOrder tl status not ok =", modifySipResponse.Message, " uccId:", request.ClientID, " sipId:", request.ID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = modifySipResponse.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("ModifySipOrder TLres =", modifySipResponse, " uccId:", request.ClientID, " sipId:", request.ID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
	apiRes.Data = modifySipResponse.Data
	apiRes.Message = modifySipResponse.Message
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj SipServiceObj) UpdateSipStatus(request models.UpdateSipStatusRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + SIPURL + SIPREVOCATION

	requestJSON, err := json.Marshal(request)
	if err != nil {
		loggerconfig.Error("UpdateSipStatus marshal request error =", err, " sipId:", request.ID, " action:", request.Action, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	// Make payload
	payload := bytes.NewBuffer(requestJSON)

	// Call API
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPut, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "UpdateSipStatus", duration, "unknown", reqH.RequestId) // Client ID not directly available

	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " UpdateSipStatus call api error =", err, " sipId:", request.ID, " action:", request.Action, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}

	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("UpdateSipStatus res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " sipId:", request.ID, " action:", request.Action, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	sipResponse := TLUpdateSipStatusResponse{}
	json.Unmarshal([]byte(string(body)), &sipResponse)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " UpdateSipStatus tl status not ok =", sipResponse.Error.Message, " sipId:", request.ID, " action:", request.Action, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = sipResponse.Error.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("UpdateSipStatus TLres =", sipResponse, " sipId:", request.ID, " action:", request.Action, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
	apiRes.Data = sipResponse.Data
	apiRes.Message = sipResponse.Data.Message
	apiRes.Status = true
	return http.StatusOK, apiRes
}
