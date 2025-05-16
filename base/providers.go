package base

import (
	businessV2 "space/business/V2"
	"space/business/backoffice"
	bondetf "space/business/bondEtf"
	bondsdetails "space/business/bondsDetails"
	"space/business/charges"
	"space/business/cmots"
	"space/business/collections"
	"space/business/finvu"
	"space/business/freshdesk"
	"space/business/funds"
	"space/business/pins"
	"space/business/pockets"
	portfolioanalyzer "space/business/portfolioAnalyzer"
	"space/business/reports"
	"space/business/scrips"
	searchscriptv2 "space/business/searchScriptV2"
	technicalindicators "space/business/technicalIndicators"
	technicalindicatorsV2 "space/business/technicalIndicatorsV2"
	"space/business/tradelab"
	upipreference "space/business/upi"
	userdetails "space/business/userDetails"
	warning "space/business/warning"
	"space/business/watchlists"
	v1 "space/controllers/api/v1"
	v2 "space/controllers/api/v2"
	v3 "space/controllers/api/v3"
	v4 "space/controllers/api/v4"
	"space/db"
	"space/helpers/cache"
	"space/models"
)

func InitProviders(mongodb db.MongoDatabase, redisCli cache.RedisCache, contractCacheCli cache.ContractCache, smartCacheCli cache.SmartCache) {
	defer models.HandlePanic()
	//build login provider
	loginProvider := BuildLoginProvider(mongodb, redisCli)
	v1.InitLoginProvider(loginProvider)

	//build order provider
	orderProvider := BuildOrderProvider(redisCli)
	v1.InitOrderProvider(orderProvider)

	orderProviderV2 := BuildOrderProviderV2(redisCli)
	v2.InitOrderProviderV2(orderProviderV2)

	//build portfolio provider
	portfolioProvider := BuildPortfolioProvider()
	v1.InitPortfolioProvider(portfolioProvider)

	portfolioProviderV2 := BuildPortfolioProviderV2()
	v2.InitPortfolioProviderV2(portfolioProviderV2)

	//build option chain provider
	optionChainProvider := BuildOptionChainProvider()
	v1.InitOptionChainProvider(optionChainProvider)

	optionChainProviderV2 := BuildOptionChainProviderV2(redisCli)
	v2.InitOptionChainProviderV2(optionChainProviderV2)

	optionChainProviderV3 := BuildOptionChainProviderV3()
	v3.InitOptionChainProviderV3(optionChainProviderV3)

	// build funds v1
	fundsProvider := BuildFundsProvider()
	v1.InitFundsProvider(fundsProvider)

	fundsProviderV2 := BuildFundsProviderV2()
	v2.InitFundsProviderV2(fundsProviderV2)

	fundsProviderV3 := BuildFundsProviderV3(redisCli)
	v3.InitFundsProviderV3(fundsProviderV3)

	// build profile
	profile := BuildProfileProvider(mongodb, redisCli)
	v1.InitProfileProvider(profile)

	profileV2 := BuildProfileProviderV2(mongodb, redisCli)
	v2.InitProfileProviderV2(profileV2)

	//build contract details provider
	contractDetailsProvider := BuildContractDetailsProvider(contractCacheCli)
	v1.InitContractDetailsProvider(contractDetailsProvider)

	contractDetailsProviderV2 := BuildContractDetailsProviderV2TL(contractCacheCli)
	v2.InitContractDetailsProviderV2TL(contractDetailsProviderV2)
	//build conditional order provider
	conditionalOrderProvider := BuildConditionalOrderProvider(redisCli)
	v1.InitConditionalOrderProvider(conditionalOrderProvider)

	conditionalOrderProviderV2 := BuildConditionalOrderProviderV2(redisCli)
	v2.InitConditionalOrderProviderV2(conditionalOrderProviderV2)

	// build ipo provider
	ipoProvider := BuildIpoProvider(mongodb, redisCli)
	v1.InitIpoProvider(ipoProvider)

	// build gainerLoser provider
	gainerLoserProvider := BuildGainerLoserProvider()
	v1.InitGainerLoserProvider(gainerLoserProvider)

	gainerLoserProviderV2 := BuildGainerLoserProviderV2()
	v2.InitGainerLoserProviderV2(gainerLoserProviderV2)

	// build charges provider
	chargesProvider := BuildChargesProvider()
	v1.InitChargesProvider(chargesProvider)

	// build basketorder provider
	basketOrderProvider := BuildBasketOrderProvider()
	v1.InitBasketOrderProvider(basketOrderProvider)

	alertsProvider := BuildAlertProvider()
	v1.InitAlertsProvider(alertsProvider)

	alertsProviderV2 := BuildAlertProviderV2()
	v2.InitAlertsProviderV2(alertsProviderV2)

	squareOffProvider := BuildSquareOffProviderr()
	v1.InitSquareOffProvider(squareOffProvider)

	cmotsProvider := BuildCmotsProvider(contractCacheCli)
	v1.InitCmotsProvider(cmotsProvider)

	userDetailsProvider := BuildUserDetailsProvider(mongodb)
	v1.InitUserDetailsProvider(userDetailsProvider)

	portfolioAnalyzerProvider := BuildPortfolioAnalyzerProvider()
	v1.InitPortfolioAnalyzerProvider(portfolioAnalyzerProvider)

	sessionInfoProvider := BuildSessionInfoProvider()
	v1.InitSessionInfoProvider(sessionInfoProvider)

	loginProviderV2 := BuildLoginProviderV2(mongodb, redisCli)
	v2.InitLoginProviderV2(loginProviderV2)

	loginProviderV3 := BuildLoginProviderV3(mongodb, redisCli)
	v3.InitLoginProviderV3(loginProviderV3)

	logoutProvider := BuildLogoutProvider()
	v2.InitLogoutProvider(logoutProvider)

	cmotsProviderV2 := BuildCmotsProviderV2(contractCacheCli)
	v2.InitCmotsProviderV2(cmotsProviderV2)

	technicalIndicatorsProvider := BuildTechnicalIndicatorsProvider()
	v1.InitTechnicalIndicatorsProvider(technicalIndicatorsProvider)

	technicalIndicatorsProviderV2 := BuildTechnicalIndicatorsProviderV2()
	v2.InitTechnicalIndicatorsV2(technicalIndicatorsProviderV2)

	notificationsProvider := BuildNotificationsProvider()
	v1.InitNotificationsProvider(notificationsProvider)

	notificationsProviderV2 := BuildNotificationsProviderV2()
	v2.InitNotificationsProviderV2(notificationsProviderV2)

	bondEtfProvider := BuildBondEtfProvider()
	v1.InitBondEtfProvider(bondEtfProvider)

	scripsProvider := BuildScripsProvider(mongodb)
	v3.InitScripsProvider(scripsProvider)

	searchScriptProviderV2 := BuildSearchScriptProviderV2(contractCacheCli, smartCacheCli)
	v2.InitContractDetailsProvider(searchScriptProviderV2)

	warningProvider := BuildWarningProvider()
	v1.InitWarningProvider(warningProvider)

	cmotsProviderV3 := BuildCmotsProviderV3(contractCacheCli)
	v3.InitCmotsProviderV3(cmotsProviderV3)

	backofficeProvider := BuildBackOfficeProvider()
	v1.InitBackOfficeProvider(backofficeProvider)

	finvuProvider := BuildFinvuProvider(mongodb, redisCli)
	v1.InitFinvuProvider(finvuProvider)

	// // build edis provider
	edisProvider := BuildEdisProvider()
	v1.InitEdisProvider(edisProvider)

	// build epledge provider
	epledgeProvider := BuildEpledgeProvider()
	v1.InitEpledgeProvider(epledgeProvider)

	cmotsProviderV4 := BuildCmotsProviderV4(contractCacheCli)
	v4.InitCmotsProviderV4(cmotsProviderV4)

	reportsProvider := BuildReportsProvider(redisCli)
	v1.InitReportsProvider(reportsProvider)
	// BuildExecutePocketV2Provider
	executePocketProviderV2 := BuildExecutePocketV2Provider(mongodb, redisCli)
	v2.InitExecutePocketV2Provider(executePocketProviderV2)

	executePocketProviderV3 := BuildExecutePocketV3Provider(mongodb, redisCli)
	v3.InitExecutePocketV3Provider(executePocketProviderV3)

	upiPreferenceProvider := BuildUpiPreferenceProvider(mongodb)
	v1.InitUpiPreferenceProvider(upiPreferenceProvider)

	collectionsProvider := BuildCollectionsProvider(mongodb)
	v1.InitCollectionProvider(collectionsProvider)

	pinsProvider := BuildPinsProvider(mongodb)
	v1.InitPinsProvider(pinsProvider)

	pinsProviderV2 := BuildPinsProvider(mongodb)
	v2.InitPinsProviderV2(pinsProviderV2)

	pocketsProvider := BuildPocketsProvider(mongodb, redisCli)
	v1.InitPocketsProvider(pocketsProvider)

	watchlistProviderV2 := BuildWatchListProviderV2(mongodb, contractCacheCli)
	v2.InitWatchListProviderV2(watchlistProviderV2)

	watchlistProvider := BuildWatchListProviderV1(mongodb, contractCacheCli)
	v1.InitWatchListProviderV1(watchlistProvider)

	watchlistProviderV3 := BuildWatchListProviderV3(mongodb, contractCacheCli)
	v3.InitWatchListProviderV3(watchlistProviderV3)

	freshdeskProvider := BuildFreshDeskProvider()
	v1.InitFreshdeskProvider(freshdeskProvider)

	bondsDetailsProvider := BuildBondsDetailsProvider(mongodb, contractCacheCli)
	v1.InitBondsDetailsProvider(bondsDetailsProvider)

	sipProvider := BuildSipProvider()
	v1.InitSipProvider(sipProvider)

}

func BuildLoginProvider(mongodb db.MongoDatabase, redisCli cache.RedisCache) models.LoginProvider {
	//based on vendor provider can be initialized here
	return tradelab.InitLogin(mongodb, redisCli)
}

func BuildLoginProviderV3(mongodb db.MongoDatabase, redisCli cache.RedisCache) models.LoginProviderV3 {
	return tradelab.InitLogin(mongodb, redisCli)
}

func BuildLogoutProvider() models.LogoutProvider {
	return tradelab.InitLogoutProvider()
}

func BuildOrderProvider(redisCli cache.RedisCache) models.OrderProvider {
	return tradelab.InitOrder(redisCli)
}

func BuildOrderProviderV2(redisCli cache.RedisCache) models.OrderProvider {
	return tradelab.InitOrder(redisCli)
}

func BuildPortfolioProvider() models.PortfolioProvider {
	return tradelab.InitPortfolio()
}

func BuildPortfolioProviderV2() models.PortfolioProvider {
	return tradelab.InitPortfolio()
}

func BuildOptionChainProvider() models.OptionChainProvider {
	return tradelab.InitOptionChain()
}

func BuildOptionChainProviderV3() models.OptionChainProvider {
	return tradelab.InitOptionChain()
}

func BuildProfileProvider(mongodb db.MongoDatabase, redisCli cache.RedisCache) models.ProfileProvider {
	return tradelab.InitProfile(mongodb, redisCli)
}

func BuildProfileProviderV2(mongodb db.MongoDatabase, redisCli cache.RedisCache) models.ProfileProvider {
	return tradelab.InitProfile(mongodb, redisCli)
}

func BuildFundsProvider() models.FetchFundsProvider {
	pgDB := db.GetPgObj()
	return tradelab.InitFetchFunds(pgDB)
}

func BuildFundsProviderV2() models.FetchFundsProvider {
	pgDB := db.GetPgObj()
	return tradelab.InitFetchFunds(pgDB)
}

func BuildFundsProviderV3(redisCli cache.RedisCache) models.FetchFundsProviderV3 {
	pgDB := db.GetPgObj()
	return funds.InitFundsV3(pgDB, redisCli)
}

func BuildScripsProvider(mongodb db.MongoDatabase) models.ScripProvider {
	return scrips.InitScrips(mongodb)
}

func BuildContractDetailsProvider(contractCacheCli cache.ContractCache) models.ContractDetailsProvider {
	return tradelab.InitContractDetails(contractCacheCli)
}

func BuildContractDetailsProviderV2TL(contractCacheCli cache.ContractCache) models.ContractDetailsProvider {
	return tradelab.InitContractDetails(contractCacheCli)
}

func BuildConditionalOrderProvider(redisCli cache.RedisCache) models.ConditionalOrderProvider {
	return tradelab.InitConditionalOrder(redisCli)
}

func BuildConditionalOrderProviderV2(redisCli cache.RedisCache) models.ConditionalOrderProvider {
	return tradelab.InitConditionalOrder(redisCli)
}

func BuildIpoProvider(mongodb db.MongoDatabase, redisCli cache.RedisCache) models.IpoProvider {
	return tradelab.InitIpoProvider(mongodb, redisCli)
}

func BuildGainerLoserProvider() models.GainerLoserProvider {
	return tradelab.InitGainerLoserProvider()
}

func BuildGainerLoserProviderV2() models.GainerLoserProvider {
	return tradelab.InitGainerLoserProvider()
}

func BuildChargesProvider() models.ChargesProvider {
	// return charges.InitChargesProvider()
	return charges.InitChargesProvider()
}

func BuildBasketOrderProvider() models.BasketOrderProvider {
	return tradelab.InitBasketOrder()
}

func BuildAlertProvider() models.AlertsProvider {
	return tradelab.InitAlertsProvider()
}

func BuildAlertProviderV2() models.AlertsProvider {
	return tradelab.InitAlertsProvider()
}

func BuildSquareOffProviderr() models.SquareOffProvider {
	return tradelab.InitSquareOffProvider()
}

func BuildCmotsProvider(contractCacheCli cache.ContractCache) models.CMOTSProvider {
	pgDB := db.GetPgObj()
	return cmots.InitCmotsProvider(pgDB, contractCacheCli)
}

func BuildUserDetailsProvider(mongodb db.MongoDatabase) models.UserDetailsProvider {
	return userdetails.InitUserDetailsProvider(mongodb)
}

func BuildPortfolioAnalyzerProvider() models.PortfolioAnalyzer {
	return portfolioanalyzer.InitPortfolioAnalyzer()
}

func BuildSessionInfoProvider() models.SessionInfoProvider {
	return tradelab.InitSessionInfoProvider()
}

func BuildLoginProviderV2(mongodb db.MongoDatabase, redisCli cache.RedisCache) models.LoginProviderV2 {
	//based on vendor provider can be initialized here
	return tradelab.InitLogin(mongodb, redisCli)
}

func BuildCmotsProviderV2(contractCacheCli cache.ContractCache) models.CMOTSProviderV2 {
	pgDB := db.GetPgObj()
	return cmots.InitCmotsProvider(pgDB, contractCacheCli)
}

func BuildTechnicalIndicatorsProvider() models.TechnicalIndicators {
	return technicalindicators.InitTechnicalIndicators()
}

func BuildTechnicalIndicatorsProviderV2() models.TechnicalIndicatorsV2Provider {
	return technicalindicatorsV2.InitTechnicalIndicatorsV2()
}

func BuildNotificationsProvider() models.NotificationsProvider {
	return tradelab.InitNotificationsObj()
}

func BuildNotificationsProviderV2() models.NotificationsProvider {
	return tradelab.InitNotificationsObj()
}

func BuildBondEtfProvider() models.BondEtfProvider {
	return bondetf.InitBondEtfObj()
}

func BuildSearchScriptProviderV2(contractCacheCli cache.ContractCache, smartCacheCli cache.SmartCache) models.ContractDetailsProviderV2 {
	return searchscriptv2.InitSearchScript(contractCacheCli, smartCacheCli)
}

func BuildWarningProvider() models.WarningProvider {
	return warning.InitWarningObj()
}

func BuildCmotsProviderV3(contractCacheCli cache.ContractCache) models.CMOTSProviderV3 {
	pgDB := db.GetPgObj()
	return cmots.InitCmotsProvider(pgDB, contractCacheCli)
}

func BuildOptionChainProviderV2(redisCli cache.RedisCache) models.OptionChainProviderV2 {
	return businessV2.InitOptionChainV2(redisCli)
}

func BuildBackOfficeProvider() models.BackofficeProvider {
	return backoffice.InitBackofficeObj()
}

func BuildFinvuProvider(mongodb db.MongoDatabase, redisCli cache.RedisCache) models.FinvuProvider {
	return finvu.InitFinvuProvider(mongodb, redisCli)
}

func BuildEdisProvider() models.EdisProvider {
	return tradelab.InitEdisProvider()
}

func BuildEpledgeProvider() models.EpledgeProvider {
	return tradelab.InitEpledgeProvider()
}

func BuildCmotsProviderV4(contractCacheCli cache.ContractCache) models.CMOTSProviderV4 {
	pgDB := db.GetPgObj()
	return cmots.InitCmotsProvider(pgDB, contractCacheCli)
}

func BuildExecutePocketV2Provider(mongodb db.MongoDatabase, redisCli cache.RedisCache) models.ExecutePocketV2 {
	return pockets.InitExecutePocketV2Provider(mongodb, redisCli)
}

func BuildExecutePocketV3Provider(mongodb db.MongoDatabase, redisCli cache.RedisCache) models.ExecutePocketV3 {
	return pockets.InitExecutePocketV3Provider(mongodb, redisCli)
}

func BuildUpiPreferenceProvider(mongodb db.MongoDatabase) models.UpiPreferenceProvider {
	return upipreference.InitUpiPreferenceProvider(mongodb)
}

func BuildReportsProvider(redisCli cache.RedisCache) models.ReportsProvider {
	return reports.InitReportsProvider(redisCli)
}

func BuildCollectionsProvider(mongodb db.MongoDatabase) models.CollectionsProvider {
	return collections.InitCollections(mongodb)
}

func BuildPinsProvider(mongodb db.MongoDatabase) models.PinsProvider {
	return pins.InitPins(mongodb)
}

func BuildPocketsProvider(mongodb db.MongoDatabase, redisCli cache.RedisCache) models.PocketsProvider {
	return pockets.InitExecutePocketV2Provider(mongodb, redisCli)
}

func BuildWatchListProviderV2(mongodb db.MongoDatabase, contractCacheCli cache.ContractCache) models.WatchListProvider {
	return watchlists.InitWatchlists(mongodb, contractCacheCli)
}

func BuildWatchListProviderV1(mongodb db.MongoDatabase, contractCacheCli cache.ContractCache) models.WatchListProvider {
	return watchlists.InitWatchlists(mongodb, contractCacheCli)
}

func BuildWatchListProviderV3(mongodb db.MongoDatabase, contractCacheCli cache.ContractCache) models.WatchListProvider {
	return watchlists.InitWatchlists(mongodb, contractCacheCli)
}

func BuildFreshDeskProvider() models.FreshdeskProvider {
	return freshdesk.InitFreshdeskProvider()
}

func BuildBondsDetailsProvider(mongodb db.MongoDatabase, contractCacheCli cache.ContractCache) models.BondsDetailsProvider {
	return bondsdetails.InitBondsDetailsProvider(mongodb, contractCacheCli)
}

func BuildSipProvider() models.SipProvider {
	return tradelab.InitSipSerivceProvider()
}
