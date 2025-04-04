package database

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"goServer/config"
	"goServer/models"
)

const (
	batchSize     = 10 // Smaller batch
	flushInterval = 500 * time.Millisecond
)

var (
	dbPools       = make(map[string]*pgxpool.Pool)
	packetChannel = make(chan *models.PacketFlow, 1000) // Buffered channel
	wg            sync.WaitGroup
)

func GetPool(jobName string) *pgxpool.Pool {
	pool, ok := dbPools[jobName]
	if !ok {
		log.Fatalf("No DB pool found for job: %s", jobName)
	}
	return pool
}

// ConnectDB initializes PostgreSQL connection
func ConnectDB(jobName string, localPort int, dbName string) {
	connString := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable",
		config.DBUser, config.DBPassword, localPort, dbName)

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to connect to database for %s: %v", jobName, err)
	}

	dbPools[jobName] = pool
	fmt.Printf("[%s] Connected to PostgreSQL: %s\n", jobName, dbName)

	var version string
	err = pool.QueryRow(context.Background(), "SELECT version();").Scan(&version)
	if err != nil {
		log.Fatalf("[%s] Query failed: %v", jobName, err)
	}

	fmt.Printf("[%s] PostgreSQL Version: %s\n", jobName, version)

	// If this job is for network capture, start the worker
	if jobName == "packets" {
		wg.Add(1)
		for i := 0; i < runtime.NumCPU(); i++ {
			go batchInsertWorker()
		}
	}
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
	pool := GetPool("packets")
	if pool == nil {
		log.Println("DB pool not initialized")
		return
	}

	if len(packets) == 0 || pool == nil {
		return
	}

	ctx := context.Background()
	tx, err := pool.Begin(ctx) // Start transaction
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {

		}
	}(tx, ctx)

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
	packetChannel <- packet
}

func InsertEmail(meta *models.EmailMetadata) error {
	pool := GetPool("emails")

	query := `
		INSERT INTO emails (
			message_id, date, subject, sender, recipient,
			received, return_path, delivered_to, body,
			dkim, spf, attachments,
			prediction, winner_model, winner_probability
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, $11, $12,
			$13, $14, $15
		) ON CONFLICT (message_id) DO NOTHING;
	`

	_, err := pool.Exec(context.Background(), query,
		meta.MessageID,
		meta.Date,
		meta.Subject,
		meta.From,
		meta.To,
		meta.Received,
		meta.ReturnPath,
		meta.DeliveredTo,
		meta.Body,
		meta.DKIMSignature,
		meta.SPFResult,
		meta.Attachments,
		nil,
		nil,
		nil,
	)

	if err != nil {
		log.Printf("Failed to insert email (%s): %v\n", meta.MessageID, err)
	}
	return err
}

// CloseDB: Gracefully shuts down database connection
func CloseDB() {
	// First, close packetChannel to stop packet batch worker
	close(packetChannel)
	wg.Wait() // Wait for all background workers to finish

	// Close all db pools
	for jobName, pool := range dbPools {
		fmt.Printf("[%s] Closing database connection...\n", jobName)
		pool.Close()
	}

	fmt.Println("All database connections closed.")
}
