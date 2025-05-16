package tradelab

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
)

type BasketOrderObj struct {
	tradeLabURL string
}

func InitBasketOrder() BasketOrderObj {
	defer models.HandlePanic()

	basketOrderObj := BasketOrderObj{
		tradeLabURL: constants.TLURL,
	}

	return basketOrderObj
}

func (obj BasketOrderObj) CreateBasket(req models.CreateBasketReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + BasketURL

	var tlCreateBasketReq TradeLabCreateBasketReq
	tlCreateBasketReq.LoginID = req.LoginID
	tlCreateBasketReq.Name = req.Name
	tlCreateBasketReq.Type = req.Type
	tlCreateBasketReq.ProductType = req.ProductType
	tlCreateBasketReq.OrderType = req.OrderType

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlCreateBasketReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "CreateBasket", duration, req.LoginID, reqH.RequestId)
	if err != nil {
		// logObj.Printf("BOOrderRes call api error =", err)
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " CreateBasketRes call api error", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("CreateBasketRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlCreateBasketRes := TradeLabBasketRes{}
	json.Unmarshal([]byte(string(body)), &tlCreateBasketRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " CreateBasketRes tl status not ok =", tlCreateBasketRes.Message, " statuscode: ", res.StatusCode, " StatusCode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		apiRes.Message = tlCreateBasketRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	// var createBasketRes models.CreateBasketRes
	// createBasketRes.LoginID = req.LoginID
	// createBasketRes.Name = req.Name
	// createBasketRes.Type = req.Type
	// createBasketRes.ProductType = req.ProductType
	// createBasketRes.OrderType = req.OrderType
	// createBasketRes.

	var createBasketRes models.BasketDataRes

	for i := 0; i < len(tlCreateBasketRes.Data); i++ {

		if !strings.EqualFold(tlCreateBasketRes.Data[i].Name, req.Name) {
			continue
		}

		var createBasketDataRes models.BasketDataRes
		createBasketDataRes.BasketID = tlCreateBasketRes.Data[i].BasketID
		createBasketDataRes.BasketType = tlCreateBasketRes.Data[i].BasketType
		createBasketDataRes.IsExecuted = tlCreateBasketRes.Data[i].IsExecuted
		createBasketDataRes.LoginID = tlCreateBasketRes.Data[i].LoginID
		createBasketDataRes.Name = tlCreateBasketRes.Data[i].Name
		createBasketDataRes.OrderType = tlCreateBasketRes.Data[i].OrderType
		createBasketDataRes.ProductType = tlCreateBasketRes.Data[i].ProductType
		createBasketDataRes.SipEligible = tlCreateBasketRes.Data[i].SipEligible
		createBasketDataRes.SipEnabled = tlCreateBasketRes.Data[i].SipEnabled

		createBasketDataOrderResAll := make([]models.BasketDataOrderRes, 0)
		for j := 0; j < len(tlCreateBasketRes.Data[i].Orders); j++ {
			var createBasketDataOrderRes models.BasketDataOrderRes
			createBasketDataOrderRes.OrderID = tlCreateBasketRes.Data[i].Orders[j].OrderID
			createBasketDataOrderRes.OrderInfo.TriggerPrice = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.TriggerPrice
			createBasketDataOrderRes.OrderInfo.UnderlyingToken = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.UnderlyingToken
			createBasketDataOrderRes.OrderInfo.Series = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.Series
			createBasketDataOrderRes.OrderInfo.UserOrderID = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.UserOrderID
			createBasketDataOrderRes.OrderInfo.Exchange = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.Exchange
			createBasketDataOrderRes.OrderInfo.SquareOff = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.SquareOff
			createBasketDataOrderRes.OrderInfo.Mode = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.Mode
			createBasketDataOrderRes.OrderInfo.RemainingQuantity = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.RemainingQuantity
			createBasketDataOrderRes.OrderInfo.AverageTradePrice = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.AverageTradePrice
			createBasketDataOrderRes.OrderInfo.TradePrice = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.TradePrice
			createBasketDataOrderRes.OrderInfo.OrderTag = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.OrderTag
			createBasketDataOrderRes.OrderInfo.OrderStatusInfo = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.OrderStatusInfo
			createBasketDataOrderRes.OrderInfo.OrderSide = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.OrderSide
			createBasketDataOrderRes.OrderInfo.SquareOffValue = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.SquareOffValue
			createBasketDataOrderRes.OrderInfo.ContractDescription = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.ContractDescription
			createBasketDataOrderRes.OrderInfo.Segment = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.Segment
			createBasketDataOrderRes.OrderInfo.ClientID = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.ClientID
			createBasketDataOrderRes.OrderInfo.TradingSymbol = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.TradingSymbol
			createBasketDataOrderRes.OrderInfo.RejectionCode = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.RejectionCode
			createBasketDataOrderRes.OrderInfo.LotSize = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.LotSize
			createBasketDataOrderRes.OrderInfo.Quantity = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.Quantity
			createBasketDataOrderRes.OrderInfo.LastActivityReference = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.LastActivityReference
			createBasketDataOrderRes.OrderInfo.NnfID = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.NnfID
			createBasketDataOrderRes.OrderInfo.ProCli = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.ProCli
			createBasketDataOrderRes.OrderInfo.Price = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.Price
			createBasketDataOrderRes.OrderInfo.OrderType = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.OrderType
			createBasketDataOrderRes.OrderInfo.Validity = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.Validity
			createBasketDataOrderRes.OrderInfo.TargetPriceType = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.TargetPriceType
			createBasketDataOrderRes.OrderInfo.InstrumentToken = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.InstrumentToken
			createBasketDataOrderRes.OrderInfo.SlTriggerPrice = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.SlTriggerPrice
			createBasketDataOrderRes.OrderInfo.IsTrailing = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.IsTrailing
			createBasketDataOrderRes.OrderInfo.SlOrderQuantity = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.SlOrderQuantity
			createBasketDataOrderRes.OrderInfo.OrderEntryTime = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.OrderEntryTime
			createBasketDataOrderRes.OrderInfo.ExchangeTime = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.ExchangeTime
			createBasketDataOrderRes.OrderInfo.LegOrderIndicator = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.LegOrderIndicator
			createBasketDataOrderRes.OrderInfo.TrailingStopLoss = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.TrailingStopLoss
			createBasketDataOrderRes.OrderInfo.LoginID = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.LoginID
			createBasketDataOrderRes.OrderInfo.OmsOrderID = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.OmsOrderID
			createBasketDataOrderRes.OrderInfo.MarketProtectionPercentage = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.MarketProtectionPercentage
			createBasketDataOrderRes.OrderInfo.ExecutionType = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.ExecutionType
			createBasketDataOrderRes.OrderInfo.DisclosedQuantity = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.DisclosedQuantity
			createBasketDataOrderRes.OrderInfo.RejectionReason = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.RejectionReason
			createBasketDataOrderRes.OrderInfo.StopLossValue = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.StopLossValue
			createBasketDataOrderRes.OrderInfo.Device = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.Device
			createBasketDataOrderRes.OrderInfo.Product = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.Product
			createBasketDataOrderRes.OrderInfo.SlOrderPrice = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.SlOrderPrice
			createBasketDataOrderRes.OrderInfo.FilledQuantity = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.FilledQuantity
			createBasketDataOrderRes.OrderInfo.ExchangeOrderID = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.ExchangeOrderID
			createBasketDataOrderRes.OrderInfo.Deposit = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.Deposit
			createBasketDataOrderRes.OrderInfo.AveragePrice = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.AveragePrice
			createBasketDataOrderRes.OrderInfo.SpreadToken = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.SpreadToken
			createBasketDataOrderRes.OrderInfo.OrderStatus = tlCreateBasketRes.Data[i].Orders[j].OrderInfo.OrderStatus
			createBasketDataOrderResAll = append(createBasketDataOrderResAll, createBasketDataOrderRes)
		}
		createBasketDataRes.Orders = createBasketDataOrderResAll

		createBasketRes = createBasketDataRes
	}

	loggerconfig.Info("CreateBasketRes tl resp=", helpers.LogStructAsJSON(createBasketRes), " StatusCode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)

	apiRes.Data = createBasketRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj BasketOrderObj) FetchBasket(req models.FetchBasketReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + BasketURL + "?login_id=" + url.QueryEscape(req.LoginID)

	//make payload
	payload := new(bytes.Buffer)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "FetchBasket", duration, req.LoginID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " FetchBasketRes call api error", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("FetchBasketRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlFetchBasketRes := TradeLabBasketRes{}
	json.Unmarshal([]byte(string(body)), &tlFetchBasketRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " FetchBasketRes tl status not ok =", tlFetchBasketRes.Message, " StatusCode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		apiRes.Message = tlFetchBasketRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var fetchBasketRes models.BasketRes

	fetchBasketDataResAll := make([]models.BasketDataRes, 0)

	for i := 0; i < len(tlFetchBasketRes.Data); i++ {
		var fetchBasketDataRes models.BasketDataRes
		fetchBasketDataRes.BasketID = tlFetchBasketRes.Data[i].BasketID
		fetchBasketDataRes.BasketType = tlFetchBasketRes.Data[i].BasketType
		fetchBasketDataRes.IsExecuted = tlFetchBasketRes.Data[i].IsExecuted
		fetchBasketDataRes.LoginID = tlFetchBasketRes.Data[i].LoginID
		fetchBasketDataRes.Name = tlFetchBasketRes.Data[i].Name
		fetchBasketDataRes.OrderType = tlFetchBasketRes.Data[i].OrderType
		fetchBasketDataRes.ProductType = tlFetchBasketRes.Data[i].ProductType
		fetchBasketDataRes.SipEligible = tlFetchBasketRes.Data[i].SipEligible
		fetchBasketDataRes.SipEnabled = tlFetchBasketRes.Data[i].SipEnabled

		fetchBasketDataOrderResAll := make([]models.BasketDataOrderRes, 0)
		for j := len(tlFetchBasketRes.Data[i].Orders) - 1; j >= 0; j-- {
			var fetchBasketDataOrderRes models.BasketDataOrderRes
			fetchBasketDataOrderRes.OrderID = tlFetchBasketRes.Data[i].Orders[j].OrderID
			fetchBasketDataOrderRes.OrderInfo.TriggerPrice = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.TriggerPrice
			fetchBasketDataOrderRes.OrderInfo.UnderlyingToken = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.UnderlyingToken
			fetchBasketDataOrderRes.OrderInfo.Series = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.Series
			fetchBasketDataOrderRes.OrderInfo.UserOrderID = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.UserOrderID
			fetchBasketDataOrderRes.OrderInfo.Exchange = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.Exchange
			fetchBasketDataOrderRes.OrderInfo.SquareOff = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.SquareOff
			fetchBasketDataOrderRes.OrderInfo.Mode = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.Mode
			fetchBasketDataOrderRes.OrderInfo.RemainingQuantity = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.RemainingQuantity
			fetchBasketDataOrderRes.OrderInfo.AverageTradePrice = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.AverageTradePrice
			fetchBasketDataOrderRes.OrderInfo.TradePrice = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.TradePrice
			fetchBasketDataOrderRes.OrderInfo.OrderTag = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.OrderTag
			fetchBasketDataOrderRes.OrderInfo.OrderStatusInfo = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.OrderStatusInfo
			fetchBasketDataOrderRes.OrderInfo.OrderSide = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.OrderSide
			fetchBasketDataOrderRes.OrderInfo.SquareOffValue = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.SquareOffValue
			fetchBasketDataOrderRes.OrderInfo.ContractDescription = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.ContractDescription
			fetchBasketDataOrderRes.OrderInfo.Segment = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.Segment
			fetchBasketDataOrderRes.OrderInfo.ClientID = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.ClientID
			fetchBasketDataOrderRes.OrderInfo.TradingSymbol = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.TradingSymbol
			fetchBasketDataOrderRes.OrderInfo.RejectionCode = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.RejectionCode
			fetchBasketDataOrderRes.OrderInfo.LotSize = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.LotSize
			fetchBasketDataOrderRes.OrderInfo.Quantity = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.Quantity
			fetchBasketDataOrderRes.OrderInfo.LastActivityReference = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.LastActivityReference
			fetchBasketDataOrderRes.OrderInfo.NnfID = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.NnfID
			fetchBasketDataOrderRes.OrderInfo.ProCli = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.ProCli
			fetchBasketDataOrderRes.OrderInfo.Price = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.Price
			fetchBasketDataOrderRes.OrderInfo.OrderType = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.OrderType
			fetchBasketDataOrderRes.OrderInfo.Validity = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.Validity
			fetchBasketDataOrderRes.OrderInfo.TargetPriceType = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.TargetPriceType
			fetchBasketDataOrderRes.OrderInfo.InstrumentToken = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.InstrumentToken
			fetchBasketDataOrderRes.OrderInfo.SlTriggerPrice = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.SlTriggerPrice
			fetchBasketDataOrderRes.OrderInfo.IsTrailing = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.IsTrailing
			fetchBasketDataOrderRes.OrderInfo.SlOrderQuantity = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.SlOrderQuantity
			fetchBasketDataOrderRes.OrderInfo.OrderEntryTime = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.OrderEntryTime
			fetchBasketDataOrderRes.OrderInfo.ExchangeTime = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.ExchangeTime
			fetchBasketDataOrderRes.OrderInfo.LegOrderIndicator = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.LegOrderIndicator
			fetchBasketDataOrderRes.OrderInfo.TrailingStopLoss = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.TrailingStopLoss
			fetchBasketDataOrderRes.OrderInfo.LoginID = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.LoginID
			fetchBasketDataOrderRes.OrderInfo.OmsOrderID = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.OmsOrderID
			fetchBasketDataOrderRes.OrderInfo.MarketProtectionPercentage = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.MarketProtectionPercentage
			fetchBasketDataOrderRes.OrderInfo.ExecutionType = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.ExecutionType
			fetchBasketDataOrderRes.OrderInfo.DisclosedQuantity = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.DisclosedQuantity
			fetchBasketDataOrderRes.OrderInfo.RejectionReason = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.RejectionReason
			fetchBasketDataOrderRes.OrderInfo.StopLossValue = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.StopLossValue
			fetchBasketDataOrderRes.OrderInfo.Device = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.Device
			fetchBasketDataOrderRes.OrderInfo.Product = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.Product
			fetchBasketDataOrderRes.OrderInfo.SlOrderPrice = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.SlOrderPrice
			fetchBasketDataOrderRes.OrderInfo.FilledQuantity = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.FilledQuantity
			fetchBasketDataOrderRes.OrderInfo.ExchangeOrderID = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.ExchangeOrderID
			fetchBasketDataOrderRes.OrderInfo.Deposit = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.Deposit
			fetchBasketDataOrderRes.OrderInfo.AveragePrice = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.AveragePrice
			fetchBasketDataOrderRes.OrderInfo.SpreadToken = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.SpreadToken
			fetchBasketDataOrderRes.OrderInfo.OrderStatus = tlFetchBasketRes.Data[i].Orders[j].OrderInfo.OrderStatus
			fetchBasketDataOrderResAll = append(fetchBasketDataOrderResAll, fetchBasketDataOrderRes)
		}
		fetchBasketDataRes.Orders = fetchBasketDataOrderResAll

		fetchBasketDataResAll = append(fetchBasketDataResAll, fetchBasketDataRes)
	}
	fetchBasketRes.Data = fetchBasketDataResAll

	loggerconfig.Info("FetchBasketRes tl resp=", helpers.LogStructAsJSON(fetchBasketRes), " StatusCode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)

	apiRes.Data = fetchBasketRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj BasketOrderObj) DeleteBasket(req models.DeleteBasketReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + BasketURL

	var tlDeleteBasketReq TradeLabDeleteBasketReq
	tlDeleteBasketReq.BasketID = req.BasketID
	tlDeleteBasketReq.Name = req.Name
	tlDeleteBasketReq.SipCount = req.SipCount

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlDeleteBasketReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodDelete, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "DeleteBasket", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		// logObj.Printf("BOOrderRes call api error =", err)
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " tlDeleteBasketRes call api error", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("DeleteBasketRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	tlDeleteBasketRes := TradeLabBasketRes{}
	json.Unmarshal([]byte(string(body)), &tlDeleteBasketRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " DeleteBasketRes tl status not ok =", tlDeleteBasketRes.Message, " StatusCode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		apiRes.Message = tlDeleteBasketRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("DeleteBasketRes tl status ok, clientID: ", reqH.ClientId, "requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)

	apiRes.Data = nil
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj BasketOrderObj) AddBasketInstrument(req models.AddBasketInstrumentReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + BasketInstrumentURL

	var tlAddBasketInstrumentReq TradeLabAddBasketInstrumentReq
	tlAddBasketInstrumentReq.BasketID = req.BasketID
	tlAddBasketInstrumentReq.Name = req.Name
	tlAddBasketInstrumentReq.OrderInfo.Exchange = req.OrderInfo.Exchange
	tlAddBasketInstrumentReq.OrderInfo.InstrumentToken = strconv.Itoa(req.OrderInfo.InstrumentToken)
	tlAddBasketInstrumentReq.OrderInfo.ClientID = req.OrderInfo.ClientID
	tlAddBasketInstrumentReq.OrderInfo.OrderType = req.OrderInfo.OrderType
	tlAddBasketInstrumentReq.OrderInfo.Price = req.OrderInfo.Price
	tlAddBasketInstrumentReq.OrderInfo.Quantity = req.OrderInfo.Quantity
	tlAddBasketInstrumentReq.OrderInfo.DisclosedQuantity = req.OrderInfo.DisclosedQuantity
	tlAddBasketInstrumentReq.OrderInfo.Validity = req.OrderInfo.Validity
	tlAddBasketInstrumentReq.OrderInfo.Product = req.OrderInfo.Product
	tlAddBasketInstrumentReq.OrderInfo.TradingSymbol = req.OrderInfo.TradingSymbol
	tlAddBasketInstrumentReq.OrderInfo.OrderSide = req.OrderInfo.OrderSide
	tlAddBasketInstrumentReq.OrderInfo.UserOrderID = req.OrderInfo.UserOrderID
	tlAddBasketInstrumentReq.OrderInfo.UnderlyingToken = req.OrderInfo.UnderlyingToken
	tlAddBasketInstrumentReq.OrderInfo.Series = req.OrderInfo.Series
	tlAddBasketInstrumentReq.OrderInfo.Device = reqH.DeviceType
	tlAddBasketInstrumentReq.OrderInfo.TriggerPrice = req.OrderInfo.TriggerPrice
	tlAddBasketInstrumentReq.OrderInfo.ExecutionType = req.OrderInfo.ExecutionType

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlAddBasketInstrumentReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "AddBasketInstrument", duration, req.OrderInfo.ClientID, reqH.RequestId)
	if err != nil {
		// logObj.Printf("BOOrderRes call api error =", err)
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " tlAddBasketInstrumentReq call api error", "clientID: ", req.OrderInfo.ClientID, err, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("AddBasketInstrumentRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", req.OrderInfo.ClientID, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	tlAddBasketInstrumentRes := TradeLabBasketInstrumentRes{}
	json.Unmarshal([]byte(string(body)), &tlAddBasketInstrumentRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " AddBasketInstrumentRes tl status not ok = ", res.StatusCode, " statuscode: ", res.StatusCode, " uccId:", req.OrderInfo.ClientID, "message = ", tlAddBasketInstrumentRes.Message, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		apiRes.Message = tlAddBasketInstrumentRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var addBasketInstrumentRes models.BasketInstrumentRes

	var addBasketInstrumentDataRes models.BasketDataRes
	addBasketInstrumentDataRes.BasketID = tlAddBasketInstrumentRes.Data.BasketID
	addBasketInstrumentDataRes.BasketType = tlAddBasketInstrumentRes.Data.BasketType
	addBasketInstrumentDataRes.IsExecuted = tlAddBasketInstrumentRes.Data.IsExecuted
	addBasketInstrumentDataRes.LoginID = tlAddBasketInstrumentRes.Data.LoginID
	addBasketInstrumentDataRes.Name = tlAddBasketInstrumentRes.Data.Name
	addBasketInstrumentDataRes.OrderType = tlAddBasketInstrumentRes.Data.OrderType
	addBasketInstrumentDataRes.ProductType = tlAddBasketInstrumentRes.Data.ProductType
	addBasketInstrumentDataRes.SipEligible = tlAddBasketInstrumentRes.Data.SipEligible
	addBasketInstrumentDataRes.SipEnabled = tlAddBasketInstrumentRes.Data.SipEnabled

	addBasketInstrumentDataOrderResAll := make([]models.BasketDataOrderRes, 0)
	for i := 0; i < len(tlAddBasketInstrumentRes.Data.Orders); i++ {
		var addBasketInstrumentDataOrderRes models.BasketDataOrderRes
		addBasketInstrumentDataOrderRes.OrderID = tlAddBasketInstrumentRes.Data.Orders[i].OrderID
		addBasketInstrumentDataOrderRes.OrderInfo.TriggerPrice = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.TriggerPrice
		addBasketInstrumentDataOrderRes.OrderInfo.UnderlyingToken = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.UnderlyingToken
		addBasketInstrumentDataOrderRes.OrderInfo.Series = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.Series
		addBasketInstrumentDataOrderRes.OrderInfo.UserOrderID = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.UserOrderID
		addBasketInstrumentDataOrderRes.OrderInfo.Exchange = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.Exchange
		addBasketInstrumentDataOrderRes.OrderInfo.SquareOff = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.SquareOff
		addBasketInstrumentDataOrderRes.OrderInfo.Mode = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.Mode
		addBasketInstrumentDataOrderRes.OrderInfo.RemainingQuantity = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.RemainingQuantity
		addBasketInstrumentDataOrderRes.OrderInfo.AverageTradePrice = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.AverageTradePrice
		addBasketInstrumentDataOrderRes.OrderInfo.TradePrice = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.TradePrice
		addBasketInstrumentDataOrderRes.OrderInfo.OrderTag = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.OrderTag
		addBasketInstrumentDataOrderRes.OrderInfo.OrderStatusInfo = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.OrderStatusInfo
		addBasketInstrumentDataOrderRes.OrderInfo.OrderSide = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.OrderSide
		addBasketInstrumentDataOrderRes.OrderInfo.SquareOffValue = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.SquareOffValue
		addBasketInstrumentDataOrderRes.OrderInfo.ContractDescription = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.ContractDescription
		addBasketInstrumentDataOrderRes.OrderInfo.Segment = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.Segment
		addBasketInstrumentDataOrderRes.OrderInfo.ClientID = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.ClientID
		addBasketInstrumentDataOrderRes.OrderInfo.TradingSymbol = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.TradingSymbol
		addBasketInstrumentDataOrderRes.OrderInfo.RejectionCode = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.RejectionCode
		addBasketInstrumentDataOrderRes.OrderInfo.LotSize = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.LotSize
		addBasketInstrumentDataOrderRes.OrderInfo.Quantity = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.Quantity
		addBasketInstrumentDataOrderRes.OrderInfo.LastActivityReference = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.LastActivityReference
		addBasketInstrumentDataOrderRes.OrderInfo.NnfID = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.NnfID
		addBasketInstrumentDataOrderRes.OrderInfo.ProCli = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.ProCli
		addBasketInstrumentDataOrderRes.OrderInfo.Price = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.Price
		addBasketInstrumentDataOrderRes.OrderInfo.OrderType = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.OrderType
		addBasketInstrumentDataOrderRes.OrderInfo.Validity = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.Validity
		addBasketInstrumentDataOrderRes.OrderInfo.TargetPriceType = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.TargetPriceType
		addBasketInstrumentDataOrderRes.OrderInfo.InstrumentToken = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.InstrumentToken
		addBasketInstrumentDataOrderRes.OrderInfo.SlTriggerPrice = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.SlTriggerPrice
		addBasketInstrumentDataOrderRes.OrderInfo.IsTrailing = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.IsTrailing
		addBasketInstrumentDataOrderRes.OrderInfo.SlOrderQuantity = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.SlOrderQuantity
		addBasketInstrumentDataOrderRes.OrderInfo.OrderEntryTime = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.OrderEntryTime
		addBasketInstrumentDataOrderRes.OrderInfo.ExchangeTime = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.ExchangeTime
		addBasketInstrumentDataOrderRes.OrderInfo.LegOrderIndicator = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.LegOrderIndicator
		addBasketInstrumentDataOrderRes.OrderInfo.TrailingStopLoss = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.TrailingStopLoss
		addBasketInstrumentDataOrderRes.OrderInfo.LoginID = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.LoginID
		addBasketInstrumentDataOrderRes.OrderInfo.OmsOrderID = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.OmsOrderID
		addBasketInstrumentDataOrderRes.OrderInfo.MarketProtectionPercentage = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.MarketProtectionPercentage
		addBasketInstrumentDataOrderRes.OrderInfo.ExecutionType = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.ExecutionType
		addBasketInstrumentDataOrderRes.OrderInfo.DisclosedQuantity = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.DisclosedQuantity
		addBasketInstrumentDataOrderRes.OrderInfo.RejectionReason = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.RejectionReason
		addBasketInstrumentDataOrderRes.OrderInfo.StopLossValue = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.StopLossValue
		addBasketInstrumentDataOrderRes.OrderInfo.Device = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.Device
		addBasketInstrumentDataOrderRes.OrderInfo.Product = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.Product
		addBasketInstrumentDataOrderRes.OrderInfo.SlOrderPrice = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.SlOrderPrice
		addBasketInstrumentDataOrderRes.OrderInfo.FilledQuantity = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.FilledQuantity
		addBasketInstrumentDataOrderRes.OrderInfo.ExchangeOrderID = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.ExchangeOrderID
		addBasketInstrumentDataOrderRes.OrderInfo.Deposit = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.Deposit
		addBasketInstrumentDataOrderRes.OrderInfo.AveragePrice = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.AveragePrice
		addBasketInstrumentDataOrderRes.OrderInfo.SpreadToken = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.SpreadToken
		addBasketInstrumentDataOrderRes.OrderInfo.OrderStatus = tlAddBasketInstrumentRes.Data.Orders[i].OrderInfo.OrderStatus
		addBasketInstrumentDataOrderResAll = append(addBasketInstrumentDataOrderResAll, addBasketInstrumentDataOrderRes)
	}
	addBasketInstrumentDataRes.Orders = addBasketInstrumentDataOrderResAll

	addBasketInstrumentRes.Data = addBasketInstrumentDataRes

	loggerconfig.Info("AddBasketInstrumentRes tl resp=", helpers.LogStructAsJSON(addBasketInstrumentRes), " uccId:", req.OrderInfo.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)

	apiRes.Data = addBasketInstrumentRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj BasketOrderObj) EditBasketInstrument(req models.EditBasketInstrumentReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + BasketInstrumentURL

	var tlEditBasketInstrumentReq TradeLabEditBasketInstrumentReq
	tlEditBasketInstrumentReq.BasketID = req.BasketID
	tlEditBasketInstrumentReq.Name = req.Name
	tlEditBasketInstrumentReq.OrderID = req.OrderID
	tlEditBasketInstrumentReq.OrderInfo.Exchange = req.OrderInfo.Exchange
	tlEditBasketInstrumentReq.OrderInfo.InstrumentToken = req.OrderInfo.InstrumentToken
	tlEditBasketInstrumentReq.OrderInfo.ClientID = req.OrderInfo.ClientID
	tlEditBasketInstrumentReq.OrderInfo.OrderType = req.OrderInfo.OrderType
	tlEditBasketInstrumentReq.OrderInfo.Price = req.OrderInfo.Price
	tlEditBasketInstrumentReq.OrderInfo.Quantity = req.OrderInfo.Quantity
	tlEditBasketInstrumentReq.OrderInfo.DisclosedQuantity = req.OrderInfo.DisclosedQuantity
	tlEditBasketInstrumentReq.OrderInfo.Validity = req.OrderInfo.Validity
	tlEditBasketInstrumentReq.OrderInfo.Product = req.OrderInfo.Product
	tlEditBasketInstrumentReq.OrderInfo.TradingSymbol = req.OrderInfo.TradingSymbol
	tlEditBasketInstrumentReq.OrderInfo.OrderSide = req.OrderInfo.OrderSide
	tlEditBasketInstrumentReq.OrderInfo.UserOrderID = req.OrderInfo.UserOrderID
	tlEditBasketInstrumentReq.OrderInfo.UnderlyingToken = req.OrderInfo.UnderlyingToken
	tlEditBasketInstrumentReq.OrderInfo.Series = req.OrderInfo.Series
	tlEditBasketInstrumentReq.OrderInfo.OmsOrderID = req.OrderInfo.OmsOrderID
	tlEditBasketInstrumentReq.OrderInfo.ExchangeOrderID = req.OrderInfo.ExchangeOrderID
	tlEditBasketInstrumentReq.OrderInfo.TriggerPrice = req.OrderInfo.TriggerPrice
	tlEditBasketInstrumentReq.OrderInfo.ExecutionType = req.OrderInfo.ExecutionType

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlEditBasketInstrumentReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPut, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "EditBasketInstrument", duration, req.OrderInfo.ClientID, reqH.RequestId)
	if err != nil {
		// logObj.Printf("BOOrderRes call api error =", err)
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " tlEditBasketInstrumentReq call api error", err, "clientID: ", req.OrderInfo.ClientID, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("EditBasketInstrumentRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", req.OrderInfo.ClientID, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	tlEditBasketInstrumentRes := TradeLabBasketInstrumentRes{}
	json.Unmarshal([]byte(string(body)), &tlEditBasketInstrumentRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " EditBasketInstrumentRes tl status not ok =", tlEditBasketInstrumentRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.OrderInfo.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		apiRes.Message = tlEditBasketInstrumentRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var editBasketInstrumentRes models.BasketInstrumentRes

	var editBasketInstrumentDataRes models.BasketDataRes
	editBasketInstrumentDataRes.BasketID = tlEditBasketInstrumentRes.Data.BasketID
	editBasketInstrumentDataRes.BasketType = tlEditBasketInstrumentRes.Data.BasketType
	editBasketInstrumentDataRes.IsExecuted = tlEditBasketInstrumentRes.Data.IsExecuted
	editBasketInstrumentDataRes.LoginID = tlEditBasketInstrumentRes.Data.LoginID
	editBasketInstrumentDataRes.Name = tlEditBasketInstrumentRes.Data.Name
	editBasketInstrumentDataRes.OrderType = tlEditBasketInstrumentRes.Data.OrderType
	editBasketInstrumentDataRes.ProductType = tlEditBasketInstrumentRes.Data.ProductType
	editBasketInstrumentDataRes.SipEligible = tlEditBasketInstrumentRes.Data.SipEligible
	editBasketInstrumentDataRes.SipEnabled = tlEditBasketInstrumentRes.Data.SipEnabled

	editBasketInstrumentDataOrderResAll := make([]models.BasketDataOrderRes, 0)
	for i := 0; i < len(tlEditBasketInstrumentRes.Data.Orders); i++ {
		var editBasketInstrumentDataOrderRes models.BasketDataOrderRes
		editBasketInstrumentDataOrderRes.OrderID = tlEditBasketInstrumentRes.Data.Orders[i].OrderID
		editBasketInstrumentDataOrderRes.OrderInfo.TriggerPrice = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.TriggerPrice
		editBasketInstrumentDataOrderRes.OrderInfo.UnderlyingToken = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.UnderlyingToken
		editBasketInstrumentDataOrderRes.OrderInfo.Series = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.Series
		editBasketInstrumentDataOrderRes.OrderInfo.UserOrderID = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.UserOrderID
		editBasketInstrumentDataOrderRes.OrderInfo.Exchange = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.Exchange
		editBasketInstrumentDataOrderRes.OrderInfo.SquareOff = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.SquareOff
		editBasketInstrumentDataOrderRes.OrderInfo.Mode = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.Mode
		editBasketInstrumentDataOrderRes.OrderInfo.RemainingQuantity = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.RemainingQuantity
		editBasketInstrumentDataOrderRes.OrderInfo.AverageTradePrice = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.AverageTradePrice
		editBasketInstrumentDataOrderRes.OrderInfo.TradePrice = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.TradePrice
		editBasketInstrumentDataOrderRes.OrderInfo.OrderTag = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.OrderTag
		editBasketInstrumentDataOrderRes.OrderInfo.OrderStatusInfo = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.OrderStatusInfo
		editBasketInstrumentDataOrderRes.OrderInfo.OrderSide = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.OrderSide
		editBasketInstrumentDataOrderRes.OrderInfo.SquareOffValue = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.SquareOffValue
		editBasketInstrumentDataOrderRes.OrderInfo.ContractDescription = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.ContractDescription
		editBasketInstrumentDataOrderRes.OrderInfo.Segment = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.Segment
		editBasketInstrumentDataOrderRes.OrderInfo.ClientID = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.ClientID
		editBasketInstrumentDataOrderRes.OrderInfo.TradingSymbol = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.TradingSymbol
		editBasketInstrumentDataOrderRes.OrderInfo.RejectionCode = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.RejectionCode
		editBasketInstrumentDataOrderRes.OrderInfo.LotSize = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.LotSize
		editBasketInstrumentDataOrderRes.OrderInfo.Quantity = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.Quantity
		editBasketInstrumentDataOrderRes.OrderInfo.LastActivityReference = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.LastActivityReference
		editBasketInstrumentDataOrderRes.OrderInfo.NnfID = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.NnfID
		editBasketInstrumentDataOrderRes.OrderInfo.ProCli = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.ProCli
		editBasketInstrumentDataOrderRes.OrderInfo.Price = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.Price
		editBasketInstrumentDataOrderRes.OrderInfo.OrderType = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.OrderType
		editBasketInstrumentDataOrderRes.OrderInfo.Validity = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.Validity
		editBasketInstrumentDataOrderRes.OrderInfo.TargetPriceType = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.TargetPriceType
		editBasketInstrumentDataOrderRes.OrderInfo.InstrumentToken = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.InstrumentToken
		editBasketInstrumentDataOrderRes.OrderInfo.SlTriggerPrice = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.SlTriggerPrice
		editBasketInstrumentDataOrderRes.OrderInfo.IsTrailing = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.IsTrailing
		editBasketInstrumentDataOrderRes.OrderInfo.SlOrderQuantity = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.SlOrderQuantity
		editBasketInstrumentDataOrderRes.OrderInfo.OrderEntryTime = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.OrderEntryTime
		editBasketInstrumentDataOrderRes.OrderInfo.ExchangeTime = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.ExchangeTime
		editBasketInstrumentDataOrderRes.OrderInfo.LegOrderIndicator = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.LegOrderIndicator
		editBasketInstrumentDataOrderRes.OrderInfo.TrailingStopLoss = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.TrailingStopLoss
		editBasketInstrumentDataOrderRes.OrderInfo.LoginID = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.LoginID
		editBasketInstrumentDataOrderRes.OrderInfo.OmsOrderID = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.OmsOrderID
		editBasketInstrumentDataOrderRes.OrderInfo.MarketProtectionPercentage = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.MarketProtectionPercentage
		editBasketInstrumentDataOrderRes.OrderInfo.ExecutionType = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.ExecutionType
		editBasketInstrumentDataOrderRes.OrderInfo.DisclosedQuantity = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.DisclosedQuantity
		editBasketInstrumentDataOrderRes.OrderInfo.RejectionReason = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.RejectionReason
		editBasketInstrumentDataOrderRes.OrderInfo.StopLossValue = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.StopLossValue
		editBasketInstrumentDataOrderRes.OrderInfo.Device = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.Device
		editBasketInstrumentDataOrderRes.OrderInfo.Product = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.Product
		editBasketInstrumentDataOrderRes.OrderInfo.SlOrderPrice = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.SlOrderPrice
		editBasketInstrumentDataOrderRes.OrderInfo.FilledQuantity = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.FilledQuantity
		editBasketInstrumentDataOrderRes.OrderInfo.ExchangeOrderID = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.ExchangeOrderID
		editBasketInstrumentDataOrderRes.OrderInfo.Deposit = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.Deposit
		editBasketInstrumentDataOrderRes.OrderInfo.AveragePrice = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.AveragePrice
		editBasketInstrumentDataOrderRes.OrderInfo.SpreadToken = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.SpreadToken
		editBasketInstrumentDataOrderRes.OrderInfo.OrderStatus = tlEditBasketInstrumentRes.Data.Orders[i].OrderInfo.OrderStatus
		editBasketInstrumentDataOrderResAll = append(editBasketInstrumentDataOrderResAll, editBasketInstrumentDataOrderRes)
	}
	editBasketInstrumentDataRes.Orders = editBasketInstrumentDataOrderResAll

	editBasketInstrumentRes.Data = editBasketInstrumentDataRes

	loggerconfig.Info("EditBasketInstrumentRes tl resp=", helpers.LogStructAsJSON(editBasketInstrumentRes), " uccId:", req.OrderInfo.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)

	apiRes.Data = editBasketInstrumentRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj BasketOrderObj) DeleteBasketInstrument(req models.DeleteBasketInstrumentReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + BasketInstrumentURL

	var tlDeleteBasketReq TradeLabDeleteBasketInstrumentReq
	tlDeleteBasketReq.BasketID = req.BasketID
	tlDeleteBasketReq.OrderID = req.OrderID
	tlDeleteBasketReq.Name = req.Name

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlDeleteBasketReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodDelete, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "DeleteBasketInstrument", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		// logObj.Printf("BOOrderRes call api error =", err)
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " tlDeleteBasketInstrumentRes call api error", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("DeleteBasketInstrumentRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	tlDeleteBasketInstrumentRes := TradeLabBasketRes{}
	json.Unmarshal([]byte(string(body)), &tlDeleteBasketInstrumentRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " DeleteBasketInstrumentRes tl status not ok =", tlDeleteBasketInstrumentRes.Message, " statuscode: ", res.StatusCode, " StatusCode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		apiRes.Message = tlDeleteBasketInstrumentRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var deleteBasketInstrumentRes models.BasketRes

	deleteBasketInstrumentDataResAll := make([]models.BasketDataRes, 0)

	for i := 0; i < len(tlDeleteBasketInstrumentRes.Data); i++ {
		var deleteBasketInstrumentDataRes models.BasketDataRes
		deleteBasketInstrumentDataRes.BasketID = tlDeleteBasketInstrumentRes.Data[i].BasketID
		deleteBasketInstrumentDataRes.BasketType = tlDeleteBasketInstrumentRes.Data[i].BasketType
		deleteBasketInstrumentDataRes.IsExecuted = tlDeleteBasketInstrumentRes.Data[i].IsExecuted
		deleteBasketInstrumentDataRes.LoginID = tlDeleteBasketInstrumentRes.Data[i].LoginID
		deleteBasketInstrumentDataRes.Name = tlDeleteBasketInstrumentRes.Data[i].Name
		deleteBasketInstrumentDataRes.OrderType = tlDeleteBasketInstrumentRes.Data[i].OrderType
		deleteBasketInstrumentDataRes.ProductType = tlDeleteBasketInstrumentRes.Data[i].ProductType
		deleteBasketInstrumentDataRes.SipEligible = tlDeleteBasketInstrumentRes.Data[i].SipEligible
		deleteBasketInstrumentDataRes.SipEnabled = tlDeleteBasketInstrumentRes.Data[i].SipEnabled

		deleteBasketInstrumentDataOrderResAll := make([]models.BasketDataOrderRes, 0)
		for j := 0; j < len(tlDeleteBasketInstrumentRes.Data[i].Orders); j++ {
			var deleteBasketInstrumentDataOrderRes models.BasketDataOrderRes
			deleteBasketInstrumentDataOrderRes.OrderID = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderID
			deleteBasketInstrumentDataOrderRes.OrderInfo.TriggerPrice = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.TriggerPrice
			deleteBasketInstrumentDataOrderRes.OrderInfo.UnderlyingToken = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.UnderlyingToken
			deleteBasketInstrumentDataOrderRes.OrderInfo.Series = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.Series
			deleteBasketInstrumentDataOrderRes.OrderInfo.UserOrderID = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.UserOrderID
			deleteBasketInstrumentDataOrderRes.OrderInfo.Exchange = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.Exchange
			deleteBasketInstrumentDataOrderRes.OrderInfo.SquareOff = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.SquareOff
			deleteBasketInstrumentDataOrderRes.OrderInfo.Mode = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.Mode
			deleteBasketInstrumentDataOrderRes.OrderInfo.RemainingQuantity = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.RemainingQuantity
			deleteBasketInstrumentDataOrderRes.OrderInfo.AverageTradePrice = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.AverageTradePrice
			deleteBasketInstrumentDataOrderRes.OrderInfo.TradePrice = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.TradePrice
			deleteBasketInstrumentDataOrderRes.OrderInfo.OrderTag = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.OrderTag
			deleteBasketInstrumentDataOrderRes.OrderInfo.OrderStatusInfo = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.OrderStatusInfo
			deleteBasketInstrumentDataOrderRes.OrderInfo.OrderSide = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.OrderSide
			deleteBasketInstrumentDataOrderRes.OrderInfo.SquareOffValue = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.SquareOffValue
			deleteBasketInstrumentDataOrderRes.OrderInfo.ContractDescription = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.ContractDescription
			deleteBasketInstrumentDataOrderRes.OrderInfo.Segment = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.Segment
			deleteBasketInstrumentDataOrderRes.OrderInfo.ClientID = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.ClientID
			deleteBasketInstrumentDataOrderRes.OrderInfo.TradingSymbol = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.TradingSymbol
			deleteBasketInstrumentDataOrderRes.OrderInfo.RejectionCode = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.RejectionCode
			deleteBasketInstrumentDataOrderRes.OrderInfo.LotSize = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.LotSize
			deleteBasketInstrumentDataOrderRes.OrderInfo.Quantity = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.Quantity
			deleteBasketInstrumentDataOrderRes.OrderInfo.LastActivityReference = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.LastActivityReference
			deleteBasketInstrumentDataOrderRes.OrderInfo.NnfID = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.NnfID
			deleteBasketInstrumentDataOrderRes.OrderInfo.ProCli = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.ProCli
			deleteBasketInstrumentDataOrderRes.OrderInfo.Price = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.Price
			deleteBasketInstrumentDataOrderRes.OrderInfo.OrderType = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.OrderType
			deleteBasketInstrumentDataOrderRes.OrderInfo.Validity = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.Validity
			deleteBasketInstrumentDataOrderRes.OrderInfo.TargetPriceType = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.TargetPriceType
			deleteBasketInstrumentDataOrderRes.OrderInfo.InstrumentToken = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.InstrumentToken
			deleteBasketInstrumentDataOrderRes.OrderInfo.SlTriggerPrice = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.SlTriggerPrice
			deleteBasketInstrumentDataOrderRes.OrderInfo.IsTrailing = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.IsTrailing
			deleteBasketInstrumentDataOrderRes.OrderInfo.SlOrderQuantity = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.SlOrderQuantity
			deleteBasketInstrumentDataOrderRes.OrderInfo.OrderEntryTime = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.OrderEntryTime
			deleteBasketInstrumentDataOrderRes.OrderInfo.ExchangeTime = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.ExchangeTime
			deleteBasketInstrumentDataOrderRes.OrderInfo.LegOrderIndicator = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.LegOrderIndicator
			deleteBasketInstrumentDataOrderRes.OrderInfo.TrailingStopLoss = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.TrailingStopLoss
			deleteBasketInstrumentDataOrderRes.OrderInfo.LoginID = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.LoginID
			deleteBasketInstrumentDataOrderRes.OrderInfo.OmsOrderID = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.OmsOrderID
			deleteBasketInstrumentDataOrderRes.OrderInfo.MarketProtectionPercentage = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.MarketProtectionPercentage
			deleteBasketInstrumentDataOrderRes.OrderInfo.ExecutionType = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.ExecutionType
			deleteBasketInstrumentDataOrderRes.OrderInfo.DisclosedQuantity = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.DisclosedQuantity
			deleteBasketInstrumentDataOrderRes.OrderInfo.RejectionReason = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.RejectionReason
			deleteBasketInstrumentDataOrderRes.OrderInfo.StopLossValue = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.StopLossValue
			deleteBasketInstrumentDataOrderRes.OrderInfo.Device = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.Device
			deleteBasketInstrumentDataOrderRes.OrderInfo.Product = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.Product
			deleteBasketInstrumentDataOrderRes.OrderInfo.SlOrderPrice = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.SlOrderPrice
			deleteBasketInstrumentDataOrderRes.OrderInfo.FilledQuantity = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.FilledQuantity
			deleteBasketInstrumentDataOrderRes.OrderInfo.ExchangeOrderID = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.ExchangeOrderID
			deleteBasketInstrumentDataOrderRes.OrderInfo.Deposit = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.Deposit
			deleteBasketInstrumentDataOrderRes.OrderInfo.AveragePrice = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.AveragePrice
			deleteBasketInstrumentDataOrderRes.OrderInfo.SpreadToken = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.SpreadToken
			deleteBasketInstrumentDataOrderRes.OrderInfo.OrderStatus = tlDeleteBasketInstrumentRes.Data[i].Orders[j].OrderInfo.OrderStatus
			deleteBasketInstrumentDataOrderResAll = append(deleteBasketInstrumentDataOrderResAll, deleteBasketInstrumentDataOrderRes)
		}
		deleteBasketInstrumentDataRes.Orders = deleteBasketInstrumentDataOrderResAll

		deleteBasketInstrumentDataResAll = append(deleteBasketInstrumentDataResAll, deleteBasketInstrumentDataRes)
	}
	deleteBasketInstrumentRes.Data = deleteBasketInstrumentDataResAll

	loggerconfig.Info("DeleteBasketInstrumentRes tl resp=", helpers.LogStructAsJSON(deleteBasketInstrumentRes), " StatusCode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)

	apiRes.Data = deleteBasketInstrumentRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj BasketOrderObj) RenameBasket(req models.RenameBasketReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + BasketURL

	var tlRenameBasketReq TradeLabRenameBasketReq
	tlRenameBasketReq.BasketID = req.BasketID
	tlRenameBasketReq.Name = req.Name

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlRenameBasketReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPut, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "RenameBasket", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		// logObj.Printf("BOOrderRes call api error =", err)
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " tlRenameBasketReq call api error", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("RenameBasketRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	tlRenameBasketRes := TradeLabBasketInstrumentRes{}
	json.Unmarshal([]byte(string(body)), &tlRenameBasketRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " RenameBasketRes tl status not ok =", tlRenameBasketRes.Message, " StatusCode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		apiRes.Message = tlRenameBasketRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var renameBasketRes models.BasketInstrumentRes

	var renameBasketDataRes models.BasketDataRes
	renameBasketDataRes.BasketID = tlRenameBasketRes.Data.BasketID
	renameBasketDataRes.BasketType = tlRenameBasketRes.Data.BasketType
	renameBasketDataRes.IsExecuted = tlRenameBasketRes.Data.IsExecuted
	renameBasketDataRes.LoginID = tlRenameBasketRes.Data.LoginID
	renameBasketDataRes.Name = tlRenameBasketRes.Data.Name
	renameBasketDataRes.OrderType = tlRenameBasketRes.Data.OrderType
	renameBasketDataRes.ProductType = tlRenameBasketRes.Data.ProductType
	renameBasketDataRes.SipEligible = tlRenameBasketRes.Data.SipEligible
	renameBasketDataRes.SipEnabled = tlRenameBasketRes.Data.SipEnabled

	renameBasketDataOrderResAll := make([]models.BasketDataOrderRes, 0)
	for i := 0; i < len(tlRenameBasketRes.Data.Orders); i++ {
		var renameBasketDataOrderRes models.BasketDataOrderRes
		renameBasketDataOrderRes.OrderID = tlRenameBasketRes.Data.Orders[i].OrderID
		renameBasketDataOrderRes.OrderInfo.TriggerPrice = tlRenameBasketRes.Data.Orders[i].OrderInfo.TriggerPrice
		renameBasketDataOrderRes.OrderInfo.UnderlyingToken = tlRenameBasketRes.Data.Orders[i].OrderInfo.UnderlyingToken
		renameBasketDataOrderRes.OrderInfo.Series = tlRenameBasketRes.Data.Orders[i].OrderInfo.Series
		renameBasketDataOrderRes.OrderInfo.UserOrderID = tlRenameBasketRes.Data.Orders[i].OrderInfo.UserOrderID
		renameBasketDataOrderRes.OrderInfo.Exchange = tlRenameBasketRes.Data.Orders[i].OrderInfo.Exchange
		renameBasketDataOrderRes.OrderInfo.SquareOff = tlRenameBasketRes.Data.Orders[i].OrderInfo.SquareOff
		renameBasketDataOrderRes.OrderInfo.Mode = tlRenameBasketRes.Data.Orders[i].OrderInfo.Mode
		renameBasketDataOrderRes.OrderInfo.RemainingQuantity = tlRenameBasketRes.Data.Orders[i].OrderInfo.RemainingQuantity
		renameBasketDataOrderRes.OrderInfo.AverageTradePrice = tlRenameBasketRes.Data.Orders[i].OrderInfo.AverageTradePrice
		renameBasketDataOrderRes.OrderInfo.TradePrice = tlRenameBasketRes.Data.Orders[i].OrderInfo.TradePrice
		renameBasketDataOrderRes.OrderInfo.OrderTag = tlRenameBasketRes.Data.Orders[i].OrderInfo.OrderTag
		renameBasketDataOrderRes.OrderInfo.OrderStatusInfo = tlRenameBasketRes.Data.Orders[i].OrderInfo.OrderStatusInfo
		renameBasketDataOrderRes.OrderInfo.OrderSide = tlRenameBasketRes.Data.Orders[i].OrderInfo.OrderSide
		renameBasketDataOrderRes.OrderInfo.SquareOffValue = tlRenameBasketRes.Data.Orders[i].OrderInfo.SquareOffValue
		renameBasketDataOrderRes.OrderInfo.ContractDescription = tlRenameBasketRes.Data.Orders[i].OrderInfo.ContractDescription
		renameBasketDataOrderRes.OrderInfo.Segment = tlRenameBasketRes.Data.Orders[i].OrderInfo.Segment
		renameBasketDataOrderRes.OrderInfo.ClientID = tlRenameBasketRes.Data.Orders[i].OrderInfo.ClientID
		renameBasketDataOrderRes.OrderInfo.TradingSymbol = tlRenameBasketRes.Data.Orders[i].OrderInfo.TradingSymbol
		renameBasketDataOrderRes.OrderInfo.RejectionCode = tlRenameBasketRes.Data.Orders[i].OrderInfo.RejectionCode
		renameBasketDataOrderRes.OrderInfo.LotSize = tlRenameBasketRes.Data.Orders[i].OrderInfo.LotSize
		renameBasketDataOrderRes.OrderInfo.Quantity = tlRenameBasketRes.Data.Orders[i].OrderInfo.Quantity
		renameBasketDataOrderRes.OrderInfo.LastActivityReference = tlRenameBasketRes.Data.Orders[i].OrderInfo.LastActivityReference
		renameBasketDataOrderRes.OrderInfo.NnfID = tlRenameBasketRes.Data.Orders[i].OrderInfo.NnfID
		renameBasketDataOrderRes.OrderInfo.ProCli = tlRenameBasketRes.Data.Orders[i].OrderInfo.ProCli
		renameBasketDataOrderRes.OrderInfo.Price = tlRenameBasketRes.Data.Orders[i].OrderInfo.Price
		renameBasketDataOrderRes.OrderInfo.OrderType = tlRenameBasketRes.Data.Orders[i].OrderInfo.OrderType
		renameBasketDataOrderRes.OrderInfo.Validity = tlRenameBasketRes.Data.Orders[i].OrderInfo.Validity
		renameBasketDataOrderRes.OrderInfo.TargetPriceType = tlRenameBasketRes.Data.Orders[i].OrderInfo.TargetPriceType
		renameBasketDataOrderRes.OrderInfo.InstrumentToken = tlRenameBasketRes.Data.Orders[i].OrderInfo.InstrumentToken
		renameBasketDataOrderRes.OrderInfo.SlTriggerPrice = tlRenameBasketRes.Data.Orders[i].OrderInfo.SlTriggerPrice
		renameBasketDataOrderRes.OrderInfo.IsTrailing = tlRenameBasketRes.Data.Orders[i].OrderInfo.IsTrailing
		renameBasketDataOrderRes.OrderInfo.SlOrderQuantity = tlRenameBasketRes.Data.Orders[i].OrderInfo.SlOrderQuantity
		renameBasketDataOrderRes.OrderInfo.OrderEntryTime = tlRenameBasketRes.Data.Orders[i].OrderInfo.OrderEntryTime
		renameBasketDataOrderRes.OrderInfo.ExchangeTime = tlRenameBasketRes.Data.Orders[i].OrderInfo.ExchangeTime
		renameBasketDataOrderRes.OrderInfo.LegOrderIndicator = tlRenameBasketRes.Data.Orders[i].OrderInfo.LegOrderIndicator
		renameBasketDataOrderRes.OrderInfo.TrailingStopLoss = tlRenameBasketRes.Data.Orders[i].OrderInfo.TrailingStopLoss
		renameBasketDataOrderRes.OrderInfo.LoginID = tlRenameBasketRes.Data.Orders[i].OrderInfo.LoginID
		renameBasketDataOrderRes.OrderInfo.OmsOrderID = tlRenameBasketRes.Data.Orders[i].OrderInfo.OmsOrderID
		renameBasketDataOrderRes.OrderInfo.MarketProtectionPercentage = tlRenameBasketRes.Data.Orders[i].OrderInfo.MarketProtectionPercentage
		renameBasketDataOrderRes.OrderInfo.ExecutionType = tlRenameBasketRes.Data.Orders[i].OrderInfo.ExecutionType
		renameBasketDataOrderRes.OrderInfo.DisclosedQuantity = tlRenameBasketRes.Data.Orders[i].OrderInfo.DisclosedQuantity
		renameBasketDataOrderRes.OrderInfo.RejectionReason = tlRenameBasketRes.Data.Orders[i].OrderInfo.RejectionReason
		renameBasketDataOrderRes.OrderInfo.StopLossValue = tlRenameBasketRes.Data.Orders[i].OrderInfo.StopLossValue
		renameBasketDataOrderRes.OrderInfo.Device = tlRenameBasketRes.Data.Orders[i].OrderInfo.Device
		renameBasketDataOrderRes.OrderInfo.Product = tlRenameBasketRes.Data.Orders[i].OrderInfo.Product
		renameBasketDataOrderRes.OrderInfo.SlOrderPrice = tlRenameBasketRes.Data.Orders[i].OrderInfo.SlOrderPrice
		renameBasketDataOrderRes.OrderInfo.FilledQuantity = tlRenameBasketRes.Data.Orders[i].OrderInfo.FilledQuantity
		renameBasketDataOrderRes.OrderInfo.ExchangeOrderID = tlRenameBasketRes.Data.Orders[i].OrderInfo.ExchangeOrderID
		renameBasketDataOrderRes.OrderInfo.Deposit = tlRenameBasketRes.Data.Orders[i].OrderInfo.Deposit
		renameBasketDataOrderRes.OrderInfo.AveragePrice = tlRenameBasketRes.Data.Orders[i].OrderInfo.AveragePrice
		renameBasketDataOrderRes.OrderInfo.SpreadToken = tlRenameBasketRes.Data.Orders[i].OrderInfo.SpreadToken
		renameBasketDataOrderRes.OrderInfo.OrderStatus = tlRenameBasketRes.Data.Orders[i].OrderInfo.OrderStatus
		renameBasketDataOrderResAll = append(renameBasketDataOrderResAll, renameBasketDataOrderRes)
	}
	renameBasketDataRes.Orders = renameBasketDataOrderResAll

	renameBasketRes.Data = renameBasketDataRes

	loggerconfig.Info("RenameBasketRes tl resp=", helpers.LogStructAsJSON(renameBasketRes), " StatusCode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)

	apiRes.Data = renameBasketRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj BasketOrderObj) ExecuteBasket(req models.ExecuteBasketReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + CONDATIONALORDERSURL

	var tlExecuteBasketReq TradeLabExecuteBasketReq
	tlExecuteBasketReq.BasketID = req.BasketID
	tlExecuteBasketReq.Name = req.Name
	tlExecuteBasketReq.ExecutionType = req.ExecutionType
	tlExecuteBasketReq.SquareOff = req.SquareOff
	tlExecuteBasketReq.ClientID = req.ClientID
	tlExecuteBasketReq.ExecutionState = req.ExecutionState

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlExecuteBasketReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ExecuteBasket", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " ExecuteBasketRes call api error", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("ExecuteBasketRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlExecuteBasketRes := TradeLabExecuteBasketRes{}
	json.Unmarshal([]byte(string(body)), &tlExecuteBasketRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " ExecuteBasketRes tl status not ok =", tlExecuteBasketRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		apiRes.Message = tlExecuteBasketRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var executeBasketRes models.ExecuteBasketRes
	executeBasketRes.Data.BasketID = tlExecuteBasketRes.Data.Data.BasketID
	executeBasketRes.Data.Message = tlExecuteBasketRes.Data.Data.Message

	loggerconfig.Info("ExecuteBasketRes tl resp=", helpers.LogStructAsJSON(executeBasketRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)

	apiRes.Data = executeBasketRes.Data
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj BasketOrderObj) UpdateBasketExecutionState(req models.UpdateBasketExecutionStateReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + UPDATEBASKETEXECUTIONSTATE

	var tlUpdateBasketExecutionStateReq TradeLabUpdateBasketExecutionStateReq
	tlUpdateBasketExecutionStateReq.BasketID = req.BasketID
	tlUpdateBasketExecutionStateReq.Name = req.Name
	tlUpdateBasketExecutionStateReq.ClientID = req.ClientID
	tlUpdateBasketExecutionStateReq.ExecutionType = req.ExecutionType
	tlUpdateBasketExecutionStateReq.SquareOff = req.SquareOff
	tlUpdateBasketExecutionStateReq.ExecutionState = req.ExecutionState

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlUpdateBasketExecutionStateReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "UpdateBasketExecutionState", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " UpdateBasketExecutionStateRes call api error", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("UpdateBasketExecutionStateRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlUpdateBasketExecutionStateRes := TradeLabBasketRes{}
	json.Unmarshal([]byte(string(body)), &tlUpdateBasketExecutionStateRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " UpdateBasketExecutionStateRes tl status not ok =", tlUpdateBasketExecutionStateRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)
		apiRes.Message = tlUpdateBasketExecutionStateRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var updateBasketExecutionStateRes models.BasketRes

	updateBasketExecutionStateResAll := make([]models.BasketDataRes, 0)

	for i := 0; i < len(tlUpdateBasketExecutionStateRes.Data); i++ {
		var updateBasketExecutionStateDataRes models.BasketDataRes
		updateBasketExecutionStateDataRes.BasketID = tlUpdateBasketExecutionStateRes.Data[i].BasketID
		updateBasketExecutionStateDataRes.BasketType = tlUpdateBasketExecutionStateRes.Data[i].BasketType
		updateBasketExecutionStateDataRes.IsExecuted = tlUpdateBasketExecutionStateRes.Data[i].IsExecuted
		updateBasketExecutionStateDataRes.LoginID = tlUpdateBasketExecutionStateRes.Data[i].LoginID
		updateBasketExecutionStateDataRes.Name = tlUpdateBasketExecutionStateRes.Data[i].Name
		updateBasketExecutionStateDataRes.OrderType = tlUpdateBasketExecutionStateRes.Data[i].OrderType
		updateBasketExecutionStateDataRes.ProductType = tlUpdateBasketExecutionStateRes.Data[i].ProductType
		updateBasketExecutionStateDataRes.SipEligible = tlUpdateBasketExecutionStateRes.Data[i].SipEligible
		updateBasketExecutionStateDataRes.SipEnabled = tlUpdateBasketExecutionStateRes.Data[i].SipEnabled

		updateBasketExecutionStateDataOrderResAll := make([]models.BasketDataOrderRes, 0)
		for j := 0; j < len(tlUpdateBasketExecutionStateRes.Data[i].Orders); j++ {
			var updateBasketExecutionStateDataOrderRes models.BasketDataOrderRes
			updateBasketExecutionStateDataOrderRes.OrderID = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderID
			updateBasketExecutionStateDataOrderRes.OrderInfo.TriggerPrice = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.TriggerPrice
			updateBasketExecutionStateDataOrderRes.OrderInfo.UnderlyingToken = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.UnderlyingToken
			updateBasketExecutionStateDataOrderRes.OrderInfo.Series = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.Series
			updateBasketExecutionStateDataOrderRes.OrderInfo.UserOrderID = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.UserOrderID
			updateBasketExecutionStateDataOrderRes.OrderInfo.Exchange = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.Exchange
			updateBasketExecutionStateDataOrderRes.OrderInfo.SquareOff = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.SquareOff
			updateBasketExecutionStateDataOrderRes.OrderInfo.Mode = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.Mode
			updateBasketExecutionStateDataOrderRes.OrderInfo.RemainingQuantity = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.RemainingQuantity
			updateBasketExecutionStateDataOrderRes.OrderInfo.AverageTradePrice = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.AverageTradePrice
			updateBasketExecutionStateDataOrderRes.OrderInfo.TradePrice = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.TradePrice
			updateBasketExecutionStateDataOrderRes.OrderInfo.OrderTag = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.OrderTag
			updateBasketExecutionStateDataOrderRes.OrderInfo.OrderStatusInfo = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.OrderStatusInfo
			updateBasketExecutionStateDataOrderRes.OrderInfo.OrderSide = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.OrderSide
			updateBasketExecutionStateDataOrderRes.OrderInfo.SquareOffValue = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.SquareOffValue
			updateBasketExecutionStateDataOrderRes.OrderInfo.ContractDescription = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.ContractDescription
			updateBasketExecutionStateDataOrderRes.OrderInfo.Segment = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.Segment
			updateBasketExecutionStateDataOrderRes.OrderInfo.ClientID = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.ClientID
			updateBasketExecutionStateDataOrderRes.OrderInfo.TradingSymbol = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.TradingSymbol
			updateBasketExecutionStateDataOrderRes.OrderInfo.RejectionCode = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.RejectionCode
			updateBasketExecutionStateDataOrderRes.OrderInfo.LotSize = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.LotSize
			updateBasketExecutionStateDataOrderRes.OrderInfo.Quantity = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.Quantity
			updateBasketExecutionStateDataOrderRes.OrderInfo.LastActivityReference = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.LastActivityReference
			updateBasketExecutionStateDataOrderRes.OrderInfo.NnfID = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.NnfID
			updateBasketExecutionStateDataOrderRes.OrderInfo.ProCli = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.ProCli
			updateBasketExecutionStateDataOrderRes.OrderInfo.Price = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.Price
			updateBasketExecutionStateDataOrderRes.OrderInfo.OrderType = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.OrderType
			updateBasketExecutionStateDataOrderRes.OrderInfo.Validity = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.Validity
			updateBasketExecutionStateDataOrderRes.OrderInfo.TargetPriceType = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.TargetPriceType
			updateBasketExecutionStateDataOrderRes.OrderInfo.InstrumentToken = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.InstrumentToken
			updateBasketExecutionStateDataOrderRes.OrderInfo.SlTriggerPrice = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.SlTriggerPrice
			updateBasketExecutionStateDataOrderRes.OrderInfo.IsTrailing = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.IsTrailing
			updateBasketExecutionStateDataOrderRes.OrderInfo.SlOrderQuantity = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.SlOrderQuantity
			updateBasketExecutionStateDataOrderRes.OrderInfo.OrderEntryTime = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.OrderEntryTime
			updateBasketExecutionStateDataOrderRes.OrderInfo.ExchangeTime = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.ExchangeTime
			updateBasketExecutionStateDataOrderRes.OrderInfo.LegOrderIndicator = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.LegOrderIndicator
			updateBasketExecutionStateDataOrderRes.OrderInfo.TrailingStopLoss = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.TrailingStopLoss
			updateBasketExecutionStateDataOrderRes.OrderInfo.LoginID = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.LoginID
			updateBasketExecutionStateDataOrderRes.OrderInfo.OmsOrderID = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.OmsOrderID
			updateBasketExecutionStateDataOrderRes.OrderInfo.MarketProtectionPercentage = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.MarketProtectionPercentage
			updateBasketExecutionStateDataOrderRes.OrderInfo.ExecutionType = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.ExecutionType
			updateBasketExecutionStateDataOrderRes.OrderInfo.DisclosedQuantity = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.DisclosedQuantity
			updateBasketExecutionStateDataOrderRes.OrderInfo.RejectionReason = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.RejectionReason
			updateBasketExecutionStateDataOrderRes.OrderInfo.StopLossValue = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.StopLossValue
			updateBasketExecutionStateDataOrderRes.OrderInfo.Device = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.Device
			updateBasketExecutionStateDataOrderRes.OrderInfo.Product = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.Product
			updateBasketExecutionStateDataOrderRes.OrderInfo.SlOrderPrice = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.SlOrderPrice
			updateBasketExecutionStateDataOrderRes.OrderInfo.FilledQuantity = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.FilledQuantity
			updateBasketExecutionStateDataOrderRes.OrderInfo.ExchangeOrderID = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.ExchangeOrderID
			updateBasketExecutionStateDataOrderRes.OrderInfo.Deposit = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.Deposit
			updateBasketExecutionStateDataOrderRes.OrderInfo.AveragePrice = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.AveragePrice
			updateBasketExecutionStateDataOrderRes.OrderInfo.SpreadToken = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.SpreadToken
			updateBasketExecutionStateDataOrderRes.OrderInfo.OrderStatus = tlUpdateBasketExecutionStateRes.Data[i].Orders[j].OrderInfo.OrderStatus
			updateBasketExecutionStateDataOrderResAll = append(updateBasketExecutionStateDataOrderResAll, updateBasketExecutionStateDataOrderRes)
		}
		updateBasketExecutionStateDataRes.Orders = updateBasketExecutionStateDataOrderResAll

		updateBasketExecutionStateResAll = append(updateBasketExecutionStateResAll, updateBasketExecutionStateDataRes)
	}
	updateBasketExecutionStateRes.Data = updateBasketExecutionStateResAll

	loggerconfig.Info("FetchBasketRes tl resp=", helpers.LogStructAsJSON(updateBasketExecutionStateRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion: ", reqH.ClientVersion)

	apiRes.Data = updateBasketExecutionStateRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}
