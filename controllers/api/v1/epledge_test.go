package v1

import (
	"encoding/json"
	"fmt"
	"io"
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
	epledgeRequestMock        func(req models.EpledgeReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	unpledgeRequestMock       func(req models.UnpledgeReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	mtfEpledgeRequest         func(req models.MTFEPledgeRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	getPledgeList             func(reqH models.ReqHeader) (int, apihelpers.APIRes)
	getCTDQuantityList        func(req models.MTFCTDDataReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	getPledgeTransactionsMock func(req models.FetchEpledgeTxnReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	mtfctdmock                func(models.MTFCTDReq, models.ReqHeader) (int, apihelpers.APIRes)
)

type epledgeMock struct{}

func (m epledgeMock) SendEpledgeRequest(req models.EpledgeReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return epledgeRequestMock(req, reqH)
}

func (m epledgeMock) UnpledgeRequest(req models.UnpledgeReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return unpledgeRequestMock(req, reqH)
}

func (m epledgeMock) MTFEpledgeRequest(req models.MTFEPledgeRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return mtfEpledgeRequest(req, reqH)
}

func (m epledgeMock) GetPledgeList(reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return getPledgeList(reqH)
}

func (m epledgeMock) GetCTDQuantityList(req models.MTFCTDDataReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return getCTDQuantityList(req, reqH)
}

func (m epledgeMock) GetPledgeTransactions(req models.FetchEpledgeTxnReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return getPledgeTransactionsMock(req, reqH)
}

func (m epledgeMock) MTFCTD(req models.MTFCTDReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return mtfctdmock(req, reqH)
}

func TestEpledgeRequest(t *testing.T) {
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

	// //1 invalid request
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/edis/epledgeRequest", nil)
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
			EpledgeRequest(tt.args.c)
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("epledge request() = %v, want %v", string(b), expected)
			}
		})
	}

	// //2 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"

	var epledgeReqBody models.EpledgeReq

	epledgeReqBody.Depository = "CDSL"
	epledgeReqBody.ClientID = "9CD12"
	epledgeReqBody.Exchange = "NSE"
	epledgeReqBody.BoId = "12049"
	epledgeReqBody.Segment = "Capital"

	isinBody := models.Isin{
		IsinName: "PNB",
		Isin:     "INE160",
		Quantity: "1",
		Price:    "67.8",
	}

	epledgeReqBody.IsinDetails = append(epledgeReqBody.IsinDetails, isinBody)

	// Encode the epledgeReqBody as a JSON string
	epledgeReqJson, err := json.Marshal(epledgeReqBody)
	if err != nil {
		// Handle the error
	}

	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/edis/epledgeRequest", strings.NewReader(string(epledgeReqJson)))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init epledge provider
	InitEpledgeProvider(epledgeMock{})

	//mock business layer response
	epledgeRequestMock = func(req models.EpledgeReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			EpledgeRequest(tt.args.c)
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("epledge request() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("epledge request() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestGetPledgeTransactions(t *testing.T) {

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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/space/v1/epledge/getPledgeTransactions", nil)
	expected := "{\"status\":false,\"message\":\"INVALID DEVICE TYPE\",\"errorcode\":\"P11034\",\"data\":null}"

	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"Invalid request", args{c: ctx}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetPledgeTransactions(tt.args.c)
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("epledge request() = %v, want %v", string(b), expected)
			}
		})
	}

	// //2 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "web"
	reqHsuccess.ClientId = "9CD12"
	reqHsuccess.RequestId = "test-request-id"

	ctx.Set("reqH", reqHsuccess)
	req := httptest.NewRequest(http.MethodGet, "/api/space/v1/epledge/getPledgeTransactions?clientId=9CD12&page=1&startDate=2022-01-01&endDate=2022-01-31", nil)
	ctx.Request = req

	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init epledge provider
	InitEpledgeProvider(epledgeMock{})

	//mock business layer response
	getPledgeTransactionsMock = func(req models.FetchEpledgeTxnReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			GetPledgeTransactions(tt.args.c)
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("epledge request() = %v, want %v", string(b), expected)
			}
		})
	}
}
