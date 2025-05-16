package tradelab

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/db"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
	"strconv"
	"strings"
	"sync"
	"time"
)

type EpledgeObj struct {
	tradeLabURL string
}

func InitEpledgeProvider() EpledgeObj {
	defer models.HandlePanic()

	EpledgeObj := EpledgeObj{
		tradeLabURL: constants.TLURL,
	}

	return EpledgeObj
}

func (obj EpledgeObj) SendEpledgeRequest(req models.EpledgeReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + EPLEDGEREQ

	var tlEpledgeReq TradeLabEpledgeReq

	tlEpledgeReq.Depository = req.Depository
	tlEpledgeReq.ClientID = req.ClientID
	tlEpledgeReq.Exchange = req.Exchange
	tlEpledgeReq.BoId = req.BoId
	tlEpledgeReq.Segment = req.Segment

	for _, isin := range req.IsinDetails {
		tradeLabIsin := TradeLabIsin{
			IsinName: isin.IsinName,
			Isin:     isin.Isin,
			Quantity: isin.Quantity,
			Price:    isin.Price,
		}
		tlEpledgeReq.IsinDetails = append(tlEpledgeReq.IsinDetails, tradeLabIsin)
	}

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlEpledgeReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "SendEpledgeRequest", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " SendEpledgeRequest call api error", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("SendEpledgeRequest res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		if tlErrorRes.Message == PledgeTimeError || tlErrorRes.Message == PledgeTimeErrorNew || strings.Contains(strings.ToLower(tlErrorRes.Message), strings.ToLower(PledgeTimeErrorOnlyText)) {
			apiRes.Message = PledgeHours
		}
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlEpledgeRes := TradeLabEpledgeRes{}
	json.Unmarshal([]byte(string(body)), &tlEpledgeRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " SendEpledgeRequest tl status not ok =", tlEpledgeRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlEpledgeRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var EpledgeRes models.EpledgeRes

	EpledgeRes.Data.DpId = tlEpledgeRes.Data.DpId
	EpledgeRes.Data.PledgedTls = tlEpledgeRes.Data.PledgedTls
	EpledgeRes.Data.ReqId = tlEpledgeRes.Data.ReqId
	EpledgeRes.Data.Version = tlEpledgeRes.Data.Version
	EpledgeRes.Message = tlEpledgeRes.Message
	EpledgeRes.Status = tlEpledgeRes.Status

	loggerconfig.Info("SendEpledgeRequest resp=", EpledgeRes, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = EpledgeRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	//Return the response immediately, asynchronously update db
	go func() {
		var batch []models.PledgeData
		batchSize := 100
		for i := 0; i < len(req.IsinDetails); i++ {
			insertInfoInDB := models.PledgeData{
				ClientID:       req.ClientID,
				SegmentID:      req.Segment,
				Timestamp:      helpers.GetCurrentTimeInIST(),
				ISIN:           req.IsinDetails[i].Isin,
				Quantity:       req.IsinDetails[i].Quantity,
				Price:          req.IsinDetails[i].Price,
				Exchange:       req.Exchange,
				BOID:           req.BoId,
				Depository:     req.Depository,
				PledgeUnpledge: constants.Pledge,
				DPID:           tlEpledgeRes.Data.DpId,
				PledgeTLS:      tlEpledgeRes.Data.PledgedTls,
				ReqID:          tlEpledgeRes.Data.ReqId,
				Version:        tlEpledgeRes.Data.Version,
				Status:         constants.Started,
			}

			batch = append(batch, insertInfoInDB)

			if len(batch) >= batchSize {
				if err := db.GetPgObj().InsertPledgeDataBatch(batch); err != nil {
					loggerconfig.Error("SendEpledgeRequest Batch insertion failed, error:", err, " requestId:", reqH.RequestId)
				}
				batch = nil
			}
		}

		if len(batch) > 0 {
			if err := db.GetPgObj().InsertPledgeDataBatch(batch); err != nil {
				loggerconfig.Error("SendEpledgeRequest Final batch insertion failed, error:", err, " requestId:", reqH.RequestId)
			}
		}
	}()

	return http.StatusOK, apiRes
}

func CallUnpledge(req TradeLabUnPledgeReq, url string, wg *sync.WaitGroup, ch chan<- models.UnpledgeRes, reqH models.ReqHeader) {
	defer wg.Done()

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(req)

	var unpledgeRes models.UnpledgeRes
	unpledgeRes.Isin = req.Isin

	//call api
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "CallUnpledge", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " UnpledgeRequest call api error", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion, "reqBody :", req)
		unpledgeRes.Status = false
		unpledgeRes.Message = err.Error()
		ch <- unpledgeRes
		return
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("UnpledgeRequest res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion, "reqBody :", req)
		unpledgeRes.Status = false
		unpledgeRes.Message = tlErrorRes.Message
		unpledgeRes.Error.Code = strconv.Itoa(tlErrorRes.ErrorCode)
		ch <- unpledgeRes
		return
	}

	tlUnpledgeRes := TradeLabUnpledgeRes{}
	json.Unmarshal([]byte(string(body)), &tlUnpledgeRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " UnpledgeRequest tl status not ok =", tlUnpledgeRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		unpledgeRes.Status = false
		unpledgeRes.Message = tlUnpledgeRes.Message
		ch <- unpledgeRes
		return
	}

	unpledgeRes.Data = tlUnpledgeRes.Data
	unpledgeRes.Message = tlUnpledgeRes.Message
	unpledgeRes.Result = tlUnpledgeRes.Result
	unpledgeRes.Status = true

	ch <- unpledgeRes
}

func (obj EpledgeObj) UnpledgeRequest(req models.UnpledgeReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + UNPLEDGEURL

	var unpledgeData []models.UnpledgeRes

	var wg sync.WaitGroup
	// Create a channel to receive results
	UnpledgeResChannel := make(chan models.UnpledgeRes, len(req.UnpledgeList))
	var tlUnPledgeReq TradeLabUnPledgeReq

	tlUnPledgeReq.ClientID = req.ClientID
	tlUnPledgeReq.TransactionType = UNPLEDGE
	if strings.ToUpper(req.Segment) == CAPITAL {
		tlUnPledgeReq.Segment = Capital
	}

	for i := 0; i < len(req.UnpledgeList); i++ {
		wg.Add(1)
		tlUnPledgeReq.Exchange = req.UnpledgeList[i].Exchange
		tlUnPledgeReq.Isin = req.UnpledgeList[i].Isin
		tlUnPledgeReq.Quantity = req.UnpledgeList[i].Quantity
		go CallUnpledge(tlUnPledgeReq, url, &wg, UnpledgeResChannel, reqH)
	}

	go func() {
		wg.Wait()
		close(UnpledgeResChannel)
	}()

	for UnPledgeResData := range UnpledgeResChannel {
		unpledgeData = append(unpledgeData, UnPledgeResData)
	}

	loggerconfig.Info("UnpledgeRequest resp=", unpledgeData, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	var apiRes apihelpers.APIRes
	apiRes.Data = unpledgeData
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	//Return the response immediately, asynchronously update db
	go func() {
		var batch []models.PledgeData
		batchSize := 100
		for i := 0; i < len(req.UnpledgeList); i++ {
			insertInfoInDB := models.PledgeData{
				ClientID:       req.ClientID,
				SegmentID:      req.Segment,
				Timestamp:      helpers.GetCurrentTimeInIST(),
				ISIN:           req.UnpledgeList[i].Isin,
				Quantity:       strconv.FormatInt(req.UnpledgeList[i].Quantity, 10),
				Price:          "-",
				Exchange:       req.UnpledgeList[i].Exchange,
				BOID:           "-",
				Depository:     "-",
				PledgeUnpledge: constants.UnPledge,
				DPID:           "-",
				PledgeTLS:      "-",
				ReqID:          "-",
				Version:        "-",
				Status:         constants.Started,
			}

			batch = append(batch, insertInfoInDB)

			if len(batch) >= batchSize {
				if err := db.GetPgObj().InsertPledgeDataBatch(batch); err != nil {
					loggerconfig.Error("UnpledgeRequest Batch insertion failed, error:", err, " requestId:", reqH.RequestId)
				}
				batch = nil
			}
		}

		if len(batch) > 0 {
			if err := db.GetPgObj().InsertPledgeDataBatch(batch); err != nil {
				loggerconfig.Error("UnpledgeRequest Final batch insertion failed, error:", err, " requestId:", reqH.RequestId)
			}
		}
	}()

	return http.StatusOK, apiRes
}

func (obj EpledgeObj) MTFEpledgeRequest(req models.MTFEPledgeRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + MTFEPLEDGE
	var reqPacket TradeLabMTFEPledgeRequest
	reqPacket.Depository = req.Depository
	reqPacket.ClientID = req.ClientID
	reqPacket.Exchange = req.Exchange
	reqPacket.BoID = req.BoID
	reqPacket.Segment = req.Segment
	reqPacket.RequestType = req.RequestType
	var isinPackets []IsinDetails
	for _, isin := range req.IsinDetails {
		var isinPacket IsinDetails
		isinPacket.IsinName = isin.IsinName
		isinPacket.Isin = isin.Isin
		isinPacket.Quantity = isin.Quantity
		isinPacket.Price = isin.Price
		isinPackets = append(isinPackets, isinPacket)
	}
	reqPacket.Order.Price = req.Order.Price
	reqPacket.Order.Device = req.Order.Device
	reqPacket.Order.Product = req.Order.Product
	reqPacket.Order.Exchange = req.Order.Exchange
	reqPacket.Order.Quantity = req.Order.Quantity
	reqPacket.Order.Validity = req.Order.Validity
	reqPacket.Order.ClientID = req.Order.ClientID
	reqPacket.Order.OrderSide = req.Order.OrderSide
	reqPacket.Order.OrderType = req.Order.OrderType
	reqPacket.Order.UserOrderID = req.Order.UserOrderID
	reqPacket.Order.InstrumentToken = req.Order.InstrumentToken
	reqPacket.Order.DisclosedQuantity = req.Order.DisclosedQuantity

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(reqPacket)

	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "MTFEpledge", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " MTFEpledgeRequest call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("MTFEpledgeRequest res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlMTFRes := TradelabMTFResponse{}
	json.Unmarshal([]byte(string(body)), &tlMTFRes)
	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " MTFEpledgeRequest tl status not ok =", tlMTFRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlMTFRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("MTFEpledgeRequest tl resp=", helpers.LogStructAsJSON(tlMTFRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = tlMTFRes.Data
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj EpledgeObj) GetPledgeList(reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes
	url := obj.tradeLabURL + MTFPLEDGELIST

	//make payload
	payload := new(bytes.Buffer)

	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "GetPledgeList", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " GetPledgeList call api error =", err, " uccId:", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("GetPledgeList res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlMTFRes := TLMtfPledgeListRes{}
	json.Unmarshal([]byte(string(body)), &tlMTFRes)
	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " GetPledgeList tl status not ok =", tlMTFRes.Message, " uccId:", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlMTFRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	var responseData models.TLMtfPledgeListRes
	var listData []models.MTFPledgeList
	for index := 0; index < len(tlMTFRes.Data.List); index++ {
		var lData models.MTFPledgeList
		lData.ClientID = tlMTFRes.Data.List[index].ClientID
		lData.Isin = tlMTFRes.Data.List[index].Isin
		lData.PledgeQuantity = tlMTFRes.Data.List[index].PledgeQuantity
		lData.ToBePledgedQuantity = tlMTFRes.Data.List[index].ToBePledgedQuantity
		lData.Segment = tlMTFRes.Data.List[index].Segment
		lData.Symbol = tlMTFRes.Data.List[index].Symbol
		lData.MtfSettlementDate = tlMTFRes.Data.List[index].MtfSettlementDate
		lData.MtfSquareOffDate = tlMTFRes.Data.List[index].MtfSquareOffDate
		lData.NseToken = tlMTFRes.Data.List[index].NseToken
		lData.BseToken = tlMTFRes.Data.List[index].BseToken
		lData.AvgPrice = tlMTFRes.Data.List[index].AvgPrice
		lData.MarginMultiplier = tlMTFRes.Data.List[index].MarginMultiplier
		lData.MarginVarElm = tlMTFRes.Data.List[index].MarginVarElm
		lData.MarginValue = tlMTFRes.Data.List[index].MarginValue
		lData.DaysTillSquareoff = tlMTFRes.Data.List[index].DaysTillSquareoff
		lData.IsLastDayOfMtf = tlMTFRes.Data.List[index].IsLastDayOfMtf
		lData.IsCfObligation = tlMTFRes.Data.List[index].IsCfObligation
		lData.CreatedAt = tlMTFRes.Data.List[index].CreatedAt
		lData.UpdatedAt = tlMTFRes.Data.List[index].UpdatedAt

		listData = append(listData, lData)
	}
	responseData.List = listData

	loggerconfig.Info("GetPledgeList tl resp=", helpers.LogStructAsJSON(responseData), " uccId:", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = responseData
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj EpledgeObj) GetCTDQuantityList(req models.MTFCTDDataReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes
	url := obj.tradeLabURL + CTDQUANTITYLIST + "?page_size=" + req.PageSize + "&page_no=" + req.PageNo + "&sort_order=" + req.SortOrder + "&first_page=" + req.FirstPage + "&client_id=" + strings.ToUpper(req.ClientId)

	//make payload
	payload := new(bytes.Buffer)

	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodGet, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "GetCTDQuantityList", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " GetCTDQuantityList call api error =", err, " uccId:", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("GetCTDQuantityList res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", reqH.ClientId, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlMTFCTDListRes := TLMTFCTDListRes{}
	json.Unmarshal([]byte(string(body)), &tlMTFCTDListRes)
	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " GetCTDQuantityList tl status not ok =", tlMTFCTDListRes.Message, " uccId:", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlMTFCTDListRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	var responseData models.MTFCTDDataRes
	responseData.TotalCount = tlMTFCTDListRes.Data.TotalCount
	var listData []models.CTDList
	for index := 0; index < len(tlMTFCTDListRes.Data.List); index++ {
		var lData models.CTDList
		lData.ClientID = tlMTFCTDListRes.Data.List[index].ClientID
		lData.Isin = tlMTFCTDListRes.Data.List[index].Isin
		lData.TotalPledgeQuantity = tlMTFCTDListRes.Data.List[index].TotalPledgeQuantity
		lData.CtdQuantity = tlMTFCTDListRes.Data.List[index].CtdQuantity
		lData.Symbol = tlMTFCTDListRes.Data.List[index].Symbol
		lData.AvgPrice = tlMTFCTDListRes.Data.List[index].AvgPrice
		lData.MarginMultiplier = tlMTFCTDListRes.Data.List[index].MarginMultiplier
		lData.CtdMarginValue = tlMTFCTDListRes.Data.List[index].CtdMarginValue
		lData.Token = tlMTFCTDListRes.Data.List[index].Token
		lData.Exchange = tlMTFCTDListRes.Data.List[index].Exchange
		lData.CreatedAt = tlMTFCTDListRes.Data.List[index].CreatedAt
		lData.UpdatedAt = tlMTFCTDListRes.Data.List[index].UpdatedAt
		lData.EdisApprovedQuantity = tlMTFCTDListRes.Data.List[index].EdisApprovedQuantity
		lData.ObligationQuantity = tlMTFCTDListRes.Data.List[index].ObligationQuantity
		lData.UsedQuantity = tlMTFCTDListRes.Data.List[index].UsedQuantity
		lData.LoginID = tlMTFCTDListRes.Data.List[index].LoginID
		lData.MarginValue = tlMTFCTDListRes.Data.List[index].MarginValue
		lData.TotalInvestedAmount = tlMTFCTDListRes.Data.List[index].TotalInvestedAmount
		lData.BrokerAmount = tlMTFCTDListRes.Data.List[index].BrokerAmount

		listData = append(listData, lData)
	}
	responseData.List = listData

	loggerconfig.Info("GetCTDQuantityList tl resp=", helpers.LogStructAsJSON(responseData), " uccId:", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = responseData
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj EpledgeObj) GetPledgeTransactions(req models.FetchEpledgeTxnReq, requestH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	pledgeData, err := db.GetPgObj().FetchPledgeTxnPaginated(req, constants.PledgeTxnPageSize)
	if err != nil {
		loggerconfig.Error("GetPledgeTransactions, Error while fetching data:", err)
		return apihelpers.SendInternalServerError()
	}

	apiRes.Data = pledgeData
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj EpledgeObj) MTFCTD(req models.MTFCTDReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + MTFCTD
	var reqPacket TLMTFCTDReq
	reqPacket.ClientID = req.ClientID
	reqPacket.LoginID = req.LoginID
	reqPacket.UserType = req.UserType

	var reqmtfctdvalues []MtfCtdValues
	for i := 0; i < len(req.MtfCtdValues); i++ {
		var reqmtfctdvalue MtfCtdValues
		reqmtfctdvalue.CtdQuantity = req.MtfCtdValues[i].CtdQuantity
		reqmtfctdvalue.Isin = req.MtfCtdValues[i].Isin
		reqmtfctdvalues = append(reqmtfctdvalues, reqmtfctdvalue)
	}
	reqPacket.MtfCtdValues = reqmtfctdvalues

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(reqPacket)

	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "MTFCTD", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " MTFCTD call api error =", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("MTFCTD res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlMTFCTDRes := TLMTFCTDRes{}
	json.Unmarshal([]byte(string(body)), &tlMTFCTDRes)
	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:P0-Critical, platform:", reqH.Platform, " MTFCTD tl status not ok =", tlMTFCTDRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlMTFCTDRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	loggerconfig.Info("MTFCTD tl resp=", helpers.LogStructAsJSON(tlMTFCTDRes), " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = tlMTFCTDRes.Data
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}
