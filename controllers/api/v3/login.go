package v3

import (
	"encoding/json"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
	"strings"

	masker "github.com/ggwhite/go-masker"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var theLoginProviderV3 models.LoginProviderV3
var maskObj *masker.Masker

func InitLoginProviderV3(provider models.LoginProviderV3) {
	defer models.HandlePanic()
	theLoginProviderV3 = provider
}

// SetPassword
// @Tags space auth V3
// @Description Set Password Request V3
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-Platform header string true "P-Platform" Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.SetPasswordRequest true "login"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/authapis/setPassword [PUT]
func SetPassword(c *gin.Context) {

	var reqParams models.SetPasswordRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("SetPassword V3 (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("SetPassword V3 (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("SetPassword V3 (controller), Error validating struct: ", err, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId: ", requestH.RequestId, " clientVersion:", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	maskedReq, err := maskObj.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("In Controller SetPassword V3 Error in masking request err: ", err, " clientId: ", requestH.ClientId, " requestid = ", requestH.RequestId)
		return
	}

	loggerconfig.Info("SetPassword V3 (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId)
	code, resp := theLoginProviderV3.SetPassword(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: SetPassword V3 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ValidateToken
// @Tags space auth V3
// @Description Set Validate Token V3
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-Platform header string true "P-Platform" Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param token query string true "token Query Parameter" dataType(string)
// @Param userId query string true "userId Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/authapis/validateToken [GET]
func ValidateToken(c *gin.Context) {
	var reqParams models.ValidateTokenRequest
	token := c.Query("token")
	userId := c.Query("userId")
	if token == "" || userId == "" {
		loggerconfig.Error("ValidateToken (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.Token = token
	reqParams.UserId = userId
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("ValidateToken V3 (controller), Empty Device Type requestId: ", requestH.RequestId, "userId: ", reqParams.UserId, "ClientID: ", requestH.ClientId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}
	loggerconfig.Info("ValidateToken V3 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "userId: ", reqParams.UserId, "ClientID: ", requestH.ClientId)
	code, resp := theLoginProviderV3.ValidateToken(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: ValidateToken v3 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ForgetResetTwoFa
// @Tags space auth V3
// @Description Forget Reset TwoFa V3
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-Platform header string true "P-Platform" Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.ForgetResetTwoFaRequest true "login"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/authapis/forgetResetTwoFa [PUT]
func ForgetResetTwoFa(c *gin.Context) {
	var reqParams models.ForgetResetTwoFaRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ForgetResetTwoFa V3 (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("ForgetResetTwoFa V3 (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientID)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}
	loggerconfig.Info("ForgetResetTwoFa V3 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
	code, resp := theLoginProviderV3.ForgetResetTwoFa(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: ForgetResetTwoFa v3 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ValidateLoginOtp
// @Tags space auth V3
// @Description ValidateLoginOtp V3
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string true "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.ValidateLoginOtpV2Req true "validateLoginOtp"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/authapis/validateLoginOtp [PUT]
func ValidateLoginOtp(c *gin.Context) {

	var reqParams models.ValidateLoginOtpV2Req
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ValidateLoginOtp V3 (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("ValidateLoginOtp V3 (controller), Empty Device Type requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("ValidateLoginOtp V3, reqParams:", helpers.LogStructAsJSON(reqParams), " requestId:", requestH.RequestId)

	code, resp := theLoginProviderV3.ValidateLoginOtpV2(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: ValidateLoginOtp v3 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// SetupBiometric
// @Tags space auth V3
// @Description SetupBiometric V3
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string true "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.SetupBiometricV2Req true "setupBiometric"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/authapis/setupBiometric [PUT]
func SetupBiometric(c *gin.Context) {

	var reqParams models.SetupBiometricV2Req
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("SetupBiometric V3 (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("SetupBiometric V3 (controller), Empty Device Type requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("SetupBiometric V3, reqParams:", helpers.LogStructAsJSON(reqParams), " requestId:", requestH.RequestId)

	code, resp := theLoginProviderV3.SetupBiometricV2(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: SetupBiometric v3 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DisableBiometric
// @Tags space auth V3
// @Description DisableBiometric V3
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string true "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.DisableBiometricV2Req true "disableBiometric"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/authapis/disableBiometric [DELETE]
func DisableBiometric(c *gin.Context) {

	var reqParams models.DisableBiometricV2Req
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("DisableBiometric V3 (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("DisableBiometric V3 (controller), Empty Device Type requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("DisableBiometric V3, reqParams:", helpers.LogStructAsJSON(reqParams), " requestId:", requestH.RequestId)

	code, resp := theLoginProviderV3.DisableBiometricV2(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: DisableBiometric v3 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// Login By Email Otp
// @Tags space auth V3
// @Description Login By Email Otp V3
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param request body models.LoginByEmailOtpReq true "loginByEmailOtp"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Success 200 {object} apihelpers.APIRes{data=models.LoginByEmailOtpRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/authapis/loginByEmailOtp [POST]
func LoginByEmailOtp(c *gin.Context) {

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	var reqParams models.LoginByEmailOtpReq

	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("LoginByEmailOtp V3 (controller), Login error decoding body, error:", err, "requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("LoginByEmailOtp V3 (controller), Empty Device Type requestId: ", requestH.RequestId, " email: ", reqParams.Email, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("LoginByEmailOtp V3 (controller), error validating request, error:", err, "requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	reqParams.Email = strings.ToLower(strings.TrimSpace(reqParams.Email))

	loggerconfig.Info("LoginByEmailOtp V3 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " email: ", reqParams.Email, "requestId:", requestH.RequestId)
	code, resp := theLoginProviderV3.LoginByEmailOtp(reqParams, requestH)
	logDetail := "email: " + reqParams.Email + " function: LoginByEmailOtp V3 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}
