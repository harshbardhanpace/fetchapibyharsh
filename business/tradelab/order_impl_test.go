package tradelab

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers/cache"
	"space/loggerconfig"
	"space/models"

	"github.com/go-redis/redismock/v9"
)

func TestOrderObj_PlaceOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.PlaceOrderRequest
		reqH models.ReqHeader
	}

	db, mock := redismock.NewClientMock()
	redisClient := &cache.RedisClient{
		Client:      db,
		OrderClient: db,
	}
	cache.SetRedisClientObj(redisClient)

	mock.ExpectIncr("userOrderId").SetVal(12334)

	// test 1 start
	field1 := fields{
		tradeLabURL: "http://test",
	}

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	req1 := models.PlaceOrderRequest{
		ClientID:          "test123",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OrderSide:         "BUY",
		OrderType:         "LIMIT",
		Price:             2400.0,
		Product:           "MIS",
		Quantity:          2,
		TriggerPrice:      2300.0,
		Validity:          "DAY",
	}

	reqH1 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2f",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH1,
	}

	res1 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.InternalServerError],
		ErrorCode: constants.InternalServerError,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var res *http.Response
		return res, errors.New("Call Api Error")
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"call api error", field1, arg1, http.StatusInternalServerError, res1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.PlaceOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.PlaceOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.PlaceOrderRequest{
		ClientID:          "test123",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OrderSide:         "BUY",
		OrderType:         "LIMIT",
		Price:             2401.0,
		Product:           "MIS",
		Quantity:          2,
		TriggerPrice:      2301.0,
		Validity:          "DAY",
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	res2 := apihelpers.APIRes{
		Status:    false,
		Message:   "error",
		ErrorCode: "123",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradeLabErrorRes
		tradelabError.Status = "error"
		tradelabError.Message = "error"
		tradelabError.ErrorCode = 123
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab error", field2, arg2, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.PlaceOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.PlaceOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.PlaceOrderRequest{
		ClientID:          "test124",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OrderSide:         "BUY",
		OrderType:         "LIMIT",
		Price:             2402.0,
		Product:           "MIS",
		Quantity:          2,
		TriggerPrice:      2302.0,
		Validity:          "DAY",
	}

	reqH3 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg3 := args{
		req:  req3,
		reqH: reqH3,
	}

	res3 := apihelpers.APIRes{
		Status:  false,
		Message: "error",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradelabPlaceOrderResponse
		tradelabError.Message = "error"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusInternalServerError
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab status not ok", field3, arg3, http.StatusInternalServerError, res3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.PlaceOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.PlaceOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.PlaceOrderRequest{
		ClientID:          "test124",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OrderSide:         "BUY",
		OrderType:         "LIMIT",
		Price:             2402.0,
		Product:           "MIS",
		Quantity:          2,
		TriggerPrice:      2302.0,
		Validity:          "DAY",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var orderSuccess models.PlaceOrderResponse
	orderSuccess.OmsOrderID = "122334"
	orderSuccess.UserOrderID = 12334
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradelabPlaceOrderResponse
		tradelabSuccess.Data.OmsOrderID = "122334"
		tradelabSuccess.Data.UserOrderID = 12334
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabSuccess)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", field4, arg4, http.StatusOK, res4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.PlaceOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.PlaceOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end
}

func TestOrderObj_PlaceAMOOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.PlaceOrderRequest
		reqH models.ReqHeader
	}

	db, mock := redismock.NewClientMock()
	redisClient := &cache.RedisClient{
		Client:      db,
		OrderClient: db,
	}
	cache.SetRedisClientObj(redisClient)

	mock.ExpectIncr("userOrderId").SetVal(12334)

	// test 1 start
	field1 := fields{
		tradeLabURL: "http://test",
	}

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	req1 := models.PlaceOrderRequest{
		ClientID:          "test123",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OrderSide:         "BUY",
		OrderType:         "LIMIT",
		Price:             2400.0,
		Product:           "MIS",
		Quantity:          2,
		TriggerPrice:      2300.0,
		Validity:          "DAY",
	}

	reqH1 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2f",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH1,
	}

	res1 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.InternalServerError],
		ErrorCode: constants.InternalServerError,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var res *http.Response
		return res, errors.New("Call Api Error")
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"call api error", field1, arg1, http.StatusInternalServerError, res1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceAMOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.PlaceAMOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.PlaceAMOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.PlaceOrderRequest{
		ClientID:          "test123",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OrderSide:         "BUY",
		OrderType:         "LIMIT",
		Price:             2401.0,
		Product:           "MIS",
		Quantity:          2,
		TriggerPrice:      2301.0,
		Validity:          "DAY",
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	res2 := apihelpers.APIRes{
		Status:    false,
		Message:   "error",
		ErrorCode: "123",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradeLabErrorRes
		tradelabError.Status = "error"
		tradelabError.Message = "error"
		tradelabError.ErrorCode = 123
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab error", field2, arg2, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceAMOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.PlaceAMOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.PlaceAMOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.PlaceOrderRequest{
		ClientID:          "test124",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OrderSide:         "BUY",
		OrderType:         "LIMIT",
		Price:             2402.0,
		Product:           "MIS",
		Quantity:          2,
		TriggerPrice:      2302.0,
		Validity:          "DAY",
	}

	reqH3 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg3 := args{
		req:  req3,
		reqH: reqH3,
	}

	res3 := apihelpers.APIRes{
		Status:  false,
		Message: "error",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradelabAMOResponse
		tradelabError.Message = "error"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusInternalServerError
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab status not ok", field3, arg3, http.StatusInternalServerError, res3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceAMOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.PlaceAMOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.PlaceAMOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.PlaceOrderRequest{
		ClientID:          "test124",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OrderSide:         "BUY",
		OrderType:         "LIMIT",
		Price:             2402.0,
		Product:           "MIS",
		Quantity:          2,
		TriggerPrice:      2302.0,
		Validity:          "DAY",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var orderSuccess models.AMOOrderResponse
	orderSuccess.OmsOrderID = "122334"
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradelabAMOResponse
		tradelabSuccess.Data.OmsOrderID = "122334"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabSuccess)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", field4, arg4, http.StatusOK, res4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceAMOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.PlaceAMOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.PlaceAMOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end
}

func TestOrderObj_ModifyOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.ModifyOrderRequest
		reqH models.ReqHeader
	}

	db, mock := redismock.NewClientMock()
	redisClient := &cache.RedisClient{
		Client:      db,
		OrderClient: db,
	}
	cache.SetRedisClientObj(redisClient)

	mock.ExpectIncr("userOrderId").SetVal(12334)

	// test 1 start
	field1 := fields{
		tradeLabURL: "http://test",
	}

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	req1 := models.ModifyOrderRequest{
		ClientID:          "test123",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OmsOrderID:        "20220922-4",
		OrderType:         "LIMIT",
		Price:             2400.0,
		Product:           "CNC",
		Quantity:          2,
		TriggerPrice:      2300.0,
		Validity:          "DAY",
	}

	reqH1 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2f",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH1,
	}

	res1 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.InternalServerError],
		ErrorCode: constants.InternalServerError,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var res *http.Response
		return res, errors.New("Call Api Error")
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"call api error", field1, arg1, http.StatusInternalServerError, res1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifyOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.ModifyOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.ModifyOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.ModifyOrderRequest{
		ClientID:          "test123",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OmsOrderID:        "20220922-4",
		OrderType:         "LIMIT",
		Price:             2401.0,
		Product:           "MIS",
		Quantity:          2,
		TriggerPrice:      2301.0,
		Validity:          "DAY",
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	res2 := apihelpers.APIRes{
		Status:    false,
		Message:   "error",
		ErrorCode: "123",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradeLabErrorRes
		tradelabError.Status = "error"
		tradelabError.Message = "error"
		tradelabError.ErrorCode = 123
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab error", field2, arg2, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifyOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.ModifyOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.ModifyOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.ModifyOrderRequest{
		ClientID:          "test124",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OmsOrderID:        "20220922-4",
		OrderType:         "LIMIT",
		Price:             2402.0,
		Product:           "MIS",
		Quantity:          2,
		TriggerPrice:      2302.0,
		Validity:          "DAY",
	}

	reqH3 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg3 := args{
		req:  req3,
		reqH: reqH3,
	}

	res3 := apihelpers.APIRes{
		Status:  false,
		Message: "error",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradelabPlaceOrderResponse
		tradelabError.Message = "error"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusInternalServerError
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab status not ok", field3, arg3, http.StatusInternalServerError, res3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifyOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.ModifyOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.ModifyOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.ModifyOrderRequest{
		ClientID:          "test124",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OmsOrderID:        "20220922-4",
		OrderType:         "LIMIT",
		Price:             2402.0,
		Product:           "MIS",
		Quantity:          2,
		TriggerPrice:      2302.0,
		Validity:          "DAY",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var orderSuccess models.ModifyOrCancelOrderResponse
	orderSuccess.OmsOrderID = "122334"
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradelabPlaceOrderResponse
		tradelabSuccess.Data.OmsOrderID = "122334"
		tradelabSuccess.Data.UserOrderID = 12334
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabSuccess)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", field4, arg4, http.StatusOK, res4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifyOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.ModifyOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.ModifyOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end
}

func TestOrderObj_ModifyAMOOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.ModifyAMORequest
		reqH models.ReqHeader
	}

	db, mock := redismock.NewClientMock()
	redisClient := &cache.RedisClient{
		Client:      db,
		OrderClient: db,
	}
	cache.SetRedisClientObj(redisClient)

	mock.ExpectIncr("userOrderId").SetVal(12334)

	// test 1 start
	field1 := fields{
		tradeLabURL: "http://test",
	}

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	req1 := models.ModifyAMORequest{
		ClientID:          "test123",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OmsOrderID:        "20220922-4",
		OrderType:         "LIMIT",
		Price:             2400.0,
		Product:           "CNC",
		Quantity:          2,
		TriggerPrice:      2300.0,
		Validity:          "DAY",
	}

	reqH1 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2f",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH1,
	}

	res1 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.InternalServerError],
		ErrorCode: constants.InternalServerError,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var res *http.Response
		return res, errors.New("Call Api Error")
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"call api error", field1, arg1, http.StatusInternalServerError, res1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifyAMOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.ModifyAMOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.ModifyAMOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.ModifyAMORequest{
		ClientID:          "test123",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OmsOrderID:        "20220922-4",
		OrderType:         "LIMIT",
		Price:             2401.0,
		Product:           "MIS",
		Quantity:          2,
		TriggerPrice:      2301.0,
		Validity:          "DAY",
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	res2 := apihelpers.APIRes{
		Status:    false,
		Message:   "error",
		ErrorCode: "123",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradeLabErrorRes
		tradelabError.Status = "error"
		tradelabError.Message = "error"
		tradelabError.ErrorCode = 123
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab error", field2, arg2, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifyAMOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.ModifyAMOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.ModifyAMOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.ModifyAMORequest{
		ClientID:          "test124",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OmsOrderID:        "20220922-4",
		OrderType:         "LIMIT",
		Price:             2402.0,
		Product:           "MIS",
		Quantity:          2,
		TriggerPrice:      2302.0,
		Validity:          "DAY",
	}

	reqH3 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg3 := args{
		req:  req3,
		reqH: reqH3,
	}

	res3 := apihelpers.APIRes{
		Status:  false,
		Message: "error",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradelabAMOResponse
		tradelabError.Message = "error"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusInternalServerError
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab status not ok", field3, arg3, http.StatusInternalServerError, res3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifyAMOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.ModifyAMOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.ModifyAMOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.ModifyAMORequest{
		ClientID:          "test124",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OmsOrderID:        "20220922-4",
		OrderType:         "LIMIT",
		Price:             2402.0,
		Product:           "MIS",
		Quantity:          2,
		TriggerPrice:      2302.0,
		Validity:          "DAY",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var orderSuccess models.AMOOrderResponse
	orderSuccess.OmsOrderID = "122334"
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradelabAMOResponse
		tradelabSuccess.Data.OmsOrderID = "122334"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabSuccess)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", field4, arg4, http.StatusOK, res4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifyAMOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.ModifyAMOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.ModifyAMOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end
}

func TestOrderObj_CancelOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.CancelOrderRequest
		reqH models.ReqHeader
	}

	db, mock := redismock.NewClientMock()
	redisClient := &cache.RedisClient{
		Client:      db,
		OrderClient: db,
	}
	cache.SetRedisClientObj(redisClient)

	mock.ExpectIncr("userOrderId").SetVal(12334)

	// test 1 start
	field1 := fields{
		tradeLabURL: "http://test",
	}

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	req1 := models.CancelOrderRequest{
		ClientID:      "test123",
		OmsOrderId:    "20220922-4",
		ExecutionType: "REGULAR",
	}

	reqH1 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2f",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH1,
	}

	res1 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.InternalServerError],
		ErrorCode: constants.InternalServerError,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var res *http.Response
		return res, errors.New("Call Api Error")
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"call api error", field1, arg1, http.StatusInternalServerError, res1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.CancelOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.CancelOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.CancelOrderRequest{
		ClientID:      "test123",
		OmsOrderId:    "20220922-5",
		ExecutionType: "REGULAR",
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	res2 := apihelpers.APIRes{
		Status:    false,
		Message:   "error",
		ErrorCode: "123",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradeLabErrorRes
		tradelabError.Status = "error"
		tradelabError.Message = "error"
		tradelabError.ErrorCode = 123
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab error", field2, arg2, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.CancelOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.CancelOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.CancelOrderRequest{
		ClientID:      "test123",
		OmsOrderId:    "20220922-6",
		ExecutionType: "REGULAR",
	}

	reqH3 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg3 := args{
		req:  req3,
		reqH: reqH3,
	}

	res3 := apihelpers.APIRes{
		Status:  false,
		Message: "error",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradelabPlaceOrderResponse
		tradelabError.Message = "error"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusInternalServerError
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab status not ok", field3, arg3, http.StatusInternalServerError, res3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.CancelOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.CancelOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.CancelOrderRequest{
		ClientID:      "test123",
		OmsOrderId:    "20220922-7",
		ExecutionType: "REGULAR",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var orderSuccess models.ModifyOrCancelOrderResponse
	orderSuccess.OmsOrderID = "122334"
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradelabPlaceOrderResponse
		tradelabSuccess.Data.OmsOrderID = "122334"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabSuccess)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", field4, arg4, http.StatusOK, res4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.CancelOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.CancelOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end
}

func TestOrderObj_CancelAMOOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.CancelOrderRequest
		reqH models.ReqHeader
	}

	db, mock := redismock.NewClientMock()
	redisClient := &cache.RedisClient{
		Client:      db,
		OrderClient: db,
	}
	cache.SetRedisClientObj(redisClient)

	mock.ExpectIncr("userOrderId").SetVal(12334)

	// test 1 start
	field1 := fields{
		tradeLabURL: "http://test",
	}

	req1 := models.CancelOrderRequest{
		ClientID:      "test123",
		OmsOrderId:    "20220922-4",
		ExecutionType: "REGULAR",
	}

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	reqH1 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2f",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH1,
	}

	res1 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.InternalServerError],
		ErrorCode: constants.InternalServerError,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var res *http.Response
		return res, errors.New("Call Api Error")
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"call api error", field1, arg1, http.StatusInternalServerError, res1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelAMOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.CancelAMOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.CancelAMOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.CancelOrderRequest{
		ClientID:      "test123",
		OmsOrderId:    "20220922-5",
		ExecutionType: "REGULAR",
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	res2 := apihelpers.APIRes{
		Status:    false,
		Message:   "error",
		ErrorCode: "123",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradeLabErrorRes
		tradelabError.Status = "error"
		tradelabError.Message = "error"
		tradelabError.ErrorCode = 123
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab error", field2, arg2, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelAMOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.CancelAMOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.CancelAMOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.CancelOrderRequest{
		ClientID:      "test123",
		OmsOrderId:    "20220922-6",
		ExecutionType: "REGULAR",
	}

	reqH3 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg3 := args{
		req:  req3,
		reqH: reqH3,
	}

	res3 := apihelpers.APIRes{
		Status:  false,
		Message: "error",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradelabCancelOrModifyResponse
		tradelabError.Message = "error"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusInternalServerError
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab status not ok", field3, arg3, http.StatusInternalServerError, res3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelAMOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.CancelAMOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.CancelAMOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.CancelOrderRequest{
		ClientID:      "test123",
		OmsOrderId:    "20220922-7",
		ExecutionType: "REGULAR",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var orderSuccess models.ModifyOrCancelOrderResponse
	orderSuccess.OmsOrderID = "122334"
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradelabCancelOrModifyResponse
		tradelabSuccess.Data.OmsOrderID = "122334"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabSuccess)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", field4, arg4, http.StatusOK, res4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelAMOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.CancelAMOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.CancelAMOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end
}

func TestOrderObj_PlaceGTTOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.CreateGTTOrderRequest
		reqH models.ReqHeader
	}

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	db, mock := redismock.NewClientMock()
	redisClient := &cache.RedisClient{
		Client:      db,
		OrderClient: db,
	}
	cache.SetRedisClientObj(redisClient)

	mock.ExpectIncr("userOrderId").SetVal(12334)

	// test 1 start
	field1 := fields{
		tradeLabURL: "http://test",
	}

	req1 := models.CreateGTTOrderRequest{
		ActionType:                 "single_order",
		ExpiryTime:                 "2022-10-17",
		ClientID:                   "Client1",
		DisclosedQuantity:          0,
		Exchange:                   "NSE",
		InstrumentToken:            "22",
		MarketProtectionPercentage: 0,
		OrderSide:                  "BUY",
		OrderType:                  "CNC",
		Quantity:                   1,
		SlOrderPrice:               0,
		SlOrderQuantity:            0,
		SlTriggerPrice:             0,
		TriggerPrice:               3000,
		UserOrderID:                12334,
	}

	reqH1 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2f",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH1,
	}

	res1 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.InternalServerError],
		ErrorCode: constants.InternalServerError,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var res *http.Response
		return res, errors.New("Call Api Error")
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"call api error", field1, arg1, http.StatusInternalServerError, res1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceGTTOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.PlaceGTTOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.PlaceGTTOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.CreateGTTOrderRequest{
		ActionType:                 "single_order",
		ExpiryTime:                 "2022-10-17",
		ClientID:                   "Client1",
		DisclosedQuantity:          0,
		Exchange:                   "NSE",
		InstrumentToken:            "22",
		MarketProtectionPercentage: 0,
		OrderSide:                  "BUY",
		OrderType:                  "CNC",
		Quantity:                   1,
		SlOrderPrice:               0,
		SlOrderQuantity:            0,
		SlTriggerPrice:             0,
		TriggerPrice:               3000,
		UserOrderID:                12334,
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	res2 := apihelpers.APIRes{
		Status:    false,
		Message:   "error",
		ErrorCode: "123",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradeLabErrorRes
		tradelabError.Status = "error"
		tradelabError.Message = "error"
		tradelabError.ErrorCode = 123
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab error", field2, arg2, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceGTTOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.PlaceGTTOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.PlaceGTTOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.CreateGTTOrderRequest{
		ActionType:                 "single_order",
		ExpiryTime:                 "2022-10-17",
		ClientID:                   "Client1",
		DisclosedQuantity:          0,
		Exchange:                   "NSE",
		InstrumentToken:            "22",
		MarketProtectionPercentage: 0,
		OrderSide:                  "BUY",
		OrderType:                  "CNC",
		Quantity:                   1,
		SlOrderPrice:               0,
		SlOrderQuantity:            0,
		SlTriggerPrice:             0,
		TriggerPrice:               3000,
		UserOrderID:                12334,
	}

	reqH3 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg3 := args{
		req:  req3,
		reqH: reqH3,
	}

	res3 := apihelpers.APIRes{
		Status:  false,
		Message: "error",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradelabGTTOrderResponse
		tradelabError.Message = "error"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusInternalServerError
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab status not ok", field3, arg3, http.StatusInternalServerError, res3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceGTTOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.PlaceGTTOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.PlaceGTTOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.CreateGTTOrderRequest{
		ActionType:                 "single_order",
		ExpiryTime:                 "2022-10-17",
		ClientID:                   "Client1",
		DisclosedQuantity:          0,
		Exchange:                   "NSE",
		InstrumentToken:            "22",
		MarketProtectionPercentage: 0,
		OrderSide:                  "BUY",
		OrderType:                  "CNC",
		Quantity:                   1,
		SlOrderPrice:               0,
		SlOrderQuantity:            0,
		SlTriggerPrice:             0,
		TriggerPrice:               3000,
		UserOrderID:                12334,
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var orderSuccess models.GTTOrderResponse
	orderSuccess.ID = "122334"
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradelabGTTOrderResponse
		tradelabSuccess.Data.ID = "122334"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabSuccess)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", field4, arg4, http.StatusOK, res4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceGTTOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.PlaceGTTOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.PlaceGTTOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end
}

func TestOrderObj_ModifyGTTOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.ModifyGTTOrderRequest
		reqH models.ReqHeader
	}

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	db, mock := redismock.NewClientMock()
	redisClient := &cache.RedisClient{
		Client:      db,
		OrderClient: db,
	}
	cache.SetRedisClientObj(redisClient)

	mock.ExpectIncr("userOrderId").SetVal(12334)

	// test 1 start
	field1 := fields{
		tradeLabURL: "http://test",
	}

	req1 := models.ModifyGTTOrderRequest{
		ExpiryTime: "2022-10-17",
		ActionType: "single_order",
		ID:         "testId",
		Order: models.ModifyGTTOrderDetails{
			ClientID:                   "Client1",
			DisclosedQuantity:          0,
			Exchange:                   "NSE",
			InstrumentToken:            "22",
			MarketProtectionPercentage: 0,
			OrderType:                  "CNC",
			Price:                      500,
			Product:                    "MIS",
			Quantity:                   1,
			SlOrderPrice:               0,
			SlOrderQuantity:            0,
			SlTriggerPrice:             0,
			TriggerPrice:               3000,
			UserOrderID:                12334,
		},
	}

	reqH1 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2f",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH1,
	}

	res1 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.InternalServerError],
		ErrorCode: constants.InternalServerError,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var res *http.Response
		return res, errors.New("Call Api Error")
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"call api error", field1, arg1, http.StatusInternalServerError, res1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifyGTTOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.ModifyGTTOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.ModifyGTTOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.ModifyGTTOrderRequest{
		ExpiryTime: "2022-10-17",
		ActionType: "single_order",
		ID:         "testId",
		Order: models.ModifyGTTOrderDetails{
			ClientID:                   "Client1",
			DisclosedQuantity:          0,
			Exchange:                   "NSE",
			InstrumentToken:            "22",
			MarketProtectionPercentage: 0,
			OrderType:                  "CNC",
			Price:                      500,
			Product:                    "MIS",
			Quantity:                   1,
			SlOrderPrice:               0,
			SlOrderQuantity:            0,
			SlTriggerPrice:             0,
			TriggerPrice:               3000,
			UserOrderID:                12334,
		},
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	res2 := apihelpers.APIRes{
		Status:    false,
		Message:   "error",
		ErrorCode: "123",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradeLabErrorRes
		tradelabError.Status = "error"
		tradelabError.Message = "error"
		tradelabError.ErrorCode = 123
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab error", field2, arg2, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifyGTTOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.ModifyGTTOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.ModifyGTTOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.ModifyGTTOrderRequest{
		ExpiryTime: "2022-10-17",
		ActionType: "single_order",
		ID:         "testId",
		Order: models.ModifyGTTOrderDetails{
			ClientID:                   "Client1",
			DisclosedQuantity:          0,
			Exchange:                   "NSE",
			InstrumentToken:            "22",
			MarketProtectionPercentage: 0,
			OrderType:                  "CNC",
			Price:                      500,
			Product:                    "MIS",
			Quantity:                   1,
			SlOrderPrice:               0,
			SlOrderQuantity:            0,
			SlTriggerPrice:             0,
			TriggerPrice:               3000,
			UserOrderID:                12334,
		},
	}

	reqH3 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg3 := args{
		req:  req3,
		reqH: reqH3,
	}

	res3 := apihelpers.APIRes{
		Status:  false,
		Message: "error",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradelabGTTOrderResponse
		tradelabError.Message = "error"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusInternalServerError
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab status not ok", field3, arg3, http.StatusInternalServerError, res3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifyGTTOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.ModifyGTTOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.ModifyGTTOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.ModifyGTTOrderRequest{
		ExpiryTime: "2022-10-17",
		ActionType: "single_order",
		ID:         "testId",
		Order: models.ModifyGTTOrderDetails{
			ClientID:                   "Client1",
			DisclosedQuantity:          0,
			Exchange:                   "NSE",
			InstrumentToken:            "22",
			MarketProtectionPercentage: 0,
			OrderType:                  "CNC",
			Price:                      500,
			Product:                    "MIS",
			Quantity:                   1,
			SlOrderPrice:               0,
			SlOrderQuantity:            0,
			SlTriggerPrice:             0,
			TriggerPrice:               3000,
			UserOrderID:                12334,
		},
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var orderSuccess models.GTTOrderResponse
	orderSuccess.ID = "122334"
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradelabGTTOrderResponse
		tradelabSuccess.Data.ID = "122334"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabSuccess)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", field4, arg4, http.StatusOK, res4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifyGTTOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.ModifyGTTOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.ModifyGTTOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end
}

func TestOrderObj_CancelGTTOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.CancelGTTOrderRequest
		reqH models.ReqHeader
	}

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	db, mock := redismock.NewClientMock()
	redisClient := &cache.RedisClient{
		Client:      db,
		OrderClient: db,
	}
	cache.SetRedisClientObj(redisClient)

	mock.ExpectIncr("userOrderId").SetVal(12334)

	// test 1 start
	field1 := fields{
		tradeLabURL: "http://test",
	}

	req1 := models.CancelGTTOrderRequest{
		ClientId: "Client1",
		Id:       "11234",
	}

	reqH1 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2f",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH1,
	}

	res1 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.InternalServerError],
		ErrorCode: constants.InternalServerError,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var res *http.Response
		return res, errors.New("Call Api Error")
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"call api error", field1, arg1, http.StatusInternalServerError, res1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelGTTOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.CancelGTTOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.CancelGTTOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.CancelGTTOrderRequest{
		ClientId: "Client1",
		Id:       "11234",
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	res2 := apihelpers.APIRes{
		Status:    false,
		Message:   "error",
		ErrorCode: "123",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradeLabErrorRes
		tradelabError.Status = "error"
		tradelabError.Message = "error"
		tradelabError.ErrorCode = 123
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab error", field2, arg2, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelGTTOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.CancelGTTOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.CancelGTTOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.CancelGTTOrderRequest{
		ClientId: "Client1",
		Id:       "11234",
	}

	reqH3 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg3 := args{
		req:  req3,
		reqH: reqH3,
	}

	res3 := apihelpers.APIRes{
		Status:  false,
		Message: "error",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradelabGTTOrderResponse
		tradelabError.Message = "error"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusInternalServerError
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab status not ok", field3, arg3, http.StatusInternalServerError, res3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelGTTOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.CancelGTTOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.CancelGTTOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.CancelGTTOrderRequest{
		ClientId: "Client1",
		Id:       "11234",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var orderSuccess models.GTTOrderResponse
	orderSuccess.ID = "122334"
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradelabGTTOrderResponse
		tradelabSuccess.Data.ID = "122334"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabSuccess)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", field4, arg4, http.StatusOK, res4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelGTTOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.CancelGTTOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.CancelGTTOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end
}

func TestOrderObj_FetchGTTOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.FetchGTTOrderRequest
		reqH models.ReqHeader
	}

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	db, mock := redismock.NewClientMock()
	redisClient := &cache.RedisClient{
		Client:      db,
		OrderClient: db,
	}
	cache.SetRedisClientObj(redisClient)

	mock.ExpectIncr("userOrderId").SetVal(12334)

	// test 1 start
	field1 := fields{
		tradeLabURL: "http://test",
	}

	req1 := models.FetchGTTOrderRequest{
		ClientId: "Client1",
	}

	reqH1 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2f",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH1,
	}

	res1 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.InternalServerError],
		ErrorCode: constants.InternalServerError,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var res *http.Response
		return res, errors.New("Call Api Error")
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"call api error", field1, arg1, http.StatusInternalServerError, res1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.FetchGTTOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.FetchGTTOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.FetchGTTOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.FetchGTTOrderRequest{
		ClientId: "Client1",
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	res2 := apihelpers.APIRes{
		Status:    false,
		Message:   "error",
		ErrorCode: "123",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradeLabErrorRes
		tradelabError.Status = "error"
		tradelabError.Message = "error"
		tradelabError.ErrorCode = 123
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab error", field2, arg2, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.FetchGTTOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.FetchGTTOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.FetchGTTOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.FetchGTTOrderRequest{
		ClientId: "Client1",
	}

	reqH3 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg3 := args{
		req:  req3,
		reqH: reqH3,
	}

	res3 := apihelpers.APIRes{
		Status:  false,
		Message: "error",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradelabFetchGTTOrderResponse
		tradelabError.Message = "error"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusInternalServerError
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab status not ok", field3, arg3, http.StatusInternalServerError, res3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.FetchGTTOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.FetchGTTOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.FetchGTTOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.FetchGTTOrderRequest{
		ClientId: "Client1",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var fetchGTTOrderRes models.FetchGTTOrderResponse

	fetchGttOrderDataResAll := make([]models.FetchGTTOrderResponseData, 0)

	for i := 0; i < 1; i++ {
		var fetchGttOrderDataRes models.FetchGTTOrderResponseData
		fetchGttOrderDataResAll = append(fetchGttOrderDataResAll, fetchGttOrderDataRes)
	}
	fetchGTTOrderRes.FetchGTTOrderData = fetchGttOrderDataResAll

	res4 := apihelpers.APIRes{
		Data:    fetchGTTOrderRes,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradelabFetchGTTOrderResponse
		tradelabSuccessData := make([]TradelabFetchGTTOrderResponseData, 1)
		tradelabSuccess.Data = tradelabSuccessData
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabSuccess)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", field4, arg4, http.StatusOK, res4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.FetchGTTOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.FetchGTTOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.FetchGTTOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end
}

func TestOrderObj_PendingOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.PendingOrderRequest
		reqH models.ReqHeader
	}

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	db, mock := redismock.NewClientMock()
	redisClient := &cache.RedisClient{
		Client:      db,
		OrderClient: db,
	}
	cache.SetRedisClientObj(redisClient)

	mock.ExpectIncr("userOrderId").SetVal(12334)

	// test 1 start
	field1 := fields{
		tradeLabURL: "http://test",
	}

	req1 := models.PendingOrderRequest{
		Type:     "pending",
		ClientID: "Client1",
	}

	reqH1 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2f",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH1,
	}

	res1 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.InternalServerError],
		ErrorCode: constants.InternalServerError,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var res *http.Response
		return res, errors.New("Call Api Error")
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"call api error", field1, arg1, http.StatusInternalServerError, res1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PendingOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.PendingOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.PendingOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.PendingOrderRequest{
		Type:     "pending",
		ClientID: "Client1",
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	res2 := apihelpers.APIRes{
		Status:    false,
		Message:   "error",
		ErrorCode: "123",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradeLabErrorRes
		tradelabError.Status = "error"
		tradelabError.Message = "error"
		tradelabError.ErrorCode = 123
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab error", field2, arg2, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PendingOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.PendingOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.PendingOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.PendingOrderRequest{
		Type:     "pending",
		ClientID: "Client1",
	}

	reqH3 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg3 := args{
		req:  req3,
		reqH: reqH3,
	}

	res3 := apihelpers.APIRes{
		Status:  false,
		Message: "error",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradelabPendingOrderResponse
		tradelabError.Message = "error"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusInternalServerError
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab status not ok", field3, arg3, http.StatusInternalServerError, res3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PendingOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.PendingOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.PendingOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.PendingOrderRequest{
		Type:     "pending",
		ClientID: "Client1",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var orderSuccess models.PendingOrderResponse
	responseOrders := make([]models.PendingOrderResponseOrders, 0)
	for i := 0; i < 1; i++ {
		var pendingOrderResOrders models.PendingOrderResponseOrders
		pendingOrderResOrders.TradingSymbol = "PNB-EQ"
		pendingOrderResOrders.AverageTradePrice = 10
		pendingOrderResOrders.Exchange = "NSE"
		pendingOrderResOrders.ProCli = "CLIENT"
		pendingOrderResOrders.MarketProtectionPercentage = 0
		pendingOrderResOrders.OrderEntryTime = 161045
		pendingOrderResOrders.Mode = "NEW"
		pendingOrderResOrders.OmsOrderID = "20220922-4"
		//pendingOrderResOrders.TrailingStopLoss = {}
		pendingOrderResOrders.Deposit = 0
		//pendingOrderResOrders.SquareOffValue = {}
		pendingOrderResOrders.DisclosedQuantity = 0
		//pendingOrderResOrders.StopLossValue = {}
		pendingOrderResOrders.Price = 500
		pendingOrderResOrders.OrderTag = ""
		pendingOrderResOrders.Device = "WEB"
		pendingOrderResOrders.RemainingQuantity = 1
		pendingOrderResOrders.LastActivityReference = 12345678
		pendingOrderResOrders.AveragePrice = 10
		pendingOrderResOrders.SquareOff = false
		pendingOrderResOrders.OrderStatusInfo = ""
		pendingOrderResOrders.Quantity = 1
		pendingOrderResOrders.ExecutionType = "REGULAR"
		pendingOrderResOrders.ClientID = "CLIENT1"
		pendingOrderResOrders.ExchangeTime = 153045
		pendingOrderResOrders.OrderSide = "BUY"
		pendingOrderResOrders.LoginID = "Client1"
		pendingOrderResOrders.Validity = "DAY"
		pendingOrderResOrders.InstrumentToken = 22
		pendingOrderResOrders.Product = "MIS"
		pendingOrderResOrders.TriggerPrice = 0
		pendingOrderResOrders.Segment = ""
		pendingOrderResOrders.TradePrice = 0
		pendingOrderResOrders.OrderType = "LIMIT"
		//pendingOrderResOrders.ContractDescription = {}
		pendingOrderResOrders.RejectionCode = 0
		pendingOrderResOrders.LegOrderIndicator = ""
		pendingOrderResOrders.ExchangeOrderID = "12345678"
		pendingOrderResOrders.OrderStatus = "CONFIRMED"
		pendingOrderResOrders.FilledQuantity = 0
		pendingOrderResOrders.TargetPriceType = "absolute"
		pendingOrderResOrders.IsTrailing = false
		pendingOrderResOrders.UserOrderID = "12334"
		pendingOrderResOrders.LotSize = 1
		pendingOrderResOrders.Series = ""
		pendingOrderResOrders.NnfID = 11111
		pendingOrderResOrders.RejectionReason = "NONE"
		responseOrders = append(responseOrders, pendingOrderResOrders)
	}
	orderSuccess.Orders = responseOrders
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradelabPendingOrderResponse
		responseOrders := make([]TradelabPendingOrderResponseOrders, 0)
		for i := 0; i < 1; i++ {
			var pendingOrderResOrders TradelabPendingOrderResponseOrders
			pendingOrderResOrders.TradingSymbol = "PNB-EQ"
			pendingOrderResOrders.AverageTradePrice = 10
			pendingOrderResOrders.Exchange = "NSE"
			pendingOrderResOrders.ProCli = "CLIENT"
			pendingOrderResOrders.MarketProtectionPercentage = 0
			pendingOrderResOrders.OrderEntryTime = 161045
			pendingOrderResOrders.Mode = "NEW"
			pendingOrderResOrders.OmsOrderID = "20220922-4"
			//pendingOrderResOrders.TrailingStopLoss = {}
			pendingOrderResOrders.Deposit = 0
			//pendingOrderResOrders.SquareOffValue = {}
			pendingOrderResOrders.DisclosedQuantity = 0
			//pendingOrderResOrders.StopLossValue = {}
			pendingOrderResOrders.Price = 500
			pendingOrderResOrders.OrderTag = ""
			pendingOrderResOrders.Device = "WEB"
			pendingOrderResOrders.RemainingQuantity = 1
			pendingOrderResOrders.LastActivityReference = 12345678
			pendingOrderResOrders.AveragePrice = 10
			pendingOrderResOrders.SquareOff = false
			pendingOrderResOrders.OrderStatusInfo = ""
			pendingOrderResOrders.Quantity = 1
			pendingOrderResOrders.ExecutionType = "REGULAR"
			pendingOrderResOrders.ClientID = "CLIENT1"
			pendingOrderResOrders.ExchangeTime = 153045
			pendingOrderResOrders.OrderSide = "BUY"
			pendingOrderResOrders.LoginID = "Client1"
			pendingOrderResOrders.Validity = "DAY"
			pendingOrderResOrders.InstrumentToken = 22
			pendingOrderResOrders.Product = "MIS"
			pendingOrderResOrders.TriggerPrice = 0
			pendingOrderResOrders.Segment = ""
			pendingOrderResOrders.TradePrice = 0
			pendingOrderResOrders.OrderType = "LIMIT"
			//pendingOrderResOrders.ContractDescription = {}
			pendingOrderResOrders.RejectionCode = 0
			pendingOrderResOrders.LegOrderIndicator = ""
			pendingOrderResOrders.ExchangeOrderID = "12345678"
			pendingOrderResOrders.OrderStatus = "CONFIRMED"
			pendingOrderResOrders.FilledQuantity = 0
			pendingOrderResOrders.TargetPriceType = "absolute"
			pendingOrderResOrders.IsTrailing = false
			pendingOrderResOrders.UserOrderID = "12334"
			pendingOrderResOrders.LotSize = 1
			pendingOrderResOrders.Series = ""
			pendingOrderResOrders.NnfID = 11111
			pendingOrderResOrders.RejectionReason = "NONE"
			responseOrders = append(responseOrders, pendingOrderResOrders)
		}
		var Data TradeLabPendingOrderResponseData
		Data.Orders = responseOrders
		tradelabSuccess.Data = Data
		tradelabSuccess.Message = "SUCCESS"
		tradelabSuccess.Status = "true"

		var res http.Response
		jsonRes, _ := json.Marshal(tradelabSuccess)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", field4, arg4, http.StatusOK, res4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PendingOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.PendingOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.PendingOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end
}

func TestOrderObj_CompletedOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.CompletedOrderRequest
		reqH models.ReqHeader
	}

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	db, mock := redismock.NewClientMock()
	redisClient := &cache.RedisClient{
		Client:      db,
		OrderClient: db,
	}
	cache.SetRedisClientObj(redisClient)

	mock.ExpectIncr("userOrderId").SetVal(12334)

	// test 1 start
	field1 := fields{
		tradeLabURL: "http://test",
	}

	req1 := models.CompletedOrderRequest{
		Type:     "completed",
		ClientID: "Client1",
	}

	reqH1 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2f",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH1,
	}

	res1 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.InternalServerError],
		ErrorCode: constants.InternalServerError,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var res *http.Response
		return res, errors.New("Call Api Error")
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"call api error", field1, arg1, http.StatusInternalServerError, res1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CompletedOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.CompletedOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.CompletedOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.CompletedOrderRequest{
		Type:     "completed",
		ClientID: "Client1",
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	res2 := apihelpers.APIRes{
		Status:    false,
		Message:   "error",
		ErrorCode: "123",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradeLabErrorRes
		tradelabError.Status = "error"
		tradelabError.Message = "error"
		tradelabError.ErrorCode = 123
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab error", field2, arg2, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CompletedOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.CompletedOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.CompletedOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.CompletedOrderRequest{
		Type:     "completed",
		ClientID: "Client1",
	}

	reqH3 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg3 := args{
		req:  req3,
		reqH: reqH3,
	}

	res3 := apihelpers.APIRes{
		Status:  false,
		Message: "error",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradelabCompletedOrderResponse
		tradelabError.Message = "error"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusInternalServerError
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab status not ok", field3, arg3, http.StatusInternalServerError, res3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CompletedOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.CompletedOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.CompletedOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.CompletedOrderRequest{
		Type:     "completed",
		ClientID: "Client1",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var orderSuccess models.CompletedOrderResponse
	responseOrders := make([]models.CompletedOrderResponseOrders, 0)
	for i := 0; i < 1; i++ {
		var completedOrderResOrders models.CompletedOrderResponseOrders
		completedOrderResOrders.TradingSymbol = "PNB-EQ"
		completedOrderResOrders.AverageTradePrice = 10
		completedOrderResOrders.Exchange = "NSE"
		completedOrderResOrders.ProCli = "CLIENT"
		completedOrderResOrders.MarketProtectionPercentage = 0
		completedOrderResOrders.OrderEntryTime = 161045
		completedOrderResOrders.Mode = "NEW"
		completedOrderResOrders.OmsOrderID = "20220922-4"
		//completedOrderResOrders.TrailingStopLoss = {}
		completedOrderResOrders.Deposit = 0
		//completedOrderResOrders.SquareOffValue = {}
		completedOrderResOrders.DisclosedQuantity = 0
		//completedOrderResOrders.StopLossValue = {}
		completedOrderResOrders.Price = 500
		completedOrderResOrders.OrderTag = ""
		completedOrderResOrders.Device = "WEB"
		completedOrderResOrders.RemainingQuantity = 1
		completedOrderResOrders.LastActivityReference = 12345678
		completedOrderResOrders.AveragePrice = 10
		completedOrderResOrders.SquareOff = false
		completedOrderResOrders.OrderStatusInfo = ""
		completedOrderResOrders.Quantity = 1
		completedOrderResOrders.ExecutionType = "REGULAR"
		completedOrderResOrders.ClientID = "CLIENT1"
		completedOrderResOrders.ExchangeTime = 153045
		completedOrderResOrders.OrderSide = "BUY"
		completedOrderResOrders.LoginID = "Client1"
		completedOrderResOrders.Validity = "DAY"
		completedOrderResOrders.InstrumentToken = 22
		completedOrderResOrders.Product = "MIS"
		completedOrderResOrders.TriggerPrice = 0
		completedOrderResOrders.Segment = ""
		completedOrderResOrders.TradePrice = 0
		completedOrderResOrders.OrderType = "LIMIT"
		//completedOrderResOrders.ContractDescription = {}
		completedOrderResOrders.RejectionCode = 0
		completedOrderResOrders.LegOrderIndicator = ""
		completedOrderResOrders.ExchangeOrderID = "12345678"
		completedOrderResOrders.OrderStatus = "CONFIRMED"
		completedOrderResOrders.FilledQuantity = 0
		completedOrderResOrders.TargetPriceType = "absolute"
		completedOrderResOrders.IsTrailing = false
		completedOrderResOrders.UserOrderID = "12334"
		completedOrderResOrders.LotSize = 1
		completedOrderResOrders.Series = ""
		completedOrderResOrders.NnfID = 11111
		completedOrderResOrders.RejectionReason = "NONE"
		responseOrders = append(responseOrders, completedOrderResOrders)
	}
	orderSuccess.Orders = responseOrders
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradelabCompletedOrderResponse
		responseOrders := make([]TradelabCompletedOrderResponseOrders, 0)
		for i := 0; i < 1; i++ {
			var completedOrderResOrders TradelabCompletedOrderResponseOrders
			completedOrderResOrders.TradingSymbol = "PNB-EQ"
			completedOrderResOrders.AverageTradePrice = 10
			completedOrderResOrders.Exchange = "NSE"
			completedOrderResOrders.ProCli = "CLIENT"
			completedOrderResOrders.MarketProtectionPercentage = 0
			completedOrderResOrders.OrderEntryTime = 161045
			completedOrderResOrders.Mode = "NEW"
			completedOrderResOrders.OmsOrderID = "20220922-4"
			//completedOrderResOrders.TrailingStopLoss = {}
			completedOrderResOrders.Deposit = 0
			//completedOrderResOrders.SquareOffValue = {}
			completedOrderResOrders.DisclosedQuantity = 0
			//completedOrderResOrders.StopLossValue = {}
			completedOrderResOrders.Price = 500
			completedOrderResOrders.OrderTag = ""
			completedOrderResOrders.Device = "WEB"
			completedOrderResOrders.RemainingQuantity = 1
			completedOrderResOrders.LastActivityReference = 12345678
			completedOrderResOrders.AveragePrice = 10
			completedOrderResOrders.SquareOff = false
			completedOrderResOrders.OrderStatusInfo = ""
			completedOrderResOrders.Quantity = 1
			completedOrderResOrders.ExecutionType = "REGULAR"
			completedOrderResOrders.ClientID = "CLIENT1"
			completedOrderResOrders.ExchangeTime = 153045
			completedOrderResOrders.OrderSide = "BUY"
			completedOrderResOrders.LoginID = "Client1"
			completedOrderResOrders.Validity = "DAY"
			completedOrderResOrders.InstrumentToken = 22
			completedOrderResOrders.Product = "MIS"
			completedOrderResOrders.TriggerPrice = 0
			completedOrderResOrders.Segment = ""
			completedOrderResOrders.TradePrice = 0
			completedOrderResOrders.OrderType = "LIMIT"
			//completedOrderResOrders.ContractDescription = {}
			completedOrderResOrders.RejectionCode = 0
			completedOrderResOrders.LegOrderIndicator = ""
			completedOrderResOrders.ExchangeOrderID = "12345678"
			completedOrderResOrders.OrderStatus = "CONFIRMED"
			completedOrderResOrders.FilledQuantity = 0
			completedOrderResOrders.TargetPriceType = "absolute"
			completedOrderResOrders.IsTrailing = false
			completedOrderResOrders.UserOrderID = "12334"
			completedOrderResOrders.LotSize = 1
			completedOrderResOrders.Series = ""
			completedOrderResOrders.NnfID = 11111
			completedOrderResOrders.RejectionReason = "NONE"
			responseOrders = append(responseOrders, completedOrderResOrders)
		}
		var Data TradelabCompletedOrderResponseData
		Data.Orders = responseOrders
		tradelabSuccess.Data = Data
		tradelabSuccess.Message = "SUCCESS"
		tradelabSuccess.Status = "true"

		var res http.Response
		jsonRes, _ := json.Marshal(tradelabSuccess)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", field4, arg4, http.StatusOK, res4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CompletedOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.CompletedOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.CompletedOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end
}

func TestOrderObj_TradeBook(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.TradeBookRequest
		reqH models.ReqHeader
	}

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	db, mock := redismock.NewClientMock()
	redisClient := &cache.RedisClient{
		Client:      db,
		OrderClient: db,
	}
	cache.SetRedisClientObj(redisClient)

	mock.ExpectIncr("userOrderId").SetVal(12334)

	// test 1 start
	field1 := fields{
		tradeLabURL: "http://test",
	}

	req1 := models.TradeBookRequest{
		ClientID: "Client1",
	}

	reqH1 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2f",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH1,
	}

	res1 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.InternalServerError],
		ErrorCode: constants.InternalServerError,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var res *http.Response
		return res, errors.New("Call Api Error")
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"call api error", field1, arg1, http.StatusInternalServerError, res1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.TradeBook(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.TradeBook() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.TradeBook() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.TradeBookRequest{
		ClientID: "Client1",
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	res2 := apihelpers.APIRes{
		Status:    false,
		Message:   "error",
		ErrorCode: "123",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradeLabErrorRes
		tradelabError.Status = "error"
		tradelabError.Message = "error"
		tradelabError.ErrorCode = 123
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab error", field2, arg2, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.TradeBook(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.TradeBook() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.TradeBook() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.TradeBookRequest{
		ClientID: "Client1",
	}

	reqH3 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg3 := args{
		req:  req3,
		reqH: reqH3,
	}

	res3 := apihelpers.APIRes{
		Status:  false,
		Message: "error",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradeLabTradeBookResponse
		tradelabError.Message = "error"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusInternalServerError
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab status not ok", field3, arg3, http.StatusInternalServerError, res3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.TradeBook(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.TradeBook() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.TradeBook() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.TradeBookRequest{
		ClientID: "Client1",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var orderSuccess models.TradeBookResponse
	responseOrders := make([]models.TradeBookResponseData, 0)
	for i := 0; i < 1; i++ {
		var tradeBookResponseData models.TradeBookResponseData
		tradeBookResponseData.BookType = "0"
		tradeBookResponseData.BrokerID = "12345"
		tradeBookResponseData.ClientID = "CLIENT1"
		tradeBookResponseData.DisclosedVol = 0
		tradeBookResponseData.DisclosedVolRemaining = 0
		tradeBookResponseData.Exchange = "NSE"
		tradeBookResponseData.ExchangeOrderID = "161045"
		tradeBookResponseData.ExchangeTime = 1605677129
		tradeBookResponseData.FilledQuantity = 1
		//tradeBookResponseData.FillNumber = {}
		tradeBookResponseData.GoodTillDate = 0
		tradeBookResponseData.InstrumentToken = 22
		tradeBookResponseData.LoginID = "CLIENT1"
		tradeBookResponseData.OmsOrderID = "20220923-4"
		tradeBookResponseData.OrderEntryTime = 1605677128
		tradeBookResponseData.OrderPrice = 20
		tradeBookResponseData.OrderSide = "BUY"
		tradeBookResponseData.OrderType = "MARKET"
		tradeBookResponseData.OriginalVol = 1
		tradeBookResponseData.Pan = "ABCD10"
		tradeBookResponseData.ProCli = 0
		tradeBookResponseData.Product = "CNC"
		tradeBookResponseData.RemainingQuantity = 1
		tradeBookResponseData.TradeNumber = "76817533"
		tradeBookResponseData.TradePrice = 20
		tradeBookResponseData.TradeQuantity = 15
		tradeBookResponseData.TradeTime = 1605677129
		tradeBookResponseData.TradingSymbol = "ACC-EQ"
		//tradeBookResponseData.TriggerPrice = {}
		//tradeBookResponseData.VLoginID = {}
		tradeBookResponseData.VolFilledToday = 1
		responseOrders = append(responseOrders, tradeBookResponseData)
	}
	orderSuccess.Trades = responseOrders
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradeLabTradeBookResponse
		responseOrders := make([]TradeLabTradeBookResponseTrades, 0)
		for i := 0; i < 1; i++ {
			var tradeBookResponseData TradeLabTradeBookResponseTrades
			tradeBookResponseData.BookType = "0"
			tradeBookResponseData.BrokerID = "12345"
			tradeBookResponseData.ClientID = "CLIENT1"
			tradeBookResponseData.DisclosedVol = 0
			tradeBookResponseData.DisclosedVolRemaining = 0
			tradeBookResponseData.Exchange = "NSE"
			tradeBookResponseData.ExchangeOrderID = "161045"
			tradeBookResponseData.ExchangeTime = 1605677129
			tradeBookResponseData.FilledQuantity = 1
			//tradeBookResponseData.FillNumber = {}
			tradeBookResponseData.GoodTillDate = 0
			tradeBookResponseData.InstrumentToken = 22
			tradeBookResponseData.LoginID = "CLIENT1"
			tradeBookResponseData.OmsOrderID = "20220923-4"
			tradeBookResponseData.OrderEntryTime = 1605677128
			tradeBookResponseData.OrderPrice = 20
			tradeBookResponseData.OrderSide = "BUY"
			tradeBookResponseData.OrderType = "MARKET"
			tradeBookResponseData.OriginalVol = 1
			tradeBookResponseData.Pan = "ABCD10"
			tradeBookResponseData.ProCli = 0
			tradeBookResponseData.Product = "CNC"
			tradeBookResponseData.RemainingQuantity = 1
			tradeBookResponseData.TradeNumber = "76817533"
			tradeBookResponseData.TradePrice = 20
			tradeBookResponseData.TradeQuantity = 15
			tradeBookResponseData.TradeTime = 1605677129
			tradeBookResponseData.TradingSymbol = "ACC-EQ"
			//tradeBookResponseData.TriggerPrice = {}
			//tradeBookResponseData.VLoginID = {}
			tradeBookResponseData.VolFilledToday = 1
			responseOrders = append(responseOrders, tradeBookResponseData)
		}
		var Data TradeLabTradeBookResponseData
		Data.Trades = responseOrders
		tradelabSuccess.Data = Data
		tradelabSuccess.Message = "SUCCESS"
		tradelabSuccess.Status = "true"

		var res http.Response
		jsonRes, _ := json.Marshal(tradelabSuccess)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", field4, arg4, http.StatusOK, res4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.TradeBook(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.TradeBook() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.TradeBook() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end
}

func TestOrderObj_OrderHistory(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.OrderHistoryRequest
		reqH models.ReqHeader
	}

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	db, mock := redismock.NewClientMock()
	redisClient := &cache.RedisClient{
		Client:      db,
		OrderClient: db,
	}
	cache.SetRedisClientObj(redisClient)

	mock.ExpectIncr("userOrderId").SetVal(12334)

	// test 1 start
	field1 := fields{
		tradeLabURL: "http://test",
	}

	req1 := models.OrderHistoryRequest{
		OmsOrderID: "20220923-5",
		ClientID:   "Client1",
	}

	reqH1 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2f",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH1,
	}

	res1 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.InternalServerError],
		ErrorCode: constants.InternalServerError,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var res *http.Response
		return res, errors.New("Call Api Error")
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"call api error", field1, arg1, http.StatusInternalServerError, res1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.OrderHistory(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.OrderHistory() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.OrderHistory() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.OrderHistoryRequest{
		OmsOrderID: "20220923-5",
		ClientID:   "Client1",
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	res2 := apihelpers.APIRes{
		Status:    false,
		Message:   "error",
		ErrorCode: "123",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradeLabErrorRes
		tradelabError.Status = "error"
		tradelabError.Message = "error"
		tradelabError.ErrorCode = 123
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab error", field2, arg2, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.OrderHistory(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.OrderHistory() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.OrderHistory() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.OrderHistoryRequest{
		OmsOrderID: "20220923-5",
		ClientID:   "Client1",
	}

	reqH3 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg3 := args{
		req:  req3,
		reqH: reqH3,
	}

	res3 := apihelpers.APIRes{
		Status:  false,
		Message: "error",
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabError TradeLabOrderHistoryResponse
		tradelabError.Message = "error"
		var res http.Response
		jsonRes, _ := json.Marshal(tradelabError)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusInternalServerError
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Tradelab status not ok", field3, arg3, http.StatusInternalServerError, res3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.OrderHistory(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.OrderHistory() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.OrderHistory() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.OrderHistoryRequest{
		OmsOrderID: "20220923-5",
		ClientID:   "Client1",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var orderSuccess models.OrderHistoryResponse
	responseOrders := make([]models.OrderHistoryResponseData, 0)
	for i := 0; i < 1; i++ {
		var orderHistoryResponseData models.OrderHistoryResponseData
		orderHistoryResponseData.AvgPrice = 40657
		orderHistoryResponseData.ClientID = "CLIENT1"
		orderHistoryResponseData.ClientOrderID = "900005"
		orderHistoryResponseData.CreatedAt = 1605683202
		orderHistoryResponseData.DisclosedQuantity = 0
		orderHistoryResponseData.Exchange = "NSE"
		orderHistoryResponseData.ExchangeOrderID = "161045"
		orderHistoryResponseData.ExchangeTime = 1605677129
		orderHistoryResponseData.FillQuantity = 1
		orderHistoryResponseData.LastModified = 1605683202477299000
		orderHistoryResponseData.LoginID = "CLIENT1"
		orderHistoryResponseData.ModifiedAt = 1605683202
		orderHistoryResponseData.OrderID = "20220923-4"
		orderHistoryResponseData.OrderMode = "NEW"
		orderHistoryResponseData.OrderSide = "SELL"
		orderHistoryResponseData.OrderType = "LIMIT"
		orderHistoryResponseData.Price = 20
		orderHistoryResponseData.Product = "ACC"
		orderHistoryResponseData.Quantity = 10
		orderHistoryResponseData.RejectReason = "NONE"
		orderHistoryResponseData.RemainingQuantity = 1
		orderHistoryResponseData.Segment = "FullOTP"
		orderHistoryResponseData.Status = "Complete"
		orderHistoryResponseData.Symbol = "ACC-EQ"
		orderHistoryResponseData.Token = 1605677129
		orderHistoryResponseData.TriggerPrice = 400
		orderHistoryResponseData.UnderlyingToken = 22
		orderHistoryResponseData.Validity = "DAY"
		responseOrders = append(responseOrders, orderHistoryResponseData)
	}
	orderSuccess.OrderHistory = responseOrders
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradeLabOrderHistoryResponse
		responseOrders := make([]TradeLabOrderHistoryResponseData, 0)
		for i := 0; i < 1; i++ {
			var orderHistoryResponseData TradeLabOrderHistoryResponseData
			orderHistoryResponseData.AvgPrice = 40657
			orderHistoryResponseData.ClientID = "CLIENT1"
			orderHistoryResponseData.ClientOrderID = "900005"
			orderHistoryResponseData.CreatedAt = 1605683202
			orderHistoryResponseData.DisclosedQuantity = 0
			orderHistoryResponseData.Exchange = "NSE"
			orderHistoryResponseData.ExchangeOrderID = "161045"
			orderHistoryResponseData.ExchangeTime = 1605677129
			orderHistoryResponseData.FillQuantity = 1
			orderHistoryResponseData.LastModified = 1605683202
			orderHistoryResponseData.LoginID = "CLIENT1"
			orderHistoryResponseData.ModifiedAt = 1605683202
			orderHistoryResponseData.OrderID = "20220923-4"
			orderHistoryResponseData.OrderMode = "NEW"
			orderHistoryResponseData.OrderSide = "SELL"
			orderHistoryResponseData.OrderType = "LIMIT"
			orderHistoryResponseData.Price = 20
			orderHistoryResponseData.Product = "ACC"
			orderHistoryResponseData.Quantity = 10
			orderHistoryResponseData.RejectReason = "NONE"
			orderHistoryResponseData.RemainingQuantity = 1
			orderHistoryResponseData.Segment = "FullOTP"
			orderHistoryResponseData.Status = "Complete"
			orderHistoryResponseData.Symbol = "ACC-EQ"
			orderHistoryResponseData.Token = 1605677129
			orderHistoryResponseData.TriggerPrice = 400
			orderHistoryResponseData.UnderlyingToken = 22
			orderHistoryResponseData.Validity = "DAY"
			responseOrders = append(responseOrders, orderHistoryResponseData)
		}
		tradelabSuccess.Data = responseOrders
		tradelabSuccess.Message = "SUCCESS"
		tradelabSuccess.Status = "true"

		var res http.Response
		jsonRes, _ := json.Marshal(tradelabSuccess)
		res.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonRes)))
		res.StatusCode = http.StatusOK
		return &res, nil
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", field4, arg4, http.StatusOK, res4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := OrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.OrderHistory(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("OrderObj.OrderHistory() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderObj.OrderHistory() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end
}

func TestOrderObj_LastTradedPrice(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.LastTradedPriceRequest
		reqH models.ReqHeader
	}

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	// Test 1: MCX Error
	field1 := fields{
		tradeLabURL: "http://test",
	}

	req1 := models.LastTradedPriceRequest{
		Exchange: "MCX",
		Segment:  "COMMODITY",
		Token:    "invalid_token",
	}

	reqH1 := models.ReqHeader{
		DeviceType:    "mobile",
		Authorization: "Bearer invalid_token",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH1,
	}

	res1 := apihelpers.APIRes{
		Status:    false,
		Message:   "Call API Error",
		ErrorCode: strconv.Itoa(http.StatusInternalServerError), // Expecting 500 from mocked error
	}

	// Mock CallAPIFunc for MCX Error
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var mcxErrorResponse TradeLabErrorRes
		mcxErrorResponse.Status = "error"
		mcxErrorResponse.Message = "Call API Error"
		mcxErrorResponse.ErrorCode = 500

		jsonRes, _ := json.Marshal(mcxErrorResponse)
		res := &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(bytes.NewReader(jsonRes)),
		}
		return res, nil
	}

	t.Run("MCX Error", func(t *testing.T) {
		obj := OrderObj{
			tradeLabURL: field1.tradeLabURL,
		}
		got, got1 := obj.LastTradedPrice(arg1.req, arg1.reqH)
		if got != http.StatusInternalServerError {
			t.Errorf("OrderObj.LastTradedPrice() got = %v, want %v", got, http.StatusInternalServerError)
		}
		if !reflect.DeepEqual(got1, res1) {
			t.Errorf("OrderObj.LastTradedPrice() got1 = %v, want1 %v", got1, res1)
		}
	})

	// Test 2: MCX Success
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.LastTradedPriceRequest{
		Exchange: "MCX",
		Segment:  "COMMODITY",
		Token:    "123456",
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "mobile",
		Authorization: "Bearer valid_token",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	res2 := apihelpers.APIRes{
		Status:  true,
		Message: "SUCCESS",
		Data: models.LastTradedPriceResponse{
			Price: 0,
		},
	}

	// Mock CallAPIFunc for MCX Success
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var mcxResponse TradeLabMCXLastTradedPrice
		mcxResponse.Data.LastTradePrice = 75000
		mcxResponse.Status = "success"
		mcxResponse.Message = "SUCCESS"

		jsonRes, _ := json.Marshal(mcxResponse)
		res := &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(jsonRes)),
		}
		return res, nil
	}

	t.Run("MCX Success", func(t *testing.T) {
		obj := OrderObj{
			tradeLabURL: field2.tradeLabURL,
		}
		got, got1 := obj.LastTradedPrice(arg2.req, arg2.reqH)
		if got != http.StatusOK {
			t.Errorf("OrderObj.LastTradedPrice() got = %v, want %v", got, http.StatusOK)
		}
		if !reflect.DeepEqual(got1, res2) {
			t.Errorf("OrderObj.LastTradedPrice() got1 = %v, want1 %v", got1, res2)
		}
	})
}
