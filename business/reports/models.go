package reports

type ShilpiGetBillDetailsCdsl []struct {
	Charges        string `json:"charges"`
	TrxDatee       string `json:"TrxDatee"`
	Qty            string `json:"qty"`
	Gst            string `json:"gst"`
	InstrumentName string `json:"instrument name"`
	Isincode       string `json:"isincode"`
	ChargesDetails string `json:"charges details"`
	TotalCharges   string `json:"total charges"`
}

type SendEmailLedger struct {
	ClientId          string `json:"clientId"`
	EncodedReportFile string `json:"encodedReportFile"`
	DateFrom          string `json:"dateFrom"`
	DateTo            string `json:"dateTo"`
	RecipientEmail    string `json:"recipientEmail"`
	ApplicantName     string `json:"applicantName"`
}

type SendEmailCommodityTradebook struct {
	ClientId          string `json:"clientId"`
	EncodedReportFile string `json:"encodedReportFile"`
	DateFrom          string `json:"dateFrom"`
	DateTo            string `json:"dateTo"`
	RecipientEmail    string `json:"recipientEmail"`
	ApplicantName     string `json:"applicantName"`
}

type SendEmailFnoTradebook struct {
	ClientId          string `json:"clientId"`
	EncodedReportFile string `json:"encodedReportFile"`
	DateFrom          string `json:"dateFrom"`
	DateTo            string `json:"dateTo"`
	RecipientEmail    string `json:"recipientEmail"`
	ApplicantName     string `json:"applicantName"`
}

type SendEmailDpCharges struct {
	ClientId          string `json:"clientId"`
	EncodedReportFile string `json:"encodedReportFile"`
	DateFrom          string `json:"dateFrom"`
	DateTo            string `json:"dateTo"`
	RecipientEmail    string `json:"recipientEmail"`
	ApplicantName     string `json:"applicantName"`
}

type SendEmailHoldingFinancial struct {
	ClientId          string `json:"clientId"`
	EncodedReportFile string `json:"encodedReportFile"`
	RecipientEmail    string `json:"recipientEmail"`
	ApplicantName     string `json:"applicantName"`
}
