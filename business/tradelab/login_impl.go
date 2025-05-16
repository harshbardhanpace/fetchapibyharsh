package tradelab

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/db"
	"space/dbops"
	"space/helpers"
	"space/helpers/cache"
	"space/loggerconfig"
	"space/models"

	// "reflect"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoginObj struct {
	tradeLabURL      string
	tradelaboAuthURL string
	mongodb          db.MongoDatabase
	redisCli         cache.RedisCache
}

func InitLogin(mongodb db.MongoDatabase, redisCli cache.RedisCache) LoginObj {
	defer models.HandlePanic()

	loginObj := LoginObj{
		tradeLabURL:      constants.TLURL,
		tradelaboAuthURL: constants.TLOauthUrl,
		mongodb:          mongodb,
		redisCli:         redisCli,
	}

	return loginObj
}

func (obj LoginObj) LoginByEmail(loginReq models.LoginByEmailRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	var apiRes apihelpers.APIRes

	var clientDetails models.MongoClientsDetails
	err := dbops.MongoDaoRepo.FindOne(constants.CLIENTDETAILS, bson.M{"email": bson.M{"$regex": "^" + loginReq.EmailId + "$", "$options": "i"}}, &clientDetails)
	if err != nil && err.Error() != "mongo: no documents in result" {
		loggerconfig.Error("Alert Severity:P0-Critical, LoginByEmail Unable to fetch the client-details Data for emailID = ", loginReq.EmailId, " error :", err, " requestid=", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	if err != nil && err.Error() == constants.MongoNoDocError {
		loggerconfig.Error("LoginByEmail, EmailId Does not exist in client-details collection emailID :", loginReq.EmailId, " error:", err, " requestId:", reqH.RequestId)
		apiRes.Status = false
		apiRes.Message = constants.ErrorCodeMap[constants.InvalidEmailId]
		apiRes.ErrorCode = constants.InvalidEmailId
		return http.StatusOK, apiRes
	}

	url := obj.tradeLabURL + LOGINURL
	var tlLoginReq TradeLabLoginReq
	tlLoginReq.Device = reqH.DeviceType
	tlLoginReq.LoginID = clientDetails.ClientID
	tlLoginReq.Password = loginReq.Password

	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlLoginReq)

	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, "")
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " LoginByEmail call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorLoginRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorLoginRes)
	if err == nil && tlErrorLoginRes.Status == TLERROR {
		loggerconfig.Error("LoginByEmail error tl api 1=", err, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorLoginRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorLoginRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlLoginRes := TradeLabLoginRes{}
	json.Unmarshal([]byte(string(body)), &tlLoginRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " LoginByEmail error tl api 2 NOT OK error=", tlLoginRes.Message, "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlLoginRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	var loginRes models.LoginResponse
	loginRes.Alert = tlLoginRes.Data.Alert
	loginRes.AuthToken = tlLoginRes.Data.AuthToken
	loginRes.LoginID = tlLoginReq.LoginID
	loginRes.ResetPassword = tlLoginRes.Data.ResetPassword
	loginRes.ResetTwoFa = tlLoginRes.Data.ResetTwoFa

	var twoFaDetails models.TwoFaDetails

	var twoFaQuestions []models.TwoFaQuestions

	for _, v := range tlLoginRes.Data.Twofa.Questions {
		var twoFaQuestion models.TwoFaQuestions
		twoFaQuestion.Question = v.Question
		twoFaQuestion.QuestionID = v.QuestionID
		twoFaQuestions = append(twoFaQuestions, twoFaQuestion)
	}

	twoFaDetails.Questions = twoFaQuestions
	twoFaDetails.TwofaToken = tlLoginRes.Data.Twofa.TwofaToken
	twoFaDetails.Type = tlLoginRes.Data.Twofa.Type

	loginRes.Twofa = twoFaDetails
	loginRes.TwofaEnabled = tlLoginRes.Data.TwofaEnabled

	loggerconfig.Info("LoginByEmail response =", helpers.LogStructAsJSON(loginRes), "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
	apiRes.Data = loginRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}

func (obj LoginObj) LoginByPass(loginReq models.LoginRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	url := obj.tradeLabURL + LOGINURL
	var tlLoginReq TradeLabLoginReq
	tlLoginReq.Device = reqH.DeviceType
	tlLoginReq.LoginID = loginReq.UserName
	tlLoginReq.Password = loginReq.Password

	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlLoginReq)

	var apiRes apihelpers.APIRes
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, "")
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " LoginByPass call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorLoginRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorLoginRes)
	if err == nil && tlErrorLoginRes.Status == TLERROR {
		loggerconfig.Error("LoginByPass error tl api 1=", err, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorLoginRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorLoginRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlLoginRes := TradeLabLoginRes{}
	json.Unmarshal([]byte(string(body)), &tlLoginRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " LoginByPass error tl api 2 NOT OK error=", tlLoginRes.Message, loginReq.UserName, "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlLoginRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	var loginRes models.LoginResponse
	loginRes.Alert = tlLoginRes.Data.Alert
	loginRes.AuthToken = tlLoginRes.Data.AuthToken
	loginRes.LoginID = tlLoginReq.LoginID
	loginRes.ResetPassword = tlLoginRes.Data.ResetPassword
	loginRes.ResetTwoFa = tlLoginRes.Data.ResetTwoFa

	var twoFaDetails models.TwoFaDetails

	var twoFaQuestions []models.TwoFaQuestions

	for _, v := range tlLoginRes.Data.Twofa.Questions {
		var twoFaQuestion models.TwoFaQuestions
		twoFaQuestion.Question = v.Question
		twoFaQuestion.QuestionID = v.QuestionID
		twoFaQuestions = append(twoFaQuestions, twoFaQuestion)
	}

	twoFaDetails.Questions = twoFaQuestions
	twoFaDetails.TwofaToken = tlLoginRes.Data.Twofa.TwofaToken
	twoFaDetails.Type = tlLoginRes.Data.Twofa.Type

	loginRes.Twofa = twoFaDetails
	loginRes.TwofaEnabled = tlLoginRes.Data.TwofaEnabled

	loggerconfig.Info("LoginByPass response =", helpers.LogStructAsJSON(loginRes), "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
	apiRes.Data = loginRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}

func (obj LoginObj) ValidateTwoFa(validateTwoFaReq models.ValidateTwoFARequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + VERIFYTWOFA

	//fill up the TL Req
	var tlValidateTwoFaReq TradeLabValidateTwoFaRequest
	tlValidateTwoFaReq.LoginID = validateTwoFaReq.LoginID

	for i := 0; i < len(validateTwoFaReq.Twofa); i++ {
		var questions TradeLabTwoFaQuestion
		questions.Answer = validateTwoFaReq.Twofa[i].Answer
		questions.QuestionID = validateTwoFaReq.Twofa[i].QuestionID
		tlValidateTwoFaReq.Twofa = append(tlValidateTwoFaReq.Twofa, questions)
	}

	tlValidateTwoFaReq.TwofaToken = validateTwoFaReq.TwofaToken
	tlValidateTwoFaReq.Type = validateTwoFaReq.Type

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlValidateTwoFaReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, "")
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ValidateTwoFa", duration, validateTwoFaReq.LoginID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " validate two fa call api error =", err, " uccId:", validateTwoFaReq.LoginID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("validate two fa res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", validateTwoFaReq.LoginID, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlValidateTwoFaRes := TradeLabValidateTwoFaResponse{}
	json.Unmarshal([]byte(string(body)), &tlValidateTwoFaRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " validate two fa tl status not ok =", tlValidateTwoFaRes.Message, " uccId:", validateTwoFaReq.LoginID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlValidateTwoFaRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var validateTwoFaRes models.ValidateTwoFAResponse
	validateTwoFaRes.ResetPassword = tlValidateTwoFaRes.Data.ResetPassword
	validateTwoFaRes.ResetTwoFa = tlValidateTwoFaRes.Data.ResetTwoFa
	validateTwoFaRes.AuthToken = tlValidateTwoFaRes.Data.AuthToken

	loggerconfig.Info("validate two fa tl resp=", helpers.LogStructAsJSON(validateTwoFaRes), " uccId:", validateTwoFaReq.LoginID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)

	apiRes.Data = validateTwoFaRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}

func (obj LoginObj) SetTwoFaPin(setTwoFaPinReq models.SetTwoFAPinRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	url := obj.tradeLabURL + SETTWOFA

	//fill up the TL Req
	var tlSetTwoFaPinReq TradelabSetTwoFaPinRequest
	tlSetTwoFaPinReq.LoginID = setTwoFaPinReq.LoginID
	tlSetTwoFaPinReq.Pin = setTwoFaPinReq.Pin
	tlSetTwoFaPinReq.TwofaType = setTwoFaPinReq.TwofaType

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlSetTwoFaPinReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "SetTwoFaPin", duration, setTwoFaPinReq.LoginID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " SetTwoFaPin call api error =", err, " uccId:", setTwoFaPinReq.LoginID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("SetTwoFaPin call api error =", err, " statuscode: ", res.StatusCode, " uccId:", setTwoFaPinReq.LoginID, " requestId:", reqH.RequestId)
		loggerconfig.Error("SetTwoFaPin call api error =", tlErrorRes.ErrorCode, " uccId:", setTwoFaPinReq.LoginID)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlSetTwoFaPinRes := TradelabSetTwoFaPinResponse{}
	json.Unmarshal([]byte(string(body)), &tlSetTwoFaPinRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " SetTwoFaPin tl status not ok =", tlSetTwoFaPinRes.Message, " uccId:", setTwoFaPinReq.LoginID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlSetTwoFaPinRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("SetTwoFaPin  tl resp=", helpers.LogStructAsJSON(tlSetTwoFaPinRes), " uccId:", setTwoFaPinReq.LoginID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)

	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj LoginObj) ForgetPassword(forgetPasswordReq models.ForgotPasswordRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + FORGOTPASSWORD

	//fill up the TL Req
	var tlForgetPasswordReq TradeLabForgetPasswordReq
	tlForgetPasswordReq.LoginID = forgetPasswordReq.LoginID
	tlForgetPasswordReq.Pan = forgetPasswordReq.Pan

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlForgetPasswordReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, "")
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ForgetPassword", duration, forgetPasswordReq.LoginID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " ForgetPassword call api error =", err, " uccId:", forgetPasswordReq.LoginID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("ForgetPassword call api error =", tlErrorRes.ErrorCode, " statuscode: ", res.StatusCode, " uccId:", forgetPasswordReq.LoginID, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlForgetPasswordRes := TradeLabForgetPasswordRes{}
	json.Unmarshal([]byte(string(body)), &tlForgetPasswordRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " ForgetPassword tl status not ok =", tlForgetPasswordRes.Message, " uccId:", forgetPasswordReq.LoginID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlForgetPasswordRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("ForgetPassword  tl resp=", helpers.LogStructAsJSON(tlForgetPasswordRes), " uccId:", forgetPasswordReq.LoginID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)

	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}

func (obj LoginObj) ForgetPasswordEmail(req models.ForgetResetEmailRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	var clientDetails models.MongoClientsDetails
	err := dbops.MongoDaoRepo.FindOne(constants.CLIENTDETAILS, bson.M{"email": bson.M{"$regex": "^" + req.EmailId + "$", "$options": "i"}}, &clientDetails)
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("Alert Severity:P0-Critical, ForgetPasswordEmail Unable to fetch the client-details Data for emailID = ", req.EmailId, " error :", err, " requestid=", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	if err != nil && err.Error() == constants.MongoNoDocError {
		loggerconfig.Error("ForgetPasswordEmail, EmailId Does not exist in client-details collection emailID :", req.EmailId, " error:", err, " requestId:", reqH.RequestId)
		apiRes.Status = false
		apiRes.Message = constants.ErrorCodeMap[constants.InvalidEmailId]
		apiRes.ErrorCode = constants.InvalidEmailId
		return http.StatusOK, apiRes
	}

	var forgetPasswordReq models.ForgotPasswordRequest
	forgetPasswordReq.LoginID = clientDetails.ClientID
	forgetPasswordReq.Pan = req.Pan

	status, res := obj.ForgetPassword(forgetPasswordReq, reqH)
	if status != http.StatusOK {
		loggerconfig.Error("ForgetPasswordEmail in ForgetPassword status != 200: ", status, " res tl: ", helpers.LogStructAsJSON(res), " uccId:", clientDetails.ClientID, " requestId:", reqH.RequestId)
		return status, res
	}

	loggerconfig.Info("ForgetPasswordEmail success uccId:", clientDetails.ClientID, " requestId:", reqH.RequestId)

	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj LoginObj) SetPassword(setPasswordReq models.SetPasswordRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + SETPASSWORD

	//fill up the TL Req
	var tlSetPasswordReq TradeLabSetPasswordReq
	tlSetPasswordReq.OldPassword = setPasswordReq.OldPass
	tlSetPasswordReq.NewPassword = setPasswordReq.NewPass

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlSetPasswordReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPut, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "SetPassword", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " SetPassword call api error =", err, " uccId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("SetPassword call api error =", tlErrorRes.ErrorCode, " statuscode: ", res.StatusCode, " uccId:", reqH.ClientId, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlSetPasswordRes := TradeLabSetPasswordRes{}
	json.Unmarshal([]byte(string(body)), &tlSetPasswordRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " ForgetPassword tl status not ok =", tlSetPasswordRes.Message, " uccId:", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlSetPasswordRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("SetPassword  tl resp=", helpers.LogStructAsJSON(tlSetPasswordRes), " uccId:", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)

	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}

func (obj LoginObj) ValidateToken(req models.ValidateTokenRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + FETCHFUNDSURL + "?client_id=" + url.QueryEscape(req.UserId) + "&type=all"
	payload := new(bytes.Buffer) // empty payload
	var apiRes apihelpers.APIRes
	authToken := constants.BEARER + " " + req.Token
	res, err := apihelpers.CallApi(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, authToken)
	if err != nil {
		loggerconfig.Error("FetchFundsResponse call api error =%v", err, " uccId:", req.UserId, " requestId:", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.InvalidToken, res.StatusCode)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("ValidateToken res error =%v", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.UserId, " requestId:", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.InvalidToken, res.StatusCode)
	}

	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj LoginObj) ForgetResetTwoFa(forgetResetTwoFaReq models.ForgetResetTwoFaRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + FORGETRESETTWOFA + "?client_id=" + url.QueryEscape(forgetResetTwoFaReq.ClientID) + "&pan=" + forgetResetTwoFaReq.Pan

	//empty payload
	payload := new(bytes.Buffer)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPut, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ForgetResetTwoFa", duration, forgetResetTwoFaReq.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " ForgetResetTwoFa call api error =", err, "clientID: ", forgetResetTwoFaReq.ClientID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("ForgetResetTwoFa call api error =", tlErrorRes.ErrorCode, " statuscode: ", res.StatusCode, " uccId:", forgetResetTwoFaReq.ClientID, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlForgetResetTwoFaRes := TradeLabForgetResetTwoFaRes{}
	json.Unmarshal([]byte(string(body)), &tlForgetResetTwoFaRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " ForgetPassword tl status not ok =", tlForgetResetTwoFaRes.Message, " uccId:", forgetResetTwoFaReq.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlForgetResetTwoFaRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("ForgetResetTwoFa  tl resp=", helpers.LogStructAsJSON(tlForgetResetTwoFaRes), " uccId:", forgetResetTwoFaReq.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)

	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}

func (obj LoginObj) ForgetResetTwoFaEmail(req models.ForgetResetEmailRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	var clientDetails models.MongoClientsDetails
	err := dbops.MongoDaoRepo.FindOne(constants.CLIENTDETAILS, bson.M{"email": bson.M{"$regex": "^" + req.EmailId + "$", "$options": "i"}}, &clientDetails)
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("Alert Severity:P0-Critical, ForgetResetTwoFaEmail Unable to fetch the client-details Data for emailID = ", req.EmailId, " error :", err, " requestid=", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	if err != nil && err.Error() == constants.MongoNoDocError {
		loggerconfig.Error("ForgetResetTwoFaEmail, EmailId Does not exist in client-details collection emailID :", req.EmailId, " error:", err, " requestId:", reqH.RequestId)
		apiRes.Status = false
		apiRes.Message = constants.ErrorCodeMap[constants.InvalidEmailId]
		apiRes.ErrorCode = constants.InvalidEmailId
		return http.StatusOK, apiRes
	}

	var forgetResetTwoFaReq models.ForgetResetTwoFaRequest
	forgetResetTwoFaReq.ClientID = clientDetails.ClientID
	forgetResetTwoFaReq.Pan = req.Pan

	status, res := obj.ForgetResetTwoFa(forgetResetTwoFaReq, reqH)
	if status != http.StatusOK {
		loggerconfig.Error("ForgetResetTwoFaEmail in ForgetResetTwoFa status != 200: ", status, " res tl:", helpers.LogStructAsJSON(res), " uccId:", clientDetails.ClientID, " requestId:", reqH.RequestId)
		return status, res
	}

	loggerconfig.Info("ForgetResetTwoFaEmail success uccId:", clientDetails.ClientID, " requestId:", reqH.RequestId)

	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj LoginObj) LoginV2(loginReq models.LoginV2Request, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	url := obj.tradeLabURL + LOGINV2URL
	var tlLoginReq TradeLabLoginV2Req
	tlLoginReq.ChannelID = loginReq.ID
	tlLoginReq.ChannelSecret = loginReq.Secret

	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlLoginReq)

	var apiRes apihelpers.APIRes
	res, err := apihelpers.CallAPIFuncV2(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, "")
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " LoginV2 call api error =", err, " uccId:", loginReq.ID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorLoginRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorLoginRes)
	if err == nil && tlErrorLoginRes.Status == TLERROR {
		loggerconfig.Error("platform:", reqH.Platform, " LoginV2 error tl api error=", err, "status code:", res.StatusCode, " uccId:", loginReq.ID, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorLoginRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorLoginRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlLoginRes := TradeLabLoginV2Response{}
	json.Unmarshal([]byte(string(body)), &tlLoginRes)

	if res.StatusCode != http.StatusOK {
		if tlLoginRes.Message != AccountFrozen {
			loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " LoginV2 error tl api 2 NOT OK error=", tlLoginRes.Message, "status code:", res.StatusCode, " uccId:", loginReq.ID, " requestId:", reqH.RequestId)
		} else {
			loggerconfig.Error("Alert Severity:P2-Mid, platform:", reqH.Platform, " LoginV2 error tl api 2 NOT OK error=", tlLoginRes.Message, "status code:", res.StatusCode, " uccId:", loginReq.ID, " requestId:", reqH.RequestId)
		}

		apiRes.Message = tlLoginRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	// added go routine to speed up login process
	go func() {
		if loginReq.ID != "" && reqH.FCMToken != "" {
			//Store FCMToken at time of Login
			clientId, err := db.GetPgObj().StoreFCMToken(loginReq.ID, reqH.FCMToken, reqH.RequestId)
			if err != nil {
				loggerconfig.Error("Alert Severity:P1-Critical, LoginV2 Error (Login) while storing FCMToken in db: ", err.Error(), "requestId: ", reqH.RequestId)
			}

			if constants.KafkaEnable {

				redisCliObj := cache.GetRedisClientObj()

				err = redisCliObj.SAdd(constants.ClientMembers+strings.ToUpper(clientId), reqH.FCMToken)
				if err != nil {
					logrus.Error("Alert Severity:P1-Critical, LoginV2 Error (Login) adding FCM token in Redis:", err, "requestId: ", reqH.RequestId)
				}

				// Push client ID into list with key as subscription_id
				subscriptionKey := constants.SubscriptionKey
				err = redisCliObj.LPush(subscriptionKey, strings.ToUpper(clientId))
				if err != nil {
					logrus.Error("Alert Severity:P1-Critical, LoginV2 Error (Login) pushing client ID in Redis list:", err, "requestId: ", reqH.RequestId)
				}

				logrus.Info("LoginV2 FCM Token details stored successfully in redis for client ID: ", loginReq.ID, "requestId: ", reqH.RequestId)

			}

			logrus.Info("LoginV2 FCM Token details stored successfully in db for client ID: ", loginReq.ID, "requestId: ", reqH.RequestId)

		}
	}()

	var loginRes models.LoginV2Response
	loginRes.Alert = tlLoginRes.Data.Alert
	loginRes.AuthToken = tlLoginRes.Data.AuthToken
	loginRes.CheckPan = tlLoginRes.Data.CheckPan
	loginRes.LoginID = tlLoginRes.Data.LoginID
	loginRes.Name = tlLoginRes.Data.Name
	loginRes.ReferenceToken = tlLoginRes.Data.ReferenceToken
	loginRes.ResetPassword = tlLoginRes.Data.ResetPassword
	loginRes.ResetTwoFa = tlLoginRes.Data.ResetTwoFa

	var twoFaDetails models.TwoFaDetails

	var twoFaQuestions []models.TwoFaQuestions

	for _, v := range tlLoginRes.Data.Twofa.Questions {
		var twoFaQuestion models.TwoFaQuestions
		twoFaQuestion.Question = v.Question
		twoFaQuestion.QuestionID = v.QuestionID
		twoFaQuestions = append(twoFaQuestions, twoFaQuestion)
	}

	twoFaDetails.Questions = twoFaQuestions
	twoFaDetails.TwofaToken = tlLoginRes.Data.Twofa.TwofaToken
	twoFaDetails.Type = tlLoginRes.Data.Twofa.Type

	loginRes.Twofa = twoFaDetails
	loginRes.TwofaEnabled = tlLoginRes.Data.TwofaEnabled

	loggerconfig.Info("LoginV2 response =", helpers.LogStructAsJSON(loginRes), " uccId:", loginReq.ID, " requestId:", reqH.RequestId)
	apiRes.Data = loginRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}

func (obj LoginObj) ValidateTwofaV2(validateTwoFaReq models.ValidateTwofaV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	url := obj.tradeLabURL + TWOFAV2URL

	var tlValidateTwofaReq TradeLabValidateTwofaV2Req

	tlValidateTwofaReq.LoginID = validateTwoFaReq.LoginID

	for i := 0; i < len(validateTwoFaReq.Twofa); i++ {
		var questions TradeLabTwoFaQuestion
		questions.Answer = validateTwoFaReq.Twofa[i].Answer
		questions.QuestionID = validateTwoFaReq.Twofa[i].QuestionID
		tlValidateTwofaReq.Twofa = append(tlValidateTwofaReq.Twofa, questions)
	}

	tlValidateTwofaReq.TwofaToken = validateTwoFaReq.TwofaToken
	tlValidateTwofaReq.Type = validateTwoFaReq.Type
	tlValidateTwofaReq.DeviceType = reqH.DeviceType

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlValidateTwofaReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFuncV2(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, "")
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ValidateTwofaV2", duration, validateTwoFaReq.LoginID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " validate two fa v2 call api error =", err, " uccId:", validateTwoFaReq.LoginID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("validate two fa v2 res error =", tlErrorRes.Message, "status code:", res.StatusCode, " uccId:", validateTwoFaReq.LoginID, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlValidateTwofaV2Res := TradeLabValidateTwofaV2Res{}
	json.Unmarshal([]byte(string(body)), &tlValidateTwofaV2Res)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " validate two fa v2 tl status not ok =", tlValidateTwofaV2Res.Message, "status code:", res.StatusCode, " uccId:", validateTwoFaReq.LoginID, " requestId:", reqH.RequestId)
		apiRes.Message = tlValidateTwofaV2Res.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var validateTwofaRes models.ValidateTwofaV2Res
	validateTwofaRes.ResetPassword = tlValidateTwofaV2Res.Data.ResetPassword
	validateTwofaRes.ResetTwoFa = tlValidateTwofaV2Res.Data.ResetTwoFa
	validateTwofaRes.AuthToken = tlValidateTwofaV2Res.Data.AuthToken

	loggerconfig.Info("validate two fa v2 tl resp=", helpers.LogStructAsJSON(validateTwofaRes), " uccId:", validateTwoFaReq.LoginID, " requestId:", reqH.RequestId)

	apiRes.Data = validateTwofaRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj LoginObj) SetupTotpV2(setupTotpReq models.SetupTotpV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	url := obj.tradeLabURL + SETTOTPV2URL

	var tlSetupTotpReq TradeLabSetupTotpV2Req

	tlSetupTotpReq.ClientID = setupTotpReq.ClientID

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlSetupTotpReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "SetupTotpV2", duration, setupTotpReq.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " setup totp v2 call api error =", err, " uccId:", tlSetupTotpReq.ClientID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("setup totp v2 res error =", tlErrorRes.Message, "status code:", res.StatusCode, " uccId:", tlSetupTotpReq.ClientID, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlSetupTotpRes := TradeLabSetupTotpV2Res{}
	json.Unmarshal([]byte(string(body)), &tlSetupTotpRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " setup totp v2 tl status not ok =", tlSetupTotpRes.Message, " uccId:", tlSetupTotpReq.ClientID, " requestId:", reqH.RequestId)
		apiRes.Message = tlSetupTotpRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var setupTotpRes models.SetupTotpV2Res
	setupTotpRes.ClientID = setupTotpReq.ClientID
	setupTotpRes.Token = tlSetupTotpRes.Data

	loggerconfig.Info("setup totp v2 tl resp=", helpers.LogStructAsJSON(setupTotpRes), "status code:", res.StatusCode, " uccId:", setupTotpRes.ClientID, " requestId:", reqH.RequestId)

	apiRes.Data = setupTotpRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj LoginObj) ChooseTwofaV2(chooseTwofaReq models.ChooseTwofaV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	url := obj.tradeLabURL + CHOOSETWOFAV2URL

	var tlChooseTwofaV2Req TradeLabChooseTwofaV2Req

	tlChooseTwofaV2Req.LoginID = chooseTwofaReq.LoginID
	tlChooseTwofaV2Req.TwofaType = chooseTwofaReq.TwofaType
	tlChooseTwofaV2Req.Totp = chooseTwofaReq.Totp

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlChooseTwofaV2Req)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ChooseTwofaV2", duration, chooseTwofaReq.LoginID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " choose totp v2 call api error =", err, " uccId:", tlChooseTwofaV2Req.LoginID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("choose totp v2 res error =", tlErrorRes.Message, "status code:", res.StatusCode, " uccId:", tlChooseTwofaV2Req.LoginID, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlChooseTwofaRes := TradeLabChooseTwofaV2Res{}
	json.Unmarshal([]byte(string(body)), &tlChooseTwofaRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " choose totp v2 tl status not ok =", tlChooseTwofaRes.Message, "status code:", res.StatusCode, " uccId:", tlChooseTwofaV2Req.LoginID, " requestId:", reqH.RequestId)
		apiRes.Message = tlChooseTwofaRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("choose totp v2 tl resp=success", " uccId:", tlChooseTwofaV2Req.LoginID, " requestId:", reqH.RequestId)

	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj LoginObj) ForgetTotpV2(forgetTotpReq models.ForgetTotpV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	url := obj.tradeLabURL + FORGETTOTPV2URL

	var tlForgetTotpV2Req TradeLabForgetTotpV2Req

	tlForgetTotpV2Req.LoginID = forgetTotpReq.LoginID
	tlForgetTotpV2Req.Pan = forgetTotpReq.Pan

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlForgetTotpV2Req)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ForgetTotpV2", duration, forgetTotpReq.LoginID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " forget Totp v2 call api error =", err, " uccId:", tlForgetTotpV2Req.LoginID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("forget Totp v2 res error =", tlErrorRes.Message, "status code:", res.StatusCode, " uccId:", tlForgetTotpV2Req.LoginID, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlForgetTotpRes := TradeLabForgetTotpV2Res{}
	json.Unmarshal([]byte(string(body)), &tlForgetTotpRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " forget Totp v2 tl status not ok =", tlForgetTotpRes.Message, "status code:", res.StatusCode, " uccId:", tlForgetTotpV2Req.LoginID, " requestId:", reqH.RequestId)
		apiRes.Message = tlForgetTotpRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("forget Totp v2 tl resp=success", " uccId:", forgetTotpReq.LoginID, " requestId:", reqH.RequestId)

	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj LoginObj) ValidateLoginOtpV2(validateLoginOtpReq models.ValidateLoginOtpV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	url := obj.tradeLabURL + VALIDATELOGINOTP

	var tlValidateLoginOtpV2Req TradeLabValidateLoginOtpV2Req

	tlValidateLoginOtpV2Req.ReferenceToken = validateLoginOtpReq.ReferenceToken
	tlValidateLoginOtpV2Req.Otp = validateLoginOtpReq.Otp

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlValidateLoginOtpV2Req)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFuncV2(http.MethodPut, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "ValidateLoginOtpV2", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " validate login otp v2 call api error =", err, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("validate login otp v2 res error =", tlErrorRes.Message, "status code:", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlValidateLoginOtpRes := TradeLabValidateLoginOtpV2Res{}
	json.Unmarshal([]byte(string(body)), &tlValidateLoginOtpRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " validate login otp v2 tl status not ok =", tlValidateLoginOtpRes.Message, "status code:", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlValidateLoginOtpRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	// added go routine to speed up login process
	go func() {
		if reqH.ClientId != "" && reqH.FCMToken != "" {
			//Store FCMToken at time of Login
			clientId, err := db.GetPgObj().StoreFCMToken(reqH.ClientId, reqH.FCMToken, reqH.RequestId)
			if err != nil {
				loggerconfig.Error("Alert Severity:P1-Critical, ValidateLoginOtpV2: Error (Login) while storing FCMToken in db: ", err.Error(), "requestId: ", reqH.RequestId)
			}

			if constants.KafkaEnable {

				redisCliObj := cache.GetRedisClientObj()

				err = redisCliObj.SAdd(constants.ClientMembers+strings.ToUpper(clientId), reqH.FCMToken)
				if err != nil {
					logrus.Error("Alert Severity:P1-Critical, ValidateLoginOtpV2: Error (Login) adding FCM token in Redis:", err, "requestId: ", reqH.RequestId)
				}

				// Push client ID into list with key as subscription_id
				subscriptionKey := constants.SubscriptionKey
				err = redisCliObj.LPush(subscriptionKey, strings.ToUpper(clientId))
				if err != nil {
					logrus.Error("Alert Severity:P1-Critical, ValidateLoginOtpV2 Error (Login) pushing client ID in Redis list:", err, "requestId: ", reqH.RequestId)
				}

				logrus.Info("ValidateLoginOtpV2 FCM Token details stored successfully in redis for client ID: ", reqH.ClientId, "requestId: ", reqH.RequestId)

			}

			logrus.Info("ValidateLoginOtpV2 FCM Token details stored successfully in db for client ID: ", reqH.ClientId, "requestId: ", reqH.RequestId)

		}
	}()

	var clientDetails models.MongoClientsDetails
	err = dbops.MongoDaoRepo.FindOne(constants.CLIENTDETAILS, bson.M{"email": strings.ToLower(reqH.ClientId)}, &clientDetails)
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("Alert Severity:P0-Critical, ValidateLoginOtpV2 Unable to fetch the client-details Data for emailID = ", reqH.ClientId, " error :", err, " requestid=", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	if err != nil && err.Error() == constants.MongoNoDocError {
		loggerconfig.Error("ValidateLoginOtpV2, EmailId Does not exist in client-details collection. emailID :", reqH.ClientId, " error:", err, " requestId:", reqH.RequestId)
		apiRes.Status = false
		apiRes.Message = constants.ErrorCodeMap[constants.InvalidEmailId]
		apiRes.ErrorCode = constants.InvalidEmailId
		return http.StatusBadRequest, apiRes
	}

	var validateLoginOtp models.LoginV2Response
	validateLoginOtp.KycUserId = clientDetails.KycUserId
	validateLoginOtp.Alert = tlValidateLoginOtpRes.Data.Alert
	validateLoginOtp.AuthToken = tlValidateLoginOtpRes.Data.AuthToken
	validateLoginOtp.CheckPan = tlValidateLoginOtpRes.Data.CheckPan
	validateLoginOtp.LoginID = tlValidateLoginOtpRes.Data.LoginID
	validateLoginOtp.Name = tlValidateLoginOtpRes.Data.Name
	validateLoginOtp.ReferenceToken = tlValidateLoginOtpRes.Data.ReferenceToken
	validateLoginOtp.ResetPassword = tlValidateLoginOtpRes.Data.ResetPassword
	validateLoginOtp.ResetTwoFa = tlValidateLoginOtpRes.Data.ResetTwoFa

	var twoFaDetails models.TwoFaDetails

	var twoFaQuestions []models.TwoFaQuestions

	for _, v := range tlValidateLoginOtpRes.Data.Twofa.Questions {
		var twoFaQuestion models.TwoFaQuestions
		twoFaQuestion.Question = v.Question
		twoFaQuestion.QuestionID = v.QuestionID
		twoFaQuestions = append(twoFaQuestions, twoFaQuestion)
	}

	twoFaDetails.Questions = twoFaQuestions
	twoFaDetails.TwofaToken = tlValidateLoginOtpRes.Data.Twofa.TwofaToken
	twoFaDetails.Type = tlValidateLoginOtpRes.Data.Twofa.Type

	validateLoginOtp.Twofa = twoFaDetails
	validateLoginOtp.TwofaEnabled = tlValidateLoginOtpRes.Data.TwofaEnabled

	loggerconfig.Info("validate login otp v2 tl resp=success", " requestId:", reqH.RequestId)

	apiRes.Data = validateLoginOtp
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj LoginObj) SetupBiometricV2(setupBiometricReq models.SetupBiometricV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	url := obj.tradeLabURL + BIOMETRIC + "?client_id=" + url.QueryEscape(setupBiometricReq.ClientID) + "&fingerprint=" + setupBiometricReq.FingerPrint

	//make payload
	payload := new(bytes.Buffer)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFuncV2(http.MethodPut, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "SetupBiometricV2", duration, setupBiometricReq.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " SetupBiometric v2 call api error =", err, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("SetupBiometric v2 res error =", tlErrorRes.Message, "status code:", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlSetupBiometricV2Res := TradeLabSetupBiometricV2Res{}
	json.Unmarshal([]byte(string(body)), &tlSetupBiometricV2Res)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " SetupBiometric v2 tl status not ok =", tlSetupBiometricV2Res.Message, "status code:", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlSetupBiometricV2Res.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("SetupBiometric v2 tl tradelab Resp:", helpers.LogStructAsJSON(tlSetupBiometricV2Res), " requestId:", reqH.RequestId)

	apiRes.Data = tlSetupBiometricV2Res.Data
	apiRes.Message = tlSetupBiometricV2Res.Message
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj LoginObj) DisableBiometricV2(disableBiometricReq models.DisableBiometricV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	url := obj.tradeLabURL + BIOMETRIC + "?client_id=" + url.QueryEscape(disableBiometricReq.ClientID) + "&fingerprint=" + disableBiometricReq.FingerPrint

	//make payload
	payload := new(bytes.Buffer)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFuncV2(http.MethodDelete, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "DisableBiometricV2", duration, disableBiometricReq.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " DisableBiometric v2 call api error =", err, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("DisableBiometric v2 res error =", tlErrorRes.Message, "status code:", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlDisableBiometricV2Res := TradeLabDisableBiometricV2Res{}
	json.Unmarshal([]byte(string(body)), &tlDisableBiometricV2Res)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " DisableBiometric v2 tl status not ok =", tlDisableBiometricV2Res.Message, "status code:", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlDisableBiometricV2Res.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("DisableBiometric v2 tl tradelab resp: ", helpers.LogStructAsJSON(tlDisableBiometricV2Res), " requestId:", reqH.RequestId)

	apiRes.Data = tlDisableBiometricV2Res.Data
	apiRes.Message = tlDisableBiometricV2Res.Message
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj LoginObj) GuestUserStatus(guestUserStatusReq models.GuestUserStatusReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	guestUserStatusReq.Email = strings.ToLower(guestUserStatusReq.Email)

	var guestUserStatusRes models.GuestUserStatusRes

	var tradingUserInfoData models.TradingUserInfoData
	var err error
	err = dbops.MongoRepo.FindOne(constants.TRADINGUSERS, bson.M{"emailno": guestUserStatusReq.Email}, &tradingUserInfoData)
	userFound := false
	if err != nil && err.Error() == "mongo: no documents in result" {
		userFound = false
	} else {
		userFound = true
	}

	if userFound {
		guestUserStatusRes.UserId = tradingUserInfoData.Userid
		guestUserStatusRes.GuestStatus = constants.TRADINGGUESTUSER
		loggerconfig.Info("GuestUserStatus Successful, response:", helpers.LogStructAsJSON(guestUserStatusReq), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		apiRes.Data = guestUserStatusRes
		apiRes.Message = "SUCCESS"
		apiRes.Status = true
		return http.StatusOK, apiRes
	}

	var dbUser models.MongoSignup
	err = dbops.MongoRepo.FindOne(constants.CLIENTCOLLECTION, bson.M{"emailid": guestUserStatusReq.Email}, &dbUser)

	if err != nil && err.Error() == "mongo: no documents in result" {
		userFound = false
	} else {
		userFound = true
	}

	if userFound {
		guestUserStatusRes.UserId = dbUser.UserId
		guestUserStatusRes.GuestStatus = constants.COMPLETEKYCUSER
		loggerconfig.Info("GuestUserStatus Successful, response:", helpers.LogStructAsJSON(guestUserStatusReq), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		apiRes.Data = guestUserStatusRes
		apiRes.Message = "SUCCESS"
		apiRes.Status = true
		return http.StatusOK, apiRes
	}

	guestUserStatusRes.GuestStatus = constants.UNIDENTIFIEDUSER
	loggerconfig.Info("GuestUserStatus Successful, response:", helpers.LogStructAsJSON(guestUserStatusReq), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)

	apiRes.Data = guestUserStatusRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj LoginObj) QRWebLogin(req models.LoginWithQRReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	err := dbops.RedisRepo.Set(req.WebsocketID, reqH.Authorization, 3*time.Minute)
	if err != nil {
		loggerconfig.Info("QRWebLogin failed to set redis! error:", err, " clientID:", reqH.ClientId, " requestID:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj LoginObj) UnblockUser(reqParams models.UnblockUserReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	url := obj.tradeLabURL + UNBLOCKUSER

	var apiRes apihelpers.APIRes

	var tlUnblockUserReq TradeLabUnblockUserReq
	tlUnblockUserReq.LoginID = reqParams.LoginID
	tlUnblockUserReq.Pan = reqParams.Pan

	// make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlUnblockUserReq)

	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	if err != nil {
		loggerconfig.Error("UnblockUser call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("UnblockUser tl res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", reqH.ClientId, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		apiRes.Data = tlErrorRes.Data
		return res.StatusCode, apiRes
	}

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("UnblockUser tl status not ok =", res.StatusCode, " uccId:", reqH.ClientId, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		apiRes.Data = tlErrorRes.Data
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("UnblockUser tl success uccId:", reqH.ClientId, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)

	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	apiRes.Data = tlErrorRes.Data

	return http.StatusOK, apiRes
}

func (obj LoginObj) ForgetPasswordV2(forgetPasswordV2Req models.ForgotPasswordV2Request, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	var forgetPasswordReq models.ForgotPasswordRequest
	var clientDetails models.MongoClientsDetails
	var apiRes apihelpers.APIRes

	forgetPasswordReq.Pan = forgetPasswordV2Req.Pan

	emailID := strings.ToLower(forgetPasswordV2Req.EmailID)

	if forgetPasswordV2Req.ClientID == "" && forgetPasswordV2Req.EmailID == "" {
		apiRes.Status = false
		apiRes.Message = constants.ErrorCodeMap[constants.EmptyCredentials]
		apiRes.ErrorCode = constants.EmptyCredentials
		return http.StatusBadRequest, apiRes
	} else if forgetPasswordV2Req.ClientID != "" {
		forgetPasswordReq.LoginID = forgetPasswordV2Req.ClientID
	} else {
		err := dbops.MongoDaoRepo.FindOne(constants.CLIENTDETAILS, bson.M{"email": emailID}, &clientDetails)
		if err != nil && err.Error() != "mongo: no documents in result" {
			loggerconfig.Error("Alert Severity:P0-Critical, ForgetPasswordV2 Unable to fetch the client-details Data for emailID = ", emailID, " requestid=", reqH.RequestId, " error :", err)
			return apihelpers.SendInternalServerError()
		}

		if err != nil && err.Error() == "mongo: no documents in result" {
			loggerconfig.Error("ForgetPasswordV2, EmailId Does not exist in client-details collection emailID :", emailID, " error:", err, " requestId:", reqH.RequestId)
			apiRes.Status = false
			apiRes.Message = constants.ErrorCodeMap[constants.InvalidEmailId]
			apiRes.ErrorCode = constants.InvalidEmailId
			return http.StatusOK, apiRes
		}

		forgetPasswordReq.LoginID = clientDetails.ClientID
	}
	return obj.ForgetPassword(forgetPasswordReq, reqH)
}

func (obj LoginObj) UnblockUserV2(unblockUserV2Req models.UnblockUserV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	var unblockUserReq models.UnblockUserReq
	var clientDetails models.MongoClientsDetails
	var apiRes apihelpers.APIRes

	unblockUserReq.Pan = unblockUserV2Req.Pan

	emailID := strings.ToLower(unblockUserV2Req.EmailID)

	if unblockUserV2Req.ClientID == "" && unblockUserV2Req.EmailID == "" {
		apiRes.Status = false
		apiRes.Message = constants.ErrorCodeMap[constants.EmptyCredentials]
		apiRes.ErrorCode = constants.EmptyCredentials
		return http.StatusBadRequest, apiRes
	} else if unblockUserV2Req.ClientID != "" {
		unblockUserReq.LoginID = unblockUserV2Req.ClientID
	} else {
		err := dbops.MongoDaoRepo.FindOne(constants.CLIENTDETAILS, bson.M{"email": emailID}, &clientDetails)
		if err != nil && err.Error() != "mongo: no documents in result" {
			loggerconfig.Error("Alert Severity:P0-Critical, ForgetPasswordV2 Unable to fetch the client-details Data for emailID = ", emailID, " requestid=", reqH.RequestId, " error :", err)
			return apihelpers.SendInternalServerError()
		}

		if err != nil && err.Error() == "mongo: no documents in result" {
			loggerconfig.Error("ForgetPasswordV2, EmailId Does not exist in client-details collection emailID :", emailID, " error:", err, " requestId:", reqH.RequestId)
			apiRes.Status = false
			apiRes.Message = constants.ErrorCodeMap[constants.InvalidEmailId]
			apiRes.ErrorCode = constants.InvalidEmailId
			return http.StatusBadRequest, apiRes
		}

		unblockUserReq.LoginID = clientDetails.ClientID
	}
	return obj.UnblockUser(unblockUserReq, reqH)
}

// Oauth2
func (obj LoginObj) CreateApp(reqParams models.CreateAppReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradelaboAuthURL + CREATEAPP

	// fill TL Req
	var tlCreateAppReq TradeLabCreateAppReq
	tlCreateAppReq.AppName = reqParams.AppName
	tlCreateAppReq.RedirectUris = reqParams.RedirectUris
	tlCreateAppReq.Scope = reqParams.Scope
	tlCreateAppReq.GrantTypes = reqParams.GrantTypes
	tlCreateAppReq.Owner = reqParams.Owner

	// make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlCreateAppReq)

	// call api
	var apiRes apihelpers.APIRes
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " Create App call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("CreateApp Error in reading response body error:", err, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	tlErrorAppRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorAppRes)
	if err == nil && tlErrorAppRes.Status == TLERROR {
		loggerconfig.Error("Alert Severity:P1-Critical, platform: ", reqH.Platform, "status: ", TLERROR, " CreateApp error tl api 2 NOT OK error=", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorAppRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorAppRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlCreateAppRes := TradeLabCreateAppRes{}
	json.Unmarshal([]byte(string(body)), &tlCreateAppRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P1-Critical, platform:", reqH.Platform, " CreateApp error tl api 2 NOT OK error=", tlCreateAppRes.Message, "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlCreateAppRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	// fill up controller response
	var createapp models.CreateAppRes
	var appdetails models.AppDetails

	// Set the fields for createapp and appdetails
	createapp.AppOwner = tlCreateAppRes.Data.AppOwner
	appdetails.AppID = tlCreateAppRes.Data.AppID
	appdetails.AppName = tlCreateAppRes.Data.AppName
	appdetails.AppSecret = tlCreateAppRes.Data.AppSecret
	appdetails.GrantTypes = tlCreateAppRes.Data.GrantTypes
	appdetails.RedirectUris = tlCreateAppRes.Data.RedirectUris
	appdetails.Scope = tlCreateAppRes.Data.Scope

	// adding uuid as unique identifier for state field
	id := uuid.New().String()
	appdetails.State = id

	// Initialize the Apps slice with the new appdetails
	createapp.Apps = []models.AppDetails{appdetails}

	// Prepare the API response
	apiRes.Data = appdetails
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	// Filter to find the existing document by AppOwner
	filter := bson.M{"appowner": createapp.AppOwner}

	var existingApp models.CreateAppRes
	err2 := dbops.MongoRepo.FindOne(constants.APPDETAILS, filter, &existingApp)
	if err2 == mongo.ErrNoDocuments {
		// If no document exists, insert the new createapp document
		err = dbops.MongoRepo.InsertOne(constants.APPDETAILS, createapp)
		if err != nil {
			loggerconfig.Error("Alert Severity:P1-Critical, platform:", reqH.Platform, "CreateApp (controller), Error inserting app-details in mongo: ", err, " requestId: ", reqH.RequestId, "ClientID: ", reqH.ClientId, " deviceId: ", reqH.DeviceId)
			return apihelpers.SendInternalServerError()
		}
	} else if err2 != nil {
		// Handle other potential errors in finding the document
		loggerconfig.Error("Alert Severity:P1-Critical, platform:", reqH.Platform, "CreateApp (controller), Error finding app-details in mongo: ", err, " requestId: ", reqH.RequestId, "ClientID: ", reqH.ClientId, " deviceId: ", reqH.DeviceId)
		return apihelpers.SendInternalServerError()
	} else {
		// Append the new app details to the existing Apps slice
		existingApp.Apps = append(existingApp.Apps, appdetails)

		// Update the existing document in MongoDB
		update := bson.M{
			"$set": bson.M{
				"apps": existingApp.Apps,
			},
		}

		err = dbops.MongoRepo.UpdateOne(constants.APPDETAILS, filter, update)
		if err != nil {
			loggerconfig.Error("CreateApp (controller), Error updating app-details in mongo: ", err, " requestId: ", reqH.RequestId, "ClientID: ", reqH.ClientId, " deviceId: ", reqH.DeviceId)
			return apihelpers.SendInternalServerError()
		}
	}
	loggerconfig.Info("CreateApp success response:", helpers.LogStructAsJSON(appdetails), "RequestId: ", reqH.RequestId)
	return http.StatusOK, apiRes
}

func (obj LoginObj) FetchApps(clientId string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	filter := bson.M{"appowner": clientId}

	var allApps models.CreateAppRes
	err := dbops.MongoRepo.FindOne(constants.APPDETAILS, filter, &allApps)
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("Alert Severity:P1-Critical, HandleAuthCode (controller), Error finding the app: ", err, " requestId: ", reqH.RequestId, "ClientID: ", reqH.ClientId, " deviceId: ", reqH.DeviceId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("FetchApps success response : ", helpers.LogStructAsJSON(allApps), " requestId:", reqH.RequestId)

	for i := len(allApps.Apps)/2 - 1; i >= 0; i-- {
		opp := len(allApps.Apps) - 1 - i
		allApps.Apps[i], allApps.Apps[opp] = allApps.Apps[opp], allApps.Apps[i]
	}

	var apiRes apihelpers.APIRes
	apiRes.Data = allApps
	apiRes.Status = true
	return http.StatusOK, apiRes

}

func (obj LoginObj) DeleteApp(appId string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradelaboAuthURL + CREATEAPP + "/" + appId
	payload := new(bytes.Buffer)
	var apiRes apihelpers.APIRes
	res, err := apihelpers.CallAPIFunc(http.MethodDelete, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-Critical, platform:", reqH.Platform, " OAUTH2 call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("FetchApp Error in reading response body error:", err, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	tlErrorDeleteRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorDeleteRes)
	if err == nil && tlErrorDeleteRes.Status == TLERROR {
		loggerconfig.Error("Alert Severity:P1-Critical, DeleteApp Error ", TLERROR, " statuscode: ", res.StatusCode, "requestId: ", reqH.RequestId)
		apiRes.Message = tlErrorDeleteRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorDeleteRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlDeleteRes := TradelabDeleteResponse{}
	json.Unmarshal([]byte(string(body)), &tlDeleteRes)

	filter := bson.M{"appowner": reqH.ClientId}

	// Define the pull operation to remove the app with the given appId from the apps array
	update := bson.M{
		"$pull": bson.M{
			"apps": bson.M{"appId": appId},
		},
	}

	// Perform the update operation
	err = dbops.MongoRepo.UpdateOne(constants.APPDETAILS, filter, update)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-Critical, DeleteApp (controller), Error updating deletion in mongo: ", err, " requestId: ", reqH.RequestId, "ClientID: ", reqH.ClientId, " deviceId: ", reqH.DeviceId)
		return apihelpers.SendInternalServerError()
	}
	loggerconfig.Info("DeleteApp success response", helpers.LogStructAsJSON(tlDeleteRes), "requestId", reqH.RequestId)

	apiRes.Data = tlDeleteRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj LoginObj) GenerateAccessToken(appState string, authCode string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	authurl := obj.tradeLabURL + GENERATEACCESSTOKEN

	filter := bson.M{"apps.state": appState}
	var result models.CreateAppRes
	err := dbops.MongoRepo.FindOne(constants.APPDETAILS, filter, &result)
	if err != nil {
		if err.Error() == constants.MongoNoDocError {
			loggerconfig.Error("GenerateAccessToken, No document found with given appState: ", appState, " requestId: ", reqH.RequestId)
			return apihelpers.SendErrorResponse(false, constants.AppDoesNotExists, http.StatusBadRequest)
		}
		loggerconfig.Error("Alert Severity:P1-High, Generate Access Token, Error finding the app: ", err, " requestId: ", reqH.RequestId, "ClientID: ", reqH.ClientId, " deviceId: ", reqH.DeviceId)
		return apihelpers.SendInternalServerError()
	}

	var finalApp models.AppDetails
	for _, app := range result.Apps {
		if app.State == appState {
			finalApp = app
			break
		}
	}

	// Basic Authorization header setup
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", finalApp.AppID, finalApp.AppSecret)))
	reqH.Authorization = fmt.Sprintf("Basic %s", auth)

	// Set up form data for OAuth2 request
	data := url.Values{}
	if len(finalApp.GrantTypes) <= 0 || len(finalApp.RedirectUris) <= 0 || authCode == "" {
		loggerconfig.Error("GenerateAccessToken, GrantTypes, RedirectUris or AuthCode is empty, grantType: ", len(finalApp.GrantTypes), " redirectUris: ", len(finalApp.RedirectUris), " authCode: ", authCode, "appState: ", appState, " requestId: ", reqH.RequestId, "ClientID: ", reqH.ClientId, " deviceId: ", reqH.DeviceId)
		return apihelpers.SendErrorResponse(false, constants.AppDoesNotExists, http.StatusBadRequest)
	}
	data.Set("grant_type", finalApp.GrantTypes[0])
	data.Set("redirect_uri", finalApp.RedirectUris[0])
	data.Set("code", authCode)
	payload := bytes.NewBufferString(data.Encode())

	// Make the API call using the helper function
	res, err := apihelpers.CallAPIFuncOauth(http.MethodPost, authurl, payload, reqH.Authorization)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-Critical, platform:", reqH.Platform, "GenerateAccessToken call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	// Parse response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("GenerateAccessToken (controller), unable to read body, err ", err, " statuscode: ", res.StatusCode, " requestId: ", reqH.RequestId, "ClientID: ", reqH.ClientId, " deviceId: ", reqH.DeviceId)
		return apihelpers.SendInternalServerError()
	}
	var apiRes apihelpers.APIRes

	tlErrorTokenRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorTokenRes)
	if err == nil && tlErrorTokenRes.Status == TLERROR {
		loggerconfig.Error("OauthAPI error tl api 1= ", tlErrorTokenRes.Status, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorTokenRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorTokenRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	var errorResponse map[string]interface{}
	err = json.Unmarshal(body, &errorResponse)
	if err != nil {
		loggerconfig.Error("GenerateAccessToken (controller), unable to unmarshal, err: ", err, " requestId: ", reqH.RequestId, "ClientID: ", reqH.ClientId, " deviceId: ", reqH.DeviceId)
		return apihelpers.SendInternalServerError()
	}

	tlOauthAccessTokenRes := TradelabOauthAccessTokenRes{}
	json.Unmarshal([]byte(string(body)), &tlOauthAccessTokenRes)

	// update expiry time of token
	tokenExpiry := tlOauthAccessTokenRes.ExpiresIn
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		loggerconfig.Error("GenerateAccessToken (controller), unable to load location for IST time, err: ", err, " requestId: ", reqH.RequestId, "ClientID: ", reqH.ClientId, " deviceId: ", reqH.DeviceId)
		return apihelpers.SendInternalServerError()
	}
	currentTime := helpers.GetCurrentTimeInIST().In(loc)
	expiryTime := currentTime.Add(time.Duration(tokenExpiry) * time.Second)
	expiryFormatted := expiryTime.Format("2006-01-02 15:04:05")

	updateFilter := bson.M{
		"apps.state": appState,
	}

	update := bson.M{
		"$set": bson.M{
			"apps.$.accessToken": tlOauthAccessTokenRes.AccessToken, // Update the accesstoken for the matched app
			"apps.$.expiryTime":  expiryFormatted,
		},
	}
	err = dbops.MongoRepo.UpdateOne(constants.APPDETAILS, updateFilter, update)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-Critical, platform:", reqH.Platform, "GenerateAccessToken call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("GenerateAccessToken success response:", helpers.LogStructAsJSON(tlOauthAccessTokenRes), "RequestId: ", reqH.RequestId)
	apiRes.Data = tlOauthAccessTokenRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}

func (obj LoginObj) HandleAuthCode(authCode string, appState string, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	filter := bson.M{"apps.state": appState}
	var result models.CreateAppRes
	err := dbops.MongoRepo.FindOne(constants.APPDETAILS, filter, &result)
	if err != nil {
		if err.Error() == constants.MongoNoDocError {
			loggerconfig.Error("HandleAuthCode, No document found with given appState: ", appState, " requestId: ", reqH.RequestId)
			return apihelpers.SendErrorResponse(false, constants.AppDoesNotExists, http.StatusBadRequest)
		}
		loggerconfig.Error("Alert Severity:P0-Critical, HandleAuthCode, Error finding the app: ", err, " requestId: ", reqH.RequestId, "ClientID: ", reqH.ClientId, " deviceId: ", reqH.DeviceId)
		return apihelpers.SendInternalServerError()
	}

	updateFilter := bson.M{
		"apps.state": appState,
	}

	update := bson.M{
		"$set": bson.M{
			"apps.$.authCode": authCode, // Update the authcode for the matched app
		},
	}

	err = dbops.MongoRepo.UpdateOne(constants.APPDETAILS, updateFilter, update)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-Critical, HandleAuthCode (controller), Error updating auth-code in mongo: ", err, " requestId: ", reqH.RequestId, "ClientID: ", reqH.ClientId, " deviceId: ", reqH.DeviceId)
		return apihelpers.SendInternalServerError()
	}

	obj.GenerateAccessToken(appState, authCode, reqH)

	loggerconfig.Info("HandleAuthCode success ", "RequestId: ", reqH.RequestId)

	var apiRes apihelpers.APIRes
	apiRes.Message = "SUCCESS"
	return http.StatusOK, apiRes
}

func (obj LoginObj) GetAccessToken(reqParams models.GetAccessTokenReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	filter := bson.M{"apps.state": reqParams.AppState}
	var result models.CreateAppRes
	err := dbops.MongoRepo.FindOne(constants.APPDETAILS, filter, &result)
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("Alert Severity:P1-Critical, GetAccessToken (controller), Error finding the app: ", err, " requestId: ", reqH.RequestId, "ClientID: ", reqH.ClientId, " deviceId: ", reqH.DeviceId)
		return apihelpers.SendInternalServerError()
	}

	var finalApp models.AppDetails
	for _, app := range result.Apps {
		if app.State == reqParams.AppState {
			finalApp = app
			break
		}
	}
	var apiRes apihelpers.APIRes

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, GetAccessToken, unable to load location for IST time, err: ", err, " requestId: ", reqH.RequestId, "ClientID: ", reqH.ClientId, " deviceId: ", reqH.DeviceId)
		return apihelpers.SendInternalServerError()
	}

	currentDate := time.Now().In(loc).Format(constants.YYYYMMDD)

	if finalApp.ExpiryTime == "" {
		apiRes.Message = constants.TokenFailure
		apiRes.Status = false
		loggerconfig.Info("GetAccessToken success Message:", apiRes.Message, " StatusCode : ", apiRes.Status, " requestId:", reqH.RequestId)
		return http.StatusOK, apiRes
	}

	expiryTime, err := time.Parse(constants.TIMEFORMAT, finalApp.ExpiryTime)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, GetAccessToken, unable to parse expiry time, err: ", err, " requestId: ", reqH.RequestId, "ClientID: ", reqH.ClientId, " deviceId: ", reqH.DeviceId)
		return apihelpers.SendInternalServerError()
	}

	expiryDate := expiryTime.Format(constants.YYYYMMDD)

	if finalApp.AccessToken == "" || currentDate >= expiryDate {
		apiRes.Message = constants.TokenFailure
		apiRes.Status = false
		loggerconfig.Info("GetAccessToken success Message:", apiRes.Message, " StatusCode : ", apiRes.Status, " requestId:", reqH.RequestId)
		return http.StatusOK, apiRes
	} else {
		apiRes.Data = finalApp.AccessToken
		apiRes.Message = constants.TokenSuccess
		apiRes.Status = true
		loggerconfig.Info("GetAccessToken success Message:", apiRes.Message, " StatusCode : ", apiRes.Status, " requestId:", reqH.RequestId)
		return http.StatusOK, apiRes
	}
}

func (obj LoginObj) GetAccessTokenV2(reqParams models.GetAccessTokenV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	filter := bson.M{"apps.state": reqParams.AppState}
	var result models.CreateAppRes
	err := dbops.MongoRepo.FindOne(constants.APPDETAILS, filter, &result)
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("Alert Severity:P1-High, GetAccessTokenV2, Error finding the app: ", err, " requestId: ", reqH.RequestId, "ClientID: ", reqH.ClientId, " deviceId: ", reqH.DeviceId)
		return apihelpers.SendInternalServerError()
	}

	var finalApp models.AppDetails
	for _, app := range result.Apps {
		if app.State == reqParams.AppState {
			finalApp = app
			break
		}
	}
	var apiRes apihelpers.APIRes

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, GetAccessTokenV2, unable to load location for IST time, err: ", err, " requestId: ", reqH.RequestId, "ClientID: ", reqH.ClientId, " deviceId: ", reqH.DeviceId)
		return apihelpers.SendInternalServerError()
	}

	currentDate := time.Now().In(loc).Format(constants.YYYYMMDD)

	if finalApp.ExpiryTime == "" {
		apiRes.Message = constants.TokenFailure
		apiRes.Status = false
		loggerconfig.Info("GetAccessToken success Message:", apiRes.Message, " StatusCode : ", apiRes.Status, " requestId:", reqH.RequestId)
		return http.StatusOK, apiRes
	}

	expiryTime, err := time.Parse(constants.TIMEFORMAT, finalApp.ExpiryTime)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, GetAccessTokenV2, unable to parse expiry time, err: ", err, " requestId: ", reqH.RequestId, "ClientID: ", reqH.ClientId, " deviceId: ", reqH.DeviceId)
		return apihelpers.SendInternalServerError()
	}

	expiryDate := expiryTime.Format(constants.YYYYMMDD)

	if finalApp.AccessToken == "" || currentDate >= expiryDate || finalApp.AccessToken == reqParams.AccessToken {
		apiRes.Message = constants.TokenFailure
		apiRes.Status = false
		loggerconfig.Info("GetAccessTokenV2 success Message:", apiRes.Message, " StatusCode : ", apiRes.Status, " requestId:", reqH.RequestId)
		return http.StatusOK, apiRes
	} else {
		apiRes.Data = finalApp.AccessToken
		apiRes.Message = constants.TokenSuccess
		apiRes.Status = true
		loggerconfig.Info("GetAccessTokenV2 success Message:", apiRes.Message, " StatusCode : ", apiRes.Status, " requestId:", reqH.RequestId)
		return http.StatusOK, apiRes
	}
}

func (obj LoginObj) LoginByEmailOtp(req models.LoginByEmailOtpReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	loggerconfig.Info("LoginByEmailOtp (Service):", "Request Packet:", helpers.LogStructAsJSON(req), " requestid=", reqH.RequestId, "platform:", reqH.Platform, "clientVersion:", reqH.ClientVersion)

	url := obj.tradeLabURL + LOGINV2URL

	var tlLoginByEmailOtpReq TradelabLoginByEmailOtpReq
	tlLoginByEmailOtpReq.ChannelID = req.Email

	// make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlLoginByEmailOtpReq)

	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFuncV2(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, "")
	duration := time.Since(start)

	helpers.RecordAPILatency(url, "LoginByEmailOtp", duration, req.Email, reqH.RequestId)

	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " LoginByEmailOtp call api error =", err, " EmailId:", req.Email, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " LoginByEmailOtp: Failed to read response body for emailId ", req.Email, " requestid=", reqH.RequestId, "clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	tlErrorLoginRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorLoginRes)
	if err == nil && tlErrorLoginRes.Status == TLERROR {
		loggerconfig.Error("platform:", reqH.Platform, " LoginByEmailOtp error tl api error=", err, "status code:", res.StatusCode, " EmailId:", req.Email, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorLoginRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorLoginRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlLoginByEmailOtpRes := TradeLabLoginV2Response{}
	json.Unmarshal([]byte(string(body)), &tlLoginByEmailOtpRes)

	if res.StatusCode != http.StatusOK {
		if tlLoginByEmailOtpRes.Message != AccountFrozen {
			loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " LoginByEmailOtp error tl api 2 NOT OK error=", tlLoginByEmailOtpRes.Message, "status code:", res.StatusCode, " emailId:", req.Email, " requestId:", reqH.RequestId)
		} else {
			loggerconfig.Error("Alert Severity:P2-Mid, platform:", reqH.Platform, " LoginByEmailOtp error tl api 2 NOT OK error=", tlLoginByEmailOtpRes.Message, "status code:", res.StatusCode, " emailId:", req.Email, " requestId:", reqH.RequestId)
		}

		apiRes.Message = tlLoginByEmailOtpRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	var loginByEmailOtpRes models.LoginByEmailOtpRes
	loginByEmailOtpRes.Alert = tlLoginByEmailOtpRes.Data.Alert
	loginByEmailOtpRes.AuthToken = tlLoginByEmailOtpRes.Data.AuthToken
	loginByEmailOtpRes.CheckPan = tlLoginByEmailOtpRes.Data.CheckPan
	loginByEmailOtpRes.LoginID = tlLoginByEmailOtpRes.Data.LoginID
	loginByEmailOtpRes.Name = tlLoginByEmailOtpRes.Data.Name
	loginByEmailOtpRes.ReferenceToken = tlLoginByEmailOtpRes.Data.ReferenceToken
	loginByEmailOtpRes.ResetPassword = tlLoginByEmailOtpRes.Data.ResetPassword
	loginByEmailOtpRes.ResetTwoFA = tlLoginByEmailOtpRes.Data.ResetTwoFa
	loginByEmailOtpRes.TwoFAEnabled = tlLoginByEmailOtpRes.Data.TwofaEnabled

	loggerconfig.Info("LoginByEmailOtp (service) Successful, response:", helpers.LogStructAsJSON(tlLoginByEmailOtpRes), "requestId:", reqH.RequestId, " platform= ", reqH.Platform, " clientVersion= ", reqH.ClientVersion)

	apiRes.Data = loginByEmailOtpRes
	apiRes.Message = tlLoginByEmailOtpRes.Message
	apiRes.Status = true
	return http.StatusOK, apiRes
}
