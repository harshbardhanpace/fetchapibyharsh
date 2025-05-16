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
	fetchDematHoldingsMock func(req models.FetchDematHoldingsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	convertPositionsMock   func(req models.ConvertPositionsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	getPositionsMock       func(req models.GetPositionRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
)

type portfolioMock struct{}

func (m portfolioMock) FetchDematHoldings(req models.FetchDematHoldingsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return fetchDematHoldingsMock(req, reqH)
}

func (m portfolioMock) ConvertPositions(req models.ConvertPositionsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return convertPositionsMock(req, reqH)
}

func (m portfolioMock) GetPositions(req models.GetPositionRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return getPositionsMock(req, reqH)
}

func TestFetchDematHoldings(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/space/v2/portfolioapis/fetchDematHoldings", nil)
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
			FetchDematHoldings(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("FetchDematHoldings() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/space/v2/portfolioapis/fetchDematHoldings?clientId=HI009", nil)
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
			FetchDematHoldings(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("FetchDematHoldings() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/space/v2/portfolioapis/fetchDematHoldings?clientId=HI009", nil)
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"
	//Init order provider
	InitPortfolioProviderV2(portfolioMock{})

	//mock business layer response
	fetchDematHoldingsMock = func(req models.FetchDematHoldingsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			FetchDematHoldings(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("FetchDematHoldings() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}
}

func TestConvertPositions(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v2/portfolioapis/convertPositions", nil)
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
			ConvertPositions(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("ConvertPositions() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	var requestPacket models.ConvertPositionsRequest
	requestPacket.ClientID = "abc"
	requestPacket.Exchange = "NSE"
	requestPacket.InstrumentToken = 22
	requestPacket.Product = "CNC"
	requestPacket.NewProduct = "MIS"
	requestPacket.Quantity = 1
	requestPacket.Validity = "DAY"
	requestPacket.OrderSide = "BUY"

	reqPacket, err := json.Marshal(requestPacket)
	if err != nil {
		// Handle the error
	}
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v2/portfolioapis/convertPositions", strings.NewReader(string(reqPacket)))
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
			ConvertPositions(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("ConvertPositions() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/space/v2/portfolioapis/convertPositions", strings.NewReader(string(reqPacket)))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"
	//Init order provider
	InitPortfolioProviderV2(portfolioMock{})

	//mock business layer response
	convertPositionsMock = func(req models.ConvertPositionsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			ConvertPositions(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("ConvertPositions() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}

}

func TestGetPositions(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/space/v2/portfolioapis/getPositions", nil)
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
			GetPositions(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("GetPositions() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/space/v2/portfolioapis/getPositions?clientId=HI009&type=live", nil)
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
			GetPositions(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("GetPositions() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/space/v2/portfolioapis/getPositions?clientId=HI009&type=live", nil)
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"
	//Init order provider
	InitPortfolioProviderV2(portfolioMock{})

	//mock business layer response
	getPositionsMock = func(req models.GetPositionRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			GetPositions(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("GetPositions() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}
}
