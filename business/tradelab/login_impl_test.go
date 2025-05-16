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

func TestLoginObj_LoginByPass(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.LoginRequest
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

	req1 := models.LoginRequest{
		UserName: "CLIENT1",
		Password: "jfksd",
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
			obj := LoginObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.LoginByPass(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("LoginObj.LoginByPass() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("LoginObj.LoginByPass() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.LoginRequest{
		UserName: "CLIENT1",
		Password: "jfksd",
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
			obj := LoginObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.LoginByPass(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("LoginObj.LoginByPass() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("LoginObj.LoginByPass() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.LoginRequest{
		UserName: "CLIENT1",
		Password: "jfksd",
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
		var tradelabError TradeLabLoginRes
		// tradelabError.Error.Code = 0
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
			obj := LoginObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.LoginByPass(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("LoginObj.LoginByPass() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("LoginObj.LoginByPass() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.LoginRequest{
		UserName: "CLIENT1",
		Password: "jfksd",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var loginRes models.LoginResponse
	loginRes.Alert = "fdsk"
	loginRes.AuthToken = "fds"
	loginRes.LoginID = "CLIENT1"
	loginRes.ResetPassword = true
	loginRes.ResetTwoFa = true

	var twoFaDetails models.TwoFaDetails

	var twoFaQuestions []models.TwoFaQuestions

	for i := 0; i < 1; i++ {
		var twoFaQuestion models.TwoFaQuestions
		twoFaQuestion.Question = "fdssd"
		twoFaQuestion.QuestionID = 2
		twoFaQuestions = append(twoFaQuestions, twoFaQuestion)
	}

	twoFaDetails.Questions = twoFaQuestions
	twoFaDetails.TwofaToken = "fsd"
	twoFaDetails.Type = "fdsf"

	loginRes.Twofa = twoFaDetails
	loginRes.TwofaEnabled = false

	res4 := apihelpers.APIRes{
		Data:    loginRes,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradeLabLoginRes
		tradelabSuccess.Data.Alert = "fdsk"
		tradelabSuccess.Data.AuthToken = "fds"
		tradelabSuccess.Data.LoginID = "dfjfs"
		tradelabSuccess.Data.ResetPassword = true
		tradelabSuccess.Data.ResetTwoFa = true

		var tradeLabTwoFaDetails TradeLabTwoFaDetails

		var tradeLabTwoFaQuestions []TradeLabTwoFaQuestions

		for i := 0; i < 1; i++ {
			var tradeLabTwoFaQuestion TradeLabTwoFaQuestions
			tradeLabTwoFaQuestion.Question = "fdssd"
			tradeLabTwoFaQuestion.QuestionID = 2
			tradeLabTwoFaQuestions = append(tradeLabTwoFaQuestions, tradeLabTwoFaQuestion)
		}

		tradeLabTwoFaDetails.Questions = tradeLabTwoFaQuestions
		tradeLabTwoFaDetails.TwofaToken = "fsd"
		tradeLabTwoFaDetails.Type = "fdsf"

		tradelabSuccess.Data.Twofa = tradeLabTwoFaDetails
		tradelabSuccess.Data.TwofaEnabled = false
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
			obj := LoginObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.LoginByPass(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("LoginObj.LoginByPass() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("LoginObj.LoginByPass() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end

}

func TestLoginObj_ValidateTwoFa(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.ValidateTwoFARequest
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

	var twofa []models.TwoQuestions

	for i := 0; i < 1; i++ {
		var questions models.TwoQuestions
		questions.Answer = "fdssdf"
		questions.QuestionID = "34"
		twofa = append(twofa, questions)
	}

	req1 := models.ValidateTwoFARequest{
		LoginID:    "CLIENT1",
		TwofaToken: "jfksd",
		Twofa:      twofa,
		Type:       "fkjds",
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
			obj := LoginObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ValidateTwoFa(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ValidateTwoFa() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ValidateTwoFa() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.ValidateTwoFARequest{
		LoginID:    "CLIENT1",
		TwofaToken: "jfksd",
		Twofa:      twofa,
		Type:       "fkjds",
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
			obj := LoginObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ValidateTwoFa(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("LoginObj.ValidateTwoFa() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("LoginObj.ValidateTwoFa() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.ValidateTwoFARequest{
		LoginID:    "CLIENT1",
		TwofaToken: "jfksd",
		Twofa:      twofa,
		Type:       "fkjds",
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
		var tradelabError TradeLabLoginRes
		// tradelabError.Error.Code = 0
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
			obj := LoginObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ValidateTwoFa(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("LoginObj.ValidateTwoFa() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("LoginObj.ValidateTwoFa() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.ValidateTwoFARequest{
		LoginID:    "CLIENT1",
		TwofaToken: "jfksd",
		Twofa:      twofa,
		Type:       "fkjds",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	var validateTwoFaRes models.ValidateTwoFAResponse
	validateTwoFaRes.ResetPassword = true
	validateTwoFaRes.ResetTwoFa = true
	validateTwoFaRes.AuthToken = "fdsf"

	res4 := apihelpers.APIRes{
		Data:    validateTwoFaRes,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var tradelabSuccess TradeLabValidateTwoFaResponse
		tradelabSuccess.Data.ResetPassword = true
		tradelabSuccess.Data.ResetTwoFa = true
		tradelabSuccess.Data.AuthToken = "fdsf"
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
			obj := LoginObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ValidateTwoFa(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ValidateTwoFa() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ValidateTwoFa() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end

}

func TestLoginObj_SetTwoFaPin(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.SetTwoFAPinRequest
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

	req1 := models.SetTwoFAPinRequest{
		LoginID:   "CLIENT1",
		Pin:       "jfksd",
		TwofaType: "fdsf",
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
			obj := LoginObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.SetTwoFaPin(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("SetTwoFaPin() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SetTwoFaPin() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.SetTwoFAPinRequest{
		LoginID:   "CLIENT1",
		Pin:       "jfksd",
		TwofaType: "fdsf",
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
			obj := LoginObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.SetTwoFaPin(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("LoginObj.SetTwoFaPin() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("LoginObj.SetTwoFaPin() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.SetTwoFAPinRequest{
		LoginID:   "CLIENT1",
		Pin:       "jfksd",
		TwofaType: "fdsf",
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
		var tradelabError TradeLabLoginRes
		// tradelabError.Error.Code = 0
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
			obj := LoginObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.SetTwoFaPin(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("LoginObj.SetTwoFaPin() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("LoginObj.SetTwoFaPin() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	// test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.SetTwoFAPinRequest{
		LoginID:   "CLIENT1",
		Pin:       "jfksd",
		TwofaType: "fdsf",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	res4 := apihelpers.APIRes{
		Data:    nil,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		tradelabSuccess := TradelabSetTwoFaPinResponse{}
		jsonRes, _ := json.Marshal(tradelabSuccess)
		var res http.Response
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
			obj := LoginObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.SetTwoFaPin(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("SetTwoFaPin() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SetTwoFaPin() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end

}

func TestLoginObj_ForgetPassword(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.ForgotPasswordRequest
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

	req1 := models.ForgotPasswordRequest{
		LoginID: "CLIENT1",
		Pan:     "jfksd",
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
			obj := LoginObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ForgetPassword(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ForgetPassword() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ForgetPassword() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.ForgotPasswordRequest{
		LoginID: "CLIENT1",
		Pan:     "jfksd",
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
			obj := LoginObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ForgetPassword(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("LoginObj.ForgetPassword() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("LoginObj.ForgetPassword() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.ForgotPasswordRequest{
		LoginID: "CLIENT1",
		Pan:     "jfksd",
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
		var tradelabError TradeLabLoginRes
		// tradelabError.Error.Code = 0
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
			obj := LoginObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ForgetPassword(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("LoginObj.ForgetPassword() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("LoginObj.ForgetPassword() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.ForgotPasswordRequest{
		LoginID: "CLIENT1",
		Pan:     "jfksd",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	res4 := apihelpers.APIRes{
		Data:    nil,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var res http.Response
		tradelabSuccess := TradelabSetTwoFaPinResponse{}
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
			obj := LoginObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.ForgetPassword(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ForgetPassword() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ForgetPassword() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end

}

func TestLoginObj_SetPassword(t *testing.T) {
	type fields struct {
		tradeLabURL string
	}
	type args struct {
		req  models.SetPasswordRequest
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

	req1 := models.SetPasswordRequest{
		NewPass: "fdsf",
		OldPass: "hgrdf",
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
			obj := LoginObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.SetPassword(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("SetPassword() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SetPassword() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	// test 2 start
	field2 := fields{
		tradeLabURL: "http://test",
	}

	req2 := models.SetPasswordRequest{
		NewPass: "fsdsdf",
		OldPass: "sdfdfs",
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
			obj := LoginObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.SetPassword(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("LoginObj.SetPassword() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("LoginObj.SetPassword() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	//test 3 start
	field3 := fields{
		tradeLabURL: "http://test",
	}

	req3 := models.SetPasswordRequest{
		NewPass: "fsdsdf",
		OldPass: "sdfdfs",
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
		var tradelabError TradeLabLoginRes
		// tradelabError.Error.Code = 0
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
			obj := LoginObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.SetPassword(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("LoginObj.SetPassword() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("LoginObj.SetPassword() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 3 end

	//test 4 start
	field4 := fields{
		tradeLabURL: "http://test",
	}

	req4 := models.SetPasswordRequest{
		NewPass: "fsdsdf",
		OldPass: "sdfdfs",
	}

	reqH4 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}

	arg4 := args{
		req:  req4,
		reqH: reqH4,
	}

	res4 := apihelpers.APIRes{
		Data:    nil,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	apihelpers.CallAPIFunc = func(methodType, url string, payload *bytes.Buffer, deviceType, deviceid, userAgent, remoteAddr, authToken string) (*http.Response, error) {
		var res http.Response
		tradelabSuccess := TradeLabSetPasswordRes{}
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
			obj := LoginObj{
				tradeLabURL: tt.fields.tradeLabURL,
			}
			got, got1 := obj.SetPassword(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("SetPassword() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SetPassword() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 4 end

}
