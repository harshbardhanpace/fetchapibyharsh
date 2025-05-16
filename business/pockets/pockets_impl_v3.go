package pockets

import (
	"context"
	"errors"
	"net/http"
	apihelpers "space/apiHelpers"
	"space/business/tradelab"
	"space/constants"
	"space/db"
	"space/dbops"
	"space/helpers"
	"space/helpers/cache"
	"space/loggerconfig"
	"space/models"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ExecutePocketV3Obj struct {
	mongodb  db.MongoDatabase
	redisCli cache.RedisCache
}

func InitExecutePocketV3Provider(mongodb db.MongoDatabase, redisCli cache.RedisCache) ExecutePocketV3Obj {
	defer models.HandlePanic()
	exectutePocketObj := ExecutePocketV3Obj{mongodb: mongodb, redisCli: redisCli}
	return exectutePocketObj
}

func (obj ExecutePocketV3Obj) BuyPocketV3(req models.ExecutePocketV3Request, pocketAction string, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	// Fetch pocket details for the requested pocket ID
	latestPocket, httpCode, errorCode := getLatestPocketDetails(req.PocketId, reqH)
	if errorCode != "" {
		return apihelpers.SendErrorResponse(false, errorCode, httpCode)
	}

	// Fetch user's pockets to check if they already own this pocket
	var userPocketHolding models.UserPocketHolding
	err := dbops.MongoRepo.FindOne(constants.USERPOCKETSCOLLECTION, bson.M{"clientId": req.ClientId}, &userPocketHolding)
	if err != nil && err.Error() == constants.MongoNoDocError {

		// User does not have this pocket, proceed with a fresh purchase
		loggerconfig.Info("BuyPocketV3, userPocketHolding doesn't have the pocket with id", req.PocketId, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return obj.BuyPocket(req, latestPocket, userPocketHolding, reqH)

	} else if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " BuyPocketV3 mongo error =", err, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	// Check if the user already has the specified pocket
	for _, userPocket := range userPocketHolding.Pockets {
		if userPocket.PocketID == latestPocket.PocketId {

			//check for mismatch in pocket stocks
			_, rebalanceResp := checkRebalanceRequirements(req.ClientId, userPocket.Version, userPocket.LotSize, latestPocket, reqH)
			loggerconfig.Info("BuyPocketV3 CheckRebalanceRequirements response:=", helpers.LogStructAsJSON(rebalanceResp), " requestId:", reqH.RequestId)

			rebalanceData, ok := rebalanceResp.Data.(models.RebalanceResponse)
			if !ok {
				loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, "BuyPocketV3 Failed to parse rebalance data for clientId:", req.ClientId, " requestId:", reqH.RequestId)
				return apihelpers.SendInternalServerError()
			}

			// Check if rebalance or repair action is required
			if len(rebalanceData.BuyRequirements) > 0 {
				var resp models.PocketOrdersRes

				return http.StatusOK, apihelpers.APIRes{
					Status:  true,
					Message: "Action Required",
					Data:    resp,
				}
			}
			break
		}
	}

	loggerconfig.Info("BuyPocketV3, user wants to buy a fresh pocket with id", req.PocketId, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	return obj.BuyPocket(req, latestPocket, userPocketHolding, reqH)
}

func (obj ExecutePocketV3Obj) BuyPocket(req models.ExecutePocketV3Request, latestPocket models.MongoLatestPocketDetails, userPocketHolding models.UserPocketHolding, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	orderObj := tradelab.InitOrder(obj.redisCli)
	orderResults, mapPocketMetaDataAgainstToken, placedOrderIds, httpCode, errorCode := placePocketOrder(req, latestPocket, constants.BUY, orderObj, reqH)
	if errorCode != "" || len(orderResults) == 0 {
		return apihelpers.SendErrorResponse(false, errorCode, httpCode)
	}

	processOrderReq := models.ProcessOrderReq{
		PocketId:          req.PocketId,
		ClientId:          req.ClientId,
		LotSize:           req.LotSize,
		Action:            constants.BUY,
		TransactionStatus: 0,
	}

	pocketDetailsV3, orderCompleted, orderCancelled := processOrderResults(orderResults, mapPocketMetaDataAgainstToken, placedOrderIds, processOrderReq, reqH)

	responseMessage := "Successfully bought the pocket"

	// if any number of orders got placed, update the user's holdings
	if len(orderCompleted) != 0 {

		//update transaction store
		err := updatePocketTransactionStore(req.ClientId, pocketDetailsV3, reqH)
		if err != nil {
			loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, "BuyPocketV3 updatePocketTransactionStore error =", err, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		//update user's holdings
		err = updateUserPocketHoldings(processOrderReq, latestPocket.PocketVersion, userPocketHolding, reqH)
		if err != nil {
			return apihelpers.SendInternalServerError()
		}
	} else { // if all the orders got cancelled
		return http.StatusOK, apihelpers.APIRes{
			Status:  false,
			Message: "Failed to buy the pocket",
			Data:    nil,
		}
	}

	if len(orderCancelled) != 0 {
		responseMessage = "Partially bought the pocket"
	}

	resp := models.PocketOrdersRes{
		OrderCompleted: orderCompleted,
		OrderCancelled: orderCancelled,
	}

	return http.StatusOK, apihelpers.APIRes{
		Status:  true,
		Message: responseMessage,
		Data:    resp,
	}
}

func (obj ExecutePocketV3Obj) CheckActionRequired(clientId, pocketId string, userVersion, UsersLotSize int, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	latestPocket, httpCode, errorCode := getLatestPocketDetails(pocketId, reqH)
	if errorCode != "" {
		return apihelpers.SendErrorResponse(false, errorCode, httpCode)
	}

	return checkRebalanceRequirements(clientId, userVersion, UsersLotSize, latestPocket, reqH)

}

func checkRebalanceRequirements(clientId string, userVersion int, UsersLotSize int, latestPocket models.MongoLatestPocketDetails, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var buyRequirements []models.PocketsMetaData
	action := constants.Repair

	//get user's all positions
	var getPositionsReq models.GetPositionRequest
	getPositionsReq.ClientID = clientId
	getPositionsReq.Type = "historical"

	statusPosition, resPositions := tradelab.GetPositionsInternal(getPositionsReq, reqH)
	if statusPosition != http.StatusOK {
		loggerconfig.Error("CheckRebalanceRequirements GetPositions status != 200", statusPosition, " uccId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	loggerconfig.Info("ExecutePocketV3 CallTLforGetPositionsInternal response:=", helpers.LogStructAsJSON(resPositions), " uccId:", reqH.ClientId, " requestId:", reqH.RequestId)

	getPositionsRes, ok := resPositions.Data.([]models.GetPositionResponseData)
	if !ok {
		loggerconfig.Error("Alert Severity:P1-High, CheckRebalanceRequirements GetPosition interface parsing error", ok, " uccId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	// Map user's current holdings by stock token : quantity
	userHoldingsMap := make(map[int]int)
	for _, holding := range getPositionsRes {
		if holding.BuyQuantity-holding.SellQuantity > 0 {
			userHoldingsMap[holding.Token] = holding.BuyQuantity - holding.SellQuantity
		}
	}

	// get all the holdings of user
	holdingReq := models.FetchDematHoldingsRequest{
		ClientID: clientId,
	}
	portfolioObj := tradelab.InitPortfolio()
	statusHoldings, resHoldings := portfolioObj.FetchDematHoldings(holdingReq, reqH)
	if statusHoldings != http.StatusOK {
		loggerconfig.Error("CheckRebalanceRequirements GetPositions status != 200", statusPosition, " uccId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	getHoldingRes, ok1 := resHoldings.Data.(models.FetchDematHoldingsResponse)
	if !ok1 {
		loggerconfig.Error("Alert Severity:P1-High, CheckRebalanceRequirements GetPosition interface parsing error", ok, " uccId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	for _, holding := range getHoldingRes.Holdings {
		if existingQty, found := userHoldingsMap[holding.Token]; found { // if there is already a holding for user add the qty
			userHoldingsMap[holding.Token] = holding.Quantity + existingQty
		} else {
			userHoldingsMap[holding.Token] = holding.Quantity
		}
	}

	// Calculate buy/sell requirements by comparing the user's pocket with the latest pocket
	for _, stock := range latestPocket.PocketTokens {

		pocketStockToken, err1 := strconv.Atoi(stock.Token)
		pocketStockQty, err2 := strconv.Atoi(stock.Qty)
		if err1 != nil || err2 != nil {
			loggerconfig.Error("CheckRebalanceRequirements, getting error while parsing pocketStockToken: ", pocketStockToken, "err: ", err1, " and reqQty: ", pocketStockQty, "err: ", err2, " for client: ", clientId)
			apihelpers.SendInternalServerError()
		}

		//multiply it with the lot user have
		pocketStockQty = pocketStockQty * UsersLotSize

		// quantity that user have
		userQty, exists := userHoldingsMap[pocketStockToken]

		if !exists {
			// If any stock is missing from user's holdings, set action to "rebalance"
			buyRequirements = append(buyRequirements, stock)
		} else if userQty < pocketStockQty {
			// If stock exists but quantity is less, add to buyRequirements with remaining qty
			remainingQty := pocketStockQty - userQty

			stockToBuy := stock
			stockToBuy.Qty = strconv.Itoa(remainingQty)

			buyRequirements = append(buyRequirements, stockToBuy)
		}

	}

	if latestPocket.PocketVersion != userVersion {
		action = constants.Rebalance
	}

	var actionMessage string
	if len(buyRequirements) > 0 {
		actionMessage = action + " required to update pocket to the latest version."
	} else {
		actionMessage = "No aciton required"
		action = constants.BUY
	}

	// Create a response summarizing the required buy/sell actions
	rebalanceResponse := models.RebalanceResponse{
		ClientId:        clientId,
		PocketId:        latestPocket.PocketId,
		BuyRequirements: buyRequirements,
		Message:         actionMessage,
		Action:          action,
		OrderSide:       constants.BUY,
		PocketVersion:   latestPocket.PocketVersion,
	}

	loggerconfig.Info("CheckRebalanceRequirements completed", " clientId:", clientId, " requestId:", reqH.RequestId, " action:", action)

	return http.StatusOK, apihelpers.APIRes{
		Data:    rebalanceResponse,
		Message: "Align your pocket with the latest version.",
		Status:  true,
	}
}

func (obj ExecutePocketV3Obj) ManageRequiredStocksForPocket(pocketBalanceReq models.RebalanceResponse, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	placedOrderIds := make(map[string]bool)
	mapPocketMetaDataAgainstToken := make(map[string]models.PocketsMetaData)
	pocketTokens := pocketBalanceReq.BuyRequirements
	orderObj := tradelab.InitOrder(obj.redisCli)

	for _, StockDetails := range pocketTokens {
		mapPocketMetaDataAgainstToken[StockDetails.Token] = StockDetails

		qty, err := strconv.Atoi(StockDetails.Qty)
		if err != nil {
			loggerconfig.Error("ManageRequiredStocksForPocket, Conversion error for qty", " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
			return http.StatusInternalServerError, apihelpers.APIRes{}
		}

		placeOrderReq := createPlaceOrderRequest(pocketBalanceReq.ClientId, StockDetails.Exchange, StockDetails.Token, "BUY")
		placeOrderReq.Quantity = qty

		orderId, httpStatus, errorCode := executePlaceOrder(orderObj, placeOrderReq, pocketBalanceReq.ClientId, reqH)
		if errorCode != "" {
			return httpStatus, apihelpers.APIRes{}
		}

		placedOrderIds[orderId] = true
	}

	if len(placedOrderIds) == 0 {
		loggerconfig.Info("ManageRequiredStocksForPocket, No orders were placed", " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return http.StatusOK, apihelpers.APIRes{}
	}

	// Get completed orders
	completedOrders, httpStatus, errorCode := getCompletedOrders(orderObj, pocketBalanceReq.ClientId, reqH)
	if errorCode != "" {
		return httpStatus, apihelpers.APIRes{}
	}

	//find orderside either rebalace or repair
	processOrderReq := models.ProcessOrderReq{
		PocketId:          pocketBalanceReq.PocketId,
		ClientId:          pocketBalanceReq.ClientId,
		LotSize:           0,
		Action:            pocketBalanceReq.Action,
		TransactionStatus: 0,
	}

	pocketDetailsV3, orderCompleted, orderCancelled := processOrderResults(completedOrders, mapPocketMetaDataAgainstToken, placedOrderIds, processOrderReq, reqH)

	// if any number of orders got placed, update the user's holdings
	responseMessage := "Successfully balanced the pocket"

	// if any number of orders got placed, update the user's transactions
	if len(orderCompleted) != 0 {

		//update transaction store
		err := updatePocketTransactionStore(pocketBalanceReq.ClientId, pocketDetailsV3, reqH)
		if err != nil {
			return apihelpers.SendInternalServerError()
		}

		//if action is rebalace, update user's pocket-version
		if pocketBalanceReq.Action == constants.Rebalance {

			var userPocketHolding models.UserPocketHolding
			err := dbops.MongoRepo.FindOne(constants.USERPOCKETSCOLLECTION, bson.M{"clientId": pocketBalanceReq.ClientId}, &userPocketHolding)
			if err != nil {
				if err.Error() == constants.MongoNoDocError {
					return http.StatusBadRequest, apihelpers.APIRes{
						Status:  false,
						Message: "Failed to rebalance the pocket",
						Data:    nil,
					}
				}
				loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " ManageRequiredStocksForPocket mongo error =", err, " requestId:", reqH.RequestId)
				return apihelpers.SendInternalServerError()
			}

			processOrderReq.Action = constants.Rebalance

			err = updateUserPocketHoldings(processOrderReq, pocketBalanceReq.PocketVersion, userPocketHolding, reqH)
			if err != nil {
				return apihelpers.SendInternalServerError()
			}
		}

	} else { // if all the orders got cancelled
		return http.StatusOK, apihelpers.APIRes{
			Status:  false,
			Message: "Failed to rebalance the pocket",
			Data:    nil,
		}
	}

	if len(orderCancelled) != 0 {
		responseMessage = "Align your pocket with the latest version."
	}

	resp := models.PocketOrdersRes{
		OrderCompleted: orderCompleted,
		OrderCancelled: orderCancelled,
	}

	return http.StatusOK, apihelpers.APIRes{
		Status:  true,
		Message: responseMessage,
		Data:    resp,
	}
}

func (obj ExecutePocketV3Obj) FetchAllPocketsV3(reqH models.ReqHeader) (int, apihelpers.APIRes) {

	var allPockets []models.MongoPocketsV3
	var count = 0
	cursor, err := dbops.MongoRepo.Find(constants.POCKETSCOLLECTIONV3, bson.M{})
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " FetchAllPocketsV3 Mongo Find() failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	defer func() {
		if err := cursor.Close(context.Background()); err != nil {
			loggerconfig.Error("Alert Severity:P1-High, FetchAllPocketsV3 failed to close cursor, error:", err, " requestId:", reqH.RequestId, "clientId:", reqH.ClientId)
		}
	}()

	if cursor == nil {
		loggerconfig.Error("FetchAllPocketsV3 cursor is nil, requestId:", reqH.RequestId, "clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	for cursor.Next(context.Background()) {
		count++
		var pockets models.MongoPocketsV3
		err := cursor.Decode(&pockets)
		if err != nil {
			loggerconfig.Error("Alert Severity:P1-High, FetchAllPocketsV3 Error in decoding for mongo cursor, err: ", err, " requestId: ", reqH.RequestId, "clientVersion:", reqH.ClientVersion)
			continue
		}
		allPockets = append(allPockets, pockets)
	}

	var apiRes apihelpers.APIRes

	if len(allPockets) == 0 {
		apiRes.Message = constants.NoPocketFound
		apiRes.Status = true
		apiRes.Data = []models.MongoPocketsV3{}
		return http.StatusOK, apiRes
	}

	apiRes.Data = allPockets
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ExecutePocketV3Obj) FetchPocketPortfolioV3(req models.FetchPocketPortfolioRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	fetchPocketPortfolioResponseV2 := models.FetchPocketPortfolioResponseV2{
		ClientId: req.ClientId,
	}

	// Fetch user pocket holdings once
	var userPocketHolding models.UserPocketHolding
	err := dbops.MongoRepo.FindOne(constants.USERPOCKETSCOLLECTION, bson.M{
		"clientId": bson.M{"$regex": "^" + req.ClientId + "$", "$options": "i"},
	}, &userPocketHolding)
	if err != nil {
		if err.Error() == constants.MongoNoDocError {
			return http.StatusOK, apihelpers.APIRes{Message: "No user pockets found for the client.", Status: true, Data: fetchPocketPortfolioResponseV2}
		}
		loggerconfig.Error("Alert Severity:P0-Critical, FetchPocketPortfolioV3: Mongo Find() failed - error:", err, " ClientID:", req.ClientId, " RequestID:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	// get user version of each pocket
	userPocketMap := make(map[string]int, len(userPocketHolding.Pockets))
	for _, userPocket := range userPocketHolding.Pockets {
		userPocketMap[userPocket.PocketID] = userPocket.Version
	}

	var pocketTransactionStoreV3 models.PocketTransactionStoreV3
	err = dbops.MongoRepo.FindOne(constants.POCKETSTRANSACTIONSDetailsV2, bson.M{"clientId": bson.M{"$regex": "^" + req.ClientId + "$", "$options": "i"}}, &pocketTransactionStoreV3)
	if err != nil {
		if err.Error() == constants.MongoNoDocError {
			return http.StatusOK, apihelpers.APIRes{Message: "No transactions found for the client.", Status: true, Data: fetchPocketPortfolioResponseV2}
		}
		loggerconfig.Error("Alert Severity:P0-Critical, FetchPocketPortfolioV3: Mongo Find() failed - error:", err, " ClientID:", req.ClientId, " RequestID:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	// organize the transactions by pocket id
	pocketData := make(map[string]*models.PocketPortfolioDetails)

	for _, purchase := range pocketTransactionStoreV3.AllPocketPurchases {

		pocketID := purchase.PocketId
		pocketQty := purchase.LotSize

		if _, exists := pocketData[pocketID]; !exists {
			pocketData[pocketID] = &models.PocketPortfolioDetails{
				PocketId:           pocketID,
				UsersPocketVersion: userPocketMap[pocketID],
			}
		}
		pocket := pocketData[pocketID]
		if purchase.TransactionStatus == 0 && strings.ToLower(purchase.Action) != constants.Repair { // Buy
			pocket.TotalBuyPockets += pocketQty
			pocket.AveragePrice += purchase.OrderCompletedPrice
		} else if purchase.TransactionStatus == 1 { // Sell
			pocket.TotalSellPockets += pocketQty
		}
	}

	for _, pocket := range pocketData {
		if pocket.TotalBuyPockets == pocket.TotalSellPockets {
			continue
		}

		// Fetch latest pocket details
		latestPocket, httpCode, errorCode := getLatestPocketDetails(pocket.PocketId, reqH)
		if errorCode != "" {
			loggerconfig.Error("Alert Severity:P0-Critical, FetchPocketPortfolioV3, error in fetching latest pockets detail - error:", constants.ErrorCodeMap[errorCode], " ClientID:", req.ClientId, " RequestID:", reqH.RequestId)
			return apihelpers.SendErrorResponse(false, errorCode, httpCode)
		}

		pocket.PocketName = latestPocket.PocketName
		pocket.PocketShortDesc = latestPocket.PocketShortDesc
		pocket.PocketLongDesc = latestPocket.PocketLongDesc
		pocket.PocketExchange = latestPocket.PocketExchange
		pocket.PocketImage = latestPocket.PocketImage
		pocket.PocketWebImage = latestPocket.PocketWebImage
		pocket.PrimaryBackgroundColor = latestPocket.PrimaryBackgroundColor
		pocket.PrimarySecondaryColor = latestPocket.PrimarySecondaryColor
		pocket.PocketBenchMark = latestPocket.PocketBenchMark
		pocket.PocketCreateTimeUnix = latestPocket.PocketCreateTimeUnix
		pocket.PocketVersion = latestPocket.PocketVersion
		pocket.LatestPocketTokens = latestPocket.PocketTokens
		pocket.AveragePrice /= float64(pocket.TotalBuyPockets)
		pocket.TotalInvestment = pocket.AveragePrice * float64(pocket.TotalBuyPockets-pocket.TotalSellPockets)
		fetchPocketPortfolioResponseV2.PortfolioDetails = append(fetchPocketPortfolioResponseV2.PortfolioDetails, *pocket)
	}

	loggerconfig.Info("FetchPocketPortfolioV3: Portfolio fetched successfully. ClientID:", req.ClientId, " RequestID:", reqH.RequestId)

	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	apiRes.Data = fetchPocketPortfolioResponseV2

	return http.StatusOK, apiRes
}

func (obj ExecutePocketV3Obj) FetchUsersPockets(clientId string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var userPocketHolding models.UserPocketHolding
	var apiRes apihelpers.APIRes

	err := dbops.MongoRepo.FindOne(constants.USERPOCKETSCOLLECTION, bson.M{
		"clientId": bson.M{"$regex": "^" + clientId + "$", "$options": "i"}, //ensures full string match and case insensitive.
	}, &userPocketHolding)
	if err != nil {
		if err.Error() == constants.MongoNoDocError {
			apiRes.Data = userPocketHolding
			apiRes.Message = constants.NoPocketFound
			apiRes.Status = true

			return http.StatusOK, apiRes
		}
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " FetchUsersPockets Mongo Find() failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("FetchUsersPockets: Successfully fetched user pockets. ClientID:", clientId, " RequestID:", reqH.RequestId)

	apiRes.Data = userPocketHolding
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ExecutePocketV3Obj) SellPocketV3(req models.ExecutePocketV3Request, pocketAction string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	const orderSide = "SELL"

	// Validate user's pocket existence
	userPocketHolding, httpCode, errorCode := validateUserPocket(req, reqH)
	if errorCode != "" {
		loggerconfig.Info("SellPocketV3, invalid pocket to sell, id: ", req.PocketId, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, errorCode, httpCode)
	}

	// Fetch latest pocket details
	latestPocket, httpCode, errorCode := getLatestPocketDetails(req.PocketId, reqH)
	if errorCode != "" {
		return apihelpers.SendErrorResponse(false, errorCode, httpCode)
	}

	// Place sell orders for pocket stocks
	orderObj := tradelab.InitOrder(obj.redisCli)
	orderResults, mapPocketMetaDataAgainstToken, placedOrderIds, httpCode, errorCode := placePocketOrder(req, latestPocket, orderSide, orderObj, reqH)
	if errorCode != "" {
		return apihelpers.SendErrorResponse(false, errorCode, httpCode)
	}

	processOrderReq := models.ProcessOrderReq{
		PocketId:          req.PocketId,
		ClientId:          req.ClientId,
		LotSize:           req.LotSize,
		Action:            orderSide,
		TransactionStatus: 1,
	}

	// Process completed and cancelled orders
	pocketDetailsV3, orderCompleted, orderCancelled := processOrderResults(orderResults, mapPocketMetaDataAgainstToken, placedOrderIds, processOrderReq, reqH)

	responseMessage := "Partially sold the pocket"

	//if any order got placed, update the transaction.
	if len(orderCompleted) != 0 {
		//update transaction store
		err := updatePocketTransactionStore(req.ClientId, pocketDetailsV3, reqH)
		if err != nil {
			return apihelpers.SendInternalServerError()
		}
	}

	// if all order got placed, update user's holding
	if len(orderCancelled) == 0 && len(orderCompleted) != 0 {
		err := updateUserPocketHoldings(processOrderReq, latestPocket.PocketVersion, userPocketHolding, reqH)
		if err != nil {
			return apihelpers.SendInternalServerError()
		}
		responseMessage = "Successfully sold the pocket"
	}

	// if no order got placed
	if len(orderCancelled) != 0 && len(orderCompleted) == 0 {
		return http.StatusOK, apihelpers.APIRes{
			Status:  false,
			Message: "failed to sell the pocket",
			Data:    nil,
		}
	}

	resp := models.PocketOrdersRes{
		OrderCompleted: orderCompleted,
		OrderCancelled: orderCancelled,
	}

	return http.StatusOK, apihelpers.APIRes{
		Status:  true,
		Message: responseMessage,
		Data:    resp,
	}
}

func validateUserPocket(req models.ExecutePocketV3Request, reqH models.ReqHeader) (models.UserPocketHolding, int, string) {
	var userPocketHolding models.UserPocketHolding
	err := dbops.MongoRepo.FindOne(constants.USERPOCKETSCOLLECTION, bson.M{
		"clientId": req.ClientId,
	}, &userPocketHolding)
	if err != nil {
		if err.Error() == constants.MongoNoDocError {
			return userPocketHolding, http.StatusBadRequest, constants.PocketDoesNotExists
		}
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " validateUserPocket Mongo Find() failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return userPocketHolding, http.StatusInternalServerError, constants.InternalServerError
	}

	pocketExists := false
	validLotSize := false

	for _, userPocket := range userPocketHolding.Pockets {
		if userPocket.PocketID == req.PocketId {
			pocketExists = true
			if userPocket.LotSize >= req.LotSize {
				validLotSize = true
				break
			}
		}
	}

	if !pocketExists {
		return userPocketHolding, http.StatusBadRequest, constants.PocketDoesNotExists
	}
	if !validLotSize {
		return userPocketHolding, http.StatusBadRequest, constants.LotSizeExceeds
	}

	return userPocketHolding, http.StatusOK, ""
}

// create function to return the Map of user's current holdings by stock token : quantity
func getAllUserHoldings(clientId string, reqH models.ReqHeader) (map[int]int, error) {
	userHoldingsMap := make(map[int]int)

	//get user's all positions
	var getPositionsReq models.GetPositionRequest
	getPositionsReq.ClientID = clientId
	getPositionsReq.Type = "historical"

	statusPosition, resPositions := tradelab.GetPositionsInternal(getPositionsReq, reqH)
	if statusPosition != http.StatusOK {
		loggerconfig.Error("getAllUserHoldings GetPositions status != 200", statusPosition, " uccId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return nil, errors.New("get positions failed")
	}

	getPositionsRes, ok := resPositions.Data.([]models.GetPositionResponseData)
	if !ok {
		loggerconfig.Error("Alert Severity:P1-High, getAllUserHoldings GetPosition interface parsing error", ok, " uccId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return nil, errors.New("interface parsing error")
	}

	for _, holding := range getPositionsRes {
		if holding.BuyQuantity-holding.SellQuantity > 0 {
			userHoldingsMap[holding.Token] = holding.BuyQuantity - holding.SellQuantity
		}
	}

	// get all the holdings of user
	holdingReq := models.FetchDematHoldingsRequest{
		ClientID: clientId,
	}
	portfolioObj := tradelab.InitPortfolio()
	statusHoldings, resHoldings := portfolioObj.FetchDematHoldings(holdingReq, reqH)
	if statusHoldings != http.StatusOK {
		loggerconfig.Error("getAllUserHoldings, FetchDematHoldings status != 200", statusHoldings, " uccId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return nil, errors.New("fetch demat holdings failed")
	}

	getHoldingRes, ok := resHoldings.Data.(models.FetchDematHoldingsResponse)
	if !ok {
		loggerconfig.Error("Alert Severity:P1-High, getAllUserHoldings, interface parsing error", "uccId:", reqH.ClientId, "requestId:", reqH.RequestId)
		return nil, errors.New("interface parsing error")
	}

	for _, holding := range getHoldingRes.Holdings {
		if existingQty, found := userHoldingsMap[holding.Token]; found { // if there is already a holding for user, add new to existing quantity
			userHoldingsMap[holding.Token] = holding.Quantity + existingQty
		} else {
			userHoldingsMap[holding.Token] = holding.Quantity
		}
	}

	return userHoldingsMap, nil
}

func getLatestPocketDetails(pocketId string, reqH models.ReqHeader) (models.MongoLatestPocketDetails, int, string) {
	var latestPocket models.MongoLatestPocketDetails
	err := dbops.MongoRepo.FindOne(constants.LATESTPOCKETDETAILS, bson.M{"pocketId": pocketId}, &latestPocket)
	if err != nil {
		if err.Error() == constants.MongoNoDocError {
			return latestPocket, http.StatusBadRequest, constants.PocketDoesNotExists
		}
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " getLatestPocketDetails Mongo Find() failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return latestPocket, http.StatusInternalServerError, constants.InternalServerError
	}

	return latestPocket, http.StatusOK, ""
}

func placePocketOrder(req models.ExecutePocketV3Request, latestPocket models.MongoLatestPocketDetails, orderSide string, orderObj tradelab.OrderObj, reqH models.ReqHeader) ([]models.CompletedOrderResponseOrders, map[string]models.PocketsMetaData, map[string]bool, int, string) {
	placedOrderIds := make(map[string]bool)
	mapStockDetailsDataAgainstToken := make(map[string]models.PocketsMetaData)
	var userHoldingsMap map[int]int
	var err error

	if strings.EqualFold(orderSide, constants.SELL) {
		userHoldingsMap, err = getAllUserHoldings(req.ClientId, reqH)
		if err != nil {
			return nil, nil, nil, http.StatusInternalServerError, constants.InternalServerError
		}
	}

	for _, StockDetails := range latestPocket.PocketTokens {
		mapStockDetailsDataAgainstToken[StockDetails.Token] = StockDetails

		//calculate the qty of order
		qty, err := strconv.Atoi(StockDetails.Qty)
		if err != nil {
			loggerconfig.Error("placePocketOrder, Conversion error for qty", " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
			return nil, nil, nil, http.StatusInternalServerError, constants.InternalServerError
		}

		if strings.EqualFold(orderSide, constants.SELL) {
			//in case of sell order, check if user has the stock
			pocketStockToken, _ := strconv.Atoi(StockDetails.Token)
			if _, ok := userHoldingsMap[pocketStockToken]; !ok {
				continue
			}
			//if user has the stock, calculate the qty to sell
			qty = min(qty*req.LotSize, userHoldingsMap[pocketStockToken])
		} else {
			//in case of buy order, calculate the qty to buy
			qty *= req.LotSize
		}

		placeOrderReq := createPlaceOrderRequest(req.ClientId, StockDetails.Exchange, StockDetails.Token, orderSide)

		placeOrderReq.Quantity = qty

		loggerconfig.Info("placePocketOrder, PlaceOrderRequest created: ", placeOrderReq, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)

		orderId, httpStatus, errorCode := executePlaceOrder(orderObj, placeOrderReq, req.ClientId, reqH)
		if errorCode != "" {
			return nil, nil, nil, httpStatus, errorCode
		}

		placedOrderIds[orderId] = true
	}

	if len(placedOrderIds) == 0 {
		loggerconfig.Info("placePocketOrder, No orders were placed", " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return []models.CompletedOrderResponseOrders{}, mapStockDetailsDataAgainstToken, placedOrderIds, http.StatusOK, ""
	}

	// Get completed orders
	completedOrders, httpStatus, errorCode := getCompletedOrders(orderObj, req.ClientId, reqH)
	if errorCode != "" {
		return nil, nil, nil, httpStatus, errorCode
	}

	return completedOrders, mapStockDetailsDataAgainstToken, placedOrderIds, http.StatusOK, ""
}

func createPlaceOrderRequest(clientId, exchange, token, orderSide string) models.PlaceOrderRequest {
	var placeOrderReq models.PlaceOrderRequest
	placeOrderReq.ClientID = clientId
	placeOrderReq.DisclosedQuantity = 0
	placeOrderReq.Exchange = exchange
	placeOrderReq.ExecutionType = strings.ToUpper(constants.REGULAR)
	placeOrderReq.InstrumentToken = token
	placeOrderReq.OrderType = strings.ToUpper(constants.MARKET)

	if strings.EqualFold(orderSide, constants.BUY) {
		placeOrderReq.OrderSide = strings.ToUpper(constants.BUY)
	} else {
		placeOrderReq.OrderSide = strings.ToUpper(constants.SELL)
	}

	placeOrderReq.Price = 0.0
	placeOrderReq.Product = strings.ToUpper(constants.CNC)
	placeOrderReq.TriggerPrice = 0.0
	placeOrderReq.Validity = strings.ToUpper(constants.IOC) // complete or reject

	return placeOrderReq
}

func executePlaceOrder(orderObj tradelab.OrderObj, placeOrderReq models.PlaceOrderRequest, clientId string, reqH models.ReqHeader) (string, int, string) {
	status, res := tradelab.OrderObj.PlaceOrder(orderObj, placeOrderReq, reqH)
	if status != http.StatusOK {
		loggerconfig.Error("executePlaceOrder in PlaceOrder status != 200", status, " uccId:", clientId, " requestId:", reqH.RequestId)
		return "", http.StatusInternalServerError, constants.InternalServerError
	}

	placeOrderRes, ok := res.Data.(models.PlaceOrderResponse)
	if !ok {
		loggerconfig.Error("Alert Severity:P1-High, executePlaceOrder in PlaceOrder interface parsing error", ok, " uccId:", clientId, " requestId:", reqH.RequestId)
		return "", http.StatusInternalServerError, constants.InternalServerError
	}

	return placeOrderRes.OmsOrderID, http.StatusOK, ""
}

func getCompletedOrders(orderObj tradelab.OrderObj, clientId string, reqH models.ReqHeader) ([]models.CompletedOrderResponseOrders, int, string) {
	var completedOrderReq models.CompletedOrderRequest
	completedOrderReq.ClientID = clientId
	completedOrderReq.Type = strings.ToLower(constants.Completed)

	// Sleep for a moment to allow the orders to be processed
	time.Sleep(constants.PocketOrderSleepSeconds * time.Second)

	statusCompletedOrder, resCompletedOrder := tradelab.OrderObj.CompletedOrder(orderObj, completedOrderReq, reqH)
	if statusCompletedOrder != http.StatusOK {
		loggerconfig.Error("getCompletedOrders in CompletedOrder status != 200", statusCompletedOrder, " uccId:", clientId, " requestId:", reqH.RequestId)
		return nil, http.StatusInternalServerError, constants.InternalServerError
	}

	loggerconfig.Info("getCompletedOrders CallTLforCompletedOrder response:=", helpers.LogStructAsJSON(resCompletedOrder), "uccId:", clientId, " requestId:", reqH.RequestId)

	completedOrderRes, ok := resCompletedOrder.Data.(models.CompletedOrderResponse)
	if !ok {
		loggerconfig.Error("Alert Severity:P1-High, getCompletedOrders in CompletedOrder interface parsing error", ok, " uccId:", clientId, " requestId:", reqH.RequestId)
		return nil, http.StatusInternalServerError, constants.InternalServerError
	}

	return completedOrderRes.Orders, http.StatusOK, ""
}

func processOrderResults(orderResults []models.CompletedOrderResponseOrders, mapPocketMetaDataAgainstToken map[string]models.PocketsMetaData, placedOrderIds map[string]bool, req models.ProcessOrderReq, reqH models.ReqHeader) (models.PurchagedPocketDetails, []models.PocketsMetaData, []models.PocketsMetaData) {
	var orderCompleted, orderCancelled []models.PocketsMetaData
	orderCompletedPrice := 0.0

	noPlacedOrderIds := len(placedOrderIds)

	for _, order := range orderResults {
		if placedOrderIds[order.OmsOrderID] { // only check the completed orders placed using pocket
			if order.OrderStatus == strings.ToUpper(constants.Complete) {
				orderCompletedPrice += order.Price * float64(order.Quantity)
				pocketData, exists := mapPocketMetaDataAgainstToken[strconv.Itoa(order.InstrumentToken)]
				if exists {
					orderCompleted = append(orderCompleted, pocketData)
				} else {
					loggerconfig.Error("processOrderResults in mapPocketTokenMetaData completed order details don't present in pocket with token", order.InstrumentToken, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
				}
			} else if order.OrderStatus == strings.ToUpper(constants.CancelConfirmed) ||
				order.OrderStatus == strings.ToUpper(constants.Rejected) {
				pocketData, exists := mapPocketMetaDataAgainstToken[strconv.Itoa(order.InstrumentToken)]
				if exists {
					orderCancelled = append(orderCancelled, pocketData)
				} else {
					loggerconfig.Error("processOrderResults in mapPocketTokenMetaData cancelled order details don't present in pocket with token", strconv.Itoa(order.InstrumentToken), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
				}
			}
			noPlacedOrderIds--
		}

		if noPlacedOrderIds == 0 {
			break
		}
	}

	pocketDetailsV3 := models.PurchagedPocketDetails{
		PocketTransactionId: uuid.New().String(),
		PocketId:            req.PocketId,
		TransactionStatus:   req.TransactionStatus, //1 for sold, 0 for buy
		LotSize:             req.LotSize,
		Action:              req.Action,
		OrderCompletedPrice: orderCompletedPrice,
		OrderCompleted:      orderCompleted,
		OrderCancelled:      orderCancelled,
	}

	return pocketDetailsV3, orderCompleted, orderCancelled
}

// updatePocketTransactionStore updates the pocket transaction store
func updatePocketTransactionStore(clientId string, pocketDetailsV3 models.PurchagedPocketDetails, reqH models.ReqHeader) error {
	var pocketTransactionStoreV3 models.PocketTransactionStoreV3
	err := dbops.MongoRepo.FindOne(constants.POCKETSTRANSACTIONSDetailsV2, bson.M{"clientid": clientId}, &pocketTransactionStoreV3)

	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " updatePocketTransactionStore, mongo error =", err, " requestId:", reqH.RequestId)
		return err
	}

	if err != nil && err.Error() == constants.MongoNoDocError {
		pocketTransactionStoreV3.ClientId = clientId
	}

	pocketTransactionStoreV3.AllPocketPurchases = append(pocketTransactionStoreV3.AllPocketPurchases, pocketDetailsV3)

	filter := bson.D{{"clientid", clientId}}
	update := bson.D{{"$set", pocketTransactionStoreV3}}
	opts := options.Update().SetUpsert(true)

	return dbops.MongoRepo.UpdateOne(constants.POCKETSTRANSACTIONSDetailsV2, filter, update, opts)
}

func updateUserPocketHoldings(req models.ProcessOrderReq, PocketVersion int, userPocketHolding models.UserPocketHolding, reqH models.ReqHeader) error {
	pocketIndex := -1

	// Find existing pocket index
	for i, pocket := range userPocketHolding.Pockets {
		if pocket.PocketID == req.PocketId {
			pocketIndex = i
			break
		}
	}

	if strings.EqualFold(req.Action, constants.BUY) {
		if pocketIndex == -1 {
			// New user or pocket not present, add new pocket
			userPocketHolding.Pockets = append(userPocketHolding.Pockets, models.UserPocket{
				PocketID: req.PocketId, Version: PocketVersion, LotSize: req.LotSize,
			})
		} else {
			// Update existing pocket
			userPocketHolding.Pockets[pocketIndex].Version = PocketVersion
			userPocketHolding.Pockets[pocketIndex].LotSize += req.LotSize
		}
	} else if strings.EqualFold(req.Action, constants.SELL) {
		if pocketIndex == -1 {
			loggerconfig.Error("ExecutePocketV3, updateUserPocketHoldings doesn't have the pocket with id", req.PocketId, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
			return errors.New(constants.PocketDoesNotExists)
		}

		if userPocketHolding.Pockets[pocketIndex].LotSize > req.LotSize {
			userPocketHolding.Pockets[pocketIndex].LotSize -= req.LotSize
		} else {
			// Remove pocket if lot size is fully sold
			userPocketHolding.Pockets = append(userPocketHolding.Pockets[:pocketIndex], userPocketHolding.Pockets[pocketIndex+1:]...)
		}
	} else if strings.EqualFold(req.Action, constants.Rebalance) {
		if pocketIndex == -1 {
			loggerconfig.Error("updateUserPocketHoldings doesn't have the pocket with id", req.PocketId, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
			return errors.New(constants.PocketDoesNotExists)
		}

		userPocketHolding.Pockets[pocketIndex].Version = PocketVersion
	}

	filter := bson.D{{"clientId", req.ClientId}}
	update := bson.D{{"$set", bson.D{{"pockets", userPocketHolding.Pockets}}}}
	opts := options.Update().SetUpsert(true)

	err := dbops.MongoRepo.UpdateOne(constants.USERPOCKETSCOLLECTION, filter, update, opts)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " updateUserPocketHoldings, mongo error =", err, " requestId:", reqH.RequestId)
		return err
	}

	return nil
}

func (obj ExecutePocketV3Obj) ExitPocketV3(req models.ExecutePocketV3Request, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	// get latest pocket details
	latestPocket, httpCode, errorCode := getLatestPocketDetails(req.PocketId, reqH)
	if errorCode != "" {
		return httpCode, apihelpers.APIRes{
			Status:  false,
			Message: errorCode,
			Data:    nil,
		}
	}

	var userPocketHolding models.UserPocketHolding
	err := dbops.MongoRepo.FindOne(constants.USERPOCKETSCOLLECTION, bson.M{
		"clientId": req.ClientId,
	}, &userPocketHolding)
	if err != nil {
		if err.Error() == constants.MongoNoDocError {
			return apihelpers.SendErrorResponse(false, constants.PocketDoesNotExists, http.StatusBadRequest)
		}
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " ExitPocketV3, Mongo Find() failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	processOrderReq := models.ProcessOrderReq{
		PocketId:          req.PocketId,
		ClientId:          req.ClientId,
		LotSize:           req.LotSize,
		Action:            constants.SELL,
		TransactionStatus: 1,
	}

	err = updateUserPocketHoldings(processOrderReq, latestPocket.PocketVersion, userPocketHolding, reqH)
	if err != nil {
		return apihelpers.SendInternalServerError()
	}

	return http.StatusOK, apihelpers.APIRes{
		Status:  true,
		Message: "Exit pocket success",
		Data:    nil,
	}
}

func (obj ExecutePocketV3Obj) GetPocketDetails(pocketId, tag string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var pocketData []models.MongoPocketsV3

	// if pocketId is provided
	if pocketId != "" {
		var pocket models.MongoPocketsV3
		err := dbops.MongoRepo.FindOne(constants.POCKETSCOLLECTIONV3, bson.M{"pocketId": pocketId}, &pocket)
		if err != nil {
			if err.Error() == constants.MongoNoDocError {
				return apihelpers.SendErrorResponse(false, constants.PocketDoesNotExists, http.StatusBadRequest)
			}
			return apihelpers.SendErrorResponse(false, err.Error(), http.StatusInternalServerError)
		}

		pocketData = append(pocketData, pocket)

		return http.StatusOK, apihelpers.APIRes{
			Status:  true,
			Message: "SUCCESS",
			Data:    pocketData,
		}

	}

	// if tag is provided
	if tag != "" {
		cursor, err := dbops.MongoRepo.Find(constants.POCKETSCOLLECTIONV3, bson.M{"tag": tag})
		if err != nil {
			return apihelpers.SendErrorResponse(false, err.Error(), http.StatusInternalServerError)
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var pocket models.MongoPocketsV3
			if err := cursor.Decode(&pocket); err != nil {
				return apihelpers.SendErrorResponse(false, err.Error(), http.StatusInternalServerError)
			}
			pocketData = append(pocketData, pocket)
		}

		if len(pocketData) == 0 {
			return apihelpers.SendErrorResponse(false, constants.PocketDoesNotExists, http.StatusBadRequest)
		}
		return http.StatusOK, apihelpers.APIRes{
			Status:  true,
			Message: "SUCCESS",
			Data:    pocketData,
		}
	}

	return http.StatusOK, apihelpers.APIRes{
		Status:  true,
		Message: "SUCCESS",
		Data:    pocketData,
	}
}
