package models

import "time"

// Fetch Funds Request
type FetchFundsRequest struct {
	ClientID string `json:"clientId"`
	Type     string `json:"type"`
}

// Fetch Funds Response
type FetchFundsResponse struct {
	ClientID string                     `json:"clientId"`
	Headers  []string                   `json:"headers"`
	Values   []FetchFundsResponseValues `json:"values"`
}

type FetchFundsResponseValues struct {
	Num0 string `json:"0"`
	Num1 string `json:"1"`
}

type CancelPayoutReq struct {
	Transactions []string `json:"transactions"`
	ClientID     string   `json:"clientId"`
	Status       string   `json:"status"`
}

type CancelPayoutReqV3 struct {
	TransactionId string `json:"transactionsId"`
	ClientID      string `json:"clientId"`
}

type AtomPayoutRequest struct {
	Amount        string `json:"amount" validate:"required,numeric"`
	ClientID      string `json:"clientId" validate:"required"`
	Ifsc          string `json:"ifsc" validate:"required"`
	AccountNumber string `json:"accountNumber" validate:"required"`
	BankName      string `json:"bankName" validate:"required"`
}

type AtomPayoutResponse struct {
	Data interface{} `json:"data"`
}

type ClientTransactionsRequest struct {
	ClientID string `json:"clientId"`
}

type ClientTransactionsResponse struct {
	Data    []ClientTransactionsResponseData `json:"data" mask:"struct"`
	Message string                           `json:"message"`
	Status  string                           `json:"status"`
}

type ClientTransactionsResponseData struct {
	AccountName                 string                                          `json:"accountName"`
	Amount                      string                                          `json:"amount"`
	BankName                    string                                          `json:"bankName"`
	BankTransactionID           string                                          `json:"bankTransactionId" mask:"id"`
	ClientID                    string                                          `json:"clientId"`
	CreatedAt                   string                                          `json:"createdAt"`
	Ifsc                        string                                          `json:"ifsc" mask:"id"`
	MerchantTransactionID       string                                          `json:"merchantTransactionId" mask:"id"`
	PaymentGatewayTransactionID string                                          `json:"paymentGatewayTransactionId" mask:"id"`
	PaymentGatewayUsername      string                                          `json:"paymentGatewayUsername" mask:"id"`
	PreviousBalance             int                                             `json:"previousBalance"`
	Status                      string                                          `json:"status"`
	StatusLifeCycle             []ClientTransactionsResponseDataStatusLifeCycle `json:"statusLifeCycle"`
	TransactionID               string                                          `json:"transactionId"`
	TransactionTimestamp        int                                             `json:"transactionTimestamp"`
	TransactionType             string                                          `json:"transactionType"`
	UpdatedAt                   string                                          `json:"updatedAt"`
	UpdatedBy                   string                                          `json:"updatedBy"`
	UserID                      string                                          `json:"userId"`
}

type ClientTransactionsResponseDataStatusLifeCycle struct {
	Status    string `json:"status"`
	UpdatedAt string `json:"updatedAt"`
	UpdatedBy string `json:"updatedBy"`
}

type PayoutRequest struct {
	Amount        string `json:"amount"`
	ClientID      string `json:"clientId"`
	Ifsc          string `json:"ifsc"`
	AccountNumber string `json:"accountNumber"`
	BankName      string `json:"bank_name"`
}

type PayoutDetails struct {
	Amount                 int64     `json:"amount"`
	ClientID               string    `json:"clientId"`
	Ifsc                   string    `json:"ifsc"`
	AccountNumber          string    `json:"accountNumber"`
	BankName               string    `json:"bankName"`
	DebitCredit            string    `json:"debitCredit"`
	TradelabFundsUpdated   bool      `json:"tradelabFundsUpdated"`
	BackofficeFundsUpdated bool      `json:"backofficeFundsUpdated"`
	TransactionType        string    `json:"transactionType"`
	TransactionId          string    `json:"transactionId"`
	TransactionStatus      string    `json:"transactionStatus"`
	Remarks                string    `json:"remarks"`
	CreateDate             time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
}

type PayoutQueueDetails struct {
	ClientID      string `json:"clientId"`
	Amount        string `json:"amount"`
	Dates         string `json:"dates"`
	BankName      string `json:"bankName"`
	TransactionID string `json:"transactionId"`
	Remark        string `json:"remark"`
	PayoutAmount  string `json:"payoutAmount"`
	IFSC          string `json:"ifsc"`
	AccountNumber string `json:"accountNumber"`
}
