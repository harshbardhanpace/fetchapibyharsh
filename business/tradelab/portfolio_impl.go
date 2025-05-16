package tradelab

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
)

type PortfolioObj struct {
	tradeLabURL string
}

var objPortfolio PortfolioObj

func InitPortfolio() PortfolioObj {
	defer models.HandlePanic()

	portfolioObj := PortfolioObj{
		tradeLabURL: constants.TLURL,
	}

	objPortfolio = portfolioObj

	return portfolioObj
}

func (obj PortfolioObj) FetchDematHoldings(req models.FetchDematHoldingsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + FETCHDEMATHOLDINGSURL + "?client_id=" + url.QueryEscape(req.ClientID)
	//make payload
	payload := new(bytes.Buffer)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "FetchDematHoldings", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " fetchDematHoldingsReq call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("fetchDematHoldingsReq res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlFetchDematHoldingsResponse := TradeLabFetchDematHoldingsResponse{}
	json.Unmarshal([]byte(string(body)), &tlFetchDematHoldingsResponse)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " fetchDematHoldingsRes tl status not ok =", tlFetchDematHoldingsResponse.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlFetchDematHoldingsResponse.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var fetchDematHoldingsResponse models.FetchDematHoldingsResponse
	fetchDematHoldingsResponseData := make([]models.FetchDematHoldingsResponseData, 0)
	for i := 0; i < len(tlFetchDematHoldingsResponse.Data.Holdings); i++ {
		var fetchDematHoldingsData models.FetchDematHoldingsResponseData
		fetchDematHoldingsData.BranchCode = tlFetchDematHoldingsResponse.Data.Holdings[i].BranchCode
		fetchDematHoldingsData.BuyAvg = tlFetchDematHoldingsResponse.Data.Holdings[i].BuyAvg
		fetchDematHoldingsData.BuyAvgMtm = tlFetchDematHoldingsResponse.Data.Holdings[i].BuyAvgMtm
		fetchDematHoldingsData.ClientID = tlFetchDematHoldingsResponse.Data.Holdings[i].ClientID
		fetchDematHoldingsData.Exchange = tlFetchDematHoldingsResponse.Data.Holdings[i].Exchange
		fetchDematHoldingsData.FreeQuantity = tlFetchDematHoldingsResponse.Data.Holdings[i].FreeQuantity
		fetchDematHoldingsData.InstrumentDetails.Exchange = tlFetchDematHoldingsResponse.Data.Holdings[i].InstrumentDetails.Exchange
		fetchDematHoldingsData.InstrumentDetails.InstrumentName = tlFetchDematHoldingsResponse.Data.Holdings[i].InstrumentDetails.InstrumentName
		fetchDematHoldingsData.InstrumentDetails.InstrumentToken = tlFetchDematHoldingsResponse.Data.Holdings[i].InstrumentDetails.InstrumentToken
		fetchDematHoldingsData.InstrumentDetails.TradingSymbol = tlFetchDematHoldingsResponse.Data.Holdings[i].InstrumentDetails.TradingSymbol
		fetchDematHoldingsData.Isin = tlFetchDematHoldingsResponse.Data.Holdings[i].Isin
		fetchDematHoldingsData.Ltp = tlFetchDematHoldingsResponse.Data.Holdings[i].Ltp
		fetchDematHoldingsData.PendingQuantity = tlFetchDematHoldingsResponse.Data.Holdings[i].PendingQuantity
		fetchDematHoldingsData.PledgeAllow = tlFetchDematHoldingsResponse.Data.Holdings[i].PledgeAllow
		fetchDematHoldingsData.PledgeQuantity = tlFetchDematHoldingsResponse.Data.Holdings[i].PledgeQuantity
		fetchDematHoldingsData.PreviousClose = tlFetchDematHoldingsResponse.Data.Holdings[i].PreviousClose
		fetchDematHoldingsData.Quantity = tlFetchDematHoldingsResponse.Data.Holdings[i].Quantity
		fetchDematHoldingsData.Symbol = tlFetchDematHoldingsResponse.Data.Holdings[i].Symbol
		fetchDematHoldingsData.T0Price = tlFetchDematHoldingsResponse.Data.Holdings[i].T0Price
		fetchDematHoldingsData.T0Quantity = tlFetchDematHoldingsResponse.Data.Holdings[i].T0Quantity
		fetchDematHoldingsData.T1Price = tlFetchDematHoldingsResponse.Data.Holdings[i].T1Price
		fetchDematHoldingsData.T1Quantity = tlFetchDematHoldingsResponse.Data.Holdings[i].T1Quantity
		fetchDematHoldingsData.T2Price = tlFetchDematHoldingsResponse.Data.Holdings[i].T2Price
		fetchDematHoldingsData.T2Quantity = tlFetchDematHoldingsResponse.Data.Holdings[i].T2Quantity
		fetchDematHoldingsData.TodayPledgeQuantity = tlFetchDematHoldingsResponse.Data.Holdings[i].TodayPledgeQuantity
		fetchDematHoldingsData.TodayUnpledgeQuantity = tlFetchDematHoldingsResponse.Data.Holdings[i].TodayUnpledgeQuantity
		fetchDematHoldingsData.Token = tlFetchDematHoldingsResponse.Data.Holdings[i].Token
		fetchDematHoldingsData.TradingSymbol = tlFetchDematHoldingsResponse.Data.Holdings[i].TradingSymbol
		fetchDematHoldingsData.TransactionType = tlFetchDematHoldingsResponse.Data.Holdings[i].TransactionType
		fetchDematHoldingsData.UsedQuantity = tlFetchDematHoldingsResponse.Data.Holdings[i].UsedQuantity
		fetchDematHoldingsData.ActualBuyAvg = tlFetchDematHoldingsResponse.Data.Holdings[i].ActualBuyAvg
		fetchDematHoldingsData.NetHoldingQty = tlFetchDematHoldingsResponse.Data.Holdings[i].NetHoldingQty
		fetchDematHoldingsData.PledgePercentage = tlFetchDematHoldingsResponse.Data.Holdings[i].PledgePercentage
		if fetchDematHoldingsData.Exchange == "NFO" || fetchDematHoldingsData.Exchange == "BFO" || fetchDematHoldingsData.Exchange == "MCX" {
			var req models.ScripInfoRequest
			req.Exchange = fetchDematHoldingsData.Exchange
			req.Token = strconv.Itoa(fetchDematHoldingsData.Token)
			_, fetchDematHoldingsData.AdditionalInfo = GetAdditionalInfo(obj.tradeLabURL, req, reqH)
		}
		fetchDematHoldingsResponseData = append(fetchDematHoldingsResponseData, fetchDematHoldingsData)
	}
	fetchDematHoldingsResponse.Holdings = fetchDematHoldingsResponseData
	loggerconfig.Info("fetchDematHoldingsRes tl resp=", helpers.LogStructAsJSON(fetchDematHoldingsResponse), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)

	apiRes.Data = fetchDematHoldingsResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj PortfolioObj) ConvertPositions(req models.ConvertPositionsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + CONVERTPOSITIONSURL

	//fill up the TL Req
	var tlConvertPositionsReq TradeLabConvertPositionsRequest
	tlConvertPositionsReq.ClientID = req.ClientID
	tlConvertPositionsReq.Exchange = req.Exchange
	tlConvertPositionsReq.InstrumentToken = req.InstrumentToken
	tlConvertPositionsReq.Product = req.Product
	tlConvertPositionsReq.NewProduct = req.NewProduct
	tlConvertPositionsReq.Quantity = req.Quantity
	tlConvertPositionsReq.Validity = req.Validity
	tlConvertPositionsReq.OrderSide = req.OrderSide

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlConvertPositionsReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPut, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ConvertPositions", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " convertPositionsRes call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("convertPositionsRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlConvertPositionsRes := TradeLabConvertPositionsResponse{}
	json.Unmarshal([]byte(string(body)), &tlConvertPositionsRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("convertPositionsRes tl status not ok =", tlConvertPositionsRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlConvertPositionsRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var convertPositionsRes models.ConvertPositionsResponse
	convertPositionsRes.Data = tlConvertPositionsRes.Data

	loggerconfig.Info("convertPositionsRes tl resp=", helpers.LogStructAsJSON(convertPositionsRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)

	apiRes.Data = convertPositionsRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj PortfolioObj) GetPositions(req models.GetPositionRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return GetPositionsInternal(req, reqH)
}

var GetPositionsInternal = func(req models.GetPositionRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := objPortfolio.tradeLabURL + GETPOSITIONURL + "?client_id=" + url.QueryEscape(req.ClientID) + "&type=" + req.Type

	//make payload
	payload := new(bytes.Buffer) // empty payload

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "GetPositions", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " GetPositionResponse call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("GetPositionResponse res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlGetPositionResponse := TradeLabGetPositionResponse{}
	json.Unmarshal([]byte(string(body)), &tlGetPositionResponse)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " GetPositionResponse tl status not ok =", tlGetPositionResponse.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlGetPositionResponse.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var getPositionRes models.GetPositionResponse
	// fmt.Println(tlGetPositionResponse.Data)
	// getPositionRes.Data = tlGetPositionResponse.Data
	// getPositionRes.Data[0].AverageBuyPrice = tlGetPositionResponse.Data[0].AverageBuyPrice

	responseData := make([]models.GetPositionResponseData, 0)

	for i := 0; i < len(tlGetPositionResponse.Data); i++ {
		var getPositionResData models.GetPositionResponseData
		getPositionResData.ActualAverageBuyPrice = tlGetPositionResponse.Data[i].ActualAverageBuyPrice
		getPositionResData.ActualAverageSellPrice = tlGetPositionResponse.Data[i].ActualAverageSellPrice
		getPositionResData.ActualCfBuyAmount = tlGetPositionResponse.Data[i].ActualCfBuyAmount
		getPositionResData.ActualCfSellAmount = tlGetPositionResponse.Data[i].ActualCfSellAmount
		getPositionResData.OtherMargin = tlGetPositionResponse.Data[i].OtherMargin
		getPositionResData.TotalPledgeCollateral = tlGetPositionResponse.Data[i].TotalPledgeCollateral
		getPositionResData.AverageBuyPrice = tlGetPositionResponse.Data[i].AverageBuyPrice
		getPositionResData.AveragePrice = tlGetPositionResponse.Data[i].AveragePrice
		getPositionResData.AverageSellPrice = tlGetPositionResponse.Data[i].AverageSellPrice
		getPositionResData.BuyAmount = tlGetPositionResponse.Data[i].BuyAmount
		getPositionResData.BuyQuantity = tlGetPositionResponse.Data[i].BuyQuantity
		getPositionResData.CfBuyAmount = tlGetPositionResponse.Data[i].CfBuyAmount
		getPositionResData.CfBuyQuantity = tlGetPositionResponse.Data[i].CfBuyQuantity
		getPositionResData.CfSellAmount = tlGetPositionResponse.Data[i].CfSellAmount
		getPositionResData.CfSellQuantity = tlGetPositionResponse.Data[i].CfSellQuantity
		getPositionResData.ClientID = tlGetPositionResponse.Data[i].ClientID
		getPositionResData.ClosePrice = tlGetPositionResponse.Data[i].ClosePrice
		getPositionResData.Exchange = tlGetPositionResponse.Data[i].Exchange
		getPositionResData.InstrumentToken = tlGetPositionResponse.Data[i].InstrumentToken
		getPositionResData.Ltp = tlGetPositionResponse.Data[i].Ltp
		getPositionResData.Multiplier = tlGetPositionResponse.Data[i].Multiplier
		getPositionResData.NetAmount = tlGetPositionResponse.Data[i].NetAmount
		getPositionResData.NetQuantity = tlGetPositionResponse.Data[i].NetQuantity
		getPositionResData.PreviousClose = tlGetPositionResponse.Data[i].PreviousClose
		getPositionResData.ProCli = tlGetPositionResponse.Data[i].ProCli
		getPositionResData.ProdType = tlGetPositionResponse.Data[i].ProdType
		getPositionResData.Product = tlGetPositionResponse.Data[i].Product
		getPositionResData.RealizedMtm = tlGetPositionResponse.Data[i].RealizedMtm
		getPositionResData.Segment = tlGetPositionResponse.Data[i].Segment
		getPositionResData.SellAmount = tlGetPositionResponse.Data[i].SellAmount
		getPositionResData.SellQuantity = tlGetPositionResponse.Data[i].SellQuantity
		getPositionResData.Symbol = tlGetPositionResponse.Data[i].Symbol
		getPositionResData.Token = tlGetPositionResponse.Data[i].Token
		getPositionResData.TradingSymbol = tlGetPositionResponse.Data[i].TradingSymbol
		getPositionResData.SellAmount = tlGetPositionResponse.Data[i].SellAmount
		getPositionResData.VLoginID = tlGetPositionResponse.Data[i].VLoginID
		if getPositionResData.Exchange == "NFO" || getPositionResData.Exchange == "BFO" || getPositionResData.Exchange == "MCX" {
			var req models.ScripInfoRequest
			req.Exchange = getPositionResData.Exchange
			req.Token = strconv.Itoa(getPositionResData.InstrumentToken)
			_, getPositionResData.AdditionalInfo = GetAdditionalInfo(objPortfolio.tradeLabURL, req, reqH)

			// getPositionResData.AverageBuyPrice = tlGetPositionResponse.Data[i].ActualAverageBuyPrice
			// getPositionResData.AverageSellPrice = tlGetPositionResponse.Data[i].ActualAverageSellPrice
		}
		responseData = append(responseData, getPositionResData)
	}
	getPositionRes.Data = responseData

	loggerconfig.Info("GetPositionResponse tl resp=", helpers.LogStructAsJSON(getPositionRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)

	apiRes.Data = getPositionRes.Data
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}
