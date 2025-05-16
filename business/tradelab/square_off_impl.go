package tradelab

import (
	"net/http"
	"strconv"
	"strings"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
)

type SquareOffObj struct {
	tradeLabURL string
}

func InitSquareOffProvider() SquareOffObj {
	defer models.HandlePanic()

	squareOffObj := SquareOffObj{
		tradeLabURL: constants.TLURL,
	}

	return squareOffObj
}

func (obj SquareOffObj) SquareOffAll(req models.SquareOffAllReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var getPositionReq models.GetPositionRequest
	getPositionReq.ClientID = req.ClientID
	getPositionReq.Type = constants.LIVE
	var portfolioObj PortfolioObj
	portfolioObj.tradeLabURL = obj.tradeLabURL
	status, res := PortfolioObj.GetPositions(portfolioObj, getPositionReq, reqH)

	if status != http.StatusOK {
		loggerconfig.Error("SquareOffAll in GetPositions status != 200", status, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	getPositionResponseData, ok := res.Data.([]models.GetPositionResponseData)
	if !ok {
		loggerconfig.Error("SquareOffAll interface parsing error", ok, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}

	var orderObj OrderObj
	orderObj.tradeLabURL = obj.tradeLabURL

	var squareOffAllRes models.SquareOffAllRes

	for i := 0; i < len(getPositionResponseData); i++ {
		var placeOrderReq models.PlaceOrderRequest
		placeOrderReq.Exchange = getPositionResponseData[i].Exchange
		placeOrderReq.InstrumentToken = strconv.Itoa(getPositionResponseData[i].InstrumentToken)
		placeOrderReq.ClientID = getPositionResponseData[i].ClientID
		placeOrderReq.OrderType = constants.MARKET
		placeOrderReq.Price = getPositionResponseData[i].AverageSellPrice
		if getPositionResponseData[i].NetQuantity < 0 {
			placeOrderReq.Quantity = -1 * getPositionResponseData[i].NetQuantity
			placeOrderReq.DisclosedQuantity = 0
			placeOrderReq.OrderSide = strings.ToUpper(constants.BUY)
		} else {
			placeOrderReq.Quantity = getPositionResponseData[i].NetQuantity
			placeOrderReq.DisclosedQuantity = getPositionResponseData[i].NetQuantity
			placeOrderReq.OrderSide = strings.ToUpper(constants.SELL)
		}
		// placeOrderReq.DisclosedQuantity = getPositionResponseData[i].DisclosedQuantity
		placeOrderReq.Validity = constants.DAY
		if getPositionResponseData[i].Product == constants.CNC {
			placeOrderReq.Product = constants.CNC
		} else {
			placeOrderReq.Product = constants.MIS // getPositionResponseData[i].Product // MIS
		}

		placeOrderReq.TriggerPrice = 0
		placeOrderReq.ExecutionType = constants.REGULAR

		statusPlaceOrder, resPlaceOrder := OrderObj.PlaceOrder(orderObj, placeOrderReq, reqH)

		if statusPlaceOrder != http.StatusOK {
			loggerconfig.Error("SquareOffAll in PlaceOrder status != 200", statusPlaceOrder, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
			return apihelpers.SendInternalServerError()
		}

		placeOrderResponse, ok := resPlaceOrder.Data.(models.PlaceOrderResponse)
		if !ok {
			loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " SquareOffAll interface parsing error", ok, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
			return apihelpers.SendInternalServerError()
		}
		squareOffAllRes.SquareOffAll = append(squareOffAllRes.SquareOffAll, placeOrderResponse)
	}

	var apiRes apihelpers.APIRes
	loggerconfig.Info("squareOffAllRes tl resp=", helpers.LogStructAsJSON(squareOffAllRes), " uccId:", req.ClientID, " StatusCode: ", status, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
	apiRes.Data = squareOffAllRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes

}
