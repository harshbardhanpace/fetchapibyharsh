package constants

var (
	UseMCXLtpUrl bool
	ShilpiURL    string

	CmURL  string
	CmAuth string

	FinvuBaseUrl          string
	FinvuRidDetails       string
	FinvuTsDetails        string
	FinvuChannelIdDetails string
	FinvuPurposeRefURI    string
	FinvuRedirectURL      string

	FreshDeskBaseUrl string
	FreshDeskApiKey  string
	FreshDeskPass    string

	TLURL      string
	TLOauthUrl string

	LocalCachingCallEnabled bool

	CheckDisplayNameFlag bool

	MongoURI               string
	MongoBase              string
	MongoPass              string
	ConfigSpaceDB          string
	ConfigDaoDB            string
	ConfigContractSearchDB string

	ServerIP           string
	Port               string
	User               string
	Password           string
	DB                 string
	DBMaxConn          int
	DBMaxIdleConn      int
	CertificateEnabled string
	CertificatePath    string

	Msg91FlowID  string
	Msg91Url     string
	AuthKeyMsg91 string

	RabbitMqUser      string
	RabbitMqPassword  string
	RabbitMqAddress   string
	RabbitMqHeartbeat int

	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSBucketName      string

	RedisURL              string
	OrderRedisURL         string
	ContractCacheAddr     string
	ContractCachePoolSize int
	ContractCachePassword string
	SmartCacheAddr        string
	SmartCacheUsername    string
	SmartCachePassword    string

	FileLoggingEnabled bool
	LogFilePath        string

	AdminSecretKey string

	TokenCacheTime int

	LedgerS3FolderName             string
	TradebookS3FolderName          string
	OpenPositionS3FolderName       string
	FnoPnlS3FolderName             string
	HoldingFinancialS3FolderName   string
	CommodityTradebookS3FolderName string
	FnoTradebookS3FolderName       string
	DpChargesS3FolderName          string

	RedisUrl      string
	OrderRedisUrl string
)

var (
	EquityDeliveryBrokerage                     float64
	EquityDeliveryBrokeragePocketful            float64
	EquityDeliverySttOrCtt                      float64
	EquityDeliveryTransactionChargeNse          float64
	EquityDeliveryTransactionChargeBse          float64
	EquityDeliverySebiCharges                   float64
	EquityDeliveryGst                           float64
	EquityDeliveryStampProcessBuy               float64
	EquityIntradayBrokerageOption1              float64
	EquityIntradayBrokerageOption2              float64
	EquityIntradaySttOrCttSell                  float64
	EquityIntradayTransactionChargeNse          float64
	EquityIntradayTransactionChargeBse          float64
	EquityIntradaySebiCharges                   float64
	EquityIntradayGst                           float64
	EquityIntradayStampProcessBuy               float64
	EquityFuturesBrokerageOption1               float64
	EquityFuturesBrokerageOption2               float64
	EquityFuturesSttOrCttSell                   float64
	EquityFuturesTransactionChargeNse           float64
	EquityFuturesSebiCharges                    float64
	EquityFuturesGst                            float64
	EquityFuturesStampProcessBuy                float64
	EquityOptionsBrokerage                      float64
	EquityOptionsSttOrCttSell                   float64
	EquityOptionsTransactionChargeNse           float64
	EquityOptionsSebiCharges                    float64
	EquityOptionsGst                            float64
	EquityOptionsStampProcessBuy                float64
	CurrencyFuturesBrokerageOption1             float64
	CurrencyFuturesBrokerageOption2             float64
	CurrencyFuturesSttOrCttSell                 float64
	CurrencyFuturesTransactionChargeNse         float64
	CurrencyFuturesTransactionChargeBse         float64
	CurrencyFuturesSebiCharges                  float64
	CurrencyFuturesGst                          float64
	CurrencyFuturesStampProcessBuy              float64
	CurrencyOptionsBrokerage                    float64
	CurrencyOptionsSttOrCttSell                 float64
	CurrencyOptionsTransactionChargeNse         float64
	CurrencyOptionsTransactionChargeBse         float64
	CurrencyOptionsSebiCharges                  float64
	CurrencyOptionsGst                          float64
	CurrencyOptionsStampProcessBuy              float64
	CommodityFuturesBrokerageOption1            float64
	CommodityFuturesBrokerageOption2            float64
	CommodityFuturesSttOrCttNonAgriSell         float64
	CommodityFuturesTransactionChargeNormal     float64
	CommodityFuturesTransactionChargeCastorseed float64
	CommodityFuturesTransactionChargeKapas      float64
	CommodityFuturesTransactionChargePepper     float64
	CommodityFuturesTransactionChargeRbdmolein  float64
	CommodityFuturesGst                         float64
	CommodityFuturesSebiChargesAgri             float64
	CommodityFuturesSebiChargesNonAgri          float64
	CommodityFuturesStampProcessBuy             float64
	CommodityOptionsBrokerage                   float64
	CommodityOptionsSttOrCttSell                float64
	CommodityOptionsTransactionCharge           float64
	CommodityOptionsSebiCharges                 float64
	CommodityOptionsGst                         float64
	CommodityOptionsStampProcessBuy             float64
	DpCharges                                   float64
	KafkaEnable                                 bool
)
