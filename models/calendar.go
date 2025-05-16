package models

type DateDetails struct {
	Date        string
	DayOfWeek   string
	Description string
	IsHoliday   bool
}

type Calendar struct {
	Date []DateDetails
}

type NSECalendarModel struct {
	Cbm []struct {
		TradingDate string `json:"tradingDate"`
		WeekDay     string `json:"weekDay"`
		Description string `json:"description"`
		SrNo        int    `json:"Sr_no"`
	} `json:"CBM"`
	Cd []struct {
		TradingDate string `json:"tradingDate"`
		WeekDay     string `json:"weekDay"`
		Description string `json:"description"`
		SrNo        int    `json:"Sr_no"`
	} `json:"CD"`
	Cm []struct {
		TradingDate string `json:"tradingDate"`
		WeekDay     string `json:"weekDay"`
		Description string `json:"description"`
		SrNo        int    `json:"Sr_no"`
	} `json:"CM"`
	Cmot []struct {
		TradingDate string `json:"tradingDate"`
		WeekDay     string `json:"weekDay"`
		Description string `json:"description"`
		SrNo        int    `json:"Sr_no"`
	} `json:"CMOT"`
	Com []struct {
		TradingDate string `json:"tradingDate"`
		WeekDay     string `json:"weekDay"`
		Description string `json:"description"`
		SrNo        int    `json:"Sr_no"`
	} `json:"COM"`
	Fo []struct {
		TradingDate string `json:"tradingDate"`
		WeekDay     string `json:"weekDay"`
		Description string `json:"description"`
		SrNo        int    `json:"Sr_no"`
	} `json:"FO"`
	Ird []struct {
		TradingDate string `json:"tradingDate"`
		WeekDay     string `json:"weekDay"`
		Description string `json:"description"`
		SrNo        int    `json:"Sr_no"`
	} `json:"IRD"`
	Mf []struct {
		TradingDate string `json:"tradingDate"`
		WeekDay     string `json:"weekDay"`
		Description string `json:"description"`
		SrNo        int    `json:"Sr_no"`
	} `json:"MF"`
	Ndm []struct {
		TradingDate string `json:"tradingDate"`
		WeekDay     string `json:"weekDay"`
		Description string `json:"description"`
		SrNo        int    `json:"Sr_no"`
	} `json:"NDM"`
	Ntrp []struct {
		TradingDate string `json:"tradingDate"`
		WeekDay     string `json:"weekDay"`
		Description string `json:"description"`
		SrNo        int    `json:"Sr_no"`
	} `json:"NTRP"`
	Slbs []struct {
		TradingDate string `json:"tradingDate"`
		WeekDay     string `json:"weekDay"`
		Description string `json:"description"`
		SrNo        int    `json:"Sr_no"`
	} `json:"SLBS"`
}
