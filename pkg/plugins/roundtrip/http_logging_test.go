package roundtrip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/vuvietnguyenit/golibs/pkg/vlog"
)

func TestGetHTTPLoggingHandlerPlugin(t *testing.T) {
	jsonLogger := vlog.NewJsonLogger(vlog.LoggerConfig{
		Level:          zerolog.InfoLevel,
		TimeFormat:     "2006-01-02T15:04:05Z07:00",
		IncludesCaller: true,
	})
	client := http.Client{
		Transport: Chain(nil, LoggingResp(jsonLogger.LegacyHandler, true, true, true)),
	}

	req, err := http.NewRequest(
		http.MethodGet,
		"http://www.google.com/robots.txt",
		nil,
	)
	if err != nil {
		t.Error(err)
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	_, err = io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		t.Log(res)
	}
	if err != nil {
		t.Error(err)
	}
	t.Logf("Content-Type: %s", res.Header.Get("Content-Type"))

}

func TestPostJsonResp(t *testing.T) {
	jsonLogger := vlog.NewJsonLogger(vlog.LoggerConfig{
		Level:          zerolog.InfoLevel,
		TimeFormat:     "2006-01-02T15:04:05Z07:00",
		IncludesCaller: true,
	})

	client := http.Client{
		Transport: Chain(nil, LoggingResp(jsonLogger.LegacyHandler, true, true, true)),
	}
	url := "https://httpbin.org/post"
	type RequestData struct {
		Name string `json:"name"`
	}
	// Create the data to send
	data := RequestData{Name: "John Doe"}
	// Marshal the data into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Error(err)
	}
	// Create a new POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Check the response status
	fmt.Println("Response Status:", resp.Status)

	// Optionally, you can read the response body
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	t.Logf("Content-Type: %s", resp.Header.Get("Content-Type"))

}

func TestPostFormDataReq(t *testing.T) {
	jsonLogger := vlog.NewJsonLogger(vlog.LoggerConfig{
		Level:          zerolog.InfoLevel,
		TimeFormat:     "2006-01-02T15:04:05Z07:00",
		IncludesCaller: true,
	})

	client := http.Client{
		Transport: Chain(nil, LoggingResp(jsonLogger.LegacyHandler, true, true, true)),
	}
	type RequestData struct {
		Name  string `form:"name"`
		Email string `form:"email"`
	}
	// Create a URL-encoded form
	formData := url.Values{}
	formData.Set("name", "John Doe")
	formData.Set("email", "john.doe@example.com")

	// Create an HTTP client

	// Create the POST request with form data
	req, err := http.NewRequest("POST", "https://httpbin.org/post", bytes.NewBufferString(formData.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}
	// Set the appropriate content type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		os.Exit(1)
	}
	// defer resp.Body.Close()

	// Print the response status
	fmt.Println("Response Status:", resp.Status)
	// Create the data to send
}
