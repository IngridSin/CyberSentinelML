package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"goServer/config"
	"log"

	"goServer/models"
)

// ConnectDB initializes PostgreSQL connection
func ConnectDB(localPort int) {
	// Connect to PostgreSQL via SSH Tunnel
	connString := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable",
		config.DBUser, config.DBPassword, localPort, config.DBName)

	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbpool.Close()

	fmt.Println("Connected to PostgreSQL through SSH tunnel!")

	// Run a test query
	var version string
	err = dbpool.QueryRow(context.Background(), "SELECT version();").Scan(&version)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	fmt.Println("PostgreSQL Version:", version)
}

// InsertPacket stores packet data into PostgreSQL
func InsertPacket(packet models.PacketFlow) {
	//query := `INSERT INTO packets (flow_id, source_ip, source_port, destination_ip, destination_port,
	//	protocol, timestamp, flow_duration, total_fwd_packets, total_bwd_packets, min_seg_size_forward, label)
	//    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	//
	//_, err := DB.Exec(context.Background(), query, packet.FlowID, packet.SourceIP,
	//	packet.SourcePort, packet.DestinationIP, packet.DestinationPort, packet.Protocol,
	//	packet.Timestamp, packet.FlowDuration, packet.TotalFwdPackets, packet.TotalBwdPackets,
	//	packet.MinSegmentSizeFwd, packet.Label)
	//
	//if err != nil {
	//	log.Println("Failed to insert packet:", err)
	//}
}
