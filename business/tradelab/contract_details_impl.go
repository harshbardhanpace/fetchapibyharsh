package tradelab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/helpers/cache"
	"space/loggerconfig"
	"space/models"
)

type ContractDetailsObj struct {
	tradeLabURL      string
	contractCacheCli cache.ContractCache
}

func InitContractDetails(contractCacheCli cache.ContractCache) ContractDetailsObj {
	defer models.HandlePanic()

	contractDetailsObj := ContractDetailsObj{
		tradeLabURL:      constants.TLURL,
		contractCacheCli: contractCacheCli,
	}

	return contractDetailsObj
}

func (obj ContractDetailsObj) SearchScrip(req models.SearchScripRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + SEARCHSCRIPURL + "?key=" + url.QueryEscape(req.Key)

	//make payload
	payload := new(bytes.Buffer)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "SearchScrip", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " searchScripReq call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("searchScripReq res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlSearchScripResponse := TradeLabSearchScripResponse{}
	json.Unmarshal([]byte(string(body)), &tlSearchScripResponse)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " searchScripRes tl status not ok =", tlSearchScripResponse.Error.Message, " StatusCode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlSearchScripResponse.Error.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var searchScripResponse models.SearchScripResponse
	searchScripResponseData := make([]models.SearchScripResponseResult, 0)
	for i := 0; i < len(tlSearchScripResponse.Result); i++ {
		var searchScripData models.SearchScripResponseResult
		searchScripData.Token = tlSearchScripResponse.Result[i].Token
		searchScripData.Exchange = tlSearchScripResponse.Result[i].Exchange
		searchScripData.Execution = tlSearchScripResponse.Result[i].Execution
		searchScripData.Company = tlSearchScripResponse.Result[i].Company
		searchScripData.Symbol = tlSearchScripResponse.Result[i].Symbol
		searchScripData.Isin = tlSearchScripResponse.Result[i].Isin
		searchScripData.TradingSymbol = tlSearchScripResponse.Result[i].TradingSymbol
		searchScripData.DisplayName = tlSearchScripResponse.Result[i].DisplayName
		searchScripData.Score = tlSearchScripResponse.Result[i].Score
		searchScripData.IsTradable = tlSearchScripResponse.Result[i].IsTradable
		searchScripData.Segment = tlSearchScripResponse.Result[i].Segment
		searchScripData.Tag = tlSearchScripResponse.Result[i].Tag
		searchScripData.Expiry = tlSearchScripResponse.Result[i].Expiry
		searchScripData.ClosePrice = tlSearchScripResponse.Result[i].ClosePrice
		searchScripData.Alternate.Token = tlSearchScripResponse.Result[i].Alternate.Token
		searchScripData.Alternate.Exchange = tlSearchScripResponse.Result[i].Alternate.Exchange
		searchScripData.Alternate.Execution = tlSearchScripResponse.Result[i].Alternate.Execution
		searchScripData.Alternate.Company = tlSearchScripResponse.Result[i].Alternate.Company
		searchScripData.Alternate.Symbol = tlSearchScripResponse.Result[i].Alternate.Symbol
		searchScripData.Alternate.TradingSymbol = tlSearchScripResponse.Result[i].Alternate.TradingSymbol
		searchScripData.Alternate.DisplayName = tlSearchScripResponse.Result[i].Alternate.DisplayName
		searchScripData.Alternate.IsTradable = tlSearchScripResponse.Result[i].Alternate.IsTradable
		searchScripData.Alternate.Segment = tlSearchScripResponse.Result[i].Alternate.Segment
		searchScripData.Alternate.Tag = tlSearchScripResponse.Result[i].Alternate.Tag
		searchScripData.Alternate.Expiry = tlSearchScripResponse.Result[i].Alternate.Expiry
		searchScripData.Alternate.ClosePrice = tlSearchScripResponse.Result[i].Alternate.ClosePrice
		searchScripResponseData = append(searchScripResponseData, searchScripData)
	}
	searchScripResponse.Result = searchScripResponseData
	loggerconfig.Info("searchScripRes tl resp=", helpers.LogStructAsJSON(searchScripResponse), " StatusCode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = searchScripResponse.Result
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj ContractDetailsObj) ScripInfo(req models.ScripInfoRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + SCRIPINFOURL + req.Exchange + "?info=" + req.Info + "&token=" + req.Token

	//make payload
	payload := new(bytes.Buffer)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ScripInfo", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " scripInfoReq call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("scripInfoReq res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, "requestId: ", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlScripInfoResponse := TradeLabScripInfoResponse{}
	json.Unmarshal([]byte(string(body)), &tlScripInfoResponse)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " scripInfoRes tl status not ok =", tlScripInfoResponse.Error.Message, " StatusCode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlScripInfoResponse.Error.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var scripInfoResponse models.ScripInfoResponse
	scripInfoResponse.Result.BoardLotQuantity = tlScripInfoResponse.Result.BoardLotQuantity
	scripInfoResponse.Result.ChangeInOi = tlScripInfoResponse.Result.ChangeInOi
	scripInfoResponse.Result.Exchange = tlScripInfoResponse.Result.Exchange
	scripInfoResponse.Result.Expiry = tlScripInfoResponse.Result.Expiry
	scripInfoResponse.Result.HigherCircuitLimit = tlScripInfoResponse.Result.HigherCircuitLimit
	scripInfoResponse.Result.InstrumentName = tlScripInfoResponse.Result.InstrumentName
	scripInfoResponse.Result.InstrumentToken = tlScripInfoResponse.Result.InstrumentToken
	scripInfoResponse.Result.Isin = tlScripInfoResponse.Result.Isin
	scripInfoResponse.Result.LowerCircuitLimit = tlScripInfoResponse.Result.LowerCircuitLimit
	scripInfoResponse.Result.Multiplier = tlScripInfoResponse.Result.Multiplier
	scripInfoResponse.Result.OpenInterest = tlScripInfoResponse.Result.OpenInterest
	scripInfoResponse.Result.OptionType = tlScripInfoResponse.Result.OptionType
	scripInfoResponse.Result.Precision = tlScripInfoResponse.Result.Precision
	scripInfoResponse.Result.Series = tlScripInfoResponse.Result.Series
	scripInfoResponse.Result.Strike = tlScripInfoResponse.Result.Strike
	scripInfoResponse.Result.Symbol = tlScripInfoResponse.Result.Symbol
	scripInfoResponse.Result.TickSize = tlScripInfoResponse.Result.TickSize
	scripInfoResponse.Result.TradingSymbol = tlScripInfoResponse.Result.TradingSymbol
	scripInfoResponse.Result.UnderlyingToken = tlScripInfoResponse.Result.UnderlyingToken
	scripInfoResponse.Result.RawExpiry = tlScripInfoResponse.Result.RawExpiry
	scripInfoResponse.Result.Freeze = tlScripInfoResponse.Result.Freeze
	scripInfoResponse.Result.InstrumentType = tlScripInfoResponse.Result.InstrumentType
	scripInfoResponse.Result.IssueRate = tlScripInfoResponse.Result.IssueRate
	scripInfoResponse.Result.IssueStartDate = tlScripInfoResponse.Result.IssueStartDate
	scripInfoResponse.Result.ListDate = tlScripInfoResponse.Result.ListDate
	scripInfoResponse.Result.MaxOrderSize = tlScripInfoResponse.Result.MaxOrderSize
	scripInfoResponse.Result.PriceNumerator = tlScripInfoResponse.Result.PriceNumerator
	scripInfoResponse.Result.PriceDenominator = tlScripInfoResponse.Result.PriceDenominator
	scripInfoResponse.Result.Comments = tlScripInfoResponse.Result.Comments
	scripInfoResponse.Result.CircuitRating = tlScripInfoResponse.Result.CircuitRating
	scripInfoResponse.Result.CompanyName = tlScripInfoResponse.Result.CompanyName
	scripInfoResponse.Result.DisplayName = tlScripInfoResponse.Result.DisplayName
	scripInfoResponse.Result.RawTickSize = tlScripInfoResponse.Result.RawTickSize
	scripInfoResponse.Result.IsIndex = tlScripInfoResponse.Result.IsIndex
	scripInfoResponse.Result.Tradable = tlScripInfoResponse.Result.Tradable
	scripInfoResponse.Result.MaxSingleQty = tlScripInfoResponse.Result.MaxSingleQty
	scripInfoResponse.Result.ExpiryString = tlScripInfoResponse.Result.ExpiryString
	scripInfoResponse.Result.LocalUpdateTime = tlScripInfoResponse.Result.LocalUpdateTime
	scripInfoResponse.Result.MarketType = tlScripInfoResponse.Result.MarketType
	scripInfoResponse.Result.PriceUnits = tlScripInfoResponse.Result.PriceUnits
	scripInfoResponse.Result.TradingUnits = tlScripInfoResponse.Result.TradingUnits
	scripInfoResponse.Result.LastTradingDate = tlScripInfoResponse.Result.LastTradingDate
	scripInfoResponse.Result.TenderPeriodEndDate = tlScripInfoResponse.Result.TenderPeriodEndDate
	scripInfoResponse.Result.DeliveryStartDate = tlScripInfoResponse.Result.DeliveryStartDate
	scripInfoResponse.Result.PriceQuotation = tlScripInfoResponse.Result.PriceQuotation
	scripInfoResponse.Result.GeneralDenominator = tlScripInfoResponse.Result.GeneralDenominator
	scripInfoResponse.Result.TenderPeriodStartDate = tlScripInfoResponse.Result.TenderPeriodStartDate
	scripInfoResponse.Result.DeliveryUnits = tlScripInfoResponse.Result.DeliveryUnits
	scripInfoResponse.Result.DeliveryEndDate = tlScripInfoResponse.Result.DeliveryEndDate
	scripInfoResponse.Result.TradingUnitFactor = tlScripInfoResponse.Result.TradingUnitFactor
	scripInfoResponse.Result.DeliveryUnitFactor = tlScripInfoResponse.Result.DeliveryUnitFactor
	scripInfoResponse.Result.BookClosureEndDate = tlScripInfoResponse.Result.BookClosureEndDate
	scripInfoResponse.Result.BookClosureStartDate = tlScripInfoResponse.Result.BookClosureStartDate
	scripInfoResponse.Result.NoDeliveryDateEnd = tlScripInfoResponse.Result.NoDeliveryDateEnd
	scripInfoResponse.Result.NoDeliveryDateStart = tlScripInfoResponse.Result.NoDeliveryDateStart
	scripInfoResponse.Result.ReAdmissionDate = tlScripInfoResponse.Result.ReAdmissionDate
	scripInfoResponse.Result.RecordDate = tlScripInfoResponse.Result.RecordDate
	scripInfoResponse.Result.Warning = tlScripInfoResponse.Result.Warning
	scripInfoResponse.Result.Dpr = tlScripInfoResponse.Result.Dpr
	scripInfoResponse.Result.TradeToTrade = tlScripInfoResponse.Result.TradeToTrade
	scripInfoResponse.Result.SurveillanceIndicator = tlScripInfoResponse.Result.SurveillanceIndicator
	scripInfoResponse.Result.PartitionID = tlScripInfoResponse.Result.PartitionID
	scripInfoResponse.Result.ProductID = tlScripInfoResponse.Result.ProductID
	scripInfoResponse.Result.ProductCategory = tlScripInfoResponse.Result.ProductCategory
	scripInfoResponse.Result.MonthIdentifier = tlScripInfoResponse.Result.MonthIdentifier
	scripInfoResponse.Result.ClosePrice = tlScripInfoResponse.Result.ClosePrice
	scripInfoResponse.Result.SpecialPreopen = tlScripInfoResponse.Result.SpecialPreopen
	scripInfoResponse.Result.AlternateExchange = tlScripInfoResponse.Result.AlternateExchange
	scripInfoResponse.Result.AlternateToken = tlScripInfoResponse.Result.AlternateToken
	scripInfoResponse.Result.Asm = tlScripInfoResponse.Result.Asm
	scripInfoResponse.Result.Gsm = tlScripInfoResponse.Result.Gsm
	scripInfoResponse.Result.Execution = tlScripInfoResponse.Result.Execution
	scripInfoResponse.Result.Symbol2 = tlScripInfoResponse.Result.Symbol2
	scripInfoResponse.Result.RawTenderPeriodStartDate = tlScripInfoResponse.Result.RawTenderPeriodStartDate
	scripInfoResponse.Result.RawTenderPeriodEndDate = tlScripInfoResponse.Result.RawTenderPeriodEndDate
	scripInfoResponse.Result.YearlyHighPrice = tlScripInfoResponse.Result.YearlyHighPrice
	scripInfoResponse.Result.YearlyLowPrice = tlScripInfoResponse.Result.YearlyLowPrice
	scripInfoResponse.Result.IssueMaturityDate = tlScripInfoResponse.Result.IssueMaturityDate
	scripInfoResponse.Result.Var = tlScripInfoResponse.Result.Var
	scripInfoResponse.Result.Exposure = tlScripInfoResponse.Result.Exposure
	scripInfoResponse.Result.Span = tlScripInfoResponse.Result.Span
	scripInfoResponse.Result.HaveFutures = tlScripInfoResponse.Result.HaveFutures
	scripInfoResponse.Result.HaveOptions = tlScripInfoResponse.Result.HaveOptions
	scripInfoResponse.Result.Tag = tlScripInfoResponse.Result.Tag
	scripInfoResponse.Result.ShortCode = tlScripInfoResponse.Result.ShortCode
	scripInfoResponse.Result.IsMisEligible = tlScripInfoResponse.Result.IsMisEligible
	scripInfoResponse.Result.IsMtfEligible = tlScripInfoResponse.Result.IsMtfEligible
	scripInfoResponse.Result.ExBonusDate = tlScripInfoResponse.Result.ExBonusDate
	scripInfoResponse.Result.ExDate = tlScripInfoResponse.Result.ExDate
	scripInfoResponse.Result.Exflag = tlScripInfoResponse.Result.Exflag
	scripInfoResponse.Result.ExRightDate = tlScripInfoResponse.Result.ExRightDate
	scripInfoResponse.Result.MtfMargin = tlScripInfoResponse.Result.MtfMargin

	loggerconfig.Info("scripInfoRes tl resp=", helpers.LogStructAsJSON(scripInfoResponse), " StatusCode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = scripInfoResponse.Result
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func ScripInfoHelper(tlBaseUrl string, req models.ScripInfoRequest, reqH models.ReqHeader) (error, TradeLabScripInfoResponse) {
	var tlScripInfoResponse TradeLabScripInfoResponse
	url := tlBaseUrl + SCRIPINFOURL + req.Exchange + "?info=" + req.Info + "&token=" + req.Token

	payload := new(bytes.Buffer)

	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ScripInfo", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " ScripInfoHelper call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return err, tlScripInfoResponse
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("ScripInfoHelper res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, "requestId: ", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return err, tlScripInfoResponse
	}

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " ScripInfoHelper tl status not ok =", tlScripInfoResponse.Error.Message, " StatusCode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return err, tlScripInfoResponse
	}

	json.Unmarshal([]byte(string(body)), &tlScripInfoResponse)
	return nil, tlScripInfoResponse
}

func GetAdditionalInfo(tlBaseUrl string, req models.ScripInfoRequest, reqH models.ReqHeader) (error, models.AdditionalInfo) {
	var additionalInfo models.AdditionalInfo
	req.Info = "scrip"
	err, res := ScripInfoHelper(tlBaseUrl, req, reqH)
	if err != nil {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " GetAdditionalInfo error:", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return err, additionalInfo
	}
	additionalInfo.Expiry = res.Result.Expiry
	additionalInfo.Isin = res.Result.Isin
	additionalInfo.UnderlyingSymbol = res.Result.Symbol
	additionalInfo.Type = checkType(res.Result.InstrumentName)
	if additionalInfo.Type == "Option" {
		additionalInfo.StrikePrice, _ = extractPrice(res.Result.TradingSymbol)
		inputString := res.Result.DisplayName
		additionalInfo.IsWeekly = (len(inputString) > 0 && inputString[len(inputString)-1] == 'W')
	}
	if strings.Contains(res.Result.TradingSymbol, "CE") {
		additionalInfo.Type = additionalInfo.Type + "_CE"
	} else if strings.Contains(res.Result.TradingSymbol, "PE") {
		additionalInfo.Type = additionalInfo.Type + "_PE"
	}

	return nil, additionalInfo
}

func checkType(s string) string {
	if strings.Contains(s, "OPT") {
		return "Option"
	} else if strings.Contains(s, "FUT") {
		return "Future"
	}
	return ""
}

func extractPrice(input string) (string, error) {
	// Define the regex pattern to extract the price
	// - [A-Z]{3}    : Matches 3 uppercase letters (any valid 3-letter sequence)
	// - (\d+)       : Captures one or more digits (the price)
	// - (PE|CE)     : Matches PE or CE after the price
	re := regexp.MustCompile(`([A-Z]{3})(\d+)(PE|CE)`)

	// Match the string using the regex pattern
	match := re.FindStringSubmatch(input)

	// If no match, return an error
	if len(match) == 0 {
		return "", fmt.Errorf("price not found in input")
	}

	// Return the price (the second captured group)
	return match[2], nil
}
