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
	placeOrderMock         func(req models.PlaceOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	modifyOrderMock        func(req models.ModifyOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	cancelOrderMock        func(req models.CancelOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	placeAMOOrderMock      func(req models.PlaceOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	modifyAMOOrderMock     func(req models.ModifyAMORequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	cancelAMOOrderMock     func(req models.CancelOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	pendingOrderMock       func(req models.PendingOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	completedOrderMock     func(req models.CompletedOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	tradeBookMock          func(req models.TradeBookRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	orderHistoryMock       func(req models.OrderHistoryRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	placeGTTOrderMock      func(req models.CreateGTTOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	modifyGTTOrderMock     func(req models.ModifyGTTOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	cancelGTTOrderMock     func(req models.CancelGTTOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	fetchGTTOrderMock      func(req models.FetchGTTOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	placeGttOCOOrderMock   func(req models.CreateGttOCORequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	marginCalculationsMock func(req models.MarginCalculationRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	lastTradedPriceMock    func(req models.LastTradedPriceRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	placeIcebergOrder      func(req models.IcebergOrderReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	modifyIcebergOrder     func(req models.ModifyIcebergOrderReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	cancelIcebergOrder     func(req models.CancelIcebergOrderReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
)

type orderMock struct{}

func (m orderMock) PlaceOrder(req models.PlaceOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return placeOrderMock(req, reqH)
}

func (m orderMock) ModifyOrder(req models.ModifyOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return modifyOrderMock(req, reqH)
}

func (m orderMock) CancelOrder(req models.CancelOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return cancelOrderMock(req, reqH)
}

func (m orderMock) PlaceAMOOrder(req models.PlaceOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return placeAMOOrderMock(req, reqH)
}

func (m orderMock) ModifyAMOOrder(req models.ModifyAMORequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return modifyAMOOrderMock(req, reqH)
}

func (m orderMock) CancelAMOOrder(req models.CancelOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return cancelAMOOrderMock(req, reqH)
}

func (m orderMock) PendingOrder(req models.PendingOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return pendingOrderMock(req, reqH)
}

func (m orderMock) CompletedOrder(req models.CompletedOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return completedOrderMock(req, reqH)
}

func (m orderMock) TradeBook(req models.TradeBookRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return tradeBookMock(req, reqH)
}

func (m orderMock) OrderHistory(req models.OrderHistoryRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return orderHistoryMock(req, reqH)
}

func (m orderMock) PlaceGTTOrder(req models.CreateGTTOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return placeGTTOrderMock(req, reqH)
}

func (m orderMock) ModifyGTTOrder(req models.ModifyGTTOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return modifyGTTOrderMock(req, reqH)
}

func (m orderMock) CancelGTTOrder(req models.CancelGTTOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return cancelGTTOrderMock(req, reqH)
}

func (m orderMock) FetchGTTOrder(req models.FetchGTTOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return fetchGTTOrderMock(req, reqH)
}

func (m orderMock) PlaceGttOCOOrder(req models.CreateGttOCORequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return placeGttOCOOrderMock(req, reqH)
}

func (m orderMock) MarginCalculations(req models.MarginCalculationRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return marginCalculationsMock(req, reqH)
}

func (m orderMock) LastTradedPrice(req models.LastTradedPriceRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return lastTradedPriceMock(req, reqH)
}

func (m orderMock) PlaceIcebergOrder(req models.IcebergOrderReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return placeIcebergOrder(req, reqH)
}
func (m orderMock) ModifyIcebergOrder(req models.ModifyIcebergOrderReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return modifyIcebergOrder(req, reqH)
}
func (m orderMock) CancelIcebergOrder(req models.CancelIcebergOrderReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return cancelIcebergOrder(req, reqH)
}

func TestPlaceOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeOrder", nil)
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
			PlaceOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("place order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeOrder", strings.NewReader("{\"clientId\":\"CLient1\",\"disclosedQuantity\":0,\"exchange\":\"NSE\",\"executionType\":\"REGULAR\",\"instrumentToken\":\"10666\",\"orderSide\":\"BUY\",\"orderType\":\"LIMIT\",\"price\":1000,\"product\":\"CNC\",\"quantity\":10,\"triggerPrice\":0,\"validity\":\"DAY\"}"))
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
			PlaceOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("place order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 bad request
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqH := models.ReqHeader{}
	reqH.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqH)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeOrder", strings.NewReader("{\"clientId\":\"CLient1\",\"disclosedQuantity\":0,\"exchange\":\"DSE\",\"executionType\":\"REGULAR\",\"instrumentToken\":\"10666\",\"orderSide\":\"BUY\",\"orderType\":\"LIMIT\",\"price\":1000,\"product\":\"CNC\",\"quantity\":10,\"triggerPrice\":0,\"validity\":\"DAY\"}"))
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
			PlaceOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("place order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeOrder", strings.NewReader("{\"clientId\":\"CLient1\",\"disclosedQuantity\":0,\"exchange\":\"BSE\",\"executionType\":\"REGULAR\",\"instrumentToken\":\"10666\",\"orderSide\":\"BUY\",\"orderType\":\"LIMIT\",\"price\":1000,\"product\":\"CNC\",\"quantity\":10,\"triggerPrice\":0,\"validity\":\"DAY\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitOrderProvider(orderMock{})

	//mock business layer response
	placeOrderMock = func(req models.PlaceOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			PlaceOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("place order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestModifyOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyOrder", nil)
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
			ModifyOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("modify order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyOrder", strings.NewReader("{\"clientId\":\"Client1\",\"instrumentToken\":\"22\",\"disclosedQuantity\":0,\"exchange\":\"NSE\",\"executionType\":\"REGULAR\",\"orderType\":\"LIMIT\",\"price\":1000,\"product\":\"CNC\",\"quantity\":10,\"triggerPrice\":0,\"omsOrderID\":\"20220920-39\",\"validity\":\"DAY\"}"))
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
			ModifyOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("modify order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 bad request
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqH := models.ReqHeader{}
	reqH.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqH)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyOrder", strings.NewReader("{\"clientId\":\"Client1\",\"instrumentToken\":\"22\",\"disclosedQuantity\":0,\"exchange\":\"dSE\",\"executionType\":\"REGULAR\",\"orderType\":\"LIMIT\",\"price\":1000,\"product\":\"CNC\",\"quantity\":10,\"triggerPrice\":0,\"omsOrderID\":\"20220920-39\",\"validity\":\"DAY\"}"))
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
			ModifyOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("modify order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyOrder", strings.NewReader("{\"clientId\":\"Client1\",\"instrumentToken\":\"22\",\"disclosedQuantity\":0,\"exchange\":\"NSE\",\"executionType\":\"REGULAR\",\"orderType\":\"LIMIT\",\"price\":1000,\"product\":\"CNC\",\"quantity\":10,\"triggerPrice\":0,\"omsOrderID\":\"20220920-39\",\"validity\":\"DAY\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitOrderProvider(orderMock{})

	//mock business layer response
	modifyOrderMock = func(req models.ModifyOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			ModifyOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("modify order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestCancelOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/cancelOrder", nil)
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
			CancelOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Cancel order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/cancelOrder", strings.NewReader("{\"clientId\":\"Client1\",\"executionType\":\"REGULAR\",\"omsOrderID\":\"20220920-39\"}"))
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
			CancelOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Cancel order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 bad request
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqH := models.ReqHeader{}
	reqH.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqH)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/cancelOrder", strings.NewReader("{\"clientId\":\"Client1\",\"executionType\":\"NO\",\"omsOrderID\":\"20220920-39\"}"))
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
			CancelOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Cancel order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/cancelOrder", strings.NewReader("{\"clientId\":\"Client1\",\"executionType\":\"REGULAR\",\"omsOrderID\":\"20220920-39\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitOrderProvider(orderMock{})

	//mock business layer response
	cancelOrderMock = func(req models.CancelOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			CancelOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Cancel order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestPlaceAMOOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeAMOOrder", nil)
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
			PlaceAMOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("place amo order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeAMOOrder", strings.NewReader("{\"clientId\":\"Client1\",\"disclosedQuantity\":0,\"exchange\":\"NSE\",\"executionType\":\"AMO\",\"instrumentToken\":\"10666\",\"orderSide\":\"BUY\",\"orderType\":\"LIMIT\",\"price\":1000,\"product\":\"CNC\",\"quantity\":10,\"triggerPrice\":0,\"validity\":\"DAY\"}"))
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
			PlaceAMOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("place amo order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 bad request
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqH := models.ReqHeader{}
	reqH.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqH)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeAMOOrder", strings.NewReader("{\"clientId\":\"CLient1\",\"disclosedQuantity\":0,\"exchange\":\"DSE\",\"executionType\":\"AMO\",\"instrumentToken\":\"10666\",\"orderSide\":\"BUY\",\"orderType\":\"LIMIT\",\"price\":1000,\"product\":\"CNC\",\"quantity\":10,\"triggerPrice\":0,\"validity\":\"DAY\"}"))
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
			PlaceAMOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("place amo order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeAMOOrder", strings.NewReader("{\"exchange\":\"NSE\",\"instrumentToken\":\"10666\",\"clientId\":\"Client1\",\"disclosedQuantity\":0,\"executionType\":\"AMO\",\"orderSide\":\"BUY\",\"orderType\":\"LIMIT\",\"price\":1000,\"product\":\"CNC\",\"quantity\":10,\"triggerPrice\":0,\"validity\":\"DAY\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitOrderProvider(orderMock{})

	//mock business layer response
	placeAMOOrderMock = func(req models.PlaceOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			PlaceAMOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("place AMO order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestModifyAMOOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyAMOOrder", nil)
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
			ModifyAMOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("modify amo order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyAMOOrder", strings.NewReader("{\"clientId\":\"Client1\",\"disclosedQuantity\":0,\"exchange\":\"NSE\",\"executionType\":\"AMO\",\"filledQuantity\":0,\"instrumentToken\":\"1330\",\"lastActivityReference\":0,\"omsOrderID\":\"20220913-31\",\"orderType\":\"LIMIT\",\"price\":2000,\"product\":\"CNC\",\"quantity\":1,\"triggerPrice\":0,\"remainingQuantity\":0,\"exchangeOrderId\":\"12345\",\"validity\":\"DAY\"}"))
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
			ModifyAMOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("modify amo order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 bad request
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqH := models.ReqHeader{}
	reqH.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqH)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyAMOOrder", strings.NewReader("{\"clientId\":\"Client1\",\"disclosedQuantity\":0,\"exchange\":\"NO\",\"executionType\":\"AMO\",\"filledQuantity\":0,\"instrumentToken\":\"1330\",\"lastActivityReference\":0,\"omsOrderID\":\"20220913-31\",\"orderType\":\"LIMIT\",\"price\":2000,\"product\":\"CNC\",\"quantity\":1,\"triggerPrice\":0,\"remainingQuantity\":0,\"exchangeOrderId\":\"12345\",\"validity\":\"DAY\"}"))
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
			ModifyAMOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("modify amo order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyAMOOrder", strings.NewReader("{\"clientId\":\"Client1\",\"disclosedQuantity\":0,\"exchange\":\"NSE\",\"executionType\":\"AMO\",\"filledQuantity\":0,\"instrumentToken\":\"1330\",\"lastActivityReference\":0,\"omsOrderID\":\"20220913-31\",\"orderType\":\"LIMIT\",\"price\":2000,\"product\":\"CNC\",\"quantity\":1,\"triggerPrice\":0,\"remainingQuantity\":0,\"exchangeOrderId\":\"12345\",\"validity\":\"DAY\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitOrderProvider(orderMock{})

	//mock business layer response
	modifyAMOOrderMock = func(req models.ModifyAMORequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			ModifyAMOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("modify amo order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestCancelAMOOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/cancelAMOOrder", nil)
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
			CancelAMOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Cancel AMO order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/cancelAMOOrder", strings.NewReader("{\"clientId\":\"Client1\",\"executionType\":\"AMO\",\"omsOrderID\":\"20220920-39\"}"))
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
			CancelAMOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Cancel AMO order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 bad request
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqH := models.ReqHeader{}
	reqH.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqH)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/cancelAMOOrder", strings.NewReader("{\"clientId\":\"Client1\",\"executionType\":\"NO\",\"omsOrderID\":\"20220920-39\"}"))
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
			CancelAMOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Cancel AMO order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/cancelAMOOrder", strings.NewReader("{\"clientId\":\"Client1\",\"executionType\":\"AMO\",\"omsOrderID\":\"20220920-39\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitOrderProvider(orderMock{})

	//mock business layer response
	cancelAMOOrderMock = func(req models.CancelOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			CancelAMOOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Cancel AMO order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestPendingOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/pendingOrder", nil)
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
			PendingOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Pending order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/pendingOrder", strings.NewReader("{\"clientId\":\"Client1\",\"type\":\"comp\"}"))
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
			PendingOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Pending order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 bad request
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqH := models.ReqHeader{}
	reqH.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqH)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/pendingOrder", strings.NewReader("{\"clientId\":\"Client1\",\"type\":\"comp\"}"))
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
			PendingOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Pending order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/pendingOrder", strings.NewReader("{\"clientId\":\"Client1\",\"type\":\"pending\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitOrderProvider(orderMock{})

	//mock business layer response
	pendingOrderMock = func(req models.PendingOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			PendingOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Pending order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestCompletedOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/completedOrder", nil)
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
			CompletedOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Completed order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/completedOrder", strings.NewReader("{\"clientId\":\"Client1\",\"type\":\"completed\"}"))
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
			CompletedOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Completed order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 bad request
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqH := models.ReqHeader{}
	reqH.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqH)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/completedOrder", strings.NewReader("{\"clientId\":\"Client1\",\"type\":\"comp\"}"))
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
			CompletedOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Completed order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/completedOrder", strings.NewReader("{\"clientId\":\"Client1\",\"type\":\"completed\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitOrderProvider(orderMock{})

	//mock business layer response
	completedOrderMock = func(req models.CompletedOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			CompletedOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Completed order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestTradeBook(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/tradeBook", nil)
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
			TradeBook(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("TradeBook() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/tradeBook", strings.NewReader("{\"clientId\":\"Client1\"}"))
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
			TradeBook(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("TradeBook() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/tradeBook", strings.NewReader("{\"clientId\":\"Client1\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitOrderProvider(orderMock{})

	//mock business layer response
	tradeBookMock = func(req models.TradeBookRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			TradeBook(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("TradeBook() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestOrderHistory(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/orderHistory", nil)
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
			OrderHistory(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("OrderHistory() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/orderHistory", strings.NewReader("{\"clientId\":\"Client1\",\"omsOrderId\":\"completed\"}"))
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
			OrderHistory(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("OrderHistory() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/orderHistory", strings.NewReader("{\"clientId\":\"Client1\",\"omsOrderId\":\"20220920-4\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitOrderProvider(orderMock{})

	//mock business layer response
	orderHistoryMock = func(req models.OrderHistoryRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			OrderHistory(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("OrderHistory() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestPlaceGTTOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/createGTTOrder", nil)
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
			CreateGTTOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Create GTT order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/createGTTOrder", strings.NewReader("{\"actionType\":\"single_order\",\"expiryTime\":\"2022-10-17\",\"order\":{\"clientId\":\"Client1\",\"disclosedQuantity\":0,\"exchange\":\"NSE\",\"instrumentToken\":\"22\",\"marketProtectionPercentage\":0,\"orderSide\":\"BUY\",\"orderType\":\"LIMIT\",\"price\":0,\"product\":\"CNC\",\"quantity\":0,\"slOrderPrice\":0,\"slOrderQuantity\":0,\"slTriggerPrice\":0,\"triggerPrice\":0,\"userOrderId\":0}}"))
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
			CreateGTTOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Create GTT order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 bad request
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqH := models.ReqHeader{}
	reqH.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqH)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/createGTTOrder", strings.NewReader("{\"actionType\":\"single\",\"expiryTime\":\"2022-10-17\",\"order\":{\"clientId\":\"Client1\",\"disclosedQuantity\":0,\"exchange\":\"NSE\",\"instrumentToken\":\"22\",\"marketProtectionPercentage\":0,\"orderSide\":\"BUY\",\"orderType\":\"LIMIT\",\"price\":0,\"product\":\"CNC\",\"quantity\":0,\"slOrderPrice\":0,\"slOrderQuantity\":0,\"slTriggerPrice\":0,\"triggerPrice\":0,\"userOrderId\":0}}"))
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
			CreateGTTOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Create GTT order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/createGTTOrder", strings.NewReader("{\"actionType\":\"single_order\",\"expiryTime\":\"2022-10-17\",\"clientId\":\"Client1\",\"disclosedQuantity\":0,\"exchange\":\"NSE\",\"instrumentToken\":\"22\",\"marketProtectionPercentage\":0,\"orderSide\":\"BUY\",\"orderType\":\"LIMIT\",\"price\":0,\"product\":\"CNC\",\"quantity\":0,\"slOrderPrice\":0,\"slOrderQuantity\":0,\"slTriggerPrice\":0,\"triggerPrice\":0,\"userOrderId\":0}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitOrderProvider(orderMock{})

	//mock business layer response
	placeGTTOrderMock = func(req models.CreateGTTOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			CreateGTTOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Create GTT order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestModifyGTTOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyGTTOrder", nil)
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
			ModifyGTTOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Modify GTT order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyGTTOrder", strings.NewReader("{\"expiryTime\":\"2022-10-17\",\"order\":{\"clientId\":\"Client1\",\"device\":\"WEB\",\"disclosedQuantity\":0,\"exchange\":\"NSE\",\"instrumentToken\":\"22\",\"marketProtectionPercentage\":0,\"orderType\":\"LIMIT\",\"price\":0,\"product\":\"CNC\",\"quantity\":0,\"slOrderPrice\":0,\"slOrderQuantity\":0,\"slTriggerPrice\":0,\"triggerPrice\":0,\"userOrderId\":0}}"))
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
			ModifyGTTOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Modify GTT order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 bad request
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqH := models.ReqHeader{}
	reqH.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqH)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyGTTOrder", strings.NewReader("{\"expiryTime\":\"2022-10-17\",\"order\":{\"clientId\":\"Client1\",\"device\":\"WEB\",\"disclosedQuantity\":0,\"exchange\":\"NO\",\"instrumentToken\":\"22\",\"marketProtectionPercentage\":0,\"orderType\":\"LIMIT\",\"price\":0,\"product\":\"CNC\",\"quantity\":0,\"slOrderPrice\":0,\"slOrderQuantity\":0,\"slTriggerPrice\":0,\"triggerPrice\":0,\"userOrderId\":0}}"))
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
			ModifyGTTOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Modify GTT order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyGTTOrder", strings.NewReader("{\"expiryTime\":\"2022-10-17\",\"clientId\":\"Client1\",\"device\":\"WEB\",\"disclosedQuantity\":0,\"exchange\":\"NSE\",\"instrumentToken\":\"22\",\"marketProtectionPercentage\":0,\"orderType\":\"LIMIT\",\"price\":0,\"product\":\"CNC\",\"quantity\":0,\"slOrderPrice\":0,\"slOrderQuantity\":0,\"slTriggerPrice\":0,\"triggerPrice\":0,\"userOrderId\":0}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitOrderProvider(orderMock{})

	//mock business layer response
	modifyGTTOrderMock = func(req models.ModifyGTTOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			ModifyGTTOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Modify GTT order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestCancelGTTOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/cancelGTTOrder", nil)
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
			CancelGTTOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Cancel GTT Order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/cancelGTTOrder", strings.NewReader("{\"clientId\":\"Client1\",\"id\":\"ID\"}"))
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
			CancelGTTOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Cancel GTT Order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/cancelGTTOrder", strings.NewReader("{\"clientId\":\"Client1\",\"id\":\"ID\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitOrderProvider(orderMock{})

	//mock business layer response
	cancelGTTOrderMock = func(req models.CancelGTTOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			CancelGTTOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Cancel GTT Order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}

func TestFetchGTTOrder(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/fetchGTTOrder", nil)
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
			FetchGTTOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Fetch GTT Order() = %v, want %v", string(b), expected)
			}
		})
	}

	//2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/fetchGTTOrder", strings.NewReader("{\"clientId\":\"Client1\"}"))
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
			FetchGTTOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Fetch GTT Order() = %v, want %v", string(b), expected)
			}
		})
	}

	//3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/fetchGTTOrder", strings.NewReader("{\"clientId\":\"Client1\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitOrderProvider(orderMock{})

	//mock business layer response
	fetchGTTOrderMock = func(req models.FetchGTTOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			FetchGTTOrder(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("Fetch GTT Order() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 400)
			}
		})
	}

}
