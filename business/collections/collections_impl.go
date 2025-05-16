package collections

import (
	"context"
	"net/http"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/db"
	"space/dbops"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CollectionsObj struct {
	mongodb db.MongoDatabase
}

func InitCollections(mongodb db.MongoDatabase) CollectionsObj {
	defer models.HandlePanic()
	collectionsObj := CollectionsObj{mongodb: mongodb}
	return collectionsObj
}

func (obj CollectionsObj) CreateCollections(createCollectionsReq models.CreateCollectionsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var Collections models.MongoCollections

	err := dbops.MongoRepo.FindOne(constants.COLLECTIONS, bson.M{"collectionName": createCollectionsReq.CollectionName}, &Collections)
	var apiRes apihelpers.APIRes
	// if Collection already exists
	if err == nil && Collections.CollectionName != "" {
		loggerconfig.Error("CreateCollections Collection Exists error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.CollectionAlreadyExists, http.StatusBadRequest)
	}

	id := uuid.New().String()

	for i := 0; i < len(createCollectionsReq.CollectionMetaData); i++ {
		idStock := uuid.New().String()
		createCollectionsReq.CollectionMetaData[i].StockId = idStock
		err, location := helpers.Base64toPng("stockid", createCollectionsReq.CollectionMetaData[i].StockId, createCollectionsReq.CollectionMetaData[i].StockImage)
		if err != nil {
			loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " CreateCollections Error uploading into s3 =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		createCollectionsReq.CollectionMetaData[i].StockImage = location
	}

	mongoCollectionDetails := &models.MongoCollections{
		CollectionName:         createCollectionsReq.CollectionName,
		CollectionShortDesc:    createCollectionsReq.CollectionShortDesc,
		CollectionLongDesc:     createCollectionsReq.CollectionLongDesc,
		CollectionExchange:     createCollectionsReq.CollectionExchange,
		CollectionMetaData:     createCollectionsReq.CollectionMetaData,
		CollectionId:           id,
		CollectionBackColor:    createCollectionsReq.CollectionBackColor,
		CollectionImage:        createCollectionsReq.CollectionImage,
		CollectionIllustration: createCollectionsReq.CollectionIllustration,
		CollectionType:         createCollectionsReq.CollectionType,
	}

	filter := bson.D{{"collectionId", id}}
	update := bson.D{{"$set", mongoCollectionDetails}}
	opts := options.Update().SetUpsert(true)
	dbops.MongoRepo.UpdateOne(constants.COLLECTIONS, filter, update, opts)
	if err != nil {
		loggerconfig.Error("CreateCollections Mongo Upsert failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var resp models.CreateCollectionsResponse
	resp.CollectionId = id
	resp.CollectionName = createCollectionsReq.CollectionName
	resp.CollectionMetaData = createCollectionsReq.CollectionMetaData

	loggerconfig.Info("CreateCollections Successful, response:", helpers.LogStructAsJSON(resp), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = resp
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}

func (obj CollectionsObj) ModifyCollections(modifyCollectionsReq models.ModifyCollectionsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var Collections models.MongoCollections
	err := dbops.MongoRepo.FindOne(constants.COLLECTIONS, bson.M{"collectionName": modifyCollectionsReq.CollectionName}, &Collections)

	var apiRes apihelpers.APIRes
	// if Collection does not exists
	if err != nil && Collections.CollectionName == "" {
		loggerconfig.Error("ModifyCollections Collection Does Not Exist error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.CollectionDoesNotExists, http.StatusBadRequest)
	}

	for i := 0; i < len(modifyCollectionsReq.CollectionMetaData); i++ {
		err, location := helpers.Base64toPng("stockid", modifyCollectionsReq.CollectionMetaData[i].StockId, modifyCollectionsReq.CollectionMetaData[i].StockImage)
		if err != nil {
			loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " ModifyCollections Error uploading into s3 =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		modifyCollectionsReq.CollectionMetaData[i].StockImage = location
	}

	mongoCollectionDetails := &models.MongoCollections{
		CollectionName:         modifyCollectionsReq.CollectionName,
		CollectionShortDesc:    modifyCollectionsReq.CollectionShortDesc,
		CollectionLongDesc:     modifyCollectionsReq.CollectionLongDesc,
		CollectionExchange:     modifyCollectionsReq.CollectionExchange,
		CollectionMetaData:     modifyCollectionsReq.CollectionMetaData,
		CollectionId:           modifyCollectionsReq.CollectionId,
		CollectionBackColor:    modifyCollectionsReq.CollectionBackColor,
		CollectionImage:        modifyCollectionsReq.CollectionImage,
		CollectionIllustration: modifyCollectionsReq.CollectionIllustration,
		CollectionType:         modifyCollectionsReq.CollectionType,
	}

	filter := bson.D{{"collectionId", modifyCollectionsReq.CollectionId}}
	update := bson.D{{"$set", mongoCollectionDetails}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.COLLECTIONS, filter, update, opts)

	if err != nil {
		loggerconfig.Error("ModifyCollections Mongo Upsert failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var resp models.ModifyCollectionsResponse
	resp.CollectionId = modifyCollectionsReq.CollectionId
	resp.CollectionMetaData = modifyCollectionsReq.CollectionMetaData

	loggerconfig.Info("ModifyCollections Successful, response:", helpers.LogStructAsJSON(resp), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = resp
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CollectionsObj) FetchCollections(fetchCollectionsReq models.FetchCollectionsDetailsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var Collections models.MongoCollections
	err := dbops.MongoRepo.FindOne(constants.COLLECTIONS, bson.M{"collectionName": fetchCollectionsReq.CollectionId}, &Collections)

	var apiRes apihelpers.APIRes
	// if Collection does not exists
	if err != nil && Collections.CollectionName == "" {
		loggerconfig.Error("FetchCollections Collection Does Not Exist error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.CollectionDoesNotExists, http.StatusBadRequest)
	}

	var resp models.FetchCollectionsDetailsResponse
	resp.CollectionId = Collections.CollectionId
	resp.CollectionExchange = Collections.CollectionExchange
	resp.CollectionLongDesc = Collections.CollectionLongDesc
	resp.CollectionShortDesc = Collections.CollectionShortDesc
	resp.CollectionName = Collections.CollectionName
	resp.CollectionMetaData = Collections.CollectionMetaData
	resp.CollectionBackColor = Collections.CollectionBackColor
	resp.CollectionImage = Collections.CollectionImage
	resp.CollectionIllustration = Collections.CollectionIllustration
	resp.CollectionType = Collections.CollectionType

	loggerconfig.Info("FetchCollections Successful, response:", helpers.LogStructAsJSON(resp), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = resp
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CollectionsObj) DeleteCollections(deleteCollectionsReq models.DeleteCollectionsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	var Collections models.MongoCollections
	err := dbops.MongoRepo.FindOne(constants.COLLECTIONS, bson.M{"collectionName": deleteCollectionsReq.CollectionId}, &Collections)

	var apiRes apihelpers.APIRes
	// if Collection does not exists
	if err != nil && Collections.CollectionName == "" {
		loggerconfig.Error("DeleteCollections Collection Does Not Exist error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.CollectionDoesNotExists, http.StatusBadRequest)
	}

	filter := bson.D{{"collectionId", deleteCollectionsReq.CollectionId}}
	opts := options.Delete()
	_, err = dbops.MongoRepo.DeleteOne(constants.COLLECTIONS, filter, opts)
	if err != nil {
		loggerconfig.Error("DeleteCollections Mongo DeleteOne failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var deleteCollectionResponse models.DeleteCollecionsResponse
	deleteCollectionResponse.CollectionMetaData = Collections.CollectionMetaData

	loggerconfig.Info("DeleteCollections Successful, response:", helpers.LogStructAsJSON(deleteCollectionResponse), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = deleteCollectionResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj CollectionsObj) FetchAllCollections(reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var allCollections models.FetchAllCollectionsDetailsResponse
	var err error
	allDoc, err := dbops.MongoRepo.Find(constants.COLLECTIONS, bson.M{})
	if err != nil {
		loggerconfig.Error("FetchAllCollections Mongo Find() failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
	}
	err = allDoc.Decode(&allCollections)

	for allDoc.Next(context.Background()) {
		var fetchCollectionsDetailsResponse models.FetchCollectionsDetailsResponse
		err := allDoc.Decode(&fetchCollectionsDetailsResponse)
		if err != nil {
			loggerconfig.Error("FetchAllCollections Mongo parsing failed error =", err, " clientId:", reqH.ClientId, " requestId:", reqH.RequestId)
		}
		allCollections.FetchAllCollectionsDetailsResponse = append(allCollections.FetchAllCollectionsDetailsResponse, fetchCollectionsDetailsResponse)
	}

	var apiRes apihelpers.APIRes

	loggerconfig.Info("FetchAllCollections Successful, response:", helpers.LogStructAsJSON(allCollections), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = allCollections
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}
