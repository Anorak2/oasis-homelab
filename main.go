package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"github.com/Anorak/oasis-homelab/go-files/conway"
	"github.com/gorilla/websocket"
)
func serveFavicon(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "assets/images/favicon.ico")
}

func serveFile(w http.ResponseWriter, r *http.Request, FilePath string) {
	filepath := path.Clean(FilePath) // clean the file path from things like ..
	extension := path.Ext(filepath) // see if there is a file extension
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
	// extract only the relevant part for us
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


var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

// used to send websocket connections to where they need to go
func wsPiper(w http.ResponseWriter, r *http.Request){
	// Upgrade the HTTP connection to a WebSocket connection
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
       fmt.Println("Error upgrading:", err)
       return
    }
    defer conn.Close() // close connection when we're done
	
	// pipe the connection to the conway specific ws handler
	conway.WsHandler(conn)
}

func main() {
	// This is the main game loop for conways, runs every 5s	
	go conway.UpdateConway()

	http.HandleFunc("/favicon.ico", serveFavicon)
	http.HandleFunc("/games/ws/conway", wsPiper)
	http.HandleFunc("/games/", gameHandler)
	http.HandleFunc("/", handlehome)
	err := http.ListenAndServe(":8080", nil)
	if(err != nil){
		fmt.Println("Error starting server:",err)
	}
}
