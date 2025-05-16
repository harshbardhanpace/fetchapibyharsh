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
	edisRequestMock  func(req models.EdisReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	generateTpinMock func(req models.TpinReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
)

type edisMock struct{}

func (m edisMock) SendEdisRequest(req models.EdisReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return edisRequestMock(req, reqH)
}

func (m edisMock) GenerateTpin(req models.TpinReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return generateTpinMock(req, reqH)
}

func TestEdisRequest(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/edis/edisRequest", nil)
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
			EdisRequest(tt.args.c)
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("edis request() = %v, want %v", string(b), expected)
			}
		})
	}

	// //2 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"

	var edisReqBody models.EdisReq
	edisReqBody.ClientID = "9CD12"

	instrumentBody := models.Instrument{
		InstrumentToken: 22,
		Exchange:        "NSE",
		Total:           40,
		Authorized:      1,
	}

	edisReqBody.Instruments = append(edisReqBody.Instruments, instrumentBody)
	edisReqBody.RequestType = "REGULAR"

	// Encode the edisReqBody as a JSON string
	edisReqJSON, err := json.Marshal(edisReqBody)
	if err != nil {
		// Handle the error
	}

	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/edis/edisRequest", strings.NewReader(string(edisReqJSON)))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init edis provider
	InitEdisProvider(edisMock{})

	//mock business layer response
	edisRequestMock = func(req models.EdisReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			EdisRequest(tt.args.c)
			b, _ := io.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("edis request() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("edis request() = %v, want %v", w.Code, 400)
			}
		})
	}

}
