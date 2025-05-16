package tradelab

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/db"
	"space/dbops"
	"space/helpers"
	"space/helpers/cache"
	"space/loggerconfig"
	"space/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProfileObj struct {
	tradeLabURL string
	redisCli    cache.RedisCache
	mongodb     db.MongoDatabase
}

func InitProfile(mongodb db.MongoDatabase, redisCli cache.RedisCache) ProfileObj {
	defer models.HandlePanic()

	profileObj := ProfileObj{
		tradeLabURL: constants.TLURL,
		redisCli:    redisCli,
		mongodb:     mongodb,
	}

	return profileObj
}

func (obj ProfileObj) GetProfile(req models.ProfileRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + PROFILEURL + "?client_id=" + url.QueryEscape(req.ClientID)

	//make payload
	payload := new(bytes.Buffer) // empty payload

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "GetProfile", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " GetProfileResponse call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("GetProfileResponse res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlGetProfileResponse := TradeLabProfileResponse{}
	json.Unmarshal([]byte(string(body)), &tlGetProfileResponse)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " GetProfileResponse tl status not ok =", tlGetProfileResponse.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlGetProfileResponse.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var profileRes models.ProfileResponse
	// profileRes.Data = tlGetProfileResponse.Data
	profileRes.Data.AccountType = tlGetProfileResponse.Data.AccountType
	profileRes.Data.BackofficeLink = tlGetProfileResponse.Data.BackofficeLink
	profileRes.Data.BankCity = tlGetProfileResponse.Data.BankCity
	profileRes.Data.BankAccountNumber = tlGetProfileResponse.Data.BankAccountNumber
	profileRes.Data.BankBranchName = tlGetProfileResponse.Data.BankBranchName
	profileRes.Data.BankName = tlGetProfileResponse.Data.BankName
	profileRes.Data.BankState = tlGetProfileResponse.Data.BankState
	profileRes.Data.BasketEnabled = tlGetProfileResponse.Data.BasketEnabled
	profileRes.Data.Branch = tlGetProfileResponse.Data.Branch
	profileRes.Data.BrokerID = tlGetProfileResponse.Data.BrokerID
	profileRes.Data.City = tlGetProfileResponse.Data.City
	profileRes.Data.ClientID = tlGetProfileResponse.Data.ClientID
	profileRes.Data.Depository = tlGetProfileResponse.Data.Depository
	profileRes.Data.Dob = tlGetProfileResponse.Data.Dob
	profileRes.Data.EmailID = tlGetProfileResponse.Data.EmailID
	profileRes.Data.ExchangeNnf = tlGetProfileResponse.Data.ExchangeNnf
	profileRes.Data.ExchangesSubscribed = tlGetProfileResponse.Data.ExchangesSubscribed
	profileRes.Data.IfscCode = tlGetProfileResponse.Data.IfscCode
	profileRes.Data.LastPasswordChangeDate = tlGetProfileResponse.Data.LastPasswordChangeDate
	profileRes.Data.Name = tlGetProfileResponse.Data.Name
	profileRes.Data.OfficeAddr = tlGetProfileResponse.Data.OfficeAddr
	profileRes.Data.PanNumber = tlGetProfileResponse.Data.PanNumber
	profileRes.Data.PermanentAddr = tlGetProfileResponse.Data.PermanentAddr
	profileRes.Data.PhoneNumber = tlGetProfileResponse.Data.PhoneNumber
	profileRes.Data.PoaEnabled = tlGetProfileResponse.Data.PoaEnabled
	profileRes.Data.PoaStatus = tlGetProfileResponse.Data.PoaStatus
	profileRes.Data.ProductsEnabled = tlGetProfileResponse.Data.ProductsEnabled
	profileRes.Data.ProfileURL = tlGetProfileResponse.Data.ProfileURL
	var profileResRole models.RoleDetails
	profileResRole.ID = tlGetProfileResponse.Data.Role.ID
	profileResRole.Name = tlGetProfileResponse.Data.Role.Name
	profileRes.Data.Role = profileResRole
	// profileRes.Data.Role = tlGetProfileResponse.Data.Role
	profileRes.Data.Sex = tlGetProfileResponse.Data.Sex
	profileRes.Data.State = tlGetProfileResponse.Data.State
	profileRes.Data.Status = tlGetProfileResponse.Data.Status
	profileRes.Data.TwofaEnabled = tlGetProfileResponse.Data.TwofaEnabled
	profileRes.Data.UserType = tlGetProfileResponse.Data.UserType
	profileRes.Data.BoID = tlGetProfileResponse.Data.BoID
	profileRes.Data.DpID = tlGetProfileResponse.Data.DpID

	segments, _ := SegmentDetailsUpdateInProfile(req.ClientID, tlGetProfileResponse, reqH)

	profileRes.Data.SegmentDetails = segments

	maskedProfileRes, err := maskObj.Struct(profileRes)
	if err != nil {
		loggerconfig.Error("GetProfileResponse Error in masking request err: ", err, " clientId: ", req.ClientID, " requestid = ", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("GetProfileResponse tl resp=", helpers.LogStructAsJSON(maskedProfileRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = profileRes.Data
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}

func (obj ProfileObj) SendAFOtp(req models.SendAFOtpReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	var reqProfile models.ProfileRequest
	reqProfile.ClientID = req.ClientID

	statusProfile, resProfile := obj.GetProfile(reqProfile, reqH)
	if statusProfile != http.StatusOK {
		loggerconfig.Error("platform:", reqH.Platform, "SendAFOtp error status: ", statusProfile, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return statusProfile, resProfile
	}
	profileRes, ok := resProfile.Data.(models.ProfileResponseData)
	if !ok {
		loggerconfig.Error("platform:", reqH.Platform, "SendAFOtp interface parsing error: ", ok, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	email := profileRes.EmailID
	phoneNo := profileRes.PhoneNumber

	var otp = ""
	emailOtp, _ := dbops.RedisRepo.Get(constants.AccountFreeze + email)

	loggerconfig.Info("platform:", reqH.Platform, " SendAFOtp OTP email in redis server old = ", emailOtp, " uccId=", req.ClientID, " requestid=", reqH.RequestId, "clientVersion:", reqH.ClientVersion)
	if emailOtp == "" {
		otp = helpers.GenerateOTP(constants.OtpLen)
		err := dbops.RedisRepo.Set(constants.AccountFreeze+email, otp, constants.OtpSetTimeRedis*time.Minute)
		if err != nil {
			loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " SendAFOtp Error storing OTP value in Redis Server=", err, " userId=", req.ClientID, " requestid=", reqH.RequestId, "clientVersion:", reqH.ClientVersion)
			return apihelpers.SendInternalServerError()
		}
		loggerconfig.Info("SendAFOtp New OTP email generated now in redis server = ", otp, " uccId=", req.ClientID, " requestid=", reqH.RequestId, "platform:", reqH.Platform, "clientVersion:", reqH.ClientVersion)
	} else {
		otp = emailOtp
	}
	loggerconfig.Info("SendAFOtp OTP final email= ", otp, " uccId=", req.ClientID, " requestid=", reqH.RequestId, "platform:", reqH.Platform, "clientVersion:", reqH.ClientVersion)

	go SendEmailOtpNotifications(otp, email, profileRes.Name, true, reqH)

	go helpers.SendSms(otp, phoneNo)

	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj ProfileObj) VerifyAFOtp(req models.VerifyAFOtpReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	var reqProfile models.ProfileRequest
	reqProfile.ClientID = req.ClientID

	statusProfile, resProfile := obj.GetProfile(reqProfile, reqH)
	if statusProfile != http.StatusOK {
		loggerconfig.Error("platform:", reqH.Platform, "VerifyAFOtp error status: ", statusProfile, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return statusProfile, resProfile
	}
	profileRes, ok := resProfile.Data.(models.ProfileResponseData)
	if !ok {
		loggerconfig.Error("platform:", reqH.Platform, "VerifyAFOtp interface parsing error: ", ok, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	email := profileRes.EmailID

	emailOtp, _ := dbops.RedisRepo.Get(constants.AccountFreeze + email)
	if emailOtp != req.Otp {
		loggerconfig.Error("platform:", reqH.Platform, "VerifyAFOtp Failed to to set email otp in Redis corresponding to EmailId", profileRes.EmailID, " userId", req.ClientID, " requestid=", reqH.RequestId, "platform:", reqH.Platform, "clientVersion:", reqH.ClientVersion)
		return apihelpers.SendErrorResponse(false, constants.InvalidEmailOtp, http.StatusBadRequest)
	}

	var mongoAccountFreeze models.MongoAccountFreeze
	mongoAccountFreeze.ClientId = req.ClientID
	mongoAccountFreeze.OtpStatus = true
	mongoAccountFreeze.FreezeStatus = false

	filter := bson.D{{"clientId", req.ClientID}}
	update := bson.D{{"$set", mongoAccountFreeze}}
	opts := options.Update().SetUpsert(true)
	err := dbops.MongoRepo.UpdateOne(constants.AccountFreezeStatus, filter, update, opts)
	if err != nil {
		loggerconfig.Error("platform:", reqH.Platform, "VerifyAFOtp Mongo Upsert failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func SendEmailOtpNotifications(otp string, emailId string, name string, isEmail bool, reqH models.ReqHeader) {
	var sendEmailOtp models.SendEmailOtp
	sendEmailOtp.Otp = otp
	sendEmailOtp.RecipientEmail = emailId
	sendEmailOtp.RecipientName = name
	sendEmailOtp.IsEmail = isEmail
	sendEmailOtp.ReqHeader = reqH

	helpers.PublishMessage(constants.TopicExchange, constants.PKTFLKyc1, sendEmailOtp)

}

func (obj ProfileObj) AccountFreeze(req models.AccountFreezeReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	var mongoAccountFreeze models.MongoAccountFreeze
	err := dbops.MongoRepo.FindOne(constants.AccountFreezeStatus, bson.M{"clientId": req.ClientID}, &mongoAccountFreeze)
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("platform:", reqH.Platform, "AccountFreeze mongo err:", err, "clientID: ", req.ClientID, "requestId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	if (err != nil && err.Error() == constants.MongoNoDocError) || !mongoAccountFreeze.OtpStatus {
		loggerconfig.Error("platform:", reqH.Platform, "AccountFreeze can't be freezed clientID: ", req.ClientID, "requestId: ", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.AccountFreezeInvalidRequest, http.StatusBadRequest)
	}

	url := obj.tradeLabURL + ACCOUNTFREEZE
	//make payload
	payload := new(bytes.Buffer) // empty payload

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "AccountFreeze", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " AccountFreeze call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("platform:", reqH.Platform, "AccountFreeze res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlAccountFreezeRes := TradeLabAccountFreezeRes{}
	json.Unmarshal([]byte(string(body)), &tlAccountFreezeRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " AccountFreeze tl status not ok =", tlAccountFreezeRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlAccountFreezeRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	mongoAccountFreeze.FreezeStatus = true

	err = dbops.MongoRepo.UpdateOne(constants.AccountFreezeStatus, bson.M{"clientId": req.ClientID}, bson.D{{"$set", mongoAccountFreeze}}, options.Update().SetUpsert(true))
	if err != nil {
		loggerconfig.Error("platform:", reqH.Platform, "AccountFreeze Mongo Upsert failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var accountFreeze models.AccountFreezeRes
	accountFreeze.Data = tlAccountFreezeRes.Data
	accountFreeze.Message = tlAccountFreezeRes.Message

	apiRes.Data = accountFreeze
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func SegmentDetailsUpdateInProfile(clientId string, tlGetProfileResponse TradeLabProfileResponse, reqH models.ReqHeader) (models.ExchangesActive, error) {
	//all fields set to "Activate"
	segments := models.ExchangesActive{
		NSE: "Activate",
		BSE: "Activate",
		NFO: "Activate",
		BFO: "Activate",
		MCX: "Activate",
	}

	exchangeMap := map[string]*string{
		"NSE": &segments.NSE,
		"BSE": &segments.BSE,
		"NFO": &segments.NFO,
		"BFO": &segments.BFO,
		"MCX": &segments.MCX,
	}

	// Update subscribed exchanges to "ACTIVATED"
	for _, exchange := range tlGetProfileResponse.Data.ExchangesSubscribed {
		if field, exists := exchangeMap[exchange]; exists {
			*field = "ACTIVATED"
		}
	}

	res, err := FetchSegmentUpdatePending(clientId, reqH)
	if err != nil {
		loggerconfig.Error("SegmentDetailsUpdateInProfile Error Fetching pending segment updates, error:", err, " clientId=", clientId, " requestId=", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return segments, nil
	}

	if len(res) > 0 {
		segments.NFO = "InProcess"
		segments.BFO = "InProcess"
		segments.MCX = "InProcess"
	}

	return segments, nil
}

func FetchSegmentUpdatePending(clientId string, reqH models.ReqHeader) ([]models.MongoClientDetailsChange, error) {
	var result []models.MongoClientDetailsChange
	filter := bson.M{
		"clientId":        clientId,
		"transactionType": bson.M{"$regex": ".*SEGMENTUPDATE$"},
	}

	res, err := dbops.MongoDaoRepo.Find(constants.ACCOUNTDETAILSUPDATE, filter)
	if err != nil {
		loggerconfig.Error("Alert Severity:P2-Mid, platform:", reqH.Platform, " FetchSegmentUpdatePending can't find request in Mongo = ", err, " clientId=", clientId, " requestId=", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return result, err
	}
	defer res.Close(context.Background())

	for res.Next(context.Background()) {
		var accountChangeDetails models.MongoClientDetailsChange
		if err := res.Decode(&accountChangeDetails); err != nil {
			loggerconfig.Error("FetchSegmentUpdatePending Decode Error = ", err, " requestId=", reqH.RequestId, " platform:", reqH.Platform, " clientVersion:", reqH.ClientVersion)
			return result, err
		}
		if accountChangeDetails.Esign.EsignStatus != "" {
			result = append(result, accountChangeDetails)
		}
	}

	return result, nil
}
