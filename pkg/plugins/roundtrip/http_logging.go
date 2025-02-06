package roundtrip

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/rs/zerolog"
	"github.com/vuvietnguyenit/golibs/pkg/httpstat"
)

const (
	ContentTypeJSON          = "application/json"
	ContentTypeXML           = "application/xml"
	ContentTypeHTML          = "text/html"
	ContentTypePlainText     = "text/plain"
	ContentTypeForm          = "application/x-www-form-urlencoded"
	ContentTypeMultipartForm = "multipart/form-data"
)

func selectLogType(contentType string, logger *zerolog.Event, bodyData []byte, key string) error {
	switch contentType {
	case ContentTypeJSON:
		logger.RawCBOR(key, bodyData)
	case ContentTypeXML:
	case ContentTypeHTML:
	case ContentTypePlainText:
	case ContentTypeForm:
		logger.RawCBOR(key, bodyData)
	case ContentTypeMultipartForm:
	default:
		return fmt.Errorf("received unknown content type")
	}
	return nil
}

func logRespBodyData(resp *http.Response, logger *zerolog.Event, isDumpBody bool) error {
	o, err := httputil.DumpResponse(resp, isDumpBody)
	if err != nil {
		return err
	}
	err = selectLogType(resp.Header.Get("Content-Type"), logger, o, "response_data")
	if err != nil {
		return err
	}
	return nil
}

// CustomTimer takes a writer and will output a request duration.
func LoggingResp(logger *zerolog.Logger, printResp bool, printDuration bool) Middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		return customRoundTripper(func(req *http.Request) (resp *http.Response, err error) {
			defer func() {
				var result httpstat.Result
				ctx := httpstat.WithHTTPStat(req.Context(), &result)
				req = req.WithContext(ctx)

				var log = logger.Info()
				log.Str("method", req.Method).
					Stringer("url", req.URL).
					Str("proto", req.Proto).
					Stringer("url", req.URL).
					Int("status", resp.StatusCode).
					Int64("response_byte", resp.ContentLength)

				if printDuration {
					result.End(time.Now())
					log.Dur("duration", result.GetTotalDur())
				}
				if printResp {
					if resp.Body != nil {
						logRespBodyData(resp, log, true)
					} else {
						logRespBodyData(resp, log, false)
					}
				}
				log.Msg("")
			}()
			return rt.RoundTrip(req)

		})
	}
}
