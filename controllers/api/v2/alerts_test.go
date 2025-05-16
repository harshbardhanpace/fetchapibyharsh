package v2

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
	mockCreateAlert  func(req models.CreateAlertsReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	mockEditAlerts   func(req models.EditAlertsReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	mockGetAlerts    func(req models.GetAlertsReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	mockPauseAlerts  func(req models.PauseAlertsReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	mockDeleteAlerts func(req models.DeleteAlertsReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
)

type alertsMock struct{}

func (m alertsMock) CreateAlert(req models.CreateAlertsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return mockCreateAlert(req, reqH)
}

func (m alertsMock) EditAlerts(req models.EditAlertsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return mockEditAlerts(req, reqH)
}

func (m alertsMock) GetAlerts(req models.GetAlertsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return mockGetAlerts(req, reqH)
}

func (m alertsMock) PauseAlerts(req models.PauseAlertsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return mockPauseAlerts(req, reqH)
}

func (m alertsMock) DeleteAlerts(req models.DeleteAlertsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return mockDeleteAlerts(req, reqH)
}

func TestEditAlerts(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v2/alerts/editAlerts", nil)
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
			EditAlerts(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("EditAlerts() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	var requestPacket models.EditAlertsReq
	requestPacket.ClientId = "abc"
	requestPacket.AlertId = "abc"
	requestPacket.Exchange = "abc"
	requestPacket.InstrumentToken = "abc"
	requestPacket.WaitTime = "abc"
	requestPacket.Condition = "abc"
	requestPacket.Frequency = "abc"
	requestPacket.StateAfterExpiry = "abc"
	requestPacket.UserMessage = "abc"
	requestPacket.Expiry = 123
	var packet []float64
	packet = append(packet, 1.0)
	packet = append(packet, 2.0)
	requestPacket.UserSetValues = packet

	reqPacket, err := json.Marshal(requestPacket)
	if err != nil {
		// Handle the error
	}
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v2/alerts/editAlerts", strings.NewReader(string(reqPacket)))
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
			EditAlerts(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("EditAlerts() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v2/alerts/editAlerts", strings.NewReader(string(reqPacket)))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"
	//Init order provider
	InitAlertsProviderV2(alertsMock{})

	//mock business layer response
	mockEditAlerts = func(req models.EditAlertsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			EditAlerts(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("EditAlerts() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}

}

func TestPauseAlerts(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v2/alerts/pauseAlerts", nil)
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
			PauseAlerts(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("PauseAlerts() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	var requestPacket models.PauseAlertsReq
	requestPacket.ClientId = "abc"
	requestPacket.AlertId = 123
	requestPacket.Status = "abc"

	reqPacket, err := json.Marshal(requestPacket)
	if err != nil {
		// Handle the error
	}
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v2/alerts/pauseAlerts", strings.NewReader(string(reqPacket)))
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
			PauseAlerts(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("PauseAlerts() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v2/alerts/pauseAlerts", strings.NewReader(string(reqPacket)))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"
	//Init order provider
	InitAlertsProviderV2(alertsMock{})

	//mock business layer response
	mockPauseAlerts = func(req models.PauseAlertsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			PauseAlerts(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("PauseAlerts() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}

}

func TestDeleteAlerts(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/api/space/v2/alerts/deleteAlerts", nil)
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
			DeleteAlerts(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("DeleteAlerts() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	var requestPacket models.PauseAlertsReq
	requestPacket.ClientId = "abc"
	requestPacket.AlertId = 123

	reqPacket, err := json.Marshal(requestPacket)
	if err != nil {
		// Handle the error
	}
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v2/alerts/deleteAlerts", strings.NewReader(string(reqPacket)))
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
			DeleteAlerts(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("DeleteAlerts() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v2/alerts/deleteAlerts", strings.NewReader(string(reqPacket)))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"
	//Init order provider
	InitAlertsProviderV2(alertsMock{})

	//mock business layer response
	mockDeleteAlerts = func(req models.DeleteAlertsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			DeleteAlerts(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("DeleteAlerts() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}

}

func TestGetAlerts(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/space/v2/alerts/getAlerts", nil)
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
			GetAlerts(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("GetAlerts() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/space/v2/alerts/getAlerts?clientId=HI009", nil)
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
			GetAlerts(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("GetAlerts() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/space/v2/alerts/getAlerts?clientId=HI009", nil)
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"
	//Init order provider
	InitAlertsProviderV2(alertsMock{})

	//mock business layer response
	mockGetAlerts = func(req models.GetAlertsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			GetAlerts(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("GetAlerts() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}
}
