package pins

import (
	"net/http"
	"sort"

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

type PinsObj struct {
	mongodb db.MongoDatabase
}

func InitPins(mongodb db.MongoDatabase) PinsObj {
	defer models.HandlePanic()
	pinsObj := PinsObj{mongodb: mongodb}
	return pinsObj
}

func (obj PinsObj) FetchPins(req models.PinsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var pinsMetaData models.MongoPinsMetadata
	err := dbops.MongoRepo.FindOne(constants.PINS, bson.M{"clientId": req.ClientId}, &pinsMetaData)
	var apiRes apihelpers.APIRes
	if err != nil && err.Error() == constants.MongoNoDocError {
		loggerconfig.Error("FetchPins", req, " mongo err:", err, "requestId: ", reqH.RequestId)
		pinsMetaData.ClientId = req.ClientId
		filter := bson.D{{"clientId", req.ClientId}}
		pinsMetaData.PinsMetaDatas = populateDefaultPins()
		update := bson.D{{"$set", pinsMetaData}}
		opts := options.Update().SetUpsert(true)
		err = dbops.MongoRepo.UpdateOne(constants.PINS, filter, update, opts)
		if err != nil {
			loggerconfig.Error("FetchPin  pins", req, " mongo err:", err, "clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
	}

	if err != nil {
		loggerconfig.Error("FetchPins", req, "  err:", err, "clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	sort.Slice(pinsMetaData.PinsMetaDatas, func(i, j int) bool {
		return pinsMetaData.PinsMetaDatas[i].PinIndex < pinsMetaData.PinsMetaDatas[j].PinIndex
	})

	var resp models.PinsResponse
	resp.ClientId = req.ClientId
	resp.PinsMetaDatas = pinsMetaData.PinsMetaDatas

	loggerconfig.Info("FetchPins Successful, response:", helpers.LogStructAsJSON(resp), " uccId:", req.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = resp
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PinsObj) UpdatePins(req models.UpdatePins, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var pinsMetaData models.MongoPinsMetadata
	err := dbops.MongoRepo.FindOne(constants.PINS, bson.M{"clientId": req.ClientId}, &pinsMetaData)
	var apiRes apihelpers.APIRes
	if err != nil && err.Error() == constants.MongoNoDocError {
		loggerconfig.Error("UpdatePins", req, " mongo err:", err, "clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.PinsDoesNotExists, http.StatusBadRequest)
	}

	if len(req.PinsMetaDatas) > len(pinsMetaData.PinsMetaDatas) {
		loggerconfig.Error(constants.ERROR, "UpdatePins: Too many pins requested. UserID:", req.ClientId, " requestid:", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.PinsSizeExceed, http.StatusBadRequest)
	}

	for i := 0; i < len(req.PinsMetaDatas); i++ {
		isValid := false
		for j := 0; j < len(pinsMetaData.PinsMetaDatas); j++ {
			if req.PinsMetaDatas[i].PinId == pinsMetaData.PinsMetaDatas[j].PinId {
				isValid = true
				break
			}
		}
		if !(req.PinsMetaDatas[i].PinIndex > 0 && req.PinsMetaDatas[i].PinIndex <= len(pinsMetaData.PinsMetaDatas)) {
			loggerconfig.Error(constants.ERROR, "UpdatePins: Invalid pins Index requested. UserID:", req.ClientId, " requestid:", reqH.RequestId)
			return apihelpers.SendErrorResponse(false, constants.PinsIndexInvalid, http.StatusBadRequest)
		}
		if !isValid {
			loggerconfig.Error(constants.ERROR, "UpdatePins: Invalid pins requested. UserID:", req.ClientId, " requestid:", reqH.RequestId)
			return apihelpers.SendErrorResponse(false, constants.PinsDoesNotExists, http.StatusBadRequest)
		}
	}

	sort.Slice(req.PinsMetaDatas, func(i, j int) bool {
		return req.PinsMetaDatas[i].PinIndex < req.PinsMetaDatas[j].PinIndex
	})

	pinsMetaDataEntries := make([]models.PinsMetaData, 0)
	for _, reqMetaData := range req.PinsMetaDatas {
		filteredPinsMetaData := make([]models.PinsMetaData, 0)
		for _, pinMetaData := range pinsMetaData.PinsMetaDatas {
			if pinMetaData.PinId != reqMetaData.PinId {
				filteredPinsMetaData = append(filteredPinsMetaData, pinMetaData)
			}
		}
		pinsMetaData.PinsMetaDatas = filteredPinsMetaData
	}

	pinsMetaDataPtr := 0
	reqMetaDataPtr := 0

	for pinsMetaDataPtr < len(pinsMetaData.PinsMetaDatas) && reqMetaDataPtr < len(req.PinsMetaDatas) {
		if req.PinsMetaDatas[reqMetaDataPtr].PinIndex-1 == len(pinsMetaDataEntries) {
			pinsMetaDataEntries = append(pinsMetaDataEntries, req.PinsMetaDatas[reqMetaDataPtr])
			reqMetaDataPtr++
		} else {
			pinsMetaDataEntries = append(pinsMetaDataEntries, pinsMetaData.PinsMetaDatas[pinsMetaDataPtr])
			pinsMetaDataPtr++
		}
	}
	for pinsMetaDataPtr < len(pinsMetaData.PinsMetaDatas) {
		pinsMetaDataEntries = append(pinsMetaDataEntries, pinsMetaData.PinsMetaDatas[pinsMetaDataPtr])
		pinsMetaDataPtr++
	}
	for reqMetaDataPtr < len(req.PinsMetaDatas) {
		pinsMetaDataEntries = append(pinsMetaDataEntries, req.PinsMetaDatas[reqMetaDataPtr])
		reqMetaDataPtr++
	}

	for i := 0; i < len(pinsMetaDataEntries); i++ {
		pinsMetaDataEntries[i].PinIndex = i + 1
	}

	req.PinsMetaDatas = pinsMetaDataEntries

	//req.PinsMetaDatas = populateDefaultPins() //run to reset to default for testing and such

	filter := bson.D{{"clientId", req.ClientId}}
	update := bson.D{{"$set", req}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.PINS, filter, update, opts)
	if err != nil {
		loggerconfig.Error("UpdatePins", helpers.LogStructAsJSON(req), " uccId:", req.ClientId, "Update err:", err, "clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PinsObj) AddPins(req models.AddPinReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var pinsMetaData models.MongoPinsMetadata
	err := dbops.MongoRepo.FindOne(constants.PINS, bson.M{"clientId": req.ClientId}, &pinsMetaData)
	var apiRes apihelpers.APIRes
	if err != nil && err.Error() == constants.MongoNoDocError {
		loggerconfig.Error("AddPins", req, " mongo err:", err, "clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.PinsDoesNotExists, http.StatusBadRequest)
	}

	if len(pinsMetaData.PinsMetaDatas) >= 6 {
		loggerconfig.Error("AddPins All Pins already full for userid:", req.ClientId, " requestid:", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.PinsCapacityFull, http.StatusOK)
	}

	var toBeAdded models.UpdatePins
	toBeAdded.ClientId = req.ClientId

	pinsMetaDataEntries := make([]models.PinsMetaData, len(pinsMetaData.PinsMetaDatas))
	copy(pinsMetaDataEntries, pinsMetaData.PinsMetaDatas)
	noOfNewStocks := len(req.StockDetailsData)

	if (noOfNewStocks + len(pinsMetaData.PinsMetaDatas)) > constants.PinsSize {
		loggerconfig.Error(constants.ERROR, "AddPins: Too many pins requested. UserID:", req.ClientId, " requestid:", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.PinsSizeExceed, http.StatusBadRequest)

	}

	for i := 0; i < noOfNewStocks; i++ {
		var instanceToBeAdded models.PinsMetaData
		instanceToBeAdded.PinId = uuid.New().String()
		instanceToBeAdded.PinIndex = len(pinsMetaDataEntries) + 1
		instanceToBeAdded.StockDet.StockId = uuid.New().String()
		instanceToBeAdded.StockDet.Company = req.StockDetailsData[i].Company
		instanceToBeAdded.StockDet.DisplayName = req.StockDetailsData[i].DisplayName
		instanceToBeAdded.StockDet.Exchange = req.StockDetailsData[i].Exchange
		instanceToBeAdded.StockDet.Expiry = req.StockDetailsData[i].Expiry
		instanceToBeAdded.StockDet.Token = req.StockDetailsData[i].Token
		instanceToBeAdded.StockDet.TradingSymbol = req.StockDetailsData[i].TradingSymbol
		instanceToBeAdded.StockDet.Symbol = req.StockDetailsData[i].Symbol
		instanceToBeAdded.StockDet.IsTradable = true
		instanceToBeAdded.StockDet.Isin = ""
		instanceToBeAdded.StockDet.Segment = req.StockDetailsData[i].Segment

		pinsMetaDataEntries = append(pinsMetaDataEntries, instanceToBeAdded)
	}

	toBeAdded.PinsMetaDatas = pinsMetaDataEntries
	filter := bson.D{{"clientId", req.ClientId}}
	update := bson.D{{"$set", toBeAdded}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.PINS, filter, update, opts)
	if err != nil {
		loggerconfig.Error("AddPins", toBeAdded, " uccId:", req.ClientId, "Update err:", err, "requestId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj PinsObj) DeletePins(req models.DeletePins, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var pinsMetaData models.MongoPinsMetadata
	err := dbops.MongoRepo.FindOne(constants.PINS, bson.M{"clientId": req.ClientId}, &pinsMetaData)
	var apiRes apihelpers.APIRes
	if err != nil && err.Error() == constants.MongoNoDocError {
		loggerconfig.Error("DeletePins", req, " mongo err:", err, "clientID: ", req.ClientId, "requestId: ", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.PinsDoesNotExists, http.StatusBadRequest)
	}

	if len(req.PinId) > len(pinsMetaData.PinsMetaDatas) {
		loggerconfig.Error(constants.ERROR, "DeletePins: Too many pins deletion requested. UserID:", req.ClientId, " requestid:", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.PinsDoesNotExists, http.StatusBadRequest)
	}

	for i := 0; i < len(req.PinId); i++ {
		isValid := false
		for j := 0; j < len(pinsMetaData.PinsMetaDatas); j++ {
			if req.PinId[i] == pinsMetaData.PinsMetaDatas[j].PinId {
				isValid = true
				break
			}
		}

		if !isValid {
			loggerconfig.Error(constants.ERROR, "DeletePins: Invalid pins requested. UserID:", req.ClientId, " requestid:", reqH.RequestId)
			return apihelpers.SendErrorResponse(false, constants.PinsDoesNotExists, http.StatusBadRequest)
		}
	}

	pinsMetaDataEntries := make([]models.PinsMetaData, 0)
	for i := 0; i < len(pinsMetaData.PinsMetaDatas); i++ {
		flag := true
		for j := 0; j < len(req.PinId); j++ {
			if req.PinId[j] == pinsMetaData.PinsMetaDatas[i].PinId {
				flag = false
				break
			}
		}
		if flag {
			pinsMetaData.PinsMetaDatas[i].PinIndex = len(pinsMetaDataEntries) + 1
			pinsMetaDataEntries = append(pinsMetaDataEntries, pinsMetaData.PinsMetaDatas[i])
		}
	}
	pinsMetaData.PinsMetaDatas = pinsMetaDataEntries

	filter := bson.D{{"clientId", req.ClientId}}
	update := bson.D{{"$set", pinsMetaData}}
	opts := options.Update().SetUpsert(true)
	err = dbops.MongoRepo.UpdateOne(constants.PINS, filter, update, opts)
	if err != nil {
		loggerconfig.Error("DeletePins", pinsMetaData, " uccId:", req.ClientId, "Update err:", err, "requestId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func populateDefaultPins() []models.PinsMetaData {
	var stockd []models.StockDetails
	// populate
	// by default add Nifty 50, NIfty Bank, Sensex ,
	//Nifty IT, NIFTY INFRA and Nifty Pharma
	var stock1 models.StockDetails
	stock1.Company = "SENSEX"
	stock1.DisplayName = "SENSEX"
	stock1.Exchange = "BSE"
	stock1.StockId = uuid.New().String()
	stock1.Token = "1"
	stock1.TradingSymbol = "SENSEX"
	stock1.Symbol = "SENSEX"
	stock1.Segment = "Indices"

	stockd = append(stockd, stock1)

	var stock2 models.StockDetails
	stock2.Company = "Nifty 50"
	stock2.DisplayName = "Nifty 50"
	stock2.Exchange = "NSE"
	stock2.StockId = uuid.New().String()
	stock2.Token = "26000"
	stock2.TradingSymbol = "Nifty 50"
	stock2.Symbol = "Nifty 50"
	stock2.Segment = "Indices"

	stockd = append(stockd, stock2)

	var stock3 models.StockDetails
	stock3.Company = "Nifty Bank"
	stock3.DisplayName = "Nifty Bank"
	stock3.Exchange = "NSE"
	stock3.StockId = uuid.New().String()
	stock3.Token = "26009"
	stock3.TradingSymbol = "Nifty Bank"
	stock3.Symbol = "Nifty Bank"
	stock3.Segment = "Indices"

	stockd = append(stockd, stock3)

	var stock4 models.StockDetails
	stock4.Company = "India VIX"
	stock4.DisplayName = "India VIX"
	stock4.Exchange = "NSE"
	stock4.StockId = uuid.New().String()
	stock4.Token = "26017"
	stock4.TradingSymbol = "India VIX"
	stock4.Symbol = "India VIX"
	stock4.Segment = "Indices"

	stockd = append(stockd, stock4)

	var stock5 models.StockDetails
	stock5.Company = "Nifty Midcap 50"
	stock5.DisplayName = "Nifty Midcap 50"
	stock5.Exchange = "NSE"
	stock5.StockId = uuid.New().String()
	stock5.Token = "26014"
	stock5.TradingSymbol = "Nifty Midcap 50"
	stock5.Symbol = "Nifty Midcap 50"
	stock5.Segment = "Indices"

	stockd = append(stockd, stock5)

	var stock6 models.StockDetails
	stock6.Company = "Nifty Fin Service"
	stock6.DisplayName = "Nifty Fin Service"
	stock6.Exchange = "NSE"
	stock6.StockId = uuid.New().String()
	stock6.Token = "26037"
	stock6.TradingSymbol = "Nifty Fin Service"
	stock6.Symbol = "Nifty Fin Service"
	stock6.Segment = "Indices"

	stockd = append(stockd, stock6)

	var pins []models.PinsMetaData

	for k, v := range stockd {
		var pin models.PinsMetaData
		pin.PinId = uuid.New().String()
		pin.PinIndex = k
		pin.StockDet = v
		pins = append(pins, pin)
	}

	return pins
}
