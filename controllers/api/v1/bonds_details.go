package v1

import (
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var BondsDataProvider models.BondsDetailsProvider

func InitBondsDetailsProvider(provider models.BondsDetailsProvider) {
	defer models.HandlePanic()
	BondsDataProvider = provider
}

// FetchBondDataByIsin
// @Tags space bonds details v1
// @Description Fetch Bond Data By Isin
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param isin query string true "isin Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.FetchBondDataByIsinResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/bondsDetails/fetchBondDataByIsin [GET]
func FetchBondDataByIsin(c *gin.Context) {

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	var BondDataReq models.FetchBondDataByIsinReq
	BondDataReq.Isin = c.Query("isin")

	if requestH.DeviceType == "" {
		loggerconfig.Error("FetchBondDataByIsin (controller), Empty Device Type requestId: ", requestH.RequestId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(BondDataReq)
	if err != nil {
		loggerconfig.Error("FetchBondDataByIsin (controller), error validating request packet error:", err, "requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("FetchBondDataByIsin (controller), reqParams: ", helpers.LogStructAsJSON(BondDataReq), " uccId: ", requestH.ClientId, " requestId:", requestH.RequestId, " deviceId: ", requestH.DeviceId, " clientVersion: ", requestH.ClientVersion)

	code, resp := BondsDataProvider.FetchBondDataByIsin(BondDataReq, requestH)

	logDetails := "isin: " + BondDataReq.Isin + " function: FetchBondDataByIsin requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetails)
}
