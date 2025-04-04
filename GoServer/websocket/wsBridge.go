package websocket

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WSBridge struct {
	Pool      *pgxpool.Pool
	Conn      *websocket.Conn
	Interval  time.Duration
	Query     string
	Transform func(row map[string]interface{}) interface{}
}

// Start begins polling the DB and sending updates over WebSocket
func (b *WSBridge) Start(ctx context.Context) {
	ticker := time.NewTicker(b.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rows, err := b.Pool.Query(ctx, b.Query)
			if err != nil {
				log.Printf("DB query error: %v", err)
				continue
			}

			cols := rows.FieldDescriptions()
			for rows.Next() {
				rowData := make(map[string]interface{})
				values, err := rows.Values()
				if err != nil {
					log.Printf("Failed to get row values: %v", err)
					continue
				}
				for i, col := range cols {
					rowData[string(col.Name)] = values[i]
				}

				payload := rowData
				if b.Transform != nil {
					//payload = b.Transform(rowData)
				}

				jsonData, err := json.Marshal(payload)
				if err != nil {
					log.Printf("Failed to marshal row: %v", err)
					continue
				}

				err = b.Conn.WriteMessage(websocket.TextMessage, jsonData)
				if err != nil {
					log.Printf("WebSocket write error: %v", err)
					return // exit if client disconnects
				}
			}
			rows.Close()

		case <-ctx.Done():
			return
		}
	}
}
