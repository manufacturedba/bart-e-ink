package transit

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

/*
 * Trip struct represents a trip record
 * A trip is a sequence of two or more stops along a route that occur at
 * specific times
 */
type Trip struct {
	route_id string
	service_id string
	trip_id string
	trip_headsign string
	direction_id string
	block_id string
	shape_id string
	trip_load_information string
	wheelchair_accessible string
	bikes_allowed string
}

/*
 * StopTime struct represents a stop time record
 * A stop time is a specific instance of a vehicle stopping at a location during
 * a trip
 */
type StopTime struct {
	trip_id string
	arrival_time string
	departure_time string
	stop_id string
	stop_sequence string
	stop_headsign string
	pickup_type string
	drop_off_type string
	shape_distance_traveled string
}

/*
 * Route struct represents a route record
 * A route is an end-to-end path that a vehicle travels along
 * A route may have multiple trips
 */
type Route struct {
	route_id string
	sign_code string
	paired_service_id string
}

// Time layout for parsing time strings
const timeLayout = "15:04:05"

/*
 * Reads out CSV formatted trips and returns a slice of Trip objects
 */
func getTripRecords() []Trip {
	tripFileContents, err := os.ReadFile("gtfs/trips.txt")
	
	if err != nil {
		log.Fatal(err)
	}
	
	csvReader := csv.NewReader(strings.NewReader(string(tripFileContents)))
	
	trip_records, err := csvReader.ReadAll()
	
	if err != nil {
		log.Fatal(err)
	}
	
	var trips []Trip
	
	for _, record := range trip_records {
		trip := Trip{
			route_id: record[0],
			service_id: record[1],
			trip_id: record[2],
			trip_headsign: record[3],
			direction_id: record[4],
			block_id: record[5],
			shape_id: record[6],
			trip_load_information: record[7],
			wheelchair_accessible: record[8],
			bikes_allowed: record[9],
		}
			
		trips = append(trips, trip)
	}
	
	return trips
}

/*
 * Reads out CSV formatted stop times and returns a slice of StopTime objects
 */
func getStopTimeRecords() []StopTime {
	stopTimeFileContents, err := os.ReadFile("gtfs/stop_times.txt")
	
	if err != nil {
		log.Fatal(err)
	}
	
	csvReader := csv.NewReader(strings.NewReader(string(stopTimeFileContents)))
	
	stopTime_records, err := csvReader.ReadAll()
	
	if err != nil {
		log.Fatal(err)
	}
	
	var stopTimes []StopTime
	
	for _, record := range stopTime_records {
		stopTime := StopTime{
			trip_id: record[0],
			arrival_time: record[1],
			departure_time: record[2],
			stop_id: record[3],
			stop_sequence: record[4],
			stop_headsign: record[5],
			pickup_type: record[6],
			drop_off_type: record[7],
			shape_distance_traveled: record[8],
		}
		
		stopTimes = append(stopTimes, stopTime)
	}
	
	return stopTimes
}

/*
 * Returns a slice of all trips filtered by a given Route
 */
func getTripsByRoute(route Route, trips []Trip) []Trip {
	var tripsByRoute []Trip
	
	for _, trip := range trips {
		if trip.route_id == route.route_id && trip.service_id == route.paired_service_id{
			tripsByRoute = append(tripsByRoute, trip)
		}
	}
	
	return tripsByRoute
}

/*
 * Returns a slice of all stop times for given Trips
 */
func getStopTimesByTrips(trips []Trip, stopTimes []StopTime) []StopTime {
	var stopTimesByTrip []StopTime
	
	for _, trip := range trips {
		for _, stopTime := range stopTimes {
			if stopTime.trip_id == trip.trip_id {
				stopTimesByTrip = append(stopTimesByTrip, stopTime)
			}
		}
	}
	
	return stopTimesByTrip
}

/*
 * Returns a slice of all stop times for a given stop ID
 */
func getStopTimesByStop(stop_id string, stopTimes []StopTime) []StopTime {
	var stopTimesByStop []StopTime
	
	for _, stopTime := range stopTimes {
		if stopTime.stop_id == stop_id {
			stopTimesByStop = append(stopTimesByStop, stopTime)
		}
	}
	
	return stopTimesByStop
}

/*
 * Converts a time string to a time.Time object
 *
 * If the time exceeds 24 hours, it will be converted to the next day
 * Trains do not start past midnight but may end past midnight
 */
func getTimeStringAsTime(timeString string) time.Time {
	// time may exceed 24 hours
	exceeds24 := strings.HasPrefix(timeString, "24:")
	
	if exceeds24 {
		timeString = strings.Replace(timeString, "24:", "00:", 1)
	}
	
	partialTime, err := time.Parse(timeLayout, timeString)
	
	if err != nil {
		log.Fatal(err)
	}
	
	now := time.Now()
	time := time.Date(now.Year(), now.Month(), now.Day(), partialTime.Hour(), partialTime.Minute(), partialTime.Second(), 0, time.Local)

	if exceeds24 {
		return time.AddDate(0, 0, 1)
	}
	
	return time
}

/*
 * Returns the stop time closest to the current time
 */
func getNextStopTime(currentTime time.Time, stopTimes []StopTime) StopTime {
	
	upcomingTimestampsToStopTimes := make(map[time.Time]StopTime, 0)
	
	for _, stopTime := range stopTimes {
		arrivalTime := getTimeStringAsTime(stopTime.arrival_time)
		
		if arrivalTime.After(currentTime) {
			upcomingTimestampsToStopTimes[arrivalTime] = stopTime
		}
	}
	
	if len(upcomingTimestampsToStopTimes) == 0 {
		return StopTime{}
	}
	
	var nearestArrivalTime time.Time
	
	for comparisonTimestamp := range upcomingTimestampsToStopTimes {
		comparisonDifference := comparisonTimestamp.Sub(currentTime)
		activeLowestDifference := nearestArrivalTime.Sub(currentTime)
		
		if nearestArrivalTime == (time.Time{}) || comparisonDifference < activeLowestDifference {
			nearestArrivalTime = comparisonTimestamp
		}
	}
	
	return upcomingTimestampsToStopTimes[nearestArrivalTime]
}

/*
 * Returns N upcoming stop times from the current time
 */
func getUpcomingStopTimes(initialTime time.Time, stopTimes []StopTime, count int) []StopTime {
	upcomingStopTimes := make([]StopTime, 0) // May not fill up to count
	activeTime := initialTime
	
	for i := 0; i < count; i++ {
		nextStopTime := getNextStopTime(activeTime, stopTimes)
		
		// Exhausted all upcoming stop times
		if nextStopTime == (StopTime{}) {
			break
		}
		
		upcomingStopTimes = append(upcomingStopTimes, nextStopTime)
		activeTime = getTimeStringAsTime(nextStopTime.arrival_time)
	}
	
	return upcomingStopTimes
}

/*
 * Formats the upcoming stop times in a display readable format
 * e.g. "5, 10, 15 MIN"
 */
func formatTimes(stopTimes []StopTime) string {
	
	now := time.Now()
	minutes := make([]string, len(stopTimes))
	
	for i, stopTime := range stopTimes {
		arrivalAsTime := getTimeStringAsTime(stopTime.arrival_time)
		minutesAway := arrivalAsTime.Sub(now).Minutes()
		minutes[i] = fmt.Sprintf("%.0f", minutesAway)
	}
	
	text := strings.Join(minutes, ",")
	
	return text + " MIN"
}

/*
 * Returns a slice of display rows for the transit display
 */
func Rows() []string {
	desiredStopTimeCount := 2
	
	// Routes hardcoded for now along with the active service schedule
	orangeSouth := Route{route_id: "4", sign_code: "BERRYESSA", paired_service_id: "2024_01_15-DX-MVS-Weekday-01"} // Richmond to Berryessa/North San Jose
	redSouth := Route{route_id: "7", sign_code: "MILLBRAE", paired_service_id: "2024_01_15-DX-MVS-Weekday-01"} // Richmond to Daly City/Millbrae
	northBerkeleyStopId := "NBRK" // North Berkeley

	routes := [2]Route{orangeSouth, redSouth}
	stops := [1]string{northBerkeleyStopId}
	rows := make([]string, 0)
	
	tripRecords := getTripRecords()
	stopTimeRecords := getStopTimeRecords()
	
	for _, route := range routes {
		trips := getTripsByRoute(route, tripRecords)
		tripStopTimes := getStopTimesByTrips(trips, stopTimeRecords)
		
		for _, stop := range stops {
			stopTimes := getStopTimesByStop(stop, tripStopTimes)
			
			maybeUpcomingStopTimes := getUpcomingStopTimes(time.Now(), stopTimes, desiredStopTimeCount)
			
			if len(maybeUpcomingStopTimes) == 0 {
				continue
			}
			
			formattedTime := formatTimes(maybeUpcomingStopTimes)
			
			rows = append(rows, fmt.Sprintf("%s  %s", route.sign_code, formattedTime))
			rows = append(rows, "6-CAR, 2-DOOR")
		}
	}
	
	return rows
}