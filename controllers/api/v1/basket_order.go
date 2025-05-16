package v1

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

var theBasketOrderProvider models.BasketOrderProvider

func InitBasketOrderProvider(provider models.BasketOrderProvider) {
	defer models.HandlePanic()
	theBasketOrderProvider = provider
}

// CreateBasket
// @Tags space basket order V1
// @Description Create Basket
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.CreateBasketReq true "Create Basket"
// @Success 200 {object} apihelpers.APIRes{data=models.BasketRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/basket/createBasket [POST]
func CreateBasket(c *gin.Context) {
	var reqParams models.CreateBasketReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("CreateBasket (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("CreateBasket (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("CreateBasket (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.SendErrorController(c, false, constants.InvalidNameLength, http.StatusBadRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.LoginID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("CreateBasket (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.LoginID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("CreateBasket (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.LoginID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("CreateBasket (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theBasketOrderProvider.CreateBasket(reqParams, requestH)
	logDetail := "clientId: " + reqParams.LoginID + " function: CreateBasket requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchBasket
// @Tags space basket order V1
// @Description Fetch Basket
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.FetchBasketReq true "Fetch Basket"
// @Success 200 {object} apihelpers.APIRes{data=models.BasketRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/basket/fetchBasket [POST]
func FetchBasket(c *gin.Context) {
	var reqParams models.FetchBasketReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FetchBasket (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchBasket (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FetchBasket (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.LoginID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("FetchBasket (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.LoginID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("FetchBasket (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.LoginID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("FetchBasket (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theBasketOrderProvider.FetchBasket(reqParams, requestH)
	logDetail := "clientId: " + reqParams.LoginID + " function: FetchBasket requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DeleteBasket
// @Tags space basket order V1
// @Description Delete Basket
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.DeleteBasketReq true "Delete Basket"
// @Success 200 {object} apihelpers.APIRes{data=models.BasketRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/basket/deleteBasket [POST]
func DeleteBasket(c *gin.Context) {
	var reqParams models.DeleteBasketReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("DeleteBasket (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("DeleteBasket (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("DeleteBasket (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("DeleteBasket (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theBasketOrderProvider.DeleteBasket(reqParams, requestH)
	logDetail := "basketId: " + reqParams.BasketID + " function: DeleteBasket requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// AddBasketInstrument
// @Tags space basket order V1
// @Description Add Basket Instrument
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.AddBasketInstrumentReq true "Add Basket Instrument"
// @Success 200 {object} apihelpers.APIRes{data=models.BasketInstrumentRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/basket/addBasketInstrument [POST]
func AddBasketInstrument(c *gin.Context) {
	var reqParams models.AddBasketInstrumentReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("AddBasketInstrument (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("AddBasketInstrument (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", reqParams.OrderInfo.ClientID, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("AddBasketInstrument (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", reqParams.OrderInfo.ClientID, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.OrderInfo.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("AddBasketInstrument (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.OrderInfo.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("AddBasketInstrument (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.OrderInfo.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("AddBasketInstrument (controller),reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.OrderInfo.ClientID, "requestId:", requestH.RequestId, " deviceId: ", requestH.DeviceId)
	code, resp := theBasketOrderProvider.AddBasketInstrument(reqParams, requestH)
	logDetail := "clientId: " + reqParams.BasketID + " function: AddBasketInstrument requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// EditBasketInstrument
// @Tags space basket order V1
// @Description Edit Basket Instrument
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.EditBasketInstrumentReq true "Edit Basket Instrument"
// @Success 200 {object} apihelpers.APIRes{data=models.BasketInstrumentRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/basket/editBasketInstrument [POST]
func EditBasketInstrument(c *gin.Context) {
	var reqParams models.EditBasketInstrumentReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("EditBasketInstrument (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("EditBasketInstrument (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", reqParams.OrderInfo.ClientID, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("EditBasketInstrument (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", reqParams.OrderInfo.ClientID, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.OrderInfo.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("EditBasketInstrument (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.OrderInfo.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("EditBasketInstrument (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.OrderInfo.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("EditBasketInstrument (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.OrderInfo.ClientID, "requestId:", requestH.RequestId, " deviceId: ", requestH.DeviceId)
	code, resp := theBasketOrderProvider.EditBasketInstrument(reqParams, requestH)
	logDetail := "basketId: " + reqParams.BasketID + " function: EditBasketInstrument requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DeleteBasketInstrument
// @Tags space basket order V1
// @Description Delete Basket Instrument
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.DeleteBasketInstrumentReq true "Delete Basket Instrument"
// @Success 200 {object} apihelpers.APIRes{data=models.BasketRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/basket/deleteBasketInstrument [POST]
func DeleteBasketInstrument(c *gin.Context) {
	var reqParams models.DeleteBasketInstrumentReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("DeleteBasketInstrument (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("DeleteBasketInstrument (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("DeleteBasketInstrument (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("DeleteBasketInstrument (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theBasketOrderProvider.DeleteBasketInstrument(reqParams, requestH)
	logDetail := "clientId: " + reqParams.BasketID + " function: DeleteBasketInstrument requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// RenameBasketInstrument
// @Tags space basket order V1
// @Description Rename Basket Instrument
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.RenameBasketReq true "Rename Basket Instrument"
// @Success 200 {object} apihelpers.APIRes{data=models.BasketInstrumentRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/basket/renameBasket [POST]
func RenameBasket(c *gin.Context) {
	var reqParams models.RenameBasketReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("RenameBasket (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("RenameBasket (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("RenameBasket (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("RenameBasket (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theBasketOrderProvider.RenameBasket(reqParams, requestH)
	logDetail := "basketID: " + reqParams.BasketID + " function: RenameBasket requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ExecuteBasket
// @Tags space basket order V1
// @Description Execute Basket
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.ExecuteBasketReq true "Execute Basket"
// @Success 200 {object} apihelpers.APIRes{data=models.ExecuteBasketRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/basket/executeBasket [POST]
func ExecuteBasket(c *gin.Context) {
	var reqParams models.ExecuteBasketReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ExecuteBasket (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("ExecuteBasket (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("ExecuteBasket (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("ExecuteBasket (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.BasketID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("ExecuteBasket (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("ExecuteBasket (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " ClientID: ", reqParams.ClientID, "requestId:", requestH.RequestId, " deviceId: ", requestH.DeviceId)
	code, resp := theBasketOrderProvider.ExecuteBasket(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: ExecuteBasket requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// UpdateBasketExecutionState
// @Tags space basket order V1
// @Description Update Basket Execution State
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.UpdateBasketExecutionStateReq true "Update Basket Execution State"
// @Success 200 {object} apihelpers.APIRes{data=models.BasketRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/basket/updateBasketExecutionState [POST]
func UpdateBasketExecutionState(c *gin.Context) {
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	var reqParams models.UpdateBasketExecutionStateReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("UpdateBasketExecutionState (controller), requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("UpdateBasketExecutionState (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("UpdateBasketExecutionState (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("UpdateBasketExecutionState (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("UpdateBasketExecutionState (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("UpdateBasketExecutionState (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " ClientID: ", reqParams.ClientID, "requestId:", requestH.RequestId, " deviceId: ", requestH.DeviceId)
	code, resp := theBasketOrderProvider.UpdateBasketExecutionState(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: UpdateBasketExecutionState requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
