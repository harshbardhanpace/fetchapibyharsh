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
	placeBOOrderMock      func(req models.PlaceBOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	modifyBOOrderMock     func(req models.ModifyBOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	cancelBOOrderMock     func(req models.ExitBOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	placeCOOrderMock      func(req models.PlaceCOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	modifyCOOrderMock     func(req models.ModifyCOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	cancelCOOrderMock     func(req models.ExitCOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	placeSpreadOrderMock  func(req models.PlaceSpreadOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	modifySpreadOrderMock func(req models.ModifySpreadOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	cancelSpreadOrderMock func(req models.ExitSpreadOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
)

type conditionalOrderMock struct{}

func (m conditionalOrderMock) PlaceBOOrder(req models.PlaceBOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return placeBOOrderMock(req, reqH)
}

func (m conditionalOrderMock) ModifyBOOrder(req models.ModifyBOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return modifyBOOrderMock(req, reqH)
}

func (m conditionalOrderMock) CancelBOOrder(req models.ExitBOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return cancelBOOrderMock(req, reqH)
}

func (m conditionalOrderMock) PlaceCOOrder(req models.PlaceCOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return placeCOOrderMock(req, reqH)
}

func (m conditionalOrderMock) ModifyCOOrder(req models.ModifyCOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return modifyCOOrderMock(req, reqH)
}

func (m conditionalOrderMock) CancelCOOrder(req models.ExitCOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return cancelCOOrderMock(req, reqH)
}

func (m conditionalOrderMock) PlaceSpreadOrder(req models.PlaceSpreadOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return placeSpreadOrderMock(req, reqH)
}

func (m conditionalOrderMock) ModifySpreadOrder(req models.ModifySpreadOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return modifySpreadOrderMock(req, reqH)
}

func (m conditionalOrderMock) CancelSpreadOrder(req models.ExitSpreadOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return cancelSpreadOrderMock(req, reqH)
}

func TestPlaceBOOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeBOOrder", nil)
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
			PlaceBOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional place bo order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeBOOrder", strings.NewReader("{\"clientId\":\"CLient1\",\"disclosedQuantity\":0,\"exchange\":\"DSE\",\"instrumentToken\":\"10666\",\"isTrailing\": true,\"orderSide\":\"BUY\",\"orderType\":\"LIMIT\",\"price\":1000,\"product\":\"CNC\",\"quantity\":10,\"squareOffValue\":1,\"stopLossValue\":1,\"trailingStopLoss\":\"0.05\",\"triggerPrice\":0,\"userOrderId\":1002,\"validity\":\"DAY\"}"))
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
			PlaceBOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional place bo order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 bad request
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqH := models.ReqHeader{}
	reqH.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqH)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeBOOrder", strings.NewReader("{\"clientId\":\"CLient1\",\"disclosedQuantity\":0,\"exchange\":\"DSE\",\"instrumentToken\":\"10666\",\"isTrailing\": true,\"orderSide\":\"BUY\",\"orderType\":\"LIMIT\",\"price\":1000,\"product\":\"CNC\",\"quantity\":10,\"squareOffValue\":1,\"stopLossValue\":1,\"trailingStopLoss\":\"0.05\",\"triggerPrice\":0,\"userOrderId\":1002,\"validity\":\"DAY\"}"))
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
			PlaceBOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional place bo order() = %v, want %v", string(b), expected)
			}
		})
	}

	//4 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeBOOrder", strings.NewReader("{\"clientId\":\"CLient1\",\"disclosedQuantity\":0,\"exchange\":\"NSE\",\"instrumentToken\":\"10666\",\"isTrailing\": true,\"orderSide\":\"BUY\",\"orderType\":\"LIMIT\",\"price\":1000,\"product\":\"CNC\",\"quantity\":10,\"squareOffValue\":1,\"stopLossValue\":1,\"trailingStopLoss\":\"0.05\",\"triggerPrice\":0,\"userOrderId\":1002,\"validity\":\"DAY\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitConditionalOrderProvider(conditionalOrderMock{})

	//mock business layer response
	placeBOOrderMock = func(req models.PlaceBOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			PlaceBOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional place bo order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("conditional place bo order() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestModifyBOOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyBOOrder", nil)
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
			ModifyBOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional modify bo order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyBOOrder", strings.NewReader("{\"clientId\":\"CLient1\",\"disclosedQuantity\":2,\"exchange\":\"DSE\",\"exchangeOrderId\": \"ionwdg123\",\"filledQuantity\": 1,\"instrumentToken\":\"10666\",\"isTrailing\": true,\"lastActivityReference\": 1325938440097498600,\"omsOrderId\": \"123445\",\"orderType\": \"LIMIT\",\"price\":34.2,\"product\":\"CNC\",\"quantity\":10,\"remainingQuantity\": 0,\"squareOffValue\":1,\"stopLossValue\":1,\"trailingStopLoss\":41,\"triggerPrice\":0,\"validity\":\"DAY\"}"))
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
			ModifyBOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional modify bo order() = %v, want %v", string(b), expected)
			}
		})
	}

	// 3 bad request
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqH := models.ReqHeader{}
	reqH.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqH)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyBOOrder", strings.NewReader("{\"clientId\":\"CLient1\",\"disclosedQuantity\":2,\"exchange\":\"DSE\",\"exchangeOrderId\": \"ionwdg123\",\"filledQuantity\": 1,\"instrumentToken\":\"10666\",\"isTrailing\": true,\"lastActivityReference\": 1325938440097498600,\"omsOrderId\": \"123445\",\"orderType\": \"LIMIT\",\"price\":34.2,\"product\":\"CNC\",\"quantity\":10,\"remainingQuantity\": 0,\"squareOffValue\":1,\"stopLossValue\":1,\"trailingStopLoss\":41,\"triggerPrice\":0,\"validity\":\"DAY\"}"))
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
			ModifyBOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional modify bo order() = %v, want %v", string(b), expected)
			}
		})
	}

	//4 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyBOOrder", strings.NewReader("{\"clientId\":\"CLient1\",\"disclosedQuantity\":2,\"exchange\":\"NSE\",\"exchangeOrderId\": \"ionwdg123\",\"filledQuantity\": 1,\"instrumentToken\":\"10666\",\"isTrailing\": true,\"lastActivityReference\": 1325938440097498600,\"omsOrderId\": \"123445\",\"orderType\": \"LIMIT\",\"price\":34.2,\"product\":\"CNC\",\"quantity\":10,\"remainingQuantity\": 0,\"squareOffValue\":1,\"stopLossValue\":1,\"trailingStopLoss\":41,\"triggerPrice\":0,\"validity\":\"DAY\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitConditionalOrderProvider(conditionalOrderMock{})

	//mock business layer response
	modifyBOOrderMock = func(req models.ModifyBOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			ModifyBOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional modify bo order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestExitBOOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/exitBOOrder", nil)
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
			ExitBOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional exit bo order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/exitBOOrder", strings.NewReader("{\"clientId\":\"CLient1\",\"exchangeOrderId\": \"ionwdg123\",\"legOrderIndicator\":\"fdsfd\",\"omsOrderId\": \"123445\",\"Status\": \"CONFIRMED\"}"))
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
			ExitBOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional exit bo order() = %v, want %v", string(b), expected)
			}
		})
	}

	//4 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/exitBOOrder", strings.NewReader("{\"clientId\":\"CLient1\",\"exchangeOrderId\": \"ionwdg123\",\"legOrderIndicator\":\"fdsfd\",\"omsOrderId\": \"123445\",\"Status\": \"CONFIRMED\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitConditionalOrderProvider(conditionalOrderMock{})

	//mock business layer response
	cancelBOOrderMock = func(req models.ExitBOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			ExitBOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional exit bo order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestPlaceCOOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeCOOrder", nil)
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
			PlaceCOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional place co order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeCOOrder", strings.NewReader("{\"clientId\": \"CLient1\",\"disclosedQuantity\": 2,\"exchange\": \"DSE\",\"instrumentToken\": \"22\",\"orderSide\": \"BUY\",\"orderType\": \"LIMIT\",\"price\": 34.2,\"product\": \"CNC\",\"quantity\": 2,\"stopLossValue\": 33,\"trailingStopLoss\": 33,\"userOrderId\": 91261928,\"validity\": \"DAY\"}"))
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
			PlaceCOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional place co order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 bad request
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqH := models.ReqHeader{}
	reqH.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqH)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeCOOrder", strings.NewReader("{\"clientId\": \"CLient1\",\"disclosedQuantity\": 2,\"exchange\": \"DSE\",\"instrumentToken\": \"22\",\"orderSide\": \"BUY\",\"orderType\": \"LIMIT\",\"price\": 34.2,\"product\": \"CNC\",\"quantity\": 2,\"stopLossValue\": 33,\"trailingStopLoss\": 33,\"userOrderId\": 91261928,\"validity\": \"DAY\"}"))
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
			PlaceCOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional place co order() = %v, want %v", string(b), expected)
			}
		})
	}

	// //4 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeCOOrder", strings.NewReader("{\"clientId\": \"CLient1\",\"disclosedQuantity\": 2,\"exchange\": \"NSE\",\"instrumentToken\": \"22\",\"orderSide\": \"BUY\",\"orderType\": \"LIMIT\",\"price\": 34.2,\"product\": \"CNC\",\"quantity\": 2,\"stopLossValue\": 33,\"trailingStopLoss\": 33,\"userOrderId\": 91261928,\"validity\": \"DAY\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitConditionalOrderProvider(conditionalOrderMock{})

	//mock business layer response
	placeCOOrderMock = func(req models.PlaceCOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			PlaceCOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional place co order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestModifyCOOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyCOOrder", nil)
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
			ModifyCOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional modify Co order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyCOOrder", strings.NewReader("{\"clientId\": \"Client1\",\"disclosedQuantity\": 2,\"exchange\": \"DSE\",\"exchangeOrderId\": \"ionwdg123\",\"filledQuantity\": 1,\"instrumentToken\": \"22\",\"lastActivityReference\": 1325938440097498600,\"omsOrderId\": \"123445\",\"orderType\": \"LIMIT\",\"price\": 34.2,\"product\": \"CNC\",\"quantity\": 2,\"remainingQuantity\": 0,\"stopLossValue\": 0,\"trailingStopLoss\": 33,\"validity\": \"DAY\"}"))
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
			ModifyCOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional modify co order() = %v, want %v", string(b), expected)
			}
		})
	}

	// 3 bad request
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqH := models.ReqHeader{}
	reqH.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqH)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyCOOrder", strings.NewReader("{\"clientId\": \"Client1\",\"disclosedQuantity\": 2,\"exchange\": \"DSE\",\"exchangeOrderId\": \"ionwdg123\",\"filledQuantity\": 1,\"instrumentToken\": \"22\",\"lastActivityReference\": 1325938440097498600,\"omsOrderId\": \"123445\",\"orderType\": \"LIMIT\",\"price\": 34.2,\"product\": \"CNC\",\"quantity\": 2,\"remainingQuantity\": 0,\"stopLossValue\": 0,\"trailingStopLoss\": 33,\"validity\": \"DAY\"}"))
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
			ModifyCOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional modify co order() = %v, want %v", string(b), expected)
			}
		})
	}

	//4 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyCOOrder", strings.NewReader("{\"clientId\": \"Client1\",\"disclosedQuantity\": 2,\"exchange\": \"NSE\",\"exchangeOrderId\": \"ionwdg123\",\"filledQuantity\": 1,\"instrumentToken\": \"22\",\"lastActivityReference\": 1325938440097498600,\"omsOrderId\": \"123445\",\"orderType\": \"LIMIT\",\"price\": 34.2,\"product\": \"CNC\",\"quantity\": 2,\"remainingQuantity\": 0,\"stopLossValue\": 0,\"trailingStopLoss\": 33,\"validity\": \"DAY\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitConditionalOrderProvider(conditionalOrderMock{})

	//mock business layer response
	modifyCOOrderMock = func(req models.ModifyCOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			ModifyCOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional modify co order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestExitCOOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/exitCOOrder", nil)
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
			ExitCOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional exit co order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/exitCOOrder", strings.NewReader("{\"clientId\": \"CLIENT1\",\"exchangeOrderId\": \"ionwdg123\",\"legOrderIndicator\": \"Entry\",\"omsOrderId\": \"123445\"}"))
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
			ExitCOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional exit co order() = %v, want %v", string(b), expected)
			}
		})
	}

	//4 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/exitCOOrder", strings.NewReader("{\"clientId\": \"CLIENT1\",\"exchangeOrderId\": \"ionwdg123\",\"legOrderIndicator\": \"Entry\",\"omsOrderId\": \"123445\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitConditionalOrderProvider(conditionalOrderMock{})

	//mock business layer response
	cancelCOOrderMock = func(req models.ExitCOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			ExitCOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional exit co order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestPlaceSpreadOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeSpreadOrder", nil)
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
			PlaceSpreadOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional place spread order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeSpreadOrder", strings.NewReader("{\"clientId\": \"CLIENT1\",\"disclosedQuantity\": 2,\"exchange\": \"DSE\",\"instrumentToken\": \"22\",\"orderSide\": \"BUY\",\"orderType\": \"LIMIT\",\"price\": 34.2,\"product\": \"CNC\",\"quantity\": 2,\"userOrderId\": 91261928,\"validity\": \"DAY\"}"))
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
			PlaceSpreadOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional place spread order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 bad request
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqH := models.ReqHeader{}
	reqH.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqH)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeSpreadOrder", strings.NewReader("{\"clientId\": \"CLIENT1\",\"disclosedQuantity\": 2,\"exchange\": \"DSE\",\"instrumentToken\": \"22\",\"orderSide\": \"BUY\",\"orderType\": \"LIMIT\",\"price\": 34.2,\"product\": \"CNC\",\"quantity\": 2,\"userOrderId\": 91261928,\"validity\": \"DAY\"}"))
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
			PlaceSpreadOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional place spread order() = %v, want %v", string(b), expected)
			}
		})
	}

	//4 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeSpreadOrder", strings.NewReader("{\"clientId\": \"CLIENT1\",\"disclosedQuantity\": 2,\"exchange\": \"NSE\",\"instrumentToken\": \"22\",\"orderSide\": \"BUY\",\"orderType\": \"LIMIT\",\"price\": 34.2,\"product\": \"CNC\",\"quantity\": 2,\"userOrderId\": 91261928,\"validity\": \"DAY\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitConditionalOrderProvider(conditionalOrderMock{})

	//mock business layer response
	placeSpreadOrderMock = func(req models.PlaceSpreadOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			PlaceSpreadOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional place spread order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestModifySpreadOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifySpreadOrder", nil)
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
			ModifySpreadOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional modify spread order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifySpreadOrder", strings.NewReader("{\"clientId\": \"CLIENT1\",\"disclosedQuantity\": 0,\"exchange\": \"DSE\",\"exchangeOrderId\": \"kjh\",\"instrumentToken\": \"22\",\"isTrailing\": true,\"omsOrderId\": \"123445\",\"orderType\": \"LIMIT\",\"price\": 34.2,\"prodType\": \"fdsf\",\"product\": \"CNC\",\"quantity\": 2,\"squareOffValue\": 40,\"stopLossValue\": 39,\"trailingStopLoss\": 33,\"triggerPrice\": 38,\"validity\": \"DAY\"}"))
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
			ModifySpreadOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional modify spread order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 bad request
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqH := models.ReqHeader{}
	reqH.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqH)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifySpreadOrder", strings.NewReader("{\"clientId\": \"CLIENT1\",\"disclosedQuantity\": 0,\"exchange\": \"DSE\",\"exchangeOrderId\": \"kjh\",\"instrumentToken\": \"22\",\"isTrailing\": true,\"omsOrderId\": \"123445\",\"orderType\": \"LIMIT\",\"price\": 34.2,\"prodType\": \"fdsf\",\"product\": \"CNC\",\"quantity\": 2,\"squareOffValue\": 40,\"stopLossValue\": 39,\"trailingStopLoss\": 33,\"triggerPrice\": 38,\"validity\": \"DAY\"}"))
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
			ModifySpreadOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional modify spread order() = %v, want %v", string(b), expected)
			}
		})
	}

	//4 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifySpreadOrder", strings.NewReader("{\"clientId\": \"CLIENT1\",\"disclosedQuantity\": 0,\"exchange\": \"NSE\",\"exchangeOrderId\": \"kjh\",\"instrumentToken\": \"22\",\"isTrailing\": true,\"omsOrderId\": \"123445\",\"orderType\": \"LIMIT\",\"price\": 34.2,\"prodType\": \"fdsf\",\"product\": \"CNC\",\"quantity\": 2,\"squareOffValue\": 40,\"stopLossValue\": 39,\"trailingStopLoss\": 33,\"triggerPrice\": 38,\"validity\": \"DAY\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitConditionalOrderProvider(conditionalOrderMock{})

	//mock business layer response
	modifySpreadOrderMock = func(req models.ModifySpreadOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			ModifySpreadOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional modify spread order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestExitSpreadOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/exitSpreadOrder", nil)
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
			ExitSpreadOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional exit spread order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/exitSpreadOrder", strings.NewReader("{\"clientId\": \"CLIENT1\",\"exchangeOrderId\": \"ionwdg123\",\"legOrderIndicator\": \"Entry\",\"omsOrderId\": \"123445\"}"))
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
			ExitSpreadOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional exit spread order() = %v, want %v", string(b), expected)
			}
		})
	}

	//4 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/exitSpreadOrder", strings.NewReader("{\"clientId\": \"CLIENT1\",\"exchangeOrderId\": \"ionwdg123\",\"legOrderIndicator\": \"Entry\",\"omsOrderId\": \"123445\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitConditionalOrderProvider(conditionalOrderMock{})

	//mock business layer response
	cancelSpreadOrderMock = func(req models.ExitSpreadOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			ExitSpreadOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("conditional exit spread order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}
