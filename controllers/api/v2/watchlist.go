package v2

import (
	"encoding/json"
	"net/http"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var watchListProviderV2 models.WatchListProvider

func InitWatchListProviderV2(provider models.WatchListProvider) {
	defer models.HandlePanic()
	watchListProviderV2 = provider
}

// FetchWatchList
// @Tags space watchlist V2
// @Description modify watch list
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.FetchWatchListV2Request true "watchlist"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchWatchListV2Response}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/watchlist/fetchWatchList [POST]
func FetchWatchList(c *gin.Context) {
	var fetchWatchListReq models.FetchWatchListV2Request

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("FetchWatchList V2 (controller), error decoding header, error:", err)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&fetchWatchListReq)
	if err != nil {
		loggerconfig.Error("FetchWatchList V2 (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(fetchWatchListReq.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("FetchWatchList V2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", fetchWatchListReq.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("FetchWatchList V2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", fetchWatchListReq.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	//call service
	code, resp := watchListProviderV2.FetchWatchListsV2(fetchWatchListReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: FetchWatchList v2 requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// AddStockToWatchList
// @Tags space watchlist V2
// @Description Add Stock To WatchList
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.AddStockToWatchListV2Request true "watchlist"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/watchlist/addStockToWatchList [POST]
func AddStockToWatchList(c *gin.Context) {
	var addWatchListReq models.AddStockToWatchListV2Request

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("AddStockToWatchList (controller), error decoding header, error:", err)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&addWatchListReq)
	if err != nil {
		loggerconfig.Error("AddStockToWatchList (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(addWatchListReq.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("AddStockToWatchList V2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", addWatchListReq.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("FetchWatchList V2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", addWatchListReq.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	validate := validator.New()
	if err := validate.Struct(addWatchListReq); err != nil {
		loggerconfig.Error("ddStockToWatchList (controller) validation error: ", err, " requestid: ", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	//call service
	code, resp := watchListProviderV2.AddStockToWatchListV2(addWatchListReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: AddStockToWatchList v2 requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// AddStockToWatchList
// @Tags space watchlist V2
// @Description Add Stock To WatchList Concise
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.AddStockToWatchListConciseV2Request true "watchlist"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/watchlist/addStockToWatchListConcise [POST]
// func AddStockToWatchListConcise(c *gin.Context) {
// 	var addWatchListConciseReq models.AddStockToWatchListConciseV2Request

// 	var reqH models.ReqHeader
// 	if err := c.ShouldBindHeader(&reqH); err != nil {
// 		fmt.Printf("error parsing header=%v\n", err)
// 	}

// 	err := json.NewDecoder(c.Request.Body).Decode(&addWatchListConciseReq)
// 	if err != nil {
// 		apihelpers.Respond(c.Writer, apihelpers.Message(1, "Invalid request"))
// 		return
// 	}

// 	var addWatchList models.AddStockToWatchListV2Request
// 	var Collections models.MongoCollections
// 	collection := models.GetMongoCollection(constants.StocksContracts)
// 	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

// 	for i := 0; i < len(addWatchListConciseReq.StockData); i++ {
// 		err = collection.FindOne(ctx, bson.M{"collectionId": modifyCollectionsReq.CollectionId}).Decode(&Collections)
// 	}

// 	if addWatchListReq.WatchListId != "wl1" && addWatchListReq.WatchListId != "wl2" &&
// 		addWatchListReq.WatchListId != "wl3" &&
// 		addWatchListReq.WatchListId != "wl4" &&
// 		addWatchListReq.WatchListId != "wl5" {
// 		apihelpers.Respond(c.Writer, apihelpers.Message(1, "Invalid WatchListId"))
// 		return
// 	}

// 	//call service
// 	code, resp := watchListProviderV2.AddStockToWatchListV2(addWatchListReq, reqH)

// 	//return response using api helper
// 	apihelpers.CustomResponse(c, code, resp)
// }

// DeleteStockInWatchList
// @Tags space watchlist V2
// @Description Delete Stock In WatchList
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.DeleteWatchListV2Request true "watchlist"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/watchlist/deleteStockInWatchList [POST]
func DeleteStockInWatchList(c *gin.Context) {
	var delWatchListReq models.DeleteWatchListV2Request

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("DeleteStockInWatchList V2 (controller), error decoding header, error:", err)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&delWatchListReq)
	if err != nil {
		loggerconfig.Error("DeleteStockInWatchList V2 (controller), error decoding body, error:", err, "clientID: ", delWatchListReq.ClientId, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(delWatchListReq.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("DeleteStockInWatchList V2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", delWatchListReq.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("DeleteStockInWatchList V2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", delWatchListReq.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	if delWatchListReq.WatchListId != "wl1" && delWatchListReq.WatchListId != "wl2" &&
		delWatchListReq.WatchListId != "wl3" &&
		delWatchListReq.WatchListId != "wl4" &&
		delWatchListReq.WatchListId != "wl5" {
		loggerconfig.Error("AddStockToWatchList (controller), invalid watchlist")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	//call service
	code, resp := watchListProviderV2.DeleteStockInWatchListV2(delWatchListReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: DeleteStockInWatchList v2 requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}

// DeleteStockInWatchListUpdated
// @Tags space watchlist V2
// @Description Delete Stock In WatchList Updated
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.DeleteWatchListV2UpdatedRequest true "watchlist"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/watchlist/deleteStockInWatchListUpdated [POST]
func DeleteStockInWatchListUpdated(c *gin.Context) {
	var delWatchListReq models.DeleteWatchListV2UpdatedRequest

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("DeleteStockInWatchListUpdated (controller), error decoding header, error:", err)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&delWatchListReq)
	if err != nil {
		loggerconfig.Error("DeleteStockInWatchListUpdated (controller), error decoding body, error:", err, "clientID: ", delWatchListReq.ClientId, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(delWatchListReq.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("DeleteStockInWatchListUpdated V2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", delWatchListReq.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("DeleteStockInWatchListUpdated V2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", delWatchListReq.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	//call service
	code, resp := watchListProviderV2.DeleteStockInWatchListV2Updated(delWatchListReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: DeleteStockInWatchListUpdated v2 requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}

// ArrangeStocksWatchList
// @Tags space watchlist V2
// @Description Arrange Stocks WatchList
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.ArrangeStocksWatchListV2Request true "watchlist"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/watchlist/arrangeStocksWatchList [POST]
func ArrangeStocksWatchList(c *gin.Context) {
	var arrangeStocksWatchListReq models.ArrangeStocksWatchListV2Request

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("Error in parsing ArrangeStocksWatchList (Controller), error = ", err, " requestId:", reqH.RequestId)
		loggerconfig.Error("ArrangeStocksWatchList (controller), error decoding header, error:", err)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&arrangeStocksWatchListReq)
	if err != nil {
		loggerconfig.Error("ArrangeStocksWatchList (controller), Error in decode ArrangeStocksWatchList (Controller), error = ", err, "clientID: ", arrangeStocksWatchListReq.ClientId, " , requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(arrangeStocksWatchListReq.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("ArrangeStocksWatchList V2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", arrangeStocksWatchListReq.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("ArrangeStocksWatchList V2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", arrangeStocksWatchListReq.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	if arrangeStocksWatchListReq.WatchListId != "wl1" && arrangeStocksWatchListReq.WatchListId != "wl2" &&
		arrangeStocksWatchListReq.WatchListId != "wl3" &&
		arrangeStocksWatchListReq.WatchListId != "wl4" &&
		arrangeStocksWatchListReq.WatchListId != "wl5" {
		loggerconfig.Error("ArrangeStocksWatchList (controller), Invalid WatchListId ArrangeStocksWatchList (Controller), clientId :", arrangeStocksWatchListReq.ClientId, " requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	//call service
	code, resp := watchListProviderV2.ArrangeStocksWatchListV2(arrangeStocksWatchListReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: ArrangeStocksWatchListV2 v2 requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}
