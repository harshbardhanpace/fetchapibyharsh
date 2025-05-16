package cache

import (
	"errors"
	"fmt"
	"space/constants"
	"strconv"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

type SmartCacheRedisClient struct {
	pool      *redis.Pool
	searchObj *redisearch.Client
}

var smartCacheClientObj SmartCache

func newPool(server, username, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			// Authenticate with the provided username and password
			if _, err := c.Do("AUTH", username, password); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func (c *SmartCacheRedisClient) InitSmartCache() error {

	c.pool = newPool(constants.SmartCacheAddr, constants.SmartCacheUsername, constants.SmartCachePassword)
	if c.pool == nil {
		return errors.New("Error connecting to smart cache")
	}
	c.searchObj = redisearch.NewClientFromPool(c.pool, "stockIndex")
	return nil
}

func (c *SmartCacheRedisClient) GetStatusRedisSmartCache() error {
	conn := c.pool.Get()
	_, err := conn.Do("PING")
	if err != nil {
		logrus.Error("GetStatusRedisSmartCache Failed to execute PING command:", err)
		return err
	}
	return nil
}

func (c *SmartCacheRedisClient) PerformNewSearch(exchange, searchTerm string, offset, capacity int, fuzzy bool) ([]string, error) {
	keys := make([]string, 0)

	if exchange != "" {

		val := "@exchange:" + exchange + "(@tradingSymbol:" + searchTerm + "*)|(@symbol:" + searchTerm + "*)|(@name:" + searchTerm + "*)|(@isin:" + searchTerm + "*)"

		logrus.Infof("query1  -%v\n", val)
		query := redisearch.NewQuery(val).
			SetSortBy("instIdentifier", false).
			Limit(offset, capacity)

		docsTradingSymbol, total, err := c.searchObj.Search(query)

		if err != nil {
			logrus.Infof("error searching via tradingSymbol -%v\n", err)
			return keys, err
		}

		if total < capacity {
			val := "@exchange:" + exchange + "(@tradingSymbol:*" + searchTerm + "*)|(@symbol:*" + searchTerm + "*)|(@name:*" + searchTerm + "*)|(@isin:*" + searchTerm + "*)"

			logrus.Infof("query2  -%v\n", val)
			query := redisearch.NewQuery(val).
				SetSortBy("instIdentifier", false).
				Limit(offset, capacity-total)

			docsNew, _, err := c.searchObj.Search(query)

			if err != nil {
				logrus.Infof("error searching via tradingSymbol -%v\n", err)
				return keys, err
			}

			docsTradingSymbol = append(docsTradingSymbol, docsNew...)

		}

		for _, doc := range docsTradingSymbol {

			// vol := doc.Properties["volume"].(string)
			// if _, ok := mapOfKeys[doc.Properties["stockKey"].(string)]; !ok {
			// 	mapOfKeys[doc.Properties["stockKey"].(string)] = atoi(vol)
			// }

			keys = append(keys, doc.Properties["stockKey"].(string))
		}

	} else {
		var docs []redisearch.Document
		var err error
		if fuzzy {

			val := "@exchange:NSE (@tradingSymbol:%%" + searchTerm + "%%)|(@symbol:%%" + searchTerm + "%%)|(@name:%%" + searchTerm + "%%)"

			logrus.Infof("query1  -%v\n", val)
			query := redisearch.NewQuery(val).
				SetSortBy("instIdentifier", false).
				Limit(offset, capacity)

			docsNSE, _, err := c.searchObj.Search(query)

			if err != nil {
				logrus.Infof("error searching via nse fuzzy -%v\n", err)
				return keys, err
			}

			for _, doc := range docsNSE {

				// vol := doc.Properties["volume"].(string)
				// if _, ok := mapOfKeys[doc.Properties["stockKey"].(string)]; !ok {
				// 	mapOfKeys[doc.Properties["stockKey"].(string)] = atoi(vol)
				// }

				keys = append(keys, doc.Properties["stockKey"].(string))
			}

			valBSE := "@exchange:BSE (@tradingSymbol:%%" + searchTerm + "%%)|(@symbol:%%" + searchTerm + "%%)|(@name:%%" + searchTerm + "%%)"

			logrus.Infof("query1  -%v\n", valBSE)
			queryBSE := redisearch.NewQuery(valBSE).
				SetSortBy("instIdentifier", false).
				Limit(offset, capacity)

			docsBSE, _, err := c.searchObj.Search(queryBSE)

			if err != nil {
				logrus.Infof("error searching via bse fuzzy -%v\n", err)
				return keys, err
			}

			for _, doc := range docsBSE {

				// vol := doc.Properties["volume"].(string)
				// if _, ok := mapOfKeys[doc.Properties["stockKey"].(string)]; !ok {
				// 	mapOfKeys[doc.Properties["stockKey"].(string)] = atoi(vol)
				// }

				keys = append(keys, doc.Properties["stockKey"].(string))
			}

		} else {
			val := "(@tradingSymbol:" + searchTerm + "*)|(@symbol:" + searchTerm + "*)|(@name:" + searchTerm + "*)|(@isin:" + searchTerm + "*)"
			if isNumber(searchTerm) {
				val = "(@strike:" + searchTerm + "*)"
				//(@symbol:*nifty* & *21700*)
				//if any number comes up search in NFO banknifty and nifty
				//	toMatchNifty := "NIFTY* &*" + searchTerm
				//	toMatchBankNifty := "BANKNIFTY* &*" + searchTerm
				//	val = "@exchange:NFO (@tradingSymbol:" + toMatchNifty + "*)|(@symbol:" + toMatchNifty + "*)|(@name:" + toMatchNifty + "*)|(@isin:" + searchTerm + "*)|(@tradingSymbol:" + toMatchNifty + "*)|(@symbol:" + toMatchBankNifty + "*)|(@name:" + toMatchBankNifty + "*)|(@isin:" + toMatchBankNifty + "*)"
				logrus.Infof("queryNormal  -%v\n", val)
				query := redisearch.NewQuery(val).
					SetSortBy("instIdentifier", false).
					Limit(offset, capacity)
				// if strings.ToUpper(searchTerm) == "NIFTY" || strings.ToUpper(searchTerm) == "BANKNIFTY" {
				// 	val = "(@tradingSymbol:" + searchTerm + "*)|(@symbol:" + searchTerm + "*)|(@name:" + searchTerm + "*)"
				// 	query = redisearch.NewQuery(val).
				// 		SetSortBy("volume", false).
				// 		Limit(offset, capacity)
				// }

				docs, _, err = c.searchObj.Search(query)

				if err != nil {
					logrus.Errorf("error querying searchTerm %v\n", err)
					return keys, err
				}
				for _, doc := range docs {
					fmt.Printf("Name: %s, Symbol: %s, Trading Symbol: %s, Strike: %s, key: %s\n",
						doc.Properties["name"], doc.Properties["symbol"],
						doc.Properties["tradingSymbol"], doc.Properties["strike"], doc.Properties["stockKey"])
					keys = append(keys, doc.Properties["stockKey"].(string))
				}
				return keys, nil
			}
			logrus.Infof("queryNormal  -%v\n", val)
			query := redisearch.NewQuery(val).
				SetSortBy("instIdentifier", false).
				Limit(offset, capacity)
			// if strings.ToUpper(searchTerm) == "NIFTY" || strings.ToUpper(searchTerm) == "BANKNIFTY" {
			// 	val = "(@tradingSymbol:" + searchTerm + "*)|(@symbol:" + searchTerm + "*)|(@name:" + searchTerm + "*)"
			// 	query = redisearch.NewQuery(val).
			// 		SetSortBy("volume", false).
			// 		Limit(offset, capacity)
			// }

			var total int
			docs, total, err = c.searchObj.Search(query)

			if err != nil {
				logrus.Errorf("error querying searchTerm %v\n", err)
				return keys, err
			}

			if total < capacity {
				val = "(@tradingSymbol:*" + searchTerm + "*)|(@symbol:*" + searchTerm + "*)|(@name:*" + searchTerm + "*)|(@isin:*" + searchTerm + "*)"
				if isNumber(searchTerm) {
					val = "(@strike:" + searchTerm + "*)"
				}

				query := redisearch.NewQuery(val).
					SetSortBy("instIdentifier", false).
					Limit(offset, capacity-total)
				docsNew, _, err := c.searchObj.Search(query)

				if err != nil {
					logrus.Errorf("error querying searchTerm %v\n", err)
					return keys, err
				}

				docs = append(docs, docsNew...)
			}
		}

		for _, doc := range docs {
			fmt.Printf("Name: %s, Symbol: %s, Trading Symbol: %s, Strike: %s, key: %s\n",
				doc.Properties["name"], doc.Properties["symbol"],
				doc.Properties["tradingSymbol"], doc.Properties["strike"], doc.Properties["stockKey"])
			keys = append(keys, doc.Properties["stockKey"].(string))
		}
	}

	return keys, nil
}

func atoi(val string) int {
	res, _ := strconv.Atoi(val)
	return res
}

func (c *SmartCacheRedisClient) ExecFTCommand(args []interface{}) ([]interface{}, error) {
	fmt.Println(args)
	conn := c.pool.Get()
	defer conn.Close()
	val, err := redis.Values(conn.Do(args[0].(string), args[1:]...))
	return val, err
}

func ParseResultValues(results []interface{}) (map[string]interface{}, error) {
	if len(results) <= 1 {
		return nil, errors.New("no result found")
	}

	document := make(map[string]interface{})

	// Extract the document ID
	if id, ok := results[0].(int64); ok {
		document["id"] = id
	}

	// Extract the name field (assuming it's a string)
	if name, ok := results[1].([]byte); ok {
		document["key"] = string(name)
	}

	// Extract the fields and their values
	fieldsData, ok := results[2].([]interface{})
	if ok {
		for i := 0; i < len(fieldsData); i += 2 {
			fieldKey, ok := fieldsData[i].([]byte)
			if !ok {
				logrus.Error("Error converting field key to []byte")
				continue
			}

			fieldValue, ok := fieldsData[i+1].([]byte)
			if !ok {
				continue
			}

			document[string(fieldKey)] = string(fieldValue)
		}
	}

	return document, nil
}

func (c *SmartCacheRedisClient) GetFromHashSetNew(hash, key string) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()
	val, err := redis.String(conn.Do("HGET", hash, key))
	if err != nil {
		return "", err
	}
	return val, nil
}

// GetSortedSetData extracts unique keys from the results of a sorted set query.
func GetSortedSetData(results []interface{}) []string {
	uniqueKeys := make([]string, 0)
	var getValueFlag bool

	// Iterate through the results
	for i := 0; i < len(results); i++ {
		switch val := results[i].(type) {
		case []interface{}:
			// Extract data from nested arrays
			for _, nestedVal := range val {
				switch nestedArray := nestedVal.(type) {
				case []byte:
					key := string(nestedArray)
					if getValueFlag {
						uniqueKeys = append(uniqueKeys, string(nestedArray))
						getValueFlag = false
					}
					if key == constants.UniqueKey {
						getValueFlag = true
					}
				case []interface{}:
					properties := make(map[string]string)
					for j := 0; j < len(nestedArray); j++ {
						field := nestedArray[j].([]interface{})
						property := string(field[0].([]byte))
						value := string(field[1].([]byte))
						properties[property] = value
					}
				}
			}
		}
	}

	return uniqueKeys
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func SetSmartCacheClienttObj(smartCacheCliObj SmartCache) {
	smartCacheClientObj = smartCacheCliObj
}

func GetSmartCacheClientObj() SmartCache {
	return smartCacheClientObj
}
