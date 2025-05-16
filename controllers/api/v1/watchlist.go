package v1

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

var watchListProvider models.WatchListProvider

func InitWatchListProviderV1(provider models.WatchListProvider) {
	defer models.HandlePanic()
	watchListProvider = provider
}

// CreateWatchList
// @Tags space watchlist V1
// @Description create watch list
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.CreateWatchListRequest true "funds"
// @Success 200 {object} apihelpers.APIRes{data=models.CreateWatchListResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/watchlist/createWatchlist [POST]
func CreateWatchList(c *gin.Context) {
	var watchListReq models.CreateWatchListRequest
	cRH, _ := c.Get("reqH")
	reqH, _ := (cRH).(models.ReqHeader)

	err := json.NewDecoder(c.Request.Body).Decode(&watchListReq)
	if err != nil {
		loggerconfig.Error("CreateWatchList (controller), Invalid Request, error:", err, " requestId: ", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(watchListReq.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("CreateWatchList (controller) CheckAuthWithClient invalid authtoken", " clientId: ", watchListReq.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("CreateWatchList (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", watchListReq.ClientId, "requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("CreateWatchList (controller), reqParams:", watchListReq, " uccId: ", watchListReq.ClientId, "requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	//call service
	code, resp := watchListProvider.CreateWatchList(watchListReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: CreateWatchList requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}

// ModifyWatchList
// @Tags space watchlist V1
// @Description Modify Watchlist
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.ModifyWatchListRequest true "funds"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/watchlist/modifyWatchList [POST]
func ModifyWatchList(c *gin.Context) {
	var modWatchListReq models.ModifyWatchListRequest
	cRH, _ := c.Get("reqH")
	reqH, _ := (cRH).(models.ReqHeader)

	err := json.NewDecoder(c.Request.Body).Decode(&modWatchListReq)
	if err != nil {
		loggerconfig.Error("ModifyWatchList (controller), Invalid Request, error:", err, "  requestId: ", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("ModifyWatchList (controller), reqParams:", modWatchListReq, "clientID: ", reqH.ClientId, "requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	//call service
	code, resp := watchListProvider.ModifyWatchList(modWatchListReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: ModifyWatchList requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}

// FetchWatchList
// @Tags space watchlist V1
// @Description Fetch WatchList
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.FetchWatchListsRequest true "funds"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchWatchListResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/watchlist/fetchWatchList [POST]
func FetchWatchList(c *gin.Context) {
	var fetchWatchListReq models.FetchWatchListsRequest
	cRH, _ := c.Get("reqH")
	reqH, _ := (cRH).(models.ReqHeader)

	err := json.NewDecoder(c.Request.Body).Decode(&fetchWatchListReq)
	if err != nil {
		loggerconfig.Error("FetchWatchList (controller), Invalid Request, error:", err, "  requestId: ", reqH.RequestId, "ClientID: ", fetchWatchListReq.ClientId, " clientVersion:", reqH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("FetchWatchList (controller), reqParams:", helpers.LogStructAsJSON(fetchWatchListReq), " uccId: ", fetchWatchListReq.ClientId, "requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(fetchWatchListReq.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("FetchWatchList (controller) CheckAuthWithClient invalid authtoken", " clientId: ", fetchWatchListReq.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("FetchWatchList (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", fetchWatchListReq.ClientId, "requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	//call service
	code, resp := watchListProvider.FetchWatchLists(fetchWatchListReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: FetchWatchList requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DeleteWatchList
// @Tags space watchlist V1
// @Description delete watch list
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.DeleteWatchListRequest true "funds"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/watchlist/deleteWatchList [POST]
func DeleteWatchList(c *gin.Context) {
	var deleteWatchListReq models.DeleteWatchListRequest
	cRH, _ := c.Get("reqH")
	reqH, _ := (cRH).(models.ReqHeader)

	err := json.NewDecoder(c.Request.Body).Decode(&deleteWatchListReq)
	if err != nil {
		loggerconfig.Error("DeleteWatchList (controller), Invalid Request, error:", err, "  requestId: ", reqH.RequestId, "ClientID: ", deleteWatchListReq.ClientId, " clientVersion:", reqH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(deleteWatchListReq.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("DeleteWatchList (controller) CheckAuthWithClient invalid authtoken", " clientId: ", deleteWatchListReq.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("DeleteWatchList (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", deleteWatchListReq.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("DeleteWatchList (controller), reqParams:", helpers.LogStructAsJSON(deleteWatchListReq), " uccId: ", deleteWatchListReq.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	//call service
	code, resp := watchListProvider.DeleteWatchList(deleteWatchListReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: DeleteWatchList requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}

// FetchWatchListDetails
// @Tags space watchlist V1
// @Description fetch watch list
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.FetchWatchListsDetailsRequest true "funds"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchWatchListsDetailsResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/watchlist/fetchWatchListDetails [POST]
func FetchWatchListDetails(c *gin.Context) {
	var fetchWatchListDetailsReq models.FetchWatchListsDetailsRequest
	cRH, _ := c.Get("reqH")
	reqH, _ := (cRH).(models.ReqHeader)

	err := json.NewDecoder(c.Request.Body).Decode(&fetchWatchListDetailsReq)
	if err != nil {
		loggerconfig.Error("FetchWatchListDetails (controller), Invalid Request, error:", err, "  requestId: ", reqH.RequestId, "ClientID: ", fetchWatchListDetailsReq.ClientId, " clientVersion:", reqH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("FetchWatchListDetails (controller), reqParams:", helpers.LogStructAsJSON(fetchWatchListDetailsReq), " uccId: ", fetchWatchListDetailsReq.ClientId, "requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(fetchWatchListDetailsReq.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("FetchWatchListDetails (controller) CheckAuthWithClient invalid authtoken", " clientId: ", fetchWatchListDetailsReq.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("FetchWatchListDetails (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", fetchWatchListDetailsReq.ClientId, "requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	//call service
	code, resp := watchListProvider.FetchWatchListDetails(fetchWatchListDetailsReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: FetchWatchListDetails requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}

// AddStockToWatchList
// @Tags space watchlist V1
// @Description add stocks to watch list
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.AddStockToWatchListsRequest true "funds"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/watchlist/addStockToWatchList [POST]
func AddStockToWatchList(c *gin.Context) {
	var addStockToWatchListDetailsReq models.AddStockToWatchListsRequest
	cRH, _ := c.Get("reqH")
	reqH, _ := (cRH).(models.ReqHeader)

	err := json.NewDecoder(c.Request.Body).Decode(&addStockToWatchListDetailsReq)
	if err != nil {
		loggerconfig.Error("AddStockToWatchList (controller), Invalid Request, error:", err, "  requestId: ", reqH.RequestId, "ClientID: ", addStockToWatchListDetailsReq.ClientId, " clientVersion:", reqH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(addStockToWatchListDetailsReq.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("AddStockToWatchList (controller) CheckAuthWithClient invalid authtoken", " clientId: ", addStockToWatchListDetailsReq.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("AddStockToWatchList (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", addStockToWatchListDetailsReq.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}
	addStockToWatchListDetailsReq.Stock.Exchange = strings.ToUpper(addStockToWatchListDetailsReq.Stock.Exchange)
	validate := validator.New()
	err = validate.Struct(addStockToWatchListDetailsReq)
	if err != nil {
		loggerconfig.Error("AddStockToWatchList (controller), Error validating struct: ", err, " requestId: ", reqH.RequestId, "clientId: ", addStockToWatchListDetailsReq.ClientId, " deviceId: ", reqH.DeviceId, " clientVersion:", reqH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("AddStockToWatchList (controller), reqParams:", helpers.LogStructAsJSON(addStockToWatchListDetailsReq), " uccId: ", addStockToWatchListDetailsReq.ClientId, "requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	//call service
	code, resp := watchListProvider.AddStockToWatchList(addStockToWatchListDetailsReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: AddStockToWatchList requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
