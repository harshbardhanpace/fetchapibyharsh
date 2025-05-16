package v1

import (
	"encoding/json"
	"strings"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
)

var pocketsProvider models.PocketsProvider

func InitPocketsProvider(provider models.PocketsProvider) {
	defer models.HandlePanic()
	pocketsProvider = provider
}

// AdminLogin
// @Tags space admin login V1
// @Description AdminLogin - Admin can login with help of userId and password
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param request body models.AdminLoginRequest true "pockets"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Success 200 {object} apihelpers.APIRes{data=models.AdminLoginResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/adminapis/adminLogin [POST]
func AdminLogin(c *gin.Context) {
	var login models.AdminLoginRequest

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("AdminLogin (controller), error decoding header, error:", err)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&login)
	if err != nil {
		loggerconfig.Error("AdminLogin (controller), error decoding body, error:", err, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if login.Password == "" || login.UserId == "" {
		apihelpers.ErrorMessage(c, constants.InvalidParameters)
		return
	}

	maskedReq, err := maskObj.Struct(login)
	if err != nil {
		loggerconfig.Error("In Controller AdminLogin Error in masking request err: ", err, " clientId: ", login.UserId, " requestid = ", reqH.RequestId)
		return
	}

	loggerconfig.Info("AdminLogin (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), "requestId:", reqH.RequestId, "userID: ", login.UserId)

	//call service
	code, resp := pocketsProvider.AdminLogin(login, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: AdminLogin requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// CreatePockets
// @Tags space admin pockets V1
// @Description CreatePockets - Create pocket by providing its details
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.CreatePocketsRequest true "pockets"
// @Success 200 {object} apihelpers.APIRes{data=models.CreatePocketsResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/adminapis/createPockets [POST]
func CreatePockets(c *gin.Context) {
	var pocketsReq models.CreatePocketsRequest

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("CreatePockets (controller), error decoding header, error:", err)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&pocketsReq)
	if err != nil {
		loggerconfig.Error("CreatePockets (controller), error decoding body, error:", err, "clientID: ", reqH.ClientId, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("CreatePockets (controller), reqParams:", helpers.LogStructAsJSON(pocketsReq), "clientID: ", reqH.ClientId, "requestId:", reqH.RequestId)

	//call service
	code, resp := pocketsProvider.CreatePockets(pocketsReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: CreatePockets requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// Pockets
// @Tags space admin pockets V1
// @Description ModifyPockets - Modify the already created pocket
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.ModifyPocketsRequest true "pockets"
// @Success 200 {object} apihelpers.APIRes{data=models.ModifyPocketsResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/adminapis/modifyPockets [POST]
func ModifyPockets(c *gin.Context) {
	var pocketsReq models.ModifyPocketsRequest

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("ModifyPockets (controller), error decoding header, error:", err)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&pocketsReq)
	if err != nil {
		loggerconfig.Error("ModifyPockets (controller), error decoding body, error:", err, "clientID: ", reqH.ClientId, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("ModifyPockets (controller), reqParams:", helpers.LogStructAsJSON(pocketsReq), "clientID: ", reqH.ClientId, "requestId:", reqH.RequestId)

	//call service
	code, resp := pocketsProvider.ModifyPockets(pocketsReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: ModifyPockets requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// Pockets
// @Tags space admin pockets V1
// @Description FetchPockets - Fetch the pocket details by providing pocketId
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.FetchPocketsDetailsRequest true "pockets"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchPocketsDetailsResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/adminapis/fetchPockets [POST]
func FetchPockets(c *gin.Context) {
	var pocketsReq models.FetchPocketsDetailsRequest

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("FetchPockets (controller), error decoding header, error:", err)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&pocketsReq)
	if err != nil {
		loggerconfig.Error("FetchPockets (controller), error decoding body, error:", err, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("FetchPockets (controller), reqParams:", helpers.LogStructAsJSON(pocketsReq), "clientID: ", reqH.ClientId, "requestId:", reqH.RequestId, "PocketId: ", pocketsReq.PocketId)

	//call service
	code, resp := pocketsProvider.FetchPockets(pocketsReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: FetchPockets requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// Pockets User
// @Tags space user pockets V1
// @Description FetchPocketsUser - User api for fetching pocket
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.FetchPocketsDetailsRequest true "pockets"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchPocketsDetailsResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/pockets/fetchPockets [POST]
func FetchPocketsUser(c *gin.Context) {
	var pocketsReq models.FetchPocketsDetailsRequest

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("FetchPocketsUser (controller), error decoding header, error:", err)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&pocketsReq)
	if err != nil {
		loggerconfig.Error("FetchPocketsUser (controller), error decoding body, error:", err, "requestId:", reqH.RequestId, "PocketId: ", pocketsReq.PocketId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("FetchPocketsUser (controller), reqParams:", helpers.LogStructAsJSON(pocketsReq), "requestId:", reqH.RequestId, "PocketId: ", pocketsReq.PocketId)

	//call service
	code, resp := pocketsProvider.FetchPockets(pocketsReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: FetchPocketsUser requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// Pockets
// @Tags space admin pockets V1
// @Description Delete Pocket
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.DeletePocketsRequest true "pockets"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/adminapis/deletePockets [POST]
func DeletePockets(c *gin.Context) {
	var pocketsReq models.DeletePocketsRequest

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("DeletePockets (controller), error decoding header, error:", err)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&pocketsReq)
	if err != nil {
		loggerconfig.Error("DeletePockets (controller), error decoding body, error:", err, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("DeletePockets (controller), reqParams:", helpers.LogStructAsJSON(pocketsReq), "requestId:", reqH.RequestId, "PocketId: ", pocketsReq.PocketId)

	//call service
	code, resp := pocketsProvider.DeletePockets(pocketsReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: DeletePockets requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// Pockets
// @Tags space admin pockets V1
// @Description Fetch All Pockets
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchAllPocketsDetailsResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/adminapis/fetchAllPockets [GET]
func FetchAllPockets(c *gin.Context) {

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("FetchAllPockets (controller), error decoding header, error:", err)
	}

	//call service
	code, resp := pocketsProvider.FetchAllPockets(reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: FetchAllPockets requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// Pockets User
// @Tags space user pockets V1
// @Description FetchAllPocketsUser - It provides details of all pockets to user
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchAllPocketsDetailsResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/pockets/fetchAllPockets [GET]
func FetchAllPocketsUser(c *gin.Context) {

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("FetchAllPocketsUser (controller), error decoding header, error:", err)
	}

	//call service
	code, resp := pocketsProvider.FetchAllPockets(reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: FetchAllPocketsUser requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchPocketPortfolio
// @Tags space client pockets V1
// @Description FetchPocketPortfolio - It will provide portfolio details of pocket
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
// @Router /api/space/v1/pockets/fetchPocketPortfolio [POST]
func FetchPocketPortfolio(c *gin.Context) {
	var req models.FetchPocketPortfolioRequest

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("FetchPocketPortfolio (controller), Error in parsing header FetchPocketPortfolio (Controller), error = ", err, " requestId:", reqH.RequestId)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		loggerconfig.Error("FetchPocketPortfolio (controller), error decoding body, error:", err, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("FetchPocketPortfolio (controller), reqParams:", helpers.LogStructAsJSON(req), " uccId: ", req.ClientId, req, "requestId:", reqH.RequestId)

	//call service
	code, resp := pocketsProvider.FetchPocketPortfolio(req, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: FetchPocketPortfolio requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// BuyPocket
// @Tags space client pockets V1
// @Description Buy a Pocket
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.ExecutePocketRequest true "pockets"
// @Success 200 {object} apihelpers.APIRes{data=models.ExecutePocketResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/pockets/buyPocket [POST]
func BuyPocket(c *gin.Context) {

	var req models.ExecutePocketRequest

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("BuyPocket (controller), Error in parsing header BuyPocket (Controller), error = ", err, " requestId:", reqH.RequestId)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		loggerconfig.Error("BuyPocket (controller), error decoding body, error:", err, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	operation := strings.ToUpper(constants.BUY)

	loggerconfig.Info("BuyPocket (controller), reqParams:", helpers.LogStructAsJSON(req), " uccId: ", req.ClientId, "requestId:", reqH.RequestId)
	//call service
	code, resp := pocketsProvider.ExecutePocket(req, operation, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: BuyPocket requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ExitPocket
// @Tags space client pockets V1
// @Description Exit a Pocket
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.ExecutePocketRequest true "pockets"
// @Success 200 {object} apihelpers.APIRes{data=models.ExecutePocketResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/pockets/exitPocket [POST]
func ExitPocket(c *gin.Context) {
	var req models.ExecutePocketRequest

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("ExitPocket (controller), Error in parsing header ExitPocket (Controller), error = ", err, " requestId:", reqH.RequestId)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		loggerconfig.Error("ExitPocket (controller), error decoding body, error:", err, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	operation := strings.ToUpper(constants.SELL)

	loggerconfig.Info("ExitPocket (controller), reqParams:", helpers.LogStructAsJSON(req), " uccId: ", req.ClientId, "requestId:", reqH.RequestId)
	//call service
	code, resp := pocketsProvider.ExecutePocket(req, operation, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: ExitPocket requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchPocketPortfolio
// @Tags space client pockets V1
// @Description Fetch Pocket for a client
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.FetchPocketTransactionReq true "pockets"
// @Success 200 {object} apihelpers.APIRes{data=models.PocketTransactionComplete}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/pockets/fetchPocketTransaction [POST]
func FetchPocketTransaction(c *gin.Context) {
	var req models.FetchPocketTransactionReq

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("FetchPocketTransaction (controller), Error in parsing header FetchPocketTransaction (Controller), error = ", err, " requestId:", reqH.RequestId)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		loggerconfig.Error("FetchPocketTransaction (controller), error decoding body, error:", err, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("FetchPocketTransaction (controller), reqParams:", helpers.LogStructAsJSON(req), " uccId: ", req.ClientId, "requestId:", reqH.RequestId)

	//call service
	code, resp := pocketsProvider.FetchPocketTransaction(req, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: FetchPocketTransaction requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// Pockets User
// @Tags space user pockets V1
// @Description Pockets Calculations
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PocketsCalculationsReq true "PocketsCalculations"
// @Success 200 {object} apihelpers.APIRes{data=models.PocketsCalculationsRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/pockets/pocketsCalculations [POST]
func PocketsCalculations(c *gin.Context) {
	var pocketsReq models.PocketsCalculationsReq

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("PocketsCalculations (controller), error decoding header, error:", err)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&pocketsReq)
	if err != nil {
		loggerconfig.Error("PocketsCalculations (controller), error decoding body, error:", err, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("PocketsCalculations (controller), reqParams:", helpers.LogStructAsJSON(pocketsReq), "clientID: ", reqH.ClientId, "requestId:", reqH.RequestId, "PcketId: ", pocketsReq.PocketId)

	//call service
	code, resp := pocketsProvider.PocketsCalculations(pocketsReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: PocketsCalculations requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// Pockets User
// @Tags space user pockets V1
// @Description Multiple And Individual Stocks Calculations
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.MultipleAndIndividualStocksCalculationsReq true "MultipleAndIndividualStocksCalculations"
// @Success 200 {object} apihelpers.APIRes{data=models.MultipleAndIndividualStocksCalculationsRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/pockets/multipleAndIndividualStocksCalculations [POST]
func MultipleAndIndividualStocksCalculations(c *gin.Context) {
	var pocketsReq models.MultipleAndIndividualStocksCalculationsReq

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("MultipleAndIndividualStocksCalculations (controller), error decoding header, error:", err)
		return
	}

	err := json.NewDecoder(c.Request.Body).Decode(&pocketsReq)
	if err != nil {
		loggerconfig.Error("MultipleAndIndividualStocksCalculations (controller), error decoding body, error:", err, "requestId:", reqH.RequestId, "ClientID: ", reqH.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("MultipleAndIndividualStocksCalculations (controller), reqParams:", helpers.LogStructAsJSON(pocketsReq), "clientID: ", reqH.ClientId, "requestId:", reqH.RequestId)

	//call service
	code, resp := pocketsProvider.MultipleAndIndividualStocksCalculations(pocketsReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: MultipleAndIndividualStocksCalculations requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// StorePocketTransaction
// @Tags space user pockets V1
// @Description StorePocketTransaction - It will store the details of transaction in db
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.StorePocketTransactionReq true "StorePocketTransaction"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/pockets/storePocketTransaction [POST]
func StorePocketTransaction(c *gin.Context) {
	var pocketsReq models.StorePocketTransactionReq

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("StorePocketTransaction (controller), error decoding header, error:", err)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&pocketsReq)
	if err != nil {
		loggerconfig.Error("StorePocketTransaction (controller), error decoding body, error:", err, "clientID: ", pocketsReq.ClientId, "requestId:", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("StorePocketTransaction (controller), reqParams:", helpers.LogStructAsJSON(pocketsReq), "clientID: ", pocketsReq.ClientId, "requestId:", reqH.RequestId)

	//call service
	code, resp := pocketsProvider.StorePocketTransaction(pocketsReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: StorePocketTransaction requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
