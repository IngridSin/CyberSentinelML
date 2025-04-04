package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type PacketUpdate struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan PacketUpdate)
	mutex     = sync.Mutex{}

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func StartWebsocket() {
	http.HandleFunc("/ws", handleConnections)

	http.HandleFunc("/api/email-stats", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == http.MethodOptions {
			return
		}
		GetEmailStatsHandler(w, r)
	})

	http.HandleFunc("/api/emails", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == http.MethodOptions {
			return
		}
		GetPaginatedEmailsHandler(w, r)
	})

	go handleMessages()

	fmt.Println("WebSocket server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer ws.Close()

	mutex.Lock()
	clients[ws] = true
	mutex.Unlock()

	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			mutex.Lock()
			delete(clients, ws)
			mutex.Unlock()
			break
		}
	}
}

func handleMessages() {
	for msg := range broadcast {
		data, err := json.Marshal(msg)
		if err != nil {
			log.Println("Failed to marshal WebSocket message:", err)
			continue
		}

		mutex.Lock()
		for conn := range clients {
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("WebSocket write error:", err)
				conn.Close()
				delete(clients, conn)
			}
		}
		mutex.Unlock()
	}
}
