package searchscriptv2

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/helpers"
	"space/helpers/cache"
	"space/loggerconfig"
	"space/models"
	"strconv"
	"strings"

	"unicode"
)

type SearchScriptV2 struct {
	contractCacheCli cache.ContractCache
	smartCacheCli    cache.SmartCache
}

func InitSearchScript(contractCacheCli cache.ContractCache, smartCacheCli cache.SmartCache) SearchScriptV2 {
	defer models.HandlePanic()
	obj := SearchScriptV2{
		contractCacheCli: contractCacheCli,
		smartCacheCli:    smartCacheCli,
	}
	return obj
}

func (obj SearchScriptV2) SearchScrip(querySearch string, reqH models.ReqHeader, offset int, capacity int, exchange string) (int, apihelpers.APIRes) {

	var apiRes apihelpers.APIRes

	querySearch = strings.ReplaceAll(querySearch, " ", "")

	containHyphen := strings.Contains(querySearch, "-")
	if containHyphen {
		symbolWithoutHypen := strings.ReplaceAll(querySearch, "-", "")
		querySearch = symbolWithoutHypen
	}

	if strings.Contains(querySearch, "&") {
		querySearch = strings.ReplaceAll(querySearch, "&", "AND")
	}

	if strings.Contains(querySearch, ".") {
		// replace . with Q since . is special character in full text search
		// querySearch = strings.ReplaceAll(querySearch, ".", "Q")
		querySearch = strings.ReplaceAll(querySearch, ".", "")
	}

	// in case special character exists return blank result
	if obj.containsSpecialCharacters(querySearch) {
		apiRes.Message = "SUCCESS"
		apiRes.Status = true

		return http.StatusOK, apiRes
	}

	//prefix search
	// searchTerm := querySearch + "*"

	uniqueKeys, err := obj.smartCacheCli.PerformNewSearch(exchange, querySearch, offset, capacity, false)
	if err != nil {
		loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " SearchScrip2: Error while prefix searching script substrig v2 for query ", querySearch, " and reqId ", reqH.RequestId, " and err ", querySearch, reqH.RequestId, err)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info(" capacity uk and capacity ", len(uniqueKeys), capacity)

	loggerconfig.Info(" uniqueKeysSubstring  and capacity ", len(uniqueKeys), capacity)

	if len(uniqueKeys) == 0 {
		uniqueKeysFuzzy, err := obj.smartCacheCli.PerformNewSearch(exchange, querySearch, offset, capacity, true)
		if err != nil {
			loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " SearchScrip2: Error while searching fuzzy script v2 for query  ", querySearch, " and reqId ", reqH.RequestId, " and err ", querySearch, reqH.RequestId, err)
			return apihelpers.SendInternalServerError()
		}
		uniqueKeys = append(uniqueKeys, uniqueKeysFuzzy...)
	}

	loggerconfig.Info(" uniqueKeysFuzzy  and capacity ", len(uniqueKeys), capacity)

	loggerconfig.Info("uniqueKeys", uniqueKeys)

	results := removeDuplicates(uniqueKeys)

	if len(results) > capacity {
		results = results[:capacity]
	}

	var stocks []models.SearchScripResponseResult
	// var stocks []models.ContractDetails
	for _, uniqueKey := range results {
		err, val := obj.contractCacheCli.GetFromHash("stock_key", uniqueKey)
		if err != nil {
			loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " SearchScrip2: Error while GetFromHashSetNew script v2 for query ", querySearch, " and reqId ", reqH.RequestId, " & uniqueKey-", uniqueKey, " and err ", querySearch, reqH.RequestId, uniqueKey, err)
			return apihelpers.SendInternalServerError()
		}

		stockDetail := models.ContractDetails{}
		if err = json.Unmarshal([]byte(val), &stockDetail); err != nil {
			loggerconfig.Info("Alert Severity:P1-High, platform:", reqH.Platform, " SearchScrip2: Error Unmarshalling stockmetadata for key  err ", err, " and reqId:", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		var searchScriptResp models.SearchScripResponseResult
		searchScriptResp.Token = stockDetail.Token1
		searchScriptResp.Exchange = stockDetail.Exchange
		searchScriptResp.Company = stockDetail.Name
		searchScriptResp.Symbol = stockDetail.Symbol
		searchScriptResp.Isin = stockDetail.Isin
		searchScriptResp.TradingSymbol = stockDetail.TradingSymbol
		searchScriptResp.DisplayName = buildDisplayName(stockDetail.Symbol, stockDetail.Exchange, stockDetail.TradingSymbol)
		searchScriptResp.Series = stockDetail.Series
		searchScriptResp.Segment = stockDetail.Exchange

		searchScriptResp.IsTradable = true

		if strings.EqualFold(searchScriptResp.Exchange, constants.NSE) || strings.EqualFold(searchScriptResp.Exchange, constants.BSE) {
			searchScriptResp.Segment = constants.SegmentEquity
		}

		if searchScriptResp.Series == constants.SegmentIndices {
			searchScriptResp.Segment = constants.SegmentIndex
			searchScriptResp.IsTradable = false
		}

		if strings.EqualFold(searchScriptResp.Exchange, constants.MCX) {
			searchScriptResp.Segment = constants.SegmentCommodity
		}

		searchScriptResp.Expiry = stockDetail.Expiry1
		searchScriptResp.Strike = stockDetail.Strike
		searchScriptResp.IsMtfEligible = stockDetail.IsMtfEligible

		if stockDetail.AlternateToken != "" {

			uniqueKey = "NSE_" + stockDetail.AlternateToken
			if searchScriptResp.Exchange == "NSE" {
				uniqueKey = "BSE_" + stockDetail.AlternateToken
			}

			err, val = obj.contractCacheCli.GetFromHash("stock_key", uniqueKey)
			if err != nil {
				loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " SearchScrip2: Error while GetFromHashSetNew script v2 for query ", querySearch, " and reqId ", reqH.RequestId, " & uniqueKey-", uniqueKey, " and err ", err)
				return apihelpers.SendInternalServerError()
			}

			stockDetailAlternate := models.ContractDetails{}
			if err = json.Unmarshal([]byte(val), &stockDetailAlternate); err != nil {
				loggerconfig.Info("Alert Severity:P1-High, platform:", reqH.Platform, " SearchScrip2: Error Unmarshalling stockmetadata for alternate token ", uniqueKey, "  err -", err, " and reqId-", reqH.RequestId)
				return apihelpers.SendInternalServerError()
			}
			searchScriptResp.Alternate.Token = stockDetailAlternate.Token1
			searchScriptResp.Alternate.Exchange = stockDetailAlternate.Exchange
			searchScriptResp.Alternate.Company = stockDetailAlternate.Name
			searchScriptResp.Alternate.Symbol = stockDetailAlternate.Symbol
			searchScriptResp.Alternate.TradingSymbol = stockDetailAlternate.TradingSymbol
			searchScriptResp.Alternate.DisplayName = stockDetailAlternate.TradingSymbol
			searchScriptResp.Alternate.IsTradable = true
			searchScriptResp.Alternate.Segment = stockDetailAlternate.Exchange
			searchScriptResp.Alternate.Expiry = stockDetailAlternate.Expiry1
			searchScriptResp.Alternate.Series = stockDetailAlternate.Series
			searchScriptResp.Alternate.Strike = stockDetailAlternate.Strike
			searchScriptResp.Alternate.IsMtfEligible = stockDetailAlternate.IsMtfEligible
		}

		stocks = append(stocks, searchScriptResp)
		//stocks = append(stocks, stockDetail)
	}

	// Return the matching stocks with metadata
	apiRes.Data = stocks
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj SearchScriptV2) PerformSearch(exchange, searchTerm string, offset int, capacity int) ([]interface{}, error) {

	var searchResults []interface{}
	var args []interface{}
	var err error

	// Perform the search
	if exchange == "" {
		//if exchange is blank override search query
		rediSearchQuery := fmt.Sprintf("%s", searchTerm)
		args = []interface{}{"FT.SEARCH", "stocks_idx", rediSearchQuery, "limit", offset, capacity, "SORTBY", "volume", "DESC"} // Command arguments
	} else {
		rediSearchQuery := fmt.Sprintf("@exchange:%s @symbol:%s", exchange, searchTerm)
		args = []interface{}{"FT.SEARCH", "stocks_idx", rediSearchQuery, "limit", offset, capacity, "SORTBY", "volume", "DESC"} // Command arguments
	}

	searchResults, err = obj.smartCacheCli.ExecFTCommand(args)
	if err != nil {
		loggerconfig.Error("SearchScrip2: Error while searching script v2 for query -%v and reqId-%v and err -%v\n", searchTerm, err)
		return searchResults, err
	}

	if len(searchResults) == 1 && exchange != "" {
		rediSearchQuery := fmt.Sprintf("@exchange:%s @tradingSymbol:%s", exchange, searchTerm)
		//"\"@exchange:%s @symbol:%s\"",
		args = []interface{}{"FT.SEARCH", "stocks_idx", rediSearchQuery, "limit", offset, capacity, "SORTBY", "volume", "DESC"} // Command arguments
		searchResults, err = obj.smartCacheCli.ExecFTCommand(args)
		if err != nil {
			loggerconfig.Error("SearchScrip2: Error while searching script v2 for query -", searchTerm, " and err ", searchTerm, err)
			return searchResults, err
		}

		if len(searchResults) == 1 {
			//search inside name
			rediSearchQuery := fmt.Sprintf("@exchange:%s @name:%s", exchange, searchTerm)
			//"\"@exchange:%s @symbol:%s\"",
			args = []interface{}{"FT.SEARCH", "stocks_idx", rediSearchQuery, "limit", offset, capacity, "SORTBY", "volume", "DESC"} // Command arguments
			searchResults, err = obj.smartCacheCli.ExecFTCommand(args)

			if err != nil {
				loggerconfig.Error("SearchScrip2: Error while searching script v2 for query ", searchTerm, " and err -", err)
				return searchResults, err
			}

		}
	}

	return searchResults, nil
}

func removeDuplicates(input []string) []string {
	uniqueMap := make(map[string]bool)
	result := []string{} // Stores the unique elements in order

	for _, item := range input {
		if _, exists := uniqueMap[item]; !exists {
			uniqueMap[item] = true
			result = append(result, item)
		}
	}

	return result
}

func buildDisplayName(symbol, exchange, tradingSymbol string) string {
	if strings.EqualFold(exchange, constants.NSE) || strings.EqualFold(exchange, constants.BSE) {
		return tradingSymbol
	}

	if tradingSymbol[len(tradingSymbol)-3:] == "FUT" {
		// log.Printf("fut %v", nfoFut)
		year := tradingSymbol[len(symbol) : len(symbol)+2]
		month := tradingSymbol[len(symbol)+2 : len(symbol)+5]
		//log.Printf("year %v", year)

		if exchange == strings.ToUpper(constants.MCX) {
			currYear := helpers.GetCurrentTimeInIST().Year()
			currYear = currYear % 100
			yearInt, _ := strconv.Atoi(year)
			if currYear < yearInt {
				return symbol + " " + year + month + " FUT"
			}
		}

		return symbol + " " + month + " FUT"
	}

	if tradingSymbol[len(tradingSymbol)-2:] == "CE" || tradingSymbol[len(tradingSymbol)-2:] == "PE" {
		// log.Printf("fut %v", nfoCE)
		year := tradingSymbol[len(symbol) : len(symbol)+2]
		month := tradingSymbol[len(symbol)+2 : len(symbol)+5]
		strike := tradingSymbol[len(symbol)+5 : len(tradingSymbol)-2]

		if containsNumeric(month) {
			modifiedMonth, err := parseDate(month)
			if err == nil {
				month = modifiedMonth
			}
		}

		return symbol + " " + month + " " + year + " " + strike + " " + tradingSymbol[len(tradingSymbol)-2:]

	}

	return tradingSymbol
}

// private methods
func (obj SearchScriptV2) containsSpecialCharacters(s string) bool {
	// Regular expression to match special characters, excluding '&' and '.'
	re := regexp.MustCompile(`[^a-zA-Z0-9.&]`)
	return re.MatchString(s)
}

func containsNumeric(str string) bool {
	for _, ch := range str {
		if unicode.IsDigit(ch) {
			return true
		}
	}
	return false
}

func getMonthName(monthID byte) (string, error) {
	switch monthID {
	case '1':
		return "JAN", nil
	case '2':
		return "FEB", nil
	case '3':
		return "MAR", nil
	case '4':
		return "APR", nil
	case '5':
		return "MAY", nil
	case '6':
		return "JUN", nil
	case '7':
		return "JUL", nil
	case '8':
		return "AUG", nil
	case '9':
		return "SEP", nil
	case 'O':
		return "OCT", nil
	case 'N':
		return "NOV", nil
	case 'D':
		return "DEC", nil
	}
	return "", errors.New("Invalid Month")

}

func parseDate(dateStr string) (string, error) {
	// Validate input length
	if len(dateStr) != 3 {
		return "", errors.New("Invalid Input")
	}

	// Get the month ID and the day from the string
	monthID := dateStr[0]
	dayStr := dateStr[1:]

	// Get the month name from the month ID
	monthName, err := getMonthName(monthID)

	if err != nil {
		return "", err
	}

	// Parse the day from the last two characters
	day, err := strconv.Atoi(dayStr)
	if err != nil || day < 1 || day > 31 {
		return "", errors.New("Invalid Day")
	}

	// Return the formatted date string
	return fmt.Sprintf("%d %s Weekly", day, monthName), nil
}
