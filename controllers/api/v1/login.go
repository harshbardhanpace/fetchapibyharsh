package v1

import (
	"encoding/json"
	"net/http"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	masker "github.com/ggwhite/go-masker"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var theLoginProvider models.LoginProvider
var maskObj *masker.Masker

func InitLoginProvider(provider models.LoginProvider) {
	defer models.HandlePanic()
	theLoginProvider = provider
	maskObj = masker.New()
}

// Login
// @Tags space auth V1
// @Description Login by username and password of user
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param request body models.LoginRequest true "login"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Success 200 {object} apihelpers.APIRes{data=models.LoginResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/authapis/login [POST]
func Login(c *gin.Context) {

	var reqParams models.LoginRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("Login (controller), Login error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("Login (controller), Empty Device Type requestId: ", requestH.RequestId, " ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	maskedReq, err := maskObj.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("In Controller Login Error in masking request err: ", err, " clientId: ", requestH.ClientId, " requestid = ", requestH.RequestId)
		return
	}

	loggerconfig.Info("Login (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := theLoginProvider.LoginByPass(reqParams, requestH)
	logDetail := "clientId: " + reqParams.UserName + " function: Login requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// LoginByEmail
// @Tags space auth V1
// @Description Login by email and password of user
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param request body models.LoginByEmailRequest true "login"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Success 200 {object} apihelpers.APIRes{data=models.LoginResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/authapis/loginByEmail [POST]
func LoginByEmail(c *gin.Context) {

	var reqParams models.LoginByEmailRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("LoginByEmail (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("LoginByEmail (controller), Empty Device Type requestId: ", requestH.RequestId, " ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}
	loggerconfig.Info("LoginByEmail (controller), Empty Device Type requestId: ", requestH.RequestId, " ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theLoginProvider.LoginByEmail(reqParams, requestH)
	logDetail := "clientId: " + reqParams.EmailId + " function: LoginByEmail requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ValidateTwoFA
// @Tags space auth V1
// @Description Validate Two Fa Answers
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-Platform header string true "P-Platform" Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.ValidateTwoFARequest true "login"
// @Success 200 {object} apihelpers.APIRes{data=models.ValidateTwoFAResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/authapis/validateTwoFa [POST]
func ValidateTwoFA(c *gin.Context) {

	var reqParams models.ValidateTwoFARequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ValidateTwoFA (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("ValidateTwoFA (controller), Empty Device Type clientID: ", reqParams.LoginID, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("ValidateTwoFA (controller), reqParams: loginId:", reqParams.LoginID, " Twofa:", reqParams.Twofa, " twoFaToken:", maskObj.Password(reqParams.TwofaToken), " type:", reqParams.Type, " clientID: ", reqParams.LoginID, " requestId:", requestH.RequestId)
	code, resp := theLoginProvider.ValidateTwoFa(reqParams, requestH)
	logDetail := "clientId: " + reqParams.LoginID + " function: ValidateTwoFA requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// SetTwoFAPin
// @Tags space auth V1
// @Description Set Two Fa Pin of user
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.SetTwoFAPinRequest true "login"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/authapis/setTwoFaPin [POST]
func SetTwoFAPin(c *gin.Context) {

	var reqParams models.SetTwoFAPinRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("SetTwoFAPin (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("SetTwoFAPin (controller), Empty Device Type clientID: ", reqParams.LoginID, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	maskedReq, err := maskObj.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("In Controller SetTwoFAPin Error in masking request err: ", err, " clientId: ", requestH.ClientId, " requestid = ", requestH.RequestId)
		return
	}

	loggerconfig.Info("SetTwoFAPin (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)
	code, resp := theLoginProvider.SetTwoFaPin(reqParams, requestH)
	logDetail := "clientId: " + reqParams.LoginID + " function: SetTwoFAPin requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ForgotPassword
// @Tags space auth V1
// @Description ForgotPassword - Used when user forget password
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-Platform header string true "P-Platform" Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.ForgotPasswordRequest true "login"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/authapis/forgotPassword [POST]
func ForgotPassword(c *gin.Context) {

	var reqParams models.ForgotPasswordRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ForgotPassword (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("ForgotPassword (controller), Empty Device Type  clientID: ", reqParams.LoginID, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	maskedReq, err := maskObj.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("In Controller ForgotPassword Error in masking request err: ", err, " clientId: ", requestH.ClientId, " requestid = ", requestH.RequestId)
		return
	}

	loggerconfig.Info("ForgotPassword (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)

	code, resp := theLoginProvider.ForgetPassword(reqParams, requestH)
	logDetail := "clientId: " + reqParams.LoginID + " function: ForgotPassword requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ForgetPasswordEmail
// @Tags space auth V1
// @Description ForgetPasswordEmail - Used when user forget password
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-Platform header string true "P-Platform" Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.ForgetResetEmailRequest true "login"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/authapis/forgetPasswordEmail [POST]
func ForgetPasswordEmail(c *gin.Context) {

	var reqParams models.ForgetResetEmailRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ForgetPasswordEmail (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("ForgetPasswordEmail (controller), Empty Device Type  clientID: ", requestH.ClientId, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("ForgetPasswordEmail (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	maskedReq, err := maskObj.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("In Controller ForgetPasswordEmail Error in masking request err: ", err, " clientId: ", requestH.ClientId, " requestid = ", requestH.RequestId)
		return
	}

	loggerconfig.Info("ForgetPasswordEmail (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)

	code, resp := theLoginProvider.ForgetPasswordEmail(reqParams, requestH)
	logDetail := "emailId: " + reqParams.EmailId + " function: ForgetPasswordEmail requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// SetPassword
// @Tags space auth V1
// @Description SetPassword - Used to set the password if it is not set
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-Platform header string true "P-Platform" Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.SetPasswordRequest true "login"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/authapis/setPassword [POST]
func SetPassword(c *gin.Context) {

	var reqParams models.SetPasswordRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("SetPassword (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("SetPassword (controller), Empty Device Type requestId: ", requestH.RequestId, " ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	maskedReq, err := maskObj.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("In Controller SetPassword Error in masking request err: ", err, " clientId: ", requestH.ClientId, " requestid = ", requestH.RequestId)
	}

	loggerconfig.Info("SetPassword (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)

	code, resp := theLoginProvider.SetPassword(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: SetPassword requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ValidateToken
// @Tags space auth V1
// @Description ValidateToken - Used to validate the token
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-Platform header string true "P-Platform" Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.ValidateTokenRequest true "login"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/authapis/validateToken [POST]
func ValidateToken(c *gin.Context) {
	var reqParams models.ValidateTokenRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ValidateToken (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("ValidateToken (controller), Empty Device Type requestId: ", requestH.RequestId, " userId: ", reqParams.UserId, " ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("ValidateToken (controller), reqParams: token:", maskObj.Password(reqParams.Token), " userId:", reqParams.UserId, " requestId:", requestH.RequestId, " userId: ", reqParams.UserId, " ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theLoginProvider.ValidateToken(reqParams, requestH)
	logDetail := "clientId: " + reqParams.UserId + " function: ValidateToken requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ForgetResetTwoFa
// @Tags space auth V1
// @Description ForgetResetTwoFa - Used when TwoFa is forgot
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-Platform header string true "P-Platform" Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.ForgetResetTwoFaRequest true "login"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/authapis/forgetResetTwoFa [POST]
func ForgetResetTwoFa(c *gin.Context) {
	var reqParams models.ForgetResetTwoFaRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ForgetResetTwoFa (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("ForgetResetTwoFa (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	maskedReq, err := maskObj.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("In Controller ForgetResetTwoFa Error in masking request err: ", err, " clientId: ", requestH.ClientId, " requestid = ", requestH.RequestId)
		return
	}

	loggerconfig.Info("ForgetResetTwoFa (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)

	code, resp := theLoginProvider.ForgetResetTwoFa(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: ForgetResetTwoFa requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// ForgetResetTwoFaEmail
// @Tags space auth V1
// @Description ForgetResetTwoFa Email - Used when TwoFa is forgot
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-Platform header string true "P-Platform" Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.ForgetResetEmailRequest true "login"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/authapis/forgetResetTwoFaEmail [POST]
func ForgetResetTwoFaEmail(c *gin.Context) {
	var reqParams models.ForgetResetEmailRequest
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("ForgetResetTwoFaEmail (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("ForgetResetTwoFaEmail (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("ForgetResetTwoFaEmail (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	maskedReq, err := maskObj.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("In Controller ForgetResetTwoFaEmail Error in masking request err: ", err, " clientId: ", requestH.ClientId, " requestid = ", requestH.RequestId)
		return
	}

	loggerconfig.Info("ForgetResetTwoFaEmail (controller), reqParams:", helpers.LogStructAsJSON(maskedReq), " uccId: ", requestH.ClientId, "requestId:", requestH.RequestId)

	code, resp := theLoginProvider.ForgetResetTwoFaEmail(reqParams, requestH)
	logDetail := "emailId: " + reqParams.EmailId + " function: ForgetResetTwoFaEmail requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GuestUserStatus
// @Tags space auth V1
// @Description GuestUserStatus - Provides user status - trading, kyc completed, unidentified
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-Platform header string true "P-Platform" Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.GuestUserStatusReq true "GuestUserStatus"
// @Success 200 {object} apihelpers.APIRes{data=models.GuestUserStatusRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/authapis/guestUserStatus [POST]
func GuestUserStatus(c *gin.Context) {
	var reqParams models.GuestUserStatusReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GuestUserStatus (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("GuestUserStatus (controller), Empty Device Type  clientID: ", requestH.ClientId, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("GuestUserStatus (controller), reqParams: email:", reqParams.Email, " clientID: ", requestH.ClientId, " requestId:", requestH.RequestId)
	code, resp := theLoginProvider.GuestUserStatus(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: GuestUserStatus requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// QRWebLogin
// @Tags space QRWebLogin V1
// @Description QRWebLogin - Login with help of QR
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.LoginWithQRReq true "login"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/qr/webLogin [POST]
func QRWebLogin(c *gin.Context) {
	var reqParams models.LoginWithQRReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("QRWebLogin (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("QRWebLogin (controller), Empty Device Type clientID: ", requestH.ClientId, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}
	loggerconfig.Info("QRWebLogin (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId, " requestId:", requestH.RequestId)

	code, resp := theLoginProvider.QRWebLogin(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: QRWebLogin requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// UnblockUser
// @Tags space auth V1
// @Description Unblock User
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.UnblockUserReq true "UnblockUser"
// @Success 200 {object} apihelpers.APIRes{data=models.UnblockUserRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/authapis/unblockUser [POST]
func UnblockUser(c *gin.Context) {
	var reqParams models.UnblockUserReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("UnblockUser (controller), error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("UnblockUser (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("UnblockUser (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("UnblockUser (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theLoginProvider.UnblockUser(reqParams, requestH)
	logDetail := "clientId: " + reqParams.LoginID + " function: UnblockUser requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// CreateApp
// @Tags space Oauth2 V1
// @Description CreateApp - create app to generate appId and secret
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.CreateAppReq true "login"
// @Success 200 {object} apihelpers.APIRes{data=models.CreateAppRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/authapis/createApp [POST]
func CreateApp(c *gin.Context) {
	var reqParams models.CreateAppReq
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("CreateApp (controller), CreateApp error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if requestH.DeviceType == "" {
		loggerconfig.Error("CreateApp (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("CreateApp (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theLoginProvider.CreateApp(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: CreateApp requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// FetchApps
// @Tags space Oauth2 V1
// @Description FetchApps - Login with help of FetchApps
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param clientId path string true "Client ID"
// @Success 200 {object} apihelpers.APIRes{data=models.CreateAppRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/authapis/fetchApps/{clientId} [GET]
func FetchApps(c *gin.Context) {
	clientId := c.Param("clientId")

	if clientId == "" {
		loggerconfig.Error("FetchApps (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchApps (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(clientId, requestH.Authorization)

	if !tokenValidStatus {
		loggerconfig.Error("FetchApps (controller) CheckAuthWithClient invalid authtoken ", "platform: ", requestH.Platform, " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}

	if !matchStatus {
		loggerconfig.Error("FetchApps (controller) CheckAuthWithClient difference in authtoken-clientId and clientId ", "platform: ", requestH.Platform, " clientId: ", requestH.ClientId, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("DeleteApp (controller), reqParams: clientId:", clientId, " clientID: ", requestH.ClientId, " requestId:", requestH.RequestId)
	code, resp := theLoginProvider.FetchApps(clientId, requestH)
	logDetail := "clientId: " + clientId + " function: FetchApps requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// DeleteApp
// @Tags space Oauth2 V1
// @Description DeleteApp - Delete Apps from developer console
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param appId path string true "App ID"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/authapis/deleteApp/{appId} [DELETE]
func DeleteApp(c *gin.Context) {
	appId := c.Param("appId")

	if appId == "" {
		loggerconfig.Error("DeleteApp (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("DeleteApp (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	loggerconfig.Info("DeleteApp (controller), reqParams: appId:", appId, " clientID: ", requestH.ClientId, " requestId:", requestH.RequestId)
	code, resp := theLoginProvider.DeleteApp(appId, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: DeleteApp requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// The method will get called by the oauth/auth API as a redirect url and auth code will be handled here.
func HandleAuthCode(c *gin.Context) {
	authCode := c.Query("code")
	appState := c.Query("state")

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if authCode == "" || appState == "" {
		loggerconfig.Error("HandleAuthCode (controller), Missing parameters", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	loggerconfig.Info("HandleAuthCode (controller), reqParams: authcode:", authCode, "appState", appState, " clientID: ", requestH.ClientId, " requestId:", requestH.RequestId)

	code, resp := theLoginProvider.HandleAuthCode(authCode, appState, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: HandleAuthCode requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GetAccessToken
// @Tags space Oauth2 V1
// @Description GetAccessToken - Polling API to get access token
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.GetAccessTokenReq true "login"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/authapis/getAccessToken [POST]
func GetAccessToken(c *gin.Context) {
	var reqParams models.GetAccessTokenReq

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)

	if err != nil {
		loggerconfig.Error("GetAccessToken (controller), GetAccessToken error decoding body, error:", err, "requestId:", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("GetAccessToken (controller), reqParams: loginId:", reqParams.AppState, " clientID: ", requestH.ClientId, " requestId:", requestH.RequestId)
	code, resp := theLoginProvider.GetAccessToken(reqParams, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: GetAccessToken requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
