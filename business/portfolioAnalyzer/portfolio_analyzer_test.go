package portfolioanalyzer

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
	"space/business/pockets"
	"space/constants"
	"space/db"
	"space/loggerconfig"
	"space/models"
)

func TestPortfolioAnalyzerObj_HoldingsWeightages(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.PortfolioAnalyzerReq
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

	field1 := fields{
		tradeLabURL: "http://test",
	}

	req1 := models.PortfolioAnalyzerReq{
		ClientId: "9CD12",
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
			obj := PAObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.HoldingsWeightages(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("PAObj.HoldingsWeightages() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PAObj.HoldingsWeightages() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	//test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.PortfolioAnalyzerReq{
		ClientId: "9CD12",
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2f",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	holdings := make([]models.Holdings, 0)
	for i := 0; i < 1; i++ {
		var holdingsData models.Holdings
		holdingsData.Token = "10666"
		holdingsData.Exchange = "NSE"
		holdingsData.Isin = "IN699879"
		holdingsData.TradingSymbol = "XYZ"
		holdingsData.SectorName = "IN699879"
		holdingsData.SectorCode = "XYZ"
		holdingsData.PercentageOfPortfolio = 100
		holdingsData.ValueOfHolding = 2
		holdings = append(holdings, holdingsData)
	}

	res2 := apihelpers.APIRes{
		Data:    holdings,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock holdings
	FetchHoldingsData = func(url string, clientID string, reqH models.ReqHeader) ([]models.HoldingsData, error) {
		var response []models.HoldingsData

		for i := 0; i < 1; i++ {
			var holding models.HoldingsData
			holding.Exchange = "NSE"
			holding.Symbol = "XYZ"
			holding.Isin = "IN699879"
			holding.LTP = 2
			holding.Quantity = 1
			holding.Token = "10666"
			response = append(response, holding)
		}

		return response, nil
	}

	db.CallFetchSector = func(isin string) (string, string, error) {
		return "XYZ", "IN699879", nil
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceId, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		holdings := make([]models.Holdings, 0)
		for i := 0; i < 1; i++ {
			var holdingsData models.Holdings
			holdingsData.Token = "10666"
			holdingsData.Exchange = "NSE"
			holdingsData.Isin = "IN699879"
			holdingsData.TradingSymbol = "XYZ"
			holdingsData.SectorName = "IN699879"
			holdingsData.SectorCode = "XYZ"
			holdingsData.PercentageOfPortfolio = 100
			holdingsData.ValueOfHolding = 2
			holdings = append(holdings, holdingsData)
		}

		var res http.Response
		jsonRes, _ := json.Marshal(holdings)
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
		{"Success", field2, arg2, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := PAObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.HoldingsWeightages(tt.args.req, tt.args.reqH)
			//got1 = res2
			if got != tt.want {
				t.Errorf("PAObj.HoldingsWeightages() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PAObj.HoldingsWeightages() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

}

func TestPortfolioAnalyzerObj_PortfolioBeta(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.PortfolioAnalyzerReq
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

	field1 := fields{
		tradeLabURL: "http://test",
	}

	req1 := models.PortfolioAnalyzerReq{
		ClientId: "9CD12",
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
			obj := PAObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PortfolioBeta(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("PAObj.PortfolioBeta() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PAObj.PortfolioBeta() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	//test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.PortfolioAnalyzerReq{
		ClientId: "9CD12",
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2f",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	//mock holdings
	FetchHoldingsData = func(url string, clientID string, reqH models.ReqHeader) ([]models.HoldingsData, error) {
		var response []models.HoldingsData

		for i := 0; i < 1; i++ {
			var holding models.HoldingsData
			holding.Exchange = "NSE"
			holding.Isin = "IN699879"
			holding.LTP = 2
			holding.Quantity = 1
			holding.Token = "10666"
			response = append(response, holding)
		}

		return response, nil
	}

	individualBeta := make([]models.IndividualBeta, 0)
	for i := 0; i < 1; i++ {
		var betaIndi models.IndividualBeta
		betaIndi.Token = "10666"
		betaIndi.Exchange = "NSE"
		betaIndi.Isin = "IN699879"
		betaIndi.Beta = 0.89
		individualBeta = append(individualBeta, betaIndi)
	}
	var portfolioBeta models.PortfolioBeta
	portfolioBeta.PortfolioBeta = 0.89
	portfolioBeta.IndividualBeta = individualBeta

	res2 := apihelpers.APIRes{
		Data:    portfolioBeta,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock beta
	CalculateBeta = func(benchmark []models.TLCandleData, pricesIth []models.TLCandleData) float64 {
		return 0.89
	}

	//mock tl candles
	pockets.CallTLforData = func(url string, reqH models.ReqHeader) ([]models.TLCandleData, bool) {
		var candles []models.TLCandleData
		return candles, true
	}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", field2, arg2, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := PAObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PortfolioBeta(tt.args.req, tt.args.reqH)
			//got1 = res2
			if got != tt.want {
				t.Errorf("PAObj.PortfolioBeta() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PAObj.PortfolioBeta() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

}

func TestPortfolioAnalyzerObj_PortfolioPE(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.PortfolioAnalyzerReq
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

	field1 := fields{
		tradeLabURL: "http://test",
	}

	req1 := models.PortfolioAnalyzerReq{
		ClientId: "9CD12",
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
			obj := PAObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PortfolioPE(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("PAObj.PortfolioPE() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PAObj.PortfolioPE() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	//test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.PortfolioAnalyzerReq{
		ClientId: "9CD12",
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2f",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	individualPE := make([]models.IndividualPE, 0)
	for i := 0; i < 1; i++ {
		var indiPE models.IndividualPE
		indiPE.Token = "10666"
		indiPE.Isin = "IN699879"
		indiPE.TradingSymbol = "XYZ"
		indiPE.Pe = 18.18
		individualPE = append(individualPE, indiPE)
	}
	var portfolioPE models.PortfolioPE
	portfolioPE.IndividualPE = individualPE
	portfolioPE.PortfolioPE = 18.18

	res2 := apihelpers.APIRes{
		Data:    portfolioPE,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock holdings
	FetchHoldingsData = func(url string, clientID string, reqH models.ReqHeader) ([]models.HoldingsData, error) {
		var response []models.HoldingsData

		for i := 0; i < 1; i++ {
			var holding models.HoldingsData
			holding.Exchange = "NSE"
			holding.Symbol = "XYZ"
			holding.Isin = "IN699879"
			holding.LTP = 2
			holding.Quantity = 1
			holding.Token = "10666"
			response = append(response, holding)
		}

		return response, nil
	}

	//mock db
	db.CallFetchPE = func(isin string) (float64, error) {
		return 18.18, nil
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceId, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		individualPE := make([]models.IndividualPE, 0)
		for i := 0; i < 1; i++ {
			var indiPE models.IndividualPE
			indiPE.Token = "10666"
			indiPE.Isin = "IN699879"
			indiPE.TradingSymbol = "XYZ"
			indiPE.Pe = 18.18
			individualPE = append(individualPE, indiPE)
		}
		var portfolioPE models.PortfolioPE
		portfolioPE.IndividualPE = individualPE
		portfolioPE.PortfolioPE = 18.18

		var res http.Response
		jsonRes, _ := json.Marshal(portfolioPE)
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
		{"Success", field2, arg2, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := PAObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PortfolioPE(tt.args.req, tt.args.reqH)
			//got1 = res2
			if got != tt.want {
				t.Errorf("PAObj.PortfolioPE() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PAObj.PortfolioPE() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

}

func TestPortfolioAnalyzerObj_PortfolioDE(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.PortfolioAnalyzerReq
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

	field1 := fields{
		tradeLabURL: "http://test",
	}

	req1 := models.PortfolioAnalyzerReq{
		ClientId: "9CD12",
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
			obj := PAObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PortfolioDE(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("PAObj.PortfolioDE() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PAObj.PortfolioDE() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	//test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.PortfolioAnalyzerReq{
		ClientId: "9CD12",
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2f",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	individualDE := make([]models.IndividualDE, 0)
	for i := 0; i < 1; i++ {
		var indiDE models.IndividualDE
		indiDE.Token = "10666"
		indiDE.Isin = "IN699879"
		indiDE.TradingSymbol = "XYZ"
		indiDE.De = 1.18
		individualDE = append(individualDE, indiDE)
	}
	var portfolioDE models.PortfolioDE
	portfolioDE.IndividualDE = individualDE
	portfolioDE.PortfolioDE = 1.18

	res2 := apihelpers.APIRes{
		Data:    portfolioDE,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock holdings
	FetchHoldingsData = func(url string, clientID string, reqH models.ReqHeader) ([]models.HoldingsData, error) {
		var response []models.HoldingsData

		for i := 0; i < 1; i++ {
			var holding models.HoldingsData
			holding.Exchange = "NSE"
			holding.Symbol = "XYZ"
			holding.Isin = "IN699879"
			holding.LTP = 2
			holding.Quantity = 1
			holding.Token = "10666"
			response = append(response, holding)
		}

		return response, nil
	}

	//mock db
	db.CallFetchDE = func(isin string) (float64, error) {
		return 1.18, nil
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceId, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		individualDE := make([]models.IndividualDE, 0)
		for i := 0; i < 1; i++ {
			var indiDE models.IndividualDE
			indiDE.Token = "10666"
			indiDE.Isin = "IN699879"
			indiDE.TradingSymbol = "XYZ"
			indiDE.De = 1.18
			individualDE = append(individualDE, indiDE)
		}
		var portfolioDE models.PortfolioDE
		portfolioDE.IndividualDE = individualDE
		portfolioDE.PortfolioDE = 1.18

		var res http.Response
		jsonRes, _ := json.Marshal(portfolioDE)
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
		{"Success", field2, arg2, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := PAObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.PortfolioDE(tt.args.req, tt.args.reqH)
			//got1 = res2
			if got != tt.want {
				t.Errorf("PAObj.PortfolioDE() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PAObj.PortfolioDE() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

}
