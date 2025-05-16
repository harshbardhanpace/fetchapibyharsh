package v3

import (
	"fmt"
	"strconv"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
)

var theScripProvider models.ScripProvider

// InitScripsProvider initializes the search service with MongoDB collection
func InitScripsProvider(scripProvider models.ScripProvider) {
	defer models.HandlePanic()
	theScripProvider = scripProvider
}

// SearchScrip V3
// @Tags space contractdetails V3
// @Description Search Scrip Version 3 with MongoDB Atlas Search
// @Param searchText query string false "searchText"
// @Param page query string false "page"
// @Param exchange query string false "exchange"
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param P-Platform header string false "P-Platform Header"
// @Param P-DeviceId header string false "P-DeviceId Header"
// @Param P-ClientType header string false "P-ClientType Header"
// @Param Authorization header string true "Authorization Header"
// @Param P-ClientPublicIP header string false "P-ClientPublicIP Header"
// @Success 200 {object} apihelpers.APIRes{data=models.SearchScripResponse}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 403 {object} apihelpers.APIRes
// @Router /api/space/v3/contractdetails/searchScrip [GET]
func SearchScrip(ctx *gin.Context) {
	var reqH models.ReqHeader
	if err := ctx.ShouldBindHeader(&reqH); err != nil {
		loggerconfig.Error("SearchScrip (controller), Error in parsing header, error = ", err, " requestId:", reqH.RequestId, " platform:", reqH.Platform, " deviceId:", reqH.DeviceId)
		return
	}

	query := ctx.Query("searchText")
	pageStr := ctx.Query("page")
	exchange := ctx.Query("exchange")

	if query == "" {
		loggerconfig.Error("SearchScrip (controller), v3 invalid query param query:", query, " requestId:", reqH.RequestId, " platform:", reqH.Platform, " deviceId:", reqH.DeviceId)
		apihelpers.ErrorMessage(ctx, constants.InvalidRequest)
		return
	}

	if exchange != "" && !constants.ValidExchangeMap[exchange] {
		loggerconfig.Error("SearchScrip (controller), v3 invalid exchange:", exchange, " requestId:", reqH.RequestId, " platform:", reqH.Platform, " deviceId:", reqH.DeviceId)
		apihelpers.ErrorMessage(ctx, constants.InvalidExchange)
		return
	}

	// Convert page to integer
	pageInt, err := strconv.Atoi(pageStr)
	if err != nil || pageInt < 1 {
		pageInt = 1 // Default to page 1 if the provided value is invalid
	}

	loggerconfig.Info("SearchScrip v3 (controller), for query:", query, " page:", pageStr, " and exchange:", exchange, " requestId:", reqH.RequestId, " platform:", reqH.Platform, " deviceId:", reqH.DeviceId)

	// Create request object
	searchRequest := models.SearchScripAPIRequest{
		SearchText: query,
		Exchange:   exchange,
		Page:       pageInt,
	}

	// Call business logic
	statusCode, apiRes := theScripProvider.SearchScrip(searchRequest, reqH)

	// Log the response
	logDetail := fmt.Sprintf("Search query: %s, page: %d, exchange: %s, requestId: %s, platform: %s, deviceId: %s, statusCode: %d",
		query, pageInt, exchange, reqH.RequestId, reqH.Platform, reqH.DeviceId, statusCode)
	loggerconfig.Info("SearchScrip v3 (controller) response, ", logDetail)

	// Send response
	apihelpers.CustomResponse(ctx, statusCode, apiRes, logDetail)
}
