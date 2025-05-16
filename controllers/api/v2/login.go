package v2

import (
	"encoding/json"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/ggwhite/go-masker"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

var theLoginProviderV2 models.LoginProviderV2
var maskObj *masker.Masker

func InitLoginProviderV2(provider models.LoginProviderV2) {
	defer models.HandlePanic()
	theLoginProviderV2 = provider
	maskObj = masker.New()
}

// Login
// @Tags space auth V2
// @Description Login by id and secret
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string true "P-DeviceId Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param DeviceToken header string false "DeviceToken"
// @Param request body models.LoginV2Request true "login"
// @Success 200 {object} apihelpers.APIRes{data=models.LoginV2Response}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/authapis/login [POST]
func Login(c *gin.Context) {

	var reqParams models.LoginV2Request
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("Login V2 (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("Login V2 (controller), Empty Device Type requestId: ", requestH.RequestId, "loginId: ", reqParams.ID)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	maskedReq, err := maskObj.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("In Controller Login V2 Error in masking request err: ", err, " clientId: ", requestH.ClientId, " requestid = ", requestH.RequestId)
		return
	}

	loggerconfig.Info("Login V2 (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)

	code, resp := theLoginProviderV2.LoginV2(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: LoginV2 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// Twofa
// @Tags space auth V2
// @Description ValidateTwofa
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string true "P-DeviceId Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.ValidateTwofaV2Req true "validateTwofa"
// @Success 200 {object} apihelpers.APIRes{data=models.ValidateTwofaV2Res}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/authapis/validateTwofa [POST]
func ValidateTwofa(c *gin.Context) {

	var reqParams models.ValidateTwofaV2Req
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ValidateTwofa V2 (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("ValidateTwofa V2 (controller), Empty Device Type requestId: ", requestH.RequestId, "loginId: ", reqParams.LoginID)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("ValidateTwofa V2 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.LoginID, " requestId:", requestH.RequestId)

	code, resp := theLoginProviderV2.ValidateTwofaV2(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: ValidateTwofaV2 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// SetupTotp
// @Tags space auth V2
// @Description SetupTotp
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.SetupTotpV2Req true "setupTotp"
// @Success 200 {object} apihelpers.APIRes{data=models.SetupTotpV2Res}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/authapis/setupTotp [POST]
func SetupTotp(c *gin.Context) {

	var reqParams models.SetupTotpV2Req
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("SetupTotp V2 (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("SetupTotp V2 (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("SetupTotp V2 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, " requestId:", requestH.RequestId)

	code, resp := theLoginProviderV2.SetupTotpV2(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: SetupTotp requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ChooseTwofa
// @Tags space auth V2
// @Description ChooseTwofa
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.ChooseTwofaV2Req true "chooseTwofa"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/authapis/chooseTwofa [POST]
func ChooseTwofa(c *gin.Context) {

	var reqParams models.ChooseTwofaV2Req
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ChooseTotp V2 (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("ChooseTotp V2 (controller), Empty Device Type requestId: ", requestH.RequestId, "loginId: ", reqParams.LoginID)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	maskedReq, err := maskObj.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("In Controller ChooseTwofa V2 Error in masking request err: ", err, " clientId: ", requestH.ClientId, " requestid = ", requestH.RequestId)
		return
	}

	loggerconfig.Info("ChooseTotp V2 (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), " uccId: ", reqParams.LoginID, " requestId:", requestH.RequestId)

	code, resp := theLoginProviderV2.ChooseTwofaV2(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: ChooseTwofa requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ForgetTotp
// @Tags space auth V2
// @Description ForgetTotp
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.ForgetTotpV2Req true "forgetTotp"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/authapis/forgetTotp [POST]
func ForgetTotp(c *gin.Context) {

	var reqParams models.ForgetTotpV2Req
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ForgetTotp V2 (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("ForgetTotp V2 (controller), Empty Device Type requestId: ", requestH.RequestId, "loginId: ", reqParams.LoginID)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	maskedReq, err := maskObj.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("In Controller ForgetTotp V2 Error in masking request err: ", err, " clientId: ", requestH.ClientId, " requestid = ", requestH.RequestId)
		return
	}

	loggerconfig.Info("ForgetTotp v2 (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)

	code, resp := theLoginProviderV2.ForgetTotpV2(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: ForgetTotp requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ValidateLoginOtp
// @Tags space auth V2
// @Description ValidateLoginOtp
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string true "P-DeviceId Header"
// @Param DeviceToken header string false "DeviceToken"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.ValidateLoginOtpV2Req true "validateLoginOtp"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/authapis/validateLoginOtp [POST]
func ValidateLoginOtp(c *gin.Context) {

	var reqParams models.ValidateLoginOtpV2Req
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ValidateLoginOtp V2 (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("ValidateLoginOtp V2 (controller), Empty Device Type requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("ValidateLoginOtp V2, reqParams ReferenceToken:", reqParams.ReferenceToken, " requestId:", requestH.RequestId)

	code, resp := theLoginProviderV2.ValidateLoginOtpV2(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: ValidateLoginOtp requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// SetupBiometric
// @Tags space auth V2
// @Description SetupBiometric
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string true "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.SetupBiometricV2Req true "setupBiometric"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/authapis/setupBiometric [POST]
func SetupBiometric(c *gin.Context) {

	var reqParams models.SetupBiometricV2Req
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ValidateLoginOtp V2 (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("ValidateLoginOtp V2 (controller), Empty Device Type requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	maskedReq, err := maskObj.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("In Controller SetupBiometric V2 Error in masking request err: ", err, " clientId: ", requestH.ClientId, " requestid = ", requestH.RequestId)
		return
	}

	loggerconfig.Info("SetupBiometric V2 (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)

	code, resp := theLoginProviderV2.SetupBiometricV2(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: SetupBiometric requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DisableBiometric
// @Tags space auth V2
// @Description DisableBiometric
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string true "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.DisableBiometricV2Req true "disableBiometric"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/authapis/disableBiometric [POST]
func DisableBiometric(c *gin.Context) {

	var reqParams models.DisableBiometricV2Req
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("DisableBiometric V2 (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("DisableBiometric V2 (controller), Empty Device Type requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	maskedReq, err := maskObj.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("In Controller DisableBiometric V2 Error in masking request err: ", err, " clientId: ", requestH.ClientId, " requestid = ", requestH.RequestId)
		return
	}

	loggerconfig.Info("DisableBiometric V2 (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)

	code, resp := theLoginProviderV2.DisableBiometricV2(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: DisableBiometric requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ForgotPassword
// @Tags space auth V2
// @Description ForgotPassword - Used when user forget password
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-Platform header string true "P-Platform" Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.ForgotPasswordV2Request true "login"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/authapis/forgotPassword [POST]
func ForgotPasswordV2(c *gin.Context) {

	var reqParams models.ForgotPasswordV2Request
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ForgotPasswordV2 (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("ForgotPasswordV2 (controller), Empty Device Type  clientID: ", reqParams.ClientID, "emailID: ", reqParams.EmailID, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	maskedReq, err := maskObj.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("In Controller ForgotPasswordV2 Error in masking request err: ", err, " clientId: ", requestH.ClientId, " requestid = ", requestH.RequestId)
		return
	}

	loggerconfig.Info("ForgotPasswordV2 (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)

	code, resp := theLoginProviderV2.ForgetPasswordV2(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: ForgotPasswordV2 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// UnblockUser
// @Tags space auth V2
// @Description Unblock User
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.UnblockUserV2Req true "UnblockUser"
// @Success 200 {object} apihelpers.APIRes{data=models.UnblockUserRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/authapis/unblockUser [POST]
func UnblockUserV2(c *gin.Context) {
	var reqParams models.UnblockUserV2Req
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("UnblockUserV2 (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("UnblockUserV2 (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("UnblockUserV2 (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	maskedReq, err := maskObj.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("In Controller UnblockUserV2 Error in masking request err: ", err, " clientId: ", requestH.ClientId, " requestid = ", requestH.RequestId)
		return
	}

	loggerconfig.Info("UnblockUserV2 (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theLoginProviderV2.UnblockUserV2(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: UnblockUserV2 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetAccessTokenV2
// @Tags space Oauth2 V2
// @Description GetAccessToken - Polling API to get access token
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion header string false "P-ClientVersion Header"
// @Param request body models.GetAccessTokenV2Req true "GetAccessTokenV2"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/authapis/getAccessToken [POST]
func GetAccessTokenV2(c *gin.Context) {
	var reqParams models.GetAccessTokenV2Req

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)

	if err != nil {
		loggerconfig.Error("GetAccessTokenV2 (controller), GetAccessTokenV2 error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetAccessTokenV2 (controller), reqParams: loginId:", reqParams.AppState, " clientID: ", requestH.ClientId, " requestId:", requestH.RequestId)
	code, resp := theLoginProviderV2.GetAccessTokenV2(reqParams, requestH)
	apihelpers.CustomResponse(c, code, resp)
}
