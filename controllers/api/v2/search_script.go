package v2

import (
	apihelpers "space/apiHelpers"
	"space/constants"
	"space/loggerconfig"
	"space/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var theContractDetailsProviderV2 models.ContractDetailsProviderV2

func InitContractDetailsProvider(provider models.ContractDetailsProviderV2) {
	defer models.HandlePanic()
	theContractDetailsProviderV2 = provider
}

// SearchScrip V2
// @Tags space contractdetails V2
// @Description Search Scrip Version 2
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
// @Router /api/space/v2/contractdetails/searchScrip [GET]
func SearchScript(ctx *gin.Context) {

	query := ctx.Query("searchText")
	pageStr := ctx.Query("page")
	exchange := ctx.Query("exchange")

	if query == "" {
		loggerconfig.Error("SearchScript (controller), v2 invalid query param query:", query)
		apihelpers.ErrorMessage(ctx, constants.InvalidRequest)
		return
	}

	if exchange != "" && !constants.ValidExchangeMap[exchange] {
		loggerconfig.Error("SearchScript (controller), v2 invalid exchange:", exchange)
		apihelpers.ErrorMessage(ctx, constants.InvalidExchange)
		return
	}

	// Convert the query parameters to integers
	pageInt, err := strconv.Atoi(pageStr)
	if err != nil || pageInt < 1 {
		pageInt = 1 // Default to page 1 if the provided value is invalid
	}

	cRH, _ := ctx.Get("reqH")
	requestH, _ := (cRH).(models.ReqHeader)

	start := (pageInt - 1) * constants.Capacity

	loggerconfig.Info(" SearchScript v2 (controller), for query:", query, " page:", pageStr, " and exchange:", exchange)

	//paginated api
	// first page = 1, return first 20 results
	// next page 2, return next 20 results
	if strings.Contains(query, "&") {
		query = strings.ReplaceAll(query, "&", "AND")
	}

	code, resp := theContractDetailsProviderV2.SearchScrip(query, requestH, start, constants.Capacity, exchange)

	logDetail := "clientId: " + requestH.ClientId + " function: logout V2 requestId: " + requestH.RequestId
	apihelpers.CustomResponse(ctx, code, resp, logDetail)
}
