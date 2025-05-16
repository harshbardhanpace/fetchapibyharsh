package freshdesk

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/loggerconfig"
	"space/models"
)

type FreshdeskObj struct {
	freshdeskbaseUrl string
	freshdeskApiKey  string
	freshdeskPass    string
}

func InitFreshdeskProvider() FreshdeskObj {
	defer models.HandlePanic()

	freshdesk := FreshdeskObj{
		freshdeskbaseUrl: constants.FreshDeskBaseUrl,
		freshdeskApiKey:  constants.FreshDeskApiKey,
		freshdeskPass:    constants.FreshDeskPass,
	}

	return freshdesk
}

func (obj FreshdeskObj) CreateFreshdeskTicket(ticket models.FreshdeskTicketReq, requestH models.ReqHeader) (int, apihelpers.APIRes) {
	var res apihelpers.APIRes

	// Convert ticket object to JSON
	ticketData, err := json.Marshal(ticket)
	if err != nil {
		loggerconfig.Error("CreateFreshdeskTicketService, Failed to marshal ticket, platform:", requestH.Platform, "err:", err, "clientId:", requestH.ClientId, "requestId:", requestH.RequestId)
		res.Status = false
		res.Message = constants.InvalidRequest
		return http.StatusBadRequest, res
	}

	// Create HTTP request to Freshdesk API
	url := obj.freshdeskbaseUrl + constants.FreshdeskTicketURL
	payload := bytes.NewBuffer(ticketData)

	response, err := apihelpers.CallApiFreshDesk(http.MethodPost, url, payload, obj.freshdeskApiKey, obj.freshdeskPass)
	if err != nil {
		if err != nil {
			loggerconfig.Error("CreateFreshdeskTicketService, API call to Freshdesk failed, platform:", requestH.Platform, "err:", err, "clientId:", requestH.ClientId, "requestId:", requestH.RequestId)
			return apihelpers.SendInternalServerError()
		}
	}
	defer response.Body.Close()

	// Handle Freshdesk API response
	if response.StatusCode == http.StatusCreated {
		var freshdeskResponse models.FreshdeskTicketResponse
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			loggerconfig.Error("Alert Severity:P1-High, CreateFreshdeskTicketService, Failed to read response body, platform:", requestH.Platform, "err:", err, "clientId:", requestH.ClientId, "requestId:", requestH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		err = json.Unmarshal(body, &freshdeskResponse)
		if err != nil {
			loggerconfig.Error("Alert Severity:P1-High, CreateFreshdeskTicketService, Failed to unmarshal Freshdesk response, platform:", requestH.Platform, "err:", err, "clientId:", requestH.ClientId, "requestId:", requestH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		// Return the response struct and success code
		res.Status = true
		res.Data = freshdeskResponse
		res.Message = "Ticket created Successfully"
	} else if response.StatusCode == http.StatusUnauthorized {
		loggerconfig.Error("Alert Severity:P0-Critical, CreateFreshdeskTicketService, Unauthorized, invalid Freshdesk credentials, status:", response.Status, "platform:", requestH.Platform, "err:", err, "clientId:", requestH.ClientId, "requestId:", requestH.RequestId)
		return apihelpers.SendInternalServerError()
	} else {
		loggerconfig.Error("Failed to create ticket, status code:", response.StatusCode, "platform:", requestH.Platform, "clientId:", requestH.ClientId, "requestId:", requestH.RequestId)
		res.Status = false
		res.Data = nil
		res.Message = "failed to create ticket"
		return response.StatusCode, res
	}

	return http.StatusOK, res
}
