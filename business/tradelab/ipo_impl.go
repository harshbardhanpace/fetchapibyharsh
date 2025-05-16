package tradelab

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"time"

	apihelpers "space/apiHelpers"
	"space/constants"
	"space/db"
	"space/dbops"
	"space/helpers"
	"space/helpers/cache"
	"space/loggerconfig"
	"space/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IpoObj struct {
	tradeLabURL string
	mongodb     db.MongoDatabase
	redisCli    cache.RedisCache
}

func InitIpoProvider(mongodb db.MongoDatabase, redisCli cache.RedisCache) IpoObj {
	defer models.HandlePanic()

	ipoObj := IpoObj{
		tradeLabURL: constants.TLURL,
		mongodb:     mongodb,
		redisCli:    redisCli,
	}

	return ipoObj
}

func (obj IpoObj) GetAllIpo(req models.GetAllIpoRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {

	var apiRes apihelpers.APIRes

	url := obj.tradeLabURL + GETIPOURL

	//make payload
	payload := new(bytes.Buffer)

	//call api
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallApi(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "GetAllIpo", duration, reqH.ClientId, reqH.RequestId)
	if err != nil {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " getAllIpoReq call api error =", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("getAllIpoRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlGetAllIpoRes := TradeLabGetAllIpoResponse{}
	json.Unmarshal([]byte(string(body)), &tlGetAllIpoRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " GetAllIpo tl status not ok =", tlGetAllIpoRes.Message, "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlGetAllIpoRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}
	redisCli := cache.GetRedisClientObj()
	if len(tlGetAllIpoRes.Data.AllIpo.Data) == 0 { //Data null
		tlGetAllIpoRes, err = getGetAllIPODataRedis(FetchAllIPODataRedisKey, obj)
		if err != nil {
			json.Unmarshal([]byte(string(body)), &tlGetAllIpoRes) //undo any change if it might have happened due to some error
			loggerconfig.Error("Alert Severity:TLP2-Mid, platform:", reqH.Platform, "  GetAllIpo fetch from cache failed, err:", err, " requestId:", reqH.RequestId)
		}
	}
	// fill up controller response
	var getIpoResponse models.GetAllIpoResponse

	allIpo := make([]models.IpoState, 0)
	openIpo := make([]models.IpoState, 0)
	upcomingIpo := make([]models.IpoState, 0)
	closeIpo := make([]models.IpoState, 0)
	var ipoWG sync.WaitGroup
	ipoWG.Add(4)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := 0; i < len(tlGetAllIpoRes.Data.AllIpo.Data); i++ {
			var allIpoRes models.IpoState
			bidEndDate := helpers.FindDateIndex(tlGetAllIpoRes.Data.AllIpo.Data[i].BiddingStartDate)
			var bidEndDateActual string
			if bidEndDate != -1 {
				for countWorkingDays := 0; countWorkingDays < 2; { //T+2 only counting working days
					if !helpers.HolidayCalendar.Date[bidEndDate].IsHoliday { //if not holiday, x and bidEndDate both are incremented, else x is not
						countWorkingDays++
					}
					bidEndDate++
				}
				bidEndDateActual = helpers.HolidayCalendar.Date[bidEndDate].Date
				originalFormat := "02-Jan-2006" //not putting in constants as it would be confusing to understand what format is being converted to which
				desiredFormat := "02-01-2006"
				bidEndDateActual, _ = helpers.ConvertFormat(bidEndDateActual, originalFormat, desiredFormat)
			}

			allIpoRes.BiddingStartDate = tlGetAllIpoRes.Data.AllIpo.Data[i].BiddingStartDate
			allIpoRes.Symbol = tlGetAllIpoRes.Data.AllIpo.Data[i].Symbol
			allIpoRes.MinBidQuantity = tlGetAllIpoRes.Data.AllIpo.Data[i].MinBidQuantity
			allIpoRes.Registrar = tlGetAllIpoRes.Data.AllIpo.Data[i].Registrar
			allIpoRes.LotSize = tlGetAllIpoRes.Data.AllIpo.Data[i].LotSize
			allIpoRes.T1ModEndDate = tlGetAllIpoRes.Data.AllIpo.Data[i].T1ModEndDate
			allIpoRes.DailyStartTime = tlGetAllIpoRes.Data.AllIpo.Data[i].DailyStartTime
			allIpoRes.T1ModStartTime = tlGetAllIpoRes.Data.AllIpo.Data[i].T1ModStartTime
			allIpoRes.BiddingEndDate = tlGetAllIpoRes.Data.AllIpo.Data[i].BiddingEndDate
			if bidEndDateActual != "" { //if there was en error then null would be return and hence original will not be changed
				allIpoRes.BiddingEndDate = bidEndDateActual
			}

			listingDate := helpers.FindDateIndex(allIpoRes.BiddingEndDate)
			var listingDateActual string
			if listingDate != -1 {
				for countWorkingDays := 0; countWorkingDays < 3; { //T+3 only counting working days
					if !helpers.HolidayCalendar.Date[listingDate].IsHoliday { //if not holiday, x and bidEndDate both are incremented, else x is not
						countWorkingDays++
					}
					listingDate++
				}
				listingDateActual = helpers.HolidayCalendar.Date[listingDate].Date
				originalFormat := "02-Jan-2006" //not putting in constants as it would be confusing to understand what format is being converted to which
				desiredFormat := "02-01-2006"
				listingDateActual, _ = helpers.ConvertFormat(listingDateActual, originalFormat, desiredFormat)
			}
			allIpoRes.T1ModEndTime = tlGetAllIpoRes.Data.AllIpo.Data[i].T1ModEndTime
			allIpoRes.DailyEndTime = tlGetAllIpoRes.Data.AllIpo.Data[i].DailyEndTime
			allIpoRes.TickSize = tlGetAllIpoRes.Data.AllIpo.Data[i].TickSize
			allIpoRes.IssueType = tlGetAllIpoRes.Data.AllIpo.Data[i].IssueType
			allIpoRes.FaceValue = tlGetAllIpoRes.Data.AllIpo.Data[i].FaceValue
			allIpoRes.MinPrice = tlGetAllIpoRes.Data.AllIpo.Data[i].MinPrice
			allIpoRes.T1ModStartDate = tlGetAllIpoRes.Data.AllIpo.Data[i].T1ModStartDate
			allIpoRes.Name = tlGetAllIpoRes.Data.AllIpo.Data[i].Name
			allIpoRes.IssueSize = tlGetAllIpoRes.Data.AllIpo.Data[i].IssueSize
			allIpoRes.MaxPrice = tlGetAllIpoRes.Data.AllIpo.Data[i].MaxPrice
			allIpoRes.CutOffPrice = tlGetAllIpoRes.Data.AllIpo.Data[i].CutOffPrice
			allIpoRes.UnixBiddingEndDate = tlGetAllIpoRes.Data.AllIpo.Data[i].UnixBiddingEndDate
			allIpoRes.UnixBiddingStartDate = tlGetAllIpoRes.Data.AllIpo.Data[i].UnixBiddingStartDate
			allIpoRes.Isin = tlGetAllIpoRes.Data.AllIpo.Data[i].Isin
			allIpoRes.AllotmentDate = tlGetAllIpoRes.Data.AllIpo.Data[i].AllotmentDate
			allIpoRes.ExchangeIssueType = tlGetAllIpoRes.Data.AllIpo.Data[i].ExchangeIssueType
			allIpoRes.AllotmentBegins = tlGetAllIpoRes.Data.AllIpo.Data[i].AllotmentBegins
			allIpoRes.RefundDate = tlGetAllIpoRes.Data.AllIpo.Data[i].RefundDate
			allIpoRes.ListingDate = tlGetAllIpoRes.Data.AllIpo.Data[i].ListingDate
			if listingDateActual != "" { //if there was en error then null would be return and hence original will not be changed
				allIpoRes.ListingDate = listingDateActual
			}
			allIpoRes.AboutCompany = tlGetAllIpoRes.Data.AllIpo.Data[i].AboutCompany
			allIpoRes.ParentCompany = tlGetAllIpoRes.Data.AllIpo.Data[i].ParentCompany
			allIpoRes.FoundedYear = tlGetAllIpoRes.Data.AllIpo.Data[i].FoundedYear
			allIpoRes.ProspectusFileURL = tlGetAllIpoRes.Data.AllIpo.Data[i].ProspectusFileURL
			allIpoRes.ManagingDirector = tlGetAllIpoRes.Data.AllIpo.Data[i].ManagingDirector
			allIpoRes.MaxLimit = tlGetAllIpoRes.Data.AllIpo.Data[i].MaxLimit
			allIpoRes.RetailDiscount = tlGetAllIpoRes.Data.AllIpo.Data[i].RetailDiscount
			allIpoRes.NseExchangeListed = tlGetAllIpoRes.Data.AllIpo.Data[i].NseExchangeListed
			allIpoRes.BseExchangeListed = tlGetAllIpoRes.Data.AllIpo.Data[i].BseExchangeListed
			allIpoRes.AmoOrderEntryTime = tlGetAllIpoRes.Data.AllIpo.Data[i].AmoOrderEntryTime
			allIpoRes.ApplicationRangeStart = tlGetAllIpoRes.Data.AllIpo.Data[i].ApplicationRangeStart
			allIpoRes.ApplicationRangeEnd = tlGetAllIpoRes.Data.AllIpo.Data[i].ApplicationRangeEnd
			allIpoRes.TotalApplicationRangeCount = tlGetAllIpoRes.Data.AllIpo.Data[i].TotalApplicationRangeCount
			allIpoRes.CategoryDetails = tlGetAllIpoRes.Data.AllIpo.Data[i].CategoryDetails

			var allIpoStateSubCategorySettings []models.IpoStateSubCategorySettings
			for j := 0; j < len(tlGetAllIpoRes.Data.AllIpo.Data[i].SubCategorySettings); j++ {
				var ipoStateSubCategorySettings models.IpoStateSubCategorySettings
				ipoStateSubCategorySettings.SubCatCode = tlGetAllIpoRes.Data.AllIpo.Data[i].SubCategorySettings[j].SubCatCode
				ipoStateSubCategorySettings.MinValue = tlGetAllIpoRes.Data.AllIpo.Data[i].SubCategorySettings[j].MinValue
				ipoStateSubCategorySettings.MaxUpiLimit = tlGetAllIpoRes.Data.AllIpo.Data[i].SubCategorySettings[j].MaxUpiLimit
				ipoStateSubCategorySettings.AllowCutOff = tlGetAllIpoRes.Data.AllIpo.Data[i].SubCategorySettings[j].AllowCutOff
				ipoStateSubCategorySettings.AllowUpi = tlGetAllIpoRes.Data.AllIpo.Data[i].SubCategorySettings[j].AllowUpi
				ipoStateSubCategorySettings.MaxValue = tlGetAllIpoRes.Data.AllIpo.Data[i].SubCategorySettings[j].MaxValue
				ipoStateSubCategorySettings.DiscountPrice = tlGetAllIpoRes.Data.AllIpo.Data[i].SubCategorySettings[j].DiscountPrice
				ipoStateSubCategorySettings.DiscountType = tlGetAllIpoRes.Data.AllIpo.Data[i].SubCategorySettings[j].DiscountType
				ipoStateSubCategorySettings.MaxPrice = tlGetAllIpoRes.Data.AllIpo.Data[i].SubCategorySettings[j].MaxPrice
				ipoStateSubCategorySettings.CaCode = tlGetAllIpoRes.Data.AllIpo.Data[i].SubCategorySettings[j].CaCode
				ipoStateSubCategorySettings.Allowed = tlGetAllIpoRes.Data.AllIpo.Data[i].SubCategorySettings[j].Allowed
				ipoStateSubCategorySettings.StartDate = tlGetAllIpoRes.Data.AllIpo.Data[i].SubCategorySettings[j].StartDate
				ipoStateSubCategorySettings.EndDate = tlGetAllIpoRes.Data.AllIpo.Data[i].SubCategorySettings[j].EndDate
				ipoStateSubCategorySettings.DisplayName = tlGetAllIpoRes.Data.AllIpo.Data[i].SubCategorySettings[j].DisplayName
				ipoStateSubCategorySettings.MinLotSize = tlGetAllIpoRes.Data.AllIpo.Data[i].SubCategorySettings[j].MinLotSize
				ipoStateSubCategorySettings.StartTime = tlGetAllIpoRes.Data.AllIpo.Data[i].SubCategorySettings[j].StartTime
				ipoStateSubCategorySettings.EndTime = tlGetAllIpoRes.Data.AllIpo.Data[i].SubCategorySettings[j].EndTime
				allIpoStateSubCategorySettings = append(allIpoStateSubCategorySettings, ipoStateSubCategorySettings)
			}

			allIpoRes.SubCategorySettings = allIpoStateSubCategorySettings

			allIpoRes.IpoAllowed = tlGetAllIpoRes.Data.AllIpo.Data[i].IpoAllowed
			allIpoRes.BseAllowed = tlGetAllIpoRes.Data.AllIpo.Data[i].BseAllowed
			allIpoRes.NseAllowed = tlGetAllIpoRes.Data.AllIpo.Data[i].NseAllowed
			allIpoRes.SubType = tlGetAllIpoRes.Data.AllIpo.Data[i].SubType
			allIpoRes.EnablePio = tlGetAllIpoRes.Data.AllIpo.Data[i].EnablePio
			allIpoRes.PioStartDate = tlGetAllIpoRes.Data.AllIpo.Data[i].PioStartDate
			allIpoRes.PioEndDate = tlGetAllIpoRes.Data.AllIpo.Data[i].PioEndDate
			allIpoRes.PioEndTime = tlGetAllIpoRes.Data.AllIpo.Data[i].PioEndTime
			allIpoRes.PioStartTime = tlGetAllIpoRes.Data.AllIpo.Data[i].PioStartTime
			allIpoRes.DematTransferDate = tlGetAllIpoRes.Data.AllIpo.Data[i].DematTransferDate
			allIpoRes.MandateEndDate = tlGetAllIpoRes.Data.AllIpo.Data[i].MandateEndDate
			allIpoRes.IsEmployeeCat = tlGetAllIpoRes.Data.AllIpo.Data[i].IsEmployeeCat
			allIpoRes.IsShareHolderCat = tlGetAllIpoRes.Data.AllIpo.Data[i].IsShareHolderCat
			allIpo = append(allIpo, allIpoRes)
		}
		getIpoResponse.AllIpo.Data = allIpo
		getIpoResponse.AllIpo.Status = tlGetAllIpoRes.Data.AllIpo.Status
	}(&ipoWG)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := 0; i < len(tlGetAllIpoRes.Data.OpenIpo); i++ {
			var openIpoRes models.IpoState

			openIpoRes.BiddingStartDate = tlGetAllIpoRes.Data.OpenIpo[i].BiddingStartDate
			openIpoRes.Symbol = tlGetAllIpoRes.Data.OpenIpo[i].Symbol
			openIpoRes.MinBidQuantity = tlGetAllIpoRes.Data.OpenIpo[i].MinBidQuantity
			openIpoRes.Registrar = tlGetAllIpoRes.Data.OpenIpo[i].Registrar
			openIpoRes.LotSize = tlGetAllIpoRes.Data.OpenIpo[i].LotSize
			openIpoRes.T1ModEndDate = tlGetAllIpoRes.Data.OpenIpo[i].T1ModEndDate
			openIpoRes.DailyStartTime = tlGetAllIpoRes.Data.OpenIpo[i].DailyStartTime
			openIpoRes.T1ModStartTime = tlGetAllIpoRes.Data.OpenIpo[i].T1ModStartTime
			openIpoRes.BiddingEndDate = tlGetAllIpoRes.Data.OpenIpo[i].BiddingEndDate
			openIpoRes.T1ModEndTime = tlGetAllIpoRes.Data.OpenIpo[i].T1ModEndTime
			openIpoRes.DailyEndTime = tlGetAllIpoRes.Data.OpenIpo[i].DailyEndTime
			openIpoRes.TickSize = tlGetAllIpoRes.Data.OpenIpo[i].TickSize
			openIpoRes.IssueType = tlGetAllIpoRes.Data.OpenIpo[i].IssueType
			openIpoRes.FaceValue = tlGetAllIpoRes.Data.OpenIpo[i].FaceValue
			openIpoRes.MinPrice = tlGetAllIpoRes.Data.OpenIpo[i].MinPrice
			openIpoRes.T1ModStartDate = tlGetAllIpoRes.Data.OpenIpo[i].T1ModStartDate
			openIpoRes.Name = tlGetAllIpoRes.Data.OpenIpo[i].Name
			openIpoRes.IssueSize = tlGetAllIpoRes.Data.OpenIpo[i].IssueSize
			issueSizeFromRedis := redisCli.GetRedis(constants.IPOPrefix + openIpoRes.Symbol)
			decryptIssueSizeFromRedis, _ := issueSizeFromRedis.Result()
			if decryptIssueSizeFromRedis != "" && openIpoRes.IssueSize == 0 {
				integerIssueSizeFromRedis, err := strconv.Atoi(decryptIssueSizeFromRedis)
				if err == nil {
					openIpoRes.IssueSize = integerIssueSizeFromRedis
				}
			}
			openIpoRes.MaxPrice = tlGetAllIpoRes.Data.OpenIpo[i].MaxPrice
			openIpoRes.CutOffPrice = tlGetAllIpoRes.Data.OpenIpo[i].CutOffPrice
			openIpoRes.UnixBiddingEndDate = tlGetAllIpoRes.Data.OpenIpo[i].UnixBiddingEndDate
			openIpoRes.UnixBiddingStartDate = tlGetAllIpoRes.Data.OpenIpo[i].UnixBiddingStartDate
			openIpoRes.Isin = tlGetAllIpoRes.Data.OpenIpo[i].Isin

			openIpoRes.AllotmentDate = tlGetAllIpoRes.Data.OpenIpo[i].AllotmentDate
			openIpoRes.ExchangeIssueType = tlGetAllIpoRes.Data.OpenIpo[i].ExchangeIssueType
			openIpoRes.AllotmentBegins = tlGetAllIpoRes.Data.OpenIpo[i].AllotmentBegins
			openIpoRes.RefundDate = tlGetAllIpoRes.Data.OpenIpo[i].RefundDate
			openIpoRes.ListingDate = tlGetAllIpoRes.Data.OpenIpo[i].ListingDate
			openIpoRes.AboutCompany = tlGetAllIpoRes.Data.OpenIpo[i].AboutCompany
			openIpoRes.ParentCompany = tlGetAllIpoRes.Data.OpenIpo[i].ParentCompany
			openIpoRes.FoundedYear = tlGetAllIpoRes.Data.OpenIpo[i].FoundedYear
			openIpoRes.ProspectusFileURL = tlGetAllIpoRes.Data.OpenIpo[i].ProspectusFileURL
			openIpoRes.ManagingDirector = tlGetAllIpoRes.Data.OpenIpo[i].ManagingDirector
			openIpoRes.MaxLimit = tlGetAllIpoRes.Data.OpenIpo[i].MaxLimit
			openIpoRes.RetailDiscount = tlGetAllIpoRes.Data.OpenIpo[i].RetailDiscount
			openIpoRes.NseExchangeListed = tlGetAllIpoRes.Data.OpenIpo[i].NseExchangeListed
			openIpoRes.BseExchangeListed = tlGetAllIpoRes.Data.OpenIpo[i].BseExchangeListed
			openIpoRes.AmoOrderEntryTime = tlGetAllIpoRes.Data.OpenIpo[i].AmoOrderEntryTime
			openIpoRes.ApplicationRangeStart = tlGetAllIpoRes.Data.OpenIpo[i].ApplicationRangeStart
			openIpoRes.ApplicationRangeEnd = tlGetAllIpoRes.Data.OpenIpo[i].ApplicationRangeEnd
			openIpoRes.TotalApplicationRangeCount = tlGetAllIpoRes.Data.OpenIpo[i].TotalApplicationRangeCount
			openIpoRes.CategoryDetails = tlGetAllIpoRes.Data.OpenIpo[i].CategoryDetails

			var allIpoStateSubCategorySettings []models.IpoStateSubCategorySettings
			for j := 0; j < len(tlGetAllIpoRes.Data.OpenIpo[i].SubCategorySettings); j++ {
				var ipoStateSubCategorySettings models.IpoStateSubCategorySettings
				ipoStateSubCategorySettings.SubCatCode = tlGetAllIpoRes.Data.OpenIpo[i].SubCategorySettings[j].SubCatCode
				ipoStateSubCategorySettings.MinValue = tlGetAllIpoRes.Data.OpenIpo[i].SubCategorySettings[j].MinValue
				ipoStateSubCategorySettings.MaxUpiLimit = tlGetAllIpoRes.Data.OpenIpo[i].SubCategorySettings[j].MaxUpiLimit
				ipoStateSubCategorySettings.AllowCutOff = tlGetAllIpoRes.Data.OpenIpo[i].SubCategorySettings[j].AllowCutOff
				ipoStateSubCategorySettings.AllowUpi = tlGetAllIpoRes.Data.OpenIpo[i].SubCategorySettings[j].AllowUpi
				ipoStateSubCategorySettings.MaxValue = tlGetAllIpoRes.Data.OpenIpo[i].SubCategorySettings[j].MaxValue
				ipoStateSubCategorySettings.DiscountPrice = tlGetAllIpoRes.Data.OpenIpo[i].SubCategorySettings[j].DiscountPrice
				ipoStateSubCategorySettings.DiscountType = tlGetAllIpoRes.Data.OpenIpo[i].SubCategorySettings[j].DiscountType
				ipoStateSubCategorySettings.MaxPrice = tlGetAllIpoRes.Data.OpenIpo[i].SubCategorySettings[j].MaxPrice
				ipoStateSubCategorySettings.CaCode = tlGetAllIpoRes.Data.OpenIpo[i].SubCategorySettings[j].CaCode
				ipoStateSubCategorySettings.Allowed = tlGetAllIpoRes.Data.OpenIpo[i].SubCategorySettings[j].Allowed
				ipoStateSubCategorySettings.StartDate = tlGetAllIpoRes.Data.OpenIpo[i].SubCategorySettings[j].StartDate
				ipoStateSubCategorySettings.EndDate = tlGetAllIpoRes.Data.OpenIpo[i].SubCategorySettings[j].EndDate
				ipoStateSubCategorySettings.DisplayName = tlGetAllIpoRes.Data.OpenIpo[i].SubCategorySettings[j].DisplayName
				ipoStateSubCategorySettings.MinLotSize = tlGetAllIpoRes.Data.OpenIpo[i].SubCategorySettings[j].MinLotSize
				ipoStateSubCategorySettings.StartTime = tlGetAllIpoRes.Data.OpenIpo[i].SubCategorySettings[j].StartTime
				ipoStateSubCategorySettings.EndTime = tlGetAllIpoRes.Data.OpenIpo[i].SubCategorySettings[j].EndTime
				allIpoStateSubCategorySettings = append(allIpoStateSubCategorySettings, ipoStateSubCategorySettings)
			}

			openIpoRes.SubCategorySettings = allIpoStateSubCategorySettings

			openIpoRes.IpoAllowed = tlGetAllIpoRes.Data.OpenIpo[i].IpoAllowed
			openIpoRes.BseAllowed = tlGetAllIpoRes.Data.OpenIpo[i].BseAllowed
			openIpoRes.NseAllowed = tlGetAllIpoRes.Data.OpenIpo[i].NseAllowed
			openIpoRes.SubType = tlGetAllIpoRes.Data.OpenIpo[i].SubType
			openIpoRes.EnablePio = tlGetAllIpoRes.Data.OpenIpo[i].EnablePio
			openIpoRes.PioStartDate = tlGetAllIpoRes.Data.OpenIpo[i].PioStartDate
			openIpoRes.PioEndDate = tlGetAllIpoRes.Data.OpenIpo[i].PioEndDate
			openIpoRes.PioEndTime = tlGetAllIpoRes.Data.OpenIpo[i].PioEndTime
			openIpoRes.PioStartTime = tlGetAllIpoRes.Data.OpenIpo[i].PioStartTime
			openIpoRes.DematTransferDate = tlGetAllIpoRes.Data.OpenIpo[i].DematTransferDate
			openIpoRes.MandateEndDate = tlGetAllIpoRes.Data.OpenIpo[i].MandateEndDate
			openIpoRes.IsEmployeeCat = tlGetAllIpoRes.Data.OpenIpo[i].IsEmployeeCat
			openIpoRes.IsShareHolderCat = tlGetAllIpoRes.Data.OpenIpo[i].IsShareHolderCat
			openIpo = append(openIpo, openIpoRes)
		}
		getIpoResponse.OpenIpo = openIpo
	}(&ipoWG)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := 0; i < len(tlGetAllIpoRes.Data.UpcomingIpo); i++ {
			var upcomingIpoRes models.IpoState

			upcomingIpoRes.BiddingStartDate = tlGetAllIpoRes.Data.UpcomingIpo[i].BiddingStartDate
			upcomingIpoRes.Symbol = tlGetAllIpoRes.Data.UpcomingIpo[i].Symbol
			upcomingIpoRes.MinBidQuantity = tlGetAllIpoRes.Data.UpcomingIpo[i].MinBidQuantity
			upcomingIpoRes.Registrar = tlGetAllIpoRes.Data.UpcomingIpo[i].Registrar
			upcomingIpoRes.LotSize = tlGetAllIpoRes.Data.UpcomingIpo[i].LotSize
			upcomingIpoRes.T1ModEndDate = tlGetAllIpoRes.Data.UpcomingIpo[i].T1ModEndDate
			upcomingIpoRes.DailyStartTime = tlGetAllIpoRes.Data.UpcomingIpo[i].DailyStartTime
			upcomingIpoRes.T1ModStartTime = tlGetAllIpoRes.Data.UpcomingIpo[i].T1ModStartTime
			upcomingIpoRes.BiddingEndDate = tlGetAllIpoRes.Data.UpcomingIpo[i].BiddingEndDate
			upcomingIpoRes.T1ModEndTime = tlGetAllIpoRes.Data.UpcomingIpo[i].T1ModEndTime
			upcomingIpoRes.DailyEndTime = tlGetAllIpoRes.Data.UpcomingIpo[i].DailyEndTime
			upcomingIpoRes.TickSize = tlGetAllIpoRes.Data.UpcomingIpo[i].TickSize
			upcomingIpoRes.IssueType = tlGetAllIpoRes.Data.UpcomingIpo[i].IssueType
			upcomingIpoRes.FaceValue = tlGetAllIpoRes.Data.UpcomingIpo[i].FaceValue
			upcomingIpoRes.MinPrice = tlGetAllIpoRes.Data.UpcomingIpo[i].MinPrice
			upcomingIpoRes.T1ModStartDate = tlGetAllIpoRes.Data.UpcomingIpo[i].T1ModStartDate
			upcomingIpoRes.Name = tlGetAllIpoRes.Data.UpcomingIpo[i].Name
			upcomingIpoRes.IssueSize = tlGetAllIpoRes.Data.UpcomingIpo[i].IssueSize
			issueSizeStr := strconv.Itoa(upcomingIpoRes.IssueSize)
			redisCli.SetRedis(constants.IPOPrefix+upcomingIpoRes.Symbol, issueSizeStr, constants.TwoMonthsInMins)
			upcomingIpoRes.MaxPrice = tlGetAllIpoRes.Data.UpcomingIpo[i].MaxPrice
			upcomingIpoRes.CutOffPrice = tlGetAllIpoRes.Data.UpcomingIpo[i].CutOffPrice
			upcomingIpoRes.UnixBiddingEndDate = tlGetAllIpoRes.Data.UpcomingIpo[i].UnixBiddingEndDate
			upcomingIpoRes.UnixBiddingStartDate = tlGetAllIpoRes.Data.UpcomingIpo[i].UnixBiddingStartDate
			upcomingIpoRes.Isin = tlGetAllIpoRes.Data.UpcomingIpo[i].Isin

			upcomingIpoRes.AllotmentDate = tlGetAllIpoRes.Data.UpcomingIpo[i].AllotmentDate
			upcomingIpoRes.ExchangeIssueType = tlGetAllIpoRes.Data.UpcomingIpo[i].ExchangeIssueType
			upcomingIpoRes.AllotmentBegins = tlGetAllIpoRes.Data.UpcomingIpo[i].AllotmentBegins
			upcomingIpoRes.RefundDate = tlGetAllIpoRes.Data.UpcomingIpo[i].RefundDate
			upcomingIpoRes.ListingDate = tlGetAllIpoRes.Data.UpcomingIpo[i].ListingDate
			upcomingIpoRes.AboutCompany = tlGetAllIpoRes.Data.UpcomingIpo[i].AboutCompany
			upcomingIpoRes.ParentCompany = tlGetAllIpoRes.Data.UpcomingIpo[i].ParentCompany
			upcomingIpoRes.FoundedYear = tlGetAllIpoRes.Data.UpcomingIpo[i].FoundedYear
			upcomingIpoRes.ProspectusFileURL = tlGetAllIpoRes.Data.UpcomingIpo[i].ProspectusFileURL
			upcomingIpoRes.ManagingDirector = tlGetAllIpoRes.Data.UpcomingIpo[i].ManagingDirector
			upcomingIpoRes.MaxLimit = tlGetAllIpoRes.Data.UpcomingIpo[i].MaxLimit
			upcomingIpoRes.RetailDiscount = tlGetAllIpoRes.Data.UpcomingIpo[i].RetailDiscount
			upcomingIpoRes.NseExchangeListed = tlGetAllIpoRes.Data.UpcomingIpo[i].NseExchangeListed
			upcomingIpoRes.BseExchangeListed = tlGetAllIpoRes.Data.UpcomingIpo[i].BseExchangeListed
			upcomingIpoRes.AmoOrderEntryTime = tlGetAllIpoRes.Data.UpcomingIpo[i].AmoOrderEntryTime
			upcomingIpoRes.ApplicationRangeStart = tlGetAllIpoRes.Data.UpcomingIpo[i].ApplicationRangeStart
			upcomingIpoRes.ApplicationRangeEnd = tlGetAllIpoRes.Data.UpcomingIpo[i].ApplicationRangeEnd
			upcomingIpoRes.TotalApplicationRangeCount = tlGetAllIpoRes.Data.UpcomingIpo[i].TotalApplicationRangeCount
			upcomingIpoRes.CategoryDetails = tlGetAllIpoRes.Data.UpcomingIpo[i].CategoryDetails

			var allIpoStateSubCategorySettings []models.IpoStateSubCategorySettings
			for j := 0; j < len(tlGetAllIpoRes.Data.UpcomingIpo[i].SubCategorySettings); j++ {
				var ipoStateSubCategorySettings models.IpoStateSubCategorySettings
				ipoStateSubCategorySettings.SubCatCode = tlGetAllIpoRes.Data.UpcomingIpo[i].SubCategorySettings[j].SubCatCode
				ipoStateSubCategorySettings.MinValue = tlGetAllIpoRes.Data.UpcomingIpo[i].SubCategorySettings[j].MinValue
				ipoStateSubCategorySettings.MaxUpiLimit = tlGetAllIpoRes.Data.UpcomingIpo[i].SubCategorySettings[j].MaxUpiLimit
				ipoStateSubCategorySettings.AllowCutOff = tlGetAllIpoRes.Data.UpcomingIpo[i].SubCategorySettings[j].AllowCutOff
				ipoStateSubCategorySettings.AllowUpi = tlGetAllIpoRes.Data.UpcomingIpo[i].SubCategorySettings[j].AllowUpi
				ipoStateSubCategorySettings.MaxValue = tlGetAllIpoRes.Data.UpcomingIpo[i].SubCategorySettings[j].MaxValue
				ipoStateSubCategorySettings.DiscountPrice = tlGetAllIpoRes.Data.UpcomingIpo[i].SubCategorySettings[j].DiscountPrice
				ipoStateSubCategorySettings.DiscountType = tlGetAllIpoRes.Data.UpcomingIpo[i].SubCategorySettings[j].DiscountType
				ipoStateSubCategorySettings.MaxPrice = tlGetAllIpoRes.Data.UpcomingIpo[i].SubCategorySettings[j].MaxPrice
				ipoStateSubCategorySettings.CaCode = tlGetAllIpoRes.Data.UpcomingIpo[i].SubCategorySettings[j].CaCode
				ipoStateSubCategorySettings.Allowed = tlGetAllIpoRes.Data.UpcomingIpo[i].SubCategorySettings[j].Allowed
				ipoStateSubCategorySettings.StartDate = tlGetAllIpoRes.Data.UpcomingIpo[i].SubCategorySettings[j].StartDate
				ipoStateSubCategorySettings.EndDate = tlGetAllIpoRes.Data.UpcomingIpo[i].SubCategorySettings[j].EndDate
				ipoStateSubCategorySettings.DisplayName = tlGetAllIpoRes.Data.UpcomingIpo[i].SubCategorySettings[j].DisplayName
				ipoStateSubCategorySettings.MinLotSize = tlGetAllIpoRes.Data.UpcomingIpo[i].SubCategorySettings[j].MinLotSize
				ipoStateSubCategorySettings.StartTime = tlGetAllIpoRes.Data.UpcomingIpo[i].SubCategorySettings[j].StartTime
				ipoStateSubCategorySettings.EndTime = tlGetAllIpoRes.Data.UpcomingIpo[i].SubCategorySettings[j].EndTime
				allIpoStateSubCategorySettings = append(allIpoStateSubCategorySettings, ipoStateSubCategorySettings)
			}

			upcomingIpoRes.SubCategorySettings = allIpoStateSubCategorySettings

			upcomingIpoRes.IpoAllowed = tlGetAllIpoRes.Data.UpcomingIpo[i].IpoAllowed
			upcomingIpoRes.BseAllowed = tlGetAllIpoRes.Data.UpcomingIpo[i].BseAllowed
			upcomingIpoRes.NseAllowed = tlGetAllIpoRes.Data.UpcomingIpo[i].NseAllowed
			upcomingIpoRes.SubType = tlGetAllIpoRes.Data.UpcomingIpo[i].SubType
			upcomingIpoRes.EnablePio = tlGetAllIpoRes.Data.UpcomingIpo[i].EnablePio
			upcomingIpoRes.PioStartDate = tlGetAllIpoRes.Data.UpcomingIpo[i].PioStartDate
			upcomingIpoRes.PioEndDate = tlGetAllIpoRes.Data.UpcomingIpo[i].PioEndDate
			upcomingIpoRes.PioEndTime = tlGetAllIpoRes.Data.UpcomingIpo[i].PioEndTime
			upcomingIpoRes.PioStartTime = tlGetAllIpoRes.Data.UpcomingIpo[i].PioStartTime
			upcomingIpoRes.DematTransferDate = tlGetAllIpoRes.Data.UpcomingIpo[i].DematTransferDate
			upcomingIpoRes.MandateEndDate = tlGetAllIpoRes.Data.UpcomingIpo[i].MandateEndDate
			upcomingIpoRes.IsEmployeeCat = tlGetAllIpoRes.Data.UpcomingIpo[i].IsEmployeeCat
			upcomingIpoRes.IsShareHolderCat = tlGetAllIpoRes.Data.UpcomingIpo[i].IsShareHolderCat
			upcomingIpo = append(upcomingIpo, upcomingIpoRes)
		}
		getIpoResponse.UpcomingIpo = upcomingIpo
	}(&ipoWG)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := 0; i < len(tlGetAllIpoRes.Data.ClosedIpo); i++ {
			var closeIpoRes models.IpoState

			closeIpoRes.BiddingStartDate = tlGetAllIpoRes.Data.ClosedIpo[i].BiddingStartDate
			closeIpoRes.Symbol = tlGetAllIpoRes.Data.ClosedIpo[i].Symbol
			closeIpoRes.MinBidQuantity = tlGetAllIpoRes.Data.ClosedIpo[i].MinBidQuantity
			closeIpoRes.Registrar = tlGetAllIpoRes.Data.ClosedIpo[i].Registrar
			closeIpoRes.LotSize = tlGetAllIpoRes.Data.ClosedIpo[i].LotSize
			closeIpoRes.T1ModEndDate = tlGetAllIpoRes.Data.ClosedIpo[i].T1ModEndDate
			closeIpoRes.DailyStartTime = tlGetAllIpoRes.Data.ClosedIpo[i].DailyStartTime
			closeIpoRes.T1ModStartTime = tlGetAllIpoRes.Data.ClosedIpo[i].T1ModStartTime
			closeIpoRes.BiddingEndDate = tlGetAllIpoRes.Data.ClosedIpo[i].BiddingEndDate
			closeIpoRes.T1ModEndTime = tlGetAllIpoRes.Data.ClosedIpo[i].T1ModEndTime
			closeIpoRes.DailyEndTime = tlGetAllIpoRes.Data.ClosedIpo[i].DailyEndTime
			closeIpoRes.TickSize = tlGetAllIpoRes.Data.ClosedIpo[i].TickSize
			closeIpoRes.IssueType = tlGetAllIpoRes.Data.ClosedIpo[i].IssueType
			closeIpoRes.FaceValue = tlGetAllIpoRes.Data.ClosedIpo[i].FaceValue
			closeIpoRes.MinPrice = tlGetAllIpoRes.Data.ClosedIpo[i].MinPrice
			closeIpoRes.T1ModStartDate = tlGetAllIpoRes.Data.ClosedIpo[i].T1ModStartDate
			closeIpoRes.Name = tlGetAllIpoRes.Data.ClosedIpo[i].Name
			closeIpoRes.IssueSize = tlGetAllIpoRes.Data.ClosedIpo[i].IssueSize
			closeIpoRes.MaxPrice = tlGetAllIpoRes.Data.ClosedIpo[i].MaxPrice
			closeIpoRes.CutOffPrice = tlGetAllIpoRes.Data.ClosedIpo[i].CutOffPrice
			closeIpoRes.UnixBiddingEndDate = tlGetAllIpoRes.Data.ClosedIpo[i].UnixBiddingEndDate
			closeIpoRes.UnixBiddingStartDate = tlGetAllIpoRes.Data.ClosedIpo[i].UnixBiddingStartDate
			closeIpoRes.Isin = tlGetAllIpoRes.Data.ClosedIpo[i].Isin

			closeIpoRes.AllotmentDate = tlGetAllIpoRes.Data.ClosedIpo[i].AllotmentDate
			closeIpoRes.ExchangeIssueType = tlGetAllIpoRes.Data.ClosedIpo[i].ExchangeIssueType
			closeIpoRes.AllotmentBegins = tlGetAllIpoRes.Data.ClosedIpo[i].AllotmentBegins
			closeIpoRes.RefundDate = tlGetAllIpoRes.Data.ClosedIpo[i].RefundDate
			closeIpoRes.ListingDate = tlGetAllIpoRes.Data.ClosedIpo[i].ListingDate
			closeIpoRes.AboutCompany = tlGetAllIpoRes.Data.ClosedIpo[i].AboutCompany
			closeIpoRes.ParentCompany = tlGetAllIpoRes.Data.ClosedIpo[i].ParentCompany
			closeIpoRes.FoundedYear = tlGetAllIpoRes.Data.ClosedIpo[i].FoundedYear
			closeIpoRes.ProspectusFileURL = tlGetAllIpoRes.Data.ClosedIpo[i].ProspectusFileURL
			closeIpoRes.ManagingDirector = tlGetAllIpoRes.Data.ClosedIpo[i].ManagingDirector
			closeIpoRes.MaxLimit = tlGetAllIpoRes.Data.ClosedIpo[i].MaxLimit
			closeIpoRes.RetailDiscount = tlGetAllIpoRes.Data.ClosedIpo[i].RetailDiscount
			closeIpoRes.NseExchangeListed = tlGetAllIpoRes.Data.ClosedIpo[i].NseExchangeListed
			closeIpoRes.BseExchangeListed = tlGetAllIpoRes.Data.ClosedIpo[i].BseExchangeListed
			closeIpoRes.AmoOrderEntryTime = tlGetAllIpoRes.Data.ClosedIpo[i].AmoOrderEntryTime
			closeIpoRes.ApplicationRangeStart = tlGetAllIpoRes.Data.ClosedIpo[i].ApplicationRangeStart
			closeIpoRes.ApplicationRangeEnd = tlGetAllIpoRes.Data.ClosedIpo[i].ApplicationRangeEnd
			closeIpoRes.TotalApplicationRangeCount = tlGetAllIpoRes.Data.ClosedIpo[i].TotalApplicationRangeCount
			closeIpoRes.CategoryDetails = tlGetAllIpoRes.Data.ClosedIpo[i].CategoryDetails

			var allIpoStateSubCategorySettings []models.IpoStateSubCategorySettings
			for j := 0; j < len(tlGetAllIpoRes.Data.ClosedIpo[i].SubCategorySettings); j++ {
				var ipoStateSubCategorySettings models.IpoStateSubCategorySettings
				ipoStateSubCategorySettings.SubCatCode = tlGetAllIpoRes.Data.ClosedIpo[i].SubCategorySettings[j].SubCatCode
				ipoStateSubCategorySettings.MinValue = tlGetAllIpoRes.Data.ClosedIpo[i].SubCategorySettings[j].MinValue
				ipoStateSubCategorySettings.MaxUpiLimit = tlGetAllIpoRes.Data.ClosedIpo[i].SubCategorySettings[j].MaxUpiLimit
				ipoStateSubCategorySettings.AllowCutOff = tlGetAllIpoRes.Data.ClosedIpo[i].SubCategorySettings[j].AllowCutOff
				ipoStateSubCategorySettings.AllowUpi = tlGetAllIpoRes.Data.ClosedIpo[i].SubCategorySettings[j].AllowUpi
				ipoStateSubCategorySettings.MaxValue = tlGetAllIpoRes.Data.ClosedIpo[i].SubCategorySettings[j].MaxValue
				ipoStateSubCategorySettings.DiscountPrice = tlGetAllIpoRes.Data.ClosedIpo[i].SubCategorySettings[j].DiscountPrice
				ipoStateSubCategorySettings.DiscountType = tlGetAllIpoRes.Data.ClosedIpo[i].SubCategorySettings[j].DiscountType
				ipoStateSubCategorySettings.MaxPrice = tlGetAllIpoRes.Data.ClosedIpo[i].SubCategorySettings[j].MaxPrice
				ipoStateSubCategorySettings.CaCode = tlGetAllIpoRes.Data.ClosedIpo[i].SubCategorySettings[j].CaCode
				ipoStateSubCategorySettings.Allowed = tlGetAllIpoRes.Data.ClosedIpo[i].SubCategorySettings[j].Allowed
				ipoStateSubCategorySettings.StartDate = tlGetAllIpoRes.Data.ClosedIpo[i].SubCategorySettings[j].StartDate
				ipoStateSubCategorySettings.EndDate = tlGetAllIpoRes.Data.ClosedIpo[i].SubCategorySettings[j].EndDate
				ipoStateSubCategorySettings.DisplayName = tlGetAllIpoRes.Data.ClosedIpo[i].SubCategorySettings[j].DisplayName
				ipoStateSubCategorySettings.MinLotSize = tlGetAllIpoRes.Data.ClosedIpo[i].SubCategorySettings[j].MinLotSize
				ipoStateSubCategorySettings.StartTime = tlGetAllIpoRes.Data.ClosedIpo[i].SubCategorySettings[j].StartTime
				ipoStateSubCategorySettings.EndTime = tlGetAllIpoRes.Data.ClosedIpo[i].SubCategorySettings[j].EndTime
				allIpoStateSubCategorySettings = append(allIpoStateSubCategorySettings, ipoStateSubCategorySettings)
			}

			closeIpoRes.SubCategorySettings = allIpoStateSubCategorySettings

			closeIpoRes.IpoAllowed = tlGetAllIpoRes.Data.ClosedIpo[i].IpoAllowed
			closeIpoRes.BseAllowed = tlGetAllIpoRes.Data.ClosedIpo[i].BseAllowed
			closeIpoRes.NseAllowed = tlGetAllIpoRes.Data.ClosedIpo[i].NseAllowed
			closeIpoRes.SubType = tlGetAllIpoRes.Data.ClosedIpo[i].SubType
			closeIpoRes.EnablePio = tlGetAllIpoRes.Data.ClosedIpo[i].EnablePio
			closeIpoRes.PioStartDate = tlGetAllIpoRes.Data.ClosedIpo[i].PioStartDate
			closeIpoRes.PioEndDate = tlGetAllIpoRes.Data.ClosedIpo[i].PioEndDate
			closeIpoRes.PioEndTime = tlGetAllIpoRes.Data.ClosedIpo[i].PioEndTime
			closeIpoRes.PioStartTime = tlGetAllIpoRes.Data.ClosedIpo[i].PioStartTime
			closeIpoRes.DematTransferDate = tlGetAllIpoRes.Data.ClosedIpo[i].DematTransferDate
			closeIpoRes.MandateEndDate = tlGetAllIpoRes.Data.ClosedIpo[i].MandateEndDate
			closeIpoRes.IsEmployeeCat = tlGetAllIpoRes.Data.ClosedIpo[i].IsEmployeeCat
			closeIpoRes.IsShareHolderCat = tlGetAllIpoRes.Data.ClosedIpo[i].IsShareHolderCat
			closeIpo = append(closeIpo, closeIpoRes)
		}
		getIpoResponse.ClosedIpo = closeIpo
	}(&ipoWG)
	ipoWG.Wait()

	//Update redis only if data from tradelab is different from existing data.
	existingData, _ := dbops.RedisRepo.Get(constants.IPO_KEY)

	if getIpoResponse.AllIpo.Data != nil { // if data from TL is not null then compare and update redis.

		if existingData != "" {
			var existingIpoRes models.GetAllIpoResponse
			err := json.Unmarshal([]byte(existingData), &existingIpoRes)
			if err == nil && reflect.DeepEqual(existingIpoRes, getIpoResponse) {
				loggerconfig.Info("GetAllIpo data is the same, not updating Redis key. clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			} else {
				jsonData, _ := json.Marshal(getIpoResponse)
				err := dbops.RedisRepo.Set(constants.IPO_KEY, string(jsonData), 0)
				if err != nil {
					loggerconfig.Error("GetAllIpo: Error setting data in Redis -", err)
				}
				loggerconfig.Info("GetAllIpo updated Redis key. clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
			}

		} else {

			jsonData, _ := json.Marshal(getIpoResponse)
			err := dbops.RedisRepo.Set(constants.IPO_KEY, string(jsonData), 0)
			if err != nil {
				loggerconfig.Error("GetAllIpo: Error setting data in Redis -", err)
			}
			loggerconfig.Info("GetAllIpo set Redis key. clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)

		}
		loggerconfig.Info("GetAllIpo tl resp=", helpers.LogStructAsJSON(getIpoResponse), "clientID: ", reqH.ClientId, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)

	} else { // if data from TL is null then return from redis
		if existingData != "" {
			var existingIpoRes models.GetAllIpoResponse
			err := json.Unmarshal([]byte(existingData), &existingIpoRes)
			if err != nil {
				loggerconfig.Error("GetAllIpo: Error unmarshaling data from Redis -", err)
				return apihelpers.SendInternalServerError()
			}

			apiRes.Data = existingIpoRes
			apiRes.Message = "SUCCESS"
			apiRes.Status = true
			return http.StatusOK, apiRes
		}
	}

	apiRes.Data = getIpoResponse
	apiRes.Message = "SUCCESS"
	apiRes.Status = true
	return http.StatusOK, apiRes
}

func (obj IpoObj) PlaceIpoOrder(placeIpoOrderRequest models.PlaceIpoOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + PLACEIPOORDER

	var tlPlaceIpoOrderReq TradeLabPlaceIpoOrderRequest
	tlPlaceIpoOrderReq.ClientID = placeIpoOrderRequest.ClientID
	tlPlaceIpoOrderReq.Symbol = placeIpoOrderRequest.Symbol
	tlPlaceIpoOrderReq.UpiID = placeIpoOrderRequest.UpiID

	tlPlaceOrderBids := make([]PlaceIpoOrderRequestBids, 0)
	for i := 0; i < len(placeIpoOrderRequest.Bids); i++ {
		var tlPlaceIpoOrderBidsReq PlaceIpoOrderRequestBids
		tlPlaceIpoOrderBidsReq.ActivityType = placeIpoOrderRequest.Bids[i].ActivityType
		tlPlaceIpoOrderBidsReq.Quantity = placeIpoOrderRequest.Bids[i].Quantity
		tlPlaceIpoOrderBidsReq.AtCutOff = placeIpoOrderRequest.Bids[i].AtCutOff
		tlPlaceIpoOrderBidsReq.Price = placeIpoOrderRequest.Bids[i].Price
		tlPlaceIpoOrderBidsReq.Amount = placeIpoOrderRequest.Bids[i].Amount
		tlPlaceOrderBids = append(tlPlaceOrderBids, tlPlaceIpoOrderBidsReq)
	}
	tlPlaceIpoOrderReq.Bids = tlPlaceOrderBids

	tlPlaceIpoOrderReq.AllotmentMode = placeIpoOrderRequest.AllotmentMode
	tlPlaceIpoOrderReq.BankAccount = placeIpoOrderRequest.BankAccount
	tlPlaceIpoOrderReq.BankCode = placeIpoOrderRequest.BankCode
	tlPlaceIpoOrderReq.Broker = placeIpoOrderRequest.Broker
	tlPlaceIpoOrderReq.CategoryCode = placeIpoOrderRequest.CategoryCode
	tlPlaceIpoOrderReq.CategoryCode = placeIpoOrderRequest.CategoryCode
	tlPlaceIpoOrderReq.ClientBenID = placeIpoOrderRequest.ClientBenID
	tlPlaceIpoOrderReq.ClientName = placeIpoOrderRequest.ClientName
	tlPlaceIpoOrderReq.DpID = placeIpoOrderRequest.DpID
	tlPlaceIpoOrderReq.Ifsc = placeIpoOrderRequest.Ifsc
	tlPlaceIpoOrderReq.LocationCode = placeIpoOrderRequest.LocationCode
	tlPlaceIpoOrderReq.NonAsba = placeIpoOrderRequest.NonAsba
	tlPlaceIpoOrderReq.Pan = placeIpoOrderRequest.Pan
	tlPlaceIpoOrderReq.Category = placeIpoOrderRequest.Category

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlPlaceIpoOrderReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallApi(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "PlaceIpoOrder", duration, placeIpoOrderRequest.ClientID, reqH.RequestId)
	defer res.Body.Close()
	if err != nil {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " placeIpoOrderReq call api error =", err, " uccId:", placeIpoOrderRequest.ClientID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("placeIpoOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", placeIpoOrderRequest.ClientID, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlPlaceIpoRes := TradeLabPlaceIpoOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlPlaceIpoRes)
	fmt.Println("tlPlaceIpoRes = ", tlPlaceIpoRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " placeIpoOrderRes tl status not ok =", tlPlaceIpoRes.Message, " uccId:", placeIpoOrderRequest.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlPlaceIpoRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	// fill up controller response
	var placeIpoOrderResponse models.PlaceIpoOrderResponse
	placeIpoOrderResponse.Data = tlPlaceIpoRes.Data

	loggerconfig.Info("placeIpoOrderRes tl resp=", helpers.LogStructAsJSON(placeIpoOrderResponse), " uccId:", placeIpoOrderRequest.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)

	apiRes.Data = placeIpoOrderResponse
	apiRes.Message = tlPlaceIpoRes.Message
	apiRes.Status = tlPlaceIpoRes.Success
	return http.StatusOK, apiRes
}

func (obj IpoObj) FetchIpoOrder(fetchIpoOrderRequest models.FetchIpoOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + FETCHIPOORDER

	var tlFetchIpoOrderRequest TradeLabFetchIpoOrderRequest
	tlFetchIpoOrderRequest.ClientID = fetchIpoOrderRequest.ClientID

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlFetchIpoOrderRequest)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallApi(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "FetchIpoOrder", duration, fetchIpoOrderRequest.ClientID, reqH.RequestId)
	defer res.Body.Close()
	if err != nil {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " fetchIpoOrderRequest call api error =", err, " uccId:", fetchIpoOrderRequest.ClientID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("fetchIpoOrder res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", fetchIpoOrderRequest.ClientID, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlFetchIpoOrderRes := TradeLabFetchIpoOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlFetchIpoOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " tlFetchIpoOrderRes tl status not ok =", tlFetchIpoOrderRes.Message, " uccId:", fetchIpoOrderRequest.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlFetchIpoOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}
	if len(tlFetchIpoOrderRes.Data) > 0 {
		err = setFetchIpoOrderDataRedis(FetchIpoOrderDataRedisKeyPrefix+fetchIpoOrderRequest.ClientID, tlFetchIpoOrderRes, obj)
		if err != nil {
			loggerconfig.Error("Alert Severity:TLP2-Mid, platform:", reqH.Platform, " tlFetchIpoOrderRes update the cache data, err:", err, " uccId:", fetchIpoOrderRequest.ClientID, " requestId:", reqH.RequestId)
		}
	} else {
		tlFetchIpoOrderRes, err = getFetchIpoOrderDataRedis(FetchIpoOrderDataRedisKeyPrefix+fetchIpoOrderRequest.ClientID, obj)
		if err != nil {
			json.Unmarshal([]byte(string(body)), &tlFetchIpoOrderRes) //undo any change if it might have happened due to some error
			loggerconfig.Error("Alert Severity:TLP2-Mid, platform:", reqH.Platform, " tlFetchIpoOrderRes fetch from cache failed, err:", err, " uccId:", fetchIpoOrderRequest.ClientID, " requestId:", reqH.RequestId)
		}
	}

	// fill up controller response
	var fetchIpoOrderResponse models.FetchIpoOrderResponse

	allResponseData := make([]models.FetchIpoOrderResponseData, 0)
	for i := 0; i < len(tlFetchIpoOrderRes.Data); i++ {
		var resData models.FetchIpoOrderResponseData
		resData.Symbol = tlFetchIpoOrderRes.Data[i].Symbol
		resData.Reason = tlFetchIpoOrderRes.Data[i].Reason
		resData.ApplicationNumber = tlFetchIpoOrderRes.Data[i].ApplicationNumber
		resData.ClientName = tlFetchIpoOrderRes.Data[i].ClientName
		resData.ChequeNumber = tlFetchIpoOrderRes.Data[i].ChequeNumber
		resData.ReferenceNumber = tlFetchIpoOrderRes.Data[i].ReferenceNumber
		resData.DpVerStatusFlag = tlFetchIpoOrderRes.Data[i].DpVerStatusFlag
		resData.SubBrokerCode = tlFetchIpoOrderRes.Data[i].SubBrokerCode
		resData.Depository = tlFetchIpoOrderRes.Data[i].Depository
		resData.ReasonCode = tlFetchIpoOrderRes.Data[i].ReasonCode
		resData.Pan = tlFetchIpoOrderRes.Data[i].Pan
		resData.Ifsc = tlFetchIpoOrderRes.Data[i].Ifsc
		resData.Timestamp = tlFetchIpoOrderRes.Data[i].Timestamp
		resData.BankAccount = tlFetchIpoOrderRes.Data[i].BankAccount
		resData.BankCode = tlFetchIpoOrderRes.Data[i].BankCode
		resData.DpVerReason = tlFetchIpoOrderRes.Data[i].DpVerReason
		resData.DpID = tlFetchIpoOrderRes.Data[i].DpID
		resData.Upi = tlFetchIpoOrderRes.Data[i].Upi
		if tlFetchIpoOrderRes.Data[i].UpiAmtBlocked != nil {
			resData.UpiAmtBlocked = fmt.Sprintf("%v", tlFetchIpoOrderRes.Data[i].UpiAmtBlocked)
		} else {
			resData.UpiAmtBlocked = tlFetchIpoOrderRes.Data[i].UpiAmtBlocked
		}
		// resData.Bids = tlFetchIpoOrderRes.Data[i].Bids

		// var allFetchIpoOrderResponseBids models.FetchIpoOrderResponseBids
		allFetchIpoOrderResponseBids := make([]models.FetchIpoOrderResponseBids, 0)
		for j := 0; j < len(tlFetchIpoOrderRes.Data[i].Bids); j++ {
			var fetchIpoOrderResponseBids models.FetchIpoOrderResponseBids
			fetchIpoOrderResponseBids.AtCutOff = tlFetchIpoOrderRes.Data[i].Bids[j].AtCutOff
			fetchIpoOrderResponseBids.Amount = tlFetchIpoOrderRes.Data[i].Bids[j].Amount
			fetchIpoOrderResponseBids.Quantity = tlFetchIpoOrderRes.Data[i].Bids[j].Quantity
			fetchIpoOrderResponseBids.BidReferenceNumber = tlFetchIpoOrderRes.Data[i].Bids[j].BidReferenceNumber
			fetchIpoOrderResponseBids.Series = tlFetchIpoOrderRes.Data[i].Bids[j].Series
			fetchIpoOrderResponseBids.Price = tlFetchIpoOrderRes.Data[i].Bids[j].Price
			fetchIpoOrderResponseBids.ActivityType = tlFetchIpoOrderRes.Data[i].Bids[j].ActivityType
			fetchIpoOrderResponseBids.Status = tlFetchIpoOrderRes.Data[i].Bids[j].Status
			allFetchIpoOrderResponseBids = append(allFetchIpoOrderResponseBids, fetchIpoOrderResponseBids)
		}
		resData.Bids = allFetchIpoOrderResponseBids

		resData.AllotmentMode = tlFetchIpoOrderRes.Data[i].AllotmentMode
		resData.DpVerFailCode = tlFetchIpoOrderRes.Data[i].DpVerFailCode
		resData.NonASBA = tlFetchIpoOrderRes.Data[i].NonASBA
		resData.UpiFlag = tlFetchIpoOrderRes.Data[i].UpiFlag
		resData.Category = tlFetchIpoOrderRes.Data[i].Category
		resData.LocationCode = tlFetchIpoOrderRes.Data[i].LocationCode
		resData.ClientBenID = tlFetchIpoOrderRes.Data[i].ClientBenID
		resData.ClientID = tlFetchIpoOrderRes.Data[i].ClientID
		resData.Status = tlFetchIpoOrderRes.Data[i].Status
		resData.Mode = tlFetchIpoOrderRes.Data[i].Mode
		resData.AllotmentStatus = tlFetchIpoOrderRes.Data[i].Allotmentstatus
		resData.AllotmentDate = tlFetchIpoOrderRes.Data[i].Allotmentdate
		resData.AllotmentUpdated = tlFetchIpoOrderRes.Data[i].Allotmentupdated
		resData.AllotmentQuantity = tlFetchIpoOrderRes.Data[i].Allotmentquantity
		resData.AllotmentPrice = tlFetchIpoOrderRes.Data[i].Allotmentprice
		resData.CategoryCode = tlFetchIpoOrderRes.Data[i].CategoryCode
		resData.CategoryDisplayName = tlFetchIpoOrderRes.Data[i].CategoryDisplayName
		resData.IsAmoOrder = tlFetchIpoOrderRes.Data[i].IsAmoOrder
		resData.PaymentMode = tlFetchIpoOrderRes.Data[i].PaymentMode
		resData.AmtBlockTime = tlFetchIpoOrderRes.Data[i].AmtBlockTime
		resData.Modify = tlFetchIpoOrderRes.Data[i].Modify
		resData.IsOrderModify = tlFetchIpoOrderRes.Data[i].IsOrderModify
		resData.IsBseIpo = tlFetchIpoOrderRes.Data[i].IsBseIpo
		resData.IsNseIpo = tlFetchIpoOrderRes.Data[i].IsNseIpo
		resData.UpiPaymentStatusMessage = tlFetchIpoOrderRes.Data[i].UpiPaymentStatusMessage
		resData.ExchangeUpdatedUpiBlockedAmount = tlFetchIpoOrderRes.Data[i].ExchangeUpdatedUpiBlockedAmount
		resData.IsPioOrder = tlFetchIpoOrderRes.Data[i].IsPioOrder
		allResponseData = append(allResponseData, resData)
	}
	fetchIpoOrderResponse.Data = allResponseData

	maskedFetchIpoOrderResponse, err := maskObj.Struct(fetchIpoOrderResponse)
	if err != nil {
		loggerconfig.Error("fetchIpoOrderRes Error in masking request err: ", err, " clientId: ", fetchIpoOrderRequest.ClientID, " requestid = ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("fetchIpoOrderRes tl resp=", helpers.LogStructAsJSON(maskedFetchIpoOrderResponse), " uccId:", fetchIpoOrderRequest.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
	apiRes.Data = fetchIpoOrderResponse
	apiRes.Message = tlFetchIpoOrderRes.Message
	apiRes.Status = tlFetchIpoOrderRes.Success
	return http.StatusOK, apiRes
}

func (obj IpoObj) CancelIpoOrder(cancelIpoOrderRequest models.CancelIpoOrderRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	url := obj.tradeLabURL + CANCELIPOORDER

	var tlCancelIpoOrderReq TradeLabCancelIpoOrderRequest
	tlCancelIpoOrderReq.ClientID = cancelIpoOrderRequest.ClientID
	tlCancelIpoOrderReq.Symbol = cancelIpoOrderRequest.Symbol
	tlCancelIpoOrderReq.ApplicationNumber = cancelIpoOrderRequest.ApplicationNumber

	tlAllCancelIpoOrderReqBids := make([]CancelIpoOrderRequestBids, 0)

	for i := 0; i < len(cancelIpoOrderRequest.Bids); i++ {
		var tlCancelIpoOrderReqBids CancelIpoOrderRequestBids
		tlCancelIpoOrderReqBids.Quantity = cancelIpoOrderRequest.Bids[i].Quantity
		tlCancelIpoOrderReqBids.AtCutOff = cancelIpoOrderRequest.Bids[i].AtCutOff
		tlCancelIpoOrderReqBids.Price = cancelIpoOrderRequest.Bids[i].Price
		tlCancelIpoOrderReqBids.Amount = cancelIpoOrderRequest.Bids[i].Amount
		tlCancelIpoOrderReqBids.BidReferenceNumber = cancelIpoOrderRequest.Bids[i].BidReferenceNumber
		tlCancelIpoOrderReqBids.Series = cancelIpoOrderRequest.Bids[i].Series
		tlCancelIpoOrderReqBids.ActivityType = cancelIpoOrderRequest.Bids[i].ActivityType
		tlCancelIpoOrderReqBids.Status = cancelIpoOrderRequest.Bids[i].Status
		tlAllCancelIpoOrderReqBids = append(tlAllCancelIpoOrderReqBids, tlCancelIpoOrderReqBids)
	}
	tlCancelIpoOrderReq.Bids = tlAllCancelIpoOrderReqBids

	tlCancelIpoOrderReq.UpiID = cancelIpoOrderRequest.UpiID
	tlCancelIpoOrderReq.AllotmentMode = cancelIpoOrderRequest.AllotmentMode
	tlCancelIpoOrderReq.AllotmentMode = cancelIpoOrderRequest.AllotmentMode
	tlCancelIpoOrderReq.BankAccount = cancelIpoOrderRequest.BankAccount
	tlCancelIpoOrderReq.BankCode = cancelIpoOrderRequest.BankCode
	tlCancelIpoOrderReq.Broker = cancelIpoOrderRequest.Broker
	tlCancelIpoOrderReq.ClientBenID = cancelIpoOrderRequest.ClientBenID
	tlCancelIpoOrderReq.ClientName = cancelIpoOrderRequest.ClientName
	tlCancelIpoOrderReq.DpID = cancelIpoOrderRequest.DpID
	tlCancelIpoOrderReq.Ifsc = cancelIpoOrderRequest.Ifsc
	tlCancelIpoOrderReq.LocationCode = cancelIpoOrderRequest.LocationCode
	tlCancelIpoOrderReq.NonAsba = cancelIpoOrderRequest.NonAsba
	tlCancelIpoOrderReq.Pan = cancelIpoOrderRequest.Pan

	//make payload
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(tlCancelIpoOrderReq)

	//call api
	var apiRes apihelpers.APIRes
	start := helpers.GetCurrentTimeInIST()
	res, err := apihelpers.CallApi(http.MethodPost, url, payload, reqH.DeviceType, reqH.DeviceId, reqH.Platform, reqH.ClientPublicIP, reqH.Authorization)
	duration := time.Since(start)
	helpers.RecordAPILatency(url, "CancelIpoOrder", duration, cancelIpoOrderRequest.ClientID, reqH.RequestId)
	defer res.Body.Close()
	if err != nil {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " cancelIpoOrderReq call api error =", err, " uccId:", cancelIpoOrderRequest.ClientID, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	body, err := io.ReadAll(res.Body)
	tlErrorRes := TradeLabErrorRes{}
	err = json.Unmarshal([]byte(string(body)), &tlErrorRes)
	if err == nil && tlErrorRes.Status == TLERROR {
		loggerconfig.Error("cancelIpoOrderRes res error =", tlErrorRes.Message, " statuscode: ", res.StatusCode, " uccId:", cancelIpoOrderRequest.ClientID, " requestId:", reqH.RequestId)
		apiRes.Message = tlErrorRes.Message
		apiRes.ErrorCode = strconv.Itoa(tlErrorRes.ErrorCode)
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	tlcancelIpoOrderRes := TradeLabCancelIpoOrderResponse{}
	json.Unmarshal([]byte(string(body)), &tlcancelIpoOrderRes)

	if res.StatusCode != http.StatusOK {
		loggerconfig.Error("Alert Severity:TLP1-High, platform:", reqH.Platform, " cancelIpoOrderRes tl status not ok =", tlcancelIpoOrderRes.Message, " uccId:", cancelIpoOrderRequest.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
		apiRes.Message = tlcancelIpoOrderRes.Message
		apiRes.Status = false
		return res.StatusCode, apiRes
	}

	// fill up controller response
	var cancelIpoOrderResponse models.CancelIpoOrderResponse
	cancelIpoOrderResponse.Data = tlcancelIpoOrderRes.Data

	loggerconfig.Info("cancelIpoOrderRes tl resp=", helpers.LogStructAsJSON(cancelIpoOrderResponse), " uccId:", cancelIpoOrderRequest.ClientID, " StatusCode: ", res.StatusCode, " requestId:", reqH.RequestId)
	apiRes.Data = cancelIpoOrderResponse
	apiRes.Message = tlcancelIpoOrderRes.Message
	apiRes.Status = tlcancelIpoOrderRes.Success
	return http.StatusOK, apiRes
}

func (obj IpoObj) FetchIpoData(fetchIpoDataRequest models.FetchIpoDataRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	var fetchIpoDataRes models.FetchIpoDataResponse
	err := dbops.MongoRepo.FindOne(constants.IPODATA, bson.M{"name": fetchIpoDataRequest.Name}, &fetchIpoDataRes)
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("FetchIpoData err in mongodb: ", err, " requestId:", reqH.RequestId)
		fmt.Println("error in fetching userId")
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("FetchIpoData resp=", helpers.LogStructAsJSON(fetchIpoDataRes), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)

	apiRes.Data = fetchIpoDataRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj IpoObj) FetchIpoGmpData(ipoGmpDataRequest models.FetchIpoDataRequest, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	var fetchIpoGmpDataRes models.FetchIpoGmpDataResponse
	err := dbops.MongoRepo.FindOne(constants.IPOGMPDATA, bson.M{"ipo_name": ipoGmpDataRequest.Name}, &fetchIpoGmpDataRes)
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("FetchIpoGmpData err in mongodb: ", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	loggerconfig.Info("FetchIpoGmpData resp=", helpers.LogStructAsJSON(fetchIpoGmpDataRes), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId)

	apiRes.Data = fetchIpoGmpDataRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj IpoObj) FetchEIpo(ipoDataRequest models.FetchEipoReq, reqH models.ReqHeader) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes
	var fetchNseIpoRes []models.FetchNseIpoResponse

	for _, symbol := range ipoDataRequest.IpoSymbol {
		var fetchEIpoDataRes []map[string]interface{}
		var ipoResElement models.FetchNseIpoResponse

		filter := bson.M{"$and": []bson.M{{"symbol": symbol}, {"status": ipoDataRequest.IpoStage}}}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cursor, err := dbops.MongoRepo.Find(constants.EIPODataCollection, filter)
		if err != nil {
			loggerconfig.Error("Alert Severity:P1-High, FetchEIpo error in MongoDB and error is: ", err, "clientID:", reqH.ClientId, "requestId:", reqH.RequestId, "platform:", reqH.Platform, "clientVersion", reqH.ClientVersion)
			return apihelpers.SendInternalServerError()
		}
		defer cursor.Close(ctx)
		err = cursor.All(ctx, &fetchEIpoDataRes)
		if err != nil && err.Error() != constants.MongoNoDocError {
			loggerconfig.Error("Alert Severity:P1-High, FetchEIpo error in MongoDB", err, "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, "platform:", reqH.Platform, "clientVersion", reqH.ClientVersion)
			return apihelpers.SendInternalServerError()
		}

		for _, eIpoData := range fetchEIpoDataRes {

			ipoResElement.IpoSymbol = extractIpoSymbol(eIpoData, reqH)

			if ipoResElement.IpoSymbol != symbol {
				loggerconfig.Info(" FetchEIpo (services): Symbol are different ", ipoResElement.IpoSymbol, symbol, "clientID:", reqH.ClientId, "requestId:", reqH.RequestId)
				continue
			}

			ipoResElement.IpoStage = extractIpoStage(eIpoData, reqH)
			ipoResElement.Series = extractIpoSeries(eIpoData, ipoDataRequest.IpoStage, reqH)
			ipoResElement.DRHPLink = extractDRHPLink(eIpoData, reqH)
			ipoResElement.BidDeatils = extractBidDetails(eIpoData, reqH)

			if ipoResElement.Series == "SME" {
				ipoResElement.SubsTimes = extractSubscriptionTimes(eIpoData, reqH)
			}
		}

		if ipoResElement.IpoSymbol == "" {
			ipoResElement.IpoSymbol = symbol
		}
		fetchNseIpoRes = append(fetchNseIpoRes, ipoResElement)
	}

	loggerconfig.Info("FetchEIpo resp=", helpers.LogStructAsJSON(fetchNseIpoRes), "clientID: ", reqH.ClientId, " requestId:", reqH.RequestId, "platform:", reqH.Platform, "clientVersion", reqH.ClientVersion)

	apiRes.Data = fetchNseIpoRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func setFetchIpoOrderDataRedis(key string, data TradeLabFetchIpoOrderResponse, obj IpoObj) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return dbops.RedisRepo.Set(key, string(jsonData), 0)
}

func getFetchIpoOrderDataRedis(key string, obj IpoObj) (TradeLabFetchIpoOrderResponse, error) {
	var data TradeLabFetchIpoOrderResponse

	jsonData, err := dbops.RedisRepo.Get(key)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func getGetAllIPODataRedis(key string, obj IpoObj) (TradeLabGetAllIpoResponse, error) {
	var data TradeLabGetAllIpoResponse

	jsonData, err := dbops.RedisRepo.Get(key)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func extractIpoSeries(eIpoData map[string]interface{}, stage string, reqH models.ReqHeader) string {
	if ipoData, ok := eIpoData["ipodata"].(map[string]interface{}); ok {
		if stage == constants.CLOSED {
			if ipoSeries, ok := ipoData["securityType"].(string); ok {
				return ipoSeries
			} else {
				loggerconfig.Info(" FetchEIpo / extractIpoSeries (services): Missing or invalid ipodata for closed Ipo in response data", "clientID:", reqH.ClientId, "requestId:", reqH.RequestId)
			}
		} else {
			if ipoSeries, ok := ipoData["series"].(string); ok {
				return ipoSeries
			} else {
				loggerconfig.Info(" FetchEIpo / extractIpoSeries (services): Missing or invalid ipodata in response data", "clientID:", reqH.ClientId, "requestId:", reqH.RequestId)
			}
		}
	}
	return ""
}

func extractIpoSymbol(eIpoData map[string]interface{}, reqH models.ReqHeader) string {
	if symbol, ok := eIpoData["symbol"].(string); ok {
		return symbol
	}
	loggerconfig.Info("FetchEIpo / extractIpoSymbol (services): Missing or invalid symbol in response data", "clientID:", reqH.ClientId, "requestId:", reqH.RequestId)
	return ""
}

func extractIpoStage(eIpoData map[string]interface{}, reqH models.ReqHeader) string {
	if ipoStage, ok := eIpoData["status"].(string); ok {
		return ipoStage
	}
	loggerconfig.Info("FetchEIpo / extractIpoStage (services): Missing or invalid status in response data", "clientID:", reqH.ClientId, "requestId:", reqH.RequestId)
	return ""
}

func extractSubscriptionTimes(eIpoData map[string]interface{}, reqH models.ReqHeader) string {
	if eIpoAggregatedData, ok := eIpoData["eipoaggregateddata"].(map[string]interface{}); ok {

		if demandGraphALL, ok := eIpoAggregatedData["demandGraphALL"].(map[string]interface{}); ok {
			if noOfTimesIssueSubscribed, ok := demandGraphALL["noOfTimesIssueSubscribed"].(string); ok {
				return noOfTimesIssueSubscribed
			} else {
				loggerconfig.Info(" FetchEIpo (services): Missing or invalid noOfTimesIssueSubscribed in demandGraphALL", "clientID:", reqH.ClientId, "requestId:", reqH.RequestId)
			}
		} else {
			loggerconfig.Info("FetchEIpo (services): Missing or invalid demandGraphALL in eipoaggregateddata", "clientID:", reqH.ClientId, "requestId:", reqH.RequestId)
		}

		if demandGraph, ok := eIpoAggregatedData["demandGraph"].(map[string]interface{}); ok {
			if noOfTimesIssueSubscribed, ok := demandGraph["noOfTimesIssueSubscribed"].(string); ok {
				return noOfTimesIssueSubscribed
			} else {
				loggerconfig.Info(" FetchEIpo (services): Missing or invalid noOfTimesIssueSubscribed in demandGraph", "clientID:", reqH.ClientId, "requestId:", reqH.RequestId)
			}
		} else {
			loggerconfig.Info("FetchEIpo (services): Missing or invalid demandGraph in eipoaggregateddata", "clientID:", reqH.ClientId, "requestId:", reqH.RequestId)
		}
	}
	return ""
}

func extractDRHPLink(eIpoData map[string]interface{}, reqH models.ReqHeader) string {
	if eIpoAggregatedData, ok := eIpoData["eipoaggregateddata"].(map[string]interface{}); ok {
		if issueInfo, ok := eIpoAggregatedData["issueInfo"].(map[string]interface{}); ok {
			if dataListRaw, ok := issueInfo["dataList"]; ok {
				if dataList, ok := dataListRaw.(primitive.A); ok {
					for _, dataItem := range dataList {
						if data, ok := dataItem.(map[string]interface{}); ok {
							if title, ok := data["title"].(string); ok && title == "Red Herring Prospectus" {
								if rhpValue, ok := data["value"].(string); ok {
									return rhpValue
								} else {
									loggerconfig.Info("FetchEIpo (services): Missing or invalid value for 'Red Herring Prospectus'", "clientID:", reqH.ClientId, "requestId:", reqH.RequestId)
								}
							}
						} else {
							loggerconfig.Info("FetchEIpo (services): data type of dataItem is incorrect", "clientID:", reqH.ClientId, "requestId:", reqH.RequestId)
						}
					}
				} else {
					loggerconfig.Info("FetchEIpo (services): data type of dataList is invalid", "clientID:", reqH.ClientId, "requestId:", reqH.RequestId)
				}
			} else {
				loggerconfig.Info("FetchEIpo (services): Missing dataList key in issueInfo", "clientID:", reqH.ClientId, "requestId:", reqH.RequestId)
			}
		} else {
			loggerconfig.Info("FetchEIpo (services): Missing or invalid issueInfo in eipoaggregateddata", "clientID:", reqH.ClientId, "requestId:", reqH.RequestId)
		}
	} else {
		loggerconfig.Info("FetchEIpo (services): Missing or invalid eipoaggregateddata in response data", "clientID:", reqH.ClientId, "requestId:", reqH.RequestId)
	}
	return ""
}

func extractBidDetails(eIpoData map[string]interface{}, reqH models.ReqHeader) interface{} {
	if eIpoAggregatedData, ok := eIpoData["eipoaggregateddata"].(map[string]interface{}); ok {
		if activeCat, ok := eIpoAggregatedData["activeCat"].(map[string]interface{}); ok {
			if dataListRaw, ok := activeCat["dataList"].(primitive.A); ok && len(dataListRaw) > 2 {
				return map[string]interface{}{
					"dataList":   filterDataList(dataListRaw),
					"updateTime": activeCat["updateTime"],
				}
			}
		}

		if bidDetails, ok := eIpoAggregatedData["bidDetails"].(map[string]interface{}); ok {
			if dataListRaw, ok := bidDetails["data"].(primitive.A); ok {
				return map[string]interface{}{
					"dataList":   filterDataList(dataListRaw),
					"updateTime": bidDetails["updateTime"],
				}
			}
		}
	}
	return nil
}

func filterDataList(dataListRaw primitive.A) []map[string]interface{} {
	var filteredDataList []map[string]interface{}
	for _, item := range dataListRaw {
		if dataItem, ok := item.(map[string]interface{}); ok {
			if srNo, ok := dataItem["srNo"]; ok {
				if srNoStr, isString := srNo.(string); isString && len(srNoStr) == 1 {
					filteredDataList = append(filteredDataList, dataItem)
				} else if srNo == nil {
					filteredDataList = append(filteredDataList, dataItem)
				}
			}
		}
	}
	return filteredDataList
}
