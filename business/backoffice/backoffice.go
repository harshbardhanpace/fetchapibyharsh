package backoffice

import (
	"bytes"
	"encoding/json"

	"io"
	"math"
	"net/http"
	"sort"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
	"strconv"
	"strings"
	"time"

	"github.com/ggwhite/go-masker"
	"github.com/sirupsen/logrus"
)

type BackofficeObj struct {
	shilpiBaseUrl string
}

var maskObj *masker.Masker

func InitBackofficeObj() BackofficeObj {
	defer models.HandlePanic()

	backofficeObj := BackofficeObj{
		shilpiBaseUrl: constants.ShilpiURL,
	}

	maskObj = masker.New()

	return backofficeObj
}

func (obj BackofficeObj) TradeConfirmationDateRange(tcDateRangeReq models.TradeConfirmationDateRangeReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	url := obj.shilpiBaseUrl + "?Requesttype=" + constants.TradeConfirmationDateRange + "&userid=" + tcDateRangeReq.UserID + "&dfdatefr=" + tcDateRangeReq.DFDateFr + "&dfdateto=" + tcDateRangeReq.DFDateTo

	// emtpy payload
	payload := new(bytes.Buffer)

	res, err := apihelpers.CallApiShilpi(http.MethodGet, url, payload)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " TradeConfirmationDateRange call api error =", err, " uccId:", tcDateRangeReq.UserID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("TradeConfirmationDateRange error in reading response body=", err, " uccId:", tcDateRangeReq.UserID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	shilpiTradeConfirmationDateRangeRes := ShilpiTradeConfirmationDateRangeRes{}
	json.Unmarshal([]byte(string(body)), &shilpiTradeConfirmationDateRangeRes)

	var tradeConfirmationDateRangeRes models.TradeConfirmationDateRangeRes

	for i := 0; i < len(shilpiTradeConfirmationDateRangeRes); i++ {
		var tradeConfirmationDateRange models.TradeConfirmationDateRange
		tradeConfirmationDateRange.SellAvgRate = shilpiTradeConfirmationDateRangeRes[i].SellAvgRate
		tradeConfirmationDateRange.Scripname = shilpiTradeConfirmationDateRangeRes[i].Scripname
		tradeConfirmationDateRange.BuyValue = shilpiTradeConfirmationDateRangeRes[i].BuyValue
		tradeConfirmationDateRange.SellValue = shilpiTradeConfirmationDateRangeRes[i].SellValue
		tradeConfirmationDateRange.SellQty = shilpiTradeConfirmationDateRangeRes[i].SellQty
		tradeConfirmationDateRange.NetQty = shilpiTradeConfirmationDateRangeRes[i].NetQty
		tradeConfirmationDateRange.Segment = shilpiTradeConfirmationDateRangeRes[i].Segment
		tradeConfirmationDateRange.NetValue = shilpiTradeConfirmationDateRangeRes[i].NetValue
		tradeConfirmationDateRange.Buyqty = shilpiTradeConfirmationDateRangeRes[i].Buyqty
		tradeConfirmationDateRange.TradeDate = shilpiTradeConfirmationDateRangeRes[i].TradeDate
		tradeConfirmationDateRange.NetAvgRate = shilpiTradeConfirmationDateRangeRes[i].NetAvgRate

		tradeConfirmationDateRangeRes.AllTradeConfirmationDateRange = append(tradeConfirmationDateRangeRes.AllTradeConfirmationDateRange, tradeConfirmationDateRange)
	}

	loggerconfig.Info("TradeConfirmationDateRange  resp=", helpers.LogStructAsJSON(tradeConfirmationDateRangeRes), " uccId:", tcDateRangeReq.UserID, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)

	var apiRes apihelpers.APIRes
	apiRes.Data = tradeConfirmationDateRangeRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func GetBillDetailsCdslData(getBillDetailsCdslReq models.GetBillDetailsCdslReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	var fetchProfileReq models.FetchProfileReq
	fetchProfileReq.UserID = getBillDetailsCdslReq.UserID

	shilpiObj := InitBackofficeObj()
	status, resFetchProfile := shilpiObj.FetchProfile(fetchProfileReq, reqH)

	if status != http.StatusOK {
		loggerconfig.Error("GetBillDetailsCdsl FetchProfile status != 200", status, " uccId:", getBillDetailsCdslReq.UserID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	fetchProfileRes, ok := resFetchProfile.Data.(models.FetchProfileRes)
	if !ok {
		loggerconfig.Error("GetBillDetailsCdsl FetchProfile interface parsing error", ok, " uccId:", getBillDetailsCdslReq.UserID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	if len(fetchProfileRes.FetchProfile) == 0 {
		loggerconfig.Error("GetBillDetailsCdsl FetchProfile interface parsing error", ok, " uccId:", getBillDetailsCdslReq.UserID, " requestId:", reqH.RequestId)
		return apihelpers.SendErrorResponse(false, constants.EmptyProfileResponse, http.StatusOK)
	}

	url := shilpiObj.shilpiBaseUrl + "?Requesttype=" + constants.GetBillDetailsCdsl + "&userid=" + fetchProfileRes.FetchProfile[0].Dpaccountno + "&dfdatefr=" + getBillDetailsCdslReq.DFDateFr + "&dfdateto=" + getBillDetailsCdslReq.DFDateTo

	// emtpy payload
	payload := new(bytes.Buffer)

	res, err := apihelpers.CallApiShilpi(http.MethodGet, url, payload)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " GetBillDetailsCdsl call api error =", err, " uccId:", getBillDetailsCdslReq.UserID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("GetBillDetailsCdsl error in reading response body=", err, " uccId:", getBillDetailsCdslReq.UserID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	shilpiGetBillDetailsCdsl := ShilpiGetBillDetailsCdsl{}
	json.Unmarshal([]byte(string(body)), &shilpiGetBillDetailsCdsl)

	var getBillDetailsCdslRes models.GetBillDetailsCdslRes

	for i := 0; i < len(shilpiGetBillDetailsCdsl); i++ {
		var getBillDetailsCdsl models.GetBillDetailsCdsl
		getBillDetailsCdsl.Charges = shilpiGetBillDetailsCdsl[i].Charges
		getBillDetailsCdsl.TrxDate = shilpiGetBillDetailsCdsl[i].TrxDate
		getBillDetailsCdsl.Qty = shilpiGetBillDetailsCdsl[i].Qty
		getBillDetailsCdsl.Gst = shilpiGetBillDetailsCdsl[i].Gst
		getBillDetailsCdsl.InstrumentName = shilpiGetBillDetailsCdsl[i].InstrumentName
		getBillDetailsCdsl.Isincode = shilpiGetBillDetailsCdsl[i].Isincode
		getBillDetailsCdsl.ChargesDetails = shilpiGetBillDetailsCdsl[i].ChargesDetails
		getBillDetailsCdsl.TotalCharges = shilpiGetBillDetailsCdsl[i].TotalCharges

		getBillDetailsCdslRes.GetBillDetailsCdsl = append(getBillDetailsCdslRes.GetBillDetailsCdsl, getBillDetailsCdsl)
	}

	loggerconfig.Info("GetBillDetailsCdsl  resp=", getBillDetailsCdslRes, " uccId:", getBillDetailsCdslReq.UserID, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)

	var apiRes apihelpers.APIRes
	apiRes.Data = getBillDetailsCdslRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj BackofficeObj) GetBillDetailsCdsl(getBillDetailsCdslReq models.GetBillDetailsCdslReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	return GetBillDetailsCdslData(getBillDetailsCdslReq, reqH)
}

func (obj BackofficeObj) LongTermShortTerm(longTermShortTermReq models.LongTermShortTermReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.shilpiBaseUrl + "?Requesttype=" + constants.GetLongTermShortTerm + "&userid=" + longTermShortTermReq.UserID

	// emtpy payload
	payload := new(bytes.Buffer)

	res, err := apihelpers.CallApiShilpi(http.MethodGet, url, payload)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " LongTermShortTerm call api error =", err, " uccId:", longTermShortTermReq.UserID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("LongTermShortTerm error in reading response body=", err, " uccId:", longTermShortTermReq.UserID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	shilpiLongTermShortTerm := ShilpiLongTermShortTerm{}
	json.Unmarshal([]byte(string(body)), &shilpiLongTermShortTerm)

	var longTermShortTermRes models.LongTermShortTermRes
	for i := 0; i < len(shilpiLongTermShortTerm); i++ {
		var longTermShortTerm models.LongTermShortTerm
		longTermShortTerm.Scripname = shilpiLongTermShortTerm[i].Scripname
		longTermShortTerm.BuyRate = shilpiLongTermShortTerm[i].BuyRate
		longTermShortTerm.SellQty = shilpiLongTermShortTerm[i].SellQty
		longTermShortTerm.Jobbing = shilpiLongTermShortTerm[i].Jobbing
		longTermShortTerm.SellRate = shilpiLongTermShortTerm[i].SellRate
		longTermShortTerm.ShortTerm = shilpiLongTermShortTerm[i].ShortTerm
		longTermShortTerm.Buyqty = shilpiLongTermShortTerm[i].Buyqty
		longTermShortTerm.LongTerm = shilpiLongTermShortTerm[i].LongTerm
		longTermShortTerm.Scripcd = shilpiLongTermShortTerm[i].Scripcd
		longTermShortTerm.Isin = shilpiLongTermShortTerm[i].Isin

		longTermShortTermRes.LongTermShortTerm = append(longTermShortTermRes.LongTermShortTerm, longTermShortTerm)
	}

	loggerconfig.Info("LongTermShortTerm  resp=", longTermShortTermRes, " uccId:", longTermShortTermReq.UserID, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)

	var apiRes apihelpers.APIRes
	apiRes.Data = longTermShortTermRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}

func (obj BackofficeObj) FetchProfile(fetchProfileReq models.FetchProfileReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.shilpiBaseUrl + "?Requesttype=" + constants.Profile + "&userid=" + strings.ToUpper(fetchProfileReq.UserID)

	// emtpy payload
	payload := new(bytes.Buffer)

	res, err := apihelpers.CallApiShilpi(http.MethodGet, url, payload)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " FetchProfile call api error =", err, " uccId:", fetchProfileReq.UserID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("FetchProfile error in reading response body=", err, " uccId:", fetchProfileReq.UserID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	shilpiFetchProfile := ShilpiFetchProfile{}
	json.Unmarshal([]byte(string(body)), &shilpiFetchProfile)

	var fetchProfileRes models.FetchProfileRes
	for i := 0; i < len(shilpiFetchProfile); i++ {
		var fetchProfile models.FetchProfile
		fetchProfile.Emailno = shilpiFetchProfile[i].Emailno
		fetchProfile.Phnos = shilpiFetchProfile[i].Phnos
		fetchProfile.Occupationcode = shilpiFetchProfile[i].Occupationcode
		fetchProfile.Poaflag = shilpiFetchProfile[i].Poaflag
		fetchProfile.Gender = shilpiFetchProfile[i].Gender
		fetchProfile.City = shilpiFetchProfile[i].City
		fetchProfile.Annualincomedate = shilpiFetchProfile[i].Annualincomedate
		fetchProfile.Panno = shilpiFetchProfile[i].Panno
		fetchProfile.Annualincomecode = shilpiFetchProfile[i].Annualincomecode
		fetchProfile.Dpid = shilpiFetchProfile[i].Dpid
		fetchProfile.Mobileno = shilpiFetchProfile[i].Mobileno
		fetchProfile.Micrcode = shilpiFetchProfile[i].Micrcode
		fetchProfile.Branch = shilpiFetchProfile[i].Branch
		fetchProfile.Ckycregnno = shilpiFetchProfile[i].Ckycregnno
		fetchProfile.Bankcity = shilpiFetchProfile[i].Bankcity
		fetchProfile.Groupclient = shilpiFetchProfile[i].Groupclient
		fetchProfile.Accountstatus = shilpiFetchProfile[i].Accountstatus
		fetchProfile.Bankadd3 = shilpiFetchProfile[i].Bankadd3
		fetchProfile.Bankadd1 = shilpiFetchProfile[i].Bankadd1
		fetchProfile.Bankadd2 = shilpiFetchProfile[i].Bankadd2
		fetchProfile.State = shilpiFetchProfile[i].State
		fetchProfile.Add1 = shilpiFetchProfile[i].Add1
		fetchProfile.Activeexchange = shilpiFetchProfile[i].Activeexchange
		fetchProfile.Add4 = shilpiFetchProfile[i].Add4
		fetchProfile.Branchname = shilpiFetchProfile[i].Branchname
		fetchProfile.Enableotradeing = shilpiFetchProfile[i].Enableotradeing
		fetchProfile.Pincode = shilpiFetchProfile[i].Pincode
		fetchProfile.Introdate = shilpiFetchProfile[i].Introdate
		fetchProfile.Motherfullname = shilpiFetchProfile[i].Motherfullname
		fetchProfile.Bankacno = shilpiFetchProfile[i].Bankacno
		fetchProfile.Ifsccode = shilpiFetchProfile[i].Ifsccode
		fetchProfile.Bankactype = shilpiFetchProfile[i].Bankactype
		fetchProfile.Dpaccountno = shilpiFetchProfile[i].Dpaccountno
		fetchProfile.Fathername = shilpiFetchProfile[i].Fathername
		fetchProfile.Nationality = shilpiFetchProfile[i].Nationality
		fetchProfile.Dob = shilpiFetchProfile[i].Dob
		fetchProfile.Name = shilpiFetchProfile[i].Name
		fetchProfile.Bankname = shilpiFetchProfile[i].Bankname
		fetchProfile.Kraupload = shilpiFetchProfile[i].Kraupload
		fetchProfile.Polexposueperson = shilpiFetchProfile[i].Polexposueperson
		fetchProfile.Fatca = shilpiFetchProfile[i].Fatca
		fetchProfile.Married = shilpiFetchProfile[i].Married

		fetchProfileRes.FetchProfile = append(fetchProfileRes.FetchProfile, fetchProfile)
	}

	maskedFetchProfileRes, err := maskObj.Struct(fetchProfileRes)
	if err != nil {
		loggerconfig.Error("FetchProfile Error in masking request err: ", err, " clientId: ", fetchProfileReq.UserID, " requestid = ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("FetchProfile  resp=", maskedFetchProfileRes, " uccId:", fetchProfileReq.UserID, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)

	var apiRes apihelpers.APIRes
	apiRes.Data = fetchProfileRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}

func (obj BackofficeObj) TradeConfirmationOnDate(tcOnDateReq models.TradeConfirmationOnDateReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.shilpiBaseUrl + "?Requesttype=" + constants.TradeConfirmation + "&userid=" + tcOnDateReq.UserId + "&tradedate=" + tcOnDateReq.TradeDate

	// emtpy payload
	payload := new(bytes.Buffer)

	res, err := apihelpers.CallApiShilpi(http.MethodGet, url, payload)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " TradeConfirmationOnDate call api error =", err, " uccId:", tcOnDateReq.UserId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("TradeConfirmationOnDate error in reading response body=", err, " uccId:", tcOnDateReq.UserId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	shilpiTradeConfirmationOnDate := ShilpiTradeConfirmationOnDate{}
	json.Unmarshal([]byte(string(body)), &shilpiTradeConfirmationOnDate)

	var tradeConfirmationOnDateRes models.TradeConfirmationOnDateRes
	for i := 0; i < len(shilpiTradeConfirmationOnDate); i++ {
		var tradeConfirmationOnDate models.TradeConfirmationOnDate
		tradeConfirmationOnDate.SellAvgRate = shilpiTradeConfirmationOnDate[i].SellAvgRate
		tradeConfirmationOnDate.Scripname = shilpiTradeConfirmationOnDate[i].Scripname
		tradeConfirmationOnDate.BuyValue = shilpiTradeConfirmationOnDate[i].BuyValue
		tradeConfirmationOnDate.SellValue = shilpiTradeConfirmationOnDate[i].SellValue
		tradeConfirmationOnDate.SellQty = shilpiTradeConfirmationOnDate[i].SellQty
		tradeConfirmationOnDate.NetQty = shilpiTradeConfirmationOnDate[i].NetQty
		tradeConfirmationOnDate.Segment = shilpiTradeConfirmationOnDate[i].Segment
		tradeConfirmationOnDate.NetValue = shilpiTradeConfirmationOnDate[i].NetValue
		tradeConfirmationOnDate.Buyqty = shilpiTradeConfirmationOnDate[i].Buyqty
		tradeConfirmationOnDate.BuyAvgRate = shilpiTradeConfirmationOnDate[i].BuyAvgRate
		tradeConfirmationOnDate.NetAvgRate = shilpiTradeConfirmationOnDate[i].NetAvgRate

		tradeConfirmationOnDateRes.TradeConfirmationOnDate = append(tradeConfirmationOnDateRes.TradeConfirmationOnDate, tradeConfirmationOnDate)
	}

	loggerconfig.Info("TradeConfirmationOnDate  resp=", tradeConfirmationOnDateRes, " uccId:", tcOnDateReq.UserId, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)

	var apiRes apihelpers.APIRes
	apiRes.Data = tradeConfirmationOnDateRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}

func (obj BackofficeObj) OpenPositions(openPositionsReq models.OpenPositionsReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.shilpiBaseUrl + "?Requesttype=" + constants.OpenPosition + "&userid=" + openPositionsReq.UserId + "&dsflag=" + openPositionsReq.DsFlag

	// emtpy payload
	payload := new(bytes.Buffer)

	res, err := apihelpers.CallApiShilpi(http.MethodGet, url, payload)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " OpenPositions call api error =", err, " uccId:", openPositionsReq.UserId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("OpenPositions error in reading response body=", err, " uccId:", openPositionsReq.UserId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	shilpiOpenPositions := ShilpiOpenPositions{}
	json.Unmarshal([]byte(string(body)), &shilpiOpenPositions)

	var openPositionsRes models.OpenPositionsRes
	for i := 0; i < len(shilpiOpenPositions); i++ {
		var openPositions models.OpenPositions
		openPositions.Segment = shilpiOpenPositions[i].Num1
		openPositions.Exchange = shilpiOpenPositions[i].Num2
		openPositions.Contract = shilpiOpenPositions[i].Num3
		openPositions.BuyQty = shilpiOpenPositions[i].Num4
		openPositions.BuyValue = shilpiOpenPositions[i].Num5
		openPositions.SellQty = shilpiOpenPositions[i].Num6
		openPositions.SellValue = shilpiOpenPositions[i].Num7
		openPositions.NetQty = shilpiOpenPositions[i].Num8
		openPositions.NetValue = shilpiOpenPositions[i].Num9
		openPositions.ClosingRate = shilpiOpenPositions[i].Num10
		openPositions.PAndL = shilpiOpenPositions[i].Num11
		openPositions.Exposure = shilpiOpenPositions[i].Num12

		openPositionsRes.OpenPositions = append(openPositionsRes.OpenPositions, openPositions)
	}

	loggerconfig.Info("OpenPositions  resp=", openPositionsRes, " uccId:", openPositionsReq.UserId, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)

	var apiRes apihelpers.APIRes
	apiRes.Data = openPositionsRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}

func (obj BackofficeObj) GetHolding(getHoldingReq models.GetHoldingReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.shilpiBaseUrl + "?Requesttype=" + constants.GetHolding + "&userid=" + getHoldingReq.UserId

	// emtpy payload
	payload := new(bytes.Buffer)

	res, err := apihelpers.CallApiShilpi(http.MethodGet, url, payload)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " GetHolding call api error =", err, " uccId:", getHoldingReq.UserId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("GetHolding error in reading response body=", err, " uccId:", getHoldingReq.UserId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	shilpiGetHolding := ShilpiGetHolding{}
	json.Unmarshal([]byte(string(body)), &shilpiGetHolding)

	var getHoldingRes models.GetHoldingRes
	for i := 0; i < len(shilpiGetHolding); i++ {
		var getHolding models.GetHolding
		getHolding.MobileNo = shilpiGetHolding[i].Num1
		getHolding.LoginId = shilpiGetHolding[i].Num2
		getHolding.IsinCode = shilpiGetHolding[i].Num3
		getHolding.IsinName = shilpiGetHolding[i].Num4
		getHolding.Holding = shilpiGetHolding[i].Num5
		getHolding.ClosePrice = shilpiGetHolding[i].Num6
		getHolding.Valuation = shilpiGetHolding[i].Num7
		getHolding.ScripCode = shilpiGetHolding[i].Num8

		getHoldingRes.GetHolding = append(getHoldingRes.GetHolding, getHolding)
	}

	maskedGetHoldingRes, err := maskObj.Struct(getHoldingRes)
	if err != nil {
		loggerconfig.Error("In Controller GetHolding Error in masking request err: ", err, " clientId: ", getHoldingReq.UserId, " requestid = ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("GetHolding  resp=", helpers.LogStructAsJSON(maskedGetHoldingRes), " uccId:", getHoldingReq.UserId, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)

	var apiRes apihelpers.APIRes
	apiRes.Data = getHoldingRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj BackofficeObj) GetMarginOnDate(getMarginOnDateReq models.GetMarginOnDateReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.shilpiBaseUrl + "?Requesttype=" + constants.GetMarginOnDate + "&userid=" + getMarginOnDateReq.UserId + "&margindate=" + getMarginOnDateReq.MarginDate

	// emtpy payload
	payload := new(bytes.Buffer)

	res, err := apihelpers.CallApiShilpi(http.MethodGet, url, payload)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " GetMarginOnDate call api error =", err, " uccId:", getMarginOnDateReq.UserId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("GetMarginOnDate error in reading response body=", err, " uccId:", getMarginOnDateReq.UserId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	shilpiGetMarginOnDate := ShilpiGetMarginOnDate{}
	json.Unmarshal([]byte(string(body)), &shilpiGetMarginOnDate)

	var getMarginOnDateRes models.GetMarginOnDateRes
	for i := 0; i < len(shilpiGetMarginOnDate); i++ {
		var getMarginOnDate models.GetMarginOnDate
		getMarginOnDate.TotalMargin = shilpiGetMarginOnDate[i].TotalMargin
		getMarginOnDate.Exposure = shilpiGetMarginOnDate[i].Exposure
		getMarginOnDate.PeakDeposit = shilpiGetMarginOnDate[i].PeakDeposit
		getMarginOnDate.VarMargin = shilpiGetMarginOnDate[i].VarMargin
		getMarginOnDate.Deposit = shilpiGetMarginOnDate[i].Deposit
		getMarginOnDate.ShortMargin = shilpiGetMarginOnDate[i].ShortMargin
		getMarginOnDate.TotalPeakmargin = shilpiGetMarginOnDate[i].TotalPeakmargin
		getMarginOnDate.PeakMarginShort = shilpiGetMarginOnDate[i].PeakMarginShort
		getMarginOnDate.Span = shilpiGetMarginOnDate[i].Span

		getMarginOnDateRes.GetMarginOnDate = append(getMarginOnDateRes.GetMarginOnDate, getMarginOnDate)
	}

	loggerconfig.Info("GetMarginOnDate  resp=", getMarginOnDateRes, " uccId:", getMarginOnDateReq.UserId, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)

	var apiRes apihelpers.APIRes
	apiRes.Data = getMarginOnDateRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj BackofficeObj) FinancialLedgerBalanceOnDate(financialLedgerBalanceOnDateReq models.FinancialLedgerBalanceOnDateReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.shilpiBaseUrl + "?Requesttype=" + constants.FinancialLedgerBalanceOnDate + "&userid=" + financialLedgerBalanceOnDateReq.UserId + "&asondate=" + financialLedgerBalanceOnDateReq.AsOnDate

	// emtpy payload
	payload := new(bytes.Buffer)

	res, err := apihelpers.CallApiShilpi(http.MethodGet, url, payload)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " FinancialLedgerBalanceOnDate call api error =", err, " uccId:", financialLedgerBalanceOnDateReq.UserId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("FinancialLedgerBalanceOnDate error in reading response body=", err, " uccId:", financialLedgerBalanceOnDateReq.UserId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	shilpiFinancialLedgerBalanceOnDate := ShilpiFinancialLedgerBalanceOnDate{}
	json.Unmarshal([]byte(string(body)), &shilpiFinancialLedgerBalanceOnDate)

	var financialLedgerBalanceOnDateRes models.FinancialLedgerBalanceOnDateRes
	for i := 0; i < len(shilpiFinancialLedgerBalanceOnDate); i++ {
		var financialLedgerBalanceOnDate models.FinancialLedgerBalanceOnDate
		financialLedgerBalanceOnDate.FinBalance = shilpiFinancialLedgerBalanceOnDate[i].Num1

		financialLedgerBalanceOnDateRes.FinancialLedgerBalanceOnDate = append(financialLedgerBalanceOnDateRes.FinancialLedgerBalanceOnDate, financialLedgerBalanceOnDate)
	}

	loggerconfig.Info("FinancialLedgerBalanceOnDate  resp=", financialLedgerBalanceOnDateRes, " uccId:", financialLedgerBalanceOnDateReq.UserId, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)

	var apiRes apihelpers.APIRes
	apiRes.Data = financialLedgerBalanceOnDateRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj BackofficeObj) GetFinancial(getFinancialReq models.GetFinancialReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.shilpiBaseUrl + "?Requesttype=" + constants.GetFinancial + "&noofentries=" + getFinancialReq.NoOfEntries + "&userid=" + getFinancialReq.UserId

	// emtpy payload
	payload := new(bytes.Buffer)

	res, err := apihelpers.CallApiShilpi(http.MethodGet, url, payload)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " Alert Severity:P0-Critical, platform:", reqH.Platform, " GetFinancial call api error =", err, " uccId:", getFinancialReq.UserId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("GetFinancial error in reading response body=", err, " uccId:", getFinancialReq.UserId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	shilpiGetFinancial := ShilpiGetFinancial{}
	json.Unmarshal([]byte(string(body)), &shilpiGetFinancial)

	var getFinancialRes models.GetFinancialRes
	for i := 0; i < len(shilpiGetFinancial); i++ {
		var getFinancial models.GetFinancial
		getFinancial.Date = shilpiGetFinancial[i].Date
		getFinancial.Narr = shilpiGetFinancial[i].Narr
		getFinancial.Segment = shilpiGetFinancial[i].Segment
		getFinancial.Exchange = shilpiGetFinancial[i].Exchange
		getFinancial.Debit = shilpiGetFinancial[i].Debit
		getFinancial.Credit = shilpiGetFinancial[i].Credit
		getFinancial.Payrefno = shilpiGetFinancial[i].Payrefno
		getFinancial.ValueDate = shilpiGetFinancial[i].ValueDate

		getFinancialRes.GetFinancial = append(getFinancialRes.GetFinancial, getFinancial)
	}

	loggerconfig.Info("GetFinancial  resp=", getFinancialRes, " uccId:", getFinancialReq.UserId, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)

	var apiRes apihelpers.APIRes
	apiRes.Data = getFinancialRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}

// Implement the sort.Interface for ScripWiseCosting based on TrxDate
type ByTrxDate []models.ScripWiseCosting

func (a ByTrxDate) Len() int           { return len(a) }
func (a ByTrxDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTrxDate) Less(i, j int) bool { return a[i].UnixTimeFormat < a[j].UnixTimeFormat }

func (obj BackofficeObj) GetScripWiseCostingData(getScripWiseCosting models.TradebookReq, reqH models.ReqHeader) (models.ScripWiseCostingRes, error) {

	var scripWiseCostingRes models.ScripWiseCostingRes
	shilpiObj := InitBackofficeObj()

	url := shilpiObj.shilpiBaseUrl + "?Requesttype=" + constants.ScripWiseCosting + "&userid=" + strings.ToUpper(getScripWiseCosting.UserID) + "&datefr=" + getScripWiseCosting.DFDateFr + "&dateto=" + getScripWiseCosting.DFDateTo + "&dfreport=D"

	// emtpy payload
	payload := new(bytes.Buffer)

	res, err := apihelpers.CallApiShilpi(http.MethodGet, url, payload)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " GetScripWiseCostingData call api error =", err, " uccId:", getScripWiseCosting.UserID, " requestId:", reqH.RequestId)
		return scripWiseCostingRes, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("GetScripWiseCostingData error in reading response body=", err, " uccId:", getScripWiseCosting.UserID, " requestId:", reqH.RequestId)
		return scripWiseCostingRes, err
	}

	shilpiScripWiseCosting := ShilpiScripWiseCosting{}
	json.Unmarshal([]byte(string(body)), &shilpiScripWiseCosting)

	for i := 0; i < len(shilpiScripWiseCosting); i++ {
		var scripWiseCosting models.ScripWiseCosting
		var otherCharges1, otherCharges2 float64
		scripWiseCosting.NetValue, err = strconv.ParseFloat(shilpiScripWiseCosting[i].NetValue, 64)
		scripWiseCosting.Brokerage, err = strconv.ParseFloat(shilpiScripWiseCosting[i].Brok, 64)
		scripWiseCosting.NetQty, err = strconv.ParseFloat(shilpiScripWiseCosting[i].NetQty, 64)
		scripWiseCosting.OrderNo = shilpiScripWiseCosting[i].OrderNo
		otherCharges1, err = strconv.ParseFloat(shilpiScripWiseCosting[i].OthChrg1, 64)
		otherCharges2, err = strconv.ParseFloat(shilpiScripWiseCosting[i].OthChrg2, 64)
		scripWiseCosting.OtherCharges = math.Floor((otherCharges1+otherCharges2)*100) / 100
		scripWiseCosting.ScripCode = shilpiScripWiseCosting[i].ScripCode
		scripWiseCosting.GST, err = strconv.ParseFloat(shilpiScripWiseCosting[i].GST, 64)
		scripWiseCosting.Stamp, err = strconv.ParseFloat(shilpiScripWiseCosting[i].Stamp, 64)
		scripWiseCosting.Price, err = strconv.ParseFloat(shilpiScripWiseCosting[i].AvgRate, 64)
		scripWiseCosting.TrxDate = shilpiScripWiseCosting[i].TrxDate
		scripWiseCosting.ISINCode = shilpiScripWiseCosting[i].ISINCode
		scripWiseCosting.SEBIFee, err = strconv.ParseFloat(shilpiScripWiseCosting[i].SEBIFee, 64)
		scripWiseCosting.NetAmt, err = strconv.ParseFloat(shilpiScripWiseCosting[i].NetAmt, 64)
		scripWiseCosting.SellValue, err = strconv.ParseFloat(shilpiScripWiseCosting[i].SellValue, 64)
		scripWiseCosting.STT, err = strconv.ParseFloat(shilpiScripWiseCosting[i].STT, 64)
		scripWiseCosting.BuyValue, err = strconv.ParseFloat(shilpiScripWiseCosting[i].BuyValue, 64)
		scripWiseCosting.ExchClg, err = strconv.ParseFloat(shilpiScripWiseCosting[i].ExchClg, 64)
		scripWiseCosting.TurnTax, err = strconv.ParseFloat(shilpiScripWiseCosting[i].TurnTax, 64)
		scripWiseCosting.Exchange = shilpiScripWiseCosting[i].Exchange
		scripWiseCosting.BrokType = shilpiScripWiseCosting[i].BrokType
		scripWiseCosting.ScripName = shilpiScripWiseCosting[i].ScripName
		if shilpiScripWiseCosting[i].BuyQty != "0" {
			scripWiseCosting.BuySellType = "Buy"
			scripWiseCosting.Quantity, err = strconv.ParseFloat(shilpiScripWiseCosting[i].BuyQty, 64)
		} else {
			scripWiseCosting.BuySellType = "Sell"
			scripWiseCosting.Quantity, err = strconv.ParseFloat(shilpiScripWiseCosting[i].SellQty, 64)
		}
		if err != nil {
			loggerconfig.Error("GetScripWiseCostingData, error in parsing string to float64 Error : ", err)
			return scripWiseCostingRes, err
		}

		// Parse the input date string using the specified layout
		parsedTime, err := time.Parse(constants.ShilpiDateFormat, scripWiseCosting.TrxDate)
		if err != nil {
			loggerconfig.Error("GetScripWiseCostingData, Error parsing shilpi date for tradebook :", err)
			return scripWiseCostingRes, err
		}

		// Convert the parsed Date to Unix time
		scripWiseCosting.UnixTimeFormat = parsedTime.Unix()

		scripWiseCostingRes.TotalBrokerage += scripWiseCosting.Brokerage
		scripWiseCostingRes.TotalGST += scripWiseCosting.GST
		scripWiseCostingRes.TotalSEBITax += scripWiseCosting.SEBIFee
		scripWiseCostingRes.TotalSTT += scripWiseCosting.STT
		scripWiseCostingRes.TotalTurnCharges += scripWiseCosting.TurnTax
		scripWiseCostingRes.TotalStampDuty += scripWiseCosting.Stamp
		scripWiseCostingRes.TotalOtherCharges += scripWiseCosting.OtherCharges
		scripWiseCostingRes.ScripWiseCosting = append(scripWiseCostingRes.ScripWiseCosting, scripWiseCosting)
	}

	scripWiseCostingRes.TotalCharges = math.Floor((scripWiseCostingRes.TotalBrokerage+scripWiseCostingRes.TotalGST+scripWiseCostingRes.TotalSEBITax+scripWiseCostingRes.TotalSTT+scripWiseCostingRes.TotalTurnCharges+scripWiseCostingRes.TotalStampDuty)*100) / 100
	scripWiseCostingRes.TotalBrokerage = math.Floor((scripWiseCostingRes.TotalBrokerage)*100) / 100
	scripWiseCostingRes.TotalGST = math.Floor((scripWiseCostingRes.TotalGST)*100) / 100
	scripWiseCostingRes.TotalSEBITax = math.Floor((scripWiseCostingRes.TotalSEBITax)*100) / 100
	scripWiseCostingRes.TotalSTT = math.Floor((scripWiseCostingRes.TotalSTT)*100) / 100
	scripWiseCostingRes.TotalTurnCharges = math.Floor((scripWiseCostingRes.TotalTurnCharges)*100) / 100
	scripWiseCostingRes.TotalStampDuty = math.Floor((scripWiseCostingRes.TotalStampDuty)*100) / 100
	scripWiseCostingRes.TotalOtherCharges = math.Floor((scripWiseCostingRes.TotalOtherCharges)*100) / 100

	sort.Sort(ByTrxDate(scripWiseCostingRes.ScripWiseCosting))

	loggerconfig.Info("GetScripWiseCostingData  resp=", scripWiseCostingRes, " uccId:", getScripWiseCosting.UserID, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)

	return scripWiseCostingRes, nil
}

// Implement the sort.Interface for ScripWiseCosting based on TrxDate
type ByTrxDateLedger []models.FinancialLedgerData

func (a ByTrxDateLedger) Len() int           { return len(a) }
func (a ByTrxDateLedger) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTrxDateLedger) Less(i, j int) bool { return a[i].UnixTimeFormat < a[j].UnixTimeFormat }

func (obj BackofficeObj) GetFinancialLedgerData(getFinancialLedgerReq models.GetFinancialLedgerDataReq, reqH models.ReqHeader) (models.FinancialLedgerRes, error) {

	var financialLedgerRes models.FinancialLedgerRes
	shilpiObj := InitBackofficeObj()

	url := shilpiObj.shilpiBaseUrl + "?Requesttype=" + constants.GetFinancialDateRange + "&userid=" + strings.ToUpper(getFinancialLedgerReq.UserID) + "&dfdatefr=" + getFinancialLedgerReq.DFDateFr + "&dfdateto=" + getFinancialLedgerReq.DFDateTo + "&ignoremarginentries=N"

	// emtpy payload
	payload := new(bytes.Buffer)

	res, err := apihelpers.CallApiShilpi(http.MethodGet, url, payload)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " GetFinancialLedgerData call api error =", err, " uccId:", getFinancialLedgerReq.UserID, " requestId:", reqH.RequestId)
		return financialLedgerRes, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("GetFinancialLedgerData error in reading response body=", err, " uccId:", getFinancialLedgerReq.UserID, " requestId:", reqH.RequestId)
		return financialLedgerRes, err
	}

	shilpiGetFinancials := ShilpiGetFinancials{}
	json.Unmarshal([]byte(string(body)), &shilpiGetFinancials)

	for i := 0; i < len(shilpiGetFinancials); i++ {
		var financialLedgerData models.FinancialLedgerData
		if i == 0 {
			financialLedgerRes.OpeningBalance, err = strconv.ParseFloat(shilpiGetFinancials[i].RuningBal, 64)
		}
		if i == len(shilpiGetFinancials)-1 {
			financialLedgerRes.ClosingBalance, err = strconv.ParseFloat(shilpiGetFinancials[i].RuningBal, 64)
		}
		financialLedgerData.TransactionDate = shilpiGetFinancials[i].ValueDate
		financialLedgerData.SettlementDate = shilpiGetFinancials[i].Date
		financialLedgerData.Credit, err = strconv.ParseFloat(shilpiGetFinancials[i].Credit, 64)
		financialLedgerData.Debit, err = strconv.ParseFloat(shilpiGetFinancials[i].Debit, 64)
		financialLedgerData.Exchange = shilpiGetFinancials[i].Exchange
		financialLedgerData.TransactionDetails = shilpiGetFinancials[i].Narr
		financialLedgerData.SettlementNumber = shilpiGetFinancials[i].PayRefNo
		financialLedgerData.NetBalance, err = strconv.ParseFloat(shilpiGetFinancials[i].RuningBal, 64)
		financialLedgerData.Segment = shilpiGetFinancials[i].Segment
		if shilpiGetFinancials[i].Narr != constants.FundPayment {
			financialLedgerRes.Inflow += financialLedgerData.Credit
		} else {
			financialLedgerRes.FundsWithdrawn += financialLedgerData.Debit
		}
		if shilpiGetFinancials[i].Narr != constants.FundReceived {
			financialLedgerRes.Outflow += financialLedgerData.Debit
		} else {
			financialLedgerRes.FundsReceived += financialLedgerData.Credit
		}
		if err != nil {
			loggerconfig.Error("GetFinancialLedgerData, error in parsing string to float64 Error : ", err)
			return financialLedgerRes, err
		}

		if i != 0 {
			// Parse the input date string using the specified layout
			parsedTime, err := time.Parse(constants.DDMMYYYY, financialLedgerData.TransactionDate)
			if err != nil {
				loggerconfig.Error("GetFinancialLedgerData, Error parsing shilpi date for tradebook :", err)
				return financialLedgerRes, err
			}

			// Convert the parsed Date to Unix time
			financialLedgerData.UnixTimeFormat = parsedTime.Unix()
		}

		financialLedgerRes.FinancialLedger = append(financialLedgerRes.FinancialLedger, financialLedgerData)
	}

	if len(financialLedgerRes.FinancialLedger) > 1 {
		sort.Sort(ByTrxDateLedger(financialLedgerRes.FinancialLedger[1:]))
	}

	loggerconfig.Info("GetFinancialLedgerData  resp=", financialLedgerRes, " uccId:", getFinancialLedgerReq.UserID, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)

	return financialLedgerRes, nil
}

func (obj BackofficeObj) GetOpenPositionData(openPositionReq models.OpenPositionReq, reqH models.ReqHeader) (models.OpenPositionRes, error) {

	var openPositionRes models.OpenPositionRes
	shilpiObj := InitBackofficeObj()

	url := shilpiObj.shilpiBaseUrl + "?Requesttype=" + constants.OpenPosition + "&userid=" + strings.ToUpper(openPositionReq.UserID) + "&dsflag=" + openPositionReq.Dsflag

	// emtpy payload
	payload := new(bytes.Buffer)

	res, err := apihelpers.CallApiShilpi(http.MethodGet, url, payload)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " GetOpenPositionData call api error =", err, " uccId:", openPositionReq.UserID, " requestId:", reqH.RequestId)
		return openPositionRes, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("GetOpenPositionData error in reading response body=", err, " uccId:", openPositionReq.UserID, " requestId:", reqH.RequestId)
		return openPositionRes, err
	}

	shilpiGetOpenPosition := ShilpiGetOpenPosition{}
	json.Unmarshal([]byte(string(body)), &shilpiGetOpenPosition)

	for i := 0; i < len(shilpiGetOpenPosition); i++ {
		if shilpiGetOpenPosition[i].NetQty != "0" {
			var openPositionData models.OpenPositionData
			openPositionData.ScripName = shilpiGetOpenPosition[i].Contract
			openPositionData.ExpiryDate = shilpiGetOpenPosition[i].ExpiryDate
			openPositionData.Exchange = shilpiGetOpenPosition[i].Exchange
			openPositionData.StrikePrice, err = strconv.ParseFloat(shilpiGetOpenPosition[i].StrikePrice, 64)
			openPositionData.OpenQuantity, err = strconv.ParseFloat(shilpiGetOpenPosition[i].NetQty, 64)
			openPositionData.ClosingPrice, err = strconv.ParseFloat(shilpiGetOpenPosition[i].ClosingRate, 64)
			if shilpiGetOpenPosition[i].OptionType == constants.CE || shilpiGetOpenPosition[i].OptionType == constants.PE {
				openPositionData.OptionType = shilpiGetOpenPosition[i].OptionType
				openPositionData.InstrumentType = constants.OPTIONS
			} else {
				openPositionData.InstrumentType = constants.FUTURES
			}
			if shilpiGetOpenPosition[i].BuyQty > shilpiGetOpenPosition[i].SellQty {
				var buyPrice, buyQty float64
				buyPrice, err = strconv.ParseFloat(shilpiGetOpenPosition[i].BuyValue, 64)
				buyQty, err = strconv.ParseFloat(shilpiGetOpenPosition[i].BuyQty, 64)
				openPositionData.AveragePrice = buyPrice / buyQty
				openPositionData.BuySellType = constants.BUY
			} else {
				var sellPrice, sellQty float64
				sellPrice, err = strconv.ParseFloat(shilpiGetOpenPosition[i].SellValue, 64)
				sellQty, err = strconv.ParseFloat(shilpiGetOpenPosition[i].SellQty, 64)
				openPositionData.AveragePrice = sellPrice / sellQty
				openPositionData.BuySellType = constants.SELL
			}
			openPositionData.UnrealisedProfitOrLoss = (openPositionData.ClosingPrice - openPositionData.AveragePrice) * openPositionData.OpenQuantity
			if err != nil {
				loggerconfig.Error("GetOpenPositionData, error in parsing string to float64 Error : ", err)
				return openPositionRes, err
			}

			if shilpiGetOpenPosition[i].Segment == constants.CUR {
				if openPositionData.InstrumentType == constants.OPTIONS {
					openPositionRes.CurrencyOptionMTM += openPositionData.UnrealisedProfitOrLoss
				} else {
					openPositionRes.CurrencyFutureMTM += openPositionData.UnrealisedProfitOrLoss
				}
				openPositionRes.CurrencyDerivative = append(openPositionRes.CurrencyDerivative, openPositionData)
			} else if shilpiGetOpenPosition[i].Segment == constants.COM {
				if openPositionData.InstrumentType == constants.OPTIONS {
					openPositionRes.CommodityOptionMTM += openPositionData.UnrealisedProfitOrLoss
				} else {
					openPositionRes.CommodityFutureMTM += openPositionData.UnrealisedProfitOrLoss
				}
				openPositionRes.CommodityDerivative = append(openPositionRes.CommodityDerivative, openPositionData)
			} else {
				if openPositionData.InstrumentType == constants.OPTIONS {
					openPositionRes.EquityOptionMTM += openPositionData.UnrealisedProfitOrLoss
				} else {
					openPositionRes.EquityFutureMTM += openPositionData.UnrealisedProfitOrLoss
				}
				openPositionRes.EquityDerivative = append(openPositionRes.EquityDerivative, openPositionData)
			}

		}
	}

	loggerconfig.Info("GetOpenPositionData  resp=", openPositionRes, " uccId:", openPositionReq.UserID, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)

	return openPositionRes, nil
}

func (obj BackofficeObj) GetFONetPositionData(getFONetPositionDataReq models.GetFONetPositionDataReq, reqH models.ReqHeader) (models.FONetPositionRes, error) {

	var foNetPositionRes models.FONetPositionRes
	shilpiObj := InitBackofficeObj()

	url := shilpiObj.shilpiBaseUrl + "?Requesttype=" + constants.NetPositionFO + "&userid=" + strings.ToUpper(getFONetPositionDataReq.UserID) + "&datefr=" + getFONetPositionDataReq.DFDateFr + "&dateto=" + getFONetPositionDataReq.DFDateTo

	// emtpy payload
	payload := new(bytes.Buffer)

	res, err := apihelpers.CallApiShilpi(http.MethodGet, url, payload)
	if err != nil {
		logrus.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " GetFONetPositionData call api error =", err, " uccId:", getFONetPositionDataReq.UserID, " requestId:", reqH.RequestId)
		return foNetPositionRes, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logrus.Error("GetFONetPositionData error in reading response body=", err, " uccId:", getFONetPositionDataReq.UserID, " requestId:", reqH.RequestId)
		return foNetPositionRes, err
	}

	shilpiFONetPosition := ShilpiFONetPosition{}
	json.Unmarshal([]byte(string(body)), &shilpiFONetPosition)

	var foNetPositionSummaryData models.FONetPositionSummaryData
	foNetPositionSummaryData.DateRange = getFONetPositionDataReq.DFDateFr + " to " + getFONetPositionDataReq.DFDateTo

	for i := 0; i < len(shilpiFONetPosition); i++ {
		if shilpiFONetPosition[i].ID == "1" {
			var foNetPositionDetailsData models.FONetPositionDetailsData
			foNetPositionDetailsData.Symbol = shilpiFONetPosition[i].Contracts
			if shilpiFONetPosition[i].OptionType == constants.PE || shilpiFONetPosition[i].OptionType == constants.CE {
				foNetPositionDetailsData.InstrumentType = constants.OPTIONS
				if shilpiFONetPosition[i].OptionType == constants.PE {
					foNetPositionDetailsData.OptionType = constants.PUT
				} else {
					foNetPositionDetailsData.OptionType = constants.CALL
				}
			} else {
				foNetPositionDetailsData.InstrumentType = constants.FUTURES
			}
			foNetPositionDetailsData.StrikePrice, err = strconv.ParseFloat(shilpiFONetPosition[i].StrikePrice, 64)
			foNetPositionDetailsData.ExpiryDate = shilpiFONetPosition[i].ExpDate
			foNetPositionDetailsData.Quantity = shilpiFONetPosition[i].SaleQty
			foNetPositionDetailsData.BuyPrice, err = strconv.ParseFloat(shilpiFONetPosition[i].BuyAvgRate, 64)
			foNetPositionDetailsData.BuyValue, err = strconv.ParseFloat(shilpiFONetPosition[i].BuyValue, 64)
			foNetPositionDetailsData.SellPrice, err = strconv.ParseFloat(shilpiFONetPosition[i].SaleAvgRate, 64)
			foNetPositionDetailsData.SellValue, err = strconv.ParseFloat(shilpiFONetPosition[i].BuyValue, 64)
			foNetPositionDetailsData.RealizedPNL, err = strconv.ParseFloat(shilpiFONetPosition[i].NetPL, 64)
			foNetPositionDetailsData.PreviousClosingPrice, err = strconv.ParseFloat(shilpiFONetPosition[i].ClosPrice, 64)
			foNetPositionDetailsData.OpenQuantity = shilpiFONetPosition[i].NetQty
			foNetPositionDetailsData.OpenValue, err = strconv.ParseFloat(shilpiFONetPosition[i].NetValue, 64)
			if err != nil {
				logrus.Error("foNetPositionDetailsData, error in parsing string to float64 Error : ", err)
				return foNetPositionRes, err
			}
			if foNetPositionDetailsData.OpenQuantity < 0 {
				foNetPositionDetailsData.UnrealizedPNL = float64(foNetPositionDetailsData.OpenQuantity) * (foNetPositionDetailsData.PreviousClosingPrice - foNetPositionDetailsData.SellPrice)
			} else {
				foNetPositionDetailsData.UnrealizedPNL = float64(foNetPositionDetailsData.OpenQuantity) * (foNetPositionDetailsData.PreviousClosingPrice - foNetPositionDetailsData.BuyPrice)
			}
			foNetPositionSummaryData.RealisedPNL += foNetPositionDetailsData.RealizedPNL
			foNetPositionSummaryData.UnRealisedPNL += foNetPositionDetailsData.UnrealizedPNL
			foNetPositionRes.FONetPositionDetails = append(foNetPositionRes.FONetPositionDetails, foNetPositionDetailsData)

		} else {
			var foNetPositionChargesData models.FONetPositionChargesData
			foNetPositionChargesData.Brockerage, err = strconv.ParseFloat(shilpiFONetPosition[i].MinBrok, 64)
			foNetPositionChargesData.ExchangeTransactionCharges, err = strconv.ParseFloat(shilpiFONetPosition[i].ExchClg, 64)
			foNetPositionChargesData.IntegratedGST, err = strconv.ParseFloat(shilpiFONetPosition[i].TotalGST, 64)
			foNetPositionChargesData.SecuritiesTransactionTax, err = strconv.ParseFloat(shilpiFONetPosition[i].STT, 64)
			foNetPositionChargesData.SEBIFees, err = strconv.ParseFloat(shilpiFONetPosition[i].SebiFee, 64)
			foNetPositionChargesData.StampDuty, err = strconv.ParseFloat(shilpiFONetPosition[i].StampDuty, 64)
			if err != nil {
				logrus.Error("foNetPositionChargesData, error in parsing string to float64 Error : ", err)
				return foNetPositionRes, err
			}
			foNetPositionChargesData.TotalCharges = foNetPositionChargesData.Brockerage + foNetPositionChargesData.ExchangeTransactionCharges + foNetPositionChargesData.IntegratedGST + foNetPositionChargesData.SecuritiesTransactionTax + foNetPositionChargesData.SEBIFees + foNetPositionChargesData.StampDuty
			foNetPositionRes.ChargesDetails = foNetPositionChargesData
			foNetPositionSummaryData.Charges = foNetPositionChargesData.TotalCharges
		}
	}
	foNetPositionSummaryData.NetPNL = foNetPositionSummaryData.RealisedPNL + foNetPositionSummaryData.UnRealisedPNL - foNetPositionSummaryData.Charges
	foNetPositionRes.Summary = foNetPositionSummaryData

	logrus.Info("GetFONetPositionData  resp=", foNetPositionRes, " uccId:", getFONetPositionDataReq.UserID, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)

	return foNetPositionRes, nil
}

func (obj BackofficeObj) GetHoldingFinancialData(getHoldingFinancialDataReq models.GetHoldingFinancialDataReq, reqH models.ReqHeader) (models.GetHoldingFinancialDataRes, error) {

	var holdingFinancialDataRes models.GetHoldingFinancialDataRes
	shilpiObj := InitBackofficeObj()

	url := shilpiObj.shilpiBaseUrl + "?Requesttype=" + constants.HoldingCumFinancial + "&userid=" + strings.ToUpper(getHoldingFinancialDataReq.UserID)

	// emtpy payload
	payload := new(bytes.Buffer)

	res, err := apihelpers.CallApiShilpi(http.MethodGet, url, payload)
	if err != nil {
		logrus.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " GetHoldingFinancialData call api error =", err, " uccId:", getHoldingFinancialDataReq.UserID, " requestId:", reqH.RequestId)
		return holdingFinancialDataRes, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logrus.Error("GetHoldingFinancialData error in reading response body=", err, " uccId:", getHoldingFinancialDataReq.UserID, " requestId:", reqH.RequestId)
		return holdingFinancialDataRes, err
	}

	shilpiHoldingFinancial := ShilpiHoldingFinancial{}
	json.Unmarshal([]byte(string(body)), &shilpiHoldingFinancial)

	var holdingSummaryData models.HoldingSummaryData

	for i := 0; i < len(shilpiHoldingFinancial); i++ {
		var holdingFinancialData models.GetHoldingFinancialData

		holdingFinancialData.Isin = shilpiHoldingFinancial[i].ISIN
		holdingFinancialData.Instrument = shilpiHoldingFinancial[i].ScripName
		holdingFinancialData.PledgedQty, err = strconv.ParseFloat(shilpiHoldingFinancial[i].PledgeStock, 64)
		holdingFinancialData.FreeQty, err = strconv.ParseFloat(shilpiHoldingFinancial[i].DPStock, 64)
		holdingFinancialData.ClosingPrice, err = strconv.ParseFloat(shilpiHoldingFinancial[i].ClosingPrice, 64)
		holdingFinancialData.AvgBuyPrice, err = strconv.ParseFloat(shilpiHoldingFinancial[i].BuyPrice, 64)
		holdingFinancialData.TotalPledgedValue = math.Floor((holdingFinancialData.PledgedQty*holdingFinancialData.ClosingPrice)*100) / 100
		pledgeBenefit, err1 := strconv.ParseFloat(shilpiHoldingFinancial[i].Variance, 64)
		if err1 != nil {
			logrus.Error("GetHoldingFinancialData, error in parsing string to float64 Error : ", err)
			return holdingFinancialDataRes, err
		}
		holdingFinancialData.HaircutPercentage = math.Ceil((100-pledgeBenefit)*100) / 100
		holdingFinancialData.MarginAvailableAfterHaircut = math.Floor((holdingFinancialData.TotalPledgedValue*pledgeBenefit)*100) / 100
		holding, err2 := strconv.ParseFloat(shilpiHoldingFinancial[i].Holding, 64)
		if err2 != nil {
			logrus.Error("GetHoldingFinancialData, error in parsing string to float64 Error : ", err)
			return holdingFinancialDataRes, err
		}
		holdingFinancialData.TotalQty = holdingFinancialData.PledgedQty + holdingFinancialData.FreeQty + holding
		holdingFinancialData.InvestmentValue = math.Ceil((holdingFinancialData.TotalQty*holdingFinancialData.AvgBuyPrice)*100) / 100
		holdingFinancialData.CurrentValue = math.Ceil((holdingFinancialData.TotalQty*holdingFinancialData.ClosingPrice)*100) / 100
		holdingFinancialData.UnrealizedProfitLoss = math.Ceil((holdingFinancialData.CurrentValue-holdingFinancialData.InvestmentValue)*100) / 100
		if holdingFinancialData.UnrealizedProfitLoss != 0 && holdingFinancialData.InvestmentValue != 0 {
			holdingFinancialData.NetChange = math.Ceil(((holdingFinancialData.UnrealizedProfitLoss/holdingFinancialData.InvestmentValue)*100)*100) / 100
		}
		holdingSummaryData.InvestedValue += holdingFinancialData.InvestmentValue
		holdingSummaryData.CurrentValue += holdingFinancialData.CurrentValue
		holdingSummaryData.UnrealisedPNL += holdingFinancialData.UnrealizedProfitLoss
		holdingSummaryData.TotalPledgeValue += holdingFinancialData.TotalPledgedValue
		holdingSummaryData.TotalMarginValueAfterHaircut += holdingFinancialData.MarginAvailableAfterHaircut

		holdingFinancialDataRes.HoldingFinancialData = append(holdingFinancialDataRes.HoldingFinancialData, holdingFinancialData)
	}
	for i := 0; i < len(holdingFinancialDataRes.HoldingFinancialData); i++ {
		if holdingFinancialDataRes.HoldingFinancialData[i].CurrentValue != 0 && holdingSummaryData.CurrentValue != 0 {
			holdingFinancialDataRes.HoldingFinancialData[i].ContributionPercentage = math.Ceil(((holdingFinancialDataRes.HoldingFinancialData[i].CurrentValue/holdingSummaryData.CurrentValue)*100)*100) / 100
		}
	}
	holdingFinancialDataRes.HoldingSummary = holdingSummaryData
	holdingSummaryData.InvestedValue = math.Ceil((holdingSummaryData.InvestedValue)*100) / 100
	holdingSummaryData.CurrentValue = math.Ceil((holdingSummaryData.CurrentValue)*100) / 100
	holdingSummaryData.UnrealisedPNL = math.Ceil((holdingSummaryData.UnrealisedPNL)*100) / 100
	holdingSummaryData.TotalPledgeValue = math.Ceil((holdingSummaryData.TotalPledgeValue)*100) / 100
	holdingSummaryData.TotalMarginValueAfterHaircut = math.Ceil((holdingSummaryData.TotalMarginValueAfterHaircut)*100) / 100

	logrus.Info("GetFONetPositionData  resp=", holdingFinancialDataRes, " uccId:", getHoldingFinancialDataReq.UserID, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)

	return holdingFinancialDataRes, nil
}

// Implement the sort.Interface for ScripWiseCosting based on TrxDate
type ByTrxDateCommodity []models.CommodityTransactionData

func (a ByTrxDateCommodity) Len() int           { return len(a) }
func (a ByTrxDateCommodity) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTrxDateCommodity) Less(i, j int) bool { return a[i].UnixTimeFormat < a[j].UnixTimeFormat }

func (obj BackofficeObj) GetCommodityTransactionData(getCommodityTransactionReq models.CommodityTradebookReq, reqH models.ReqHeader) (models.CommodityTransactionRes, error) {

	var commodityTransactionRes models.CommodityTransactionRes
	shilpiObj := InitBackofficeObj()

	url := shilpiObj.shilpiBaseUrl + "?Requesttype=" + constants.GetCommodityTransaction + "&userid=" + strings.ToUpper(getCommodityTransactionReq.UserID) + "&datefr=" + getCommodityTransactionReq.DFDateFr + "&dateto=" + getCommodityTransactionReq.DFDateTo

	// emtpy payload
	payload := new(bytes.Buffer)

	res, err := apihelpers.CallApiShilpi(http.MethodGet, url, payload)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " GetCommodityTransactionData call api error =", err, " uccId:", getCommodityTransactionReq.UserID, " requestId:", reqH.RequestId)
		return commodityTransactionRes, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("GetCommodityTransactionData error in reading response body=", err, " uccId:", getCommodityTransactionReq.UserID, " requestId:", reqH.RequestId)
		return commodityTransactionRes, err
	}

	shilpiCommodityTransaction := ShilpiCommodityTransaction{}
	json.Unmarshal([]byte(string(body)), &shilpiCommodityTransaction)

	for i := 0; i < len(shilpiCommodityTransaction); i++ {
		var commodityTransactionData models.CommodityTransactionData

		commodityTransactionData.Symbol = shilpiCommodityTransaction[i].Symbol
		commodityTransactionData.InstrumentType = shilpiCommodityTransaction[i].InstrumentType
		commodityTransactionData.ExpiryDate = shilpiCommodityTransaction[i].ExpiryDate
		commodityTransactionData.OptionType = shilpiCommodityTransaction[i].OptionType
		commodityTransactionData.StrikePrice, _ = strconv.ParseFloat(shilpiCommodityTransaction[i].StrikePrice, 64)
		commodityTransactionData.TradeDate = shilpiCommodityTransaction[i].TradeDate
		commodityTransactionData.TradePrice, _ = strconv.ParseFloat(shilpiCommodityTransaction[i].TradePrice, 64)
		commodityTransactionData.Unit, _ = strconv.ParseFloat(shilpiCommodityTransaction[i].Unit, 64)
		commodityTransactionData.TradeQty, _ = strconv.ParseFloat(shilpiCommodityTransaction[i].TradeQty, 64)
		commodityTransactionData.Brokerage, _ = strconv.ParseFloat(shilpiCommodityTransaction[i].Brokerage, 64)
		commodityTransactionData.GST, _ = strconv.ParseFloat(shilpiCommodityTransaction[i].GST, 64)
		commodityTransactionData.CTT, _ = strconv.ParseFloat(shilpiCommodityTransaction[i].CTT, 64)
		commodityTransactionData.SEBITax, _ = strconv.ParseFloat(shilpiCommodityTransaction[i].SEBITax, 64)
		commodityTransactionData.TurnoverTax, _ = strconv.ParseFloat(shilpiCommodityTransaction[i].TurnoverTax, 64)
		commodityTransactionData.StampDuty, _ = strconv.ParseFloat(shilpiCommodityTransaction[i].StampDuty, 64)
		commodityTransactionData.CLGTax, _ = strconv.ParseFloat(shilpiCommodityTransaction[i].CLGTax, 64)
		commodityTransactionData.Segment = constants.Commodity
		commodityTransactionData.Exchange = shilpiCommodityTransaction[i].Exchange
		commodityTransactionData.OrderNo = shilpiCommodityTransaction[i].OrderNo
		commodityTransactionData.TradeNo = shilpiCommodityTransaction[i].TradeNo
		commodityTransactionData.TradeTime = shilpiCommodityTransaction[i].TradeTime

		if shilpiCommodityTransaction[i].BuySellInd == 1 {
			commodityTransactionData.BuySellInd = constants.BUY
		} else if shilpiCommodityTransaction[i].BuySellInd == 2 {
			commodityTransactionData.BuySellInd = constants.SELL
		}

		// Parse the input date string using the specified layout
		parsedTime, err := time.Parse(constants.ShilpiDateFormatWithTime, commodityTransactionData.TradeDate)
		if err != nil {
			loggerconfig.Error("GetCommodityTransactionData, Error parsing shilpi date for tradebook :", err, "reqId: ", reqH.RequestId)
			return commodityTransactionRes, err
		}

		// Convert the parsed Date to Unix time
		commodityTransactionData.UnixTimeFormat = parsedTime.Unix()

		commodityTransactionRes.TotalBrokerage += commodityTransactionData.Brokerage
		commodityTransactionRes.TotalGST += commodityTransactionData.GST
		commodityTransactionRes.TotalSEBITax += commodityTransactionData.SEBITax
		commodityTransactionRes.TotalCTT += commodityTransactionData.CTT
		commodityTransactionRes.TotalTurnCharges += commodityTransactionData.TurnoverTax
		commodityTransactionRes.TotalStampDuty += commodityTransactionData.StampDuty
		commodityTransactionRes.CommodityTransactions = append(commodityTransactionRes.CommodityTransactions, commodityTransactionData)
	}

	commodityTransactionRes.TotalCharges = math.Ceil((commodityTransactionRes.TotalBrokerage+commodityTransactionRes.TotalGST+commodityTransactionRes.TotalSEBITax+commodityTransactionRes.TotalCTT+commodityTransactionRes.TotalTurnCharges+commodityTransactionRes.TotalStampDuty)*100) / 100
	commodityTransactionRes.TotalBrokerage = math.Ceil((commodityTransactionRes.TotalBrokerage)*100) / 100
	commodityTransactionRes.TotalGST = math.Ceil((commodityTransactionRes.TotalGST)*100) / 100
	commodityTransactionRes.TotalSEBITax = math.Ceil((commodityTransactionRes.TotalSEBITax)*100) / 100
	commodityTransactionRes.TotalCTT = math.Ceil((commodityTransactionRes.TotalCTT)*100) / 100
	commodityTransactionRes.TotalTurnCharges = math.Ceil((commodityTransactionRes.TotalTurnCharges)*100) / 100
	commodityTransactionRes.TotalStampDuty = math.Ceil((commodityTransactionRes.TotalStampDuty)*100) / 100

	sort.Sort(ByTrxDateCommodity(commodityTransactionRes.CommodityTransactions))

	loggerconfig.Info("GetCommodityTransactionData  resp=", commodityTransactionRes, " uccId:", getCommodityTransactionReq.UserID, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)

	return commodityTransactionRes, nil
}

func (obj BackofficeObj) GetFNOTransactionData(getFNOTransactionReq models.FNOTradebookReq, reqH models.ReqHeader) (models.FNOTransactionRes, error) {

	var fnoTransactionRes models.FNOTransactionRes
	shilpiObj := InitBackofficeObj()

	url := shilpiObj.shilpiBaseUrl + "?Requesttype=" + constants.GetFNOTransaction + "&userid=" + strings.ToUpper(getFNOTransactionReq.UserID) + "&datefr=" + getFNOTransactionReq.DFDateFr + "&dateto=" + getFNOTransactionReq.DFDateTo

	// emtpy payload
	payload := new(bytes.Buffer)

	res, err := apihelpers.CallApiShilpi(http.MethodGet, url, payload)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " GetFNOTransactionData call api error =", err, " uccId:", getFNOTransactionReq.UserID, " requestId:", reqH.RequestId)
		return fnoTransactionRes, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerconfig.Error("GetFNOTransactionData error in reading response body=", err, " uccId:", getFNOTransactionReq.UserID, " requestId:", reqH.RequestId)
		return fnoTransactionRes, err
	}

	shilpiFNOTransaction := ShilpiFNOTransaction{}
	json.Unmarshal([]byte(string(body)), &shilpiFNOTransaction)

	for i := 0; i < len(shilpiFNOTransaction); i++ {
		var fnoTransactionData models.FNOTransactionData

		fnoTransactionData.Symbol = shilpiFNOTransaction[i].Symbol
		fnoTransactionData.InstrumentType = shilpiFNOTransaction[i].InstrumentType
		fnoTransactionData.ExpiryDate = shilpiFNOTransaction[i].ExpiryDate
		fnoTransactionData.OptionType = shilpiFNOTransaction[i].OptionType
		fnoTransactionData.StrikePrice, _ = strconv.ParseFloat(shilpiFNOTransaction[i].StrikePrice, 64)
		fnoTransactionData.TradeDate = shilpiFNOTransaction[i].TradeDate
		fnoTransactionData.TradeQty, _ = strconv.ParseFloat(shilpiFNOTransaction[i].TradeQty, 64)
		fnoTransactionData.TradePrice, _ = strconv.ParseFloat(shilpiFNOTransaction[i].TradePrice, 64)
		fnoTransactionData.Brokerage, _ = strconv.ParseFloat(shilpiFNOTransaction[i].Brokerage, 64)
		fnoTransactionData.GST, _ = strconv.ParseFloat(shilpiFNOTransaction[i].GST, 64)
		fnoTransactionData.STT, _ = strconv.ParseFloat(shilpiFNOTransaction[i].STT, 64)
		fnoTransactionData.SEBITax, _ = strconv.ParseFloat(shilpiFNOTransaction[i].SebiTax, 64)
		fnoTransactionData.TurnoverTax, _ = strconv.ParseFloat(shilpiFNOTransaction[i].TurnoverTax, 64)
		fnoTransactionData.StampDuty, _ = strconv.ParseFloat(shilpiFNOTransaction[i].StampDuty, 64)
		fnoTransactionData.ClearingCharges, _ = strconv.ParseFloat(shilpiFNOTransaction[i].ClgTax, 64)
		fnoTransactionData.IPFTax, _ = strconv.ParseFloat(shilpiFNOTransaction[i].IPFTax, 64)
		fnoTransactionData.Segment = constants.FutnOpt
		fnoTransactionData.Exchange = shilpiFNOTransaction[i].Exchange
		fnoTransactionData.OrderNo = shilpiFNOTransaction[i].OrderNo
		fnoTransactionData.TradeNo = shilpiFNOTransaction[i].TradeNo
		fnoTransactionData.TradeTime = shilpiFNOTransaction[i].TradeTime

		if shilpiFNOTransaction[i].BuySellInd == 1 {
			fnoTransactionData.BuySellInd = constants.BUY
		} else if shilpiFNOTransaction[i].BuySellInd == 2 {
			fnoTransactionData.BuySellInd = constants.SELL
		}

		// Parse the input date string using the specified layout
		parsedTime, err := time.Parse(constants.ShilpiDateFormatWithTime, fnoTransactionData.TradeDate)
		if err != nil {
			loggerconfig.Error("GetFNOTransactionData, Error parsing shilpi date for tradebook :", err, "reqId: ", reqH.RequestId)
			return fnoTransactionRes, err
		}

		// Convert the parsed Date to Unix time
		fnoTransactionData.UnixTimeFormat = parsedTime.Unix()

		fnoTransactionRes.TotalBrokerage += fnoTransactionData.Brokerage
		fnoTransactionRes.TotalGST += fnoTransactionData.GST
		fnoTransactionRes.TotalSEBITax += fnoTransactionData.SEBITax
		fnoTransactionRes.TotalSTT += fnoTransactionData.STT
		fnoTransactionRes.TotalTurnCharges += fnoTransactionData.TurnoverTax
		fnoTransactionRes.TotalStampDuty += fnoTransactionData.StampDuty
		fnoTransactionRes.TotalClearingCharges += fnoTransactionData.ClearingCharges
		fnoTransactionRes.FNOTransactions = append(fnoTransactionRes.FNOTransactions, fnoTransactionData)
	}

	fnoTransactionRes.TotalCharges = math.Ceil((fnoTransactionRes.TotalBrokerage+fnoTransactionRes.TotalGST+fnoTransactionRes.TotalSEBITax+fnoTransactionRes.TotalSTT+fnoTransactionRes.TotalTurnCharges+fnoTransactionRes.TotalStampDuty+fnoTransactionRes.TotalClearingCharges)*100) / 100
	fnoTransactionRes.TotalBrokerage = math.Ceil((fnoTransactionRes.TotalBrokerage)*100) / 100
	fnoTransactionRes.TotalGST = math.Ceil((fnoTransactionRes.TotalGST)*100) / 100
	fnoTransactionRes.TotalSEBITax = math.Ceil((fnoTransactionRes.TotalSEBITax)*100) / 100
	fnoTransactionRes.TotalSTT = math.Ceil((fnoTransactionRes.TotalSTT)*100) / 100
	fnoTransactionRes.TotalTurnCharges = math.Ceil((fnoTransactionRes.TotalTurnCharges)*100) / 100
	fnoTransactionRes.TotalStampDuty = math.Ceil((fnoTransactionRes.TotalStampDuty)*100) / 100
	fnoTransactionRes.TotalClearingCharges = math.Ceil((fnoTransactionRes.TotalClearingCharges)*100) / 100

	// sort.Sort(ByTrxDateCommodity(fnoTransactionRes.FNOTransactions))

	loggerconfig.Info("GetFNOTransactionData  resp=", fnoTransactionRes, " uccId:", getFNOTransactionReq.UserID, " StatusCode : ", res.StatusCode, " requestId:", reqH.RequestId)

	return fnoTransactionRes, nil
}
