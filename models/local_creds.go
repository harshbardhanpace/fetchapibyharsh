package models

type LocalCreds struct {
	Local struct {
		AdminSecretKey string `json:"adminSecretKey"`
		CmotsAuthToken string `json:"cmotsAuthToken"`
		Mongo          struct {
			MongoBase           string `json:"mongo_base"`
			MongoPassword       string `json:"mongo_password"`
			MongoURI            string `json:"mongo_uri"`
			MongoDao            string `json:"mongo_dao"`
			MongoSpace          string `json:"mongo_space"`
			MongoContractSearch string `json:"mongo_contract_search"`
		} `json:"mongo"`
		Postgres struct {
			Host        string `json:"host"`
			Port        int    `json:"port"`
			User        string `json:"user"`
			Password    string `json:"password"`
			Db          string `json:"db"`
			MaxConn     int    `json:"max_conn"`
			MaxIdleConn int    `json:"max_idle_conn"`
		} `json:"postgres"`
		SmartCache struct {
			Username                  string `json:"username"`
			Password                  string `json:"password"`
			Addr                      string `json:"addr"`
			PoolSize                  int    `json:"poolSize"`
			MindIdleConnection        int    `json:"mindIdleConnection"`
			PoolTimeoutInMilliSecond  int    `json:"poolTimeoutInMilliSecond"`
			ReadTimeoutInMilliSecond  int    `json:"readTimeoutInMilliSecond"`
			WriteTimeoutInMilliSecond int    `json:"writeTimeoutInMilliSecond"`
		} `json:"smartCache"`
		ContractCache struct {
			Addr                      string `json:"addr"`
			PoolSize                  int    `json:"poolSize"`
			Password                  string `json:"password"`
			MindIdleConnection        int    `json:"mindIdleConnection"`
			PoolTimeoutInMilliSecond  int    `json:"poolTimeoutInMilliSecond"`
			ReadTimeoutInMilliSecond  int    `json:"readTimeoutInMilliSecond"`
			WriteTimeoutInMilliSecond int    `json:"writeTimeoutInMilliSecond"`
		} `json:"contractCache"`
		RabbitMq struct {
			Addr     string `json:"addr"`
			User     string `json:"user"`
			Password string `json:"password"`
		} `json:"rabbitMq"`
		AwsS3CredConfig struct {
			MyRegion         string `json:"myRegion"`
			AccessKeyID      string `json:"AccessKeyId"`
			SecretAccessKey  string `json:"SecretAccessKey"`
			BucketNamePocket string `json:"BucketNamePocket"`
			BucketName       string `json:"BucketName"`
		} `json:"awsS3CredConfig"`
		FinvuPass       string `json:"FinvuPass"`
		FinvuRid        string `json:"FinvuRid"`
		FreshDeskApiKey string `json:"freshDeskApiKey"`
		FreshDeskPass   string `json:"freshDeskPass"`
		AuthKeyMsg91    string `json:"authKeyMsg91"`
	} `json:"local"`
}
