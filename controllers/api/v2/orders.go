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

var theOrderProviderV2 models.OrderProvider
var theConditionalOrderProviderV2 models.ConditionalOrderProvider

func InitOrderProviderV2(provider models.OrderProvider) {
	defer models.HandlePanic()
	theOrderProviderV2 = provider
}

func InitConditionalOrderProviderV2(provider models.ConditionalOrderProvider) {
	defer models.HandlePanic()
	theConditionalOrderProviderV2 = provider
}

// PendingOrder
// @Tags space order V2
// @Description Fetch Pending Order V2 - A Buy or Sell order that was placed but yet to be executed is pending order
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param clientId query string true "clientId Query Parameter" dataType(string)
// @Param type query string true "type Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.PendingOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/orderapis/pendingOrder [GET]
func PendingOrder(c *gin.Context) {
	var reqParams models.PendingOrderRequest

	clientID := c.Query("clientId")
	type1 := c.Query("type")
	if clientID == "" || type1 == "" {
		loggerconfig.Error("PendingOrder V2 (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.ClientID = clientID
	reqParams.Type = type1
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("PendingOrder V2 (controller), Empty Device Type clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	if reqParams.ClientID == "" {
		loggerconfig.Error("PendingOrder V2 (controller), Empty Client Id clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidClient)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("PendingOrder V2 (controller), Error validating struct: ", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("PendingOrder V2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("PendingOrder V2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("PendingOrder V2 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
	code, resp := theOrderProviderV2.PendingOrder(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}

// CompletedOrder
// @Tags space order V2
// @Description Fetch Completed Order V2 - It Fetch Completed Orders
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param clientId query string true "clientId Query Parameter" dataType(string)
// @Param type query string true "type Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.CompletedOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/orderapis/completedOrder [GET]
func CompletedOrder(c *gin.Context) {
	var reqParams models.CompletedOrderRequest

	clientID := c.Query("clientId")
	type1 := c.Query("type")
	if clientID == "" || type1 == "" {
		loggerconfig.Error("CompletedOrder V2 (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.ClientID = clientID
	reqParams.Type = type1

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("CompletedOrder V2 (controller), Empty Device Type clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}
	if reqParams.ClientID == "" {
		loggerconfig.Error("CompletedOrder V2 (controller), Empty Client Id clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidClient)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("CompletedOrder V2 (controller), Error validating struct: ", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("CompletedOrder V2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("CompletedOrder V2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("CompletedOrder V2 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
	code, resp := theOrderProviderV2.CompletedOrder(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}

// TradeBook
// @Tags space order V2
// @Description Trade Book V2 - Trade book is a list that reflects executed or completed trades.
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param clientId query string true "clientId Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.TradeBookResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/orderapis/tradeBook [GET]
func TradeBook(c *gin.Context) {
	var reqParams models.TradeBookRequest
	clientID := c.Query("clientId")
	if clientID == "" {
		loggerconfig.Error("TradeBook V2 (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.ClientID = clientID

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("TradeBook V2 (controller), Empty Device Type clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("TradeBook V2 (controller), Error validating struct: ", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("TradeBook V2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("TradeBook V2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("TradeBook V2 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
	code, resp := theOrderProviderV2.TradeBook(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// OrderHistory
// @Tags space order V2
// @Description Order History V2 - It comprises of completed, rejected, cancelled, and failed orders.
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param clientId query string true "clientId Query Parameter" dataType(string)
// @Param omsOrderId query string true "omsOrderId Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.OrderHistoryResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/orderapis/orderHistory [GET]
func OrderHistory(c *gin.Context) {
	var reqParams models.OrderHistoryRequest
	clientID := c.Query("clientId")
	omsOrderId := c.Query("omsOrderId")
	if clientID == "" || omsOrderId == "" {
		loggerconfig.Error("CompletedOrder V2 (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.ClientID = clientID
	reqParams.OmsOrderID = omsOrderId

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("OrderHistory V2 (controller), Empty Device Type clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("OrderHistory V2 (controller), Error validating struct: ", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("OrderHistory V2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("OrderHistory V2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("OrderHistory V2 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
	code, resp := theOrderProviderV2.OrderHistory(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// PlaceOrder
// @Tags space order V2
// @Description Place Order V2 - An order to buy or sell a stock
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.OrderReq true "order"
// @Success 200 {object} apihelpers.APIRes{data=models.PlaceOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/orderapis/placeOrder [POST]
func PlaceOrder(c *gin.Context) {

	var reqParams models.OrderReq

	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("PlaceOrder (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("PlaceOrder (controller), Empty Device Type clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceId)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("PlaceOrder (controller), Error validating struct: ", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("PlaceOrder V2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("PlaceOrder V2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	values, ok := reqParams.RequestPacket.(map[string]interface{})

	if !ok {
		loggerconfig.Error("PlaceOrder (controller), RequestPacket is not a valid map[string]interface{}")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	jsonString, err := json.Marshal(values)
	if err != nil {
		loggerconfig.Error("PlaceOrder (controller), error converting map to JSON:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	switch reqParams.OrderType {
	case "REGULAR":
		placeOrderReq := models.PlaceOrderRequest{}

		err = json.Unmarshal(jsonString, &placeOrderReq)
		if err != nil {
			loggerconfig.Error("PlaceOrder (controller), error decoding body, error:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		err = validate.Struct(placeOrderReq)
		if err != nil {
			loggerconfig.Error("PlaceOrder (controller), Error validating PlaceOrderRequest struct: ", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		loggerconfig.Info("PlaceOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		code, resp := theOrderProviderV2.PlaceOrder(placeOrderReq, requestH)
		logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
		apihelpers.CustomResponse(c, code, resp, logDetail)
		return

	case "BO":
		placeBOOrderReq := models.PlaceBOOrderRequest{}

		err = json.Unmarshal(jsonString, &placeBOOrderReq)
		if err != nil {
			loggerconfig.Error("PlaceOrder (controller), error decoding body, error:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		err = validate.Struct(placeBOOrderReq)
		if err != nil {
			loggerconfig.Error("PlaceOrder (controller), Error validating PlaceOrderRequest struct: ", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		loggerconfig.Info("PlaceOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		code, resp := theConditionalOrderProviderV2.PlaceBOOrder(placeBOOrderReq, requestH)
		logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
		apihelpers.CustomResponse(c, code, resp, logDetail)
		return

	case "AMO":
		placeAMOOrderReq := models.PlaceOrderRequest{}
		err = json.Unmarshal(jsonString, &placeAMOOrderReq)
		if err != nil {
			loggerconfig.Error("PlaceOrder (controller), error decoding body, error:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		err = validate.Struct(placeAMOOrderReq)
		if err != nil {
			loggerconfig.Error("PlaceOrder (controller), Error validating PlaceOrderRequest struct: ", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		loggerconfig.Info("PlaceOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		code, resp := theOrderProviderV2.PlaceAMOOrder(placeAMOOrderReq, requestH)
		logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
		apihelpers.CustomResponse(c, code, resp, logDetail)
		return

	case "CO":
		placeCOOrderReq := models.PlaceCOOrderRequest{}
		err = json.Unmarshal(jsonString, &placeCOOrderReq)
		if err != nil {
			loggerconfig.Error("PlaceOrder (controller), error decoding body, error:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		err = validate.Struct(placeCOOrderReq)
		if err != nil {
			loggerconfig.Error("PlaceOrder (controller), Error validating PlaceOrderRequest struct: ", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		loggerconfig.Info("PlaceOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		code, resp := theConditionalOrderProviderV2.PlaceCOOrder(placeCOOrderReq, requestH)
		logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
		apihelpers.CustomResponse(c, code, resp, logDetail)
		return

	case "Spread":
		placeSpreadOrderReq := models.PlaceSpreadOrderRequest{}
		err = json.Unmarshal(jsonString, &placeSpreadOrderReq)
		if err != nil {
			loggerconfig.Error("PlaceOrder (controller), error decoding body, error:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		err = validate.Struct(placeSpreadOrderReq)
		if err != nil {
			loggerconfig.Error("PlaceOrder (controller), Error validating PlaceOrderRequest struct: ", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		loggerconfig.Info("PlaceOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		code, resp := theConditionalOrderProviderV2.PlaceSpreadOrder(placeSpreadOrderReq, requestH)
		logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
		apihelpers.CustomResponse(c, code, resp, logDetail)
		return

	case "GTT":
		placeGttOrderReq := models.CreateGTTOrderRequest{}
		err = json.Unmarshal(jsonString, &placeGttOrderReq)
		if err != nil {
			loggerconfig.Error("PlaceOrder (controller), error decoding body, error:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		err = validate.Struct(placeGttOrderReq)
		if err != nil {
			loggerconfig.Error("PlaceOrder (controller), Error validating PlaceOrderRequest struct: ", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		loggerconfig.Info("PlaceOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		code, resp := theOrderProviderV2.PlaceGTTOrder(placeGttOrderReq, requestH)
		logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
		apihelpers.CustomResponse(c, code, resp, logDetail)
		return

	default:
		loggerconfig.Error("CancelOrder (controller), Error Invalid order type clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
}

// modifyOrder
// @Tags space order V2
// @Description Modify Order - It will modify the already placed order
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.OrderReq true "order"
// @Success 200 {object} apihelpers.APIRes{data=models.ModifyOrCancelOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/orderapis/modifyOrder [POST]
func ModifyOrder(c *gin.Context) {

	var reqParams models.OrderReq

	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ModifyOrder (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("ModifyOrder (controller), Empty Device Type clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceId)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("ModifyOrder (controller), Error validating struct: ", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("ModifyOrder V2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("ModifyOrder V2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	values, ok := reqParams.RequestPacket.(map[string]interface{})

	if !ok {
		loggerconfig.Error("PlaceOrder (controller), RequestPacket is not a valid map[string]interface{}")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	jsonString, err := json.Marshal(values)
	if err != nil {
		loggerconfig.Error("PlaceOrder (controller), error converting map to JSON:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	switch reqParams.OrderType {
	case "REGULAR":
		modifyOrderReq := models.ModifyOrderRequest{}

		err = json.Unmarshal(jsonString, &modifyOrderReq)
		if err != nil {
			loggerconfig.Error("PlaceOrder (controller), error decoding body, error:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		err = validate.Struct(modifyOrderReq)
		if err != nil {
			loggerconfig.Error("PlaceOrder (controller), Error validating PlaceOrderRequest struct: ", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		loggerconfig.Info("ModifyOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		code, resp := theOrderProviderV2.ModifyOrder(modifyOrderReq, requestH)
		logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
		apihelpers.CustomResponse(c, code, resp, logDetail)

	case "AMO":
		modifyAMOReq := models.ModifyAMORequest{}

		err = json.Unmarshal(jsonString, &modifyAMOReq)
		if err != nil {
			loggerconfig.Error("ModifyOrder (controller), error decoding body:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		err = validate.Struct(modifyAMOReq)
		if err != nil {
			loggerconfig.Error("ModifyOrder (controller), Error validating ModifyAMORequest struct:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		loggerconfig.Info("ModifyOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		code, resp := theOrderProviderV2.ModifyAMOOrder(modifyAMOReq, requestH)
		logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
		apihelpers.CustomResponse(c, code, resp, logDetail)

	case "BO":
		modifyBOReq := models.ModifyBOOrderRequest{}

		err = json.Unmarshal(jsonString, &modifyBOReq)
		if err != nil {
			loggerconfig.Error("ModifyOrder (controller), error decoding body:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		err = validate.Struct(modifyBOReq)
		if err != nil {
			loggerconfig.Error("ModifyOrder (controller), Error validating ModifyBOOrderRequest struct:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		loggerconfig.Info("ModifyOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		code, resp := theConditionalOrderProviderV2.ModifyBOOrder(modifyBOReq, requestH)
		logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
		apihelpers.CustomResponse(c, code, resp, logDetail)

	case "CO":
		modifyCOReq := models.ModifyCOOrderRequest{}

		err = json.Unmarshal(jsonString, &modifyCOReq)
		if err != nil {
			loggerconfig.Error("ModifyOrder (controller), error decoding body:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		err = validate.Struct(modifyCOReq)
		if err != nil {
			loggerconfig.Error("ModifyOrder (controller), Error validating ModifyCOOrderRequest struct:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		loggerconfig.Info("ModifyOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		code, resp := theConditionalOrderProviderV2.ModifyCOOrder(modifyCOReq, requestH)
		logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
		apihelpers.CustomResponse(c, code, resp, logDetail)

	case "Spread":
		modifySpreadReq := models.ModifySpreadOrderRequest{}

		err = json.Unmarshal(jsonString, &modifySpreadReq)
		if err != nil {
			loggerconfig.Error("ModifyOrder (controller), error decoding body:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		err = validate.Struct(modifySpreadReq)
		if err != nil {
			loggerconfig.Error("ModifyOrder (controller), Error validating ModifySpreadOrderRequest struct:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		loggerconfig.Info("ModifyOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		code, resp := theConditionalOrderProviderV2.ModifySpreadOrder(modifySpreadReq, requestH)
		logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
		apihelpers.CustomResponse(c, code, resp, logDetail)

	case "GTT":
		modifyGTTReq := models.ModifyGTTOrderRequest{}

		err = json.Unmarshal(jsonString, &modifyGTTReq)
		if err != nil {
			loggerconfig.Error("ModifyOrder (controller), error decoding body:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		err = validate.Struct(modifyGTTReq)
		if err != nil {
			loggerconfig.Error("ModifyOrder (controller), Error validating ModifyGTTOrderRequest struct:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		loggerconfig.Info("ModifyOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		code, resp := theOrderProviderV2.ModifyGTTOrder(modifyGTTReq, requestH)
		logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
		apihelpers.CustomResponse(c, code, resp, logDetail)

	default:
		loggerconfig.Error("ModifyOrder (controller), Error Invalid order type clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

}

// CancelOrder
// @Tags space order V2
// @Description Cancel Order V2 - Cancel the already placed order
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.CancelOrderRequest true "order"
// @Success 200 {object} apihelpers.APIRes{data=models.ModifyOrCancelOrderResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/orderapis/CancelOrder [POST]
func CancelOrder(c *gin.Context) {
	var reqParams models.OrderReq

	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("CancelOrder (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("CancelOrder (controller), Empty Device Id clientID: ", reqParams.ClientID, "requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceId)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("CancelOrder (controller), Error validating struct: ", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("CancelOrder V2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("CancelOrder V2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	values, ok := reqParams.RequestPacket.(map[string]interface{})

	if !ok {
		loggerconfig.Error("PlaceOrder (controller), RequestPacket is not a valid map[string]interface{}")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	jsonString, err := json.Marshal(values)
	if err != nil {
		loggerconfig.Error("PlaceOrder (controller), error converting map to JSON:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	switch reqParams.OrderType {
	case "REGULAR":
		cancelOrderReq := models.CancelOrderRequest{}

		err = json.Unmarshal(jsonString, &cancelOrderReq)
		if err != nil {
			loggerconfig.Error("PlaceOrder (controller), error decoding body, error:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		err = validate.Struct(cancelOrderReq)
		if err != nil {
			loggerconfig.Error("PlaceOrder (controller), Error validating PlaceOrderRequest struct: ", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}
		loggerconfig.Info("CancelOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		code, resp := theOrderProviderV2.CancelOrder(cancelOrderReq, requestH)
		logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
		apihelpers.CustomResponse(c, code, resp, logDetail)
		return

	case "AMO":
		cancelAMOReq := models.CancelOrderRequest{}

		err = json.Unmarshal(jsonString, &cancelAMOReq)
		if err != nil {
			loggerconfig.Error("CancelOrder (controller), error decoding body, error:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		err = validate.Struct(cancelAMOReq)
		if err != nil {
			loggerconfig.Error("CancelOrder (controller), Error validating CancelAMORequest struct:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		loggerconfig.Info("CancelOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		code, resp := theOrderProviderV2.CancelOrder(cancelAMOReq, requestH)
		logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
		apihelpers.CustomResponse(c, code, resp, logDetail)
		return

	case "BO":
		exitBOReq := models.ExitBOOrderRequest{}

		err = json.Unmarshal(jsonString, &exitBOReq)
		if err != nil {
			loggerconfig.Error("CancelOrder (controller), error decoding body, error:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		err = validate.Struct(exitBOReq)
		if err != nil {
			loggerconfig.Error("CancelOrder (controller), Error validating ExitBOOrderRequest struct:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		loggerconfig.Info("CancelOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		code, resp := theConditionalOrderProviderV2.CancelBOOrder(exitBOReq, requestH)
		logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
		apihelpers.CustomResponse(c, code, resp, logDetail)
		return

	case "CO":
		exitCOReq := models.ExitCOOrderRequest{}

		err = json.Unmarshal(jsonString, &exitCOReq)
		if err != nil {
			loggerconfig.Error("CancelOrder (controller), error decoding body, error:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		err = validate.Struct(exitCOReq)
		if err != nil {
			loggerconfig.Error("CancelOrder (controller), Error validating ExitCOOrderRequest struct:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		loggerconfig.Info("CancelOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		code, resp := theConditionalOrderProviderV2.CancelCOOrder(exitCOReq, requestH)
		logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
		apihelpers.CustomResponse(c, code, resp, logDetail)
		return

	case "Spread":
		exitSpreadReq := models.ExitSpreadOrderRequest{}

		err = json.Unmarshal(jsonString, &exitSpreadReq)
		if err != nil {
			loggerconfig.Error("CancelOrder (controller), error decoding body, error:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		err = validate.Struct(exitSpreadReq)
		if err != nil {
			loggerconfig.Error("CancelOrder (controller), Error validating ExitSpreadOrderRequest struct:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		loggerconfig.Info("CancelOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		code, resp := theConditionalOrderProviderV2.CancelSpreadOrder(exitSpreadReq, requestH)
		logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
		apihelpers.CustomResponse(c, code, resp, logDetail)
		return

	case "GTT":
		cancelGTTReq := models.CancelGTTOrderRequest{}

		err = json.Unmarshal(jsonString, &cancelGTTReq)
		if err != nil {
			loggerconfig.Error("CancelOrder (controller), error decoding body, error:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		err = validate.Struct(cancelGTTReq)
		if err != nil {
			loggerconfig.Error("CancelOrder (controller), Error validating CancelGTTOrderRequest struct:", err, "clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
			apihelpers.ErrorMessage(c, constants.InvalidRequest)
			return
		}

		loggerconfig.Info("CancelOrder (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		code, resp := theOrderProviderV2.CancelGTTOrder(cancelGTTReq, requestH)
		logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
		apihelpers.CustomResponse(c, code, resp, logDetail)
		return

	default:
		loggerconfig.Error("CancelOrder (controller), Error Invalid order type clientID: ", reqParams.ClientID, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
}
