package tradelab

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/loggerconfig"
	"space/models"
	"strconv"
	"time"
)

type EdisObj struct {
	tradeLabURL string
}

func InitEdisProvider() EdisObj {
	defer models.HandlePanic()

	edisObj := EdisObj{
		tradeLabURL: constants.TLURL,
	}

	return edisObj
}

func (obj EdisObj) SendEdisRequest(req models.EdisReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + EDISREQ

	var tlEdisReq TradeLabEdisReq
	tlEdisReq.ClientID = req.ClientID

	for _, instrument := range req.Instruments {
		tradeLabInstrument := TradeLabInstrument{
			InstrumentToken: instrument.InstrumentToken,
			Exchange:        instrument.Exchange,
			Total:           instrument.Total,
			Authorized:      instrument.Authorized,
		}
		tlEdisReq.Instruments = append(tlEdisReq.Instruments, tradeLabInstrument)
	}

	tlEdisReq.RequestType = req.RequestType

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlEdisReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "SendEdisRequest", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " SendEdisRequest call api error", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("SendEdisRequest res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlEdisRes := TradeLabEdisRes{}
	json.Unmarshal([]byte(string(body)), &tlEdisRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " SendEdisRequest status not ok =", tlEdisRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlEdisRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var edisRes models.EdisRes
	edisRes.Data.Depository = tlEdisRes.Data.Depository
	edisRes.Data.DpId = tlEdisRes.Data.DpId
	edisRes.Data.EncryptedDtls = tlEdisRes.Data.EncryptedDtls
	edisRes.Data.RequestId = tlEdisRes.Data.RequestId
	edisRes.Data.Version = tlEdisRes.Data.Version
	edisRes.Html = tlEdisRes.Html
	edisRes.Message = tlEdisRes.Message
	edisRes.Status = tlEdisRes.Status

	loggerconfig.Info("SendEdisRequest resp=", edisRes, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = edisRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj EdisObj) GenerateTpin(req models.TpinReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + TPINREQ

	var tlTpinReq TradeLabTpinReq

	tlTpinReq.Boid = req.Boid
	tlTpinReq.Pan = req.Pan
	tlTpinReq.ReqFlag = req.ReqFlag
	tlTpinReq.ReqTime = req.ReqTime

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlTpinReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallAPIFunc(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "GenerateTpin", duration, req.ClientID, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " GenerateTpin call api error", err, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("GenerateTpin res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", req.ClientID, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlTpinRes := TradeLabTpinRes{}
	json.Unmarshal([]byte(string(body)), &tlTpinRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " GenerateTpin status not ok =", tlTpinRes.Message, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)
		apiRes.Message = tlTpinRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	//fill up controller response
	var tpinRes models.TpinRes
	tpinRes.TlRespMessage = tlTpinRes.Message

	loggerconfig.Info("GenerateTpin resp=", tpinRes, " uccId:", req.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId, " clientVersion:", reqH.ClientVersion)

	apiRes.Data = tpinRes
	if tlTpinRes.Success {
		apiRes.Message = "SUCCESS"
		apiRes.Status = true
		return http.StatusOK, apiRes
	} else {
		apiRes.Message = "FAILURE"
		apiRes.Status = false
		return http.StatusOK, apiRes
	}

}
