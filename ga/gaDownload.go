package githubarchive

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
)

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

// GetDownloadLink returns string depending
func GetDownloadLink(year, month, day, hour int) (constructURL string, filePath string) {
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
	filePath = yearString + "-"
	filePath += monthString + "-"
	filePath += dayString + "-"
	filePath += hourString
	filePath += endURL

	// construct URL
	constructURL = githubArchiveURL + filePath

	return constructURL, filePath

}

// IsValidDate returns true if matches ^[0-9]{10,10}$
// vlidates 2021010101
func IsValidDate(date string) bool {

	match, err := regexp.MatchString("^[0-9]{10,11}$", date)
	if err != nil {
		panic(err)
	}

	return match
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// check if file exists
	if _, err := os.Stat(filepath); err == nil {
		return errors.New("file already exists")

	}

	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")

	// Close the file without defer so it can happen before Rename()
	out.Close()

	// rename path
	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}

func gUnzip(filepath string) error {
	cmd := exec.Command("gunzip", "-d", filepath)
	err := cmd.Run()

	fmt.Println("done unzipping ", filepath)
	return err
}
