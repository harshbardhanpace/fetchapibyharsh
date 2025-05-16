package v3

import (
	"encoding/json"
	"net/http"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
)

var theExecutePocketV3Provider models.ExecutePocketV3

func InitExecutePocketV3Provider(provider models.ExecutePocketV3) {
	defer models.HandlePanic()
	theExecutePocketV3Provider = provider
}

// BuyPocket
// @Tags space client pockets V3
// @Description Buy a Pocket
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.ExecutePocketV3Request true "pockets"
// @Success 200 {object} apihelpers.APIRes{data=models.PocketOrdersRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/pockets/buyPocket [POST]
func BuyPocketV3(c *gin.Context) {

	var req models.ExecutePocketV3Request

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("BuyPocket V3 (controller), Error in parsing header BuyPocket (Controller), error = ", err, " requestId:", reqH.RequestId)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		loggerconfig.Error("BuyPocket V3 (controller), error decoding body, error:", err, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("BuyPocket V3 (controller), reqParams:", helpers.LogStructAsJSON(req), " uccId: ", req.ClientId, "requestId:", reqH.RequestId)

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(req.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("BuyPocket V3 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", req.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("BuyPocket V3 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", req.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	code, resp := theExecutePocketV3Provider.BuyPocketV3(req, "", reqH)

	logDetail := "clientId: " + reqH.ClientId + " function: BuyPocketV3 requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// CheckActionRequired
// @Tags space client pockets V3
// @Description Check if any action required for pocket
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.CheckActionRequiredReq true "checkActionRequired"
// @Success 200 {object} apihelpers.APIRes{data=models.RebalanceResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/pockets/action/checkActionRequired [POST]
func CheckActionRequired(c *gin.Context) {
	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("checkActionRequired (controller), Error in parsing header BuyPocket (Controller), error = ", err, " requestId:", reqH.RequestId)
		return
	}

	var req models.CheckActionRequiredReq
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		loggerconfig.Error("checkActionRequired (controller), error decoding body, error:", err, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("CheckActionRequired (controller), reqParams:", helpers.LogStructAsJSON(req), " clientId: ", reqH.ClientId, "requestId:", reqH.RequestId)

	code, resp := theExecutePocketV3Provider.CheckActionRequired(req.ClientId, req.PocketId, req.UserVersion, req.UsersLotSize, reqH)

	logDetail := "clientId: " + reqH.ClientId + " function: CheckActionRequired v3 requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// AdjustStocksForPocket
// @Tags space client pockets V3
// @Description complete the action for balancing/repairing of pockets
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.RebalanceResponse true "adjust-pocket"
// @Success 200 {object} apihelpers.APIRes{data=models.PocketOrdersRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/pockets/action/adjustStocks [POST]
func AdjustStocksForPocket(c *gin.Context) {

	var req models.RebalanceResponse

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("AdjustStocksForPocket (controller), Error in parsing header BuyPocket (Controller), error = ", err, " requestId:", reqH.RequestId)
		return
	}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		loggerconfig.Error("AdjustStocksForPocket (controller), error decoding body, error:", err, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("AdjustStocksForPocket (controller), reqParams:", helpers.LogStructAsJSON(req), " clientId: ", req.ClientId, "requestId:", reqH.RequestId)

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(req.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("AdjustStocksForPocket (controller) CheckAuthWithClient invalid authtoken", " clientId: ", req.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("AdjustStocksForPocket (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", req.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	code, resp := theExecutePocketV3Provider.ManageRequiredStocksForPocket(req, reqH)

	logDetail := "clientId: " + reqH.ClientId + " function: AdjustStocksForPocket v3 requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchAllPockets
// @Tags space client pockets V3
// @Description Fetch All Pockets
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Success 200 {object} apihelpers.APIRes{data=[]models.MongoPocketsV3}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/pockets/fetchAllPockets [GET]
func FetchAllPocketsV3(c *gin.Context) {
	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("FetchAllPockets V3 (controller), Error in parsing header BuyPocket (Controller), error = ", err, " requestId:", reqH.RequestId)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqH.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("FetchAllPockets V3 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqH.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("FetchAllPockets V3 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqH.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	code, resp := theExecutePocketV3Provider.FetchAllPocketsV3(reqH)

	logDetail := "clientId: " + reqH.ClientId + " function: FetchAllPockets V3 requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchPocketPortfolio
// @Tags space client pockets V3
// @Description FetchPocketPortfolio - It will provide portfolio details of pocket
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.FetchPocketPortfolioRequest true "pockets"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchPocketPortfolioResponseV2}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/pockets/fetchPocketPortfolio [POST]
func FetchPocketPortfolioV3(c *gin.Context) {
	var req models.FetchPocketPortfolioRequest

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("FetchPocketPortfolio V3 (controller), Error in parsing header FetchPocketPortfolio (Controller), error = ", err, " requestId:", reqH.RequestId)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		loggerconfig.Error("FetchPocketPortfolio V3 (controller), error decoding body, error:", err, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("FetchPocketPortfolio V3 (controller), reqParams:", helpers.LogStructAsJSON(req), " uccId: ", req.ClientId, req, "requestId:", reqH.RequestId)

	//call service
	code, resp := theExecutePocketV3Provider.FetchPocketPortfolioV3(req, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: FetchPocketPortfolio requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchUsersPockets
// @Tags space client pockets V3
// @Description Fetch All Pockets against clientID
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param clientId query string true "clientId"
// @Success 200 {object} apihelpers.APIRes{data=models.UserPocketHolding}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/pockets/fetchUsersPockets [GET]
func FetchUsersPockets(c *gin.Context) {
	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("AdjustStocksForPocket (controller), Error in parsing header BuyPocket (Controller), error = ", err, " requestId:", reqH.RequestId)
		return
	}
	clientID := c.Query("clientId")
	if clientID == "" {
		loggerconfig.Error("AdjustStocksForPocket V2 (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqH.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("AdjustStocksForPocket (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqH.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("AdjustStocksForPocket (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqH.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	logDetail := "clientId: " + reqH.ClientId + " function: AdjustStocksForPocket v3 requestId: " + reqH.RequestId

	code, resp := theExecutePocketV3Provider.FetchUsersPockets(clientID, reqH)
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// SellPocket
// @Tags space client pockets V3
// @Description Sell a Pocket
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.ExecutePocketV3Request true "pockets"
// @Success 200 {object} apihelpers.APIRes{data=models.PocketOrdersRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/pockets/sellPocket [POST]
func SellPocketV3(c *gin.Context) {
	var req models.ExecutePocketV3Request

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("SellPocket V3 (controller), Error in parsing header SellPocket (Controller), error = ", err, " requestId:", reqH.RequestId)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		loggerconfig.Error("SellPocket V3 (controller), error decoding body, error:", err, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("SellPocket V3 (controller), reqParams:", helpers.LogStructAsJSON(req), " uccId: ", req.ClientId, "requestId:", reqH.RequestId)

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(req.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("SellPocket V3 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", req.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("SellPocket V3 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", req.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	code, resp := theExecutePocketV3Provider.SellPocketV3(req, "", reqH)

	logDetail := "clientId: " + reqH.ClientId + " function: SellPocketV3 requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ExitPocket
// @Tags space client pockets V3
// @Description Exit a Pocket
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.ExecutePocketV3Request true "pockets"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/pockets/exitPocket [POST]
func ExitPocketV3(c *gin.Context) {
	var req models.ExecutePocketV3Request

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("ExitPocket V3 (controller), Error in parsing header ExitPocket (Controller), error = ", err, " requestId:", reqH.RequestId)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		loggerconfig.Error("ExitPocket V3 (controller), error decoding body, error:", err, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("ExitPocket V3 (controller), reqParams:", helpers.LogStructAsJSON(req), " uccId: ", req.ClientId, "requestId:", reqH.RequestId)

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(req.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("ExitPocket V3 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", req.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("ExitPocket V3 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", req.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	code, resp := theExecutePocketV3Provider.ExitPocketV3(req, reqH)

	logDetail := "clientId: " + reqH.ClientId + " function: ExitPocketV3 requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GETPOCKET
// @Tags space client pockets V3
// @Description Fetch a Pocket against PocketId
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param pocketId query string false "pocketId"
// @Param tag query string false "tag"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/pockets/getPocketDetails [GET]
func GetPocketDetails(c *gin.Context) {
	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("GetPocketDetails  (controller), Error in parsing header ExitPocket (Controller), error = ", err, " requestId:", reqH.RequestId)
		return
	}

	pocketId := c.Query("pocketId")
	tag := c.Query("tag")
	if pocketId == "" && tag == "" {
		loggerconfig.Error("GetPocketDetails  (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetPocketDetails  (controller), pocketId: ", pocketId, " uccId: ", reqH.ClientId, "requestId:", reqH.RequestId)

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqH.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("GetPocketDetails  (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqH.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("GetPocketDetails  (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqH.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	code, resp := theExecutePocketV3Provider.GetPocketDetails(pocketId, tag, reqH)

	logDetail := "clientId: " + reqH.ClientId + " function: GetPocketDetails requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
