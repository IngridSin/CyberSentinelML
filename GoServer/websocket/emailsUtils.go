package websocket

import (
	"context"
	"encoding/json"
	"goServer/database"
	"goServer/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

var lastBroadcast time.Time

// BroadcastStats sends the latest dashboard stats to connected clients via WebSocket
func BroadcastStats() {
	if time.Since(lastBroadcast) < 3*time.Second {
		return
	}
	lastBroadcast = time.Now()

	stats := GenerateEmailDashboardStats()
	broadcast <- PacketUpdate{
		Type:    "EMAIL_DASHBOARD_STATS",
		Payload: stats,
	}
}

// GenerateEmailDashboardStats retrieves and returns latest dashboard stats from DB
func GenerateEmailDashboardStats() models.EmailDashboardStats {
	db := database.GetPool("emails")
	ctx := context.Background()

	var total, phishing int
	var lastTime time.Time
	var last models.PhishingDetail

	// Total email + phishing counts
	err := db.QueryRow(ctx, `
		SELECT COUNT(*), COUNT(*) FILTER (WHERE prediction = 1)
		FROM emails;
	`).Scan(&total, &phishing)
	if err != nil {
		log.Println("Failed to get email counts:", err)
	}

	// Last phishing time
	_ = db.QueryRow(ctx, `
		SELECT date FROM emails
		WHERE prediction = 1
		ORDER BY date DESC LIMIT 1;
	`).Scan(&lastTime)

	// Last phishing email details
	err = db.QueryRow(ctx, `
	SELECT subject, body, message_id, sender, recipient, return_path, dkim, spf, date
	FROM emails
	WHERE prediction = 1
	ORDER BY date DESC LIMIT 1
	`).Scan(
		&last.Subject,
		&last.Body,
		&last.MessageID,
		&last.Sender,
		&last.Recipient,
		&last.ReturnPath,
		&last.DKIM,
		&last.SPF,
		&last.Timestamp,
	)
	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed to get last phishing email:", err)
	}

	return models.EmailDashboardStats{
		Type:              "EMAIL_DASHBOARD_STATS",
		TotalEmails:       total,
		PhishingEmails:    phishing,
		LastPhishingTime:  lastTime,
		LastPhishingEmail: last,
	}
}

func GetEmailStatsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		// Handle preflight request
		return
	}

	stats := GenerateEmailDashboardStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func GetPaginatedEmailsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	ctx := context.Background()
	db := database.GetPool("emails")

	// Defaults
	page := 1
	pageSize := 10

	// Get query params
	queryPage := r.URL.Query().Get("page")
	querySize := r.URL.Query().Get("pageSize")

	if queryPage != "" {
		if p, err := strconv.Atoi(queryPage); err == nil && p > 0 {
			page = p
		}
	}
	if querySize != "" {
		if s, err := strconv.Atoi(querySize); err == nil && s > 0 && s <= 100 {
			pageSize = s
		}
	}

	offset := (page - 1) * pageSize

	// Total email count
	var total int
	err := db.QueryRow(ctx, `SELECT COUNT(*) FROM emails`).Scan(&total)
	if err != nil {
		http.Error(w, "Failed to count emails", http.StatusInternalServerError)
		log.Println("DB count error:", err)
		return
	}

	// Fetch paginated emails
	rows, err := db.Query(ctx, `
		SELECT 
		  id, message_id, subject, body, sender, prediction, date,
		  winner_probability, header_valid, risk_score,
		  top_5_words_from_nb,
		  top_5_words_from_rf,
		  top_5_words_from_xgb,
		  top_5_words_from_knn,
		  top_5_words_from_logreg
		FROM emails
		ORDER BY date DESC
		LIMIT $1 OFFSET $2;
	`, pageSize, offset)
	if err != nil {
		http.Error(w, "Failed to query emails", http.StatusInternalServerError)
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	var emails []models.EmailRow
	for rows.Next() {
		var e models.EmailRow
		err := rows.Scan(
			&e.ID,
			&e.MessageID,
			&e.Subject,
			&e.Body,
			&e.Sender,
			&e.Prediction,
			&e.Date,
			&e.WinnerProbability,
			&e.HeaderValid,
			&e.RiskScore,
			&e.TopWordsNBRaw,
			&e.TopWordsRFRaw,
			&e.TopWordsXGBRaw,
			&e.TopWordsKNNRaw,
			&e.TopWordsLogregRaw,
		)
		json.Unmarshal(e.TopWordsNBRaw, &e.TopWordsNB)
		json.Unmarshal(e.TopWordsRFRaw, &e.TopWordsRF)
		json.Unmarshal(e.TopWordsXGBRaw, &e.TopWordsXGB)
		json.Unmarshal(e.TopWordsKNNRaw, &e.TopWordsKNN)
		json.Unmarshal(e.TopWordsLogregRaw, &e.TopWordsLogreg)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}
		emails = append(emails, e)
	}

	res := models.PaginatedEmailsResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Emails:   emails,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
