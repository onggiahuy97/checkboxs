package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8081", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/grid", getGrid)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	test_connection()
}

func test_connection() {
	s := "ws://localhost:8081"
	ws, _, err := websocket.DefaultDialer.Dial(s, nil)
	if err != nil {
		fmt.Println("Error")
	}

	defer ws.Close()

	// send 10 randome json messages as { "row": random(1, 15), "col": random(1, 15) }
	for i := 0; i < 10; i++ {
		err := ws.WriteJSON(map[string]int{"row": rand.Intn(15) + 1, "col": rand.Intn(15) + 1})
		if err != nil {
			fmt.Println("Error")
		}
	}

}
