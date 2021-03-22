/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/johnDorian/clatest/client"
	"github.com/spf13/cobra"
)

var from, to, extact, format string
var latest = false
var RequestURI = "https://disease.sh/v3/covid-19/historical/%v?lastdays=%v"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "clatest",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		client := client.NewClient(RequestURI)

		fromDate := parseDate(from)
		toDate := parseDate(to)
		if extact != "" {
			exactDate := parseDate(extact)
			fromDate = exactDate
			toDate = exactDate
		}

		res, err := client.Get(strings.Join(args[:], " "), fromDate, toDate, latest)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		res.TimeSeries.Print(os.Stdout, format)

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	rootCmd.PersistentFlags().StringVarP(&from, "from", "f", yesterday, "first date to download data for")
	rootCmd.PersistentFlags().StringVarP(&to, "to", "t", today, "last date to download data for")
	rootCmd.PersistentFlags().StringVarP(&extact, "on", "o", "", "A single date to get")
	rootCmd.PersistentFlags().StringVar(&format, "format", "markdown", "Output format (markdown, csv, tab)")

}

func parseDate(date string) time.Time {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return parsedDate
}
