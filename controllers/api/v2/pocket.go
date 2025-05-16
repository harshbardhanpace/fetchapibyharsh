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
)

var theExecutePocketV2Provider models.ExecutePocketV2

func InitExecutePocketV2Provider(provider models.ExecutePocketV2) {
	defer models.HandlePanic()
	theExecutePocketV2Provider = provider
}

// BuyPocket
// @Tags space client pockets V2
// @Description Buy a Pocket
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.ExecutePocketV2Request true "pockets"
// @Success 200 {object} apihelpers.APIRes{data=models.PocketDetailsV2}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/pockets/buyPocket [POST]
func BuyPocket(c *gin.Context) {

	var req models.ExecutePocketV2Request

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("BuyPocket V2 (controller), Error in parsing header BuyPocket (Controller), error = ", err, " requestId:", reqH.RequestId)
		return
	}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		loggerconfig.Error("BuyPocket V2 (controller), error decoding body, error:", err, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(req.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("BuyPocket V2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", req.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("BuyPocket V2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", req.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("BuyPocket V2 (controller), reqParams:", helpers.LogStructAsJSON(req), " uccId: ", req.ClientId, "requestId:", reqH.RequestId)
	//call service
	code, resp := theExecutePocketV2Provider.BuyPocketV2(req, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: BuyPocketV2 requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ExitPocket
// @Tags space client pockets V2
// @Description Exit a Pocket
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.ExecutePocketV2Request true "pockets"
// @Success 200 {object} apihelpers.APIRes{data=models.PocketDetailsV2}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/pockets/exitPocket [POST]
func ExitPocket(c *gin.Context) {
	var req models.ExecutePocketV2Request

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("ExitPocket V2 (controller), Error in parsing header ExitPocket (Controller), error = ", err, " requestId:", reqH.RequestId)
		return
	}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		loggerconfig.Error("ExitPocket V2 (controller), error decoding body, error:", err, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(req.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("ExitPocket V2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", req.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("ExitPocket V2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", req.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("ExitPocket V2 (controller), reqParams:", helpers.LogStructAsJSON(req), " uccId: ", req.ClientId, "requestId:", reqH.RequestId)
	//call service
	code, resp := theExecutePocketV2Provider.ExitPocketV2(req, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: ExitPocketV2 requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchPocketPortfolio
// @Tags space client pockets V2
// @Description Fetch Pocket Portfolio
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.FetchPocketPortfolioRequest true "pockets"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchPocketPortfolioResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/pockets/fetchPocketPortfolio [POST]
func FetchPocketPortfolio(c *gin.Context) {
	var req models.FetchPocketPortfolioRequest

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("FetchPocketPortfolio V2 (controller), Error in parsing header ExitPocket (Controller), error = ", err, " requestId:", reqH.RequestId)
		return
	}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		loggerconfig.Error("FetchPocketPortfolio V2 (controller), error decoding body, error:", err, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(req.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("FetchPocketPortfolio V2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", req.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("FetchPocketPortfolio V2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", req.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("FetchPocketPortfolio V2 (controller), reqParams:", helpers.LogStructAsJSON(req), " uccId: ", req.ClientId, "requestId:", reqH.RequestId)
	//call service
	code, resp := theExecutePocketV2Provider.FetchPocketPortfolioV2(req, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: FetchPocketPortfolio V2 requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
