package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"github.com/Anorak/oasis-homelab/go-files/conway"
	"time"
	"github.com/gorilla/websocket"
	"sync"
	"encoding/json"
)


func serveFile(w http.ResponseWriter, r *http.Request, FilePath string) {
	filepath := path.Clean(FilePath)
	extension := path.Ext(filepath) 
	switch extension {
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		default:
			w.Header().Set("Content-Type", "text/html")
			filepath += ".html"
		}

	// Check if the file exists
    if _, err := os.Stat(filepath); os.IsNotExist(err) {
        http.NotFound(w, r)
        return
    }
	// actually serve the file
	http.ServeFile(w, r, filepath)
}

func gameHandler(w http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/games/"):]
	// construct filepath
	filepath := "assets/games/" + title
	serveFile(w, req, filepath)
}

func handlehome(w http.ResponseWriter, req *http.Request){
	// construct filepath
	filepath := "assets/games/conway" 
	serveFile(w, req, filepath)
}

func updateConway(){
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			//apply the rules to the board once
			conway.UpdateBoard()
			message := conway.GetBoard()
			// message each connection
			con_mu.Lock()
			for conn := range activeConnections {
				if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					fmt.Println("Error sending periodic message:", err)
					conn.Close() // Close connection if it fails to send a message
					delete(activeConnections, conn) // Remove broken connection
				}
			}
			con_mu.Unlock()
		}
	}
}

var(
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	activeConnections = make(map[*websocket.Conn]bool)
	con_mu sync.Mutex
)

type jsonData struct {
	Row int
	Column int
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
       fmt.Println("Error upgrading:", err)
       return
    }
    defer conn.Close()

	con_mu.Lock()
	activeConnections[conn] = true
	con_mu.Unlock()

	// Listen for incoming messages from this specific connection
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		log.Printf("Received: %s\n", message)
		// Echo the message back to the client
		if string(message) == "req-b"{
			boardState := conway.GetBoard()
			if err := conn.WriteMessage(websocket.TextMessage, []byte(boardState)); err != nil {
				log.Println("Error writing message:", err)
				break
			}
		} else {
			var payload jsonData
			err := json.Unmarshal([]byte(message), &payload)
			if err != nil {
				log.Printf("Error unmarshalling json")
			} else{
				conway.ChangeSquare(int(payload.Row), int(payload.Column))
			}
		}
	}

	// Once the connection is closed, remove it from the active connections map
	con_mu.Lock()
	delete(activeConnections, conn)
	con_mu.Unlock()
}

func main() {
	go updateConway()
	http.HandleFunc("/games/ws/conway", wsHandler)
	http.HandleFunc("/games/", gameHandler)
	http.HandleFunc("/", handlehome)
	err := http.ListenAndServe(":8080", nil)
	if(err != nil){
		fmt.Println("Error starting server:",err)
	}
}
