package v1

import (
	"encoding/json"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var theBondEtfProvider models.BondEtfProvider

func InitBondEtfProvider(provider models.BondEtfProvider) {
	defer models.HandlePanic()
	theBondEtfProvider = provider
}

// FetchBondData
// @Tags space basket order V1
// @Description Fetch Bond Data
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.FetchBondDataReq true "Fetch Bond Data"
// @Success 200 {object} apihelpers.APIRes{data=models.NseBondStoreDbData}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/bondEtf/fetchBondData [POST]
func FetchBondData(c *gin.Context) {
	var reqParams models.FetchBondDataReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("FetchBondData (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchBondData (controller), Empty Device Type requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("FetchBondData (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " clientVersion: ", requestH.ClientVersion)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	if reqParams.Isin == "" {
		loggerconfig.Error("FetchBondData (controller), Empty isin requestid:", requestH.RequestId, "ClientID: ", requestH.ClientId)
		apihelpers.ErrorMessage(c, constants.EmptyIsin)
		return
	}

	loggerconfig.Info("FetchBondData (controller), reqParams:", helpers.LogStructAsJSON(reqParams), "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " clientVersion: ", requestH.ClientVersion)
	code, resp := theBondEtfProvider.FetchBondData(reqParams, requestH)
	logDetail := "isin: " + reqParams.Isin + " function: PendingOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
