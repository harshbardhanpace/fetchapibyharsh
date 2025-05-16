package v2

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
	searchScripMock func(req models.SearchScripRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	scripInfoMock   func(req models.ScripInfoRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
)

type contractDetailsMock struct{}

func (m contractDetailsMock) SearchScrip(req models.SearchScripRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return searchScripMock(req, reqH)
}

func (m contractDetailsMock) ScripInfo(req models.ScripInfoRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return scripInfoMock(req, reqH)
}

func TestSearchScrip(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/v2/contractdetails/searchScripTL", nil)
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
			SearchScrip(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("search scrip() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/v2/contractdetails/searchScripTL?key=abc", nil)
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
			SearchScrip(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("search scrip() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/v2/contractdetails/searchScripTL?key=abc", nil)
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitContractDetailsProviderV2TL(contractDetailsMock{})

	//mock business layer response
	searchScripMock = func(req models.SearchScripRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			SearchScrip(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("search scrip() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}
}

func TestScripInfo(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/v2/contractdetails/scripInfoTL", nil)
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
			ScripInfo(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("scrip Info() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/v2/contractdetails/scripInfoTL?exchange=NSE&info=scrip&token=11001", nil)
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
			ScripInfo(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("scrip Info() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 bad request
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqH := models.ReqHeader{}
	reqH.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqH)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/v2/contractdetails/scripInfoTL?exchange=NSE&info=scrip", nil)
	expected = "{\"status\":false,\"message\":\"INVALID REQUEST\",\"errorcode\":\"P11017\",\"data\":null}"
	tests = []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"Bad  request", args{c: ctx}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ScripInfo(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("scrip Info() = %v, want %v", string(b), expected)
			}
		})
	}

	//4 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/v2/contractdetails/scripInfoTL?exchange=NSE&info=scrip&token=11001", nil)
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitContractDetailsProviderV2TL(contractDetailsMock{})

	//mock business layer response
	scripInfoMock = func(req models.ScripInfoRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			ScripInfo(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("scrip Info() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}
}
