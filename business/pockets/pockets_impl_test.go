package pockets

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/loggerconfig"
	"space/models"
	"testing"
)

func Test_FetchPockets(t *testing.T) {
	type args struct {
		req  models.FetchPocketsDetailsRequest
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

	req1 := models.FetchPocketsDetailsRequest{
		PocketId: "2645d918-1216-4f20-8858-4059eef55a09",
	}

	reqH1 := models.ReqHeader{
		DeviceType:    "web",
		Authorization: "Bearer 12e23r2f",
	}

	arg1 := args{
		req:  req1,
		reqH: reqH1,
	}

	res1 := apihelpers.APIRes{
		Status:    false,
		Message:   constants.ErrorCodeMap[constants.PocketDoesNotExists],
		ErrorCode: constants.PocketDoesNotExists,
	}

	//mock call api
	CallFetchPocketMongo = func(req models.FetchPocketsDetailsRequest, obj ExecutePocketV2Obj) (models.MongoPockets, error) {
		var pockets models.MongoPockets
		return pockets, errors.New("Call Api Error")
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
			obj := ExecutePocketV2Obj{}
			got, got1 := obj.FetchPockets(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("FetchPockets() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("FetchPockets() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	//test 2 start
	req2 := models.FetchPocketsDetailsRequest{
		PocketId: "2645d918-1216-4f20-8858-4059eef55a09",
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "web",
		Authorization: "Bearer 12e23r2f",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	var resp models.FetchPocketsDetailsResponse
	resp.PocketExchange = "NSE"
	resp.PocketImage = "https://spacepocket.s3.ap-south-1.amazonaws.com/Financials.png"
	resp.PocketLongDesc = "--"
	resp.PocketName = "Financials"
	resp.PocketShortDesc = "-"
	resp.PocketId = "2645d918-1216-4f20-8858-4059eef55a09"

	res2 := apihelpers.APIRes{
		Data:    resp,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	CallFetchPocketMongo = func(req models.FetchPocketsDetailsRequest, obj ExecutePocketV2Obj) (models.MongoPockets, error) {
		var pockets models.MongoPockets
		pockets.PocketExchange = "NSE"
		pockets.PocketId = "2645d918-1216-4f20-8858-4059eef55a09"
		pockets.PocketImage = "https://spacepocket.s3.ap-south-1.amazonaws.com/Financials.png"
		pockets.PocketLongDesc = "--"
		pockets.PocketName = "Financials"
		pockets.PocketShortDesc = "-"
		return pockets, nil
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
			obj := ExecutePocketV2Obj{}
			got, got1 := obj.FetchPockets(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("FetchPockets() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("FetchPockets() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end
}
