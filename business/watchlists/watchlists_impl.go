package watchlists

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	apihelpers "space/apiHelpers"
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
)

var stockMap sync.Map

func (obj WatchlistsObj) PopulateIsinMappingInLocalCache() {
	loggerconfig.Info("PopulateIsinMappingInLocalCache started...")

	// clear existing map if any key is there (only if this method is called elsewhere)
	stockMap.Range(func(key, value interface{}) bool {
		stockMap.Delete(key)
		return true
	})

	// Get the current time in the Asia/Kolkata timezone
	currentTime := helpers.GetCurrentTimeInIST()

	data, err := obj.contractCacheCli.GetAllFromHashWithPipeline("isin_data", 1000)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, PopulateIsinMappingInLocalCache: Error while GetAllFromHash for bucket isin_data, err ", err)
		return
	}

	// Loop over each key-value pair
	for key, val := range data {
		stockDetail := models.ContractDetails{}

		// Unmarshal JSON data into stockDetail
		if err = json.Unmarshal([]byte(val), &stockDetail); err != nil {
			loggerconfig.Info("Alert Severity:P1-High, PopulateIsinMappingInLocalCache: Error Unmarshalling stockmetadata for key ", key, " err ", err)
			continue // Skip this key-value pair and continue with the next one
		}
		var stock models.StockDetailsV2
		stock.IsinStockId = key
		stock.Isin = stockDetail.Isin
		stock.Exchange = stockDetail.Exchange
		stock.DisplayName = stockDetail.TradingSymbol
		stock.Symbol = stockDetail.Symbol
		stock.Company = stockDetail.Name
		stock.IsTradable = stockDetail.IsTradable
		stock.Token = stockDetail.Token1
		stock.TradingSymbol = stockDetail.TradingSymbol
		stock.StockId = stockDetail.Exchange + "-" + stockDetail.Token1

		if strings.EqualFold(stock.Exchange, constants.NSE) || strings.EqualFold(stock.Exchange, constants.BSE) {
			stock.Segment = constants.SegmentEquity
		}

		if stockDetail.Series == constants.SegmentIndices {
			stock.Segment = constants.SegmentIndex
		}

		if strings.EqualFold(stock.Exchange, constants.MCX) {
			stock.Segment = constants.SegmentCommodity
		}

		var stockDetailsWithExpiry models.StockDetailsWithExpiry

		stockDetailsWithExpiry.Expiry = helpers.Next830AMUnix(currentTime)
		stockDetailsWithExpiry.Stock = stock

		// Store in the sync.Map
		stockMap.Store(stock.StockId, stockDetailsWithExpiry)
	}

	size := 0
	stockMap.Range(func(_, _ interface{}) bool {
		size++
		return true
	})

	loggerconfig.Info("PopulateIsinMappingInLocalCache Completed... the size of the map : ", size)
}

type WatchlistsObj struct {
	contractCacheCli cache.ContractCache
	mongodb          db.MongoDatabase
	displayNameCheck bool
}

func InitWatchlists(mongodb db.MongoDatabase, contractCacheCli cache.ContractCache) WatchlistsObj {
	defer models.HandlePanic()
	watchlists := WatchlistsObj{
		mongodb:          mongodb,
		contractCacheCli: contractCacheCli,
		displayNameCheck: constants.CheckDisplayNameFlag,
	}

	return watchlists
}

func (obj WatchlistsObj) CreateWatchList(req models.CreateWatchListRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var watchLists models.MongoWatchLists
	var err error
	err = dbops.MongoRepo.FindOne(constants.WATCHLISTSCOLLECTION, bson.M{"watchListName": req.WatchListName, "clientId": req.ClientId}, &watchLists)

	var apiRes apihelpers.APIRes
	// if pocket already exists
	if err == nil && watchLists.WatchListName != "" {
		return apihelpers.SendErrorResponse(false, constants.WatchListsAlreadyExists, http.StatusBadRequest)
	}

	id := uuid.New().String()

	mongoWatchListDetails := &models.MongoWatchLists{
		ClientId:               req.ClientId,
		WatchListName:          req.WatchListName,
		WatchListNameLongDesc:  req.WatchListNameLongDesc,
		WatchListNameShortDesc: req.WatchListNameShortDesc,
		WatchListId:            id,
	}

	filter := bson.D{{"watchListId", id}}
	update := bson.D{{"$set", mongoWatchListDetails}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.WATCHLISTSCOLLECTION, filter, update, opts)
	if err != nil {
		loggerconfig.Error("CreateWatchList Mongo Upsert failed error =", err, " clientId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	var resp models.CreateWatchListResponse
	resp.WatchListId = id

	loggerconfig.Info("CreateWatchList Successful, response:", helpers.LogStructAsJSON(resp), "clientID: ", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
	apiRes.Data = resp
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj WatchlistsObj) ModifyWatchList(req models.ModifyWatchListRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var watchLists models.MongoWatchLists
	var err error
	err = dbops.MongoRepo.FindOne(constants.POCKETSCOLLECTION, bson.M{"WatchListId": req.WatchListId}, &watchLists)

	var apiRes apihelpers.APIRes
	// if pocket does not exists
	if err != nil && watchLists.WatchListName == "" {
		return apihelpers.SendErrorResponse(false, constants.WatchListsDoesNotExists, http.StatusBadRequest)
	}

	mongoWatchListDetails := &models.MongoWatchLists{
		WatchListName:          req.WatchListName,
		WatchListNameLongDesc:  req.WatchListNameLongDesc,
		WatchListNameShortDesc: req.WatchListNameShortDesc,
		WatchListId:            req.WatchListId,
	}

	filter := bson.D{{"watchListId", req.WatchListId}}
	update := bson.D{{"$set", mongoWatchListDetails}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.WATCHLISTSCOLLECTION, filter, update, opts)
	if err != nil {
		loggerconfig.Error("ModifyWatchList Mongo Upsert failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	var resp models.CreateWatchListResponse
	resp.WatchListId = req.WatchListId

	loggerconfig.Info("ModifyWatchList Successful, response:", helpers.LogStructAsJSON(resp), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
	apiRes.Data = resp
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj WatchlistsObj) FetchWatchLists(req models.FetchWatchListsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var err error
	cur, err := dbops.MongoRepo.Find(constants.WATCHLISTSCOLLECTION, bson.M{"clientId": req.ClientId})
	if err != nil {
		loggerconfig.Error("FetchWatchLists Mongo Find() failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var apiRes apihelpers.APIRes

	var resp models.FetchWatchListResponse
	for cur.Next(context.Background()) {
		var mongoWatchLists models.MongoWatchLists
		err := cur.Decode(&mongoWatchLists)
		if err != nil {
			loggerconfig.Error("FetchWatchLists Mongo Cursor Parsing failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
			return apihelpers.SendInternalServerError()
		}
		var metaData models.WatchListMetadata
		metaData.WatchListId = mongoWatchLists.WatchListId
		metaData.WatchListName = mongoWatchLists.WatchListName
		metaData.WatchListNameLongDesc = mongoWatchLists.WatchListNameLongDesc
		metaData.WatchListNameShortDesc = mongoWatchLists.WatchListNameShortDesc
		resp.WatchListIds = append(resp.WatchListIds, metaData)
	}
	loggerconfig.Info("FetchWatchLists Successful, response:", helpers.LogStructAsJSON(resp), "clientID: ", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
	apiRes.Data = resp
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj WatchlistsObj) DeleteWatchList(req models.DeleteWatchListRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var watchLists models.MongoWatchLists
	var err error
	err = dbops.MongoRepo.FindOne(constants.WATCHLISTSCOLLECTION, bson.M{"clientId": req.ClientId, "watchListId": req.WatchListId}, &watchLists)

	var apiRes apihelpers.APIRes
	// if pocket does not exists
	if err != nil && watchLists.WatchListName == "" {
		return apihelpers.SendErrorResponse(false, constants.WatchListsDoesNotExists, http.StatusBadRequest)
	}

	filter := bson.D{{"WatchListId", req.WatchListId}}
	opts := options.Delete()
	_, err = dbops.MongoRepo.DeleteOne(constants.WATCHLISTSCOLLECTION, filter, opts)
	if err != nil {
		loggerconfig.Error("DeleteWatchList Mongo DeleteOne() failed error =", err, " clientId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj WatchlistsObj) FetchWatchListDetails(req models.FetchWatchListsDetailsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var err error
	cur, err := dbops.MongoRepo.Find(constants.WATCHLISTSTOCKSCOLLECTION, bson.M{"clientId": req.ClientId, "watchListId": req.WatchListId})
	if err != nil {
		loggerconfig.Error("FetchWatchListDetails Mongo FIND() failed error =", err, " clientId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	var apiRes apihelpers.APIRes

	var resp models.FetchWatchListsDetailsResponse
	for cur.Next(context.Background()) {
		var mongoWatchLists models.MongoStocksWatchLists
		err := cur.Decode(&mongoWatchLists)
		if err != nil {
			loggerconfig.Error("FetchWatchListDetails Mongo Cursor Parsing failed error =", err, " clientId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
			return apihelpers.SendInternalServerError()
		}
		var metaData models.StockDetails
		metaData.Company = mongoWatchLists.Stock.Company
		metaData.DisplayName = mongoWatchLists.Stock.DisplayName
		metaData.Exchange = mongoWatchLists.Stock.Exchange
		metaData.Expiry = mongoWatchLists.Stock.Expiry
		metaData.IsTradable = mongoWatchLists.Stock.IsTradable
		metaData.Isin = mongoWatchLists.Stock.Isin
		metaData.Segment = mongoWatchLists.Stock.Segment
		metaData.Symbol = mongoWatchLists.Stock.Symbol
		metaData.Token = mongoWatchLists.Stock.Token
		metaData.TradingSymbol = mongoWatchLists.Stock.TradingSymbol
		resp.Stocks = append(resp.Stocks, metaData)

	}
	loggerconfig.Info("FetchWatchListDetails Successful, response:", helpers.LogStructAsJSON(resp), "clientID: ", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
	apiRes.Data = resp
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj WatchlistsObj) AddStockToWatchList(req models.AddStockToWatchListsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var stockLists models.AddStockToWatchListsRequest
	var err error
	err = dbops.MongoRepo.FindOne(constants.WATCHLISTSTOCKSCOLLECTION, bson.M{"clientId": req.ClientId, "watchListId": req.WatchListId}, &stockLists)

	var apiRes apihelpers.APIRes
	if err == nil && stockLists.WatchListId != "" {
		return apihelpers.SendErrorResponse(false, constants.WatchListsAlreadyExists, http.StatusBadRequest)
	}

	mongoStockListDetails := &models.MongoStocksWatchLists{
		ClientId:    req.ClientId,
		WatchListId: req.WatchListId,
		Stock:       req.Stock,
	}

	filter := bson.M{"clientId": req.ClientId, "watchListId": req.WatchListId}
	update := bson.D{{"$set", mongoStockListDetails}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.WATCHLISTSTOCKSCOLLECTION, filter, update, opts)
	if err != nil {
		loggerconfig.Error("AddStockToWatchList Mongo Upsert failed error =", err, " clientId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj WatchlistsObj) FetchWatchListsV2(req models.FetchWatchListV2Request, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	var stockLists models.MongoNewWatchLists
	var err error
	err = dbops.MongoRepo.FindOne(constants.WATCHLISTSTOCKSCOLLECTIONNEW, bson.M{"clientId": req.ClientId}, &stockLists)
	var apiRes apihelpers.APIRes
	if err != nil && err.Error() == constants.MongoNoDocError {
		stockLists.ClientId = req.ClientId
		filter := bson.D{{"clientId", req.ClientId}}
		update := bson.D{{"$set", stockLists}}
		opts := options.Update().SetUpsert(true)
		err = dbops.MongoRepo.UpdateOne(constants.WATCHLISTSTOCKSCOLLECTIONNEW, filter, update, opts)
		if err != nil {
			loggerconfig.Error("FetchWatchListsV2 Mongo Upsert failed error =", err, " clientId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
			return apihelpers.SendInternalServerError()
		}
	}

	if err != nil {
		loggerconfig.Error("FetchWatchListsV2  error =", err, " clientId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	var resp models.FetchWatchListV2Response

	resp.ClientId = req.ClientId

	var watchLists []models.WatchListsDetails

	var watchList models.WatchListsDetails
	watchList.WatchListId = "wl1"
	watchList.Stocks = stockLists.WatchList1
	watchLists = append(watchLists, watchList)

	watchList.WatchListId = "wl2"
	watchList.Stocks = stockLists.WatchList2
	watchLists = append(watchLists, watchList)

	watchList.WatchListId = "wl3"
	watchList.Stocks = stockLists.WatchList3
	watchLists = append(watchLists, watchList)

	watchList.WatchListId = "wl4"
	watchList.Stocks = stockLists.WatchList4
	watchLists = append(watchLists, watchList)

	watchList.WatchListId = "wl5"
	watchList.Stocks = stockLists.WatchList5
	watchLists = append(watchLists, watchList)

	resp.WatchLists = watchLists
	loggerconfig.Info("FetchWatchListsV2 Successful, response:", helpers.LogStructAsJSON(resp), "clientID: ", req.ClientId, " requestId:", reqH.RequestId, " deviceType: ", reqH.DeviceType, " clientVersion:", reqH.ClientVersion)
	apiRes.Data = resp
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}

func (obj WatchlistsObj) AddStockToWatchListV2(req models.AddStockToWatchListV2Request, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var stockLists models.MongoNewWatchLists
	var err error
	err = dbops.MongoRepo.FindOne(constants.WATCHLISTSTOCKSCOLLECTIONNEW, bson.M{"clientId": req.ClientId}, &stockLists)
	var apiRes apihelpers.APIRes
	if err != nil && err.Error() == constants.MongoNoDocError {
		return apihelpers.SendErrorResponse(false, constants.WatchListsDoesNotExists, http.StatusBadRequest)
	}

	stockId := req.Stock.Exchange + "-" + req.Stock.Token

	// display name check
	if obj.displayNameCheck {
		// a unique key for stock
		stockUniqueKey := strings.ToUpper(req.Stock.Exchange) + "_" + req.Stock.Token
		err, val := obj.contractCacheCli.GetFromHash("stock_key", stockUniqueKey)
		if err != nil {
			loggerconfig.Error("Alert Severity:P1-High, AddStockToWatchListV2: Error while GetFromHashSetNew script v2 for query ", stockId, " and reqId ", reqH.RequestId, " and client Version:", reqH.ClientVersion, " and err ", err)
			return apihelpers.SendInternalServerError()
		}

		stockDetail := models.ContractDetails{}
		if err = json.Unmarshal([]byte(val), &stockDetail); err != nil {
			loggerconfig.Info("Alert Severity:P1-High, AddStockToWatchListV2: Error Unmarshalling stockmetadata for key  err ", err, " and reqId:", reqH.RequestId, " and client Version:", reqH.ClientVersion)
			return apihelpers.SendInternalServerError()
		}

		validStockDisplayName := strings.Contains(strings.ToUpper(stockDetail.Name), strings.ToUpper(req.Stock.DisplayName))

		if !validStockDisplayName {
			apiRes.Status = false
			apiRes.Message = constants.ErrorCodeMap[constants.InvalidDisplayName]
			apiRes.ErrorCode = constants.InvalidDisplayName
			return http.StatusBadRequest, apiRes
		}
	}

	duplicateStock := false

	for i := 0; i < len(req.WatchListId); i++ {

		if req.WatchListId[i] == "wl1" {
			duplicateStock = duplicateStockInWatchlist(stockId, stockLists.WatchList1)
			if duplicateStock {
				break
			}
			req.Stock.StockId = stockId
			stockLists.WatchList1 = append([]models.StockDetails{req.Stock}, stockLists.WatchList1...) // prepending the new stock in watchlist
		}

		if req.WatchListId[i] == "wl2" {
			duplicateStock = duplicateStockInWatchlist(stockId, stockLists.WatchList2)
			if duplicateStock {
				break
			}
			req.Stock.StockId = stockId
			stockLists.WatchList2 = append([]models.StockDetails{req.Stock}, stockLists.WatchList2...) // prepending the new stock in watchlist
		}

		if req.WatchListId[i] == "wl3" {
			duplicateStock = duplicateStockInWatchlist(stockId, stockLists.WatchList3)
			if duplicateStock {
				break
			}
			req.Stock.StockId = stockId
			stockLists.WatchList3 = append([]models.StockDetails{req.Stock}, stockLists.WatchList3...) // prepending the new stock in watchlist
		}

		if req.WatchListId[i] == "wl4" {
			duplicateStock = duplicateStockInWatchlist(stockId, stockLists.WatchList4)
			if duplicateStock {
				break
			}
			req.Stock.StockId = stockId
			stockLists.WatchList4 = append([]models.StockDetails{req.Stock}, stockLists.WatchList4...) // prepending the new stock in watchlist
		}

		if req.WatchListId[i] == "wl5" {
			duplicateStock = duplicateStockInWatchlist(stockId, stockLists.WatchList5)
			if duplicateStock {
				break
			}
			req.Stock.StockId = stockId
			stockLists.WatchList5 = append([]models.StockDetails{req.Stock}, stockLists.WatchList5...) // prepending the new stock in watchlist
		}

	}

	if duplicateStock {
		loggerconfig.Error("AddStockToWatchListV2 stock already present in watchlist", " clientId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Status = false
		apiRes.Message = constants.ErrorCodeMap[constants.WatchListsStockAlreadyExists]
		apiRes.ErrorCode = constants.WatchListsStockAlreadyExists
		return http.StatusBadRequest, apiRes
	}

	filter := bson.D{{"clientId", req.ClientId}}
	update := bson.D{{"$set", stockLists}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.WATCHLISTSTOCKSCOLLECTIONNEW, filter, update, opts)
	if err != nil {
		loggerconfig.Error("AddStockToWatchListV2 Mongo Upsert failed error =", err, " clientId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	loggerconfig.Info("AddStockToWatchListV2, stock is added successfully for req packet: ", helpers.LogStructAsJSON(req), "req requestId:", reqH.RequestId, " deviceType: ", reqH.DeviceType, " clientVersion:", reqH.ClientVersion)

	return http.StatusOK, apiRes

}

func duplicateStockInWatchlist(currStockId string, addedStocks []models.StockDetails) bool {
	for i := 0; i < len(addedStocks); i++ {
		if strings.EqualFold(currStockId, addedStocks[i].StockId) {
			return true
		}
	}
	// unique stock
	return false
}

func (obj WatchlistsObj) DeleteStockInWatchListV2(req models.DeleteWatchListV2Request, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var stockLists models.MongoNewWatchLists
	err := dbops.MongoRepo.FindOne(constants.WATCHLISTSTOCKSCOLLECTIONNEW, bson.M{"clientId": req.ClientId}, &stockLists)
	var apiRes apihelpers.APIRes
	if err != nil && err.Error() == constants.MongoNoDocError {
		return apihelpers.SendErrorResponse(false, constants.WatchListsDoesNotExists, http.StatusBadRequest)
	}

	stockIdMap := addStockIdsInMap(req.StockId)

	if req.WatchListId == "wl1" {
		var stockDetails []models.StockDetails
		for i := 0; i < len(stockLists.WatchList1); i++ { // looping on all stocks of current watchlist
			if _, present := stockIdMap[stockLists.WatchList1[i].StockId]; !present { // if block is only executed if current stockid is not present in req packet from frontend.
				stockDetails = append(stockDetails, stockLists.WatchList1[i]) // adding stocks to stockDetails from current watchlist
			}
		}
		stockLists.WatchList1 = stockDetails // reassigning watchlist stocks
	}

	if req.WatchListId == "wl2" {
		var stockDetails []models.StockDetails
		for i := 0; i < len(stockLists.WatchList2); i++ { // looping on all stocks of current watchlist
			if _, present := stockIdMap[stockLists.WatchList2[i].StockId]; !present { // if block is only executed if current stockid is not present in req packet from frontend.
				stockDetails = append(stockDetails, stockLists.WatchList2[i]) // adding stocks to stockDetails from current watchlist
			}
		}
		stockLists.WatchList2 = stockDetails // reassigning watchlist stocks
	}

	if req.WatchListId == "wl3" {
		var stockDetails []models.StockDetails
		for i := 0; i < len(stockLists.WatchList3); i++ { // looping on all stocks of current watchlist
			if _, present := stockIdMap[stockLists.WatchList3[i].StockId]; !present { // if block is only executed if current stockid is not present in req packet from frontend.
				stockDetails = append(stockDetails, stockLists.WatchList3[i]) // adding stocks to stockDetails from current watchlist
			}
		}
		stockLists.WatchList3 = stockDetails // reassigning watchlist stocks
	}

	if req.WatchListId == "wl4" {
		var stockDetails []models.StockDetails
		for i := 0; i < len(stockLists.WatchList4); i++ { // looping on all stocks of current watchlist
			if _, present := stockIdMap[stockLists.WatchList4[i].StockId]; !present { // if block is only executed if current stockid is not present in req packet from frontend.
				stockDetails = append(stockDetails, stockLists.WatchList4[i]) // adding stocks to stockDetails from current watchlist
			}
		}
		stockLists.WatchList4 = stockDetails // reassigning watchlist stocks
	}

	if req.WatchListId == "wl5" {
		var stockDetails []models.StockDetails
		for i := 0; i < len(stockLists.WatchList5); i++ { // looping on all stocks of current watchlist
			if _, present := stockIdMap[stockLists.WatchList5[i].StockId]; !present { // if block is only executed if current stockid is not present in req packet from frontend.
				stockDetails = append(stockDetails, stockLists.WatchList5[i]) // adding stocks to stockDetails from current watchlist
			}
		}
		stockLists.WatchList5 = stockDetails // reassigning watchlist stocks
	}

	filter := bson.D{{"clientId", req.ClientId}}
	update := bson.D{{"$set", stockLists}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.WATCHLISTSTOCKSCOLLECTIONNEW, filter, update, opts)
	if err != nil {
		loggerconfig.Error("DeleteStockInWatchListV2 Mongo Upsert failed error =", err, " clientId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	loggerconfig.Info("DeleteStockInWatchListV2, stock is deleted successfully for req packet: ", helpers.LogStructAsJSON(req), "req requestId:", reqH.RequestId, " deviceType: ", reqH.DeviceType, " clientVersion:", reqH.ClientVersion)
	return http.StatusOK, apiRes

}

func (obj WatchlistsObj) DeleteStockInWatchListV2Updated(req models.DeleteWatchListV2UpdatedRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var stockLists models.MongoNewWatchLists
	err := dbops.MongoRepo.FindOne(constants.WATCHLISTSTOCKSCOLLECTIONNEW, bson.M{"clientId": req.ClientId}, &stockLists)
	var apiRes apihelpers.APIRes
	if err != nil && err.Error() == constants.MongoNoDocError {
		return apihelpers.SendErrorResponse(false, constants.WatchListsDoesNotExists, http.StatusBadRequest)
	}

	var stockIdMap map[string]int

	for k := 0; k < len(req.WatchListId); k++ {
		switch req.WatchListId[k].WatchlistId {
		case "wl1":
			stockIdMap = addStockIdsInMap(req.WatchListId[k].StockId)
			noOfStocksInWL := len(req.WatchListId[k].StockId)
			var stockDetails []models.StockDetails
			for i := 0; i < len(stockLists.WatchList1); i++ {
				if _, present := stockIdMap[stockLists.WatchList1[i].StockId]; !present {
					stockDetails = append(stockDetails, stockLists.WatchList1[i])
				} else {
					noOfStocksInWL--
				}
			}
			if noOfStocksInWL > 0 {
				return apihelpers.SendErrorResponse(false, constants.InvalidRequest, http.StatusBadRequest)
			}
			stockLists.WatchList1 = stockDetails
		case "wl2":
			stockIdMap = addStockIdsInMap(req.WatchListId[k].StockId)
			noOfStocksInWL := len(req.WatchListId[k].StockId)
			var stockDetails []models.StockDetails
			for i := 0; i < len(stockLists.WatchList2); i++ {
				if _, present := stockIdMap[stockLists.WatchList2[i].StockId]; !present {
					stockDetails = append(stockDetails, stockLists.WatchList2[i])
				} else {
					noOfStocksInWL--
				}
			}
			if noOfStocksInWL > 0 {
				return apihelpers.SendErrorResponse(false, constants.InvalidRequest, http.StatusBadRequest)
			}
			stockLists.WatchList2 = stockDetails
		case "wl3":
			stockIdMap = addStockIdsInMap(req.WatchListId[k].StockId)
			noOfStocksInWL := len(req.WatchListId[k].StockId)
			var stockDetails []models.StockDetails
			for i := 0; i < len(stockLists.WatchList3); i++ {
				if _, present := stockIdMap[stockLists.WatchList3[i].StockId]; !present {
					stockDetails = append(stockDetails, stockLists.WatchList3[i])
				} else {
					noOfStocksInWL--
				}
			}
			if noOfStocksInWL > 0 {
				return apihelpers.SendErrorResponse(false, constants.InvalidRequest, http.StatusBadRequest)
			}
			stockLists.WatchList3 = stockDetails
		case "wl4":
			stockIdMap = addStockIdsInMap(req.WatchListId[k].StockId)
			noOfStocksInWL := len(req.WatchListId[k].StockId)
			var stockDetails []models.StockDetails
			for i := 0; i < len(stockLists.WatchList4); i++ {
				if _, present := stockIdMap[stockLists.WatchList4[i].StockId]; !present {
					stockDetails = append(stockDetails, stockLists.WatchList4[i])
				} else {
					noOfStocksInWL--
				}
			}
			if noOfStocksInWL > 0 {
				return apihelpers.SendErrorResponse(false, constants.InvalidRequest, http.StatusBadRequest)
			}
			stockLists.WatchList4 = stockDetails
		case "wl5":
			stockIdMap = addStockIdsInMap(req.WatchListId[k].StockId)
			noOfStocksInWL := len(req.WatchListId[k].StockId)
			var stockDetails []models.StockDetails
			for i := 0; i < len(stockLists.WatchList5); i++ {
				if _, present := stockIdMap[stockLists.WatchList5[i].StockId]; !present {
					stockDetails = append(stockDetails, stockLists.WatchList5[i])
				} else {
					noOfStocksInWL--
				}
			}
			if noOfStocksInWL > 0 {
				return apihelpers.SendErrorResponse(false, constants.InvalidRequest, http.StatusBadRequest)
			}
			stockLists.WatchList5 = stockDetails
		default:
			return apihelpers.SendErrorResponse(false, constants.InvalidRequest, http.StatusBadRequest)
		}
	}

	filter := bson.D{{"clientId", req.ClientId}}
	update := bson.D{{"$set", stockLists}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.WATCHLISTSTOCKSCOLLECTIONNEW, filter, update, opts)
	if err != nil {
		loggerconfig.Error("DeleteStockInWatchListV2 Mongo Upsert failed error =", err, " clientId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}

func (obj WatchlistsObj) ArrangeStocksWatchListV2(req models.ArrangeStocksWatchListV2Request, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	var stockLists models.MongoNewWatchLists
	var err error
	err = dbops.MongoRepo.FindOne(constants.WATCHLISTSTOCKSCOLLECTIONNEW, bson.M{"clientId": req.ClientId}, &stockLists)
	if err != nil && err.Error() == constants.MongoNoDocError {
		loggerconfig.Error("ArrangeStocksWatchListV2 Failed to locate docs in Mongo corresponding to clientId:", req.ClientId, " requestid=", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendErrorResponse(false, constants.WatchListsDoesNotExists, http.StatusBadRequest)
	}

	if req.WatchListId == "wl1" {
		stockIdMap := addStocksInMap(stockLists.WatchList1)
		stockLists.WatchList1 = nil
		for i := 0; i < len(req.StockIds); i++ {
			stockLists.WatchList1 = append(stockLists.WatchList1, stockIdMap[req.StockIds[i]])
		}
	} else if req.WatchListId == "wl2" {
		stockIdMap := addStocksInMap(stockLists.WatchList2)
		stockLists.WatchList2 = nil
		for i := 0; i < len(req.StockIds); i++ {
			stockLists.WatchList2 = append(stockLists.WatchList2, stockIdMap[req.StockIds[i]])
		}
	} else if req.WatchListId == "wl3" {
		stockIdMap := addStocksInMap(stockLists.WatchList3)
		stockLists.WatchList3 = nil
		for i := 0; i < len(req.StockIds); i++ {
			stockLists.WatchList3 = append(stockLists.WatchList3, stockIdMap[req.StockIds[i]])
		}
	} else if req.WatchListId == "wl4" {
		stockIdMap := addStocksInMap(stockLists.WatchList4)
		stockLists.WatchList4 = nil
		for i := 0; i < len(req.StockIds); i++ {
			stockLists.WatchList4 = append(stockLists.WatchList4, stockIdMap[req.StockIds[i]])
		}
	} else if req.WatchListId == "wl5" {
		stockIdMap := addStocksInMap(stockLists.WatchList5)
		stockLists.WatchList5 = nil
		for i := 0; i < len(req.StockIds); i++ {
			stockLists.WatchList5 = append(stockLists.WatchList5, stockIdMap[req.StockIds[i]])
		}
	}

	filter := bson.D{{"clientId", req.ClientId}}
	update := bson.D{{"$set", stockLists}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.WATCHLISTSTOCKSCOLLECTIONNEW, filter, update, opts)
	if err != nil {
		loggerconfig.Error("ArrangeStocksWatchListV2 Failed to update MongoNewWatchLists to mongo, error = ", err, "clientID: ", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("ArrangeStocksWatchListV2 packet=", helpers.LogStructAsJSON(stockLists), " uccId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}

func addStocksInMap(stocks []models.StockDetails) map[string]models.StockDetails {
	stockIdIndexMap := make(map[string]models.StockDetails)

	for _, stock := range stocks {
		stockIdIndexMap[stock.StockId] = stock
	}
	return stockIdIndexMap
}

func removeIndex(s []models.StockDetails, index int) []models.StockDetails {
	return append(s[:index], s[index+1:]...)
}

func addStockIdsInMap(stockIds []string) map[string]int {
	stockIdLocMap := make(map[string]int)

	for i, stockId := range stockIds {
		stockIdLocMap[stockId] = i
	}
	return stockIdLocMap
}

func (obj WatchlistsObj) AddWatchlistStockData(watchlistId string, watchListStocks []models.StockDetailsV2, wg *sync.WaitGroup, resultChan chan<- models.WatchListsDetailsV2, errorChan chan<- error) {
	defer wg.Done()

	var watchList models.WatchListsDetailsV2
	watchList.WatchListId = watchlistId

	for _, stockV2 := range watchListStocks {
		if stockV2.Isin != "" {
			if stockV2.IsinStockId == "" {
				stockV2.IsinStockId = stockV2.Exchange + "-" + stockV2.Isin
			}
			currentTime := helpers.GetCurrentTimeInIST()
			unixCurrentTime := currentTime.Unix()
			value, found := stockMap.Load(stockV2.IsinStockId)

			if found {
				stockDetails := value.(models.StockDetailsWithExpiry)
				if unixCurrentTime < stockDetails.Expiry {
					stockV2 = stockDetails.Stock
					watchList.Stocks = append(watchList.Stocks, stockV2)
					continue
				}
			}

			err, val := obj.contractCacheCli.GetFromHash("isin_data", stockV2.IsinStockId)
			if err != nil {
				loggerconfig.Error("Alert Severity:P2-Mid,FetchWatchListsV3 AddWatchlistStockData: Error while GetFromHashSetNew script v2 for isinStockID ", stockV2.IsinStockId, " and err ", err)
				errorChan <- err
				continue
			}

			stockDetail := models.ContractDetails{}
			if err = json.Unmarshal([]byte(val), &stockDetail); err != nil {
				loggerconfig.Info("Alert Severity:P1-High,FetchWatchListsV3 AddWatchlistStockData: Error Unmarshalling stockmetadata for key  err ", err)
				errorChan <- err
				continue
			}

			stockV2.DisplayName = stockDetail.Symbol
			stockV2.Symbol = stockDetail.Symbol
			stockV2.Company = stockDetail.Name
			stockV2.IsTradable = stockDetail.IsTradable
			stockV2.Token = stockDetail.Token1
			stockV2.TradingSymbol = stockDetail.TradingSymbol

			if strings.EqualFold(stockV2.Exchange, constants.NSE) || strings.EqualFold(stockV2.Exchange, constants.BSE) {
				stockV2.Segment = constants.SegmentEquity
			}

			if stockDetail.Series == constants.SegmentIndices {
				stockV2.Segment = constants.SegmentIndex
			}

			if strings.EqualFold(stockV2.Exchange, constants.MCX) {
				stockV2.Segment = constants.SegmentCommodity
			}

			var stockDetailsWithExpiry models.StockDetailsWithExpiry

			expiryUnixTime := helpers.Next830AMUnix(currentTime)
			stockDetailsWithExpiry.Expiry = expiryUnixTime
			stockDetailsWithExpiry.Stock = stockV2
			stockMap.Store(stockV2.IsinStockId, stockDetailsWithExpiry)
		}
		watchList.Stocks = append(watchList.Stocks, stockV2)
	}

	// Send the result back through the channel
	resultChan <- watchList
}

func (obj WatchlistsObj) FetchWatchListsV3(req models.FetchWatchListV3Request, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	var stockLists models.MongoNewWatchListsV2
	var err error
	err = dbops.MongoRepo.FindOne(constants.WATCHLISTSTOCKSCOLLECTIONNEW, bson.M{"clientId": req.ClientId}, &stockLists)
	var apiRes apihelpers.APIRes
	if err != nil && err.Error() == constants.MongoNoDocError {
		stockLists.ClientId = req.ClientId
		filter := bson.D{{"clientId", req.ClientId}}
		update := bson.D{{"$set", stockLists}}
		opts := options.Update().SetUpsert(true)
		err = dbops.MongoRepo.UpdateOne(constants.WATCHLISTSTOCKSCOLLECTIONNEW, filter, update, opts)
		if err != nil {
			loggerconfig.Error("Alert Severity:P1-High, FetchWatchListsV3 Mongo Upsert failed error =", err, " clientId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
			return apihelpers.SendInternalServerError()
		}
	}

	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, FetchWatchListsV3 mongo error =", err, " clientId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	var resp models.FetchWatchListV3Response

	resp.ClientId = req.ClientId

	watchLists := make([]models.WatchListsDetailsV2, 5)

	var wg sync.WaitGroup

	// Channels for results and errors
	resultChan := make(chan models.WatchListsDetailsV2, 5)
	errChan := make(chan error, 10)

	// Start a goroutine to process the watchList
	wg.Add(5)
	go obj.AddWatchlistStockData("wl1", stockLists.WatchList1, &wg, resultChan, errChan)
	go obj.AddWatchlistStockData("wl2", stockLists.WatchList2, &wg, resultChan, errChan)
	go obj.AddWatchlistStockData("wl3", stockLists.WatchList3, &wg, resultChan, errChan)
	go obj.AddWatchlistStockData("wl4", stockLists.WatchList4, &wg, resultChan, errChan)
	go obj.AddWatchlistStockData("wl5", stockLists.WatchList5, &wg, resultChan, errChan)

	// Wait for the goroutine to finish
	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
	}()

	for err := range errChan {
		loggerconfig.Error("Alert Severity:P2-Mid, FetchWatchListsV3 Error in fetching contractDetails from redis or Unmarshalling and reqId ", reqH.RequestId, " and client Version:", reqH.ClientVersion, " and err ", err)
	}

	for result := range resultChan {
		// watchLists = append(watchLists, result)
		switch result.WatchListId {
		case "wl1":
			watchLists[0] = result
		case "wl2":
			watchLists[1] = result
		case "wl3":
			watchLists[2] = result
		case "wl4":
			watchLists[3] = result
		case "wl5":
			watchLists[4] = result
		default:
			// Handle unknown watchlistId case, if needed
			loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " FetchWatchListsV3 Invalid watchlistId: ", result.WatchListId, " ReqID:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		}

	}

	resp.WatchLists = watchLists
	loggerconfig.Info("FetchWatchListsV3 Successful, response:", helpers.LogStructAsJSON(resp), "clientID: ", req.ClientId, " requestId:", reqH.RequestId, " deviceType: ", reqH.DeviceType, " clientVersion:", reqH.ClientVersion)
	apiRes.Data = resp
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}
func duplicateStockInWatchlistV2(currStockId string, addedStocks []models.StockDetailsV2) bool {
	for i := 0; i < len(addedStocks); i++ {
		if addedStocks[i].Isin == "" && strings.EqualFold(currStockId, addedStocks[i].StockId) {
			return true
		} else if addedStocks[i].Isin != "" && strings.EqualFold(currStockId, addedStocks[i].IsinStockId) {
			return true
		}
	}
	// unique stock
	return false
}

func (obj WatchlistsObj) AddStockToWatchListV3(req models.AddStockToWatchListV3Request, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var stockLists models.MongoNewWatchListsV2
	stockLists.ClientId = req.ClientId
	var err error
	err = dbops.MongoRepo.FindOne(constants.WATCHLISTSTOCKSCOLLECTIONNEW, bson.M{"clientId": req.ClientId}, &stockLists)
	var apiRes apihelpers.APIRes
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("Alert Severity:P1-High, AddStockToWatchListV3 Mongo error =", err, " clientId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	stockId := req.Stock.Exchange + "-" + req.Stock.Token

	// display name check
	if obj.displayNameCheck {
		// a unique key for stock
		stockUniqueKey := strings.ToUpper(req.Stock.Exchange) + "_" + req.Stock.Token
		err, val := obj.contractCacheCli.GetFromHash("stock_key", stockUniqueKey)
		if err != nil {
			loggerconfig.Error("Alert Severity:P1-High, AddStockToWatchListV3: Error while GetFromHashSetNew script v2 for query ", stockId, " and reqId ", reqH.RequestId, " and client Version:", reqH.ClientVersion, " and err ", err)
			return apihelpers.SendInternalServerError()
		}

		stockDetail := models.ContractDetails{}
		if err = json.Unmarshal([]byte(val), &stockDetail); err != nil {
			loggerconfig.Info("Alert Severity:P1-High, AddStockToWatchListV3: Error Unmarshalling stockmetadata for key  err ", err, " and reqId:", reqH.RequestId, " and client Version:", reqH.ClientVersion)
			return apihelpers.SendInternalServerError()
		}

		validStockDisplayName := strings.Contains(strings.ToUpper(stockDetail.Name), strings.ToUpper(req.Stock.DisplayName))

		if !validStockDisplayName {
			loggerconfig.Info("AddStockToWatchListV3, Invalid stock display name stockUniqueKey:", stockUniqueKey)
			apiRes.Status = false
			apiRes.Message = constants.ErrorCodeMap[constants.InvalidDisplayName]
			apiRes.ErrorCode = constants.InvalidDisplayName
			return http.StatusBadRequest, apiRes
		}
	}

	duplicateStock := false

	var uniqueStockId string
	if req.Stock.Isin == "" {
		uniqueStockId = req.Stock.Exchange + "-" + req.Stock.Token
		req.Stock.StockId = uniqueStockId
	} else {
		var stockDetails models.StockDetailsV2
		uniqueStockId = req.Stock.Exchange + "-" + req.Stock.Isin
		stockDetails.Isin = req.Stock.Isin
		stockDetails.Exchange = req.Stock.Exchange
		stockDetails.IsinStockId = req.Stock.Exchange + "-" + req.Stock.Isin
		req.Stock = stockDetails
	}

	for i := 0; i < len(req.WatchListId); i++ {

		if req.WatchListId[i] == "wl1" {
			duplicateStock = duplicateStockInWatchlistV2(uniqueStockId, stockLists.WatchList1)
			if duplicateStock {
				break
			}
			stockLists.WatchList1 = append([]models.StockDetailsV2{req.Stock}, stockLists.WatchList1...) // prepending the new stock in watchlist
		}

		if req.WatchListId[i] == "wl2" {
			duplicateStock = duplicateStockInWatchlistV2(uniqueStockId, stockLists.WatchList2)
			if duplicateStock {
				break
			}
			stockLists.WatchList2 = append([]models.StockDetailsV2{req.Stock}, stockLists.WatchList2...) // prepending the new stock in watchlist
		}

		if req.WatchListId[i] == "wl3" {
			duplicateStock = duplicateStockInWatchlistV2(uniqueStockId, stockLists.WatchList3)
			if duplicateStock {
				break
			}
			stockLists.WatchList3 = append([]models.StockDetailsV2{req.Stock}, stockLists.WatchList3...) // prepending the new stock in watchlist
		}

		if req.WatchListId[i] == "wl4" {
			duplicateStock = duplicateStockInWatchlistV2(uniqueStockId, stockLists.WatchList4)
			if duplicateStock {
				break
			}
			stockLists.WatchList4 = append([]models.StockDetailsV2{req.Stock}, stockLists.WatchList4...) // prepending the new stock in watchlist
		}

		if req.WatchListId[i] == "wl5" {
			duplicateStock = duplicateStockInWatchlistV2(uniqueStockId, stockLists.WatchList5)
			if duplicateStock {
				break
			}
			stockLists.WatchList5 = append([]models.StockDetailsV2{req.Stock}, stockLists.WatchList5...) // prepending the new stock in watchlist
		}

	}

	if duplicateStock {
		loggerconfig.Error("AddStockToWatchListV3 stock already present in watchlist", " clientId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Status = false
		apiRes.Message = constants.ErrorCodeMap[constants.WatchListsStockAlreadyExists]
		apiRes.ErrorCode = constants.WatchListsStockAlreadyExists
		return http.StatusBadRequest, apiRes
	}

	filter := bson.D{{"clientId", stockLists.ClientId}}
	update := bson.D{{"$set", stockLists}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.WATCHLISTSTOCKSCOLLECTIONNEW, filter, update, opts)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, AddStockToWatchListV3 Mongo Upsert failed error =", err, " clientId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	loggerconfig.Info("AddStockToWatchListV3, stock is added successfully for req packet: ", helpers.LogStructAsJSON(req), "req requestId:", reqH.RequestId, " deviceType: ", reqH.DeviceType, " clientVersion:", reqH.ClientVersion)
	return http.StatusOK, apiRes
}

func (obj WatchlistsObj) DeleteStockInWatchListV3(req models.DeleteWatchListV3Request, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var stockLists models.MongoNewWatchListsV2
	err := dbops.MongoRepo.FindOne(constants.WATCHLISTSTOCKSCOLLECTIONNEW, bson.M{"clientId": req.ClientId}, &stockLists)
	var apiRes apihelpers.APIRes
	if err != nil && err.Error() == constants.MongoNoDocError {
		return apihelpers.SendErrorResponse(false, constants.WatchListsDoesNotExists, http.StatusBadRequest)
	}
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, DeleteStockInWatchListV3 Mongo error =", err, " clientId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	stockIdMap := addStockIdsInMap(req.StockId)

	if req.WatchListId == "wl1" {
		var stockDetails []models.StockDetailsV2
		for i := 0; i < len(stockLists.WatchList1); i++ {
			if stockLists.WatchList1[i].Isin == "" {
				if _, present := stockIdMap[stockLists.WatchList1[i].StockId]; !present { // if block is only executed if current stockid is not present in req packet from frontend.
					stockDetails = append(stockDetails, stockLists.WatchList1[i]) // adding stocks to stockDetails from current watchlist
				}
			} else {
				if _, present := stockIdMap[stockLists.WatchList1[i].IsinStockId]; !present {
					stockDetails = append(stockDetails, stockLists.WatchList1[i])
				}
			}
		}
		stockLists.WatchList1 = stockDetails // reassigning watchlist stocks
	}

	if req.WatchListId == "wl2" {
		var stockDetails []models.StockDetailsV2
		for i := 0; i < len(stockLists.WatchList2); i++ {
			if stockLists.WatchList2[i].Isin == "" {
				if _, present := stockIdMap[stockLists.WatchList2[i].StockId]; !present {
					stockDetails = append(stockDetails, stockLists.WatchList2[i])
				}
			} else {
				if _, present := stockIdMap[stockLists.WatchList2[i].IsinStockId]; !present {
					stockDetails = append(stockDetails, stockLists.WatchList2[i])
				}
			}
		}
		stockLists.WatchList2 = stockDetails
	}

	if req.WatchListId == "wl3" {
		var stockDetails []models.StockDetailsV2
		for i := 0; i < len(stockLists.WatchList3); i++ {
			if stockLists.WatchList3[i].Isin == "" {
				if _, present := stockIdMap[stockLists.WatchList3[i].StockId]; !present {
					stockDetails = append(stockDetails, stockLists.WatchList3[i])
				}
			} else {
				if _, present := stockIdMap[stockLists.WatchList3[i].IsinStockId]; !present {
					stockDetails = append(stockDetails, stockLists.WatchList3[i])
				}
			}
		}
		stockLists.WatchList3 = stockDetails
	}

	if req.WatchListId == "wl4" {
		var stockDetails []models.StockDetailsV2
		for i := 0; i < len(stockLists.WatchList4); i++ {
			if stockLists.WatchList4[i].Isin == "" {
				if _, present := stockIdMap[stockLists.WatchList4[i].StockId]; !present {
					stockDetails = append(stockDetails, stockLists.WatchList4[i])
				}
			} else {
				if _, present := stockIdMap[stockLists.WatchList4[i].IsinStockId]; !present {
					stockDetails = append(stockDetails, stockLists.WatchList4[i])
				}
			}
		}
		stockLists.WatchList4 = stockDetails
	}

	if req.WatchListId == "wl5" {
		var stockDetails []models.StockDetailsV2
		for i := 0; i < len(stockLists.WatchList5); i++ {
			if stockLists.WatchList5[i].Isin == "" {
				if _, present := stockIdMap[stockLists.WatchList5[i].StockId]; !present {
					stockDetails = append(stockDetails, stockLists.WatchList5[i])
				}
			} else {
				if _, present := stockIdMap[stockLists.WatchList5[i].IsinStockId]; !present {
					stockDetails = append(stockDetails, stockLists.WatchList5[i])
				}
			}
		}
		stockLists.WatchList5 = stockDetails
	}

	filter := bson.D{{"clientId", req.ClientId}}
	update := bson.D{{"$set", stockLists}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.WATCHLISTSTOCKSCOLLECTIONNEW, filter, update, opts)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, DeleteStockInWatchListV3 Mongo Upsert failed error =", err, " clientId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	loggerconfig.Info("DeleteStockInWatchListV3, stock is deleted successfully for req packet: ", helpers.LogStructAsJSON(req), "req requestId:", reqH.RequestId, " deviceType: ", reqH.DeviceType, " clientVersion:", reqH.ClientVersion)

	return http.StatusOK, apiRes

}

func (obj WatchlistsObj) ArrangeStocksWatchListV3(req models.ArrangeStocksWatchListV3Request, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	var stockLists models.MongoNewWatchListsV2
	var err error
	err = dbops.MongoRepo.FindOne(constants.WATCHLISTSTOCKSCOLLECTIONNEW, bson.M{"clientId": req.ClientId}, &stockLists)
	if err != nil && err.Error() == constants.MongoNoDocError {
		loggerconfig.Error("ArrangeStocksWatchListV3 Failed to locate docs in Mongo corresponding to clientId:", req.ClientId, " requestid=", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendErrorResponse(false, constants.WatchListsDoesNotExists, http.StatusBadRequest)
	}
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, ArrangeStocksWatchListV3 Mongo error =", err, " clientId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	if req.WatchListId == "wl1" {
		stockIdMap := addStocksInMapV2(stockLists.WatchList1)
		stockLists.WatchList1 = nil
		for i := 0; i < len(req.StockIds); i++ {
			stockLists.WatchList1 = append(stockLists.WatchList1, stockIdMap[req.StockIds[i]])
		}
	} else if req.WatchListId == "wl2" {
		stockIdMap := addStocksInMapV2(stockLists.WatchList2)
		stockLists.WatchList2 = nil
		for i := 0; i < len(req.StockIds); i++ {
			stockLists.WatchList2 = append(stockLists.WatchList2, stockIdMap[req.StockIds[i]])
		}
	} else if req.WatchListId == "wl3" {
		stockIdMap := addStocksInMapV2(stockLists.WatchList3)
		stockLists.WatchList3 = nil
		for i := 0; i < len(req.StockIds); i++ {
			stockLists.WatchList3 = append(stockLists.WatchList3, stockIdMap[req.StockIds[i]])
		}
	} else if req.WatchListId == "wl4" {
		stockIdMap := addStocksInMapV2(stockLists.WatchList4)
		stockLists.WatchList4 = nil
		for i := 0; i < len(req.StockIds); i++ {
			stockLists.WatchList4 = append(stockLists.WatchList4, stockIdMap[req.StockIds[i]])
		}
	} else if req.WatchListId == "wl5" {
		stockIdMap := addStocksInMapV2(stockLists.WatchList5)
		stockLists.WatchList5 = nil
		for i := 0; i < len(req.StockIds); i++ {
			stockLists.WatchList5 = append(stockLists.WatchList5, stockIdMap[req.StockIds[i]])
		}
	}

	filter := bson.D{{"clientId", req.ClientId}}
	update := bson.D{{"$set", stockLists}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.WATCHLISTSTOCKSCOLLECTIONNEW, filter, update, opts)
	if err != nil {
		loggerconfig.Error("ArrangeStocksWatchListV3 Failed to update MongoNewWatchLists to mongo, error = ", err, "clientID: ", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("ArrangeStocksWatchListV3 packet=", helpers.LogStructAsJSON(stockLists), " uccId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}

func addStocksInMapV2(stocks []models.StockDetailsV2) map[string]models.StockDetailsV2 {
	stockIdIndexMap := make(map[string]models.StockDetailsV2)

	for _, stock := range stocks {
		if stock.Isin != "" {
			stockIdIndexMap[stock.IsinStockId] = stock
		} else {
			stockIdIndexMap[stock.StockId] = stock
		}
	}
	return stockIdIndexMap
}

func (obj WatchlistsObj) DeleteStockInWatchListV3Updated(req models.DeleteWatchListV3UpdatedRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var stockLists models.MongoNewWatchListsV2
	err := dbops.MongoRepo.FindOne(constants.WATCHLISTSTOCKSCOLLECTIONNEW, bson.M{"clientId": req.ClientId}, &stockLists)
	var apiRes apihelpers.APIRes
	if err != nil && err.Error() == constants.MongoNoDocError {
		return apihelpers.SendErrorResponse(false, constants.WatchListsDoesNotExists, http.StatusBadRequest)
	}

	var stockIdMap map[string]int

	for watchListIndex := 0; watchListIndex < len(req.WatchListId); watchListIndex++ {
		switch req.WatchListId[watchListIndex].WatchlistId {
		case "wl1":
			stockIdMap = addStockIdsInMap(req.WatchListId[watchListIndex].StockId)
			noOfStocksInWL := len(req.WatchListId[watchListIndex].StockId)
			var stockDetails []models.StockDetailsV2
			for stockIndex := 0; stockIndex < len(stockLists.WatchList1); stockIndex++ {
				if stockLists.WatchList1[stockIndex].Isin == "" {
					if _, present := stockIdMap[stockLists.WatchList1[stockIndex].StockId]; !present {
						stockDetails = append(stockDetails, stockLists.WatchList1[stockIndex])
					} else {
						noOfStocksInWL--
					}
				} else {
					if _, present := stockIdMap[stockLists.WatchList1[stockIndex].IsinStockId]; !present {
						stockDetails = append(stockDetails, stockLists.WatchList1[stockIndex])
					} else {
						noOfStocksInWL--
					}
				}
			}
			if noOfStocksInWL > 0 {
				return apihelpers.SendErrorResponse(false, constants.InvalidRequest, http.StatusBadRequest)
			}
			stockLists.WatchList1 = stockDetails
		case "wl2":
			stockIdMap = addStockIdsInMap(req.WatchListId[watchListIndex].StockId)
			noOfStocksInWL := len(req.WatchListId[watchListIndex].StockId)
			var stockDetails []models.StockDetailsV2
			for stockIndex := 0; stockIndex < len(stockLists.WatchList2); stockIndex++ {
				if stockLists.WatchList2[stockIndex].Isin == "" {
					if _, present := stockIdMap[stockLists.WatchList2[stockIndex].StockId]; !present {
						stockDetails = append(stockDetails, stockLists.WatchList2[stockIndex])
					} else {
						noOfStocksInWL--
					}
				} else {
					if _, present := stockIdMap[stockLists.WatchList2[stockIndex].IsinStockId]; !present {
						stockDetails = append(stockDetails, stockLists.WatchList2[stockIndex])
					} else {
						noOfStocksInWL--
					}
				}
			}
			if noOfStocksInWL > 0 {
				return apihelpers.SendErrorResponse(false, constants.InvalidRequest, http.StatusBadRequest)
			}
			stockLists.WatchList2 = stockDetails
		case "wl3":
			stockIdMap = addStockIdsInMap(req.WatchListId[watchListIndex].StockId)
			noOfStocksInWL := len(req.WatchListId[watchListIndex].StockId)
			var stockDetails []models.StockDetailsV2
			for stockIndex := 0; stockIndex < len(stockLists.WatchList3); stockIndex++ {
				if stockLists.WatchList3[stockIndex].Isin == "" {
					if _, present := stockIdMap[stockLists.WatchList3[stockIndex].StockId]; !present {
						stockDetails = append(stockDetails, stockLists.WatchList3[stockIndex])
					} else {
						noOfStocksInWL--
					}
				} else {
					if _, present := stockIdMap[stockLists.WatchList3[stockIndex].IsinStockId]; !present {
						stockDetails = append(stockDetails, stockLists.WatchList3[stockIndex])
					} else {
						noOfStocksInWL--
					}
				}
			}
			if noOfStocksInWL > 0 {
				return apihelpers.SendErrorResponse(false, constants.InvalidRequest, http.StatusBadRequest)
			}
			stockLists.WatchList3 = stockDetails
		case "wl4":
			stockIdMap = addStockIdsInMap(req.WatchListId[watchListIndex].StockId)
			noOfStocksInWL := len(req.WatchListId[watchListIndex].StockId)
			var stockDetails []models.StockDetailsV2
			for stockIndex := 0; stockIndex < len(stockLists.WatchList4); stockIndex++ {
				if stockLists.WatchList4[stockIndex].Isin == "" {
					if _, present := stockIdMap[stockLists.WatchList4[stockIndex].StockId]; !present {
						stockDetails = append(stockDetails, stockLists.WatchList4[stockIndex])
					} else {
						noOfStocksInWL--
					}
				} else {
					if _, present := stockIdMap[stockLists.WatchList4[stockIndex].IsinStockId]; !present {
						stockDetails = append(stockDetails, stockLists.WatchList4[stockIndex])
					} else {
						noOfStocksInWL--
					}
				}
			}
			if noOfStocksInWL > 0 {
				return apihelpers.SendErrorResponse(false, constants.InvalidRequest, http.StatusBadRequest)
			}
			stockLists.WatchList4 = stockDetails
		case "wl5":
			stockIdMap = addStockIdsInMap(req.WatchListId[watchListIndex].StockId)
			noOfStocksInWL := len(req.WatchListId[watchListIndex].StockId)
			var stockDetails []models.StockDetailsV2
			for stockIndex := 0; stockIndex < len(stockLists.WatchList5); stockIndex++ {
				if stockLists.WatchList5[stockIndex].Isin == "" {
					if _, present := stockIdMap[stockLists.WatchList5[stockIndex].StockId]; !present {
						stockDetails = append(stockDetails, stockLists.WatchList5[stockIndex])
					} else {
						noOfStocksInWL--
					}
				} else {
					if _, present := stockIdMap[stockLists.WatchList5[stockIndex].IsinStockId]; !present {
						stockDetails = append(stockDetails, stockLists.WatchList5[stockIndex])
					} else {
						noOfStocksInWL--
					}
				}
			}
			if noOfStocksInWL > 0 {
				return apihelpers.SendErrorResponse(false, constants.InvalidRequest, http.StatusBadRequest)
			}
			stockLists.WatchList5 = stockDetails
		default:
			return apihelpers.SendErrorResponse(false, constants.InvalidRequest, http.StatusBadRequest)
		}
	}

	filter := bson.D{{"clientId", req.ClientId}}
	update := bson.D{{"$set", stockLists}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.WATCHLISTSTOCKSCOLLECTIONNEW, filter, update, opts)
	if err != nil {
		loggerconfig.Error("DeleteStockInWatchListV3Updated Mongo Upsert failed error =", err, " clientId:", req.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}
