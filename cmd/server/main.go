package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/bart-e-ink/internal/transit"
)


func handler(w http.ResponseWriter, r *http.Request) {
	displayRows := transit.Rows()
	encodedRows, _ := json.Marshal(displayRows)
    fmt.Fprintf(w, "%s", encodedRows)
}

func main() {
	port := flag.String("port", "8080", "Listening port for the server")
	
	flag.Parse()
    
	http.HandleFunc("/", handler)
	log.Println("Server started on port " + *port)
	log.Println("Press Ctrl+C to stop the server")
	log.Fatal(http.ListenAndServe(":" + *port, nil))
}