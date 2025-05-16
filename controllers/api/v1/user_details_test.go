package v1

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	apihelpers "space/apiHelpers"
	"space/loggerconfig"
	"space/models"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

var (
	getAllBankAccountsMock        func(req models.GetAllBankAccountsReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	getAllBankAccountsUpdatedMock func(req models.GetAllBankAccountsUpdatedReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	getUserIdMock                 func(req models.GetUserIdReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	userNotificationsMock         func(reqH models.ReqHeader) (int, apihelpers.APIRes)
	getClientStatusMock           func(emailId string, reqH models.ReqHeader) (int, apihelpers.APIRes)
)

type userDetailsMock struct{}

func (m userDetailsMock) GetAllBankAccounts(req models.GetAllBankAccountsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return getAllBankAccountsMock(req, reqH)
}

func (m userDetailsMock) GetAllBankAccountsUpdated(req models.GetAllBankAccountsUpdatedReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return getAllBankAccountsUpdatedMock(req, reqH)
}

func (m userDetailsMock) GetUserId(req models.GetUserIdReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return getUserIdMock(req, reqH)
}

func (m userDetailsMock) UserNotifications(req models.UserNotificationsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return userNotificationsMock(reqH)
}

func (m userDetailsMock) GetClientStatus(emailId string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return getClientStatusMock(emailId, reqH)
}

func TestUserNotifications(t *testing.T) {
	type args struct {
		c *gin.Context
	}

	gin.SetMode(gin.ReleaseMode)

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	//1 invalid request
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/space/v1/userDetails/userNotifications", nil)
	expected := "{\"status\":false,\"message\":\"INVALID REQUEST\",\"errorcode\":\"P11017\",\"data\":null}"
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"Invalid request", args{c: ctx}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UserNotifications(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("UserNotifications() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/space/v1/userDetails/userNotifications?clientId=test&page=sd&size=sdf", nil)
	expected = "{\"status\":false,\"message\":\"INVALID DEVICE TYPE\",\"errorcode\":\"P11034\",\"data\":null}"
	tests = []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"Blank device id", args{c: ctx}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UserNotifications(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("UserNotifications() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/space/v1/userDetails/userNotifications?clientId=test&page=sd&size=sdf", nil)
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitUserDetailsProvider(userDetailsMock{})

	//mock business layer response
	userNotificationsMock = func(reqH models.ReqHeader) (int, apihelpers.APIRes) {
		var res apihelpers.APIRes
		res.Status = true
		res.Message = "SUCCESS"
		return http.StatusOK, res
	}

	tests = []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"Success", args{c: ctx}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UserNotifications(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("UserNotifications() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("UserNotifications() = %v, want %v", w.Code, 200)
			}
		})
	}

}
