package userdetails

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	apihelpers "space/apiHelpers"
	"space/business/tradelab"
	"space/constants"
	"space/db"
	"space/dbops"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/ggwhite/go-masker"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserDetailsObj struct {
	tradeLabURL string
	mongodb     db.MongoDatabase
}

var maskObj *masker.Masker

func InitUserDetailsProvider(mongodb db.MongoDatabase) UserDetailsObj {
	defer models.HandlePanic()

	userDetailsObj := UserDetailsObj{
		tradeLabURL: constants.TLURL,
		mongodb:     mongodb,
	}

	maskObj = masker.New()

	return userDetailsObj
}

func (obj UserDetailsObj) GetAllBankAccounts(req models.GetAllBankAccountsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	var bankAccounts models.MongoAllBankAccountDetails
	err := dbops.MongoRepo.FindOne(constants.ALLBANKACCOUNTS, bson.M{"userid": req.UserId}, &bankAccounts)
	if err != nil {
		loggerconfig.Error("GetAllBankAccounts Bank Accounts not found corresponding to userid:", req.UserId, " error:", err, " requestId:", reqH.RequestId)
		apiRes.Status = false
		apiRes.Message = constants.ErrorCodeMap[constants.DetailsDoesNotExsists]
		apiRes.ErrorCode = constants.DetailsDoesNotExsists
		return http.StatusOK, apiRes
	}

	maskedBankAccounts, err := maskObj.Struct(bankAccounts)
	if err != nil {
		loggerconfig.Error("GetAllBankAccounts Error in masking request err: ", err, " clientId: ", req.UserId, " requestid = ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("GetAllBankAccountsRes resp=", helpers.LogStructAsJSON(maskedBankAccounts), " requestId:", reqH.RequestId)
	apiRes.Status = true
	apiRes.Message = "SUCCESS"
	apiRes.Data = bankAccounts

	return http.StatusOK, apiRes
}

func (obj UserDetailsObj) GetAllBankAccountsUpdated(req models.GetAllBankAccountsUpdatedReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + tradelab.ALLBANKACCOUNTURL + "?client_id=" + url.QueryEscape(req.ClientId)

	//make payload
	payload := new(bytes.Buffer) // empty payload

	//call api
	var apiRes apihelpers.APIRes
	res, err := apihelpers.CallApi(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	defer res.Body.Close()
	if err != nil {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " GetAllBankAccountsUpdated call api error =%v", err, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	body, err := io.ReadAll(res.Body)
	tlErrorRes := tradelab.TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == tradelab.TLERROR {
		loggerconfig.Error("GetAllBankAccountsUpdated res error =%v", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientId, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlAllBankAccounts := tradelab.TradeLabAllBankAccountsUpdatedRes{}
	json.Unmarshal([]byte(string(body)), &tlAllBankAccounts)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " GetAllBankAccountsUpdated tl status not ok =%v", tlAllBankAccounts.Message, " uccId:", req.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlAllBankAccounts.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var allBankAccountDetails models.GetAllBankAccountsUpdatedRes
	allBankAccountDetails.ClientId = tlAllBankAccounts.Data.ClientId

	var allBankAccount []models.BankAccount
	for i := 0; i < len(tlAllBankAccounts.Data.BankAccounts); i++ {
		var bankAccount models.BankAccount
		bankAccount.AccountType = tlAllBankAccounts.Data.BankAccounts[i].AccountType
		bankAccount.BankAccountNumber = tlAllBankAccounts.Data.BankAccounts[i].BankAccountNumber
		bankAccount.BankBranchName = tlAllBankAccounts.Data.BankAccounts[i].BankBranchName
		bankAccount.BankID = tlAllBankAccounts.Data.BankAccounts[i].BankID
		bankAccount.BankName = tlAllBankAccounts.Data.BankAccounts[i].BankName
		bankAccount.City = tlAllBankAccounts.Data.BankAccounts[i].City
		bankAccount.Ifsc = tlAllBankAccounts.Data.BankAccounts[i].Ifsc
		bankAccount.PanNumber = tlAllBankAccounts.Data.BankAccounts[i].PanNumber
		bankAccount.State = tlAllBankAccounts.Data.BankAccounts[i].State
		allBankAccount = append(allBankAccount, bankAccount)
	}

	allBankAccountDetails.BankAccounts = allBankAccount

	maskedAllBankAccountDetails, err := maskObj.Struct(allBankAccountDetails)
	if err != nil {
		loggerconfig.Error("GetAllBankAccountsUpdated Error in masking request err: ", err, " clientId: ", req.ClientId, " requestid = ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("GetAllBankAccountsUpdated resp=", helpers.LogStructAsJSON(maskedAllBankAccountDetails), " uccId:", req.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)

	apiRes.Data = allBankAccountDetails
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}

// Constraint - Admin should daily upload trading user - an api is already there is pace-rbac to upload csv file of trading user daily
func (obj UserDetailsObj) GetUserId(req models.GetUserIdReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	var getUserIdRes models.GetUserIdRes
	var tradingUserInfoData models.TradingUserInfoData
	var email string
	if strings.EqualFold(req.IdType, constants.EMAILID) {
		// if email is idType then we can query directly to get userId but to check correct email is send we are checking in trading user collection
		err := dbops.MongoRepo.FindOne(constants.TRADINGUSERS, bson.M{"emailno": req.Id}, &tradingUserInfoData)
		if err != nil {
			if err.Error() == constants.MongoNoDocError {
				loggerconfig.Error("GetUserId Failed to locate docs in Mongo corresponding to emailId:", req.Id, " requestid=", reqH.RequestId)
				apiRes.Status = false
				apiRes.Message = constants.ErrorCodeMap[constants.InvalidEmailId]
				apiRes.ErrorCode = constants.InvalidEmailId
				return http.StatusBadRequest, apiRes
			}
			loggerconfig.Error("GetUserId Mongo connection Error = ", err, " emailId:", req.Id, " requestid=", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		email = tradingUserInfoData.Emailno
	} else if strings.EqualFold(req.IdType, constants.CLIENTID) {
		err := dbops.MongoRepo.FindOne(constants.TRADINGUSERS, bson.M{"userid": req.Id}, &tradingUserInfoData)
		if err != nil {
			if err.Error() == constants.MongoNoDocError {
				loggerconfig.Error("GetUserId Failed to locate docs in Mongo corresponding to clientId:", req.Id, " requestid=", reqH.RequestId)
				apiRes.Status = false
				apiRes.Message = constants.ErrorCodeMap[constants.InvalidClient]
				apiRes.ErrorCode = constants.InvalidClient
				return http.StatusBadRequest, apiRes
			}
			loggerconfig.Error("GetUserId Mongo connection Error = ", err, " clientId:", req.Id, " requestid=", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		fmt.Println("emailId: ", email)
		email = tradingUserInfoData.Emailno
	}

	var dbUser models.MongoSignup
	err := dbops.MongoRepo.FindOne(constants.CLIENTCOLLECTION, bson.M{"emailid": email}, &dbUser)
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("GetUserId Mongo connection Error = ", err, " id:", req.Id, " requestid=", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	getUserIdRes.Id = req.Id
	getUserIdRes.UserId = dbUser.UserId

	apiRes.Data = getUserIdRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj UserDetailsObj) UserNotifications(req models.UserNotificationsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes
	var mongoNotificationStore []models.MongoNotificationStore

	// Set default values for page and pageSize if they are not provided
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	skip := (req.Page - 1) * req.PageSize

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(req.PageSize))
	findOptions.SetSort(bson.D{{"storedat", -1}})

	filter := bson.M{"clientid": strings.ToUpper(req.ClientId)}

	// Perform the query
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := db.CallFindAllMongo(obj.mongodb, constants.CLIENTNOTIFICATION, filter, findOptions) // dbops.MongoRepo.Find(constants.CLIENTNOTIFICATION, filter, findOptions)
	if err != nil {
		loggerconfig.Error("UserNotifications Mongo query Error = ", err, " clientid: ", req.ClientId, " requestid=", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	err = readValueMongoCurs(req, ctx, cursor, &mongoNotificationStore, reqH)
	if err != nil {
		return apihelpers.SendInternalServerError()
	}

	apiRes.Data = mongoNotificationStore
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

var readValueMongoCurs = func(req models.UserNotificationsReq, ctx context.Context, cursor *mongo.Cursor, mongoNotificationStore *[]models.MongoNotificationStore, reqH models.ReqHeader) error {
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var notification models.MongoNotificationStore
		err := cursor.Decode(&notification)
		if err != nil {
			loggerconfig.Error("ReadValueMongoCurs UserNotifications Mongo decode Error = ", err, " clientid: ", req.ClientId, " requestid=", reqH.RequestId)
			return errors.New(constants.ErrorCodeMap[constants.InternalServerError]) // apihelpers.SendInternalServerError()
		}
		*mongoNotificationStore = append(*mongoNotificationStore, notification)
	}

	// Check for errors during iteration
	if err := cursor.Err(); err != nil {
		loggerconfig.Error("ReadValueMongoCurs UserNotifications Mongo cursor Error = ", err, " clientid: ", req.ClientId, " requestid=", reqH.RequestId)
		return errors.New(constants.ErrorCodeMap[constants.InternalServerError])
	}

	return nil

}

func (obj UserDetailsObj) GetClientStatus(emailId string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var clientDetails models.MongoClientsDetails
	pattern := "^" + regexp.QuoteMeta(emailId) + "$"
	err := dbops.MongoDaoRepo.FindOne(constants.CLIENTDETAILS, bson.M{
		"email": bson.M{
			"$regex":   pattern,
			"$options": "i",
		},
	}, &clientDetails)
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("Alert Severity:P0-Critical, GetClientStatus Unable to fetch the client-details Data for emailID = ", emailId, " error :", err, " requestid=", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var apiRes apihelpers.APIRes
	var res models.GetClientStatusRes

	if err != nil && err.Error() == constants.MongoNoDocError {
		loggerconfig.Info("GetClientStatus, EmailId Does not exist in client-details collection emailID :", emailId, " requestId:", reqH.RequestId)
		res.UserStatus = constants.NonUserValue // Not a trading user
	} else {
		if strings.EqualFold(clientDetails.ClientType, constants.GUESTUSERTYPE) {
			res.UserStatus = constants.GuestUserValue // Guest user
			res.KycUserId  = clientDetails.KycUserId
		} else {
			res.UserStatus = constants.TradingUserValue // Trading user
			res.KycUserId  = clientDetails.KycUserId
		}
	}

	loggerconfig.Info("GetClientStatus resp=", helpers.LogStructAsJSON(res), " emailId:", emailId, " requestId:", reqH.RequestId)

	apiRes.Data = res
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}
