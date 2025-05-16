package db

import (
	"space/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database interface {
	FetchOverviewData(req models.GetOverviewReq) ([]models.GetOverviewRes, error)
	FetchFinancialsData(req models.FetchFinancialsReq) (models.FetchFinancialsRes, error)
	FetchQuarterlyData(req models.FetchFinancialsDetailedReq) ([]models.QuarterlyData, error)
	FetchPeersData(req models.FetchPeersReq) ([]models.FetchPeerData, error)
	FetchShareHoldingPatternsData(req models.ShareHoldingPatternsReq) ([]models.ShareHoldingPatternsRes, error)
	FetchRatiosCompareData(req models.RatiosCompareReq) ([]models.RatiosCompareRes, error)
	FetchSector(isin string) (string, string, error)
	FetchPE(isin string) (float64, error)
	FetchDE(isin string) (float64, error)
	FetchTechnicalIndicatorsData(req models.FetchTechnicalIndicatorsReq) (models.FetchTechnicalIndicatorsRes, error)
	FetchHighPledgePromoterHoldingMatchData(allIsin models.AllIsin) (models.AllHighPledgePromoterHolding, error)
	FetchAdditionalSurveillanceMeasureData(allIsin models.AllIsin) (models.AllAdditionalSurveillanceMeasure, error)
	FetchGradedSurveillanceMeasureData(allIsin models.AllIsin) (models.AllGradedSurveillanceMeasure, error)
	FetchRoeData(allIsin models.AllIsin) (models.AllLowRoe, error)
	FetchTokenAndSymbol(stringOfKeys string, fetchBy string) ([]models.FetchTokenAndSymbol, error)
	FetchLowProfitGrowthData(allIsin models.AllIsin) (models.AllProfitabilityGrowthDb, error)
	FetchSectorListData(sectorCode string) ([]models.SectorList, error)
	FetchSectorWiseCompanyData(sectorCode string) ([]models.SectorWiseCompany, error)
	FetchCompanyCategory(stringOfisin string) ([]models.CompanyCategory, error)
	FetchDailyAnnouncement(stringOfCoCode string) ([]models.DailyAnnouncement, error)
	FetchBoardMeeting(stringOfCoCode string) ([]models.BoardMeetingForthComing, error)
	FetchChangedName(stringOfCoCode string) ([]models.ChangeOfName, error)
	FetchSplits(stringOfCoCode string) ([]models.Splits, error)
	FetchMerger(stringOfCoCode string) ([]models.MergerDemerger, error)
	FetchDividend(stringOfCoCode string) ([]models.DividendAnnouncementData, error)
	FetchBulkDeals(stringOfCoCode string) ([]models.BulkDeals, error)
	FetchBlockDeals(stringOfCoCode string) ([]models.BlockDeals, error)
	FetchBonus(stringOfCoCode string) ([]models.Bonus, error)
	FetchPLStatementData(isin string) (models.PLStatementResponse, error)
	FetchBalanceSheetsData(isin string) (models.BalanceSheetsResponse, error)
	FetchCashFlowData(isin string) (models.CashflowResponse, error)
	FetchCompanyMasterData(allIsin models.AllIsin) (models.CompanyMasterDb, error)
	GetPostgresStatus() error
	FetchDeclineInPromoterHoldingData(allIsin models.AllIsin) (models.AllDeclineInPromoterHoldingDb, error)
	FetchInterestCoverageRatioData(allIsin models.AllIsin) (models.AllInterestCoverageRatioDb, error)
	DeclineInRevenueAndProfitData(allIsin models.AllIsin) (models.AllDeclineInRevenueAndProfitDb, error)
	LowNetWorthData(allIsin models.AllIsin) (models.AllNetWorthDb, error)
	DeclineInRevenueData(allIsin models.AllIsin) (models.AllDeclineInRevenueDb, error)
	PromoterPledgeData(allIsin models.AllIsin) (models.AllPromoterPledgeDataDb, error)
	PennyStocksData(allIsin models.AllIsin) (models.AllPennyStocksDataDb, error)
	StockReturnData(allIsin models.AllIsin) (models.AllStockReturn, error)
	FetchFinancialsDataV2(req models.FetchFinancialsReq) (models.FetchFinancialsV2Res, error)
	FetchChangeInInstitutionalHoldingData(allIsin models.AllIsin) (models.AllChangeInInstitutionalHoldingDb, error)
	FetchRoeAndStockReturnData(allIsin models.AllIsin) (models.AllRoeAndStockReturnDb, error)
	FetchNseBondData(isin string) (models.NseBondStoreDbData, error)
	NudgeCheck(isin string) (bool, bool, error)
	FetchPeersV2Data(req models.FetchPeersV2Req) ([]models.FetchPeerV2Data, error)
	FetchFinancialsDataV3(req models.FetchFinancialsReq) (models.FetchFinancialsV3Res, error)
	FetchFinancialsDataV4(req models.FetchFinancialsReq) (models.FetchFinancialsV4Res, error)
	FetchCompanyMaster() ([]models.CompanyDetails, error)
	InsertTransactionData(payoutDetails models.PayoutDetails) error
	UpdateTransactionData(transactionID string, updates map[string]interface{}) error
	GetTransactionData(transactionID string) (*models.PayoutDetails, error)
	CheckExistingPayoutRequest(clientID string) (bool, error)
	InsertPledgeData(pledgeData models.PledgeData) (int64, error)
	FetchCorporateAnnouncements(req models.FetchCorporateActionsIndividualReq) ([]models.CorporateAnnouncements, error)
	FetchCorporateAnnouncementsAll(req models.FetchCorporateActionsAllReq) ([]models.CorporateAnnouncements, error)
	FetchSectorWiseCompanyDataV2(sectorCode []string) ([]models.SectorWiseCompanyV2, error)
	GetSectorWiseCompanyList(page int, sectorName string) ([]models.SectorWiseCompany, error)
}

type MongoDatabase interface {
	InitMongoClient(env string) error
	GetMongoStatus() error
	FindOneMongo(collectionName string, filter interface{}, result interface{}) error
	FindManyMongo(collectionName string, filter interface{}, results interface{}) error
	UpdateOneMongo(collectionName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) error
	UpdateOneMongoDao(collectionName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) error
	DeleteOneMongo(collectionName string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	FindAllMongo(collectionName string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	InsertOneMongo(collectionName string, document interface{}) (*mongo.InsertOneResult, error)
	DeleteMany(collectionName string, filter interface{}, opts ...*options.DeleteOptions) error
	InsertMany(collectionName string, documents []interface{}, opts ...*options.InsertManyOptions) error
	FindOneMongoDao(collectionName string, filter interface{}, result interface{}) error
}
