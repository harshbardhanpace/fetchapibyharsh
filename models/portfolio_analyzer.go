package models

type Holdings struct {
	Exchange              string  `json:"exchange"`
	Isin                  string  `json:"isin"`
	TradingSymbol         string  `json:"tradingSymbol"`
	Token                 string  `json:"token"`
	SectorName            string  `json:"sectorName"`
	SectorCode            string  `json:"sectorCode"`
	PercentageOfPortfolio float64 `json:"percentageOfPortfolio"`
	ValueOfHolding        float64 `json:"valueOfHolding"`
}

type PortfolioBeta struct {
	PortfolioBeta       float64          `json:"portfolioBeta"`
	PortfolioTotalValue float64          `json:"portfolioTotalValue"`
	IndividualBeta      []IndividualBeta `json:"individualBeta"`
}

type IndividualBeta struct {
	Exchange string  `json:"exchange"`
	Isin     string  `json:"isin"`
	Symbol   string  `json:"symbol"`
	Token    string  `json:"token"`
	Beta     float64 `json:"beta"`
}

type PortfolioPE struct {
	PortfolioPE  float64        `json:"portfolioPE"`
	IndividualPE []IndividualPE `json:"individualPE"`
}

type IndividualPE struct {
	Isin          string  `json:"isin"`
	TradingSymbol string  `json:"tradingSymbol"`
	Token         string  `json:"token"`
	Pe            float64 `json:"PE"`
}

type PortfolioDE struct {
	PortfolioDE  float64        `json:"portfolioDE"`
	IndividualDE []IndividualDE `json:"individualDE"`
}

type IndividualDE struct {
	Isin          string  `json:"isin"`
	TradingSymbol string  `json:"tradingSymbol"`
	Token         string  `json:"token"`
	De            float64 `json:"DE"`
}

type RedFlags struct {
	Holdings []RFHoldings `json:"holdings"`
}

type RFHoldings struct {
	Isin          string `json:"isin"`
	TradingSymbol string `json:"tradingSymbol"`
	Token         string `json:"token"`
}

type HighDefaultProbability struct {
	Holdings []HDPHoldings `json:"holdings"`
}

type HDPHoldings struct {
	Isin                 string  `json:"isin"`
	TradingSymbol        string  `json:"tradingSymbol"`
	Token                string  `json:"token"`
	InterseCoverageRatio float64 `json:"interseCoverageRatio"`
	Ebitda               float64 `json:"ebitda"`
	InteresetExpense     float64 `json:"interesetExpense"`
}

type LowROE struct {
	Holdings []LowROEHoldings `json:"holdings"`
}

type LowROEHoldings struct {
	Isin          string  `json:"isin"`
	TradingSymbol string  `json:"tradingSymbol"`
	Token         string  `json:"token"`
	Roe           float64 `json:"roe"`
}

type LowProfitGrowth struct {
	Holdings []LPGHoldings `json:"holdings"`
}

type LPGHoldings struct {
	Isin            string  `json:"isin"`
	TradingSymbol   string  `json:"tradingSymbol"`
	Token           string  `json:"token"`
	AvgProfitGrowth float64 `json:"avgProfitGrowth"`
}

type PortfolioAnalyzerReq struct {
	ClientId string `json:"clientId" binding:"required" example:"abc123"`
}

type HoldingsData struct {
	LTP      float64
	Quantity int
	Isin     string
	Symbol   string
	Token    string
	Exchange string
}

type HighPledgePromoterHoldingRes struct {
	Holdings []HoldingRedFlagData `json:"holdings"`
}

type HoldingRedFlagData struct {
	Isin          string `json:"isin"`
	TradingSymbol string `json:"tradingSymbol"`
	Token         string `json:"token"`
}

type RoeHoldingRedFlagData struct {
	Isin          string  `json:"isin"`
	TradingSymbol string  `json:"tradingSymbol"`
	Token         string  `json:"token"`
	Roe           float64 `json:"roe"`
}

type AllHighPledgePromoterHolding struct {
	HighPledgePromoterHoldingAll []StockDetailDb `json:"highPledgePromoterHoldingAll"`
}

type StockDetailDb struct {
	CompanyName string `json:"companyName"`
	Isin        string `json:"Isin"`
}

type RoeStockDetailDb struct {
	CompanyName string  `json:"companyName"`
	Isin        string  `json:"Isin"`
	Roe         float64 `json:"roe"`
}

type AllIsin struct {
	Isin []string `json:"isin"`
}

type AllAdditionalSurveillanceMeasure struct {
	AdditionalSurveillanceMeasureAll []StockDetailDb `json:"additionalSurveillanceMeasureAll"`
}

type AdditionalSurveillanceMeasureRes struct {
	Holdings []HoldingRedFlagData `json:"holdings"`
}

type AllGradedSurveillanceMeasure struct {
	AllGradedSurveillanceMeasureAll []StockDetailDb `json:"additionalSurveillanceMeasureAll"`
}

type GradedSurveillanceMeasureRes struct {
	Holdings []HoldingRedFlagData `json:"holdings"`
}

type AllLowRoe struct {
	LowRoeAll []RoeStockDetailDb `json:"lowRoeAll"`
}

type LowROERes struct {
	Holdings []RoeHoldingRedFlagData `json:"holdings"`
}

type AllProfitabilityGrowthDb struct {
	ProfitabilityGrowthAll []ProfitabilityGrowthDb `json:"profitabilityGrowthAll"`
}

type ProfitabilityGrowthDb struct {
	CompanyName string  `json:"companyName"`
	Isin        string  `json:"Isin"`
	YZero       float64 `json:"yZero"`
	YFour       float64 `json:"yFour"`
}

type ProfitabilityGrowthRedFlagData struct {
	Isin          string  `json:"isin"`
	TradingSymbol string  `json:"tradingSymbol"`
	Token         string  `json:"token"`
	NetProfit     float64 `json:"netProfit"`
}

type ProfitabilityGrowthRedFlagRes struct {
	Holdings []ProfitabilityGrowthRedFlagData `json:"holdings"`
}

type HoldingStockContributionRes struct {
	Holdings                       []HoldingStockContributionData `json:"holdings"`
	TopFiveHoldingPrice            float64                        `json:"topFiveHoldingPrice"`
	TopFiveCombinedPercentageShare float64                        `json:"topFiveCombinedPercentageShare"`
}

type HoldingStockContributionData struct {
	Isin               string  `json:"isin"`
	TradingSymbol      string  `json:"tradingSymbol"`
	Token              string  `json:"token"`
	StockInvestedPrice float64 `json:"stockInvestedPrice"`
	PercentageShare    float64 `json:"percentageShare"`
}

type InvestmentSectorRes struct {
	InvestmentSector []InvestmentSectorData `json:"investmentSector"`
}

type InvestmentSectorData struct {
	SectorCode         string  `json:"sectorCode"`
	SectorName         string  `json:"sectorName"`
	InvestedValue      float64 `json:"investedValue"`
	InvestedPercentage float64 `json:"investedPercentage"`
}

type CompanyMasterDb struct {
	CompanyMasterAll []CompanyMasterDbData `json:"companyMasterDbData"`
}

type CompanyMasterDbData struct {
	Bsecode    string `json:"bsecode"`
	Nsesymbol  string `json:"nsesymbol"`
	Isin       string `json:"isin"`
	Sectorcode string `json:"sectorcode"`
	Sectorname string `json:"sectorname"`
}

type InvestmentData struct {
	Isin     string  `json:"isin"`
	LTP      float64 `json:"ltp"`
	Quantity int     `json:"quantity"`
}

type DeclineInPromoterHoldingDb struct {
	Isin                string  `json:"isin"`
	CurrentQuarterTPPS  float64 `json:"currentQuarterTPPS"`
	PreviousQuarterTPPS float64 `json:"previousQuarterTPPS"`
}

type AllDeclineInPromoterHoldingDb struct {
	DeclineInPromoterHolding []DeclineInPromoterHoldingDb `json:"DeclineInPromoterHolding"`
}

type DeclineInPromoterHoldingData struct {
	Isin                  string  `json:"isin"`
	TotalPromoterPerShare float64 `json:"totalPromoterPerShare"`
}

type DeclineInPromoterHoldingRedFlagData struct {
	Isin          string  `json:"isin"`
	TradingSymbol string  `json:"tradingSymbol"`
	Token         string  `json:"token"`
	NetDecline    float64 `json:"netDecline"`
}

type DeclineInPromoterHoldingRedFlagRes struct {
	Holdings []DeclineInPromoterHoldingRedFlagData `json:"holdings"`
}

type InterestCoverageRatioDb struct {
	Isin               string  `json:"isin"`
	FinanceCost        float64 `json:"financeCost"`
	ProfitBeforeTax    float64 `json:"profitBeforeTax"`
	InterestCoverRatio float64 `json:"interestCoverRatio"`
}

type AllInterestCoverageRatioDb struct {
	InterestCoverageRatioData []InterestCoverageRatioDb `json:"interestCoverageRatioDb"`
}

type InterestCoverageRatioRedFlagRes struct {
	Holdings []HoldingRedFlagData `json:"holdings"`
}

type RevenueAndProfitData struct {
	Isin string  `json:"isin"`
	Y0   float64 `json:"y0"`
	Y1   float64 `json:"y1"`
	Y2   float64 `json:"y2"`
}

type DeclineInRevenueAndProfitDb struct {
	Isin string `json:"isin"`
}
type AllDeclineInRevenueAndProfitDb struct {
	Holding []DeclineInRevenueAndProfitDb `json:"holdings"`
}

type DeclineInRevenueAndProfitRedFlagRes struct {
	Holdings []HoldingRedFlagData `json:"holdings"`
}

type LowNetWorthDb struct {
	Isin string `json:"isin"`
}
type AllLowNetWorthDb struct {
	Holding []DeclineInRevenueAndProfitDb `json:"holdings"`
}

type LowNetWorthRedFlagRes struct {
	Holdings []HoldingRedFlagData `json:"holdings"`
}

type NetWorthDb struct {
	Isin     string  `json:"isin"`
	NetWorth float64 `josn:"netWorth"`
}
type AllNetWorthDb struct {
	NetWorthData []NetWorthDb `json:"netWorthData"`
}

type NetWorthDataRedFlagData struct {
	Isin          string  `json:"isin"`
	TradingSymbol string  `json:"tradingSymbol"`
	Token         string  `json:"token"`
	NetWorth      float64 `json:"netWorth"`
}

type LowNetWorthDataRedFlagRes struct {
	Holdings []NetWorthDataRedFlagData `json:"holdings"`
}

type DeclineInRevenueDb struct {
	Isin string `json:"isin"`
}
type AllDeclineInRevenueDb struct {
	Holding []DeclineInRevenueDb `json:"holdings"`
}

type RevenueData struct {
	Isin string  `json:"isin"`
	Y0   float64 `json:"y0"`
	Y1   float64 `json:"y1"`
	Y2   float64 `json:"y2"`
	Y3   float64 `json:"y3"`
	Y4   float64 `json:"y4"`
}

type DeclineInRevenueRedFlagRes struct {
	Holdings []HoldingRedFlagData `json:"holdings"`
}

type PromoterPledgeDataDb struct {
	Isin string `json:"isin"`
}
type AllPromoterPledgeDataDb struct {
	Holding []PromoterPledgeDataDb `json:"holdings"`
}

type PromoterPledgeData struct {
	Isin                         string  `json:"isin"`
	TotalPromoterPerPledgeShares float64 `josn:"totalPromoterPerPledgeShares"`
}

type PromoterPledgeRedFlagRes struct {
	Holdings []HoldingRedFlagData `json:"holdings"`
}

type AllPennyStocksDataDb struct {
	Holding []PennyStocksData `json:"holdings"`
}

type PennyStocksData struct {
	Isin      string  `json:"isin"`
	MarketCap float64 `josn:"marketCap"`
}

type PennyStocksRedFlagRes struct {
	Holdings []PennyStocksHoldingRedFlagData `json:"holdings"`
}

type PennyStocksHoldingRedFlagData struct {
	Isin          string  `json:"isin"`
	TradingSymbol string  `json:"tradingSymbol"`
	Token         string  `json:"token"`
	MarketCap     float64 `josn:"marketCap"`
}

type StockReturneRes struct {
	Holdings []HoldingRedFlagData `json:"holdings"`
}

type AllStockReturn struct {
	StockReturn []StockReturn `json:"stockReturn"`
}

type StockReturn struct {
	Isin       string  `json:"isin"`
	ReturnRate float64 `json:"returnRate"`
	McapType   string  `json:"mcapType"`
}

type NiftyIndexData struct {
	TokenId         string  `json:"tokenId"`
	MarketCap       string  `json:"marketCap"`
	MarketCapReturn float64 `json:"marketCapReturn"`
}

type NiftyVsPortfolioReq struct {
	ClientId                  string  `json:"clientId" binding:"required" example:"abc123"`
	MovementOfNiftyPercentage float64 `json:"movementOfNiftyPercentage"`
}

type NiftyVsPortfolioRes struct {
	PortfolioMovement                  float64              `json:"portfolioMovement"`
	ExpectedGainAmountForNiftyMovement float64              `json:"expectedGainAmountForNiftyMovement"`
	IndividualMovement                 []IndividualMovement `json:"individualMovement"`
}

type IndividualMovement struct {
	Exchange string  `json:"exchange"`
	Isin     string  `json:"isin"`
	Symbol   string  `json:"symbol"`
	Token    string  `json:"token"`
	Movement float64 `json:"movement"`
}

type ChangeInInstitutionalHoldingData struct {
	Isin                 string  `json:"isin"`
	Yrc                  int     `json:"yrc"`
	InstitutionalHolding float64 `json:"institutionalHolding"`
}

type ChangeInInstitutionalHoldingDb struct {
	Isin                            string  `json:"isin"`
	DiffrenceInInstitutionalHolding float64 `json:"diffrenceInInstitutionalHolding"`
}

type AllChangeInInstitutionalHoldingDb struct {
	ChangeInInstitutionalHoldingDb []ChangeInInstitutionalHoldingDb `json:"changeInInstitutionalHoldingDb"`
}

type ChangeInInstitutionalHoldingRedFlagData struct {
	Isin                            string  `json:"isin"`
	TradingSymbol                   string  `json:"tradingSymbol"`
	Token                           string  `json:"token"`
	DiffrenceInInstitutionalHolding float64 `json:"diffrenceInInstitutionalHolding"`
}

type ChangeInInstitutionalHoldingRedFlagRes struct {
	Holdings []ChangeInInstitutionalHoldingRedFlagData `json:"holdings"`
}

type RoeAndStockReturnData struct {
	Cocode                 int     `json:"cocode"`
	Isin                   string  `json:"isin"`
	Industrycode           string  `json:"industrycode"`
	Ltp                    float64 `json:"ltp"`
	Return3Yrs             float64 `json:"return3Yrs"`
	Y0ProfitAES            float64 `json:"y0ProfitAES"`
	Y1ProfitAES            float64 `json:"y1ProfitAES"`
	Y2ProfitAES            float64 `json:"y2ProfitAES"`
	Y0TotalShareholderFund float64 `json:"y0TotalShareholderFund"`
	Y1TotalShareholderFund float64 `json:"y1TotalShareholderFund"`
	Y2TotalShareholderFund float64 `json:"y2TotalShareholderFund"`
}

type AllRoeAndStockReturnDb struct {
	RoeAndStockReturndb []RoeAndStockReturnData `json:"roeAndStockReturnData"`
}

type RoeAndStockReturnHoldingRedFlagData struct {
	Isin          string  `json:"isin"`
	TradingSymbol string  `json:"tradingSymbol"`
	Token         string  `json:"token"`
	AvgRoe        float64 `json:"avgRoe"`
	Return3Yrs    float64 `json:"return3Yrs"`
}

type RoeAndStockReturnHoldingRedFlagRes struct {
	RoeAndStockReturnHolding []RoeAndStockReturnHoldingRedFlagData `json:"roeAndStockReturnHolding"`
}

type IlliquidStocksHolding struct {
	Isin           string  `json:"isin"`
	TradingSymbol  string  `json:"tradingSymbol"`
	Token          string  `json:"token"`
	Exchange       string  `json:"exchange"`
	AvgVolumeMonth float64 `json:"avgVolumeMonth"`
	AvgVolumeDay   float64 `json:"avgVolumeDay"`
}

type IlliquidStocksResponse struct {
	Holding []IlliquidStocksHolding `json:"holding"`
}
