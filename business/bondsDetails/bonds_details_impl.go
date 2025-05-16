package bondsdetails

import (
	"encoding/json"
	"fmt"
	"net/http"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/db"
	"space/dbops"
	"space/helpers/cache"
	"space/loggerconfig"
	"space/models"

	"go.mongodb.org/mongo-driver/bson"
)

type BondsDetailsObj struct {
	contractCacheCli cache.ContractCache
	mongodb          db.MongoDatabase
}

func InitBondsDetailsProvider(mongodb db.MongoDatabase, contractCacheCli cache.ContractCache) BondsDetailsObj {
	defer models.HandlePanic()
	return BondsDetailsObj{
		mongodb:          mongodb,
		contractCacheCli: contractCacheCli,
	}
}

func (obj BondsDetailsObj) FetchBondDataByIsin(req models.FetchBondDataByIsinReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	var bondDataResponse models.FetchBondDataByIsinResponse
	var apiRes apihelpers.APIRes

	err := dbops.MongoRepo.FindOne(constants.BondsDataCollection, bson.M{"isin": req.Isin}, &bondDataResponse)
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("Alert Severity:P1-High, FetchBondDataByIsin  Error while fetching data from mongo. Error:", err, "isin:", req.Isin,  "clientID:", reqH.ClientId, "requestId:", reqH.RequestId, "platform:", reqH.Platform, "clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	fetchStockDetail := func(prefix, isin string) (models.ContractDetails, error) {
		key := fmt.Sprintf("%s-%s", prefix, isin)
		err, val := obj.contractCacheCli.GetFromHash("isin_data", key)
		if err != nil {
			return models.ContractDetails{}, err
		}
		var detail models.ContractDetails
		if err = json.Unmarshal([]byte(val), &detail); err != nil {
			return models.ContractDetails{}, err
		}
		return detail, nil
	}

	nseDetail, err := fetchStockDetail("NSE", req.Isin)
	if err != nil {
		loggerconfig.Error("FetchBondDataByIsin: Failed to fetch NSE stock details Error:", err, "isin:", req.Isin,  "clientID:", reqH.ClientId, "requestId:", reqH.RequestId, "platform:", reqH.Platform, "clientVersion:", reqH.ClientVersion)
	}

	bseDetail, err := fetchStockDetail("BSE", req.Isin)
	if err != nil {
		loggerconfig.Error("FetchBondDataByIsin: Failed to fetch BSE stock details Error:", err, "isin:", req.Isin,  "clientID:", reqH.ClientId, "requestId:", reqH.RequestId, "platform:", reqH.Platform, "clientVersion:", reqH.ClientVersion)
	}

	bondDataResponse.NseToken = nseDetail.Token1
	bondDataResponse.BseToken = bseDetail.Token1

	apiRes.Data = bondDataResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}
