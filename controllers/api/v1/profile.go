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

var theProfileProvider models.ProfileProvider

func InitProfileProvider(provider models.ProfileProvider) {
	defer models.HandlePanic()
	theProfileProvider = provider
}

func GetProfileProvider() models.ProfileProvider {
	return theProfileProvider
}

// Profile
// @Tags space profile V1
// @Description GetProfile - Provide profile details of user
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.ProfileRequest true "profile"
// @Success 200 {object} apihelpers.APIRes{data=models.ProfileResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/user/profile/getProfile [POST]
func GetProfile(c *gin.Context) {
	var reqParams models.ProfileRequest
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	errr := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if errr != nil {
		loggerconfig.Error("GetProfile (controller), error decoding body, error:", errr, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetProfile (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetProfile (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("GetProfile (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("GetProfile (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("GetProfile (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
	code, resp := theProfileProvider.GetProfile(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}

// SendAFOtp
// @Tags space profile V1
// @Description SendAFOtp - Send Account Freeze Otp
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.SendAFOtpReq true "profile"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/user/profile/sendAFOtp [POST]
func SendAFOtp(c *gin.Context) {
	var reqParams models.SendAFOtpReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	errr := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if errr != nil {
		loggerconfig.Error("SendAFOtp (controller), error decoding body, error:", errr, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.ClientID = strings.ToUpper(reqParams.ClientID)

	if requestH.DeviceType == "" {
		loggerconfig.Error("SendAFOtp (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("SendAFOtp (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("SendAFOtp (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("SendAFOtp (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("SendAFOtp (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
	code, resp := theProfileProvider.SendAFOtp(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// VerifyAFOtp
// @Tags space profile V1
// @Description VerifyAFOtp - Verify Account Freeze Otp
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.VerifyAFOtpReq true "profile"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/user/profile/verifyAFOtp [POST]
func VerifyAFOtp(c *gin.Context) {
	var reqParams models.VerifyAFOtpReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	errr := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if errr != nil {
		loggerconfig.Error("VerifyAFOtp (controller), error decoding body, error:", errr, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	reqParams.ClientID = strings.ToUpper(reqParams.ClientID)

	if requestH.DeviceType == "" {
		loggerconfig.Error("VerifyAFOtp (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("VerifyAFOtp (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("VerifyAFOtp (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("VerifyAFOtp (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("VerifyAFOtp (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
	code, resp := theProfileProvider.VerifyAFOtp(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// AccountFreeze
// @Tags space profile V1
// @Description AccountFreeze - Account Freeze of user
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.AccountFreezeReq true "profile"
// @Success 200 {object} apihelpers.APIRes{data=models.AccountFreezeRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/user/profile/accountFreeze [POST]
func AccountFreeze(c *gin.Context) {
	var reqParams models.AccountFreezeReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	errr := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if errr != nil {
		loggerconfig.Error("AccountFreeze (controller), error decoding body, error:", errr, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	reqParams.ClientID = strings.ToUpper(reqParams.ClientID)

	if requestH.DeviceType == "" {
		loggerconfig.Error("AccountFreeze (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("AccountFreeze (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("AccountFreeze (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("AccountFreeze (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("AccountFreeze (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
	code, resp := theProfileProvider.AccountFreeze(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}
