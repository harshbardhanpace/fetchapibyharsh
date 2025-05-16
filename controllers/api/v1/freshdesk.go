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

var theFreshdeskProvider models.FreshdeskProvider

func InitFreshdeskProvider(provider models.FreshdeskProvider) {
	defer models.HandlePanic()
	theFreshdeskProvider = provider
}

// CreateFreshdeskTicket
// @Tags Freshdesk
// @Description Create a new Freshdesk ticket
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param Authorization header string true "Authorization Header"
// @Param request body models.FreshdeskTicketReq true "Ticket Request Body"
// @Success 200 {object} apihelpers.APIRes{data=models.FreshdeskTicketResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/support/ticket/create [POST]
func CreateFreshdeskTicket(c *gin.Context) {
	var reqParams models.FreshdeskTicketReq

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("CreateFreshdeskTicket (controller), Error decoding request body, platform:", requestH.Platform, "err:", err, "clientId:", requestH.ClientId, "requestId:", requestH.RequestId)

		apihelpers.ErrorMessage(c, "Invalid request")
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("CreateFreshdeskTicket (controller), Empty Device Type, platform:", requestH.Platform, "clientId:", requestH.ClientId, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("CreateFreshdeskTicket (controller), Validation error, platform:", requestH.Platform, "err:", err, "clientId:", requestH.ClientId, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if reqParams.Phone != "" {
		if len(reqParams.Phone) != 10 || !helpers.IsAllDigits(reqParams.Phone) {
			loggerconfig.Error("CreateFreshdeskTicket (controller), Invalid Phone number, platform:", requestH.Platform, "clientId:", requestH.ClientId, "requestId:", requestH.RequestId)
			apihelpers.ErrorMessage(c, "Phone number must be 10 digits")
			return
		}
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(requestH.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("CreateFreshdeskTicket (controller), CheckAuthWithClient invalid authtoken, platform:", requestH.Platform, "clientId:", requestH.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("CreateFreshdeskTicket (controller), CheckAuthWithClient mismatch in authtoken-clientId and clientId, platform:", requestH.Platform, "clientId:", requestH.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	// Call the service to create a Freshdesk ticket
	code, resp := theFreshdeskProvider.CreateFreshdeskTicket(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: CreateFreshdeskTicket requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
