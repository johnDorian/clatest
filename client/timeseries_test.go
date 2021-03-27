package client

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrder(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		in       TimeSeries
		expected TimeSeries
	}{
		{
			in: TimeSeries{
				[]Day{
					{Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
					{Date: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)},
					{Date: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)},
				},
			},
			expected: TimeSeries{
				[]Day{
					{Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
					{Date: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)},
					{Date: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)},
				},
			},
		},
		{
			in: TimeSeries{
				[]Day{
					{Date: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)},
					{Date: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)},
					{Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
				},
			},
			expected: TimeSeries{
				[]Day{
					{Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
					{Date: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)},
					{Date: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)},
				},
			},
		},
	}

	for _, test := range tests {
		test.in.Order()
		assert.Equal(test.expected, test.in)
	}

}

//Filter(from, to time.Time, latest bool) {
func TestFilter(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		from     time.Time
		to       time.Time
		latest   bool
		data     TimeSeries
		expected TimeSeries
	}{
		{

			from:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			to:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			latest: false,
			data: TimeSeries{
				[]Day{
					{Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
					{Date: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)},
					{Date: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)},
				},
			},
			expected: TimeSeries{
				[]Day{
					{Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
				},
			},
		},
		{
			from:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			to:     time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
			latest: false,
			data: TimeSeries{
				[]Day{
					{Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
					{Date: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)},
					{Date: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)},
				},
			},
			expected: TimeSeries{
				[]Day{
					{Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
					{Date: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)},
				},
			},
		},
		{
			from:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			to:     time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
			latest: true,
			data: TimeSeries{
				[]Day{
					{Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
					{Date: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)},
					{Date: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)},
				},
			},
			expected: TimeSeries{
				[]Day{
					{Date: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)},
				},
			},
		},
	}

	for _, test := range tests {
		test.data.Filter(test.from, test.to, test.latest)
		assert.Equal(test.expected, test.data)
	}

}

func TestWriteCSV(t *testing.T) {
	assert := assert.New(t)
	_ = assert
	buf := new(bytes.Buffer)
	writeCSV([][]string{{"hello", "world"}}, []string{"hello", "world"}, buf)
	expected := "hello,world\nhello,world\n"
	assert.Equal(expected, buf.String())
}

func TestWriteMarkdown(t *testing.T) {
	assert := assert.New(t)
	buf := new(bytes.Buffer)
	writeMarkdown([][]string{{"hello", "world"}}, []string{"hello", "world"}, buf)
	expected := "  HELLO | WORLD  \n--------|--------\n  hello | world  \n"
	assert.Equal(expected, buf.String())
}

func TestToStringArray(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		in       TimeSeries
		expected [][]string
	}{
		{
			in: TimeSeries{
				[]Day{
					{Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Cases: 1, Deaths: 2, Recovered: 3},
					{Date: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC), Cases: 4, Deaths: 5, Recovered: 6},
					{Date: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC), Cases: 7, Deaths: 8, Recovered: 9},
				},
			},
			expected: [][]string{
				{"2021-01-01", "1", "2", "3"},
				{"2021-01-02", "4", "5", "6"},
				{"2021-01-03", "7", "8", "9"},
			},
		},
	}
	for _, test := range tests {
		assert.Equal(test.in.toStringArray(), test.expected)
	}
}

func TestPrint(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		in       TimeSeries
		format   string
		expected string
	}{
		{
			in: TimeSeries{
				[]Day{
					{Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Cases: 1, Deaths: 2, Recovered: 3},
					{Date: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC), Cases: 4, Deaths: 5, Recovered: 6},
					{Date: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC), Cases: 7, Deaths: 8, Recovered: 9},
				},
			},
			format:   "markdown",
			expected: "  DATE       | CASES | DEATHS | RECOVERED  \n-------------|-------|--------|------------\n  2021-01-01 | 1     | 2      | 3          \n  2021-01-02 | 4     | 5      | 6          \n  2021-01-03 | 7     | 8      | 9          \n",
		},
		{
			in: TimeSeries{
				[]Day{
					{Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Cases: 1, Deaths: 2, Recovered: 3},
					{Date: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC), Cases: 4, Deaths: 5, Recovered: 6},
					{Date: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC), Cases: 7, Deaths: 8, Recovered: 9},
				},
			},
			format:   "non_supported_format",
			expected: "  DATE       | CASES | DEATHS | RECOVERED  \n-------------|-------|--------|------------\n  2021-01-01 | 1     | 2      | 3          \n  2021-01-02 | 4     | 5      | 6          \n  2021-01-03 | 7     | 8      | 9          \n",
		},

		{
			in: TimeSeries{
				[]Day{
					{Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Cases: 1, Deaths: 2, Recovered: 3},
					{Date: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC), Cases: 4, Deaths: 5, Recovered: 6},
					{Date: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC), Cases: 7, Deaths: 8, Recovered: 9},
				},
			},
			format:   "csv",
			expected: "Date,Cases,Deaths,Recovered\n2021-01-01,1,2,3\n2021-01-02,4,5,6\n2021-01-03,7,8,9\n",
		},
	}
	for _, test := range tests {
		buf := new(bytes.Buffer)
		test.in.Print(buf, test.format)

		assert.Equal(test.expected, buf.String())

	}
}
