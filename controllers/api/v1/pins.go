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
)

var pinProvider models.PinsProvider

func InitPinsProvider(provider models.PinsProvider) {
	defer models.HandlePanic()
	pinProvider = provider
}

// FetchPins
// @Tags space pins V1
// @Description Fetch Pins
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.PinsRequest true "funds"
// @Success 200 {object} apihelpers.APIRes{data=models.PinsResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/pins/fetchPins [POST]
func FetchPins(c *gin.Context) {
	var fetchPinsReq models.PinsRequest

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("FetchPins (controller), error decoding header, error:", err)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&fetchPinsReq)
	if err != nil {
		loggerconfig.Error("FetchPins (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(fetchPinsReq.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("FetchPins (controller) CheckAuthWithClient invalid authtoken", " clientId: ", fetchPinsReq.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("FetchPins (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", fetchPinsReq.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("FetchPins (controller), reqParams:", helpers.LogStructAsJSON(fetchPinsReq), " uccId: ", fetchPinsReq.ClientId, "requestId:", reqH.RequestId)
	//call service
	code, resp := pinProvider.FetchPins(fetchPinsReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: FetchPins requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// UpdatePins
// @Tags space pins V1
// @Description Update Pins
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.UpdatePins true "pins"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/pins/updatePins [POST]
func UpdatePins(c *gin.Context) {
	var modifyPinsReq models.UpdatePins

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("UpdatePins (controller), error decoding header, error:", err)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&modifyPinsReq)
	if err != nil {
		loggerconfig.Error("UpdatePins (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(modifyPinsReq.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("UpdatePins (controller) CheckAuthWithClient invalid authtoken", " clientId: ", modifyPinsReq.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("UpdatePins (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", modifyPinsReq.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("UpdatePins (controller), reqParams:", helpers.LogStructAsJSON(modifyPinsReq), " uccId: ", modifyPinsReq.ClientId, "requestId:", reqH.RequestId)
	//call service
	code, resp := pinProvider.UpdatePins(modifyPinsReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: UpdatePins requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// AddPins
// @Tags space pins V1
// @Description Update Pins
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.AddPinReq true "pins"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/pins/addPins [POST]
func AddPins(c *gin.Context) {
	var addPinsReq models.AddPinReq

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("AddPins (controller), error decoding header, error:", err)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&addPinsReq)
	if err != nil {
		loggerconfig.Error("AddPins (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(addPinsReq.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("AddPins (controller) CheckAuthWithClient invalid authtoken", " clientId: ", addPinsReq.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("AddPins (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", addPinsReq.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("AddPins (controller), reqParams:", helpers.LogStructAsJSON(addPinsReq), " uccId: ", addPinsReq.ClientId, "requestId:", reqH.RequestId)
	//call service
	code, resp := pinProvider.AddPins(addPinsReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: AddPins requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DeletePins
// @Tags space pins V1
// @Description Update Pins
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.DeletePins true "pins"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/pins/deletePins [POST]
func DeletePins(c *gin.Context) {
	var deletePinsReq models.DeletePins

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("DeletePins (controller), error decoding header, error:", err)
	}

	err := json.NewDecoder(c.Request.Body).Decode(&deletePinsReq)
	if err != nil {
		loggerconfig.Error("DeletePins (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(deletePinsReq.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("DeletePins (controller) CheckAuthWithClient invalid authtoken", " clientId: ", deletePinsReq.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("DeletePins (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", deletePinsReq.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("DeletePins (controller), reqParams:", helpers.LogStructAsJSON(deletePinsReq), " uccId: ", deletePinsReq.ClientId, "requestId:", reqH.RequestId)
	//call service
	code, resp := pinProvider.DeletePins(deletePinsReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: DeletePins requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
