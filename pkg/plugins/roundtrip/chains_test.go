package roundtrip_test

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/vuvietnguyenit/golibs/pkg/plugins/roundtrip"
	"github.com/vuvietnguyenit/golibs/pkg/vlog"
)

func ExampleChain_logResponseandAddHeader() {
	// Declare logger you want to pass to roundtrip
	jsonLogger := vlog.NewJsonLogger(vlog.LoggerConfig{
		Level:          zerolog.InfoLevel,
		TimeFormat:     "2006-01-02T15:04:05Z07:00",
		IncludesCaller: true,
	})
	client := http.Client{
		Transport: roundtrip.Chain(nil, roundtrip.LoggingResp(
			jsonLogger,
			false, // include response data returned by upstream
			false, // print duration
			false, // Disable raw
		), // Add log response middleware to chain
			roundtrip.AddHeader("author", "vunv"), // Add add custom headers on each request to chain
		),
	}
	url := "https://jsonplaceholder.typicode.com/posts/1"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	fmt.Print(resp.Header.Get("author"))
	// Output: vunv
}
