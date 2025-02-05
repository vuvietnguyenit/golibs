package roundtrip

import "net/http"

type Middleware func(http.RoundTripper) http.RoundTripper
type customRoundTripper func(*http.Request) (*http.Response, error)

func (rt customRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}

// A wrapper function that allows chaining multiple HTTP RoundTrips together. This enables the use of multiple middleware RoundTrips simultaneously.
func Chain(rt http.RoundTripper, middlewares ...Middleware) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}
	for _, m := range middlewares {
		rt = m(rt)
	}
	return rt
}
