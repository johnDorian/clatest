package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type APIClient struct {
	Client     *http.Client
	RequestURI string
}

type RawData struct {
	Cases     map[string]int `json:"cases"`
	Deaths    map[string]int `json:"deaths"`
	Recovered map[string]int `json:"recovered"`
}

type APIResponse struct {
	Country    string   `json:"country"`
	Province   []string `json:"province"`
	RawData    RawData  `json:"timeline"`
	TimeSeries TimeSeries
}

var (
	ErrorBadDateFormat = "Incorrect Date Format"
)

func NewClient(RequestURI string) *APIClient {
	httpClient := &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: 10 * time.Second,
		},
	}
	return &APIClient{Client: httpClient, RequestURI: RequestURI}
}

func (c *APIClient) Get(country string, from, to time.Time, latest bool) (APIResponse, error) {
	var data APIResponse
	totalDays := calcDays(from)

	resp, err := c.Client.Get(fmt.Sprintf(c.RequestURI, country, totalDays))

	if err != nil {
		return data, err
	}

	if resp.StatusCode != 200 {
		return data, parseErrorMessage(resp)
	}

	d := json.NewDecoder(resp.Body)
	err = d.Decode(&data)

	if err != nil {
		return data, err
	}
	data.FormatResponse(from, to, latest)

	return data, nil

}

func (r *APIResponse) FormatResponse(from, to time.Time, latest bool) error {
	var timeSeries TimeSeries
	for date := range r.RawData.Cases {
		formattedTime, err := cleanReturnedDate(date)
		if err != nil {
			return err
		}
		timeSeries.Data = append(timeSeries.Data, Day{
			Country:   r.Country,
			Date:      formattedTime,
			Cases:     r.RawData.Cases[date],
			Deaths:    r.RawData.Deaths[date],
			Recovered: r.RawData.Recovered[date],
		})
	}

	r.TimeSeries = timeSeries
	r.TimeSeries.Order()
	r.TimeSeries.Filter(from, to, latest)
	return nil

}

func cleanReturnedDate(date string) (time.Time, error) {
	dateParts := strings.Split(date, "/")
	if len(dateParts) != 3 {
		return time.Time{}, errors.New(ErrorBadDateFormat)
	}
	cleanDate := fmt.Sprintf("20%v-%v%v-%v%v", dateParts[2], strings.Repeat("0", 2-len(dateParts[0])), dateParts[0], strings.Repeat("0", 2-len(dateParts[1])), dateParts[1])
	formattedTime, err := time.Parse("2006-01-02", cleanDate)

	return formattedTime, err

}

func calcDays(from time.Time) int {
	return int(time.Now().Sub(from).Hours()/24) + 1
}

func parseErrorMessage(resp *http.Response) error {
	var errMessage struct {
		Message string `json:"message"`
	}
	errdecoder := json.NewDecoder(resp.Body)
	err := errdecoder.Decode(&errMessage)
	if err != nil {
		return err
	}
	return errors.New(errMessage.Message)
}
