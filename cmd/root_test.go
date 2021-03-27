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
package cmd

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/johnDorian/clatest/client"
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

func TestRunCMD(t *testing.T) {
	assert := assert.New(t)
	responseJSON := new(client.APIResponse)
	err := json.Unmarshal([]byte(responseData), responseJSON)
	if err != nil {
		log.Fatalln(err)
	}

	tests := []struct {
		country  string
		from     string
		to       string
		exact    string
		format   string
		expected string
	}{
		{
			country:  "australia",
			from:     "2021-01-01",
			to:       "2021-01-01",
			exact:    "2021-03-25",
			format:   "markdown",
			expected: "  DATE       | CASES | DEATHS | RECOVERED  \n-------------|-------|--------|------------\n  2021-03-25 | 29239 | 909    | 22991      \n",
		},
		{
			country:  "australia",
			from:     "2021-01-01",
			to:       "2021-01-01",
			exact:    "2021-03-25",
			format:   "csv",
			expected: "Date,Cases,Deaths,Recovered\n2021-03-25,29239,909,22991\n",
		},
		{
			country:  "azzz",
			from:     "2021-01-01",
			to:       "2021-01-01",
			exact:    "2021-03-25",
			format:   "csv",
			expected: "",
		},
		{
			country:  "australia",
			from:     "2021-01-01",
			to:       "2021-01-01",
			exact:    "2021-03-25",
			format:   "badformat",
			expected: "  DATE       | CASES | DEATHS | RECOVERED  \n-------------|-------|--------|------------\n  2021-03-25 | 29239 | 909    | 22991      \n",
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

	for _, test := range tests {
		buf := new(bytes.Buffer)
		run_cmd(test.country, server.URL+"/%v%v", test.from, test.to, test.exact, test.format, buf)
		assert.Equal(test.expected, buf.String())

	}

}
