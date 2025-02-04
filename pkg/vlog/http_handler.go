package vlog

import (
	"io"
	"net/http"

	"github.com/rs/zerolog"
)

type HTTPPrinter struct {
	zerolog *zerolog.Logger
	req     *http.Request
	resp    *http.Response
}
type HTTPHandler struct {
	printer *HTTPPrinter
}

func (h *HTTPPrinter) printReq(e *zerolog.Event) (*zerolog.Event, error) {
	// Process body data
	requestBody, err := io.ReadAll(h.req.Body)
	if err != nil {
		return nil, err
	}
	return e.Str("method", h.req.Method).
		Str("type", "http_request").
		Stringer("url", h.req.URL).
		RawJSON("body", requestBody).
		Int("content_length", int(h.req.ContentLength)), nil
}

func (h *HTTPPrinter) printResp(e *zerolog.Event) (*zerolog.Event, error) {
	return e.Str("method", h.resp.Request.Method).
		Stringer("url", h.resp.Request.URL).
		Str("type", "http_response").
		Int("status", h.resp.StatusCode).
		Int("size", int(h.resp.ContentLength)), nil
}

func (h *HTTPPrinter) Info() *zerolog.Event {
	var e *zerolog.Event
	logger := h.zerolog.Info()
	if h.req == nil {
		e, _ = h.printResp(logger)
	} else {
		e, _ = h.printReq(logger)
	}
	return e
}

func (h *HTTPPrinter) Error() *zerolog.Event {
	return h.zerolog.Error()
}

func (h *HTTPHandler) FromReq(r http.Request) *HTTPPrinter {
	h.printer.req = &r
	h.printer.resp = nil
	return h.printer
}

func (h *HTTPHandler) FromResp(r http.Response) *HTTPPrinter {
	h.printer.resp = &r
	h.printer.req = nil
	return h.printer
}
