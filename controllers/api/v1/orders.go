package v1

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

var theOrderProvider models.OrderProvider

func InitOrderProvider(provider models.OrderProvider) {
	defer models.HandlePanic()
	theOrderProvider = provider
}

// PlaceOrder
// @Tags space order V1
// @Description PlaceOrder - An order to buy or sell a stock
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.PlaceOrderRequest true "order"
// @Success 200 {object} apihelpers.APIRes{data=models.PlaceOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/orderapis/placeOrder [POST]
func PlaceOrder(c *gin.Context) {
	var reqParams models.PlaceOrderRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("PlaceOrder (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("PlaceOrder (controller), Empty Device Type clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("PlaceOrder (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("PlaceOrder (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("PlaceOrder CheckAuthWithClient difference in authtoken-clientId and clientId", err, " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
	}

	loggerconfig.Info("PlaceOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theOrderProvider.PlaceOrder(reqParams, requestH)

	logDetail := "clientId: " + reqParams.ClientID + " function: PlaceOrder requestId:" + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ModifyOrder
// @Tags space order V1
// @Description ModifyOrder - It will modify the already placed order
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.ModifyOrderRequest true "order"
// @Success 200 {object} apihelpers.APIRes{data=models.ModifyOrCancelOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/orderapis/modifyOrder [POST]
func ModifyOrder(c *gin.Context) {
	var reqParams models.ModifyOrderRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ModifyOrder (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("ModifyOrder (controller), Empty Device Type clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("ModifyOrder (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("ModifyOrder (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("ModifyOrder CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("ModifyOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theOrderProvider.ModifyOrder(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: ModifyOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// CancelOrder
// @Tags space order V1
// @Description CancelOrder - Cancel the already placed order
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.CancelOrderRequest true "order"
// @Success 200 {object} apihelpers.APIRes{data=models.ModifyOrCancelOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/orderapis/cancelOrder [POST]
func CancelOrder(c *gin.Context) {
	var reqParams models.CancelOrderRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("CancelOrder (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("CancelOrder (controller), Empty Device Type clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("CancelOrder (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("CancelOrder (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("CancelOrder CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("CancelOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theOrderProvider.CancelOrder(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: CancelOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// PlaceAMOOrder
// @Tags space order V1
// @Description Place AMO Order - It is After Market Order (AMO) that can be used to place orders outside of regular trading hours
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.PlaceOrderRequest true "AMO order"
// @Success 200 {object} apihelpers.APIRes{data=models.AMOOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/orderapis/placeAMOOrder [POST]
func PlaceAMOOrder(c *gin.Context) {
	var reqParams models.PlaceOrderRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("PlaceAMOOrder (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("PlaceAMOOrder (controller), Empty Device Type clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("PlaceAMOOrder (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("PlaceAMOOrder (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("PlaceAMOOrder CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("PlaceAMOOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theOrderProvider.PlaceAMOOrder(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: PlaceAMOOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ModifyAMOOrder
// @Tags space order V1
// @Description ModifyAMOOrder - Modify the already placed AMO Order
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.ModifyOrderRequest true "amo order"
// @Success 200 {object} apihelpers.APIRes{data=models.AMOOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/orderapis/modifyAMOOrder [POST]
func ModifyAMOOrder(c *gin.Context) {
	var reqParams models.ModifyAMORequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ModifyAMOOrder (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("ModifyAMOOrder (controller), Empty Device Type clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("ModifyAMOOrder (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("ModifyAMOOrder (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("ModifyAMOOrder CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("ModifyAMOOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theOrderProvider.ModifyAMOOrder(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: ModifyAMOOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// CancelAMOOrder
// @Tags space order V1
// @Description CancelAMOOrder - Cancel already placed AMO order
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.CancelOrderRequest true "amo order"
// @Success 200 {object} apihelpers.APIRes{data=models.AMOOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/orderapis/cancelAMOOrder [POST]
func CancelAMOOrder(c *gin.Context) {
	var reqParams models.CancelOrderRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("CancelAMOOrder (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("CancelAMOOrder (controller), Empty Device Type clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("CancelAMOOrder (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("CancelAMOOrder (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("CancelAMOOrder CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("CancelAMOOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theOrderProvider.CancelOrder(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: CancelAMOOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// PendingOrder
// @Tags space order V1
// @Description Fetch Pending Order - A Buy or Sell order that was placed but yet to be executed is pending order
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.PendingOrderRequest true "order"
// @Success 200 {object} apihelpers.APIRes{data=models.PendingOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/orderapis/pendingOrder [POST]
func PendingOrder(c *gin.Context) {
	var reqParams models.PendingOrderRequest

	errr := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if errr != nil {
		loggerconfig.Error("PendingOrder (controller), error decoding body, error:", errr)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("PendingOrder (controller), Empty Device Type clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	if reqParams.ClientID == "" {
		loggerconfig.Error("PendingOrder (controller), Empty Client Id clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidClient)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("PendingOrder (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("PendingOrder (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("PendingOrder CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("PendingOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theOrderProvider.PendingOrder(reqParams, requestH)

	logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}

// CompletedOrder
// @Tags space order V1
// @Description CompletedOrder - It Fetch Completed Orders
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.CompletedOrderRequest true "order"
// @Success 200 {object} apihelpers.APIRes{data=models.CompletedOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/orderapis/completedOrder [POST]
func CompletedOrder(c *gin.Context) {
	var reqParams models.CompletedOrderRequest

	errr := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if errr != nil {
		loggerconfig.Error("CompletedOrder (controller), error decoding body, error:", errr)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("CompletedOrder (controller), Empty Device Type clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}
	if reqParams.ClientID == "" {
		loggerconfig.Error("CompletedOrder (controller), Empty Client Id clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidClient)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("CompletedOrder (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("CompletedOrder (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("CompletedOrder CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("CompletedOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theOrderProvider.CompletedOrder(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: CompletedOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}

// TradeBook
// @Tags space order V1
// @Description TradeBook - Trade book is a list that reflects executed or completed trades.
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.TradeBookRequest true "trade book"
// @Success 200 {object} apihelpers.APIRes{data=models.TradeBookResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/orderapis/tradeBook [POST]
func TradeBook(c *gin.Context) {
	var reqParams models.TradeBookRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("TradeBook (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("TradeBook (controller), Empty Device Type clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("TradeBook (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("TradeBook (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("TradeBook CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("TradeBook (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theOrderProvider.TradeBook(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: TradeBook requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// OrderHistory
// @Tags space order V1
// @Description OrderHistory - It comprises of completed, rejected, cancelled, and failed orders.
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.OrderHistoryRequest true "order history"
// @Success 200 {object} apihelpers.APIRes{data=models.OrderHistoryResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/orderapis/orderHistory [POST]
func OrderHistory(c *gin.Context) {
	var reqParams models.OrderHistoryRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("OrderHistory (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("OrderHistory (controller), Empty Device Type clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("OrderHistory (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("OrderHistory (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("OrderHistory CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("OrderHistory (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theOrderProvider.OrderHistory(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: OrderHistory requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// CreateGTTOrder
// @Tags space gtt V1
// @Description CreateGTTOrder - It will create GTT order which provides feature that allows investors to buy and sell as per their predetermined price
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.CreateGTTOrderRequest true "gtt order"
// @Success 200 {object} apihelpers.APIRes{data=models.GTTOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/orderapis/createGTTOrder [POST]
func CreateGTTOrder(c *gin.Context) {
	var reqParams models.CreateGTTOrderRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("CreateGTTOrder (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("CreateGTTOrder (controller), Empty Device Type clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("CreateGTTOrder (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("CreateGTTOrder (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("CreateGTTOrder CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("CreateGTTOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theOrderProvider.PlaceGTTOrder(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: CreateGTTOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ModifyGTTOrder
// @Tags space gtt V1
// @Description ModifyGTTOrder - Modify the placed GTT Order
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.ModifyGTTOrderRequest true "gtt order"
// @Success 200 {object} apihelpers.APIRes{data=models.GTTOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/orderapis/modifyGTTOrder [POST]
func ModifyGTTOrder(c *gin.Context) {
	var reqParams models.ModifyGTTOrderRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ModifyGTTOrder (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("ModifyGTTOrder (controller), Empty Device Type clientID: ", reqParams.Order.ClientID, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("ModifyGTTOrder (controller), Error validating struct: ", err, "clientID: ", reqParams.Order.ClientID, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.Order.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("ModifyGTTOrder (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.Order.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("ModifyGTTOrder CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.Order.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("ModifyGTTOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "clientID: ", reqParams.Order.ClientID, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theOrderProvider.ModifyGTTOrder(reqParams, requestH)
	logDetail := "clientId: " + reqParams.Order.ClientID + " function: ModifyGTTOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// CancelGTTOrder
// @Tags space gtt V1
// @Description CancelGTTOrder - Cancel the Placed GTT Order
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.CancelGTTOrderRequest true "gtt order"
// @Success 200 {object} apihelpers.APIRes{data=models.ModifyGTTOrderRequest}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/orderapis/cancelGTTOrder [POST]
func CancelGTTOrder(c *gin.Context) {
	var reqParams models.CancelGTTOrderRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("CancelGTTOrder (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("CancelGTTOrder (controller), Empty Device Type clientID: ", reqParams.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("CancelGTTOrder (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("CancelGTTOrder (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("CancelGTTOrder CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("CancelGTTOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theOrderProvider.CancelGTTOrder(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientId + " function: CancelGTTOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchGTTOrder
// @Tags space gtt V1
// @Description FetchGTTOrder - Provide list of GTT Order placed
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.FetchGTTOrderRequest true "gtt order"
// @Success 200 {object} apihelpers.APIRes{data=models.FetchGTTOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/orderapis/fetchGTTOrder [POST]
func FetchGTTOrder(c *gin.Context) {
	var reqParams models.FetchGTTOrderRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FetchGTTOrder (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchGTTOrder (controller), Empty Device Type clientID: ", reqParams.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FetchGTTOrder (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("FetchGTTOrder (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("FetchGTTOrder CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("FetchGTTOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theOrderProvider.FetchGTTOrder(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientId + " function: FetchGTTOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// PlaceGTTOcoOrder
// @Tags space gtt V1
// @Description PlaceGttOCOOrder - It will create GTT OCO order
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.CreateGttOCORequest true "gtt oco order"
// @Success 200 {object} apihelpers.APIRes{data=models.GTTOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/orderapis/PlaceGTTOcoOrder [POST]
func PlaceGttOCOOrder(c *gin.Context) {
	var reqParams models.CreateGttOCORequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("PlaceGttOCOOrder (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("PlaceGttOCOOrder (controller), Empty Device Type clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("PlaceGttOCOOrder (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("PlaceGttOCOOrder (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("PlaceGttOCOOrder CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("PlaceGttOCOOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theOrderProvider.PlaceGttOCOOrder(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: CreateGTTOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// MarginCalculations
// @Tags Margin Calculations V1
// @Description MarginCalculations - It provides the total amount of funds that can be utilized
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.MarginCalculationRequest true "Margin Calculations"
// @Success 200 {object} apihelpers.APIRes{data=models.MarginResultData}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/orderapis/marginCalculations [POST]
func MarginCalculations(c *gin.Context) {
	var reqParams models.MarginCalculationRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("MarginCalculations (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("MarginCalculations (controller), Empty Device Type clientID: ", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	for i := 0; i < len(reqParams.Data); i++ {

		// Validate Exchange
		exchange := reqParams.Data[i].Exchange
		if !constants.ValidExchangeMap[strings.ToUpper(exchange)] {
			logrus.Error("MarginCalculations (controller), invalid exchange provided: ", exchange, " clientId:", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
			apihelpers.SendErrorController(c, false, constants.InvalidRequest, http.StatusBadRequest)
			return
		}

		price, err := strconv.ParseFloat(reqParams.Data[i].Price, 64)
		if err != nil {
			logrus.Error("MarginCalculations (controller), can't convert price to float err: ", err, " clientId:", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
			apihelpers.SendErrorController(c, false, constants.InvalidRequest, http.StatusBadRequest)
			return
		}
		if price <= 0 { // in margin calcualtion we will always get some value in price so <= to 0 is used
			logrus.Error("MarginCalculations (controller), invalid value of price provided: ", price, " clientId:", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
			apihelpers.SendErrorController(c, false, constants.InvalidRequest, http.StatusBadRequest)
			return
		}

		quantity, err := strconv.Atoi(reqParams.Data[i].Quantity)
		if err != nil {
			logrus.Error("MarginCalculations (controller), can't convert quantity to int err: ", err, " clientId:", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
			apihelpers.SendErrorController(c, false, constants.InvalidRequest, http.StatusBadRequest)
			return
		}
		if quantity <= 0 {
			logrus.Error("MarginCalculations (controller), invalid value of quantity provided: ", quantity, " clientId:", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
			apihelpers.SendErrorController(c, false, constants.InvalidRequest, http.StatusBadRequest)
			return
		}

		if strings.ToUpper(reqParams.Data[i].Mode) != constants.New {
			logrus.Error("MarginCalculations (controller), invalid mode provided: ", reqParams.Data[i].Mode, " clientId:", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
			apihelpers.SendErrorController(c, false, constants.InvalidRequest, http.StatusBadRequest)
			return
		}
	}

	loggerconfig.Info("MarginCalculations (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theOrderProvider.MarginCalculations(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: MarginCalculations requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// LastTradedPrice
// @Tags space order V1
// @Description LastTradedPrice - It will provide last price at which traded occured
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.LastTradedPriceRequest true "Last Traded Price"
// @Success 200 {object} apihelpers.APIRes{data=models.LastTradedPriceResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/orderapis/lastTradedPrice [POST]
func LastTradedPrice(c *gin.Context) {

	var reqParams models.LastTradedPriceRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("LastTradedPrice (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("LastTradedPrice (controller), Empty Device Type clientID: ", requestH.ClientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}
	loggerconfig.Info("LastTradedPrice (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theOrderProvider.LastTradedPrice(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: LastTradedPrice requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// PlaceIcebergOrder
// @Tags space order V1
// @Description PlaceIcebergOrder - It will place iceberg order
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.IcebergOrderReq true "Place Iceberg Order"
// @Success 200 {object} apihelpers.APIRes{data=models.IcebergOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/orderapis/createIcebergOrder [POST]
func CreateIcebergOrder(c *gin.Context) {

	var reqParams models.IcebergOrderReq
	err := c.ShouldBindJSON(&reqParams)
	if err != nil {
		loggerconfig.Error("PlaceIcebergOrder (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("LastTradedPrice (controller), Empty Device Type clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("PlaceIcebergOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "ClientID: ", reqParams.ClientID, " deviceId: ", requestH.DeviceId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theOrderProvider.PlaceIcebergOrder(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: PlaceIcebergOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ModifyIcebergOrder
// @Tags space order V1
// @Description ModifyIcebergOrder - It will modify an existing iceberg order
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion header string false "P-ClientVersion Header"
// @Param request body models.ModifyIcebergOrderReq true "Modify Iceberg Order"
// @Success 200 {object} apihelpers.APIRes{data=models.IcebergOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/orderapis/modifyIcebergOrder [PUT]
func ModifyIcebergOrder(c *gin.Context) {
	var reqParams models.ModifyIcebergOrderReq
	err := c.ShouldBindJSON(&reqParams)
	if err != nil {
		loggerconfig.Error("ModifyIcebergOrder (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	if requestH.DeviceType == "" {
		loggerconfig.Error("ModifyIcebergOrder (controller), Empty Device Type clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}
	loggerconfig.Info("ModifyIcebergOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "ClientID: ", reqParams.ClientID, " deviceId: ", requestH.DeviceId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
	code, resp := theOrderProvider.ModifyIcebergOrder(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: ModifyIcebergOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// CancelIcebergOrder
// @Tags space order V1
// @Description CancelIcebergOrder - It will cancel an existing iceberg order
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion header string false "P-ClientVersion Header"
// @Param omsOrderId query string true "Order ID"
// @Param clientId query string true "Client ID"
// @Param executiontype query string true "Execution Type"
// @Success 200 {object} apihelpers.APIRes{data=models.IcebergCanelOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/orderapis/cancelIcebergOrder [DELETE]
func CancelIcebergOrder(c *gin.Context) {
	orderId := c.Query("omsOrderId")
	clientId := c.Query("clientId")
	executionType := c.Query("executiontype")

	if orderId == "" || clientId == "" || executionType == "" {
		loggerconfig.Error("CancelIcebergOrder (controller), missing required parameters. OrderID: ", orderId, "ClientID: ", clientId, "ExecutionType: ", executionType)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	if requestH.DeviceType == "" {
		loggerconfig.Error("CancelIcebergOrder (controller), Empty Device Type clientID: ", clientId, "requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("CancelIcebergOrder (controller), orderId:", orderId, "clientId:", clientId, "executionType:", executionType, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, "requestId:", requestH.RequestId, " clientVersion:", requestH.ClientVersion)

	cancelReq := models.CancelIcebergOrderReq{
		OmsOrderID:    orderId,
		ClientId:      clientId,
		ExecutionType: executionType,
	}

	code, resp := theOrderProvider.CancelIcebergOrder(cancelReq, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: CancelIcebergOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
