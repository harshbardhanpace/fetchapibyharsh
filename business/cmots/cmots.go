package cmots

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/db"
	"space/helpers"
	"space/helpers/cache"
	"space/loggerconfig"
	"space/models"
)

type CmotsObj struct {
	CmURL         string
	CmAuth        string
	Db            db.Database
	contractCache cache.ContractCache
}

func InitCmotsProvider(dbInstance db.Database, contractCacheCli cache.ContractCache) CmotsObj {
	defer models.HandlePanic()
	cmotsObj := CmotsObj{}
	env := constants.Env
	cmotsObj.CmURL = constants.CmURL
	cmotsObj.CmAuth = constants.CmAuth
	if env == constants.LocalEnv {
		cmotsObj.CmAuth = loggerconfig.LocalCreds.Local.CmotsAuthToken
	}
	cmotsObj.Db = dbInstance
	cmotsObj.contractCache = contractCacheCli
	return cmotsObj
}

func (obj CmotsObj) GetOverview(req models.GetOverviewReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	// dbResponse, err := db.GetPgObj().FetchOverviewData(req)
	dbResponse, err := obj.Db.FetchOverviewData(req)

	if err != nil {
		loggerconfig.Error("GetOverview FetchOverviewData failed, clientID: ", reqH.ClientId, "reqId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("GetOverview Successful, response:", helpers.LogStructAsJSON(dbResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = dbResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CmotsObj) FetchFinancials(req models.FetchFinancialsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	dbResponse, err := obj.Db.FetchFinancialsData(req)

	if err != nil {
		loggerconfig.Error("FetchFinancials FetchFinancialsData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("FetchFinancials Successful, response:", helpers.LogStructAsJSON(dbResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = dbResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CmotsObj) FetchFinancialsDetailed(req models.FetchFinancialsDetailedReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	dbResponse, err := obj.Db.FetchQuarterlyData(req)

	if err != nil {
		loggerconfig.Error("FetchFinancialsDetailed FetchQuarterlyData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("FetchFinancialsDetailed Successful, response:", helpers.LogStructAsJSON(dbResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = dbResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CmotsObj) FetchPeers(req models.FetchPeersReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	dbResponse, err := obj.Db.FetchPeersData(req)
	if err != nil {
		loggerconfig.Error("FetchPeers FetchPeersData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	if len(dbResponse) > 0 && req.Sector == constants.Sector {
		var dataFinal []models.FetchPeerData
		dataFinal = append(dataFinal, dbResponse[0])
		for i := 1; i < len(dbResponse); i++ {
			if dbResponse[0].SectorCode == dbResponse[i].SectorCode {
				dataFinal = append(dataFinal, dbResponse[i])
			}
		}
		loggerconfig.Info("FetchPeers Successful, response:", helpers.LogStructAsJSON(dbResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		apiRes.Data = dataFinal
		apiRes.Message = "SUCCESS"
		apiRes.Status = true

		return http.StatusOK, apiRes
	}

	loggerconfig.Info("FetchPeers Successful, response:", helpers.LogStructAsJSON(dbResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = dbResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CmotsObj) ShareHoldingPatterns(req models.ShareHoldingPatternsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	dbResponse, err := obj.Db.FetchShareHoldingPatternsData(req)
	if err != nil {
		loggerconfig.Error("ShareHoldingPatterns FetchShareHoldingPatternsData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("ShareHoldingPatterns Successful, response:", helpers.LogStructAsJSON(dbResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = dbResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CmotsObj) RatiosCompare(req models.RatiosCompareReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	dbResponse, err := obj.Db.FetchRatiosCompareData(req)
	if err != nil {
		loggerconfig.Error("RatiosCompare FetchRatiosCompareData failed, clientID: ", reqH.ClientId, "reqId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("RatiosCompare Successful, response:", helpers.LogStructAsJSON(dbResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = dbResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CmotsObj) FetchTechnicalIndicators(req models.FetchTechnicalIndicatorsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	dbResponse, err := obj.Db.FetchTechnicalIndicatorsData(req)
	if err != nil {
		loggerconfig.Error("FetchTechnicalIndicators FetchTechnicalIndicatorsData failed, clientID: ", reqH.ClientId, "reqId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("FetchTechnicalIndicators Successful, response:", helpers.LogStructAsJSON(dbResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = dbResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CmotsObj) StocksOnNews(req models.StocksOnNewsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes
	var stocksOnNewsResponse models.StocksOnNewsResponse
	var coCodeAndResponsePointer []models.CoCodeAndResponsePointer

	cmURL := obj.CmURL
	cmAuth := obj.CmAuth
	payload := new(bytes.Buffer)

	switch strings.ToLower(req.Filter) {
	case constants.Daily:
		URL := cmURL + constants.DAILYANNOUNCEMENT
		var dailyAnnouncementResponse models.DailyAnnouncementResponse
		res, err := apihelpers.CallCmotsApi(http.MethodGet, URL, payload, cmAuth)
		if err != nil {
			loggerconfig.Error("DailyAnnouncement CMOTS fetch error:", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		json.Unmarshal([]byte(string(body)), &dailyAnnouncementResponse)
		if err != nil {
			loggerconfig.Error("DailyAnnouncement CMOTS unmarshalling error:", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		for i := 0; i < len(dailyAnnouncementResponse.Data); i++ {
			var coCodeMapping models.CoCodeAndResponsePointer
			var stocksOnNewsResponseData models.StocksOnNewsResponseData
			coCodeMapping.CoCode = dailyAnnouncementResponse.Data[i].CoCode
			coCodeMapping.ResponsePointer = &stocksOnNewsResponseData
			stocksOnNewsResponseData.CompanyName = dailyAnnouncementResponse.Data[i].CoName
			stocksOnNewsResponseData.TradingSymbol = dailyAnnouncementResponse.Data[i].Symbol
			stocksOnNewsResponseData.Remark = dailyAnnouncementResponse.Data[i].Memo
			stocksOnNewsResponseData.Exchange = constants.EXCHANGE
			coCodeAndResponsePointer = append(coCodeAndResponsePointer, coCodeMapping)
		}
	case constants.Meeting:
		URL := cmURL + constants.BOARDMEETING
		var boardMeetingForthComingResponse models.BoardMeetingForthComingResponse
		res, err := apihelpers.CallCmotsApi(http.MethodGet, URL, payload, cmAuth)
		if err != nil {
			loggerconfig.Error("BoardMeeting CMOTS fetch error:", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		json.Unmarshal([]byte(string(body)), &boardMeetingForthComingResponse)
		if err != nil {
			loggerconfig.Error("BoardMeeting CMOTS unmarshalling error:", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		for i := 0; i < len(boardMeetingForthComingResponse.Data); i++ {
			var coCodeMapping models.CoCodeAndResponsePointer
			var stocksOnNewsResponseData models.StocksOnNewsResponseData
			coCodeMapping.CoCode = boardMeetingForthComingResponse.Data[i].CoCode
			coCodeMapping.ResponsePointer = &stocksOnNewsResponseData
			stocksOnNewsResponseData.CompanyName = boardMeetingForthComingResponse.Data[i].CoName
			stocksOnNewsResponseData.TradingSymbol = boardMeetingForthComingResponse.Data[i].Symbol
			stocksOnNewsResponseData.Remark = boardMeetingForthComingResponse.Data[i].Note
			stocksOnNewsResponseData.Exchange = constants.EXCHANGE
			coCodeAndResponsePointer = append(coCodeAndResponsePointer, coCodeMapping)
		}

	case constants.Name:
		URL := cmURL + constants.CHANGEOFNAME
		var changeOfNameResponse models.ChangeOfNameResponse
		res, err := apihelpers.CallCmotsApi(http.MethodGet, URL, payload, cmAuth)
		if err != nil {
			loggerconfig.Error("ChangeOfName CMOTS fetch error:", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		json.Unmarshal([]byte(string(body)), &changeOfNameResponse)
		if err != nil {
			loggerconfig.Error("ChangeOfName CMOTS unmarshalling error:", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		for i := 0; i < len(changeOfNameResponse.Data); i++ {
			var coCodeMapping models.CoCodeAndResponsePointer
			var stocksOnNewsResponseData models.StocksOnNewsResponseData
			coCodeMapping.CoCode = changeOfNameResponse.Data[i].CoCode
			coCodeMapping.ResponsePointer = &stocksOnNewsResponseData
			stocksOnNewsResponseData.CompanyName = changeOfNameResponse.Data[i].CoName
			stocksOnNewsResponseData.TradingSymbol = changeOfNameResponse.Data[i].Symbol
			stocksOnNewsResponseData.Remark = changeOfNameResponse.Data[i].Oldname + " changed their name to " + changeOfNameResponse.Data[i].CoName
			stocksOnNewsResponseData.Exchange = constants.EXCHANGE
			coCodeAndResponsePointer = append(coCodeAndResponsePointer, coCodeMapping)
		}
	case constants.Splits:
		URL := cmURL + constants.SPLITS + strconv.Itoa(req.Limit) + "/all/all"
		var splitsResponse models.SplitsResponse
		res, err := apihelpers.CallCmotsApi(http.MethodGet, URL, payload, cmAuth)
		if err != nil {
			loggerconfig.Error("Splits CMOTS fetch error:", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		json.Unmarshal([]byte(string(body)), &splitsResponse)
		if err != nil {
			loggerconfig.Error("Splits CMOTS unmarshalling error:", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		for i := 0; i < len(splitsResponse.Data); i++ {
			var date string
			if splitsResponse.Data[i].SplitDate != "" {
				date = "on date: " + splitsResponse.Data[i].SplitDate[:10]
			}
			var coCodeMapping models.CoCodeAndResponsePointer
			var stocksOnNewsResponseData models.StocksOnNewsResponseData
			coCodeMapping.CoCode = splitsResponse.Data[i].CoCode
			coCodeMapping.ResponsePointer = &stocksOnNewsResponseData
			stocksOnNewsResponseData.CompanyName = splitsResponse.Data[i].CoName
			stocksOnNewsResponseData.TradingSymbol = splitsResponse.Data[i].Symbol
			stocksOnNewsResponseData.Summary = splitsResponse.Data[i].SplitRatio
			stocksOnNewsResponseData.Remark = splitsResponse.Data[i].CoName + " has decided to split its stocks " + date + " by a split ratio of " + splitsResponse.Data[i].SplitRatio + ", " + splitsResponse.Data[i].Remark
			stocksOnNewsResponseData.Exchange = constants.EXCHANGE
			coCodeAndResponsePointer = append(coCodeAndResponsePointer, coCodeMapping)
		}
	case constants.Merger:
		URL := cmURL + constants.MERGER + strconv.Itoa(req.Limit) + "/all/all"
		var mergerDemergerResponse models.MergerDemergerResponse
		res, err := apihelpers.CallCmotsApi(http.MethodGet, URL, payload, cmAuth)
		if err != nil {
			loggerconfig.Error("Merger CMOTS fetch error:", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		json.Unmarshal([]byte(string(body)), &mergerDemergerResponse)
		if err != nil {
			loggerconfig.Error("Merger CMOTS unmarshalling error:", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		for i := 0; i < len(mergerDemergerResponse.Data); i++ {
			var mergerDemergerType string
			var date string
			if mergerDemergerResponse.Data[i].MergerDemergerDate != "" {
				date = " on date: " + mergerDemergerResponse.Data[i].MergerDemergerDate[:10]
			}
			if strings.ToLower(mergerDemergerResponse.Data[i].Type) == constants.Merger {
				mergerDemergerType = constants.Merged
			} else {
				mergerDemergerType = constants.Demerged
			}
			var coCodeMapping models.CoCodeAndResponsePointer
			var stocksOnNewsResponseData models.StocksOnNewsResponseData
			coCodeMapping.CoCode = mergerDemergerResponse.Data[i].CoCode
			coCodeMapping.ResponsePointer = &stocksOnNewsResponseData
			stocksOnNewsResponseData.CompanyName = mergerDemergerResponse.Data[i].CoName
			stocksOnNewsResponseData.Exchange = constants.EXCHANGE
			stocksOnNewsResponseData.Remark = mergerDemergerResponse.Data[i].CoName + " has " + mergerDemergerType + " into " + mergerDemergerResponse.Data[i].MergedIntoName + date + " with a merging ratio of " + mergerDemergerResponse.Data[i].MgrRatio
			coCodeAndResponsePointer = append(coCodeAndResponsePointer, coCodeMapping)
		}
	case constants.Div:
		URL := cmURL + constants.DIVIDENDANNOUNCEMENT + strconv.Itoa(req.Limit) + "/all/all"
		var dividendAnnouncementDataResponse models.DividendAnnouncementDataResponse
		res, err := apihelpers.CallCmotsApi(http.MethodGet, URL, payload, cmAuth)
		if err != nil {
			loggerconfig.Error("DividendAnnouncement CMOTS fetch error:", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		json.Unmarshal([]byte(string(body)), &dividendAnnouncementDataResponse)
		if err != nil {
			loggerconfig.Error("DividendAnnouncement CMOTS unmarshalling error:", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		for i := 0; i < len(dividendAnnouncementDataResponse.Data); i++ {
			var coCodeMapping models.CoCodeAndResponsePointer
			var stocksOnNewsResponseData models.StocksOnNewsResponseData
			coCodeMapping.CoCode = dividendAnnouncementDataResponse.Data[i].CoCode
			coCodeMapping.ResponsePointer = &stocksOnNewsResponseData
			stocksOnNewsResponseData.CompanyName = dividendAnnouncementDataResponse.Data[i].CoName
			stocksOnNewsResponseData.TradingSymbol = dividendAnnouncementDataResponse.Data[i].Symbol
			stocksOnNewsResponseData.Exchange = constants.EXCHANGE
			stocksOnNewsResponseData.Remark = dividendAnnouncementDataResponse.Data[i].Description
			coCodeAndResponsePointer = append(coCodeAndResponsePointer, coCodeMapping)
		}
	case constants.Bulk:
		URL := cmURL + constants.BULKDEALS + strconv.Itoa(req.Limit)
		var bulkDealsResponse models.BulkDealsResponse
		res, err := apihelpers.CallCmotsApi(http.MethodGet, URL, payload, cmAuth)
		if err != nil {
			loggerconfig.Error("BulkDeals CMOTS fetch error:", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		json.Unmarshal([]byte(string(body)), &bulkDealsResponse)
		if err != nil {
			loggerconfig.Error("BulkDeals CMOTS unmarshalling error:", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		for i := 0; i < len(bulkDealsResponse.Data); i++ {
			var coCodeMapping models.CoCodeAndResponsePointer
			var stocksOnNewsResponseData models.StocksOnNewsResponseData
			coCodeMapping.CoCode = bulkDealsResponse.Data[i].CoCode
			coCodeMapping.ResponsePointer = &stocksOnNewsResponseData
			stocksOnNewsResponseData.TradingSymbol = bulkDealsResponse.Data[i].Scripname
			stocksOnNewsResponseData.Exchange = constants.EXCHANGE
			coCodeAndResponsePointer = append(coCodeAndResponsePointer, coCodeMapping)
		}
	case constants.Block:
		URL := cmURL + constants.BLOCKDEALS + strconv.Itoa(req.Limit)
		var blockDealsResponse models.BlockDealsResponse
		res, err := apihelpers.CallCmotsApi(http.MethodGet, URL, payload, cmAuth)
		if err != nil {
			loggerconfig.Error("BlockDeals CMOTS fetch error:", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		json.Unmarshal([]byte(string(body)), &blockDealsResponse)
		if err != nil {
			loggerconfig.Error("BlockDeals CMOTS unmarshalling error:", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		for i := 0; i < len(blockDealsResponse.Data); i++ {
			var coCodeMapping models.CoCodeAndResponsePointer
			var stocksOnNewsResponseData models.StocksOnNewsResponseData
			coCodeMapping.CoCode = blockDealsResponse.Data[i].CoCode
			coCodeMapping.ResponsePointer = &stocksOnNewsResponseData
			stocksOnNewsResponseData.TradingSymbol = blockDealsResponse.Data[i].ScripName
			stocksOnNewsResponseData.Exchange = constants.EXCHANGE
			coCodeAndResponsePointer = append(coCodeAndResponsePointer, coCodeMapping)
		}
	case constants.Bonus:

	default:
		loggerconfig.Error("StocksOnNews, Invalid request, clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var stringOfCoCode string
	for i := 0; i < len(coCodeAndResponsePointer); i++ {
		coCode := strconv.Itoa(int(coCodeAndResponsePointer[i].CoCode))
		if i > 0 {
			stringOfCoCode = stringOfCoCode + "," + coCode
		} else {
			stringOfCoCode = coCode
		}
	}

	dbResponse, err := obj.Db.FetchTokenAndSymbol(stringOfCoCode, constants.CoCode)

	if err != nil {
		loggerconfig.Error("StocksOnNews,  FetchTokenAndSymbol failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	coCodeDbResponseMapping := make(map[int][]string)

	for i := 0; i < len(dbResponse); i++ {
		coCodeDbResponseMapping[int(dbResponse[i].CoCode)] = append(coCodeDbResponseMapping[int(dbResponse[i].CoCode)], dbResponse[i].CompanyName, dbResponse[i].Token, dbResponse[i].TradingSymbol)
	}

	for i := 0; i < len(coCodeAndResponsePointer); i++ {
		if dbData, isPresent := coCodeDbResponseMapping[int(coCodeAndResponsePointer[i].CoCode)]; isPresent {
			if coCodeDbResponseMapping[int(coCodeAndResponsePointer[i].CoCode)][constants.ZERO] != "" {
				coCodeAndResponsePointer[i].ResponsePointer.CompanyName = dbData[constants.ZERO]
				coCodeAndResponsePointer[i].ResponsePointer.Token = dbData[constants.ONE]
				coCodeAndResponsePointer[i].ResponsePointer.TradingSymbol = dbData[constants.TWO]
			}
		}
		stocksOnNewsResponse.Data = append(stocksOnNewsResponse.Data, *coCodeAndResponsePointer[i].ResponsePointer)
	}

	loggerconfig.Info("StocksOnNews Successful, response:", helpers.LogStructAsJSON(dbResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = stocksOnNewsResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CmotsObj) FetchSectorList(sectorCode string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	dbResponse, err := obj.Db.FetchSectorListData(sectorCode)
	if err != nil {
		loggerconfig.Error("FetchSectorList FetchSectorListData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("FetchSectorList Successful, response:", helpers.LogStructAsJSON(dbResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = dbResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CmotsObj) FetchSectorWiseCompany(req models.FetchSectorWiseCompanyReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	dbResponse, err := obj.Db.FetchSectorWiseCompanyData(req.SectCode)
	if err != nil {
		loggerconfig.Error("FetchSectorList FetchSectorListData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("FetchSectorList Successful, response:", helpers.LogStructAsJSON(dbResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = dbResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CmotsObj) FetchCompanyCategory(req models.FetchCompanyCategoryReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	stringOfisin := ""
	for i := 0; i < len(req.IsinList); i++ {
		if req.IsinList[i] != "" {
			stringOfisin += "'" + req.IsinList[i] + "',"
		}
	}
	if len(stringOfisin) < 2 || len(req.IsinList) == 0 {
		loggerconfig.Error("FetchCompanyCategory Invalid input:", stringOfisin, " clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	stringOfisin = stringOfisin[:len(stringOfisin)-1]

	dbResponse, err := obj.Db.FetchCompanyCategory(stringOfisin)
	if err != nil {
		loggerconfig.Error("FetchCompanyCategory FetchCompanyCategoryData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("FetchCompanyCategory Successful, response:", helpers.LogStructAsJSON(dbResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = dbResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CmotsObj) StocksOnNewsV2(req models.StocksOnNewsV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	stringOfisin := ""
	for i := 0; i < len(req.IsinList); i++ {
		if req.IsinList[i] != "" {
			stringOfisin += "'" + req.IsinList[i] + "',"
		}
	}
	if len(stringOfisin) < 2 || len(req.IsinList) == 0 {
		loggerconfig.Error("FetchTokenAndSymbol Invalid input:", stringOfisin, " clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	stringOfisin = stringOfisin[:len(stringOfisin)-1]

	fetchCoCode, err := obj.Db.FetchTokenAndSymbol(stringOfisin, constants.Isin)
	if err != nil {
		loggerconfig.Error("StocksOnNewsV2 FetchTokenAndSymbol FetchTokenAndSymbolData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	coCodeDbResponseMapping := make(map[int][]string)

	for i := 0; i < len(fetchCoCode); i++ {
		coCodeDbResponseMapping[int(fetchCoCode[i].CoCode)] = append(coCodeDbResponseMapping[int(fetchCoCode[i].CoCode)], fetchCoCode[i].CompanyName, fetchCoCode[i].Token, fetchCoCode[i].TradingSymbol)
	}

	var stringOfCoCode string
	for i := 0; i < len(fetchCoCode); i++ {
		coCode := strconv.Itoa(int(fetchCoCode[i].CoCode))
		if i > 0 {
			stringOfCoCode = stringOfCoCode + "," + coCode
		} else {
			stringOfCoCode = coCode
		}
	}

	switch strings.ToLower(req.Filter) {
	case constants.Daily:
		var stocksOnNewsResponseDataAll []models.StocksOnNewsResponseData
		dbResponse, err := obj.Db.FetchDailyAnnouncement(stringOfCoCode)
		if err != nil {
			loggerconfig.Error("StocksOnNewsV2 FetchDailyAnnouncement FetchDailyAnnouncementData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		for i := 0; i < len(dbResponse); i++ {
			var stocksOnNewsResponseData models.StocksOnNewsResponseData
			stocksOnNewsResponseData.CompanyName = dbResponse[i].CoName
			stocksOnNewsResponseData.TradingSymbol = dbResponse[i].Symbol
			stocksOnNewsResponseData.Remark = dbResponse[i].Memo
			stocksOnNewsResponseData.Exchange = constants.EXCHANGE
			stocksOnNewsResponseData.Token = coCodeDbResponseMapping[int(dbResponse[i].CoCode)][constants.ONE]
			stocksOnNewsResponseDataAll = append(stocksOnNewsResponseDataAll, stocksOnNewsResponseData)
		}
		apiRes.Data = stocksOnNewsResponseDataAll
	case constants.Meeting:
		var stocksOnNewsResponseDataAll []models.StocksOnNewsResponseData
		dbResponse, err := obj.Db.FetchBoardMeeting(stringOfCoCode)
		if err != nil {
			loggerconfig.Error("StocksOnNewsV2 FetchBoardMeeting FetchBoardMeetingData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		for i := 0; i < len(dbResponse); i++ {
			var stocksOnNewsResponseData models.StocksOnNewsResponseData
			stocksOnNewsResponseData.CompanyName = dbResponse[i].CoName
			stocksOnNewsResponseData.TradingSymbol = dbResponse[i].Symbol
			stocksOnNewsResponseData.Remark = dbResponse[i].Note
			stocksOnNewsResponseData.Exchange = constants.EXCHANGE
			stocksOnNewsResponseData.Token = coCodeDbResponseMapping[int(dbResponse[i].CoCode)][constants.ONE]
			stocksOnNewsResponseDataAll = append(stocksOnNewsResponseDataAll, stocksOnNewsResponseData)
		}
		apiRes.Data = stocksOnNewsResponseDataAll

	case constants.Name:
		var stocksOnNewsResponseDataAll []models.StocksOnNewsResponseData
		dbResponse, err := obj.Db.FetchChangedName(stringOfCoCode)
		if err != nil {
			loggerconfig.Error("StocksOnNewsV2 FetchChangedName FetchChangedNameData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		for i := 0; i < len(dbResponse); i++ {
			var stocksOnNewsResponseData models.StocksOnNewsResponseData
			stocksOnNewsResponseData.CompanyName = dbResponse[i].CoName
			stocksOnNewsResponseData.TradingSymbol = dbResponse[i].Symbol
			stocksOnNewsResponseData.Remark = dbResponse[i].Oldname + " changed their name to " + dbResponse[i].CoName
			stocksOnNewsResponseData.Exchange = constants.EXCHANGE
			stocksOnNewsResponseData.Token = coCodeDbResponseMapping[int(dbResponse[i].CoCode)][constants.ONE]
			stocksOnNewsResponseDataAll = append(stocksOnNewsResponseDataAll, stocksOnNewsResponseData)
		}
		apiRes.Data = stocksOnNewsResponseDataAll
	case constants.Splits:
		var stocksOnNewsResponseDataAll []models.StocksOnNewsResponseData
		dbResponse, err := obj.Db.FetchSplits(stringOfCoCode)
		if err != nil {
			loggerconfig.Error("StocksOnNewsV2 FetchSplits FetchSplitsData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		for i := 0; i < len(dbResponse); i++ {
			var date string
			if dbResponse[i].SplitDate != "" {
				date = "on date: " + dbResponse[i].SplitDate[:10]
			}
			var stocksOnNewsResponseData models.StocksOnNewsResponseData
			stocksOnNewsResponseData.CompanyName = dbResponse[i].CoName
			stocksOnNewsResponseData.TradingSymbol = dbResponse[i].Symbol
			stocksOnNewsResponseData.Remark = dbResponse[i].CoName + " has decided to split its stocks " + date + " by a split ratio of " + dbResponse[i].SplitRatio + ", " + dbResponse[i].Remark
			stocksOnNewsResponseData.Exchange = constants.EXCHANGE
			stocksOnNewsResponseData.Token = coCodeDbResponseMapping[int(dbResponse[i].CoCode)][constants.ONE]
			stocksOnNewsResponseData.Summary = dbResponse[i].SplitRatio
			stocksOnNewsResponseDataAll = append(stocksOnNewsResponseDataAll, stocksOnNewsResponseData)
		}
		apiRes.Data = stocksOnNewsResponseDataAll
	case constants.Merger:
		var stocksOnNewsResponseDataAll []models.StocksOnNewsResponseData
		dbResponse, err := obj.Db.FetchMerger(stringOfCoCode)
		if err != nil {
			loggerconfig.Error("StocksOnNewsV2 FetchMerger FetchMergerData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		for i := 0; i < len(dbResponse); i++ {
			var mergerDemergerType string
			var date string
			if dbResponse[i].MergerDemergerDate != "" {
				date = " on date: " + dbResponse[i].MergerDemergerDate[:10]
			}
			if strings.ToLower(dbResponse[i].Type) == constants.Merger {
				mergerDemergerType = constants.Merged
			} else {
				mergerDemergerType = constants.Demerged
			}
			var stocksOnNewsResponseData models.StocksOnNewsResponseData
			stocksOnNewsResponseData.CompanyName = dbResponse[i].CoName
			stocksOnNewsResponseData.TradingSymbol = coCodeDbResponseMapping[int(dbResponse[i].CoCode)][constants.TWO]
			stocksOnNewsResponseData.Remark = dbResponse[i].CoName + " has " + mergerDemergerType + " into " + dbResponse[i].MergedIntoName + date + " with a merging ratio of " + dbResponse[i].MgrRatio
			stocksOnNewsResponseData.Exchange = constants.EXCHANGE
			stocksOnNewsResponseData.Token = coCodeDbResponseMapping[int(dbResponse[i].CoCode)][constants.ONE]
			stocksOnNewsResponseDataAll = append(stocksOnNewsResponseDataAll, stocksOnNewsResponseData)
		}
		apiRes.Data = stocksOnNewsResponseDataAll
	case constants.Div:
		var stocksOnNewsResponseDataAll []models.StocksOnNewsResponseData
		dbResponse, err := obj.Db.FetchDividend(stringOfCoCode)
		if err != nil {
			loggerconfig.Error("StocksOnNewsV2 FetchDividend FetchDividendData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		for i := 0; i < len(dbResponse); i++ {
			var stocksOnNewsResponseData models.StocksOnNewsResponseData
			stocksOnNewsResponseData.CompanyName = dbResponse[i].CoName
			stocksOnNewsResponseData.TradingSymbol = dbResponse[i].Symbol
			stocksOnNewsResponseData.Remark = dbResponse[i].Description
			stocksOnNewsResponseData.Exchange = constants.EXCHANGE
			stocksOnNewsResponseData.Token = coCodeDbResponseMapping[int(dbResponse[i].CoCode)][constants.ONE]
			stocksOnNewsResponseDataAll = append(stocksOnNewsResponseDataAll, stocksOnNewsResponseData)
		}
		apiRes.Data = stocksOnNewsResponseDataAll

	case constants.Bulk:
		var stocksOnNewsResponseDataAll []models.StocksOnNewsResponseData
		dbResponse, err := obj.Db.FetchBulkDeals(stringOfCoCode)
		if err != nil {
			loggerconfig.Error("StocksOnNewsV2 FetchBulkDeals FetchBulkDealsData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		for i := 0; i < len(dbResponse); i++ {
			var stocksOnNewsResponseData models.StocksOnNewsResponseData
			stocksOnNewsResponseData.CompanyName = coCodeDbResponseMapping[int(dbResponse[i].CoCode)][constants.ZERO]
			stocksOnNewsResponseData.TradingSymbol = dbResponse[i].Scripname
			stocksOnNewsResponseData.Exchange = constants.EXCHANGE
			stocksOnNewsResponseData.Token = coCodeDbResponseMapping[int(dbResponse[i].CoCode)][constants.ONE]
			if dbResponse[i].Buysell == "S" {
				stocksOnNewsResponseData.Summary = strconv.FormatFloat(dbResponse[i].Qtyshares, 'f', -1, 64) + " quantity has been sold by " + dbResponse[i].Clientname + " at " + strconv.FormatFloat(dbResponse[i].AvgPrice, 'f', -1, 64)
			} else {
				stocksOnNewsResponseData.Summary = strconv.FormatFloat(dbResponse[i].Qtyshares, 'f', -1, 64) + " quantity has been bought by " + dbResponse[i].Clientname + " at " + strconv.FormatFloat(dbResponse[i].AvgPrice, 'f', -1, 64)
			}
			stocksOnNewsResponseDataAll = append(stocksOnNewsResponseDataAll, stocksOnNewsResponseData)
		}
		apiRes.Data = stocksOnNewsResponseDataAll

	case constants.Block:
		var stocksOnNewsResponseDataAll []models.StocksOnNewsResponseData
		dbResponse, err := obj.Db.FetchBlockDeals(stringOfCoCode)
		if err != nil {
			loggerconfig.Error("StocksOnNewsV2 FetchBlockDeals FetchBlockDealsData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		for i := 0; i < len(dbResponse); i++ {
			var stocksOnNewsResponseData models.StocksOnNewsResponseData
			stocksOnNewsResponseData.CompanyName = coCodeDbResponseMapping[int(dbResponse[i].CoCode)][constants.ZERO]
			stocksOnNewsResponseData.TradingSymbol = dbResponse[i].ScripName
			stocksOnNewsResponseData.Exchange = constants.EXCHANGE
			stocksOnNewsResponseData.Token = coCodeDbResponseMapping[int(dbResponse[i].CoCode)][constants.ONE]
			if dbResponse[i].Buysell == "S" {
				stocksOnNewsResponseData.Summary = strconv.FormatFloat(dbResponse[i].Qtyshares, 'f', -1, 64) + " quantity has been sold by " + dbResponse[i].ClientName + " at " + strconv.FormatFloat(dbResponse[i].AvgPrice, 'f', -1, 64)
			} else {
				stocksOnNewsResponseData.Summary = strconv.FormatFloat(dbResponse[i].Qtyshares, 'f', -1, 64) + " quantity has been bought by " + dbResponse[i].ClientName + " at " + strconv.FormatFloat(dbResponse[i].AvgPrice, 'f', -1, 64)
			}
			stocksOnNewsResponseDataAll = append(stocksOnNewsResponseDataAll, stocksOnNewsResponseData)
		}
		apiRes.Data = stocksOnNewsResponseDataAll

	case constants.Bonus:
		var stocksOnNewsResponseDataAll []models.StocksOnNewsResponseData
		dbResponse, err := obj.Db.FetchBonus(stringOfCoCode)
		if err != nil {
			loggerconfig.Error("StocksOnNewsV2 FetchBonus FetchBonusData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		for i := 0; i < len(dbResponse); i++ {
			var stocksOnNewsResponseData models.StocksOnNewsResponseData
			stocksOnNewsResponseData.CompanyName = coCodeDbResponseMapping[int(dbResponse[i].CoCode)][constants.ZERO]
			stocksOnNewsResponseData.TradingSymbol = coCodeDbResponseMapping[int(dbResponse[i].CoCode)][constants.TWO]
			stocksOnNewsResponseData.Exchange = constants.EXCHANGE
			stocksOnNewsResponseData.Token = coCodeDbResponseMapping[int(dbResponse[i].CoCode)][constants.ONE]
			stocksOnNewsResponseData.Remark = dbResponse[i].Remark
			stocksOnNewsResponseData.Summary = dbResponse[i].BonusDate
			stocksOnNewsResponseDataAll = append(stocksOnNewsResponseDataAll, stocksOnNewsResponseData)
		}
		apiRes.Data = stocksOnNewsResponseDataAll
	default:
		loggerconfig.Error("StocksOnNews, Invalid request, clientID: ", reqH.ClientId, "requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CmotsObj) StocksAnalyzer(req models.StocksAnalyzerReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	if len(req.Isin) == constants.ZERO {
		loggerconfig.Info("StocksAnalyzer,  Isin is empty , clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
		apiRes.Status = false
		apiRes.Message = constants.ErrorCodeMap[constants.EmptyIsin]
		apiRes.ErrorCode = constants.EmptyIsin
		return http.StatusBadRequest, apiRes
	}

	var stocksAnalyzerRes models.StocksAnalyzerRes

	var wg sync.WaitGroup

	totalCurrCallInStocksAnalyzer := 3
	wg.Add(totalCurrCallInStocksAnalyzer)

	//PLStatement
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		plstatementData, err := obj.Db.FetchPLStatementData(req.Isin)
		if err != nil {
			loggerconfig.Error("StocksAnalyzer FetchPLStatementData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
		} else {
			stocksAnalyzerRes.PLStatementRes = plstatementData
		}
	}(&wg)

	//BalanceSheets
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		balanceSheets, err := obj.Db.FetchBalanceSheetsData(req.Isin)
		if err != nil {
			loggerconfig.Error("StocksAnalyzer FetchBalanceSheetsData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
		} else {
			stocksAnalyzerRes.BalanceSheetRes = balanceSheets
		}
	}(&wg)

	// cashflow
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		cashflow, err := obj.Db.FetchCashFlowData(req.Isin)
		if err != nil {
			loggerconfig.Error("StocksAnalyzer FetchCashFlowData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
		} else {
			stocksAnalyzerRes.CashFlowRes = cashflow
		}
	}(&wg)

	wg.Wait()

	loggerconfig.Info("StocksAnalyzer Successful, response:", helpers.LogStructAsJSON(stocksAnalyzerRes), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = stocksAnalyzerRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

var CallFetchFinancialsDataV2 = func(obj CmotsObj, req models.FetchFinancialsReq) (models.FetchFinancialsV2Res, error) {
	return obj.Db.FetchFinancialsDataV2(req)
}

func (obj CmotsObj) FetchFinancialsV2(req models.FetchFinancialsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	dbResponse, err := CallFetchFinancialsDataV2(obj, req)
	if err != nil {
		loggerconfig.Error("FetchFinancialsV2 FetchFinancialsDataV2 failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("FetchFinancialsV2 Successful, response:", helpers.LogStructAsJSON(dbResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = dbResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CmotsObj) FetchPeersV2(req models.FetchPeersV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	dbResponse, err := obj.Db.FetchPeersV2Data(req)
	if err != nil {
		loggerconfig.Error("FetchPeersV2 FetchPeersV2Data failed, err= ", err, "clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	if len(dbResponse) > 0 && req.Sector == constants.Sector {
		var dataFinal models.FetchPeersV2Res
		dataFinal.CompanyList = append(dataFinal.CompanyList, dbResponse[0])
		for i := 1; i < len(dbResponse); i++ {
			if dbResponse[0].SectorCode == dbResponse[i].SectorCode {
				dataFinal.CompanyList = append(dataFinal.CompanyList, dbResponse[i])
			}
		}
		loggerconfig.Info("FetchPeersV2 Successful, response:", helpers.LogStructAsJSON(dbResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		apiRes.Data = dataFinal
		apiRes.Message = "SUCCESS"
		apiRes.Status = true

		return http.StatusOK, apiRes
	}

	loggerconfig.Info("FetchPeersV2 Successful, response:", helpers.LogStructAsJSON(dbResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = dbResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CmotsObj) FetchFinancialsV3(req models.FetchFinancialsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	dbResponse, err := obj.Db.FetchFinancialsDataV3(req)
	if err != nil {
		loggerconfig.Error("FetchFinancialsV3 FetchFinancialsDataV3 failed, clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("FetchFinancialsV3 Successful, response:", helpers.LogStructAsJSON(dbResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = dbResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CmotsObj) FetchFinancialsV4(req models.FetchFinancialsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	dbResponse, err := obj.Db.FetchFinancialsDataV4(req)
	if err != nil {
		loggerconfig.Error("FetchFinancialsV4 FetchFinancialsDataV3 failed, clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("FetchFinancialsV4 Successful, response:", helpers.LogStructAsJSON(dbResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = dbResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CmotsObj) CorporateActionsIndividual(req models.FetchCorporateActionsIndividualReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	dbResponse, err := obj.Db.FetchCorporateAnnouncements(req)
	if err != nil && err.Error() == constants.NoRowPG {
		loggerconfig.Error("CorporateActionsIndividual failed to find data for, reqPacket:", req, " clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		apiRes.Message = "Data Not Found"
		apiRes.Status = false
		return http.StatusOK, apiRes
	}
	if err != nil {
		loggerconfig.Error("CorporateActionsIndividual failed, reqPacket:", req, " clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("CorporateActionsIndividual, response:", helpers.LogStructAsJSON(dbResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = dbResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CmotsObj) CorporateActionsAll(req models.FetchCorporateActionsAllReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	dbResponse, err := obj.Db.FetchCorporateAnnouncementsAll(req)
	if err != nil && err.Error() == constants.NoRowPG {
		loggerconfig.Error("CorporateActionsIndividual failed to find data for, reqPacket:", req, " clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		apiRes.Message = "Data Not Found"
		apiRes.Status = false
		return http.StatusOK, apiRes
	}
	if err != nil {
		loggerconfig.Error("CorporateActionsAll failed, reqPacket:", req, " clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("CorporateActionsAll, response:", helpers.LogStructAsJSON(dbResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = dbResponse

	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CmotsObj) FetchSectorListV2(sectorCode string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	dbResponse, err := obj.Db.FetchSectorListData(sectorCode)
	if err != nil {
		loggerconfig.Error("FetchSectorListV2, FetchSectorListData failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	//Mapping of sectornames to the categories defined in constants
	sectorToCategory := make(map[string]string)
	for category, sectorNames := range constants.SectorMapping {
		for _, sectorName := range sectorNames {
			sectorToCategory[sectorName] = category
		}
	}

	//Mapping of categories to sector codes
	categoryToSectCodes := make(map[string][]string)
	for _, sector := range dbResponse {
		sectName := sector.SectName
		if category, exists := sectorToCategory[sectName]; exists {
			categoryToSectCodes[category] = append(categoryToSectCodes[category], sector.SectCode)
		}
	}

	var finalResponse []map[string]interface{}
	for category, sectCodes := range categoryToSectCodes {
		finalResponse = append(finalResponse, map[string]interface{}{
			"sectCode": sectCodes,
			"sectName": category,
		})
	}

	loggerconfig.Info("FetchSectorListV2 Success response:", helpers.LogStructAsJSON(finalResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = finalResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CmotsObj) FetchSectorWiseCompanyV2(req models.FetchSectorWiseCompanyReqV2, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	dbResponse, err := obj.Db.FetchSectorWiseCompanyDataV2(req.SectCode)
	if err != nil || len(dbResponse) == 0 {
		loggerconfig.Error("FetchSectorWiseCompanyV2, FetchSectorWiseCompanyDataV2 failed, clientID: ", reqH.ClientId, " reqId:", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.InvalidParameters, http.StatusBadRequest)
	}

	var nseIndices []models.Index
	var bseIndices []models.Index

	for _, code := range req.SectCode {
		if indexData, exists := constants.SectorToIndices[code]; exists {
			for _, nseIndex := range indexData.NSEIndices {
				nseIndices = append(nseIndices, models.Index{
					Name:  nseIndex.Name,
					Token: nseIndex.Token,
				})
			}
			for _, bseIndex := range indexData.BSEIndices {
				bseIndices = append(bseIndices, models.Index{
					Name:  bseIndex.Name,
					Token: bseIndex.Token,
				})
			}
		}
	}

	for i := range dbResponse[0].Companies {
		company := &dbResponse[0].Companies[i]
		isinNse := fmt.Sprintf("NSE-%s", company.Isin)
		err, val := obj.contractCache.GetFromHash("isin_data", isinNse)
		if err != nil {
			loggerconfig.Error("FetchSectorWiseCompanyV2, Error fetching data from Redis for ISIN", isinNse, "error: ", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			continue
		}

		var stockDetail models.ContractDetails
		if err := json.Unmarshal([]byte(val), &stockDetail); err != nil {
			loggerconfig.Error("FetchSectorWiseCompanyV2, Error unmarshalling stock metadata for ISIN", isinNse, " error: ", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			continue
		}

		company.Token1 = stockDetail.Token1
		company.Exchange1 = stockDetail.Exchange

		isinBse := fmt.Sprintf("BSE-%s", company.Isin)

		err, val = obj.contractCache.GetFromHash("isin_data", isinBse)
		if err != nil {
			loggerconfig.Error("FetchSectorWiseCompanyV2, Error fetching data from Redis for ISIN", isinNse, "error: ", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			continue
		}

		if err := json.Unmarshal([]byte(val), &stockDetail); err != nil {
			loggerconfig.Error("FetchSectorWiseCompanyV2, Error unmarshalling stock metadata for ISIN", isinBse, " error: ", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			continue
		}

		company.Token2 = stockDetail.Token1
		company.Exchange2 = stockDetail.Exchange

	}

	sectorWiseData := models.SectorWiseCompanyV2{
		Companies:  dbResponse[0].Companies,
		NSEIndices: nseIndices,
		BSEIndices: bseIndices,
	}

	loggerconfig.Info("FetchSectorWiseCompanyV2 Success response:", helpers.LogStructAsJSON(dbResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = sectorWiseData
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CmotsObj) GetSectorWiseStockList(page int, sectorCode, sectorName string, requestH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	//fetch sector name if sector code is given
	if sectorCode != "" {
		sectors, err := obj.Db.FetchSectorListData(sectorCode)
		if err != nil || len(sectors) == 0 {
			loggerconfig.Error("GetSectorWiseStockList, FetchSectorListData failed for sectorCode, clientID: ", requestH.ClientId, " reqId:", requestH.RequestId)
			return apihelpers.SendErrorResponse(false, constants.InvalidParameters, http.StatusBadRequest)
		}
		sectorName = sectors[0].SectName
	}

	companyList, err := obj.Db.GetSectorWiseCompanyList(page, sectorName)
	if err != nil {
		loggerconfig.Error("GetSectorWiseStockList, GetSectorWiseCompanyList failed, clientID: ", requestH.ClientId, " reqId:", requestH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.InvalidParameters, http.StatusBadRequest)
	}

	loggerconfig.Info("GetSectorWiseStockList Success response:", len(companyList), "clientID: ", requestH.ClientId, " requestId:", requestH.RequestId)
	apiRes.Data = companyList
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}
