package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"space/constants"
	"space/dbops"
	"space/loggerconfig"
	"space/models"

	"go.mongodb.org/mongo-driver/bson"

	_ "github.com/lib/pq"
)

type Postgres struct {
	conn *sql.DB
}

// pgObj
var pgObj *Postgres

// singleTon object
var oncePg sync.Once

func GetPgObj() *Postgres {
	oncePg.Do(func() {
		pgObj = &Postgres{}
		pgObj.postgresConnect()
	})

	return pgObj
}

func GetPostgresObject() *Postgres {
	return pgObj
}

func ReconnectToPostgres() error {
	pgObj = &Postgres{}
	pgObj.postgresConnect()
	ctx := context.Background()
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("ReconnectToPostgres Error ping : " + err.Error())
		return err
	}
	return nil
}

// PostgresConnection Connection Open
func (pgObj *Postgres) postgresConnect() {

	loggerconfig.GetConfig().SetConfigName("config")
	loggerconfig.GetConfig().AddConfigPath("./config")
	loggerconfig.Start()

	var dbinfo string
	dbinfo = fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s TimeZone=Asia/Kolkata sslrootcert=%s",
		constants.ServerIP, constants.Port, constants.User,
		constants.Password, constants.DB, constants.CertificateEnabled, constants.CertificatePath)

	var err error

	// Create connection pool
	pgObj.conn, err = sql.Open("postgres", dbinfo+" connect_timeout=1")
	if err != nil {
		loggerconfig.Error("Error creating pg connection pool :" + err.Error())
		return
	}

	pgObj.conn.SetMaxOpenConns(constants.DBMaxConn)
	pgObj.conn.SetMaxIdleConns(constants.DBMaxIdleConn)
	pgObj.conn.SetConnMaxLifetime(24 * time.Hour)
	// pgObj.conn.SetConnMaxLifetime(5 * time.Second)

	ctx := context.Background()

	err = pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error(" Error ping : " + err.Error())
		return
	}
}

func FetchOverviewQuery(searchBy string, ttmTypeCS string, cpkpTypeCS string) string {
	queryStatement := fmt.Sprintf(`SELECT ttmdata.mcap, ttmdata.pbttm, ttmdata.pettm, ttmdata.sectorpe, 
	ttmdata.roettm, ttmdata.epsttm, ttmdata.dividendyield, ttmdata.bookvalue, ttmdata.debttoequity,
	companypeerkeyparams.netprofit, companyprofile.lname, companyprofile.hsesname,
	companyprofile.incdt, companyprofile.indlname, companyprofile.auditor,
	companyprofile.chairman, companyprofile.cosec, companyprofile.internet
	FROM ttmdata 
	INNER JOIN companymaster ON ttmdata.cocode = companymaster.cocode 
	INNER JOIN companypeerkeyparams ON ttmdata.cocode = companypeerkeyparams.cocode
	INNER JOIN companyprofile ON ttmdata.cocode = companyprofile.cocode
	WHERE companymaster.%s = $1 AND ttmdata.typecs=%s AND companypeerkeyparams.typecs=%s;`, searchBy, ttmTypeCS, cpkpTypeCS)
	return queryStatement
}

func (pgObj *Postgres) FetchOverviewData(req models.GetOverviewReq) ([]models.GetOverviewRes, error) {
	ctx := context.Background()
	var dbResponse []models.GetOverviewRes
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchOverviewData Error if database is alive :", err.Error())
		return dbResponse, err
	}
	var queryStatementConsolidate string
	var queryStatementStandalone string
	var resConsolidate, resStandalone *sql.Rows
	if req.Isin != "" {
		loggerconfig.Info("FetchOverviewData isin =", req.Isin)
		queryStatementConsolidate = FetchOverviewQuery(constants.ISIN, constants.TypeCapsConsolidate, constants.TypeSmallConsolidate)
		queryStatementStandalone = FetchOverviewQuery(constants.ISIN, constants.TypeCapsStandalone, constants.TypeSmallStandalone)
		resConsolidate, err = dbops.PostgresRepo.Fetch(queryStatementConsolidate, req.Isin)
		if err != nil {
			loggerconfig.Error("FetchOverviewData Error isin fetching data Consolidate:", err.Error())
			return dbResponse, err
		}
		resStandalone, err = dbops.PostgresRepo.Fetch(queryStatementStandalone, req.Isin)
		if err != nil {
			loggerconfig.Error("FetchOverviewData Error isin fetching data Standalone:", err.Error())
			return dbResponse, err
		}

	} else if req.Exchange == "" {
		loggerconfig.Error("FetchOverviewData Error no exchange entered in request!")
		return dbResponse, err
	} else if strings.ToLower(req.Exchange) == constants.NSE && req.NseSymbol != "" {
		loggerconfig.Info("FetchOverviewData NseSymbol =", req.NseSymbol)
		queryStatementConsolidate = FetchOverviewQuery(constants.NSESymbol, constants.TypeCapsConsolidate, constants.TypeSmallConsolidate)
		queryStatementStandalone = FetchOverviewQuery(constants.NSESymbol, constants.TypeCapsStandalone, constants.TypeSmallStandalone)
		resConsolidate, err = dbops.PostgresRepo.Fetch(queryStatementConsolidate, req.NseSymbol)
		if err != nil {
			loggerconfig.Error("FetchOverviewData Error NseSymbol fetching data Consolidate:", err.Error())
			return dbResponse, err
		}
		resStandalone, err = dbops.PostgresRepo.Fetch(queryStatementStandalone, req.NseSymbol)
		if err != nil {
			loggerconfig.Error("FetchOverviewData Error NseSymbol fetching data Standalone:", err.Error())
			return dbResponse, err
		}
	} else if strings.ToLower(req.Exchange) == constants.BSE && req.BseToken != "" {
		loggerconfig.Info("FetchOverviewData BseToken =", req.BseToken)
		queryStatementConsolidate = FetchOverviewQuery(constants.BSECode, constants.TypeCapsConsolidate, constants.TypeSmallConsolidate)
		queryStatementStandalone = FetchOverviewQuery(constants.BSECode, constants.TypeCapsStandalone, constants.TypeSmallStandalone)
		resConsolidate, err = dbops.PostgresRepo.Fetch(queryStatementConsolidate, req.BseToken)
		if err != nil {
			loggerconfig.Error("FetchOverviewData Error NseSymbol fetching data Consolidate:", err.Error())
			return dbResponse, err
		}
		resStandalone, err = dbops.PostgresRepo.Fetch(queryStatementStandalone, req.BseToken)
		if err != nil {
			loggerconfig.Error("FetchOverviewData Error NseSymbol fetching data Standalone:", err.Error())
			return dbResponse, err
		}
	} else {
		loggerconfig.Error("FetchOverviewData Error no NSE Symbol or BSE Token entered in request!")
		return dbResponse, err
	}

	var row models.GetOverviewData
	defer resConsolidate.Close()
	defer resStandalone.Close()

	if resConsolidate.Next() {
		err = resConsolidate.Scan(&row.MarketCap, &row.PbRatio, &row.PeRatio, &row.IndustryPE, &row.Roe, &row.Eps, &row.DivYield, &row.BookValue, &row.DebttoEquity, &row.NetProfit, &row.CompanyName, &row.HseSName, &row.EstablishedYear, &row.Industry, &row.Auditor, &row.Chairman, &row.CoSec, &row.Website)
		if err != nil {
			return dbResponse, err
		}
		var entryRow models.GetOverviewRes
		entryRow.MarketCap = row.MarketCap
		entryRow.PbRatio = row.PbRatio
		entryRow.PeRatio = row.PeRatio
		entryRow.IndustryPE = row.IndustryPE
		entryRow.Roe = row.Roe
		entryRow.Eps = row.Eps
		entryRow.DivYield = row.DivYield
		entryRow.BookValue = row.BookValue
		entryRow.DebttoEquity = row.DebttoEquity
		entryRow.NetProfit = row.NetProfit
		entryRow.AboutCompany = "CompanyName:" + row.CompanyName + ", HseSName:" + row.HseSName + ", EstablishedYear:" + row.EstablishedYear + ", Industry:" + row.Industry + ", Auditor:" + row.Auditor + ", Chairman:" + row.Chairman + ", CoSec:" + row.CoSec + ", Website:" + row.Website

		dbResponse = append(dbResponse, entryRow)
		return dbResponse, nil
	}

	if resStandalone.Next() {
		err = resStandalone.Scan(&row.MarketCap, &row.PbRatio, &row.PeRatio, &row.IndustryPE, &row.Roe, &row.Eps, &row.DivYield, &row.BookValue, &row.DebttoEquity, &row.NetProfit, &row.CompanyName, &row.HseSName, &row.EstablishedYear, &row.Industry, &row.Auditor, &row.Chairman, &row.CoSec, &row.Website)
		if err != nil {
			return dbResponse, err
		}
		var entryRow models.GetOverviewRes
		entryRow.MarketCap = row.MarketCap
		entryRow.PbRatio = row.PbRatio
		entryRow.PeRatio = row.PeRatio
		entryRow.IndustryPE = row.IndustryPE
		entryRow.Roe = row.Roe
		entryRow.Eps = row.Eps
		entryRow.DivYield = row.DivYield
		entryRow.BookValue = row.BookValue
		entryRow.DebttoEquity = row.DebttoEquity
		entryRow.NetProfit = row.NetProfit
		entryRow.AboutCompany = "CompanyName:" + row.CompanyName + ", HseSName:" + row.HseSName + ", EstablishedYear:" + row.EstablishedYear + ", Industry:" + row.Industry + ", Auditor:" + row.Auditor + ", Chairman:" + row.Chairman + ", CoSec:" + row.CoSec + ", Website:" + row.Website

		dbResponse = append(dbResponse, entryRow)
	}

	return dbResponse, nil
}

func (pgObj *Postgres) FetchFinancialsData(req models.FetchFinancialsReq) (models.FetchFinancialsRes, error) {
	ctx := context.Background()
	var dbResponse models.FetchFinancialsRes
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchFinancialsData Error if database is alive :", err.Error())
		return dbResponse, err
	}
	var queryStatement string
	var res *sql.Rows
	if req.Isin != "" {
		loggerconfig.Info("FetchFinancialsData isin =", req.Isin)
		queryStatement = fmt.Sprintf(`SELECT PL.*
		FROM plstatement AS PL
		JOIN companymaster AS CM ON PL.cocode = CM.cocode
		WHERE PL.columnname = %s 
		AND CM.isin = $1;`, constants.TotalExpenses)
		loggerconfig.Info(queryStatement)
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.Isin)
		if err != nil {
			loggerconfig.Error("FetchFinancialsData Error isin fetching data:", err.Error())
			return dbResponse, err
		}
		var totalExpenses models.FetchFinancialsData
		defer res.Close()
		for res.Next() {
			err = res.Scan(&totalExpenses.CoCode, &totalExpenses.TypeCS, &totalExpenses.Columnname, &totalExpenses.Rid, &totalExpenses.Yc0, &totalExpenses.Yc1, &totalExpenses.Yc2, &totalExpenses.Yc3, &totalExpenses.Yc4, &totalExpenses.Rowno)
			loggerconfig.Info("error:", err)
			if err != nil {
				return dbResponse, err
			}
		}

		queryStatement = fmt.Sprintf(`SELECT PL.*
		FROM plstatement AS PL
		JOIN companymaster AS CM ON PL.cocode = CM.cocode
		WHERE PL.columnname = %s
		  AND CM.isin = $1;`, constants.TotalRevenue)
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.Isin)
		if err != nil {
			loggerconfig.Error("FetchFinancialsData Error isin fetching data:", err.Error())
			return dbResponse, err
		}
		var totalRevenue models.FetchFinancialsData
		defer res.Close()
		for res.Next() {
			err = res.Scan(&totalRevenue.CoCode, &totalRevenue.TypeCS, &totalRevenue.Columnname, &totalRevenue.Rid, &totalRevenue.Yc0, &totalRevenue.Yc1, &totalRevenue.Yc2, &totalRevenue.Yc3, &totalRevenue.Yc4, &totalRevenue.Rowno)
			if err != nil {
				return dbResponse, err
			}
		}

		var netProfit models.FetchFinancialsData

		netProfit.Columnname = constants.NetProfit
		netProfit.CoCode = totalRevenue.CoCode
		netProfit.TypeCS = totalRevenue.TypeCS
		netProfit.Rid = totalRevenue.Rid
		netProfit.Rowno = totalRevenue.Rowno
		netProfit.Yc0 = totalRevenue.Yc0 - totalExpenses.Yc0
		netProfit.Yc1 = totalRevenue.Yc1 - totalExpenses.Yc1
		netProfit.Yc2 = totalRevenue.Yc2 - totalExpenses.Yc2
		netProfit.Yc3 = totalRevenue.Yc0 - totalExpenses.Yc3
		netProfit.Yc4 = totalRevenue.Yc0 - totalExpenses.Yc4

		queryStatement = fmt.Sprintf(`SELECT BS.*
		FROM balancesheets AS BS
		JOIN companymaster AS CM ON BS.cocode = CM.cocode
		WHERE BS.columnname = %s
		  AND CM.isin = $1;`, constants.NetCurrentAssets)
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.Isin)
		if err != nil {
			loggerconfig.Error("FetchFinancialsData Error isin fetching data:", err.Error())
			return dbResponse, err
		}
		var netCurrentAssets models.FetchFinancialsData
		defer res.Close()
		for res.Next() {
			err = res.Scan(&netCurrentAssets.CoCode, &netCurrentAssets.TypeCS, &netCurrentAssets.Rid, &netCurrentAssets.Columnname, &netCurrentAssets.Yc0, &netCurrentAssets.Yc1, &netCurrentAssets.Yc2, &netCurrentAssets.Yc3, &netCurrentAssets.Yc4, &netCurrentAssets.Rowno)
			if err != nil {
				return dbResponse, err
			}
		}

		queryStatement = `SELECT CR.*
		FROM cashflowratios AS CR
		JOIN companymaster AS CM ON CR.cocode = CM.cocode
		WHERE CM.isin = $1;`
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.Isin)
		if err != nil {
			loggerconfig.Error("FetchFinancialsData Error isin fetching data:", err.Error())
			return dbResponse, err
		}
		var cfr models.CashFlowRatios
		var cfrArray []models.CashFlowRatios
		defer res.Close()
		for res.Next() {
			err = res.Scan(&cfr.CoCode, &cfr.TypeCS, &cfr.Yrc, &cfr.CashFlowPerShare, &cfr.PricetoCashFlowRatio, &cfr.FreeCashFlowperShare, &cfr.PricetoFreeCashFlow, &cfr.FreeCashFlowYield, &cfr.Salestocashflowratio)
			if err != nil {
				return dbResponse, err
			}
			cfrArray = append(cfrArray, cfr)
		}

		dbResponse.Revenue = totalRevenue
		dbResponse.NetProfit = netProfit
		dbResponse.BalanceSheet = netCurrentAssets
		dbResponse.Cashflow = cfrArray

	} else if req.Exchange == "" {
		loggerconfig.Error("FetchFinancialsData Error no exchange entered in request!")
		return dbResponse, err
	} else if strings.ToLower(req.Exchange) == constants.NSE && req.NseSymbol != "" {
		loggerconfig.Info("FetchFinancialsData NseSymbol =", req.NseSymbol)
		queryStatement = fmt.Sprintf(`SELECT PL.*
		FROM plstatement AS PL
		JOIN companymaster AS CM ON PL.cocode = CM.cocode
		WHERE PL.columnname = %s 
		AND CM.nsesymbol = $1;`, constants.TotalExpenses)
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.NseSymbol)
		if err != nil {
			loggerconfig.Error("FetchFinancialsData Error NSESymbol fetching data:", err.Error())
			return dbResponse, err
		}
		var totalExpenses models.FetchFinancialsData
		defer res.Close()
		for res.Next() {
			err = res.Scan(&totalExpenses.CoCode, &totalExpenses.TypeCS, &totalExpenses.Columnname, &totalExpenses.Rid, &totalExpenses.Yc0, &totalExpenses.Yc1, &totalExpenses.Yc2, &totalExpenses.Yc3, &totalExpenses.Yc4, &totalExpenses.Rowno)
			if err != nil {
				return dbResponse, err
			}
		}

		queryStatement = fmt.Sprintf(`SELECT PL.*
		FROM plstatement AS PL
		JOIN companymaster AS CM ON PL.cocode = CM.cocode
		WHERE PL.columnname = %s
		  AND CM.nsesymbol = $1;`, constants.TotalRevenue)
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.NseSymbol)
		if err != nil {
			loggerconfig.Error("FetchFinancialsData Error NSESymbol fetching data:", err.Error())
			return dbResponse, err
		}
		var totalRevenue models.FetchFinancialsData
		defer res.Close()
		for res.Next() {
			err = res.Scan(&totalRevenue.CoCode, &totalRevenue.TypeCS, &totalRevenue.Columnname, &totalRevenue.Rid, &totalRevenue.Yc0, &totalRevenue.Yc1, &totalRevenue.Yc2, &totalRevenue.Yc3, &totalRevenue.Yc4, &totalRevenue.Rowno)
			if err != nil {
				return dbResponse, err
			}
		}

		var netProfit models.FetchFinancialsData

		netProfit.Columnname = constants.NetProfit
		netProfit.CoCode = totalRevenue.CoCode
		netProfit.TypeCS = totalRevenue.TypeCS
		netProfit.Rid = totalRevenue.Rid
		netProfit.Rowno = totalRevenue.Rowno
		netProfit.Yc0 = totalRevenue.Yc0 - totalExpenses.Yc0
		netProfit.Yc1 = totalRevenue.Yc1 - totalExpenses.Yc1
		netProfit.Yc2 = totalRevenue.Yc2 - totalExpenses.Yc2
		netProfit.Yc3 = totalRevenue.Yc0 - totalExpenses.Yc3
		netProfit.Yc4 = totalRevenue.Yc0 - totalExpenses.Yc4

		queryStatement = fmt.Sprintf(`SELECT BS.*
		FROM balancesheets AS BS
		JOIN companymaster AS CM ON BS.cocode = CM.cocode
		WHERE BS.columnname = %s
		  AND CM.nsesymbol = $1;`, constants.NetCurrentAssets)
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.NseSymbol)
		if err != nil {
			loggerconfig.Error("FetchFinancialsData Error NSESymbol fetching data:", err.Error())
			return dbResponse, err
		}
		var netCurrentAssets models.FetchFinancialsData
		defer res.Close()
		for res.Next() {
			err = res.Scan(&netCurrentAssets.CoCode, &netCurrentAssets.TypeCS, &netCurrentAssets.Rid, &netCurrentAssets.Columnname, &netCurrentAssets.Yc0, &netCurrentAssets.Yc1, &netCurrentAssets.Yc2, &netCurrentAssets.Yc3, &netCurrentAssets.Yc4, &netCurrentAssets.Rowno)
			if err != nil {
				return dbResponse, err
			}
		}

		queryStatement = `SELECT CR.*
		FROM cashflowratios AS CR
		JOIN companymaster AS CM ON CR.cocode = CM.cocode
		WHERE CM.nsesymbol = $1;`
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.NseSymbol)
		if err != nil {
			loggerconfig.Error("FetchFinancialsData Error NSESymbol fetching data:", err.Error())
			return dbResponse, err
		}
		var cfr models.CashFlowRatios
		var cfrArray []models.CashFlowRatios
		defer res.Close()
		for res.Next() {
			err = res.Scan(&cfr.CoCode, &cfr.TypeCS, &cfr.Yrc, &cfr.CashFlowPerShare, &cfr.PricetoCashFlowRatio, &cfr.FreeCashFlowperShare, &cfr.PricetoFreeCashFlow, &cfr.FreeCashFlowYield, &cfr.Salestocashflowratio)
			if err != nil {
				return dbResponse, err
			}
			cfrArray = append(cfrArray, cfr)
		}

		dbResponse.Revenue = totalRevenue
		dbResponse.NetProfit = netProfit
		dbResponse.BalanceSheet = netCurrentAssets
		dbResponse.Cashflow = cfrArray

	} else if strings.ToLower(req.Exchange) == constants.BSE && req.BseToken != "" {
		loggerconfig.Info("FetchFinancialsData BseToken =", req.BseToken)
		queryStatement = fmt.Sprintf(`SELECT PL.*
		FROM plstatement AS PL
		JOIN companymaster AS CM ON PL.cocode = CM.cocode
		WHERE PL.columnname = %s 
		AND CM.bsecode = $1;`, constants.TotalExpenses)
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.BseToken)
		if err != nil {
			loggerconfig.Error("FetchFinancialsData Error Bsecode fetching data:", err.Error())
			return dbResponse, err
		}
		var totalExpenses models.FetchFinancialsData
		defer res.Close()
		for res.Next() {
			err = res.Scan(&totalExpenses.CoCode, &totalExpenses.TypeCS, &totalExpenses.Columnname, &totalExpenses.Rid, &totalExpenses.Yc0, &totalExpenses.Yc1, &totalExpenses.Yc2, &totalExpenses.Yc3, &totalExpenses.Yc4, &totalExpenses.Rowno)
			if err != nil {
				return dbResponse, err
			}
		}

		queryStatement = fmt.Sprintf(`SELECT PL.*
		FROM plstatement AS PL
		JOIN companymaster AS CM ON PL.cocode = CM.cocode
		WHERE PL.columnname = %s
		  AND CM.bsecode = $1;`, constants.TotalRevenue)
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.BseToken)
		if err != nil {
			loggerconfig.Error("FetchFinancialsData Error Bsecode fetching data:", err.Error())
			return dbResponse, err
		}
		var totalRevenue models.FetchFinancialsData
		defer res.Close()
		for res.Next() {
			err = res.Scan(&totalRevenue.CoCode, &totalRevenue.TypeCS, &totalRevenue.Columnname, &totalRevenue.Rid, &totalRevenue.Yc0, &totalRevenue.Yc1, &totalRevenue.Yc2, &totalRevenue.Yc3, &totalRevenue.Yc4, &totalRevenue.Rowno)
			if err != nil {
				return dbResponse, err
			}
		}

		var netProfit models.FetchFinancialsData

		netProfit.Columnname = constants.NetProfit
		netProfit.CoCode = totalRevenue.CoCode
		netProfit.TypeCS = totalRevenue.TypeCS
		netProfit.Rid = totalRevenue.Rid
		netProfit.Rowno = totalRevenue.Rowno
		netProfit.Yc0 = totalRevenue.Yc0 - totalExpenses.Yc0
		netProfit.Yc1 = totalRevenue.Yc1 - totalExpenses.Yc1
		netProfit.Yc2 = totalRevenue.Yc2 - totalExpenses.Yc2
		netProfit.Yc3 = totalRevenue.Yc0 - totalExpenses.Yc3
		netProfit.Yc4 = totalRevenue.Yc0 - totalExpenses.Yc4

		queryStatement = fmt.Sprintf(`SELECT BS.*
		FROM balancesheets AS BS
		JOIN companymaster AS CM ON BS.cocode = CM.cocode
		WHERE BS.columnname = %s
		  AND CM.bsecode = $1;`, constants.NetCurrentAssets)
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.BseToken)
		if err != nil {
			loggerconfig.Error("FetchFinancialsData Error Bsecode fetching data:", err.Error())
			return dbResponse, err
		}
		var netCurrentAssets models.FetchFinancialsData
		defer res.Close()
		for res.Next() {
			err = res.Scan(&netCurrentAssets.CoCode, &netCurrentAssets.TypeCS, &netCurrentAssets.Rid, &netCurrentAssets.Columnname, &netCurrentAssets.Yc0, &netCurrentAssets.Yc1, &netCurrentAssets.Yc2, &netCurrentAssets.Yc3, &netCurrentAssets.Yc4, &netCurrentAssets.Rowno)
			if err != nil {
				return dbResponse, err
			}
		}

		queryStatement = `SELECT CR.*
		FROM cashflowratios AS CR
		JOIN companymaster AS CM ON CR.cocode = CM.cocode
		WHERE CM.bsecode = $1;`
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.BseToken)
		if err != nil {
			loggerconfig.Error("FetchFinancialsData Error Bsecode fetching data:", err.Error())
			return dbResponse, err
		}
		var cfr models.CashFlowRatios
		var cfrArray []models.CashFlowRatios
		defer res.Close()
		for res.Next() {
			err = res.Scan(&cfr.CoCode, &cfr.TypeCS, &cfr.Yrc, &cfr.CashFlowPerShare, &cfr.PricetoCashFlowRatio, &cfr.FreeCashFlowperShare, &cfr.PricetoFreeCashFlow, &cfr.FreeCashFlowYield, &cfr.Salestocashflowratio)
			if err != nil {
				return dbResponse, err
			}
			cfrArray = append(cfrArray, cfr)
		}

		dbResponse.Revenue = totalRevenue
		dbResponse.NetProfit = netProfit
		dbResponse.BalanceSheet = netCurrentAssets
		dbResponse.Cashflow = cfrArray
	} else {
		loggerconfig.Error("FetchFinancialsData Error no NSE Symbol or BSE Token entered in request!")
		return dbResponse, err
	}

	return dbResponse, nil
}

func FetchQuarterlyDataQuery(searchBy string) string {
	queryStatement := fmt.Sprintf(`SELECT QR.*
	FROM quarterlyresults AS QR
	JOIN companymaster AS CM ON QR.cocode = CM.cocode
	WHERE CM.%s = $1 AND typecs='C';`, searchBy)
	return queryStatement
}

func (pgObj *Postgres) FetchQuarterlyData(req models.FetchFinancialsDetailedReq) ([]models.QuarterlyData, error) {
	ctx := context.Background()
	var dbResponse []models.QuarterlyData
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchQuarterlyData Error if database is alive :", err.Error())
		return dbResponse, err
	}
	var queryStatement string
	var res *sql.Rows
	if req.Isin != "" {
		loggerconfig.Info("FetchQuarterlyData isin =", req.Isin)
		queryStatement = FetchQuarterlyDataQuery(constants.ISIN)
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.Isin)
		if err != nil {
			loggerconfig.Error("FetchQuarterlyData Error isin fetching data:", err.Error())
			return dbResponse, err
		}
	} else if req.Exchange == "" {
		loggerconfig.Error("FetchQuarterlyData Error no exchange entered in request!")
		return dbResponse, err
	} else if strings.ToLower(req.Exchange) == constants.NSE && req.NseSymbol != "" {
		loggerconfig.Info("FetchQuarterlyData NseSymbol =", req.NseSymbol)
		queryStatement = FetchQuarterlyDataQuery(constants.NSESymbol)
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.NseSymbol)
		if err != nil {
			loggerconfig.Error("FetchQuarterlyData Error NSESymbol fetching data:", err.Error())
			return dbResponse, err
		}
	} else if strings.ToLower(req.Exchange) == constants.BSE && req.BseToken != "" {
		loggerconfig.Info("FetchQuarterlyData BseToken =", req.BseToken)
		queryStatement = FetchQuarterlyDataQuery(constants.BSECode)
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.BseToken)
		if err != nil {
			loggerconfig.Error("FetchQuarterlyData Error client+transaction fetching data:", err.Error())
			return dbResponse, err
		}
	} else {
		loggerconfig.Error("FetchQuarterlyData Error no NSE Symbol or BSE Token entered in request!")
		return dbResponse, err
	}

	var row models.QuarterlyData
	defer res.Close()
	for res.Next() {
		err = res.Scan(&row.CoCode, &row.Type, &row.Rid, &row.Columnname, &row.Y202212, &row.Y202209, &row.Y202206, &row.Y202203, &row.Y202112, &row.Y202109, &row.Y202106, &row.Y202103, &row.Y202012, &row.Y202009, &row.Y202006, &row.Y202003, &row.Y201912, &row.Y201909, &row.Y201906, &row.Y201903, &row.Y201812, &row.Y201809, &row.Y201806, &row.Y201803, &row.Rowno)
		if err != nil {
			return dbResponse, err
		}
		dbResponse = append(dbResponse, row)
	}

	return dbResponse, nil
}

func FetchPeersDataQuery(fetchByVal string, searchBy string) string {
	queryStatement := fmt.Sprintf(`SELECT companymaster.companyname, companymaster.sectorcode,companymaster.nsesymbol, companymaster.bsecode, %s
	FROM ttmdata
	INNER JOIN companymaster
	ON ttmdata.cocode = companymaster.cocode
	WHERE companymaster.%s = $1 AND typecs='C';`, fetchByVal, searchBy)
	return queryStatement
}

func (pgObj *Postgres) FetchPeersData(req models.FetchPeersReq) ([]models.FetchPeerData, error) {
	ctx := context.Background()
	var dbResponse []models.FetchPeerData
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchPeersData Error if database is alive :", err.Error())
		return dbResponse, err
	}
	var fetchByVal string
	if strings.ToLower(req.FetchBy) != constants.Mcap {
		fetchByVal = "ttmdata." + strings.ToLower(req.FetchBy) + "ttm"
	} else {
		fetchByVal = "ttmdata." + strings.ToLower(req.FetchBy)
	}

	var queryStatement string
	var res *sql.Rows
	if req.Isin != "" {
		loggerconfig.Info("FetchPeersData isin =", req.Isin)
		queryStatement = FetchPeersDataQuery(fetchByVal, constants.ISIN)
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.Isin)
		if err != nil {
			loggerconfig.Error("FetchPeersData Error isin fetching data:", err.Error())
			return dbResponse, err
		}
	} else if req.Exchange == "" {
		loggerconfig.Error("FetchPeersData Error no exchange entered in request!")
		return dbResponse, err
	} else if strings.ToLower(req.Exchange) == constants.NSE && req.NseSymbol != "" {
		loggerconfig.Info("FetchPeersData NseSymbol =", req.NseSymbol)
		queryStatement = FetchPeersDataQuery(fetchByVal, constants.NSESymbol)
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.NseSymbol)
		if err != nil {
			loggerconfig.Error("FetchPeersData Error NSESymbol fetching data:", err.Error())
			return dbResponse, err
		}
	} else if strings.ToLower(req.Exchange) == constants.BSE && req.BseToken != "" {
		loggerconfig.Info("FetchPeersData BseToken =", req.BseToken)
		queryStatement = FetchPeersDataQuery(fetchByVal, constants.BSECode)
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.BseToken)
		if err != nil {
			loggerconfig.Error("FetchPeersData Error client+transaction fetching data:", err.Error())
			return dbResponse, err
		}
	} else {
		loggerconfig.Error("FetchPeersData Error no NSE Symbol or BSE Token entered in request!")
		return dbResponse, err
	}
	var row models.FetchPeerData
	defer res.Close()
	for res.Next() {
		err = res.Scan(&row.Company, &row.SectorCode, &row.TradingSymbol, &row.Token, &row.Filter)
		if err != nil {
			return dbResponse, err
		}
		row.Exchange = constants.BSE
	}
	dbResponse = append(dbResponse, row)

	loggerconfig.Info("FetchPeersData FetchBy =", req.FetchBy)
	queryStatement = fmt.Sprintf(`(SELECT companymaster.companyname, companymaster.sectorcode,companymaster.nsesymbol, companymaster.bsecode, %s
		FROM ttmdata
		INNER JOIN companymaster ON ttmdata.cocode = companymaster.cocode
		WHERE companymaster.sectorcode = $1  AND %s > $2 AND typecs='C'
		ORDER BY %s ASC
		LIMIT 4) 
		UNION ALL 
		(SELECT companymaster.companyname, companymaster.sectorcode,companymaster.nsesymbol, companymaster.bsecode, %s
		FROM ttmdata
		INNER JOIN companymaster ON ttmdata.cocode = companymaster.cocode
		WHERE companymaster.sectorcode = $1  AND %s < $2 AND typecs='C'
		ORDER BY %s DESC
		LIMIT 4)`, fetchByVal, fetchByVal, fetchByVal, fetchByVal, fetchByVal, fetchByVal)
	loggerconfig.Info(queryStatement)
	res, err = dbops.PostgresRepo.Fetch(queryStatement, row.SectorCode, row.Filter)
	if err != nil {
		loggerconfig.Error("FetchPeersData Error client+transaction fetching data:", err.Error())
		return dbResponse, err
	}
	defer res.Close()
	var rowsLowerThanFilterValue []models.FetchPeerData

	for res.Next() {
		err = res.Scan(&row.Company, &row.SectorCode, &row.TradingSymbol, &row.Token, &row.Filter)
		if err != nil {
			return dbResponse, err
		}
		row.Exchange = constants.BSE
		if dbResponse[0].Filter < row.Filter {
			dbResponse = append(dbResponse, row)
		} else {
			rowsLowerThanFilterValue = append(rowsLowerThanFilterValue, row)
		}
	}

	for len(dbResponse)+len(rowsLowerThanFilterValue) > 5 {
		if len(dbResponse) > 3 {
			dbResponse = dbResponse[:len(dbResponse)-1]
		}
		if len(rowsLowerThanFilterValue) > 2 {
			rowsLowerThanFilterValue = rowsLowerThanFilterValue[:len(rowsLowerThanFilterValue)-1]
		}
	}
	dbResponse = append(dbResponse, rowsLowerThanFilterValue...)

	return dbResponse, nil
}

func FetchShareHoldingPatternsDataQueryYRC(searchBy string) string {
	queryStatement := fmt.Sprintf(`SELECT DISTINCT SHP.*,SHPD.Ppimf, SHPD.Ppifii, SHPD.ppsubtot 
	FROM shareholdingpattern AS SHP
	JOIN companymaster AS CM ON SHP.cocode = CM.cocode
	INNER JOIN shareholdingpatterndetails AS SHPD ON (SHP.cocode = SHPD.cocode AND SHP.yrc=SHPD.yrc)
	WHERE SHP.yrc=$1 AND CM.%s = $2;`, searchBy)
	return queryStatement
}

func FetchShareHoldingPatternsDataQueryNoYRC(searchBy string) string {
	queryStatement := fmt.Sprintf(`SELECT DISTINCT SHP.*,SHPD.Ppimf, SHPD.Ppifii, SHPD.Ppsubtot 
	FROM shareholdingpattern AS SHP
	JOIN companymaster AS CM ON SHP.cocode = CM.cocode
	INNER JOIN shareholdingpatterndetails AS SHPD ON (SHP.cocode = SHPD.cocode AND SHP.yrc=SHPD.yrc)
	WHERE (
		SHP.yrc= (SELECT MAX(YRC) FROM shareholdingpattern shp 
				  INNER JOIN companymaster cm ON shp.cocode=cm.cocode
				  WHERE cm.%s=$1
				 )
	) AND CM.%s = $2;`, searchBy, searchBy)
	return queryStatement
}

func (pgObj *Postgres) FetchShareHoldingPatternsData(req models.ShareHoldingPatternsReq) ([]models.ShareHoldingPatternsRes, error) {
	ctx := context.Background()
	var dbResponse []models.ShareHoldingPatternsRes
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchShareHoldingPatternsData Error if database is alive :", err.Error())
		return dbResponse, err
	}
	var queryStatement string
	var res *sql.Rows
	if req.Isin != "" {
		loggerconfig.Info("FetchShareHoldingPatternsData isin =", req.Isin)
		if req.Yrc != 0 {
			queryStatement = FetchShareHoldingPatternsDataQueryYRC(constants.ISIN)
			res, err = dbops.PostgresRepo.Fetch(queryStatement, req.Yrc, req.Isin)
		} else {
			queryStatement = FetchShareHoldingPatternsDataQueryNoYRC(constants.ISIN)
			res, err = dbops.PostgresRepo.Fetch(queryStatement, req.Isin, req.Isin)
		}

		if err != nil {
			loggerconfig.Error("FetchShareHoldingPatternsData Error isin fetching data:", err.Error())
			return dbResponse, err
		}
	} else if req.Exchange == "" {
		loggerconfig.Error("FetchShareHoldingPatternsData Error no exchange entered in request!")
		return dbResponse, err
	} else if strings.ToLower(req.Exchange) == constants.NSE && req.NseSymbol != "" {
		loggerconfig.Info("FetchShareHoldingPatternsData NseSymbol =", req.NseSymbol)
		if req.Yrc != 0 {
			queryStatement = FetchShareHoldingPatternsDataQueryYRC(constants.NSESymbol)
			res, err = dbops.PostgresRepo.Fetch(queryStatement, req.Yrc, req.NseSymbol)
		} else {
			queryStatement = FetchShareHoldingPatternsDataQueryNoYRC(constants.NSESymbol)
			res, err = dbops.PostgresRepo.Fetch(queryStatement, req.NseSymbol, req.NseSymbol)
		}
		if err != nil {
			loggerconfig.Error("FetchShareHoldingPatternsData Error NSESymbol fetching data:", err.Error())
			return dbResponse, err
		}
	} else if strings.ToLower(req.Exchange) == constants.BSE && req.BseToken != "" {
		loggerconfig.Info("FetchShareHoldingPatternsData BseToken =", req.BseToken)
		if req.Yrc != 0 {
			queryStatement = FetchShareHoldingPatternsDataQueryYRC(constants.BSECode)
			res, err = dbops.PostgresRepo.Fetch(queryStatement, req.Yrc, req.BseToken)
		} else {
			queryStatement = FetchShareHoldingPatternsDataQueryNoYRC(constants.BSECode)
			res, err = dbops.PostgresRepo.Fetch(queryStatement, req.BseToken)
		}
		if err != nil {
			loggerconfig.Error("FetchShareHoldingPatternsData Error client+transaction fetching data:", err.Error())
			return dbResponse, err
		}
	} else {
		loggerconfig.Error("FetchShareHoldingPatternsData Error no NSE Symbol or BSE Token entered in request!")
		return dbResponse, err
	}
	var row models.ShareHoldingPatternsRes
	defer res.Close()
	for res.Next() {
		err = res.Scan(&row.CoCode, &row.Yrc, &row.TotalPromoterShares, &row.TotalPromoterPerShares, &row.TotalPromoterPledgeShares, &row.TotalPromoterPerPledgeShares, &row.TotalNoofShareholders, &row.PPIMF, &row.PPIFII, &row.PPSUBTOT)
		if err != nil {
			return dbResponse, err
		}
		dbResponse = append(dbResponse, row)
		dbResponse[len(dbResponse)-1].Other = constants.Hundred - (row.TotalPromoterPerShares + row.PPIMF + row.PPIFII + row.PPSUBTOT)
	}

	return dbResponse, nil
}
func FetchRatiosCompareDataQuery(searchBy string, stringOf string) string {
	queryStatement := fmt.Sprintf(`SELECT CM.companyname, TTM.mcap, TTM.pbttm, TTM.pettm, TTM.sectorpe, TTM.roettm, TTM.epsttm, TTM.dividendyield, TTM.bookvalue 
	FROM ttmdata AS TTM
	JOIN companymaster AS CM ON TTM.cocode = CM.cocode
	WHERE CM.%s IN (%s);`, searchBy, stringOf)
	return queryStatement
}

func (pgObj *Postgres) FetchRatiosCompareData(req models.RatiosCompareReq) ([]models.RatiosCompareRes, error) {
	ctx := context.Background()
	var dbResponse []models.RatiosCompareRes
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchRatiosCompareData Error if database is alive :", err.Error())
		return dbResponse, err
	}

	if len(req.ReqData) > 0 {
		if req.ReqData[0].Isin != "" {
			stringOfIsin := ""
			var queryStatement string
			var res *sql.Rows
			for i := 0; i < len(req.ReqData); i++ {
				stringOfIsin += "'" + req.ReqData[i].Isin + "',"
			}
			stringOfIsin = stringOfIsin[:len(stringOfIsin)-1]

			loggerconfig.Info("FetchRatiosCompareData isin =", stringOfIsin, " length:", len(stringOfIsin))
			queryStatement = FetchRatiosCompareDataQuery(constants.ISIN, stringOfIsin)
			res, err = dbops.PostgresRepo.Fetch(queryStatement)
			if err != nil {
				loggerconfig.Error("FetchRatiosCompareData Error isin fetching data:", err.Error())
				return dbResponse, err
			}

			var row models.RatiosCompareRes
			defer res.Close()
			for res.Next() {
				err = res.Scan(&row.Company, &row.MarketCap, &row.PbRatio, &row.PeRatio, &row.IndustryPE, &row.Roe, &row.Eps, &row.DivYield, &row.BookValue)
				if err != nil {
					return dbResponse, err
				}
				dbResponse = append(dbResponse, row)
			}
			return dbResponse, nil
		}
	}

	stringOfNseSymbols := ""
	var queryStatement string
	var res *sql.Rows
	for i := 0; i < len(req.ReqData); i++ {
		if req.ReqData[i].Exchange == "" {
			loggerconfig.Error("FetchRatiosCompareData Error no exchange entered in request!")
			return dbResponse, err
		} else {
			stringOfNseSymbols += "'" + req.ReqData[i].Symbol + "',"
		}
	}
	if len(stringOfNseSymbols) > 0 {
		stringOfNseSymbols = stringOfNseSymbols[:len(stringOfNseSymbols)-1]
	}

	if len(stringOfNseSymbols) == 0 {
		return dbResponse, errors.New("FetchRatiosCompareData length of string is 0")
	}

	loggerconfig.Info("FetchRatiosCompareData NseSymbols =", stringOfNseSymbols, " length:", len(stringOfNseSymbols))
	queryStatement = FetchRatiosCompareDataQuery(constants.NSESymbol, stringOfNseSymbols)
	res, err = dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchRatiosCompareData Error NSESymbol fetching data:", err.Error())
		return dbResponse, err
	}

	var row models.RatiosCompareRes
	defer res.Close()
	for res.Next() {
		err = res.Scan(&row.Company, &row.MarketCap, &row.PbRatio, &row.PeRatio, &row.IndustryPE, &row.Roe, &row.Eps, &row.DivYield, &row.BookValue)
		if err != nil {
			return dbResponse, err
		}
		dbResponse = append(dbResponse, row)
	}

	return dbResponse, nil
}

var CallFetchSector = func(isin string) (string, string, error) {
	return pgObj.FetchSector(isin)
}

func (pgObj *Postgres) FetchSector(isin string) (string, string, error) {
	ctx := context.Background()
	var sectorCode string
	var sectorName string
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchSector Error if database is alive :", err.Error())
		return sectorCode, sectorName, err
	}
	var queryStatement string
	var res *sql.Rows
	if isin != "" {
		loggerconfig.Info("FetchSector isin =", isin)
		queryStatement = `SELECT companymaster.sectorcode, companymaster.sectorname
		FROM companymaster
		WHERE companymaster.isin = $1;`
		res, err = dbops.PostgresRepo.Fetch(queryStatement, isin)
		if err != nil {
			loggerconfig.Error("FetchSector Error isin fetching data:", err.Error())
			return sectorCode, sectorName, err
		}
	}
	defer res.Close()

	for res.Next() {
		err = res.Scan(&sectorCode, &sectorName)
		if err != nil {
			return sectorCode, sectorName, err
		}
	}

	return sectorCode, sectorName, err
}

var CallFetchPE = func(isin string) (float64, error) {
	return pgObj.FetchPE(isin)
}

func (pgObj *Postgres) FetchPE(isin string) (float64, error) {
	ctx := context.Background()
	var pe float64

	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchPE Error if database is alive :", err.Error())
		return pe, err
	}

	var queryStatement string
	var res *sql.Rows
	if isin != "" {
		loggerconfig.Info("FetchPE isin =", isin)
		queryStatement = `SELECT cpr.pe
		FROM companypeerratios AS cpr
		JOIN companymaster AS CM ON cpr.cocode = CM.cocode
		WHERE CM.isin = $1 AND typecs='c';`
		res, err = dbops.PostgresRepo.Fetch(queryStatement, isin)
		if err != nil {
			loggerconfig.Error("FetchPE Error isin fetching data:", err.Error())
			return pe, err
		}
	}
	defer res.Close()

	for res.Next() {
		err = res.Scan(&pe)
		if err != nil {
			return pe, err
		}
	}

	return pe, err
}

var CallFetchDE = func(isin string) (float64, error) {
	return pgObj.FetchDE(isin)
}

func (pgObj *Postgres) FetchDE(isin string) (float64, error) {
	ctx := context.Background()
	var de float64

	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchPE Error if database is alive :", err.Error())
		return de, err
	}

	var queryStatement string
	var res *sql.Rows
	if isin != "" {
		loggerconfig.Info("FetchPE isin =", isin)
		queryStatement = `SELECT cpr.de
		FROM companypeerratios AS cpr
		JOIN companymaster AS CM ON cpr.cocode = CM.cocode
		WHERE CM.isin = $1 AND typecs='c';`
		res, err = dbops.PostgresRepo.Fetch(queryStatement, isin)
		if err != nil {
			loggerconfig.Error("FetchPE Error isin fetching data:", err.Error())
			return de, err
		}
	}
	defer res.Close()

	for res.Next() {
		err = res.Scan(&de)
		if err != nil {
			return de, err
		}
	}

	return de, err
}

func (pgObj *Postgres) FetchTechnicalIndicatorsData(req models.FetchTechnicalIndicatorsReq) (models.FetchTechnicalIndicatorsRes, error) {
	ctx := context.Background()
	var dbResponse models.FetchTechnicalIndicatorsRes
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchTechnicalIndicatorsData Error if database is alive :", err.Error())
		return dbResponse, err
	}
	var queryStatement string
	var res *sql.Rows
	if req.Isin != "" {
		loggerconfig.Info("FetchTechnicalIndicatorsData isin =", req.Isin)
		queryStatement = `SELECT SE.*,PF.* from SmaEma AS SE 
		INNER JOIN CompanyMaster AS CM ON cm.cocode=SE.cocode
		LEFT JOIN PivotFibonacci AS PF ON SE.cocode=PF.cocode
		WHERE CM.isin= $1 AND frequency=$2;`
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.Isin, req.Frequency)
		if err != nil {
			loggerconfig.Error("FetchTechnicalIndicatorsData Error isin fetching data:", err.Error())
			return dbResponse, err
		}
	}
	defer res.Close()
	var row models.FetchTechnicalIndicatorsRes
	for res.Next() {
		err = res.Scan(&row.CoCode, &row.Rsi, &row.MACD1226Days, &row.Avg20Days, &row.Avg50Days, &row.Avg100Days, &row.Avg200Days, &row.Avg10Days, &row.EMA10Day, &row.EMA20Day, &row.EMA50Day, &row.MACD12269Days, &row.Exchange, &row.CoCode, &row.Frequency, &row.CoName, &row.Currprice, &row.PivotPoint, &row.S1, &row.S2, &row.S3, &row.R1, &row.R2, &row.R3, &row.CurrTime, &row.ExchangeAlt)
		if err != nil {
			loggerconfig.Error("FetchTechnicalIndicatorsData Error Scan data:", err.Error())
			return dbResponse, err
		}
		fmt.Println("the details are : ", row)
		dbResponse = row
	}

	return dbResponse, nil
}

func (pgObj *Postgres) FetchHighPledgePromoterHoldingMatchData(allIsin models.AllIsin) (models.AllHighPledgePromoterHolding, error) {
	var allHighPledgePromoterHolding models.AllHighPledgePromoterHolding
	ctx := context.Background()
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchHighPledgePromoterHoldingMatchData Error if database is alive :", err.Error())
		return allHighPledgePromoterHolding, err
	}

	if allIsin.Isin == nil || len(allIsin.Isin) == 0 {
		return allHighPledgePromoterHolding, fmt.Errorf("FetchHighPledgePromoterHoldingMatchData: input ISIN list is empty or nil")
	}

	stringOfIsin := ""
	for i := 0; i < len(allIsin.Isin); i++ {
		stringOfIsin += "'" + allIsin.Isin[i] + "',"
	}
	stringOfIsin = stringOfIsin[:len(stringOfIsin)-1]

	queryStatement := fmt.Sprintf(`SELECT * FROM highpledgedpromoterholdings where isin IN (%s);`, stringOfIsin)
	res, err := dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchHighPledgePromoterHoldingMatchData Error just client fetching data:", err.Error())
		return allHighPledgePromoterHolding, err
	}
	defer res.Close()

	var highPledgePromoterHolding models.StockDetailDb
	for res.Next() {
		err = res.Scan(&highPledgePromoterHolding.CompanyName, &highPledgePromoterHolding.Isin)
		if err != nil {
			loggerconfig.Error("FetchHighPledgePromoterHoldingMatchData Error in scanning fetched companymaster response from postgres :", err)
			break
		}
		allHighPledgePromoterHolding.HighPledgePromoterHoldingAll = append(allHighPledgePromoterHolding.HighPledgePromoterHoldingAll, highPledgePromoterHolding)
	}
	return allHighPledgePromoterHolding, nil
}

func (pgObj *Postgres) FetchAdditionalSurveillanceMeasureData(allIsin models.AllIsin) (models.AllAdditionalSurveillanceMeasure, error) {
	var allAdditionalSurveillanceMeasure models.AllAdditionalSurveillanceMeasure
	ctx := context.Background()
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchAdditionalSurveillanceMeasureData Error if database is alive :", err.Error())
		return allAdditionalSurveillanceMeasure, err
	}

	if allIsin.Isin == nil || len(allIsin.Isin) == 0 {
		return allAdditionalSurveillanceMeasure, fmt.Errorf("FetchAdditionalSurveillanceMeasureData: input ISIN list is empty or nil")
	}

	stringOfIsin := ""
	for i := 0; i < len(allIsin.Isin); i++ {
		stringOfIsin += "'" + allIsin.Isin[i] + "',"
	}
	stringOfIsin = stringOfIsin[:len(stringOfIsin)-1]

	queryStatement := fmt.Sprintf(`SELECT companyname, isin FROM additional_surveillance_measure_list where isin IN (%s);`, stringOfIsin)
	res, err := dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchAdditionalSurveillanceMeasureData Error just client fetching data:", err.Error())
		return allAdditionalSurveillanceMeasure, err
	}
	defer res.Close()

	var additionalSurveillanceMeasure models.StockDetailDb
	for res.Next() {
		err = res.Scan(&additionalSurveillanceMeasure.CompanyName, &additionalSurveillanceMeasure.Isin)
		if err != nil {
			loggerconfig.Error("FetchAdditionalSurveillanceMeasureData Error in scanning fetched companymaster response from postgres :", err)
			break
		}
		allAdditionalSurveillanceMeasure.AdditionalSurveillanceMeasureAll = append(allAdditionalSurveillanceMeasure.AdditionalSurveillanceMeasureAll, additionalSurveillanceMeasure)
	}
	return allAdditionalSurveillanceMeasure, nil
}

func (pgObj *Postgres) FetchGradedSurveillanceMeasureData(allIsin models.AllIsin) (models.AllGradedSurveillanceMeasure, error) {
	var allGradedSurveillanceMeasure models.AllGradedSurveillanceMeasure
	ctx := context.Background()
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchGradedSurveillanceMeasureData Error if database is alive :", err.Error())
		return allGradedSurveillanceMeasure, err
	}

	if allIsin.Isin == nil || len(allIsin.Isin) == 0 {
		return allGradedSurveillanceMeasure, fmt.Errorf("FetchGradedSurveillanceMeasureData: input ISIN list is empty or nil")
	}

	stringOfIsin := ""
	for i := 0; i < len(allIsin.Isin); i++ {
		stringOfIsin += "'" + allIsin.Isin[i] + "',"
	}
	stringOfIsin = stringOfIsin[:len(stringOfIsin)-1]

	queryStatement := fmt.Sprintf(`SELECT companyname, isin FROM graded_surveillance_measure where isin IN (%s);`, stringOfIsin)
	res, err := dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchGradedSurveillanceMeasureData Error just client fetching data:", err.Error())
		return allGradedSurveillanceMeasure, err
	}
	defer res.Close()

	var gradedSurveillanceMeasure models.StockDetailDb
	for res.Next() {
		err = res.Scan(&gradedSurveillanceMeasure.CompanyName, &gradedSurveillanceMeasure.Isin)
		if err != nil {
			loggerconfig.Error("FetchAdditionalSurveillanceMeasureData Error in scanning fetched companymaster response from postgres :", err)
			break
		}
		allGradedSurveillanceMeasure.AllGradedSurveillanceMeasureAll = append(allGradedSurveillanceMeasure.AllGradedSurveillanceMeasureAll, gradedSurveillanceMeasure)
	}

	return allGradedSurveillanceMeasure, nil
}

func (pgObj *Postgres) FetchRoeData(allIsin models.AllIsin) (models.AllLowRoe, error) {
	var allLowRoe models.AllLowRoe
	ctx := context.Background()
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchGradedSurveillanceMeasureData Error if database is alive :", err.Error())
		return allLowRoe, err
	}

	if allIsin.Isin == nil || len(allIsin.Isin) == 0 {
		return allLowRoe, fmt.Errorf("FetchRoeData: input ISIN list is empty or nil")
	}

	stringOfIsin := ""
	for i := 0; i < len(allIsin.Isin); i++ {
		stringOfIsin += "'" + allIsin.Isin[i] + "',"
	}
	stringOfIsin = stringOfIsin[:len(stringOfIsin)-1]

	queryStatement := fmt.Sprintf(`SELECT CM.companyname, CM.isin, TTM.roettm
	FROM ttmdata AS TTM
	JOIN companymaster AS CM ON TTM.cocode = CM.cocode
	WHERE CM.isin IN (%s);`, stringOfIsin)
	res, err := dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchRoeData Error just client fetching data:", err.Error())
		return allLowRoe, err
	}
	defer res.Close()

	var roeDetails models.RoeStockDetailDb
	for res.Next() {
		err = res.Scan(&roeDetails.CompanyName, &roeDetails.Isin, &roeDetails.Roe)
		if err != nil {
			loggerconfig.Error("FetchRoeData Error in scanning fetched companymaster response from postgres :", err)
			break
		}
		allLowRoe.LowRoeAll = append(allLowRoe.LowRoeAll, roeDetails)
	}

	return allLowRoe, nil
}

func (pgObj *Postgres) FetchTokenAndSymbol(stringOfKeys string, fetchBy string) ([]models.FetchTokenAndSymbol, error) {
	ctx := context.Background()
	var dbResponse []models.FetchTokenAndSymbol

	if len(stringOfKeys) == 0 {
		loggerconfig.Error("FetchTokenAndSymbol stringOfKeys is Empty :")
		return dbResponse, nil
	}

	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchTokenAndSymbol Error if database is alive :", err.Error())
		return dbResponse, err
	}
	var queryStatement string
	var res *sql.Rows

	queryStatement = fmt.Sprintf(`select CM.bsecode,CM.nsesymbol,CM.companyname,CM.cocode from companymaster AS CM where %s IN (%s);`, fetchBy, stringOfKeys)
	res, err = dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchTokenAndSymbol Error fetching data:", err.Error())
		return dbResponse, err
	}
	defer res.Close()
	var row models.FetchTokenAndSymbol
	for res.Next() {
		err = res.Scan(&row.Token, &row.TradingSymbol, &row.CompanyName, &row.CoCode)
		if err != nil {
			loggerconfig.Error("FetchTokenAndSymbol Error Scan data:", err.Error())
			return dbResponse, err
		}
		dbResponse = append(dbResponse, row)
	}
	return dbResponse, nil
}

func (pgObj *Postgres) FetchLowProfitGrowthData(allIsin models.AllIsin) (models.AllProfitabilityGrowthDb, error) {
	var allProfitabilityGrowthDb models.AllProfitabilityGrowthDb
	ctx := context.Background()
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchLowProfitGrowthData Error if database is alive :", err.Error())
		return allProfitabilityGrowthDb, err
	}

	if allIsin.Isin == nil || len(allIsin.Isin) == 0 {
		return allProfitabilityGrowthDb, fmt.Errorf("FetchLowProfitGrowthData: input ISIN list is empty or nil")
	}

	stringOfIsin := ""
	for i := 0; i < len(allIsin.Isin); i++ {
		stringOfIsin += "'" + allIsin.Isin[i] + "',"
	}
	if len(stringOfIsin) > 0 {
		stringOfIsin = stringOfIsin[:len(stringOfIsin)-1]
	}

	queryStatement := fmt.Sprintf(`SELECT CM.companyname, CM.isin, PLS.y0, PLS.y4
	FROM plstatement AS PLS
	JOIN companymaster AS CM ON PLS.cocode = CM.cocode AND PLS.columnname = 'Profit After Tax'
	WHERE CM.isin IN (%s);`, stringOfIsin)
	res, err := dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchLowProfitGrowthData Error just client fetching data:", err.Error())
		return allProfitabilityGrowthDb, err
	}
	defer res.Close()

	var profitabilityGrowth models.ProfitabilityGrowthDb
	for res.Next() {
		err = res.Scan(&profitabilityGrowth.CompanyName, &profitabilityGrowth.Isin, &profitabilityGrowth.YZero, &profitabilityGrowth.YFour)
		if err != nil {
			loggerconfig.Error("FetchLowProfitGrowthData Error in scanning fetched companymaster response from postgres :", err)
			break
		}
		allProfitabilityGrowthDb.ProfitabilityGrowthAll = append(allProfitabilityGrowthDb.ProfitabilityGrowthAll, profitabilityGrowth)
	}

	return allProfitabilityGrowthDb, nil
}

func (pgObj *Postgres) FetchSectorListData(sectorCode string) ([]models.SectorList, error) {
	ctx := context.Background()
	var dbResponse []models.SectorList
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchSectorListData Error if database is alive :", err.Error())
		return dbResponse, err
	}

	var queryStatement string
	if sectorCode != "" {
		queryStatement = fmt.Sprintf(`SELECT * FROM sectorlist WHERE sectcode = '%s';`, sectorCode)
	} else {
		queryStatement = `SELECT * FROM sectorlist;`
	}

	// Execute query
	res, err := dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchSectorListData Error fetching data:", err.Error())
		return dbResponse, err
	}
	defer res.Close()

	for res.Next() {
		var row models.SectorList
		err = res.Scan(&row.SectCode, &row.SectName)
		if err != nil {
			loggerconfig.Error("FetchSectorListData Error Scan data:", err.Error())
			return dbResponse, err
		}
		dbResponse = append(dbResponse, row)
	}

	// Check for errors after iteration
	if err := res.Err(); err != nil {
		loggerconfig.Error("FetchSectorListData Error iterating rows:", err.Error())
		return nil, err
	}

	return dbResponse, nil
}

func (pgObj *Postgres) FetchSectorWiseCompanyData(sectorCode string) ([]models.SectorWiseCompany, error) {
	ctx := context.Background()
	var dbResponse []models.SectorWiseCompany
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchSectorWiseCompanyData Error if database is alive :", err.Error())
		return dbResponse, err
	}
	var queryStatement string
	var res *sql.Rows

	queryStatement = `SELECT sectname FROM sectorlist WHERE sectcode=$1;`
	res, err = dbops.PostgresRepo.Fetch(queryStatement, sectorCode)
	if err != nil {
		loggerconfig.Error("FetchSectorWiseCompanyData Error in fetching sector Name data:", err.Error())
		return dbResponse, err
	}
	var sectorName string
	for res.Next() {
		err = res.Scan(&sectorName)
		if err != nil {
			loggerconfig.Error("FetchSectorWiseCompanyData Error Scan data:", err.Error())
			return dbResponse, err
		}
	}
	var res2 *sql.Rows
	queryStatement = `SELECT * FROM sectorwisecompany WHERE LOWER(sectname)=$1;`
	res2, err = dbops.PostgresRepo.Fetch(queryStatement, strings.ToLower(sectorName))
	if err != nil {
		loggerconfig.Error("FetchSectorWiseCompanyData Error in fetching sector Name data:", err.Error())
		return dbResponse, err
	}

	defer res.Close()
	defer res2.Close()
	var row models.SectorWiseCompany
	for res2.Next() {
		err = res2.Scan(&row.CoCode, &row.CoName, &row.Lname, &row.ScCode, &row.Symbol, &row.SectName, &row.Isin)
		if err != nil {
			loggerconfig.Error("FetchSectorWiseCompanyData Error Scan data:", err.Error())
			return dbResponse, err
		}
		dbResponse = append(dbResponse, row)
	}
	return dbResponse, nil
}

func (pgObj *Postgres) FetchCompanyCategory(stringOfisin string) ([]models.CompanyCategory, error) {
	ctx := context.Background()
	var dbResponse []models.CompanyCategory
	if len(stringOfisin) == 0 {
		loggerconfig.Info("FetchCompanyCategory isinList is Empty :")
		return dbResponse, nil
	}
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchCompanyCategory Error if database is alive :", err.Error())
		return dbResponse, err
	}
	var queryStatement string
	var res *sql.Rows

	queryStatement = fmt.Sprintf(`SELECT CM.isin,CM.mcaptype,CM.industryname,CM.sectorname FROM companymaster AS CM WHERE CM.isin IN (%s);`, stringOfisin)
	res, err = dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchCompanyCategory Error fetching data:", err.Error())
		return dbResponse, err
	}
	defer res.Close()
	var row models.CompanyCategory
	for res.Next() {
		err = res.Scan(&row.Isin, &row.McapType, &row.IndustryName, &row.SectorName)
		if err != nil {
			loggerconfig.Error("FetchCompanyCategory Error Scan data:", err.Error())
			return dbResponse, err
		}
		dbResponse = append(dbResponse, row)
	}
	return dbResponse, nil
}

func (pgObj *Postgres) FetchDailyAnnouncement(stringOfCoCode string) ([]models.DailyAnnouncement, error) {
	ctx := context.Background()
	var dbResponse []models.DailyAnnouncement
	if len(stringOfCoCode) == 0 {
		loggerconfig.Info("FetchDailyAnnouncement isinList is Empty :")
		return dbResponse, nil
	}
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchDailyAnnouncement Error if database is alive :", err.Error())
		return dbResponse, err
	}
	var queryStatement string
	var res *sql.Rows

	queryStatement = fmt.Sprintf(`select DA.cocode, DA.symbol, DA.lname, DA.caption, DA.date1, DA.memo from dailyannouncement AS DA WHERE DA.cocode IN (%s);`, stringOfCoCode)
	res, err = dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchDailyAnnouncement Error fetching data:", err.Error())
		return dbResponse, err
	}
	defer res.Close()
	var row models.DailyAnnouncement
	for res.Next() {
		err = res.Scan(&row.CoCode, &row.Symbol, &row.CoName, &row.Caption, &row.Date, &row.Memo)
		if err != nil {
			loggerconfig.Error("FetchDailyAnnouncement Error Scan data:", err.Error())
			return dbResponse, err
		}
		loggerconfig.Info("FetchDailyAnnouncement Successful, dbResponse :", dbResponse)
		dbResponse = append(dbResponse, row)
	}
	return dbResponse, nil
}

func (pgObj *Postgres) FetchBoardMeeting(stringOfCoCode string) ([]models.BoardMeetingForthComing, error) {
	ctx := context.Background()
	var dbResponse []models.BoardMeetingForthComing
	if len(stringOfCoCode) == 0 {
		loggerconfig.Error("FetchBoardMeeting isinList is Empty :")
		return dbResponse, nil
	}
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchBoardMeeting Error if database is alive :", err.Error())
		return dbResponse, err
	}
	var queryStatement string
	var res *sql.Rows

	queryStatement = fmt.Sprintf(`SELECT BA.cocode,BA.coname,BA.symbol,BA.date1, BA.note from boardmeetingforthcoming AS BA WHERE BA.cocode IN (%s);`, stringOfCoCode)
	res, err = dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchBoardMeeting Error fetching data:", err.Error())
		return dbResponse, err
	}
	defer res.Close()
	var row models.BoardMeetingForthComing
	for res.Next() {
		err = res.Scan(&row.CoCode, &row.CoName, &row.Symbol, &row.Date, &row.Note)
		if err != nil {
			loggerconfig.Error("FetchBoardMeeting Error Scan data:", err.Error())
			return dbResponse, err
		}
		dbResponse = append(dbResponse, row)
	}
	loggerconfig.Info("FetchBoardMeeting Successful, dbResponse :", dbResponse)
	return dbResponse, nil
}

func (pgObj *Postgres) FetchChangedName(stringOfCoCode string) ([]models.ChangeOfName, error) {
	ctx := context.Background()
	var dbResponse []models.ChangeOfName
	if len(stringOfCoCode) == 0 {
		loggerconfig.Error("FetchChangedName isinList is Empty :")
		return dbResponse, nil
	}
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchChangedName Error if database is alive :", err.Error())
		return dbResponse, err
	}
	var queryStatement string
	var res *sql.Rows

	queryStatement = fmt.Sprintf(`SELECT CN.oldname, CN.newname, CN.cocode, CN.symbol FROM changeofname AS CN WHERE CN.cocode IN (%s);`, stringOfCoCode)
	res, err = dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchChangedName Error fetching data:", err.Error())
		return dbResponse, err
	}
	defer res.Close()
	var row models.ChangeOfName
	for res.Next() {
		err = res.Scan(&row.Oldname, &row.CoName, &row.CoCode, &row.Symbol)
		if err != nil {
			loggerconfig.Error("FetchChangedName Error Scan data:", err.Error())
			return dbResponse, err
		}
		dbResponse = append(dbResponse, row)
	}
	loggerconfig.Info("FetchChangedName Successful, dbResponse :", dbResponse)
	return dbResponse, nil
}

func (pgObj *Postgres) FetchSplits(stringOfCoCode string) ([]models.Splits, error) {
	ctx := context.Background()
	var dbResponse []models.Splits
	if len(stringOfCoCode) == 0 {
		loggerconfig.Error("FetchSplits isinList is Empty :")
		return dbResponse, nil
	}
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchSplits Error if database is alive :", err.Error())
		return dbResponse, err
	}
	var queryStatement string
	var res *sql.Rows

	queryStatement = fmt.Sprintf(`SELECT S.cocode, S.coname, S.symbol, S.splitratio, S.remark, S.splitdate FROM splits AS S WHERE S.cocode IN (%s);`, stringOfCoCode)
	res, err = dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchSplits Error fetching data:", err.Error())
		return dbResponse, err
	}
	defer res.Close()
	var row models.Splits
	for res.Next() {
		err = res.Scan(&row.CoCode, &row.CoName, &row.Symbol, &row.SplitRatio, &row.Remark, &row.SplitDate)
		if err != nil {
			loggerconfig.Error("FetchSplits Error Scan data:", err.Error())
			return dbResponse, err
		}
		dbResponse = append(dbResponse, row)
	}
	loggerconfig.Info("FetchSplits Successful, dbResponse :", dbResponse)
	return dbResponse, nil
}

func (pgObj *Postgres) FetchMerger(stringOfCoCode string) ([]models.MergerDemerger, error) {
	ctx := context.Background()
	var dbResponse []models.MergerDemerger
	if len(stringOfCoCode) == 0 {
		loggerconfig.Error("FetchMerger isinList is Empty :")
		return dbResponse, nil
	}
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchMerger Error if database is alive :", err.Error())
		return dbResponse, err
	}
	var queryStatement string
	var res *sql.Rows

	queryStatement = fmt.Sprintf(`select MD.cocode, MD.coname, MD.mergedintoname, MD.mgrratio, MD.mergerdemergerdate from mergerdemerger AS MD WHERE MD.cocode IN (%s);`, stringOfCoCode)
	res, err = dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchMerger Error fetching data:", err.Error())
		return dbResponse, err
	}
	defer res.Close()
	var row models.MergerDemerger
	for res.Next() {
		err = res.Scan(&row.CoCode, &row.CoName, &row.MergedIntoName, &row.MgrRatio, &row.MergerDemergerDate)
		if err != nil {
			loggerconfig.Error("FetchMerger Error Scan data:", err.Error())
			return dbResponse, err
		}
		dbResponse = append(dbResponse, row)
	}
	loggerconfig.Info("FetchMerger Successful, dbResponse :", dbResponse)
	return dbResponse, nil
}

func (pgObj *Postgres) FetchDividend(stringOfCoCode string) ([]models.DividendAnnouncementData, error) {
	ctx := context.Background()
	var dbResponse []models.DividendAnnouncementData
	if len(stringOfCoCode) == 0 {
		loggerconfig.Error("FetchDividend isinList is Empty :")
		return dbResponse, nil
	}
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchDividend Error if database is alive :", err.Error())
		return dbResponse, err
	}
	var queryStatement string
	var res *sql.Rows

	queryStatement = fmt.Sprintf(`SELECT DA.cocode, DA.coname, DA.symbol, DA.description FROM dividendannouncementdata AS DA WHERE DA.cocode IN (%s);`, stringOfCoCode)
	res, err = dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchDividend Error fetching data:", err.Error())
		return dbResponse, err
	}
	defer res.Close()
	var row models.DividendAnnouncementData
	for res.Next() {
		err = res.Scan(&row.CoCode, &row.CoName, &row.Symbol, &row.Description)
		if err != nil {
			loggerconfig.Error("FetchDividend Error Scan data:", err.Error())
			return dbResponse, err
		}
		dbResponse = append(dbResponse, row)
	}
	loggerconfig.Info("FetchDividend Successful, dbResponse :", dbResponse)
	return dbResponse, nil
}

func (pgObj *Postgres) FetchBulkDeals(stringOfCoCode string) ([]models.BulkDeals, error) {
	ctx := context.Background()
	var dbResponse []models.BulkDeals
	if len(stringOfCoCode) == 0 {
		loggerconfig.Error("FetchBulkDeals isinList is Empty :")
		return dbResponse, nil
	}
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchBulkDeals Error if database is alive :", err.Error())
		return dbResponse, err
	}
	var queryStatement string
	var res *sql.Rows

	queryStatement = fmt.Sprintf(`SELECT cocode, scripname, clientname, qtyshares, avgprice FROM bulkdeals WHERE cocode IN (%s);`, stringOfCoCode)
	res, err = dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchBulkDeals Error fetching data:", err.Error())
		return dbResponse, err
	}
	defer res.Close()
	var bulkDeal models.BulkDeals
	for res.Next() {
		err = res.Scan(&bulkDeal.CoCode, &bulkDeal.Scripname, &bulkDeal.Clientname, &bulkDeal.Qtyshares, &bulkDeal.AvgPrice)
		if err != nil {
			loggerconfig.Error("FetchBulkDeals Error Scan data:", err.Error())
			return dbResponse, err
		}
		dbResponse = append(dbResponse, bulkDeal)
	}
	loggerconfig.Info("FetchBulkDeals Successful, dbResponse :", dbResponse)
	return dbResponse, nil
}

func (pgObj *Postgres) FetchBlockDeals(stringOfCoCode string) ([]models.BlockDeals, error) {
	ctx := context.Background()
	var dbResponse []models.BlockDeals
	if len(stringOfCoCode) == 0 {
		loggerconfig.Error("FetchBlockDeals isinList is Empty :")
		return dbResponse, nil
	}
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchBlockDeals Error if database is alive :", err.Error())
		return dbResponse, err
	}
	var queryStatement string
	var res *sql.Rows

	queryStatement = fmt.Sprintf(`SELECT cocode,scripname, clientname, qtyshares, avgprice FROM blockdeals WHERE cocode IN (%s);`, stringOfCoCode)
	res, err = dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchBlockDeals Error fetching data:", err.Error())
		return dbResponse, err
	}
	defer res.Close()
	var blockDeal models.BlockDeals
	for res.Next() {
		err = res.Scan(&blockDeal.CoCode, &blockDeal.ScripName, &blockDeal.ClientName, &blockDeal.Qtyshares, &blockDeal.AvgPrice)
		if err != nil {
			loggerconfig.Error("FetchBlockDeals Error Scan data:", err.Error())
			return dbResponse, err
		}
		dbResponse = append(dbResponse, blockDeal)
	}
	loggerconfig.Info("FetchBlockDeals Successful, dbResponse :", dbResponse)
	return dbResponse, nil
}

func (pgObj *Postgres) FetchBonus(stringOfCoCode string) ([]models.Bonus, error) {
	ctx := context.Background()
	var dbResponse []models.Bonus
	if len(stringOfCoCode) == 0 {
		loggerconfig.Error("FetchBonus isinList is Empty :")
		return dbResponse, nil
	}
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchBonus Error if database is alive :", err.Error())
		return dbResponse, err
	}
	var queryStatement string
	var res *sql.Rows

	queryStatement = fmt.Sprintf(`SELECT cocode, remark, bonusdate FROM bonus WHERE cocode IN (%s);`, stringOfCoCode)
	res, err = dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchBonus Error fetching data:", err.Error())
		return dbResponse, err
	}
	defer res.Close()
	var row models.Bonus
	for res.Next() {
		err = res.Scan(&row.CoCode, &row.Remark, &row.BonusDate)
		if err != nil {
			loggerconfig.Error("FetchBonus Error Scan data:", err.Error())
			return dbResponse, err
		}
		dbResponse = append(dbResponse, row)
	}
	loggerconfig.Info("FetchBonus Successful, dbResponse :", dbResponse)
	return dbResponse, nil
}

func (pgObj *Postgres) FetchPLStatementData(isin string) (models.PLStatementResponse, error) {
	ctx := context.Background()
	var dbResponse models.PLStatementResponse
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchPLStatementData Error if database is alive :", err.Error())
		return dbResponse, err
	}
	var queryStatement string
	var res *sql.Rows

	queryStatement = `SELECT columnname,y0 FROM plstatement PLS 
	INNER JOIN companymaster CM ON CM.cocode=PLS.cocode 
	WHERE isin=$1 AND typeCS='c' AND columnname IN 
	('Total Revenue', 'Total Expenses', 'Finance Costs', 'Taxation', 'Profit After Tax', 'Earning Per Share - Diluted');`
	res, err = dbops.PostgresRepo.Fetch(queryStatement, isin)
	if err != nil {
		loggerconfig.Error("FetchPLStatementData Error fetching data:", err.Error())
		return dbResponse, err
	}
	defer res.Close()
	var row models.PLStatement
	for res.Next() {
		err = res.Scan(&row.ColumnName, &row.Y0)
		if err != nil {
			loggerconfig.Error("FetchPLStatementData Error Scan data:", err.Error())
			return dbResponse, err
		}
		dbResponse.Data = append(dbResponse.Data, row)
	}
	loggerconfig.Info("FetchPLStatementData Successful, dbResponse :", dbResponse)
	return dbResponse, nil
}

func (pgObj *Postgres) FetchBalanceSheetsData(isin string) (models.BalanceSheetsResponse, error) {
	ctx := context.Background()
	var dbResponse models.BalanceSheetsResponse
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchBalanceSheetsData Error if database is alive :", err.Error())
		return dbResponse, err
	}
	var queryStatement string
	var res *sql.Rows

	queryStatement = `SELECT columnname,y0 FROM balancesheets BS
	INNER JOIN companymaster CM ON CM.cocode=BS.cocode 
	WHERE isin=$1 AND typecs='C'
	AND ((rid, rowno) IN (VALUES (68,73), (73,78),(1,2)) OR columnname IN ('TOTAL ASSETS','Total Current Liabilities','Total Current Assets','TOTAL EQUITY AND LIABILITIES'));`
	res, err = dbops.PostgresRepo.Fetch(queryStatement, isin)
	if err != nil {
		loggerconfig.Error("FetchBalanceSheetsData Error fetching data:", err.Error())
		return dbResponse, err
	}
	defer res.Close()

	var row models.BalanceSheets
	var otherLiabilities float64
	var otherAssets float64
	for res.Next() {
		err = res.Scan(&row.ColumnName, &row.Y0)
		if err != nil {
			loggerconfig.Error("FetchBalanceSheetsData Error Scan data:", err.Error())
			return dbResponse, err
		}

		if row.ColumnName == constants.TotalEquityAndLiabilities {
			otherLiabilities += row.Y0
		} else if row.ColumnName == constants.ShareCapital || row.ColumnName == constants.TotalCurrentLiabilities || row.ColumnName == constants.ReservesAndSurplus {
			otherLiabilities -= row.Y0
		} else if row.ColumnName == constants.TotalAssets {
			otherAssets += row.Y0
		} else {
			otherAssets -= row.Y0
		}

		if row.ColumnName == constants.ReservesAndSurplus {
			row.ColumnName = row.ColumnName[constants.THREE:]
		} else if row.ColumnName == constants.ShareCapital {
			row.ColumnName = row.ColumnName[:constants.THIRTEEN]
		}
		dbResponse.Data = append(dbResponse.Data, row)
	}
	var InsertOtherAssets, InsertOtherLiabiities models.BalanceSheets

	InsertOtherAssets.ColumnName = constants.OtherAssets
	InsertOtherAssets.Y0 = otherAssets
	dbResponse.Data = append(dbResponse.Data, InsertOtherAssets)

	InsertOtherLiabiities.ColumnName = constants.OtherLiabiities
	InsertOtherLiabiities.Y0 = otherLiabilities
	dbResponse.Data = append(dbResponse.Data, InsertOtherLiabiities)

	loggerconfig.Info("FetchBalanceSheetsData Successful, dbResponse :", dbResponse)
	return dbResponse, nil
}

func (pgObj *Postgres) FetchCashFlowData(isin string) (models.CashflowResponse, error) {
	ctx := context.Background()
	var dbResponse models.CashflowResponse
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchCashFlowData  Error if database is alive :", err.Error())
		return dbResponse, err
	}

	var wg sync.WaitGroup

	totalCurrCallInFetchCashFlowData := 4
	wg.Add(totalCurrCallInFetchCashFlowData)

	// CFO
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		var queryStatement string
		var res *sql.Rows

		queryStatement = `SELECT columnname,y0 FROM cashflow CF
		INNER JOIN companymaster CM ON CM.cocode=CF.cocode 
		WHERE isin=$1 AND typecs='c'
		AND (rid, rowno) IN (VALUES (40,47), (14,19),(15,20),(16,21),(28,34));`
		res, err = dbops.PostgresRepo.Fetch(queryStatement, isin)
		if err != nil {
			loggerconfig.Error("FetchCashFlowData  Error fetching data:", err.Error())
		} else {
			defer res.Close()
			var row models.CashFlow
			for res.Next() {
				err = res.Scan(&row.ColumnName, &row.Y0)
				if err != nil {
					loggerconfig.Error("FetchCashFlowData  Error Scan data:", err.Error())
				} else {
					dbResponse.CFO = append(dbResponse.CFO, row)
				}
			}
		}

	}(&wg)

	//CFI
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		var queryStatement1 string
		var res1 *sql.Rows

		queryStatement1 = `SELECT columnname,y0 FROM cashflow CF
		INNER JOIN companymaster CM ON CM.cocode=CF.cocode 
		WHERE isin=$1 AND typecs='c'
		AND (rid, rowno) IN (VALUES (28,34), (43,54),(46,58),(47,59),(49,61),(5,9),(0,49));`
		res1, err = dbops.PostgresRepo.Fetch(queryStatement1, isin)
		if err != nil {
			loggerconfig.Error("FetchCashFlowData  Error fetching data:", err.Error())
		} else {
			defer res1.Close()
			var row1 models.CashFlow
			var cfiData models.CashFlow
			cfiData.ColumnName = constants.OtherInvestingItems
			for res1.Next() {
				err = res1.Scan(&row1.ColumnName, &row1.Y0)
				if err != nil {
					loggerconfig.Error("FetchCashFlowData  Error Scan data:", err.Error())
				} else {
					if row1.ColumnName == constants.CFI {
						cfiData.Y0 += row1.Y0
					} else {
						dbResponse.CFI = append(dbResponse.CFI, row1)
						cfiData.Y0 -= row1.Y0
					}
				}
			}
		}
	}(&wg)

	// CFF
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		var queryStatement2 string
		var res2 *sql.Rows

		queryStatement2 = `SELECT columnname,y0 FROM cashflow CF
		INNER JOIN companymaster CM ON CM.cocode=CF.cocode 
		WHERE isin=$1 AND typecs='c'
	    AND (rid, rowno) IN (VALUES (60,76), (62,78),(71,88),(72,89),(76,93),(74,91),(73,90),(0,73));`
		res2, err = dbops.PostgresRepo.Fetch(queryStatement2, isin)
		if err != nil {
			loggerconfig.Error("FetchCashFlowData  Error fetching data:", err.Error())
		} else {
			defer res2.Close()
			var row2 models.CashFlow
			var cffData models.CashFlow
			var repaymentOfBorrowings models.CashFlow
			cffData.ColumnName = constants.OtherInvestingItems
			repaymentOfBorrowings.ColumnName = constants.RepaymentOfBorrowings
			for res2.Next() {
				err = res2.Scan(&row2.ColumnName, &row2.Y0)
				if err != nil {
					loggerconfig.Error("FetchCashFlowData  Error Scan data:", err.Error())
				} else {
					if row2.ColumnName == constants.CFF {
						cffData.Y0 += row2.Y0
					} else if row2.ColumnName != constants.OfLongTermBorrowing && row2.ColumnName != constants.OfShortTermBorrowing {
						dbResponse.CFF = append(dbResponse.CFF, row2)
						cffData.Y0 -= row2.Y0
					} else if row2.ColumnName == constants.OfLongTermBorrowing {
						repaymentOfBorrowings.Y0 += row2.Y0
					} else {
						repaymentOfBorrowings.Y0 -= row2.Y0
					}
				}
			}
			cffData.Y0 -= repaymentOfBorrowings.Y0
			if len(dbResponse.CFF) != constants.ZERO {
				dbResponse.CFF = append(dbResponse.CFF, repaymentOfBorrowings)
				dbResponse.CFF = append(dbResponse.CFF, cffData)
			}
		}

	}(&wg)

	// Netcashflow
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		var queryStatement3 string
		var res3 *sql.Rows

		queryStatement3 = `SELECT columnname,y0 FROM cashflow CF
		INNER JOIN companymaster CM ON CM.cocode=CF.cocode 
		WHERE isin=$1 AND typecs='c'
	    AND (rid, rowno) IN (VALUES (79,98));`
		res3, err = dbops.PostgresRepo.Fetch(queryStatement3, isin)
		if err != nil {
			loggerconfig.Error("FetchCashFlowData  Error fetching data:", err.Error())
		} else {
			defer res3.Close()
			var row3 models.CashFlow
			for res3.Next() {
				err = res3.Scan(&row3.ColumnName, &row3.Y0)
				if err != nil {
					loggerconfig.Error("FetchCashFlowData  Error Scan data:", err.Error())
				} else {
					dbResponse.NetCashFlow = row3
				}
			}
		}

	}(&wg)

	wg.Wait()

	loggerconfig.Info("FetchCashFlowData  Successful, dbResponse :", dbResponse)
	return dbResponse, nil
}

func (pgObj *Postgres) FetchCompanyMasterData(allIsin models.AllIsin) (models.CompanyMasterDb, error) {
	var allCompanyMasterData models.CompanyMasterDb
	ctx := context.Background()
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchCompanyMasterData Error if database is alive :", err.Error())
		return allCompanyMasterData, err
	}

	if allIsin.Isin == nil || len(allIsin.Isin) == 0 {
		return allCompanyMasterData, fmt.Errorf("FetchCompanyMasterData: input ISIN list is empty or nil")
	}

	stringOfIsin := ""
	for i := 0; i < len(allIsin.Isin); i++ {
		stringOfIsin += "'" + allIsin.Isin[i] + "',"
	}
	if len(stringOfIsin) > 0 {
		stringOfIsin = stringOfIsin[:len(stringOfIsin)-1]
	}

	queryStatement := fmt.Sprintf(`SELECT CM.bsecode, CM.nsesymbol, CM.isin, CM.sectorcode, CM.sectorname
	FROM companymaster AS CM
	WHERE CM.isin IN (%s);`, stringOfIsin)
	res, err := dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchCompanyMasterData Error just client fetching data:", err.Error())
		return allCompanyMasterData, err
	}
	defer res.Close()

	var companyMasterDbData models.CompanyMasterDbData
	for res.Next() {
		err = res.Scan(&companyMasterDbData.Bsecode, &companyMasterDbData.Nsesymbol, &companyMasterDbData.Isin, &companyMasterDbData.Sectorcode, &companyMasterDbData.Sectorname)
		if err != nil {
			loggerconfig.Error("FetchCompanyMasterData Error in scanning fetched companymaster response from postgres :", err)
			break
		}
		allCompanyMasterData.CompanyMasterAll = append(allCompanyMasterData.CompanyMasterAll, companyMasterDbData)
	}

	return allCompanyMasterData, nil
}

func (pgObj *Postgres) GetPostgresStatus() error {

	ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		return err
	}

	return nil

}

func (pgObj *Postgres) FetchDeclineInPromoterHoldingData(allIsin models.AllIsin) (models.AllDeclineInPromoterHoldingDb, error) {
	var allDeclineInPromoterHoldingDb models.AllDeclineInPromoterHoldingDb
	ctx := context.Background()
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchDeclineInPromoterHoldingData Error if database is alive :", err.Error())
		return allDeclineInPromoterHoldingDb, err
	}

	if allIsin.Isin == nil || len(allIsin.Isin) == 0 {
		return allDeclineInPromoterHoldingDb, fmt.Errorf("FetchDeclineInPromoterHoldingData: input ISIN list is empty or nil")
	}

	stringOfIsin := ""
	for i := 0; i < len(allIsin.Isin); i++ {
		if len(allIsin.Isin[i]) == constants.ZERO {
			continue
		}
		stringOfIsin += "'" + allIsin.Isin[i] + "',"
	}
	if len(stringOfIsin) != 0 {
		stringOfIsin = stringOfIsin[:len(stringOfIsin)-1]
	}

	queryStatement := fmt.Sprintf(`SELECT cm.isin,totalpromoterpershares
	FROM shareholdingpatterndetails spd
	JOIN companymaster cm ON spd.cocode = cm.cocode
	WHERE cm.isin IN (%s)
	AND spd.yrc IN (
	  SELECT MAX(yrc)
	  FROM shareholdingpatterndetails
	  WHERE cocode = spd.cocode
	  UNION
	  SELECT MAX(yrc)
	  FROM shareholdingpatterndetails
	  WHERE cocode = spd.cocode
	  AND yrc < (
		SELECT MAX(yrc)
		FROM shareholdingpatterndetails
		WHERE cocode = spd.cocode
	  )
	) ORDER BY cm.cocode,spd.yrc DESC;`, stringOfIsin)
	res, err := dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchDeclineInPromoterHoldingData Error just client fetching data:", err.Error())
		return allDeclineInPromoterHoldingDb, err
	}
	defer res.Close()

	var declineInPromoterHoldingDb models.DeclineInPromoterHoldingDb
	var quarter = constants.CurrentQuarter
	for res.Next() {

		if quarter == constants.CurrentQuarter {
			err = res.Scan(&declineInPromoterHoldingDb.Isin, &declineInPromoterHoldingDb.CurrentQuarterTPPS)
			if err != nil {
				loggerconfig.Error("FetchDeclineInPromoterHoldingData Error in scanning fetched companymaster response from postgres :", err)
				break
			}
			quarter = constants.PreviousQuarter
		} else {
			err = res.Scan(&declineInPromoterHoldingDb.Isin, &declineInPromoterHoldingDb.PreviousQuarterTPPS)
			if err != nil {
				loggerconfig.Error("FetchDeclineInPromoterHoldingData Error in scanning fetched companymaster response from postgres :", err)
				break
			}
			allDeclineInPromoterHoldingDb.DeclineInPromoterHolding = append(allDeclineInPromoterHoldingDb.DeclineInPromoterHolding, declineInPromoterHoldingDb)
			quarter = constants.CurrentQuarter
		}
	}
	loggerconfig.Info("FetchDeclineInPromoterHoldingData SuccessFul, dbResponse", allDeclineInPromoterHoldingDb)
	return allDeclineInPromoterHoldingDb, nil
}

func (pgObj *Postgres) FetchInterestCoverageRatioData(allIsin models.AllIsin) (models.AllInterestCoverageRatioDb, error) {
	var allInterestCoverageRatioDb models.AllInterestCoverageRatioDb
	ctx := context.Background()
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchInterestCoverageRatioData Error if database is alive :", err.Error())
		return allInterestCoverageRatioDb, err
	}

	if allIsin.Isin == nil || len(allIsin.Isin) == 0 {
		return allInterestCoverageRatioDb, fmt.Errorf("FetchInterestCoverageRatioData: input ISIN list is empty or nil")
	}

	stringOfIsin := ""
	for i := 0; i < len(allIsin.Isin); i++ {
		if len(allIsin.Isin[i]) == constants.ZERO {
			continue
		}
		stringOfIsin += "'" + allIsin.Isin[i] + "',"
	}
	if len(stringOfIsin) != 0 {
		stringOfIsin = stringOfIsin[:len(stringOfIsin)-1]
	}

	queryStatement := fmt.Sprintf(`SELECT cm.isin, pls.y0 ,interestcoverratio FROM plstatement pls
	INNER JOIN companymaster cm ON cm.cocode=pls.cocode
	INNER JOIN industryratios ir ON ir.indcode=cm.industrycode
	WHERE cm.isin IN 
		(SELECT cm.isin
		FROM plstatement pls
		INNER JOIN companymaster cm ON cm.cocode=pls.cocode
		WHERE pls.columnname IN ('Finance Costs','Profit Before Tax') 
		AND cm.isin 
		IN (%s)
		GROUP BY cm.isin
		HAVING COUNT(DISTINCT columnname) = 2) 
	AND typecs='c' AND columnname IN ('Finance Costs','Profit Before Tax')
	AND ir.yrc=(SELECT MAX(yrc) FROM industryratios);`, stringOfIsin)
	res, err := dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchInterestCoverageRatioData Error just client fetching data:", err.Error())
		return allInterestCoverageRatioDb, err
	}
	defer res.Close()

	var interestCoverageRatioDb models.InterestCoverageRatioDb
	var column = constants.ONE
	for res.Next() {

		if column == constants.ONE {
			err = res.Scan(&interestCoverageRatioDb.Isin, &interestCoverageRatioDb.FinanceCost, &interestCoverageRatioDb.InterestCoverRatio)
			if err != nil {
				loggerconfig.Error("FetchInterestCoverageRatioData Error in scanning fetched companymaster response from postgres :", err)
				break
			}
			column = constants.ZERO
		} else {
			err = res.Scan(&interestCoverageRatioDb.Isin, &interestCoverageRatioDb.ProfitBeforeTax, &interestCoverageRatioDb.InterestCoverRatio)
			if err != nil {
				loggerconfig.Error("FetchDeclineInPromoterHoldingData Error in scanning fetched companymaster response from postgres :", err)
				break
			}
			allInterestCoverageRatioDb.InterestCoverageRatioData = append(allInterestCoverageRatioDb.InterestCoverageRatioData, interestCoverageRatioDb)
			column = constants.ONE
		}
	}
	loggerconfig.Info("FetchInterestCoverageRatioData SuccessFul, dbResponse", allInterestCoverageRatioDb)
	return allInterestCoverageRatioDb, nil
}

func (pgObj *Postgres) DeclineInRevenueAndProfitData(allIsin models.AllIsin) (models.AllDeclineInRevenueAndProfitDb, error) {
	var allDeclineInRevenueAndProfitDb models.AllDeclineInRevenueAndProfitDb
	ctx := context.Background()
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("DeclineInRevenueAndProfitData Error if database is alive :", err.Error())
		return allDeclineInRevenueAndProfitDb, err
	}

	if allIsin.Isin == nil || len(allIsin.Isin) == 0 {
		return allDeclineInRevenueAndProfitDb, fmt.Errorf("DeclineInRevenueAndProfitData: input ISIN list is empty or nil")
	}

	stringOfIsin := ""
	for i := 0; i < len(allIsin.Isin); i++ {
		if len(allIsin.Isin[i]) == constants.ZERO {
			continue
		}
		stringOfIsin += "'" + allIsin.Isin[i] + "',"
	}
	if len(stringOfIsin) != 0 {
		stringOfIsin = stringOfIsin[:len(stringOfIsin)-1]
	}

	queryStatement := fmt.Sprintf(`SELECT cm.isin,pls.y0,pls.y1,pls.y2 FROM plstatement pls
	INNER JOIN companymaster cm ON cm.cocode=pls.cocode
		WHERE cm.isin IN (SELECT cm.isin
		FROM plstatement pls
		INNER JOIN companymaster cm on cm.cocode=pls.cocode
		WHERE pls.columnname IN ('Total Revenue','Profit After Tax') 
		AND cm.isin 
		IN (%s)
		GROUP BY cm.isin
		HAVING COUNT(DISTINCT columnname) = 2) 
	AND typecs='c' AND columnname IN ('Total Revenue','Profit After Tax');`, stringOfIsin)
	res, err := dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("DeclineInRevenueAndProfitData Error just client fetching data:", err.Error())
		return allDeclineInRevenueAndProfitDb, err
	}
	defer res.Close()

	var column = constants.ONE
	var numberOfConditionSatified = constants.ZERO
	for res.Next() {
		var revenueAndProfitData models.RevenueAndProfitData
		err = res.Scan(&revenueAndProfitData.Isin, &revenueAndProfitData.Y0, &revenueAndProfitData.Y1, &revenueAndProfitData.Y2)
		if err != nil {
			loggerconfig.Error("DeclineInRevenueAndProfitData Error in scanning fetched plStatement response from postgres :", err)
			continue
		}
		if column == constants.ONE {
			if revenueAndProfitData.Y0 < revenueAndProfitData.Y1 && revenueAndProfitData.Y1 < revenueAndProfitData.Y2 {
				numberOfConditionSatified = constants.ONE
			}
			column = constants.ZERO
		} else {
			if numberOfConditionSatified == constants.ONE && revenueAndProfitData.Y0 < revenueAndProfitData.Y1 && revenueAndProfitData.Y1 < revenueAndProfitData.Y2 {
				var declineInRevenueAndProfitDb models.DeclineInRevenueAndProfitDb
				declineInRevenueAndProfitDb.Isin = revenueAndProfitData.Isin
				allDeclineInRevenueAndProfitDb.Holding = append(allDeclineInRevenueAndProfitDb.Holding, declineInRevenueAndProfitDb)
			}
			numberOfConditionSatified = constants.ZERO
			column = constants.ONE
		}
	}
	loggerconfig.Info("FetchInterestCoverageRatioData SuccessFul, dbResponse", allDeclineInRevenueAndProfitDb)
	return allDeclineInRevenueAndProfitDb, nil
}

func (pgObj *Postgres) LowNetWorthData(allIsin models.AllIsin) (models.AllNetWorthDb, error) {
	var allNetWorthData models.AllNetWorthDb
	ctx := context.Background()
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("LowNetWorthData Error if database is alive :", err.Error())
		return allNetWorthData, err
	}

	if allIsin.Isin == nil || len(allIsin.Isin) == 0 {
		return allNetWorthData, fmt.Errorf("LowNetWorthData: input ISIN list is empty or nil")
	}

	stringOfIsin := ""
	for i := 0; i < len(allIsin.Isin); i++ {
		if len(allIsin.Isin[i]) == constants.ZERO {
			continue
		}
		stringOfIsin += "'" + allIsin.Isin[i] + "',"
	}
	if len(stringOfIsin) != 0 {
		stringOfIsin = stringOfIsin[:len(stringOfIsin)-1]
	}

	queryStatement := fmt.Sprintf(`SELECT subquery.isin, SUM(subquery.y0) AS Networth 
	FROM (
	  SELECT cm.isin, BS.columnname, BS.y0
	  FROM balancesheets BS
	  INNER JOIN companymaster cm ON cm.cocode = BS.cocode
	  WHERE cm.isin IN (%s)
		AND BS.columnname IN ('TOTAL ASSETS', 'Total Current Liabilities', 'Total Non Current Liabilities')
		AND BS.typecs = 'C'
	) AS subquery
	GROUP BY subquery.isin;`, stringOfIsin)
	res, err := dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("LowNetWorthData Error just client fetching data:", err.Error())
		return allNetWorthData, err
	}
	defer res.Close()

	for res.Next() {
		var netWorthData models.NetWorthDb
		err = res.Scan(&netWorthData.Isin, &netWorthData.NetWorth)
		if err != nil {
			loggerconfig.Error("LowNetWorthData Error in scanning fetched plStatement response from postgres :", err)
			break
		}
		if netWorthData.NetWorth < 0 {
			allNetWorthData.NetWorthData = append(allNetWorthData.NetWorthData, netWorthData)
		}
	}
	loggerconfig.Info("LowNetWorthData SuccessFul, dbResponse", allNetWorthData)
	return allNetWorthData, nil
}

func (pgObj *Postgres) DeclineInRevenueData(allIsin models.AllIsin) (models.AllDeclineInRevenueDb, error) {
	var allDeclineInRevenueDb models.AllDeclineInRevenueDb
	ctx := context.Background()
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("DeclineInRevenueData Error if database is alive :", err.Error())
		return allDeclineInRevenueDb, err
	}

	if allIsin.Isin == nil || len(allIsin.Isin) == 0 {
		return allDeclineInRevenueDb, fmt.Errorf("DeclineInRevenueData: input ISIN list is empty or nil")
	}
	stringOfIsin := ""
	for i := 0; i < len(allIsin.Isin); i++ {
		if len(allIsin.Isin[i]) == constants.ZERO {
			continue
		}
		stringOfIsin += "'" + allIsin.Isin[i] + "',"
	}
	if len(stringOfIsin) != 0 {
		stringOfIsin = stringOfIsin[:len(stringOfIsin)-1]
	}

	queryStatement := fmt.Sprintf(`SELECT cm.isin,pls.y0,pls.y1,pls.y2,pls.y3,pls.y4 FROM plstatement pls
	INNER JOIN companymaster cm ON cm.cocode=pls.cocode
	WHERE typecs='c' AND columnname = 'Total Revenue'
	AND	cm.isin IN (%s);`, stringOfIsin)
	res, err := dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("DeclineInRevenueData Error just client fetching data:", err.Error())
		return allDeclineInRevenueDb, err
	}
	defer res.Close()

	for res.Next() {
		var revenueData models.RevenueData
		err = res.Scan(&revenueData.Isin, &revenueData.Y0, &revenueData.Y1, &revenueData.Y2, &revenueData.Y3, &revenueData.Y4)
		if err != nil {
			loggerconfig.Error("DeclineInRevenueData Error in scanning fetched plStatement response from postgres :", err)
			continue
		}
		if revenueData.Y0 < revenueData.Y1 && revenueData.Y1 < revenueData.Y2 && revenueData.Y2 < revenueData.Y3 && revenueData.Y3 < revenueData.Y4 {
			var declineInRevenueDb models.DeclineInRevenueDb
			declineInRevenueDb.Isin = revenueData.Isin
			allDeclineInRevenueDb.Holding = append(allDeclineInRevenueDb.Holding, declineInRevenueDb)
		}
	}
	loggerconfig.Info("DeclineInRevenueData SuccessFul, dbResponse", allDeclineInRevenueDb)
	return allDeclineInRevenueDb, nil
}

func (pgObj *Postgres) PromoterPledgeData(allIsin models.AllIsin) (models.AllPromoterPledgeDataDb, error) {
	var allPromoterPledgeDataDb models.AllPromoterPledgeDataDb
	ctx := context.Background()
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("PromoterPledgeData Error if database is alive :", err.Error())
		return allPromoterPledgeDataDb, err
	}

	if allIsin.Isin == nil || len(allIsin.Isin) == 0 {
		return allPromoterPledgeDataDb, fmt.Errorf("PromoterPledgeData: input ISIN list is empty or nil")
	}
	stringOfIsin := ""
	for i := 0; i < len(allIsin.Isin); i++ {
		if len(allIsin.Isin[i]) == constants.ZERO {
			continue
		}
		stringOfIsin += "'" + allIsin.Isin[i] + "',"
	}
	if len(stringOfIsin) != 0 {
		stringOfIsin = stringOfIsin[:len(stringOfIsin)-1]
	}

	queryStatement := fmt.Sprintf(`SELECT cm.isin, spd.totalpromoterperpledgeshares
	FROM companymaster cm
	INNER JOIN shareholdingpatterndetails spd ON cm.cocode = spd.cocode
	INNER JOIN (
		SELECT cocode, MAX(yrc) AS max_yrc
		FROM shareholdingpatterndetails
		GROUP BY cocode
	) max_yrc_table ON spd.cocode = max_yrc_table.cocode AND spd.yrc = max_yrc_table.max_yrc
	WHERE cm.isin IN (%s);`, stringOfIsin)
	res, err := dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("PromoterPledgeData Error just client fetching data:", err.Error())
		return allPromoterPledgeDataDb, err
	}
	defer res.Close()

	for res.Next() {
		var promoterPledgeDataData models.PromoterPledgeData
		err = res.Scan(&promoterPledgeDataData.Isin, &promoterPledgeDataData.TotalPromoterPerPledgeShares)
		if err != nil {
			loggerconfig.Error("PromoterPledgeData Error in scanning fetched plStatement response from postgres :", err)
			continue
		}
		if promoterPledgeDataData.TotalPromoterPerPledgeShares > constants.FIFTY {
			var promoterPledgeDataDb models.PromoterPledgeDataDb
			promoterPledgeDataDb.Isin = promoterPledgeDataData.Isin
			allPromoterPledgeDataDb.Holding = append(allPromoterPledgeDataDb.Holding, promoterPledgeDataDb)
		}
	}
	loggerconfig.Info("PromoterPledgeData SuccessFul, dbResponse", allPromoterPledgeDataDb)
	return allPromoterPledgeDataDb, nil
}

func (pgObj *Postgres) PennyStocksData(allIsin models.AllIsin) (models.AllPennyStocksDataDb, error) {
	var allPennyStocksDataDb models.AllPennyStocksDataDb
	ctx := context.Background()
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("PennyStocksData Error if database is alive :", err.Error())
		return allPennyStocksDataDb, err
	}

	if allIsin.Isin == nil || len(allIsin.Isin) == 0 {
		return allPennyStocksDataDb, fmt.Errorf("PennyStocksData: input ISIN list is empty or nil")
	}
	stringOfIsin := ""
	for i := 0; i < len(allIsin.Isin); i++ {
		if len(allIsin.Isin[i]) == constants.ZERO {
			continue
		}
		stringOfIsin += "'" + allIsin.Isin[i] + "',"
	}
	if len(stringOfIsin) != 0 {
		stringOfIsin = stringOfIsin[:len(stringOfIsin)-1]
	}

	queryStatement := fmt.Sprintf(`SELECT cpkp.isin, marketcap FROM CompanyPeerKeyParams cpkp
	INNER JOIN companymaster cm ON cm.cocode=cpkp.cocode
	WHERE cpkp.isin IN (%s)
	AND typecs='c';`, stringOfIsin)
	res, err := dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("PennyStocksData Error just client fetching data:", err.Error())
		return allPennyStocksDataDb, err
	}
	defer res.Close()

	for res.Next() {
		var pennyStocksData models.PennyStocksData
		err = res.Scan(&pennyStocksData.Isin, &pennyStocksData.MarketCap)
		if err != nil {
			loggerconfig.Error("PennyStocksData Error in scanning fetched plStatement response from postgres :", err)
			continue
		}
		if pennyStocksData.MarketCap < constants.FIFTY {
			allPennyStocksDataDb.Holding = append(allPennyStocksDataDb.Holding, pennyStocksData)
		}
	}
	loggerconfig.Info("PennyStocksData SuccessFul, dbResponse", allPennyStocksDataDb)
	return allPennyStocksDataDb, nil
}

func (pgObj *Postgres) StockReturnData(allIsin models.AllIsin) (models.AllStockReturn, error) {
	var allStockReturn models.AllStockReturn
	ctx := context.Background()
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("StockReturnData Error if database is alive :", err.Error())
		return allStockReturn, err
	}

	queryStatementNifyIndex := `SELECT * FROM niftyindex;`
	resNiftyIndex, err := dbops.PostgresRepo.Fetch(queryStatementNifyIndex)
	if err != nil {
		loggerconfig.Error("StockReturnData error in fetching nifty index data:", err.Error())
		return allStockReturn, err
	}
	defer resNiftyIndex.Close()
	smallCapReturn := 0.0
	midCapReturn := 0.0
	largeCapReturn := 0.0
	for resNiftyIndex.Next() {
		var niftyIndexData models.NiftyIndexData
		err = resNiftyIndex.Scan(&niftyIndexData.TokenId, &niftyIndexData.MarketCap, &niftyIndexData.MarketCapReturn)
		if err != nil {
			continue
		}
		if niftyIndexData.MarketCap == constants.NSESMALLCAP {
			smallCapReturn = niftyIndexData.MarketCapReturn
		} else if niftyIndexData.MarketCap == constants.NSEMIDCAP {
			midCapReturn = niftyIndexData.MarketCapReturn
		} else if niftyIndexData.MarketCap == constants.NSELARGECAP {
			largeCapReturn = niftyIndexData.MarketCapReturn
		}
	}

	if allIsin.Isin == nil || len(allIsin.Isin) == 0 {
		return allStockReturn, fmt.Errorf("StockReturnData: input ISIN list is empty or nil")
	}
	stringOfIsin := ""
	for i := 0; i < len(allIsin.Isin); i++ {
		if len(allIsin.Isin[i]) == constants.ZERO {
			continue
		}
		stringOfIsin += "'" + allIsin.Isin[i] + "',"
	}
	if len(stringOfIsin) != 0 {
		stringOfIsin = stringOfIsin[:len(stringOfIsin)-1]
	}

	queryStatement := fmt.Sprintf(`SELECT isin, mcaptype, returnrate FROM stockreturnrate
	WHERE isin IN (%s);`, stringOfIsin)
	res, err := dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("StockReturnData Error in fetching stockreturnrate data:", err.Error())
		return allStockReturn, err
	}
	defer res.Close()

	for res.Next() {
		var stockReturn models.StockReturn
		err = res.Scan(&stockReturn.Isin, &stockReturn.McapType, &stockReturn.ReturnRate)
		if err != nil {
			loggerconfig.Error("StockReturnData Error in scanning fetched stockreturnrate response from postgres :", err)
			continue
		}
		if len(stockReturn.Isin) == 0 {
			continue
		}
		if (stockReturn.McapType == constants.SMALLMCAPTYPE && stockReturn.ReturnRate < smallCapReturn) || (stockReturn.McapType == constants.MIDMCAPTYPE && stockReturn.ReturnRate < midCapReturn) || (stockReturn.McapType == constants.LARGEMCAPTYPE && stockReturn.ReturnRate < largeCapReturn) {
			allStockReturn.StockReturn = append(allStockReturn.StockReturn, stockReturn)
		}
	}
	loggerconfig.Info("StockReturnData SuccessFul with low return", allStockReturn)
	return allStockReturn, nil
}

func FetchFinancialsDataV2Query(searchBy string) string {
	queryStatement := fmt.Sprintf(`SELECT cocode, columnname, y0, y1, y2, y3, y4
	FROM (
	  SELECT cm.cocode, columnname, y0, y1, y2, y3, y4
	  FROM plstatement pl
	  INNER JOIN companymaster cm ON pl.cocode = cm.cocode
	  WHERE cm.%s =$1
	  AND typecs = 'c' AND columnname IN ('Profit After Tax', 'Total Revenue')
	
	  UNION ALL
	
	  SELECT cm.cocode, columnname, y0,  y1,  y2,  y3,  y4
	  FROM balancesheets bs
	  INNER JOIN companymaster cm ON bs.cocode = cm.cocode
	  WHERE cm.%s =$2
	  AND typecs = 'C' AND columnname IN ('TOTAL ASSETS', 'Total Non Current Liabilities','Total Current Liabilities')
	
	  UNION ALL
	
	  SELECT cm.cocode, columnname,  y0,  y1,  y2,  y3,  y4
	  FROM cashflow cf
	  INNER JOIN companymaster cm ON cf.cocode = cm.cocode
	  WHERE cm.%s =$3
	  AND typecs = 'c' AND columnname = 'Net Inc./(Dec.) in Cash and Cash Equivalent'
	) AS subquery
	ORDER BY cocode, columnname;`, searchBy, searchBy, searchBy)
	return queryStatement
}

func (pgObj *Postgres) FetchFinancialsDataV2(req models.FetchFinancialsReq) (models.FetchFinancialsV2Res, error) {
	ctx := context.Background()
	var dbResponse models.FetchFinancialsV2Res
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchFinancialsDataV2 Error if database is alive :", err.Error())
		return dbResponse, err
	}

	if req.Isin != "" {
		loggerconfig.Info("FetchFinancialsDataV2 isin =", req.Isin)
		queryStatement := FetchFinancialsDataV2Query(constants.ISIN)
		res, err := dbops.PostgresRepo.Fetch(queryStatement, req.Isin, req.Isin, req.Isin)
		if err != nil {
			loggerconfig.Error("FetchFinancialsDataV2 Error just client fetching data:", err.Error())
			return dbResponse, err
		}
		defer res.Close()

		for res.Next() {
			var row models.FetchFinancialsDataV2
			err = res.Scan(&row.CoCode, &row.ColumnName, &row.Y0, &row.Y1, &row.Y2, &row.Y3, &row.Y4)
			if err != nil {
				loggerconfig.Error("FetchFinancialsDataV2 Error in scanning fetched response from postgres :", err)
				continue
			}

			if row.ColumnName == "Profit After Tax" {
				dbResponse.NetProfit = row

			} else if row.ColumnName == "Total Revenue" {
				dbResponse.Revenue = row

			} else if row.ColumnName == "TOTAL ASSETS" {
				dbResponse.BalanceSheet.TotalAssets = row

			} else if row.ColumnName == "Total Current Liabilities" || row.ColumnName == "Total Non Current Liabilities" {
				dbResponse.BalanceSheet.TotalLiabilities.CoCode = row.CoCode
				dbResponse.BalanceSheet.TotalLiabilities.Y0 = row.Y0
				dbResponse.BalanceSheet.TotalLiabilities.Y1 = row.Y1
				dbResponse.BalanceSheet.TotalLiabilities.Y2 = row.Y2
				dbResponse.BalanceSheet.TotalLiabilities.Y3 = row.Y3
				dbResponse.BalanceSheet.TotalLiabilities.Y4 = row.Y4

			} else if row.ColumnName == "Net Inc./(Dec.) in Cash and Cash Equivalent" {
				dbResponse.Cashflow = row

			}
		}
		dbResponse.BalanceSheet.TotalLiabilities.ColumnName = "Total Liabilities"

	} else if req.Exchange == "" {
		loggerconfig.Error("FetchFinancialsData Error no exchange entered in request!")
		return dbResponse, err
	} else if strings.ToLower(req.Exchange) == constants.BSE && req.BseToken != "" {
		loggerconfig.Info("FetchFinancialsDataV2 BseToken =", req.BseToken)
		queryStatement := FetchFinancialsDataV2Query(constants.BSECode)
		res, err := dbops.PostgresRepo.Fetch(queryStatement, req.BseToken, req.BseToken, req.BseToken)
		if err != nil {
			loggerconfig.Error("FetchFinancialsDataV2 Error just client fetching data:", err.Error())
			return dbResponse, err
		}
		defer res.Close()

		for res.Next() {
			var row models.FetchFinancialsDataV2
			err = res.Scan(&row.CoCode, &row.ColumnName, &row.Y0, &row.Y1, &row.Y2, &row.Y3, &row.Y4)
			if err != nil {
				loggerconfig.Error("FetchFinancialsDataV2 Error in scanning fetched response from postgres :", err)
				continue
			}

			if row.ColumnName == "Profit After Tax" {
				dbResponse.NetProfit = row

			} else if row.ColumnName == "Total Revenue" {
				dbResponse.Revenue = row

			} else if row.ColumnName == "TOTAL ASSETS" {
				dbResponse.BalanceSheet.TotalAssets = row

			} else if row.ColumnName == "Total Current Liabilities" || row.ColumnName == "Total Non Current Liabilities" {
				dbResponse.BalanceSheet.TotalLiabilities.CoCode = row.CoCode
				dbResponse.BalanceSheet.TotalLiabilities.Y0 = row.Y0
				dbResponse.BalanceSheet.TotalLiabilities.Y1 = row.Y1
				dbResponse.BalanceSheet.TotalLiabilities.Y2 = row.Y2
				dbResponse.BalanceSheet.TotalLiabilities.Y3 = row.Y3
				dbResponse.BalanceSheet.TotalLiabilities.Y4 = row.Y4

			} else if row.ColumnName == "Net Inc./(Dec.) in Cash and Cash Equivalent" {
				dbResponse.Cashflow = row

			}
		}
		dbResponse.BalanceSheet.TotalLiabilities.ColumnName = "Total Liabilities"

	} else if strings.ToLower(req.Exchange) == constants.NSE && req.NseSymbol != "" {
		loggerconfig.Info("FetchFinancialsData NseSymbol =", req.NseSymbol)
		queryStatement := FetchFinancialsDataV2Query(constants.NSESymbol)
		res, err := dbops.PostgresRepo.Fetch(queryStatement, req.NseSymbol, req.NseSymbol, req.NseSymbol)
		if err != nil {
			loggerconfig.Error("FetchFinancialsDataV2 Error just client fetching data:", err.Error())
			return dbResponse, err
		}
		defer res.Close()

		for res.Next() {
			var row models.FetchFinancialsDataV2
			err = res.Scan(&row.CoCode, &row.ColumnName, &row.Y0, &row.Y1, &row.Y2, &row.Y3, &row.Y4)
			if err != nil {
				loggerconfig.Error("FetchFinancialsDataV2 Error in scanning fetched response from postgres :", err)
				continue
			}

			if row.ColumnName == "Profit After Tax" {
				dbResponse.NetProfit = row

			} else if row.ColumnName == "Total Revenue" {
				dbResponse.Revenue = row

			} else if row.ColumnName == "TOTAL ASSETS" {
				dbResponse.BalanceSheet.TotalAssets = row

			} else if row.ColumnName == "Total Current Liabilities" || row.ColumnName == "Total Non Current Liabilities" {
				dbResponse.BalanceSheet.TotalLiabilities.CoCode = row.CoCode
				dbResponse.BalanceSheet.TotalLiabilities.Y0 = row.Y0
				dbResponse.BalanceSheet.TotalLiabilities.Y1 = row.Y1
				dbResponse.BalanceSheet.TotalLiabilities.Y2 = row.Y2
				dbResponse.BalanceSheet.TotalLiabilities.Y3 = row.Y3
				dbResponse.BalanceSheet.TotalLiabilities.Y4 = row.Y4

			} else if row.ColumnName == "Net Inc./(Dec.) in Cash and Cash Equivalent" {
				dbResponse.Cashflow = row

			}
		}
		dbResponse.BalanceSheet.TotalLiabilities.ColumnName = "Total Liabilities"

	} else {
		loggerconfig.Error("FetchFinancialsDataV2 Error no NSE Symbol or BSE Token entered in request!")
		return dbResponse, err
	}

	loggerconfig.Info("FetchFinancialsDataV2 SuccessFul, dbResponse", dbResponse)
	return dbResponse, nil
}

func (pgObj *Postgres) FetchChangeInInstitutionalHoldingData(allIsin models.AllIsin) (models.AllChangeInInstitutionalHoldingDb, error) {
	var allChangeInInstitutionalHoldingDb models.AllChangeInInstitutionalHoldingDb
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchChangeInInstitutionalHoldingData Error if database is alive :", err.Error())
		return allChangeInInstitutionalHoldingDb, err
	}

	if allIsin.Isin == nil || len(allIsin.Isin) == 0 {
		return allChangeInInstitutionalHoldingDb, fmt.Errorf("FetchChangeInInstitutionalHoldingData: input ISIN list is empty or nil")
	}
	stringOfIsin := ""
	for i := 0; i < len(allIsin.Isin); i++ {
		if len(allIsin.Isin[i]) == constants.ZERO {
			continue
		}
		stringOfIsin += "'" + allIsin.Isin[i] + "',"
	}
	if len(stringOfIsin) != 0 {
		stringOfIsin = stringOfIsin[:len(stringOfIsin)-1]
	}

	queryStatement := fmt.Sprintf(`SELECT cm.isin, shpa.yrc, (FII_Holding + DII_Holding) AS Total_Holding
	FROM shareholdingpatternaggregate shpa
	JOIN companymaster cm ON shpa.cocode = cm.cocode
	WHERE cm.isin IN (%s)
	AND shpa.yrc IN (
	  SELECT MAX(yrc)
	  FROM shareholdingpatternaggregate
	  WHERE cocode = shpa.cocode
	  UNION
	  SELECT MAX(yrc)
	  FROM shareholdingpatternaggregate
	  WHERE cocode = shpa.cocode
	  AND yrc < (
		SELECT MAX(yrc)
		FROM shareholdingpatternaggregate
		WHERE cocode = shpa.cocode
	  )
	)`, stringOfIsin)
	res, err := dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchChangeInInstitutionalHoldingData Error just client fetching data:", err.Error())
		return allChangeInInstitutionalHoldingDb, err
	}
	defer res.Close()

	var quarter = constants.CurrentQuarter
	var diffrenceInInstitutionalHolding float64
	for res.Next() {
		var changeInInstitutionalHoldingData models.ChangeInInstitutionalHoldingData
		err = res.Scan(&changeInInstitutionalHoldingData.Isin, &changeInInstitutionalHoldingData.Yrc, &changeInInstitutionalHoldingData.InstitutionalHolding)
		if err != nil {
			loggerconfig.Error("FetchChangeInInstitutionalHoldingData Error in scanning fetched companymaster response from postgres :", err)
			break
		}
		if quarter == constants.CurrentQuarter {
			diffrenceInInstitutionalHolding = changeInInstitutionalHoldingData.InstitutionalHolding
			quarter = constants.PreviousQuarter
		} else {
			diffrenceInInstitutionalHolding -= changeInInstitutionalHoldingData.InstitutionalHolding
			if diffrenceInInstitutionalHolding > constants.FIVE {
				var changeInInstitutionalHoldingCur models.ChangeInInstitutionalHoldingDb
				changeInInstitutionalHoldingCur.Isin = changeInInstitutionalHoldingData.Isin
				changeInInstitutionalHoldingCur.DiffrenceInInstitutionalHolding = diffrenceInInstitutionalHolding
				allChangeInInstitutionalHoldingDb.ChangeInInstitutionalHoldingDb = append(allChangeInInstitutionalHoldingDb.ChangeInInstitutionalHoldingDb, changeInInstitutionalHoldingCur)
				quarter = constants.CurrentQuarter
			}
		}
	}
	loggerconfig.Info("FetchChangeInInstitutionalHoldingData SuccessFul, dbResponse", allChangeInInstitutionalHoldingDb)
	return allChangeInInstitutionalHoldingDb, nil
}

func (pgObj *Postgres) FetchRoeAndStockReturnData(allIsin models.AllIsin) (models.AllRoeAndStockReturnDb, error) {
	var allRoeAndStockReturnDb models.AllRoeAndStockReturnDb
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchRoeAndStockReturnData Error if database is alive :", err.Error())
		return allRoeAndStockReturnDb, err
	}

	if allIsin.Isin == nil || len(allIsin.Isin) == 0 {
		return allRoeAndStockReturnDb, fmt.Errorf("FetchRoeAndStockReturnData: input ISIN list is empty or nil")
	}
	stringOfIsin := ""
	for i := 0; i < len(allIsin.Isin); i++ {
		if len(allIsin.Isin[i]) == constants.ZERO {
			continue
		}
		stringOfIsin += "'" + allIsin.Isin[i] + "',"
	}
	if len(stringOfIsin) != 0 {
		stringOfIsin = stringOfIsin[:len(stringOfIsin)-1]
	}

	queryStatement := fmt.Sprintf(`SELECT cm.cocode,cm.industrycode,cm.isin,ltp,r3yearprice,
	pls.y0,pls.y1,pls.y2,
	bs.y0,bs.y1,bs.y2
	FROM companypeerperformance cpp
	INNER JOIN plstatement pls ON pls.cocode=cpp.cocode
	INNER JOIN balancesheets bs ON bs.cocode=cpp.cocode
	INNER JOIN companymaster cm ON cm.cocode=cpp.cocode
	WHERE pls.columnname='Profit Attributable to Shareholders' 
	AND bs.columnname='Total Shareholder''s Fund'
	AND bs.typecs='C' AND pls.typecs='c'
	AND cm.isin IN (%s);`, stringOfIsin)
	res, err := dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchRoeAndStockReturnData Error just client fetching data:", err.Error())
		return allRoeAndStockReturnDb, err
	}
	defer res.Close()
	for res.Next() {
		var roeAndStockReturnData models.RoeAndStockReturnData
		err = res.Scan(&roeAndStockReturnData.Cocode, &roeAndStockReturnData.Industrycode, &roeAndStockReturnData.Isin, &roeAndStockReturnData.Ltp, &roeAndStockReturnData.Return3Yrs, &roeAndStockReturnData.Y0ProfitAES, &roeAndStockReturnData.Y1ProfitAES, &roeAndStockReturnData.Y2ProfitAES, &roeAndStockReturnData.Y0TotalShareholderFund, &roeAndStockReturnData.Y1TotalShareholderFund, &roeAndStockReturnData.Y2TotalShareholderFund)
		if err != nil {
			loggerconfig.Error("FetchRoeAndStockReturnData Error in scanning fetched response from postgres :", err)
			break
		}

		if roeAndStockReturnData.Industrycode != constants.BankIndustryCode {
			allRoeAndStockReturnDb.RoeAndStockReturndb = append(allRoeAndStockReturnDb.RoeAndStockReturndb, roeAndStockReturnData)
		}
	}
	loggerconfig.Info("FetchRoeAndStockReturnData SuccessFul, dbResponse", allRoeAndStockReturnDb)
	return allRoeAndStockReturnDb, nil
}

func (pgObj *Postgres) FetchNseBondData(isin string) (models.NseBondStoreDbData, error) {
	var nseBondStoreDbData models.NseBondStoreDbData

	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchNseBondData Error if database is alive :", err.Error())
		return nseBondStoreDbData, err
	}

	inputIsin := "'" + isin + "'"

	queryStatement := fmt.Sprintf(`SELECT * from bondsdata WHERE isin = %s;`, inputIsin)
	res, err := dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchNseBondData Error just client fetching data:", err.Error())
		return nseBondStoreDbData, err
	}
	defer res.Close()

	resultCount := 0

	for res.Next() {
		var nseBondStoreDbDataScan models.NseBondStoreDbData
		err = res.Scan(&nseBondStoreDbDataScan.Symbol, &nseBondStoreDbDataScan.Series, &nseBondStoreDbDataScan.BondType, &nseBondStoreDbDataScan.Open, &nseBondStoreDbDataScan.High, &nseBondStoreDbDataScan.Low, &nseBondStoreDbDataScan.LtP, &nseBondStoreDbDataScan.Close, &nseBondStoreDbDataScan.Per, &nseBondStoreDbDataScan.Qty, &nseBondStoreDbDataScan.TrdVal, &nseBondStoreDbDataScan.Coupr, &nseBondStoreDbDataScan.CreditRating, &nseBondStoreDbDataScan.RatingAgency, &nseBondStoreDbDataScan.FaceValue, &nseBondStoreDbDataScan.NxtipDate, &nseBondStoreDbDataScan.MaturityDate, &nseBondStoreDbDataScan.BYield, &nseBondStoreDbDataScan.CompanyName, &nseBondStoreDbDataScan.Industry, &nseBondStoreDbDataScan.IsFNOSec, &nseBondStoreDbDataScan.IsCASec, &nseBondStoreDbDataScan.IsSLBSec, &nseBondStoreDbDataScan.IsDebtSec, &nseBondStoreDbDataScan.IsSuspended, &nseBondStoreDbDataScan.IsETFSec, &nseBondStoreDbDataScan.IsDelisted, &nseBondStoreDbDataScan.Isin)
		if err != nil {
			loggerconfig.Error("FetchNseBondData Error in scanning fetched response from postgres :", err)
			break
		}
		resultCount++
		nseBondStoreDbData = nseBondStoreDbDataScan
	}

	if resultCount == 0 {
		loggerconfig.Error("FetchNseBondData, error in fetching data: ", err)
		return nseBondStoreDbData, errors.New(constants.InvalidBondIsin)
	}

	loggerconfig.Info("FetchNseBondData SuccessFul, dbResponse", nseBondStoreDbData)
	return nseBondStoreDbData, nil
}

func (pgObj *Postgres) NudgeCheck(isin string) (bool, bool, error) {
	var asmPresent, gsmPresent bool

	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("NudgeCheck Error if database is alive :", err.Error())
		return false, false, err
	}

	inputIsin := "'" + isin + "'"

	queryStatement := fmt.Sprintf(`SELECT isin FROM additional_surveillance_measure_list WHERE isin = %s;`, inputIsin)
	resAsm, err := dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("NudgeCheck Error just client fetching data:", err.Error())
		return false, false, err
	}
	defer resAsm.Close()

	for resAsm.Next() {
		var dbIsinAsm string
		err = resAsm.Scan(&dbIsinAsm)
		if err != nil {
			loggerconfig.Error("NudgeCheck Error in scanning fetched response from postgres :", err)
			break
		}
		if dbIsinAsm == isin {
			asmPresent = true
		}
	}

	queryStatement = fmt.Sprintf(`SELECT isin FROM graded_surveillance_measure WHERE isin = %s;`, inputIsin)
	resGsm, err := dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("NudgeCheck Error just client fetching data:", err.Error())
		return false, false, err
	}
	defer resGsm.Close()

	for resGsm.Next() {
		var dbIsinGsm string
		err = resGsm.Scan(&dbIsinGsm)
		if err != nil {
			loggerconfig.Error("NudgeCheck Error in scanning fetched response from postgres :", err)
			break
		}
		if dbIsinGsm == isin {
			gsmPresent = true
		}
	}

	return asmPresent, gsmPresent, nil
}

func FetchPeersDataV2Query(searchBy string) string {
	queryStatement := fmt.Sprintf(`SELECT companymaster.companyname, companymaster.sectorcode,companymaster.nsesymbol, companymaster.bsecode, cpr.mcap, cpr.pe,
	cpr.pb, ttmdata.epsttm, ttmdata.roettm
	FROM ttmdata
	INNER JOIN companymaster ON ttmdata.cocode = companymaster.cocode
	INNER JOIN companypeerratios cpr ON (cpr.cocode,cpr.typecs)=(ttmdata.cocode,LOWER(ttmdata.typecs))
	WHERE companymaster.%s = $1 AND ttmdata.typecs='C';`, searchBy)
	return queryStatement
}

func (pgObj *Postgres) FetchPeersV2Data(req models.FetchPeersV2Req) ([]models.FetchPeerV2Data, error) {
	ctx := context.Background()
	var dbResponse []models.FetchPeerV2Data
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchPeersV2Data Error if database is alive :", err.Error())
		return dbResponse, err
	}

	var queryStatement string
	var res *sql.Rows
	if req.Isin != "" {
		loggerconfig.Info("FetchPeersV2Data isin =", req.Isin)
		queryStatement = FetchPeersDataV2Query(constants.ISIN)
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.Isin)
		if err != nil {
			loggerconfig.Error("FetchPeersV2Data Error isin fetching data:", err.Error())
			return dbResponse, err
		}
	} else if req.Exchange == "" {
		loggerconfig.Error("FetchPeersV2Data Error no exchange entered in request!")
		return dbResponse, err
	} else if strings.ToLower(req.Exchange) == constants.NSE && req.NseSymbol != "" {
		loggerconfig.Info("FetchPeersV2Data NseSymbol =", req.NseSymbol)
		queryStatement = FetchPeersDataV2Query(constants.NSESymbol)
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.NseSymbol)
		if err != nil {
			loggerconfig.Error("FetchPeersV2Data Error NSESymbol fetching data:", err.Error())
			return dbResponse, err
		}
	} else if strings.ToLower(req.Exchange) == constants.BSE && req.BseToken != "" {
		loggerconfig.Info("FetchPeersV2Data BseToken =", req.BseToken)
		queryStatement = FetchPeersDataV2Query(constants.BSECode)
		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.BseToken)
		if err != nil {
			loggerconfig.Error("FetchPeersV2Data Error client+transaction fetching data:", err.Error())
			return dbResponse, err
		}
	} else {
		loggerconfig.Error("FetchPeersV2Data Error no NSE Symbol or BSE Token entered in request!")
		return dbResponse, err
	}
	var row models.FetchPeerV2Data
	defer res.Close()
	for res.Next() {
		err = res.Scan(&row.Company, &row.SectorCode, &row.TradingSymbol, &row.Token, &row.Mcap, &row.PeRatio, &row.PbRatio, &row.Eps, &row.Roe)
		if err != nil {
			return dbResponse, err
		}
		row.Exchange = constants.BSE
	}
	dbResponse = append(dbResponse, row)

	queryStatement = `(SELECT companymaster.companyname, companymaster.sectorcode,companymaster.nsesymbol, companymaster.bsecode, cpr.mcap , cpr.pe,
			cpr.pb, ttmdata.epsttm, ttmdata.roettm
			FROM ttmdata
			INNER JOIN companymaster ON ttmdata.cocode = companymaster.cocode
			INNER JOIN companypeerratios cpr ON (cpr.cocode,cpr.typecs)=(ttmdata.cocode,LOWER(ttmdata.typecs))
			WHERE companymaster.sectorcode = $1  AND cpr.mcap > $2 AND ttmdata.typecs='C'
			ORDER BY mcap ASC
			LIMIT 4) 
			UNION ALL 
			(SELECT companymaster.companyname, companymaster.sectorcode,companymaster.nsesymbol, companymaster.bsecode, cpr.mcap, cpr.pe,
			cpr.pb, ttmdata.epsttm, ttmdata.roettm
			FROM ttmdata
			INNER JOIN companymaster ON ttmdata.cocode = companymaster.cocode
			INNER JOIN companypeerratios cpr ON (cpr.cocode,cpr.typecs)=(ttmdata.cocode,LOWER(ttmdata.typecs))
			WHERE companymaster.sectorcode = $1  AND cpr.mcap < $2 AND ttmdata.typecs='C'
			ORDER BY mcap DESC
			LIMIT 4)`
	loggerconfig.Info(queryStatement)
	res, err = dbops.PostgresRepo.Fetch(queryStatement, row.SectorCode, row.Mcap)
	if err != nil {
		loggerconfig.Error("FetchPeersV2Data Error client+transaction fetching data:", err.Error())
		return dbResponse, err
	}
	defer res.Close()
	var rowsLowerThanFilterValue []models.FetchPeerV2Data

	for res.Next() {
		err = res.Scan(&row.Company, &row.SectorCode, &row.TradingSymbol, &row.Token, &row.Mcap, &row.PeRatio, &row.PbRatio, &row.Eps, &row.Roe)
		if err != nil {
			return dbResponse, err
		}
		row.Exchange = constants.BSE
		if dbResponse[0].Mcap < row.Mcap {
			dbResponse = append(dbResponse, row)
		} else {
			rowsLowerThanFilterValue = append(rowsLowerThanFilterValue, row)
		}
	}

	for len(dbResponse)+len(rowsLowerThanFilterValue) > 5 {
		if len(dbResponse) > 3 {
			dbResponse = dbResponse[:len(dbResponse)-1]
		}
		if len(rowsLowerThanFilterValue) > 2 {
			rowsLowerThanFilterValue = rowsLowerThanFilterValue[:len(rowsLowerThanFilterValue)-1]
		}
	}
	dbResponse = append(dbResponse, rowsLowerThanFilterValue...)

	return dbResponse, nil
}

func FetchFinancialsDataV3Query(searchBy string) string {
	queryStatement := fmt.Sprintf(`SELECT cocode, columnname, y0, y1, y2, y3, y4, yrc0, yrc1, yrc2, yrc3, yrc4
	FROM (
	  SELECT cm.cocode, columnname, y0, y1, y2, y3, y4,yrc0,yrc1,yrc2,yrc3,yrc4
	  FROM plstatement pl
	  INNER JOIN companymaster cm ON pl.cocode = cm.cocode
	  INNER JOIN plstatement_yrc_mapping plsyrc ON (plsyrc.cocode,plsyrc.typecs)=(pl.cocode,pl.typecs)
	  WHERE cm.%s =$1
	  AND pl.typecs = 'c' AND columnname IN ('Profit After Tax', 'Total Revenue')
	
	  UNION ALL
	
	  SELECT cm.cocode, columnname, y0,  y1,  y2,  y3,  y4, yrc0, yrc1, yrc2, yrc3, yrc4
	  FROM balancesheets bs
	  INNER JOIN companymaster cm ON bs.cocode = cm.cocode
	  INNER JOIN balancesheets_yrc_mapping bsyrc ON (bsyrc.cocode,bsyrc.typecs)=(bs.cocode,bs.typecs)
	  WHERE cm.%s =$2
	  AND bs.typecs = 'C' AND columnname IN ('TOTAL EQUITY AND LIABILITIES','Total Equity','TOTAL ASSETS')
	
	  UNION ALL 
	
	  SELECT cm.cocode, columnname,  y0,  y1,  y2,  y3,  y4, yrc0, yrc1, yrc2, yrc3, yrc4
	  FROM cashflow cf
	  INNER JOIN companymaster cm ON cf.cocode = cm.cocode
	  INNER JOIN cashflow_yrc_mapping cfyrc ON (cfyrc.cocode,cfyrc.typecs)=(cf.cocode,cf.typecs)
	  WHERE cm.%s =$3
	  AND cf.typecs = 'c' AND columnname = 'Net Inc./(Dec.) in Cash and Cash Equivalent'
	) AS subquery
	ORDER BY cocode, columnname;`, searchBy, searchBy, searchBy)
	return queryStatement
}

func (pgObj *Postgres) FetchFinancialsDataV3(req models.FetchFinancialsReq) (models.FetchFinancialsV3Res, error) {
	ctx := context.Background()
	var dbResponse models.FetchFinancialsV3Res
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchFinancialsDataV3 Error if database is alive :", err.Error())
		return dbResponse, err
	}

	var res *sql.Rows
	var queryStatement string

	if req.Isin != "" {
		loggerconfig.Info("FetchFinancialsDataV3 isin =", req.Isin)
		queryStatement = FetchFinancialsDataV3Query(constants.ISIN)

		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.Isin, req.Isin, req.Isin)
		if err != nil {
			loggerconfig.Error("FetchFinancialsDataV3 Error just client fetching data:", err.Error())
			return dbResponse, err
		}
		defer res.Close()

	} else if req.Exchange == "" {
		loggerconfig.Error("FetchFinancialsData Error no exchange entered in request!")
		return dbResponse, err
	} else if strings.ToLower(req.Exchange) == constants.BSE && req.BseToken != "" {
		loggerconfig.Info("FetchFinancialsDataV3 BseToken =", req.BseToken)
		queryStatement = FetchFinancialsDataV3Query(constants.BSECode)

		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.BseToken, req.BseToken, req.BseToken)
		if err != nil {
			loggerconfig.Error("FetchFinancialsDataV3 Error just client fetching data:", err.Error())
			return dbResponse, err
		}
		defer res.Close()

	} else if strings.ToLower(req.Exchange) == constants.NSE && req.NseSymbol != "" {
		loggerconfig.Info("FetchFinancialsData NseSymbol =", req.NseSymbol)
		queryStatement = FetchFinancialsDataV3Query(constants.NSESymbol)

		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.NseSymbol, req.NseSymbol, req.NseSymbol)
		if err != nil {
			loggerconfig.Error("FetchFinancialsDataV3 Error just client fetching data:", err.Error())
			return dbResponse, err
		}
		defer res.Close()

	} else {
		loggerconfig.Error("FetchFinancialsDataV3 Error no NSE Symbol or BSE Token entered in request!")
		return dbResponse, err
	}

	for res.Next() {
		var row models.FetchFinancialsDataV3
		err = res.Scan(&row.CoCode, &row.ColumnName, &row.Y0, &row.Y1, &row.Y2, &row.Y3, &row.Y4, &row.Yrc0, &row.Yrc1, &row.Yrc2, &row.Yrc3, &row.Yrc4)
		if err != nil {
			loggerconfig.Error("FetchFinancialsDataV3 Error in scanning fetched response from postgres :", err)
			continue
		}
		if row.ColumnName == constants.ProfitAfterTax {
			dbResponse.NetProfit = row

		} else if row.ColumnName == constants.TotalRevenue[1:len(constants.TotalRevenue)-1] {
			dbResponse.Revenue = row

		} else if row.ColumnName == constants.TotalAssets {
			dbResponse.BalanceSheet.TotalAssets = row

		} else if row.ColumnName == constants.TotalEquity {
			dbResponse.BalanceSheet.TotalLiabilities.Y0 -= row.Y0
			dbResponse.BalanceSheet.TotalLiabilities.Y1 -= row.Y1
			dbResponse.BalanceSheet.TotalLiabilities.Y2 -= row.Y2
			dbResponse.BalanceSheet.TotalLiabilities.Y3 -= row.Y3
			dbResponse.BalanceSheet.TotalLiabilities.Y4 -= row.Y4
			dbResponse.BalanceSheet.TotalLiabilities.Yrc0 = row.Yrc0
			dbResponse.BalanceSheet.TotalLiabilities.Yrc1 = row.Yrc1
			dbResponse.BalanceSheet.TotalLiabilities.Yrc2 = row.Yrc2
			dbResponse.BalanceSheet.TotalLiabilities.Yrc3 = row.Yrc3
			dbResponse.BalanceSheet.TotalLiabilities.Yrc4 = row.Yrc4

		} else if row.ColumnName == constants.TotalEquityAndLiabilities {
			dbResponse.BalanceSheet.TotalLiabilities.CoCode = row.CoCode
			dbResponse.BalanceSheet.TotalLiabilities.Y0 += row.Y0
			dbResponse.BalanceSheet.TotalLiabilities.Y1 += row.Y1
			dbResponse.BalanceSheet.TotalLiabilities.Y2 += row.Y2
			dbResponse.BalanceSheet.TotalLiabilities.Y3 += row.Y3
			dbResponse.BalanceSheet.TotalLiabilities.Y4 += row.Y4
			dbResponse.BalanceSheet.TotalLiabilities.Yrc0 = row.Yrc0
			dbResponse.BalanceSheet.TotalLiabilities.Yrc1 = row.Yrc1
			dbResponse.BalanceSheet.TotalLiabilities.Yrc2 = row.Yrc2
			dbResponse.BalanceSheet.TotalLiabilities.Yrc3 = row.Yrc3
			dbResponse.BalanceSheet.TotalLiabilities.Yrc4 = row.Yrc4

		} else if row.ColumnName == constants.NetIncDecInCashAndCashEquivalent {
			dbResponse.Cashflow = row

		}
	}
	dbResponse.BalanceSheet.TotalLiabilities.ColumnName = constants.TotalLiabilities

	loggerconfig.Info("FetchFinancialsDataV3 SuccessFul, dbResponse", dbResponse)
	return dbResponse, nil
}

func (pgObj *Postgres) StoreFCMToken(id string, fcmtoken string, requestId string) (string, error) {
	id = strings.ToUpper(id)

	var clientDetails models.MongoClientsDetails
	// userId will be email id in many cases(as login Id supports both email and loginId)
	err := dbops.MongoDaoRepo.FindOne(constants.CLIENTDETAILS, bson.M{"email": bson.M{"$regex": "^" + id + "$", "$options": "i"}}, &clientDetails)
	if err != nil && err.Error() != constants.MongoNoDocError {
		loggerconfig.Error("StoreFCMToken Unable to fetch the client-details Data for id = ", id, " error :", err, " requestid=", requestId)
	}

	queryStatement := `
	INSERT INTO client_fcm("client_id", "fcm_token")
	VALUES($1, $2)
	ON CONFLICT ON CONSTRAINT pk_clientid_fcmtoken DO NOTHING;`

	var clientId string

	if clientDetails.ClientID != "" {
		clientId = strings.ToUpper(clientDetails.ClientID)
	} else {
		clientId = strings.ToUpper(id)
	}

	rows, err := dbops.PostgresRepo.Insert(queryStatement, clientId, fcmtoken)
	if err != nil {
		loggerconfig.Error("StoreFCMToken Error inserting data:", err.Error(), " requestid=", requestId)
		return clientId, err
	}
	defer rows.Close()

	loggerconfig.Info("StoreFCMToken successful clientId: ", id, " requestid=", requestId)
	return clientId, nil
}

func (pgObj *Postgres) DeleteFCMToken(userID string, fcmtoken string) error {
	userID = strings.ToUpper(userID)

	// Delete the record from the client_fcm table
	deleteQuery := "DELETE FROM client_fcm WHERE client_id = $1 AND fcm_token = $2"
	res, err := dbops.PostgresRepo.Delete(deleteQuery, userID, fcmtoken)
	if err != nil {
		loggerconfig.Error("DeleteFCMToken Error deleting FCM token: ", err)
		return err
	}
	defer res.Close()
	loggerconfig.Info("DeleteFCMToken FCM token deleted for user ", userID)

	return nil
}

func FetchFinancialsDataV4Query(searchBy string) string {
	queryStatement := fmt.Sprintf(`SELECT cocode, typecs, columnname, y0, y1, y2, y3, y4, yrc0, yrc1, yrc2, yrc3, yrc4
	FROM (
	  SELECT cm.cocode, pl.typecs AS typecs, columnname, y0, y1, y2, y3, y4,yrc0,yrc1,yrc2,yrc3,yrc4
	  FROM plstatement pl
	  INNER JOIN companymaster cm ON pl.cocode = cm.cocode
	  INNER JOIN plstatement_yrc_mapping plsyrc ON (plsyrc.cocode,plsyrc.typecs)=(pl.cocode,pl.typecs)
	  WHERE cm.%s =$1
	  AND columnname IN ('Profit After Tax', 'Total Revenue')
	
	  UNION ALL
	
	  SELECT cm.cocode, bs.typecs AS typecs, columnname, y0,  y1,  y2,  y3,  y4, yrc0, yrc1, yrc2, yrc3, yrc4
	  FROM balancesheets bs
	  INNER JOIN companymaster cm ON bs.cocode = cm.cocode
	  INNER JOIN balancesheets_yrc_mapping bsyrc ON (bsyrc.cocode,bsyrc.typecs)=(bs.cocode,bs.typecs)
	  WHERE cm.%s =$2
	  AND columnname IN ('TOTAL EQUITY AND LIABILITIES','Total Equity','TOTAL ASSETS')
	
	  UNION ALL
	
	  SELECT cm.cocode, cf.typecs AS typecs, columnname,  y0,  y1,  y2,  y3,  y4, yrc0, yrc1, yrc2, yrc3, yrc4
	  FROM cashflow cf
	  INNER JOIN companymaster cm ON cf.cocode = cm.cocode
	  INNER JOIN cashflow_yrc_mapping cfyrc ON (cfyrc.cocode,cfyrc.typecs)=(cf.cocode,cf.typecs)
	  WHERE cm.%s =$3
	  AND columnname = 'Net Inc./(Dec.) in Cash and Cash Equivalent'
	) AS subquery
	ORDER BY cocode, columnname;`, searchBy, searchBy, searchBy)
	return queryStatement
}

func (pgObj *Postgres) FetchFinancialsDataV4(req models.FetchFinancialsReq) (models.FetchFinancialsV4Res, error) {
	ctx := context.Background()
	var dbResponse models.FetchFinancialsV4Res
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchFinancialsDataV4 Error if database is alive :", err.Error())
		return dbResponse, err
	}

	var res *sql.Rows
	var queryStatement string

	if req.Isin != "" {
		loggerconfig.Info("FetchFinancialsDataV4 isin =", req.Isin)
		queryStatement = FetchFinancialsDataV4Query(constants.ISIN)

		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.Isin, req.Isin, req.Isin)
		if err != nil {
			loggerconfig.Error("FetchFinancialsDataV4 Error just client fetching data:", err.Error())
			return dbResponse, err
		}
		defer res.Close()

	} else if req.Exchange == "" {
		loggerconfig.Error("FetchFinancialsData Error no exchange entered in request!")
		return dbResponse, err
	} else if strings.ToLower(req.Exchange) == constants.BSE && req.BseToken != "" {
		loggerconfig.Info("FetchFinancialsDataV4 BseToken =", req.BseToken)
		queryStatement = FetchFinancialsDataV4Query(constants.BSECode)

		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.BseToken, req.BseToken, req.BseToken)
		if err != nil {
			loggerconfig.Error("FetchFinancialsDataV4 Error just client fetching data:", err.Error())
			return dbResponse, err
		}
		defer res.Close()

	} else if strings.ToLower(req.Exchange) == constants.NSE && req.NseSymbol != "" {
		loggerconfig.Info("FetchFinancialsDataV4 NseSymbol =", req.NseSymbol)
		queryStatement = FetchFinancialsDataV4Query(constants.NSESymbol)

		res, err = dbops.PostgresRepo.Fetch(queryStatement, req.NseSymbol, req.NseSymbol, req.NseSymbol)
		if err != nil {
			loggerconfig.Error("FetchFinancialsDataV4 Error just client fetching data:", err.Error())
			return dbResponse, err
		}
		defer res.Close()

	} else {
		loggerconfig.Error("FetchFinancialsDataV4 Error no NSE Symbol or BSE Token entered in request!")
		return dbResponse, err
	}

	for res.Next() {
		var row models.FetchFinancialsDataV4
		err = res.Scan(&row.CoCode, &row.TypeCS, &row.ColumnName, &row.Y0, &row.Y1, &row.Y2, &row.Y3, &row.Y4, &row.Yrc0, &row.Yrc1, &row.Yrc2, &row.Yrc3, &row.Yrc4)
		if err != nil {
			loggerconfig.Error("FetchFinancialsDataV4 Error in scanning fetched response from postgres :", err)
			continue
		}
		if row.ColumnName == constants.ProfitAfterTax {
			if row.TypeCS == "s" || row.TypeCS == "S" {
				dbResponse.NetProfitStandalone = row
			} else {
				dbResponse.NetProfitConsolidated = row
			}

		} else if row.ColumnName == constants.TotalRevenue[1:len(constants.TotalRevenue)-1] {
			if row.TypeCS == "s" || row.TypeCS == "S" {
				dbResponse.RevenueStandalone = row
			} else {
				dbResponse.RevenueConsolidated = row
			}

		} else if row.ColumnName == constants.TotalAssets {
			if row.TypeCS == "s" || row.TypeCS == "S" {
				dbResponse.BalanceSheetStandalone.TotalAssets = row
			} else {
				dbResponse.BalanceSheetConsolidated.TotalAssets = row
			}

		} else if row.ColumnName == constants.TotalEquity {
			if row.TypeCS == "s" || row.TypeCS == "S" {
				dbResponse.BalanceSheetStandalone.TotalLiabilities.CoCode = row.CoCode
				dbResponse.BalanceSheetStandalone.TotalLiabilities.TypeCS = row.TypeCS
				dbResponse.BalanceSheetStandalone.TotalLiabilities.ColumnName = constants.TotalLiabilities
				dbResponse.BalanceSheetStandalone.TotalLiabilities.Y0 -= row.Y0
				dbResponse.BalanceSheetStandalone.TotalLiabilities.Y1 -= row.Y1
				dbResponse.BalanceSheetStandalone.TotalLiabilities.Y2 -= row.Y2
				dbResponse.BalanceSheetStandalone.TotalLiabilities.Y3 -= row.Y3
				dbResponse.BalanceSheetStandalone.TotalLiabilities.Y4 -= row.Y4
				dbResponse.BalanceSheetStandalone.TotalLiabilities.Yrc0 = row.Yrc0
				dbResponse.BalanceSheetStandalone.TotalLiabilities.Yrc1 = row.Yrc1
				dbResponse.BalanceSheetStandalone.TotalLiabilities.Yrc2 = row.Yrc2
				dbResponse.BalanceSheetStandalone.TotalLiabilities.Yrc3 = row.Yrc3
				dbResponse.BalanceSheetStandalone.TotalLiabilities.Yrc4 = row.Yrc4
			} else {
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.CoCode = row.CoCode
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.TypeCS = row.TypeCS
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.ColumnName = constants.TotalLiabilities
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.Y0 -= row.Y0
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.Y1 -= row.Y1
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.Y2 -= row.Y2
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.Y3 -= row.Y3
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.Y4 -= row.Y4
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.Yrc0 = row.Yrc0
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.Yrc1 = row.Yrc1
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.Yrc2 = row.Yrc2
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.Yrc3 = row.Yrc3
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.Yrc4 = row.Yrc4
			}

		} else if row.ColumnName == constants.TotalEquityAndLiabilities {

			if row.TypeCS == "s" || row.TypeCS == "S" {
				dbResponse.BalanceSheetStandalone.TotalLiabilities.Y0 += row.Y0
				dbResponse.BalanceSheetStandalone.TotalLiabilities.Y1 += row.Y1
				dbResponse.BalanceSheetStandalone.TotalLiabilities.Y2 += row.Y2
				dbResponse.BalanceSheetStandalone.TotalLiabilities.Y3 += row.Y3
				dbResponse.BalanceSheetStandalone.TotalLiabilities.Y4 += row.Y4
				dbResponse.BalanceSheetStandalone.TotalLiabilities.Yrc0 = row.Yrc0
				dbResponse.BalanceSheetStandalone.TotalLiabilities.Yrc1 = row.Yrc1
				dbResponse.BalanceSheetStandalone.TotalLiabilities.Yrc2 = row.Yrc2
				dbResponse.BalanceSheetStandalone.TotalLiabilities.Yrc3 = row.Yrc3
				dbResponse.BalanceSheetStandalone.TotalLiabilities.Yrc4 = row.Yrc4
			} else {
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.Y0 += row.Y0
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.Y1 += row.Y1
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.Y2 += row.Y2
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.Y3 += row.Y3
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.Y4 += row.Y4
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.Yrc0 = row.Yrc0
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.Yrc1 = row.Yrc1
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.Yrc2 = row.Yrc2
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.Yrc3 = row.Yrc3
				dbResponse.BalanceSheetConsolidated.TotalLiabilities.Yrc4 = row.Yrc4
			}

		} else if row.ColumnName == constants.NetIncDecInCashAndCashEquivalent {
			if row.TypeCS == "s" || row.TypeCS == "S" {
				dbResponse.CashflowStandalone = row
			} else {
				dbResponse.CashflowConsolidated = row
			}

		}
	}

	loggerconfig.Info("FetchFinancialsDataV4 SuccessFul, dbResponse", dbResponse)
	return dbResponse, nil
}

func (pgObj *Postgres) FetchCompanyMaster() ([]models.CompanyDetails, error) {
	var dbResponse []models.CompanyDetails
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchCompanyMaster  Error if database is alive:", err.Error())
		return dbResponse, err
	}
	var queryStatement string
	var res *sql.Rows
	queryStatement = `SELECT cocode,bsecode,nsesymbol,isin FROM companymaster ;`
	res, err = dbops.PostgresRepo.Fetch(queryStatement)
	if err != nil {
		loggerconfig.Error("FetchCompanyMaster Error just client fetching data:", err.Error())
		return dbResponse, err
	}

	defer res.Close()
	var row models.CompanyDetails
	for res.Next() {
		err = res.Scan(&row.CoCode, &row.Bsecode, &row.Nsesymbol, &row.Isin)
		if err != nil {
			break
		}
		dbResponse = append(dbResponse, row)

	}

	return dbResponse, nil
}

func (pgObj *Postgres) InsertTransactionData(payoutDetails models.PayoutDetails) error {
	queryStatement := `
	INSERT INTO transaction_info_v2 (amount_in_paisa, client_id, ifsc, customer_account_number, customer_bank_name, debit_credit, tradelab_funds_updated, backoffice_funds_updated, transaction_type, transaction_id, transaction_status, remarks, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	RETURNING id;`

	rows, err := dbops.PostgresRepo.Insert(queryStatement,
		payoutDetails.Amount,
		payoutDetails.ClientID,
		payoutDetails.Ifsc,
		payoutDetails.AccountNumber,
		payoutDetails.BankName,
		payoutDetails.DebitCredit,
		payoutDetails.TradelabFundsUpdated,
		payoutDetails.BackofficeFundsUpdated,
		payoutDetails.TransactionType,
		payoutDetails.TransactionId,
		payoutDetails.TransactionStatus,
		payoutDetails.Remarks,
		payoutDetails.CreateDate,
		payoutDetails.UpdatedAt,
	)

	if err != nil {
		loggerconfig.Error("InsertPledgeData Error while inserting data:", err)
		return err
	}
	defer rows.Close()

	var insertedID int64
	if rows.Next() {
		if err := rows.Scan(&insertedID); err != nil {
			loggerconfig.Error("InsertPledgeData Error while scanning inserted ID:", err)
			return err
		}
	} else {
		loggerconfig.Error("InsertPledgeData Error: no rows returned")
		return fmt.Errorf("no rows returned")
	}
	loggerconfig.Info("Inserted data with ID:", insertedID)
	return nil
}

func (pgObj *Postgres) UpdateTransactionData(transactionID string, updates map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("UpdateTransactionData Error if database is alive:", err.Error())
		return err
	}

	setClause := []string{}
	args := []interface{}{}
	argID := 1

	for field, value := range updates {
		setClause = append(setClause, fmt.Sprintf("%s = $%d", field, argID))
		args = append(args, value)
		argID++
	}

	args = append(args, transactionID)
	queryStatement := fmt.Sprintf("UPDATE transaction_info_v2 SET %s WHERE transaction_id = $%d", strings.Join(setClause, ", "), argID)

	res, err := dbops.PostgresRepo.Update(queryStatement, args...)
	if err != nil {
		loggerconfig.Error("UpdateTransactionData Error while updating data:", err.Error())
		return err
	}
	defer res.Close()
	loggerconfig.Info("Updated data with transaction ID:", transactionID)
	return nil
}

func (pgObj *Postgres) GetTransactionData(transactionID string) (*models.PayoutDetails, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("GetTransactionData Error if database is alive:", err.Error())
		return nil, err
	}

	queryStatement := `
	SELECT amount_in_paisa, client_id, ifsc, customer_account_number, customer_bank_name, debit_credit, tradelab_funds_updated, backoffice_funds_updated, transaction_type, transaction_id, transaction_status, remarks, created_at, updated_at
	FROM transaction_info_v2
	WHERE transaction_id = $1;`

	rows, err := dbops.PostgresRepo.Fetch(queryStatement, transactionID)
	if err != nil {
		loggerconfig.Error("GetTransactionData Error while executing query:", err.Error())
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("no transaction found for ID: %s", transactionID)
	}

	var payoutDetails models.PayoutDetails
	err = rows.Scan(
		&payoutDetails.Amount,
		&payoutDetails.ClientID,
		&payoutDetails.Ifsc,
		&payoutDetails.AccountNumber,
		&payoutDetails.BankName,
		&payoutDetails.DebitCredit,
		&payoutDetails.TradelabFundsUpdated,
		&payoutDetails.BackofficeFundsUpdated,
		&payoutDetails.TransactionType,
		&payoutDetails.TransactionId,
		&payoutDetails.TransactionStatus,
		&payoutDetails.Remarks,
		&payoutDetails.CreateDate,
		&payoutDetails.UpdatedAt,
	)
	if err != nil {
		loggerconfig.Error("GetTransactionData Error while scanning row:", err.Error())
		return nil, err
	}

	return &payoutDetails, nil
}

func (pgObj *Postgres) CheckExistingPayoutRequest(clientID string) (bool, error) {
	query := `
    SELECT COUNT(*)
    FROM transaction_info_v2
    WHERE client_id = $1 AND transaction_status IN ($2,$3) AND transaction_type = $4;
    `

	rows, err := dbops.PostgresRepo.Fetch(query, clientID, constants.PROCESS, constants.Pending, constants.Payout)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return false, err
		}
	}

	return count > 0, nil
}

func (pgObj *Postgres) InsertPledgeData(pledgeData models.PledgeData) (int64, error) {
	queryStatement := `
	INSERT INTO pledge_data (client_id, segment_id, timestamp, isin, quantity, price, exchange, bo_id, depository, pledge_unpledge, dp_id, pledge_tls, req_id, version, status)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	RETURNING id;`

	rows, err := dbops.PostgresRepo.Insert(queryStatement,
		pledgeData.ClientID,
		pledgeData.SegmentID,
		pledgeData.Timestamp,
		pledgeData.ISIN,
		pledgeData.Quantity,
		pledgeData.Price,
		pledgeData.Exchange,
		pledgeData.BOID,
		pledgeData.Depository,
		pledgeData.PledgeUnpledge,
		pledgeData.DPID,
		pledgeData.PledgeTLS,
		pledgeData.ReqID,
		pledgeData.Version,
		pledgeData.Status,
	)
	if err != nil {
		loggerconfig.Error("InsertPledgeData Error while inserting data:", err)
		return 0, err
	}
	defer rows.Close()

	var insertedID int64
	if rows.Next() {
		if err := rows.Scan(&insertedID); err != nil {
			loggerconfig.Error("InsertPledgeData Error while scanning inserted ID:", err)
			return 0, err
		}
	} else {
		loggerconfig.Error("InsertPledgeData Error: no rows returned")
		return 0, fmt.Errorf("no rows returned")
	}

	loggerconfig.Info("InsertPledgeData Inserted data with ID: ", insertedID)
	return insertedID, nil
}

func (pgObj *Postgres) InsertPledgeDataBatch(pledgeDataBatch []models.PledgeData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := pgObj.conn.PingContext(ctx); err != nil {
		loggerconfig.Error("InsertPledgeDataBatch Error if database is alive:", err)
		return err
	}

	query := `
	INSERT INTO pledge_data (client_id, segment_id, timestamp, isin, quantity, price, exchange, bo_id, depository, pledge_unpledge, dp_id, pledge_tls, req_id, version, status)
	VALUES `

	valueStrings := []string{}
	valueArgs := []interface{}{}
	for i, pd := range pledgeDataBatch {
		valueStrings = append(valueStrings, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*15+1, i*15+2, i*15+3, i*15+4, i*15+5, i*15+6, i*15+7, i*15+8, i*15+9, i*15+10, i*15+11, i*15+12, i*15+13, i*15+14, i*15+15,
		))
		valueArgs = append(valueArgs,
			pd.ClientID,
			pd.SegmentID,
			pd.Timestamp,
			pd.ISIN,
			pd.Quantity,
			pd.Price,
			pd.Exchange,
			pd.BOID,
			pd.Depository,
			pd.PledgeUnpledge,
			pd.DPID,
			pd.PledgeTLS,
			pd.ReqID,
			pd.Version,
			pd.Status,
		)
	}

	query += strings.Join(valueStrings, ",")
	query += " RETURNING id;"

	rows, err := dbops.PostgresRepo.Insert(query, valueArgs...)
	if err != nil {
		loggerconfig.Error("InsertPledgeDataBatch Error while inserting batch data:", err)
		return err
	}
	defer rows.Close()

	var insertedIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			loggerconfig.Error("InsertPledgeDataBatch Error while scanning inserted IDs:", err)
			return err
		}
		insertedIDs = append(insertedIDs, id)
	}
	loggerconfig.Info("InsertPledgeDataBatch Inserted data with IDs: ", insertedIDs)
	return nil
}

func FetchCorporateAnnouncementsQuery(searchBy string) string {
	queryStatement := fmt.Sprintf(`
		SELECT 
			ca.*, 
			cm.isin, 
			cm.bsecode, 
			cm.nsesymbol
		FROM corporate_announcements ca
		INNER JOIN companymaster cm 
			ON ca.cocode = cm.cocode
		WHERE cm.%s = $1
		%s
		ORDER BY ca.announcementdate DESC
		LIMIT 10 OFFSET $%d;
	`, searchBy, "%s", "%d")
	return queryStatement
}

func (pgObj *Postgres) FetchCorporateAnnouncements(req models.FetchCorporateActionsIndividualReq) ([]models.CorporateAnnouncements, error) {
	ctx := context.Background()
	var dbResponse []models.CorporateAnnouncements
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchCorporateAnnouncements Error checking database connection:", err.Error())
		return dbResponse, err
	}

	var queryStatement string
	var res *sql.Rows
	var args []interface{}
	var timeFilterClause string
	var paramCount int = 2

	//time filtering if provided
	if req.StartTimeUnix != 0 && req.EndTimeUnix != 0 {
		timeFilterClause = "AND ca.announcementdate >= $2 AND ca.announcementdate <= $3"
		paramCount = 4
	}

	//offset based on page number (10 items per page)
	offset := (req.PageNo - 1) * 10
	if req.PageNo <= 0 {
		offset = 0 //first page if invalid page number
	}

	if req.Isin != "" {
		loggerconfig.Info("FetchCorporateAnnouncements by ISIN =", req.Isin)
		queryStatement = fmt.Sprintf(FetchCorporateAnnouncementsQuery("isin"), timeFilterClause, paramCount)
		args = append(args, req.Isin)
	} else if req.NseSymbol != "" {
		loggerconfig.Info("FetchCorporateAnnouncements by NSE Symbol =", req.NseSymbol)
		queryStatement = fmt.Sprintf(FetchCorporateAnnouncementsQuery("nsesymbol"), timeFilterClause, paramCount)
		args = append(args, req.NseSymbol)
	} else if req.BseCode != "" {
		loggerconfig.Info("FetchCorporateAnnouncements by BSE Code =", req.BseCode)
		queryStatement = fmt.Sprintf(FetchCorporateAnnouncementsQuery("bsecode"), timeFilterClause, paramCount)
		args = append(args, req.BseCode)
	} else {
		loggerconfig.Error("No ISIN, NSE Symbol, or BSE Code provided in the request!")
		return dbResponse, errors.New("invalid request parameters")
	}

	if req.StartTimeUnix != 0 && req.EndTimeUnix != 0 {
		startTime := time.Unix(req.StartTimeUnix, 0).Format("2006-01-02T00:00:00")
		endTime := time.Unix(req.EndTimeUnix, 0).Format("2006-01-02T00:00:00")
		args = append(args, startTime, endTime)
	}

	args = append(args, offset)

	//query with all parameters
	res, err = dbops.PostgresRepo.Fetch(queryStatement, args...)

	if err != nil {
		loggerconfig.Error("Error FetchCorporateAnnouncements:", err.Error())
		return dbResponse, err
	}

	defer res.Close()
	for res.Next() {
		var announcement models.CorporateAnnouncements
		if err := res.Scan(
			&announcement.Cocode, &announcement.Coname, &announcement.Purpose,
			&announcement.Ratio, &announcement.AnnouncementDate, &announcement.ExecutionDate,
			&announcement.Isin, &announcement.BseCode, &announcement.NseSymbol,
		); err != nil {
			loggerconfig.Error("FetchCorporateAnnouncements Error scanning row:", err.Error())
			return dbResponse, err
		}
		dbResponse = append(dbResponse, announcement)
	}

	return dbResponse, nil
}

func FetchCorporateAnnouncementsQueryAll() string {
	return `
		SELECT 
			ca.*, 
			cm.isin, 
			cm.bsecode, 
			cm.nsesymbol
		FROM corporate_announcements ca
		INNER JOIN companymaster cm 
			ON ca.cocode = cm.cocode
		%s
		ORDER BY ca.announcementdate DESC
		LIMIT 10 OFFSET $%d;
	`
}

func (pgObj *Postgres) FetchCorporateAnnouncementsAll(req models.FetchCorporateActionsAllReq) ([]models.CorporateAnnouncements, error) {
	ctx := context.Background()
	var dbResponse []models.CorporateAnnouncements
	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchCorporateAnnouncementsAll Error checking database connection:", err.Error())
		return dbResponse, err
	}

	var queryStatement string
	var res *sql.Rows
	var args []interface{}
	var timeFilterClause string
	var paramCount int = 1

	//time filtering if provided
	if req.StartTimeUnix != 0 && req.EndTimeUnix != 0 {
		timeFilterClause = "WHERE ca.announcementdate >= $1 AND ca.announcementdate <= $2"
		paramCount = 3
	}

	//offset based on page number (10 items per page)
	offset := (req.PageNo - 1) * 10
	if req.PageNo <= 0 {
		offset = 0 //first page if invalid page number
	}
	queryStatement = fmt.Sprintf(FetchCorporateAnnouncementsQueryAll(), timeFilterClause, paramCount)

	if req.StartTimeUnix != 0 && req.EndTimeUnix != 0 {
		startTime := time.Unix(req.StartTimeUnix, 0).Format("2006-01-02T00:00:00")
		endTime := time.Unix(req.EndTimeUnix, 0).Format("2006-01-02T00:00:00")
		args = append(args, startTime, endTime)
	}
	args = append(args, offset)

	//query with all parameters
	res, err = dbops.PostgresRepo.Fetch(queryStatement, args...)

	if err != nil {
		loggerconfig.Error("FetchCorporateAnnouncementsAll Error fetching corporate announcements:", err.Error())
		return dbResponse, err
	}

	defer res.Close()
	for res.Next() {
		var announcement models.CorporateAnnouncements
		if err := res.Scan(
			&announcement.Cocode, &announcement.Coname, &announcement.Purpose,
			&announcement.Ratio, &announcement.AnnouncementDate, &announcement.ExecutionDate,
			&announcement.Isin, &announcement.BseCode, &announcement.NseSymbol,
		); err != nil {
			loggerconfig.Error("FetchCorporateAnnouncementsAll Error scanning row:", err.Error())
			return dbResponse, err
		}
		dbResponse = append(dbResponse, announcement)
	}

	return dbResponse, nil
}

func (pgObj *Postgres) FetchPledgeTxnPaginated(req models.FetchEpledgeTxnReq, pageSize int) ([]models.PledgeData, error) {
	var pledgeData []models.PledgeData

	isSameDay := req.StartDate == req.EndDate

	var query string
	var rows *sql.Rows
	var err error

	if isSameDay {
		query = `
			SELECT id, client_id, segment_id, "timestamp", isin, quantity, price, exchange, bo_id, depository, pledge_unpledge, dp_id, pledge_tls, req_id, version, status
			FROM pledge_data
			WHERE client_id = $1
			AND "timestamp"::date = $2
			ORDER BY "timestamp" DESC
			LIMIT $3 OFFSET $4
		`
		rows, err = dbops.PostgresRepo.Fetch(query, req.ClientID, req.StartDate, pageSize, (req.Page-1)*pageSize)

	} else {
		req.EndDate += " 23:59:59" //this will add 24 hours in end date so that all transactions of that day(end date) will be fetched

		query = `
			SELECT id, client_id, segment_id, "timestamp", isin, quantity, price, exchange, bo_id, depository, pledge_unpledge, dp_id, pledge_tls, req_id, version, status
			FROM pledge_data
			WHERE client_id = $1
			AND "timestamp" >= $2 AND "timestamp" <= $3
			ORDER BY "timestamp" DESC
			LIMIT $4 OFFSET $5
		`
		rows, err = dbops.PostgresRepo.Fetch(query, req.ClientID, req.StartDate, req.EndDate, pageSize, (req.Page-1)*pageSize)
	}

	if err != nil {
		loggerconfig.Error("fetchPledgeDataPaginated Error while fetching data:", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var pledge models.PledgeData
		if err := rows.Scan(
			&pledge.ID,
			&pledge.ClientID,
			&pledge.SegmentID,
			&pledge.Timestamp,
			&pledge.ISIN,
			&pledge.Quantity,
			&pledge.Price,
			&pledge.Exchange,
			&pledge.BOID,
			&pledge.Depository,
			&pledge.PledgeUnpledge,
			&pledge.DPID,
			&pledge.PledgeTLS,
			&pledge.ReqID,
			&pledge.Version,
			&pledge.Status,
		); err != nil {
			loggerconfig.Error("fetchPledgeDataPaginatedAgainstClientID Error while scanning rows:", err)
			return nil, err
		}

		pledgeData = append(pledgeData, pledge)
	}
	return pledgeData, err
}

func (pgObj *Postgres) FetchSectorWiseCompanyDataV2(sectorCode []string) ([]models.SectorWiseCompanyV2, error) {
	ctx := context.Background()
	var dbResponse []models.SectorWiseCompanyV2

	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("FetchSectorWiseCompanyDataV2 Error if database is alive :", err.Error())
		return dbResponse, err
	}

	placeholders := make([]string, len(sectorCode))
	args := make([]interface{}, len(sectorCode))
	for i, code := range sectorCode {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = code
	}

	queryStatement := fmt.Sprintf(`SELECT sectname FROM sectorlist WHERE sectcode IN (%s);`, strings.Join(placeholders, ", "))
	res, err := dbops.PostgresRepo.Fetch(queryStatement, args...)
	if err != nil {
		loggerconfig.Error("FetchSectorWiseCompanyDataV2 Error in fetching sector Name data:", err.Error())
		return dbResponse, err
	}
	defer res.Close()

	var sectorNames []string
	for res.Next() {
		var sectorName string
		err = res.Scan(&sectorName)
		if err != nil {
			loggerconfig.Error("FetchSectorWiseCompanyDataV2 Error Scan data:", err.Error())
			return dbResponse, err
		}
		sectorNames = append(sectorNames, sectorName)
	}

	placeholders = make([]string, len(sectorNames))
	args = make([]interface{}, len(sectorNames))
	for i, name := range sectorNames {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = strings.ToLower(name)
	}

	queryStatement = fmt.Sprintf(`SELECT * FROM sectorwisecompany WHERE LOWER(sectname) IN (%s);`, strings.Join(placeholders, ", "))
	res2, err := dbops.PostgresRepo.Fetch(queryStatement, args...)
	if err != nil {
		loggerconfig.Error("FetchSectorWiseCompanyDataV2 Error in fetching sector company data:", err.Error())
		return dbResponse, err
	}
	defer res2.Close()

	var companies []models.SectorWiseCompanyDetails
	for res2.Next() {
		var row models.SectorWiseCompanyDetails
		if err := res2.Scan(&row.CoCode, &row.CoName, &row.Lname, &row.ScCode, &row.Symbol, &row.SectName, &row.Isin); err != nil {
			loggerconfig.Error("Error scanning sector-wise company data:", err.Error())
			return nil, err
		}
		companies = append(companies, row)
	}

	sectorWiseData := models.SectorWiseCompanyV2{
		Companies: companies,
	}

	dbResponse = append(dbResponse, sectorWiseData)

	return dbResponse, nil
}

func (pgObj *Postgres) GetSectorWiseCompanyList(page int, sectorName string) ([]models.SectorWiseCompany, error) {

	ctx := context.Background()
	var dbResponse []models.SectorWiseCompany

	err := pgObj.conn.PingContext(ctx)
	if err != nil {
		loggerconfig.Error("GetSectorWiseCompanyList Error if database is alive :", err.Error())
		return dbResponse, err
	}

	queryStatement := `SELECT * FROM sectorwisecompany WHERE sectname = $1 LIMIT $2 OFFSET $3`
	res, err := dbops.PostgresRepo.Fetch(queryStatement, sectorName, constants.Capacity, (page-1)*constants.Capacity)
	if err != nil {
		loggerconfig.Error("GetSectorWiseCompanyList Error in fetching sector-wise company data:", err.Error())
		return dbResponse, err
	}
	defer res.Close()

	for res.Next() {
		var row models.SectorWiseCompany
		if err := res.Scan(&row.CoCode, &row.CoName, &row.Lname, &row.ScCode, &row.Symbol, &row.SectName, &row.Isin); err != nil {
			loggerconfig.Error("Error scanning sector-wise company data:", err.Error())
			return nil, err
		}

		dbResponse = append(dbResponse, row)
	}

	return dbResponse, nil
}
