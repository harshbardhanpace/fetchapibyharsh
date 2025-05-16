package upipreference

import (
	"net/http"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/db"
	"space/dbops"
	"space/loggerconfig"
	"space/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UpiPreferenceObj struct {
	mongodb db.MongoDatabase
}

func InitUpiPreferenceProvider(mongodb db.MongoDatabase) UpiPreferenceObj {
	defer models.HandlePanic()
	upiPreferenceObj := UpiPreferenceObj{mongodb: mongodb}
	return upiPreferenceObj
}

func (obj UpiPreferenceObj) SetUpiPreference(req models.SetUpiPreferenceReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	//call api
	var apiRes apihelpers.APIRes

	upiPreference, err := CallFetchUpiPreferenceMongo(req.ClientId, obj)
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("SetUpiPreference", req, " mongo err:", err, " clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.InvalidRequest, http.StatusBadRequest)
	}

	if upiPreference.ClientId == "" {
		upiPreference.ClientId = req.ClientId
	}

	var duplicateUpi bool
	for i := 0; i < len(upiPreference.UpiIds); i++ {
		if req.UpiId == upiPreference.UpiIds[i] {
			duplicateUpi = true
			break
		}
	}

	if duplicateUpi {
		loggerconfig.Error("SetUpiPreference user provided duplicate upiid clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.DuplicateUpi, http.StatusBadRequest)
	}

	// prepending the new upiId
	upiPreference.UpiIds = append([]string{req.UpiId}, upiPreference.UpiIds...)

	err = CallUpdateUpiPreferenceMongo(req.ClientId, upiPreference, obj)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " SetUpiPreference Mongo Upsert failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj UpiPreferenceObj) FetchUpiPreference(req models.FetchUpiPreferenceReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	//call api
	var apiRes apihelpers.APIRes

	upiPreference, err := CallFetchUpiPreferenceMongo(req.ClientId, obj)
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("FetchUpiPreference", req, " mongo err:", err, " clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.InvalidRequest, http.StatusBadRequest)
	}

	var fetchUpiPreferenceRes models.FetchUpiPreferenceRes
	fetchUpiPreferenceRes.ClientId = req.ClientId
	fetchUpiPreferenceRes.UpiIds = upiPreference.UpiIds

	apiRes.Data = fetchUpiPreferenceRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj UpiPreferenceObj) DeleteUpiPreference(req models.DeleteUpiPreferenceReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	//call api
	var apiRes apihelpers.APIRes

	upiPreference, err := CallFetchUpiPreferenceMongo(req.ClientId, obj)
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("DeleteUpiPreference", req, " mongo err:", err, " clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.InvalidRequest, http.StatusBadRequest)
	}

	mapUpiId := make(map[string]bool)
	for i := 0; i < len(req.UpiIds); i++ {
		_, present := mapUpiId[req.UpiIds[i]]
		if present {
			loggerconfig.Error("DeleteUpiPreference user provided duplicate upiid clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
			return apihelpers.SendErrorResponse(false, constants.DuplicateUpi, http.StatusBadRequest)
		}
		mapUpiId[req.UpiIds[i]] = true
	}

	var updatedUpiPreference []string
	var upiRemovedCount int
	for i := 0; i < len(upiPreference.UpiIds); i++ {
		_, present := mapUpiId[upiPreference.UpiIds[i]]
		if present {
			upiRemovedCount++
			continue
		}
		updatedUpiPreference = append(updatedUpiPreference, upiPreference.UpiIds[i])
	}

	if upiRemovedCount < len(req.UpiIds) {
		loggerconfig.Error("DeleteUpiPreference Some Upi don't exist, requestId: ", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.UpiDontExist, http.StatusBadRequest)
	}

	upiPreference.UpiIds = updatedUpiPreference

	err = CallUpdateUpiPreferenceMongo(req.ClientId, upiPreference, obj)
	if err != nil {
		loggerconfig.Error("DeleteUpiPreference Mongo Upsert failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

var CallFetchUpiPreferenceMongo = func(clientId string, obj UpiPreferenceObj) (models.UpiPreference, error) {
	return obj.FetchUpiPreferenceMongo(clientId)
}

func (obj UpiPreferenceObj) FetchUpiPreferenceMongo(clientId string) (models.UpiPreference, error) {
	var upiPreference models.UpiPreference
	var err error
	err = dbops.MongoRepo.FindOne(constants.UPIPREFERENCE, bson.M{"clientId": clientId}, &upiPreference)

	return upiPreference, err
}

var CallUpdateUpiPreferenceMongo = func(clientId string, upiPreference models.UpiPreference, obj UpiPreferenceObj) error {
	return obj.UpdateUpiPreferenceMongo(clientId, upiPreference)
}

func (obj UpiPreferenceObj) UpdateUpiPreferenceMongo(clientId string, upiPreference models.UpiPreference) error {
	filter := bson.D{{"clientId", clientId}}
	update := bson.D{{"$set", upiPreference}}
	opts := options.Update().SetUpsert(true)
	err := dbops.MongoRepo.UpdateOne(constants.UPIPREFERENCE, filter, update, opts)

	return err
}
