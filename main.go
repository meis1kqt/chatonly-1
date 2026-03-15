package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)


var (
	mu sync.Mutex
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients = make(map[*websocket.Conn]bool)
)

func main(){
	fmt.Print("started")

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsHadnler)
	http.ListenAndServe(":8080", mux)
}

func wsHadnler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w,r,nil)
	if err != nil {
		return
	}	

	defer conn.Close()

	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			break
		}
		broadcast(messageType, data)
	}


	mu.Lock()
	delete(clients, conn)
	mu.Unlock()

}


func broadcast(messageType int, data []byte){
	mu.Lock()
	defer mu.Unlock()

	for conn := range clients {
		conn.WriteMessage(messageType, data)
	}
}	


