package v2

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
	placeBOOrderMock       func(req models.PlaceBOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	modifyBOOrderMock      func(req models.ModifyBOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	cancelBOOrderMock      func(req models.ExitBOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	placeCOOrderMock       func(req models.PlaceCOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	modifyCOOrderMock      func(req models.ModifyCOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	cancelCOOrderMock      func(req models.ExitCOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	placeSpreadOrderMock   func(req models.PlaceSpreadOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	modifySpreadOrderMock  func(req models.ModifySpreadOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
	cancelSpreadOrderMock  func(req models.ExitSpreadOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes)
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

func TestPlaceAnOrder(t *testing.T) {
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

	//3 success

	orderTypes := []string{"REGULAR", "AMO", "BO", "CO", "Spread", "GTT", "X"}
	// orderTypes := []string{"REGULAR"}

	for _, orderType := range orderTypes {
		switch orderType {
		case "REGULAR":
			w = httptest.NewRecorder()
			ctx, _ = gin.CreateTestContext(w)
			reqHsuccess := models.ReqHeader{}
			reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
			ctx.Set("reqH", reqHsuccess)
			requestBody := "{\"clientId\":\"9CD12\",\"orderType\":\"REGULAR\",\"requestPacket\":{\"clientId\":\"CLient1\",\"disclosedQuantity\":0,\"exchange\":\"BSE\",\"executionType\":\"REGULAR\",\"instrumentToken\":\"10666\",\"orderSide\":\"BUY\",\"orderType\":\"LIMIT\",\"price\":1000,\"product\":\"CNC\",\"quantity\":10,\"triggerPrice\":0,\"validity\":\"DAY\"}}"
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeOrder", strings.NewReader(requestBody))
			expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

			InitOrderProviderV2(orderMock{})

			//mock business layer response
			placeOrderMock = func(req models.PlaceOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
				var res apihelpers.APIRes
				res.Status = true
				res.Message = "SUCCESS"
				return http.StatusOK, res
			}

			t.Run("Success", func(t *testing.T) {
				PlaceOrder(ctx)
				b, _ := ioutil.ReadAll(w.Body)
				if strings.TrimSuffix(string(b), "\n") != expected {
					t.Errorf("place order() = %v, want %v", string(b), expected)
				}
				if w.Code != 200 {
					t.Errorf("logoutController() = %v, want %v", w.Code, 400)
				}
			})
			continue

		case "AMO":
			w = httptest.NewRecorder()
			ctx, _ = gin.CreateTestContext(w)
			reqHsuccess := models.ReqHeader{}
			reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
			ctx.Set("reqH", reqHsuccess)
			requestBody := "{\"clientId\":\"9CD12\",\"orderType\":\"AMO\",\"requestPacket\":{\"exchange\":\"NSE\",\"instrumentToken\":\"10666\",\"clientId\":\"Client1\",\"disclosedQuantity\":0,\"executionType\":\"AMO\",\"orderSide\":\"BUY\",\"orderType\":\"LIMIT\",\"price\":1000,\"product\":\"CNC\",\"quantity\":10,\"triggerPrice\":0,\"validity\":\"DAY\"}}"
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeOrder", strings.NewReader(requestBody))
			expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

			InitOrderProviderV2(orderMock{})

			//mock business layer response
			placeAMOOrderMock = func(req models.PlaceOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
				var res apihelpers.APIRes
				res.Status = true
				res.Message = "SUCCESS"
				return http.StatusOK, res
			}

			t.Run("Success", func(t *testing.T) {
				PlaceOrder(ctx)
				b, _ := ioutil.ReadAll(w.Body)
				if strings.TrimSuffix(string(b), "\n") != expected {
					t.Errorf("place order() = %v, want %v", string(b), expected)
				}
				if w.Code != 200 {
					t.Errorf("logoutController() = %v, want %v", w.Code, 400)
				}
			})
			continue

		case "BO":
			w = httptest.NewRecorder()
			ctx, _ = gin.CreateTestContext(w)
			reqHsuccess := models.ReqHeader{}
			reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
			ctx.Set("reqH", reqHsuccess)
			requestBody := "{\"clientId\":\"9CD12\",\"orderType\":\"BO\",\"requestPacket\":{\"clientId\":\"CLient1\",\"disclosedQuantity\":0,\"exchange\":\"NSE\",\"instrumentToken\":\"10666\",\"isTrailing\": true,\"orderSide\":\"BUY\",\"orderType\":\"LIMIT\",\"price\":1000,\"product\":\"CNC\",\"quantity\":10,\"squareOffValue\":1,\"stopLossValue\":1,\"trailingStopLoss\":\"0.05\",\"triggerPrice\":0,\"userOrderId\":1002,\"validity\":\"DAY\"}}"
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeOrder", strings.NewReader(requestBody))
			expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

			InitConditionalOrderProviderV2(conditionalOrderMock{})

			//mock business layer response
			placeBOOrderMock = func(req models.PlaceBOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
				var res apihelpers.APIRes
				res.Status = true
				res.Message = "SUCCESS"
				return http.StatusOK, res
			}

			t.Run("Success", func(t *testing.T) {
				PlaceOrder(ctx)
				b, _ := ioutil.ReadAll(w.Body)
				if strings.TrimSuffix(string(b), "\n") != expected {
					t.Errorf("place order() = %v, want %v", string(b), expected)
				}
				if w.Code != 200 {
					t.Errorf("logoutController() = %v, want %v", w.Code, 400)
				}
			})
			continue

		case "CO":
			w = httptest.NewRecorder()
			ctx, _ = gin.CreateTestContext(w)
			reqHsuccess := models.ReqHeader{}
			reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
			ctx.Set("reqH", reqHsuccess)
			requestBody := "{\"clientId\":\"9CD12\",\"orderType\":\"CO\",\"requestPacket\":{\"clientId\": \"CLient1\",\"disclosedQuantity\": 2,\"exchange\": \"NSE\",\"instrumentToken\": \"22\",\"orderSide\": \"BUY\",\"orderType\": \"LIMIT\",\"price\": 34.2,\"product\": \"CNC\",\"quantity\": 2,\"stopLossValue\": 33,\"trailingStopLoss\": 33,\"userOrderId\": 91261928,\"validity\": \"DAY\"}}"
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeOrder", strings.NewReader(requestBody))
			expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

			InitConditionalOrderProviderV2(conditionalOrderMock{})

			//mock business layer response
			placeCOOrderMock = func(req models.PlaceCOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
				var res apihelpers.APIRes
				res.Status = true
				res.Message = "SUCCESS"
				return http.StatusOK, res
			}

			t.Run("Success", func(t *testing.T) {
				PlaceOrder(ctx)
				b, _ := ioutil.ReadAll(w.Body)
				if strings.TrimSuffix(string(b), "\n") != expected {
					t.Errorf("place order() = %v, want %v", string(b), expected)
				}
				if w.Code != 200 {
					t.Errorf("logoutController() = %v, want %v", w.Code, 400)
				}
			})
			continue

		case "Spread":
			w = httptest.NewRecorder()
			ctx, _ = gin.CreateTestContext(w)
			reqHsuccess := models.ReqHeader{}
			reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
			ctx.Set("reqH", reqHsuccess)
			requestBody := "{\"clientId\":\"9CD12\",\"orderType\":\"Spread\",\"requestPacket\":{\"clientId\": \"CLIENT1\",\"disclosedQuantity\": 2,\"exchange\": \"NSE\",\"instrumentToken\": \"22\",\"orderSide\": \"BUY\",\"orderType\": \"LIMIT\",\"price\": 34.2,\"product\": \"CNC\",\"quantity\": 2,\"userOrderId\": 91261928,\"validity\": \"DAY\"}}"
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeOrder", strings.NewReader(requestBody))
			expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

			InitConditionalOrderProviderV2(conditionalOrderMock{})

			//mock business layer response
			placeSpreadOrderMock = func(req models.PlaceSpreadOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
				var res apihelpers.APIRes
				res.Status = true
				res.Message = "SUCCESS"
				return http.StatusOK, res
			}

			t.Run("Success", func(t *testing.T) {
				PlaceOrder(ctx)
				b, _ := ioutil.ReadAll(w.Body)
				if strings.TrimSuffix(string(b), "\n") != expected {
					t.Errorf("place order() = %v, want %v", string(b), expected)
				}
				if w.Code != 200 {
					t.Errorf("logoutController() = %v, want %v", w.Code, 400)
				}
			})
			continue

		default:
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			reqHsuccess := models.ReqHeader{}
			reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
			ctx.Set("reqH", reqHsuccess)
			requestBody := "{\"clientId\":\"9CD12\",\"orderType\":\"X\",\"requestPacket\":{\"clientId\": \"CLIENT1\",\"disclosedQuantity\": 2,\"exchange\": \"NSE\",\"instrumentToken\": \"22\",\"orderSide\": \"BUY\",\"orderType\": \"LIMIT\",\"price\": 34.2,\"product\": \"CNC\",\"quantity\": 2,\"userOrderId\": 91261928,\"validity\": \"DAY\"}}"
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/placeOrder", strings.NewReader(requestBody))
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
		}
	}

}

func TestModifyAnOrder(t *testing.T) {
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

	//3 success

	orderTypes := []string{"REGULAR", "AMO", "BO", "CO", "Spread", "GTT"}
	// orderTypes := []string{"REGULAR"}

	for _, orderType := range orderTypes {
		switch orderType {
		case "REGULAR":
			w = httptest.NewRecorder()
			ctx, _ = gin.CreateTestContext(w)
			reqHsuccess := models.ReqHeader{}
			reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
			ctx.Set("reqH", reqHsuccess)
			requestBody := "{\"clientId\":\"9CD12\",\"orderType\":\"REGULAR\",\"requestPacket\":{\"exchange\":\"NSE\",\"instrumentToken\":\"14366\",\"clientId\":\"9CD12\",\"orderType\":\"LIMIT\",\"price\":14.5,\"quantity\":1,\"disclosedQuantity\":0,\"validity\":\"DAY\",\"product\":\"MIS\",\"omsOrderId\":\"20231127-71\",\"exchangeOrderId\":\"\",\"filledQuantity\":0,\"remainingQuantity\":0,\"lastActivityReference\":0,\"triggerPrice\":0,\"executionType\":\"REGULAR\"}}"
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyOrder", strings.NewReader(requestBody))
			expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

			InitOrderProviderV2(orderMock{})

			//mock business layer response
			modifyOrderMock = func(req models.ModifyOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
				var res apihelpers.APIRes
				res.Status = true
				res.Message = "SUCCESS"
				return http.StatusOK, res
			}
			t.Run("Success", func(t *testing.T) {
				ModifyOrder(ctx)
				b, _ := ioutil.ReadAll(w.Body)
				if strings.TrimSuffix(string(b), "\n") != expected {
					t.Errorf("modify order() = %v, want %v", string(b), expected)
				}
				if w.Code != 200 {
					t.Errorf("logoutController() = %v, want %v", w.Code, 400)
				}
			})
			continue

		case "AMO":
			w = httptest.NewRecorder()
			ctx, _ = gin.CreateTestContext(w)
			reqHsuccess := models.ReqHeader{}
			reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
			ctx.Set("reqH", reqHsuccess)
			requestBody := "{\"clientId\":\"9CD12\",\"orderType\":\"AMO\",\"requestPacket\":{\"clientId\":\"Client1\",\"disclosedQuantity\":0,\"exchange\":\"NSE\",\"executionType\":\"AMO\",\"filledQuantity\":0,\"instrumentToken\":\"1330\",\"lastActivityReference\":0,\"omsOrderID\":\"20220913-31\",\"orderType\":\"LIMIT\",\"price\":2000,\"product\":\"CNC\",\"quantity\":1,\"triggerPrice\":0,\"remainingQuantity\":0,\"exchangeOrderId\":\"12345\",\"validity\":\"DAY\"}}"
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyOrder", strings.NewReader(requestBody))
			expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

			InitOrderProviderV2(orderMock{})

			//mock business layer response
			modifyAMOOrderMock = func(req models.ModifyAMORequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
				var res apihelpers.APIRes
				res.Status = true
				res.Message = "SUCCESS"
				return http.StatusOK, res
			}
			t.Run("Success", func(t *testing.T) {
				ModifyOrder(ctx)
				b, _ := ioutil.ReadAll(w.Body)
				if strings.TrimSuffix(string(b), "\n") != expected {
					t.Errorf("modify order() = %v, want %v", string(b), expected)
				}
				if w.Code != 200 {
					t.Errorf("logoutController() = %v, want %v", w.Code, 400)
				}
			})
			continue

		case "BO":
			w = httptest.NewRecorder()
			ctx, _ = gin.CreateTestContext(w)
			reqHsuccess := models.ReqHeader{}
			reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
			ctx.Set("reqH", reqHsuccess)
			requestBody := "{\"clientId\":\"9CD12\",\"orderType\":\"BO\",\"requestPacket\":{\"clientId\":\"CLient1\",\"disclosedQuantity\":2,\"exchange\":\"NSE\",\"exchangeOrderId\": \"ionwdg123\",\"filledQuantity\": 1,\"instrumentToken\":\"10666\",\"isTrailing\": true,\"lastActivityReference\": 1325938440097498600,\"omsOrderId\": \"123445\",\"orderType\": \"LIMIT\",\"price\":34.2,\"product\":\"CNC\",\"quantity\":10,\"remainingQuantity\": 0,\"squareOffValue\":1,\"stopLossValue\":1,\"trailingStopLoss\":41,\"triggerPrice\":0,\"validity\":\"DAY\"}}"

			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyOrder", strings.NewReader(requestBody))
			expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

			InitConditionalOrderProviderV2(conditionalOrderMock{})

			//mock business layer response
			modifyBOOrderMock = func(req models.ModifyBOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
				var res apihelpers.APIRes
				res.Status = true
				res.Message = "SUCCESS"
				return http.StatusOK, res
			}
			t.Run("Success", func(t *testing.T) {
				ModifyOrder(ctx)
				b, _ := ioutil.ReadAll(w.Body)
				if strings.TrimSuffix(string(b), "\n") != expected {
					t.Errorf("modify order() = %v, want %v", string(b), expected)
				}
				if w.Code != 200 {
					t.Errorf("logoutController() = %v, want %v", w.Code, 400)
				}
			})
			continue

		case "CO":
			w = httptest.NewRecorder()
			ctx, _ = gin.CreateTestContext(w)
			reqHsuccess := models.ReqHeader{}
			reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
			ctx.Set("reqH", reqHsuccess)
			requestBody := "{\"clientId\":\"9CD12\",\"orderType\":\"CO\",\"requestPacket\":{\"clientId\": \"Client1\",\"disclosedQuantity\": 2,\"exchange\": \"NSE\",\"exchangeOrderId\": \"ionwdg123\",\"filledQuantity\": 1,\"instrumentToken\": \"22\",\"lastActivityReference\": 1325938440097498600,\"omsOrderId\": \"123445\",\"orderType\": \"LIMIT\",\"price\": 34.2,\"product\": \"CNC\",\"quantity\": 2,\"remainingQuantity\": 0,\"stopLossValue\": 0,\"trailingStopLoss\": 33,\"validity\": \"DAY\"}}"
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyOrder", strings.NewReader(requestBody))
			expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

			InitConditionalOrderProviderV2(conditionalOrderMock{})

			//mock business layer response
			modifyCOOrderMock = func(req models.ModifyCOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
				var res apihelpers.APIRes
				res.Status = true
				res.Message = "SUCCESS"
				return http.StatusOK, res
			}
			t.Run("Success", func(t *testing.T) {
				ModifyOrder(ctx)
				b, _ := ioutil.ReadAll(w.Body)
				if strings.TrimSuffix(string(b), "\n") != expected {
					t.Errorf("modify order() = %v, want %v", string(b), expected)
				}
				if w.Code != 200 {
					t.Errorf("logoutController() = %v, want %v", w.Code, 400)
				}
			})
			continue

		case "Spread":
			w = httptest.NewRecorder()
			ctx, _ = gin.CreateTestContext(w)
			reqHsuccess := models.ReqHeader{}
			reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
			ctx.Set("reqH", reqHsuccess)
			requestBody := "{\"clientId\":\"9CD12\",\"orderType\":\"Spread\",\"requestPacket\":{\"clientId\": \"CLIENT1\",\"disclosedQuantity\": 0,\"exchange\": \"NSE\",\"exchangeOrderId\": \"kjh\",\"instrumentToken\": \"22\",\"isTrailing\": true,\"omsOrderId\": \"123445\",\"orderType\": \"LIMIT\",\"price\": 34.2,\"prodType\": \"fdsf\",\"product\": \"CNC\",\"quantity\": 2,\"squareOffValue\": 40,\"stopLossValue\": 39,\"trailingStopLoss\": 33,\"triggerPrice\": 38,\"validity\": \"DAY\"}}"
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/modifyOrder", strings.NewReader(requestBody))
			expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

			InitConditionalOrderProviderV2(conditionalOrderMock{})

			//mock business layer response
			modifySpreadOrderMock = func(req models.ModifySpreadOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
				var res apihelpers.APIRes
				res.Status = true
				res.Message = "SUCCESS"
				return http.StatusOK, res
			}
			t.Run("Success", func(t *testing.T) {
				ModifyOrder(ctx)
				b, _ := ioutil.ReadAll(w.Body)
				if strings.TrimSuffix(string(b), "\n") != expected {
					t.Errorf("modify order() = %v, want %v", string(b), expected)
				}
				if w.Code != 200 {
					t.Errorf("logoutController() = %v, want %v", w.Code, 400)
				}
			})
			continue

		default:
			w = httptest.NewRecorder()
			ctx, _ = gin.CreateTestContext(w)
			reqHsuccess := models.ReqHeader{}
			reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
			ctx.Set("reqH", reqHsuccess)
			requestBody := "{\"clientId\":\"9CD12\",\"orderType\":\"X\",\"requestPacket\":{\"clientId\":\"\",\"id\":\"\"}}"
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/ModifyOrder", strings.NewReader(requestBody))
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
						t.Errorf("Cancel order() = %v, want %v", string(b), expected)
					}
				})
			}
		}

	}

}

func TestCancelAnOrder(t *testing.T) {
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

	//3 Success
	orderTypes := []string{"REGULAR", "AMO", "BO", "CO", "Spread", "GTT", "X"}

	for _, orderType := range orderTypes {
		switch orderType {
		case "REGULAR":
			w = httptest.NewRecorder()
			ctx, _ = gin.CreateTestContext(w)
			reqHsuccess := models.ReqHeader{}
			reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
			ctx.Set("reqH", reqHsuccess)
			requestBody := "{\"clientId\":\"9CD12\",\"orderType\":\"REGULAR\",\"requestPacket\":{\"clientId\":\"Client1\",\"executionType\":\"REGULAR\",\"omsOrderID\":\"20220920-39\"}}"
			// requestBody := generateOrderRequestJSON(orderType)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/cancelOrder", strings.NewReader(requestBody))
			expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

			InitOrderProviderV2(orderMock{})

			cancelOrderMock = func(req models.CancelOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
				var res apihelpers.APIRes
				res.Status = true
				res.Message = "SUCCESS"
				return http.StatusOK, res
			}

			t.Run("Success", func(t *testing.T) {
				CancelOrder(ctx)
				b, _ := ioutil.ReadAll(w.Body)
				if strings.TrimSuffix(string(b), "\n") != expected {
					t.Errorf("Cancel order() = %v, want %v", string(b), expected)
				}
				if w.Code != http.StatusOK {
					t.Errorf("Cancel order() = %v, want %v", w.Code, http.StatusOK)
				}
			})
			continue

		case "AMO":
			w = httptest.NewRecorder()
			ctx, _ = gin.CreateTestContext(w)
			reqHsuccess := models.ReqHeader{}
			reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
			ctx.Set("reqH", reqHsuccess)
			requestBody := "{\"clientId\":\"9CD12\",\"orderType\":\"AMO\",\"requestPacket\":{\"clientId\":\"9CD12\",\"omsOrderId\":\"\",\"executionType\":\"REGULAR\"}}"
			// requestBody := generateOrderRequestJSON(orderType)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/cancelOrder", strings.NewReader(requestBody))
			expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

			InitOrderProviderV2(orderMock{})

			cancelAMOOrderMock = func(req models.CancelOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
				var res apihelpers.APIRes
				res.Status = true
				res.Message = "SUCCESS"
				return http.StatusOK, res
			}

			t.Run("Success", func(t *testing.T) {
				CancelOrder(ctx)
				b, _ := ioutil.ReadAll(w.Body)
				if strings.TrimSuffix(string(b), "\n") != expected {
					t.Errorf("Cancel order() = %v, want %v", string(b), expected)
				}
				if w.Code != http.StatusOK {
					t.Errorf("Cancel order() = %v, want %v", w.Code, http.StatusOK)
				}
			})
			continue

		case "BO":
			w = httptest.NewRecorder()
			ctx, _ = gin.CreateTestContext(w)
			reqHsuccess := models.ReqHeader{}
			reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
			ctx.Set("reqH", reqHsuccess)
			requestBody := "{\"clientId\":\"9CD12\",\"orderType\":\"BO\",\"requestPacket\":{\"clientId\":\"\",\"exchangeOrderId\":\"\",\"legOrderIndicator\":\"\",\"omsOrderId\":\"\",\"status\":\"\"}}"
			// requestBody := generateOrderRequestJSON(orderType)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/cancelOrder", strings.NewReader(requestBody))
			expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

			InitConditionalOrderProviderV2(conditionalOrderMock{})

			cancelBOOrderMock = func(req models.ExitBOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
				var res apihelpers.APIRes
				res.Status = true
				res.Message = "SUCCESS"
				return http.StatusOK, res
			}

			t.Run("Success", func(t *testing.T) {
				CancelOrder(ctx)
				b, _ := ioutil.ReadAll(w.Body)
				if strings.TrimSuffix(string(b), "\n") != expected {
					t.Errorf("Cancel order() = %v, want %v", string(b), expected)
				}
				if w.Code != http.StatusOK {
					t.Errorf("Cancel order() = %v, want %v", w.Code, http.StatusOK)
				}
			})
			continue

		case "CO":
			w = httptest.NewRecorder()
			ctx, _ = gin.CreateTestContext(w)
			reqHsuccess := models.ReqHeader{}
			reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
			ctx.Set("reqH", reqHsuccess)
			requestBody := "{\"clientId\":\"9CD12\",\"orderType\":\"CO\",\"requestPacket\":{\"clientId\":\"\",\"exchangeOrderId\":\"\",\"legOrderIndicator\":\"AMO\",\"omsOrderId\":\"\"}}"
			// requestBody := generateOrderRequestJSON(orderType)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/cancelOrder", strings.NewReader(requestBody))
			expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

			InitConditionalOrderProviderV2(conditionalOrderMock{})

			//mock business layer response
			cancelCOOrderMock = func(req models.ExitCOOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
				var res apihelpers.APIRes
				res.Status = true
				res.Message = "SUCCESS"
				return http.StatusOK, res
			}

			t.Run("Success", func(t *testing.T) {
				CancelOrder(ctx)
				b, _ := ioutil.ReadAll(w.Body)
				if strings.TrimSuffix(string(b), "\n") != expected {
					t.Errorf("Cancel order() = %v, want %v", string(b), expected)
				}
				if w.Code != http.StatusOK {
					t.Errorf("Cancel order() = %v, want %v", w.Code, http.StatusOK)
				}
			})
			continue

		case "Spread":
			w = httptest.NewRecorder()
			ctx, _ = gin.CreateTestContext(w)
			reqHsuccess := models.ReqHeader{}
			reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
			ctx.Set("reqH", reqHsuccess)
			requestBody := "{\"clientId\":\"9CD12\",\"orderType\":\"Spread\",\"requestPacket\":{\"clientId\":\"\",\"exchangeOrderId\":\"\",\"legOrderIndicator\":\"AMO\",\"omsOrderId\":\"\",\"status\":\"\"}}"
			// requestBody := generateOrderRequestJSON(orderType)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/cancelOrder", strings.NewReader(requestBody))
			expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

			InitConditionalOrderProviderV2(conditionalOrderMock{})

			//mock business layer response
			cancelSpreadOrderMock = func(req models.ExitSpreadOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
				var res apihelpers.APIRes
				res.Status = true
				res.Message = "SUCCESS"
				return http.StatusOK, res
			}

			t.Run("Success", func(t *testing.T) {
				CancelOrder(ctx)
				b, _ := ioutil.ReadAll(w.Body)
				if strings.TrimSuffix(string(b), "\n") != expected {
					t.Errorf("Cancel order() = %v, want %v", string(b), expected)
				}
				if w.Code != http.StatusOK {
					t.Errorf("Cancel order() = %v, want %v", w.Code, http.StatusOK)
				}
			})
			continue

		case "GTT":
			w = httptest.NewRecorder()
			ctx, _ = gin.CreateTestContext(w)
			reqHsuccess := models.ReqHeader{}
			reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
			ctx.Set("reqH", reqHsuccess)
			requestBody := "{\"clientId\":\"9CD12\",\"orderType\":\"GTT\",\"requestPacket\":{\"clientId\":\"\",\"id\":\"\"}}"
			// requestBody := generateOrderRequestJSON(orderType)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/cancelOrder", strings.NewReader(requestBody))
			expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

			InitOrderProviderV2(orderMock{})

			//mock business layer response
			cancelGTTOrderMock = func(req models.CancelGTTOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
				var res apihelpers.APIRes
				res.Status = true
				res.Message = "SUCCESS"
				return http.StatusOK, res
			}

			t.Run("Success", func(t *testing.T) {
				CancelOrder(ctx)
				b, _ := ioutil.ReadAll(w.Body)
				if strings.TrimSuffix(string(b), "\n") != expected {
					t.Errorf("Cancel order() = %v, want %v", string(b), expected)
				}
				if w.Code != http.StatusOK {
					t.Errorf("Cancel order() = %v, want %v", w.Code, http.StatusOK)
				}
			})
			continue

		default:
			w = httptest.NewRecorder()
			ctx, _ = gin.CreateTestContext(w)
			reqHsuccess := models.ReqHeader{}
			reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
			ctx.Set("reqH", reqHsuccess)
			requestBody := "{\"clientId\":\"9CD12\",\"orderType\":\"X\",\"requestPacket\":{\"clientId\":\"\",\"id\":\"\"}}"
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orderapis/cancelOrder", strings.NewReader(requestBody))
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
		}

	}

}
