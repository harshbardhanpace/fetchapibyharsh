package v2

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers/cache"
	"space/loggerconfig"
	"space/models"
)

type OptionChainV2Obj struct {
	tradeLabURL string
	redisCli    cache.RedisCache
}

func InitOptionChainV2(redisCli cache.RedisCache) OptionChainV2Obj {
	defer models.HandlePanic()

	optionChainObj := OptionChainV2Obj{
		tradeLabURL: constants.TLURL,
		redisCli:    redisCli,
	}

	return optionChainObj
}

func (obj OptionChainV2Obj) FetchOptionChainV2(req models.FetchOptionChainV2Request, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	callData, _ := obj.redisCli.FetchByScoreWithRange(constants.NFO_+strconv.Itoa(req.Token)+constants.Call_, req.Price, req.Num)
	fetchOptionChainCall, err := parseRedisPacket(callData)
	if err != nil {
		loggerconfig.Error("FetchOptionChainV2 Call error fetching options data from redis; error:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	putData, _ := obj.redisCli.FetchByScoreWithRange(constants.NFO_+strconv.Itoa(req.Token)+constants.Put_, req.Price, req.Num)
	fetchOptionChainPut, err := parseRedisPacket(putData)
	if err != nil {
		loggerconfig.Error("FetchOptionChainV2 Put error fetching options data from redis; error:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	lowerLen := len(fetchOptionChainPut)
	if len(fetchOptionChainPut) != len(fetchOptionChainCall) && len(fetchOptionChainPut) > len(fetchOptionChainCall) {
		lowerLen = len(fetchOptionChainCall)
	}

	var optionChainResV2 []models.FetchOptionChainResponseData

	for i := 0; i < lowerLen; i++ {
		var optionChainExpiry models.FetchOptionChainResponseData
		optionChainExpiry.ExpiryDate = fetchOptionChainPut[i].ExpiryDate
		for j := 0; j < len(fetchOptionChainPut[i].Strikes); j++ {
			var err error
			var strikes models.FetchOptionChainResponseDataStrikes
			strikes.StrikePrice, err = strconv.ParseFloat(fetchOptionChainPut[i].Strikes[j].StrikePrice, 64)
			if err != nil {
				loggerconfig.Error("FetchOptionChainV2 string to float conversion failed i=", i, " j=", j, "error:", err, "; clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
				continue
			}
			strikes.CallOption.Token = fetchOptionChainCall[i].Strikes[j].Token
			strikes.CallOption.Exchange = fetchOptionChainCall[i].Strikes[j].Exchange
			strikes.CallOption.Company = fetchOptionChainCall[i].Strikes[j].Company
			strikes.CallOption.Symbol = fetchOptionChainCall[i].Strikes[j].Symbol
			strikes.CallOption.TradingSymbol = fetchOptionChainCall[i].Strikes[j].TradingSymbol
			strikes.CallOption.DisplayName = fetchOptionChainCall[i].Strikes[j].DisplayName
			strikes.CallOption.StrikePrice = strikes.StrikePrice
			strikes.CallOption.ExpiryRaw = fetchOptionChainCall[i].Strikes[j].ExpiryRaw
			strikes.CallOption.ClosePrice = fetchOptionChainCall[i].Strikes[j].ClosePrice

			strikes.PutOption.Token = fetchOptionChainPut[i].Strikes[j].Token
			strikes.PutOption.Exchange = fetchOptionChainPut[i].Strikes[j].Exchange
			strikes.PutOption.Company = fetchOptionChainPut[i].Strikes[j].Company
			strikes.PutOption.Symbol = fetchOptionChainPut[i].Strikes[j].Symbol
			strikes.PutOption.TradingSymbol = fetchOptionChainPut[i].Strikes[j].TradingSymbol
			strikes.PutOption.DisplayName = fetchOptionChainPut[i].Strikes[j].DisplayName
			strikes.PutOption.StrikePrice = strikes.StrikePrice
			strikes.PutOption.ExpiryRaw = fetchOptionChainPut[i].Strikes[j].ExpiryRaw
			strikes.PutOption.ClosePrice = fetchOptionChainPut[i].Strikes[j].ClosePrice
			optionChainExpiry.Strikes = append(optionChainExpiry.Strikes, strikes)
		}
		optionChainResV2 = append(optionChainResV2, optionChainExpiry)
	}

	loggerconfig.Info("FetchOptionChainV2 completed execution successfully", " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)

	apiRes.Data = optionChainResV2
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj OptionChainV2Obj) FetchOptionChainByExpiryV2(req models.FetchOptionChainByExpiryV2Request, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	calls, err := obj.redisCli.FetchByRoundedScoreWithRangeAndExpiry(constants.NFO_+strconv.Itoa(req.Token)+constants.Call_, req.Price, req.Num, req.Expiry)
	if err != nil {
		loggerconfig.Error("FetchOptionChainByExpiryV2 Call error fetching options data from redis; error:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	optionDataListCall, err := parseRedisPacketByExpiry(calls)
	if err != nil {
		loggerconfig.Error("FetchOptionChainByExpiryV2 Call error unmarshalling data from redis to json; error:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	puts, err := obj.redisCli.FetchByRoundedScoreWithRangeAndExpiry(constants.NFO_+strconv.Itoa(req.Token)+constants.Put_, req.Price, req.Num, req.Expiry)
	if err != nil {
		loggerconfig.Error("FetchOptionChainByExpiryV2 Put error fetching options data from redis; error:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	optionDataListPut, err := parseRedisPacketByExpiry(puts)
	if err != nil {
		loggerconfig.Error("FetchOptionChainByExpiryV2 Put error unmarshalling data from redis to json; error:", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var optionChainRes models.FetchOptionChainResponseData
	optionChainRes.ExpiryDate = req.Expiry

	lowerLen := len(optionDataListPut)
	if len(optionDataListPut) != len(optionDataListCall) && len(optionDataListPut) > len(optionDataListCall) {
		lowerLen = len(optionDataListCall)
	}

	var strikes []models.FetchOptionChainResponseDataStrikes
	for i := 0; i < lowerLen; i++ {
		var option models.FetchOptionChainResponseDataStrikes
		option.StrikePrice, err = strconv.ParseFloat(optionDataListPut[i].StrikePrice, 64)
		if err != nil {
			loggerconfig.Error("FetchOptionChainByExpiryV2 string to float conversion failed; clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
			continue
		}
		option.CallOption.Token = optionDataListCall[i].Token
		option.CallOption.Exchange = optionDataListCall[i].Exchange
		option.CallOption.Company = optionDataListCall[i].Company
		option.CallOption.Symbol = optionDataListCall[i].Symbol
		option.CallOption.TradingSymbol = optionDataListCall[i].TradingSymbol
		option.CallOption.DisplayName = optionDataListCall[i].DisplayName
		option.CallOption.StrikePrice = option.StrikePrice
		option.CallOption.ExpiryRaw = optionDataListCall[i].ExpiryRaw
		option.CallOption.ClosePrice = optionDataListCall[i].ClosePrice

		option.PutOption.Token = optionDataListPut[i].Token
		option.PutOption.Exchange = optionDataListPut[i].Exchange
		option.PutOption.Company = optionDataListPut[i].Company
		option.PutOption.Symbol = optionDataListPut[i].Symbol
		option.PutOption.TradingSymbol = optionDataListPut[i].TradingSymbol
		option.PutOption.DisplayName = optionDataListPut[i].DisplayName
		option.PutOption.StrikePrice = option.StrikePrice
		option.PutOption.ExpiryRaw = optionDataListPut[i].ExpiryRaw
		option.PutOption.ClosePrice = optionDataListPut[i].ClosePrice
		strikes = append(strikes, option)
	}
	optionChainRes.Strikes = strikes

	loggerconfig.Info("FetchOptionChainByExpiryV2 completed execution successfully", " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)

	apiRes.Data = optionChainRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func parseRedisPacket(redisData []string) ([]models.OptionDataByExpiry, error) {
	var optionDataSlice []models.OptionData
	var optionDataByExpiryRes []models.OptionDataByExpiry
	for _, data := range redisData {
		var optionData models.OptionData
		err := json.Unmarshal([]byte(data), &optionData)
		if err != nil {
			return optionDataByExpiryRes, err
		}
		optionDataSlice = append(optionDataSlice, optionData)
	}

	optionDataByExpiry := make(map[string][]models.OptionData)
	for _, option := range optionDataSlice {
		optionDataByExpiry[option.ExpiryRaw] = append(optionDataByExpiry[option.ExpiryRaw], option)
	}

	for expiry, strikes := range optionDataByExpiry {
		sort.Slice(strikes, func(i, j int) bool {
			strikePriceI, _ := strconv.ParseFloat(strikes[i].StrikePrice, 64)
			strikePriceJ, _ := strconv.ParseFloat(strikes[j].StrikePrice, 64)
			return strikePriceI < strikePriceJ
		})
		responseData := models.OptionDataByExpiry{
			ExpiryDate: expiry,
			Strikes:    strikes,
		}
		optionDataByExpiryRes = append(optionDataByExpiryRes, responseData)
	}

	sort.Slice(optionDataByExpiryRes, func(i, j int) bool {
		return compareExpiryDates(strings.ToLower(optionDataByExpiryRes[i].ExpiryDate), strings.ToLower(optionDataByExpiryRes[j].ExpiryDate))
	})

	return optionDataByExpiryRes, nil
}

func parseRedisPacketByExpiry(redisData []string) ([]models.OptionData, error) {
	var optionDataListCall []models.OptionData
	for _, jsonString := range redisData {
		var optionDataCalls models.OptionData
		err := json.Unmarshal([]byte(jsonString), &optionDataCalls)
		if err != nil {
			loggerconfig.Error("parseRedisPacketByExpiry error unmarshalling data from redis to json; error:", err)
			return optionDataListCall, err
		} else {
			optionDataListCall = append(optionDataListCall, optionDataCalls)
		}
	}
	return optionDataListCall, nil
}

func compareExpiryDates(a, b string) bool {
	layout := constants.DateTimeFormat
	timeA, _ := time.Parse(layout, a)
	timeB, _ := time.Parse(layout, b)

	return timeA.Before(timeB)
}
