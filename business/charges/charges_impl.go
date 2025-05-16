package charges

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	apihelpers "space/apiHelpers"
	"space/business/tradelab"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
)

type ChargesObj struct {
	// FundsObj     tradelab.FundsObj // Embed the dependency
	// PortfolioObj tradelab.PortfolioObj
	// OrderObj     tradelab.OrderObj
}

var objCharges ChargesObj

func InitChargesProvider() ChargesObj {
	defer models.HandlePanic()
	chargesObj := ChargesObj{}
	objCharges = chargesObj

	return chargesObj
}

func (obj ChargesObj) CombineBrokerCharges(combineBrokerChargesReq models.CombineBrokerChargesReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return CombineBrokerCharges(combineBrokerChargesReq, obj, reqH)
}

var CombineBrokerCharges = func(combineBrokerChargesReq models.CombineBrokerChargesReq, obj ChargesObj, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	var combineBrokerChargesRes models.CombineBrokerChargesRes

	for i := 0; i < len(combineBrokerChargesReq.BrokerCharges); i++ {
		status, res := BrokerChargesInternal(combineBrokerChargesReq.BrokerCharges[i], reqH)
		if status != http.StatusOK {
			loggerconfig.Error("CombineBrokerCharges in BrokerCharges status != 200", status, " uccId:", combineBrokerChargesReq.ClientID, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		brokerChargesRes, ok := res.Data.(models.BrokerChargesRes)
		if !ok {
			loggerconfig.Error("CombineBrokerCharges in BrokerCharges interface parsing error", ok, " uccId:", combineBrokerChargesReq.ClientID, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		combineBrokerChargesRes.BrokerCharges = append(combineBrokerChargesRes.BrokerCharges, brokerChargesRes)
	}

	loggerconfig.Info("CombineBrokerCharges response=", helpers.LogStructAsJSON(combineBrokerChargesRes), " uccId:", combineBrokerChargesReq.ClientID, " requestId:", reqH.RequestId)
	apiRes.Data = combineBrokerChargesRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}

func (obj ChargesObj) BrokerCharges(brokerChargesReq models.BrokerChargesReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	if brokerChargesReq.Product != "" && brokerChargesReq.Product == "MTF" {
		var apiRes apihelpers.APIRes
		var brokerChargesRes models.BrokerChargesRes
		brokerChargesRes = SegmentEquity(brokerChargesReq)
		inititalPrice := brokerChargesReq.Price * float64(brokerChargesReq.Quantity)
		brokerage := inititalPrice * constants.EquityDeliveryBrokeragePocketful
		brokerChargesRes.Brokerage = ceilToTwoDecimalPlaces(brokerage)
		if brokerChargesRes.Brokerage > 20.0 { //maximum of Rs.20 in brokerage
			brokerChargesRes.Brokerage = 20.0
		}
		brokerChargesRes.SttOrCtt = ceilToTwoDecimalPlaces(brokerChargesRes.SttOrCtt)
		brokerChargesRes.TransactionCharges = ceilToTwoDecimalPlaces(brokerChargesRes.TransactionCharges)
		brokerChargesRes.SebiCharges = ceilToTwoDecimalPlaces(brokerChargesRes.SebiCharges)
		brokerChargesRes.Gst = ceilToTwoDecimalPlaces(brokerChargesRes.Gst)
		brokerChargesRes.StampCharges = ceilToTwoDecimalPlaces(brokerChargesRes.StampCharges)
		brokerChargesRes.StampCharges = math.Floor(brokerChargesRes.StampCharges)

		total := brokerChargesRes.Brokerage + brokerChargesRes.SttOrCtt + brokerChargesRes.TransactionCharges + brokerChargesRes.SebiCharges + brokerChargesRes.Gst + math.Floor(brokerChargesRes.StampCharges)
		brokerChargesRes.TotalCharge = ceilToTwoDecimalPlaces(total)

		loggerconfig.Info("BrokerCharges brokarage charages response for MTF=", helpers.LogStructAsJSON(brokerChargesRes), " uccId:", brokerChargesReq.ClientID, " requestId:", reqH.RequestId)
		apiRes.Data = brokerChargesRes
		apiRes.Message = "SUCCESS"
		apiRes.Status = true
		return http.StatusOK, apiRes
	}
	return BrokerChargesInternal(brokerChargesReq, reqH)
}

var BrokerChargesInternal = func(brokerChargesReq models.BrokerChargesReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	var brokerChargesRes models.BrokerChargesRes

	if strings.ToLower(brokerChargesReq.Segment) == constants.EQUITY {
		brokerChargesRes = SegmentEquity(brokerChargesReq)
	} else if strings.ToLower(brokerChargesReq.Segment) == constants.CURRENCY {
		brokerChargesRes = SegmentCurrency(brokerChargesReq)
	} else if strings.ToLower(brokerChargesReq.Segment) == constants.COMMODITY {
		brokerChargesRes = SegmentCommodity(brokerChargesReq)
	}

	brokerChargesRes.Brokerage = ceilToTwoDecimalPlaces(brokerChargesRes.Brokerage)
	brokerChargesRes.SttOrCtt = ceilToTwoDecimalPlaces(brokerChargesRes.SttOrCtt)
	brokerChargesRes.TransactionCharges = ceilToTwoDecimalPlaces(brokerChargesRes.TransactionCharges)
	brokerChargesRes.SebiCharges = ceilToTwoDecimalPlaces(brokerChargesRes.SebiCharges)
	brokerChargesRes.Gst = ceilToTwoDecimalPlaces(brokerChargesRes.Gst)
	brokerChargesRes.StampCharges = ceilToTwoDecimalPlaces(brokerChargesRes.StampCharges)

	total := brokerChargesRes.Brokerage + brokerChargesRes.SttOrCtt + brokerChargesRes.TransactionCharges + brokerChargesRes.SebiCharges + brokerChargesRes.Gst + math.Floor(brokerChargesRes.StampCharges)
	brokerChargesRes.TotalCharge = ceilToTwoDecimalPlaces(total)

	loggerconfig.Info("brokarage charages response=", helpers.LogStructAsJSON(brokerChargesRes), " uccId:", brokerChargesReq.ClientID, " requestId:", reqH.RequestId)
	apiRes.Data = brokerChargesRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}

func SegmentEquity(brokerChargesReq models.BrokerChargesReq) models.BrokerChargesRes {

	inititalPrice := brokerChargesReq.Price * float64(brokerChargesReq.Quantity)

	var brokerage float64
	var sttOrCtt float64
	var transactionCharges float64
	var gst float64
	var sebiCharges float64
	var stampCharges float64

	var totalPrice float64
	if brokerChargesReq.SubSegment == constants.DELIVERY {
		// Brokerage
		brokerage = constants.EquityDeliveryBrokerage

		// STT/CTT on both buy and sell
		sttOrCtt = reduceTenDecimalPlaces(inititalPrice * constants.EquityDeliverySttOrCtt)

		// Transaction Charges
		if strings.ToLower(brokerChargesReq.Exchange) == constants.NSE {
			transactionCharges = reduceTenDecimalPlaces(inititalPrice * constants.EquityDeliveryTransactionChargeNse)
		} else if strings.ToLower(brokerChargesReq.Exchange) == constants.BSE {
			transactionCharges = reduceTenDecimalPlaces(inititalPrice * constants.EquityDeliveryTransactionChargeBse)
		}

		// Sebi Charges (check may contain errors)
		sebiCharges = reduceTenDecimalPlaces(inititalPrice * constants.EquityDeliverySebiCharges)

		// GST
		gst = reduceTenDecimalPlaces((brokerage + sebiCharges + transactionCharges) * constants.EquityDeliveryGst)

		// Stamp Charges
		if strings.ToLower(brokerChargesReq.Process) == constants.BUY {
			stampCharges = reduceTenDecimalPlaces(inititalPrice * constants.EquityDeliveryStampProcessBuy)
		}

	} else if brokerChargesReq.SubSegment == constants.INTRADAY {

		// Brokerage
		brokerageOption1 := inititalPrice * constants.EquityIntradayBrokerageOption1
		brokerageOption2 := constants.EquityIntradayBrokerageOption2
		brokerage = reduceTenDecimalPlaces(math.Min(brokerageOption1, brokerageOption2))

		// STT/CTT on sell side
		if strings.ToLower(brokerChargesReq.Process) == constants.SELL {
			sttOrCtt = reduceTenDecimalPlaces(inititalPrice * constants.EquityIntradaySttOrCttSell)
		}

		// Transaction Charges
		if strings.ToLower(brokerChargesReq.Exchange) == constants.NSE {
			transactionCharges = reduceTenDecimalPlaces(inititalPrice * constants.EquityIntradayTransactionChargeNse)
		} else if strings.ToLower(brokerChargesReq.Exchange) == constants.BSE {
			transactionCharges = reduceTenDecimalPlaces(inititalPrice * constants.EquityIntradayTransactionChargeBse)
		}

		// Sebi Charges (check may contains error)
		sebiCharges = reduceTenDecimalPlaces(inititalPrice * constants.EquityIntradaySebiCharges)

		// GST
		gst = reduceTenDecimalPlaces(((brokerage + sebiCharges + transactionCharges) * constants.EquityIntradayGst))

		// Stamp Charges
		if strings.ToLower(brokerChargesReq.Process) == constants.BUY {
			stampCharges = reduceTenDecimalPlaces(inititalPrice * constants.EquityIntradayStampProcessBuy)
		}

	} else if brokerChargesReq.SubSegment == constants.FUTURES {

		// Brokerage
		brokerageOption1 := inititalPrice * constants.EquityFuturesBrokerageOption1
		brokerageOption2 := constants.EquityFuturesBrokerageOption2
		brokerage = reduceTenDecimalPlaces(math.Min(brokerageOption1, brokerageOption2))

		// STT/CTT on sell side
		if strings.ToLower(brokerChargesReq.Process) == constants.SELL {
			sttOrCtt = reduceTenDecimalPlaces(inititalPrice * constants.EquityFuturesSttOrCttSell)
		}

		// Transaction Charges
		if strings.ToLower(brokerChargesReq.Exchange) == constants.NSE {
			transactionCharges = reduceTenDecimalPlaces(inititalPrice * constants.EquityFuturesTransactionChargeNse)
		}

		// Sebi Charges (check may contains error)
		sebiCharges = reduceTenDecimalPlaces(inititalPrice * constants.EquityFuturesSebiCharges)

		// GST
		gst = reduceTenDecimalPlaces((brokerage + sebiCharges + transactionCharges) * constants.EquityFuturesGst)

		// Stamp Charges
		if strings.ToLower(brokerChargesReq.Process) == constants.BUY {
			stampCharges = reduceTenDecimalPlaces(inititalPrice * constants.EquityFuturesStampProcessBuy)
		}

	} else if brokerChargesReq.SubSegment == constants.OPTIONS {

		// Brokerage
		brokerage = constants.EquityOptionsBrokerage

		// STT/CTT on sell side
		if strings.ToLower(brokerChargesReq.Process) == constants.SELL {
			sttOrCtt = reduceTenDecimalPlaces(inititalPrice * constants.EquityOptionsSttOrCttSell)
		}

		// Transaction Charges (check for premium)
		if strings.ToLower(brokerChargesReq.Exchange) == constants.NSE {
			transactionCharges = reduceTenDecimalPlaces(inititalPrice * constants.EquityOptionsTransactionChargeNse)
		}

		// Sebi Charges (check may contains error)
		sebiCharges = reduceTenDecimalPlaces(inititalPrice * constants.EquityOptionsSebiCharges)

		// GST
		gst = reduceTenDecimalPlaces((brokerage + sebiCharges + transactionCharges) * constants.EquityOptionsGst)

		// Stamp Charges
		if strings.ToLower(brokerChargesReq.Process) == constants.BUY {
			stampCharges = reduceTenDecimalPlaces(inititalPrice * constants.EquityOptionsStampProcessBuy)
		}
	}

	totalPrice = reduceTenDecimalPlaces(brokerage + sttOrCtt + transactionCharges + gst + sebiCharges + stampCharges)

	var brokerChargesRes models.BrokerChargesRes
	brokerChargesRes.Price = inititalPrice + totalPrice
	brokerChargesRes.Brokerage = brokerage
	brokerChargesRes.SttOrCtt = sttOrCtt
	brokerChargesRes.TransactionCharges = transactionCharges
	brokerChargesRes.SebiCharges = sebiCharges
	brokerChargesRes.Gst = gst
	brokerChargesRes.StampCharges = stampCharges

	return brokerChargesRes
}

func SegmentCurrency(brokerChargesReq models.BrokerChargesReq) models.BrokerChargesRes {
	inititalPrice := brokerChargesReq.Price * float64(brokerChargesReq.Quantity)

	var brokerage float64
	var sttOrCtt float64
	var transactionCharges float64
	var gst float64
	var sebiCharges float64
	var stampCharges float64

	var totalPrice float64
	if brokerChargesReq.SubSegment == constants.FUTURES {

		// Brokerage
		brokerageOption1 := inititalPrice * constants.CurrencyFuturesBrokerageOption1
		brokerageOption2 := constants.CurrencyFuturesBrokerageOption2
		brokerage = reduceTenDecimalPlaces(math.Min(brokerageOption1, brokerageOption2))

		// STT/CTT
		sttOrCtt = constants.CurrencyFuturesSttOrCttSell

		// Transaction Charges
		if strings.ToLower(brokerChargesReq.Exchange) == constants.NSE {
			transactionCharges = reduceTenDecimalPlaces(inititalPrice * constants.CurrencyFuturesTransactionChargeNse)
		} else if strings.ToLower(brokerChargesReq.Exchange) == constants.BSE {
			transactionCharges = reduceTenDecimalPlaces(inititalPrice * constants.CurrencyFuturesTransactionChargeBse)
		}

		// Sebi Charges (check may contains error)
		sebiCharges = reduceTenDecimalPlaces(inititalPrice * constants.CurrencyFuturesSebiCharges)

		// GST
		gst = reduceTenDecimalPlaces((brokerage + transactionCharges) * constants.CurrencyFuturesGst)

		// Stamp Charges
		if strings.ToLower(brokerChargesReq.Process) == constants.BUY {
			stampCharges = reduceTenDecimalPlaces(inititalPrice * constants.CurrencyFuturesStampProcessBuy)
		}

	} else if brokerChargesReq.SubSegment == constants.OPTIONS {
		// Brokerage
		brokerage = constants.CurrencyOptionsBrokerage

		// STT/CTT
		sttOrCtt = constants.CurrencyOptionsSttOrCttSell

		// Transaction Charges
		if strings.ToLower(brokerChargesReq.Exchange) == constants.NSE {
			transactionCharges = reduceTenDecimalPlaces(inititalPrice * constants.CurrencyOptionsTransactionChargeNse)
		} else if strings.ToLower(brokerChargesReq.Exchange) == constants.BSE {
			transactionCharges = reduceTenDecimalPlaces(inititalPrice * constants.CurrencyOptionsTransactionChargeBse)
		}

		// Sebi Charges (check may contains error)
		sebiCharges = reduceTenDecimalPlaces(inititalPrice * constants.CurrencyOptionsSebiCharges)

		// GST Charges
		gst = reduceTenDecimalPlaces((brokerage + transactionCharges) * constants.CurrencyOptionsGst)

		// Stamp Charges
		if strings.ToLower(brokerChargesReq.Process) == constants.BUY {
			stampCharges = reduceTenDecimalPlaces(inititalPrice * constants.CurrencyOptionsStampProcessBuy)
		}
	}

	totalPrice = reduceTenDecimalPlaces(brokerage + sttOrCtt + transactionCharges + gst + sebiCharges + stampCharges)

	var brokerChargesRes models.BrokerChargesRes
	brokerChargesRes.Price = inititalPrice + totalPrice
	brokerChargesRes.Brokerage = brokerage
	brokerChargesRes.SttOrCtt = sttOrCtt
	brokerChargesRes.TransactionCharges = transactionCharges
	brokerChargesRes.SebiCharges = sebiCharges
	brokerChargesRes.Gst = gst
	brokerChargesRes.StampCharges = stampCharges

	return brokerChargesRes
}

func SegmentCommodity(brokerChargesReq models.BrokerChargesReq) models.BrokerChargesRes {
	inititalPrice := brokerChargesReq.Price * float64(brokerChargesReq.Quantity)

	var brokerage float64
	var sttOrCtt float64
	var transactionCharges float64
	var gst float64
	var sebiCharges float64
	var stampCharges float64

	var totalPrice float64

	if brokerChargesReq.SubSegment == constants.FUTURES {

		// Brokerage
		brokerageOption1 := inititalPrice * constants.CommodityFuturesBrokerageOption1
		brokerageOption2 := constants.CommodityFuturesBrokerageOption2
		brokerage = reduceTenDecimalPlaces(math.Min(brokerageOption1, brokerageOption2))

		// STT/CTT
		if !brokerChargesReq.Agri && strings.ToLower(brokerChargesReq.Process) == constants.SELL {
			sttOrCtt = reduceTenDecimalPlaces(inititalPrice * constants.CommodityFuturesSttOrCttNonAgriSell)
		}

		// Transaction Charges
		if brokerChargesReq.GroupInfo == constants.NORMALCOMMODITY { // Group 1
			transactionCharges = reduceTenDecimalPlaces(inititalPrice * constants.CommodityFuturesTransactionChargeNormal)
		} else if brokerChargesReq.GroupInfo == constants.CASTORSEEDCOMMODITY { // Group 2 CASTORSEED
			transactionCharges = reduceTenDecimalPlaces(inititalPrice * constants.CommodityFuturesTransactionChargeCastorseed)
		} else if brokerChargesReq.GroupInfo == constants.KAPASCOMMODITY { // Group 3 KAPAS
			transactionCharges = reduceTenDecimalPlaces(inititalPrice * constants.CommodityFuturesTransactionChargeKapas)
		} else if brokerChargesReq.GroupInfo == constants.PEPPERCOMMODITY { // Group 4 PEPPER
			transactionCharges = reduceTenDecimalPlaces(inititalPrice * constants.CommodityFuturesTransactionChargePepper)
		} else if brokerChargesReq.GroupInfo == constants.RBDPMOLEINCOMMODITY { // Group 5 RBDPMOLEIN
			transactionCharges = reduceTenDecimalPlaces(inititalPrice * constants.CommodityFuturesTransactionChargeRbdmolein)
		}

		// GST
		gst = reduceTenDecimalPlaces((brokerage + transactionCharges) * constants.CommodityFuturesGst)

		// Sebi Charges (check may contains error)
		if brokerChargesReq.Agri {
			sebiCharges = reduceTenDecimalPlaces(inititalPrice * constants.CommodityFuturesSebiChargesAgri)
		} else {
			sebiCharges = reduceTenDecimalPlaces(inititalPrice * constants.CommodityFuturesSebiChargesNonAgri)
		}

		// Stamp Charges
		if strings.ToLower(brokerChargesReq.Process) == constants.BUY {
			stampCharges = reduceTenDecimalPlaces(inititalPrice * constants.CommodityFuturesStampProcessBuy)
		}

	} else if brokerChargesReq.SubSegment == constants.OPTIONS {
		// Brokerage
		brokerage = constants.CommodityOptionsBrokerage

		// STT/CTT
		if strings.ToLower(brokerChargesReq.Process) == constants.SELL {
			sttOrCtt = reduceTenDecimalPlaces(inititalPrice * constants.CommodityOptionsSttOrCttSell)
		}

		// Transaction Charges
		transactionCharges = reduceTenDecimalPlaces(inititalPrice * constants.CommodityOptionsTransactionCharge)

		// Sebi Charges (check may contains error)
		sebiCharges = reduceTenDecimalPlaces(inititalPrice * constants.CommodityOptionsSebiCharges)

		// GST Charges
		gst = reduceTenDecimalPlaces((brokerage + transactionCharges) * constants.CommodityOptionsGst)

		// Stamp Charges
		stampCharges = reduceTenDecimalPlaces(inititalPrice * constants.CommodityOptionsStampProcessBuy)
	}

	totalPrice = reduceTenDecimalPlaces(brokerage + sttOrCtt + transactionCharges + gst + sebiCharges + stampCharges)

	var brokerChargesRes models.BrokerChargesRes
	brokerChargesRes.Price = inititalPrice + totalPrice
	brokerChargesRes.Brokerage = brokerage
	brokerChargesRes.SttOrCtt = sttOrCtt
	brokerChargesRes.TransactionCharges = transactionCharges
	brokerChargesRes.SebiCharges = sebiCharges
	brokerChargesRes.Gst = gst
	brokerChargesRes.StampCharges = stampCharges

	return brokerChargesRes
}

func reduceTenDecimalPlaces(val float64) float64 {
	return math.Floor(val*10000000000) / 10000000000
}

func ceilToTwoDecimalPlaces(num float64) float64 {
	result := math.Ceil(num*100) / 100
	return result
}

func (obj ChargesObj) FundsPayout(req models.FundsPayoutReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return FundsPayoutInternal(req, reqH)
}

var FundsPayoutInternal = func(req models.FundsPayoutReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	var fetchFundsReq models.FetchFundsRequest
	fetchFundsReq.ClientID = req.ClientID
	fetchFundsReq.Type = constants.FundsTypeAll

	statusFunds, resFunds := tradelab.FetchFundsInternal(fetchFundsReq, reqH)

	if statusFunds != http.StatusOK {
		loggerconfig.Error("FundsPayout FetchFunds status != 200", statusFunds, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	fetchFundsRes, ok := resFunds.Data.(models.FetchFundsResponse)
	if !ok {
		loggerconfig.Error("FundsPayout FetchFunds interface parsing error", ok, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var openingBalance float64
	var marginUsed float64
	var payIn float64
	var lossOnClosedPositions float64
	var pnlOnClosedPositions float64
	var cmPnlOnClosedPositions float64
	var totalPnlOnClosedPositions float64
	var equityCreditSell float64

	// Extract funds data (opening balance, margin used, pay-in)
	var chargesOnTrades float64 // only executed
	for i := 0; i < len(fetchFundsRes.Values); i++ {
		if fetchFundsRes.Values[i].Num0 == constants.OpeningBalance {
			openingBalance, _ = strconv.ParseFloat(fetchFundsRes.Values[i].Num1, 64)
		} else if fetchFundsRes.Values[i].Num0 == constants.MarginUsed {
			marginUsed, _ = strconv.ParseFloat(fetchFundsRes.Values[i].Num1, 64)
		} else if fetchFundsRes.Values[i].Num0 == constants.Payin {
			payIn, _ = strconv.ParseFloat(fetchFundsRes.Values[i].Num1, 64)
		} else if fetchFundsRes.Values[i].Num0 == constants.EquityCreditSell {
			equityCreditSell, _ = strconv.ParseFloat(fetchFundsRes.Values[i].Num1, 64)
		}
	}

	// Fetch position data (to calculate profit/loss)
	var getPositionsReq models.GetPositionRequest
	getPositionsReq.ClientID = req.ClientID
	getPositionsReq.Type = constants.Historical

	statusPosition, resPositions := tradelab.GetPositionsInternal(getPositionsReq, reqH)
	if statusPosition != http.StatusOK {
		loggerconfig.Error("FundsPayout GetPositions status != 200", statusPosition, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	getPositionsRes, ok := resPositions.Data.([]models.GetPositionResponseData)
	if !ok {
		loggerconfig.Error("FundsPayout GetPosition interface parsing error", ok, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	dpCharges := 0.0

	// Calculate profit and loss from positions
	for i := 0; i < len(getPositionsRes); i++ {
		if getPositionsRes[i].Exchange != strings.ToUpper(constants.NFO) && getPositionsRes[i].Exchange != strings.ToUpper(constants.MCX) && getPositionsRes[i].NetQuantity == 0 {
			cmPnlOnClosedPositions += getPositionsRes[i].NetAmount
		}
		if getPositionsRes[i].NetQuantity == 0 {
			pnlOnClosedPositions += getPositionsRes[i].NetAmount
			if getPositionsRes[i].NetAmount < 0 {
				lossOnClosedPositions += getPositionsRes[i].NetAmount
			}
		}
		if getPositionsRes[i].NetQuantity < 0 && strings.EqualFold(getPositionsRes[i].Product, constants.CNC) {
			dpCharges += constants.DpCharges
		}
	}

	// Calculate total profit/loss based on conditions
	if pnlOnClosedPositions < 0 { // if pnl is negative, then this will be the total pnl
		totalPnlOnClosedPositions = pnlOnClosedPositions
	} else { // if pnl is positive, then 95% of cmPnl(cash-market) will be the total pnl
		totalPnlOnClosedPositions = 0.95 * cmPnlOnClosedPositions
	}

	// Fetch completed orders for brokerage charges
	var completedOrderRequest models.CompletedOrderRequest
	completedOrderRequest.ClientID = req.ClientID
	completedOrderRequest.Type = strings.ToLower(constants.Completed)

	statusOrder, resOrder := tradelab.CompletedOrderInternal(completedOrderRequest, reqH)
	if statusOrder != http.StatusOK {
		loggerconfig.Error("FundsPayout CompletedOrder status != 200", statusOrder, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	completedOrderRes, ok := resOrder.Data.(models.CompletedOrderResponse)
	if !ok {
		loggerconfig.Error("FundsPayout CompletedOrder interface parsing error", ok, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	allBrokerageChargesReq, _ := createBrokerageChargesRequest(completedOrderRes, reqH)

	for i := 0; i < len(allBrokerageChargesReq); i++ {
		statusBrokerageCal, resBrokerageCal := BrokerChargesInternal(allBrokerageChargesReq[i], reqH)
		if statusBrokerageCal != http.StatusOK {
			loggerconfig.Error("FundsPayout BrokerCharges status != 200", statusBrokerageCal, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		brokerChargesRes, ok := resBrokerageCal.Data.(models.BrokerChargesRes)
		if !ok {
			loggerconfig.Error("FundsPayout BrokerCharges interface parsing error", ok, " uccId:", req.ClientID, " requestId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		chargesOnTrades += brokerChargesRes.Brokerage + brokerChargesRes.SttOrCtt + brokerChargesRes.TransactionCharges + brokerChargesRes.SebiCharges + brokerChargesRes.Gst + brokerChargesRes.StampCharges
	}
	chargesOnTrades += dpCharges

	// Calculate payout amounts
	calculatedPayoutAmount := openingBalance + payIn - marginUsed - math.Abs(lossOnClosedPositions) - chargesOnTrades

	//default payout amount is opening Balance
	payoutAmount := openingBalance
	//if calculted amount is less than opening balance then payout amount is calculated amount.
	if calculatedPayoutAmount < openingBalance {
		payoutAmount = calculatedPayoutAmount
	}

	// If the payoutAmount is negative, set it to 0.
	if payoutAmount < 0 {
		payoutAmount = 0
	}

	// calculation for extra payout amount
	calculatedExtraPayoutAmount := openingBalance + payIn - marginUsed + totalPnlOnClosedPositions + equityCreditSell - chargesOnTrades
	if calculatedExtraPayoutAmount < 0 {
		calculatedExtraPayoutAmount = 0
	}

	var fundsPayoutRes models.FundsPayoutRes
	fundsPayoutRes.ClientID = req.ClientID
	fundsPayoutRes.PayoutAmount = payoutAmount
	fundsPayoutRes.OpeningBalance = openingBalance
	fundsPayoutRes.MarginUsed = marginUsed
	fundsPayoutRes.LossOnClosedPositions = lossOnClosedPositions
	fundsPayoutRes.ChargesOnTrades = chargesOnTrades
	fundsPayoutRes.Payin = payIn
	fundsPayoutRes.ExtraPayoutAmount = calculatedExtraPayoutAmount
	fundsPayoutRes.ProfitNLoss = pnlOnClosedPositions
	fundsPayoutRes.CmPnl = cmPnlOnClosedPositions
	fundsPayoutRes.UserPnl = totalPnlOnClosedPositions
	fundsPayoutRes.EquityCreditSell = equityCreditSell

	loggerconfig.Info("FundsPayout resp=", helpers.LogStructAsJSON(fundsPayoutRes), " userId: ", req.ClientID, " requestId:", reqH.RequestId)

	apiRes.Data = fundsPayoutRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}

func createBrokerageChargesRequest(completedOrderRes models.CompletedOrderResponse, reqH models.ReqHeader) ([]models.BrokerChargesReq, error) {
	var allBrokerageChargesReq []models.BrokerChargesReq

	for i := 0; i < len(completedOrderRes.Orders); i++ {
		if !strings.EqualFold(completedOrderRes.Orders[i].OrderStatus, constants.Completed) && !strings.EqualFold(completedOrderRes.Orders[i].OrderStatus, constants.Complete) {
			continue
		}
		var brokerageChargeReq models.BrokerChargesReq
		brokerageChargeReq.ClientID = completedOrderRes.Orders[i].ClientID
		brokerageChargeReq.Price = completedOrderRes.Orders[i].AveragePrice
		brokerageChargeReq.Quantity = completedOrderRes.Orders[i].Quantity
		brokerageChargeReq.Process = completedOrderRes.Orders[i].OrderSide
		brokerageChargeReq.Exchange = completedOrderRes.Orders[i].Exchange

		segment := constants.EQUITY
		groupInfo := 0
		agriType := false

		switch completedOrderRes.Orders[i].Exchange {
		case strings.ToUpper(constants.CDS):
			segment = strings.ToLower(constants.CURRENCY)
		case strings.ToUpper(constants.NFO):
			segment = strings.ToLower(constants.EQUITY)
		case strings.ToUpper(constants.MCX):
			groupInfo = constants.CommodityMap[completedOrderRes.Orders[i].Segment]
			if groupInfo != 0 {
				agriType = true
			}
			segment = strings.ToLower(constants.COMMODITY)
		}

		brokerageChargeReq.Segment = segment
		brokerageChargeReq.GroupInfo = groupInfo
		brokerageChargeReq.Agri = agriType

		subSegment := ""
		if completedOrderRes.Orders[i].Exchange == strings.ToUpper(constants.NSE) || completedOrderRes.Orders[i].Exchange == strings.ToUpper(constants.BSE) {
			switch completedOrderRes.Orders[i].Product {
			case strings.ToUpper(constants.CNC):
				subSegment = strings.ToLower(constants.DELIVERY)
			case strings.ToUpper(constants.MIS):
				subSegment = strings.ToLower(constants.INTRADAY)
			}
		} else {
			if strings.HasSuffix(completedOrderRes.Orders[i].TradingSymbol, "FUT") {
				subSegment = strings.ToLower(constants.FUTURES)
			} else {
				subSegment = strings.ToLower(constants.OPTIONS)
			}
		}
		brokerageChargeReq.SubSegment = subSegment

		allBrokerageChargesReq = append(allBrokerageChargesReq, brokerageChargeReq)
	}

	return allBrokerageChargesReq, nil

}
