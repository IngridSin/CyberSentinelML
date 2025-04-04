package models

import (
	"encoding/json"
	"time"
)

// EmailMetadata holds parsed email info
type EmailMetadata struct {
	MessageID     string
	Date          time.Time
	Subject       string
	From          string
	To            string
	Received      string
	ReturnPath    string
	DeliveredTo   string
	Body          string
	DKIMSignature string
	SPFResult     string
	Attachments   []string
}

type EmailDashboardStats struct {
	Type              string         `json:"type"`
	TotalEmails       int            `json:"total_emails"`
	PhishingEmails    int            `json:"phishing_emails"`
	LastPhishingTime  time.Time      `json:"last_phishing_time"`
	LastPhishingEmail PhishingDetail `json:"last_phishing_email"`
}

type PhishingDetail struct {
	Subject    string    `json:"subject"`
	Body       string    `json:"body"`
	MessageID  string    `json:"message_id"`
	Sender     string    `json:"sender"`
	Recipient  string    `json:"recipient"`
	ReturnPath string    `json:"return_path"`
	DKIM       string    `json:"dkim"`
	SPF        string    `json:"spf"`
	Timestamp  time.Time `json:"timestamp"`
}

type EmailRow struct {
	ID                int             `json:"id"`
	MessageID         string          `json:"message_id"`
	Subject           string          `json:"subject"`
	Body              string          `json:"body"`
	Sender            string          `json:"sender"`
	Prediction        int             `json:"prediction"`
	Date              time.Time       `json:"date"`
	WinnerProbability float64         `json:"winner_probability"`
	HeaderValid       bool            `json:"header_valid"`
	RiskScore         int             `json:"risk_score"`
	TopWordsNBRaw     json.RawMessage `json:"-"`
	TopWordsRFRaw     json.RawMessage `json:"-"`
	TopWordsXGBRaw    json.RawMessage `json:"-"`
	TopWordsKNNRaw    json.RawMessage `json:"-"`
	TopWordsLogregRaw json.RawMessage `json:"-"`

	TopWordsNB     map[string]float64 `json:"top_5_words_from_nb"`
	TopWordsRF     map[string]float64 `json:"top_5_words_from_rf"`
	TopWordsXGB    map[string]float64 `json:"top_5_words_from_xgb"`
	TopWordsKNN    map[string]float64 `json:"top_5_words_from_knn"`
	TopWordsLogreg map[string]float64 `json:"top_5_words_from_logreg"`
}

type PaginatedEmailsResponse struct {
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
	Total    int        `json:"total"`
	Emails   []EmailRow `json:"emails"`
}
