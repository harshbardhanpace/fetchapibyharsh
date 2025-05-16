package finvu

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/db"
	"space/dbops"
	"space/helpers"
	"space/helpers/cache"
	"space/loggerconfig"
	"space/models"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FinvuObj struct {
	finvuURL       string
	FinvuRid       string
	FinvuTs        string
	FinvuChannelId string
	mongodb        db.MongoDatabase
	redisCli       cache.RedisCache
}

func InitFinvuProvider(mongodb db.MongoDatabase, redisCli cache.RedisCache) FinvuObj {
	defer models.HandlePanic()

	finvuObj := FinvuObj{
		finvuURL:       constants.FinvuBaseUrl,
		FinvuRid:       constants.FinvuRidDetails,
		FinvuTs:        constants.FinvuTsDetails,
		FinvuChannelId: constants.FinvuChannelIdDetails,
		mongodb:        mongodb,
		redisCli:       redisCli,
	}

	return finvuObj
}

func (obj FinvuObj) FinvuConsentRequestPlus(req models.CreateConsentRequestPlusReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	consentsRequestPlusReq := models.ConsentsRequestPlusReq{}

	// Assign values to its fields
	consentsRequestPlusReq.Header.Rid = obj.FinvuRid
	consentsRequestPlusReq.Header.Ts = obj.FinvuTs
	consentsRequestPlusReq.Header.ChannelID = obj.FinvuChannelId
	consentsRequestPlusReq.Body.CustID = req.CustomerId
	consentsRequestPlusReq.Body.ConsentDescription = constants.FinvuConsentDescription
	consentsRequestPlusReq.Body.TemplateName = constants.FinvuTemplateName
	consentsRequestPlusReq.Body.UserSessionID = constants.FinvuUserSessionID
	consentsRequestPlusReq.Body.RedirectURL = constants.FinvuRedirectURL
	consentsRequestPlusReq.Body.Fip = constants.FinvuFip

	consentsRequestPlusReq.Body.ConsentDetails.Customer.ID = req.CustomerId
	consentsRequestPlusReq.Body.ConsentDetails.DataConsumer.ID = constants.FinvuDataConsumerId
	consentsRequestPlusReq.Body.ConsentDetails.Purpose.Code = constants.FinvuPurposeCode
	consentsRequestPlusReq.Body.ConsentDetails.Purpose.RefURI = constants.FinvuPurposeRefURI
	consentsRequestPlusReq.Body.ConsentDetails.Purpose.Text = constants.FinvuPurposeText
	consentsRequestPlusReq.Body.ConsentDetails.Purpose.Category.Type = constants.FinvuCategoryType
	consentsRequestPlusReq.Body.ConsentDetails.ConsentMode = constants.FinvuConsentMode
	consentsRequestPlusReq.Body.ConsentDetails.ConsentTypes = constants.FinvuConsentTypes
	consentsRequestPlusReq.Body.ConsentDetails.FiTypes = constants.FinvuFiTypes
	consentsRequestPlusReq.Body.ConsentDetails.FetchType = constants.FinvuFetchType
	consentsRequestPlusReq.Body.ConsentDetails.Frequency.Value = constants.FinvuFrequencyValue
	consentsRequestPlusReq.Body.ConsentDetails.Frequency.Unit = constants.FinvuFrequencyUnit
	consentsRequestPlusReq.Body.ConsentDetails.DataLife.Value = constants.FinvuDataLifeValue
	consentsRequestPlusReq.Body.ConsentDetails.DataLife.Unit = constants.FinvuDataLifeUnit
	consentsRequestPlusReq.Body.ConsentDetails.ConsentStart = constants.FinvuConsentStart
	consentsRequestPlusReq.Body.ConsentDetails.ConsentExpiry = constants.FinvuConsentExpiry
	consentsRequestPlusReq.Body.ConsentDetails.FIDataRange.From = constants.FinvuFIDataRangeFrom
	consentsRequestPlusReq.Body.ConsentDetails.FIDataRange.To = constants.FinvuFIDataRangeTo
	consentsRequestPlusReq.Body.AaID = constants.FinvuAaID

	//call api
	var apiRes apihelpers.APIRes

	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(consentsRequestPlusReq)

	url := obj.finvuURL + "/" + constants.FinvuConsentRequestPlus
	auth, err := dbops.RedisRepo.Get(constants.FinvuAuthTokenKey)

	res, err := apihelpers.CallApiFinvu(http.MethodPost, url, payload, auth)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " FinvuConsentRequestPlus call api error", err)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("FinvuConsentRequestPlus Error in reading response body", err)
		return apihelpers.SendInternalServerError()
	}

	if string(responseData) == constants.FinvuInvalidTokenMessage {
		loggerconfig.Error("FinvuConsentRequestPlus responseData message: ", constants.FinvuInvalidTokenMessage, " clientID: ", req.ClientId, " requestId:", reqH.RequestId)
		apiRes.Status = false
		apiRes.Message = constants.ErrorCodeMap[constants.TokenExpired]
		apiRes.ErrorCode = constants.TokenExpired
		return http.StatusUnauthorized, apiRes
	}

	if string(responseData) == constants.FinvuInternalServerError {
		loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " FinvuConsentRequestPlus responseData message: ", constants.FinvuInternalServerError, " clientID: ", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	finvuErrorRes := models.FinvuErrorRes{}
	err = json.Unmarshal(responseData, &finvuErrorRes)
	if err == nil && finvuErrorRes.Errors != nil {
		loggerconfig.Error("FinvuConsentRequestPlus res error =", finvuErrorRes.Errors[0].ErrorMsg, " clientID: ", req.ClientId, " requestId:", reqH.RequestId)
		apiRes.Message = finvuErrorRes.Errors[0].ErrorMsg
		apiRes.ErrorCode = strconv.Itoa(finvuErrorRes.Errors[0].ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	consentsRequestPlusRes := models.ConsentsRequestPlusRes{}
	json.Unmarshal([]byte(string(responseData)), &consentsRequestPlusRes)

	var consentDetailsStore models.ConsentDetailsStore
	consentDetailsStore.ClientId = req.ClientId
	consentDetailsStore.ConsentHandle = consentsRequestPlusRes.Body.ConsentHandle
	consentDetailsStore.CustomerId = req.CustomerId
	consentDetailsStore.ConsentStatus = constants.FinvuConsentStatusRequested

	filter := bson.D{{"clientId", req.ClientId}}
	update := bson.D{{"$set", consentDetailsStore}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.Finvu, filter, update, opts)
	if err != nil {
		loggerconfig.Error("FinvuConsentRequestPlus Mongo Upsert failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var consentsRequestPlusResFrontRes models.ConsentsRequestPlusResFrontRes
	consentsRequestPlusResFrontRes.URL = consentsRequestPlusRes.Body.URL

	apiRes.Data = consentsRequestPlusResFrontRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj FinvuObj) FinvuGetBankStatement(req models.FinvuGetBankStatementReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	var apiRes apihelpers.APIRes

	var consentDetailsStore models.ConsentDetailsStore
	err := dbops.MongoRepo.FindOne(constants.Finvu, bson.M{"clientId": req.ClientId}, &consentDetailsStore)
	if err != nil && err.Error() == constants.MongoNoDocError {
		loggerconfig.Error("FinvuGetBankStatement", req, " mongo err:", err, " clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.InvalidRequest, http.StatusBadRequest)
	}

	checkConsentStatusUrl := obj.finvuURL + "/" + constants.FinvuConsentStatus + "/" + consentDetailsStore.ConsentHandle + "/" + consentDetailsStore.CustomerId

	auth, err := dbops.RedisRepo.Get(constants.FinvuAuthTokenKey)
	payload := new(bytes.Buffer)

	res, err := apihelpers.CallApiFinvu(http.MethodGet, checkConsentStatusUrl, payload, auth)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " FinvuGetBankStatement call api error", err)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("FinvuGetBankStatement Error in reading response body", err)
		return apihelpers.SendInternalServerError()
	}

	if string(responseData) == constants.FinvuInvalidTokenMessage {
		loggerconfig.Error("FinvuGetBankStatement responseData message: ", constants.FinvuInvalidTokenMessage, " clientID: ", req.ClientId, " requestId:", reqH.RequestId)
		apiRes.Status = false
		apiRes.Message = constants.ErrorCodeMap[constants.TokenExpired]
		apiRes.ErrorCode = constants.TokenExpired
		return http.StatusUnauthorized, apiRes
	}

	if string(responseData) == constants.FinvuInternalServerError {
		loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " FinvuGetBankStatement responseData message: ", constants.FinvuInternalServerError, " clientID: ", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	checkConsentStatusRes := models.CheckConsentStatusRes{}
	json.Unmarshal([]byte(string(responseData)), &checkConsentStatusRes)

	if checkConsentStatusRes.Body.ConsentStatus != constants.FinvuConsentStatusACCEPTED {
		apiRes.Message = consentDetailsStore.ConsentStatus
		apiRes.Status = true
		return http.StatusOK, apiRes
	}

	consentDetailsStore.ConsentStatus = checkConsentStatusRes.Body.ConsentStatus

	filter := bson.D{{"clientId", req.ClientId}}
	update := bson.D{{"$set", consentDetailsStore}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.Finvu, filter, update, opts)
	if err != nil {
		loggerconfig.Error("FinvuConsentRequestPlus Mongo Upsert failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	getConsentDetailsByIdStatus, getConsentDetailsByIdApiRes := getConsentDetailsById(obj, checkConsentStatusRes.Body.ConsentID, auth, req.ClientId, reqH)
	if getConsentDetailsByIdStatus != http.StatusOK {
		loggerconfig.Error("getConsentDetailsById BrokerCharges status != 200", getConsentDetailsByIdStatus, " clientId:", req.ClientId, " requestId:", reqH.RequestId)
		return getConsentDetailsByIdStatus, getConsentDetailsByIdApiRes
	}
	getConsentDetailsById, ok := getConsentDetailsByIdApiRes.Data.(models.GetConsentAStatusById)
	if !ok {
		loggerconfig.Error("FinvuGetBankStatement getConsentDetailsById interface parsing error", ok, " clientId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var fiRequestReq models.FiRequestReq
	fiRequestReq.Header.Rid = obj.FinvuRid
	fiRequestReq.Header.Ts = obj.FinvuTs
	fiRequestReq.Header.ChannelID = obj.FinvuChannelId
	fiRequestReq.Body.CustID = consentDetailsStore.CustomerId
	fiRequestReq.Body.ConsentHandleID = consentDetailsStore.ConsentHandle
	fiRequestReq.Body.ConsentID = checkConsentStatusRes.Body.ConsentID
	fiRequestReq.Body.DateTimeRangeFrom = getConsentDetailsById.Body.ConsentDetail.FIDataRange.From
	fiRequestReq.Body.DateTimeRangeTo = getConsentDetailsById.Body.ConsentDetail.FIDataRange.To

	fiRequestStatus, fiRequestApiRes := fiRequest(obj, fiRequestReq, auth, req.ClientId, reqH)
	if fiRequestStatus != http.StatusOK {
		loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " FinvuGetBankStatement fiRequest status != 200", fiRequestStatus, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return getConsentDetailsByIdStatus, getConsentDetailsByIdApiRes
	}
	fiRequest, ok := fiRequestApiRes.Data.(models.FiRequestRes)
	if !ok {
		loggerconfig.Error("FinvuGetBankStatement fiRequest interface parsing error", ok, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var fiStatusRes models.FiStatusRes
	for i := 0; i < 5; i++ {
		fiStatusResStatus, fiStatusResApiRes := fiStatus(obj, checkConsentStatusRes.Body.ConsentID, fiRequest.Body.SessionID, consentDetailsStore.ConsentHandle, consentDetailsStore.CustomerId, auth, req.ClientId, reqH)
		if fiStatusResStatus != http.StatusOK {
			loggerconfig.Error("FinvuGetBankStatement fiStatus status != 200", fiStatusResStatus, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
			return getConsentDetailsByIdStatus, getConsentDetailsByIdApiRes
		}
		var ok bool
		fiStatusRes, ok = fiStatusResApiRes.Data.(models.FiStatusRes)
		if !ok {
			loggerconfig.Error("FinvuGetBankStatement fiRequest interface parsing error", ok, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		if fiStatusRes.Body.FiRequestStatus == constants.FinvuFiRequestStatusREADY { // Status can be PENDING or READY
			break
		}
	}

	fiDataFetchPdfStatus, fiDataFetchPdfApiRes := fiDataFetchPdf(obj, consentDetailsStore.ConsentHandle, fiRequest.Body.SessionID, auth, req.ClientId, reqH)
	if fiDataFetchPdfStatus != http.StatusOK {
		loggerconfig.Error("fiDataFetchPdf fiStatus status != 200", fiDataFetchPdfStatus, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return fiDataFetchPdfStatus, fiDataFetchPdfApiRes
	}

	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func fiDataFetchPdf(obj FinvuObj, consentHandleId, sessionId, auth string, clientId string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	fiDataFetchUrl := obj.finvuURL + "/" + constants.FinvuFIDataFetch + "/" + consentHandleId + "/" + sessionId

	var apiRes apihelpers.APIRes

	payload := new(bytes.Buffer)

	res, err := apihelpers.CallApiFinvuBankStatementPdf(http.MethodGet, fiDataFetchUrl, payload, auth)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " fiDataFetchPdf call api error", err, " clientID: ", clientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("fiDataFetchPdf Error in reading response body", err)
		return apihelpers.SendInternalServerError()
	}

	if string(responseData) == constants.FinvuInvalidTokenMessage {
		loggerconfig.Error("fiDataFetchPdf responseData message: ", constants.FinvuInvalidTokenMessage, " clientID: ", clientId, " requestId:", reqH.RequestId)
		apiRes.Status = false
		apiRes.Message = constants.ErrorCodeMap[constants.TokenExpired]
		apiRes.ErrorCode = constants.TokenExpired
		return http.StatusUnauthorized, apiRes
	}

	if string(responseData) == constants.FinvuInternalServerError {
		loggerconfig.Error("fiDataFetchPdf responseData message: ", constants.FinvuInternalServerError, " clientID: ", clientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	finvuErrorRes := models.FinvuErrorRes{}
	err = json.Unmarshal(responseData, &finvuErrorRes)
	if err == nil && finvuErrorRes.Errors != nil {
		loggerconfig.Error("fiDataFetchPdf res error =", finvuErrorRes.Errors[0].ErrorMsg, " clientID: ", clientId, " requestId:", reqH.RequestId)
		var apiRes apihelpers.APIRes
		apiRes.Message = finvuErrorRes.Errors[0].ErrorMsg
		apiRes.ErrorCode = strconv.Itoa(finvuErrorRes.Errors[0].ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	bankStatementName := constants.BankStatement + sessionId + ".pdf"
	if err := os.WriteFile(bankStatementName, responseData, 0644); err != nil {
		loggerconfig.Error("fiDataFetchPdf Error write esign doc = ", err, " clientID: ", clientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	s3FrontPartKey := constants.BankStatement + "/" + clientId + "/" + bankStatementName
	upFile1, err := os.Open(bankStatementName)
	defer upFile1.Close()

	if upFile1 != nil {
		// Get the file info
		upFileInfo1, _ := upFile1.Stat()
		var fileSize1 int64 = upFileInfo1.Size()
		fileBuffer1 := make([]byte, fileSize1)
		upFile1.Read(fileBuffer1)

		var s3UserBankStatementLocation string
		err, s3UserBankStatementLocation = helpers.UploadToS3(s3FrontPartKey, bytes.NewReader(fileBuffer1))
		if err != nil {
			loggerconfig.Error("fiDataFetchPdf Upload Failed BankStatement err: ", err, " clientID: ", clientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		} else {
			loggerconfig.Info("fiDataFetchPdf Upload Location BankStatement = ", s3UserBankStatementLocation, " clientID: ", clientId, " requestId:", reqH.RequestId)
			bankStatementDetails := &models.MongoBankStatementDetails{
				UserId:                  clientId,
				BankStatementS3Location: s3UserBankStatementLocation,
				Verified:                false,
				Rejection:               "",
			}

			//upsert in mongo
			filter := bson.D{{"userid", clientId}}
			update := bson.D{{"$set", bankStatementDetails}}
			opts := options.Update().SetUpsert(true)
			err = dbops.MongoDaoRepo.UpdateOne(constants.BANKSTAEMENTCOLLECTION, filter, update, opts)
			if err != nil {
				loggerconfig.Error("fiDataFetchPdf UploadOtherMetadata Failed to upload Bank Statement Details to Mongo, error = ", err, " clientID: ", clientId, " requestId:", reqH.RequestId)
				return apihelpers.SendInternalServerError()
			}
		}
	} else {
		loggerconfig.Error("fiDataFetchPdf BankStatement File Read Failed! err: ", err, " clientID: ", clientId, " requestId:", reqH.RequestId)
	}
	upFile1.Close()

	os.Remove(bankStatementName)

	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}

func fiStatus(obj FinvuObj, consentId, sessionId, consentHandleId, custId, auth, clientId string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	fiStatusUrl := obj.finvuURL + "/" + constants.FinvuFIStatus + "/" + consentId + "/" + sessionId + "/" + consentHandleId + "/" + custId

	var apiRes apihelpers.APIRes

	payload := new(bytes.Buffer)

	time.Sleep(300 * time.Millisecond) // tested by dry running on different time, got to know that at 300 miliseconds status might come to ready

	res, err := apihelpers.CallApiFinvu(http.MethodGet, fiStatusUrl, payload, auth)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " fiStatus call api error", err)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("fiStatus Error in reading response body", err, " clientID: ", clientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	if string(responseData) == constants.FinvuInvalidTokenMessage {
		loggerconfig.Error("fiStatus responseData message: ", constants.FinvuInvalidTokenMessage, " clientID: ", clientId, " requestId:", reqH.RequestId)
		apiRes.Status = false
		apiRes.Message = constants.ErrorCodeMap[constants.TokenExpired]
		apiRes.ErrorCode = constants.TokenExpired
		return http.StatusUnauthorized, apiRes
	}

	if string(responseData) == constants.FinvuInternalServerError {
		loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " fiStatus responseData message: ", constants.FinvuInternalServerError, " clientID: ", clientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	finvuErrorRes := models.FinvuErrorRes{}
	err = json.Unmarshal(responseData, &finvuErrorRes)
	if err == nil && finvuErrorRes.Errors != nil {
		loggerconfig.Error("fiStatus res error =", finvuErrorRes.Errors[0].ErrorMsg, " clientID: ", clientId, " requestId:", reqH.RequestId)
		apiRes.Message = finvuErrorRes.Errors[0].ErrorMsg
		apiRes.ErrorCode = strconv.Itoa(finvuErrorRes.Errors[0].ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	fiStatusRes := models.FiStatusRes{}
	json.Unmarshal([]byte(string(responseData)), &fiStatusRes)

	apiRes.Data = fiStatusRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func fiRequest(obj FinvuObj, fiRequestReq models.FiRequestReq, auth string, clientId string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	fiRequestUrl := obj.finvuURL + "/" + constants.FinvuFIRequest

	fiRequestRes := models.FiRequestRes{}

	var apiRes apihelpers.APIRes

	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(fiRequestReq)

	res, err := apihelpers.CallApiFinvu(http.MethodPost, fiRequestUrl, payload, auth)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " fiRequest call api error", err, " clientID: ", clientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("fiRequest Error in reading response body", err, " clientID: ", clientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	if string(responseData) == constants.FinvuInvalidTokenMessage {
		loggerconfig.Error("fiStatus responseData message: ", constants.FinvuInvalidTokenMessage, " clientID: ", clientId, " requestId:", reqH.RequestId)
		apiRes.Status = false
		apiRes.Message = constants.ErrorCodeMap[constants.TokenExpired]
		apiRes.ErrorCode = constants.TokenExpired
		return http.StatusUnauthorized, apiRes
	}

	if string(responseData) == constants.FinvuInternalServerError {
		loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " fiStatus responseData message: ", constants.FinvuInternalServerError, " clientID: ", clientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	finvuErrorRes := models.FinvuErrorRes{}
	err = json.Unmarshal(responseData, &finvuErrorRes)
	if err == nil && finvuErrorRes.Errors != nil {
		loggerconfig.Error("fiStatus res error =", finvuErrorRes.Errors[0].ErrorMsg, " clientID: ", clientId, " requestId:", reqH.RequestId)
		apiRes.Message = finvuErrorRes.Errors[0].ErrorMsg
		apiRes.ErrorCode = strconv.Itoa(finvuErrorRes.Errors[0].ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	json.Unmarshal([]byte(string(responseData)), &fiRequestRes)

	apiRes.Data = fiRequestRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}

func getConsentDetailsById(obj FinvuObj, consentId string, auth string, clientId string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	getConsentAStatusByIdUrl := obj.finvuURL + "/" + constants.FinvuConsent + "/" + consentId

	var apiRes apihelpers.APIRes
	payload := new(bytes.Buffer)

	res, err := apihelpers.CallApiFinvu(http.MethodGet, getConsentAStatusByIdUrl, payload, auth)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " getConsentDetailsById call api error", err, " clientID: ", clientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("getConsentDetailsById Error in reading response body", err, " clientID: ", clientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	if string(responseData) == constants.FinvuInvalidTokenMessage {
		loggerconfig.Error("getConsentDetailsById responseData message: ", constants.FinvuInvalidTokenMessage, " clientID: ", clientId, " requestId:", reqH.RequestId)
		apiRes.Status = false
		apiRes.Message = constants.ErrorCodeMap[constants.TokenExpired]
		apiRes.ErrorCode = constants.TokenExpired
		return http.StatusUnauthorized, apiRes
	}

	if string(responseData) == constants.FinvuInternalServerError {
		loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " getConsentDetailsById responseData message: ", constants.FinvuInternalServerError, " clientID: ", clientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	finvuErrorRes := models.FinvuErrorRes{}
	err = json.Unmarshal(responseData, &finvuErrorRes)
	if err == nil && finvuErrorRes.Errors != nil {
		loggerconfig.Error("getConsentDetailsById res error =", finvuErrorRes.Errors[0].ErrorMsg, " clientID: ", clientId, " requestId:", reqH.RequestId)
		apiRes.Message = finvuErrorRes.Errors[0].ErrorMsg
		apiRes.ErrorCode = strconv.Itoa(finvuErrorRes.Errors[0].ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	getConsentAStatusById := models.GetConsentAStatusById{}
	json.Unmarshal([]byte(string(responseData)), &getConsentAStatusById)

	apiRes.Data = getConsentAStatusById
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}
