package health

import (
	"space/constants"
	"space/db"
	"space/helpers/cache"
)

func GetHealthInfo() map[string]string {

	errorMap := make(map[string]string)

	redisObj := cache.GetRedisClientObj()
	// check redis health
	if err := redisObj.GetClientStatus(); err != nil {
		errorMap[constants.REDIS] = err.Error()
	}

	if err := redisObj.GetOrderClientStatus(); err != nil {
		errorMap[constants.REDIS] = err.Error()
	}

	//check mongo health
	obj := db.GetMongoDBObj()
	if err := obj.GetMongoStatus(); err != nil {
		errorMap[constants.MONGO] = err.Error()
	}

	//check postgres health
	if err := db.GetPgObj().GetPostgresStatus(); err != nil {
		errorMap[constants.POSTGRES] = err.Error()
	}

	if len(errorMap) > 0 {
		return errorMap
	}

	return nil
}
