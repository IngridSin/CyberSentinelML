package database

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"goServer/config"
	"goServer/models"
)

const (
	batchSize     = 50              // Number of packets per batch
	flushInterval = 2 * time.Second // Max wait time before flushing batch
)

var (
	dbPool        *pgxpool.Pool
	packetChannel = make(chan *models.PacketFlow, 1000) // Buffered channel
	wg            sync.WaitGroup
)

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

	// Start the batch worker
	wg.Add(1)
	go batchInsertWorker()
}

// batchInsertWorker: Collects packets and inserts them in bulk
func batchInsertWorker() {
	defer wg.Done()
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	var batch []*models.PacketFlow

	for {
		select {
		case packet := <-packetChannel:
			batch = append(batch, packet)
			if len(batch) >= batchSize {
				insertBatch(batch)
				batch = nil // Reset batch after insert
			}
		case <-ticker.C:
			// Flush remaining packets on interval timeout
			if len(batch) > 0 {
				insertBatch(batch)
				batch = nil
			}
		}
	}
}

// insertBatch: Inserts packets in bulk using pgx.Batch
func insertBatch(packets []*models.PacketFlow) {
	if len(packets) == 0 || dbPool == nil {
		return
	}

	ctx := context.Background()
	tx, err := dbPool.Begin(ctx) // Start transaction
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return
	}
	defer tx.Rollback(ctx) // Ensure rollback on failure

	batch := &pgx.Batch{}

	fullTable := fmt.Sprintf("%s.%s", config.DBSchema, config.DBTable)

	query := fmt.Sprintf(`
		INSERT INTO %s (
			flow_id, source_ip, source_port, destination_ip, destination_port, protocol, timestamp,
			flow_duration, total_fwd_packets, total_fwd_bytes, fin_flag_count, syn_flag_count, 
			rst_flag_count, psh_flag_count, ack_flag_count, urg_flag_count
		) VALUES (
			$1, $2, $3, $4, $5, $6, NOW(), $7, $8, $9, $10, $11, $12, $13, $14, $15
		)
		ON CONFLICT (flow_id) DO UPDATE SET 
			flow_duration = EXCLUDED.flow_duration, 
			total_fwd_packets = EXCLUDED.total_fwd_packets, 
			total_fwd_bytes = EXCLUDED.total_fwd_bytes, 
			fin_flag_count = EXCLUDED.fin_flag_count, 
			syn_flag_count = EXCLUDED.syn_flag_count, 
			rst_flag_count = EXCLUDED.rst_flag_count, 
			psh_flag_count = EXCLUDED.psh_flag_count, 
			ack_flag_count = EXCLUDED.ack_flag_count, 
			urg_flag_count = EXCLUDED.urg_flag_count;
	`, fullTable)

	// Add packets to batch
	for _, packet := range packets {
		duration := packet.EndTime.Sub(packet.StartTime).Seconds()
		batch.Queue(query,
			packet.FlowID, packet.SourceIP, packet.SourcePort,
			packet.DestinationIP, packet.DestinationPort, packet.Protocol, duration,
			packet.TotalFwdPackets, packet.TotalFwdBytes, packet.TCPFlags["FIN"],
			packet.TCPFlags["SYN"], packet.TCPFlags["RST"], packet.TCPFlags["PSH"],
			packet.TCPFlags["ACK"], packet.TCPFlags["URG"])
	}

	// Execute batch
	br := tx.SendBatch(ctx, batch)
	err = br.Close()
	if err != nil {
		log.Printf("Batch insert failed: %v", err)
		return
	}

	err = tx.Commit(ctx) // Commit transaction
	if err != nil {
		log.Printf("Transaction commit failed: %v", err)
	}
}

// InsertPacket InsertPacket: Sends packet to batch queue
func InsertPacket(packet *models.PacketFlow) {
	if dbPool == nil {
		log.Println("Database connection is not initialized")
		return
	}
	packetChannel <- packet
}

// CloseDB: Gracefully shuts down database connection
func CloseDB() {
	close(packetChannel) // Close channel before waiting
	wg.Wait()            // Wait for batch worker to finish
	dbPool.Close()
	fmt.Println("Database connection closed.")
}
