package cmd

import (
	"errors"
	"fmt"
	"strconv"

	ga "github.com/ralsuwaidi/GAGo/ga"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var (
	// unzip .gz if true
	unzip bool

	downloadCmd = &cobra.Command{
		Use:   "download 20210101013",
		Short: "Downloads a single GitHub Archive file",
		Long:  `Downloads a single Github Archive gz file.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a date argument")
			}
			// make sure date is valid
			if ga.IsValidDate(args[0]) {
				return nil
			}
			return fmt.Errorf("invalid date: %s", args[0])
		},
		Run: func(cmd *cobra.Command, args []string) {

			// convert string to int
			year, err := strconv.Atoi(args[0][0:4])
			if err != nil {
				panic(err)
			}
			month, err := strconv.Atoi(args[0][4:6])
			if err != nil {
				panic(err)
			}
			day, err := strconv.Atoi(args[0][6:8])
			if err != nil {
				panic(err)
			}
			hour, err := strconv.Atoi(args[0][8:])
			if err != nil {
				panic(err)
			}
			ga.InsertPost("title", "ee")
			// download
			gaURL, filePath := ga.GetDownloadLink(year, month, day, hour)
			err = ga.DownloadFile(filePath, gaURL)
			if err != nil {
				panic(err)
			}

			// unzip if flag exist
			if unzip {
				ga.GUnzip(filePath)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	downloadCmd.Flags().BoolVarP(&unzip, "unzip", "u", false, "Unzip .gz file")
}
