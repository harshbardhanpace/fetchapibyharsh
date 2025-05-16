package scrips

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/db"
	"space/dbops"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ScripObj struct {
	mongoDb db.MongoDatabase
}

func InitScrips(mongoDb db.MongoDatabase) ScripObj {
	defer models.HandlePanic()
	ScripObj := ScripObj{mongoDb: mongoDb}
	return ScripObj
}

func (obj ScripObj) SearchScrip(req models.SearchScripAPIRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	loggerconfig.Info("SearchScrip (service), for searchText:", req.SearchText, " page:", req.Page, " and exchange:", req.Exchange, " requestId:", reqH.RequestId, " platform:", reqH.Platform, " deviceId:", reqH.DeviceId)

	var apiRes apihelpers.APIRes

	// Set fixed limit per page
	limit := constants.ContractSearchPageLimit

	// Calculate offset based on page
	page := req.Page
	if page < 1 {
		page = 1 // Default to page 1 if invalid
	}
	offset := (page - 1) * limit

	// Create the MongoDB pipeline
	pipeline := mongo.Pipeline{}

	// Add search stage
	searchStage := bson.D{
		{Key: "$search", Value: bson.D{
			{Key: "index", Value: "instruments_search"},
			{Key: "text", Value: bson.D{
				{Key: "query", Value: req.SearchText},
				{Key: "path", Value: bson.D{{Key: "wildcard", Value: "*"}}},
			}},
		}},
	}
	pipeline = append(pipeline, searchStage)

	// Add exchange filter if provided
	if req.Exchange != "" {
		matchStage := bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "Exchange", Value: req.Exchange},
			}},
		}
		pipeline = append(pipeline, matchStage)
	}

	// Add sorting
	sortStage := bson.D{
		{Key: "$sort", Value: bson.D{
			{Key: "InternalInstrumentIdentifier", Value: -1},
		}},
	}
	pipeline = append(pipeline, sortStage)

	// Add pagination
	skipStage := bson.D{{Key: "$skip", Value: offset}}
	limitStage := bson.D{{Key: "$limit", Value: limit}}
	pipeline = append(pipeline, skipStage, limitStage)

	// Get the appropriate collection
	collection := dbops.MongoContractSearchRepo.GetMongoCollection(constants.INSTRUMENTS_COLLECTION)

	// Execute the query
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, SearchScrip (service), aggregate error:", err, " requestId:", reqH.RequestId, " platform:", reqH.Platform, " deviceId:", reqH.DeviceId)
		return apihelpers.SendInternalServerError()
	}
	defer cursor.Close(context.Background())

	// Parse results
	var results []models.ContractDetails
	if err := cursor.All(context.Background(), &results); err != nil {
		loggerconfig.Error("Alert Severity:P0-Critical, SearchScrip (service), cursor.All error:", err, " requestId:", reqH.RequestId, " platform:", reqH.Platform, " deviceId:", reqH.DeviceId)
		return apihelpers.SendInternalServerError()
	}

	// Convert ContractDetails to SearchScripResponseResult
	var searchResults []models.SearchScripResponseResult
	for _, stockDetail := range results {
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
			// Create a filter for the alternate token
			alternateFilter := bson.D{
				{Key: "Token1", Value: stockDetail.AlternateToken},
				{Key: "Exchange", Value: bson.D{
					{Key: "$ne", Value: stockDetail.Exchange},
				}},
			}

			// Find the alternate token details
			var alternateDetail models.ContractDetails
			err := collection.FindOne(context.Background(), alternateFilter).Decode(&alternateDetail)
			if err != nil && err != mongo.ErrNoDocuments {
				loggerconfig.Error("Alert Severity:P1-High, platform:", reqH.Platform, " SearchScrip: Error while finding alternate token for query ", req.SearchText, " and reqId ", reqH.RequestId, " & alternateToken-", stockDetail.AlternateToken, " and err ", err)
				// Though it should be logged, it should not be returned as error
			}

			if err == nil {
				searchScriptResp.Alternate.Token = alternateDetail.Token1
				searchScriptResp.Alternate.Exchange = alternateDetail.Exchange
				searchScriptResp.Alternate.Company = alternateDetail.Name
				searchScriptResp.Alternate.Symbol = alternateDetail.Symbol
				searchScriptResp.Alternate.TradingSymbol = alternateDetail.TradingSymbol
				searchScriptResp.Alternate.DisplayName = alternateDetail.TradingSymbol
				searchScriptResp.Alternate.IsTradable = true
				searchScriptResp.Alternate.Segment = alternateDetail.Exchange
				searchScriptResp.Alternate.Expiry = alternateDetail.Expiry1
				searchScriptResp.Alternate.Series = alternateDetail.Series
				searchScriptResp.Alternate.Strike = alternateDetail.Strike
				searchScriptResp.Alternate.IsMtfEligible = alternateDetail.IsMtfEligible
			}
		}

		searchResults = append(searchResults, searchScriptResp)
	}

	// Log success
	logDetail := fmt.Sprintf("Search query: %s, exchange: %s, page: %d, found %d results, requestId: %s, platform: %s, deviceId: %s",
		req.SearchText, req.Exchange, page, len(searchResults), reqH.RequestId, reqH.Platform, reqH.DeviceId)
	loggerconfig.Info("SearchScrip (service) successful, ", logDetail)

	// Return success response
	apiRes.Data = searchResults
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

// Helper function to build display name
func buildDisplayName(symbol, exchange, tradingSymbol string) string {
	if strings.EqualFold(exchange, constants.NSE) || strings.EqualFold(exchange, constants.BSE) {
		return tradingSymbol
	}

	if tradingSymbol[len(tradingSymbol)-3:] == "FUT" {
		year := tradingSymbol[len(symbol) : len(symbol)+2]
		month := tradingSymbol[len(symbol)+2 : len(symbol)+5]

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

// Helper function to check if string contains numeric characters
func containsNumeric(str string) bool {
	for _, ch := range str {
		if unicode.IsDigit(ch) {
			return true
		}
	}
	return false
}

// Helper function to get month name
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
	return "", errors.New("invalid month")
}

// Helper function to parse date
func parseDate(dateStr string) (string, error) {
	if len(dateStr) != 3 {
		return "", errors.New("invalid input")
	}

	monthID := dateStr[0]
	dayStr := dateStr[1:]

	monthName, err := getMonthName(monthID)
	if err != nil {
		return "", err
	}

	day, err := strconv.Atoi(dayStr)
	if err != nil || day < 1 || day > 31 {
		return "", errors.New("invalid day")
	}

	return fmt.Sprintf("%d %s Weekly", day, monthName), nil
}
