package v2

import (
	"net/http"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var theProfileProviderV2 models.ProfileProvider

func InitProfileProviderV2(provider models.ProfileProvider) {
	defer models.HandlePanic()
	theProfileProviderV2 = provider
}

// Profile
// @Tags space profile V2
// @Description Profile V2 - Provide profile details of user
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param clientId query string true "clientId Query Parameter" dataType(string)
// @Success 200 {object} apihelpers.APIRes{data=models.ProfileResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/user/profile/getProfile [GET]
func GetProfile(c *gin.Context) {
	var reqParams models.ProfileRequest
	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)
	clientID := c.Query("clientId")
	if clientID == "" {
		loggerconfig.Error("GetProfile V2 (controller), error parsing the query params in Get request, not found error!")
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}
	reqParams.ClientID = clientID

	if requestH.DeviceType == "" {
		loggerconfig.Error("GetProfile V2 (controller), Empty Device Type requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientID)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceType)
		return
	}

	validate := validator.New()
	err := validate.Struct(reqParams)
	if err != nil {
		loggerconfig.Error("GetProfile V2 (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId, "clientId: ", reqParams.ClientID)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	matchStatus, tokenValidStatus := helpers.CheckAuthWithClient(reqParams.ClientID, requestH.Authorization)
	if !tokenValidStatus {
		loggerconfig.Error("GetProfile V2 (controller) CheckAuthWithClient invalid authtoken", " clientId: ", reqParams.ClientID, " requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.InvalidToken, http.StatusUnauthorized)
		return
	}
	if !matchStatus {
		loggerconfig.Error("GetProfile V2 (controller) CheckAuthWithClient difference in authtoken-clientId and clientId", " clientId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
		apihelpers.SendErrorController(c, false, constants.MismatchAuthClient, http.StatusForbidden)
		return
	}

	loggerconfig.Info("GetProfile V2 (controller), reqParams:", helpers.LogStructAsJSON(reqParams), " uccId: ", reqParams.ClientID, "requestId:", requestH.RequestId)
	code, resp := theProfileProviderV2.GetProfile(reqParams, requestH)
	logDetail := "clientId: " + reqParams.ClientID + " function: PendingOrder requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}
