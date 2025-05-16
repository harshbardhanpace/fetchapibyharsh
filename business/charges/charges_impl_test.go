package charges

import (
	"fmt"
	"net/http"
	"reflect"
	apihelpers "space/apiHelpers"
	"space/business/tradelab"
	"space/constants"
	"space/loggerconfig"
	"space/models"
	"testing"
)

func TestChargesObj_FundsPayout(t *testing.T) {
	type fields struct {
	}
	type args struct {
		req  models.FundsPayoutReq
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

	// Create mock objects and data for the test
	mockFetchFunds := func(fetchFundsReq models.FetchFundsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
		return http.StatusOK, apihelpers.APIRes{}
	}

	mockGetPositions := func(getPositionsReq models.GetPositionRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
		return http.StatusOK, apihelpers.APIRes{}
	}

	mockCompletedOrder := func(completedOrderRequest models.CompletedOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
		return http.StatusOK, apihelpers.APIRes{}
	}

	mockBrokerCharges := func(brokerChargesReq models.BrokerChargesReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
		return http.StatusOK, apihelpers.APIRes{}
	}

	// Replace the actual functions with mock functions
	originalFetchFunds := tradelab.FetchFundsInternal
	originalGetPositions := tradelab.GetPositionsInternal
	originalCompletedOrder := tradelab.CompletedOrderInternal
	originalBrokerCharges := BrokerChargesInternal

	tradelab.FetchFundsInternal = mockFetchFunds
	tradelab.GetPositionsInternal = mockGetPositions
	tradelab.CompletedOrderInternal = mockCompletedOrder
	BrokerChargesInternal = mockBrokerCharges

	defer func() {
		// Restoring the original functions after testing
		tradelab.FetchFundsInternal = originalFetchFunds
		tradelab.GetPositionsInternal = originalGetPositions
		tradelab.CompletedOrderInternal = originalCompletedOrder
		BrokerChargesInternal = originalBrokerCharges
	}()

	field1 := fields{}

	req1 := models.FundsPayoutReq{
		ClientID: "CLIENT1",
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
			obj := ChargesObj{}
			got, got1 := obj.FundsPayout(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ChargesObj.FundsPayout() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ChargesObj.FundsPayout() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	// Create mock objects and data for the test
	mockFetchFunds = func(fetchFundsReq models.FetchFundsRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
		dummyResponse := models.FetchFundsResponse{
			ClientID: "12345",
			Headers:  []string{"Header1", "Header2"},
			Values: []models.FetchFundsResponseValues{
				{
					Num0: constants.OpeningBalance,
					Num1: "4423.23",
				},
				{
					Num0: constants.MarginUsed,
					Num1: "32.42",
				},
				// Add more values as needed
			},
		}
		var apiRes apihelpers.APIRes
		apiRes.Data = dummyResponse
		apiRes.Message = "SUCCESS"
		apiRes.Status = true
		return http.StatusOK, apiRes
	}

	mockGetPositions = func(getPositionsReq models.GetPositionRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
		var allGetPositionResponseData []models.GetPositionResponseData
		var getPositionResponseData models.GetPositionResponseData
		getPositionResponseData.NetQuantity = 0
		getPositionResponseData.NetAmount = 234
		allGetPositionResponseData = append(allGetPositionResponseData, getPositionResponseData)

		var apiRes apihelpers.APIRes
		apiRes.Data = allGetPositionResponseData
		apiRes.Message = "SUCCESS"
		apiRes.Status = true
		return http.StatusOK, apiRes
	}

	mockCompletedOrder = func(completedOrderRequest models.CompletedOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
		var allCompletedOrderResponse models.CompletedOrderResponse
		var completedOrderResponseOrders models.CompletedOrderResponseOrders
		completedOrderResponseOrders.OrderStatus = constants.Completed
		completedOrderResponseOrders.ClientID = "CLIENT1"
		completedOrderResponseOrders.Price = 23
		completedOrderResponseOrders.Quantity = 2
		completedOrderResponseOrders.OrderSide = "BUY"
		completedOrderResponseOrders.ClientID = "Client1"
		completedOrderResponseOrders.Exchange = "NSE"
		completedOrderResponseOrders.TradingSymbol = "FUT"
		allCompletedOrderResponse.Orders = append(allCompletedOrderResponse.Orders, completedOrderResponseOrders)
		var apiRes apihelpers.APIRes
		apiRes.Data = allCompletedOrderResponse
		apiRes.Message = "SUCCESS"
		apiRes.Status = true
		return http.StatusOK, apiRes
	}

	mockBrokerCharges = func(brokerChargesReq models.BrokerChargesReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
		var brokerChargesRes models.BrokerChargesRes
		brokerChargesRes.Price = 23
		brokerChargesRes.Brokerage = 2.2
		brokerChargesRes.SttOrCtt = 2
		brokerChargesRes.TransactionCharges = 1
		brokerChargesRes.SebiCharges = 1
		brokerChargesRes.Gst = 1
		brokerChargesRes.StampCharges = 1
		var apiRes apihelpers.APIRes
		apiRes.Data = brokerChargesRes
		apiRes.Message = "SUCCESS"
		apiRes.Status = true
		return http.StatusOK, apiRes
	}

	tradelab.FetchFundsInternal = mockFetchFunds
	tradelab.GetPositionsInternal = mockGetPositions
	tradelab.CompletedOrderInternal = mockCompletedOrder
	BrokerChargesInternal = mockBrokerCharges

	req2 := models.FundsPayoutReq{
		ClientID: "CLIENT1",
	}
	reqH2 := models.ReqHeader{
		DeviceType:    "scascc",
		Authorization: "Bearer 12e23r2d",
	}
	arg2 := args{
		req:  req2,
		reqH: reqH2,
	}
	var fundsPayoutRes models.FundsPayoutRes
	fundsPayoutRes.ClientID = "CLIENT1"
	fundsPayoutRes.PayoutAmount = 4148.61
	fundsPayoutRes.OpeningBalance = 4423.23
	fundsPayoutRes.MarginUsed = 32.42
	fundsPayoutRes.LossOnClosedPositions = 234
	fundsPayoutRes.ChargesOnTrades = 8.2
	res2 := apihelpers.APIRes{
		Data:    fundsPayoutRes,
		Message: "SUCCESS",
		Status:  true,
	}

	field := fields{}

	tests = []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  apihelpers.APIRes
	}{
		// TODO: Add test cases.
		{"Success", field, arg2, http.StatusOK, res2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := ChargesObj{}
			got, got1 := obj.FundsPayout(tt.args.req, tt.args.reqH)
			if got != tt.want {
				t.Errorf("ChargesObj.FundsPayout() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ChargesObj.FundsPayout() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

}
