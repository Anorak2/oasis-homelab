package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"github.com/Anorak/oasis-homelab/go-files/conway"
	"time"
)


func serveFile(w http.ResponseWriter, r *http.Request, FilePath string) {
	filepath := path.Clean(FilePath)
	extension := path.Ext(filepath) 
	switch extension {
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
			fmt.Println("hit")
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

func updateConway(){
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			conway.UpdateBoard()
		}
	}
}


func main() {
	go updateConway()
	http.HandleFunc("/games/conway/post", conway.HandlePost)
	http.HandleFunc("/games/conway/get", conway.HandleGet)
	http.HandleFunc("/games/", gameHandler)
	err := http.ListenAndServe(":8080", nil)
	if(err != nil){
		fmt.Println("Error starting server:",err)
	}
}
