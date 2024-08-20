package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	grid [15][15]bool

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true }, // Adjust this for CORS if necessary
	}

	clients   []*websocket.Conn
	clientsMu sync.Mutex
)

func main() {
	fmt.Println("Starting server on :8080")
	http.HandleFunc("/ws", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()

	clientsMu.Lock()
	clients = append(clients, ws)
	clientsMu.Unlock()

	go sendPings(ws)

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			removeClient(ws)
			return
		}

		var index struct {
			Row int `json:"row"`
			Col int `json:"col"`
		}
		if err := json.Unmarshal(msg, &index); err != nil {
			fmt.Println("Unmarshal error:", err)
			continue
		}

		if index.Row >= 0 && index.Row < 10 && index.Col >= 0 && index.Col < 10 {
			grid[index.Row][index.Col] = !grid[index.Row][index.Col]
		}

		sendGridToAllClients()
	}
}

func sendGridToAllClients() {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	gridJSON, err := json.Marshal(grid)
	if err != nil {
		fmt.Println("Marshal error:", err)
		return
	}

	for _, client := range clients {
		err := client.WriteMessage(websocket.TextMessage, gridJSON)
		if err != nil {
			fmt.Println("WriteMessage error:", err)
			removeClient(client)
		}
	}
}

func removeClient(ws *websocket.Conn) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for i, client := range clients {
		if client == ws {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

func sendPings(ws *websocket.Conn) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
			fmt.Println("Ping error:", err)
			removeClient(ws)
			return
		}
	}
}
