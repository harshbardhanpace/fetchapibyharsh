package models

import (
	apihelpers "space/apiHelpers"

	"github.com/gin-gonic/gin"
)

type LoginProvider interface {
	LoginByPass(LoginRequest, ReqHeader) (int, apihelpers.APIRes)
	LoginByEmail(LoginByEmailRequest, ReqHeader) (int, apihelpers.APIRes)
	ValidateTwoFa(ValidateTwoFARequest, ReqHeader) (int, apihelpers.APIRes)
	SetTwoFaPin(SetTwoFAPinRequest, ReqHeader) (int, apihelpers.APIRes)
	ForgetPassword(ForgotPasswordRequest, ReqHeader) (int, apihelpers.APIRes)
	ForgetPasswordEmail(ForgetResetEmailRequest, ReqHeader) (int, apihelpers.APIRes)
	SetPassword(SetPasswordRequest, ReqHeader) (int, apihelpers.APIRes)
	ValidateToken(ValidateTokenRequest, ReqHeader) (int, apihelpers.APIRes)
	ForgetResetTwoFa(ForgetResetTwoFaRequest, ReqHeader) (int, apihelpers.APIRes)
	ForgetResetTwoFaEmail(ForgetResetEmailRequest, ReqHeader) (int, apihelpers.APIRes)
	GuestUserStatus(GuestUserStatusReq, ReqHeader) (int, apihelpers.APIRes)
	QRWebLogin(LoginWithQRReq, ReqHeader) (int, apihelpers.APIRes)
	UnblockUser(UnblockUserReq, ReqHeader) (int, apihelpers.APIRes)
	CreateApp(CreateAppReq, ReqHeader) (int, apihelpers.APIRes)
	FetchApps(string, ReqHeader) (int, apihelpers.APIRes)
	DeleteApp(string, ReqHeader) (int, apihelpers.APIRes)
	HandleAuthCode(string, string, ReqHeader) (int, apihelpers.APIRes)
	GetAccessToken(GetAccessTokenReq, ReqHeader) (int, apihelpers.APIRes)
	GenerateAccessToken(string, string, ReqHeader) (int, apihelpers.APIRes)
}

type LoginProviderV2 interface {
	LoginV2(LoginV2Request, ReqHeader) (int, apihelpers.APIRes)
	ValidateTwofaV2(ValidateTwofaV2Req, ReqHeader) (int, apihelpers.APIRes)
	SetupTotpV2(SetupTotpV2Req, ReqHeader) (int, apihelpers.APIRes)
	ChooseTwofaV2(ChooseTwofaV2Req, ReqHeader) (int, apihelpers.APIRes)
	ForgetTotpV2(ForgetTotpV2Req, ReqHeader) (int, apihelpers.APIRes)
	ValidateLoginOtpV2(ValidateLoginOtpV2Req, ReqHeader) (int, apihelpers.APIRes)
	SetupBiometricV2(SetupBiometricV2Req, ReqHeader) (int, apihelpers.APIRes)
	DisableBiometricV2(DisableBiometricV2Req, ReqHeader) (int, apihelpers.APIRes)
	ForgetPasswordV2(ForgotPasswordV2Request, ReqHeader) (int, apihelpers.APIRes)
	UnblockUserV2(UnblockUserV2Req, ReqHeader) (int, apihelpers.APIRes)
	GetAccessTokenV2(GetAccessTokenV2Req, ReqHeader) (int, apihelpers.APIRes)
}

type LoginProviderV3 interface {
	SetPassword(SetPasswordRequest, ReqHeader) (int, apihelpers.APIRes)
	ValidateToken(ValidateTokenRequest, ReqHeader) (int, apihelpers.APIRes)
	ForgetResetTwoFa(ForgetResetTwoFaRequest, ReqHeader) (int, apihelpers.APIRes)
	ValidateLoginOtpV2(ValidateLoginOtpV2Req, ReqHeader) (int, apihelpers.APIRes)
	SetupBiometricV2(SetupBiometricV2Req, ReqHeader) (int, apihelpers.APIRes)
	DisableBiometricV2(DisableBiometricV2Req, ReqHeader) (int, apihelpers.APIRes)
	LoginByEmailOtp(LoginByEmailOtpReq, ReqHeader) (int, apihelpers.APIRes)
}

type OrderProvider interface {
	PlaceOrder(PlaceOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	ModifyOrder(ModifyOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	CancelOrder(CancelOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	PlaceAMOOrder(PlaceOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	ModifyAMOOrder(ModifyAMORequest, ReqHeader) (int, apihelpers.APIRes)
	CancelAMOOrder(CancelOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	PendingOrder(PendingOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	CompletedOrder(CompletedOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	TradeBook(TradeBookRequest, ReqHeader) (int, apihelpers.APIRes)
	OrderHistory(OrderHistoryRequest, ReqHeader) (int, apihelpers.APIRes)
	PlaceGTTOrder(CreateGTTOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	ModifyGTTOrder(ModifyGTTOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	CancelGTTOrder(CancelGTTOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	FetchGTTOrder(FetchGTTOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	PlaceGttOCOOrder(CreateGttOCORequest, ReqHeader) (int, apihelpers.APIRes)
	MarginCalculations(MarginCalculationRequest, ReqHeader) (int, apihelpers.APIRes)
	LastTradedPrice(LastTradedPriceRequest, ReqHeader) (int, apihelpers.APIRes)
	PlaceIcebergOrder(IcebergOrderReq, ReqHeader) (int, apihelpers.APIRes)
	ModifyIcebergOrder(ModifyIcebergOrderReq, ReqHeader) (int, apihelpers.APIRes)
	CancelIcebergOrder(CancelIcebergOrderReq, ReqHeader) (int, apihelpers.APIRes)
}

type PortfolioProvider interface {
	FetchDematHoldings(FetchDematHoldingsRequest, ReqHeader) (int, apihelpers.APIRes)
	ConvertPositions(ConvertPositionsRequest, ReqHeader) (int, apihelpers.APIRes)
	GetPositions(GetPositionRequest, ReqHeader) (int, apihelpers.APIRes)
}

type OptionChainProvider interface {
	FetchOptionChain(FetchOptionChainRequest, ReqHeader) (int, apihelpers.APIRes)
	FetchFuturesChain(FetchFuturesChainReq, ReqHeader) (int, apihelpers.APIRes)
}

type ProfileProvider interface {
	GetProfile(ProfileRequest, ReqHeader) (int, apihelpers.APIRes)
	SendAFOtp(SendAFOtpReq, ReqHeader) (int, apihelpers.APIRes)
	VerifyAFOtp(VerifyAFOtpReq, ReqHeader) (int, apihelpers.APIRes)
	AccountFreeze(AccountFreezeReq, ReqHeader) (int, apihelpers.APIRes)
}

type FetchFundsProvider interface {
	FetchFunds(FetchFundsRequest, ReqHeader) (int, apihelpers.APIRes)
	CancelPayout(CancelPayoutReq, ReqHeader) (int, apihelpers.APIRes)
	Payout(AtomPayoutRequest, ReqHeader) (int, apihelpers.APIRes)
	ClientTransactions(ClientTransactionsRequest, ReqHeader) (int, apihelpers.APIRes)
}

type FetchFundsProviderV3 interface {
	CancelPayout(CancelPayoutReqV3, ReqHeader) (int, apihelpers.APIRes)
	Payout(AtomPayoutRequest, ReqHeader) (int, apihelpers.APIRes)
}

type ContractDetailsProvider interface {
	SearchScrip(SearchScripRequest, ReqHeader) (int, apihelpers.APIRes)
	ScripInfo(ScripInfoRequest, ReqHeader) (int, apihelpers.APIRes)
}

type ScripProvider interface {
	SearchScrip(SearchScripAPIRequest, ReqHeader) (int, apihelpers.APIRes)
}

type ContractDetailsProviderV2 interface {
	SearchScrip(string, ReqHeader, int, int, string) (int, apihelpers.APIRes)
	PerformSearch(string, string, int, int) ([]interface{}, error)
}

type ConditionalOrderProvider interface {
	PlaceBOOrder(PlaceBOOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	ModifyBOOrder(ModifyBOOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	CancelBOOrder(ExitBOOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	PlaceCOOrder(PlaceCOOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	ModifyCOOrder(ModifyCOOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	CancelCOOrder(ExitCOOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	PlaceSpreadOrder(PlaceSpreadOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	ModifySpreadOrder(ModifySpreadOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	CancelSpreadOrder(ExitSpreadOrderRequest, ReqHeader) (int, apihelpers.APIRes)
}

type IpoProvider interface {
	GetAllIpo(GetAllIpoRequest, ReqHeader) (int, apihelpers.APIRes)
	PlaceIpoOrder(PlaceIpoOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	FetchIpoOrder(FetchIpoOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	CancelIpoOrder(CancelIpoOrderRequest, ReqHeader) (int, apihelpers.APIRes)
	FetchIpoData(FetchIpoDataRequest, ReqHeader) (int, apihelpers.APIRes)
	FetchIpoGmpData(FetchIpoDataRequest, ReqHeader) (int, apihelpers.APIRes)
	FetchEIpo(FetchEipoReq, ReqHeader) (int, apihelpers.APIRes)
}

type GainerLoserProvider interface {
	GainerLoserNse(ReqHeader) (int, apihelpers.APIRes)
	GainerLoserNiftyFifty(Req GainersLosersMostActiveVolumeReq, reqH ReqHeader) (int, apihelpers.APIRes)
	MostActiveVolumeNSE(ReqHeader) (int, apihelpers.APIRes)
	MostActiveVolumeDataNifty50(Req GainersLosersMostActiveVolumeReq, reqH ReqHeader) (int, apihelpers.APIRes)
	ChartData(ChartDataReq, ReqHeader) (int, apihelpers.APIRes)
	ReturnOnInvestment(ReturnOnInvestmentReq, ReqHeader) (int, apihelpers.APIRes)
	FetchHistoricPerformance(HistoricPerformaceReq, ReqHeader) (int, apihelpers.APIRes)
	FetchAllHistoricPerformance(AllHistoricPerformaceReq, ReqHeader) (int, apihelpers.APIRes)
}

type ChargesProvider interface {
	BrokerCharges(BrokerChargesReq, ReqHeader) (int, apihelpers.APIRes)
	CombineBrokerCharges(CombineBrokerChargesReq, ReqHeader) (int, apihelpers.APIRes)
	FundsPayout(FundsPayoutReq, ReqHeader) (int, apihelpers.APIRes)
}

type TestProvider interface {
	TestingApi(TestingReq, ReqHeader) (int, apihelpers.APIRes)
}

type BasketOrderProvider interface {
	CreateBasket(CreateBasketReq, ReqHeader) (int, apihelpers.APIRes)
	FetchBasket(FetchBasketReq, ReqHeader) (int, apihelpers.APIRes)
	DeleteBasket(DeleteBasketReq, ReqHeader) (int, apihelpers.APIRes)
	AddBasketInstrument(AddBasketInstrumentReq, ReqHeader) (int, apihelpers.APIRes)
	EditBasketInstrument(EditBasketInstrumentReq, ReqHeader) (int, apihelpers.APIRes)
	DeleteBasketInstrument(DeleteBasketInstrumentReq, ReqHeader) (int, apihelpers.APIRes)
	RenameBasket(RenameBasketReq, ReqHeader) (int, apihelpers.APIRes)
	ExecuteBasket(ExecuteBasketReq, ReqHeader) (int, apihelpers.APIRes)
	UpdateBasketExecutionState(UpdateBasketExecutionStateReq, ReqHeader) (int, apihelpers.APIRes)
}

type BackofficeProvider interface {
	TradeConfirmationDateRange(TradeConfirmationDateRangeReq, ReqHeader) (int, apihelpers.APIRes)
	GetBillDetailsCdsl(GetBillDetailsCdslReq, ReqHeader) (int, apihelpers.APIRes)
	LongTermShortTerm(LongTermShortTermReq, ReqHeader) (int, apihelpers.APIRes)
	FetchProfile(FetchProfileReq, ReqHeader) (int, apihelpers.APIRes)
	TradeConfirmationOnDate(TradeConfirmationOnDateReq, ReqHeader) (int, apihelpers.APIRes)
	OpenPositions(OpenPositionsReq, ReqHeader) (int, apihelpers.APIRes)
	GetHolding(GetHoldingReq, ReqHeader) (int, apihelpers.APIRes)
	GetMarginOnDate(GetMarginOnDateReq, ReqHeader) (int, apihelpers.APIRes)
	FinancialLedgerBalanceOnDate(FinancialLedgerBalanceOnDateReq, ReqHeader) (int, apihelpers.APIRes)
	GetFinancial(GetFinancialReq, ReqHeader) (int, apihelpers.APIRes)
	GetScripWiseCostingData(TradebookReq, ReqHeader) (ScripWiseCostingRes, error)
	GetFinancialLedgerData(GetFinancialLedgerDataReq, ReqHeader) (FinancialLedgerRes, error)
	GetOpenPositionData(OpenPositionReq, ReqHeader) (OpenPositionRes, error)
	GetFONetPositionData(GetFONetPositionDataReq, ReqHeader) (FONetPositionRes, error)
	GetHoldingFinancialData(GetHoldingFinancialDataReq, ReqHeader) (GetHoldingFinancialDataRes, error)
	GetCommodityTransactionData(CommodityTradebookReq, ReqHeader) (CommodityTransactionRes, error)
	GetFNOTransactionData(FNOTradebookReq, ReqHeader) (FNOTransactionRes, error)
}

type AlertsProvider interface {
	CreateAlert(CreateAlertsReq, ReqHeader) (int, apihelpers.APIRes)
	EditAlerts(EditAlertsReq, ReqHeader) (int, apihelpers.APIRes)
	GetAlerts(GetAlertsReq, ReqHeader) (int, apihelpers.APIRes)
	PauseAlerts(PauseAlertsReq, ReqHeader) (int, apihelpers.APIRes)
	DeleteAlerts(DeleteAlertsReq, ReqHeader) (int, apihelpers.APIRes)
}

type SquareOffProvider interface {
	SquareOffAll(SquareOffAllReq, ReqHeader) (int, apihelpers.APIRes)
}

type CMOTSProvider interface {
	GetOverview(GetOverviewReq, ReqHeader) (int, apihelpers.APIRes)
	FetchFinancials(FetchFinancialsReq, ReqHeader) (int, apihelpers.APIRes)
	FetchFinancialsDetailed(FetchFinancialsDetailedReq, ReqHeader) (int, apihelpers.APIRes)
	FetchPeers(FetchPeersReq, ReqHeader) (int, apihelpers.APIRes)
	ShareHoldingPatterns(ShareHoldingPatternsReq, ReqHeader) (int, apihelpers.APIRes)
	RatiosCompare(RatiosCompareReq, ReqHeader) (int, apihelpers.APIRes)
	FetchTechnicalIndicators(FetchTechnicalIndicatorsReq, ReqHeader) (int, apihelpers.APIRes)
	StocksOnNews(StocksOnNewsReq, ReqHeader) (int, apihelpers.APIRes)
	FetchSectorList(sectorCode string, requestH ReqHeader) (int, apihelpers.APIRes)
	FetchSectorWiseCompany(FetchSectorWiseCompanyReq, ReqHeader) (int, apihelpers.APIRes)
	FetchCompanyCategory(FetchCompanyCategoryReq, ReqHeader) (int, apihelpers.APIRes)
	StocksAnalyzer(StocksAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	CorporateActionsIndividual(FetchCorporateActionsIndividualReq, ReqHeader) (int, apihelpers.APIRes)
	CorporateActionsAll(FetchCorporateActionsAllReq, ReqHeader) (int, apihelpers.APIRes)
	GetSectorWiseStockList(page int, sectorCode, sectorName string, requestH ReqHeader) (int, apihelpers.APIRes)
}

type CMOTSProviderV2 interface {
	StocksOnNewsV2(StocksOnNewsV2Req, ReqHeader) (int, apihelpers.APIRes)
	FetchFinancialsV2(FetchFinancialsReq, ReqHeader) (int, apihelpers.APIRes)
	FetchPeersV2(FetchPeersV2Req, ReqHeader) (int, apihelpers.APIRes)
	FetchSectorListV2(sectorCode string, requestH ReqHeader) (int, apihelpers.APIRes)
	FetchSectorWiseCompanyV2(FetchSectorWiseCompanyReqV2, ReqHeader) (int, apihelpers.APIRes)
}

type UserDetailsProvider interface {
	GetAllBankAccounts(GetAllBankAccountsReq, ReqHeader) (int, apihelpers.APIRes)
	GetAllBankAccountsUpdated(GetAllBankAccountsUpdatedReq, ReqHeader) (int, apihelpers.APIRes)
	GetUserId(GetUserIdReq, ReqHeader) (int, apihelpers.APIRes)
	UserNotifications(UserNotificationsReq, ReqHeader) (int, apihelpers.APIRes)
	GetClientStatus(string, ReqHeader) (int, apihelpers.APIRes)
}

type PortfolioAnalyzer interface {
	HoldingsWeightages(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	PortfolioBeta(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	PortfolioPE(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	PortfolioDE(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	HighPledgedPromoterHoldings(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	AdditionalSurveillanceMeasure(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	GradedSurveillanceMeasure(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	HighDefaultProbability(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	LowROE(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	LowProfitGrowth(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	HoldingStockContribution(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	InvestmentSector(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	DeclineInPromoterHolding(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	InterestCoverageRatio(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	DeclineInRevenueAndProfit(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	LowNetWorth(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	DeclineInRevenue(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	PromoterPledge(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	PennyStocks(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	StockReturn(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	NiftyVsPortfolio(NiftyVsPortfolioReq, ReqHeader) (int, apihelpers.APIRes)
	ChangeInInstitutionalHolding(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	RoeAndStockReturn(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
	IlliquidStocks(PortfolioAnalyzerReq, ReqHeader) (int, apihelpers.APIRes)
}

type SessionInfoProvider interface {
	SessionInfo(SessionInfoReq, ReqHeader) (int, apihelpers.APIRes)
}

type TechnicalIndicators interface {
	TechnicalIndicatorsValues(TechnicalIndicatorsValuesReq, ReqHeader) (int, apihelpers.APIRes)
}

type NotificationsProvider interface {
	FetchAdminMessages(FetchAdminMessageRequest, ReqHeader) (int, apihelpers.APIRes)
	NotificationUpdates(NotificationUpdatesReq, ReqHeader) (int, apihelpers.APIRes)
}

type BondEtfProvider interface {
	FetchBondData(FetchBondDataReq, ReqHeader) (int, apihelpers.APIRes)
}

type WarningProvider interface {
	NudgeAlert(NudgeAlertReq, ReqHeader) (int, apihelpers.APIRes)
}

type CMOTSProviderV3 interface {
	FetchFinancialsV3(FetchFinancialsReq, ReqHeader) (int, apihelpers.APIRes)
}

type CMOTSProviderV4 interface {
	FetchFinancialsV4(FetchFinancialsReq, ReqHeader) (int, apihelpers.APIRes)
}

type OptionChainProviderV2 interface {
	FetchOptionChainV2(FetchOptionChainV2Request, ReqHeader) (int, apihelpers.APIRes)
	FetchOptionChainByExpiryV2(FetchOptionChainByExpiryV2Request, ReqHeader) (int, apihelpers.APIRes)
}

type LogoutProvider interface {
	LogoutSingleDevice(ReqHeader) (int, apihelpers.APIRes)
}

type FinvuProvider interface {
	FinvuConsentRequestPlus(CreateConsentRequestPlusReq, ReqHeader) (int, apihelpers.APIRes)
	FinvuGetBankStatement(FinvuGetBankStatementReq, ReqHeader) (int, apihelpers.APIRes)
}

type EdisProvider interface {
	SendEdisRequest(EdisReq, ReqHeader) (int, apihelpers.APIRes)
	GenerateTpin(TpinReq, ReqHeader) (int, apihelpers.APIRes)
}

type EpledgeProvider interface {
	SendEpledgeRequest(EpledgeReq, ReqHeader) (int, apihelpers.APIRes)
	UnpledgeRequest(UnpledgeReq, ReqHeader) (int, apihelpers.APIRes)
	MTFEpledgeRequest(MTFEPledgeRequest, ReqHeader) (int, apihelpers.APIRes)
	GetPledgeList(ReqHeader) (int, apihelpers.APIRes)
	GetCTDQuantityList(MTFCTDDataReq, ReqHeader) (int, apihelpers.APIRes)
	GetPledgeTransactions(FetchEpledgeTxnReq, ReqHeader) (int, apihelpers.APIRes)
	MTFCTD(MTFCTDReq, ReqHeader) (int, apihelpers.APIRes)
}

type ReportsProvider interface {
	ViewDPCharges(DPChargesReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
	DownloadDPCharges(DPChargesReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
	SendEmailDPCharges(DPChargesReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
	ViewTradebook(TradebookReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
	DownloadTradebook(TradebookReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
	ViewLedger(LedgerReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
	DownloadLedger(LedgerReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
	ViewOpenPosition(OpenPositionReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
	DownloadOpenPosition(OpenPositionReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
	ViewFnoPnl(FnoPnlReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
	DownloadFnoPnl(FnoPnlReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
	ViewHoldingFinancial(GetHoldingFinancialDataReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
	DownloadHoldingFinancial(GetHoldingFinancialDataReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
	SendEmailHoldingFinancial(GetHoldingFinancialDataReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
	SendEmailLedger(LedgerReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
	ViewCommodityTradebook(CommodityTradebookReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
	DownloadCommodityTradebook(CommodityTradebookReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
	SendEmailCommodityTradebook(CommodityTradebookReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
	ViewFnoTradebook(FNOTradebookReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
	DownloadFnoTradebook(FNOTradebookReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
	SendEmailFnoTradebook(FNOTradebookReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
}
type ExecutePocketV2 interface {
	BuyPocketV2(ExecutePocketV2Request, ReqHeader) (int, apihelpers.APIRes)
	ExitPocketV2(ExecutePocketV2Request, ReqHeader) (int, apihelpers.APIRes)
	FetchPocketPortfolioV2(FetchPocketPortfolioRequest, ReqHeader) (int, apihelpers.APIRes)
}

type UpiPreferenceProvider interface {
	SetUpiPreference(SetUpiPreferenceReq, ReqHeader) (int, apihelpers.APIRes)
	FetchUpiPreference(FetchUpiPreferenceReq, ReqHeader) (int, apihelpers.APIRes)
	DeleteUpiPreference(DeleteUpiPreferenceReq, ReqHeader) (int, apihelpers.APIRes)
}

type CollectionsProvider interface {
	CreateCollections(CreateCollectionsRequest, ReqHeader) (int, apihelpers.APIRes)
	ModifyCollections(ModifyCollectionsRequest, ReqHeader) (int, apihelpers.APIRes)
	FetchCollections(FetchCollectionsDetailsRequest, ReqHeader) (int, apihelpers.APIRes)
	DeleteCollections(DeleteCollectionsRequest, ReqHeader) (int, apihelpers.APIRes)
	FetchAllCollections(ReqHeader) (int, apihelpers.APIRes)
}

type PinsProvider interface {
	AddPins(AddPinReq, ReqHeader) (int, apihelpers.APIRes)
	DeletePins(DeletePins, ReqHeader) (int, apihelpers.APIRes)
	FetchPins(PinsRequest, ReqHeader) (int, apihelpers.APIRes)
	UpdatePins(UpdatePins, ReqHeader) (int, apihelpers.APIRes)
}

type PocketsProvider interface {
	AdminLogin(AdminLoginRequest, ReqHeader) (int, apihelpers.APIRes)
	CreatePockets(CreatePocketsRequest, ReqHeader) (int, apihelpers.APIRes)
	ModifyPockets(ModifyPocketsRequest, ReqHeader) (int, apihelpers.APIRes)
	FetchPockets(FetchPocketsDetailsRequest, ReqHeader) (int, apihelpers.APIRes)
	DeletePockets(DeletePocketsRequest, ReqHeader) (int, apihelpers.APIRes)
	FetchAllPockets(ReqHeader) (int, apihelpers.APIRes)
	PocketsCalculations(PocketsCalculationsReq, ReqHeader) (int, apihelpers.APIRes)
	MultipleAndIndividualStocksCalculations(MultipleAndIndividualStocksCalculationsReq, ReqHeader) (int, apihelpers.APIRes)
	ExecutePocket(ExecutePocketRequest, string, ReqHeader) (int, apihelpers.APIRes)
	FetchPocketPortfolio(FetchPocketPortfolioRequest, ReqHeader) (int, apihelpers.APIRes)
	FetchPocketTransaction(FetchPocketTransactionReq, ReqHeader) (int, apihelpers.APIRes)
	StorePocketTransaction(StorePocketTransactionReq, ReqHeader) (int, apihelpers.APIRes)
}

type WatchListProvider interface {
	CreateWatchList(req CreateWatchListRequest, reqH ReqHeader) (int, apihelpers.APIRes)
	ModifyWatchList(req ModifyWatchListRequest, reqH ReqHeader) (int, apihelpers.APIRes)
	FetchWatchLists(req FetchWatchListsRequest, reqH ReqHeader) (int, apihelpers.APIRes)
	DeleteWatchList(req DeleteWatchListRequest, reqH ReqHeader) (int, apihelpers.APIRes)
	FetchWatchListDetails(req FetchWatchListsDetailsRequest, reqH ReqHeader) (int, apihelpers.APIRes)
	AddStockToWatchList(req AddStockToWatchListsRequest, reqH ReqHeader) (int, apihelpers.APIRes)
	FetchWatchListsV2(req FetchWatchListV2Request, reqH ReqHeader) (int, apihelpers.APIRes)
	AddStockToWatchListV2(req AddStockToWatchListV2Request, reqH ReqHeader) (int, apihelpers.APIRes)
	DeleteStockInWatchListV2(req DeleteWatchListV2Request, reqH ReqHeader) (int, apihelpers.APIRes)
	DeleteStockInWatchListV2Updated(req DeleteWatchListV2UpdatedRequest, reqH ReqHeader) (int, apihelpers.APIRes)
	ArrangeStocksWatchListV2(req ArrangeStocksWatchListV2Request, reqH ReqHeader) (int, apihelpers.APIRes)
	FetchWatchListsV3(req FetchWatchListV3Request, reqH ReqHeader) (int, apihelpers.APIRes)
	AddStockToWatchListV3(req AddStockToWatchListV3Request, reqH ReqHeader) (int, apihelpers.APIRes)
	DeleteStockInWatchListV3(req DeleteWatchListV3Request, reqH ReqHeader) (int, apihelpers.APIRes)
	PopulateIsinMappingInLocalCache()
	ArrangeStocksWatchListV3(req ArrangeStocksWatchListV3Request, reqH ReqHeader) (int, apihelpers.APIRes)
	DeleteStockInWatchListV3Updated(req DeleteWatchListV3UpdatedRequest, reqH ReqHeader) (int, apihelpers.APIRes)
}

type FreshdeskProvider interface {
	CreateFreshdeskTicket(req FreshdeskTicketReq, reqH ReqHeader) (int, apihelpers.APIRes)
}

type ExecutePocketV3 interface {
	BuyPocketV3(ExecutePocketV3Request, string, ReqHeader) (int, apihelpers.APIRes)
	CheckActionRequired(clientId, pocketId string, userVersion, usersLotSize int, reqH ReqHeader) (int, apihelpers.APIRes)
	ManageRequiredStocksForPocket(pocketBalanceReq RebalanceResponse, reqH ReqHeader) (int, apihelpers.APIRes)
	FetchAllPocketsV3(ReqHeader) (int, apihelpers.APIRes)
	FetchPocketPortfolioV3(FetchPocketPortfolioRequest, ReqHeader) (int, apihelpers.APIRes)
	FetchUsersPockets(clientId string, reqH ReqHeader) (int, apihelpers.APIRes)
	SellPocketV3(req ExecutePocketV3Request, pocketAction string, reqH ReqHeader) (int, apihelpers.APIRes)
	ExitPocketV3(req ExecutePocketV3Request, reqH ReqHeader) (int, apihelpers.APIRes)
	GetPocketDetails(pocketId, tag string, reqH ReqHeader) (int, apihelpers.APIRes)
}

type TechnicalIndicatorsV2Provider interface {
	GetSMA(req GetSMAReq, reqH ReqHeader) (int, apihelpers.APIRes)
	GetEMA(req GetEMAReq, reqH ReqHeader) (int, apihelpers.APIRes)
	GetHullMA(req GetHullMAReq, reqH ReqHeader) (int, apihelpers.APIRes)
	GetVWMA(req GetVWMAReq, reqH ReqHeader) (int, apihelpers.APIRes)
	GetRSI(req GetRSIReq, reqH ReqHeader) (int, apihelpers.APIRes)
	GetCCI(req GetCCIReq, reqH ReqHeader) (int, apihelpers.APIRes)
	GetMACD(req GetMACDReq, reqH ReqHeader) (int, apihelpers.APIRes)
	GetStochastic(req GetStochasticReq, reqH ReqHeader) (int, apihelpers.APIRes)
	GetIchimokuBaseLine(req GetIchimokuBaseLineReq, reqH ReqHeader) (int, apihelpers.APIRes)
	GetADX(req GetADXReq, reqH ReqHeader) (int, apihelpers.APIRes)
	GetAwesomeOscillator(req GetAwesomeOscillatorReq, reqH ReqHeader) (int, apihelpers.APIRes)
	GetMomentum(req GetMomentumReq, reqH ReqHeader) (int, apihelpers.APIRes)
	GetStochRSIFast(req GetStochRSIFastReq, reqH ReqHeader) (int, apihelpers.APIRes)
	GetWilliamsRange(req GetWilliamsRangeReq, reqH ReqHeader) (int, apihelpers.APIRes)
	GetUltimateOscillator(req GetUltimateOscillatorReq, reqH ReqHeader) (int, apihelpers.APIRes)
	GetAllTechnicalIndicators(req GetAllTechnicalIndicatorsReq, reqH ReqHeader) (int, apihelpers.APIRes)
}

type BondsDetailsProvider interface {
	FetchBondDataByIsin(FetchBondDataByIsinReq, ReqHeader) (int, apihelpers.APIRes)
}

type SipProvider interface {
	GetStockSips(string, ReqHeader) (int, apihelpers.APIRes)
	PlaceSipOrder(PlaceSipRequest, ReqHeader) (int, apihelpers.APIRes)
	DeleteSipOrder(string, string, ReqHeader) (int, apihelpers.APIRes)
	ModifySipOrder(ModifySipRequest, ReqHeader) (int, apihelpers.APIRes)
	UpdateSipStatus(UpdateSipStatusRequest, ReqHeader) (int, apihelpers.APIRes)
}

// blockdeals interface...

type BlockDealService interface {
	CreateBlockDeal(c *gin.Context, blockDeal BlockDeal) error
	GetAllBlockDeals(c *gin.Context) ([]BlockDeal, error)
	GetBlockDealByID(c *gin.Context, id int) (BlockDeal, error)
	UpdateBlockDeal(cocode int, blockDeal BlockDeal) error
	DeleteBlockDeal(c *gin.Context, id int) error
}
