package roundtrip

import (
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/rs/zerolog"
)

const (
	ContentTypeJSON          = "application/json"
	ContentTypeXML           = "application/xml"
	ContentTypeHTML          = "text/html"
	ContentTypePlainText     = "text/plain"
	ContentTypeForm          = "application/x-www-form-urlencoded"
	ContentTypeMultipartForm = "multipart/form-data"
)

func selectLogType(contentType string, logger *zerolog.Event, bodyData []byte, key string, raw bool) error {
	if raw {
		logger.Bytes(key, bodyData)
		return nil
	}
	switch contentType {
	case ContentTypeJSON:
		logger.RawCBOR(key, bodyData)
	case ContentTypeXML:
	case ContentTypeHTML:

	case ContentTypePlainText:
		logger.Bytes(key, bodyData)

	case ContentTypeForm:
		logger.RawCBOR(key, bodyData)
	case ContentTypeMultipartForm:
	default:
		logger.Str(key, string(bodyData))
	}
	return nil
}

func logRespBodyData(resp *http.Response, logger *zerolog.Event, isDumpBody bool, rawBody *bool) error {
	o, err := httputil.DumpResponse(resp, isDumpBody)
	if err != nil {
		return err
	}
	err = selectLogType(resp.Header.Get("Content-Type"), logger, o, "response_data", *rawBody)
	if err != nil {
		return err
	}
	return nil
}

func AddHeader(key, value string) Middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		return customRoundTripper(func(req *http.Request) (*http.Response, error) {
			header := req.Header
			if header == nil {
				header = make(http.Header)
			}

			header.Set(key, value)
			return rt.RoundTrip(req)
		})
	}
}

// It is an HTTP RoundTripper that enables logging request and response data in JSON format using the zerolog library: https://github.com/rs/zerolog. If rawRespData is enabled, all response body data returned from upstream will not be encoded (CBOR encoded).
func LoggingResp(logger *zerolog.Logger, printResp bool, printDuration bool, rawRespData bool) Middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		return customRoundTripper(func(req *http.Request) (resp *http.Response, err error) {
			startTime := time.Now()
			defer func() {
				var log = logger.Info()
				log.Str("method", req.Method).
					Stringer("url", req.URL).
					Str("proto", req.Proto).
					Stringer("url", req.URL).
					Int("status", resp.StatusCode).
					Int64("response_byte", resp.ContentLength)

				if printDuration {
					log.Dur("duration", time.Since(startTime))
				}
				if printResp {
					if resp.Body != nil {
						logRespBodyData(resp, log, true, &rawRespData)
					} else {
						logRespBodyData(resp, log, false, nil)
					}
				}
				log.Msg("")
			}()
			return rt.RoundTrip(req)

		})
	}
}
