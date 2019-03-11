package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	chrm "go-screenshot/chrome"
	"os"
)

var (
	chrome     chrm.Engine
	db         chrm.Storage
	dbLocation string

	// 'global' flags
	waitTimeout   int
	resolution    string
	chromeTimeout int
	chromePath    string
	userAgent     string

	// screenshot command flags
	screenshotURL         string
	screenshotDestination string
)

var rootCmd = &cobra.Command{
	Use:   "go-screenshot",
	Short: "go-screenshot is a coding challenge for Detectify",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		// Init Google Chrome
		chrome = chrm.ChromeEngine(resolution, chromeTimeout, chromePath, userAgent)
		chrome.Setup()

		// Setup the destination directory
		if err := chrome.SetScreenshotPath(screenshotDestination); err != nil {
			fmt.Println("Error in setting destination screenshot path.")
		}
		// open the database
		db = chrm.Storage{}
		db.Open(dbLocation)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	// Global flags
	rootCmd.PersistentFlags().IntVarP(&waitTimeout, "timeout", "T", 3, "Time in seconds to wait for a HTTP connection")
	rootCmd.PersistentFlags().IntVarP(&chromeTimeout, "chrome-timeout", "", 90, "Time in seconds to wait for Google Chrome to finish a screenshot")
	rootCmd.PersistentFlags().StringVarP(&chromePath, "chrome-path", "", "", "Full path to the Chrome executable to use. By default, gowitness will search for Google Chrome")
	rootCmd.PersistentFlags().StringVarP(&userAgent, "user-agent", "", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.50 Safari/537.36", "Alernate UserAgent string to use for Google Chrome")
	rootCmd.PersistentFlags().StringVarP(&resolution, "resolution", "R", "1440,900", "screenshot resolution")
	rootCmd.PersistentFlags().StringVarP(&screenshotDestination, "destination", "d", "./output-storage/images", "Destination directory for screenshots")
	rootCmd.PersistentFlags().StringVarP(&dbLocation, "db", "D", "./output-storage/db", "Destination for the gowitness database")
}
