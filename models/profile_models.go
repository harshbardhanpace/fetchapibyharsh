package models

// Profile Request
type ProfileRequest struct {
	ClientID string `json:"clientId" binding:"required" example:"Client1"`
}

// Profile Response
type ProfileResponse struct {
	Data ProfileResponseData `json:"data"`
}

type SendAFOtpReq struct {
	ClientID string `json:"clientId" binding:"required" example:"Client1"`
}

type SendEmailOtp struct {
	Otp            string    `json:"otp"`
	RecipientEmail string    `json:"recipientEmail"`
	RecipientName  string    `json:"recipientName"`
	IsEmail        bool      `json:"isEmail"`
	ReqHeader      ReqHeader `json:"reqHeader"`
}

type VerifyAFOtpReq struct {
	ClientID string `json:"clientId" binding:"required" example:"Client1"`
	Otp      string `json:"otp"`
}

type AccountFreezeReq struct {
	ClientID string `json:"clientId" binding:"required" example:"Client1"`
}

type AccountFreezeRes struct {
	Data struct {
	} `json:"data"`
	Message string `json:"message"`
}

type MSGSendSmsRequest struct {
	FlowID   string `json:"flow_id"`
	Sender   string `json:"sender"`
	ShortURL string `json:"short_url"`
	Mobiles  string `json:"mobiles"`
	Var      string `json:"var"`
}

type RoleDetails struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProfileDataResp struct {
	ClientID  string `json:"clientId" mask:"id"`
	Name      string `json:"name"`
	EmailID   string `json:"emailId"`
	BoID      string `json:"boId"`
	PanNumber string `json:"panNumber" mask:"id"`
}

type ExchangeNnf struct {
	Bse int `json:"BSE"`
	Mcx int `json:"MCX"`
	Nfo int `json:"NFO"`
	Nse int `json:"NSE"`
}

type ProfileResponseData struct {
	Branch                 string          `json:"branch"`
	BankBranchName         string          `json:"bankBranchName"`
	OfficeAddr             string          `json:"officeAddress"`
	DpID                   []string        `json:"dpId"`
	City                   string          `json:"city"`
	PermanentAddr          string          `json:"permanentAddress"`
	BankName               string          `json:"bankName"`
	BankAccountNumber      string          `json:"bankAccountNumber"`
	PanNumber              string          `json:"panNumber"`
	Role                   RoleDetails     `json:"role"`
	EmailID                string          `json:"emailId"`
	BrokerID               string          `json:"brokerId"`
	ClientID               string          `json:"clientId"`
	BankState              string          `json:"bankState"`
	AccountType            string          `json:"accountType"`
	Status                 string          `json:"status"`
	UserType               string          `json:"userType"`
	LastPasswordChangeDate int             `json:"lastPasswordChangeDate"`
	BoID                   []string        `json:"boId"`
	BasketEnabled          bool            `json:"basketEnabled"`
	TwofaEnabled           bool            `json:"twofaEnabled"`
	Name                   string          `json:"name"`
	Depository             string          `json:"depository"`
	ExchangeNnf            ExchangeNnf     `json:"exchangeNnf"`
	PoaStatus              bool            `json:"poaStatus"`
	BankCity               string          `json:"bankCity"`
	IfscCode               string          `json:"ifscCode"`
	Dob                    string          `json:"dob"`
	ExchangesSubscribed    []string        `json:"exchangesSubscribed"`
	Sex                    string          `json:"sex"`
	PoaEnabled             bool            `json:"poaEnabled"`
	BackofficeLink         string          `json:"backofficeLink"`
	State                  string          `json:"state"`
	PhoneNumber            string          `json:"phoneNumber"`
	ProductsEnabled        []string        `json:"productsEnabled"`
	ProfileURL             string          `json:"profileUrl"`
	SegmentDetails         ExchangesActive `json:"segmentDetails"`
}

type ExchangesActive struct {
	NSE string `json:"nse"`
	BSE string `json:"bse"`
	NFO string `json:"nfo"`
	BFO string `json:"bfo"`
	MCX string `json:"mcx"`
}
