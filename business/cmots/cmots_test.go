package cmots

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

func TestCmotsObj_FetchFinancialsV2(t *testing.T) {
	type fields struct {
		CmURL  string
		CmAuth string
	}
	type args struct {
		req  models.FetchFinancialsReq
		reqH models.ReqHeader
	}

	field1 := fields{
		CmURL:  "http://test",
		CmAuth: "Bearer Auth",
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

	req1 := models.FetchFinancialsReq{
		Exchange:  "BSE",
		BseToken:  "500410",
		NseSymbol: "",
		Isin:      "",
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
		Message:   constants.ErrorCodeMap[constants.InternalServerError],
		ErrorCode: constants.InternalServerError,
	}

	//mock call api
	CallFetchFinancialsDataV2 = func(obj CmotsObj, req models.FetchFinancialsReq) (models.FetchFinancialsV2Res, error) {
		var res models.FetchFinancialsV2Res
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
			obj := CmotsObj{
				CmURL:  tt.fields.CmURL,
				CmAuth: tt.fields.CmAuth,
			}
			got, got1 := obj.FetchFinancialsV2(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("CmotsObj.FetchFinancialsV2() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CmotsObj.FetchFinancialsV2() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	//test 1 end

	//test 2 start
	field2 := fields{
		CmURL:  "http://test",
		CmAuth: "Bearer Auth",
	}

	req2 := models.FetchFinancialsReq{
		Exchange:  "BSE",
		BseToken:  "500410",
		NseSymbol: "",
		Isin:      "",
	}

	reqH2 := models.ReqHeader{
		DeviceType:    "web",
		Authorization: "Bearer 12e23r2f",
	}

	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}

	var financials models.FetchFinancialsV2Res
	financials.NetProfit.CoCode = 69
	financials.NetProfit.ColumnName = "69"
	financials.NetProfit.Y0 = 69
	financials.NetProfit.Y1 = 69
	financials.NetProfit.Y2 = 69
	financials.NetProfit.Y3 = 69
	financials.NetProfit.Y4 = 69

	financials.Revenue.CoCode = 69
	financials.Revenue.ColumnName = "69"
	financials.Revenue.Y0 = 69
	financials.Revenue.Y1 = 69
	financials.Revenue.Y2 = 69
	financials.Revenue.Y3 = 69
	financials.Revenue.Y4 = 69

	financials.Cashflow.CoCode = 69
	financials.Cashflow.ColumnName = "69"
	financials.Cashflow.Y0 = 69
	financials.Cashflow.Y1 = 69
	financials.Cashflow.Y2 = 69
	financials.Cashflow.Y3 = 69
	financials.Cashflow.Y4 = 69

	financials.BalanceSheet.TotalAssets.CoCode = 69
	financials.BalanceSheet.TotalAssets.ColumnName = "69"
	financials.BalanceSheet.TotalAssets.Y0 = 69
	financials.BalanceSheet.TotalAssets.Y1 = 69
	financials.BalanceSheet.TotalAssets.Y2 = 69
	financials.BalanceSheet.TotalAssets.Y3 = 69
	financials.BalanceSheet.TotalAssets.Y4 = 69

	financials.BalanceSheet.TotalLiabilities.CoCode = 69
	financials.BalanceSheet.TotalLiabilities.ColumnName = "69"
	financials.BalanceSheet.TotalLiabilities.Y0 = 69
	financials.BalanceSheet.TotalLiabilities.Y1 = 69
	financials.BalanceSheet.TotalLiabilities.Y2 = 69
	financials.BalanceSheet.TotalLiabilities.Y3 = 69
	financials.BalanceSheet.TotalLiabilities.Y4 = 69

	res2 := apihelpers.APIRes{
		Data:    financials,
		Message: "SUCCESS",
		Status:  true,
	}

	//mock call api
	CallFetchFinancialsDataV2 = func(obj CmotsObj, req models.FetchFinancialsReq) (models.FetchFinancialsV2Res, error) {
		var financial models.FetchFinancialsV2Res
		financial.NetProfit.CoCode = 69
		financial.NetProfit.ColumnName = "69"
		financial.NetProfit.Y0 = 69
		financial.NetProfit.Y1 = 69
		financial.NetProfit.Y2 = 69
		financial.NetProfit.Y3 = 69
		financial.NetProfit.Y4 = 69
		financial.Revenue.CoCode = 69
		financial.Revenue.ColumnName = "69"
		financial.Revenue.Y0 = 69
		financial.Revenue.Y1 = 69
		financial.Revenue.Y2 = 69
		financial.Revenue.Y3 = 69
		financial.Revenue.Y4 = 69
		financial.Cashflow.CoCode = 69
		financial.Cashflow.ColumnName = "69"
		financial.Cashflow.Y0 = 69
		financial.Cashflow.Y1 = 69
		financial.Cashflow.Y2 = 69
		financial.Cashflow.Y3 = 69
		financial.Cashflow.Y4 = 69
		financial.BalanceSheet.TotalAssets.CoCode = 69
		financial.BalanceSheet.TotalAssets.ColumnName = "69"
		financial.BalanceSheet.TotalAssets.Y0 = 69
		financial.BalanceSheet.TotalAssets.Y1 = 69
		financial.BalanceSheet.TotalAssets.Y2 = 69
		financial.BalanceSheet.TotalAssets.Y3 = 69
		financial.BalanceSheet.TotalAssets.Y4 = 69
		financial.BalanceSheet.TotalLiabilities.CoCode = 69
		financial.BalanceSheet.TotalLiabilities.ColumnName = "69"
		financial.BalanceSheet.TotalLiabilities.Y0 = 69
		financial.BalanceSheet.TotalLiabilities.Y1 = 69
		financial.BalanceSheet.TotalLiabilities.Y2 = 69
		financial.BalanceSheet.TotalLiabilities.Y3 = 69
		financial.BalanceSheet.TotalLiabilities.Y4 = 69

		return financial, nil
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
			obj := CmotsObj{
				CmURL:  tt.fields.CmURL,
				CmAuth: tt.fields.CmAuth,
			}
			got, got1 := obj.FetchFinancialsV2(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("CmotsObj.FetchFinancialsV2() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CmotsObj.FetchFinancialsV2() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	//test 2 end
}
