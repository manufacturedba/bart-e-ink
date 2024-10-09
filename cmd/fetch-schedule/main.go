// https://www.bart.gov/dev/schedules/google_transit.zip
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


const bartScheduleURL = "https://www.bart.gov/dev/schedules/google_transit.zip"
const writeDirectory = "gtfs"

type schedule struct {
	archive *os.File
	file []*zip.File
}

func (s *schedule) Close() error {
	return s.archive.Close()
}

func isExisting() bool {
	_, err := os.Stat(writeDirectory)
	
	if err != nil {
		return err == os.ErrExist
	}
	
	return true
}

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
			
			defer archiveFile.Close()
			
			fileName := path.Join(writeDirectory, file.Name)
			bytes := make([]byte, file.UncompressedSize64)
			
			_, err = archiveFile.Read(bytes)
			
			if err != nil && err != io.EOF {
				fmt.Println("could not scan decompressed file: ", err)
			}
			
			err = os.WriteFile(fileName, bytes, 0666)
			
			if err != nil {
				fmt.Println("could not write schedule file: ", err)
			}
		}()
	}
	
	fmt.Println("schedule written to " + writeDirectory)
}