package testing

import (
	"net/http"

	apihelpers "space/apiHelpers"
	"space/models"

	"github.com/pocketful-tech/throttling"
)

type TestingObj struct {
	testingThrottleObj *throttling.APIThrottler
}

func InitTesting(testingThrottleObject *throttling.APIThrottler) TestingObj {
	testingObj := TestingObj{}
	testingObj.testingThrottleObj = testingThrottleObject
	return testingObj
}

func (obj TestingObj) TestingApi(req models.TestingReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	obj.testingThrottleObj.Throttle()

	var testingRes models.TestingRes
	testingRes.DummyResData = "This is dummy data"
	apiRes.Data = testingRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}
