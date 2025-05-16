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

func TestProfileObj_GetProfile(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.ProfileRequest
		reqH models.ReqHeader
	}

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

	req1 := models.ProfileRequest{
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
			obj := ProfileObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.GetProfile(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ProfileObj.GetProfile() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ProfileObj.GetProfile() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.ProfileRequest{
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
			obj := ProfileObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.GetProfile(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ProfileObj.GetProfile() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ProfileObj.GetProfile() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.ProfileRequest{
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
		var tradelabError TradeLabProfileResponse
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
			obj := ProfileObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.GetProfile(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ProfileObj.GetProfile() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ProfileObj.GetProfile() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.ProfileRequest{
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

	var orderSuccess models.ProfileResponse
	orderSuccess.Data.AccountType = ""
	orderSuccess.Data.BankAccountNumber = "1234567890"
	orderSuccess.Data.BankBranchName = "B1"
	orderSuccess.Data.BankName = "ABC Bank of India"
	orderSuccess.Data.Branch = "Bangalore"
	orderSuccess.Data.BrokerID = "ABC"
	orderSuccess.Data.City = "Bangalore"
	orderSuccess.Data.ClientID = "Client1"
	orderSuccess.Data.Dob = "01/01/2000"
	orderSuccess.Data.EmailID = "client1@xyz.com"
	//orderSuccess.Data.ExchangeNnf = {}
	//orderSuccess.Data.ExchangesSubscribed = []
	orderSuccess.Data.IfscCode = "ABCN001"
	orderSuccess.Data.Name = "Demo Client"
	orderSuccess.Data.OfficeAddr = "abc@gmail.com 9876543210"
	orderSuccess.Data.PanNumber = "ABCDEFGHI"
	orderSuccess.Data.PermanentAddr = "STATE"
	orderSuccess.Data.PhoneNumber = "9876543210"
	//orderSuccess.Data.ProductsEnabled = []
	orderSuccess.Data.Role.ID = 69
	orderSuccess.Data.Role.Name = "CLIENT"
	orderSuccess.Data.Sex = "M"
	orderSuccess.Data.State = "Karnataka"
	orderSuccess.Data.Status = "Activated"
	orderSuccess.Data.TwofaEnabled = true
	orderSuccess.Data.UserType = "Non-Institutional"
	res4 := apihelpers.APIRes{
		Data:    orderSuccess.Data,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradeLabProfileResponse
		tradelabSuccess.Data.AccountType = ""
		tradelabSuccess.Data.BankAccountNumber = "1234567890"
		tradelabSuccess.Data.BankBranchName = "B1"
		tradelabSuccess.Data.BankName = "ABC Bank of India"
		tradelabSuccess.Data.Branch = "Bangalore"
		tradelabSuccess.Data.BrokerID = "ABC"
		tradelabSuccess.Data.City = "Bangalore"
		tradelabSuccess.Data.ClientID = "Client1"
		tradelabSuccess.Data.Dob = "01/01/2000"
		tradelabSuccess.Data.EmailID = "client1@xyz.com"
		//tradelabSuccess.Data.ExchangeNnf = {}
		//tradelabSuccess.Data.ExchangesSubscribed = []
		tradelabSuccess.Data.IfscCode = "ABCN001"
		tradelabSuccess.Data.Name = "Demo Client"
		tradelabSuccess.Data.OfficeAddr = "abc@gmail.com 9876543210"
		tradelabSuccess.Data.PanNumber = "ABCDEFGHI"
		tradelabSuccess.Data.PermanentAddr = "STATE"
		tradelabSuccess.Data.PhoneNumber = "9876543210"
		//tradelabSuccess.Data.ProductsEnabled = []
		tradelabSuccess.Data.Role.ID = 69
		tradelabSuccess.Data.Role.Name = "CLIENT"
		tradelabSuccess.Data.Sex = "M"
		tradelabSuccess.Data.State = "Karnataka"
		tradelabSuccess.Data.Status = "Activated"
		tradelabSuccess.Data.TwofaEnabled = true
		tradelabSuccess.Data.UserType = "Non-Institutional"
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
			obj := ProfileObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.GetProfile(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ProfileObj.GetProfile() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ProfileObj.GetProfile() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end
}
