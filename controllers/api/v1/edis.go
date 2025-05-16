package v1

import (
	"encoding/json"
	"net/http"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var theEdisProvider models.EdisProvider

func InitEdisProvider(provider models.EdisProvider) {
	defer models.HandlePanic()
	theEdisProvider = provider
}

// EdisRequest
// @Tags space edis V1
// @Description Edis Request
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.EdisReq true "EdisRequest"
// @Success 200 {object} apihelpers.APIRes{data=models.EdisRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/edis/edisRequest [POST]
func EdisRequest(c *gin.Context) {
	var reqParams models.EdisReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("EdisRequest (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")

	requestH, _ := (cRH).(models.ReqHeader)

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("EdisRequest (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("EdisRequest (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("EdisRequest (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("EdisRequest (controller), reqParams:", reqParams, "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theEdisProvider.SendEdisRequest(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: EdisRequest requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}

// GenerateTpin
// @Tags space edis V1
// @Description Request
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param P-ClientVersion  header string false "P-ClientVersion Header"
// @Param request body models.TpinReq true "GenerateTpin"
// @Success 200 {object} apihelpers.APIRes{data=models.TpinRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/edis/generateTpin [POST]
func GenerateTpin(c *gin.Context) {
	var reqParams models.TpinReq
	err := json.NewDecoder(c.Request.Body).Decode(&reqParams)
	if err != nil {
		loggerconfig.Error("GenerateTpin (controller), error decoding body, error:", err)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")

	requestH, _ := (cRH).(models.ReqHeader)

	validate := validator.New()
	err = validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GenerateTpin (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("GenerateTpin (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("GenerateTpin (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("GenerateTpin (controller), reqParams:", reqParams, "requestId:", requestH.RequestId, "ClientID: ", requestH.ClientId, " deviceId: ", requestH.DeviceId)
	code, resp := theEdisProvider.GenerateTpin(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: GenerateTpin requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
