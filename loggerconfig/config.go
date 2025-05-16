package loggerconfig

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"space/constants"
	"space/models"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	viper *viper.Viper
}

var configInstance *Config // package private singleton instance of the configuration
var singleton sync.Once    // package private singleton helper utility

func GetConfig() *viper.Viper {
	// create an instance if not available
	singleton.Do(func() {
		configInstance = &Config{viper.New()}
	})

	return configInstance.viper
}

func Start() {
	// Find and read the config file
	err := GetConfig().ReadInConfig()
	fmt.Printf("err =%v\n", err)
	if err != nil {
		// Handle errors reading the config file
		log.Printf(" Error Reading config = %v\n", err)
	}
}

var LocalCreds models.LocalCreds

func FetchLocalCreds() {
	file, err := os.Open("resources/localcreds.json")
	if err != nil {
		fmt.Println("FetchLocalCreds Error Reading local creds file:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&LocalCreds); err != nil {
		fmt.Println("FetchLocalCreds Error decoding local creds file:", err)
		return
	}
}

func SetEnv(env string) {
	constants.Env = env
}
