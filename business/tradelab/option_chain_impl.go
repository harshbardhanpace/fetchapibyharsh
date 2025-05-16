package tradelab

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type OptionChainObj struct {
	tradeLabURL string
}

func InitOptionChain() OptionChainObj {
	defer models.HandlePanic()

	optionChainObj := OptionChainObj{
		tradeLabURL: constants.TLURL,
	}

	return optionChainObj
}

func (obj OptionChainObj) FetchOptionChain(req models.FetchOptionChainRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	price := fmt.Sprintf("%f", req.Price)

	url := obj.tradeLabURL + FETCHOPTIONCHAINURL + "/" + req.Exchange + "?token=" + strconv.Itoa(req.Token) + "&num=" + strconv.Itoa(req.Num) + "&price=" + price

	//make payload
	payload := new(bytes.Buffer)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallApi(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "FetchOptionChain", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " fetchOptionChainReq call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("fetchOptionChainReq res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlFetchOptionChainResponse := TradeLabFetchOptionChainResponse{}
	json.Unmarshal([]byte(string(body)), &tlFetchOptionChainResponse)

	if res.StatusCode != http.StatusOK {
		if len(tlFetchOptionChainResponse.Error.Message) >= 23 && tlFetchOptionChainResponse.Error.Message[:23] != InvalidToken {
			loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " fetchOptionChainRes tl status not ok =", tlFetchOptionChainResponse.Error.Message, "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		} else {
			loggerconfig.Error("Alert Severity:P2-Mid, platform:", reqH.Platform, " fetchOptionChainRes tl status not ok =", tlFetchOptionChainResponse.Error.Message, "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		}

		apiRes.Message = tlFetchOptionChainResponse.Error.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var fetchOptionChainResponse models.FetchOptionChainResponse
	fetchOptionChainResponseData := make([]models.FetchOptionChainResponseData, 0)
	for i := 0; i < len(tlFetchOptionChainResponse.Result); i++ {
		var fetchOptionChainDataOuter models.FetchOptionChainResponseData
		fetchOptionChainDataOuter.ExpiryDate = tlFetchOptionChainResponse.Result[i].ExpiryDate
		for j := 0; j < len(tlFetchOptionChainResponse.Result[i].Strikes); j++ {
			var fetchOptionChainStrikeData models.FetchOptionChainResponseDataStrikes
			fetchOptionChainStrikeData.StrikePrice = tlFetchOptionChainResponse.Result[i].Strikes[j].StrikePrice
			fetchOptionChainStrikeData.CallOption.Token = tlFetchOptionChainResponse.Result[i].Strikes[j].CallOption.Token
			fetchOptionChainStrikeData.CallOption.Exchange = tlFetchOptionChainResponse.Result[i].Strikes[j].CallOption.Exchange
			fetchOptionChainStrikeData.CallOption.Company = tlFetchOptionChainResponse.Result[i].Strikes[j].CallOption.Company
			fetchOptionChainStrikeData.CallOption.Symbol = tlFetchOptionChainResponse.Result[i].Strikes[j].CallOption.Symbol
			fetchOptionChainStrikeData.CallOption.TradingSymbol = tlFetchOptionChainResponse.Result[i].Strikes[j].CallOption.TradingSymbol
			fetchOptionChainStrikeData.CallOption.DisplayName = tlFetchOptionChainResponse.Result[i].Strikes[j].CallOption.DisplayName
			fetchOptionChainStrikeData.CallOption.StrikePrice = tlFetchOptionChainResponse.Result[i].Strikes[j].CallOption.StrikePrice
			fetchOptionChainStrikeData.CallOption.ExpiryRaw = tlFetchOptionChainResponse.Result[i].Strikes[j].CallOption.ExpiryRaw
			fetchOptionChainStrikeData.CallOption.ClosePrice = tlFetchOptionChainResponse.Result[i].Strikes[j].CallOption.ClosePrice
			fetchOptionChainStrikeData.PutOption.Token = tlFetchOptionChainResponse.Result[i].Strikes[j].PutOption.Token
			fetchOptionChainStrikeData.PutOption.Exchange = tlFetchOptionChainResponse.Result[i].Strikes[j].PutOption.Exchange
			fetchOptionChainStrikeData.PutOption.Company = tlFetchOptionChainResponse.Result[i].Strikes[j].PutOption.Company
			fetchOptionChainStrikeData.PutOption.Symbol = tlFetchOptionChainResponse.Result[i].Strikes[j].PutOption.Symbol
			fetchOptionChainStrikeData.PutOption.TradingSymbol = tlFetchOptionChainResponse.Result[i].Strikes[j].PutOption.TradingSymbol
			fetchOptionChainStrikeData.PutOption.DisplayName = tlFetchOptionChainResponse.Result[i].Strikes[j].PutOption.DisplayName
			fetchOptionChainStrikeData.PutOption.StrikePrice = tlFetchOptionChainResponse.Result[i].Strikes[j].PutOption.StrikePrice
			fetchOptionChainStrikeData.PutOption.ExpiryRaw = tlFetchOptionChainResponse.Result[i].Strikes[j].PutOption.ExpiryRaw
			fetchOptionChainStrikeData.PutOption.ClosePrice = tlFetchOptionChainResponse.Result[i].Strikes[j].PutOption.ClosePrice
			fetchOptionChainDataOuter.Strikes = append(fetchOptionChainDataOuter.Strikes, fetchOptionChainStrikeData)
		}
		fetchOptionChainResponseData = append(fetchOptionChainResponseData, fetchOptionChainDataOuter)
	}
	fetchOptionChainResponse.Result = fetchOptionChainResponseData
	loggerconfig.Info("fetchOptionChainRes tl resp=", helpers.LogStructAsJSON(fetchOptionChainResponse), "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)

	apiRes.Data = fetchOptionChainResponse.Result
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj OptionChainObj) FetchFuturesChain(req models.FetchFuturesChainReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + FETCHFUTURESCHAINURL + "?token=" + req.Token

	//make payload
	payload := new(bytes.Buffer)

	//call api
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallApi(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "FetchFuturesChain", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " FetchFuturesChain call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("FetchFuturesChain res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlFetchFuturesChainResponse := TradeLabFuturesChain{}
	json.Unmarshal([]byte(string(body)), &tlFetchFuturesChainResponse)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " FetchFuturesChain tl status not ok =", tlFetchFuturesChainResponse.Error.Message, "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlFetchFuturesChainResponse.Error.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var fetchFuturesChainResponse models.FetchFuturesChainRes
	fetchFuturesChainResponseData := make([]models.FuturesChain, 0)
	for i := 0; i < len(tlFetchFuturesChainResponse.Result); i++ {
		var futuresChainData models.FuturesChain
		futuresChainData.ExpiryDate = tlFetchFuturesChainResponse.Result[i].ExpiryDate
		futuresChainData.Strikes.Token = tlFetchFuturesChainResponse.Result[i].Strikes.Token
		futuresChainData.Strikes.Exchange = tlFetchFuturesChainResponse.Result[i].Strikes.Exchange
		futuresChainData.Strikes.Company = tlFetchFuturesChainResponse.Result[i].Strikes.Company
		futuresChainData.Strikes.Symbol = tlFetchFuturesChainResponse.Result[i].Strikes.Symbol
		futuresChainData.Strikes.TradingSymbol = tlFetchFuturesChainResponse.Result[i].Strikes.TradingSymbol
		futuresChainData.Strikes.DisplayName = tlFetchFuturesChainResponse.Result[i].Strikes.DisplayName
		futuresChainData.Strikes.ExpiryRaw = tlFetchFuturesChainResponse.Result[i].Strikes.ExpiryRaw
		futuresChainData.Strikes.ClosePrice = tlFetchFuturesChainResponse.Result[i].Strikes.ClosePrice
		fetchFuturesChainResponseData = append(fetchFuturesChainResponseData, futuresChainData)
	}
	fetchFuturesChainResponse.Result = fetchFuturesChainResponseData
	loggerconfig.Info("FetchFuturesChain tl resp=", helpers.LogStructAsJSON(tlFetchFuturesChainResponse), "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)

	apiRes.Data = fetchFuturesChainResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}
