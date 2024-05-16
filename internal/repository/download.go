package repository

import (
	"log"
	"os"
	"time"
)

const filename = "data.json"

func isFileStale(filename string, maxAge time.Duration) bool {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return true
		}
		log.Println("Error stating file:", err)
		return true
	}
	return time.Since(fileInfo.ModTime()) > maxAge
}

func Download() []byte {

	const maxAge = 24 * time.Hour

	if !isFileStale(filename, maxAge) {
		return get()
	}

	url := "https://boszhan.kz/api/label"
	data, err := send("GET", url, nil)
	if err != nil {
		log.Println("Error sending request:", err)
		return get()
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		log.Println("Error writing to file:", err)
		return get()
	}

	log.Println("Data saved to file")
	return data
}

func get() []byte {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Println("Error reading from file:", err)
		return nil
	}
	return data
}
