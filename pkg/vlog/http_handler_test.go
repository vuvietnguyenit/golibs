package vlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/rs/zerolog"
)

func TestNewHTTPHandler(t *testing.T) {
	jsonLogger := NewJsonLogger(LoggerConfig{
		Level:          zerolog.DebugLevel,
		TimeFormat:     "2006-01-02T15:04:05Z07:00",
		IncludesCaller: true,
	})
	t.Log(jsonLogger)
}

func TestLogReqPOST(t *testing.T) {
	jsonLogger := NewJsonLogger(LoggerConfig{
		Level:          zerolog.DebugLevel,
		TimeFormat:     "2006-01-02T15:04:05Z07:00",
		IncludesCaller: true,
	})

	data := map[string]interface{}{
		"title":  "foo",
		"body":   "bar",
		"userId": 1,
	}
	// Marshal the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Error(err)
	}

	// Make the HTTP POST request
	url := "https://jsonplaceholder.typicode.com/posts" // Example URL

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "curl/7.88.1")
	jsonLogger.HTTPHandler.FromReq(*req).Info().Msg("req sample")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close() // Ensure the response body is closed
	jsonLogger.HTTPHandler.FromResp(*resp).Info().Msg("response sample")

	// Print the response status
	fmt.Println("Response Status:", resp.Status)
}
