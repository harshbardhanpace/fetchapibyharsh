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
)

var theWarningProvider models.WarningProvider

func InitWarningProvider(provider models.WarningProvider) {
	defer models.HandlePanic()
	theWarningProvider = provider
}

// NudgeAlert
// @Tags space warning V1
// @Description nudlge alert
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.NudgeAlertReq true "NudgeAlert"
// @Success 200 {object} apihelpers.APIRes{data=models.NudgeAlertRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/warning/nudgeAlert [POST]
func NudgeAlert(c *gin.Context) {
	var nudgeAlertReq models.NudgeAlertReq

	var reqH models.ReqHeader
	if err := c.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("NudgeAlert (controller), error decoding header, error:", err)
		apihelpers.ErrorMessage(c, constants.DecodingHeaderError)
		return
	}

	err := json.NewDecoder(c.Request.Body).Decode(&nudgeAlertReq)
	if err != nil {
		loggerconfig.Error("NudgeAlert (controller), Invalid Request, error:", err, " requestId: ", reqH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(nudgeAlertReq.ClientId, reqH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("GetAllBankAccounts (controller) CheckAuthWithClient invalid authtoken", " clientId: ", nudgeAlertReq.ClientId, " requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("NudgeAlert (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", nudgeAlertReq.ClientId, "requestId:", reqH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("NudgeAlert (controller), reqParams:", helpers.LogStructAsJSON(nudgeAlertReq), " uccId: ", nudgeAlertReq.ClientId, "requestId:", reqH.RequestId, "clientId: ", nudgeAlertReq.ClientId)

	//call service
	code, resp := theWarningProvider.NudgeAlert(nudgeAlertReq, reqH)

	//return response using api helper
	logDetail := "clientId: " + reqH.ClientId + " function: NudgeAlert requestId: " + reqH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}
