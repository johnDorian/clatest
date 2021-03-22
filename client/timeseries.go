package client

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/olekukonko/tablewriter"
)

var (
	header = []string{"Date", "Cases", "Deaths", "Recovered"}
)

type Day struct {
	Country   string
	Date      time.Time
	Cases     int
	Deaths    int
	Recovered int
}

type TimeSeries struct {
	Data []Day
}

func (ts TimeSeries) Order() {
	//Now it's time to sort it
	sort.Slice(ts.Data, func(i, j int) bool {
		return ts.Data[i].Date.Before(ts.Data[j].Date)
	})
}

func (ts *TimeSeries) Filter(from, to time.Time, latest bool) {
	if latest {
		ts.Data = ts.Data[(len(ts.Data) - 1):]
		return
	}
	filteredTS := []Day{}
	for _, obs := range ts.Data {
		if (obs.Date.After(from) || obs.Date.Equal(from)) && (obs.Date.Before(to) || obs.Date.Equal(to)) {
			filteredTS = append(filteredTS, obs)

		}
	}

	ts.Data = filteredTS
}

//TODO: move this to the data
func (ts *TimeSeries) Print(output *os.File, format string) {
	exportData := ts.toStringArray()
	switch format {
	case "csv":
		writeCSV(exportData, header, output)
	default:
		writeMarkdown(exportData, header, output)
	}

}

func (ts TimeSeries) toStringArray() [][]string {
	var strData [][]string
	for _, obs := range ts.Data {
		strData = append(strData, []string{
			obs.Date.Format("2006-01-02"),
			fmt.Sprintf("%v", obs.Cases),
			fmt.Sprintf("%v", obs.Deaths),
			fmt.Sprintf("%v", obs.Recovered),
		})
	}
	return strData
}

func writeMarkdown(ts [][]string, header []string, output *os.File) {
	table := tablewriter.NewWriter(output)
	table.SetHeader(header)
	table.SetCenterSeparator("|")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false) // Set Border to false
	table.AppendBulk(ts)
	table.Render()
}

func writeCSV(ts [][]string, header []string, output *os.File) error {

	writer := csv.NewWriter(output)
	err := writer.Write(header)
	if err != nil {
		return err
	}
	err = writer.WriteAll(ts)
	if err != nil {
		return err
	}
	return nil

}
