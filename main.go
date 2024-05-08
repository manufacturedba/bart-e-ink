package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

const tripUpdateURL = "https://api.bart.gov/gtfsrt/tripupdate.aspx"

func main() {
	client := new(http.Client)
	
	req, err := http.NewRequest("GET", tripUpdateURL, nil)
	
	if err != nil {
		fmt.Println("Error fetching feed: ", err)
	}
	
	req.Header.Add("Accept-Encoding", "gzip")
	
	resp, err := client.Do(req)
	
	if err != nil {
		fmt.Println("Error fetching feed: ", err)
	}
	
	defer resp.Body.Close()
	
	fmt.Println(resp.Header.Get("Content-Encoding"))
	
	zr, err := gzip.NewReader(resp.Body)
	
	if err != nil {
		fmt.Println("Error reading feed: ", err)
	}
	
	defer zr.Close()
		
	bytestring, err := io.ReadAll(zr)
	
	if err != nil {
		fmt.Println("Error reading feed: ", err)
	}
	
	feed := gtfs.FeedMessage{}
	err = proto.Unmarshal(bytestring, &feed)
	
	fmt.Println(feed.Header.String())
	if err != nil {
		fmt.Println("Error parsing feed: ", err)
	}
	
	for _, entity := range feed.Entity {
        tripUpdate := entity.GetTripUpdate()
        trip := tripUpdate.GetTrip()
        fmt.Printf("Trip ID: %s\n", trip.GetTripId())
    }
}
