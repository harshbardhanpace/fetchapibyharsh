package models

import "time"

type GetOverviewReq struct {
	Exchange  string `json:"exchange"`
	Isin      string `json:"isin"`
	NseSymbol string `json:"nseSymbol"`
	BseToken  string `json:"bseToken"`
}

type GetOverviewRes struct {
	MarketCap    float64 `json:"marketCap"`
	PbRatio      float64 `json:"pbRatio"`
	PeRatio      float64 `json:"peRatio"`
	IndustryPE   float64 `json:"industryPE"`
	Roe          float64 `json:"roe"`
	Eps          float64 `json:"eps"`
	DivYield     float64 `json:"divYield"`
	BookValue    float64 `json:"bookValue"`
	DebttoEquity float64 `json:"debttoEquity"`
	NetProfit    float64 `json:"netProfit"`
	AboutCompany string  `json:"aboutCompany"`
}

type GetOverviewData struct {
	MarketCap       float64 `json:"marketCap"`
	PbRatio         float64 `json:"pbRatio"`
	PeRatio         float64 `json:"peRatio"`
	IndustryPE      float64 `json:"industryPE"`
	Roe             float64 `json:"roe"`
	Eps             float64 `json:"eps"`
	DivYield        float64 `json:"divYield"`
	BookValue       float64 `json:"bookValue"`
	DebttoEquity    float64 `json:"debttoEquity"`
	NetProfit       float64 `json:"netProfit"`
	CompanyName     string  `json:"companyName"`
	HseSName        string  `json:"hseSName"`
	EstablishedYear string  `json:"establishedYear"`
	Industry        string  `json:"industry"`
	Auditor         string  `json:"auditor"`
	Chairman        string  `json:"chairman"`
	CoSec           string  `json:"coSec"`
	Website         string  `json:"website"`
}

type FetchFinancialsReq struct {
	Exchange  string `json:"exchange"`
	Isin      string `json:"isin"`
	NseSymbol string `json:"nseSymbol"`
	BseToken  string `json:"bseToken"`
}

type FetchFinancialsRes struct {
	NetProfit    FetchFinancialsData `json:"netProfit"`
	Revenue      FetchFinancialsData `json:"revenue"`
	BalanceSheet FetchFinancialsData `json:"balanceSheet"`
	Cashflow     []CashFlowRatios    `json:"cashflow"`
}

type FetchFinancialsData struct {
	CoCode     int     `json:"CoCode"`
	TypeCS     string  `json:"typeCS"`
	Columnname string  `json:"COLUMNNAME"`
	Rid        int     `json:"RID"`
	Yc0        float64 `json:"YC0"`
	Yc1        float64 `json:"YC1"`
	Yc2        float64 `json:"YC2"`
	Yc3        float64 `json:"YC3"`
	Yc4        float64 `json:"YC4"`
	Rowno      int     `json:"rowno"`
}

type CashFlowRatios struct {
	CoCode               float64 `json:"co_code"`
	TypeCS               string  `json:"typeCS"`
	Yrc                  float64 `json:"YRC"`
	CashFlowPerShare     float64 `json:"CashFlowPerShare"`
	PricetoCashFlowRatio float64 `json:"PricetoCashFlowRatio"`
	FreeCashFlowperShare float64 `json:"FreeCashFlowperShare"`
	PricetoFreeCashFlow  float64 `json:"PricetoFreeCashFlow"`
	FreeCashFlowYield    float64 `json:"FreeCashFlowYield"`
	Salestocashflowratio float64 `json:"Salestocashflowratio"`
}

type FetchFinancialsDetailedReq struct {
	Exchange  string `json:"exchange"`
	Isin      string `json:"isin"`
	NseSymbol string `json:"nseSymbol"`
	BseToken  string `json:"bseToken"`
}

type FetchFinancialsDetailedRes struct {
	QuarterlyData []QuarterlyData `json:"quarterlyData"`
}

type QuarterlyData struct {
	CoCode     int     `json:"CoCode"`
	Type       string  `json:"Type"`
	Rid        int     `json:"RID"`
	Columnname string  `json:"COLUMNNAME"`
	Y202212    float64 `json:"Y202212"`
	Y202209    float64 `json:"Y202209"`
	Y202206    float64 `json:"Y202206"`
	Y202203    float64 `json:"Y202203"`
	Y202112    float64 `json:"Y202112"`
	Y202109    float64 `json:"Y202109"`
	Y202106    float64 `json:"Y202106"`
	Y202103    float64 `json:"Y202103"`
	Y202012    float64 `json:"Y202012"`
	Y202009    float64 `json:"Y202009"`
	Y202006    float64 `json:"Y202006"`
	Y202003    float64 `json:"Y202003"`
	Y201912    float64 `json:"Y201912"`
	Y201909    float64 `json:"Y201909"`
	Y201906    float64 `json:"Y201906"`
	Y201903    float64 `json:"Y201903"`
	Y201812    float64 `json:"Y201812"`
	Y201809    float64 `json:"Y201809"`
	Y201806    float64 `json:"Y201806"`
	Y201803    float64 `json:"Y201803"`
	Rowno      int     `json:"rowno"`
}

type FetchPeersReq struct {
	Exchange  string `json:"exchange"`
	Isin      string `json:"isin"`
	NseSymbol string `json:"nseSymbol"`
	BseToken  string `json:"bseToken"`
	FetchBy   string `json:"fetchBy"`
	Sector    string `json:"sector"`
}

type TTMData struct {
	CoCode        int     `json:"co_code"`
	PeTtm         float64 `json:"PE_TTM"`
	DividendYield float64 `json:"DividendYield"`
	RoeTtm        float64 `json:"ROE_TTM"`
	RoceTtm       float64 `json:"ROCE_TTM"`
	FaceValue     int     `json:"FaceValue"`
	BookValue     float64 `json:"BookValue"`
	Mcap          float64 `json:"MCAP"`
	PbTtm         float64 `json:"PB_TTM"`
	EpsTtm        float64 `json:"EPS_TTM"`
	DebttoEquity  float64 `json:"DebttoEquity"`
	EbitdaGrowth  float64 `json:"EbitdaGrowth"`
	SectorPE      float64 `json:"SectorPE"`
}

type FetchPeersRes struct {
	CompanyList []FetchPeerData `json:"companyList"`
}

type FetchPeerData struct {
	Company       string  `json:"company"`
	Exchange      string  `json:"exchange"`
	Filter        float64 `json:"filter"`
	Token         string  `json:"token"`
	TradingSymbol string  `json:"tradingSymbol"`
	SectorCode    string  `json:"sectorCode"`
}

type ShareHoldingPatternsReq struct {
	Exchange  string `json:"exchange"`
	Isin      string `json:"isin"`
	NseSymbol string `json:"nseSymbol"`
	BseToken  string `json:"bseToken"`
	Yrc       int    `json:"yrc"`
}

type ShareHoldingPatternsRes struct {
	CoCode                       int     `json:"co_code"`
	Yrc                          int     `json:"YRC"`
	TotalPromoterShares          float64 `json:"totalPromoterShares"`
	TotalPromoterPerShares       float64 `json:"totalPromoterPerShares"`
	TotalPromoterPledgeShares    float64 `json:"totalPromoterPledgeShares"`
	TotalPromoterPerPledgeShares float64 `json:"totalPromoterPerPledgeShares"`
	TotalNoofShareholders        float64 `json:"totalNoofShareholders"`
	PPIMF                        float64 `json:"ppimf"`
	PPIFII                       float64 `json:"ppifii"`
	PPSUBTOT                     float64 `json:"ppsubtot"`
	Other                        float64 `json:"other"`
}

type RatiosCompareReq struct {
	ReqData []RatioCompareData `json:"reqData"`
}

type RatioCompareData struct {
	Exchange string `json:"exchange"`
	Isin     string `json:"isin"`
	Symbol   string `json:"symbol"`
}

type RatiosCompareRes struct {
	Company    string  `json:"company"`
	MarketCap  float64 `json:"marketCap"`
	PbRatio    float64 `json:"pbRatio"`
	PeRatio    float64 `json:"peRatio"`
	IndustryPE float64 `json:"industryPE"`
	Roe        float64 `json:"roe"`
	Eps        float64 `json:"eps"`
	DivYield   float64 `json:"divYield"`
	BookValue  float64 `json:"bookValue"`
}

type FetchTechnicalIndicatorsReq struct {
	Isin      string `json:"isin"`
	Frequency string `json:"frequency" validate:"oneof=monthly daily weekly"`
}

type FetchTechnicalIndicatorsRes struct {
	CoCode        int     `json:"co_code"`
	Rsi           float64 `json:"RSI"`
	MACD1226Days  float64 `json:"MACD_12_26_Days"`
	Avg20Days     float64 `json:"Avg_20Days"`
	Avg50Days     float64 `json:"Avg_50Days"`
	Avg100Days    float64 `json:"Avg_100Days"`
	Avg200Days    float64 `json:"Avg_200Days"`
	Avg10Days     float64 `json:"Avg_10Days"`
	EMA10Day      float64 `json:"EMA10Day"`
	EMA20Day      float64 `json:"EMA20Day"`
	EMA50Day      float64 `json:"EMA50Day"`
	MACD12269Days float64 `json:"MACD_12_26_9_Days"`
	Exchange      string  `json:"exchange"`
	Frequency     string  `json:"frequency"`
	CoName        string  `json:"co_name"`
	Currprice     float64 `json:"currprice"`
	PivotPoint    float64 `json:"pivotPoint"`
	S1            float64 `json:"S1"`
	S2            float64 `json:"S2"`
	S3            float64 `json:"S3"`
	R1            float64 `json:"R1"`
	R2            float64 `json:"R2"`
	R3            float64 `json:"R3"`
	CurrTime      string  `json:"currTime"`
	ExchangeAlt   string  `json:"exchange_alt"`
}

type StocksOnNewsReq struct {
	Filter string `json:"filter"`
	Limit  int    `json:"limit"`
}

type StocksOnNewsResponseData struct {
	Exchange      string `json:"exchange"`
	Token         string `json:"token"`
	TradingSymbol string `json:"tradingSymbol"`
	CompanyName   string `json:"companyName"`
	Remark        string `json:"remark"`
	Summary       string `jaon:"summary"`
}

type StocksOnNewsResponse struct {
	Data []StocksOnNewsResponseData
}

type CoCodeAndResponsePointer struct {
	CoCode          float64
	ResponsePointer *StocksOnNewsResponseData
}

type DailyAnnouncementResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    []DailyAnnouncement `json:"data"`
}

type DailyAnnouncement struct {
	CoCode  float64 `json:"co_code"`
	Symbol  string  `json:"symbol"`
	CoName  string  `json:"lname"`
	Caption string  `json:"caption"`
	Date    string  `json:"date"`
	Memo    string  `json:"memo"`
}

type BoardMeetingForthComingResponse struct {
	Success bool                      `json:"success"`
	Message string                    `json:"message"`
	Data    []BoardMeetingForthComing `json:"data"`
}

type BoardMeetingForthComing struct {
	CoCode float64 `json:"co_code"`
	CoName string  `json:"co_name"`
	Symbol string  `json:"symbol"`
	Date   string  `json:"date"`
	Note   string  `json:"Note"`
}

type ChangeOfNameResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    []ChangeOfName `json:"data"`
}

type ChangeOfName struct {
	Oldname string  `json:"OLDNAME"`
	CoName  string  `json:"newname"`
	Srdt    string  `json:"srdt"`
	CoCode  float64 `json:"co_code"`
	Symbol  string  `json:"symbol"`
}

type SplitsResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    []Splits `json:"data"`
}

type Splits struct {
	CoName              string  `json:"co_name"`
	CoCode              float64 `json:"co_code"`
	Symbol              string  `json:"symbol"`
	Isin                string  `json:"isin"`
	AnnouncementDate    string  `json:"AnnouncementDate"`
	Recorddate          string  `json:"recorddate"`
	SplitDate           string  `json:"SplitDate"`
	FVBefore            float64 `json:"FVBefore"`
	FVAfter             float64 `json:"FVAfter"`
	Remark              string  `json:"remark"`
	Description         string  `json:"Description"`
	SplitRatio          string  `json:"SplitRatio"`
	NoDeliveryStartDate string  `json:"NoDeliveryStartDate"`
	NoDeliveryEndDate   string  `json:"NoDeliveryEndDate"`
}

type MergerDemergerResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    []MergerDemerger `json:"data"`
}

type MergerDemerger struct {
	CoName             string  `json:"co_name"`
	CoCode             float64 `json:"co_code"`
	Isin               string  `json:"ISIN"`
	AnnouncementDate   string  `json:"AnnouncementDate"`
	MergerDemergerDate string  `json:"Merger_Demerger_Date"`
	Recorddate         string  `json:"recorddate"`
	MgrRatio           string  `json:"mgrRatio"`
	MergedIntoCode     int     `json:"MergedInto_Code"`
	MergedIntoISIN     string  `json:"MergedInto_ISIN"`
	MergedIntoName     string  `json:"MergedIntoName"`
	Type               string  `json:"Type"`
	Instname           string  `json:"INSTNAME"`
}

type DividendAnnouncementDataResponse struct {
	Success bool                       `json:"success"`
	Message string                     `json:"message"`
	Data    []DividendAnnouncementData `json:"data"`
}

type DividendAnnouncementData struct {
	CoName             string  `json:"co_name"`
	CoCode             float64 `json:"co_code"`
	Symbol             string  `json:"symbol"`
	Isin               string  `json:"isin"`
	AnnouncementDate   string  `json:"AnnouncementDate"`
	DivDate            string  `json:"DivDate"`
	RecordDate         string  `json:"RecordDate"`
	DivAmount          float64 `json:"DivAmount"`
	DivPer             int     `json:"DivPer"`
	DividendType       string  `json:"DividendType"`
	Description        string  `json:"Description"`
	DividendPayoutDate string  `json:"DividendPayoutDate"`
}

type BulkDealsResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    []BulkDeals `json:"data"`
}

type BulkDeals struct {
	Scripcode  string  `json:"scripcode"`
	Serial     int     `json:"Serial"`
	CoCode     float64 `json:"CO_CODE"`
	Date       string  `json:"Date"`
	Scripname  string  `json:"scripname"`
	Clientname string  `json:"clientname"`
	Buysell    string  `json:"buysell"`
	Qtyshares  float64 `json:"qtyshares"`
	AvgPrice   float64 `json:"avg_price"`
}

type BlockDealsResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    []BlockDeals `json:"data"`
}

type BlockDeals struct {
	Scripcode  string  `json:"scripcode"`
	Serial     int     `json:"Serial"`
	CoCode     float64 `json:"CO_CODE"`
	Date       string  `json:"Date"`
	ScripName  string  `json:"ScripName"`
	ClientName string  `json:"ClientName"`
	Buysell    string  `json:"BUYSELL"`
	Qtyshares  float64 `json:"QTYSHARES"`
	AvgPrice   float64 `json:"AVG_PRICE"`
}

type AGMResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []AGM  `json:"data"`
}

type AGM struct {
	CoName           string  `json:"co_name"`
	CoCode           float64 `json:"co_code"`
	Symbol           string  `json:"symbol"`
	Isin             string  `json:"isin"`
	AnnouncementDate string  `json:"AnnouncementDate"`
	GMdate           string  `json:"GMdate"`
	Purpose          string  `json:"Purpose"`
	BCloserStartDate string  `json:"BCloserStartDate"`
	BCloserEndDate   string  `json:"BCloserEndDate"`
	Description      string  `json:"Description"`
}

type EGMResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []EGM  `json:"data"`
}

type EGM struct {
	CoName           string  `json:"co_name"`
	CoCode           float64 `json:"co_code"`
	Symbol           string  `json:"symbol"`
	Isin             string  `json:"isin"`
	AnnouncementDate string  `json:"AnnouncementDate"`
	GMdate           string  `json:"GMdate"`
	Purpose          string  `json:"Purpose"`
	BCloserStartDate string  `json:"BCloserStartDate"`
	BCloserEndDate   string  `json:"BCloserEndDate"`
	Description      string  `json:"Description"`
}

type BonusResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Data    []Bonus `json:"data"`
}

type Bonus struct {
	CoName           string  `json:"co_name"`
	CoCode           float64 `json:"co_code"`
	Symbol           string  `json:"symbol"`
	Isin             string  `json:"isin"`
	AnnouncementDate string  `json:"AnnouncementDate"`
	RecordDate       string  `json:"RecordDate"`
	BonusDate        string  `json:"BonusDate"`
	BonusRatio       string  `json:"BonusRatio"`
	Remark           string  `json:"remark"`
	Description      string  `json:"Description"`
}

type FetchTokenAndSymbol struct {
	CoCode        float64 `json:"cocode"`
	Token         string  `json:"bsecode"`
	TradingSymbol string  `json:"nsesymbol"`
	CompanyName   string  `json:"companyname"`
}

type SectorList struct {
	SectCode string `json:"sectCode"`
	SectName string `json:"sectName"`
}

type FetchSectorWiseCompanyReq struct {
	SectCode string `json:"sectCode"`
}

type SectorWiseCompany struct {
	CoCode   int    `json:"coCode"`
	CoName   string `json:"coName"`
	Lname    string `json:"lname"`
	ScCode   string `json:"scCode"`
	Symbol   string `json:"symbol"`
	SectName string `json:"sectName"`
	Isin     string `json:"isin"`
}

type FetchCompanyCategoryReq struct {
	IsinList []string `json:"isinList"`
}

type CompanyCategory struct {
	Isin         string `json:"isin"`
	McapType     string `json:"mcapType"`
	SectorName   string `json:"sectorName"`
	IndustryName string `json:"industryName"`
}

type StocksOnNewsV2Req struct {
	Filter   string   `json:"filter"`
	IsinList []string `json:"isinList"`
}

type StocksAnalyzerReq struct {
	Isin string `json:"isin"`
}

type PLStatement struct {
	ColumnName string  `json:"columnName"`
	Y0         float64 `json:"y0"`
}

type PLStatementResponse struct {
	Data []PLStatement `json:"data"`
}

type BalanceSheets struct {
	ColumnName string  `json:"columnName"`
	Y0         float64 `json:"y0"`
}

type BalanceSheetsResponse struct {
	Data []BalanceSheets `json:"data"`
}

type CashFlow struct {
	ColumnName string  `json:"columnName"`
	Y0         float64 `json:"y0"`
}

type CashflowResponse struct {
	CFO         []CashFlow `json:"cfo"`
	CFI         []CashFlow `json:"cfi"`
	CFF         []CashFlow `json:"cff"`
	NetCashFlow CashFlow   `json:"netCashFlow"`
}

type StocksAnalyzerRes struct {
	PLStatementRes  PLStatementResponse   `json:"plStatementRes"`
	BalanceSheetRes BalanceSheetsResponse `json:"balanceSheetRes"`
	CashFlowRes     CashflowResponse      `json:"cashFlowRes"`
}

type FetchFinancialsDataV2 struct {
	CoCode     int     `json:"CoCode"`
	ColumnName string  `json:"ColumnName"`
	Y0         float64 `json:"Y0"`
	Y1         float64 `json:"Y1"`
	Y2         float64 `json:"Y2"`
	Y3         float64 `json:"Y3"`
	Y4         float64 `json:"Y4"`
}

type BalanceSheetData struct {
	TotalAssets      FetchFinancialsDataV2 `json:"totalAssets"`
	TotalLiabilities FetchFinancialsDataV2 `json:"totalLiabilities"`
}

type FetchFinancialsV2Res struct {
	NetProfit    FetchFinancialsDataV2 `json:"netProfit"`
	Revenue      FetchFinancialsDataV2 `json:"revenue"`
	BalanceSheet BalanceSheetData      `json:"balanceSheet"`
	Cashflow     FetchFinancialsDataV2 `json:"cashflow"`
}

type FetchPeersV2Req struct {
	Exchange  string `json:"exchange"`
	Isin      string `json:"isin"`
	NseSymbol string `json:"nseSymbol"`
	BseToken  string `json:"bseToken"`
	Sector    string `json:"sector"`
}

type FetchPeersV2Res struct {
	CompanyList []FetchPeerV2Data `json:"companyList"`
}

type FetchPeerV2Data struct {
	Company       string  `json:"company"`
	Exchange      string  `json:"exchange"`
	Token         string  `json:"token"`
	TradingSymbol string  `json:"tradingSymbol"`
	SectorCode    string  `json:"sectorCode"`
	Mcap          float64 `json:"MCAP"`
	PbRatio       float64 `json:"pbRatio"`
	PeRatio       float64 `json:"peRatio"`
	Roe           float64 `json:"roe"`
	Eps           float64 `json:"eps"`
}

type FetchFinancialsDataV3 struct {
	CoCode     int     `json:"coCode"`
	ColumnName string  `json:"columnName"`
	Y0         float64 `json:"y0"`
	Y1         float64 `json:"y1"`
	Y2         float64 `json:"y2"`
	Y3         float64 `json:"y3"`
	Y4         float64 `json:"y4"`
	Yrc0       int     `json:"yrc0"`
	Yrc1       int     `json:"yrc1"`
	Yrc2       int     `json:"yrc2"`
	Yrc3       int     `json:"yrc3"`
	Yrc4       int     `json:"yrc4"`
}

type BalanceSheetDataV3 struct {
	TotalAssets      FetchFinancialsDataV3 `json:"totalAssets"`
	TotalLiabilities FetchFinancialsDataV3 `json:"totalLiabilities"`
}
type FetchFinancialsV3Res struct {
	NetProfit    FetchFinancialsDataV3 `json:"netProfit"`
	Revenue      FetchFinancialsDataV3 `json:"revenue"`
	BalanceSheet BalanceSheetDataV3    `json:"balanceSheet"`
	Cashflow     FetchFinancialsDataV3 `json:"cashflow"`
}

type FetchFinancialsDataV4 struct {
	CoCode     int     `json:"coCode"`
	TypeCS     string  `json:"typeCS"`
	ColumnName string  `json:"columnName"`
	Y0         float64 `json:"y0"`
	Y1         float64 `json:"y1"`
	Y2         float64 `json:"y2"`
	Y3         float64 `json:"y3"`
	Y4         float64 `json:"y4"`
	Yrc0       int     `json:"yrc0"`
	Yrc1       int     `json:"yrc1"`
	Yrc2       int     `json:"yrc2"`
	Yrc3       int     `json:"yrc3"`
	Yrc4       int     `json:"yrc4"`
}

type BalanceSheetDataV4 struct {
	TotalAssets      FetchFinancialsDataV4 `json:"totalAssets"`
	TotalLiabilities FetchFinancialsDataV4 `json:"totalLiabilities"`
}

type FetchFinancialsV4Res struct {
	NetProfitConsolidated    FetchFinancialsDataV4 `json:"netProfitConsolidated"`
	RevenueConsolidated      FetchFinancialsDataV4 `json:"revenueConsolidated"`
	BalanceSheetConsolidated BalanceSheetDataV4    `json:"balanceSheetConsolidated"`
	CashflowConsolidated     FetchFinancialsDataV4 `json:"cashflowConsolidated"`
	NetProfitStandalone      FetchFinancialsDataV4 `json:"netProfitStandalone"`
	RevenueStandalone        FetchFinancialsDataV4 `json:"revenueStandalone"`
	BalanceSheetStandalone   BalanceSheetDataV4    `json:"balanceSheetStandalone"`
	CashflowStandalone       FetchFinancialsDataV4 `json:"cashflowStandalone"`
}

type CompanyDetails struct {
	CoCode    int    `json:"coCode"`
	Bsecode   string `json:"bsecode"`
	Nsesymbol string `json:"nsesymbol"`
	Isin      string `json:"isin"`
}

type CorporateAnnouncements struct {
	Cocode           int       `json:"cocode"`
	Coname           string    `json:"coname"`
	Purpose          string    `json:"purpose"`
	Ratio            string    `json:"ratio"`
	AnnouncementDate time.Time `json:"announcementdate"`
	ExecutionDate    time.Time `json:"executiondate"`
	Isin             string    `json:"isin"`
	BseCode          string    `json:"bsecode"`
	NseSymbol        string    `json:"nsesymbol"`
}

type FetchCorporateActionsIndividualReq struct {
	Isin          string `json:"isin"`
	NseSymbol     string `json:"nseSymbol"`
	BseCode       string `json:"bseCode"`
	StartTimeUnix int64  `json:"startTimeUnix"`
	EndTimeUnix   int64  `json:"endTimeUnix"`
	PageNo        int    `json:"pageno"`
}

type FetchCorporateActionsAllReq struct {
	StartTimeUnix int64 `json:"startTimeUnix"`
	EndTimeUnix   int64 `json:"endTimeUnix"`
	PageNo        int   `json:"pageno"`
}

type FetchSectorWiseCompanyReqV2 struct {
	SectCode []string `json:"sectCode"`
}

type Index struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type SectorWiseCompanyDetails struct {
	CoCode    string   `json:"coCode"`
    CoName    string   `json:"coName"`
    Lname     string   `json:"lname"`
    ScCode    string   `json:"scCode"`
    Symbol    string   `json:"symbol"`
    SectName  string   `json:"sectName"`
    Isin      string   `json:"isin"`
	Exchange1 string   `json:"exchange1"`
	Exchange2 string   `json:"exchange2"`
    Token1    string   `json:"token1"`
	Token2    string   `json:"token2"`
}

type SectorWiseCompanyV2 struct {
	Companies  []SectorWiseCompanyDetails `json:"companies"`
	NSEIndices []Index `json:"nseIndices"`
	BSEIndices []Index `json:"bseIndices"`
}
