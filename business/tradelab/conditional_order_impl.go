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
	"space/dbops"
	"space/helpers"
	"space/helpers/cache"
	"space/loggerconfig"
	"space/models"
)

type ConditionalOrderObj struct {
	tradeLabURL string
	redisCli    cache.RedisCache
}

func InitConditionalOrder(redisCli cache.RedisCache) ConditionalOrderObj {
	defer models.HandlePanic()

	orderObj := ConditionalOrderObj{
		tradeLabURL: constants.TLURL,
		redisCli:    redisCli,
	}

	return orderObj
}

func (obj ConditionalOrderObj) PlaceBOOrder(req models.PlaceBOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + CONDATIONALORDERSURL

	var tlPlaceBOOrderReq TradeLabPlaceBOOrderRequest
	tlPlaceBOOrderReq.ClientID = req.ClientID
	tlPlaceBOOrderReq.Device = reqH.DeviceType
	tlPlaceBOOrderReq.DisclosedQuantity = req.DisclosedQuantity
	tlPlaceBOOrderReq.Exchange = req.Exchange
	tlPlaceBOOrderReq.ExecutionType = "BO"
	tlPlaceBOOrderReq.InstrumentToken = req.InstrumentToken
	tlPlaceBOOrderReq.IsTrailing = req.IsTrailing
	tlPlaceBOOrderReq.OrderSide = req.OrderSide
	tlPlaceBOOrderReq.OrderType = req.OrderType
	tlPlaceBOOrderReq.Price = float64(req.Price)
	tlPlaceBOOrderReq.Product = req.Product
	tlPlaceBOOrderReq.Quantity = req.Quantity
	tlPlaceBOOrderReq.SquareOffValue = req.SquareOffValue
	tlPlaceBOOrderReq.StopLossValue = req.StopLossValue
	tlPlaceBOOrderReq.TrailingStopLoss = req.TrailingStopLoss
	tlPlaceBOOrderReq.TriggerPrice = req.TriggerPrice
	tlPlaceBOOrderReq.UserOrderID = int(dbops.RedisRepo.Increment(USERORDERIDKEY))
	tlPlaceBOOrderReq.Validity = req.Validity

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlPlaceBOOrderReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "PlaceBOOrder", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		// logObj.Printf("BOOrderRes call api error =", err)
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " BOOrderRes call api error", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("BOOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlPlaceBOOrderRes := TradeLabPlaceBOOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlPlaceBOOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " BOOrderRes tl status not ok =", tlPlaceBOOrderRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlPlaceBOOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var boOrderRes models.BOOrderResponse
	boOrderRes.BasketID = tlPlaceBOOrderRes.Data.Data.BasketID
	boOrderRes.Message = tlPlaceBOOrderRes.Data.Data.Message

	loggerconfig.Info("BOOrderRes tl resp=", helpers.LogStructAsJSON(boOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = boOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ConditionalOrderObj) ModifyBOOrder(req models.ModifyBOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + CONDATIONALORDERSURL

	var tlModifyBOOrderReq TradeLabModifyBOOrderRequest
	tlModifyBOOrderReq.Exchange = req.Exchange
	tlModifyBOOrderReq.InstrumentToken = req.InstrumentToken
	tlModifyBOOrderReq.ClientID = req.ClientID
	tlModifyBOOrderReq.OrderType = req.OrderType
	tlModifyBOOrderReq.Price = req.Price
	tlModifyBOOrderReq.Quantity = req.Quantity
	tlModifyBOOrderReq.DisclosedQuantity = req.DisclosedQuantity
	tlModifyBOOrderReq.Validity = req.Validity
	tlModifyBOOrderReq.Product = req.Product
	tlModifyBOOrderReq.OmsOrderID = req.OmsOrderID
	tlModifyBOOrderReq.ExchangeOrderID = req.ExchangeOrderID
	tlModifyBOOrderReq.FilledQuantity = req.FilledQuantity
	tlModifyBOOrderReq.RemainingQuantity = req.RemainingQuantity
	tlModifyBOOrderReq.LastActivityReference = req.LastActivityReference
	tlModifyBOOrderReq.TriggerPrice = req.TriggerPrice
	tlModifyBOOrderReq.StopLossValue = req.StopLossValue
	tlModifyBOOrderReq.SquareOffValue = req.SquareOffValue
	tlModifyBOOrderReq.TrailingStopLoss = req.TrailingStopLoss
	tlModifyBOOrderReq.IsTrailing = req.IsTrailing
	tlModifyBOOrderReq.ExecutionType = "BO"

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlModifyBOOrderReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPut, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ModifyBOOrder", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("modifyBOOrderRes call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("modifyBOOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlBOModifyOrderRes := TradeLabModOrExitBOOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlBOModifyOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("modifyBOOrderRes tl status not ok =", tlBOModifyOrderRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlBOModifyOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var modifyBOOrderRes models.SpreadOrderResponse
	modifyBOOrderRes.BasketID = tlBOModifyOrderRes.Data.BasketID
	modifyBOOrderRes.Message = tlBOModifyOrderRes.Data.Message

	loggerconfig.Info("modifyBOOrderRes tl resp=", helpers.LogStructAsJSON(modifyBOOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = modifyBOOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ConditionalOrderObj) CancelBOOrder(req models.ExitBOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + CONDATIONALORDERSURL + "/" + req.OmsOrderID

	//fill up the TL Req
	var tlCancelBOOrderReq TradeLabCancelBOOrderRequest
	tlCancelBOOrderReq.ClientID = req.ClientID
	tlCancelBOOrderReq.ExchangeOrderID = req.ExchangeOrderID
	tlCancelBOOrderReq.ExecutionType = "BO"
	tlCancelBOOrderReq.LegOrderIndicator = req.LegOrderIndicator
	tlCancelBOOrderReq.OmsOrderID = req.OmsOrderID
	tlCancelBOOrderReq.Status = req.Status

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlCancelBOOrderReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodDelete, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "CancelBOOrder", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " cancelBOOrderRes call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("cancelBOOrderRes res error =", tlErrorRes.Message, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlCancelBOOrderRes := TradeLabModOrExitBOOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlCancelBOOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " cancelBOOrderRes tl status not ok =", tlCancelBOOrderRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlCancelBOOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var cancelBOOrderRes models.BOOrderResponse
	cancelBOOrderRes.BasketID = tlCancelBOOrderRes.Data.BasketID
	cancelBOOrderRes.Message = tlCancelBOOrderRes.Data.Message

	loggerconfig.Info("cancelBOOrderRes tl resp=", helpers.LogStructAsJSON(cancelBOOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = cancelBOOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ConditionalOrderObj) PlaceCOOrder(req models.PlaceCOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + CONDATIONALORDERSURL

	var tlPlaceCOOrderReq TradeLabPlaceCOOrderRequest
	tlPlaceCOOrderReq.Exchange = req.Exchange
	tlPlaceCOOrderReq.InstrumentToken = req.InstrumentToken
	tlPlaceCOOrderReq.ClientID = req.ClientID
	tlPlaceCOOrderReq.OrderType = req.OrderType
	tlPlaceCOOrderReq.Price = float64(req.Price)
	tlPlaceCOOrderReq.Quantity = req.Quantity
	tlPlaceCOOrderReq.DisclosedQuantity = req.DisclosedQuantity
	tlPlaceCOOrderReq.Validity = req.Validity
	tlPlaceCOOrderReq.Product = req.Product
	tlPlaceCOOrderReq.OrderSide = req.OrderSide
	tlPlaceCOOrderReq.Device = reqH.DeviceType
	tlPlaceCOOrderReq.UserOrderID = int(dbops.RedisRepo.Increment(USERORDERIDKEY))
	tlPlaceCOOrderReq.ExecutionType = "CO"
	tlPlaceCOOrderReq.StopLossValue = req.StopLossValue
	tlPlaceCOOrderReq.TrailingStopLoss = req.TrailingStopLoss

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlPlaceCOOrderReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "PlaceCOOrder", duration, req.ClientID, reqH.RequestId)

	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " COOrderRes call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("COOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlPlaceCOOrderRes := TradeLabPlaceCOOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlPlaceCOOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " COOrderRes tl status not ok =", tlPlaceCOOrderRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlPlaceCOOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var coOrderRes models.COOrderResponse
	coOrderRes.BasketID = tlPlaceCOOrderRes.Data.Data.BasketID
	coOrderRes.Message = tlPlaceCOOrderRes.Data.Data.Message

	loggerconfig.Info("COOrderRes tl resp=", helpers.LogStructAsJSON(coOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = coOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ConditionalOrderObj) ModifyCOOrder(req models.ModifyCOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + CONDATIONALORDERSURL

	var tlModifyCOOrderReq TradeLabModifyCOOrderRequest
	tlModifyCOOrderReq.ClientID = req.ClientID
	tlModifyCOOrderReq.DisclosedQuantity = req.DisclosedQuantity
	tlModifyCOOrderReq.Exchange = req.Exchange
	tlModifyCOOrderReq.ExchangeOrderID = req.ExchangeOrderID
	tlModifyCOOrderReq.ExecutionType = "CO"
	tlModifyCOOrderReq.FilledQuantity = req.FilledQuantity
	tlModifyCOOrderReq.InstrumentToken = req.InstrumentToken
	tlModifyCOOrderReq.LastActivityReference = req.LastActivityReference
	tlModifyCOOrderReq.OmsOrderID = req.OmsOrderID
	tlModifyCOOrderReq.OrderType = req.OrderType
	tlModifyCOOrderReq.Price = req.Price
	tlModifyCOOrderReq.Product = req.Product
	tlModifyCOOrderReq.Quantity = req.Quantity
	tlModifyCOOrderReq.RemainingQuantity = req.RemainingQuantity
	tlModifyCOOrderReq.StopLossValue = req.StopLossValue
	tlModifyCOOrderReq.TrailingStopLoss = req.TrailingStopLoss
	tlModifyCOOrderReq.Validity = req.Validity
	tlModifyCOOrderReq.LegOrderIndicator = req.LegOrderIndicator
	tlModifyCOOrderReq.TriggerPrice = req.TriggerPrice

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlModifyCOOrderReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPut, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ModifyCOOrder", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " modifyCOOrderRes call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("modifyCOOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlCOModifyOrderRes := TradeLabModifyOrExitCOOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlCOModifyOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " modifyCOOrderRes tl status not ok =", tlCOModifyOrderRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlCOModifyOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var modifyCOOrderRes models.COOrderResponse
	modifyCOOrderRes.BasketID = tlCOModifyOrderRes.Data.BasketID
	modifyCOOrderRes.Message = tlCOModifyOrderRes.Data.Message

	loggerconfig.Info("modifyCOOrderRes tl resp=", helpers.LogStructAsJSON(modifyCOOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = modifyCOOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ConditionalOrderObj) CancelCOOrder(req models.ExitCOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + CONDATIONALORDERSURL + "/" + req.OmsOrderID

	//fill up the TL Req
	var tlCancelCOOrderReq TradeLabCancelCOOrderRequest
	tlCancelCOOrderReq.ClientID = req.ClientID
	tlCancelCOOrderReq.ExchangeOrderID = req.ExchangeOrderID
	tlCancelCOOrderReq.ExecutionType = "CO"
	tlCancelCOOrderReq.LegOrderIndicator = req.LegOrderIndicator
	tlCancelCOOrderReq.OmsOrderID = req.OmsOrderID

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlCancelCOOrderReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodDelete, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "CancelCOOrder", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " cancelCOOrderRes call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("cancelCOOrderRes res error =", tlErrorRes.Message, " uccId:", req.ClientID, " statuscode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlCancelCOOrderRes := TradeLabModifyOrExitCOOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlCancelCOOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " cancelBOOrderRes tl status not ok =", tlCancelCOOrderRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlCancelCOOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var cancelCOOrderRes models.COOrderResponse
	cancelCOOrderRes.BasketID = tlCancelCOOrderRes.Data.BasketID
	cancelCOOrderRes.Message = tlCancelCOOrderRes.Data.Message

	loggerconfig.Info("cancelBOOrderRes tl resp=", helpers.LogStructAsJSON(cancelCOOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = cancelCOOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ConditionalOrderObj) PlaceSpreadOrder(req models.PlaceSpreadOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + CONDATIONALORDERSURL

	//fill up the TL Req
	var tlPlaceSpreadOrderReq TradelabPlaceSpreadOrderRequest
	tlPlaceSpreadOrderReq.ClientID = req.ClientID
	tlPlaceSpreadOrderReq.Device = reqH.DeviceType
	tlPlaceSpreadOrderReq.DisclosedQuantity = req.DisclosedQuantity
	tlPlaceSpreadOrderReq.Exchange = req.Exchange
	tlPlaceSpreadOrderReq.ExecutionType = "SPD"
	tlPlaceSpreadOrderReq.InstrumentToken = req.InstrumentToken
	tlPlaceSpreadOrderReq.OrderSide = req.OrderSide
	tlPlaceSpreadOrderReq.OrderType = req.OrderType
	tlPlaceSpreadOrderReq.Price = float64(req.Price)
	tlPlaceSpreadOrderReq.Product = req.Product
	tlPlaceSpreadOrderReq.Quantity = req.Quantity
	tlPlaceSpreadOrderReq.UserOrderID = int(dbops.RedisRepo.Increment(USERORDERIDKEY))
	tlPlaceSpreadOrderReq.Validity = req.Validity

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlPlaceSpreadOrderReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "PlaceSpreadOrder", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " spreadOrderRes call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("spreadOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlPlaceSpreadOrderRes := TradelabSpreadOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlPlaceSpreadOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " spreadOrderRes tl status not ok =", tlPlaceSpreadOrderRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlPlaceSpreadOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var spreadOrderRes models.SpreadOrderResponse
	spreadOrderRes.BasketID = tlPlaceSpreadOrderRes.Data.Data.BasketID
	spreadOrderRes.Message = tlPlaceSpreadOrderRes.Data.Data.Message

	loggerconfig.Info("spreadOrderRes tl resp=", helpers.LogStructAsJSON(spreadOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = spreadOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ConditionalOrderObj) ModifySpreadOrder(req models.ModifySpreadOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + CONDATIONALORDERSURL

	//fill up the TL Req
	var tlSpreadModifyOrderReq TradelabModifySpreadOrderRequest
	tlSpreadModifyOrderReq.ClientID = req.ClientID
	tlSpreadModifyOrderReq.DisclosedQuantity = req.DisclosedQuantity
	tlSpreadModifyOrderReq.ExchangeOrderID = req.ExchangeOrderID
	tlSpreadModifyOrderReq.Exchange = req.Exchange
	tlSpreadModifyOrderReq.InstrumentToken = req.InstrumentToken
	tlSpreadModifyOrderReq.IsTrailing = req.IsTrailing
	tlSpreadModifyOrderReq.ExecutionType = "SPD"
	tlSpreadModifyOrderReq.OrderType = req.OrderType
	tlSpreadModifyOrderReq.Price = float64(req.Price)
	tlSpreadModifyOrderReq.ProdType = req.ProdType
	tlSpreadModifyOrderReq.Product = req.Product
	tlSpreadModifyOrderReq.Quantity = req.Quantity
	tlSpreadModifyOrderReq.TriggerPrice = float64(req.TriggerPrice)
	tlSpreadModifyOrderReq.OmsOrderID = req.OmsOrderID
	tlSpreadModifyOrderReq.Validity = req.Validity
	tlSpreadModifyOrderReq.SquareOffValue = req.SquareOffValue
	tlSpreadModifyOrderReq.StopLossValue = req.StopLossValue
	tlSpreadModifyOrderReq.TrailingStopLoss = req.TrailingStopLoss

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlSpreadModifyOrderReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPut, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ModifySpreadOrder", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " modifySpreadOrderRes call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("modifySpreadOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlSpreadModifyOrderRes := TradelabSpreadOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlSpreadModifyOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " modifySpreadOrderRes tl status not ok =", tlSpreadModifyOrderRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlSpreadModifyOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var modifySpreadOrderRes models.SpreadOrderResponse
	modifySpreadOrderRes.BasketID = tlSpreadModifyOrderRes.Data.Data.BasketID
	modifySpreadOrderRes.Message = tlSpreadModifyOrderRes.Data.Data.Message

	loggerconfig.Info("modifySpreadOrderRes tl resp=", helpers.LogStructAsJSON(modifySpreadOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = modifySpreadOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ConditionalOrderObj) CancelSpreadOrder(req models.ExitSpreadOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + CONDATIONALORDERSURL + "/" + req.OmsOrderID

	//fill up the TL Req
	var tlCancelSpreadOrderReq TradelabCancelSpreadOrderRequest
	tlCancelSpreadOrderReq.ClientID = req.ClientID
	tlCancelSpreadOrderReq.LegOrderIndicator = req.LegOrderIndicator
	tlCancelSpreadOrderReq.OmsOrderID = req.OmsOrderID
	tlCancelSpreadOrderReq.Status = req.Status
	tlCancelSpreadOrderReq.ExecutionType = "SPD"
	tlCancelSpreadOrderReq.ExchangeOrderID = req.ExchangeOrderID

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlCancelSpreadOrderReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodDelete, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "CancelSpreadOrder", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " cancelSpreadOrderRes call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("cancelSpreadOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlCancelSpreadOrderRes := TradelabExitSpreadOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlCancelSpreadOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " cancelSpreadOrderRes tl status not ok =", tlCancelSpreadOrderRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlCancelSpreadOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var cancelSpreadOrderRes models.SpreadOrderResponse
	cancelSpreadOrderRes.BasketID = tlCancelSpreadOrderRes.Data.BasketID
	cancelSpreadOrderRes.Message = tlCancelSpreadOrderRes.Data.Message

	loggerconfig.Info("cancelSpreadOrderRes tl resp=", helpers.LogStructAsJSON(cancelSpreadOrderRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = cancelSpreadOrderRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}
