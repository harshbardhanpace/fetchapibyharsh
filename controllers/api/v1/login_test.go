package v1

import (
	"fmt"
	"io/ioutil"
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
	loginMock                 func(req models.LoginRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	loginByEmailMock          func(req models.LoginByEmailRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	validateTokenMock         func(req models.ValidateTokenRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	validateTwoFAMock         func(req models.ValidateTwoFARequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	setTwoFAPinMock           func(req models.SetTwoFAPinRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	forgotPasswordMock        func(req models.ForgotPasswordRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	setPasswordMock           func(req models.SetPasswordRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	forgetResetTwoFaMock      func(req models.ForgetResetTwoFaRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	GuestUserStatusMock       func(req models.GuestUserStatusReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	qrWebLogin                func(req models.LoginWithQRReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	unblockUserMock           func(req models.UnblockUserReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	createAppMock             func(req models.CreateAppReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	deleteAppMock             func(appId string, reqH models.ReqHeader) (int, apihelpers.APIRes)
	fetchAppMock              func(clientId string, reqH models.ReqHeader) (int, apihelpers.APIRes)
	handleAuthCode            func(authCode string, appState string, reqH models.ReqHeader) (int, apihelpers.APIRes)
	generateAccessToken       func(appState string, authCode string, reqH models.ReqHeader) (int, apihelpers.APIRes)
	getAccessToken            func(reqParams models.GetAccessTokenReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	forgetResetTwoFaEmailMock func(req models.ForgetResetEmailRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	forgetPasswordEmailMock   func(req models.ForgetResetEmailRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
)

type loginApisMock struct{}

func (m loginApisMock) LoginByPass(req models.LoginRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return loginMock(req, reqH)
}

func (m loginApisMock) LoginByEmail(req models.LoginByEmailRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return loginByEmailMock(req, reqH)
}

func (m loginApisMock) ValidateToken(req models.ValidateTokenRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return validateTokenMock(req, reqH)
}

func (m loginApisMock) ValidateTwoFa(req models.ValidateTwoFARequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return validateTwoFAMock(req, reqH)
}

func (m loginApisMock) SetTwoFaPin(req models.SetTwoFAPinRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return setTwoFAPinMock(req, reqH)
}

func (m loginApisMock) ForgetPassword(req models.ForgotPasswordRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return forgotPasswordMock(req, reqH)
}

func (m loginApisMock) SetPassword(req models.SetPasswordRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return setPasswordMock(req, reqH)
}

func (m loginApisMock) ForgetResetTwoFa(req models.ForgetResetTwoFaRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return forgetResetTwoFaMock(req, reqH)
}

func (m loginApisMock) GuestUserStatus(req models.GuestUserStatusReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return GuestUserStatusMock(req, reqH)
}

func (m loginApisMock) QRWebLogin(req models.LoginWithQRReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return qrWebLogin(req, reqH)
}

func (m loginApisMock) UnblockUser(req models.UnblockUserReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return unblockUserMock(req, reqH)
}

func (m loginApisMock) CreateApp(req models.CreateAppReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return createAppMock(req, reqH)
}

func (m loginApisMock) FetchApps(clientId string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return fetchAppMock(clientId, reqH)
}

func (m loginApisMock) DeleteApp(appId string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return deleteAppMock(appId, reqH)
}

func (m loginApisMock) GenerateAccessToken(appState string, authCode string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return generateAccessToken(appState, authCode, reqH)
}

func (m loginApisMock) GetAccessToken(reqParams models.GetAccessTokenReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return getAccessToken(reqParams, reqH)
}

func (m loginApisMock) HandleAuthCode(authCode string, appState string, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return handleAuthCode(authCode, appState, reqH)
}

//

func (m loginApisMock) ForgetResetTwoFaEmail(req models.ForgetResetEmailRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return forgetResetTwoFaEmailMock(req, reqH)
}

func (m loginApisMock) ForgetPasswordEmail(req models.ForgetResetEmailRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return forgetPasswordEmailMock(req, reqH)
}

func TestLogin(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/authapis/login", nil)
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
			Login(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("login() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/authapis/login", strings.NewReader("{\"userName\":\"CLIENT1\",\"password\":\"gkffd\"}"))
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
			Login(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Login() = %v, want %v", string(b), expected)
			}
		})
	}

	//4 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/authapis/login", strings.NewReader("{\"userName\":\"CLIENT1\",\"password\":\"gkffd\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitLoginProvider(loginApisMock{})

	//mock business layer response
	loginMock = func(req models.LoginRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			Login(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("login() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("login() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestValidateTwoFA(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/authapis/validateTwoFa", nil)
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
			ValidateTwoFA(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("validateTwoFA() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/authapis/validateTwoFa", strings.NewReader("{\"userName\":\"CLIENT1\",\"password\":\"gkffd\"}"))
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
			ValidateTwoFA(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("validateTwoFA() = %v, want %v", string(b), expected)
			}
		})
	}

	//4 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/authapis/validateTwoFa", strings.NewReader("{\"loginID\":\"CLIENT1\", \"twofa\": [{\"questionId\":\"22\",\"answer\":\"123456\"}], \"twofaToken\":\"22\", \"type\":\"PIN\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitLoginProvider(loginApisMock{})
	//mock business layer response
	validateTwoFAMock = func(req models.ValidateTwoFARequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			ValidateTwoFA(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("login() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("login() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestSetTwoFaPin(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/authapis/setTwoFaPin", nil)
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
			SetTwoFAPin(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("setTwoFaPin() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/authapis/setTwoFaPin", strings.NewReader("{\"loginID\":\"CLIENT1\",\"pin\":\"gkffd\",\"twoFaType\":\"gkffd\"}"))
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
			SetTwoFAPin(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("setTwoFAPin() = %v, want %v", string(b), expected)
			}
		})
	}

	//4 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/authapis/setTwoFaPin", strings.NewReader("{\"loginID\":\"CLIENT1\",\"pin\":\"gkffd\",\"twoFaType\":\"gkffd\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitLoginProvider(loginApisMock{})
	//mock business layer response
	setTwoFAPinMock = func(req models.SetTwoFAPinRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			SetTwoFAPin(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("setTwoFAPin() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("setTwoFAPin() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestForgotPassword(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/authapis/forgotPassword", nil)
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
			ForgotPassword(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("forgotPassword() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/authapis/forgotPassword", strings.NewReader("{\"loginID\":\"CLIENT1\",\"pan\":\"gkffd\"}"))
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
			ForgotPassword(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("forgotPassword() = %v, want %v", string(b), expected)
			}
		})
	}

	//4 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/authapis/forgotPassword", strings.NewReader("{\"loginID\":\"CLIENT1\",\"pan\":\"gkffd\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitLoginProvider(loginApisMock{})
	//mock business layer response
	forgotPasswordMock = func(req models.ForgotPasswordRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			ForgotPassword(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("forgotPassword() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("forgotPassword() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestSetPassword(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/authapis/setPassword", nil)
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
			SetPassword(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("setPassword() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/authapis/setPassword", strings.NewReader("{\"newPassword\":\"fdsf\",\"emailToken\":\"sdf\"}"))
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
			SetPassword(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("setPassword() = %v, want %v", string(b), expected)
			}
		})
	}

	//4 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/authapis/setPassword", strings.NewReader("{\"newPassword\":\"fdsf\",\"emailToken\":\"sdf\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitLoginProvider(loginApisMock{})
	//mock business layer response
	setPasswordMock = func(req models.SetPasswordRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			SetPassword(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("setPassword() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("setPassword() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestQRWebLogin(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/space/v1/authapis/webLogin", nil)
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
			QRWebLogin(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("TestQRWebLogin() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/space/v1/authapis/webLogin", strings.NewReader("{\"websocketID\":\"uuidstring\"}"))
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
			SetPassword(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("setPassword() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/space/v1/authapis/webLogin", strings.NewReader("{\"websocketID\":\"uuidstring\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitLoginProvider(loginApisMock{})

	//mock business layer response
	qrWebLogin = func(req models.LoginWithQRReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			QRWebLogin(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("TestQRWebLogin() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("TestQRWebLogin() = %v, want %v", w.Code, 400)
			}
		})
	}

}
