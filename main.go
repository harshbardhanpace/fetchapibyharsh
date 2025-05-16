package main

import (
	"fmt"
	"log"
	"os"
	"space/base"
	srv "space/business/blockdeals"
	"space/constants"
	"space/db"
	"space/dbops"
	"space/health"
	"space/loggerconfig"
	"space/utils"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	// "business/charges"
	_ "space/docs"
	"space/helpers"
	"space/helpers/cache"
	"space/models"
	"space/routers"

	v1 "space/controllers/api/v1"

	// test "space/business/testing"

	testing "space/business/testing"

	"github.com/pocketful-tech/throttling"
)

// @title Space Middleware service
// @version 1.0
// @description ## Middleware API's
// @description Authorization e.g. => 'Authorization':'Bearer XXXX'
// @description P-Operating-System e.g. => 'P-Operating-System':'Ubuntu-18.08'
// @description P-Appname e.g. => 'P-Appname':'space'
// @description P-DeviceType e.g. => 'P-DeviceType':'andriod'
// @termsOfService https://swagger.io/terms/

//Execution starts from main function
//-------------------------------------------------------------------

// Execution starts from main function
func main() {
	defer models.HandlePanic()

	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	r := routers.SetupRouter()

	port := os.Getenv("port")

	// For run on requested port
	if len(os.Args) > 1 {
		reqPort := os.Args[1]
		if reqPort != "" {
			port = reqPort
		}
	}

	location, err := time.LoadLocation(constants.ASIAKOLKATA)
	if err != nil {
		logrus.Panic("Alert Severity:P0-Critical, unable to load the location Asia/Kolkata ", err)
	}

	constants.LocationKolkata = location

	loggerconfig.GetConfig().AddConfigPath("./resources")

	loggerconfig.Start()

	env := os.Getenv("GO_ENV")
	if env == "local" {
		loggerconfig.FetchLocalCreds()
	}
	loggerconfig.SetEnv(env)
	utils.InitConfig(env)

	loggerconfig.LogrusInitialize()

	factory := dbops.NewDatabaseConnectorFactory()
	redisUrl := loggerconfig.GetConfig().GetString("config.normal." + env + ".redisUrl")
	if len(redisUrl) == 0 || !strings.Contains(redisUrl, ":") {
		log.Fatal("invalid redis URL format: must contain a non-empty string with a ':' separator")
	}
	redishost, redisport := strings.Split(redisUrl, ":")[0], strings.Split(redisUrl, ":")[1]
	redisConfig := dbops.DatabaseConfig{
		Type:     dbops.RedisType,
		Host:     redishost,
		Port:     redisport,
		Database: "RedisMain",
	}

	redisConnector, err := factory.GetConnector(redisConfig)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisConnector.Close()

	dbops.RedisRepo = redisConnector.GetRepository().(dbops.RedisRepository)

	orderRedisUrl := loggerconfig.GetConfig().GetString("config.normal." + env + ".orderRedisUrl")
	if len(orderRedisUrl) == 0 || !strings.Contains(orderRedisUrl, ":") {
		log.Fatal("invalid redis URL format: must contain a non-empty string with a ':' separator")
	}
	roorderRedishost, orderRedisport := strings.Split(orderRedisUrl, ":")[0], strings.Split(orderRedisUrl, ":")[1]
	orderRedisConfig := dbops.DatabaseConfig{
		Type:     dbops.RedisType,
		Host:     roorderRedishost,
		Port:     orderRedisport,
		Database: "RedisMain",
	}

	orderRedisConnector, err := factory.GetConnector(orderRedisConfig)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer orderRedisConnector.Close()

	dbops.OrderRedisRepo = orderRedisConnector.GetRepository().(dbops.RedisRepository)

	mongoURI := loggerconfig.GetConfig().GetString("config.secret." + env + ".mongo.mongo_uri")
	mongoBase := loggerconfig.GetConfig().GetString("config.secret." + env + ".mongo.mongo_base")
	mongoPass := loggerconfig.GetConfig().GetString("config.secret." + env + ".mongo.mongo_password")
	dbName := loggerconfig.GetConfig().GetString("config.secret." + env + ".mongo.mongo_space")
	if env == constants.LocalEnv {
		mongoURI = loggerconfig.LocalCreds.Local.Mongo.MongoURI
		mongoBase = loggerconfig.LocalCreds.Local.Mongo.MongoBase
		mongoPass = loggerconfig.LocalCreds.Local.Mongo.MongoPassword
		dbName = loggerconfig.LocalCreds.Local.Mongo.MongoSpace
	}
	mongoConfig := dbops.DatabaseConfig{
		Type:     dbops.MongoType,
		Host:     mongoBase,
		Port:     mongoURI,
		Username: "",
		Password: mongoPass,
		Database: dbName,
	}

	mongoConnector, err := factory.GetConnector(mongoConfig)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB Space: %v", err)
	}
	defer mongoConnector.Close()

	dbops.MongoRepo = mongoConnector.GetRepository().(dbops.MongoRepository)

	dbDaoName := loggerconfig.GetConfig().GetString("config.secret." + env + ".mongo.mongo_dao")
	if env == constants.LocalEnv {
		dbDaoName = loggerconfig.LocalCreds.Local.Mongo.MongoDao
	}

	mongoConfig = dbops.DatabaseConfig{
		Type:     dbops.MongoType,
		Host:     mongoBase,
		Port:     mongoURI,
		Username: "",
		Password: mongoPass,
		Database: dbDaoName,
	}

	mongoDaoConnector, err := factory.GetConnector(mongoConfig)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB DAO: %v", err)
	}
	defer mongoDaoConnector.Close()

	dbops.MongoDaoRepo = mongoDaoConnector.GetRepository().(dbops.MongoRepository)

	dbContractSearchName := constants.ConfigContractSearchDB
	if env == constants.LocalEnv {
		dbContractSearchName = loggerconfig.LocalCreds.Local.Mongo.MongoContractSearch
	}

	mongoConfig = dbops.DatabaseConfig{
		Type:     dbops.MongoType,
		Host:     mongoBase,
		Port:     mongoURI,
		Username: "",
		Password: mongoPass,
		Database: dbContractSearchName,
	}

	mongoContractSearchConnector, err := factory.GetConnector(mongoConfig)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, Failed to connect to MongoDB Contract Search:", err)
		// Continue without search service - other services can still work
	} else {
		defer mongoContractSearchConnector.Close()
		dbops.MongoContractSearchRepo = mongoContractSearchConnector.GetRepository().(dbops.MongoRepository)
	}

	configPath := "config.secret." + env + ".postgres."
	postgresserverIp := loggerconfig.GetConfig().GetString(configPath + "host")
	postgresport := loggerconfig.GetConfig().GetString(configPath + "port")
	postgresuser := loggerconfig.GetConfig().GetString(configPath + "user")
	postgrespassword := loggerconfig.GetConfig().GetString(configPath + "password")
	postgresdb := loggerconfig.GetConfig().GetString(configPath + "db")
	postgresdbMaxConn := loggerconfig.GetConfig().GetInt(configPath + "max_conn")
	postgresdbMaxIdleConn := loggerconfig.GetConfig().GetInt(configPath + "max_idle_conn")

	if env == constants.LocalEnv {
		postgresserverIp = loggerconfig.LocalCreds.Local.Postgres.Host
		postgresport = strconv.Itoa(loggerconfig.LocalCreds.Local.Postgres.Port)
		postgresuser = loggerconfig.LocalCreds.Local.Postgres.User
		postgrespassword = loggerconfig.LocalCreds.Local.Postgres.Password
		postgresdb = loggerconfig.LocalCreds.Local.Postgres.Db
		postgresdbMaxConn = loggerconfig.LocalCreds.Local.Postgres.MaxConn
		postgresdbMaxIdleConn = loggerconfig.LocalCreds.Local.Postgres.MaxIdleConn
	}

	postgresConfig := dbops.DatabaseConfig{
		Type:               dbops.PostgresType,
		Host:               postgresserverIp,
		Port:               postgresport,
		Username:           postgresuser,
		Password:           postgrespassword,
		Database:           postgresdb,
		MaxConnection:      postgresdbMaxConn,
		MaxIdleConnections: postgresdbMaxIdleConn,
	}

	postgresConnector, err := factory.GetConnector(postgresConfig)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer postgresConnector.Close()

	dbops.PostgresRepo = postgresConnector.GetRepository().(dbops.PostgresRepository)

	postgresRepo := postgresConnector.GetRepository().(*dbops.PostgresRepositoryImpl)
	blockdealsobj := srv.NewBlockDealService(postgresRepo.GetDB())
	v1.NewBlockDealController(blockdealsobj)

	var mongodb db.MongoDatabase = &db.MongoDb{}
	db.SetMongoDBObj(mongodb)

	err = mongodb.InitMongoClient(env)
	if err != nil {
		loggerconfig.Panic("Alert Severity:P0-Critical, unable to init mongo client=", err)
	}

	// utils.GetLoggerObj(utils.LOG_FILE_NAME)

	elasticApmName := os.Getenv(constants.ElasticApmServiceName)
	log.Printf("the elastic apm name : %s ", elasticApmName)
	elasticApmEnvironment := os.Getenv(constants.ElasticApmEnvironment)
	log.Printf("the elastic apm environment : %s ", elasticApmEnvironment)
	elasticApmUrl := os.Getenv(constants.ElasticApmServerURL)
	log.Printf("the elastic apm url : %s ", elasticApmUrl)

	loggerconfig.Info("starting space main")

	helpers.StartAwsSession()

	err = helpers.InitializeRabbitMq()
	if err != nil {
		loggerconfig.Panic("Alert Severity:P0-Critical, Failed to initialize RabbitMQ:", err)
	}

	testingReqPerSec := loggerconfig.GetConfig().GetInt("config.normal." + env + ".testingReqPerSec")
	testing := loggerconfig.GetConfig().GetString("config.normal." + env + ".testing")
	space := loggerconfig.GetConfig().GetString("config.normal." + env + ".space")
	testingThrottleObj := throttling.NewAPIThrottler(testingReqPerSec, space, testing)

	//cache
	var redisClient cache.RedisCache = &cache.RedisClient{}
	redisErr := redisClient.InitRedis(constants.Redis)
	if redisErr != nil {
		loggerconfig.Panic("Alert Severity:P0-Critical, Failed to initialize Redis:", redisErr.Error())
	}
	cache.SetRedisClientObj(redisClient)

	//contract_cache
	var contractCacheClient cache.ContractCache = &cache.ContractCacheRedisClient{}
	contractRedisErr := contractCacheClient.ContractCacheInit()
	if contractRedisErr != nil {
		loggerconfig.Panic("Alert Severity:P0-Critical, Failed to initialize ContractCache:", contractRedisErr.Error())
	}
	cache.SetContractCacheClienttObj(contractCacheClient)

	//Smart_cache
	var smartCacheClient cache.SmartCache = &cache.SmartCacheRedisClient{}
	smartRedisErr := smartCacheClient.InitSmartCache()
	if smartRedisErr != nil {
		loggerconfig.Panic("Alert Severity:P0-Critical, Failed to initialize ContractCache:", smartRedisErr.Error())
	}
	cache.SetSmartCacheClienttObj(smartCacheClient)

	err = helpers.ProcessAndUploadCalendar()
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, ProcessAndUploadCalendar, error:", err)
	}

	//initalise provider
	base.InitProviders(mongodb, redisClient, contractCacheClient, smartCacheClient)

	//build login provider
	testProvider := BuildTestProvider(testingThrottleObj)
	v1.InitTestProvider(testProvider)

	attempts := loggerconfig.GetConfig().GetInt("config.normal." + env + ".reconnecttries")
	waitTime := loggerconfig.GetConfig().GetInt("config.normal." + env + ".reconnectwait")
	constants.LatencyThresholdLow = loggerconfig.GetConfig().GetInt64("config.normal." + env + ".latencyThresholdLow")
	constants.LatencyThresholdHigh = loggerconfig.GetConfig().GetInt64("config.normal." + env + ".latencyThresholdHigh")

	go health.CheckConnection(env, attempts, waitTime)

	if port == "" {
		port = "8082" //localhost
	}
	type Job interface {
		Run()
	}

	r.Run(":" + port)

}

func BuildTestProvider(testingThrottleObj *throttling.APIThrottler) models.TestProvider {
	//based on vendor provider can be initialized here
	return testing.InitTesting(testingThrottleObj)
}
