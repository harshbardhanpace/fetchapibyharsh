package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"space/constants"
	"space/db"
	"space/loggerconfig"

	"space/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var HolidayCalendar models.Calendar

func GetNSECalendar() (*models.NSECalendarModel, error) {
	//Not Putting in constants as these have no reusablility and makes not difference if these are in constants or directly here
	url := "https://www.nseindia.com/api/holiday-master?type=trading"

	//Create a transport that disables HTTP/2
	transport := &http.Transport{
		DisableKeepAlives: true,
		ForceAttemptHTTP2: false,
	}

	//Create a client with the custom transport
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 10,
	}

	//Retry mechanism
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}

		//Not Putting in constants as these have no reusablility and makes not difference if these are in constants or directly here
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Referer", "https://www.nseindia.com/")
		req.Header.Set("Cookie", "*abck=570FBDA88A8DA38FC1AC8FA8166CCEFA~-1~YAAQLG/ZFygiJQSRAQAAvMjAKwwDnvMbjP2IB5Uaij/7RE6paxKv2xQdJ+xyJGTSbJNIVqHfW+VuFpwJ0uiL9Gs5+SgHQJQYAD2QSMWG1YZAkXVn2PXVh+Qw9XeJpLjcO6dtav1+N8A9fp7Y1JYmI7xJadg2OhvR4MAWRoZiw7XFm4i31Qg6TNURIpiD5hLwfIQWtJlV/kzjn25V8TIT+J9RQPx8OmvQW74P5nO73onV+BsaZBXR0j2w4NL7J8tbStJGFoS3g9YEMPH0tW+AZ8nBeQiDXp/5yXXAQh+nx5QpNKzm1M79K+optc0Bq91fIbJkD2Nv3CGvYAJb5IaEre3/XE/QPTCYk8KGpOD0LNfDiKJRTLxIndNMZQ==~-1~-1~-1; ak*bmsc=E4A1C7AF000821154DEBFCB1FCA131E0~000000000000000000000000000000~YAAQLG/ZFykiJQSRAQAAvMjAKxjrEoZZEVnyRyKZzCLsoh2f9eNV0UGHmOC3rU17r1x+IPXsLuePX+XuNBtG/9Wzl1/YPMIotL9wFjIr1m83WuDavHffq/rp/mXSuOrNHeSpczU35XmREKzQyJHPpxMgjRAkOF5R9UDF9qMrgetdwpe6+u2y07Znql+UrXU0RnAXWkRwdNnhz/pPdPZuBporCsIra8rd3jwLs4MfNxBGbDHnGvhbg1mQ8ThmeaiT02FAJlsPdLQVh8M8bNtAoCcuLQXanfeDjrviuUpdMJiIgLocxzhOWA8D+1MJyhQv5kU/p9gl6QZ5ONqvH2bLhYApRdKGR4FbDT64BN1L7dVI34YDFaxsrtHsS4DoTQ==; bm_sz=7BDAE15F230D619AAF8F3E3A52872688~YAAQLG/ZFyoiJQSRAQAAvMjAKxgBlEOvcI2XRiY+gssEYncbGaaJufg71QmkUqmja9EhsV0QqChzg8rIm08o4Vt1etWUkS4QbuMgZ/LndsPoQ2phZfSaJv/CUi/26KzSssYZ2P4xYO514vXnDNL1tdcJyN227DDApwWIouSYL76U7eq342zJyMKo0kzoHiVGUfZOkWnvXfyc5P+Lx/cwvMpI1452Hv/eU7XUNzdOlqqMp0cW98Buh5SBxaoQTzG0KIZI08GCAjcMYpZNdxDgdnJRih/2fptAXIdMHSwzr3vAtWKGQmWFvwU6NAC1pm5dbb85+shx/U14FwaApH2tPMWUazrbDez7Qn44tln4~4404792~3356216")

		resp, err := client.Do(req)
		if err != nil {
			loggerconfig.Info("Attempt:", attempt+1, " failed:", err)
			if attempt == maxRetries-1 {
				return nil, fmt.Errorf("error sending request after %d attempts: %v", maxRetries, err)
			}
			time.Sleep(time.Second * 2)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v", err)
		}

		var nseCalendar models.NSECalendarModel
		err = json.Unmarshal(body, &nseCalendar)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
		}

		return &nseCalendar, nil
	}

	return nil, fmt.Errorf("failed to get NSE calendar after %d attempts", maxRetries)
}

func GetCurrentYearCalendar() models.Calendar {
	currentYear := GetCurrentTimeInIST().Year()
	startDate := time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(currentYear+1, 1, 1, 0, 0, 0, 0, time.Local)

	calendar := models.Calendar{}

	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		dayOfWeek := d.Weekday().String()
		var desc string
		isHoliday := false
		if dayOfWeek == "Saturday" || dayOfWeek == "Sunday" {
			desc = "Weekend"
			isHoliday = true
		}
		dateDetails := models.DateDetails{
			Date:        d.Format("02-Jan-2006"),
			DayOfWeek:   dayOfWeek,
			Description: desc,
			IsHoliday:   isHoliday,
		}
		calendar.Date = append(calendar.Date, dateDetails)
	}

	return calendar
}

func MergeCalendars(calendar models.Calendar, nseCalendar models.NSECalendarModel) (models.Calendar, error) {
	nseHolidays := make(map[string]string)
	for _, holiday := range nseCalendar.Cm {
		date, err := time.Parse("02-Jan-2006", holiday.TradingDate)
		if err != nil {
			return models.Calendar{}, fmt.Errorf("error parsing NSE date: %v", err)
		}
		nseHolidays[date.Format("02-Jan-2006")] = holiday.Description
	}

	for i, date := range calendar.Date {
		if description, isHoliday := nseHolidays[date.Date]; isHoliday {
			calendar.Date[i].IsHoliday = true
			calendar.Date[i].Description = description
		}
	}

	return calendar, nil
}

func fetchAndUploadCalendar(mongoObj db.MongoDatabase) error {
	// Generate the current year calendar
	calendar := GetCurrentYearCalendar()
	nseCalendar, err := GetNSECalendar()
	if err != nil {
		return fmt.Errorf("failed to get NSE calendar: %v", err)
	}
	mergedCalendar, err := MergeCalendars(calendar, *nseCalendar)
	if err != nil {
		return fmt.Errorf("failed to merge calendars: %v", err)
	}
	HolidayCalendar = mergedCalendar

	docs := make([]interface{}, len(mergedCalendar.Date))
	for i, dateDetail := range calendar.Date {
		docs[i] = dateDetail
	}
	err = mongoObj.InsertMany(constants.HOLIDAYCALENDAR, docs)
	if err != nil {
		return fmt.Errorf("failed to upload calendar: %v", err)
	}
	HolidayCalendar.Date = SortAndRemoveDuplicates(HolidayCalendar.Date)
	loggerconfig.Info("ProcessAndUploadCalendar fetchAndUploadCalendar Merged calendar uploaded successfully.")
	loggerconfig.Info("ProcessAndUploadCalendar fetchAndUploadCalendar, holiday calendar:", HolidayCalendar)
	return nil
}

func ProcessAndUploadCalendar() error {
	mongoObj := db.GetMongoDBObj()
	/*
		Logic:
		1. Fetch from mongo.
			a. Null Data, then get current calendar, and nse holiday calendar. Make merged calendar and upload
			b. Data Found
				A. Data of current year and thus just return
				B. Data not of current year, hence delete old, and do as step 1. a.
	*/
	existingData, err := FindMongoCalendar()
	if err != nil {
		return fmt.Errorf("failed to fetch existing data: %v", err)
	}
	loggerconfig.Info("ProcessAndUploadCalendar data in MongoDB:", HolidayCalendar)
	if len(existingData.Date) == 0 {
		loggerconfig.Info("ProcessAndUploadCalendar No data exists for the current year, uploading new data.")
		return fetchAndUploadCalendar(mongoObj)
	} else {
		loggerconfig.Info("existingData[0].Date[0].Date:", existingData.Date[0].Date)
		if len(existingData.Date) > 0 && !IsCurrentYear(existingData.Date[0].Date) {
			loggerconfig.Info("ProcessAndUploadCalendar Data exists but is outdated, deleting old data and uploading new data.")
			err = mongoObj.DeleteMany(constants.HOLIDAYCALENDAR, bson.M{})
			if err != nil {
				return fmt.Errorf("failed to delete outdated data: %v", err)
			}
			return fetchAndUploadCalendar(mongoObj)
		} else {
			HolidayCalendar = existingData
			loggerconfig.Info("ProcessAndUploadCalendar Current year data already exists in the collection.")
		}
	}
	HolidayCalendar.Date = SortAndRemoveDuplicates(HolidayCalendar.Date)
	loggerconfig.Info("ProcessAndUploadCalendar, length:", len(HolidayCalendar.Date), " holiday calendar:", HolidayCalendar)
	return nil
}

func FindMongoCalendar() (models.Calendar, error) {
	var result models.Calendar
	mongoObj := db.GetMongoDBObj()
	cursor, err := mongoObj.FindAllMongo(constants.HOLIDAYCALENDAR, bson.M{})
	if err != nil {
		return result, err
	}
	defer cursor.Close(context.Background())

	var dates []models.DateDetails
	for cursor.Next(context.Background()) {
		var calendar models.DateDetails
		if err := cursor.Decode(&calendar); err != nil {
			return result, err
		}
		dates = append(dates, calendar)
	}

	if err := cursor.Err(); err != nil {
		return result, err
	}

	result.Date = dates
	return result, nil
}

func IsCurrentYear(dateStr string) bool {
	dateFormat := "02-Jan-2006"

	parsedTime, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return false
	}

	currentYear := GetCurrentTimeInIST().Year()
	dateYear := parsedTime.Year()

	return dateYear == currentYear
}

func FindDateIndex(date string) int {
	var err error
	originalFormat := "02-01-2006" //not putting in constants as it would be confusing to understand what format is being converted to which
	desiredFormat := "02-Jan-2006"
	date, err = ConvertFormat(date, originalFormat, desiredFormat)
	if err != nil {
		return -1
	}
	dateFormat := "02-Jan-2006"
	left, right := 0, len(HolidayCalendar.Date)-1

	for left <= right {
		mid := left + (right-left)/2
		midDate, err := time.Parse(dateFormat, HolidayCalendar.Date[mid].Date)
		if err != nil {
			return -1
		}

		searchDate, err := time.Parse(dateFormat, date)
		if err != nil {
			return -1
		}

		if searchDate.Equal(midDate) {
			return mid
		} else if searchDate.Before(midDate) {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}

	return -1
}

func ConvertFormat(dateStr, originalFormat, desiredFormat string) (string, error) {
	t, err := time.Parse(originalFormat, dateStr)
	if err != nil {
		return "", err
	}

	formattedDate := t.Format(desiredFormat)
	return formattedDate, nil
}

// SortAndRemoveDuplicates sorts the DateDetails slice and removes duplicates based on the Date field
func SortAndRemoveDuplicates(dates []models.DateDetails) []models.DateDetails {
	layout := "02-Jan-2006"

	sort.Slice(dates, func(i, j int) bool {
		date1, err1 := time.Parse(layout, dates[i].Date)
		date2, err2 := time.Parse(layout, dates[j].Date)

		if err1 != nil || err2 != nil {
			loggerconfig.Error("SortAndRemoveDuplicates Error parsing date:", err1, err2)
			return false
		}

		return date1.Before(date2)
	})

	// Remove duplicates if any exist
	result := []models.DateDetails{}
	seenDates := map[string]bool{}

	for _, dateDetails := range dates {
		if !seenDates[dateDetails.Date] {
			seenDates[dateDetails.Date] = true
			result = append(result, dateDetails)
		}
	}

	return result
}

func Next830AMUnix(currentTime time.Time) int64 {
	// Create the target time for 8:30 AM on the current day
	today830AM := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 8, 30, 0, 0, currentTime.Location())

	// If current time is before 8:30 AM, return today's 8:30 AM in Unix time
	// Otherwise, return the Unix time for 8:30 AM of the next day
	if currentTime.Before(today830AM) {
		return today830AM.Unix()
	}
	return today830AM.Add(24 * time.Hour).Unix()
}

func ConvertDateFormat(inputDate, CurrentFormat, resultantFormat string) (string, error) {
	// Parse the input date using the provided input format
	parsedDate, err := time.Parse(CurrentFormat, inputDate)
	if err != nil {
		return "", fmt.Errorf("error parsing date: %v", err)
	}

	// Format the parsed date to the desired output format
	formattedDate := parsedDate.Format(resultantFormat)
	return formattedDate, nil
}
