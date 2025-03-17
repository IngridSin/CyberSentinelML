package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"goServer/config"
	"goServer/models"
)

var dbPool *pgxpool.Pool

// ConnectDB initializes PostgreSQL connection
func ConnectDB(localPort int) {
	connString := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable",
		config.DBUser, config.DBPassword, localPort, config.DBName)

	var err error
	dbPool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL through SSH tunnel!")

	// Run a test query
	var version string
	err = dbPool.QueryRow(context.Background(), "SELECT version();").Scan(&version)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	fmt.Println("PostgreSQL Version:", version)
}

// InsertPacket stores packet data into PostgreSQL
func InsertPacket(packet *models.PacketFlow) {
	if dbPool == nil {
		log.Println("Database connection is not initialized")
		return
	}

	duration := packet.EndTime.Sub(packet.StartTime).Seconds()

	// Construct full table reference with schema
	fullTable := fmt.Sprintf("%s.%s", config.DBSchema, config.DBTable)

	query := fmt.Sprintf(`
		INSERT INTO %s (
			flow_id, source_ip, source_port, destination_ip, destination_port, protocol, timestamp,
			flow_duration, total_fwd_packets, total_fwd_bytes, fin_flag_count, syn_flag_count, 
			rst_flag_count, psh_flag_count, ack_flag_count, urg_flag_count
		) VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (flow_id) DO UPDATE
		SET flow_duration = $7, total_fwd_packets = $8, total_fwd_bytes = $9, 
			fin_flag_count = $10, syn_flag_count = $11, rst_flag_count = $12, 
			psh_flag_count = $13, ack_flag_count = $14, urg_flag_count = $15;
	`, fullTable)

	_, err := dbPool.Exec(context.Background(), query, packet.FlowID, packet.SourceIP, packet.SourcePort,
		packet.DestinationIP, packet.DestinationPort, packet.Protocol, duration,
		packet.TotalFwdPackets, packet.TotalFwdBytes, packet.TCPFlags["FIN"], packet.TCPFlags["SYN"],
		packet.TCPFlags["RST"], packet.TCPFlags["PSH"], packet.TCPFlags["ACK"], packet.TCPFlags["URG"])

	if err != nil {
		log.Printf("Failed to insert/update packet: %v", err)
	}
}
