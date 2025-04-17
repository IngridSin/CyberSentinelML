package websocket

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"goServer/database"
	"goServer/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

var lastNetworkBroadcast time.Time

func BroadcastNetworkStats() {
	if time.Since(lastNetworkBroadcast) < 5*time.Second {
		return
	}
	lastNetworkBroadcast = time.Now()

	stats := GenerateNetworkDashboardStats()
	broadcast <- PacketUpdate{
		Type:    "NETWORK_DASHBOARD_STATS",
		Payload: stats,
	}
}

func GetNetworkStatsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	stats := GenerateNetworkDashboardStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func GenerateNetworkDashboardStats() models.NetworkDashboardStats {
	db := database.GetPool("packets")
	ctx := context.Background()

	var total, malicious int
	var lastTime time.Time
	var last models.FlowDetail

	_ = db.QueryRow(ctx, `
		SELECT COUNT(*), COUNT(*) FILTER (WHERE xgboost_prediction = 1)
		FROM test_schema.test_network_flow;
	`).Scan(&total, &malicious)

	_ = db.QueryRow(ctx, `
		SELECT timestamp FROM test_schema.test_network_flow
		WHERE xgboost_prediction = 1
		ORDER BY timestamp DESC LIMIT 1;
	`).Scan(&lastTime)

	err := db.QueryRow(ctx, `
		SELECT flow_id, source_ip, destination_ip, protocol, xgboost_prediction, timestamp
		FROM test_schema.test_network_flow
		WHERE xgboost_prediction = 1
		ORDER BY timestamp DESC LIMIT 1;
	`).Scan(
		&last.FlowID,
		&last.SrcIP,
		&last.DstIP,
		&last.Protocol,
		&last.RiskScore,
		&last.Timestamp,
	)
	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed to get last malicious flow:", err)
	}

	return models.NetworkDashboardStats{
		Type:              "NETWORK_DASHBOARD_STATS",
		TotalFlows:        total,
		MaliciousFlows:    malicious,
		LastMaliciousTime: lastTime,
		LastMaliciousFlow: last,
	}
}

func GetPaginatedPacketsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	ctx := context.Background()
	db := database.GetPool("packets")

	// Default pagination params
	page := 1
	pageSize := 10

	// Parse query parameters
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}
	if s, err := strconv.Atoi(r.URL.Query().Get("pageSize")); err == nil && s > 0 && s <= 100 {
		pageSize = s
	}

	offset := (page - 1) * pageSize

	// Count total records
	var total int
	err := db.QueryRow(ctx, `SELECT COUNT(*) FROM test_schema.test_network_flow`).Scan(&total)
	if err != nil {
		http.Error(w, "Failed to count packets", http.StatusInternalServerError)
		log.Println("DB count error:", err)
		return
	}

	// Query paginated data
	rows, err := db.Query(ctx, `
		SELECT 
			 flow_id, source_ip, destination_ip, source_port, destination_port, protocol, 
			flow_duration, total_fwd_packets, total_backward_packets, 
			flow_bytes_per_sec, flow_packets_per_sec,
			xgboost_prediction, timestamp 
		FROM test_schema.test_network_flow
		ORDER BY timestamp DESC
		LIMIT $1 OFFSET $2
	`, pageSize, offset)
	if err != nil {
		http.Error(w, "Failed to fetch packets", http.StatusInternalServerError)
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	var packets []models.NetworkPacket
	for rows.Next() {
		var p models.NetworkPacket
		err := rows.Scan(
			&p.FlowID, &p.SourceIP, &p.DestinationIP, &p.SourcePort, &p.DestinationPort, &p.Protocol,
			&p.FlowDuration, &p.TotalFwdPackets, &p.TotalBwdPackets,
			&p.BytesPerSecond, &p.PacketsPerSecond,
			&p.Prediction, &p.Timestamp,
		)
		if err != nil {
			log.Println("Row scan error:", err)
			continue
		}
		packets = append(packets, p)
	}

	res := models.PaginatedPacketsResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Packets:  packets,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func GetMaliciousPacketsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	ctx := context.Background()
	db := database.GetPool("packets")

	// Default pagination values
	page := 1
	pageSize := 10

	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}
	if s, err := strconv.Atoi(r.URL.Query().Get("pageSize")); err == nil && s > 0 && s <= 100 {
		pageSize = s
	}
	offset := (page - 1) * pageSize

	// Count total malicious packets
	var total int
	err := db.QueryRow(ctx, `SELECT COUNT(*) FROM test_schema.test_network_flow WHERE xgboost_prediction = 1`).Scan(&total)
	if err != nil {
		http.Error(w, "Failed to count malicious packets", http.StatusInternalServerError)
		log.Println("DB count error:", err)
		return
	}

	// Fetch malicious packet rows
	rows, err := db.Query(ctx, `
		SELECT 
			timestamp, flow_id, source_ip, destination_ip, source_port, destination_port, protocol, 
			flow_duration, total_fwd_packets, total_backward_packets, 
			flow_bytes_per_sec, flow_packets_per_sec,
			xgboost_prediction 
		FROM test_schema.test_network_flow
		WHERE xgboost_prediction = 1
		ORDER BY timestamp DESC
		LIMIT $1 OFFSET $2
	`, pageSize, offset)
	if err != nil {
		http.Error(w, "Failed to query malicious packets", http.StatusInternalServerError)
		log.Println("DB query error:", err)
		return
	}
	defer rows.Close()

	var packets []models.NetworkPacket
	for rows.Next() {
		var p models.NetworkPacket
		err := rows.Scan(
			&p.Timestamp,
			&p.FlowID, &p.SourceIP, &p.DestinationIP, &p.SourcePort, &p.DestinationPort, &p.Protocol,
			&p.FlowDuration, &p.TotalFwdPackets, &p.TotalBwdPackets,
			&p.BytesPerSecond, &p.PacketsPerSecond,
			&p.Prediction,
		)
		if err != nil {
			log.Println("Row scan error:", err)
			continue
		}
		packets = append(packets, p)
	}

	res := models.PaginatedPacketsResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Packets:  packets,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
