package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {

		downloadLink, filePath := GetDownloadLink(2021, 1, 10, i)

		go DownloadFile(filePath, downloadLink)

	}

	time.Sleep(50 * time.Second)
}

// GetDownloadLink returns string depending
func GetDownloadLink(year, month, day, hour int) (string, string) {
	githubArchiveURL := "https://data.gharchive.org/"
	endURL := ".json.gz"

	//convert to string
	yearString := strconv.Itoa(year)
	monthString := strconv.Itoa(month)
	dayString := strconv.Itoa(day)
	hourString := strconv.Itoa(hour)

	if len([]rune(dayString)) == 1 {
		dayString = "0" + dayString
	}

	if len([]rune(monthString)) == 1 {
		monthString = "0" + monthString
	}

	// create file path name
	filePath := yearString + "-"
	filePath += monthString + "-"
	filePath += dayString + "-"
	filePath += hourString
	filePath += endURL

	// construct URL
	constructURL := githubArchiveURL + filePath

	return constructURL, filePath

}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("done downloading ", filepath)

	go gUnzip(filepath)

}

func gUnzip(filepath string) {
	cmd := exec.Command("gunzip", "-d", filepath)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("done unzipping ", filepath)
}
