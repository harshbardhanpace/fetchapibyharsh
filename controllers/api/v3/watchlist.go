package v3

import (
	"encoding/json"
	"net/http"
	"strings"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var watchListProviderV3 models.WatchListProvider

func InitWatchListProviderV3(provider models.WatchListProvider) {
	defer models.HandlePanic()
	watchListProviderV3 = provider
	go watchListProviderV3.PopulateIsinMappingInLocalCache()
}

// FetchWatchList
// @Tags space watchlist V3
// @Description modify watch list
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param clientId query string true "clientId Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.FetchWatchListV3Response}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/watchlist/fetchWatchList [GET]
func FetchWatchList(c *gin.Context) {
	var reqParams models.FetchWatchListV3Request

	clientID := c.Query("clientId")
	if clientID == "" {
		loggerconfig.Error("FetchWatchList (controller), error parsing the query params in Get request")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.ClientId = strings.ToUpper(clientID)

	cRH, _ := c.Get("reqH")
	reqH, _ := (cRH).(models.ReqHeader)

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("FetchWatchList V3 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}

	if !matchStatus {
		loggerconfig.Error("FetchWatchList V3 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	//call service
	code, resp := watchListProviderV3.FetchWatchListsV3(reqParams, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: FetchWatchList v3 requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// AddStockToWatchList
// @Tags space watchlist V3
// @Description Add Stock To WatchList
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.AddStockToWatchListV3Request true "watchlist"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/watchlist/addStockToWatchList [POST]
func AddStockToWatchList(c *gin.Context) {
	var addWatchListReq models.AddStockToWatchListV3Request

	cRH, _ := c.Get("reqH")
	reqH, _ := (cRH).(models.ReqHeader)

	err := json.NewDecoder(c.Request.Body).Decode(&addWatchListReq)
	if err != nil {
		loggerconfig.Error("AddStockToWatchList V3 (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	addWatchListReq.ClientId = strings.ToUpper(addWatchListReq.ClientId)
	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(addWatchListReq.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("AddStockToWatchList V3 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", addWatchListReq.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("FetchWatchList V3 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", addWatchListReq.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	validate := validator.New()
	if err := validate.Struct(addWatchListReq); err != nil {
		loggerconfig.Error("ddStockToWatchList V3 (controller) validation error: ", err, " requestid: ", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	//call service
	code, resp := watchListProviderV3.AddStockToWatchListV3(addWatchListReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: AddStockToWatchList v3 requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DeleteStockInWatchList
// @Tags space watchlist V3
// @Description User can delete multiple stocks from a single watchlist
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.DeleteWatchListV3Request true "watchlist"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/watchlist/deleteStockInWatchList [DELETE]
func DeleteStockInWatchList(c *gin.Context) {
	var delWatchListReq models.DeleteWatchListV3Request

	cRH, _ := c.Get("reqH")
	reqH, _ := (cRH).(models.ReqHeader)

	err := json.NewDecoder(c.Request.Body).Decode(&delWatchListReq)
	if err != nil {
		loggerconfig.Error("DeleteStockInWatchList V3 (controller), error decoding body, error:", err, "clientID: ", delWatchListReq.ClientId, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	delWatchListReq.ClientId = strings.ToUpper(delWatchListReq.ClientId)
	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(delWatchListReq.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("DeleteStockInWatchList V3 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", delWatchListReq.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("DeleteStockInWatchList V3 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", delWatchListReq.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	_, found := constants.WatchlistIds[delWatchListReq.WatchListId]
	if !found {
		loggerconfig.Error("AddStockToWatchList (controller), invalid watchlist")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if delWatchListReq.StockId == nil || len(delWatchListReq.StockId) == 0 || delWatchListReq.StockId[0] == "" {
		loggerconfig.Error("DeleteStockInWatchList V3 (controller), stocks are empty")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	//call service
	code, resp := watchListProviderV3.DeleteStockInWatchListV3(delWatchListReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: DeleteStockInWatchList v3 requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}

// ArrangeStocksWatchList
// @Tags space watchlist V3
// @Description Arrange Stocks WatchList
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.ArrangeStocksWatchListV3Request true "watchlist"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/watchlist/arrangeStocksWatchList [POST]
func ArrangeStocksWatchList(c *gin.Context) {
	var arrangeStocksWatchListReq models.ArrangeStocksWatchListV3Request

	cRH, _ := c.Get("reqH")
	reqH, _ := (cRH).(models.ReqHeader)

	err := json.NewDecoder(c.Request.Body).Decode(&arrangeStocksWatchListReq)
	if err != nil {
		loggerconfig.Error("ArrangeStocksWatchList V3 (controller), Error in decode ArrangeStocksWatchList (Controller), error = ", err, "clientID: ", arrangeStocksWatchListReq.ClientId, " , requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(arrangeStocksWatchListReq.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("ArrangeStocksWatchList V3 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", arrangeStocksWatchListReq.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("ArrangeStocksWatchList V3 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", arrangeStocksWatchListReq.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	_, found := constants.WatchlistIds[arrangeStocksWatchListReq.WatchListId]
	if !found {
		loggerconfig.Error("ArrangeStocksWatchList V3 (controller), Invalid WatchListId ArrangeStocksWatchList (Controller), clientId :", arrangeStocksWatchListReq.ClientId, " requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	//call service
	code, resp := watchListProviderV3.ArrangeStocksWatchListV3(arrangeStocksWatchListReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: ArrangeStocksWatchList v3 requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}

// DeleteStockInWatchListUpdated
// @Tags space watchlist V3
// @Description User can delete multiple stocks from multiple watchlist
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.DeleteWatchListV3UpdatedRequest true "watchlist"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/watchlist/deleteStockInWatchListUpdated [DELETE]
func DeleteStockInWatchListUpdated(c *gin.Context) {
	var delWatchListReq models.DeleteWatchListV3UpdatedRequest

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("DeleteStockInWatchListUpdated V3 (controller), error decoding header, error:", err)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&delWatchListReq)
	if err != nil {
		loggerconfig.Error("DeleteStockInWatchListUpdated V3 (controller), error decoding body, error:", err, "clientID: ", delWatchListReq.ClientId, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(delWatchListReq.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("DeleteStockInWatchListUpdated V3 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", delWatchListReq.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("DeleteStockInWatchListUpdated V3 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", delWatchListReq.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	//call service
	code, resp := watchListProviderV3.DeleteStockInWatchListV3Updated(delWatchListReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: DeleteStockInWatchListUpdated v3 requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}
