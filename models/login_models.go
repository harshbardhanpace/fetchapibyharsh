package models

// LoginRequest login request
type LoginRequest struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required" mask:"id"`
}

type LoginByEmailRequest struct {
	EmailId  string `json:"emailId" binding:"required"`
	Password string `json:"password" binding:"required" mask:"id"`
}

type TwoFaQuestions struct {
	Question   string `json:"question"`
	QuestionID int    `json:"questionId"`
}

type TwoFaDetails struct {
	Questions  []TwoFaQuestions `json:"questions"`
	TwofaToken string           `json:"twoFaToken"`
	Type       string           `json:"type"`
}

// LoginResponse login response
type LoginResponse struct {
	Alert         string       `json:"alert"`
	AuthToken     string       `json:"authToken"`
	LoginID       string       `json:"loginId"`
	ResetPassword bool         `json:"resetPassword"`
	ResetTwoFa    bool         `json:"resetTwoFa"`
	Twofa         TwoFaDetails `json:"twofa"`
	TwofaEnabled  bool         `json:"twoFaEnabled"`
}

type TwoQuestions struct {
	QuestionID string `json:"questionId" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
}

// ValidateTwoFARequest validate two fa request
type ValidateTwoFARequest struct {
	LoginID    string         `json:"loginId" binding:"required"`
	Twofa      []TwoQuestions `json:"twoFa"`
	TwofaToken string         `json:"twoFaToken" binding:"required"`
	Type       string         `json:"type" binding:"required"`
}

// ValidateTwoFAResponse validate two fa response
type ValidateTwoFAResponse struct {
	AuthToken     string `json:"authToken"`
	ResetPassword bool   `json:"resetPassword"`
	ResetTwoFa    bool   `json:"resetTwoFa"`
}

// SetTwoFAPinRequest set two fa pin request
type SetTwoFAPinRequest struct {
	LoginID   string `json:"loginId"`
	Pin       string `json:"pin" mask:"id"`
	TwofaType string `json:"twoFaType"`
}

type ForgotPasswordRequest struct {
	LoginID string `json:"loginId"`
	Pan     string `json:"pan" mask:"id"`
}

type SetPasswordRequest struct {
	NewPass string `json:"newPass" binding:"required" mask:"id"`
	OldPass string `json:"oldPass" binding:"required" mask:"id"`
}

type ValidateTokenRequest struct {
	Token  string `json:"token"`
	UserId string `json:"userId"`
}

type ForgetResetTwoFaRequest struct {
	ClientID string `json:"clientId"`
	Pan      string `json:"pan" mask:"id"`
}

type ForgetResetEmailRequest struct {
	EmailId string `json:"emailId" validate:"required"`
	Pan     string `json:"pan" mask:"id" validate:"required"`
}

type LoginV2Request struct {
	ID     string `json:"id" binding:"required"`
	Secret string `json:"secret" binding:"required" mask:"id"`
}

type LoginV2Response struct {
	Alert          string       `json:"alert"`
	AuthToken      string       `json:"authToken"`
	CheckPan       bool         `json:"checkPan"`
	LoginID        string       `json:"loginId"`
	Name           string       `json:"name"`
	ReferenceToken string       `json:"referenceToken"`
	ResetPassword  bool         `json:"resetPassword"`
	ResetTwoFa     bool         `json:"resetTwoFa"`
	Twofa          TwoFaDetails `json:"twofa"`
	TwofaEnabled   bool         `json:"twofaEnabled"`
	KycUserId      string       `json:"kycUserId"`
}

type ValidateTwofaV2Req struct {
	LoginID    string         `json:"loginId"`
	Twofa      []TwoQuestions `json:"twofa"`
	TwofaToken string         `json:"twofaToken"`
	Type       string         `json:"type"`
	DeviceType string         `json:"deviceType"`
}

type ValidateTwofaV2Res struct {
	AuthToken     string `json:"authToken"`
	ResetPassword bool   `json:"resetPassword"`
	ResetTwoFa    bool   `json:"resetTwoFa"`
}

type SetupTotpV2Req struct {
	ClientID string `json:"clientId"`
}

type SetupTotpV2Res struct {
	ClientID string `json:"clientId"`
	Token    string `json:"token"`
}

type ChooseTwofaV2Req struct {
	LoginID   string `json:"loginId"`
	TwofaType string `json:"twofaType"`
	Totp      string `json:"totp" mask:"password"`
}

type ForgetTotpV2Req struct {
	LoginID string `json:"loginId"`
	Pan     string `json:"pan" mask:"id"`
}

type ValidateLoginOtpV2Req struct {
	ReferenceToken string `json:"referenceToken"`
	Otp            string `json:"otp"`
}

type SetupBiometricV2Req struct {
	ClientID    string `json:"clientId"`
	FingerPrint string `json:"fingerPrint" mask:"id"`
}

type DisableBiometricV2Req struct {
	ClientID    string `json:"clientId"`
	FingerPrint string `json:"fingerPrint" mask:"id"`
}

type Auth struct {
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	ErrorCode int         `json:"error_code"`
	Data      interface{} `json:"data"`
}

type GuestUserStatusReq struct {
	Email string `json:"email"`
}

type GuestUserStatusRes struct {
	GuestStatus int    `json:"guestStatus"`
	UserId      string `json:"userId"`
}

type TradingUserInfoData struct {
	Userid       string `json:"userId"`
	UserName     string `json:"username"`
	Branchind    string `json:"branchInd"`
	Introdate    string `json:"introdate"`
	Rmtlcode     string `json:"rmtlcode"`
	Delalercode  string `json:"delalercode"`
	Subbranchind string `json:"subbranchind"`
	Dob          string `json:"dob" mask:"id"`
	Mobilenos    string `json:"mobilenos" mask:"id"`
	Phnos        string `json:"phnos" mask:"id"`
	Emailno      string `json:"emailno"`
	Itaxno       string `json:"itaxno"`
	ParAdd1      string `json:"parAdd1" mask:"id"`
	City         string `json:"city" mask:"id"`
	Lstate       string `json:"lstate" mask:"id"`
	NseBseNfo    string `json:"nseBseNfo"`
	Dpcode       string `json:"dpcode"`
	Dpaccountno  string `json:"dpaccountno"`
	Poaflag      string `json:"poaflag"`
}

type MongoSignup struct {
	EmailId           string `json:"emailId"`
	MobileNo          string `json:"mobileNo"`
	MobileOtpVerified bool   `json:"mobileVerified"`
	EmailVerified     bool   `json:"emailVerified"`
	UserId            string `json:"userid"`
	AdminPassword     string `json:"adminPassword"`
	UserType          string `json:"userType"`
	Name              string `json:"name"`
	KYCStartTime      string `json:"kycStartTime"`
	LastModified      string `json:"lastModified"`
	AdminModifiedAt   string `json:"adminModifiedAt"`
	ReferralCode      string `json:"referralCode"`
	KycStatus         string `json:"kycStatus"`
	CreatedBy         string `json:"createdBy"`
}

type LoginWithQRReq struct {
	WebsocketID string `json:"websocketID"`
}

type ForgotPasswordV2Request struct {
	ClientID string `json:"ClientId"`
	EmailID  string `json:"emailId"`
	Pan      string `json:"pan" mask:"id"`
}

type CreateAppReq struct {
	AppName      string   `json:"appName"`
	RedirectUris []string `json:"redirectUris"`
	Scope        string   `json:"scope"`
	GrantTypes   []string `json:"grantTypes"`
	Owner        string   `json:"owner"`
}

type CreateAppRes struct {
	AppOwner string       `json:"appOwner"`
	Apps     []AppDetails `json:"apps"`
}

type AppDetails struct {
	AppID              string   `json:"appId" bson:"appId"`
	AppName            string   `json:"appName" bson:"appName"`
	AppSecret          string   `json:"appSecret" bson:"appSecret"`
	AppSecretExpiresAt int      `json:"appSecretExpiresAt" bson:"appSecretExpiresAt"`
	GrantTypes         []string `json:"grantTypes" bson:"grantTypes"`
	RedirectUris       []string `json:"redirectUris" bson:"redirectUris"`
	Scope              string   `json:"scope" bson:"scope"`
	State              string   `json:"state" bson:"state"`
	AuthCode           string   `json:"authCode" bson:"authCode"`
	AccessToken        string   `json:"accessToken" bson:"accessToken"`
	ExpiryTime         string   `json:"expiryTime" bson:"expiryTime"`
}

type GetAccessTokenReq struct {
	AppState string `json:"state"`
}

type GetAccessTokenV2Req struct {
	AppState    string `json:"state"`
	AccessToken string `json:accessToken`
}

type LoginByEmailOtpReq struct {
	Email string `json:"email" validate:"required"`
}

type LoginByEmailOtpRes struct {
	Alert          string                 `json:"alert"`
	AuthToken      string                 `json:"authToken"`
	CheckPan       bool                   `json:"checkPan"`
	LoginID        string                 `json:"loginId"`
	Name           string                 `json:"name"`
	ReferenceToken string                 `json:"referenceToken"`
	ResetPassword  bool                   `json:"resetPassword"`
	ResetTwoFA     bool                   `json:"resetTwoFA"`
	TwoFA          map[string]interface{} `json:"twoFA"`
	TwoFAEnabled   bool                   `json:"twoFAEnabled"`
}
