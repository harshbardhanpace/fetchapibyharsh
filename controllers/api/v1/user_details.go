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
)

var theUserDetailsProvider models.UserDetailsProvider

func InitUserDetailsProvider(provider models.UserDetailsProvider) {
	defer models.HandlePanic()
	theUserDetailsProvider = provider
}

// GetAllBankAccounts
// @Tags space user details V1
// @Description Get All Bank Accounts
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GetAllBankAccountsUpdatedReq true "user details"
// @Success 200 {object} apihelpers.APIRes{data=models.GetAllBankAccountsUpdatedRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/userDetails/getAllBankAccounts [POST]
func GetAllBankAccounts(c *gin.Context) {
	var reqParams models.GetAllBankAccountsUpdatedReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetAllBankAccounts (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetAllBankAccounts (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetAllBankAccounts (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientId, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("GetAllBankAccounts (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("GetAllBankAccounts (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientId, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("GetAllBankAccounts (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId)
	code, resp := theUserDetailsProvider.GetAllBankAccountsUpdated(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: GetAllBankAccounts requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetUserId
// @Tags space user details V1
// @Description Get User Id
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GetUserIdReq true "Get User Id"
// @Success 200 {object} apihelpers.APIRes{data=models.GetUserIdRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/userDetails/getUserId [POST]
func GetUserId(c *gin.Context) {
	var reqParams models.GetUserIdReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GetUserId (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetUserId (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetUserId (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetUserId (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theUserDetailsProvider.GetUserId(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: GetUserId requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// UserNotifications
// @Tags space user details V1
// @Description User Notifications
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param clientId query string true "clientId" dataType(string)
// @Param page query string false "page" dataType(int)
// @Param size query string false "size" dataType(int)
// @Success 200 {object} apihelpers.APIRes{data=models.MongoNotificationStore}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/userDetails/userNotifications [GET]
func UserNotifications(c *gin.Context) {

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	var reqParams models.UserNotificationsReq

	clientID := c.Query("clientId")
	page := c.Query("page")
	size := c.Query("size")
	if clientID == "" {
		loggerconfig.Error("UserNotifications (controller), error parsing the query params in Get request, clientId not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.ClientId = strings.ToUpper(clientID)

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		loggerconfig.Error("UserNotifications (controller), error converting page to int type, Invalid page received: ", page, " requestId:", requestH.RequestId, "ClientID: ", reqParams.ClientId)
		pageInt = 1
	}
	reqParams.Page = pageInt

	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		loggerconfig.Error("UserNotifications (controller), error converting size to int type, Invalid size received: ", size, " requestId:", requestH.RequestId, "ClientID: ", reqParams.ClientId)
		sizeInt = 10
	}
	reqParams.PageSize = sizeInt

	if requestH.DeviceType == "" {
		loggerconfig.Error("UserNotifications (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("UserNotifications (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theUserDetailsProvider.UserNotifications(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: UserNotifications requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetClientStatus
// @Tags space user details V1
// @Description Get Client Status
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param emailId query string true "emailId" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.GetClientStatusRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/userDetails/getClientStatus [GET]
func GetClientStatus(c *gin.Context) {
	emailId := c.Query("emailId")
	if emailId == "" {
		loggerconfig.Error("GetClientStatus (controller), error parsing the query params in Get request, emailId not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetClientStatus (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("GetClientStatus (controller), reqParams emailId:", emailId, "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theUserDetailsProvider.GetClientStatus(emailId, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: GetClientStatus requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
