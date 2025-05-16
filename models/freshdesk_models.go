package models

import "time"

type FreshdeskTicketReq struct {
	Email       string `json:"email" validate:"required,email"`
	Subject     string `json:"subject" validate:"required"`
	Description string `json:"description" validate:"required"`
	Priority    int    `json:"priority" validate:"required,min=1,max=4"` // Required: Priority (1: Low, 2: Medium, 3: High, 4: Urgent)
	Status      int    `json:"status" validate:"required,min=2,max=5"`   // Required: Status (2: Open, 3: Pending, 4: Resolved, 5: Closed)

	Name         string            `json:"name,omitempty"`
	RequesterID  int64             `json:"requester_id,omitempty"`
	Phone        string            `json:"phone,omitempty"`
	CcEmails     []string          `json:"cc_emails,omitempty"`
	CustomFields map[string]string `json:"custom_fields,omitempty"`
	GroupID      int64             `json:"group_id,omitempty"`
	ProductID    int64             `json:"product_id,omitempty"`
	Tags         []string          `json:"tags,omitempty"`
	// CompanyID    int64             `json:"company_id,omitempty"` //The "multiple_user_companies" feature needs to be enabled in Freshdesk to support the company_id field
}

type FreshdeskTicketResponse struct {
	CcEmails        []string          `json:"cc_emails"`
	FwdEmails       []string          `json:"fwd_emails"`
	ReplyCcEmails   []string          `json:"reply_cc_emails"`
	TicketCcEmails  []string          `json:"ticket_cc_emails"`
	FrEscalated     bool              `json:"fr_escalated"`
	Spam            bool              `json:"spam"`
	EmailConfigID   *int              `json:"email_config_id,omitempty"`
	GroupID         *int              `json:"group_id,omitempty"`
	Priority        int               `json:"priority"`
	RequesterID     int64             `json:"requester_id"`
	ResponderID     *int64            `json:"responder_id,omitempty"`
	Source          int               `json:"source"`
	CompanyID       *int64            `json:"company_id,omitempty"`
	Status          int               `json:"status"`
	Subject         string            `json:"subject"`
	SupportEmail    *string           `json:"support_email,omitempty"`
	ToEmails        *string           `json:"to_emails,omitempty"`
	ProductID       *int64            `json:"product_id,omitempty"`
	ID              int               `json:"id"`
	Type            *string           `json:"type,omitempty"`
	DueBy           time.Time         `json:"due_by"`
	FrDueBy         time.Time         `json:"fr_due_by"`
	IsEscalated     bool              `json:"is_escalated"`
	Description     string            `json:"description"`
	DescriptionText string            `json:"description_text"`
	CustomFields    map[string]string `json:"custom_fields"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Tags            []string          `json:"tags"`
	Attachments     []interface{}     `json:"attachments"`
}
