package v2

import (
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// import (
// 	u "GoProject/apiHelpers"
// 	v2s "GoProject/services/api/v2"
// 	"encoding/json"
// 	"github.com/gin-gonic/gin"
// )

// //UserList function will give you the list of users
// func UserList(c *gin.Context) {
// 	var userService v2s.UserService

// 	//decode the request body into struct and failed if any error occur
// 	err := json.NewDecoder(c.Request.Body).Decode(&userService.User)
// 	if err != nil {
// 		u.Respond(c.Writer, u.Message(1, "Invalid request"))
// 		return
// 	}

// 	//call service
// 	resp := userService.UserList()

// 	//return response using api helper
// 	u.Respond(c.Writer, resp)

// }

var LogoutProvider models.LogoutProvider

func InitLogoutProvider(provider models.LogoutProvider) {
	LogoutProvider = provider
}

// Logout
// @Tags space auth V2
// @Description Logout from single device
// @Param ClientId header string true "ClientId"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param Authorization header string true "Authorization Header"
// @Param DeviceToken header string false "DeviceToken"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v2/user/logout [DELETE]
func Logout(c *gin.Context) {

	cRH, _ := c.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	if requestH.DeviceType == "" {
		logrus.Error("logout V2 (controller), Empty Device type requestId: ", requestH.RequestId)
		apihelpers.ErrorMessage(c, constants.InvalidDeviceId)
		return
	}

	if requestH.ClientId == "" {
		logrus.Error("GetOverview (controller), Empty ClientId  requestId: ", requestH.RequestId, " deviceId: ", requestH.DeviceId, " deviceId: ", requestH.DeviceId)
		apihelpers.ErrorMessage(c, constants.InvalidRequest)
		return
	}

	logrus.Info("logout V2 (controller), requestId:", requestH.RequestId)

	code, resp := LogoutProvider.LogoutSingleDevice(requestH)
	logDetail := "clientId: " + requestH.ClientId + " function: logout V2 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(c, code, resp, logDetail)

}
