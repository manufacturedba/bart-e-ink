package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

// TODO
// - To incorporate this into a larger application, it'll need to avoid any
// panics/exits and return original errors
// - For both CLI and API, it needs to accept both the write directory and URL as
// arguments, but what is already here can remain as defaults
// - It should accept a force argument in cases where the data has gone stale
// or was downloaded partially/corrupted
// - Actually a concept of staleness would be good and incorporated with the
// existing directory check. Maybe use the timestamp of any file since I don't
// want to rely on knowing all files that should exist
const bartScheduleURL = "https://www.bart.gov/dev/schedules/google_transit.zip"
const writeDirectory = "gtfs"

// A schedule is a collection of Files that make up a BART schedule
// The archive is the original zip file that contains the schedule
// This struct is contextually used to hold an already opened file so 
// it can be closed later.
type schedule struct {
	archive *os.File
	file []*zip.File
}

// Proxy to the Close method of the zip file
func (s *schedule) Close() error {
	return s.archive.Close()
}

// Establish if the output directory already exists
func isExisting() bool {
	_, err := os.Stat(writeDirectory)
	
	if err != nil {
		return err == os.ErrExist
	}
	
	return true
}

// Fetches the latest zip from the BART website and streams it to a temporary
// file for latest use.
func latestArchive() (*os.File, error) {	
	response, err := http.Get(bartScheduleURL)
	
	if err != nil {
		return nil, err
	}
	
	defer response.Body.Close()

	tempFile, err := os.CreateTemp("", "*")
	
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(tempFile, response.Body)
	
	if err != nil {
		return nil, err
	}
	
	return tempFile, nil
}

// Returns a schedule struct from the latest BART schedule download for 
// file processing
func newSchedule() (*schedule, error) {
	archive, err := latestArchive()

	s := &schedule{archive: archive}
	
	if err != nil {
		return nil, err
	}
	
	stat, err := archive.Stat()
	
	if err != nil {
		return nil, err
	}
	
	archiveReader, err := zip.NewReader(io.ReaderAt(archive), stat.Size())
		
	if err != nil {
		return nil, err
	}
	
	s.file = archiveReader.File
	
	return s, nil
}

func main() {
	if isExisting() {
		log.Println("Existing schedule directory found. Skipping.")
		return;
	}
	
	schedule, err := newSchedule()
		
	if err != nil {
		log.Fatal("could not get latest schedule: ", err)
	}
	
	defer schedule.Close()
	
	err = os.Mkdir(writeDirectory, 0777)
	
	if err != nil && !os.IsExist(err) {
		log.Fatal("could not create directory for schedule files", err)
	}
	
	for _, file := range schedule.file {
		func() {
			archiveFile, err := file.Open()
		
			if err != nil {
				fmt.Println("could not decompress file: ", err)
			}
						
			fileName := path.Join(writeDirectory, file.Name)
			outputFile, err := os.Create(fileName)
			
			defer archiveFile.Close()
			defer outputFile.Close()
			
			if err != nil {
				fmt.Println("could not create schedule file: ", err)
			}
			
			_, err = io.Copy(outputFile, archiveFile)
			
			if err != nil {
				fmt.Println("could not write schedule file: ", err)
			}
		}()
	}
	
	log.Println("Schedule written to " + writeDirectory)
}