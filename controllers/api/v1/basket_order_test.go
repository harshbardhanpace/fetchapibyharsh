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
	createBasketMock               func(req models.CreateBasketReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	fetchBasketMock                func(req models.FetchBasketReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	deleteBasketMock               func(req models.DeleteBasketReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	addBasketInstrumentMock        func(req models.AddBasketInstrumentReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	editBasketInstrumentMock       func(req models.EditBasketInstrumentReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	deleteBasketInstrumentMock     func(req models.DeleteBasketInstrumentReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	renameBasketMock               func(req models.RenameBasketReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	executeBasketMock              func(req models.ExecuteBasketReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
	updateBasketExecutionStateMock func(req models.UpdateBasketExecutionStateReq, reqH models.ReqHeader) (int, apihelpers.APIRes)
)

type basketOrderMock struct{}

func (m basketOrderMock) CreateBasket(req models.CreateBasketReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return createBasketMock(req, reqH)
}

func (m basketOrderMock) FetchBasket(req models.FetchBasketReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return fetchBasketMock(req, reqH)
}

func (m basketOrderMock) DeleteBasket(req models.DeleteBasketReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return deleteBasketMock(req, reqH)
}

func (m basketOrderMock) AddBasketInstrument(req models.AddBasketInstrumentReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return addBasketInstrumentMock(req, reqH)
}

func (m basketOrderMock) EditBasketInstrument(req models.EditBasketInstrumentReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return editBasketInstrumentMock(req, reqH)
}

func (m basketOrderMock) DeleteBasketInstrument(req models.DeleteBasketInstrumentReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return deleteBasketInstrumentMock(req, reqH)
}
func (m basketOrderMock) RenameBasket(req models.RenameBasketReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return renameBasketMock(req, reqH)
}
func (m basketOrderMock) ExecuteBasket(req models.ExecuteBasketReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return executeBasketMock(req, reqH)
}
func (m basketOrderMock) UpdateBasketExecutionState(req models.UpdateBasketExecutionStateReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return updateBasketExecutionStateMock(req, reqH)
}

func TestCreateBasket(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/createBasket", nil)
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
			CreateBasket(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("create basket() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/createBasket", strings.NewReader("{\"LoginID\":\"Login1\",\"Name\":\"XYZ\",\"Type\":\"NORMAL\",\"ProductType\":\"ALL\",\"OrderType\":\"ALL\"}"))
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
			CreateBasket(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("create basket() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/createBasket", strings.NewReader("{\"LoginID\":\"Login1\",\"Name\":\"XYZ\",\"Type\":\"NORMAL\",\"ProductType\":\"ALL\",\"OrderType\":\"ALL\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitBasketOrderProvider(basketOrderMock{})

	//mock business layer response
	createBasketMock = func(req models.CreateBasketReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			CreateBasket(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("create basket() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}

}

func TestFetchBasket(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/fetchBasket", nil)
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
			FetchBasket(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("fetch basket() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/fetchBasket", strings.NewReader("{\"LoginID\":\"Login1\"}"))
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
			FetchBasket(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("fetch basket() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/fetchBasket", strings.NewReader("{\"LoginID\":\"Login1\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitBasketOrderProvider(basketOrderMock{})

	//mock business layer response
	fetchBasketMock = func(req models.FetchBasketReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			FetchBasket(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("fetch basket() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}

}

func TestDeleteBasket(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/deleteBasket", nil)
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
			DeleteBasket(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("delete basket() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/deleteBasket", strings.NewReader("{\"BasketID\":\"xyzgsfajd\",\"Name\":\"xyz\",\"SipCount\":0}"))
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
			DeleteBasket(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("delete basket() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/deleteBasket", strings.NewReader("{\"BasketID\":\"xyzgsfajd\",\"Name\":\"xyz\",\"SipCount\":0}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitBasketOrderProvider(basketOrderMock{})

	//mock business layer response
	deleteBasketMock = func(req models.DeleteBasketReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			DeleteBasket(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("delete basket() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}

}

func TestAddBasketInstrument(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/addBasketInstrument", nil)
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
			ExecuteBasket(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("add basket instrument() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/addBasketInstrument", strings.NewReader("{\"BasketID\":\"xyzgsfajd\",\"Name\":\"xyz\",\"OrderInfo\":{\"Exchange\":\"nse\",\"InstrumentToken\":\"gfa\",\"ClientID\":\"client1\",\"OrderType\":\"regular\",\"Price\":1111,\"Quantity\":1,\"DisclosedQuantity\":11,\"Validity\":\"tue\",\"Product\":\"good\",\"TradingSymbol\":\"axt\",\"OrderSide\":\"done\",\"UserOrderID\":13242,\"UnderlyingToken\":\"ok\",\"Series\":\"agsa\",\"TriggerPrice\":12,\"ExecutionType\":\"pending\"}}"))
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
			AddBasketInstrument(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("add basket instrument() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/addBasketInstrument", strings.NewReader("{\"BasketID\":\"xyzgsfajd\",\"Name\":\"xyz\",\"OrderInfo\":{\"Exchange\":\"nse\",\"InstrumentToken\":\"gfa\",\"ClientID\":\"client1\",\"OrderType\":\"regular\",\"Price\":1111,\"Quantity\":1,\"DisclosedQuantity\":11,\"Validity\":\"tue\",\"Product\":\"good\",\"TradingSymbol\":\"axt\",\"OrderSide\":\"done\",\"UserOrderID\":13242,\"UnderlyingToken\":\"ok\",\"Series\":\"agsa\",\"TriggerPrice\":12,\"ExecutionType\":\"pending\"}}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitBasketOrderProvider(basketOrderMock{})

	//mock business layer response
	addBasketInstrumentMock = func(req models.AddBasketInstrumentReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			AddBasketInstrument(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("add basket instrument() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}

}

func TestEditBasketInstrument(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/editBasketInstrument", nil)
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
			EditBasketInstrument(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("edit basket instrument() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/editBasketInstrument", strings.NewReader("{\"BasketID\":\"xyzgsfajd\",\"Name\":\"xyz\",\"OrderID\":\"abc\",\"OrderInfo\":{\"Exchange\":\"nse\",\"InstrumentToken\":\"gfa\",\"ClientID\":\"client1\",\"OrderType\":\"regular\",\"Price\":1111,\"Quantity\":1,\"DisclosedQuantity\":11,\"Validity\":\"tue\",\"Product\":\"good\",\"TradingSymbol\":\"axt\",\"OrderSide\":\"done\",\"UserOrderID\":13242,\"UnderlyingToken\":\"ok\",\"Series\":\"agsa\",\"OmsOrderID\":\"abf\",\"ExchangeOrderID\":\"fhgs\",\"TriggerPrice\":12,\"ExecutionType\":\"pending\"}}"))
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
			EditBasketInstrument(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("edit basket instrument() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/editBasketInstrument", strings.NewReader("{\"BasketID\":\"xyzgsfajd\",\"Name\":\"xyz\",\"OrderID\":\"abc\",\"OrderInfo\":{\"Exchange\":\"nse\",\"InstrumentToken\":\"gfa\",\"ClientID\":\"client1\",\"OrderType\":\"regular\",\"Price\":1111,\"Quantity\":1,\"DisclosedQuantity\":11,\"Validity\":\"tue\",\"Product\":\"good\",\"TradingSymbol\":\"axt\",\"OrderSide\":\"done\",\"UserOrderID\":13242,\"UnderlyingToken\":\"ok\",\"Series\":\"agsa\",\"OmsOrderID\":\"abf\",\"ExchangeOrderID\":\"fhgs\",\"TriggerPrice\":12,\"ExecutionType\":\"pending\"}}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitBasketOrderProvider(basketOrderMock{})

	//mock business layer response
	editBasketInstrumentMock = func(req models.EditBasketInstrumentReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			EditBasketInstrument(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("edit basket instrument() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}

}

func TestDeleteBasketInstrument(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/deleteBasketInstrument", nil)
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
			DeleteBasketInstrument(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("delete basket instrument() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/deleteBasketInstrument", strings.NewReader("{\"BasketID\":\"xyzgsfajd\",\"OrderID\":\"hdsdse\",\"Name\":\"xyz\"}"))
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
			DeleteBasketInstrument(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("delete basket instrument() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/deleteBasketInstrument", strings.NewReader("{\"BasketID\":\"xyzgsfajd\",\"OrderID\":\"hdsdse\",\"Name\":\"xyz\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitBasketOrderProvider(basketOrderMock{})

	//mock business layer response
	deleteBasketInstrumentMock = func(req models.DeleteBasketInstrumentReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			DeleteBasketInstrument(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("delete basket instrument() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}

}

func TestRenameBasket(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/renameBasket", nil)
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
			RenameBasket(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("rename basket() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/renameBasket", strings.NewReader("{\"BasketID\":\"xyzgsfajd\",\"Name\":\"xyz\"}"))
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
			RenameBasket(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("delete basket instrument() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/renameBasket", strings.NewReader("{\"BasketID\":\"xyzgsfajd\",\"Name\":\"xyz\"}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitBasketOrderProvider(basketOrderMock{})

	//mock business layer response
	renameBasketMock = func(req models.RenameBasketReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			RenameBasket(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("rename basket() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}

}

func TestExecuteBasket(t *testing.T) {
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
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/executeBasket", nil)
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
			ExecuteBasket(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("execute basket() = %v, want %v", string(b), expected)
			}
		})
	}

	// 2 invalid device id
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/executeBasket", strings.NewReader("{\"BasketID\":\"xyzgsfajd\",\"Name\":\"xyz\",\"ExecutionType\":\"done\",\"SquareOff\":true,\"ClientID\":\"client1\",\"ExecutionState\":true}"))
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
			ExecuteBasket(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("execute basket() = %v, want %v", string(b), expected)
			}
		})
	}

	// //3 success
	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	reqHsuccess := models.ReqHeader{}
	reqHsuccess.DeviceType = "kugbkjbwVJKABWVLAWKJ"
	ctx.Set("reqH", reqHsuccess)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/basket/executeBasket", strings.NewReader("{\"BasketID\":\"xyzgsfajd\",\"Name\":\"xyz\",\"ExecutionType\":\"done\",\"SquareOff\":true,\"ClientID\":\"client1\",\"ExecutionState\":true}"))
	expected = "{\"status\":true,\"message\":\"SUCCESS\",\"errorcode\":\"\",\"data\":null}"

	//Init order provider
	InitBasketOrderProvider(basketOrderMock{})

	//mock business layer response
	executeBasketMock = func(req models.ExecuteBasketReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
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
			ExecuteBasket(tt.args.c)
			b, _ := ioutil.ReadAll(w.Body)
			if strings.TrimSuffix(string(b), "\n") != expected {
				t.Errorf("execute basket() = %v, want %v", string(b), expected)
			}
			if w.Code != 200 {
				t.Errorf("logoutController() = %v, want %v", w.Code, 200)
			}
		})
	}

}
