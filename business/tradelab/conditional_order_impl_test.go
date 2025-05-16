package tradelab

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers/cache"
	"space/loggerconfig"
	"space/models"

	"github.com/go-redis/redismock/v9"
)

func TestConditionalOrderObj_PlaceBOOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.PlaceBOOrderRequest
		reqH models.ReqHeader
	}

	db, mock := redismock.NewClientMock()
	redisClient := &cache.RedisClient{
		Client:      db,
		OrderClient: db,
	}
	cache.SetRedisClientObj(redisClient)

	mock.ExpectIncr("userOrderId").SetVal(12334)

	loggerconfig.Info = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Error = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	loggerconfig.Warn = func(args ...interface{}) {
		fmt.Println("DO Nothing")
	}

	field1 := fields{
		tradeLabURL: "http://test",
	}

	req1 := models.PlaceBOOrderRequest{
		ClientID:          "CLIENT1",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		IsTrailing:        true,
		OrderSide:         "BUY",
		OrderType:         "LIMIT",
		Price:             34.2,
		Product:           "CNC",
		Quantity:          2,
		SquareOffValue:    45,
		StopLossValue:     23,
		TrailingStopLoss:  "22",
		TriggerPrice:      44,
		UserOrderID:       1027109,
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
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceId, userAgent, remoteAddr, authToken string) (*http.Response, error) {
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceBOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.PlaceBOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.PlaceBOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.PlaceBOOrderRequest{
		ClientID:          "CLIENT1",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		IsTrailing:        true,
		OrderSide:         "BUY",
		OrderType:         "LIMIT",
		Price:             34.2,
		Product:           "CNC",
		Quantity:          2,
		SquareOffValue:    45,
		StopLossValue:     23,
		TrailingStopLoss:  "22",
		TriggerPrice:      44,
		UserOrderID:       1027109,
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceBOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.PlaceBOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.PlaceBOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.PlaceBOOrderRequest{
		ClientID:          "CLIENT1",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		IsTrailing:        true,
		OrderSide:         "BUY",
		OrderType:         "LIMIT",
		Price:             34.2,
		Product:           "CNC",
		Quantity:          2,
		SquareOffValue:    45,
		StopLossValue:     23,
		TrailingStopLoss:  "22",
		TriggerPrice:      44,
		UserOrderID:       1027109,
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
		var tradelabError TradeLabPlaceBOOrderResponse
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceBOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.PlaceBOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.PlaceBOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.PlaceBOOrderRequest{
		ClientID:          "CLIENT1",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		IsTrailing:        true,
		OrderSide:         "BUY",
		OrderType:         "LIMIT",
		Price:             34.2,
		Product:           "CNC",
		Quantity:          2,
		SquareOffValue:    45,
		StopLossValue:     23,
		TrailingStopLoss:  "22",
		TriggerPrice:      44,
		UserOrderID:       1027109,
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

	var orderSuccess models.BOOrderResponse
	orderSuccess.BasketID = "122334"
	orderSuccess.Message = "Placed"
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradeLabPlaceBOOrderResponse
		tradelabSuccess.Data.Data.BasketID = "122334"
		tradelabSuccess.Data.Data.Message = "Placed"
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceBOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.PlaceBOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.PlaceBOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end

}

func TestConditionalOrderObj_ModifyBOOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.ModifyBOOrderRequest
		reqH models.ReqHeader
	}

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

	req1 := models.ModifyBOOrderRequest{
		ClientID:              "CLIENT1",
		DisclosedQuantity:     2,
		Exchange:              "NSE",
		ExchangeOrderID:       "ionwdg123",
		FilledQuantity:        1,
		InstrumentToken:       "22",
		IsTrailing:            true,
		LastActivityReference: 1325938440097498600,
		OmsOrderID:            "123445",
		OrderType:             "LIMIT",
		Price:                 34.2,
		Product:               "CNC",
		Quantity:              2,
		RemainingQuantity:     0,
		SquareOffValue:        40,
		StopLossValue:         39,
		TrailingStopLoss:      41,
		TriggerPrice:          38,
		Validity:              "DAY",
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifyBOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.ModifyBOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.ModifyBOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.ModifyBOOrderRequest{
		ClientID:              "CLIENT1",
		DisclosedQuantity:     2,
		Exchange:              "NSE",
		ExchangeOrderID:       "ionwdg123",
		FilledQuantity:        1,
		InstrumentToken:       "22",
		IsTrailing:            true,
		LastActivityReference: 1325938440097498600,
		OmsOrderID:            "123445",
		OrderType:             "LIMIT",
		Price:                 34.2,
		Product:               "CNC",
		Quantity:              2,
		RemainingQuantity:     0,
		SquareOffValue:        40,
		StopLossValue:         39,
		TrailingStopLoss:      41,
		TriggerPrice:          38,
		Validity:              "DAY",
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifyBOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.ModifyBOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.ModifyBOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	// //test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.ModifyBOOrderRequest{
		ClientID:              "CLIENT1",
		DisclosedQuantity:     2,
		Exchange:              "NSE",
		ExchangeOrderID:       "ionwdg123",
		FilledQuantity:        1,
		InstrumentToken:       "22",
		IsTrailing:            true,
		LastActivityReference: 1325938440097498600,
		OmsOrderID:            "123445",
		OrderType:             "LIMIT",
		Price:                 34.2,
		Product:               "CNC",
		Quantity:              2,
		RemainingQuantity:     0,
		SquareOffValue:        40,
		StopLossValue:         39,
		TrailingStopLoss:      41,
		TriggerPrice:          38,
		Validity:              "DAY",
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
		var tradelabError TradeLabModOrExitBOOrderResponse
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifyBOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.ModifyBOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.ModifyBOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.ModifyBOOrderRequest{
		ClientID:              "CLIENT1",
		DisclosedQuantity:     2,
		Exchange:              "NSE",
		ExchangeOrderID:       "ionwdg123",
		FilledQuantity:        1,
		InstrumentToken:       "22",
		IsTrailing:            true,
		LastActivityReference: 1325938440097498600,
		OmsOrderID:            "123445",
		OrderType:             "LIMIT",
		Price:                 34.2,
		Product:               "CNC",
		Quantity:              2,
		RemainingQuantity:     0,
		SquareOffValue:        40,
		StopLossValue:         39,
		TrailingStopLoss:      41,
		TriggerPrice:          38,
		Validity:              "DAY",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var orderSuccess models.SpreadOrderResponse
	orderSuccess.BasketID = "122334"
	orderSuccess.Message = "Placed"
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradeLabModOrExitBOOrderResponse
		tradelabSuccess.Data.BasketID = "122334"
		tradelabSuccess.Data.Message = "Placed"
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifyBOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.ModifyBOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.ModifyBOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end

}

func TestConditionalOrderObj_CancelBOOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.ExitBOOrderRequest
		reqH models.ReqHeader
	}

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

	req1 := models.ExitBOOrderRequest{
		ClientID:          "CLIENT1",
		ExchangeOrderID:   "ionwdg123",
		LegOrderIndicator: "Entry",
		OmsOrderID:        "123445",
		Status:            "CONFIRMED",
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelBOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.CancelBOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.CancelBOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.ExitBOOrderRequest{
		ClientID:          "CLIENT1",
		ExchangeOrderID:   "ionwdg123",
		LegOrderIndicator: "Entry",
		OmsOrderID:        "123445",
		Status:            "CONFIRMED",
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelBOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.CancelBOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.CancelBOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.ExitBOOrderRequest{
		ClientID:          "CLIENT1",
		ExchangeOrderID:   "ionwdg123",
		LegOrderIndicator: "Entry",
		OmsOrderID:        "123445",
		Status:            "CONFIRMED",
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
		var tradelabError TradeLabModOrExitBOOrderResponse
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelBOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.CancelBOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.CancelBOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.ExitBOOrderRequest{
		ClientID:          "CLIENT1",
		ExchangeOrderID:   "ionwdg123",
		LegOrderIndicator: "Entry",
		OmsOrderID:        "123445",
		Status:            "CONFIRMED",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var orderSuccess models.BOOrderResponse
	orderSuccess.BasketID = "122334"
	orderSuccess.Message = "Placed"
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradeLabModOrExitBOOrderResponse
		tradelabSuccess.Data.BasketID = "122334"
		tradelabSuccess.Data.Message = "Placed"
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelBOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.CancelBOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.CancelBOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end

}

func TestConditionalOrderObj_PlaceCOOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.PlaceCOOrderRequest
		reqH models.ReqHeader
	}

	db, mock := redismock.NewClientMock()
	redisClient := &cache.RedisClient{
		Client:      db,
		OrderClient: db,
	}
	cache.SetRedisClientObj(redisClient)

	mock.ExpectIncr("userOrderId").SetVal(12334)

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

	req1 := models.PlaceCOOrderRequest{
		ClientID:          "CLIENT1",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OrderSide:         "BUY",
		OrderType:         "LIMIT",
		Price:             34.2,
		Product:           "CNC",
		Quantity:          2,
		StopLossValue:     33,
		TrailingStopLoss:  33,
		UserOrderID:       91261928,
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceCOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.PlaceCOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.PlaceCOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.PlaceCOOrderRequest{
		ClientID:          "CLIENT1",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OrderSide:         "BUY",
		OrderType:         "LIMIT",
		Price:             34.2,
		Product:           "CNC",
		Quantity:          2,
		StopLossValue:     33,
		TrailingStopLoss:  33,
		UserOrderID:       91261928,
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceCOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.PlaceOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.PlaceOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.PlaceCOOrderRequest{
		ClientID:          "CLIENT1",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OrderSide:         "BUY",
		OrderType:         "LIMIT",
		Price:             34.2,
		Product:           "CNC",
		Quantity:          2,
		StopLossValue:     33,
		TrailingStopLoss:  33,
		UserOrderID:       91261928,
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
		var tradelabError TradeLabPlaceCOOrderResponse
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceCOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.PlaceCOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.PlaceCOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.PlaceCOOrderRequest{
		ClientID:          "CLIENT1",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OrderSide:         "BUY",
		OrderType:         "LIMIT",
		Price:             34.2,
		Product:           "CNC",
		Quantity:          2,
		StopLossValue:     33,
		TrailingStopLoss:  33,
		UserOrderID:       91261928,
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

	var orderSuccess models.COOrderResponse
	orderSuccess.BasketID = "122334"
	orderSuccess.Message = "Placed"
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradeLabPlaceCOOrderResponse
		tradelabSuccess.Data.Data.BasketID = "122334"
		tradelabSuccess.Data.Data.Message = "Placed"
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceCOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.PlaceCOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.PlaceCOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end

}

func TestConditionalOrderObj_ModifyCOOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.ModifyCOOrderRequest
		reqH models.ReqHeader
	}

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

	req1 := models.ModifyCOOrderRequest{
		ClientID:              "CLIENT1",
		DisclosedQuantity:     2,
		Exchange:              "NSE",
		ExchangeOrderID:       "ionwdg123",
		FilledQuantity:        1,
		InstrumentToken:       "22",
		LastActivityReference: 1325938440097498600,
		OmsOrderID:            "123445",
		OrderType:             "LIMIT",
		Price:                 34.2,
		Product:               "CNC",
		Quantity:              2,
		RemainingQuantity:     0,
		StopLossValue:         0,
		TrailingStopLoss:      33,
		Validity:              "DAY",
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifyCOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.ModifyCOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.ModifyCOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.ModifyCOOrderRequest{
		ClientID:              "CLIENT1",
		DisclosedQuantity:     2,
		Exchange:              "NSE",
		ExchangeOrderID:       "ionwdg123",
		FilledQuantity:        1,
		InstrumentToken:       "22",
		LastActivityReference: 1325938440097498600,
		OmsOrderID:            "123445",
		OrderType:             "LIMIT",
		Price:                 34.2,
		Product:               "CNC",
		Quantity:              2,
		RemainingQuantity:     0,
		StopLossValue:         0,
		TrailingStopLoss:      33,
		Validity:              "DAY",
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifyCOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.ModifyCOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.ModifyCOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	// //test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.ModifyCOOrderRequest{
		ClientID:              "CLIENT1",
		DisclosedQuantity:     2,
		Exchange:              "NSE",
		ExchangeOrderID:       "ionwdg123",
		FilledQuantity:        1,
		InstrumentToken:       "22",
		LastActivityReference: 1325938440097498600,
		OmsOrderID:            "123445",
		OrderType:             "LIMIT",
		Price:                 34.2,
		Product:               "CNC",
		Quantity:              2,
		RemainingQuantity:     0,
		StopLossValue:         0,
		TrailingStopLoss:      33,
		Validity:              "DAY",
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
		var tradelabError TradeLabModifyOrExitCOOrderResponse
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifyCOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.ModifyCOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.ModifyCOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.ModifyCOOrderRequest{
		ClientID:              "CLIENT1",
		DisclosedQuantity:     2,
		Exchange:              "NSE",
		ExchangeOrderID:       "ionwdg123",
		FilledQuantity:        1,
		InstrumentToken:       "22",
		LastActivityReference: 1325938440097498600,
		OmsOrderID:            "123445",
		OrderType:             "LIMIT",
		Price:                 34.2,
		Product:               "CNC",
		Quantity:              2,
		RemainingQuantity:     0,
		StopLossValue:         0,
		TrailingStopLoss:      33,
		Validity:              "DAY",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var orderSuccess models.COOrderResponse
	orderSuccess.BasketID = "122334"
	orderSuccess.Message = "Placed"
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradeLabModifyOrExitCOOrderResponse
		tradelabSuccess.Data.BasketID = "122334"
		tradelabSuccess.Data.Message = "Placed"
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifyCOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.ModifyCOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.ModifyCOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end

}

func TestConditionalOrderObj_CancelCOOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.ExitCOOrderRequest
		reqH models.ReqHeader
	}

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

	req1 := models.ExitCOOrderRequest{
		ClientID:          "CLIENT1",
		ExchangeOrderID:   "ionwdg123",
		LegOrderIndicator: "Entry",
		OmsOrderID:        "123445",
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelCOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.CancelCOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.CancelCOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.ExitCOOrderRequest{
		ClientID:          "CLIENT1",
		ExchangeOrderID:   "ionwdg123",
		LegOrderIndicator: "Entry",
		OmsOrderID:        "123445",
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelCOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.CancelCOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.CancelCOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.ExitCOOrderRequest{
		ClientID:          "CLIENT1",
		ExchangeOrderID:   "ionwdg123",
		LegOrderIndicator: "Entry",
		OmsOrderID:        "123445",
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
		var tradelabError TradeLabModifyOrExitCOOrderResponse
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelCOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.CancelCOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.CancelCOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.ExitCOOrderRequest{
		ClientID:          "CLIENT1",
		ExchangeOrderID:   "ionwdg123",
		LegOrderIndicator: "Entry",
		OmsOrderID:        "123445",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var orderSuccess models.COOrderResponse
	orderSuccess.BasketID = "122334"
	orderSuccess.Message = "Placed"
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradeLabModifyOrExitCOOrderResponse
		tradelabSuccess.Data.BasketID = "122334"
		tradelabSuccess.Data.Message = "Placed"
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelCOOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.CancelCOOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.CancelCOOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end

}

func TestConditionalOrderObj_PlaceSpreadOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.PlaceSpreadOrderRequest
		reqH models.ReqHeader
	}

	db, mock := redismock.NewClientMock()
	redisClient := &cache.RedisClient{
		Client:      db,
		OrderClient: db,
	}
	cache.SetRedisClientObj(redisClient)

	mock.ExpectIncr("userOrderId").SetVal(12334)

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

	req1 := models.PlaceSpreadOrderRequest{
		ClientID:          "CLIENT1",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OrderSide:         "BUY",
		OrderType:         "LIMIT",
		Price:             34.2,
		Product:           "CNC",
		Quantity:          2,
		UserOrderID:       91261928,
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceSpreadOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.PlaceSpreadOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.PlaceSpreadOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.PlaceSpreadOrderRequest{
		ClientID:          "CLIENT1",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OrderSide:         "BUY",
		OrderType:         "LIMIT",
		Price:             34.2,
		Product:           "CNC",
		Quantity:          2,
		UserOrderID:       91261928,
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceSpreadOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.PlaceSpreadOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.PlaceSpreadOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.PlaceSpreadOrderRequest{
		ClientID:          "CLIENT1",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OrderSide:         "BUY",
		OrderType:         "LIMIT",
		Price:             34.2,
		Product:           "CNC",
		Quantity:          2,
		UserOrderID:       91261928,
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
		var tradelabError TradelabSpreadOrderResponse
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceSpreadOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.PlaceSpreadOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.PlaceSpreadOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.PlaceSpreadOrderRequest{
		ClientID:          "CLIENT1",
		DisclosedQuantity: 2,
		Exchange:          "NSE",
		InstrumentToken:   "22",
		OrderSide:         "BUY",
		OrderType:         "LIMIT",
		Price:             34.2,
		Product:           "CNC",
		Quantity:          2,
		UserOrderID:       91261928,
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

	var orderSuccess models.SpreadOrderResponse
	orderSuccess.BasketID = "122334"
	orderSuccess.Message = "Placed"
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradelabSpreadOrderResponse
		tradelabSuccess.Data.Data.BasketID = "122334"
		tradelabSuccess.Data.Data.Message = "Placed"
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PlaceSpreadOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.PlaceSpreadOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.PlaceSpreadOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end

}

func TestConditionalOrderObj_ModifySpreadOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.ModifySpreadOrderRequest
		reqH models.ReqHeader
	}

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

	req1 := models.ModifySpreadOrderRequest{
		ClientID:          "CLIENT1",
		DisclosedQuantity: 0,
		Exchange:          "NSE",
		ExchangeOrderID:   "string",
		InstrumentToken:   "22",
		IsTrailing:        true,
		OmsOrderID:        "123445",
		OrderType:         "LIMIT",
		Price:             34.2,
		ProdType:          "string",
		Product:           "CNC",
		Quantity:          2,
		SquareOffValue:    40,
		StopLossValue:     39,
		TrailingStopLoss:  33,
		TriggerPrice:      38,
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifySpreadOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.ModifySpreadOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.ModifySpreadOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.ModifySpreadOrderRequest{
		ClientID:          "CLIENT1",
		DisclosedQuantity: 0,
		Exchange:          "NSE",
		ExchangeOrderID:   "string",
		InstrumentToken:   "22",
		IsTrailing:        true,
		OmsOrderID:        "123445",
		OrderType:         "LIMIT",
		Price:             34.2,
		ProdType:          "string",
		Product:           "CNC",
		Quantity:          2,
		SquareOffValue:    40,
		StopLossValue:     39,
		TrailingStopLoss:  33,
		TriggerPrice:      38,
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifySpreadOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.ModifySpreadOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.ModifySpreadOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	// //test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.ModifySpreadOrderRequest{
		ClientID:          "CLIENT1",
		DisclosedQuantity: 0,
		Exchange:          "NSE",
		ExchangeOrderID:   "string",
		InstrumentToken:   "22",
		IsTrailing:        true,
		OmsOrderID:        "123445",
		OrderType:         "LIMIT",
		Price:             34.2,
		ProdType:          "string",
		Product:           "CNC",
		Quantity:          2,
		SquareOffValue:    40,
		StopLossValue:     39,
		TrailingStopLoss:  33,
		TriggerPrice:      38,
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
		var tradelabError TradelabSpreadOrderResponse
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifySpreadOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.ModifySpreadOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.ModifySpreadOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.ModifySpreadOrderRequest{
		ClientID:          "CLIENT1",
		DisclosedQuantity: 0,
		Exchange:          "NSE",
		ExchangeOrderID:   "string",
		InstrumentToken:   "22",
		IsTrailing:        true,
		OmsOrderID:        "123445",
		OrderType:         "LIMIT",
		Price:             34.2,
		ProdType:          "string",
		Product:           "CNC",
		Quantity:          2,
		SquareOffValue:    40,
		StopLossValue:     39,
		TrailingStopLoss:  33,
		TriggerPrice:      38,
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

	var orderSuccess models.SpreadOrderResponse
	orderSuccess.BasketID = "122334"
	orderSuccess.Message = "Placed"
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradelabSpreadOrderResponse
		tradelabSuccess.Data.Data.BasketID = "122334"
		tradelabSuccess.Data.Data.Message = "Placed"
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ModifySpreadOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.ModifySpreadOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.ModifySpreadOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end

}

func TestConditionalOrderObj_CancelSpreadOrder(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.ExitSpreadOrderRequest
		reqH models.ReqHeader
	}

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

	req1 := models.ExitSpreadOrderRequest{
		ClientID:          "CLIENT1",
		ExchangeOrderID:   "ionwdg123",
		LegOrderIndicator: "Entry",
		OmsOrderID:        "123445",
		Status:            "CONFIRMED",
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelSpreadOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.CancelSpreadOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.CancelSpreadOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.ExitSpreadOrderRequest{
		ClientID:          "CLIENT1",
		ExchangeOrderID:   "ionwdg123",
		LegOrderIndicator: "Entry",
		OmsOrderID:        "123445",
		Status:            "CONFIRMED",
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelSpreadOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.CancelSpreadOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.CancelSpreadOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.ExitSpreadOrderRequest{
		ClientID:          "CLIENT1",
		ExchangeOrderID:   "ionwdg123",
		LegOrderIndicator: "Entry",
		OmsOrderID:        "123445",
		Status:            "CONFIRMED",
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
		var tradelabError TradelabExitSpreadOrderResponse
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelSpreadOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.CancelSpreadOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.CancelSpreadOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.ExitSpreadOrderRequest{
		ClientID:          "CLIENT1",
		ExchangeOrderID:   "ionwdg123",
		LegOrderIndicator: "Entry",
		OmsOrderID:        "123445",
		Status:            "CONFIRMED",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var orderSuccess models.SpreadOrderResponse
	orderSuccess.BasketID = "122334"
	orderSuccess.Message = "Placed"
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradelabExitSpreadOrderResponse
		tradelabSuccess.Data.BasketID = "122334"
		tradelabSuccess.Data.Message = "Placed"
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
			obj := ConditionalOrderObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.CancelSpreadOrder(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ConditionalOrderObj.CancelSpreadOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConditionalOrderObj.CancelSpreadOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end

}
