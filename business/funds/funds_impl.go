package funds

import (
	"encoding/json"
	"fmt"
	"net/http"
	apihelpers "space/apiHelpers"
	"space/business/charges"
	"space/constants"
	"space/db"
	"space/dbops"
	"space/helpers"
	"space/helpers/cache"
	"space/loggerconfig"
	"space/models"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type FundsObjV3 struct {
	Db       db.Database
	RedisCli cache.RedisCache
}

func InitFundsV3(db db.Database, redisCli cache.RedisCache) FundsObjV3 {

	fundsObjV2 := FundsObjV3{
		Db:       db,
		RedisCli: redisCli,
	}
	return fundsObjV2
}

func (obj FundsObjV3) Payout(req models.AtomPayoutRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	// Check for existing payout requests
	existingRequest, err := obj.Db.CheckExistingPayoutRequest(req.ClientID)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, PayoutV3 Error checking existing payout request: %v", err, " clientId:", req.ClientID, " requestId:", reqH.RequestId)
		return http.StatusInternalServerError, apihelpers.APIRes{Status: false, Message: fmt.Sprintf("Error checking existing payout request: %v", err)}
	}

	if existingRequest {
		loggerconfig.Info("Payout request already exists with status 'process' or 'pending' for clientId:", req.ClientID, " requestId:", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.ExistPayoutRequest, http.StatusBadRequest)
	}

	//create request packate
	var payloadFill models.PayoutRequest
	payloadFill.Amount = req.Amount
	payloadFill.ClientID = req.ClientID
	payloadFill.Ifsc = req.Ifsc
	payloadFill.AccountNumber = req.AccountNumber
	payloadFill.BankName = req.BankName

	// Call FundsPayoutInternal to get the payoutAmount
	payreq := models.FundsPayoutReq{
		ClientID: req.ClientID,
	}

	_, fundPayoutRes := charges.FundsPayoutInternal(payreq, reqH)
	fundsPayoutData, ok := fundPayoutRes.Data.(models.FundsPayoutRes)
	if !ok {
		loggerconfig.Error("Alert Severity:P0-Critical, PayoutV3 Error parsing FundsPayoutInternal response", " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		return http.StatusInternalServerError, apihelpers.APIRes{Status: false, Message: "Error parsing FundsPayoutInternal response"}
	}

	payoutAmount := fundsPayoutData.PayoutAmount

	// Check if requested amount is greater than payout amount
	floatRequestAmount, err := strconv.ParseFloat(payloadFill.Amount, 64)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, PayoutV3 Error parsing request amount: %v", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		return http.StatusInternalServerError, apihelpers.APIRes{Status: false, Message: fmt.Sprintf("Error parsing request amount: %v", err)}
	}

	if floatRequestAmount > payoutAmount {
		loggerconfig.Info("Requested amount is greater than the available payout amount for clientId:", req.ClientID, " requestId:", reqH.RequestId)
		return http.StatusBadRequest, apihelpers.APIRes{Status: false, Message: "Requested amount exceeds the available payout amount."}
	}

	//insert with pending status
	var payoutDetails models.PayoutDetails

	floatAmount, _ := strconv.ParseFloat(payloadFill.Amount, 64)
	intAmount := int64(floatAmount * 100)
	payoutDetails.Amount = intAmount
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

	//set into redis
	payinRedis, err := dbops.RedisRepo.Get(constants.CalculateWB)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, PayoutV3 Error fetching 'payinredis' value from Redis: %v", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		return http.StatusInternalServerError, apihelpers.APIRes{Status: false, Message: fmt.Sprintf("Error fetching 'payinredis' value: %v", err)}
	}

	if payinRedis == "1" {
		payoutDetails.TransactionStatus = constants.PROCESS
	} else {
		payoutAmount = 0
		payoutDetails.TransactionStatus = constants.Pending
	}

	currentTime := helpers.GetCurrentTimeInIST()
	var redisKey string
	if currentTime.Hour() < constants.PayoutTime {
		redisKey = constants.PartialPayoutKey + currentTime.Format(constants.YYYYMMDD)
	} else {
		redisKey = constants.PartialPayoutKey + currentTime.AddDate(0, 0, 1).Format(constants.YYYYMMDD)
	}

	payoutDetailsStruct := models.PayoutQueueDetails{
		ClientID:      req.ClientID,
		Amount:        payloadFill.Amount,
		Dates:         currentTime.Format(time.RFC3339),
		BankName:      payloadFill.BankName,
		TransactionID: payoutDetails.TransactionId,
		Remark:        payoutDetails.Remarks,
		PayoutAmount:  fmt.Sprintf("%f", payoutAmount),
		IFSC:          payoutDetails.Ifsc,
		AccountNumber: payoutDetails.AccountNumber,
	}

	// Convert struct to JSON string
	payoutDetailsJSON, err := json.Marshal(payoutDetailsStruct)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, PayoutV3 Error marshaling payout details struct to JSON: %v", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		return http.StatusInternalServerError, apihelpers.APIRes{Status: false, Message: fmt.Sprintf("Error marshaling payout details struct to JSON: %v", err)}
	}

	// Insert JSON string into Redis
	err = dbops.RedisRepo.HSet(redisKey, payloadFill.ClientID, string(payoutDetailsJSON))
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, PayoutV3 Error inserting payout details into Redis: %v", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		return http.StatusInternalServerError, apihelpers.APIRes{Status: false, Message: fmt.Sprintf("Error inserting payout details into Redis: %v", err)}
	}

	err = obj.Db.InsertTransactionData(payoutDetails)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, PayoutV3 Error inserting data into DB: %v", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId)

		// Remove the data from Redis
		errRemove := dbops.RedisRepo.HDelete(redisKey, payloadFill.ClientID)
		if errRemove != nil {
			loggerconfig.Error("Alert Severity:P0-Critical, PayoutV3 Error removing payout details from Redis: %v", errRemove, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		}

		return http.StatusInternalServerError, apihelpers.APIRes{Status: false, Message: fmt.Sprintf("Error inserting data into DB: %v", err)}
	}

	return http.StatusOK, apihelpers.APIRes{Status: true, Message: "Payout request processed successfully"}
}

func (obj FundsObjV3) CancelPayout(req models.CancelPayoutReqV3, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	//fetch data against transactionID
	res, err := obj.Db.GetTransactionData(req.TransactionId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " CancelPayout insert error =%v", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		return http.StatusInternalServerError, apihelpers.APIRes{Status: false, Message: fmt.Sprintf("Error while fetching transactions: %v", err)}
	}

	if res.TransactionStatus == constants.PROCEED {
		return http.StatusOK, apihelpers.APIRes{Status: true, Message: constants.ErrorCodeMap[constants.InvalidPayoutRequest]}
	}

	updates := map[string]interface{}{
		"transaction_status": constants.CANCELLED,
		"updated_at":         helpers.GetCurrentTimeInIST(),
	}

	// Simulate updating data in the database
	err = obj.Db.UpdateTransactionData(req.TransactionId, updates)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " CancelPayout update error =%v", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		return http.StatusInternalServerError, apihelpers.APIRes{Status: false, Message: fmt.Sprintf("Error cancelling payout: %v", err)}
	}

	// Remove from Redis
	currentTime := helpers.GetCurrentTimeInIST()
	var redisKey string
	if currentTime.Hour() < constants.PayoutTime {
		redisKey = constants.PartialPayoutKey + currentTime.Format(constants.YYYYMMDD)
	} else {
		redisKey = constants.PartialPayoutKey + currentTime.AddDate(0, 0, 1).Format(constants.YYYYMMDD)
	}

	err = dbops.RedisRepo.HDelete(redisKey, req.ClientID)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " CancelPayout error removing from Redis: %v", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		return http.StatusInternalServerError, apihelpers.APIRes{Status: false, Message: fmt.Sprintf("Error removing payout details from Redis: %v", err)}
	}

	return http.StatusOK, apihelpers.APIRes{Status: true, Message: "Payout request cancelled successfully"}
}
