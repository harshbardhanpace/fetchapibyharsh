package pockets

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	apihelpers "space/apiHelpers"
	"space/business/tradelab"
	"space/constants"
	"space/db"
	"space/dbops"
	"space/helpers"
	"space/helpers/cache"
	"space/loggerconfig"
	"space/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type ExecutePocketV2Obj struct {
	mongodb  db.MongoDatabase
	redisCli cache.RedisCache
}

func InitExecutePocketV2Provider(mongodb db.MongoDatabase, redisCli cache.RedisCache) ExecutePocketV2Obj {
	defer models.HandlePanic()
	exectutePocketObj := ExecutePocketV2Obj{mongodb: mongodb, redisCli: redisCli}
	return exectutePocketObj
}

func (obj ExecutePocketV2Obj) AdminLogin(loginReq models.AdminLoginRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var dbUser models.MongoAdmin
	err := dbops.MongoRepo.FindOne(constants.ADMINCOLLECTION, bson.M{"userId": loginReq.UserId}, &dbUser)

	//user does not exists
	var apiRes apihelpers.APIRes
	if err != nil && err.Error() == constants.MongoNoDocError {
		return apihelpers.SendErrorResponse(false, constants.InvalidUserIdOrPass, http.StatusBadRequest)
	}

	//some internal server error
	if err != nil {
		loggerconfig.Error("AdminLogin error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	//validate password here
	userPass := []byte(loginReq.Password)
	dbPass := []byte(dbUser.Password)
	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)
	if passErr != nil {
		return apihelpers.SendErrorResponse(false, constants.AdminInvalidCreds, http.StatusBadRequest)
	}

	//save the jwt token in redis
	jwtToken, err := helpers.GenerateJWT(loginReq.UserId)
	if err != nil {
		return apihelpers.SendErrorResponse(false, constants.InternalServerError, http.StatusForbidden)
	}

	err = dbops.RedisRepo.Set(loginReq.UserId, jwtToken, 360*time.Minute)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " AdminLogin Set Redis failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	//return jwt token in response body
	var adminLogRes models.AdminLoginResponse
	adminLogRes.AuthToken = jwtToken

	loggerconfig.Info("AdminLogin Successful, response:", helpers.LogStructAsJSON(adminLogRes), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = adminLogRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}

func (obj ExecutePocketV2Obj) CreatePockets(createPocketsReq models.CreatePocketsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var pockets models.MongoPockets
	err := dbops.MongoRepo.FindOne(constants.POCKETSCOLLECTION, bson.M{"pocketName": createPocketsReq.PocketName}, &pockets)

	var apiRes apihelpers.APIRes
	// if pocket already exists
	if err == nil && pockets.PocketName != "" {
		return apihelpers.SendErrorResponse(false, constants.PocketAlreadyExists, http.StatusBadRequest)
	}

	id := uuid.New().String()

	err, location := helpers.Base64toPng("pocketid", id, createPocketsReq.PocketImage)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " CreateCollections Error uploading into s3 =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	createPocketsReq.PocketImage = location

	for i := 0; i < len(createPocketsReq.PocketTokens); i++ {
		idStock := uuid.New().String()
		createPocketsReq.PocketTokens[i].StockId = idStock
		err, location := helpers.Base64toPng("stockid", createPocketsReq.PocketTokens[i].StockId, createPocketsReq.PocketTokens[i].StockImage)
		if err != nil {
			loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " CreateCollections Error uploading into s3 =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		createPocketsReq.PocketTokens[i].StockImage = location
	}

	mongoPocketDetails := &models.MongoPockets{
		PocketName:      createPocketsReq.PocketName,
		PocketShortDesc: createPocketsReq.PocketShortDesc,
		PocketLongDesc:  createPocketsReq.PocketLongDesc,
		PocketExchange:  createPocketsReq.PocketExchange,
		PocketTokens:    createPocketsReq.PocketTokens,
		PocketImage:     createPocketsReq.PocketImage,
		PocketId:        id,
	}

	filter := bson.D{{"pocketId", id}}
	update := bson.D{{"$set", mongoPocketDetails}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.POCKETSCOLLECTION, filter, update, opts)
	if err != nil {
		loggerconfig.Error("CreatePockets Mongo Upsert failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var resp models.CreatePocketsResponse
	resp.PocketId = id
	resp.PocketName = createPocketsReq.PocketName
	resp.PocketTokens = createPocketsReq.PocketTokens

	loggerconfig.Info("CreatePockets Successful, response:", helpers.LogStructAsJSON(resp), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = resp
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}

func (obj ExecutePocketV2Obj) ModifyPockets(modifyPocketsReq models.ModifyPocketsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var pockets models.MongoPockets
	err := dbops.MongoRepo.FindOne(constants.POCKETSCOLLECTION, bson.M{"pocketId": modifyPocketsReq.PocketId}, &pockets)

	var apiRes apihelpers.APIRes
	// if pocket does not exists
	if err != nil && pockets.PocketName == "" {
		return apihelpers.SendErrorResponse(false, constants.PocketDoesNotExists, http.StatusBadRequest)
	}

	// we are uploading image with the same in in modify so the already existing image will be overwritten by new image
	err, location := helpers.Base64toPng("pocketid", modifyPocketsReq.PocketId, modifyPocketsReq.PocketImage)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " ModifyPockets Error uploading into s3 =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	modifyPocketsReq.PocketImage = location

	for i := 0; i < len(modifyPocketsReq.PocketTokens); i++ {
		err, location := helpers.Base64toPng("stockid", modifyPocketsReq.PocketTokens[i].StockId, modifyPocketsReq.PocketTokens[i].StockImage)
		if err != nil {
			loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " ModifyPockets Error uploading into s3 =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		modifyPocketsReq.PocketTokens[i].StockImage = location
	}

	// assuming that admin will send currect stock id with each stock
	mongoPocketDetails := &models.MongoPockets{
		PocketName:      modifyPocketsReq.PocketName,
		PocketShortDesc: modifyPocketsReq.PocketShortDesc,
		PocketLongDesc:  modifyPocketsReq.PocketLongDesc,
		PocketExchange:  modifyPocketsReq.PocketExchange,
		PocketTokens:    modifyPocketsReq.PocketTokens,
		PocketImage:     modifyPocketsReq.PocketImage,
		PocketId:        modifyPocketsReq.PocketId,
	}

	filter := bson.D{{"pocketId", modifyPocketsReq.PocketId}}
	update := bson.D{{"$set", mongoPocketDetails}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.POCKETSCOLLECTION, filter, update, opts)
	if err != nil {
		loggerconfig.Error("ModifyPockets Mongo Upsert failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var resp models.ModifyPocketsResponse
	resp.PocketId = modifyPocketsReq.PocketId
	resp.PocketTokens = modifyPocketsReq.PocketTokens

	loggerconfig.Info("ModifyPockets Successful, response:", helpers.LogStructAsJSON(resp), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = resp
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

var CallFetchPocketMongo = func(fetchPocketsReq models.FetchPocketsDetailsRequest, obj ExecutePocketV2Obj) (models.MongoPockets, error) {
	return obj.FetchPocketMongo(fetchPocketsReq)
}

func (obj ExecutePocketV2Obj) FetchPocketMongo(fetchPocketsReq models.FetchPocketsDetailsRequest) (models.MongoPockets, error) {
	var pockets models.MongoPockets
	err := dbops.MongoRepo.FindOne(constants.POCKETSCOLLECTION, bson.M{"pocketId": fetchPocketsReq.PocketId}, &pockets)
	return pockets, err
}

func (obj ExecutePocketV2Obj) FetchPockets(fetchPocketsReq models.FetchPocketsDetailsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes
	pockets, err := CallFetchPocketMongo(fetchPocketsReq, obj)
	// if pocket does not exists
	if err != nil && pockets.PocketName == "" {
		return apihelpers.SendErrorResponse(false, constants.PocketDoesNotExists, http.StatusBadRequest)
	}

	var resp models.FetchPocketsDetailsResponse
	resp.PocketId = pockets.PocketId
	resp.PocketExchange = pockets.PocketExchange
	resp.PocketLongDesc = pockets.PocketLongDesc
	resp.PocketShortDesc = pockets.PocketShortDesc
	resp.PocketName = pockets.PocketName
	resp.PocketImage = pockets.PocketImage
	resp.PocketWebImage = pockets.PocketWebImage
	resp.PrimaryBackgroundColor = pockets.PrimaryBackgroundColor
	resp.PrimarySecondaryColor = pockets.PrimarySecondaryColor
	resp.PocketTokens = pockets.PocketTokens

	loggerconfig.Info("FetchPockets Successful, response:", helpers.LogStructAsJSON(resp), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = resp
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ExecutePocketV2Obj) DeletePockets(deletePocketsReq models.DeletePocketsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	var pockets models.MongoPockets
	err := dbops.MongoRepo.FindOne(constants.POCKETSCOLLECTION, bson.M{"pocketId": deletePocketsReq.PocketId}, &pockets)

	var apiRes apihelpers.APIRes
	// if pocket does not exists
	if err != nil && pockets.PocketName == "" {
		return apihelpers.SendErrorResponse(false, constants.PocketDoesNotExists, http.StatusBadRequest)
	}

	filter := bson.D{{"pocketId", deletePocketsReq.PocketId}}
	opts := options.Delete()
	_, err = dbops.MongoRepo.DeleteOne(constants.POCKETSCOLLECTION, filter, opts)
	if err != nil {
		loggerconfig.Error("DeletePockets Mongo Upsert failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var resp models.DeletePocketsResponse
	resp.PocketTokens = pockets.PocketTokens

	loggerconfig.Info("DeletePockets Successful, response:", helpers.LogStructAsJSON(resp), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = resp
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ExecutePocketV2Obj) FetchAllPockets(reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var allPockets models.FetchAllPocketsDetailsResponse

	allDoc, err := dbops.MongoRepo.Find(constants.POCKETSCOLLECTION, bson.M{})
	if err != nil {
		loggerconfig.Error("FetchAllPockets Error Mongo Find() =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
	}
	err = allDoc.Decode(&allPockets)

	for allDoc.Next(context.Background()) {
		var fetchPocketsDetailsResponse models.FetchPocketsDetailsResponse
		err := allDoc.Decode(&fetchPocketsDetailsResponse)
		if err != nil {
			loggerconfig.Error("FetchAllPockets Error Parsing Mongo response error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		}

		allPockets.FetchAllPocketsDetailsResponse = append(allPockets.FetchAllPocketsDetailsResponse, fetchPocketsDetailsResponse)
	}

	var apiRes apihelpers.APIRes

	loggerconfig.Info("FetchAllPockets Successful, response:", helpers.LogStructAsJSON(allPockets), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = allPockets
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ExecutePocketV2Obj) PocketsCalculations(req models.PocketsCalculationsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	var pockets models.MongoPockets
	err := dbops.MongoRepo.FindOne(constants.POCKETSCOLLECTION, bson.M{"pocketId": req.PocketId}, &pockets)
	if err != nil && pockets.PocketName == "" {
		return apihelpers.SendErrorResponse(false, constants.PocketDoesNotExists, http.StatusBadRequest)
	}

	EndTime := UnixTime(helpers.GetCurrentTimeInIST().String(), 0)
	StartTime := UnixTime(helpers.GetCurrentTimeInIST().String(), req.TimeInterval)

	urlBenchmark := constants.TLURL + tradelab.Charts + "?exchange=" + req.BenchmarkExchange + "&token=" + req.BenchmarkToken + "&candletype=3" + "&starttime=" + StartTime + "&endtime=" + EndTime + "&data_duration=1"

	var status bool
	benchmarkData, status := CallTLforData(urlBenchmark, reqH)
	if !status {
		return apihelpers.SendErrorResponse(false, constants.TLChartDataFetchFailed, http.StatusInternalServerError)
	}

	benchmarkGain := make([]models.CalculatedChartData, 0)
	for i := len(benchmarkData) - 1; i > 0; i-- {
		var calculatedChartData models.CalculatedChartData
		calculatedChartData.PercentageGain = ((benchmarkData[i].Close / benchmarkData[0].Close * 100) - 100)
		calculatedChartData.Date = benchmarkData[i].Timestamp
		benchmarkGain = append(benchmarkGain, calculatedChartData)
	}
	var pocketsCalcRes models.PocketsCalculationsRes
	pocketsCalcRes.Benchmark = benchmarkGain

	var individualData [][]models.TLCandleData
	for i := 0; i < len(pockets.PocketTokens); i++ {
		url := constants.TLURL + tradelab.Charts + "?exchange=" + pockets.PocketTokens[i].Exchange + "&token=" + pockets.PocketTokens[i].Token + "&candletype=3" + "&starttime=" + StartTime + "&endtime=" + EndTime + "&data_duration=1"
		Data, status := CallTLforData(url, reqH)
		individualData = append(individualData, Data)
		if !status {
			return apihelpers.SendErrorResponse(false, constants.TLChartDataFetchFailed, http.StatusInternalServerError)
		}
	}

	var finalTotalPrice []models.TLCandleData
	for i := 0; i < len(individualData[0]); i++ {
		var totalPrice []models.TotalPrice
		for j := 0; j < len(individualData); j++ {
			var indiPrice models.TotalPrice
			indiPrice.TokenPrice = individualData[j][i].Close
			indiPrice.Quantity, _ = strconv.Atoi(pockets.PocketTokens[j].Qty)
			totalPrice = append(totalPrice, indiPrice)
		}
		var finalTPinstance models.TLCandleData
		finalTPinstance.Close = TotalPrice(totalPrice)
		finalTPinstance.Timestamp = individualData[0][i].Timestamp
		finalTotalPrice = append(finalTotalPrice, finalTPinstance)
	}

	performancePocket := make([]models.CalculatedChartData, 0)
	for i := len(finalTotalPrice) - 1; i > 0; i-- {
		var calculatedChartData models.CalculatedChartData
		calculatedChartData.PercentageGain = ((finalTotalPrice[i].Close / finalTotalPrice[0].Close * 100) - 100)
		calculatedChartData.Date = finalTotalPrice[i].Timestamp
		performancePocket = append(performancePocket, calculatedChartData)
	}
	pocketsCalcRes.Pocket = performancePocket

	loggerconfig.Info("PocketsCalculations Successful, response:", helpers.LogStructAsJSON(pocketsCalcRes), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	apiRes.Data = pocketsCalcRes

	return http.StatusOK, apiRes
}

func (obj ExecutePocketV2Obj) MultipleAndIndividualStocksCalculations(req models.MultipleAndIndividualStocksCalculationsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	var listOfTokens []models.TokenExchange
	var tokenExchange models.TokenExchange

	for i := 0; i < len(req.Stocks); i++ {
		tokenExchange.Exchange = req.Stocks[i].StockExchange
		tokenExchange.Token = req.Stocks[i].StockToken
		listOfTokens = append(listOfTokens, tokenExchange)
	}

	EndTime := UnixTime(helpers.GetCurrentTimeInIST().String(), 0)
	StartTime := UnixTime(helpers.GetCurrentTimeInIST().String(), req.TimeInterval)
	urlBenchmark := constants.TLURL + tradelab.Charts + "?exchange=" + req.BenchmarkExchange + "&token=" + req.BenchmarkToken + "&candletype=3" + "&starttime=" + StartTime + "&endtime=" + EndTime + "&data_duration=1"
	var status bool
	benchmarkData, status := CallTLforData(urlBenchmark, reqH)
	if !status {
		return apihelpers.SendErrorResponse(false, constants.TLChartDataFetchFailed, http.StatusInternalServerError)
	}

	var stocksCalcRes models.MultipleAndIndividualStocksCalculationsRes

	var individualData [][]models.TLCandleData
	minLen := len(benchmarkData)
	for i := 0; i < len(req.Stocks); i++ {
		url := constants.TLURL + tradelab.Charts + "?exchange=" + req.Stocks[i].StockExchange + "&token=" + req.Stocks[i].StockToken + "&candletype=3" + "&starttime=" + StartTime + "&endtime=" + EndTime + "&data_duration=1"
		Data, status := CallTLforData(url, reqH)
		if len(Data) < minLen {
			minLen = len(Data)
		}
		individualData = append(individualData, Data)
		if !status {
			return apihelpers.SendErrorResponse(false, constants.TLChartDataFetchFailed, http.StatusInternalServerError)
		}
	}
	benchmarkData = benchmarkData[len(benchmarkData)-minLen:]
	benchmarkGain := make([]models.CalculationDataTemp, 0)
	for i := 0; i < len(benchmarkData); i++ {
		var calculatedChartData models.CalculationDataTemp
		calculatedChartData.PercentageGain = ((benchmarkData[i].Close / benchmarkData[0].Close * 100) - 100)
		calculatedChartData.Date = benchmarkData[i].Timestamp
		benchmarkGain = append(benchmarkGain, calculatedChartData)
	}

	var finalTotalPrice []models.TLCandleData
	for k := 0; k < len(individualData); k++ {
		individualData[k] = individualData[k][len(individualData[k])-minLen:]
	}

	for i := 0; i < len(individualData[0]); i++ {
		var totalPrice []models.TotalPrice
		for j := 0; j < len(individualData); j++ {
			var indiPrice models.TotalPrice
			indiPrice.TokenPrice = individualData[j][i].Close
			indiPrice.Quantity = req.Stocks[j].StockQuantity
			totalPrice = append(totalPrice, indiPrice)
		}
		var finalTPinstance models.TLCandleData
		finalTPinstance.Close = TotalPrice(totalPrice)
		finalTPinstance.Timestamp = individualData[len(individualData)-1][i].Timestamp
		finalTotalPrice = append(finalTotalPrice, finalTPinstance)
	}
	performanceStocks := make([]models.CalculatedChartData, 0)

	for i := len(finalTotalPrice) - 1; i > 0; i-- {
		var calculatedChartData models.CalculatedChartData
		calculatedChartData.PercentageGain = ((finalTotalPrice[i].Close / finalTotalPrice[0].Close * 100) - 100)
		calculatedChartData.Date = finalTotalPrice[i].Timestamp
		performanceStocks = append(performanceStocks, calculatedChartData)
	}

	j := 0
	for i := len(benchmarkGain) - 1; i >= 0 && j < len(performanceStocks); i-- {
		if performanceStocks[j].Date == benchmarkGain[i].Date {
			performanceStocks[j].PercentageGainBenchmark = benchmarkGain[i].PercentageGain
			j++
		}
	}
	stocksCalcRes.MultipleOrIndividualStocks = performanceStocks

	loggerconfig.Info("MultipleAndIndividualStocksCalculations Successful, response:", helpers.LogStructAsJSON(stocksCalcRes), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	apiRes.Data = stocksCalcRes

	return http.StatusOK, apiRes
}

func UnixTime(timeStamp string, interval int) string {
	var yearMinus, monthMinus, daysMinus int
	if interval >= 365 {
		yearMinus = interval / 365
		interval = interval % 365
	}
	if interval >= 30 {
		monthMinus = interval / 30
		interval = interval % 30
	}
	daysMinus = interval
	year, _ := strconv.Atoi(timeStamp[0:4])
	month, _ := strconv.Atoi(timeStamp[5:7])
	day, _ := strconv.Atoi(timeStamp[8:10])
	year = year - yearMinus
	month = month - monthMinus
	day = day - daysMinus
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Unix()
	return strconv.FormatInt(date, 10)
}

var CallTLforData = func(url string, reqH models.ReqHeader) ([]models.TLCandleData, bool) {
	return CallTLforDataActual(url, reqH)
}

func CallTLforDataActual(url string, reqH models.ReqHeader) ([]models.TLCandleData, bool) {
	candleData := make([]models.TLCandleData, 0)

	payload := new(bytes.Buffer)
	//call api
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " CallTLforData call api error =", err, " requestId:", reqH.RequestId)
		return candleData, false
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	tlErrorRes := tradelab.TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == tradelab.TLERROR {
		loggerconfig.Error("CallTLforData res error =", tlErrorRes.Message, " requestId:", reqH.RequestId)
		return candleData, false
	}

	tlChartDataBenchmark := tradelab.TradeLabChartDataResponse{}
	json.Unmarshal([]byte(string(body)), &tlChartDataBenchmark)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " CallTLforData tl status not ok =", tlChartDataBenchmark.Status, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)
		return candleData, false
	}
	loggerconfig.Info("CallTLforDataActual TradelabResponse:=", helpers.LogStructAsJSON(tlChartDataBenchmark))

	start := 0
	for i := 0; i < len(tlChartDataBenchmark.Data.Candles); i++ {
		var candleDataEntry models.TLCandleData
		write := true
		for j := 0; j < len(tlChartDataBenchmark.Data.Candles[i]); j++ {
			if j == 0 {
				candleDataEntry.Timestamp = tlChartDataBenchmark.Data.Candles[i][j].(string)
				if i > 0 {
					if tlChartDataBenchmark.Data.Candles[i][j].(string) == tlChartDataBenchmark.Data.Candles[start][j].(string) {
						write = false
					} else {
						start++
					}
				}
			} else if j == 4 {
				candleDataEntry.Close = tlChartDataBenchmark.Data.Candles[i][j].(float64)
			}
		}
		if write {
			candleData = append(candleData, candleDataEntry)
		}
	}
	loggerconfig.Info("CallTLforDataActual candleData:=", helpers.LogStructAsJSON(candleData))
	return candleData, true
}

func TotalPrice(PriceXQtyArray []models.TotalPrice) float64 {
	var totalPrice float64
	totalPrice = 0
	for i := 0; i < len(PriceXQtyArray); i++ {
		totalPrice = totalPrice + PriceXQtyArray[i].TokenPrice*float64(PriceXQtyArray[i].Quantity)
	}
	return totalPrice
}

// BUY and SELL both handled by operation keyword
func (obj ExecutePocketV2Obj) ExecutePocket(req models.ExecutePocketRequest, operation string, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	var apiRes apihelpers.APIRes

	var pocketTransactionComplete models.PocketTransactionComplete
	err := dbops.MongoRepo.FindOne(constants.POCKETSTRANSACTIONS, bson.M{"clientId": req.ClientId}, &pocketTransactionComplete)
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("ExecutePocket BuyPocket", req, " mongo err:", err, "clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.PocketDoesNotExists, http.StatusBadRequest)
	}

	pocketCounter := 0
	if pocketTransactionComplete.PocketCounter == 0 {
		pocketCounter = 1
	} else {
		pocketCounter = pocketTransactionComplete.PocketCounter + 1
	}

	var pockets models.MongoPockets
	err = dbops.MongoRepo.FindOne(constants.POCKETSCOLLECTION, bson.M{"pocketId": req.PocketId}, &pockets)

	if err != nil && err.Error() == constants.MongoNoDocError {
		loggerconfig.Error("ExecutePocket BuyPocket", req, " mongo err:", err, "clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.PocketDoesNotExists, http.StatusBadRequest)
	}

	var tlCreateBasketReq models.CreateBasketReq
	tlCreateBasketReq.LoginID = req.ClientId
	tlCreateBasketReq.Name = pockets.PocketName + strconv.Itoa(pocketCounter) // "Buy" + pockets.PocketName + helpers.GetCurrentTimeInIST().Format("2017-09-07 17:06:06")
	tlCreateBasketReq.Type = constants.BASKETTYPE
	tlCreateBasketReq.ProductType = constants.BASKETPRODUCTTYPE
	tlCreateBasketReq.OrderType = constants.BASKETPRODUCTTYPE

	basketOrderObj := tradelab.InitBasketOrder()
	status, res := tradelab.BasketOrderObj.CreateBasket(basketOrderObj, tlCreateBasketReq, reqH)
	if status != http.StatusOK {
		loggerconfig.Error("ExecutePocket BuyPocket in CreateBasket status != 200", status, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	createBasketData, ok := res.Data.(models.BasketDataRes)
	if !ok {
		loggerconfig.Error("ExecutePocket BuyPocket interface parsing error", ok, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var basketId string
	if createBasketData.Name == tlCreateBasketReq.Name {
		basketId = createBasketData.BasketID
	}

	var reqAddBasketInstrument models.AddBasketInstrumentReq
	reqAddBasketInstrument.BasketID = basketId
	reqAddBasketInstrument.Name = tlCreateBasketReq.Name

	// var resAddBasketInstrumentAll apihelpers.APIRes

	for i := 0; i < len(pockets.PocketTokens); i++ {

		reqAddBasketInstrument.OrderInfo.Exchange = pockets.PocketTokens[i].Exchange
		instrumentToken, _ := strconv.Atoi(pockets.PocketTokens[i].Token)
		reqAddBasketInstrument.OrderInfo.InstrumentToken = instrumentToken
		reqAddBasketInstrument.OrderInfo.ClientID = req.ClientId
		reqAddBasketInstrument.OrderInfo.OrderType = constants.ORDETTYPE
		quantity, _ := strconv.Atoi(pockets.PocketTokens[i].Qty)
		reqAddBasketInstrument.OrderInfo.Quantity = quantity
		reqAddBasketInstrument.OrderInfo.DisclosedQuantity = constants.ORDERDISCLOSEDQUANTITY
		reqAddBasketInstrument.OrderInfo.Validity = constants.ORDERVALIDITY
		reqAddBasketInstrument.OrderInfo.Product = constants.ORDERPRODUCT
		reqAddBasketInstrument.OrderInfo.TradingSymbol = pockets.PocketTokens[i].TradingSymbol
		reqAddBasketInstrument.OrderInfo.OrderSide = operation
		reqAddBasketInstrument.OrderInfo.UnderlyingToken = pockets.PocketTokens[i].Token
		reqAddBasketInstrument.OrderInfo.Series = constants.ORDERSERIES               //pockets.PocketTokens[i].Exchange
		reqAddBasketInstrument.OrderInfo.ExecutionType = constants.ORDEREXECUTIONTYPE // pockets.PocketTokens[i].Exchange
		reqAddBasketInstrument.OrderInfo.UserOrderID = constants.USERORERID

		statusAddBasketInstrument, resAddBasketInstrument := tradelab.BasketOrderObj.AddBasketInstrument(basketOrderObj, reqAddBasketInstrument, reqH)

		if statusAddBasketInstrument != http.StatusOK {
			loggerconfig.Error("ExecutePocket BuyPocket in AddBasketInstrument status != 200", status, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		loggerconfig.Info("ExecutePocket BuyPocket AddBasketInstrument response for instrument token =", pockets.PocketTokens[i].Token, " is ", resAddBasketInstrument, " uccId:", req.ClientId, " requestId:", reqH.RequestId)

	}

	var executionBasketReq models.ExecuteBasketReq
	executionBasketReq.BasketID = basketId
	executionBasketReq.ClientID = req.ClientId
	executionBasketReq.ExecutionType = constants.BASKETEXECUTETYPE
	executionBasketReq.SquareOff = constants.BASKETSQUAREOFF
	executionBasketReq.Name = tlCreateBasketReq.Name
	executionBasketReq.ExecutionState = constants.BASKETEXECUTIONSTATE

	statusExecuteBasket, resExecuteBasket := tradelab.BasketOrderObj.ExecuteBasket(basketOrderObj, executionBasketReq, reqH)
	if statusExecuteBasket != http.StatusOK {
		loggerconfig.Error("ExecutePocket BuyPocket in ExecuteBasket status != 200", status, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	resExecuteBasketData, ok := resExecuteBasket.Data.(models.ExecuteBasketRes)
	if !ok {
		loggerconfig.Error("ExecutePocket BuyPocket ExecuteBasket interface parsing error", ok, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var pocketInfo models.PocketInfo
	pocketInfo.PocketID = pockets.PocketId
	pocketInfo.PocketName = pockets.PocketName
	pocketInfo.TransactionDate = helpers.GetCurrentTimeInIST()
	pocketInfo.PocketExecutedPrice = req.PocketPrice
	pocketInfo.TransactionID = uuid.New().String()
	if operation == strings.ToUpper(constants.BUY) {
		pocketInfo.TransactionStatus = 0
	} else { // SELL
		pocketInfo.TransactionStatus = 1
	}
	pocketInfo.BasketName = tlCreateBasketReq.Name

	pocketTransactionComplete.ClientId = req.ClientId
	pocketTransactionComplete.PocketCounter = pocketCounter
	pocketTransactionComplete.AllPocketPurchases = append(pocketTransactionComplete.AllPocketPurchases, pocketInfo)

	filter := bson.D{{"clientId", req.ClientId}}
	update := bson.D{{"$set", pocketTransactionComplete}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.POCKETSTRANSACTIONS, filter, update, opts)
	if err != nil {
		loggerconfig.Error("ExecutePocket Mongo Upsert failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var executePocketResponse models.ExecutePocketResponse
	executePocketResponse.BasketID = resExecuteBasketData.Data.BasketID
	executePocketResponse.Message = resExecuteBasketData.Data.Message

	loggerconfig.Info("ExecutePocket response is = ", helpers.LogStructAsJSON(executePocketResponse), " uccId:", req.ClientId, " requestId:", reqH.RequestId)

	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	apiRes.Data = executePocketResponse

	return http.StatusOK, apiRes
}

func (obj ExecutePocketV2Obj) FetchPocketPortfolio(req models.FetchPocketPortfolioRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	var apiRes apihelpers.APIRes

	var pocketTransactionComplete models.PocketTransactionComplete
	err := dbops.MongoRepo.FindOne(constants.POCKETSTRANSACTIONS, bson.M{"clientId": req.ClientId}, &pocketTransactionComplete)
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("FetchPocketPortfolio", req, " mongo err:", err, "clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.PocketDoesNotExists, http.StatusBadRequest)
	}

	mapIdPocket := make(map[string][]models.PocketInfo)
	for i := 0; i < len(pocketTransactionComplete.AllPocketPurchases); i++ {
		id := pocketTransactionComplete.AllPocketPurchases[i].PocketID
		mapIdPocket[id] = append(mapIdPocket[id], pocketTransactionComplete.AllPocketPurchases[i])
	}

	var fetchPocketPortfolioResponse models.FetchPocketPortfolioResponse
	fetchPocketPortfolioResponse.ClientId = req.ClientId

	for keys := range mapIdPocket {
		var pocket models.Pocket
		pocket.PocketID = keys
		totalBuyPockets := 0
		totalSellPockets := 0
		totalPrice := 0.0
		calPrice := 0.0
		totalBuyPrice := 0.0
		totalSellPrice := 0.0
		for i := 0; i < len(mapIdPocket[keys]); i++ {
			pocketQty := 0
			if mapIdPocket[keys][i].Qty == 0 {
				pocketQty = 1
			} else {
				pocketQty = mapIdPocket[keys][i].Qty
			}
			// mapIdPocket[keys][i].Qty == 0 ? 1 : mapIdPocket[keys][i].Qty
			calPrice += mapIdPocket[keys][i].PocketExecutedPrice * (float64)(pocketQty)
			if mapIdPocket[keys][i].TransactionStatus == 0 { // Buy
				totalBuyPockets += pocketQty
				totalPrice += mapIdPocket[keys][i].PocketExecutedPrice * (float64)(pocketQty)
				totalBuyPrice += mapIdPocket[keys][i].PocketExecutedPrice * (float64)(pocketQty)
			} else {
				totalSellPockets += pocketQty
				totalPrice -= mapIdPocket[keys][i].PocketExecutedPrice * (float64)(pocketQty) // Sell
				totalSellPrice += mapIdPocket[keys][i].PocketExecutedPrice * (float64)(pocketQty)
			}
		}
		if totalBuyPockets == totalSellPockets {
			continue // don't send this pocket to frontend
		}
		pocket.TotalBuyPockets = totalBuyPockets
		pocket.TotalSellPockets = totalSellPockets
		pocket.AveragePrice = totalBuyPrice / float64(totalBuyPockets)
		pocket.TotalInvestment = pocket.AveragePrice * float64(totalBuyPockets-totalSellPockets)

		// pocket image
		var storedPocket models.MongoPockets
		err := dbops.MongoRepo.FindOne(constants.POCKETSCOLLECTION, bson.M{"pocketId": pocket.PocketID}, &storedPocket)
		if err != nil && err.Error() != constants.MongoNoDocError {
			loggerconfig.Error("FetchPocketPortfolio Stored Pocket fetch failed with pocket id:", pocket.PocketID, " mongo err:", err, "clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
			return apihelpers.SendErrorResponse(false, constants.PocketDoesNotExists, http.StatusBadRequest)
		}
		pocket.PocketName = storedPocket.PocketName
		pocket.PocketImage = storedPocket.PocketImage

		fetchPocketPortfolioResponse.PortfolioDetails = append(fetchPocketPortfolioResponse.PortfolioDetails, pocket)
	}

	loggerconfig.Info("FetchPocketPortfolio response is = ", helpers.LogStructAsJSON(fetchPocketPortfolioResponse), " uccId:", req.ClientId, " requestId:", reqH.RequestId)

	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	apiRes.Data = fetchPocketPortfolioResponse

	return http.StatusOK, apiRes
}

func (obj ExecutePocketV2Obj) FetchPocketTransaction(req models.FetchPocketTransactionReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes
	var pocketTransactionComplete models.PocketTransactionComplete
	err := dbops.MongoRepo.FindOne(constants.POCKETSTRANSACTIONS, bson.M{"clientId": req.ClientId}, &pocketTransactionComplete)
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("FetchPocketTransaction", req, " mongo err:", err, "clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.PocketDoesNotExists, http.StatusBadRequest)
	}

	for i := 0; i < len(pocketTransactionComplete.AllPocketPurchases); i++ {
		// pocket image
		var storedPocket models.MongoPockets
		var err error
		err = dbops.MongoRepo.FindOne(constants.POCKETSCOLLECTION, bson.M{"pocketId": pocketTransactionComplete.AllPocketPurchases[i].PocketID}, &storedPocket)
		if err != nil && err.Error() != constants.MongoNoDocError {
			loggerconfig.Error("FetchPocketTransaction Stored Pocket fetch failed with pocket id:", pocketTransactionComplete.AllPocketPurchases[i].PocketID, " mongo err:", err, "clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
			return apihelpers.SendErrorResponse(false, constants.PocketDoesNotExists, http.StatusBadRequest)
		}
		pocketTransactionComplete.AllPocketPurchases[i].PocketImage = storedPocket.PocketImage
	}

	loggerconfig.Info("FetchPocketTransaction response is = ", helpers.LogStructAsJSON(pocketTransactionComplete), " uccId:", req.ClientId, " requestId:", reqH.RequestId)

	apiRes.Data = pocketTransactionComplete
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ExecutePocketV2Obj) StorePocketTransaction(req models.StorePocketTransactionReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	var apiRes apihelpers.APIRes

	var pocketTransactionComplete models.PocketTransactionComplete
	err := dbops.MongoRepo.FindOne(constants.POCKETSTRANSACTIONS, bson.M{"clientId": req.ClientId}, &pocketTransactionComplete)
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("StorePocketTransaction", req, " mongo err:", err, "clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.PocketDoesNotExists, http.StatusBadRequest)
	}

	pocketCounter := 0
	if pocketTransactionComplete.PocketCounter == 0 {
		if req.Qty == 0 {
			pocketCounter = 1
		} else {
			pocketCounter = req.Qty
		}
	} else {
		if req.Qty == 0 {
			pocketCounter = pocketTransactionComplete.PocketCounter + 1
		} else {
			pocketCounter = pocketTransactionComplete.PocketCounter + req.Qty
		}
	}

	var pockets models.MongoPockets
	err = dbops.MongoRepo.FindOne(constants.POCKETSCOLLECTION, bson.M{"pocketId": req.PocketId}, &pockets)

	if err != nil && err.Error() == constants.MongoNoDocError {
		loggerconfig.Error("StorePocketTransaction", req, " mongo err:", err, "clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.PocketDoesNotExists, http.StatusBadRequest)
	}

	var pocketInfo models.PocketInfo
	pocketInfo.PocketID = pockets.PocketId
	pocketInfo.PocketName = pockets.PocketName
	pocketInfo.PocketImage = pockets.PocketImage
	pocketInfo.PocketWebImage = pockets.PocketWebImage
	pocketInfo.PrimaryBackgroundColor = pockets.PrimaryBackgroundColor
	pocketInfo.PrimarySecondaryColor = pockets.PrimarySecondaryColor
	pocketInfo.TransactionDate = helpers.GetCurrentTimeInIST()
	pocketInfo.PocketExecutedPrice = req.PocketPrice
	pocketInfo.TransactionID = uuid.New().String()
	if req.OrderSide == strings.ToUpper(constants.BUY) {
		pocketInfo.TransactionStatus = 0
	} else { // SELL
		pocketInfo.TransactionStatus = 1
	}

	if req.Qty == 0 {
		pocketInfo.Qty = 1
	} else {
		pocketInfo.Qty = req.Qty
	}

	pocketTransactionComplete.ClientId = req.ClientId
	pocketTransactionComplete.PocketCounter = pocketCounter
	pocketTransactionComplete.AllPocketPurchases = append(pocketTransactionComplete.AllPocketPurchases, pocketInfo)

	filter := bson.D{{"clientId", req.ClientId}}
	update := bson.D{{"$set", pocketTransactionComplete}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.POCKETSTRANSACTIONS, filter, update, opts)
	if err != nil {
		loggerconfig.Error("StorePocketTransaction Mongo Upsert failed error =", err, " clientId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("StorePocketTransaction api success uccId: ", req.ClientId, " requestId: ", reqH.RequestId)

	apiRes.Data = nil
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}

func (obj ExecutePocketV2Obj) BuyPocketV2(req models.ExecutePocketV2Request, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	loggerconfig.Info("BuyPocketV2 req: ", req, " clientId: ", req.ClientId, " reqH: ", reqH)

	return executePocketV2(req, obj, constants.BUY, reqH)
}

func executePocketV2(req models.ExecutePocketV2Request, obj ExecutePocketV2Obj, orderSide string, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	var apiRes apihelpers.APIRes

	var pockets models.MongoPockets
	err := dbops.MongoRepo.FindOne(constants.POCKETSCOLLECTION, bson.M{"pocketId": req.PocketId}, &pockets)

	if err != nil && err.Error() == constants.MongoNoDocError {
		loggerconfig.Error("ExecutePocketV2 req: ", req, " mongo err:", err, "clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.PocketDoesNotExists, http.StatusBadRequest)
	}

	mapPocketTokenMetaData := make(map[string]models.PocketsMetaData)
	for i := 0; i < len(pockets.PocketTokens); i++ {
		mapPocketTokenMetaData[pockets.PocketTokens[i].Token] = pockets.PocketTokens[i]
	}

	var allPlaceOrderReq []models.PlaceOrderRequest
	for i := 0; i < len(pockets.PocketTokens); i++ {
		var placeOrderReq models.PlaceOrderRequest
		placeOrderReq.ClientID = req.ClientId
		placeOrderReq.DisclosedQuantity = 0
		placeOrderReq.Exchange = pockets.PocketTokens[i].Exchange
		placeOrderReq.ExecutionType = strings.ToUpper(constants.REGULAR)
		placeOrderReq.InstrumentToken = pockets.PocketTokens[i].Token
		placeOrderReq.OrderType = strings.ToUpper(constants.MARKET)
		if strings.EqualFold(orderSide, constants.BUY) {
			placeOrderReq.OrderSide = strings.ToUpper(constants.BUY)
		} else {
			placeOrderReq.OrderSide = strings.ToUpper(constants.SELL)
		}
		placeOrderReq.Price = 0.0
		placeOrderReq.Product = strings.ToUpper(constants.CNC)
		qty, _ := strconv.Atoi(pockets.PocketTokens[i].Qty)
		placeOrderReq.Quantity = qty * req.LotSize
		placeOrderReq.TriggerPrice = 0.0
		placeOrderReq.Validity = strings.ToUpper(constants.IOC) // complete or reject

		allPlaceOrderReq = append(allPlaceOrderReq, placeOrderReq)
	}

	var orderIdArr []string

	orderObj := tradelab.InitOrder(obj.redisCli)
	for i := 0; i < len(allPlaceOrderReq); i++ {
		status, res := tradelab.OrderObj.PlaceOrder(orderObj, allPlaceOrderReq[i], reqH)
		if status != http.StatusOK {
			loggerconfig.Error("ExecutePocketV2 in PlaceOrder status != 200", status, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		placeOrderRes, ok := res.Data.(models.PlaceOrderResponse)
		if !ok {
			loggerconfig.Error("ExecutePocketV2 in PlaceOrder interface parsing error", ok, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		orderIdArr = append(orderIdArr, placeOrderRes.OmsOrderID)
	}

	var completedOrderReq models.CompletedOrderRequest
	completedOrderReq.ClientID = req.ClientId
	completedOrderReq.Type = strings.ToLower(constants.Completed)
	time.Sleep(constants.PocketOrderSleepSeconds * time.Second) // sleeping for a second because tradelab give above order update in completed order api after some time
	statusCompletedOrder, resCompletedOrder := tradelab.OrderObj.CompletedOrder(orderObj, completedOrderReq, reqH)
	if statusCompletedOrder != http.StatusOK {
		loggerconfig.Error("ExecutePocketV2 in CompletedOrder status != 200", statusCompletedOrder, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	completedOrderRes, ok := resCompletedOrder.Data.(models.CompletedOrderResponse)
	if !ok {
		loggerconfig.Error("ExecutePocketV2 in CompletedOrder interface parsing error", ok, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	mapOrderIdStatus := make(map[string]models.CompletedOrderResponseOrders)
	for i := 0; i < len(completedOrderRes.Orders); i++ {
		mapOrderIdStatus[completedOrderRes.Orders[i].OmsOrderID] = completedOrderRes.Orders[i]
	}

	var pocketTransactionStoreV2 models.PocketTransactionStoreV2
	err = dbops.MongoRepo.FindOne(constants.POCKETSTRANSACTIONSDetails, bson.M{"clientid": req.ClientId}, &pocketTransactionStoreV2)

	if err != nil && err.Error() != constants.MongoNoDocError {
		return apihelpers.SendErrorResponse(false, constants.InvalidUserIdOrPass, http.StatusBadRequest)
	}

	if err != nil && err.Error() == constants.MongoNoDocError {
		pocketTransactionStoreV2.ClientId = req.ClientId
	}

	var pocketDetailsV2 models.PocketDetailsV2
	pocketDetailsV2.PocketTransactionId = uuid.New().String()
	pocketDetailsV2.PocketID = req.PocketId
	if strings.EqualFold(orderSide, constants.BUY) {
		pocketDetailsV2.TransactionStatus = 0
	} else {
		pocketDetailsV2.TransactionStatus = 1
	}
	pocketDetailsV2.LotSize = req.LotSize

	var orderCompleted []models.PocketsMetaData
	var orderCancelled []models.PocketsMetaData

	orderCompletedPrice := 0.0

	for i := 0; i < len(orderIdArr); i++ {
		completedOrderResOrders, orderPresent := mapOrderIdStatus[orderIdArr[i]]
		// If the key exists
		if orderPresent && completedOrderResOrders.OrderStatus == strings.ToUpper(constants.Complete) {
			orderCompletedPrice += completedOrderResOrders.Price * float64(completedOrderResOrders.Quantity)
			pocketData, pocketPresent := mapPocketTokenMetaData[strconv.Itoa(completedOrderResOrders.InstrumentToken)]
			if !pocketPresent {
				loggerconfig.Error("ExecutePocketV2 in mapPocketTokenMetaData completed order details don't present in pocket with token", strconv.Itoa(completedOrderResOrders.InstrumentToken), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
			}
			orderCompleted = append(orderCompleted, pocketData)
		} else if orderPresent && (completedOrderResOrders.OrderStatus == strings.ToUpper(constants.CancelConfirmed) || completedOrderResOrders.OrderStatus == strings.ToUpper(constants.Rejected)) {
			pocketData, pocketPresent := mapPocketTokenMetaData[strconv.Itoa(completedOrderResOrders.InstrumentToken)]
			if !pocketPresent {
				loggerconfig.Error("ExecutePocketV2 in mapPocketTokenMetaData cancelled order details don't present in pocket with token", strconv.Itoa(completedOrderResOrders.InstrumentToken), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
			}
			orderCancelled = append(orderCancelled, pocketData)
		} else {
			loggerconfig.Error("ExecutePocketV2 unknown status code for order: ", helpers.LogStructAsJSON(completedOrderResOrders))
		}
	}

	pocketDetailsV2.OrderCompleted = orderCompleted
	pocketDetailsV2.OrderCancelled = orderCancelled
	pocketDetailsV2.OrderCompletedPrice = orderCompletedPrice

	pocketTransactionStoreV2.AllPocketPurchages = append(pocketTransactionStoreV2.AllPocketPurchages, pocketDetailsV2)

	filter := bson.D{{"clientid", req.ClientId}}
	update := bson.D{{"$set", pocketTransactionStoreV2}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.POCKETSTRANSACTIONSDetails, filter, update, opts)
	if err != nil {
		loggerconfig.Error("ExecutePocketV2 pocketTransactionStoreV2 Mongo Upsert failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	apiRes.Data = pocketDetailsV2
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ExecutePocketV2Obj) ExitPocketV2(req models.ExecutePocketV2Request, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	loggerconfig.Info("ExitPocketV2 req: ", req, " clientId: ", req.ClientId, " reqH: ", reqH)

	return executePocketV2(req, obj, constants.SELL, reqH)

}

func (obj ExecutePocketV2Obj) FetchPocketPortfolioV2(req models.FetchPocketPortfolioRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	var fetchPocketPortfolioResponse models.FetchPocketPortfolioResponse
	fetchPocketPortfolioResponse.ClientId = req.ClientId

	var pocketTransactionStoreV2 models.PocketTransactionStoreV2
	err := dbops.MongoRepo.FindOne(constants.POCKETSTRANSACTIONSDetails, bson.M{"clientid": req.ClientId}, &pocketTransactionStoreV2)

	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("FetchPocketPortfolioV2 fetch failed mongo err:", err, "clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.InvalidUserId, http.StatusBadRequest)
	}

	mapIdPocket := make(map[string][]models.PocketDetailsV2)
	for i := 0; i < len(pocketTransactionStoreV2.AllPocketPurchages); i++ {
		id := pocketTransactionStoreV2.AllPocketPurchages[i].PocketID
		mapIdPocket[id] = append(mapIdPocket[id], pocketTransactionStoreV2.AllPocketPurchages[i])
	}

	for keys := range mapIdPocket {
		var pocket models.Pocket
		pocket.PocketID = keys
		totalBuyPockets := 0
		totalSellPockets := 0
		totalPrice := 0.0
		totalBuyPrice := 0.0
		totalSellPrice := 0.0
		for i := 0; i < len(mapIdPocket[keys]); i++ {
			if len(mapIdPocket[keys][i].OrderCancelled) > 0 {
				continue
			}
			pocketQty := mapIdPocket[keys][i].LotSize
			if mapIdPocket[keys][i].TransactionStatus == 0 { // Buy
				totalBuyPockets += pocketQty
				totalPrice += mapIdPocket[keys][i].OrderCompletedPrice
				totalBuyPrice += mapIdPocket[keys][i].OrderCompletedPrice
			} else {
				totalSellPockets += pocketQty
				totalPrice -= mapIdPocket[keys][i].OrderCompletedPrice * (float64)(pocketQty) // Sell
				totalSellPrice += mapIdPocket[keys][i].OrderCompletedPrice * (float64)(pocketQty)
			}
		}
		if totalBuyPockets == totalSellPockets {
			continue // don't send this pocket to frontend
		}
		pocket.TotalBuyPockets = totalBuyPockets
		pocket.TotalSellPockets = totalSellPockets
		pocket.AveragePrice = totalBuyPrice / float64(totalBuyPockets)
		pocket.TotalInvestment = pocket.AveragePrice * float64(totalBuyPockets-totalSellPockets)

		// pocket image
		var storedPocket models.MongoPockets
		var err error
		err = dbops.MongoRepo.FindOne(constants.POCKETSCOLLECTION, bson.M{"pocketId": pocket.PocketID}, &storedPocket)
		if err != nil && err.Error() != constants.MongoNoDocError {
			loggerconfig.Error("FetchPocketPortfolioV2 Stored Pocket fetch failed with pocket id:", pocket.PocketID, " mongo err:", err, "clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
			return apihelpers.SendErrorResponse(false, constants.PocketDoesNotExists, http.StatusBadRequest)
		}

		pocket.PocketName = storedPocket.PocketName
		pocket.PocketImage = storedPocket.PocketImage
		pocket.PocketWebImage = storedPocket.PocketWebImage
		pocket.PrimaryBackgroundColor = storedPocket.PrimaryBackgroundColor
		pocket.PrimarySecondaryColor = storedPocket.PrimarySecondaryColor

		fetchPocketPortfolioResponse.PortfolioDetails = append(fetchPocketPortfolioResponse.PortfolioDetails, pocket)
	}

	loggerconfig.Info("FetchPocketPortfolioV2 response is = ", helpers.LogStructAsJSON(fetchPocketPortfolioResponse), " uccId:", req.ClientId, " requestId:", reqH.RequestId)

	apiRes.Data = fetchPocketPortfolioResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}
