package upipreference

import (
	"errors"
	"net/http"
	"reflect"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/models"
	"testing"
)

func Test_SetUpiPreference(t *testing.T) {
	type args struct {
		req  models.SetUpiPreferenceReq
		reqH models.ReqHeader
	}

	req1 := models.SetUpiPreferenceReq{
		ClientId: "CLIENT1",
		UpiId:    "1234",
	}

	reqH := models.ReqHeader{
		DeviceType:    "web",
		Authorization: "Bearer 12e23r2f",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH,
	}

	res1 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.InvalidRequest],
		ErrorCode: constants.InvalidRequest,
	}

	//mock call api
	CallFetchUpiPreferenceMongo = func(clientId string, obj UpiPreferenceObj) (models.UpiPreference, error) {
		var upiPreference models.UpiPreference
		return upiPreference, errors.New("Error in connection")
	}

	tests := []struct {
		name  string
		args  args
		want  int
		want1 apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"call api error", arg1, http.StatusBadRequest, res1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := UpiPreferenceObj{}
			got, got1 := obj.SetUpiPreference(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("SetUpiPreference() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SetUpiPreference() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	req2 := models.SetUpiPreferenceReq{
		ClientId: "CLIENT1",
		UpiId:    "1234",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH,
	}

	res2 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.InternalServerError],
		ErrorCode: constants.InternalServerError,
	}

	CallFetchUpiPreferenceMongo = func(clientId string, obj UpiPreferenceObj) (models.UpiPreference, error) {
		var upiPreference models.UpiPreference
		upiPreference.ClientId = "CLIENT1"
		upiPreference.UpiIds = []string{"5678"}
		return upiPreference, nil
	}

	CallUpdateUpiPreferenceMongo = func(clientId string, upiPreference models.UpiPreference, obj UpiPreferenceObj) error {
		return errors.New("Error in connection")
	}

	tests = []struct {
		name  string
		args  args
		want  int
		want1 apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"INTERNAL SERVER ERROR", arg2, http.StatusInternalServerError, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := UpiPreferenceObj{}
			got, got1 := obj.SetUpiPreference(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("SetUpiPreference() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SetUpiPreference() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	req3 := models.SetUpiPreferenceReq{
		ClientId: "CLIENT1",
		UpiId:    "1234",
	}

	arg3 := args{
		req:  req3,
		reqH: reqH,
	}

	res3 := apihelpers.APIRes{
		Data:    nil,
		Message: "SUCCESS",
		Status:  true,
	}

	CallFetchUpiPreferenceMongo = func(clientId string, obj UpiPreferenceObj) (models.UpiPreference, error) {
		var upiPreference models.UpiPreference
		upiPreference.ClientId = "CLIENT1"
		upiPreference.UpiIds = []string{"5678"}
		return upiPreference, nil
	}

	CallUpdateUpiPreferenceMongo = func(clientId string, upiPreference models.UpiPreference, obj UpiPreferenceObj) error {
		return nil
	}

	tests = []struct {
		name  string
		args  args
		want  int
		want1 apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", arg3, http.StatusOK, res3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := UpiPreferenceObj{}
			got, got1 := obj.SetUpiPreference(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("SetUpiPreference() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SetUpiPreference() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

}

func Test_FetchUpiPreference(t *testing.T) {
	type args struct {
		req  models.FetchUpiPreferenceReq
		reqH models.ReqHeader
	}

	req1 := models.FetchUpiPreferenceReq{
		ClientId: "CLIENT1",
	}

	reqH := models.ReqHeader{
		DeviceType:    "web",
		Authorization: "Bearer 12e23r2f",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH,
	}

	res1 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.InvalidRequest],
		ErrorCode: constants.InvalidRequest,
	}

	//mock call api
	CallFetchUpiPreferenceMongo = func(clientId string, obj UpiPreferenceObj) (models.UpiPreference, error) {
		var upiPreference models.UpiPreference
		return upiPreference, errors.New("Error in connection")
	}

	tests := []struct {
		name  string
		args  args
		want  int
		want1 apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"call api error", arg1, http.StatusBadRequest, res1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := UpiPreferenceObj{}
			got, got1 := obj.FetchUpiPreference(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("FetchUpiPreference() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("FetchUpiPreference() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	req2 := models.FetchUpiPreferenceReq{
		ClientId: "CLIENT1",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH,
	}

	var fetchUpiPreferenceRes models.FetchUpiPreferenceRes
	fetchUpiPreferenceRes.ClientId = "CLIENT1"
	fetchUpiPreferenceRes.UpiIds = []string{"5678"}

	res2 := apihelpers.APIRes{
		Data:    fetchUpiPreferenceRes,
		Message: "SUCCESS",
		Status:  true,
	}

	CallFetchUpiPreferenceMongo = func(clientId string, obj UpiPreferenceObj) (models.UpiPreference, error) {
		var upiPreference models.UpiPreference
		upiPreference.ClientId = "CLIENT1"
		upiPreference.UpiIds = []string{"5678"}
		return upiPreference, nil
	}

	tests = []struct {
		name  string
		args  args
		want  int
		want1 apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", arg2, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := UpiPreferenceObj{}
			got, got1 := obj.FetchUpiPreference(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("FetchUpiPreference() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("FetchUpiPreference() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

}

func Test_DeleteUpiPreference(t *testing.T) {
	type args struct {
		req  models.DeleteUpiPreferenceReq
		reqH models.ReqHeader
	}

	req1 := models.DeleteUpiPreferenceReq{
		ClientId: "CLIENT1",
		UpiIds:   []string{"1234"},
	}

	reqH := models.ReqHeader{
		DeviceType:    "web",
		Authorization: "Bearer 12e23r2f",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH,
	}

	res1 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.InvalidRequest],
		ErrorCode: constants.InvalidRequest,
	}

	//mock call api
	CallFetchUpiPreferenceMongo = func(clientId string, obj UpiPreferenceObj) (models.UpiPreference, error) {
		var upiPreference models.UpiPreference
		return upiPreference, errors.New("Error in connection")
	}

	tests := []struct {
		name  string
		args  args
		want  int
		want1 apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"call api error", arg1, http.StatusBadRequest, res1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := UpiPreferenceObj{}
			got, got1 := obj.DeleteUpiPreference(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("DeleteUpiPreference() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("DeleteUpiPreference() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	req2 := models.DeleteUpiPreferenceReq{
		ClientId: "CLIENT1",
		UpiIds:   []string{"1234"},
	}

	arg2 := args{
		req:  req2,
		reqH: reqH,
	}

	res2 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.UpiDontExist],
		ErrorCode: constants.UpiDontExist,
	}

	CallFetchUpiPreferenceMongo = func(clientId string, obj UpiPreferenceObj) (models.UpiPreference, error) {
		var upiPreference models.UpiPreference
		upiPreference.ClientId = "CLIENT1"
		upiPreference.UpiIds = []string{"5678"}
		return upiPreference, nil
	}

	CallUpdateUpiPreferenceMongo = func(clientId string, upiPreference models.UpiPreference, obj UpiPreferenceObj) error {
		return errors.New("Error in connection")
	}

	tests = []struct {
		name  string
		args  args
		want  int
		want1 apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Upi Don't Exist", arg2, http.StatusBadRequest, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := UpiPreferenceObj{}
			got, got1 := obj.DeleteUpiPreference(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("DeleteUpiPreference() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("DeleteUpiPreference() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end

	req3 := models.DeleteUpiPreferenceReq{
		ClientId: "CLIENT1",
		UpiIds:   []string{"1234"},
	}

	arg3 := args{
		req:  req3,
		reqH: reqH,
	}

	res3 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.InternalServerError],
		ErrorCode: constants.InternalServerError,
	}

	CallFetchUpiPreferenceMongo = func(clientId string, obj UpiPreferenceObj) (models.UpiPreference, error) {
		var upiPreference models.UpiPreference
		upiPreference.ClientId = "CLIENT1"
		upiPreference.UpiIds = []string{"1234"}
		return upiPreference, nil
	}

	CallUpdateUpiPreferenceMongo = func(clientId string, upiPreference models.UpiPreference, obj UpiPreferenceObj) error {
		return errors.New("Error in connection")
	}

	tests = []struct {
		name  string
		args  args
		want  int
		want1 apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"INTERNAL SERVER ERROR", arg3, http.StatusInternalServerError, res3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := UpiPreferenceObj{}
			got, got1 := obj.DeleteUpiPreference(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("DeleteUpiPreference() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("DeleteUpiPreference() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	req4 := models.DeleteUpiPreferenceReq{
		ClientId: "CLIENT1",
		UpiIds:   []string{"1234"},
	}

	arg4 := args{
		req:  req4,
		reqH: reqH,
	}

	res4 := apihelpers.APIRes{
		Data:    nil,
		Message: "SUCCESS",
		Status:  true,
	}

	CallFetchUpiPreferenceMongo = func(clientId string, obj UpiPreferenceObj) (models.UpiPreference, error) {
		var upiPreference models.UpiPreference
		upiPreference.ClientId = "CLIENT1"
		upiPreference.UpiIds = []string{"1234"}
		return upiPreference, nil
	}

	CallUpdateUpiPreferenceMongo = func(clientId string, upiPreference models.UpiPreference, obj UpiPreferenceObj) error {
		return nil
	}

	tests = []struct {
		name  string
		args  args
		want  int
		want1 apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", arg4, http.StatusOK, res4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := UpiPreferenceObj{}
			got, got1 := obj.DeleteUpiPreference(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("DeleteUpiPreference() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("DeleteUpiPreference() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

}
