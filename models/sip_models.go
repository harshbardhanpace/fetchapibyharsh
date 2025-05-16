package models

type PlaceSipRequest struct {
	ActionType string   `json:"action_type" binding:"required"`
	Name       string   `json:"name" binding:"required"`
	Baskets    []string `json:"baskets"`
	ClientID   string   `json:"client_id" binding:"required"`
	Schedules  Schedule `json:"schedules"`
}

// add validation of freq must be either monthly, weekly or daily
type Schedule struct {
	TimeSquare []TimeSquare `json:"time_square"`
	Frequency  string       `json:"frequency" binding:"required,oneof=Monthly Weekly Daily"`
}

type TimeSquare struct {
	Day     string `json:"day"`
	Time    string `json:"time"`
	Expiry  string `json:"expiry"`
	Weekday string `json:"weekday"`
}

// add not empty validation
type ModifySipRequest struct {
	ActionType string   `json:"action_type" binding:"required"`
	Name       string   `json:"name" binding:"required"`
	Baskets    []string `json:"baskets"`
	ID         string   `json:"id" binding:"required"`
	ClientID   string   `json:"client_id" binding:"required"`
	Schedules  Schedule `json:"schedules"`
}

// one of is not working why?
type UpdateSipStatusRequest struct {
	ID     string `json:"id" binding:"required"`
	Action string `json:"action" binding:"required,oneof=Active Paused"`
}
