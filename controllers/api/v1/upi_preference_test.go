package v1

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
	setUpiPreferenceMock     func(req models.SetUpiPreferenceReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	fetchUpiPreferenceeMock  func(req models.FetchUpiPreferenceReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	deleteUpiPreferenceeMock func(req models.DeleteUpiPreferenceReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
)

type upiPreferenceMock struct{}

func (m upiPreferenceMock) SetUpiPreference(req models.SetUpiPreferenceReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return setUpiPreferenceMock(req, reqH)
}

func (m upiPreferenceMock) FetchUpiPreference(req models.FetchUpiPreferenceReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return fetchUpiPreferenceeMock(req, reqH)
}

func (m upiPreferenceMock) DeleteUpiPreference(req models.DeleteUpiPreferenceReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return deleteUpiPreferenceeMock(req, reqH)
}

func TestSetUpiPreference(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/space/v1/upi/setUpiPreference", nil)
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
			SetUpiPreference(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("SetUpiPreference() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/space/v1/upi/setUpiPreference", strings.NewReader("{\"clientId\":\"CLient1\", \"upiId\":\"12345\"}"))
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
			SetUpiPreference(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("SetUpiPreference() = %v, want %v", string(b), expected)
			}
		})
	}

	//4 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/space/v1/upi/setUpiPreference", strings.NewReader("{\"clientId\":\"CLient1\", \"upiId\":\"12345\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitUpiPreferenceProvider(upiPreferenceMock{})

	//mock business layer response
	setUpiPreferenceMock = func(req models.SetUpiPreferenceReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			SetUpiPreference(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("SetUpiPreference() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("SetUpiPreference() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestFetchUpiPreference(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/space/v1/upi/fetchUpiPreference", nil)
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
			FetchUpiPreference(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("FetchUpiPreference() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/space/v1/upi/fetchUpiPreference?clientId=CLIENT1", nil)
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
			FetchUpiPreference(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("FetchUpiPreference() = %v, want %v", string(b), expected)
			}
		})
	}

	//4 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/space/v1/upi/fetchUpiPreference?clientId=CLIENT1", nil)
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitUpiPreferenceProvider(upiPreferenceMock{})

	//mock business layer response
	fetchUpiPreferenceeMock = func(req models.FetchUpiPreferenceReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			FetchUpiPreference(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("FetchUpiPreference() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("FetchUpiPreference() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestDeleteUpiPreference(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/api/space/v1/upi/deleteUpiPreference", nil)
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
			DeleteUpiPreference(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("DeleteUpiPreference() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	var deleteUpiReq models.DeleteUpiPreferenceReq
	deleteUpiReq.ClientId = "CLIENT1"
	deleteUpiReq.UpiIds = []string{"1234", "5678"}
	deleteUpiReqMarshall, _ := json.Marshal(deleteUpiReq)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/api/space/v1/upi/deleteUpiPreference", strings.NewReader(string(deleteUpiReqMarshall)))
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
			DeleteUpiPreference(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("DeleteUpiPreference() = %v, want %v", string(b), expected)
			}
		})
	}

	//4 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)

	ctx.Request = httptest.NewRequest(http.MethodDelete, "/api/space/v1/upi/deleteUpiPreference", strings.NewReader(string(deleteUpiReqMarshall)))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitUpiPreferenceProvider(upiPreferenceMock{})

	//mock business layer response
	deleteUpiPreferenceeMock = func(req models.DeleteUpiPreferenceReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			DeleteUpiPreference(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("DeleteUpiPreference() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("DeleteUpiPreference() = %v, want %v", w.Code, 400)
			}
		})
	}

}
