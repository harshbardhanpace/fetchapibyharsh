package bondetf

import (
	"net/http"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/db"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
)

type BondEtfObj struct {
}

func InitBondEtfObj() BondEtfObj {
	defer models.HandlePanic()
	bondEtfObj := BondEtfObj{}

	return bondEtfObj
}

func (obj BondEtfObj) FetchBondData(bondDataReq models.FetchBondDataReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	nseBondStoreDb, err := db.GetPgObj().FetchNseBondData(bondDataReq.Isin)
	if err != nil && err.Error() == constants.InvalidBondIsin {
		loggerconfig.Error("FetchBondData, error in fetching FetchNseBondData ", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId," clientVersion: ", reqH.ClientVersion)
		return apihelpers.SendErrorResponse(false, constants.InvalidIsin, http.StatusOK) // api is success but isin is invalid
	}
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " FetchBondData, error in fetching FetchNseBondData ", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId," clientVersion: ", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("FetchBondData Successful, response:", helpers.LogStructAsJSON(nseBondStoreDb), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId," clientVersion: ", reqH.ClientVersion)
	apiRes.Data = nseBondStoreDb
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}
