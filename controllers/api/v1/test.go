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

var theTestingProvider models.TestProvider

func InitTestProvider(provider models.TestProvider) {
	defer models.HandlePanic()
	theTestingProvider = provider
}

// TestingApi
// @Tags space Testing
// @Description Testing Api
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Param request body models.TestingReq true "TestingApi"
// @Success 200 {object} apihelpers.APIRes{data=models.TestingRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v1/test/testingRes [POST]
func TestingApi(c *gin.Context) {
	var testingReq models.TestingReq

	errr := json.NewDecoder(c.Request.Body).Decode(&testingReq)
	if errr != nil {
		loggerconfig.Error("TestingApi (controller), error decoding body, error:", errr)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	validate := validator.New()
	err := validate.Struct(testingReq)
	if err != nil {
		loggerconfig.Error("TestingApi (controller), Error validating struct: ", err, " requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	loggerconfig.Info("TestingApi (controller), reqParams:", helpers.LogStructAsJSON(testingReq), " uccId: ")

	code, resp := theTestingProvider.TestingApi(testingReq, requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: TestingApi requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)
}
