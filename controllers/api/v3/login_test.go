package v3

import (
	"encoding/json"
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
	setPasswordMock        func(models.SetPasswordRequest, models.ReqHeader) (int, apihelpers.APIRes)
	validateTokenMock      func(models.ValidateTokenRequest, models.ReqHeader) (int, apihelpers.APIRes)
	forgetResetTwoFaMock   func(models.ForgetResetTwoFaRequest, models.ReqHeader) (int, apihelpers.APIRes)
	validateLoginOtpV2Mock func(models.ValidateLoginOtpV2Req, models.ReqHeader) (int, apihelpers.APIRes)
	setupBiometricV2Mock   func(models.SetupBiometricV2Req, models.ReqHeader) (int, apihelpers.APIRes)
	disableBiometricV2Mock func(models.DisableBiometricV2Req, models.ReqHeader) (int, apihelpers.APIRes)
)

type loginV3Mock struct{}

func (m loginV3Mock) SetPassword(req models.SetPasswordRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return setPasswordMock(req, reqH)
}

func (m loginV3Mock) ValidateToken(req models.ValidateTokenRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return validateTokenMock(req, reqH)
}

func (m loginV3Mock) ForgetResetTwoFa(req models.ForgetResetTwoFaRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return forgetResetTwoFaMock(req, reqH)
}

func (m loginV3Mock) ValidateLoginOtpV2(req models.ValidateLoginOtpV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return validateLoginOtpV2Mock(req, reqH)
}

func (m loginV3Mock) SetupBiometricV2(req models.SetupBiometricV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return setupBiometricV2Mock(req, reqH)
}

func (m loginV3Mock) DisableBiometricV2(req models.DisableBiometricV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return disableBiometricV2Mock(req, reqH)
}

func (m loginV3Mock) LoginByEmailOtp(req models.LoginByEmailOtpReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return 0, apihelpers.APIRes{
		Status:  true,
		Message: "SUCCESS",
		Data:    nil,
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
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v3/authapis/setPassword", nil)
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
				t.Errorf("SetPassword() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	var requestPacket models.SetPasswordRequest
	requestPacket.NewPass = "abcd"
	requestPacket.OldPass = "abc"

	reqPacket, _ := json.Marshal(requestPacket)

	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v3/authapis/setPassword", strings.NewReader(string(reqPacket)))
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
			SetPassword(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("SetPassword() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v3/authapis/setPassword", strings.NewReader(string(reqPacket)))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"
	//Init order provider
	InitLoginProviderV3(loginV3Mock{})

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
				t.Errorf("SetPassword() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}

}

func TestValidateToken(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/space/v3/authapis/validateToken", nil)
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
			ValidateToken(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("ValidateToken() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/space/v3/authapis/validateToken?token=101&userId=abc", nil)
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
			ValidateToken(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("ValidateToken() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/space/v3/authapis/validateToken?token=101&userId=abc", nil)
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"
	//Init order provider
	InitLoginProviderV3(loginV3Mock{})

	//mock business layer response
	validateTokenMock = func(req models.ValidateTokenRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			ValidateToken(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("ValidateToken() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}
}

func TestForgetResetTwoFa(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v3/authapis/forgetResetTwoFa", nil)
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
			ForgetResetTwoFa(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("ForgetResetTwoFa() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	var requestPacket models.ForgetResetTwoFaRequest
	requestPacket.ClientID = "abc"
	requestPacket.Pan = "abc"

	reqPacket, _ := json.Marshal(requestPacket)

	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v3/authapis/forgetResetTwoFa", strings.NewReader(string(reqPacket)))
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
			ForgetResetTwoFa(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("ForgetResetTwoFa() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v3/authapis/forgetResetTwoFa", strings.NewReader(string(reqPacket)))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"
	//Init order provider
	InitLoginProviderV3(loginV3Mock{})

	//mock business layer response
	forgetResetTwoFaMock = func(req models.ForgetResetTwoFaRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			ForgetResetTwoFa(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("ForgetResetTwoFa() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}

}

func TestValidateLoginOtp(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v3/authapis/validateLoginOtp", nil)
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
			ValidateLoginOtp(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("ValidateLoginOtp() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	var requestPacket models.ValidateLoginOtpV2Req
	requestPacket.ReferenceToken = "abc"
	requestPacket.Otp = "abc"

	reqPacket, _ := json.Marshal(requestPacket)

	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v3/authapis/validateLoginOtp", strings.NewReader(string(reqPacket)))
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
			ValidateLoginOtp(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("ValidateLoginOtp() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v3/authapis/validateLoginOtp", strings.NewReader(string(reqPacket)))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"
	//Init order provider
	InitLoginProviderV3(loginV3Mock{})

	//mock business layer response
	validateLoginOtpV2Mock = func(req models.ValidateLoginOtpV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			ValidateLoginOtp(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("ValidateLoginOtp() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}

}

func TestSetupBiometric(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v3/authapis/setupBiometric", nil)
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
			SetupBiometric(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("SetupBiometric() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	var requestPacket models.SetupBiometricV2Req
	requestPacket.ClientID = "abc"
	requestPacket.FingerPrint = "abc"

	reqPacket, _ := json.Marshal(requestPacket)

	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v3/authapis/setupBiometric", strings.NewReader(string(reqPacket)))
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
			SetupBiometric(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("SetupBiometric() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v3/authapis/setupBiometric", strings.NewReader(string(reqPacket)))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"
	//Init order provider
	InitLoginProviderV3(loginV3Mock{})

	//mock business layer response
	setupBiometricV2Mock = func(req models.SetupBiometricV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			SetupBiometric(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("SetupBiometric() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}

}

func TestDisableBiometric(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/api/space/v3/authapis/disableBiometric", nil)
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
			DisableBiometric(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("DisableBiometric() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	var requestPacket models.DisableBiometricV2Req
	requestPacket.ClientID = "abc"
	requestPacket.FingerPrint = "abc"

	reqPacket, _ := json.Marshal(requestPacket)

	ctx.Request = httptest.NewRequest(http.MethodDelete, "/api/space/v3/authapis/disableBiometric", strings.NewReader(string(reqPacket)))
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
			DisableBiometric(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("DisableBiometric() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/api/space/v3/authapis/disableBiometric", strings.NewReader(string(reqPacket)))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"
	//Init order provider
	InitLoginProviderV3(loginV3Mock{})

	//mock business layer response
	disableBiometricV2Mock = func(req models.DisableBiometricV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			DisableBiometric(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("DisableBiometric() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}

}
