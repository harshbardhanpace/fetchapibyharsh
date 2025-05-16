package alert

import (
	"net/http"
	apihelpers "space/apiHelpers"
	"space/db"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
)

type WarningObj struct {
}

func InitWarningObj() WarningObj {
	defer models.HandlePanic()
	warningObj := WarningObj{}

	return warningObj
}

func (obj WarningObj) NudgeAlert(nudgeAlertReq models.NudgeAlertReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	asmPresent, gsmPresent, err := db.GetPgObj().NudgeCheck(nudgeAlertReq.Isin)
	if err != nil {
		loggerconfig.Error("NudgeAlert, error in calling NudgeCheck ", err, "clientID: ", nudgeAlertReq.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var nudgeAlertRes models.NudgeAlertRes
	nudgeAlertRes.AsmPresent = asmPresent
	nudgeAlertRes.GsmPresent = gsmPresent

	loggerconfig.Info("NudgeAlert Successful, response:", helpers.LogStructAsJSON(nudgeAlertRes), "clientID: ", nudgeAlertReq.ClientId, " requestId:", reqH.RequestId)
	apiRes.Data = nudgeAlertRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}
