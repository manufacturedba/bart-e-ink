package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bart-e-ink/internal/transit"
)


func handler(w http.ResponseWriter, r *http.Request) {
	displayRows := transit.Rows()
	encodedRows, _ := json.Marshal(displayRows)
    fmt.Fprintf(w, "%s", encodedRows)
}

func main() {
	val, ok := os.LookupEnv("PORT")
	var port string
	
	if !ok {
		port = "8080"
	} else {
		port = val
	}
	
	http.HandleFunc("/", handler)
	log.Println("Server started on port " + port)
	log.Println("Press Ctrl+C to stop the server")
	log.Fatal(http.ListenAndServe(":" + port, nil))
}