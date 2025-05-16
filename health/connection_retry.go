package health

import (
	"space/constants"
	"space/db"
	"space/dbops"
	"space/helpers"
	"space/helpers/cache"
	"space/loggerconfig"
	"space/models"
	"strings"
	"sync"
	"time"
)

func CheckConnection(env string, attempts, waitTime int) {
	loggerconfig.Info("Started Connection check and retry service, env:", env, " attempts:", attempts, " wait time:", waitTime)
	var wg sync.WaitGroup
	pgObj := db.GetPostgresObject()
	redisCliObj := cache.GetRedisClientObj()
	contractCacheClientObj := cache.GetContractCacheClientObj()
	smartCacheClientObj := cache.GetSmartCacheClientObj()

	// Wrap function to handle panic
	withRecovery := func(passFunc func()) func() {
		return func() {
			defer models.HandlePanic()
			passFunc()
		}
	}

	checkPostgres := withRecovery(func() {
		defer wg.Done()
		for attempt := 1; attempt <= attempts; attempt++ {
			err := pgObj.GetPostgresStatus()
			if err == nil {
				break
			}
			loggerconfig.Error("Alert Severity:P3-Moderate, CheckConnection GetPostgresStatus Connection Failed error:", err, " attempt number:", attempt)
			if attempt == attempts {
				err = db.ReconnectToPostgres()
				if err != nil {
					loggerconfig.Error("Alert Severity:P0-Critical, CheckConnection ReconnectToPostgres error:", err)
				}
			}
		}
	})

	// mongoObj := db.GetMongoDBObj()

	// checkMongo := func() {
	// 	defer wg.Done()
	// 	for attempt := 1; attempt <= attempts; attempt++ {
	// 		err := mongoObj.GetMongoStatus()
	// 		if err == nil {
	// 			break
	// 		}
	// 		loggerconfig.Error("Alert Severity:P3-Moderate, CheckConnection GetMongoStatus Connection Failed error:", err, " attempt number:", attempt)
	// 		if attempt == attempts {
	// 			err = mongoObj.InitMongoClient(env)
	// 			if err != nil {
	// 				loggerconfig.Error("Alert Severity:P0-Critical, CheckConnection InitMongoClient error:", err)
	// 			}
	// 		}
	// 	}
	// }

	checkRedisMain := withRecovery(func() {
		defer wg.Done()
		for attempt := 1; attempt <= attempts; attempt++ {
			err := dbops.RedisRepo.GetStatus()
			if err == nil {
				break
			}
			loggerconfig.Error("Alert Severity:P2-Mid, CheckConnection GetStatusRedisMain Connection Failed error:", err, " attempt number:", attempt)

			if attempt == attempts {
				loggerconfig.Error("Alert Severity:P0-Critical, CheckConnection RedisInit failed after max attempts. Resetting Redis connection.")

				//Reset Redis connection
				factory := dbops.NewDatabaseConnectorFactory()
				redisUrl := loggerconfig.GetConfig().GetString("config.normal." + env + ".redisUrl")
				if len(redisUrl) == 0 || !strings.Contains(redisUrl, ":") {
					loggerconfig.Error("invalid redis URL format: must contain a non-empty string with a ':' separator")
					return
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
					loggerconfig.Error("Alert Severity:P0-Critical, Failed to reconnect Redis:", err)
					return
				}
				defer redisConnector.Close()

				//Update the RedisRepo pointer
				dbops.RedisRepo = redisConnector.GetRepository().(dbops.RedisRepository)
				loggerconfig.Info("Successfully reinitialized Redis connection.")
			}
		}
	})

	checkRedisOrder := withRecovery(func() {
		defer wg.Done()
		for attempt := 1; attempt <= attempts; attempt++ {
			err := redisCliObj.GetOrderClientStatus()
			if err == nil {
				break
			}
			loggerconfig.Error("Alert Severity:P3-Moderate, CheckConnection GetStatusRedisOrder Connection Failed error:", err, " attempt number:", attempt)
			if attempt == attempts {
				err = redisCliObj.InitRedis(constants.OrderRedis)
				if err != nil {
					loggerconfig.Error("Alert Severity:P0-Critical, CheckConnection OrderRedisinit error:", err)
				}
			}
		}
	})

	checkRedisContractCache := withRecovery(func() {
		defer wg.Done()
		for attempt := 1; attempt <= attempts; attempt++ {
			err := contractCacheClientObj.GetContractCacheStatus()
			if err == nil {
				break
			}
			loggerconfig.Error("Alert Severity:P3-Moderate, CheckConnection GetStatusRedisContractCache Connection Failed error:", err, " attempt number:", attempt)
			if attempt == attempts {
				err = contractCacheClientObj.ContractCacheInit()
				if err != nil {
					loggerconfig.Error("Alert Severity:P0-Critical, CheckConnection ContractCacheInit error:", err)
				}
			}
		}
	})

	checkRedisSmartCache := withRecovery(func() {
		defer wg.Done()
		for attempt := 1; attempt <= attempts; attempt++ {
			err := smartCacheClientObj.GetStatusRedisSmartCache()
			if err == nil {
				break
			}
			loggerconfig.Error("Alert Severity:P3-Moderate, CheckConnection GetStatusRedisSmartCache Connection Failed error:", err, " attempt number:", attempt)
			if attempt == attempts {
				err = smartCacheClientObj.InitSmartCache()
				if err != nil {
					loggerconfig.Error("Alert Severity:P0-Critical, CheckConnection InitSmartCache error:", err)
				}
			}
		}
	})

	checkRabbitMqConnection := withRecovery(func() {
		defer wg.Done()
		for attempt := 1; attempt <= attempts; attempt++ {
			if helpers.CheckConnection() {
				break
			}
			loggerconfig.Error("Alert Severity:P3-Moderate, CheckConnection CheckRabbitMqConnection Connection Failed, attempt number:", attempt)
			if attempt == attempts {
				err := helpers.InitializeRabbitMq()
				if err != nil {
					loggerconfig.Error("Alert Severity:P1-High, Scheduled health check failed to restore RabbitMQ connection: ", err)
				} else {
					loggerconfig.Info("Scheduled health check restored RabbitMQ connection")
				}
			}
		}
	})

	for {
		wg.Add(6) // Number of tasks to wait for // make it 6 after uncommented mongo retry

		go checkPostgres()
		// go checkMongo()
		go checkRedisMain()
		go checkRedisOrder()
		go checkRedisContractCache()
		go checkRedisSmartCache()
		go checkRabbitMqConnection()

		wg.Wait() // Wait for all Goroutines to finish

		time.Sleep(time.Duration(waitTime) * time.Second)
	}
}
