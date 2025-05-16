package models

type MongoClientDetailsChange struct {
	ClientId                        string                    `json:"clientId"`
	EmailId                         string                    `json:"emailId"`
	MobileNo                        string                    `json:"mobileNo"`
	EmailIdVerified                 bool                      `json:"emailIdVerified"`
	MobileNoVerified                bool                      `json:"mobileNoVerified"`
	IPV                             MongoUserVideoDetails     `json:"ipv"`
	Esign                           MongoESignDetails         `json:"esign"`
	BankAccounts                    []MongoBankAccountDetails `json:"bankaccounts"`
	NomineeDetails                  NomineeV2                 `json:"nomineeDetails"`
	SegmentDetails                  MongoSegmentDetailsV2     `json:"segmentDetails"`
	TransactionId                   string                    `json:"transactionId"`
	TransactionStatus               string                    `json:"transactionStatus"`
	TransactionType                 string                    `json:"transactionType"`
	TransactionStartTimestamp       string                    `json:"transactionStartTimestamp"`
	TransactionLastUpdatedTimestamp string                    `json:"transactionLastUpdatedTimestamp"`
	NomineeDetailsMultiple          []NomineeV2               `json:"nomineeDetailsMultiple"`
}

type NomineeV2 struct {
	NomineeDob              string          `json:"nomineeDob" mask:"id"`
	NomineeName             string          `json:"nomineeName"`
	NomineeRelationship     string          `json:"nomineeRelationship"`
	PercentageShare         string          `json:"percentageShare"`
	NomineeIdType           string          `json:"nomineeIdType"`
	NomineeId               string          `json:"nomineeId" mask:"id"`
	NomineeAddress          string          `json:"nomineeAddress"`
	NomineeMobile           string          `json:"nomineeMobile"`
	NomineeEmail            string          `json:"nomineeEmail"`
	NomineeGuardianForMinor GuardianDetails `json:"nomineeGuardianForMinor"`
}

type GuardianDetails struct {
	Name                    string `json:"name"`
	DOB                     string `json:"dob"`
	Address                 string `json:"address"`
	MobileNumber            string `json:"mobileNumber"`
	RelationshipWithNominee string `json:"relationshipWithNominee"`
	IDType                  string `json:"idType"`
	IDNumber                string `json:"idNumber"`
}

type MongoSegmentDetailsV2 struct {
	UserId                   string `json:"userid"`
	Equity                   bool   `json:"equity"`
	FutureAndOptions         bool   `json:"futureAndOptions"`
	Commodities              bool   `json:"commodities"`
	LastUpdatedUnixTimestamp int64  `json:"lastUpdatedUnixTimestamp"`
	CreatedAtUnixTimestamp   int64  `json:"createdAtUnixTimestamp"`
	ProofType                string `json:"proofType"`
	ProofLocation            string `json:"proofLocation"`
	ProofDocumentPassword    string `json:"proofDocumentPassword"`
}

type MongoESignDetails struct {
	UserId                   string `json:"userid"`
	EsignStatus              string `json:"eSignStatus"`
	DocId                    string `json:"documentId"`
	CallBack                 string `json:"callBack"`
	LastUpdatedUnixTimestamp int64  `json:"lastUpdatedUnixTimestamp"`
	CreatedAtUnixTimestamp   int64  `json:"createdAtUnixTimestamp"`
}

type MongoUserVideoDetails struct {
	UserId                   string `json:"userid"`
	UserVideoVerified        bool   `json:"userVideoVerified"`
	UserVideoS3Location      string `json:"userVideoS3Location"`
	OtpProvidedToUser        string `json:"otpProvidedToUser"`
	OtpExtractedFromVideo    string `json:"otpExtractedFromVideo"`
	Verified                 bool   `json:"verified"`
	Rejection                string `json:"Rejection"`
	LastUpdatedUnixTimestamp int64  `json:"lastUpdatedUnixTimestamp"`
	CreatedAtUnixTimestamp   int64  `json:"createdAtUnixTimestamp"`
}
