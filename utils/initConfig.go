package utils

import (
	"space/constants"
	"space/loggerconfig"
	"strconv"
)

func InitConfig(env string) {
	normalPath := "config.normal." + env
	secretPath := "config.secret." + env

	// order
	constants.UseMCXLtpUrl = loggerconfig.GetConfig().GetBool(normalPath + ".useMCXLtpUrl")
	constants.KafkaEnable = loggerconfig.GetConfig().GetBool(normalPath + ".kafkaEnable")

	// Shilpi and CMOTS settings
	constants.ShilpiURL = loggerconfig.GetConfig().GetString(normalPath + ".shilpiBaseUrl")
	constants.CmURL = loggerconfig.GetConfig().GetString(normalPath + ".cmotsbaseurl")
	constants.CmAuth = loggerconfig.GetConfig().GetString(secretPath + ".cmotsAuthToken")
	if env == constants.LocalEnv {
		constants.CmAuth = loggerconfig.LocalCreds.Local.CmotsAuthToken
	}

	// Finvu settings
	constants.FinvuBaseUrl = loggerconfig.GetConfig().GetString(normalPath + ".FinvuBaseUrl")
	constants.FinvuRidDetails = loggerconfig.GetConfig().GetString(secretPath + ".FinvuRid")
	constants.FinvuTsDetails = loggerconfig.GetConfig().GetString(normalPath + ".FinvuTs")
	constants.FinvuChannelIdDetails = loggerconfig.GetConfig().GetString(normalPath + ".FinvuChannelId")
	constants.FinvuPurposeRefURI = loggerconfig.GetConfig().GetString(normalPath + ".FinvuPurposeRefURI")
	constants.FinvuRedirectURL = loggerconfig.GetConfig().GetString(normalPath + ".FinvuRedirectURL")

	// Freshdesk settings
	constants.FreshDeskBaseUrl = loggerconfig.GetConfig().GetString(normalPath + ".freshdeskBaseUrl")
	constants.FreshDeskApiKey = loggerconfig.GetConfig().GetString(secretPath + ".freshDeskApiKey")
	constants.FreshDeskPass = loggerconfig.GetConfig().GetString(secretPath + ".freshDeskPass")
	if env == constants.LocalEnv {
		constants.FreshDeskApiKey = loggerconfig.LocalCreds.Local.FreshDeskApiKey
		constants.FreshDeskPass = loggerconfig.LocalCreds.Local.FreshDeskPass
	}

	// MongoDB settings
	constants.MongoURI = loggerconfig.GetConfig().GetString(secretPath + ".mongo.mongo_uri")
	constants.MongoBase = loggerconfig.GetConfig().GetString(secretPath + ".mongo.mongo_base")
	constants.MongoPass = loggerconfig.GetConfig().GetString(secretPath + ".mongo.mongo_password")
	constants.ConfigSpaceDB = loggerconfig.GetConfig().GetString(secretPath + ".mongo.mongo_space")
	constants.ConfigDaoDB = loggerconfig.GetConfig().GetString(secretPath + ".mongo.mongo_dao")
	constants.ConfigContractSearchDB = loggerconfig.GetConfig().GetString(secretPath + ".mongo.mongo_contract_search")
	if env == constants.LocalEnv {
		constants.MongoURI = loggerconfig.LocalCreds.Local.Mongo.MongoURI
		constants.MongoBase = loggerconfig.LocalCreds.Local.Mongo.MongoBase
		constants.MongoPass = loggerconfig.LocalCreds.Local.Mongo.MongoPassword
		constants.ConfigSpaceDB = loggerconfig.LocalCreds.Local.Mongo.MongoSpace
		constants.ConfigDaoDB = loggerconfig.LocalCreds.Local.Mongo.MongoDao
		constants.ConfigContractSearchDB = loggerconfig.LocalCreds.Local.Mongo.MongoContractSearch
	}

	// Database settings
	if env == constants.LocalEnv {
		constants.ServerIP = loggerconfig.LocalCreds.Local.Postgres.Host
		constants.Port = strconv.Itoa(loggerconfig.LocalCreds.Local.Postgres.Port)
		constants.User = loggerconfig.LocalCreds.Local.Postgres.User
		constants.Password = loggerconfig.LocalCreds.Local.Postgres.Password
		constants.DB = loggerconfig.LocalCreds.Local.Postgres.Db
		constants.DBMaxConn = loggerconfig.LocalCreds.Local.Postgres.MaxConn
		constants.DBMaxIdleConn = loggerconfig.LocalCreds.Local.Postgres.MaxIdleConn
	} else {
		constants.ServerIP = loggerconfig.GetConfig().GetString(secretPath + ".postgres.host")
		constants.Port = loggerconfig.GetConfig().GetString(secretPath + ".postgres.port")
		constants.User = loggerconfig.GetConfig().GetString(secretPath + ".postgres.user")
		constants.Password = loggerconfig.GetConfig().GetString(secretPath + ".postgres.password")
		constants.DB = loggerconfig.GetConfig().GetString(secretPath + ".postgres.db")
		constants.DBMaxConn = loggerconfig.GetConfig().GetInt(secretPath + ".postgres.max_conn")
		constants.DBMaxIdleConn = loggerconfig.GetConfig().GetInt(secretPath + ".postgres.max_idle_conn")
	}

	constants.CertificateEnabled = loggerconfig.GetConfig().GetString(normalPath + ".certificateEnabled")
	constants.CertificatePath = loggerconfig.GetConfig().GetString(normalPath + ".postgresPenFilePath")

	constants.TLURL = loggerconfig.GetConfig().GetString(normalPath + ".tradelabBaseUrl")
	constants.TLOauthUrl = loggerconfig.GetConfig().GetString("config.normal." + env + ".oauthUrl")

	constants.DpChargesS3FolderName = loggerconfig.GetConfig().GetString(normalPath + constants.ReportsFolderName + ".DpCharges")
	constants.TradebookS3FolderName = loggerconfig.GetConfig().GetString(normalPath + constants.ReportsFolderName + ".tradebook")
	constants.LedgerS3FolderName = loggerconfig.GetConfig().GetString(normalPath + constants.ReportsFolderName + ".ledger")
	constants.OpenPositionS3FolderName = loggerconfig.GetConfig().GetString(normalPath + constants.ReportsFolderName + ".OpenPosition")
	constants.FnoPnlS3FolderName = loggerconfig.GetConfig().GetString(normalPath + constants.ReportsFolderName + ".FnoPln")
	constants.HoldingFinancialS3FolderName = loggerconfig.GetConfig().GetString(normalPath + constants.ReportsFolderName + ".HoldingFinancial")
	constants.CommodityTradebookS3FolderName = loggerconfig.GetConfig().GetString(normalPath + constants.ReportsFolderName + ".CommodityTradebook")
	constants.FnoTradebookS3FolderName = loggerconfig.GetConfig().GetString(normalPath + constants.ReportsFolderName + ".FnoTradebook")

	constants.LocalCachingCallEnabled = loggerconfig.GetConfig().GetBool(normalPath + ".LocalCachingCallEnabled")
	constants.CheckDisplayNameFlag = loggerconfig.GetConfig().GetBool(normalPath + ".displayNameCheck")

	// Redis Cache settings
	constants.RedisUrl = loggerconfig.GetConfig().GetString(normalPath + ".redisUrl")
	constants.OrderRedisUrl = loggerconfig.GetConfig().GetString(normalPath + ".orderRedisUrl")

	constants.ContractCacheAddr = loggerconfig.GetConfig().GetString(secretPath + ".contractCache.addr")
	constants.ContractCachePoolSize = loggerconfig.GetConfig().GetInt(secretPath + ".contractCache.poolSize")
	constants.ContractCachePassword = loggerconfig.GetConfig().GetString(secretPath + ".contractCache.password")
	if env == constants.LocalEnv {
		constants.ContractCacheAddr = loggerconfig.LocalCreds.Local.ContractCache.Addr
		constants.ContractCachePoolSize = loggerconfig.LocalCreds.Local.ContractCache.PoolSize
		constants.ContractCachePassword = loggerconfig.LocalCreds.Local.ContractCache.Password
	}

	constants.SmartCacheAddr = loggerconfig.GetConfig().GetString(secretPath + ".smartCache.addr")
	constants.SmartCacheUsername = loggerconfig.GetConfig().GetString(secretPath + ".smartCache.userName")
	constants.SmartCachePassword = loggerconfig.GetConfig().GetString(secretPath + ".smartCache.password")
	if env == constants.LocalEnv {
		constants.SmartCacheAddr = loggerconfig.LocalCreds.Local.SmartCache.Addr
		constants.SmartCacheUsername = loggerconfig.LocalCreds.Local.SmartCache.Username
		constants.SmartCachePassword = loggerconfig.LocalCreds.Local.SmartCache.Password
	}

	constants.AWSRegion = loggerconfig.GetConfig().GetString(secretPath + ".awsS3CredConfig.myRegion")
	constants.AWSAccessKeyID = loggerconfig.GetConfig().GetString(secretPath + ".awsS3CredConfig.AccessKeyId")
	constants.AWSSecretAccessKey = loggerconfig.GetConfig().GetString(secretPath + ".awsS3CredConfig.SecretAccessKey")
	constants.AWSBucketName = loggerconfig.GetConfig().GetString(secretPath + ".awsS3CredConfig.BucketNamePocket")
	if env == constants.LocalEnv {
		constants.AWSRegion = loggerconfig.LocalCreds.Local.AwsS3CredConfig.MyRegion
		constants.AWSAccessKeyID = loggerconfig.LocalCreds.Local.AwsS3CredConfig.AccessKeyID
		constants.AWSSecretAccessKey = loggerconfig.LocalCreds.Local.AwsS3CredConfig.SecretAccessKey
		constants.AWSBucketName = loggerconfig.LocalCreds.Local.AwsS3CredConfig.BucketName

	}

	constants.Msg91FlowID = loggerconfig.GetConfig().GetString(normalPath + ".otpTemplate")
	constants.Msg91Url = loggerconfig.GetConfig().GetString(normalPath + ".msg91SendSmsUrl")
	constants.AuthKeyMsg91 = loggerconfig.GetConfig().GetString(secretPath + ".authKeyMsg91")
	if env == constants.LocalEnv {
		constants.AuthKeyMsg91 = loggerconfig.LocalCreds.Local.AuthKeyMsg91
	}

	constants.RabbitMqUser = loggerconfig.GetConfig().GetString(secretPath + ".rabbitMq.user")
	constants.RabbitMqPassword = loggerconfig.GetConfig().GetString(secretPath + ".rabbitMq.password")
	constants.RabbitMqAddress = loggerconfig.GetConfig().GetString(secretPath + ".rabbitMq.addr")
	constants.RabbitMqHeartbeat = loggerconfig.GetConfig().GetInt(normalPath + ".rabbitMQHeartbeat")
	if env == constants.LocalEnv {
		constants.RabbitMqUser = loggerconfig.LocalCreds.Local.RabbitMq.User
		constants.RabbitMqPassword = loggerconfig.LocalCreds.Local.RabbitMq.Password
		constants.RabbitMqAddress = loggerconfig.LocalCreds.Local.RabbitMq.Addr
	}
	constants.TokenCacheTime = loggerconfig.GetConfig().GetInt(normalPath + ".tokenCacheTime")

	constants.AdminSecretKey = loggerconfig.GetConfig().GetString(secretPath + ".adminSecretKey")
	if env == constants.LocalEnv {
		constants.AdminSecretKey = loggerconfig.LocalCreds.Local.AdminSecretKey
	}

	constants.FileLoggingEnabled = loggerconfig.GetConfig().GetBool(normalPath + ".fileLoggingEnabled")
	constants.LogFilePath = loggerconfig.GetConfig().GetString(normalPath + ".logFilePath")

	//init all the changes variables
	initCharges()
}

func initCharges() {
	constants.EquityDeliveryBrokerage = getVal("EquityDeliveryBrokerage")
	constants.EquityDeliveryBrokeragePocketful = getVal("EquityDeliveryBrokeragePocketful")
	constants.EquityDeliverySttOrCtt = getVal("EquityDeliverySttOrCtt")
	constants.EquityDeliveryTransactionChargeNse = getVal("EquityDeliveryTransactionChargeNse")
	constants.EquityDeliveryTransactionChargeBse = getVal("EquityDeliveryTransactionChargeBse")
	constants.EquityDeliverySebiCharges = getVal("EquityDeliverySebiCharges")
	constants.EquityDeliveryGst = getVal("EquityDeliveryGst")
	constants.EquityDeliveryStampProcessBuy = getVal("EquityDeliveryStampProcessBuy")

	constants.EquityIntradayBrokerageOption1 = getVal("EquityIntradayBrokerageOption1")
	constants.EquityIntradayBrokerageOption2 = getVal("EquityIntradayBrokerageOption2")
	constants.EquityIntradaySttOrCttSell = getVal("EquityIntradaySttOrCttSell")
	constants.EquityIntradayTransactionChargeNse = getVal("EquityIntradayTransactionChargeNse")
	constants.EquityIntradayTransactionChargeBse = getVal("EquityIntradayTransactionChargeBse")
	constants.EquityIntradaySebiCharges = getVal("EquityIntradaySebiCharges")
	constants.EquityIntradayGst = getVal("EquityIntradayGst")
	constants.EquityIntradayStampProcessBuy = getVal("EquityIntradayStampProcessBuy")

	constants.EquityFuturesBrokerageOption1 = getVal("EquityFuturesBrokerageOption1")
	constants.EquityFuturesBrokerageOption2 = getVal("EquityFuturesBrokerageOption2")
	constants.EquityFuturesSttOrCttSell = getVal("EquityFuturesSttOrCttSell")
	constants.EquityFuturesTransactionChargeNse = getVal("EquityFuturesTransactionChargeNse")
	constants.EquityFuturesSebiCharges = getVal("EquityFuturesSebiCharges")
	constants.EquityFuturesGst = getVal("EquityFuturesGst")
	constants.EquityFuturesStampProcessBuy = getVal("EquityFuturesStampProcessBuy")

	constants.EquityOptionsBrokerage = getVal("EquityOptionsBrokerage")
	constants.EquityOptionsSttOrCttSell = getVal("EquityOptionsSttOrCttSell")
	constants.EquityOptionsTransactionChargeNse = getVal("EquityOptionsTransactionChargeNse")
	constants.EquityOptionsSebiCharges = getVal("EquityOptionsSebiCharges")
	constants.EquityOptionsGst = getVal("EquityOptionsGst")
	constants.EquityOptionsStampProcessBuy = getVal("EquityOptionsStampProcessBuy")

	constants.CurrencyFuturesBrokerageOption1 = getVal("CurrencyFuturesBrokerageOption1")
	constants.CurrencyFuturesBrokerageOption2 = getVal("CurrencyFuturesBrokerageOption2")
	constants.CurrencyFuturesSttOrCttSell = getVal("CurrencyFuturesSttOrCttSell")
	constants.CurrencyFuturesTransactionChargeNse = getVal("CurrencyFuturesTransactionChargeNse")
	constants.CurrencyFuturesTransactionChargeBse = getVal("CurrencyFuturesTransactionChargeBse")
	constants.CurrencyFuturesSebiCharges = getVal("CurrencyFuturesSebiCharges")
	constants.CurrencyFuturesGst = getVal("CurrencyFuturesGst")
	constants.CurrencyFuturesStampProcessBuy = getVal("CurrencyFuturesStampProcessBuy")

	constants.CurrencyOptionsBrokerage = getVal("CurrencyOptionsBrokerage")
	constants.CurrencyOptionsSttOrCttSell = getVal("CurrencyOptionsSttOrCttSell")
	constants.CurrencyOptionsTransactionChargeNse = getVal("CurrencyOptionsTransactionChargeNse")
	constants.CurrencyOptionsTransactionChargeBse = getVal("CurrencyOptionsTransactionChargeBse")
	constants.CurrencyOptionsSebiCharges = getVal("CurrencyOptionsSebiCharges")
	constants.CurrencyOptionsGst = getVal("CurrencyOptionsGst")
	constants.CurrencyOptionsStampProcessBuy = getVal("CurrencyOptionsStampProcessBuy")

	constants.CommodityFuturesBrokerageOption1 = getVal("CommodityFuturesBrokerageOption1")
	constants.CommodityFuturesBrokerageOption2 = getVal("CommodityFuturesBrokerageOption2")
	constants.CommodityFuturesSttOrCttNonAgriSell = getVal("CommodityFuturesSttOrCttNonAgriSell")
	constants.CommodityFuturesTransactionChargeNormal = getVal("CommodityFuturesTransactionChargeNormal")
	constants.CommodityFuturesTransactionChargeCastorseed = getVal("CommodityFuturesTransactionChargeCastorseed")
	constants.CommodityFuturesTransactionChargeKapas = getVal("CommodityFuturesTransactionChargeKapas")
	constants.CommodityFuturesTransactionChargePepper = getVal("CommodityFuturesTransactionChargePepper")
	constants.CommodityFuturesTransactionChargeRbdmolein = getVal("CommodityFuturesTransactionChargeRbdmolein")
	constants.CommodityFuturesGst = getVal("CommodityFuturesGst")
	constants.CommodityFuturesSebiChargesAgri = getVal("CommodityFuturesSebiChargesAgri")
	constants.CommodityFuturesSebiChargesNonAgri = getVal("CommodityFuturesSebiChargesNonAgri")
	constants.CommodityFuturesStampProcessBuy = getVal("CommodityFuturesStampProcessBuy")

	constants.CommodityOptionsBrokerage = getVal("CommodityOptionsBrokerage")
	constants.CommodityOptionsSttOrCttSell = getVal("CommodityOptionsSttOrCttSell")
	constants.CommodityOptionsTransactionCharge = getVal("CommodityOptionsTransactionCharge")
	constants.CommodityOptionsSebiCharges = getVal("CommodityOptionsSebiCharges")
	constants.CommodityOptionsGst = getVal("CommodityOptionsGst")
	constants.CommodityOptionsStampProcessBuy = getVal("CommodityOptionsStampProcessBuy")
	constants.DpCharges = getVal("DpCharges")
}

func getVal(name string) float64 {
	name = "." + name
	return loggerconfig.GetConfig().GetFloat64("config.normal." + constants.Env + name)
}
