package v1

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apihelpers "space/apiHelpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
)

var (
	getProfileMock    func(req models.ProfileRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	sendAFOtpMock     func(req models.SendAFOtpReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	verifyAFOtpMock   func(req models.VerifyAFOtpReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	accountFreezeMock func(req models.AccountFreezeReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
)

type profileMock struct{}

func (m profileMock) GetProfile(req models.ProfileRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return getProfileMock(req, reqH)
}

func (m profileMock) SendAFOtp(req models.SendAFOtpReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return sendAFOtpMock(req, reqH)
}

func (m profileMock) VerifyAFOtp(req models.VerifyAFOtpReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return verifyAFOtpMock(req, reqH)
}

func (m profileMock) AccountFreeze(req models.AccountFreezeReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return accountFreezeMock(req, reqH)
}

func TestGetProfile(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/user/profile/getProfile", nil)
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
			GetProfile(tt.args.c)
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("GetProfile() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/user/profile/getProfile", strings.NewReader("{\"clientId\":\"Client1\"}"))
	expected = "{\"status\":false,\"message\":\"INVALID DEVICE ID\",\"errorcode\":\"P11034\",\"data\":null}"
	tests = []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"Blank device id", args{c: ctx}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetProfile(tt.args.c)
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("GetProfile() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/user/profile/getProfile", strings.NewReader("{\"clientId\":\"Client1\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitProfileProvider(profileMock{})

	//mock business layer response
	getProfileMock = func(req models.ProfileRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			GetProfile(tt.args.c)
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("GetProfile() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}
}
