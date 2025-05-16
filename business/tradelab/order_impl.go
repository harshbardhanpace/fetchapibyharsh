package tradelab

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/dbops"
	"space/helpers"
	"space/helpers/cache"
	"space/loggerconfig"
	"space/models"
)

type OrderObj struct {
	tradeLabURL string
	redisCli    cache.RedisCache
}

var objOrder OrderObj

func InitOrder(redisCli cache.RedisCache) OrderObj {
	defer models.HandlePanic()

	orderObj := OrderObj{
		tradeLabURL: constants.TLURL,
		redisCli:    redisCli,
	}

	objOrder = orderObj

	return orderObj
}

func (obj OrderObj) PlaceOrder(req models.PlaceOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + PLACEORDERURL

	//fill up the TL Req
	var tlPlaceOrderReq TradelabPlaceOrderRequest
	tlPlaceOrderReq.Amo = false
	tlPlaceOrderReq.ClientID = req.ClientID
	tlPlaceOrderReq.Device = reqH.DeviceType
	tlPlaceOrderReq.DisclosedQuantity = req.DisclosedQuantity
	tlPlaceOrderReq.Exchange = req.Exchange
	tlPlaceOrderReq.ExecutionType = "REGULAR"
	tlPlaceOrderReq.InstrumentToken = req.InstrumentToken
	tlPlaceOrderReq.OrderSide = req.OrderSide
	tlPlaceOrderReq.OrderType = req.OrderType
	tlPlaceOrderReq.Price = float64(req.Price)
	tlPlaceOrderReq.Product = req.Product
	tlPlaceOrderReq.Quantity = req.Quantity
	tlPlaceOrderReq.TriggerPrice = float64(req.TriggerPrice)
	tlPlaceOrderReq.UserOrderID = int(dbops.RedisRepo.Increment(USERORDERIDKEY))
	tlPlaceOrderReq.Validity = req.Validity

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlPlaceOrderReq)

	if req.Product == "MTF" {
		var tlMTFPlaceOrderReq TradelabPlaceOrderMTFRequest
		tlMTFPlaceOrderReq.ClientID = req.ClientID
		tlMTFPlaceOrderReq.Device = reqH.DeviceType
		tlMTFPlaceOrderReq.DisclosedQuantity = req.DisclosedQuantity
		tlMTFPlaceOrderReq.Exchange = req.Exchange
		tlMTFPlaceOrderReq.NoOfLegs = req.NoOfLegs
		tlMTFPlaceOrderReq.InstrumentToken = req.InstrumentToken
		tlMTFPlaceOrderReq.OrderSide = req.OrderSide
		tlMTFPlaceOrderReq.OrderType = req.OrderType
		tlMTFPlaceOrderReq.Price = strconv.FormatFloat(req.Price, 'f', -1, 64)
		tlMTFPlaceOrderReq.Product = req.Product
		tlMTFPlaceOrderReq.Quantity = req.Quantity
		tlMTFPlaceOrderReq.TriggerPrice = float64(req.TriggerPrice)
		tlMTFPlaceOrderReq.UserOrderId = int(dbops.RedisRepo.Increment(USERORDERIDKEY))
		tlMTFPlaceOrderReq.Validity = req.Validity
		payload = new(bytes.Buffer)
		json.NewEncoder(payload).Encode(tlMTFPlaceOrderReq)
	}

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "PlaceOrder", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " placeOrderRes call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("placeOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlPlaceOrderRes := TradelabPlaceOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlPlaceOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " placeOrderRes tl status not ok =", tlPlaceOrderRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlPlaceOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var placeOrderRes models.PlaceOrderResponse
	placeOrderRes.OmsOrderID = tlPlaceOrderRes.Data.OmsOrderID
	placeOrderRes.UserOrderID = tlPlaceOrderRes.Data.UserOrderID

	loggerconfig.Info("placeOrderRes tl resp=", helpers.LogStructAsJSON(placeOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = placeOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj OrderObj) PlaceAMOOrder(req models.PlaceOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + PLACEORDERURL

	//fill up the TL Req
	var tlAMOPlaceOrderReq TradelabPlaceAMORequest
	tlAMOPlaceOrderReq.ClientID = req.ClientID
	tlAMOPlaceOrderReq.DisclosedQuantity = req.DisclosedQuantity
	tlAMOPlaceOrderReq.Exchange = req.Exchange
	tlAMOPlaceOrderReq.Device = reqH.DeviceType
	tlAMOPlaceOrderReq.ExecutionType = "AMO"
	tlAMOPlaceOrderReq.InstrumentToken = req.InstrumentToken
	tlAMOPlaceOrderReq.OrderSide = req.OrderSide
	tlAMOPlaceOrderReq.OrderType = req.OrderType
	tlAMOPlaceOrderReq.Price = float64(req.Price)
	tlAMOPlaceOrderReq.Product = req.Product
	tlAMOPlaceOrderReq.Quantity = req.Quantity
	tlAMOPlaceOrderReq.TriggerPrice = float64(req.TriggerPrice)
	tlAMOPlaceOrderReq.UserOrderID = int(dbops.RedisRepo.Increment(USERORDERIDKEY))
	tlAMOPlaceOrderReq.Validity = req.Validity

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlAMOPlaceOrderReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "PlaceAMOOrder", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " placeAMOOrderRes call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("placeAMOOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlAMOPlaceOrderRes := TradelabAMOResponse{}
	json.Unmarshal([]byte(string(body)), &tlAMOPlaceOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " placeAMOOrderRes tl status not ok =", tlAMOPlaceOrderRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlAMOPlaceOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var placeAMOOrderRes models.AMOOrderResponse
	placeAMOOrderRes.OmsOrderID = tlAMOPlaceOrderRes.Data.OmsOrderID

	loggerconfig.Info("placeAMOOrderRes tl resp=", helpers.LogStructAsJSON(placeAMOOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = placeAMOOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj OrderObj) ModifyOrder(req models.ModifyOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + PLACEORDERURL

	//fill up the TL Req
	var tlModifyOrderReq TradelabModifyOrderRequest
	tlModifyOrderReq.ClientID = req.ClientID
	tlModifyOrderReq.DisclosedQuantity = req.DisclosedQuantity
	tlModifyOrderReq.Exchange = req.Exchange
	tlModifyOrderReq.ExecutionType = "REGULAR"
	tlModifyOrderReq.OrderType = req.OrderType
	tlModifyOrderReq.Price = float64(req.Price)
	tlModifyOrderReq.Product = req.Product
	tlModifyOrderReq.Quantity = req.Quantity
	tlModifyOrderReq.TriggerPrice = float64(req.TriggerPrice)
	tlModifyOrderReq.OmsOrderID = req.OmsOrderID
	tlModifyOrderReq.Validity = req.Validity
	tlModifyOrderReq.ExchangeOrderID = req.ExchangeOrderID
	tlModifyOrderReq.FilledQuantity = req.FilledQuantity
	tlModifyOrderReq.RemainingQuantity = req.RemainingQuantity
	tlModifyOrderReq.LastActivityReference = req.LastActivityReference

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlModifyOrderReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPut, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ModifyOrder", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " modifyOrderRes call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("modifyOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlModifyOrderRes := TradelabPlaceOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlModifyOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " modifyOrderRes tl status not ok =", tlModifyOrderRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlModifyOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var modifyOrderRes models.ModifyOrCancelOrderResponse
	modifyOrderRes.OmsOrderID = tlModifyOrderRes.Data.OmsOrderID

	loggerconfig.Info("modifyOrderRes tl resp=", helpers.LogStructAsJSON(modifyOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = modifyOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj OrderObj) ModifyAMOOrder(req models.ModifyAMORequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + PLACEORDERURL

	//fill up the TL Req
	var tlModifyAMOOrderReq TradelabModifyAMORequest
	tlModifyAMOOrderReq.ClientID = req.ClientID
	tlModifyAMOOrderReq.DisclosedQuantity = req.DisclosedQuantity
	tlModifyAMOOrderReq.Exchange = req.Exchange
	tlModifyAMOOrderReq.ExecutionType = "AMO"
	tlModifyAMOOrderReq.OrderType = req.OrderType
	tlModifyAMOOrderReq.Price = float64(req.Price)
	tlModifyAMOOrderReq.Product = req.Product
	tlModifyAMOOrderReq.Quantity = req.Quantity
	tlModifyAMOOrderReq.TriggerPrice = float64(req.TriggerPrice)
	tlModifyAMOOrderReq.ExchangeOrderID = req.ExchangeOrderID
	tlModifyAMOOrderReq.FilledQuantity = req.FilledQuantity
	tlModifyAMOOrderReq.RemainingQuantity = req.RemainingQuantity
	tlModifyAMOOrderReq.LastActivityReference = req.LastActivityReference
	tlModifyAMOOrderReq.OmsOrderID = req.OmsOrderID
	tlModifyAMOOrderReq.Validity = req.Validity

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlModifyAMOOrderReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPut, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ModifyAMOOrder", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " modifyOrderRes call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("modifyOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlModifyAMOOrderRes := TradelabAMOResponse{}
	json.Unmarshal([]byte(string(body)), &tlModifyAMOOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " modifyAMOOrderRes tl status not ok =", tlModifyAMOOrderRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlModifyAMOOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var modifyAMOOrderRes models.AMOOrderResponse
	modifyAMOOrderRes.OmsOrderID = tlModifyAMOOrderRes.Data.OmsOrderID

	loggerconfig.Info("modifyOrderRes tl resp=", helpers.LogStructAsJSON(modifyAMOOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = modifyAMOOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj OrderObj) CancelOrder(req models.CancelOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + PLACEORDERURL + "/" + req.OmsOrderId + "?client_id=" + url.QueryEscape(req.ClientID) + "&execution_type=" + req.ExecutionType

	//make payload
	payload := new(bytes.Buffer)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodDelete, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "CancelOrder", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " cancelOrderRes call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("cancelOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlCancelOrderRes := TradelabPlaceOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlCancelOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " cancelOrderRes tl status not ok =", tlCancelOrderRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlCancelOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var cancelOrderRes models.ModifyOrCancelOrderResponse
	cancelOrderRes.OmsOrderID = tlCancelOrderRes.Data.OmsOrderID

	loggerconfig.Info("cancelOrderRes tl resp=", helpers.LogStructAsJSON(cancelOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = cancelOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj OrderObj) CancelAMOOrder(req models.CancelOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + PLACEORDERURL + "/" + req.OmsOrderId + "?client_id=" + url.QueryEscape(req.ClientID) + "&execution_type=" + req.ExecutionType

	//fill up the TL Req
	var tlCancelAMOOrderReq TradelabDeleteAMORequest
	tlCancelAMOOrderReq.ClientID = req.ClientID
	tlCancelAMOOrderReq.ExecutionType = "AMO"

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlCancelAMOOrderReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodDelete, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "CancelAMOOrder", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " cancelAMOOrderRes call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("cancelAMOOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlCancelAMOOrderRes := TradelabCancelOrModifyResponse{}
	json.Unmarshal([]byte(string(body)), &tlCancelAMOOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " cancelAMOOrderRes tl status not ok =", tlCancelAMOOrderRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlCancelAMOOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var cancelAMOOrderRes models.ModifyOrCancelOrderResponse
	cancelAMOOrderRes.OmsOrderID = tlCancelAMOOrderRes.Data.OmsOrderID

	loggerconfig.Info("cancelAMOOrderRes tl resp=", helpers.LogStructAsJSON(tlCancelAMOOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = cancelAMOOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj OrderObj) PendingOrder(req models.PendingOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	if constants.LocalCachingCallEnabled {
		return fetchPendingOrderFromCache(req, reqH)
	}

	url := obj.tradeLabURL + PENDINGORDERURL + "?type=" + req.Type + "&client_id=" + url.QueryEscape(req.ClientID)
	//make payload
	payload := new(bytes.Buffer) // empty payload

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "PendingOrder", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " pendingOrderRes call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("pendingOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlPendingOrderRes := TradelabPendingOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlPendingOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " pendingOrderRes tl status not ok =", tlPendingOrderRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlPendingOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var pendingOrderRes models.PendingOrderResponse
	responseOrders := make([]models.PendingOrderResponseOrders, 0)

	for i := 0; i < len(tlPendingOrderRes.Data.Orders); i++ {
		var pendingOrderResOrders models.PendingOrderResponseOrders
		pendingOrderResOrders.TradingSymbol = tlPendingOrderRes.Data.Orders[i].TradingSymbol
		pendingOrderResOrders.AverageTradePrice = tlPendingOrderRes.Data.Orders[i].AverageTradePrice
		pendingOrderResOrders.Exchange = tlPendingOrderRes.Data.Orders[i].Exchange
		pendingOrderResOrders.ProCli = tlPendingOrderRes.Data.Orders[i].ProCli
		pendingOrderResOrders.MarketProtectionPercentage = tlPendingOrderRes.Data.Orders[i].MarketProtectionPercentage
		pendingOrderResOrders.OrderEntryTime = tlPendingOrderRes.Data.Orders[i].OrderEntryTime
		pendingOrderResOrders.Mode = tlPendingOrderRes.Data.Orders[i].Mode
		pendingOrderResOrders.OmsOrderID = tlPendingOrderRes.Data.Orders[i].OmsOrderID
		pendingOrderResOrders.TrailingStopLoss = tlPendingOrderRes.Data.Orders[i].TrailingStopLoss
		pendingOrderResOrders.Deposit = tlPendingOrderRes.Data.Orders[i].Deposit
		pendingOrderResOrders.SquareOffValue = tlPendingOrderRes.Data.Orders[i].SquareOffValue
		pendingOrderResOrders.DisclosedQuantity = tlPendingOrderRes.Data.Orders[i].DisclosedQuantity
		pendingOrderResOrders.StopLossValue = tlPendingOrderRes.Data.Orders[i].StopLossValue
		pendingOrderResOrders.Price = tlPendingOrderRes.Data.Orders[i].Price
		pendingOrderResOrders.OrderTag = tlPendingOrderRes.Data.Orders[i].OrderTag
		pendingOrderResOrders.Device = tlPendingOrderRes.Data.Orders[i].Device
		pendingOrderResOrders.RemainingQuantity = tlPendingOrderRes.Data.Orders[i].RemainingQuantity
		pendingOrderResOrders.LastActivityReference = tlPendingOrderRes.Data.Orders[i].LastActivityReference
		pendingOrderResOrders.AveragePrice = tlPendingOrderRes.Data.Orders[i].AveragePrice
		pendingOrderResOrders.SquareOff = tlPendingOrderRes.Data.Orders[i].SquareOff
		pendingOrderResOrders.OrderStatusInfo = tlPendingOrderRes.Data.Orders[i].OrderStatusInfo
		pendingOrderResOrders.Quantity = tlPendingOrderRes.Data.Orders[i].Quantity
		pendingOrderResOrders.ExecutionType = tlPendingOrderRes.Data.Orders[i].ExecutionType
		pendingOrderResOrders.ClientID = tlPendingOrderRes.Data.Orders[i].ClientID
		pendingOrderResOrders.ExchangeTime = tlPendingOrderRes.Data.Orders[i].ExchangeTime
		pendingOrderResOrders.OrderSide = tlPendingOrderRes.Data.Orders[i].OrderSide
		pendingOrderResOrders.LoginID = tlPendingOrderRes.Data.Orders[i].LoginID
		pendingOrderResOrders.Validity = tlPendingOrderRes.Data.Orders[i].Validity
		pendingOrderResOrders.InstrumentToken = tlPendingOrderRes.Data.Orders[i].InstrumentToken
		pendingOrderResOrders.Product = tlPendingOrderRes.Data.Orders[i].Product
		pendingOrderResOrders.TriggerPrice = tlPendingOrderRes.Data.Orders[i].TriggerPrice
		pendingOrderResOrders.Segment = tlPendingOrderRes.Data.Orders[i].Segment
		pendingOrderResOrders.TradePrice = tlPendingOrderRes.Data.Orders[i].TradePrice
		pendingOrderResOrders.OrderType = tlPendingOrderRes.Data.Orders[i].OrderType
		// pendingOrderResOrders.ContractDescription = tlPendingOrderRes.Data.Orders[i].ContractDescription
		pendingOrderResOrders.RejectionCode = tlPendingOrderRes.Data.Orders[i].RejectionCode
		pendingOrderResOrders.LegOrderIndicator = tlPendingOrderRes.Data.Orders[i].LegOrderIndicator
		pendingOrderResOrders.ExchangeOrderID = tlPendingOrderRes.Data.Orders[i].ExchangeOrderID
		pendingOrderResOrders.OrderStatus = tlPendingOrderRes.Data.Orders[i].OrderStatus
		pendingOrderResOrders.FilledQuantity = tlPendingOrderRes.Data.Orders[i].FilledQuantity
		pendingOrderResOrders.TargetPriceType = tlPendingOrderRes.Data.Orders[i].TargetPriceType
		pendingOrderResOrders.IsTrailing = tlPendingOrderRes.Data.Orders[i].IsTrailing
		pendingOrderResOrders.UserOrderID = tlPendingOrderRes.Data.Orders[i].UserOrderID
		pendingOrderResOrders.LotSize = tlPendingOrderRes.Data.Orders[i].LotSize
		pendingOrderResOrders.Series = tlPendingOrderRes.Data.Orders[i].Series
		pendingOrderResOrders.NnfID = tlPendingOrderRes.Data.Orders[i].NnfID
		pendingOrderResOrders.RejectionReason = tlPendingOrderRes.Data.Orders[i].RejectionReason
		if pendingOrderResOrders.Exchange == "NFO" || pendingOrderResOrders.Exchange == "BFO" || pendingOrderResOrders.Exchange == "MCX" {
			var req models.ScripInfoRequest
			req.Exchange = pendingOrderResOrders.Exchange
			req.Token = strconv.Itoa(pendingOrderResOrders.InstrumentToken)
			_, pendingOrderResOrders.AdditionalInfo = GetAdditionalInfo(obj.tradeLabURL, req, reqH)
		}
		responseOrders = append(responseOrders, pendingOrderResOrders)
	}
	pendingOrderRes.Orders = responseOrders

	loggerconfig.Info("pendingOrderRes tl resp=", helpers.LogStructAsJSON(pendingOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = pendingOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}

func fetchPendingOrderFromCache(req models.PendingOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var pendingOrderRes models.PendingOrderResponse
	responseOrders := make([]models.PendingOrderResponseOrders, 0)
	orderIdsPending, _ := getOrderIds(req.ClientID, constants.Pending, reqH)
	for orderId, _ := range orderIdsPending {
		tlOrderPacket, _ := getOrderStoredRedis(req.ClientID, orderId, reqH) // // ignoring error as it is printed in the func and we don't want to stop the execution if failed to read from redis.

		var pendingOrderResOrders models.PendingOrderResponseOrders
		pendingOrderResOrders.TradingSymbol = tlOrderPacket.TradingSymbol
		pendingOrderResOrders.AverageTradePrice = tlOrderPacket.AverageTradePrice
		pendingOrderResOrders.Exchange = tlOrderPacket.Exchange
		pendingOrderResOrders.ProCli = tlOrderPacket.ProCli
		pendingOrderResOrders.MarketProtectionPercentage = int(tlOrderPacket.MarketProtectionPercentage)
		pendingOrderResOrders.OrderEntryTime = tlOrderPacket.OrderEntryTime
		pendingOrderResOrders.Mode = tlOrderPacket.Mode
		pendingOrderResOrders.OmsOrderID = tlOrderPacket.OmsOrderID
		pendingOrderResOrders.TrailingStopLoss = tlOrderPacket.TrailingStopLoss
		pendingOrderResOrders.Deposit = int(tlOrderPacket.Deposit)
		pendingOrderResOrders.SquareOffValue = tlOrderPacket.SquareOffValue
		pendingOrderResOrders.DisclosedQuantity = tlOrderPacket.DisclosedQuantity
		pendingOrderResOrders.StopLossValue = tlOrderPacket.StopLossValue
		pendingOrderResOrders.Price = tlOrderPacket.Price
		// pendingOrderResOrders.OrderTag = tlOrderPacket.OrderTag  // tradelab give its value empty
		pendingOrderResOrders.Device = tlOrderPacket.Device
		pendingOrderResOrders.RemainingQuantity = tlOrderPacket.RemainingQuantity
		pendingOrderResOrders.LastActivityReference = tlOrderPacket.LastActivityReference
		pendingOrderResOrders.AveragePrice = tlOrderPacket.AveragePrice
		// pendingOrderResOrders.SquareOff = tlOrderPacket.SquareOff  // tradelab give its value empty
		pendingOrderResOrders.OrderStatusInfo = tlOrderPacket.OrderStatus
		pendingOrderResOrders.Quantity = tlOrderPacket.Quantity
		pendingOrderResOrders.ExecutionType = tlOrderPacket.ExecutionType
		pendingOrderResOrders.ClientID = tlOrderPacket.ClientID
		pendingOrderResOrders.ExchangeTime = tlOrderPacket.ExchangeTime
		pendingOrderResOrders.OrderSide = tlOrderPacket.OrderSide
		pendingOrderResOrders.LoginID = tlOrderPacket.LoginID
		pendingOrderResOrders.Validity = tlOrderPacket.Validity
		pendingOrderResOrders.InstrumentToken = tlOrderPacket.InstrumentToken
		pendingOrderResOrders.Product = tlOrderPacket.Product
		pendingOrderResOrders.TriggerPrice = tlOrderPacket.TriggerPrice
		// pendingOrderResOrders.Segment = tlOrderPacket.Segment  // tradelab give its value empty
		pendingOrderResOrders.TradePrice = tlOrderPacket.TradePrice
		pendingOrderResOrders.OrderType = tlOrderPacket.OrderType
		//pendingOrderResOrders.ContractDescription = tlOrderPacket.ContractDescription
		rejectionCode, _ := strconv.Atoi(tlOrderPacket.RejectionCode)
		pendingOrderResOrders.RejectionCode = rejectionCode
		pendingOrderResOrders.LegOrderIndicator = tlOrderPacket.LegOrderIndicator
		pendingOrderResOrders.ExchangeOrderID = tlOrderPacket.ExchangeOrderID
		pendingOrderResOrders.OrderStatus = tlOrderPacket.OrderStatus
		pendingOrderResOrders.FilledQuantity = tlOrderPacket.FilledQuantity
		// pendingOrderResOrders.TargetPriceType = tlOrderPacket.TargetPriceType  // for every order tradelab gives its value as ABSOLUTE
		pendingOrderResOrders.TargetPriceType = constants.AbsoluteTargetPriceType
		pendingOrderResOrders.IsTrailing = tlOrderPacket.IsTrailing
		pendingOrderResOrders.UserOrderID = tlOrderPacket.UserOrderID
		pendingOrderResOrders.LotSize = int(tlOrderPacket.LotSize)
		pendingOrderResOrders.Series = tlOrderPacket.Series
		pendingOrderResOrders.NnfID = tlOrderPacket.Nnfid
		pendingOrderResOrders.RejectionReason = tlOrderPacket.RejectionReason
		responseOrders = append(responseOrders, pendingOrderResOrders)
	}
	pendingOrderRes.Orders = responseOrders

	loggerconfig.Info("pendingOrderRes from redis=", helpers.LogStructAsJSON(pendingOrderRes), " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	var apiRes apihelpers.APIRes
	apiRes.Data = pendingOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}

func getOrderStoredRedis(clientId string, orderId string, reqH models.ReqHeader) (TLOrderUpdatePacket, error) {
	var tlOrderPacket TLOrderUpdatePacket
	val, err := dbops.OrderRedisRepo.HGet(clientId, constants.Order+orderId)
	if err != nil {
		loggerconfig.Error("readMessages, Failed to read from redis:", err, " requestId:", reqH.RequestId)
		return tlOrderPacket, errors.New(constants.RedisReadFailed)
	}
	err = json.Unmarshal([]byte(val), &tlOrderPacket)
	if err != nil {
		loggerconfig.Error("readMessages, Packet unable to unmarshal:", val, " err: ", err, " requestId:", reqH.RequestId)
		return tlOrderPacket, errors.New(constants.UnmarshalFailed)
	}
	loggerconfig.Info("getOrderStoredRedis tlOrderPacket from redis=", helpers.LogStructAsJSON(tlOrderPacket), " uccId:", clientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
	return tlOrderPacket, nil
}

func getTradeStoredRedis(clientId string, orderId string, reqH models.ReqHeader) (TlTradeUpdate, error) {
	var tlTradeUpdate TlTradeUpdate
	val, err := dbops.OrderRedisRepo.HGet(clientId, constants.Trade+orderId)
	if err != nil {
		loggerconfig.Error("getTradeStoredRedis readMessages, Failed to read from redis:", err, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return tlTradeUpdate, errors.New(constants.RedisReadFailed)
	}
	err = json.Unmarshal([]byte(val), &tlTradeUpdate)
	if err != nil {
		loggerconfig.Error("getTradeStoredRedis readMessages, Packet unable to unmarshal:", val, " err: ", err, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return tlTradeUpdate, errors.New(constants.UnmarshalFailed)
	}
	loggerconfig.Info("getTradeStoredRedis tlTradeUpdate from redis=", helpers.LogStructAsJSON(tlTradeUpdate), " uccId:", clientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
	return tlTradeUpdate, nil
}

func getOrderIds(clientId string, orderStatus string, reqH models.ReqHeader) (map[string]bool, error) {
	orderIds := make(map[string]bool)
	storedOrderMap, err := dbops.OrderRedisRepo.HGet(clientId, orderStatus)
	if err != nil {
		loggerconfig.Error("readMessages Failed to read from redis:", err, "for order status: ", orderStatus, " clientId: ", clientId, " requestId: ", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return orderIds, errors.New(constants.RedisReadFailed)
	}
	err = json.Unmarshal([]byte(storedOrderMap), &orderIds)
	if err != nil {
		loggerconfig.Error("readMessages, Packet unable to unmarshal orderIds", err, "for order status: ", orderStatus, " clientId: ", clientId, " requestId: ", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return orderIds, errors.New(constants.UnmarshalFailed)
	}
	loggerconfig.Info("GetOrderIds orderIds from redis=", orderIds, " uccId:", clientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
	return orderIds, nil
}

func (obj OrderObj) CompletedOrder(req models.CompletedOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return CompletedOrderInternal(req, reqH)
}

var CompletedOrderInternal = func(req models.CompletedOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	if constants.LocalCachingCallEnabled {
		return fetchCompletedOrderFromCache(req, reqH)
	}

	url := objOrder.tradeLabURL + COMPLETEDORDERURL + "?type=" + req.Type + "&client_id=" + url.QueryEscape(req.ClientID)

	//make payload
	payload := new(bytes.Buffer) // empty payload

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "CompletedOrder", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " completedOrderRes call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("completedOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlCompletedOrderRes := TradelabCompletedOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlCompletedOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " completedOrderRes tl status not ok =", tlCompletedOrderRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlCompletedOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var completedOrderRes models.CompletedOrderResponse

	responseOrders := make([]models.CompletedOrderResponseOrders, 0)

	for i := 0; i < len(tlCompletedOrderRes.Data.Orders); i++ {
		var completedOrderResOrders models.CompletedOrderResponseOrders
		completedOrderResOrders.TradingSymbol = tlCompletedOrderRes.Data.Orders[i].TradingSymbol
		completedOrderResOrders.AverageTradePrice = tlCompletedOrderRes.Data.Orders[i].AverageTradePrice
		completedOrderResOrders.Exchange = tlCompletedOrderRes.Data.Orders[i].Exchange
		completedOrderResOrders.ProCli = tlCompletedOrderRes.Data.Orders[i].ProCli
		completedOrderResOrders.MarketProtectionPercentage = tlCompletedOrderRes.Data.Orders[i].MarketProtectionPercentage
		completedOrderResOrders.OrderEntryTime = tlCompletedOrderRes.Data.Orders[i].OrderEntryTime
		completedOrderResOrders.Mode = tlCompletedOrderRes.Data.Orders[i].Mode
		completedOrderResOrders.OmsOrderID = tlCompletedOrderRes.Data.Orders[i].OmsOrderID
		completedOrderResOrders.TrailingStopLoss = tlCompletedOrderRes.Data.Orders[i].TrailingStopLoss
		completedOrderResOrders.Deposit = tlCompletedOrderRes.Data.Orders[i].Deposit
		completedOrderResOrders.SquareOffValue = tlCompletedOrderRes.Data.Orders[i].SquareOffValue
		completedOrderResOrders.DisclosedQuantity = tlCompletedOrderRes.Data.Orders[i].DisclosedQuantity
		completedOrderResOrders.StopLossValue = tlCompletedOrderRes.Data.Orders[i].StopLossValue
		completedOrderResOrders.Price = tlCompletedOrderRes.Data.Orders[i].Price
		completedOrderResOrders.OrderTag = tlCompletedOrderRes.Data.Orders[i].OrderTag
		completedOrderResOrders.Device = tlCompletedOrderRes.Data.Orders[i].Device
		completedOrderResOrders.RemainingQuantity = tlCompletedOrderRes.Data.Orders[i].RemainingQuantity
		completedOrderResOrders.LastActivityReference = tlCompletedOrderRes.Data.Orders[i].LastActivityReference
		completedOrderResOrders.AveragePrice = tlCompletedOrderRes.Data.Orders[i].AveragePrice
		completedOrderResOrders.SquareOff = tlCompletedOrderRes.Data.Orders[i].SquareOff
		completedOrderResOrders.OrderStatusInfo = tlCompletedOrderRes.Data.Orders[i].OrderStatusInfo
		completedOrderResOrders.Quantity = tlCompletedOrderRes.Data.Orders[i].Quantity
		completedOrderResOrders.ExecutionType = tlCompletedOrderRes.Data.Orders[i].ExecutionType
		completedOrderResOrders.ClientID = tlCompletedOrderRes.Data.Orders[i].ClientID
		completedOrderResOrders.ExchangeTime = tlCompletedOrderRes.Data.Orders[i].ExchangeTime
		completedOrderResOrders.OrderSide = tlCompletedOrderRes.Data.Orders[i].OrderSide
		completedOrderResOrders.LoginID = tlCompletedOrderRes.Data.Orders[i].LoginID
		completedOrderResOrders.Validity = tlCompletedOrderRes.Data.Orders[i].Validity
		completedOrderResOrders.InstrumentToken = tlCompletedOrderRes.Data.Orders[i].InstrumentToken
		completedOrderResOrders.Product = tlCompletedOrderRes.Data.Orders[i].Product
		completedOrderResOrders.TriggerPrice = tlCompletedOrderRes.Data.Orders[i].TriggerPrice
		completedOrderResOrders.Segment = tlCompletedOrderRes.Data.Orders[i].Segment
		completedOrderResOrders.TradePrice = tlCompletedOrderRes.Data.Orders[i].TradePrice
		completedOrderResOrders.OrderType = tlCompletedOrderRes.Data.Orders[i].OrderType
		//completedOrderResOrders.ContractDescription = tlCompletedOrderRes.Data.Orders[i].ContractDescription
		completedOrderResOrders.RejectionCode = tlCompletedOrderRes.Data.Orders[i].RejectionCode
		completedOrderResOrders.LegOrderIndicator = tlCompletedOrderRes.Data.Orders[i].LegOrderIndicator
		completedOrderResOrders.ExchangeOrderID = tlCompletedOrderRes.Data.Orders[i].ExchangeOrderID
		completedOrderResOrders.OrderStatus = tlCompletedOrderRes.Data.Orders[i].OrderStatus
		completedOrderResOrders.FilledQuantity = tlCompletedOrderRes.Data.Orders[i].FilledQuantity
		completedOrderResOrders.TargetPriceType = tlCompletedOrderRes.Data.Orders[i].TargetPriceType
		completedOrderResOrders.IsTrailing = tlCompletedOrderRes.Data.Orders[i].IsTrailing
		completedOrderResOrders.UserOrderID = tlCompletedOrderRes.Data.Orders[i].UserOrderID
		completedOrderResOrders.LotSize = tlCompletedOrderRes.Data.Orders[i].LotSize
		completedOrderResOrders.Series = tlCompletedOrderRes.Data.Orders[i].Series
		completedOrderResOrders.NnfID = tlCompletedOrderRes.Data.Orders[i].NnfID
		completedOrderResOrders.RejectionReason = tlCompletedOrderRes.Data.Orders[i].RejectionReason
		if completedOrderResOrders.Exchange == "NFO" || completedOrderResOrders.Exchange == "BFO" || completedOrderResOrders.Exchange == "MCX" {
			var req models.ScripInfoRequest
			req.Exchange = completedOrderResOrders.Exchange
			req.Token = strconv.Itoa(completedOrderResOrders.InstrumentToken)
			_, completedOrderResOrders.AdditionalInfo = GetAdditionalInfo(objOrder.tradeLabURL, req, reqH)
		}
		responseOrders = append(responseOrders, completedOrderResOrders)
	}
	completedOrderRes.Orders = responseOrders

	loggerconfig.Info("completedOrderRes tl resp=", helpers.LogStructAsJSON(completedOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = completedOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}

func fetchCompletedOrderFromCache(req models.CompletedOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var completedOrderRes models.CompletedOrderResponse
	responseOrders := make([]models.CompletedOrderResponseOrders, 0)
	orderIdsCompleted, _ := getOrderIds(req.ClientID, constants.Completed, reqH)
	for orderId, _ := range orderIdsCompleted {
		tlOrderPacket, err := getOrderStoredRedis(req.ClientID, orderId, reqH)
		if err != nil {
			// complete error message is prented on func getOrderStoredRedis
			return apihelpers.SendInternalServerError()
		}

		var completedOrderResOrders models.CompletedOrderResponseOrders
		completedOrderResOrders.TradingSymbol = tlOrderPacket.TradingSymbol
		completedOrderResOrders.AverageTradePrice = tlOrderPacket.AverageTradePrice
		completedOrderResOrders.Exchange = tlOrderPacket.Exchange
		completedOrderResOrders.ProCli = tlOrderPacket.ProCli
		completedOrderResOrders.MarketProtectionPercentage = int(tlOrderPacket.MarketProtectionPercentage)
		completedOrderResOrders.OrderEntryTime = tlOrderPacket.OrderEntryTime
		completedOrderResOrders.Mode = tlOrderPacket.Mode
		completedOrderResOrders.OmsOrderID = tlOrderPacket.OmsOrderID
		completedOrderResOrders.TrailingStopLoss = tlOrderPacket.TrailingStopLoss
		completedOrderResOrders.Deposit = int(tlOrderPacket.Deposit)
		completedOrderResOrders.SquareOffValue = tlOrderPacket.SquareOffValue
		completedOrderResOrders.DisclosedQuantity = tlOrderPacket.DisclosedQuantity
		completedOrderResOrders.StopLossValue = tlOrderPacket.StopLossValue
		completedOrderResOrders.Price = tlOrderPacket.Price
		// completedOrderResOrders.OrderTag = tlOrderPacket.OrderTag  // tradelab give its value empty
		completedOrderResOrders.Device = tlOrderPacket.Device
		completedOrderResOrders.RemainingQuantity = tlOrderPacket.RemainingQuantity
		completedOrderResOrders.LastActivityReference = int(tlOrderPacket.LastActivityReference)
		completedOrderResOrders.AveragePrice = tlOrderPacket.AveragePrice
		// completedOrderResOrders.SquareOff = tlOrderPacket.SquareOff  // tradelab give its value empty
		completedOrderResOrders.OrderStatusInfo = tlOrderPacket.OrderStatus
		completedOrderResOrders.Quantity = tlOrderPacket.Quantity
		completedOrderResOrders.ExecutionType = tlOrderPacket.ExecutionType
		completedOrderResOrders.ClientID = tlOrderPacket.ClientID
		completedOrderResOrders.ExchangeTime = tlOrderPacket.ExchangeTime
		completedOrderResOrders.OrderSide = tlOrderPacket.OrderSide
		completedOrderResOrders.LoginID = tlOrderPacket.LoginID
		completedOrderResOrders.Validity = tlOrderPacket.Validity
		completedOrderResOrders.InstrumentToken = tlOrderPacket.InstrumentToken
		completedOrderResOrders.Product = tlOrderPacket.Product
		completedOrderResOrders.TriggerPrice = tlOrderPacket.TriggerPrice
		// completedOrderResOrders.Segment = tlOrderPacket.Segment  // tradelab give its value empty
		completedOrderResOrders.TradePrice = tlOrderPacket.TradePrice
		completedOrderResOrders.OrderType = tlOrderPacket.OrderType
		//completedOrderResOrders.ContractDescription = tlOrderPacket.ContractDescription
		rejectionCode, _ := strconv.Atoi(tlOrderPacket.RejectionCode)
		completedOrderResOrders.RejectionCode = rejectionCode
		completedOrderResOrders.LegOrderIndicator = tlOrderPacket.LegOrderIndicator
		completedOrderResOrders.ExchangeOrderID = tlOrderPacket.ExchangeOrderID
		completedOrderResOrders.OrderStatus = tlOrderPacket.OrderStatus
		completedOrderResOrders.FilledQuantity = tlOrderPacket.FilledQuantity
		// completedOrderResOrders.TargetPriceType = tlOrderPacket.TargetPriceType  // for every order tradelab gives its value as ABSOLUTE
		completedOrderResOrders.TargetPriceType = constants.AbsoluteTargetPriceType
		completedOrderResOrders.IsTrailing = tlOrderPacket.IsTrailing
		completedOrderResOrders.UserOrderID = tlOrderPacket.UserOrderID
		completedOrderResOrders.LotSize = int(tlOrderPacket.LotSize)
		completedOrderResOrders.Series = tlOrderPacket.Series
		completedOrderResOrders.NnfID = tlOrderPacket.Nnfid
		completedOrderResOrders.RejectionReason = tlOrderPacket.RejectionReason
		responseOrders = append(responseOrders, completedOrderResOrders)
	}
	completedOrderRes.Orders = responseOrders

	loggerconfig.Info("completedOrderRes from redis=", helpers.LogStructAsJSON(completedOrderRes), " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	var apiRes apihelpers.APIRes
	apiRes.Data = completedOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj OrderObj) TradeBook(req models.TradeBookRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	if constants.LocalCachingCallEnabled {
		return fetchTradeBookFromCache(req, reqH)
	}
	url := obj.tradeLabURL + TRADEURL + "?client_id=" + url.QueryEscape(req.ClientID)

	//make payload
	payload := new(bytes.Buffer)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "TradeBook", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " tradeBookReq call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("tradeBookReq res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlTradeBookResponse := TradeLabTradeBookResponse{}
	json.Unmarshal([]byte(string(body)), &tlTradeBookResponse)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " tradeBookRes tl status not ok =", tlTradeBookResponse.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlTradeBookResponse.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var tradeBookResponse models.TradeBookResponse
	tradeBookResponseData := make([]models.TradeBookResponseData, 0)
	for i := 0; i < len(tlTradeBookResponse.Data.Trades); i++ {
		var tradeBookData models.TradeBookResponseData
		tradeBookData.BookType = tlTradeBookResponse.Data.Trades[i].BookType
		tradeBookData.BrokerID = tlTradeBookResponse.Data.Trades[i].BrokerID
		tradeBookData.ClientID = tlTradeBookResponse.Data.Trades[i].ClientID
		tradeBookData.DisclosedVol = tlTradeBookResponse.Data.Trades[i].DisclosedVol
		tradeBookData.DisclosedVolRemaining = tlTradeBookResponse.Data.Trades[i].DisclosedVolRemaining
		tradeBookData.Exchange = tlTradeBookResponse.Data.Trades[i].Exchange
		tradeBookData.ExchangeOrderID = tlTradeBookResponse.Data.Trades[i].ExchangeOrderID
		tradeBookData.ExchangeTime = tlTradeBookResponse.Data.Trades[i].ExchangeTime
		tradeBookData.FillNumber = tlTradeBookResponse.Data.Trades[i].FillNumber
		tradeBookData.FilledQuantity = tlTradeBookResponse.Data.Trades[i].FilledQuantity
		tradeBookData.GoodTillDate = tlTradeBookResponse.Data.Trades[i].GoodTillDate
		tradeBookData.InstrumentToken = tlTradeBookResponse.Data.Trades[i].InstrumentToken
		tradeBookData.LoginID = tlTradeBookResponse.Data.Trades[i].LoginID
		tradeBookData.OmsOrderID = tlTradeBookResponse.Data.Trades[i].OmsOrderID
		tradeBookData.OrderEntryTime = tlTradeBookResponse.Data.Trades[i].OrderEntryTime
		tradeBookData.OrderPrice = tlTradeBookResponse.Data.Trades[i].OrderPrice
		tradeBookData.OrderSide = tlTradeBookResponse.Data.Trades[i].OrderSide
		tradeBookData.OrderType = tlTradeBookResponse.Data.Trades[i].OrderType
		tradeBookData.OriginalVol = tlTradeBookResponse.Data.Trades[i].OriginalVol
		tradeBookData.Pan = tlTradeBookResponse.Data.Trades[i].Pan
		tradeBookData.ProCli = tlTradeBookResponse.Data.Trades[i].ProCli
		tradeBookData.Product = tlTradeBookResponse.Data.Trades[i].Product
		tradeBookData.RemainingQuantity = tlTradeBookResponse.Data.Trades[i].RemainingQuantity
		tradeBookData.TradeNumber = tlTradeBookResponse.Data.Trades[i].TradeNumber
		tradeBookData.TradePrice = tlTradeBookResponse.Data.Trades[i].TradePrice
		tradeBookData.TradeQuantity = tlTradeBookResponse.Data.Trades[i].TradeQuantity
		tradeBookData.TradeTime = tlTradeBookResponse.Data.Trades[i].TradeTime
		tradeBookData.TradingSymbol = tlTradeBookResponse.Data.Trades[i].TradingSymbol
		tradeBookData.TriggerPrice = tlTradeBookResponse.Data.Trades[i].TriggerPrice
		tradeBookData.VLoginID = tlTradeBookResponse.Data.Trades[i].VLoginID
		tradeBookData.VolFilledToday = tlTradeBookResponse.Data.Trades[i].VolFilledToday
		if tradeBookData.Exchange == "NFO" || tradeBookData.Exchange == "BFO" || tradeBookData.Exchange == "MCX" {
			var req models.ScripInfoRequest
			req.Exchange = tradeBookData.Exchange
			req.Token = strconv.Itoa(tradeBookData.InstrumentToken)
			_, tradeBookData.AdditionalInfo = GetAdditionalInfo(obj.tradeLabURL, req, reqH)
		}
		tradeBookResponseData = append(tradeBookResponseData, tradeBookData)
	}
	tradeBookResponse.Trades = tradeBookResponseData

	maskedTradeBookResponse, err := maskObj.Struct(tradeBookResponse)
	if err != nil {
		loggerconfig.Error("TradeBook Error in masking request err: ", err, " clientId: ", req.ClientID, " requestid = ", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("TradeBook tl resp=", helpers.LogStructAsJSON(maskedTradeBookResponse), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = tradeBookResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func fetchTradeBookFromCache(req models.TradeBookRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var tradeBookResponse models.TradeBookResponse
	tradeBookResponseData := make([]models.TradeBookResponseData, 0)
	orderIdsTradeCompleted, _ := getOrderIds(req.ClientID, constants.TradeCompleted, reqH)
	for orderId, _ := range orderIdsTradeCompleted {
		tlOrderPacket, _ := getOrderStoredRedis(req.ClientID, orderId, reqH) // ignoring error as it is printed in the func and we don't want to stop the execution if failed to read from redis.

		tlTradePacket, _ := getTradeStoredRedis(req.ClientID, orderId, reqH) // ignoring error as it is printed in the func and we don't want to stop the execution if failed to read from redis.

		var tradeBookData models.TradeBookResponseData
		// tradeBookData.BookType = tlOrderPacket.BookType
		tradeBookData.BookType = constants.BookType
		tradeBookData.BrokerID = constants.BrokerId
		tradeBookData.ClientID = tlOrderPacket.ClientID
		tradeBookData.DisclosedVol = 0
		tradeBookData.DisclosedVolRemaining = 0
		tradeBookData.Exchange = tlOrderPacket.Exchange
		tradeBookData.ExchangeOrderID = tlOrderPacket.ExchangeOrderID
		tradeBookData.ExchangeTime = tlOrderPacket.ExchangeTime
		tradeBookData.FillNumber = ""
		tradeBookData.FilledQuantity = tlOrderPacket.FilledQuantity
		// tradeBookData.GoodTillDate = tlOrderPacket.GoodTillDate
		tradeBookData.GoodTillDate = 0
		tradeBookData.InstrumentToken = tlOrderPacket.InstrumentToken
		tradeBookData.LoginID = tlOrderPacket.LoginID
		tradeBookData.OmsOrderID = tlOrderPacket.OmsOrderID
		tradeBookData.OrderEntryTime = tlOrderPacket.OrderEntryTime
		tradeBookData.OrderPrice = 0
		tradeBookData.OrderSide = tlOrderPacket.OrderSide
		tradeBookData.OrderType = tlOrderPacket.OrderType
		tradeBookData.OriginalVol = tlOrderPacket.Quantity
		// tradeBookData.Pan = tlOrderPacket.Pan
		// tradeBookData.ProCli = tlOrderPacket.ProCli
		tradeBookData.Product = tlOrderPacket.Product
		tradeBookData.RemainingQuantity = tlOrderPacket.RemainingQuantity
		tradeBookData.TradeNumber = tlTradePacket.TradeID
		tradeBookData.TradePrice = tlTradePacket.TradePrice
		tradeBookData.TradeQuantity = tlOrderPacket.FilledQuantity
		tradeBookData.TradeTime = tlOrderPacket.ExchangeTime
		tradeBookData.TradingSymbol = tlOrderPacket.TradingSymbol
		tradeBookData.TriggerPrice = tlOrderPacket.TriggerPrice
		tradeBookData.VLoginID = tlTradePacket.VLoginID
		tradeBookData.VolFilledToday = tlTradePacket.VolumeFilledToday
		tradeBookResponseData = append(tradeBookResponseData, tradeBookData)
	}
	tradeBookResponse.Trades = tradeBookResponseData

	loggerconfig.Info("tradeBookResponse from redis=", helpers.LogStructAsJSON(tradeBookResponse), " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	var apiRes apihelpers.APIRes
	apiRes.Data = tradeBookResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj OrderObj) OrderHistory(req models.OrderHistoryRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + ORDERHISTORY + req.OmsOrderID + "/history" + "?client_id=" + url.QueryEscape(req.ClientID)

	//make payload
	payload := new(bytes.Buffer)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "OrderHistory", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " orderHistoryReq call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("orderHistoryReq res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlOrderHistoryResponse := TradeLabOrderHistoryResponse{}
	json.Unmarshal([]byte(string(body)), &tlOrderHistoryResponse)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " tlOrderHistoryRes tl status not ok =", tlOrderHistoryResponse.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlOrderHistoryResponse.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var orderHistoryResponse models.OrderHistoryResponse
	orderHistoryResponseData := make([]models.OrderHistoryResponseData, 0)
	for i := 0; i < len(tlOrderHistoryResponse.Data); i++ {
		var orderHistoryData models.OrderHistoryResponseData
		orderHistoryData.AvgPrice = tlOrderHistoryResponse.Data[i].AvgPrice
		orderHistoryData.ClientID = tlOrderHistoryResponse.Data[i].ClientID
		orderHistoryData.ClientOrderID = tlOrderHistoryResponse.Data[i].ClientOrderID
		orderHistoryData.CreatedAt = tlOrderHistoryResponse.Data[i].CreatedAt
		orderHistoryData.DisclosedQuantity = tlOrderHistoryResponse.Data[i].DisclosedQuantity
		orderHistoryData.Exchange = tlOrderHistoryResponse.Data[i].Exchange
		orderHistoryData.ExchangeOrderID = tlOrderHistoryResponse.Data[i].ExchangeOrderID
		orderHistoryData.ExchangeTime = tlOrderHistoryResponse.Data[i].ExchangeTime
		orderHistoryData.FillQuantity = tlOrderHistoryResponse.Data[i].FillQuantity
		orderHistoryData.LastModified = tlOrderHistoryResponse.Data[i].LastModified
		orderHistoryData.LoginID = tlOrderHistoryResponse.Data[i].LoginID
		orderHistoryData.ModifiedAt = tlOrderHistoryResponse.Data[i].ModifiedAt
		orderHistoryData.OrderID = tlOrderHistoryResponse.Data[i].OrderID
		orderHistoryData.OrderMode = tlOrderHistoryResponse.Data[i].OrderMode
		orderHistoryData.OrderSide = tlOrderHistoryResponse.Data[i].OrderSide
		orderHistoryData.OrderType = tlOrderHistoryResponse.Data[i].OrderType
		orderHistoryData.Price = tlOrderHistoryResponse.Data[i].Price
		orderHistoryData.Product = tlOrderHistoryResponse.Data[i].Product
		orderHistoryData.Quantity = tlOrderHistoryResponse.Data[i].Quantity
		orderHistoryData.RejectReason = tlOrderHistoryResponse.Data[i].RejectReason
		orderHistoryData.RemainingQuantity = tlOrderHistoryResponse.Data[i].RemainingQuantity
		orderHistoryData.Segment = tlOrderHistoryResponse.Data[i].Segment
		orderHistoryData.Status = tlOrderHistoryResponse.Data[i].Status
		orderHistoryData.Symbol = tlOrderHistoryResponse.Data[i].Symbol
		orderHistoryData.Token = tlOrderHistoryResponse.Data[i].Token
		orderHistoryData.TriggerPrice = tlOrderHistoryResponse.Data[i].TriggerPrice
		orderHistoryData.UnderlyingToken = tlOrderHistoryResponse.Data[i].UnderlyingToken
		orderHistoryData.Validity = tlOrderHistoryResponse.Data[i].Validity
		orderHistoryResponseData = append(orderHistoryResponseData, orderHistoryData)
	}
	orderHistoryResponse.OrderHistory = orderHistoryResponseData
	loggerconfig.Info("orderHistoryRes tl resp=", helpers.LogStructAsJSON(orderHistoryResponse), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = orderHistoryResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj OrderObj) PlaceGTTOrder(req models.CreateGTTOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + GTTURL

	//fill up the TL Req
	var tlPlaceGTTOrderReq TradeLabCreateGTTOrderRequest

	if req.ExpiryTime != "" {
		tlPlaceGTTOrderReq.ExpiryTime = req.ExpiryTime
	} else {
		currDate := helpers.GetCurrentTimeInIST().UTC()
		expiryDate := currDate.AddDate(1, 0, 0).Format(constants.YYYYMMDD) // currently expiry time is 1 year.
		tlPlaceGTTOrderReq.ExpiryTime = expiryDate
	}

	tlPlaceGTTOrderReq.ActionType = req.ActionType
	tlPlaceGTTOrderReq.Order.ClientID = req.ClientID
	tlPlaceGTTOrderReq.Order.Device = reqH.DeviceType
	tlPlaceGTTOrderReq.Order.DisclosedQuantity = req.DisclosedQuantity
	tlPlaceGTTOrderReq.Order.Exchange = req.Exchange
	tlPlaceGTTOrderReq.Order.InstrumentToken = req.InstrumentToken
	tlPlaceGTTOrderReq.Order.MarketProtectionPercentage = req.MarketProtectionPercentage
	tlPlaceGTTOrderReq.Order.OrderSide = req.OrderSide
	tlPlaceGTTOrderReq.Order.OrderType = req.OrderType
	tlPlaceGTTOrderReq.Order.Price = req.Price
	tlPlaceGTTOrderReq.Order.Product = req.Product
	tlPlaceGTTOrderReq.Order.Quantity = req.Quantity
	tlPlaceGTTOrderReq.Order.SlOrderPrice = req.SlOrderPrice
	tlPlaceGTTOrderReq.Order.SlOrderQuantity = req.SlOrderQuantity
	tlPlaceGTTOrderReq.Order.SlTriggerPrice = req.SlTriggerPrice
	tlPlaceGTTOrderReq.Order.TriggerPrice = req.TriggerPrice
	tlPlaceGTTOrderReq.Order.UserOrderID = req.UserOrderID

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlPlaceGTTOrderReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "PlaceGTTOrder", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " createGttOrderRes call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("createGttOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlPlaceGTTOrderRes := TradelabGTTOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlPlaceGTTOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " createGttOrderRes tl status not ok =", tlPlaceGTTOrderRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlPlaceGTTOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var gttOrderRes models.GTTOrderResponse
	gttOrderRes.ID = tlPlaceGTTOrderRes.Data.ID

	loggerconfig.Info("createGttOrderRes tl resp=", helpers.LogStructAsJSON(gttOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = gttOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj OrderObj) ModifyGTTOrder(req models.ModifyGTTOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + GTTURL

	//fill up the TL Req
	var tlModifyGTTOrderReq TradelabModifyGTTOrderRequest

	tlModifyGTTOrderReq.ActionType = req.ActionType
	tlModifyGTTOrderReq.ExpiryTime = req.ExpiryTime
	tlModifyGTTOrderReq.ID = req.ID
	tlModifyGTTOrderReq.Order.ClientID = req.Order.ClientID
	tlModifyGTTOrderReq.Order.Device = reqH.DeviceType
	tlModifyGTTOrderReq.Order.DisclosedQuantity = req.Order.DisclosedQuantity
	tlModifyGTTOrderReq.Order.Exchange = req.Order.Exchange
	tlModifyGTTOrderReq.Order.InstrumentToken = req.Order.InstrumentToken
	tlModifyGTTOrderReq.Order.MarketProtectionPercentage = req.Order.MarketProtectionPercentage
	tlModifyGTTOrderReq.Order.OrderSide = req.Order.OrderSide
	tlModifyGTTOrderReq.Order.OrderType = req.Order.OrderType
	tlModifyGTTOrderReq.Order.Price = req.Order.Price
	tlModifyGTTOrderReq.Order.Product = req.Order.Product
	tlModifyGTTOrderReq.Order.Quantity = req.Order.Quantity
	tlModifyGTTOrderReq.Order.SlOrderPrice = req.Order.SlOrderPrice
	tlModifyGTTOrderReq.Order.SlOrderQuantity = req.Order.SlOrderQuantity
	tlModifyGTTOrderReq.Order.SlTriggerPrice = req.Order.SlTriggerPrice
	tlModifyGTTOrderReq.Order.TriggerPrice = req.Order.TriggerPrice
	tlModifyGTTOrderReq.Order.UserOrderID = req.Order.UserOrderID

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlModifyGTTOrderReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPut, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ModifyGTTOrder", duration, req.Order.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " ModifyGTTOrderRes call api error =", err, " uccId:", req.Order.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("ModifyGTTOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.Order.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlModifyGTTOrderRes := TradelabGTTOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlModifyGTTOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " ModifyGTTOrderRes tl status not ok =", tlModifyGTTOrderRes.Message, " uccId:", req.Order.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlModifyGTTOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var gttOrderRes models.GTTOrderResponse
	gttOrderRes.ID = tlModifyGTTOrderRes.Data.ID

	loggerconfig.Info("modifyGTTOrderRes tl resp=", helpers.LogStructAsJSON(gttOrderRes), " uccId:", req.Order.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = gttOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj OrderObj) FetchGTTOrder(req models.FetchGTTOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + GTTURL + "/" + url.QueryEscape(req.ClientId)

	//make payload
	payload := new(bytes.Buffer)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "FetchGTTOrder", duration, req.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " fetchGTTOrderRes call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("fetchGTTOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlFetchGTTResponse := TradelabFetchGTTOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlFetchGTTResponse)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " fetchGTTOrderRes tl status not ok =", tlFetchGTTResponse.Message, " uccId:", req.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlFetchGTTResponse.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	var fetchGTTOrderRes models.FetchGTTOrderResponse

	fetchGttOrderDataResAll := make([]models.FetchGTTOrderResponseData, 0)

	for i := 0; i < len(tlFetchGTTResponse.Data); i++ {
		var fetchGttOrderDataRes models.FetchGTTOrderResponseData
		fetchGttOrderDataRes.ActionType = tlFetchGTTResponse.Data[i].ActionType
		fetchGttOrderDataRes.ClientID = tlFetchGTTResponse.Data[i].ClientID
		fetchGttOrderDataRes.CreatedAt = tlFetchGTTResponse.Data[i].CreatedAt
		fetchGttOrderDataRes.ExpiryTime = tlFetchGTTResponse.Data[i].ExpiryTime
		fetchGttOrderDataRes.ID = tlFetchGTTResponse.Data[i].ID
		fetchGttOrderDataRes.LoginID = tlFetchGTTResponse.Data[i].LoginID
		fetchGttOrderDataRes.Order.DisclosedQty = tlFetchGTTResponse.Data[i].Order.DisclosedQty
		fetchGttOrderDataRes.Order.Exchange = tlFetchGTTResponse.Data[i].Order.Exchange
		fetchGttOrderDataRes.Order.ExecutionType = tlFetchGTTResponse.Data[i].Order.ExecutionType
		fetchGttOrderDataRes.Order.Mode = tlFetchGTTResponse.Data[i].Order.Mode
		fetchGttOrderDataRes.Order.OrderSide = tlFetchGTTResponse.Data[i].Order.OrderSide
		fetchGttOrderDataRes.Order.OrderType = tlFetchGTTResponse.Data[i].Order.OrderType
		fetchGttOrderDataRes.Order.Price = tlFetchGTTResponse.Data[i].Order.Price
		fetchGttOrderDataRes.Order.ProCli = tlFetchGTTResponse.Data[i].Order.ProCli
		fetchGttOrderDataRes.Order.ProdType = tlFetchGTTResponse.Data[i].Order.ProdType
		fetchGttOrderDataRes.Order.Quantity = tlFetchGTTResponse.Data[i].Order.Quantity
		fetchGttOrderDataRes.Order.Segment = tlFetchGTTResponse.Data[i].Order.Segment
		fetchGttOrderDataRes.Order.SlOrderPrice = tlFetchGTTResponse.Data[i].Order.SlOrderPrice
		fetchGttOrderDataRes.Order.SlOrderQuantity = tlFetchGTTResponse.Data[i].Order.SlOrderQuantity
		fetchGttOrderDataRes.Order.SlTriggerPrice = tlFetchGTTResponse.Data[i].Order.SlTriggerPrice
		fetchGttOrderDataRes.Order.SquareOffPrice = tlFetchGTTResponse.Data[i].Order.SquareOffPrice
		fetchGttOrderDataRes.Order.Token = tlFetchGTTResponse.Data[i].Order.Token
		fetchGttOrderDataRes.Order.TradingSymbol = tlFetchGTTResponse.Data[i].Order.TradingSymbol
		fetchGttOrderDataRes.Order.TrailingStopLoss = tlFetchGTTResponse.Data[i].Order.TrailingStopLoss
		fetchGttOrderDataRes.Order.TriggerPrice = tlFetchGTTResponse.Data[i].Order.TriggerPrice
		fetchGttOrderDataRes.Order.Validity = tlFetchGTTResponse.Data[i].Order.Validity
		fetchGttOrderDataRes.Order.VendorCode = tlFetchGTTResponse.Data[i].Order.VendorCode
		fetchGttOrderDataRes.RejectCode = tlFetchGTTResponse.Data[i].RejectCode
		fetchGttOrderDataRes.RejectReason = tlFetchGTTResponse.Data[i].RejectReason
		fetchGttOrderDataRes.Status = tlFetchGTTResponse.Data[i].Status
		fetchGttOrderDataRes.Type = tlFetchGTTResponse.Data[i].Type
		fetchGttOrderDataRes.UpdatedAt = tlFetchGTTResponse.Data[i].UpdatedAt
		fetchGttOrderDataResAll = append(fetchGttOrderDataResAll, fetchGttOrderDataRes)
	}
	fetchGTTOrderRes.FetchGTTOrderData = fetchGttOrderDataResAll

	loggerconfig.Info("fetchGTTOrderRes tl resp=", helpers.LogStructAsJSON(fetchGTTOrderRes), " uccId:", req.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = fetchGTTOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}

func (obj OrderObj) CancelGTTOrder(req models.CancelGTTOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + CANCELGTTURL + "/" + req.ClientId + "/" + req.Id

	//make payload
	payload := new(bytes.Buffer)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodDelete, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "CancelGTTOrder", duration, req.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " cancelGTTOrderRes call api error =", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("cancelGTTOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlCancelGTTOrderRes := TradelabGTTOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlCancelGTTOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " cancelGTTOrderRes tl status not ok =", tlCancelGTTOrderRes.Message, " uccId:", req.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlCancelGTTOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var cancelGTTOrderRes models.GTTOrderResponse
	cancelGTTOrderRes.ID = tlCancelGTTOrderRes.Data.ID

	loggerconfig.Info("cancelGTTOrderRes tl resp=", helpers.LogStructAsJSON(cancelGTTOrderRes), " uccId:", req.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = cancelGTTOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj OrderObj) PlaceGttOCOOrder(req models.CreateGttOCORequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + GTTURL

	//fill up the TL Req
	var tlPlaceGTTOCOOrderReq TradeLabCreateGTTOrderRequest

	if req.ExpiryTime != "" {
		tlPlaceGTTOCOOrderReq.ExpiryTime = req.ExpiryTime
	} else {
		currDate := helpers.GetCurrentTimeInIST().UTC()
		expiryDate := currDate.AddDate(1, 0, 0).Format(constants.YYYYMMDD) //expiry time : 1 year (for now).
		tlPlaceGTTOCOOrderReq.ExpiryTime = expiryDate
	}

	tlPlaceGTTOCOOrderReq.ActionType = constants.OCO
	tlPlaceGTTOCOOrderReq.Order.ClientID = req.ClientID
	tlPlaceGTTOCOOrderReq.Order.Device = reqH.DeviceType
	tlPlaceGTTOCOOrderReq.Order.DisclosedQuantity = req.DisclosedQuantity
	tlPlaceGTTOCOOrderReq.Order.Exchange = req.Exchange
	tlPlaceGTTOCOOrderReq.Order.InstrumentToken = req.InstrumentToken
	tlPlaceGTTOCOOrderReq.Order.MarketProtectionPercentage = req.MarketProtectionPercentage
	tlPlaceGTTOCOOrderReq.Order.OrderSide = req.OrderSide
	tlPlaceGTTOCOOrderReq.Order.OrderType = req.OrderType
	tlPlaceGTTOCOOrderReq.Order.Price = req.Price
	tlPlaceGTTOCOOrderReq.Order.Product = req.Product
	tlPlaceGTTOCOOrderReq.Order.Quantity = req.Quantity
	tlPlaceGTTOCOOrderReq.Order.SlOrderPrice = req.SlOrderPrice
	tlPlaceGTTOCOOrderReq.Order.SlOrderQuantity = req.SlOrderQuantity
	tlPlaceGTTOCOOrderReq.Order.SlTriggerPrice = req.SlTriggerPrice
	tlPlaceGTTOCOOrderReq.Order.TriggerPrice = req.TriggerPrice
	tlPlaceGTTOCOOrderReq.Order.UserOrderID = req.UserOrderID

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlPlaceGTTOCOOrderReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "PlaceGTTOCOOrder", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " PlaceGttOCOOrder call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("PlaceGttOCOOrder res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlPlaceGTTOrderRes := TradelabGTTOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlPlaceGTTOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " PlaceGttOCOOrder tl status not ok =", tlPlaceGTTOrderRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlPlaceGTTOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var gttOrderRes models.GTTOrderResponse
	gttOrderRes.ID = tlPlaceGTTOrderRes.Data.ID

	loggerconfig.Info("PlaceGttOCOOrder tl resp=", helpers.LogStructAsJSON(gttOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = gttOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj OrderObj) MarginCalculations(req models.MarginCalculationRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + MARGINCALCULATION

	//fill up the TL Req
	var tlMarginCalculationRequest []models.MarginCalculationRequestData
	for i := 0; i < len(req.Data); i++ {
		var marginCalc models.MarginCalculationRequestData
		marginCalc.Exchange = req.Data[i].Exchange
		marginCalc.Mode = req.Data[i].Mode
		marginCalc.Price = req.Data[i].Price
		marginCalc.Product = req.Data[i].Product
		marginCalc.Quantity = req.Data[i].Quantity
		marginCalc.Segment = req.Data[i].Segment
		marginCalc.Series = req.Data[i].Series
		marginCalc.Side = req.Data[i].Side
		marginCalc.Symbol = req.Data[i].Symbol
		marginCalc.Token = req.Data[i].Token
		marginCalc.Underlying = req.Data[i].Underlying
		tlMarginCalculationRequest = append(tlMarginCalculationRequest, marginCalc)
	}

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlMarginCalculationRequest)

	fmt.Printf(" payload and margincal -%v \n", payload)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "MarginCalculations", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " MarginCalculationsRes call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("MarginCalculationsRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlMarginCalculationResponse := MarginCalculationResponse{}
	json.Unmarshal([]byte(string(body)), &tlMarginCalculationResponse)

	if res.StatusCode != http.StatusOK {
		if tlMarginCalculationResponse.Error.Message != RequestDataNotFound {
			loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " MarginCalculationsRes tl status not ok =", tlMarginCalculationResponse.Error.Message, "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		} else {
			loggerconfig.Error("Alert Severity:P2-Mid, platform:", reqH.Platform, " MarginCalculationsRes tl status not ok =", tlMarginCalculationResponse.Error.Message, "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		}

		apiRes.Message = tlMarginCalculationResponse.Error.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var marginCalculationResponse models.MarginResultData
	marginCalculationResponse.CombinedMargin.DeliveryMargin = tlMarginCalculationResponse.Result.CombinedMargin.DeliveryMargin
	marginCalculationResponse.CombinedMargin.Span = tlMarginCalculationResponse.Result.CombinedMargin.Span
	marginCalculationResponse.CombinedMargin.SomtierMargin = tlMarginCalculationResponse.Result.CombinedMargin.SomtierMargin
	marginCalculationResponse.CombinedMargin.AdditionalMargin = tlMarginCalculationResponse.Result.CombinedMargin.AdditionalMargin
	marginCalculationResponse.CombinedMargin.SpanSpreadMargin = tlMarginCalculationResponse.Result.CombinedMargin.SpanSpreadMargin
	marginCalculationResponse.CombinedMargin.VarMargin = tlMarginCalculationResponse.Result.CombinedMargin.VarMargin
	marginCalculationResponse.CombinedMargin.ExposureMargin = tlMarginCalculationResponse.Result.CombinedMargin.ExposureMargin
	marginCalculationResponse.CombinedMargin.PremiumMargin = tlMarginCalculationResponse.Result.CombinedMargin.PremiumMargin
	marginCalculationResponse.CombinedMargin.PremiumBenefit = tlMarginCalculationResponse.Result.CombinedMargin.PremiumBenefit
	marginCalculationResponse.CombinedMargin.ExtremeLossMargin = tlMarginCalculationResponse.Result.CombinedMargin.ExtremeLossMargin
	marginCalculationResponse.CombinedMargin.MaxSpan = tlMarginCalculationResponse.Result.CombinedMargin.MaxSpan
	marginCalculationResponse.CombinedMargin.NetSpan = tlMarginCalculationResponse.Result.CombinedMargin.NetSpan
	marginCalculationResponse.CombinedMargin.NetSpanArray = tlMarginCalculationResponse.Result.CombinedMargin.NetSpanArray
	marginCalculationResponse.CombinedMargin.CompositeDelta = tlMarginCalculationResponse.Result.CombinedMargin.CompositeDelta
	marginCalculationResponse.CombinedMargin.FutureBuyQuantity = tlMarginCalculationResponse.Result.CombinedMargin.FutureBuyQuantity
	marginCalculationResponse.CombinedMargin.FutureSellQuantity = tlMarginCalculationResponse.Result.CombinedMargin.FutureSellQuantity
	marginCalculationResponse.CombinedMargin.OptionSellQuantity = tlMarginCalculationResponse.Result.CombinedMargin.OptionSellQuantity
	marginCalculationResponse.CombinedMargin.OptionBuyQuantity = tlMarginCalculationResponse.Result.CombinedMargin.OptionBuyQuantity
	marginCalculationResponse.CombinedMargin.UnderlyingToken = tlMarginCalculationResponse.Result.CombinedMargin.UnderlyingToken
	marginCalculationResponse.CombinedMargin.SomRate = tlMarginCalculationResponse.Result.CombinedMargin.SomRate
	marginCalculationResponse.CombinedMargin.SpreadRate = tlMarginCalculationResponse.Result.CombinedMargin.SpreadRate

	responseIndividualMarginValues := make([]models.IndividualMarginValuesData, 0)

	for i := 0; i < len(tlMarginCalculationResponse.Result.IndividualMarginValues); i++ {
		var individualMarginVals models.IndividualMarginValuesData
		individualMarginVals.DeliveryMargin = tlMarginCalculationResponse.Result.IndividualMarginValues[i].DeliveryMargin
		individualMarginVals.Span = tlMarginCalculationResponse.Result.IndividualMarginValues[i].Span
		individualMarginVals.SomtierMargin = tlMarginCalculationResponse.Result.IndividualMarginValues[i].SomtierMargin
		individualMarginVals.AdditionalMargin = tlMarginCalculationResponse.Result.IndividualMarginValues[i].AdditionalMargin
		individualMarginVals.SpanSpreadMargin = tlMarginCalculationResponse.Result.IndividualMarginValues[i].SpanSpreadMargin
		individualMarginVals.VarMargin = tlMarginCalculationResponse.Result.IndividualMarginValues[i].VarMargin
		individualMarginVals.ExposureMargin = tlMarginCalculationResponse.Result.IndividualMarginValues[i].ExposureMargin
		individualMarginVals.PremiumMargin = tlMarginCalculationResponse.Result.IndividualMarginValues[i].PremiumMargin
		individualMarginVals.PremiumBenefit = tlMarginCalculationResponse.Result.IndividualMarginValues[i].PremiumBenefit
		individualMarginVals.ExtremeLossMargin = tlMarginCalculationResponse.Result.IndividualMarginValues[i].ExtremeLossMargin
		individualMarginVals.MaxSpan = tlMarginCalculationResponse.Result.IndividualMarginValues[i].MaxSpan
		individualMarginVals.NetSpan = tlMarginCalculationResponse.Result.IndividualMarginValues[i].NetSpan
		individualMarginVals.NetSpanArray = tlMarginCalculationResponse.Result.IndividualMarginValues[i].NetSpanArray
		individualMarginVals.CompositeDelta = tlMarginCalculationResponse.Result.IndividualMarginValues[i].CompositeDelta
		individualMarginVals.FutureBuyQuantity = tlMarginCalculationResponse.Result.IndividualMarginValues[i].FutureBuyQuantity
		individualMarginVals.FutureSellQuantity = tlMarginCalculationResponse.Result.IndividualMarginValues[i].FutureSellQuantity
		individualMarginVals.OptionSellQuantity = tlMarginCalculationResponse.Result.IndividualMarginValues[i].OptionSellQuantity
		individualMarginVals.OptionBuyQuantity = tlMarginCalculationResponse.Result.IndividualMarginValues[i].OptionBuyQuantity
		individualMarginVals.UnderlyingToken = tlMarginCalculationResponse.Result.IndividualMarginValues[i].UnderlyingToken
		individualMarginVals.SomRate = tlMarginCalculationResponse.Result.IndividualMarginValues[i].SomRate
		individualMarginVals.SpreadRate = tlMarginCalculationResponse.Result.IndividualMarginValues[i].SpreadRate

		responseIndividualMarginValues = append(responseIndividualMarginValues, individualMarginVals)
	}
	marginCalculationResponse.IndividualMarginValues = responseIndividualMarginValues

	loggerconfig.Info("MarginCalculationsRes tl resp=", helpers.LogStructAsJSON(tlMarginCalculationResponse), "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = marginCalculationResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj OrderObj) LastTradedPrice(req models.LastTradedPriceRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + LastTradedPrice

	var exchange string
	var segment string

	if strings.ToLower(req.Exchange) == constants.NFO {
		if strings.ToLower(req.Segment) == constants.FUTOPT {
			exchange = "NSE"
			segment = "FutOpt"
		}
	} else if strings.ToLower(req.Exchange) == constants.CDS {
		if strings.ToLower(req.Segment) == constants.CURRENCY {
			exchange = "NSE"
			segment = "Currency"
		}
	} else if strings.ToLower(req.Exchange) == constants.NSE {
		if strings.ToLower(req.Segment) == constants.EQUITY {
			exchange = "NSE"
			segment = "Capital"
		}
	} else if strings.ToLower(req.Exchange) == constants.MCX {
		if strings.ToLower(req.Segment) == constants.COMMODITY {
			exchange = "MCX"
			segment = "FutOpt"
		}
	} else if strings.ToLower(req.Exchange) == constants.BSE {
		if strings.ToLower(req.Segment) == constants.EQUITY {
			exchange = "BSE"
			segment = "Capital"
		}
	} else if strings.ToLower(req.Exchange) == constants.BFO {
		if strings.ToLower(req.Segment) == constants.FUTOPT {
			exchange = "BSE"
			segment = "FutOpt"
		}
	}

	if strings.ToLower(req.Exchange) == constants.MCX && constants.UseMCXLtpUrl {
		url += "/" + exchange + "/" + req.Token
	} else {
		url += "/" + exchange + "/" + segment + "?token=" + req.Token + "&key=last_trade_price"
	}

	// empty payload
	payload := new(bytes.Buffer)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "LastTradedPrice", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " LastTradedPriceRes call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("LastTradedPrice res error =", tlErrorRes.Message, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	if strings.ToLower(req.Exchange) == constants.MCX && constants.UseMCXLtpUrl {

		mcxResponse := TradeLabMCXLastTradedPrice{}
		err = json.Unmarshal([]byte(string(body)), &mcxResponse)
		if err != nil {
			loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " LastTradedPrice tl status not ok =", mcxResponse.Message, "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
			return apihelpers.SendInternalServerError()
		}

		if res.StatusCode != http.StatusOK {
			loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " LastTradedPrice tl status not ok =", mcxResponse.Message, "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
			apiRes.Message = mcxResponse.Message
			apiRes.Status = false
			return res.StatusCode, apiRes
		}

		var lastTradedPriceRes models.LastTradedPriceResponse
		lastTradedPriceRes.Price = mcxResponse.Data.LastTradePrice

		apiRes.Data = lastTradedPriceRes
		apiRes.Message = "SUCCESS"
		apiRes.Status = true
		return http.StatusOK, apiRes

	} else {

		tlLastTradedPriceResponse := TradeLabLastTradedPrice{}
		json.Unmarshal([]byte(string(body)), &tlLastTradedPriceResponse)

		if res.StatusCode != http.StatusOK {
			loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " LastTradedPrice tl status not ok =", tlLastTradedPriceResponse.Message, "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
			apiRes.Message = tlLastTradedPriceResponse.Message
			apiRes.Status = false
			return res.StatusCode, apiRes
		}

		//fill up controller response
		var lastTradedPriceRes models.LastTradedPriceResponse
		lastTradedPriceRes.Price = tlLastTradedPriceResponse.Data

		loggerconfig.Info("LastTradedPrice tl resp=", helpers.LogStructAsJSON(lastTradedPriceRes), "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

		apiRes.Data = lastTradedPriceRes
		apiRes.Message = "SUCCESS"
		apiRes.Status = true
		return http.StatusOK, apiRes
	}
}

func (obj OrderObj) PlaceIcebergOrder(req models.IcebergOrderReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + PLACEICEBERGORDERURL

	// assign all the values of IcebergOrderReq into TLIcebergOrderReq in given sequence
	var tlPlaceIcebergOrder TLIcebergOrderReq
	tlPlaceIcebergOrder.Order.Exchange = req.Exchange
	tlPlaceIcebergOrder.Order.InstrumentToken = req.InstrumentToken
	tlPlaceIcebergOrder.Order.OrderSide = req.OrderSide
	tlPlaceIcebergOrder.Order.OrderType = req.OrderType
	tlPlaceIcebergOrder.Order.Quantity = req.Quantity
	tlPlaceIcebergOrder.Order.DisclosedQuantity = req.DisclosedQuantity
	tlPlaceIcebergOrder.Order.Validity = req.Validity
	tlPlaceIcebergOrder.Order.Product = req.Product
	tlPlaceIcebergOrder.Order.NoOfLegs = req.NoOfLegs
	tlPlaceIcebergOrder.Order.GttPrice = req.GttPrice
	tlPlaceIcebergOrder.Order.Price = req.Price
	tlPlaceIcebergOrder.Order.ClientId = req.ClientID
	tlPlaceIcebergOrder.Order.UserOrderId = req.UserOrderId
	tlPlaceIcebergOrder.Order.Device = reqH.DeviceType
	tlPlaceIcebergOrder.Order.ExecutionType = req.ExecutionType

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlPlaceIcebergOrder)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "PlaceIcebergOrder", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " PlaceIcebergOrder call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	if res == nil {
		loggerconfig.Error("PlaceIcebergOrder: API response is nil")
		return apihelpers.SendInternalServerError()
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("PlaceIcebergOrder res error =", tlErrorRes.Message, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlPlaceIcebergOrderRes := TradelabGTTOrderResponse{} //need to update
	json.Unmarshal([]byte(string(body)), &tlPlaceIcebergOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " PlaceIcebergOrder tl status not ok =", tlPlaceIcebergOrderRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlPlaceIcebergOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("PlaceIcebergOrder tl resp=", helpers.LogStructAsJSON(tlPlaceIcebergOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	var iceberOrderRes models.IcebergOrderResponse
	iceberOrderRes.ID = tlPlaceIcebergOrderRes.Data.ID

	apiRes.Status = false
	apiRes.Message = "SUCCESS"
	apiRes.Data = iceberOrderRes

	return http.StatusOK, apiRes

}

func (obj OrderObj) ModifyIcebergOrder(req models.ModifyIcebergOrderReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + PLACEORDERURL

	// assign all the values of modifyIcebergOrder into TLModifyIcebergOrder in given sequence
	var tlModifyIcebergOrder TLModifyIcebergOrderReq
	tlModifyIcebergOrder.Exchange = req.Exchange
	tlModifyIcebergOrder.InstrumentToken = req.InstrumentToken
	tlModifyIcebergOrder.ClientId = req.ClientID
	tlModifyIcebergOrder.OrderType = req.OrderType
	tlModifyIcebergOrder.Price = req.Price
	tlModifyIcebergOrder.Quantity = req.Quantity
	tlModifyIcebergOrder.DisclosedQuantity = req.DisclosedQuantity
	tlModifyIcebergOrder.Validity = req.Validity
	tlModifyIcebergOrder.Product = req.Product
	tlModifyIcebergOrder.GttPrice = req.GttPrice
	tlModifyIcebergOrder.OmsOrderId = req.OmsOrderId
	tlModifyIcebergOrder.ExchangeOrderId = req.ExchangeOrderId
	tlModifyIcebergOrder.RemainingQuantity = req.RemainingQuantity
	tlModifyIcebergOrder.LastActivityReference = req.LastActivityReference
	tlModifyIcebergOrder.TriggerPrice = req.TriggerPrice
	tlModifyIcebergOrder.ExecutionType = req.ExecutionType
	tlModifyIcebergOrder.MarketProtectionPercentage = req.MarketProtectionPercentage

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlModifyIcebergOrder)
	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPut, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ModifyIcebergOrder", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " ModifyIcebergOrder call api error =", err, "clientID: ", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	if res == nil {
		loggerconfig.Error("ModifyIcebergOrder: API response is nil")
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("ModifyIcebergOrder res error =", tlErrorRes.Message, "clientID: ", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}
	tlModifyIcebergOrderRes := TradelabGTTOrderResponse{} // Assuming we'll use a generic order response structure
	json.Unmarshal([]byte(string(body)), &tlModifyIcebergOrderRes)
	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " ModifyIcebergOrder tl status not ok =", tlModifyIcebergOrderRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlModifyIcebergOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}
	loggerconfig.Info("ModifyIcebergOrder tl resp=", helpers.LogStructAsJSON(tlModifyIcebergOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	var iceberOrderRes models.IcebergOrderResponse
	iceberOrderRes.ID = tlModifyIcebergOrderRes.Data.ID

	apiRes.Status = true
	apiRes.Message = "SUCCESS"
	apiRes.Data = iceberOrderRes
	return http.StatusOK, apiRes
}

func (obj OrderObj) CancelIcebergOrder(req models.CancelIcebergOrderReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + PLACEORDERURL + "/" + req.OmsOrderID + "?client_id=" + req.ClientId + "&execution_type=" + req.ExecutionType

	//empty payload
	payload := new(bytes.Buffer) //test if we got executionType

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodDelete, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "CancelIcebergOrder", duration, req.ClientId, reqH.RequestId)

	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " CancelIcebergOrder call api error =", err, "clientID: ", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	if res == nil {
		loggerconfig.Error("CancelIcebergOrder: API response is nil")
		return apihelpers.SendInternalServerError()
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("CancelIcebergOrder res error =", tlErrorRes.Message, "clientID: ", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlCancelOrderRes := TradelabCancelOrModifyResponse{} // Using the same response structure
	json.Unmarshal([]byte(string(body)), &tlCancelOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " CancelIcebergOrder tl status not ok =", tlCancelOrderRes.Message, " uccId:", req.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlCancelOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("CancelIcebergOrder tl resp=", helpers.LogStructAsJSON(tlCancelOrderRes), " uccId:", req.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	var icebergCancelOrderRes models.IcebergCanelOrderResponse
	icebergCancelOrderRes.OMSOrderID = tlCancelOrderRes.Data.OmsOrderID

	apiRes.Status = true
	apiRes.Message = "SUCCESS"
	apiRes.Data = icebergCancelOrderRes
	return http.StatusOK, apiRes
}
