package v2

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
	loginV2Mock            func(req models.LoginV2Request, reqH models.ReqHeader) (int, apihelpers.APIRes)
	validateTwofaV2Mock    func(req models.ValidateTwofaV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes)
	setupTotpV2Mock        func(req models.SetupTotpV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes)
	chooseTwofaV2Mock      func(req models.ChooseTwofaV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes)
	forgetTotpV2Mock       func(req models.ForgetTotpV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes)
	validateLoginOtpV2Mock func(req models.ValidateLoginOtpV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes)
	setupBiometricV2Mock   func(req models.SetupBiometricV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes)
	disableBiometricV2Mock func(req models.DisableBiometricV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes)
	forgotPasswordV2Mock   func(req models.ForgotPasswordV2Request, reqH models.ReqHeader) (int, apihelpers.APIRes)
	unblockUserV2Mock      func(req models.UnblockUserV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes)
	getAccessTokenV2Mock   func(req models.GetAccessTokenV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes)
)

type loginApisMock struct{}

func (m loginApisMock) LoginV2(req models.LoginV2Request, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return loginV2Mock(req, reqH)
}

func (m loginApisMock) ValidateTwofaV2(req models.ValidateTwofaV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return validateTwofaV2Mock(req, reqH)
}

func (m loginApisMock) SetupTotpV2(req models.SetupTotpV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return setupTotpV2Mock(req, reqH)
}

func (m loginApisMock) ChooseTwofaV2(req models.ChooseTwofaV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return chooseTwofaV2Mock(req, reqH)
}

func (m loginApisMock) ForgetTotpV2(req models.ForgetTotpV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return forgetTotpV2Mock(req, reqH)
}

func (m loginApisMock) ValidateLoginOtpV2(req models.ValidateLoginOtpV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return validateLoginOtpV2Mock(req, reqH)
}

func (m loginApisMock) SetupBiometricV2(req models.SetupBiometricV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return setupBiometricV2Mock(req, reqH)
}

func (m loginApisMock) DisableBiometricV2(req models.DisableBiometricV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return disableBiometricV2Mock(req, reqH)
}

func (m loginApisMock) ForgetPasswordV2(req models.ForgotPasswordV2Request, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return forgotPasswordV2Mock(req, reqH)
}

func (m loginApisMock) UnblockUserV2(req models.UnblockUserV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return unblockUserV2Mock(req, reqH)
}

func (m loginApisMock) GetAccessTokenV2(req models.GetAccessTokenV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes){
	return getAccessTokenV2Mock(req, reqH)
}

func TestLoginV2Mock(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/login", nil)
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
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("login() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/login", strings.NewReader("{\"channelId\":\"chanel\",\"channelSecret\":\"123345\",\"type\":\"ty\"}"))
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
			b, _ := io.ReadAll(w.Body)
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/login", strings.NewReader("{\"channelId\":\"chanel\",\"channelSecret\":\"123345\",\"type\":\"ty\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	//Init order provider
	InitLoginProviderV2(loginApisMock{})

	//mock business layer response

	loginV2Mock = func(req models.LoginV2Request, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("forgetTotp() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("forgetTotp() = %v, want %v", w.Code, 200)
			}
		})
	}
}

func TestValidateTwofaV2Mock(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/validateTwofa", nil)
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
			ValidateTwofa(tt.args.c)
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("login() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/validateTwofa", strings.NewReader("{\"loginId\":\"login1\",\"twofa\":[{\"questionId\":\"22\",\"answer\":\"15241\"}],\"twofaToken\":\"ttoken\",\"type\":\"ty\",\"deviceType\":\"de11\"}"))
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
			ValidateTwofa(tt.args.c)
			b, _ := io.ReadAll(w.Body)
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/validateTwofa", strings.NewReader("{\"loginId\":\"login1\",\"twofa\":[{\"questionId\":\"22\",\"answer\":\"15241\"}],\"twofaToken\":\"ttoken\",\"type\":\"ty\",\"deviceType\":\"de11\"}"))

	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	//Init order provider
	InitLoginProviderV2(loginApisMock{})

	//mock business layer response

	validateTwofaV2Mock = func(req models.ValidateTwofaV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			ValidateTwofa(tt.args.c)
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("forgetTotp() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("forgetTotp() = %v, want %v", w.Code, 200)
			}
		})
	}
}

func TestSetupTotpV2Mock(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/setupTotp", nil)
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
			SetupTotp(tt.args.c)
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("login() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/setupTotp", strings.NewReader("{\"clientId\":\"client1\"}"))
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
			SetupTotp(tt.args.c)
			b, _ := io.ReadAll(w.Body)
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/setupTotp", strings.NewReader("{\"clientId\":\"client1\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	//Init order provider
	InitLoginProviderV2(loginApisMock{})

	//mock business layer response

	setupTotpV2Mock = func(req models.SetupTotpV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			SetupTotp(tt.args.c)
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("forgetTotp() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("forgetTotp() = %v, want %v", w.Code, 200)
			}
		})
	}
}

func TestChooseTwofaV2Mock(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/chooseTwofa", nil)
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
			ChooseTwofa(tt.args.c)
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("login() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/chooseTwofa", strings.NewReader("{\"loginId\":\"login1\",\"twofaType\":\"none\",\"totp\":\"1212\"}"))
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
			ChooseTwofa(tt.args.c)
			b, _ := io.ReadAll(w.Body)
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/chooseTwofa", strings.NewReader("{\"loginId\":\"login1\",\"twofa\":[{\"questionId\":\"22\",\"answer\":\"15241\"}],\"twofaToken\":\"ttoken\",\"type\":\"ty\",\"deviceType\":\"de11\"}"))

	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	//Init order provider
	InitLoginProviderV2(loginApisMock{})

	//mock business layer response

	chooseTwofaV2Mock = func(req models.ChooseTwofaV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			ChooseTwofa(tt.args.c)
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("forgetTotp() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("forgetTotp() = %v, want %v", w.Code, 200)
			}
		})
	}
}

func TestForgetTotpV2(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/forgetTotp", nil)
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
			ForgetTotp(tt.args.c)
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("login() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/forgetTotp", strings.NewReader("{\"loginId\":\"CLIENT1\",\"pan\":\"gkffd\"}"))
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
			ForgetTotp(tt.args.c)
			b, _ := io.ReadAll(w.Body)
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/forgetTotp", strings.NewReader("{\"loginId\":\"CLIENT1\",\"pan\":\"gkffd\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	//Init order provider
	InitLoginProviderV2(loginApisMock{})

	//mock business layer response

	forgetTotpV2Mock = func(req models.ForgetTotpV2Req, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			ForgetTotp(tt.args.c)
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("forgetTotp() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("forgetTotp() = %v, want %v", w.Code, 200)
			}
		})
	}

}

func TestValidateLoginOtpV2Mock(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/validateLoginOtp", nil)
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
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("login() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/validateLoginOtp", strings.NewReader("{\"referenceToken\":\"refktn\",\"otp\":\"1111\"}"))
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
			ValidateLoginOtp(tt.args.c)
			b, _ := io.ReadAll(w.Body)
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/validateLoginOtp", strings.NewReader("{\"referenceToken\":\"refktn\",\"otp\":\"1111\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	//Init order provider
	InitLoginProviderV2(loginApisMock{})

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
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("forgetTotp() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("forgetTotp() = %v, want %v", w.Code, 200)
			}
		})
	}
}

func TestSetupBiometricV2Mock(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/setupBiometric", nil)
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
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("login() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/setupBiometric", strings.NewReader("{\"clientId\":\"client1\",\"fingerprint\":\"123345\"}"))
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
			SetupBiometric(tt.args.c)
			b, _ := io.ReadAll(w.Body)
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/setupBiometric", strings.NewReader("{\"clientId\":\"client1\",\"fingerprint\":\"123345\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	//Init order provider
	InitLoginProviderV2(loginApisMock{})

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
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("forgetTotp() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("forgetTotp() = %v, want %v", w.Code, 200)
			}
		})
	}
}

func TestDisableBiometricV2Mock(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/disableBiometric", nil)
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
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("login() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/disableBiometric", strings.NewReader("{\"clientId\":\"client1\"}"))
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
			DisableBiometric(tt.args.c)
			b, _ := io.ReadAll(w.Body)
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v2/authapis/disableBiometric", strings.NewReader("{\"clientId\":\"client1\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	//Init order provider
	InitLoginProviderV2(loginApisMock{})

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
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("forgetTotp() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("forgetTotp() = %v, want %v", w.Code, 200)
			}
		})
	}
}
