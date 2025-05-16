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
	brokerChargesMock        func(req models.BrokerChargesReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	combineBrokerChargesMock func(req models.CombineBrokerChargesReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	fundsPayoutMock          func(req models.FundsPayoutReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
)

type chargesMock struct{}

func (m chargesMock) BrokerCharges(req models.BrokerChargesReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return brokerChargesMock(req, reqH)
}

func (m chargesMock) CombineBrokerCharges(req models.CombineBrokerChargesReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return combineBrokerChargesMock(req, reqH)
}

func (m chargesMock) FundsPayout(req models.FundsPayoutReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return fundsPayoutMock(req, reqH)
}

func TestFundsPayout(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/charges/fundsPayout", nil)
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
			FundsPayout(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("FundsPayout() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/charges/fundsPayout", strings.NewReader("{\"clientId\":\"CLient1\"}"))
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
			FundsPayout(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("FundsPayout() = %v, want %v", string(b), expected)
			}
		})
	}

	//4 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/charges/fundsPayout", strings.NewReader("{\"clientId\":\"CLient1\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitChargesProvider(chargesMock{})

	//mock business layer response
	fundsPayoutMock = func(req models.FundsPayoutReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			FundsPayout(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("FundsPayout() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("FundsPayout() = %v, want %v", w.Code, 400)
			}
		})
	}

}
