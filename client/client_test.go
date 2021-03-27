/*
Copyright © 2021 Jason Lessels <jlessels@gmail.com>

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
package client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var responseData = `{
	"country":"Australia",
	"province":[
		"australian capital territory",
		"new south wales",
		"northern territory",
		"queensland",
		"south australia",
		"tasmania",
		"victoria",
		"western australia"
		],
		"timeline":{
			"cases":{
				"3/16/21":29154,
				"3/17/21":29166,
				"3/18/21":29183,
				"3/19/21":29192,
				"3/20/21":29196,
				"3/21/21":29206,
				"3/22/21":29211,
				"3/23/21":29221,
				"3/24/21":29230,
				"3/25/21":29239
				},
			"deaths":{
				"3/16/21":909,
				"3/17/21":909,
				"3/18/21":909,
				"3/19/21":909,
				"3/20/21":909,
				"3/21/21":909,
				"3/22/21":909,
				"3/23/21":909,
				"3/24/21":909,
				"3/25/21":909
				},
			"recovered":{
				"3/16/21":22960,
				"3/17/21":22962,
				"3/18/21":22963,
				"3/19/21":22965,
				"3/20/21":22966,
				"3/21/21":22971,
				"3/22/21":22977,
				"3/23/21":22982,
				"3/24/21":22988,
				"3/25/21":22991
			}
		}
	}
`

func TestGet(t *testing.T) {
	assert := assert.New(t)
	responseJSON := new(APIResponse)
	err := json.Unmarshal([]byte(responseData), responseJSON)
	if err != nil {
		log.Fatalln(err)
	}

	tests := []struct {
		country        string
		from           time.Time
		to             time.Time
		latest         bool
		expectedStatus int
		expectedData   interface{}
	}{
		{
			country:        "australia",
			from:           time.Now(),
			to:             time.Now(),
			latest:         true,
			expectedStatus: 200,
			expectedData: []Day{
				{
					Country:   "Australia",
					Date:      time.Date(2021, 3, 25, 0, 0, 0, 0, time.UTC),
					Cases:     29239,
					Deaths:    909,
					Recovered: 22991,
				},
			},
		},
		{
			country:        "australia",
			from:           time.Date(2021, 3, 24, 0, 0, 0, 0, time.UTC),
			to:             time.Date(2021, 3, 24, 0, 0, 0, 0, time.UTC),
			latest:         false,
			expectedStatus: 200,
			expectedData: []Day{
				{
					Country:   "Australia",
					Date:      time.Date(2021, 3, 24, 0, 0, 0, 0, time.UTC),
					Cases:     29230,
					Deaths:    909,
					Recovered: 22988,
				},
			},
		},
		{
			country:        "azzz",
			expectedStatus: 502,
			expectedData:   "country not found",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, _ := json.Marshal(responseJSON)
		if !strings.Contains(strings.ToLower(r.URL.Path), "australia") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"message":"country not found"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL + "/%v%v")

	for _, test := range tests {
		apiResponse, err := client.Get(test.country, test.from, test.to, test.latest)
		switch test.expectedStatus {
		case 200:
			assert.Equal(apiResponse.Country, responseJSON.Country)
			assert.Equal(apiResponse.RawData.Cases, responseJSON.RawData.Cases)
			assert.Equal(apiResponse.RawData.Deaths, responseJSON.RawData.Deaths)
			assert.Equal(apiResponse.RawData.Recovered, responseJSON.RawData.Recovered)
			assert.Equal(apiResponse.TimeSeries.Data, test.expectedData)
		case 502:
			assert.Equal(err.Error(), test.expectedData)
		}

	}

}
func TestCleanReturnedDate(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		in          string
		expected    time.Time
		expectedErr error
	}{
		{in: "21-1-1", expected: time.Time{}, expectedErr: ErrorBadDateFormat},
		{in: "1/1/21", expected: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), expectedErr: nil},
	}

	for _, test := range tests {
		formatted, err := cleanReturnedDate(test.in)
		if test.expectedErr != nil {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
		assert.Equal(test.expected, formatted)
	}

}

func TestCalcDays(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		in       time.Time
		expected int
	}{
		{in: time.Now(), expected: 1},
		{in: time.Now().Add(-24 * time.Hour), expected: 2},
	}
	for _, test := range tests {
		days := calcDays(test.in)
		assert.Equal(test.expected, days)
	}
}

func TestNewClient(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		in       string
		expected string
	}{
		{in: "http://localhost:8080/%v", expected: "http://localhost:8080/%v"},
	}
	for _, test := range tests {
		client := NewClient(test.in)
		assert.Equal(test.expected, client.RequestURL)
	}

}

func TestParseErrorMessage(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		in          string
		expected    string
		expectError bool
	}{
		{in: `{"message": "hello world"}`, expected: "hello world", expectError: false},
		{in: `"good bye"}`, expected: "", expectError: true},
		{in: `{"no messsage": "hello world"}`, expected: "", expectError: false},
	}

	for _, test := range tests {
		mockResp := &http.Response{
			Body: ioutil.NopCloser(strings.NewReader(test.in)),
		}
		err := parseErrorMessage(mockResp)
		if test.expectError {
			assert.Error(err)
			continue
		}
		assert.EqualError(err, test.expected)
	}
}
