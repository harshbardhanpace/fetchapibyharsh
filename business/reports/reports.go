package reports

import (
	"encoding/json"
	"fmt"
	"net/http"
	apihelpers "space/apiHelpers"
	"space/business/backoffice"
	"space/constants"
	v1 "space/controllers/api/v1"
	"space/dbops"
	"space/helpers"
	"space/helpers/cache"
	"space/loggerconfig"
	"space/models"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx/v3"
)

type ReportsObj struct {
	redisCli cache.RedisCache
}

func InitReportsProvider(redisCli cache.RedisCache) ReportsObj {
	defer models.HandlePanic()
	ReportsObj := ReportsObj{redisCli: redisCli}

	return ReportsObj
}

func (obj ReportsObj) ViewDPCharges(viewDPchargesReq models.DPChargesReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {
	var viewDPchargesRes models.ViewDPchargesRes
	var err error

	getBillDetailsCdslReq := models.GetBillDetailsCdslReq{}
	getBillDetailsCdslReq.UserID = viewDPchargesReq.UserID
	getBillDetailsCdslReq.DFDateFr = viewDPchargesReq.DFDateFr
	getBillDetailsCdslReq.DFDateTo = viewDPchargesReq.DFDateTo

	fileName := constants.DpChargesReport + strings.ToUpper(viewDPchargesReq.UserID) + "_" + strings.ToUpper(viewDPchargesReq.DFDateFr) + "_to_" + strings.ToUpper(viewDPchargesReq.DFDateTo)
	storedReportData, err := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" || err != nil {
		err = json.Unmarshal([]byte(storedReportData), &viewDPchargesRes)
		if err != nil {
			loggerconfig.Error("ViewDPCharges, error in unmarshalling storedReportData : ", err, " for fileName fileName: ", fileName, "reqId: ", reqH.RequestId)
		}
	}
	if storedReportData == "" {
		loggerconfig.Error("readMessages Failed to read from redis:", err, "for fileName fileName: ", fileName)

		code, Resp := backoffice.GetBillDetailsCdslData(getBillDetailsCdslReq, reqH)
		if code != http.StatusOK {
			loggerconfig.Error("ViewDPCharges there is some error in fetching data from Shilpi api function: code : ", code)
		}

		var getBillDetailsCdslRes models.GetBillDetailsCdslRes
		var ok bool

		if getBillDetailsCdslRes, ok = Resp.Data.(models.GetBillDetailsCdslRes); ok {
		} else {
			loggerconfig.Error("ViewDPCharges, Invalid data format or not of type GetBillDetailsCdslRes", "reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		for i := 0; i < len(getBillDetailsCdslRes.GetBillDetailsCdsl); i++ {
			var dpcharges models.DPcharges

			dpcharges.Charges = getBillDetailsCdslRes.GetBillDetailsCdsl[i].Charges
			dpcharges.Qty = getBillDetailsCdslRes.GetBillDetailsCdsl[i].Qty
			dpcharges.Gst = getBillDetailsCdslRes.GetBillDetailsCdsl[i].Gst
			dpcharges.InstrumentName = getBillDetailsCdslRes.GetBillDetailsCdsl[i].InstrumentName
			dpcharges.ChargesDetails = getBillDetailsCdslRes.GetBillDetailsCdsl[i].ChargesDetails
			dpcharges.TotalCharges = getBillDetailsCdslRes.GetBillDetailsCdsl[i].TotalCharges

			charges, err := strconv.ParseFloat(dpcharges.Charges, 64)
			if err != nil {
				loggerconfig.Error("ViewDPCharges, error in parsing charges string to float64 Error:", err, "reqId: ", reqH.RequestId)
				return apihelpers.SendInternalServerError()
			}

			gst, err := strconv.ParseFloat(dpcharges.Gst, 64)
			if err != nil {
				loggerconfig.Error("ViewDPCharges, error in parsing gst string to float64 Error:", err, "reqId: ", reqH.RequestId)
				return apihelpers.SendInternalServerError()
			}
			viewDPchargesRes.Charges += charges
			viewDPchargesRes.GST += gst

			viewDPchargesRes.DPChargesList = append(viewDPchargesRes.DPChargesList, dpcharges)
		}

		viewDPchargesRes.TotalCharges = viewDPchargesRes.Charges + viewDPchargesRes.GST
		viewDPchargesRes.UserDetails = profileData
		reportData, err := json.Marshal(viewDPchargesRes)
		if err != nil {
			loggerconfig.Error("Error in marshalling viewDpChargesRes", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("ViewDPCharges Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}

	var apiRes apihelpers.APIRes
	apiRes.Data = viewDPchargesRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

// CreateExcelDPChargesReport generates the Excel file for DP charges data
func CreateExcelDPChargesReport(data models.ViewDPchargesRes) (*xlsx.File, error) {
	// Create a new workbook
	file := xlsx.NewFile()

	// Create a new sheet
	sheet, err := file.AddSheet("DPCharges")
	if err != nil {
		return nil, err
	}

	// Add the title row for the report
	addTitleRowDPCharges(sheet)

	// Add user details below the title
	writeUserDetails(sheet, data.UserDetails, constants.DPChargesReportName)

	// Add an empty row for spacing between sections
	sheet.AddRow()

	// Write the summary of DP charges
	writeDPChargesSummary(sheet, data)

	// Add space between sections
	sheet.AddRow()

	// Write the detailed DP charges data
	writeDPChargesData(sheet, data.DPChargesList)

	return file, nil
}

// addTitleRowDPCharges adds the main title to the top of the Excel sheet for DP Charges Report
func addTitleRowDPCharges(sheet *xlsx.Sheet) {
	// Create the title row
	titleRow := sheet.AddRow()
	titleRow.SetHeight(20) // Optionally set height for emphasis
	titleCell := titleRow.AddCell()
	titleCell.Merge(3, 0) // Merge first 4 cells for the title
	titleCell.SetString("DP Charges")
	titleCell.GetStyle().Font.Bold = true
	titleCell.GetStyle().Alignment.Horizontal = "center"

	// Add an empty row after the title for spacing
	sheet.AddRow()
}

// writeDPChargesSummary writes the summary of the DP charges in the Excel file
func writeDPChargesSummary(sheet *xlsx.Sheet, data models.ViewDPchargesRes) {
	// Create headers for the summary
	summaryHeaders := []string{"DP Charges Summary", "Amount"}
	headerRow := sheet.AddRow()
	for _, header := range summaryHeaders {
		headerRow.AddCell().SetString(header)
	}

	// Define the order and labels of the summary data
	summaryData := []struct {
		Label string
		Value float64
	}{
		{"Charges", data.Charges},
		{"GST", data.GST},
		{"Total Charges", data.TotalCharges},
	}

	// Populate the summary data
	for _, item := range summaryData {
		row := sheet.AddRow()
		row.AddCell().SetString(item.Label)
		row.AddCell().SetString(strconv.FormatFloat(item.Value, 'f', 2, 64))
	}
}

// writeDPChargesData writes the detailed DP charges data in the Excel file
func writeDPChargesData(sheet *xlsx.Sheet, dpChargesList []models.DPcharges) {
	// Create headers for the DP charges data
	dpChargesHeaders := []string{"Instrument Name", "Charges Details", "Quantity", "Charges", "GST", "Total Charges"}
	headerRow := sheet.AddRow()
	for _, header := range dpChargesHeaders {
		headerRow.AddCell().SetString(header)
	}

	// Write each DP charge as a row
	for _, dpCharge := range dpChargesList {
		dataRow := sheet.AddRow()
		dataRow.AddCell().SetString(dpCharge.InstrumentName)
		dataRow.AddCell().SetString(dpCharge.ChargesDetails)
		dataRow.AddCell().SetString(dpCharge.Qty)
		dataRow.AddCell().SetString(dpCharge.Charges)
		dataRow.AddCell().SetString(dpCharge.Gst)
		dataRow.AddCell().SetString(dpCharge.TotalCharges)
	}
}

func (obj ReportsObj) DownloadDPCharges(downloadDPChargesReq models.DPChargesReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {
	var viewDPchargesRes models.ViewDPchargesRes
	var err error

	getBillDetailsCdslReq := models.GetBillDetailsCdslReq{}
	getBillDetailsCdslReq.UserID = downloadDPChargesReq.UserID
	getBillDetailsCdslReq.DFDateFr = downloadDPChargesReq.DFDateFr
	getBillDetailsCdslReq.DFDateTo = downloadDPChargesReq.DFDateTo

	fileName := constants.DpChargesReport + strings.ToUpper(downloadDPChargesReq.UserID) + "_" + strings.ToUpper(downloadDPChargesReq.DFDateFr) + "_to_" + strings.ToUpper(downloadDPChargesReq.DFDateTo)
	storedReportData, _ := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" {
		err = json.Unmarshal([]byte(storedReportData), &viewDPchargesRes)
		if err != nil {
			loggerconfig.Error("DownloadDPCharges, error in unmarshalling storedReportData : ", err, " for fileName fileName: ", fileName)
		}
	}
	if storedReportData == "" {
		loggerconfig.Error("DownloadDPCharges, readMessages Failed to read from redis:", err, "for fileName fileName: ", fileName)

		code, Resp := backoffice.GetBillDetailsCdslData(getBillDetailsCdslReq, reqH)
		if code != http.StatusOK {
			loggerconfig.Error("DownloadDPCharges, there is some error in fetching data from Shilpi api function: code : ", code)
		}

		var getBillDetailsCdslRes models.GetBillDetailsCdslRes
		var ok bool

		if getBillDetailsCdslRes, ok = Resp.Data.(models.GetBillDetailsCdslRes); ok {
		} else {
			loggerconfig.Error("DownloadDPCharges, Invalid data format or not of type GetBillDetailsCdslRes", "reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		for i := 0; i < len(getBillDetailsCdslRes.GetBillDetailsCdsl); i++ {
			var dpcharges models.DPcharges

			dpcharges.Charges = getBillDetailsCdslRes.GetBillDetailsCdsl[i].Charges
			dpcharges.Qty = getBillDetailsCdslRes.GetBillDetailsCdsl[i].Qty
			dpcharges.Gst = getBillDetailsCdslRes.GetBillDetailsCdsl[i].Gst
			dpcharges.InstrumentName = getBillDetailsCdslRes.GetBillDetailsCdsl[i].InstrumentName
			dpcharges.ChargesDetails = getBillDetailsCdslRes.GetBillDetailsCdsl[i].ChargesDetails
			dpcharges.TotalCharges = getBillDetailsCdslRes.GetBillDetailsCdsl[i].TotalCharges

			charges, err := strconv.ParseFloat(dpcharges.Charges, 64)
			if err != nil {
				loggerconfig.Error("DownloadDPCharges, error in parsing charges string to float64 Error:", err, "reqId: ", reqH.RequestId)
				return apihelpers.SendInternalServerError()
			}

			gst, err := strconv.ParseFloat(dpcharges.Gst, 64)
			if err != nil {
				loggerconfig.Error("DownloadDPCharges, error in parsing gst string to float64 Error:", err, "reqId: ", reqH.RequestId)
				return apihelpers.SendInternalServerError()
			}
			viewDPchargesRes.Charges += charges
			viewDPchargesRes.GST += gst

			viewDPchargesRes.DPChargesList = append(viewDPchargesRes.DPChargesList, dpcharges)
		}

		viewDPchargesRes.TotalCharges = viewDPchargesRes.Charges + viewDPchargesRes.GST
		viewDPchargesRes.UserDetails = profileData
		reportData, err := json.Marshal(viewDPchargesRes)
		if err != nil {
			loggerconfig.Error("DownloadDPCharges, Error in marshalling viewDpChargesRes", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("DownloadDPCharges Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}

	file, err := CreateExcelDPChargesReport(viewDPchargesRes)
	if err != nil {
		loggerconfig.Error("DownloadDPCharges, Error in getting excel file, error: ", err, "reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	s3FileName := fileName + ".xlsx"
	expiryHours := int64(24)

	// Generate a pre-signed URL for the XLSX file
	url, err := helpers.UploadFileToS3AndGetPresignedURL(constants.DpChargesS3FolderName, s3FileName, file, expiryHours)
	if err != nil {
		loggerconfig.Error("DownloadDPCharges, failed to generate pre-signed URL:", err, "reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var downloadDPChargesRes models.DownloadDPChargesRes
	downloadDPChargesRes.DownloadUrl = url

	var apiRes apihelpers.APIRes
	apiRes.Data = downloadDPChargesRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ReportsObj) SendEmailDPCharges(sendEmailDPChargesReq models.DPChargesReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {
	var viewDPchargesRes models.ViewDPchargesRes
	var err error

	getBillDetailsCdslReq := models.GetBillDetailsCdslReq{}
	getBillDetailsCdslReq.UserID = sendEmailDPChargesReq.UserID
	getBillDetailsCdslReq.DFDateFr = sendEmailDPChargesReq.DFDateFr
	getBillDetailsCdslReq.DFDateTo = sendEmailDPChargesReq.DFDateTo

	fileName := constants.DpChargesReport + strings.ToUpper(sendEmailDPChargesReq.UserID) + "_" + strings.ToUpper(sendEmailDPChargesReq.DFDateFr) + "_to_" + strings.ToUpper(sendEmailDPChargesReq.DFDateTo)
	storedReportData, _ := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" {
		err = json.Unmarshal([]byte(storedReportData), &viewDPchargesRes)
		if err != nil {
			loggerconfig.Error("SendEmailDPCharges, error in unmarshalling storedReportData : ", err, " for fileName fileName: ", fileName)
		}
	}
	if storedReportData == "" {
		loggerconfig.Error("SendEmailDPCharges, readMessages Failed to read from redis:", err, "for fileName fileName: ", fileName)

		code, Resp := backoffice.GetBillDetailsCdslData(getBillDetailsCdslReq, reqH)
		if code != http.StatusOK {
			loggerconfig.Error("SendEmailDPCharges, there is some error in fetching data from Shilpi api function: code : ", code)
		}

		var getBillDetailsCdslRes models.GetBillDetailsCdslRes
		var ok bool

		if getBillDetailsCdslRes, ok = Resp.Data.(models.GetBillDetailsCdslRes); ok {
		} else {
			loggerconfig.Error("SendEmailDPCharges, Invalid data format or not of type GetBillDetailsCdslRes", "reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		for i := 0; i < len(getBillDetailsCdslRes.GetBillDetailsCdsl); i++ {
			var dpcharges models.DPcharges

			dpcharges.Charges = getBillDetailsCdslRes.GetBillDetailsCdsl[i].Charges
			dpcharges.Qty = getBillDetailsCdslRes.GetBillDetailsCdsl[i].Qty
			dpcharges.Gst = getBillDetailsCdslRes.GetBillDetailsCdsl[i].Gst
			dpcharges.InstrumentName = getBillDetailsCdslRes.GetBillDetailsCdsl[i].InstrumentName
			dpcharges.ChargesDetails = getBillDetailsCdslRes.GetBillDetailsCdsl[i].ChargesDetails
			dpcharges.TotalCharges = getBillDetailsCdslRes.GetBillDetailsCdsl[i].TotalCharges

			charges, err := strconv.ParseFloat(dpcharges.Charges, 64)
			if err != nil {
				loggerconfig.Error("SendEmailDPCharges, error in parsing charges string to float64 Error:", err, " reqId: ", reqH.RequestId)
				return apihelpers.SendInternalServerError()
			}

			gst, err := strconv.ParseFloat(dpcharges.Gst, 64)
			if err != nil {
				loggerconfig.Error("SendEmailDPCharges, error in parsing gst string to float64 Error:", err, " reqId: ", reqH.RequestId)
				return apihelpers.SendInternalServerError()
			}
			viewDPchargesRes.Charges += charges
			viewDPchargesRes.GST += gst

			viewDPchargesRes.DPChargesList = append(viewDPchargesRes.DPChargesList, dpcharges)
		}

		viewDPchargesRes.TotalCharges = viewDPchargesRes.Charges + viewDPchargesRes.GST
		viewDPchargesRes.UserDetails = profileData
		reportData, err := json.Marshal(viewDPchargesRes)
		if err != nil {
			loggerconfig.Error("SendEmailDPCharges, Error in marshalling viewDpChargesRes", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("SendEmailDPCharges Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}

	file, err := CreateExcelDPChargesReport(viewDPchargesRes)
	if err != nil {
		loggerconfig.Error("SendEmailDPCharges, Error in getting excel file, error: ", err, " reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	encodedString, err := helpers.EncodeExcelToBase64(file)
	if err != nil {
		loggerconfig.Error("SendEmailDPCharges, Error in creating base64 format from Excel, error:", err, " reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var sendEmailDpChargesDetails SendEmailDpCharges

	sendEmailDpChargesDetails.ClientId = profileData.ClientID
	sendEmailDpChargesDetails.ApplicantName = profileData.Name
	sendEmailDpChargesDetails.RecipientEmail = profileData.EmailID
	sendEmailDpChargesDetails.DateFrom = sendEmailDPChargesReq.DFDateFr
	sendEmailDpChargesDetails.DateTo = sendEmailDPChargesReq.DFDateTo
	sendEmailDpChargesDetails.EncodedReportFile = encodedString

	helpers.PublishMessage(constants.TopicExchange, constants.KeyDpChargesReport, sendEmailDpChargesDetails)

	var apiRes apihelpers.APIRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}

func (obj ReportsObj) ViewTradebook(tradebookReq models.TradebookReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {

	var scripWiseCosting models.ScripWiseCostingRes
	var err error

	fileName := constants.TradebookReport + strings.ToUpper(tradebookReq.UserID) + "_" + strings.ToUpper(tradebookReq.DFDateFr) + "_to_" + strings.ToUpper(tradebookReq.DFDateTo)
	storedReportData, _ := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" {
		err = json.Unmarshal([]byte(storedReportData), &scripWiseCosting)
		if err != nil {
			loggerconfig.Error("ViewTradebook, error in unmarshalling storedReportData : ", err, " for fileName fileName: ", fileName)
		}
	}
	if storedReportData == "" {

		theBackofficeProvider := v1.GetBackOfficeProvider()
		scripWiseCosting, err = theBackofficeProvider.GetScripWiseCostingData(tradebookReq, reqH)
		if err != nil {
			loggerconfig.Error("ViewTradebook, there is some error in fetching data from Shilpi api function: Error : ", err, " reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		scripWiseCosting.UserDetails = profileData
		reportData, err := json.Marshal(scripWiseCosting)
		if err != nil {
			loggerconfig.Error("ViewTradebook, Error in marshalling viewDpChargesRes", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("ViewTradebook, Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}

	var apiRes apihelpers.APIRes
	apiRes.Data = scripWiseCosting
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

// makeExcelScripWiseData creates an Excel file with the ScripWiseData
func makeExcelScripWiseData(data models.ScripWiseCostingRes) (*xlsx.File, error) {
	// Create a new workbook
	file := xlsx.NewFile()

	// Create a new sheet
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		loggerconfig.Error("Error creating sheet:", err)
		return nil, err
	}

	// Write summary part to Excel
	writeSummary(sheet, 1, 1, data)

	// Write ScripWiseCosting part to Excel
	writeScripWiseCosting(sheet, 5, 1, data)

	return file, nil
}

// writeSummary writes the summary part to the Excel sheet
func writeSummary(sheet *xlsx.Sheet, row, col int, data models.ScripWiseCostingRes) {
	summaryHeaders := []string{"TotalBrokerage", "TotalGST", "TotalSEBITax", "TotalSTT", "TotalTurnCharges", "TotalStampDuty", "TotalOtherCharges", "TotalCharges"}

	// Write header row
	headerRow := sheet.AddRow()
	for _, header := range summaryHeaders {
		cell := headerRow.AddCell()
		cell.SetString(header)
	}

	// Write data row
	dataRow := sheet.AddRow()
	summaryData := []float64{data.TotalBrokerage, data.TotalGST, data.TotalSEBITax, data.TotalSTT, data.TotalTurnCharges, data.TotalStampDuty, data.TotalOtherCharges, data.TotalCharges}
	for _, value := range summaryData {
		cell := dataRow.AddCell()
		cell.SetFloat(value)
	}
}

// writeScripWiseCosting writes the ScripWiseCosting part to the Excel sheet
func writeScripWiseCosting(sheet *xlsx.Sheet, row, col int, data models.ScripWiseCostingRes) {
	scripHeaders := []string{"Brokerage", "OrderNo", "OtherCharges", "GST", "Stamp", "Price", "TrxDate", "ISINCode", "SEBIFee", "STT", "TurnTax", "Exchange", "BrokType", "ScripName", "BuySellType", "Quantity"}

	// Write header row
	headerRow := sheet.AddRow()
	for _, header := range scripHeaders {
		cell := headerRow.AddCell()
		cell.SetString(header)
	}

	// Write data rows
	for _, scripData := range data.ScripWiseCosting {
		dataRow := sheet.AddRow()
		for _, value := range []interface{}{
			scripData.Brokerage, scripData.OrderNo, scripData.OtherCharges, scripData.GST, scripData.Stamp, scripData.Price, scripData.TrxDate, scripData.ISINCode, scripData.SEBIFee, scripData.STT, scripData.TurnTax, scripData.Exchange, scripData.BrokType, scripData.ScripName, scripData.BuySellType, scripData.Quantity,
		} {
			cell := dataRow.AddCell()
			cell.SetValue(value)
		}
	}
}

func (obj ReportsObj) DownloadTradebook(tradebookReq models.TradebookReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {
	var scripWiseCosting models.ScripWiseCostingRes
	var err error

	fileName := constants.TradebookReport + strings.ToUpper(tradebookReq.UserID) + "_" + strings.ToUpper(tradebookReq.DFDateFr) + "_to_" + strings.ToUpper(tradebookReq.DFDateTo)
	storedReportData, _ := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" {
		err = json.Unmarshal([]byte(storedReportData), &scripWiseCosting)
		if err != nil {
			loggerconfig.Error("DownloadTradebook, error in unmarshalling storedReportData : ", err, " for fileName fileName: ", fileName)
		}
	}
	if storedReportData == "" {

		theBackofficeProvider := v1.GetBackOfficeProvider()
		scripWiseCosting, err = theBackofficeProvider.GetScripWiseCostingData(tradebookReq, reqH)
		if err != nil {
			loggerconfig.Error("DownloadTradebook, there is some error in fetching data from Shilpi api function: Error : ", err, " reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		scripWiseCosting.UserDetails = profileData
		reportData, err := json.Marshal(scripWiseCosting)
		if err != nil {
			loggerconfig.Error("DownloadTradebook, Error in marshalling viewDpChargesRes", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("DownloadTradebook, Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}

	file, err := makeExcelScripWiseData(scripWiseCosting)
	if err != nil {
		loggerconfig.Error("DownloadTradebook, Error in getting excel file, error: ", err, " reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	s3FileName := fileName + ".xlsx"
	expiryHours := int64(24)

	// Generate a pre-signed URL for the XLSX file
	url, err := helpers.UploadFileToS3AndGetPresignedURL(constants.TradebookS3FolderName, s3FileName, file, expiryHours)
	if err != nil {
		loggerconfig.Error("DownloadTradebook, failed to generate pre-signed URL:", err, " reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var downloadTradebookRes models.DownloadTradebookRes
	downloadTradebookRes.DownloadUrl = url

	var apiRes apihelpers.APIRes
	apiRes.Data = downloadTradebookRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ReportsObj) ViewLedger(ledgerReq models.LedgerReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {

	var financialLedgerData models.FinancialLedgerRes
	var err error

	fileName := constants.LedgerReport + strings.ToUpper(ledgerReq.UserID) + "_" + strings.ToUpper(ledgerReq.DFDateFr) + "_to_" + strings.ToUpper(ledgerReq.DFDateTo)
	storedReportData, _ := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" {
		err = json.Unmarshal([]byte(storedReportData), &financialLedgerData)
		if err != nil {
			loggerconfig.Error("ViewLedger, error in unmarshalling storedReportData : ", err, " for fileName fileName: ", fileName)
		}
	}
	if storedReportData == "" {

		var getScripWiseCosting models.GetFinancialLedgerDataReq
		getScripWiseCosting.UserID = ledgerReq.UserID
		getScripWiseCosting.DFDateFr = ledgerReq.DFDateFr
		getScripWiseCosting.DFDateTo = ledgerReq.DFDateTo
		theBackofficeProvider := v1.GetBackOfficeProvider()
		financialLedgerData, err = theBackofficeProvider.GetFinancialLedgerData(getScripWiseCosting, reqH)
		if err != nil {
			loggerconfig.Error("ViewLedger, there is some error in fetching data from Shilpi api function: Error : ", err, " reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		financialLedgerData.UserDetails = profileData
		reportData, err := json.Marshal(financialLedgerData)
		if err != nil {
			loggerconfig.Error("ViewLedger, Error in marshalling viewDpChargesRes", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("ViewLedger, Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}

	var apiRes apihelpers.APIRes
	apiRes.Data = financialLedgerData
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

// makeExcelFinancialLedgerData creates an Excel file with the FinancialLedgerData
func makeExcelFinancialLedgerData(data models.FinancialLedgerRes) (*xlsx.File, error) {
	// Create a new workbook
	file := xlsx.NewFile()

	// Create a new sheet
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return nil, err
	}

	// Write summary part to Excel
	writeFinancialLedgerSummary(sheet, 1, 1, data)

	// Write FinancialLedgerData part to Excel
	writeFinancialLedgerData(sheet, 3, 1, data)

	return file, nil
}

// writeFinancialLedgerSummary writes the summary part to the Excel sheet
func writeFinancialLedgerSummary(sheet *xlsx.Sheet, row, col int, data models.FinancialLedgerRes) {
	summaryHeaders := []string{"OpeningBalance", "Inflow", "Outflow", "FundsReceived", "FundsWithdrawn", "ClosingBalance"}

	// Write header row
	headerRow := sheet.AddRow()
	for _, header := range summaryHeaders {
		cell := headerRow.AddCell()
		cell.SetString(header)
	}

	// Write data row
	dataRow := sheet.AddRow()
	summaryData := []float64{data.OpeningBalance, data.Inflow, data.Outflow, data.FundsReceived, data.FundsWithdrawn, data.ClosingBalance}
	for _, value := range summaryData {
		cell := dataRow.AddCell()
		cell.SetFloat(value)
	}
}

// writeFinancialLedgerData writes the FinancialLedgerData part to the Excel sheet
func writeFinancialLedgerData(sheet *xlsx.Sheet, row, col int, data models.FinancialLedgerRes) {
	financialLedgerHeaders := []string{"TransactionDate", "TransactionDetails", "Segment", "Exchange", "Debit", "Credit", "SettlementNumber", "NetBalance", "SettlementDate"}

	// Write header row
	headerRow := sheet.AddRow()
	for _, header := range financialLedgerHeaders {
		cell := headerRow.AddCell()
		cell.SetString(header)
	}

	// Write data rows
	for _, ledgerData := range data.FinancialLedger {
		dataRow := sheet.AddRow()
		for _, value := range []interface{}{
			ledgerData.TransactionDate, ledgerData.TransactionDetails, ledgerData.Segment, ledgerData.Exchange, ledgerData.Debit, ledgerData.Credit, ledgerData.SettlementNumber, ledgerData.NetBalance, ledgerData.SettlementDate,
		} {
			cell := dataRow.AddCell()
			cell.SetValue(value)
		}
	}
}

func (obj ReportsObj) DownloadLedger(ledgerReq models.LedgerReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {

	var financialLedgerData models.FinancialLedgerRes
	var err error

	fileName := constants.LedgerReport + strings.ToUpper(ledgerReq.UserID) + "_" + strings.ToUpper(ledgerReq.DFDateFr) + "_to_" + strings.ToUpper(ledgerReq.DFDateTo)
	storedReportData, _ := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" {
		err = json.Unmarshal([]byte(storedReportData), &financialLedgerData)
		if err != nil {
			loggerconfig.Error("DownloadLedger, error in unmarshalling storedReportData : ", err, " for fileName fileName: ", fileName)
		}
	}
	if storedReportData == "" {

		var getScripWiseCosting models.GetFinancialLedgerDataReq
		getScripWiseCosting.UserID = ledgerReq.UserID
		getScripWiseCosting.DFDateFr = ledgerReq.DFDateFr
		getScripWiseCosting.DFDateTo = ledgerReq.DFDateTo
		theBackofficeProvider := v1.GetBackOfficeProvider()
		financialLedgerData, err = theBackofficeProvider.GetFinancialLedgerData(getScripWiseCosting, reqH)
		if err != nil {
			loggerconfig.Error("DownloadLedger, there is some error in fetching data from Shilpi api function: Error : ", err, " reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		financialLedgerData.UserDetails = profileData
		reportData, err := json.Marshal(financialLedgerData)
		if err != nil {
			loggerconfig.Error("DownloadLedger, Error in marshalling viewDpChargesRes", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("DownloadLedger, Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}

	file, err := makeExcelFinancialLedgerData(financialLedgerData)
	if err != nil {
		loggerconfig.Error("DownloadLedger, Error in getting excel file, error: ", err, " reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	s3FileName := fileName + ".xlsx"
	expiryHours := int64(24)
	// Generate a pre-signed URL for the XLSX file
	url, err := helpers.UploadFileToS3AndGetPresignedURL(constants.LedgerS3FolderName, s3FileName, file, expiryHours)
	if err != nil {
		loggerconfig.Error("DownloadLedger, failed to generate pre-signed URL:", err, " reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var downloadLedgerRes models.DownloadLedgerRes
	downloadLedgerRes.DownloadUrl = url

	var apiRes apihelpers.APIRes
	apiRes.Data = downloadLedgerRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func makeExcelOpenPositionData(fileName string, data models.OpenPositionRes) (*xlsx.File, error) {
	// Create a new Excel file
	file := xlsx.NewFile()

	// Create a new sheet
	sheet, err := file.AddSheet("OpenPositions")
	if err != nil {
		return nil, err
	}

	// Add summary section
	addSummarySection(sheet, "Summary", data)

	// Add section for Equity Derivative
	addDerivativeSection(sheet, "Equity", data.EquityDerivative)

	// Add section for Currency Derivative
	addDerivativeSection(sheet, "Currency", data.CurrencyDerivative)

	// Add section for Commodity Derivative
	addDerivativeSection(sheet, "Commodity", data.CommodityDerivative)

	loggerconfig.Info("Excel file", fileName, " created successfully")
	return file, err
}

func addSummarySection(sheet *xlsx.Sheet, heading string, data models.OpenPositionRes) {
	// Add heading for the section
	headingRow := sheet.AddRow()
	headingCell := headingRow.AddCell()
	headingCell.SetString(heading)
	headingCell.Merge(0, 1) // Merge cells for better formatting

	// Add column headers
	headerRow := sheet.AddRow()
	headerRow.AddCell().SetString("Field")
	headerRow.AddCell().SetString("Value")

	// Add data for each field in the summary section
	addSummaryRow(sheet, "EquityFutureMTM", data.EquityFutureMTM)
	addSummaryRow(sheet, "CurrencyFutureMTM", data.CurrencyFutureMTM)
	addSummaryRow(sheet, "CommodityFutureMTM", data.CommodityFutureMTM)
	addSummaryRow(sheet, "EquityOptionMTM", data.EquityOptionMTM)
	addSummaryRow(sheet, "CurrencyOptionMTM", data.CurrencyOptionMTM)
	addSummaryRow(sheet, "CommodityOptionMTM", data.CommodityOptionMTM)
}

func addSummaryRow(sheet *xlsx.Sheet, fieldName string, value float64) {
	row := sheet.AddRow()
	row.AddCell().SetString(fieldName)
	row.AddCell().SetFloat(value)
}

func addDerivativeSection(sheet *xlsx.Sheet, heading string, positions []models.OpenPositionData) {
	// Add heading for the section
	headingRow := sheet.AddRow()
	headingCell := headingRow.AddCell()
	headingCell.SetString("Open Position for " + heading + " Derivative")
	headingCell.Merge(9, 0) // Merge cells for better formatting

	// Add column headers
	headerRow := sheet.AddRow()
	headerRow.AddCell().SetString("InstrumentType")
	headerRow.AddCell().SetString("OptionType")
	headerRow.AddCell().SetString("BuySellType")
	headerRow.AddCell().SetString("StrikePrice")
	headerRow.AddCell().SetString("ExpiryDate")
	headerRow.AddCell().SetString("Exchange")
	headerRow.AddCell().SetString("OpenQuantity")
	headerRow.AddCell().SetString("AveragePrice")
	headerRow.AddCell().SetString("ClosingPrice")
	headerRow.AddCell().SetString("UnrealisedProfitOrLoss")

	// Add data for each position in the section
	for _, position := range positions {
		// Add data for each position in the section
		row := sheet.AddRow()
		row.AddCell().SetString(position.InstrumentType)
		row.AddCell().SetString(position.OptionType)
		row.AddCell().SetString(position.BuySellType)
		row.AddCell().SetFloat(position.StrikePrice)
		row.AddCell().SetString(position.ExpiryDate)
		row.AddCell().SetString(position.Exchange)
		row.AddCell().SetFloat(position.OpenQuantity)
		row.AddCell().SetFloat(position.AveragePrice)
		row.AddCell().SetFloat(position.ClosingPrice)
		row.AddCell().SetFloat(position.UnrealisedProfitOrLoss)
	}
}

func (obj ReportsObj) ViewOpenPosition(openPositionReq models.OpenPositionReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {

	var OpenPositionData models.OpenPositionRes
	var err error

	fileName := constants.OpenPositionReport + strings.ToUpper(openPositionReq.UserID) + "_" + helpers.GetCurrentTimeInIST().AddDate(0, 0, -1).Format(constants.ShilpiDateFormat)
	storedReportData, _ := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" {
		err = json.Unmarshal([]byte(storedReportData), &OpenPositionData)
		if err != nil {
			loggerconfig.Error("ViewOpenPosition, error in unmarshalling storedReportData : ", err, " for fileName: ", fileName)
		}
	}
	if storedReportData == "" {

		theBackofficeProvider := v1.GetBackOfficeProvider()
		OpenPositionData, err = theBackofficeProvider.GetOpenPositionData(openPositionReq, reqH)
		if err != nil {
			loggerconfig.Error("ViewOpenPosition, there is some error in fetching data from Shilpi api function: Error : ", err, " reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		OpenPositionData.UserDetails = profileData
		reportData, err := json.Marshal(OpenPositionData)
		if err != nil {
			loggerconfig.Error("ViewOpenPosition, Error in marshalling viewDpChargesRes", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("ViewOpenPosition, Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}

	var apiRes apihelpers.APIRes
	apiRes.Data = OpenPositionData
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ReportsObj) DownloadOpenPosition(openPositionReq models.OpenPositionReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {

	var OpenPositionData models.OpenPositionRes
	var err error

	fileName := constants.OpenPositionReport + strings.ToUpper(openPositionReq.UserID) + "_" + helpers.GetCurrentTimeInIST().AddDate(0, 0, -1).Format(constants.ShilpiDateFormat)
	storedReportData, _ := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" {
		err = json.Unmarshal([]byte(storedReportData), &OpenPositionData)
		if err != nil {
			loggerconfig.Error("DownloadOpenPosition, error in unmarshalling storedReportData : ", err, " for fileName: ", fileName)
		}
	}
	if storedReportData == "" {

		theBackofficeProvider := v1.GetBackOfficeProvider()
		OpenPositionData, err = theBackofficeProvider.GetOpenPositionData(openPositionReq, reqH)
		if err != nil {
			loggerconfig.Error("DownloadOpenPosition there is some error in fetching data from Shilpi api function: Error : ", err, " reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		OpenPositionData.UserDetails = profileData
		reportData, err := json.Marshal(OpenPositionData)
		if err != nil {
			loggerconfig.Error("DownloadOpenPosition, Error in marshalling viewDpChargesRes", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("DownloadOpenPosition, Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}

	file, err := makeExcelOpenPositionData("output.xlsx", OpenPositionData)
	if err != nil {
		loggerconfig.Error("DownloadOpenPosition, Error creating Excel file:", err, " reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	s3FileName := fileName + ".xlsx"
	expiryHours := int64(24)

	// Generate a pre-signed URL for the XLSX file
	url, err := helpers.UploadFileToS3AndGetPresignedURL(constants.OpenPositionS3FolderName, s3FileName, file, expiryHours)
	if err != nil {
		loggerconfig.Error("DownloadOpenPosition, failed to generate pre-signed URL:", err, " reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	var downloadOpenPositionRes models.DownloadOpenPositionRes
	downloadOpenPositionRes.DownloadUrl = url

	var apiRes apihelpers.APIRes
	apiRes.Data = downloadOpenPositionRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func createFONetPositionExcel(fileName string, data models.FONetPositionRes) (*xlsx.File, error) {
	// Create a new Excel file
	file := xlsx.NewFile()

	// Create a new sheet
	sheet, err := file.AddSheet("FONetPosition")
	if err != nil {
		return nil, err
	}

	// Add PnL Summary section
	addPnLSummarySection(sheet, "Summary", data.Summary)

	// Add Charges Details section
	addChargesDetailsSection(sheet, "Charges Details", data.ChargesDetails)

	// Add FONetPositionDetails section
	addFONetPositionDetailsSection(sheet, "FONetPositionDetails", data.FONetPositionDetails)

	loggerconfig.Info("Excel file ", fileName, " created successfully\n", fileName)
	return file, err
}

func addPnLSummarySection(sheet *xlsx.Sheet, heading string, data models.FONetPositionSummaryData) {
	// Add heading for the section
	headingRow := sheet.AddRow()
	headingCell := headingRow.AddCell()
	headingCell.SetString(heading)
	headingCell.Merge(5, 0) // Merge cells for better formatting
	headingCell.GetStyle().Font.Bold = true

	// Add data for PnL Summary
	addPnLSummaryRow(sheet, "Date Range", data.DateRange)
	addPnLSummaryRow(sheet, "Charges", data.Charges)
	addPnLSummaryRow(sheet, "Realised PNL", data.RealisedPNL)
	addPnLSummaryRow(sheet, "Unrealised PNL", data.UnRealisedPNL)
	addPnLSummaryRow(sheet, "Net PNL", data.NetPNL)

	sheet.AddRow()
	sheet.AddRow()
}

func addPnLSummaryRow(sheet *xlsx.Sheet, fieldName string, value interface{}) {
	// Add data for PnL Summary
	row := sheet.AddRow()
	row.AddCell().SetString(fieldName)
	row.AddCell().SetString(fmt.Sprintf("%v", value))
}

func addChargesDetailsSection(sheet *xlsx.Sheet, heading string, data models.FONetPositionChargesData) {
	// Add heading for the section
	headingRow := sheet.AddRow()
	headingCell := headingRow.AddCell()
	headingCell.SetString(heading)
	headingCell.Merge(5, 0) // Merge cells for better formatting
	headingCell.GetStyle().Font.Bold = true

	// Add data for Charges Details
	addChargesDetailsRow(sheet, "Brokerage", data.Brockerage)
	addChargesDetailsRow(sheet, "Exchange Transaction Charges", data.ExchangeTransactionCharges)
	addChargesDetailsRow(sheet, "Clearing Charges", data.ClearingCharges)
	addChargesDetailsRow(sheet, "Integrated GST", data.IntegratedGST)
	addChargesDetailsRow(sheet, "Securities Transaction Tax", data.SecuritiesTransactionTax)
	addChargesDetailsRow(sheet, "SEBI Fees", data.SEBIFees)
	addChargesDetailsRow(sheet, "Stamp Duty", data.StampDuty)
	addChargesDetailsRow(sheet, "Total Charges", data.TotalCharges)

	sheet.AddRow()
	sheet.AddRow()
}

func addChargesDetailsRow(sheet *xlsx.Sheet, fieldName string, value interface{}) {
	// Add data for Charges Details
	row := sheet.AddRow()
	row.AddCell().SetString(fieldName)
	row.AddCell().SetString(fmt.Sprintf("%v", value))
}

func addFONetPositionDetailsSection(sheet *xlsx.Sheet, heading string, data []models.FONetPositionDetailsData) {
	// Add heading for the section
	headingRow := sheet.AddRow()
	headingCell := headingRow.AddCell()
	headingCell.SetString(heading)
	headingCell.Merge(14, 0) // Merge cells for better formatting
	headingCell.GetStyle().Font.Bold = true

	// Add column headers
	headerRow := sheet.AddRow()
	headerRow.AddCell().SetString("Symbol")
	headerRow.AddCell().SetString("Instrument Type")
	headerRow.AddCell().SetString("Option Type")
	headerRow.AddCell().SetString("Strike Price")
	headerRow.AddCell().SetString("Expiry Date")
	headerRow.AddCell().SetString("Quantity")
	headerRow.AddCell().SetString("Buy Value")
	headerRow.AddCell().SetString("Buy Price")
	headerRow.AddCell().SetString("Sell Value")
	headerRow.AddCell().SetString("Sell Price")
	headerRow.AddCell().SetString("Realized PNL")
	headerRow.AddCell().SetString("Previous Closing Price")
	headerRow.AddCell().SetString("Open Quantity")
	headerRow.AddCell().SetString("Open Value")
	headerRow.AddCell().SetString("Unrealized PNL")

	// Add data for FONetPositionDetails
	for _, position := range data {
		addFONetPositionDetailsRow(sheet, position)
	}
}

func addFONetPositionDetailsRow(sheet *xlsx.Sheet, data models.FONetPositionDetailsData) {
	// Add data for FONetPositionDetails
	row := sheet.AddRow()
	row.AddCell().SetString(data.Symbol)
	row.AddCell().SetString(data.InstrumentType)
	row.AddCell().SetString(data.OptionType)
	row.AddCell().SetString(fmt.Sprintf("%v", data.StrikePrice))
	row.AddCell().SetString(data.ExpiryDate)
	row.AddCell().SetString(fmt.Sprintf("%v", data.Quantity))
	row.AddCell().SetString(fmt.Sprintf("%v", data.BuyValue))
	row.AddCell().SetString(fmt.Sprintf("%v", data.BuyPrice))
	row.AddCell().SetString(fmt.Sprintf("%v", data.SellValue))
	row.AddCell().SetString(fmt.Sprintf("%v", data.SellPrice))
	row.AddCell().SetString(fmt.Sprintf("%v", data.RealizedPNL))
	row.AddCell().SetString(fmt.Sprintf("%v", data.PreviousClosingPrice))
	row.AddCell().SetString(fmt.Sprintf("%v", data.OpenQuantity))
	row.AddCell().SetString(fmt.Sprintf("%v", data.OpenValue))
	row.AddCell().SetString(fmt.Sprintf("%v", data.UnrealizedPNL))
}

func (obj ReportsObj) ViewFnoPnl(fnoPnlReq models.FnoPnlReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {

	var fONetPositionData models.FONetPositionRes
	var err error

	fileName := constants.FnoPnlReport + strings.ToUpper(fnoPnlReq.UserID) + "_" + strings.ToUpper(fnoPnlReq.DFDateFr) + "_to_" + strings.ToUpper(fnoPnlReq.DFDateTo)
	storedReportData, _ := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" {
		err = json.Unmarshal([]byte(storedReportData), &fONetPositionData)
		if err != nil {
			loggerconfig.Error("ViewDPCharges, error in unmarshalling storedReportData : ", err, " for fileName: ", fileName)
		}
	}
	if storedReportData == "" {

		var getFONetPositionDataReq models.GetFONetPositionDataReq
		getFONetPositionDataReq.UserID = fnoPnlReq.UserID
		getFONetPositionDataReq.DFDateFr = fnoPnlReq.DFDateFr
		getFONetPositionDataReq.DFDateTo = fnoPnlReq.DFDateTo
		theBackofficeProvider := v1.GetBackOfficeProvider()
		fONetPositionData, err = theBackofficeProvider.GetFONetPositionData(getFONetPositionDataReq, reqH)
		if err != nil {
			loggerconfig.Error("ViewFnoPnl there is some error in fetching data from Shilpi api function: Error : ", err, " reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		fONetPositionData.UserDetails = profileData
		reportData, err := json.Marshal(fONetPositionData)
		if err != nil {
			loggerconfig.Error("Error in marshalling viewDpChargesRes", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("DownloadTradebook Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}
	var apiRes apihelpers.APIRes
	apiRes.Data = fONetPositionData
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ReportsObj) DownloadFnoPnl(fnoPnlReq models.FnoPnlReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {

	var fONetPositionData models.FONetPositionRes
	var err error
	fileName := constants.FnoPnlReport + strings.ToUpper(fnoPnlReq.UserID) + "_" + strings.ToUpper(fnoPnlReq.DFDateFr) + "_to_" + strings.ToUpper(fnoPnlReq.DFDateTo)
	storedReportData, _ := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" {
		err = json.Unmarshal([]byte(storedReportData), &fONetPositionData)
		if err != nil {
			loggerconfig.Error("DownloadFnoPnl, error in unmarshalling storedReportData : ", err, " for fileName: ", fileName)
		}
	}
	if storedReportData == "" {

		var getFONetPositionDataReq models.GetFONetPositionDataReq
		getFONetPositionDataReq.UserID = fnoPnlReq.UserID
		getFONetPositionDataReq.DFDateFr = fnoPnlReq.DFDateFr
		getFONetPositionDataReq.DFDateTo = fnoPnlReq.DFDateTo
		theBackofficeProvider := v1.GetBackOfficeProvider()
		fONetPositionData, err = theBackofficeProvider.GetFONetPositionData(getFONetPositionDataReq, reqH)
		if err != nil {
			loggerconfig.Error("DownloadFnoPnl, there is some error in fetching data from Shilpi api function: Error : ", err, " reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		fONetPositionData.UserDetails = profileData
		reportData, err := json.Marshal(fONetPositionData)
		if err != nil {
			loggerconfig.Error("DownloadFnoPnl, Error in marshalling viewDpChargesRes", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("DownloadFnoPnl, Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}

	file, err := createFONetPositionExcel("output.xlsx", fONetPositionData)
	if err != nil {
		loggerconfig.Error("DownloadFnoPnl, Error creating Excel file:", err, " reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	s3FileName := fileName + ".xlsx"
	expiryHours := int64(24)

	// Generate a pre-signed URL for the XLSX file
	url, err := helpers.UploadFileToS3AndGetPresignedURL(constants.FnoPnlS3FolderName, s3FileName, file, expiryHours)
	if err != nil {
		loggerconfig.Error("DownloadFnoPnl, failed to generate pre-signed URL:", err, " reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	var downloadFnoPnlRes models.DownloadFnoPnlRes
	downloadFnoPnlRes.DownloadUrl = url

	var apiRes apihelpers.APIRes
	apiRes.Data = downloadFnoPnlRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ReportsObj) ViewHoldingFinancial(holdingFinancialDataReq models.GetHoldingFinancialDataReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {

	var holdingFinancialData models.GetHoldingFinancialDataRes
	var err error

	fileName := constants.HoldingFinancialReport + strings.ToUpper(holdingFinancialDataReq.UserID) + "_" + helpers.GetCurrentTimeInIST().AddDate(0, 0, -1).Format(constants.ShilpiDateFormat)
	storedReportData, _ := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" {
		err = json.Unmarshal([]byte(storedReportData), &holdingFinancialData)
		if err != nil {
			loggerconfig.Error("ViewHoldingFinancial, error in unmarshalling storedReportData : ", err, " for fileName: ", fileName)
		}
	}
	if storedReportData == "" {

		theBackofficeProvider := v1.GetBackOfficeProvider()
		holdingFinancialData, err = theBackofficeProvider.GetHoldingFinancialData(holdingFinancialDataReq, reqH)
		if err != nil {
			loggerconfig.Error("ViewHoldingFinancial, there is some error in fetching data from Shilpi api function: Error : ", err, " reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		holdingFinancialData.UserDetails = profileData
		reportData, err := json.Marshal(holdingFinancialData)
		if err != nil {
			loggerconfig.Error("ViewHoldingFinancial, Error in marshalling holdingFinancialData ", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("ViewHoldingFinancial, Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}

	var apiRes apihelpers.APIRes
	apiRes.Data = holdingFinancialData
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

// CreateExcelHoldingReport generates the Excel file for holding financial data
func CreateExcelHoldingReport(data models.GetHoldingFinancialDataRes) (*xlsx.File, error) {
	// Create a new workbook
	file := xlsx.NewFile()

	// Create a new sheet
	sheet, err := file.AddSheet("HoldingReport")
	if err != nil {
		return nil, err
	}

	// Add the title row for the report
	addTitleRowHolding(sheet)

	// Add user details below the title
	writeUserDetails(sheet, data.UserDetails, constants.HoldingReportName)

	// Add an empty row for spacing between sections
	sheet.AddRow()

	// Write the summary of holdings
	writeHoldingSummary(sheet, data.HoldingSummary)

	// Add space between sections
	sheet.AddRow()

	// Write the holdings data
	writeHoldingFinancialData(sheet, data.HoldingFinancialData)

	return file, nil
}

// addTitleRowHolding adds the main title to the top of the Excel sheet for Holding Report
func addTitleRowHolding(sheet *xlsx.Sheet) {
	// Create the title row
	titleRow := sheet.AddRow()
	titleRow.SetHeight(20) // Optionally set height for emphasis
	titleCell := titleRow.AddCell()
	titleCell.Merge(3, 0) // Merge first 4 cells for the title
	titleCell.SetString("Holding Report")
	titleCell.GetStyle().Font.Bold = true
	titleCell.GetStyle().Alignment.Horizontal = "center"

	// Add an empty row after the title for spacing
	sheet.AddRow()
}

// writeHoldingSummary writes the summary of the holding financial data in the Excel file
func writeHoldingSummary(sheet *xlsx.Sheet, summary models.HoldingSummaryData) {
	// Create headers for the summary
	summaryHeaders := []string{"Holding Summary", "Amount"}
	headerRow := sheet.AddRow()
	for _, header := range summaryHeaders {
		headerRow.AddCell().SetString(header)
	}

	// Define the order and labels of the summary data
	summaryData := []struct {
		Label string
		Value float64
	}{
		{"Invested Value", summary.InvestedValue},
		{"Current Value", summary.CurrentValue},
		{"Unrealised PNL", summary.UnrealisedPNL},
		{"Total Pledge Value", summary.TotalPledgeValue},
		{"Total Margin Value After Haircut", summary.TotalMarginValueAfterHaircut},
	}

	// Populate the summary data
	for _, item := range summaryData {
		row := sheet.AddRow()
		row.AddCell().SetString(item.Label)
		row.AddCell().SetString(strconv.FormatFloat(item.Value, 'f', 2, 64))
	}
}

// writeHoldingFinancialData writes the detailed holdings data in the Excel file
func writeHoldingFinancialData(sheet *xlsx.Sheet, holdings []models.GetHoldingFinancialData) {
	// Create headers for the holding data
	holdingHeaders := []string{
		"ISIN", "Instrument", "Pledged Qty", "Free Qty", "Total Qty", "Total Pledged Value",
		"Haircut Percentage", "Margin Available After Haircut", "Avg Buy Price", "Closing Price",
		"Investment Value", "Current Value", "Contribution Percentage", "Unrealized P&L", "Net Change",
	}
	headerRow := sheet.AddRow()
	for _, header := range holdingHeaders {
		headerRow.AddCell().SetString(header)
	}

	// Write each holding as a row
	for _, holding := range holdings {
		dataRow := sheet.AddRow()
		dataRow.AddCell().SetString(holding.Isin)
		dataRow.AddCell().SetString(holding.Instrument)
		dataRow.AddCell().SetString(strconv.FormatFloat(holding.PledgedQty, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(holding.FreeQty, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(holding.TotalQty, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(holding.TotalPledgedValue, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(holding.HaircutPercentage, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(holding.MarginAvailableAfterHaircut, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(holding.AvgBuyPrice, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(holding.ClosingPrice, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(holding.InvestmentValue, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(holding.CurrentValue, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(holding.ContributionPercentage, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(holding.UnrealizedProfitLoss, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(holding.NetChange, 'f', 2, 64))
	}
}

func (obj ReportsObj) DownloadHoldingFinancial(holdingFinancialDataReq models.GetHoldingFinancialDataReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {

	var holdingFinancialData models.GetHoldingFinancialDataRes
	var err error

	fileName := constants.HoldingFinancialReport + strings.ToUpper(holdingFinancialDataReq.UserID) + "_" + helpers.GetCurrentTimeInIST().AddDate(0, 0, -1).Format(constants.ShilpiDateFormat)
	storedReportData, _ := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" {
		err = json.Unmarshal([]byte(storedReportData), &holdingFinancialData)
		if err != nil {
			loggerconfig.Error("DownloadHoldingFinancial, error in unmarshalling storedReportData : ", err, " for fileName: ", fileName)
		}
	}
	if storedReportData == "" {

		theBackofficeProvider := v1.GetBackOfficeProvider()
		holdingFinancialData, err = theBackofficeProvider.GetHoldingFinancialData(holdingFinancialDataReq, reqH)
		if err != nil {
			loggerconfig.Error("DownloadHoldingFinancial, there is some error in fetching data from Shilpi api function: Error : ", err, " reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		holdingFinancialData.UserDetails = profileData
		reportData, err := json.Marshal(holdingFinancialData)
		if err != nil {
			loggerconfig.Error("DownloadHoldingFinancial, Error in marshalling viewDpChargesRes", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("DownloadHoldingFinancial, Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}

	file, err := CreateExcelHoldingReport(holdingFinancialData)
	if err != nil {
		loggerconfig.Error("DownloadHoldingFinancial, Error creating Excel file: ", err, " reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	s3FileName := fileName + ".xlsx"
	expiryHours := int64(24)

	// Generate a pre-signed URL for the XLSX file
	url, err := helpers.UploadFileToS3AndGetPresignedURL(constants.HoldingFinancialS3FolderName, s3FileName, file, expiryHours)
	if err != nil {
		loggerconfig.Error("DownloadOpenPosition, failed to generate pre-signed URL:", err, " reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}
	var downloadHoldingFinancialRes models.DownloadHoldingFinancialRes
	downloadHoldingFinancialRes.DownloadUrl = url

	var apiRes apihelpers.APIRes
	apiRes.Data = downloadHoldingFinancialRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ReportsObj) SendEmailHoldingFinancial(holdingFinancialDataReq models.GetHoldingFinancialDataReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {

	var holdingFinancialData models.GetHoldingFinancialDataRes
	var err error

	fileName := constants.HoldingFinancialReport + strings.ToUpper(holdingFinancialDataReq.UserID) + "_" + helpers.GetCurrentTimeInIST().AddDate(0, 0, -1).Format(constants.ShilpiDateFormat)
	storedReportData, _ := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" {
		err = json.Unmarshal([]byte(storedReportData), &holdingFinancialData)
		if err != nil {
			loggerconfig.Error("DownloadHoldingFinancial, error in unmarshalling storedReportData : ", err, " for fileName: ", fileName)
		}
	}
	if storedReportData == "" {

		theBackofficeProvider := v1.GetBackOfficeProvider()
		holdingFinancialData, err = theBackofficeProvider.GetHoldingFinancialData(holdingFinancialDataReq, reqH)
		if err != nil {
			loggerconfig.Error("DownloadHoldingFinancial, there is some error in fetching data from Shilpi api function: Error : ", err, " reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		holdingFinancialData.UserDetails = profileData
		reportData, err := json.Marshal(holdingFinancialData)
		if err != nil {
			loggerconfig.Error("DownloadHoldingFinancial, Error in marshalling viewDpChargesRes", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("DownloadHoldingFinancial, Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}

	file, err := CreateExcelHoldingReport(holdingFinancialData)
	if err != nil {
		loggerconfig.Error("DownloadHoldingFinancial, Error creating Excel file: ", err, " reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	encodedString, err := helpers.EncodeExcelToBase64(file)
	if err != nil {
		loggerconfig.Error("SendEmailLedger, Error in creating base64 format from Excel, error:", err, " reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var sendEmailHoldingFinancial SendEmailHoldingFinancial

	sendEmailHoldingFinancial.ClientId = profileData.ClientID
	sendEmailHoldingFinancial.ApplicantName = profileData.Name
	sendEmailHoldingFinancial.RecipientEmail = profileData.EmailID
	sendEmailHoldingFinancial.EncodedReportFile = encodedString

	helpers.PublishMessage(constants.TopicExchange, constants.KeyHoldingFinancialReport, sendEmailHoldingFinancial)

	var apiRes apihelpers.APIRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes

}

func (obj ReportsObj) SendEmailLedger(ledgerReq models.LedgerReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {

	var financialLedgerData models.FinancialLedgerRes
	var err error

	fileName := constants.LedgerReport + strings.ToUpper(ledgerReq.UserID) + "_" + strings.ToUpper(ledgerReq.DFDateFr) + "_to_" + strings.ToUpper(ledgerReq.DFDateTo)
	storedReportData, _ := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" {
		err = json.Unmarshal([]byte(storedReportData), &financialLedgerData)
		if err != nil {
			loggerconfig.Error("SendEmailLedger, error in unmarshalling storedReportData : ", err, " for fileName fileName: ", fileName)
		}
	}
	if storedReportData == "" {

		var getScripWiseCosting models.GetFinancialLedgerDataReq
		getScripWiseCosting.UserID = ledgerReq.UserID
		getScripWiseCosting.DFDateFr = ledgerReq.DFDateFr
		getScripWiseCosting.DFDateTo = ledgerReq.DFDateTo
		theBackofficeProvider := v1.GetBackOfficeProvider()
		financialLedgerData, err = theBackofficeProvider.GetFinancialLedgerData(getScripWiseCosting, reqH)
		if err != nil {
			loggerconfig.Error("SendEmailLedger, there is some error in fetching data from Shilpi api function: Error : ", err, " reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}

		financialLedgerData.UserDetails = profileData
		reportData, err := json.Marshal(financialLedgerData)
		if err != nil {
			loggerconfig.Error("SendEmailLedger, Error in marshalling viewDpChargesRes", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("SendEmailLedger, Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}

	file, err := makeExcelFinancialLedgerData(financialLedgerData)
	if err != nil {
		loggerconfig.Error("SendEmailLedger, Error in getting excel file, error: ", err, " reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	encodedString, err := helpers.EncodeExcelToBase64(file)
	if err != nil {
		loggerconfig.Error("SendEmailLedger, Error in creating base64 format from Excel, error:", err, " reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var sendEmailLedgerDetails SendEmailLedger

	sendEmailLedgerDetails.ClientId = profileData.ClientID
	sendEmailLedgerDetails.ApplicantName = profileData.Name
	sendEmailLedgerDetails.RecipientEmail = profileData.EmailID
	sendEmailLedgerDetails.DateFrom = ledgerReq.DFDateFr
	sendEmailLedgerDetails.DateTo = ledgerReq.DFDateTo
	sendEmailLedgerDetails.EncodedReportFile = encodedString

	helpers.PublishMessage(constants.TopicExchange, constants.KeyLedgerReport, sendEmailLedgerDetails)

	var apiRes apihelpers.APIRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

// ViewCommodityTradebook(CommodityTradebookReq, ReqHeader, ProfileDataResp) (int, apihelpers.APIRes)
func (obj ReportsObj) ViewCommodityTradebook(commodityTradebookReq models.CommodityTradebookReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {

	var commodityTransaction models.CommodityTransactionRes
	var err error

	fileName := constants.CommodityTradebookReport + strings.ToUpper(commodityTradebookReq.UserID) + "_" + strings.ToUpper(commodityTradebookReq.DFDateFr) + "_to_" + strings.ToUpper(commodityTradebookReq.DFDateTo)
	storedReportData, _ := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" {

		err = json.Unmarshal([]byte(storedReportData), &commodityTransaction)
		if err != nil {
			loggerconfig.Error("ViewCommodityTradebook, error in unmarshalling storedReportData : ", err, " for fileName fileName: ", fileName)
		}
	}
	if storedReportData == "" {
		theBackofficeProvider := v1.GetBackOfficeProvider()
		commodityTransaction, err = theBackofficeProvider.GetCommodityTransactionData(commodityTradebookReq, reqH)
		if err != nil {
			loggerconfig.Error("ViewCommodityTradebook, there is some error in fetching data from Shilpi api function: Error : ", err, " reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		commodityTransaction.UserDetails = profileData
		reportData, err := json.Marshal(commodityTransaction)
		if err != nil {
			loggerconfig.Error("ViewCommodityTradebook, Error in marshalling viewDpChargesRes", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("ViewCommodityTradebook, Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}

	var apiRes apihelpers.APIRes
	apiRes.Data = commodityTransaction
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

// CreateExcelCommodityTradebook generates the Excel file for commodity transactions and charges
func CreateExcelCommodityTradebook(data models.CommodityTransactionRes) (*xlsx.File, error) {
	// Create a new workbook
	file := xlsx.NewFile()

	// Create a new sheet
	sheet, err := file.AddSheet("TradebookCommodity")
	if err != nil {
		return nil, err
	}

	addTitleRow(sheet)

	writeUserDetails(sheet, data.UserDetails, constants.CommodityTradebookReportName)

	sheet.AddRow()

	writeCommodityChargesSummary(sheet, data)

	sheet.AddRow()

	writeCommodityTransactionsData(sheet, data.CommodityTransactions)

	return file, nil
}

func addTitleRow(sheet *xlsx.Sheet) {
	// Create the title row
	titleRow := sheet.AddRow()
	titleRow.SetHeight(20)
	titleCell := titleRow.AddCell()
	titleCell.Merge(3, 0)
	titleCell.SetString(constants.CommodityTradebookReportName)
	titleCell.GetStyle().Font.Bold = true
	titleCell.GetStyle().Alignment.Horizontal = "center"

	sheet.AddRow()
}

func writeUserDetails(sheet *xlsx.Sheet, userDetails models.ProfileDataResp, reportName string) {
	// Define the user detail labels and values
	userInfo := []struct {
		Label string
		Value string
	}{
		{"Client Name", userDetails.Name},
		{"PAN", userDetails.PanNumber},
		{"Client ID", userDetails.ClientID},
		{"Report", reportName},
	}

	// Write each user detail as a row
	for _, info := range userInfo {
		row := sheet.AddRow()
		row.AddCell().SetString(info.Label)
		row.AddCell().SetString(info.Value)
	}

	// Add an empty row after user details for spacing
	sheet.AddRow()
}

// writeCommodityChargesSummary writes the total charges in the Excel file
func writeCommodityChargesSummary(sheet *xlsx.Sheet, data models.CommodityTransactionRes) {
	// Create headers for the summary
	chargesHeaders := []string{"Trade Charges", "Amount"}
	headerRow := sheet.AddRow()
	for _, header := range chargesHeaders {
		headerRow.AddCell().SetString(header)
	}

	// Define the order and labels of the charges
	charges := []struct {
		Label string
		Value float64
	}{
		{"Brokerage", data.TotalBrokerage},
		{"GST", data.TotalGST},
		{"SEBI Tax", data.TotalSEBITax},
		{"STT (CTT)", data.TotalCTT},
		{"Exchange Turnover Charges", data.TotalTurnCharges},
		{"Stamp Duty", data.TotalStampDuty},
		{"Total Charges", data.TotalCharges},
	}

	// Populate the charge data
	for _, charge := range charges {
		row := sheet.AddRow()
		row.AddCell().SetString(charge.Label)
		row.AddCell().SetString(strconv.FormatFloat(charge.Value, 'f', 2, 64))
	}
}

// writeCommodityTransactionsData writes the commodity transactions in the Excel file
func writeCommodityTransactionsData(sheet *xlsx.Sheet, transactions []models.CommodityTransactionData) {
	// Create headers for the transaction data
	transactionHeaders := []string{
		"Symbol", "Instrument Type", "Expiry Date", "Option Type", "Strike Price", "Transaction Date",
		"Buy/Sell", "Trade Price", "Market Lot", "Qty", "Brokerage", "GST", "CTT", "SEBI Fees",
		"Turnover Fees", "Stamp Duty", "Clearing Charge", "Segment", "Exchange", "Order No", "Trade ID", "Trade Time",
	}
	headerRow := sheet.AddRow()
	for _, header := range transactionHeaders {
		headerRow.AddCell().SetString(header)
	}

	// Write each transaction as a row
	for _, trx := range transactions {
		dataRow := sheet.AddRow()
		dataRow.AddCell().SetString(trx.Symbol)
		dataRow.AddCell().SetString(trx.InstrumentType)
		dataRow.AddCell().SetString(trx.ExpiryDate)
		dataRow.AddCell().SetString(trx.OptionType)
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.StrikePrice, 'f', 2, 64))
		dataRow.AddCell().SetString(trx.TradeDate)
		dataRow.AddCell().SetString(trx.BuySellInd) // Already handled conversion
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.TradePrice, 'f', 2, 64))
		dataRow.AddCell().SetString(trx.MarketLot)
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.TradeQty, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.Brokerage, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.GST, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.CTT, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.SEBITax, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.TurnoverTax, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.StampDuty, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.CLGTax, 'f', 2, 64))
		dataRow.AddCell().SetString(trx.Segment)
		dataRow.AddCell().SetString(trx.Exchange)
		dataRow.AddCell().SetString(trx.OrderNo)
		dataRow.AddCell().SetString(trx.TradeNo)
		dataRow.AddCell().SetString(trx.TradeTime)
	}
}

func (obj ReportsObj) DownloadCommodityTradebook(commodityTradebookReq models.CommodityTradebookReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {
	var commodityTransaction models.CommodityTransactionRes
	var err error

	fileName := constants.CommodityTradebookReport + strings.ToUpper(commodityTradebookReq.UserID) + "_" + strings.ToUpper(commodityTradebookReq.DFDateFr) + "_to_" + strings.ToUpper(commodityTradebookReq.DFDateTo)
	storedReportData, _ := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" {

		err = json.Unmarshal([]byte(storedReportData), &commodityTransaction)
		if err != nil {
			loggerconfig.Error("DownloadCommodityTradebook, error in unmarshalling storedReportData : ", err, " for fileName fileName: ", fileName)
		}
	}
	if storedReportData == "" {

		theBackofficeProvider := v1.GetBackOfficeProvider()
		commodityTransaction, err = theBackofficeProvider.GetCommodityTransactionData(commodityTradebookReq, reqH)
		if err != nil {
			loggerconfig.Error("DownloadCommodityTradebook, there is some error in fetching data from Shilpi api function: Error : ", err, " reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		commodityTransaction.UserDetails = profileData
		reportData, err := json.Marshal(commodityTransaction)
		if err != nil {
			loggerconfig.Error("DownloadCommodityTradebook, Error in marshalling viewDpChargesRes", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("DownloadCommodityTradebook, Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}

	file, err := CreateExcelCommodityTradebook(commodityTransaction)
	if err != nil {
		loggerconfig.Error("DownloadCommodityTradebook, Error in getting excel file, error: ", err, " reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	s3FileName := fileName + ".xlsx"
	expiryHours := int64(24)

	// Generate a pre-signed URL for the XLSX file
	url, err := helpers.UploadFileToS3AndGetPresignedURL(constants.CommodityTradebookS3FolderName, s3FileName, file, expiryHours)
	if err != nil {
		loggerconfig.Error("DownloadCommodityTradebook, failed to generate pre-signed URL:", err, " reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var downloadCommodityTradebookRes models.DownloadCommodityTradebookRes
	downloadCommodityTradebookRes.DownloadUrl = url

	var apiRes apihelpers.APIRes
	apiRes.Data = downloadCommodityTradebookRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ReportsObj) SendEmailCommodityTradebook(commodityTradebookReq models.CommodityTradebookReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {

	var commodityTransaction models.CommodityTransactionRes
	var err error

	fileName := constants.CommodityTradebookReport + strings.ToUpper(commodityTradebookReq.UserID) + "_" + strings.ToUpper(commodityTradebookReq.DFDateFr) + "_to_" + strings.ToUpper(commodityTradebookReq.DFDateTo)
	storedReportData, _ := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" {

		err = json.Unmarshal([]byte(storedReportData), &commodityTransaction)
		if err != nil {
			loggerconfig.Error("SendEmailCommodityTradebook, error in unmarshalling storedReportData : ", err, " for fileName fileName: ", fileName)
		}
	}
	if storedReportData == "" {

		theBackofficeProvider := v1.GetBackOfficeProvider()
		commodityTransaction, err = theBackofficeProvider.GetCommodityTransactionData(commodityTradebookReq, reqH)
		if err != nil {
			loggerconfig.Error("SendEmailCommodityTradebook, there is some error in fetching data from Shilpi api function: Error : ", err, " reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		commodityTransaction.UserDetails = profileData
		reportData, err := json.Marshal(commodityTransaction)
		if err != nil {
			loggerconfig.Error("SendEmailCommodityTradebook, Error in marshalling viewDpChargesRes", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("SendEmailCommodityTradebook, Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}

	file, err := CreateExcelCommodityTradebook(commodityTransaction)
	if err != nil {
		loggerconfig.Error("SendEmailCommodityTradebook, Error in getting excel file, error: ", err, " reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	encodedString, err := helpers.EncodeExcelToBase64(file)
	if err != nil {
		loggerconfig.Error("SendEmailLedger, Error in creating base64 format from Excel, error:", err, " reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var sendEmailCommodityTradebookDetails SendEmailCommodityTradebook

	sendEmailCommodityTradebookDetails.ClientId = profileData.ClientID
	sendEmailCommodityTradebookDetails.ApplicantName = profileData.Name
	sendEmailCommodityTradebookDetails.RecipientEmail = profileData.EmailID
	sendEmailCommodityTradebookDetails.DateFrom = commodityTradebookReq.DFDateFr
	sendEmailCommodityTradebookDetails.DateTo = commodityTradebookReq.DFDateTo
	sendEmailCommodityTradebookDetails.EncodedReportFile = encodedString

	helpers.PublishMessage(constants.TopicExchange, constants.KeyCommodityTradebookReport, sendEmailCommodityTradebookDetails)

	var apiRes apihelpers.APIRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ReportsObj) ViewFnoTradebook(fnoTradebookReq models.FNOTradebookReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {

	var fnoTransaction models.FNOTransactionRes
	var err error

	fileName := constants.FnoTradebookReport + strings.ToUpper(fnoTradebookReq.UserID) + "_" + strings.ToUpper(fnoTradebookReq.DFDateFr) + "_to_" + strings.ToUpper(fnoTradebookReq.DFDateTo)
	storedReportData, _ := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" {

		err = json.Unmarshal([]byte(storedReportData), &fnoTransaction)
		if err != nil {
			loggerconfig.Error("ViewFnoTradebook, error in unmarshalling storedReportData : ", err, " for fileName fileName: ", fileName)
		}
	}
	if storedReportData == "" {
		theBackofficeProvider := v1.GetBackOfficeProvider()
		fnoTransaction, err = theBackofficeProvider.GetFNOTransactionData(fnoTradebookReq, reqH)
		if err != nil {
			loggerconfig.Error("ViewFnoTradebook, there is some error in fetching data from Shilpi api function: Error : ", err, " reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		fnoTransaction.UserDetails = profileData
		reportData, err := json.Marshal(fnoTransaction)
		if err != nil {
			loggerconfig.Error("ViewFnoTradebook, Error in marshalling viewDpChargesRes", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("ViewFnoTradebook, Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}

	var apiRes apihelpers.APIRes
	apiRes.Data = fnoTransaction
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

// CreateExcelFNOTradebook generates the Excel file for FNO transactions and charges
func CreateExcelFNOTradebook(data models.FNOTransactionRes) (*xlsx.File, error) {
	// Create a new workbook
	file := xlsx.NewFile()

	// Create a new sheet
	sheet, err := file.AddSheet("FNOTradebook")
	if err != nil {
		return nil, err
	}

	// Add the title row for the report
	addTitleRowFNO(sheet)

	// Add user details below the title
	writeUserDetails(sheet, data.UserDetails, constants.FnoTradebookReportName)

	// Add an empty row for spacing between sections
	sheet.AddRow()

	// Write the summary of charges
	writeFNOChargesSummary(sheet, data)

	// Add space between sections
	sheet.AddRow()

	// Write the transaction data
	writeFNOTransactionsData(sheet, data.FNOTransactions)

	return file, nil
}

// addTitleRowFNO adds the main title to the top of the Excel sheet for FNO
func addTitleRowFNO(sheet *xlsx.Sheet) {
	// Create the title row
	titleRow := sheet.AddRow()
	titleRow.SetHeight(20) // Optionally set height for emphasis
	titleCell := titleRow.AddCell()
	titleCell.Merge(3, 0) // Merge first 4 cells for the title
	titleCell.SetString("FNO Tradebook and Charges")
	titleCell.GetStyle().Font.Bold = true
	titleCell.GetStyle().Alignment.Horizontal = "center"

	// Add an empty row after the title for spacing
	sheet.AddRow()
}

// writeFNOChargesSummary writes the total charges in the Excel file for FNO
func writeFNOChargesSummary(sheet *xlsx.Sheet, data models.FNOTransactionRes) {
	// Create headers for the summary
	chargesHeaders := []string{"Trade Charges", "Amount"}
	headerRow := sheet.AddRow()
	for _, header := range chargesHeaders {
		headerRow.AddCell().SetString(header)
	}

	// Define the order and labels of the charges
	charges := []struct {
		Label string
		Value float64
	}{
		{"Brokerage", data.TotalBrokerage},
		{"GST", data.TotalGST},
		{"SEBI Tax", data.TotalSEBITax},
		{"STT", data.TotalSTT},
		{"Exchange Turnover Charges", data.TotalTurnCharges},
		{"Stamp Duty", data.TotalStampDuty},
		{"Clearing Charges", data.TotalClearingCharges},
		{"Total Charges", data.TotalCharges},
	}

	// Populate the charge data
	for _, charge := range charges {
		row := sheet.AddRow()
		row.AddCell().SetString(charge.Label)
		row.AddCell().SetString(strconv.FormatFloat(charge.Value, 'f', 2, 64))
	}
}

// writeFNOTransactionsData writes the FNO transactions in the Excel file
func writeFNOTransactionsData(sheet *xlsx.Sheet, transactions []models.FNOTransactionData) {
	// Create headers for the transaction data
	transactionHeaders := []string{
		"Symbol", "Instrument Type", "Expiry Date", "Option Type", "Strike Price", "Transaction Date",
		"Buy/Sell", "Trade Price", "Qty", "Brokerage", "GST", "STT", "SEBI Fees",
		"Turnover Fees", "Stamp Duty", "IPF Tax", "Clearing Charges", "Segment", "Exchange", "Order No", "Trade ID", "Trade Time",
	}
	headerRow := sheet.AddRow()
	for _, header := range transactionHeaders {
		headerRow.AddCell().SetString(header)
	}

	// Write each transaction as a row
	for _, trx := range transactions {
		dataRow := sheet.AddRow()
		dataRow.AddCell().SetString(trx.Symbol)
		dataRow.AddCell().SetString(trx.InstrumentType)
		dataRow.AddCell().SetString(trx.ExpiryDate)
		dataRow.AddCell().SetString(trx.OptionType)
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.StrikePrice, 'f', 2, 64))
		dataRow.AddCell().SetString(trx.TradeDate)
		dataRow.AddCell().SetString(trx.BuySellInd) // Already handled conversion
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.TradePrice, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.TradeQty, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.Brokerage, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.GST, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.STT, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.SEBITax, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.TurnoverTax, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.StampDuty, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.IPFTax, 'f', 2, 64))
		dataRow.AddCell().SetString(strconv.FormatFloat(trx.ClearingCharges, 'f', 2, 64))
		dataRow.AddCell().SetString(trx.Segment)
		dataRow.AddCell().SetString(trx.Exchange)
		dataRow.AddCell().SetString(trx.OrderNo)
		dataRow.AddCell().SetString(trx.TradeNo)
		dataRow.AddCell().SetString(trx.TradeTime)
	}
}

func (obj ReportsObj) DownloadFnoTradebook(fnoTradebookReq models.FNOTradebookReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {

	var fnoTransaction models.FNOTransactionRes
	var err error

	fileName := constants.FnoTradebookReport + strings.ToUpper(fnoTradebookReq.UserID) + "_" + strings.ToUpper(fnoTradebookReq.DFDateFr) + "_to_" + strings.ToUpper(fnoTradebookReq.DFDateTo)
	storedReportData, _ := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" {

		err = json.Unmarshal([]byte(storedReportData), &fnoTransaction)
		if err != nil {
			loggerconfig.Error("DownloadFnoTradebook, error in unmarshalling storedReportData : ", err, " for fileName fileName: ", fileName)
		}
	}
	if storedReportData == "" {
		theBackofficeProvider := v1.GetBackOfficeProvider()
		fnoTransaction, err = theBackofficeProvider.GetFNOTransactionData(fnoTradebookReq, reqH)
		if err != nil {
			loggerconfig.Error("DownloadFnoTradebook, there is some error in fetching data from Shilpi api function: Error : ", err, " reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		fnoTransaction.UserDetails = profileData
		reportData, err := json.Marshal(fnoTransaction)
		if err != nil {
			loggerconfig.Error("DownloadFnoTradebook, Error in marshalling viewDpChargesRes", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("DownloadFnoTradebook, Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}

	file, err := CreateExcelFNOTradebook(fnoTransaction)
	if err != nil {
		loggerconfig.Error("DownloadFnoTradebook, Error in getting excel file, error: ", err, "reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	s3FileName := fileName + ".xlsx"
	expiryHours := int64(24)

	// Generate a pre-signed URL for the XLSX file
	url, err := helpers.UploadFileToS3AndGetPresignedURL(constants.FnoTradebookS3FolderName, s3FileName, file, expiryHours)
	if err != nil {
		loggerconfig.Error("DownloadFnoTradebook, failed to generate pre-signed URL:", err, "reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var downloadFnoTradebookRes models.DownloadFnoTradebookRes
	downloadFnoTradebookRes.DownloadUrl = url

	var apiRes apihelpers.APIRes
	apiRes.Data = downloadFnoTradebookRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}

func (obj ReportsObj) SendEmailFnoTradebook(fnoTradebookReq models.FNOTradebookReq, reqH models.ReqHeader, profileData models.ProfileDataResp) (int, apihelpers.APIRes) {

	var fnoTransaction models.FNOTransactionRes
	var err error

	fileName := constants.FnoTradebookReport + strings.ToUpper(fnoTradebookReq.UserID) + "_" + strings.ToUpper(fnoTradebookReq.DFDateFr) + "_to_" + strings.ToUpper(fnoTradebookReq.DFDateTo)
	storedReportData, _ := dbops.RedisRepo.Get(fileName)
	if storedReportData != "" {

		err = json.Unmarshal([]byte(storedReportData), &fnoTransaction)
		if err != nil {
			loggerconfig.Error("DownloadFnoTradebook, error in unmarshalling storedReportData : ", err, " for fileName fileName: ", fileName)
		}
	}
	if storedReportData == "" {
		theBackofficeProvider := v1.GetBackOfficeProvider()
		fnoTransaction, err = theBackofficeProvider.GetFNOTransactionData(fnoTradebookReq, reqH)
		if err != nil {
			loggerconfig.Error("DownloadFnoTradebook, there is some error in fetching data from Shilpi api function: Error : ", err, "reqId: ", reqH.RequestId)
			return apihelpers.SendInternalServerError()
		}
		fnoTransaction.UserDetails = profileData
		reportData, err := json.Marshal(fnoTransaction)
		if err != nil {
			loggerconfig.Error("DownloadFnoTradebook, Error in marshalling viewDpChargesRes", err)
		} else {
			err = dbops.RedisRepo.Set(fileName, string(reportData), constants.ReportsCachingTTL*time.Minute)
			if err != nil {
				loggerconfig.Error("DownloadFnoTradebook, Report Data not written to redis:", fileName, " Failed to write to redis:", err)
			}
		}
	}

	file, err := CreateExcelFNOTradebook(fnoTransaction)
	if err != nil {
		loggerconfig.Error("DownloadFnoTradebook, Error in getting excel file, error: ", err, "reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	encodedString, err := helpers.EncodeExcelToBase64(file)
	if err != nil {
		loggerconfig.Error("SendEmailLedger, Error in creating base64 format from Excel, error:", err, "reqId: ", reqH.RequestId)
		return apihelpers.SendInternalServerError()
	}

	var sendEmailFnoTradebookDetails SendEmailFnoTradebook

	sendEmailFnoTradebookDetails.ClientId = profileData.ClientID
	sendEmailFnoTradebookDetails.ApplicantName = profileData.Name
	sendEmailFnoTradebookDetails.RecipientEmail = profileData.EmailID
	sendEmailFnoTradebookDetails.DateFrom = fnoTradebookReq.DFDateFr
	sendEmailFnoTradebookDetails.DateTo = fnoTradebookReq.DFDateTo
	sendEmailFnoTradebookDetails.EncodedReportFile = encodedString

	helpers.PublishMessage(constants.TopicExchange, constants.KeyFnoTradebookReport, sendEmailFnoTradebookDetails)

	var apiRes apihelpers.APIRes
	apiRes.Message = "SUCCESS"
	apiRes.Status = true

	return http.StatusOK, apiRes
}
