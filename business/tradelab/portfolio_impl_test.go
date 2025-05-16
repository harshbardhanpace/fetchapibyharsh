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

func TestPortfolioObj_FetchDematHoldings(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.FetchDematHoldingsRequest
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

	req1 := models.FetchDematHoldingsRequest{
		ClientID: "test123",
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
			obj := PortfolioObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.FetchDematHoldings(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("PortfolioObj.FetchDematHoldings() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortfolioObj.FetchDematHoldings() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.FetchDematHoldingsRequest{
		ClientID: "test123",
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
			obj := PortfolioObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.FetchDematHoldings(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("PortfolioObj.FetchDematHoldings() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortfolioObj.FetchDematHoldings() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.FetchDematHoldingsRequest{
		ClientID: "test124",
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
		var tradelabError TradeLabFetchDematHoldingsResponse
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
			obj := PortfolioObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.FetchDematHoldings(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("PortfolioObj.FetchDematHoldings() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortfolioObj.FetchDematHoldings() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.FetchDematHoldingsRequest{
		ClientID: "test124",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var orderSuccess models.FetchDematHoldingsResponse
	responseOrders := make([]models.FetchDematHoldingsResponseData, 0)
	for i := 0; i < 1; i++ {
		var fetchDHResponseData models.FetchDematHoldingsResponseData
		fetchDHResponseData.BranchCode = ""
		fetchDHResponseData.BuyAvg = 250.25
		fetchDHResponseData.BuyAvgMtm = 240
		fetchDHResponseData.ClientID = "Client1"
		fetchDHResponseData.Exchange = "NSE"
		fetchDHResponseData.FreeQuantity = 55
		fetchDHResponseData.InstrumentDetails.Exchange = 1
		fetchDHResponseData.InstrumentDetails.InstrumentName = "EQ"
		fetchDHResponseData.InstrumentDetails.InstrumentToken = 3045
		fetchDHResponseData.InstrumentDetails.TradingSymbol = "SBI-EQ"
		fetchDHResponseData.Isin = "INE062A01020"
		fetchDHResponseData.Ltp = 245.45
		fetchDHResponseData.PendingQuantity = 0
		fetchDHResponseData.PledgeQuantity = 0
		fetchDHResponseData.PreviousClose = 246
		fetchDHResponseData.Quantity = 10
		fetchDHResponseData.Symbol = "SBIN"
		fetchDHResponseData.T0Price = 0
		fetchDHResponseData.T0Quantity = 0
		fetchDHResponseData.T1Price = 0
		fetchDHResponseData.T1Quantity = 0
		fetchDHResponseData.T2Price = 0
		fetchDHResponseData.T2Quantity = 0
		fetchDHResponseData.TodayPledgeQuantity = 0
		fetchDHResponseData.Token = 3045
		fetchDHResponseData.TradingSymbol = "SBIN"
		fetchDHResponseData.TransactionType = ""
		fetchDHResponseData.UsedQuantity = 0
		responseOrders = append(responseOrders, fetchDHResponseData)
	}
	orderSuccess.Holdings = responseOrders
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradeLabFetchDematHoldingsResponse
		responseOrders := make([]TradelabFetchDematHoldingsHoldings, 0)
		for i := 0; i < 1; i++ {
			var fetchDHResponseData TradelabFetchDematHoldingsHoldings
			fetchDHResponseData.BranchCode = ""
			fetchDHResponseData.BuyAvg = 250.25
			fetchDHResponseData.BuyAvgMtm = 240
			fetchDHResponseData.ClientID = "Client1"
			fetchDHResponseData.Exchange = "NSE"
			fetchDHResponseData.FreeQuantity = 55
			fetchDHResponseData.InstrumentDetails.Exchange = 1
			fetchDHResponseData.InstrumentDetails.InstrumentName = "EQ"
			fetchDHResponseData.InstrumentDetails.InstrumentToken = 3045
			fetchDHResponseData.InstrumentDetails.TradingSymbol = "SBI-EQ"
			fetchDHResponseData.Isin = "INE062A01020"
			fetchDHResponseData.Ltp = 245.45
			fetchDHResponseData.PendingQuantity = 0
			fetchDHResponseData.PledgeQuantity = 0
			fetchDHResponseData.PreviousClose = 246
			fetchDHResponseData.Quantity = 10
			fetchDHResponseData.Symbol = "SBIN"
			fetchDHResponseData.T0Price = 0
			fetchDHResponseData.T0Quantity = 0
			fetchDHResponseData.T1Price = 0
			fetchDHResponseData.T1Quantity = 0
			fetchDHResponseData.T2Price = 0
			fetchDHResponseData.T2Quantity = 0
			fetchDHResponseData.TodayPledgeQuantity = 0
			fetchDHResponseData.Token = 3045
			fetchDHResponseData.TradingSymbol = "SBIN"
			fetchDHResponseData.TransactionType = ""
			fetchDHResponseData.UsedQuantity = 0
			responseOrders = append(responseOrders, fetchDHResponseData)
		}
		var Data TradelabFetchDematHoldingsHoldingsData
		Data.Holdings = responseOrders
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
			obj := PortfolioObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.FetchDematHoldings(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("PortfolioObj.FetchDematHoldings() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortfolioObj.FetchDematHoldings() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end
}

func TestPortfolioObj_ConvertPositions(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.ConvertPositionsRequest
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

	req1 := models.ConvertPositionsRequest{
		ClientID:        "test123",
		Exchange:        "NSE",
		InstrumentToken: 22,
		Product:         "CNC",
		NewProduct:      "MIS",
		Quantity:        10,
		Validity:        "DAY",
		OrderSide:       "BUY",
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
			obj := PortfolioObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ConvertPositions(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("PortfolioObj.ConvertPositions() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortfolioObj.ConvertPositions() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.ConvertPositionsRequest{
		ClientID:        "test123",
		Exchange:        "NSE",
		InstrumentToken: 22,
		Product:         "CNC",
		NewProduct:      "MIS",
		Quantity:        10,
		Validity:        "DAY",
		OrderSide:       "BUY",
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
			obj := PortfolioObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ConvertPositions(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("PortfolioObj.ConvertPositions() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortfolioObj.ConvertPositions() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.ConvertPositionsRequest{
		ClientID:        "test123",
		Exchange:        "NSE",
		InstrumentToken: 22,
		Product:         "CNC",
		NewProduct:      "MIS",
		Quantity:        10,
		Validity:        "DAY",
		OrderSide:       "BUY",
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
		var tradelabError TradeLabConvertPositionsResponse
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
			obj := PortfolioObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ConvertPositions(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("PortfolioObj.ConvertPositions() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortfolioObj.ConvertPositions() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.ConvertPositionsRequest{
		ClientID:        "test123",
		Exchange:        "NSE",
		InstrumentToken: 22,
		Product:         "CNC",
		NewProduct:      "MIS",
		Quantity:        10,
		Validity:        "DAY",
		OrderSide:       "BUY",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var orderSuccess models.ConvertPositionsResponse
	//orderSuccess.Data = {}
	res4 := apihelpers.APIRes{
		Data:    orderSuccess,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradeLabConvertPositionsResponse
		//tradelabSuccess.Data = {}
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
			obj := PortfolioObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ConvertPositions(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("PortfolioObj.ConvertPositions() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortfolioObj.ConvertPositions() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end
}

func TestPortfolioObj_GetPositions(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.GetPositionRequest
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

	req1 := models.GetPositionRequest{
		ClientID: "test123",
		Type:     "live",
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
			obj := PortfolioObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.GetPositions(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("PortfolioObj.GetPositions() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortfolioObj.GetPositions() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.GetPositionRequest{
		ClientID: "test123",
		Type:     "live",
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
			obj := PortfolioObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.GetPositions(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("PortfolioObj.GetPositions() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortfolioObj.GetPositions() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.GetPositionRequest{
		ClientID: "test124",
		Type:     "live",
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
		var tradelabError TradeLabGetPositionResponse
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
			obj := PortfolioObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.GetPositions(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("PortfolioObj.GetPositions() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortfolioObj.GetPositions() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.GetPositionRequest{
		ClientID: "test124",
		Type:     "live",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	responseOrders := make([]models.GetPositionResponseData, 0)
	for i := 0; i < 1; i++ {
		var getPositionResponseData models.GetPositionResponseData
		getPositionResponseData.AverageBuyPrice = 14.7
		getPositionResponseData.AveragePrice = 14
		getPositionResponseData.AverageSellPrice = 0
		getPositionResponseData.BuyAmount = 15
		getPositionResponseData.BuyQuantity = 100
		getPositionResponseData.CfBuyAmount = 0
		getPositionResponseData.CfBuyQuantity = 0
		getPositionResponseData.CfSellAmount = 0
		getPositionResponseData.CfSellQuantity = 0
		getPositionResponseData.ClientID = "Client1"
		getPositionResponseData.ClosePrice = 0
		getPositionResponseData.Exchange = "NSE"
		getPositionResponseData.InstrumentToken = 11915
		getPositionResponseData.Ltp = 14.6
		getPositionResponseData.Multiplier = 1
		getPositionResponseData.NetAmount = 100
		getPositionResponseData.NetQuantity = 110
		getPositionResponseData.PreviousClose = 14.5
		getPositionResponseData.ProCli = "CLIENT"
		getPositionResponseData.ProdType = 2
		getPositionResponseData.Product = "CNC"
		getPositionResponseData.RealizedMtm = 0
		getPositionResponseData.Segment = "Capital"
		getPositionResponseData.SellAmount = 0
		getPositionResponseData.SellQuantity = 0
		getPositionResponseData.Symbol = "YESBANK"
		getPositionResponseData.Token = 11915
		//getPositionResponseData.TradingSymbol = "YESBANK-EQ"
		getPositionResponseData.VLoginID = "Client1"
		responseOrders = append(responseOrders, getPositionResponseData)
	}
	res4 := apihelpers.APIRes{
		Data:    responseOrders,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradeLabGetPositionResponse
		responseOrders := make([]TradeLabGetPositionResponseData, 0)
		for i := 0; i < 1; i++ {
			var getPositionResponseData TradeLabGetPositionResponseData
			getPositionResponseData.AverageBuyPrice = 14.7
			getPositionResponseData.AveragePrice = 14
			getPositionResponseData.AverageSellPrice = 0
			getPositionResponseData.BuyAmount = 15
			getPositionResponseData.BuyQuantity = 100
			getPositionResponseData.CfBuyAmount = 0
			getPositionResponseData.CfBuyQuantity = 0
			getPositionResponseData.CfSellAmount = 0
			getPositionResponseData.CfSellQuantity = 0
			getPositionResponseData.ClientID = "Client1"
			getPositionResponseData.ClosePrice = 0
			getPositionResponseData.Exchange = "NSE"
			getPositionResponseData.InstrumentToken = 11915
			getPositionResponseData.Ltp = 14.6
			getPositionResponseData.Multiplier = 1
			getPositionResponseData.NetAmount = 100
			getPositionResponseData.NetQuantity = 110
			getPositionResponseData.PreviousClose = 14.5
			getPositionResponseData.ProCli = "CLIENT"
			getPositionResponseData.ProdType = 2
			getPositionResponseData.Product = "CNC"
			getPositionResponseData.RealizedMtm = 0
			getPositionResponseData.Segment = "Capital"
			getPositionResponseData.SellAmount = 0
			getPositionResponseData.SellQuantity = 0
			getPositionResponseData.Symbol = "YESBANK"
			getPositionResponseData.Token = 11915
			//getPositionResponseData.TradingSymbol = "YESBANK-EQ"
			getPositionResponseData.VLoginID = "Client1"
			responseOrders = append(responseOrders, getPositionResponseData)
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
			obj := PortfolioObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.GetPositions(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("PortfolioObj.GetPositions() got = %v, want = %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortfolioObj.GetPositions() got1 = %v, want = %v", got1, tt.want1)
			}
		})
	}
	//test 4 end
}
