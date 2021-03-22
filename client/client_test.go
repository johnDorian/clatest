package client

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCleanReturnedDate(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		in          string
		expected    time.Time
		expectedErr string
	}{
		{in: "21-1-1", expected: time.Time{}, expectedErr: ErrorBadDateFormat},
		{in: "1/1/21", expected: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), expectedErr: ""},
	}

	for _, test := range tests {
		formatted, err := cleanReturnedDate(test.in)
		if test.expectedErr == "" {
			assert.NoError(err)
		} else {
			assert.EqualError(err, test.expectedErr)
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
		assert.Equal(test.expected, client.RequestURI)
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
