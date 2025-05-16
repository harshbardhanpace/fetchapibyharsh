package models

import "time"

type FinvuLoginReq struct {
	Header struct {
		Rid       string `json:"rid"`
		Ts        string `json:"ts"`
		ChannelID string `json:"channelId"`
	} `json:"header"`
	Body struct {
		UserID   string `json:"userId"`
		Password string `json:"password"`
	} `json:"body"`
}

type FinvuLoginRes struct {
	Header struct {
		Rid       string    `json:"rid"`
		Ts        time.Time `json:"ts"`
		ChannelID string    `json:"channelId"`
	} `json:"header"`
	Body struct {
		Token string `json:"token"`
	} `json:"body"`
}

type FinvuGetBankStatementReq struct {
	ClientId string `json:"clientId"`
}

type CreateConsentRequestPlusReq struct {
	ClientId   string `json:"clientId"`
	CustomerId string `json:"customerId"`
}

type ConsentsRequestPlusReq struct {
	Header struct {
		Rid       string `json:"rid"`
		Ts        string `json:"ts"`
		ChannelID string `json:"channelId"`
	} `json:"header"`
	Body struct {
		CustID             string   `json:"custId"`
		ConsentDescription string   `json:"consentDescription"`
		TemplateName       string   `json:"templateName"`
		UserSessionID      string   `json:"userSessionId"`
		RedirectURL        string   `json:"redirectUrl"`
		Fip                []string `json:"fip"`
		ConsentDetails     struct {
			Customer struct {
				ID string `json:"id"`
			} `json:"Customer"`
			DataConsumer struct {
				ID string `json:"id"`
			} `json:"DataConsumer"`
			Purpose struct {
				Code     string `json:"code"`
				RefURI   string `json:"refUri"`
				Text     string `json:"text"`
				Category struct {
					Type string `json:"type"`
				} `json:"Category"`
			} `json:"Purpose"`
			ConsentMode  string   `json:"consentMode"`
			ConsentTypes []string `json:"consentTypes"`
			FiTypes      []string `json:"fiTypes"`
			FetchType    string   `json:"fetchType"`
			Frequency    struct {
				Value int    `json:"value"`
				Unit  string `json:"unit"`
			} `json:"Frequency"`
			DataLife struct {
				Value int    `json:"value"`
				Unit  string `json:"unit"`
			} `json:"DataLife"`
			ConsentStart  time.Time `json:"consentStart"`
			ConsentExpiry string    `json:"consentExpiry"`
			FIDataRange   struct {
				From time.Time `json:"from"`
				To   time.Time `json:"to"`
			} `json:"FIDataRange"`
		} `json:"ConsentDetails"`
		AaID string `json:"aaId"`
	} `json:"body"`
}

type ConsentsRequestPlusRes struct {
	Header struct {
		Rid       string `json:"rid"`
		Ts        string `json:"ts"`
		ChannelID string `json:"channelId"`
	} `json:"header"`
	Body struct {
		EncryptedRequest string `json:"encryptedRequest"`
		RequestDate      string `json:"requestDate"`
		EncryptedFiuID   string `json:"encryptedFiuId"`
		ConsentHandle    string `json:"ConsentHandle"`
		URL              string `json:"url"`
	} `json:"body"`
}

type ConsentsRequestPlusResFrontRes struct {
	URL string `json:"url"`
}

type FiRequestReq struct {
	Header struct {
		Rid       string `json:"rid"`
		Ts        string `json:"ts"`
		ChannelID string `json:"channelId"`
	} `json:"header"`
	Body struct {
		CustID            string    `json:"custId"`
		ConsentHandleID   string    `json:"consentHandleId"`
		ConsentID         string    `json:"consentId"`
		DateTimeRangeFrom time.Time `json:"dateTimeRangeFrom"`
		DateTimeRangeTo   time.Time `json:"dateTimeRangeTo"`
	} `json:"body"`
}

type FiRequestRes struct {
	Header struct {
		Rid       string `json:"rid"`
		Ts        string `json:"ts"`
		ChannelID string `json:"channelId"`
	} `json:"header"`
	Body struct {
		Ver             string    `json:"ver"`
		Timestamp       time.Time `json:"timestamp"`
		Txnid           string    `json:"txnid"`
		ConsentID       string    `json:"consentId"`
		SessionID       string    `json:"sessionId"`
		ConsentHandleID any       `json:"consentHandleId"`
	} `json:"body"`
}

type FiStatusRes struct {
	Header struct {
		Rid       string `json:"rid"`
		Ts        string `json:"ts"`
		ChannelID string `json:"channelId"`
	} `json:"header"`
	Body struct {
		FiRequestStatus string `json:"fiRequestStatus"`
	} `json:"body"`
}

type CheckConsentStatusRes struct {
	Header struct {
		Rid       string    `json:"rid"`
		Ts        time.Time `json:"ts"`
		ChannelID any       `json:"channelId"`
	} `json:"header"`
	Body struct {
		ConsentStatus string `json:"consentStatus"`
		ConsentID     string `json:"consentId"`
	} `json:"body"`
}

type GetConsentAStatusById struct {
	Header struct {
		Rid       string    `json:"rid"`
		Ts        time.Time `json:"ts"`
		ChannelID any       `json:"channelId"`
	} `json:"header"`
	Body struct {
		ConsentID       string `json:"consentId"`
		Status          string `json:"status"`
		CreateTimestamp string `json:"createTimestamp"`
		ConsentDetail   struct {
			ConsentStart  string   `json:"consentStart"`
			ConsentExpiry string   `json:"consentExpiry"`
			ConsentMode   string   `json:"consentMode"`
			FetchType     string   `json:"fetchType"`
			ConsentTypes  []string `json:"consentTypes"`
			FiTypes       []string `json:"fiTypes"`
			DataConsumer  struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"DataConsumer"`
			DataProvider struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"DataProvider"`
			Customer struct {
				ID string `json:"id"`
			} `json:"Customer"`
			Accounts []struct {
				FiType          string `json:"fiType"`
				FipID           string `json:"fipId"`
				AccType         string `json:"accType"`
				LinkRefNumber   string `json:"linkRefNumber"`
				MaskedAccNumber string `json:"maskedAccNumber"`
			} `json:"Accounts"`
			Purpose struct {
				Code     string `json:"code"`
				RefURI   string `json:"refUri"`
				Text     string `json:"text"`
				Category struct {
					Type string `json:"type"`
				} `json:"Category"`
			} `json:"Purpose"`
			FIDataRange struct {
				From time.Time `json:"from"` //time.Time
				To   time.Time `json:"to"`
			} `json:"FIDataRange"`
			DataLife struct {
				Unit  string `json:"unit"`
				Value int    `json:"value"`
			} `json:"DataLife"`
			Frequency struct {
				Unit  string `json:"unit"`
				Value int    `json:"value"`
			} `json:"Frequency"`
			DataFilter []struct {
				Type     string `json:"type"`
				Operator string `json:"operator"`
				Value    string `json:"value"`
			} `json:"DataFilter"`
		} `json:"ConsentDetail"`
		ConsentUse struct {
			LogURI          string `json:"logUri"`
			Count           int    `json:"count"`
			LastUseDateTime string `json:"lastUseDateTime"`
		} `json:"ConsentUse"`
	} `json:"body"`
}

type FinvuErrorRes struct {
	Header struct {
		Rid       string    `json:"rid"`
		Ts        time.Time `json:"ts"`
		ChannelID string    `json:"channelId"`
	} `json:"header"`
	Errors []struct {
		ErrorCode int    `json:"errorCode"`
		ErrorMsg  string `json:"errorMsg"`
	} `json:"errors"`
}
