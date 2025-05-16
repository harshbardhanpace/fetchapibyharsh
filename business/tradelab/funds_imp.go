package tradelab

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/db"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/ggwhite/go-masker"
	"github.com/google/uuid"
)

type FundsObj struct {
	tradeLabURL string
	Db          db.Database
}

var objFunds FundsObj
var maskObj *masker.Masker

func InitFetchFunds(dbInstance db.Database) FundsObj {
	defer models.HandlePanic()

	fundsObj := FundsObj{
		tradeLabURL: constants.TLURL,
		Db:          dbInstance,
	}

	objFunds = fundsObj

	maskObj = masker.New()

	return fundsObj
}

func (obj FundsObj) FetchFunds(req models.FetchFundsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return FetchFundsInternal(req, reqH)
}

var FetchFundsInternal = func(req models.FetchFundsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := objFunds.tradeLabURL + FETCHFUNDSURL + "?client_id=" + url.QueryEscape(req.ClientID) + "&type=" + req.Type

	//make payload
	payload := new(bytes.Buffer) // empty payload

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallApi(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "FetchFunds", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " FetchFundsResponse call api error =%v", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("FetchFundsResponse res error =%v", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlFetchFundsResponse := TradeLabFetchFundsResponse{}
	json.Unmarshal([]byte(string(body)), &tlFetchFundsResponse)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " FetchFundsResponse tl status not ok =%v", tlFetchFundsResponse.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlFetchFundsResponse.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}
	//fill up controller response
	var fetchFundsRes models.FetchFundsResponse
	// fetchFundsRes.Data = tlFetchFundsResponse.Data
	fetchFundsRes.ClientID = tlFetchFundsResponse.Data.ClientID
	fetchFundsRes.Headers = tlFetchFundsResponse.Data.Headers
	responseValues := make([]models.FetchFundsResponseValues, 0)
	for i := 0; i < len(tlFetchFundsResponse.Data.Values); i++ {
		var fetchFundsResponseValues models.FetchFundsResponseValues
		fetchFundsResponseValues.Num0 = tlFetchFundsResponse.Data.Values[i][0]
		fetchFundsResponseValues.Num1 = tlFetchFundsResponse.Data.Values[i][1]
		responseValues = append(responseValues, fetchFundsResponseValues)
	}
	fetchFundsRes.Values = responseValues

	loggerconfig.Info("FetchFundsResponse tl resp=%v", helpers.LogStructAsJSON(fetchFundsRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = fetchFundsRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}

func (obj FundsObj) CancelPayout(req models.CancelPayoutReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + CANCELPAYOUTFUNDSURL

	var tlCancelPayoutReq TradeLabCancelPayoutReq
	tlCancelPayoutReq.Transactions = req.Transactions
	tlCancelPayoutReq.UserID = req.ClientID
	tlCancelPayoutReq.Status = req.Status

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlCancelPayoutReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallApi(http.MethodPut, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "CancelPayout", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " CancelPayoutResponse call api error =%v", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("CancelPayoutResponse res error =%v", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlCancelPayoutRes := TradeLabCancelPayoutRes{}
	json.Unmarshal([]byte(string(body)), &tlCancelPayoutRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " CancelPayoutResponse tl status not ok =%v", tlCancelPayoutRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlCancelPayoutRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("CancelPayoutResponse tl resp=%v", helpers.LogStructAsJSON(tlCancelPayoutRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}

func (obj FundsObj) Payout(req models.AtomPayoutRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + Payout

	var payloadFill TradeLabAtomPayoutRequest
	payloadFill.Amount = req.Amount
	payloadFill.ClientID = req.ClientID
	payloadFill.Ifsc = req.Ifsc
	payloadFill.AccountNumber = req.AccountNumber
	payloadFill.BankName = req.BankName

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(payloadFill)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallApiTradeLab(http.MethodPost, url, payload, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "Payout", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " Payout call api error =%v", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("Payout res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlAtomPayoutResponse := TradeLabAtomPayoutResponse{}
	json.Unmarshal([]byte(string(body)), &tlAtomPayoutResponse)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " Payout tl status not ok =%v", tlAtomPayoutResponse.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlAtomPayoutResponse.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var payoutRes models.AtomPayoutResponse
	payoutRes.Data = tlAtomPayoutResponse.Data

	if strings.EqualFold(tlAtomPayoutResponse.Status, constants.Success) {
		var payoutDetails models.PayoutDetails
		floatAmount, _ := strconv.ParseFloat(payloadFill.Amount, 64)

		payoutDetails.Amount = int64(floatAmount)
		payoutDetails.ClientID = payloadFill.ClientID
		payoutDetails.Ifsc = payloadFill.Ifsc
		payoutDetails.AccountNumber = payloadFill.AccountNumber
		payoutDetails.BankName = payloadFill.BankName
		payoutDetails.DebitCredit = constants.Debit
		payoutDetails.TradelabFundsUpdated = false
		payoutDetails.BackofficeFundsUpdated = false
		payoutDetails.Remarks = constants.Payout + "|" + reqH.Platform
		payoutDetails.CreateDate = helpers.GetCurrentTimeInIST()
		payoutDetails.UpdatedAt = helpers.GetCurrentTimeInIST()
		payoutDetails.TransactionType = constants.Payout
		payoutDetails.TransactionId = uuid.New().String()
		payoutDetails.TransactionStatus = constants.Success

		obj.Db.InsertTransactionData(payoutDetails)
	}

	loggerconfig.Info("Payout tl resp=%v", helpers.LogStructAsJSON(payoutRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = payoutRes.Data
	apiRes.Message = tlAtomPayoutResponse.Message
	apiRes.Status = true
	return http.StatusOK, apiRes

}

func (obj FundsObj) ClientTransactions(req models.ClientTransactionsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + ClientTransactions + "?client_id=" + url.QueryEscape(req.ClientID)

	//make payload
	payload := new(bytes.Buffer) // empty payload

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallApi(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " ClientTransactions call api error =%v", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("ClientTransactions res error =%v", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlClientTransactionsResponse := TradeLabClientTransactionsResponse{}
	json.Unmarshal([]byte(string(body)), &tlClientTransactionsResponse)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " ClientTransactions tl status not ok =%v", tlClientTransactionsResponse.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlClientTransactionsResponse.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var clientTransactionsRes models.ClientTransactionsResponse
	clientTransactionsRes.Message = tlClientTransactionsResponse.Message
	clientTransactionsRes.Status = tlClientTransactionsResponse.Status

	individualResponseDataOuter := make([]models.ClientTransactionsResponseData, 0)
	for i := 0; i < len(tlClientTransactionsResponse.Data); i++ {
		var individualResponseData models.ClientTransactionsResponseData
		individualResponseData.AccountName = tlClientTransactionsResponse.Data[i].AccountName
		individualResponseData.Amount = tlClientTransactionsResponse.Data[i].Amount
		individualResponseData.BankName = tlClientTransactionsResponse.Data[i].BankName
		individualResponseData.BankTransactionID = tlClientTransactionsResponse.Data[i].BankTransactionID
		individualResponseData.ClientID = tlClientTransactionsResponse.Data[i].ClientID
		individualResponseData.CreatedAt = tlClientTransactionsResponse.Data[i].CreatedAt
		individualResponseData.Ifsc = tlClientTransactionsResponse.Data[i].Ifsc
		individualResponseData.MerchantTransactionID = tlClientTransactionsResponse.Data[i].MerchantTransactionID
		individualResponseData.PaymentGatewayTransactionID = tlClientTransactionsResponse.Data[i].PaymentGatewayTransactionID
		individualResponseData.PaymentGatewayUsername = tlClientTransactionsResponse.Data[i].PaymentGatewayUsername
		individualResponseData.PreviousBalance = tlClientTransactionsResponse.Data[i].PreviousBalance
		individualResponseData.Status = tlClientTransactionsResponse.Data[i].Status
		individualResponseData.TransactionID = tlClientTransactionsResponse.Data[i].TransactionID
		individualResponseData.TransactionTimestamp = tlClientTransactionsResponse.Data[i].TransactionTimestamp
		individualResponseData.TransactionType = tlClientTransactionsResponse.Data[i].TransactionType
		individualResponseData.UpdatedAt = tlClientTransactionsResponse.Data[i].UpdatedAt
		individualResponseData.UpdatedBy = tlClientTransactionsResponse.Data[i].UpdatedBy
		individualResponseData.UserID = tlClientTransactionsResponse.Data[i].UserID
		individualStatusLifeCycleOuter := make([]models.ClientTransactionsResponseDataStatusLifeCycle, 0)
		for j := 0; j < len(tlClientTransactionsResponse.Data[i].StatusLifeCycle); j++ {
			var individualStatusLifeCycle models.ClientTransactionsResponseDataStatusLifeCycle
			individualStatusLifeCycle.Status = tlClientTransactionsResponse.Data[i].StatusLifeCycle[j].Status
			individualStatusLifeCycle.UpdatedAt = tlClientTransactionsResponse.Data[i].StatusLifeCycle[j].UpdatedAt
			individualStatusLifeCycle.UpdatedBy = tlClientTransactionsResponse.Data[i].StatusLifeCycle[j].UpdatedBy
			individualStatusLifeCycleOuter = append(individualStatusLifeCycleOuter, individualStatusLifeCycle)
		}
		individualResponseData.StatusLifeCycle = individualStatusLifeCycleOuter
		individualResponseDataOuter = append(individualResponseDataOuter, individualResponseData)
	}
	clientTransactionsRes.Data = individualResponseDataOuter

	maskedClientTransactionsRes, err := maskObj.Struct(clientTransactionsRes)
	if err != nil {
		loggerconfig.Error("ClientTransactions Error in masking request err: ", err, " clientId: ", req.ClientID, " requestid = ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("ClientTransactions tl resp=", helpers.LogStructAsJSON(maskedClientTransactionsRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = clientTransactionsRes.Data
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}
