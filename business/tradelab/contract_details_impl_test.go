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
	"space/loggerconfig"
	"space/models"
)

func TestContractDetailsObj_SearchScrip(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.SearchScripRequest
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

	req1 := models.SearchScripRequest{
		Key: "fewsdh",
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
			obj := ContractDetailsObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.SearchScrip(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ContractDetailsObj.SearchScrip() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ContractDetailsObj.SearchScrip() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.SearchScripRequest{
		Key: "fewsdh",
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
			obj := ContractDetailsObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.SearchScrip(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ContractDetailsObj.SearchScrip() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ContractDetailsObj.SearchScrip() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.SearchScripRequest{
		Key: "fewsdh",
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
		var tradelabError TradeLabSearchScripResponse
		// tradelabError.Error.Code = 0
		tradelabError.Error.Message = "error"
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
			obj := ContractDetailsObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.SearchScrip(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ContractDetailsObj.SearchScrip() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ContractDetailsObj.SearchScrip() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.SearchScripRequest{
		Key: "fewsdh",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	searchScripResponseData := make([]models.SearchScripResponseResult, 0)
	for i := 0; i < 1; i++ {
		var searchScripData models.SearchScripResponseResult
		searchScripData.Token = "fdsfsd"
		searchScripData.Exchange = "NSE"
		searchScripData.Company = "fdsfd"
		searchScripData.Symbol = "fsdffds"
		searchScripData.TradingSymbol = "fdsfs"
		searchScripData.DisplayName = "geasf"
		searchScripData.Score = 11.1
		searchScripData.IsTradable = true
		searchScripData.ClosePrice = "gnjdsfd"
		searchScripData.Alternate.Token = "esae"
		searchScripData.Alternate.Exchange = "DSE"
		searchScripData.Alternate.Company = "jsdf"
		searchScripData.Alternate.Symbol = "kjse"
		searchScripData.Alternate.TradingSymbol = "sefd"
		searchScripData.Alternate.DisplayName = "seds"
		searchScripData.Alternate.IsTradable = true
		searchScripData.Alternate.ClosePrice = "sdhf"
		searchScripResponseData = append(searchScripResponseData, searchScripData)
	}

	res4 := apihelpers.APIRes{
		Data:    searchScripResponseData,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradeLabSearchScripResponse
		searchScripResponseData := make([]TradeLabSearchScripResponseResult, 0)
		for i := 0; i < 1; i++ {
			var searchScripData TradeLabSearchScripResponseResult
			searchScripData.Token = "fdsfsd"
			searchScripData.Exchange = "NSE"
			searchScripData.Company = "fdsfd"
			searchScripData.Symbol = "fsdffds"
			searchScripData.TradingSymbol = "fdsfs"
			searchScripData.DisplayName = "geasf"
			searchScripData.Score = 11.1
			searchScripData.IsTradable = true
			searchScripData.ClosePrice = "gnjdsfd"
			searchScripData.Alternate.Token = "esae"
			searchScripData.Alternate.Exchange = "DSE"
			searchScripData.Alternate.Company = "jsdf"
			searchScripData.Alternate.Symbol = "kjse"
			searchScripData.Alternate.TradingSymbol = "sefd"
			searchScripData.Alternate.DisplayName = "seds"
			searchScripData.Alternate.IsTradable = true
			searchScripData.Alternate.ClosePrice = "sdhf"
			searchScripResponseData = append(searchScripResponseData, searchScripData)
		}

		tradelabSuccess.Result = searchScripResponseData
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
			obj := ContractDetailsObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.SearchScrip(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ContractDetailsObj.SearchScrip() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ContractDetailsObj.SearchScrip() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end

}

func TestContractDetailsObj_ScripInfo(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.ScripInfoRequest
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

	req1 := models.ScripInfoRequest{
		Exchange: "NSE",
		Info:     "fdsf",
		Token:    "hfgyd",
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
			obj := ContractDetailsObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ScripInfo(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ContractDetailsObj.ScripInfo() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ContractDetailsObj.ScripInfo() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.ScripInfoRequest{
		Exchange: "NSE",
		Info:     "fdsf",
		Token:    "hfgyd",
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
			obj := ContractDetailsObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ScripInfo(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ContractDetailsObj.ScripInfo() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ContractDetailsObj.ScripInfo() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.ScripInfoRequest{
		Exchange: "NSE",
		Info:     "fdsf",
		Token:    "hfgyd",
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
		var tradelabError TradeLabScripInfoResponse
		// tradelabError.Error.Code = 0
		tradelabError.Error.Message = "error"
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
			obj := ContractDetailsObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ScripInfo(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ContractDetailsObj.ScripInfo() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ContractDetailsObj.ScripInfo() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.ScripInfoRequest{
		Exchange: "NSE",
		Info:     "fdsf",
		Token:    "hfgyd",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var scripInfoResponse models.ScripInfoResponse
	scripInfoResponse.Result.BoardLotQuantity = 654
	scripInfoResponse.Result.ChangeInOi = 654
	scripInfoResponse.Result.Exchange = 654
	scripInfoResponse.Result.Expiry = 654
	scripInfoResponse.Result.HigherCircuitLimit = 654.11
	scripInfoResponse.Result.InstrumentName = "fdjif"
	scripInfoResponse.Result.InstrumentToken = 654
	scripInfoResponse.Result.Isin = "gjkf"
	scripInfoResponse.Result.LowerCircuitLimit = 654
	scripInfoResponse.Result.Multiplier = 234
	scripInfoResponse.Result.OpenInterest = 432
	scripInfoResponse.Result.OptionType = "dg"
	scripInfoResponse.Result.Precision = 43
	scripInfoResponse.Result.Series = "fed"
	scripInfoResponse.Result.Strike = 432
	scripInfoResponse.Result.Symbol = "fndsj"
	scripInfoResponse.Result.TickSize = 43
	scripInfoResponse.Result.TradingSymbol = "fjkds"
	scripInfoResponse.Result.UnderlyingToken = 432
	scripInfoResponse.Result.RawExpiry = 43
	scripInfoResponse.Result.Freeze = 43
	scripInfoResponse.Result.InstrumentType = "fds"
	scripInfoResponse.Result.IssueRate = 434
	scripInfoResponse.Result.IssueStartDate = "fgd3"
	scripInfoResponse.Result.ListDate = "fd"
	scripInfoResponse.Result.MaxOrderSize = 32
	scripInfoResponse.Result.PriceNumerator = 543
	scripInfoResponse.Result.PriceDenominator = 32
	scripInfoResponse.Result.Comments = "fes"
	scripInfoResponse.Result.CircuitRating = "fes"
	scripInfoResponse.Result.CompanyName = "fes"
	scripInfoResponse.Result.DisplayName = "fes"
	scripInfoResponse.Result.RawTickSize = 432
	scripInfoResponse.Result.IsIndex = true
	scripInfoResponse.Result.Tradable = false
	scripInfoResponse.Result.MaxSingleQty = 432
	scripInfoResponse.Result.ExpiryString = "fes"
	scripInfoResponse.Result.LocalUpdateTime = "fes"
	scripInfoResponse.Result.MarketType = "fes"
	scripInfoResponse.Result.PriceUnits = "fes"
	scripInfoResponse.Result.TradingUnits = "fes"
	scripInfoResponse.Result.LastTradingDate = "fes"
	scripInfoResponse.Result.TenderPeriodEndDate = "fes"
	scripInfoResponse.Result.DeliveryStartDate = "fes"
	scripInfoResponse.Result.PriceQuotation = 432
	scripInfoResponse.Result.GeneralDenominator = "fes"
	scripInfoResponse.Result.TenderPeriodStartDate = "fes"
	scripInfoResponse.Result.DeliveryUnits = "fes"
	scripInfoResponse.Result.DeliveryEndDate = "fes"
	scripInfoResponse.Result.TradingUnitFactor = 423
	scripInfoResponse.Result.DeliveryUnitFactor = 423
	scripInfoResponse.Result.BookClosureEndDate = "fes"
	scripInfoResponse.Result.BookClosureStartDate = "fes"
	scripInfoResponse.Result.NoDeliveryDateEnd = "fes"
	scripInfoResponse.Result.NoDeliveryDateStart = "fes"
	scripInfoResponse.Result.ReAdmissionDate = "fes"
	scripInfoResponse.Result.RecordDate = "fes"
	scripInfoResponse.Result.Warning = "fes"
	scripInfoResponse.Result.Dpr = "fes"
	scripInfoResponse.Result.TradeToTrade = false
	scripInfoResponse.Result.SurveillanceIndicator = 432
	scripInfoResponse.Result.PartitionID = 42
	scripInfoResponse.Result.ProductID = 432
	scripInfoResponse.Result.ProductCategory = "fsd"
	scripInfoResponse.Result.MonthIdentifier = 234
	scripInfoResponse.Result.ClosePrice = "234"
	scripInfoResponse.Result.SpecialPreopen = 234
	scripInfoResponse.Result.AlternateExchange = "NSE"
	scripInfoResponse.Result.AlternateToken = 432
	scripInfoResponse.Result.Asm = "fdsse"
	scripInfoResponse.Result.Gsm = "fdssz"
	scripInfoResponse.Result.Execution = "fers"
	scripInfoResponse.Result.Symbol2 = "fds"
	scripInfoResponse.Result.RawTenderPeriodStartDate = "fds"
	scripInfoResponse.Result.RawTenderPeriodEndDate = "43"
	scripInfoResponse.Result.YearlyHighPrice = "332"
	scripInfoResponse.Result.YearlyLowPrice = "543"
	scripInfoResponse.Result.IssueMaturityDate = 43
	scripInfoResponse.Result.Var = "fdfw"
	scripInfoResponse.Result.Exposure = "fd"
	scripInfoResponse.Result.Span = []int{}
	scripInfoResponse.Result.HaveFutures = true
	scripInfoResponse.Result.HaveOptions = true
	scripInfoResponse.Result.Tag = "fds"

	res4 := apihelpers.APIRes{
		Data:    scripInfoResponse.Result,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var scripInfoResponse TradeLabScripInfoResponse
		scripInfoResponse.Result.BoardLotQuantity = 654
		scripInfoResponse.Result.ChangeInOi = 654
		scripInfoResponse.Result.Exchange = 654
		scripInfoResponse.Result.Expiry = 654
		scripInfoResponse.Result.HigherCircuitLimit = 654.11
		scripInfoResponse.Result.InstrumentName = "fdjif"
		scripInfoResponse.Result.InstrumentToken = 654
		scripInfoResponse.Result.Isin = "gjkf"
		scripInfoResponse.Result.LowerCircuitLimit = 654
		scripInfoResponse.Result.Multiplier = 234
		scripInfoResponse.Result.OpenInterest = 432
		scripInfoResponse.Result.OptionType = "dg"
		scripInfoResponse.Result.Precision = 43
		scripInfoResponse.Result.Series = "fed"
		scripInfoResponse.Result.Strike = 432
		scripInfoResponse.Result.Symbol = "fndsj"
		scripInfoResponse.Result.TickSize = 43
		scripInfoResponse.Result.TradingSymbol = "fjkds"
		scripInfoResponse.Result.UnderlyingToken = 432
		scripInfoResponse.Result.RawExpiry = 43
		scripInfoResponse.Result.Freeze = 43
		scripInfoResponse.Result.InstrumentType = "fds"
		scripInfoResponse.Result.IssueRate = 434
		scripInfoResponse.Result.IssueStartDate = "fgd3"
		scripInfoResponse.Result.ListDate = "fd"
		scripInfoResponse.Result.MaxOrderSize = 32
		scripInfoResponse.Result.PriceNumerator = 543
		scripInfoResponse.Result.PriceDenominator = 32
		scripInfoResponse.Result.Comments = "fes"
		scripInfoResponse.Result.CircuitRating = "fes"
		scripInfoResponse.Result.CompanyName = "fes"
		scripInfoResponse.Result.DisplayName = "fes"
		scripInfoResponse.Result.RawTickSize = 432
		scripInfoResponse.Result.IsIndex = true
		scripInfoResponse.Result.Tradable = false
		scripInfoResponse.Result.MaxSingleQty = 432
		scripInfoResponse.Result.ExpiryString = "fes"
		scripInfoResponse.Result.LocalUpdateTime = "fes"
		scripInfoResponse.Result.MarketType = "fes"
		scripInfoResponse.Result.PriceUnits = "fes"
		scripInfoResponse.Result.TradingUnits = "fes"
		scripInfoResponse.Result.LastTradingDate = "fes"
		scripInfoResponse.Result.TenderPeriodEndDate = "fes"
		scripInfoResponse.Result.DeliveryStartDate = "fes"
		scripInfoResponse.Result.PriceQuotation = 432
		scripInfoResponse.Result.GeneralDenominator = "fes"
		scripInfoResponse.Result.TenderPeriodStartDate = "fes"
		scripInfoResponse.Result.DeliveryUnits = "fes"
		scripInfoResponse.Result.DeliveryEndDate = "fes"
		scripInfoResponse.Result.TradingUnitFactor = 423
		scripInfoResponse.Result.DeliveryUnitFactor = 423
		scripInfoResponse.Result.BookClosureEndDate = "fes"
		scripInfoResponse.Result.BookClosureStartDate = "fes"
		scripInfoResponse.Result.NoDeliveryDateEnd = "fes"
		scripInfoResponse.Result.NoDeliveryDateStart = "fes"
		scripInfoResponse.Result.ReAdmissionDate = "fes"
		scripInfoResponse.Result.RecordDate = "fes"
		scripInfoResponse.Result.Warning = "fes"
		scripInfoResponse.Result.Dpr = "fes"
		scripInfoResponse.Result.TradeToTrade = false
		scripInfoResponse.Result.SurveillanceIndicator = 432
		scripInfoResponse.Result.PartitionID = 42
		scripInfoResponse.Result.ProductID = 432
		scripInfoResponse.Result.ProductCategory = "fsd"
		scripInfoResponse.Result.MonthIdentifier = 234
		scripInfoResponse.Result.ClosePrice = "234"
		scripInfoResponse.Result.SpecialPreopen = 234
		scripInfoResponse.Result.AlternateExchange = "NSE"
		scripInfoResponse.Result.AlternateToken = 432
		scripInfoResponse.Result.Asm = "fdsse"
		scripInfoResponse.Result.Gsm = "fdssz"
		scripInfoResponse.Result.Execution = "fers"
		scripInfoResponse.Result.Symbol2 = "fds"
		scripInfoResponse.Result.RawTenderPeriodStartDate = "fds"
		scripInfoResponse.Result.RawTenderPeriodEndDate = "43"
		scripInfoResponse.Result.YearlyHighPrice = "332"
		scripInfoResponse.Result.YearlyLowPrice = "543"
		scripInfoResponse.Result.IssueMaturityDate = 43
		scripInfoResponse.Result.Var = "fdfw"
		scripInfoResponse.Result.Exposure = "fd"
		scripInfoResponse.Result.Span = []int{}
		scripInfoResponse.Result.HaveFutures = true
		scripInfoResponse.Result.HaveOptions = true
		scripInfoResponse.Result.Tag = "fds"

		var res http.Response
		jsonRes, _ := json.Marshal(scripInfoResponse)
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
			obj := ContractDetailsObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ScripInfo(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ContractDetailsObj.ScripInfo() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ContractDetailsObj.ScripInfo() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end

}

func TestExtractPrice(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"SBIN24DEC1070CE", "1070", false},           // 4-digit price
		{"SBIN24DEC1PE", "1", false},                 // 1-digit price
		{"SBIN24DEC999999999PE", "999999999", false}, // 9-digit price
		{"SBIN24JAN9999CE", "9999", false},           // 4-digit price
		{"SBIN24NOV1234PE", "1234", false},           // 4-digit price
		{"SBIN24ABC1234PE", "1234", false},           // 3-letter substring (ABC)
		{"SBIN24XYZ5678CE", "5678", false},           // 3-letter substring (XYZ)
		{"SBIN24A12BC999PE", "", true},               // Non-month 3-letter substring
		{"SBIN24ZZZ000PE", "000", false},             // Another 3-letter substring (ZZZ)
		{"SBIN24XYZ5678XYZ", "", true},               // No "PE" or "CE" at the end
		{"", "", true},                               // Empty string
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			price, err := extractPrice(tt.input)

			if tt.hasError {
				// Expect an error
				if err == nil {
					t.Errorf("expected error for input %q, but got none", tt.input)
				}
			} else {
				// Expect no error, and the correct price
				if err != nil {
					t.Errorf("unexpected error for input %q: %v", tt.input, err)
				}
				if price != tt.expected {
					t.Errorf("for input %q, expected price %q but got %q", tt.input, tt.expected, price)
				}
			}
		})
	}
}
