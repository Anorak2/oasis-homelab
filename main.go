package main
import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)
func NewRouter() *mux.Router { 
	r := mux.NewRouter()
	r.HandleFunc("/hello", handler).Methods("GET")

	staticFileDirectory := http.Dir("assets")
	// staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	 r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "assets/test.html")
	}).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(staticFileDirectory)).Methods("GET")
	// r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

	return r
}

func main(){
	r := NewRouter()
	http.ListenAndServe(":8080", r)
}

func handler(w http.ResponseWriter, r * http.Request){
	fmt.Fprint(w, "Hello Worlds!")
}
