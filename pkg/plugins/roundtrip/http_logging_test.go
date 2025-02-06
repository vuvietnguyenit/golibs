package roundtrip_test

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
	"github.com/vuvietnguyenit/golibs/pkg/plugins/roundtrip"
	"github.com/vuvietnguyenit/golibs/pkg/vlog"
)

func TestGetHTTPLoggingHandlerPlugin(t *testing.T) {
	jsonLogger := vlog.NewJsonLogger(vlog.LoggerConfig{
		Level:          zerolog.InfoLevel,
		TimeFormat:     "2006-01-02T15:04:05Z07:00",
		IncludesCaller: true,
	})
	client := http.Client{
		Transport: roundtrip.Chain(nil, roundtrip.LoggingResp(jsonLogger, true, true, false)),
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
		Transport: roundtrip.Chain(nil, roundtrip.LoggingResp(jsonLogger, true, true, false)),
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

type RequestData struct {
	Name  string `form:"name"`
	Email string `form:"email"`
}

func TestPostFormDataReq(t *testing.T) {
	jsonLogger := vlog.NewJsonLogger(vlog.LoggerConfig{
		Level:          zerolog.InfoLevel,
		TimeFormat:     "2006-01-02T15:04:05Z07:00",
		IncludesCaller: true,
	})

	client := http.Client{
		Transport: roundtrip.Chain(nil, roundtrip.LoggingResp(jsonLogger, true, true, false)),
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

func ExampleLoggingResp_withPostJson() {
	// Declare logger you want to pass to roundtrip
	jsonLogger := vlog.NewJsonLogger(vlog.LoggerConfig{
		Level:          zerolog.InfoLevel,
		TimeFormat:     "2006-01-02T15:04:05Z07:00",
		IncludesCaller: true,
	})
	client := http.Client{
		Transport: roundtrip.Chain(nil, roundtrip.LoggingResp(
			jsonLogger,
			true,  // include response data returned by upstream
			true,  // print duration
			false, // If body data is returned, print it as raw data or encode it. If this value = True, result will be not encode body data
		)),
	}

	type RequestData struct {
		Name string `json:"name"`
	}
	// Create the data to send
	data := RequestData{Name: "John Doe"}
	// Marshal the data into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
	}
	// Create a new POST request
	req, err := http.NewRequest("POST", "https://httpbin.org/post", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		os.Exit(1)
	}
	// Print the response status
	fmt.Println("Response Status:", resp.Status)
	// No // Output:
	// {"level":"info","method":"POST","url":"https://httpbin.org/post","proto":"HTTP/1.1","url":"https://httpbin.org/post","status":200,"response_byte":453,"duration":0,"response_data":"data:application/cbor;base64,SFRUUC8yLjAgMjAwIE9LDQpDb250ZW50LUxlbmd0aDogNDUzDQpBY2Nlc3MtQ29udHJvbC1BbGxvdy1DcmVkZW50aWFsczogdHJ1ZQ0KQWNjZXNzLUNvbnRyb2wtQWxsb3ctT3JpZ2luOiAqDQpDb250ZW50LVR5cGU6IGFwcGxpY2F0aW9uL2pzb24NCkRhdGU6IFRodSwgMDYgRmViIDIwMjUgMTI6MzY6MTQgR01UDQpTZXJ2ZXI6IGd1bmljb3JuLzE5LjkuMA0KDQp7CiAgImFyZ3MiOiB7fSwgCiAgImRhdGEiOiAie1wibmFtZVwiOlwiSm9obiBEb2VcIn0iLCAKICAiZmlsZXMiOiB7fSwgCiAgImZvcm0iOiB7fSwgCiAgImhlYWRlcnMiOiB7CiAgICAiQWNjZXB0LUVuY29kaW5nIjogImd6aXAiLCAKICAgICJDb250ZW50LUxlbmd0aCI6ICIxOSIsIAogICAgIkNvbnRlbnQtVHlwZSI6ICJhcHBsaWNhdGlvbi9qc29uIiwgCiAgICAiSG9zdCI6ICJodHRwYmluLm9yZyIsIAogICAgIlVzZXItQWdlbnQiOiAiR28taHR0cC1jbGllbnQvMi4wIiwgCiAgICAiWC1BbXpuLVRyYWNlLUlkIjogIlJvb3Q9MS02N2E0YWNiZS02ZTU0OTA3ZTcwZDEyMGU2NjFkNjAzMjUiCiAgfSwgCiAgImpzb24iOiB7CiAgICAibmFtZSI6ICJKb2huIERvZSIKICB9LCAKICAib3JpZ2luIjogIjEyMy4zMC4xNzUuMTIiLCAKICAidXJsIjogImh0dHBzOi8vaHR0cGJpbi5vcmcvcG9zdCIKfQo=","time":"2025-02-06T19:36:14+07:00","caller":"http_logging.go:84"}
	// Response Status: 200 OK

}

func ExampleLoggingResp_withPostFormData() {
	// Declare logger you want to pass to roundtrip
	jsonLogger := vlog.NewJsonLogger(vlog.LoggerConfig{
		Level:          zerolog.InfoLevel,
		TimeFormat:     "2006-01-02T15:04:05Z07:00",
		IncludesCaller: true,
	})
	client := http.Client{
		Transport: roundtrip.Chain(nil, roundtrip.LoggingResp(
			jsonLogger,
			true, // include response data returned by upstream
			true, // print duration
			true, // If body data is returned, print it as raw data or encode it. If this value = True, result will be not encode body data
		)),
	}

	// Create a URL-encoded form
	formData := url.Values{}
	formData.Set("name", "John Doe")
	formData.Set("email", "john.doe@example.com")

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
	// Print the response status
	fmt.Println("Response Status:", resp.Status)
	// Output:
	// {"level":"info","method":"POST","url":"https://httpbin.org/post","proto":"HTTP/1.1","url":"https://httpbin.org/post","status":200,"response_byte":487,"duration":957.81838,"response_data":"data:application/cbor;base64,SFRUUC8yLjAgMjAwIE9LDQpDb250ZW50LUxlbmd0aDogNDg3DQpBY2Nlc3MtQ29udHJvbC1BbGxvdy1DcmVkZW50aWFsczogdHJ1ZQ0KQWNjZXNzLUNvbnRyb2wtQWxsb3ctT3JpZ2luOiAqDQpDb250ZW50LVR5cGU6IGFwcGxpY2F0aW9uL2pzb24NCkRhdGU6IFRodSwgMDYgRmViIDIwMjUgMTQ6NTE6MjAgR01UDQpTZXJ2ZXI6IGd1bmljb3JuLzE5LjkuMA0KDQp7CiAgImFyZ3MiOiB7fSwgCiAgImRhdGEiOiAiIiwgCiAgImZpbGVzIjoge30sIAogICJmb3JtIjogewogICAgImVtYWlsIjogImpvaG4uZG9lQGV4YW1wbGUuY29tIiwgCiAgICAibmFtZSI6ICJKb2huIERvZSIKICB9LCAKICAiaGVhZGVycyI6IHsKICAgICJBY2NlcHQtRW5jb2RpbmciOiAiZ3ppcCIsIAogICAgIkNvbnRlbnQtTGVuZ3RoIjogIjQyIiwgCiAgICAiQ29udGVudC1UeXBlIjogImFwcGxpY2F0aW9uL3gtd3d3LWZvcm0tdXJsZW5jb2RlZCIsIAogICAgIkhvc3QiOiAiaHR0cGJpbi5vcmciLCAKICAgICJVc2VyLUFnZW50IjogIkdvLWh0dHAtY2xpZW50LzIuMCIsIAogICAgIlgtQW16bi1UcmFjZS1JZCI6ICJSb290PTEtNjdhNGNjNjgtMWViYzgyMzExZjQyZjIzNjViZDk2MjMzIgogIH0
	// sIAogICJqc29uIjogbnVsbCwgCiAgIm9yaWdpbiI6ICIxMjMuMzAuMTc1LjEyIiwgCiAgInVybCI6ICJodHRwczovL2h0dHBiaW4ub3JnL3Bvc3QiCn0K","time":"2025-02-06T21:51:20+07:00","caller":"http_logging.go:79"}
	// Response Status: 200 OK

}

func ExampleLoggingResp_withPostFormDataEncoded() {
	// Declare logger you want to pass to roundtrip
	jsonLogger := vlog.NewJsonLogger(vlog.LoggerConfig{
		Level:          zerolog.InfoLevel,
		TimeFormat:     "2006-01-02T15:04:05Z07:00",
		IncludesCaller: true,
	})
	client := http.Client{
		Transport: roundtrip.Chain(nil, roundtrip.LoggingResp(
			jsonLogger,
			true,  // include response data returned by upstream
			true,  // print duration
			false, // Disable raw
		)),
	}

	// Create a URL-encoded form
	formData := url.Values{}
	formData.Set("name", "John Doe")
	formData.Set("email", "john.doe@example.com")

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
	// Print the response status
	fmt.Println("Response Status:", resp.Status)
	// Output:
	// {"level":"info","method":"POST","url":"https://httpbin.org/post","proto":"HTTP/1.1","url":"https://httpbin.org/post","status":200,"response_byte":487,"duration":0,"response_data":"data:application/cbor;base64,SFRUUC8yLjAgMjAwIE9LDQpDb250ZW50LUxlbmd0aDogNDg3DQpBY2Nlc3MtQ29udHJvbC1BbGxvdy1DcmVkZW50aWFsczogdHJ1ZQ0KQWNjZXNzLUNvbnRyb2wtQWxsb3ctT3JpZ2luOiAqDQpDb250ZW50LVR5cGU6IGFwcGxpY2F0aW9uL2pzb24NCkRhdGU6IFRodSwgMDYgRmViIDIwMjUgMTI6MjU6MTcgR01UDQpTZXJ2ZXI6IGd1bmljb3JuLzE5LjkuMA0KDQp7CiAgImFyZ3MiOiB7fSwgCiAgImRhdGEiOiAiIiwgCiAgImZpbGVzIjoge30sIAogICJmb3JtIjogewogICAgImVtYWlsIjogImpvaG4uZG9lQGV4YW1wbGUuY29tIiwgCiAgICAibmFtZSI6ICJKb2huIERvZSIKICB9LCAKICAiaGVhZGVycyI6IHsKICAgICJBY2NlcHQtRW5jb2RpbmciOiAiZ3ppcCIsIAogICAgIkNvbnRlbnQtTGVuZ3RoIjogIjQyIiwgCiAgICAiQ29udGVudC1UeXBlIjogImFwcGxpY2F0aW9uL3gtd3d3LWZvcm0tdXJsZW5jb2RlZCIsIAogICAgIkhvc3QiOiAiaHR0cGJpbi5vcmciLCAKICAgICJVc2VyLUFnZW50IjogIkdvLWh0dHAtY2xpZW50LzIuMCIsIAogICAgIlgtQW16bi1UcmFjZS1JZCI6ICJSb290PTEtNjdhNGFhMmQtMmQ1YWI1NjQ1YTU2M2MxMTE4NzhkZGFmIgogIH0sIAogICJqc29uIjogbnVsbCwgCiAgIm9yaWdpbiI6ICIxMjMuMzAuMTc1LjEyIiwgCiAgInVybCI6ICJodHRwczovL2h0dHBiaW4ub3JnL3Bvc3QiCn0K","time":"2025-02-06T19:25:17+07:00","caller":"http_logging.go:84"}
	// Response Status: 200 OK

}
