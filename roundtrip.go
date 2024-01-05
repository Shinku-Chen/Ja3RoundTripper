package Ja3RoundTripper

import (
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Ja3RoundTripper struct {
	Proxy           func(*http.Request) (*url.URL, error)
	UserAgent       string
	Ja3             string
	HeaderOrder     []string `json:"headerOrder"`
	Timeout         int      `json:"timeout"`
	DisableRedirect bool     `json:"disableRedirect"`
}

func (receiver *Ja3RoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	client := cycletls.Init()

	options := cycletls.Options{
		Ja3:             receiver.Ja3,
		UserAgent:       receiver.UserAgent,
		HeaderOrder:     receiver.HeaderOrder,
		Timeout:         receiver.Timeout,
		DisableRedirect: receiver.DisableRedirect,
		Headers:         make(map[string]string),
	}

	if req.Body != nil {
		if b, e := io.ReadAll(req.Body); e == nil {
			options.Body = string(b)
		} else {
			return nil, e
		}
	}
	if receiver.Proxy != nil {
		if p, e := receiver.Proxy(req); e == nil && p != nil {
			options.Proxy = p.String()
		}
	}

	for s, hs := range req.Header {
		options.Headers[s] = strings.Join(hs, ",")
	}

	response, err := client.Do(req.URL.String(), options, req.Method)
	if err != nil {
		return nil, err
	}

	headers := make(http.Header)

	for s, s2 := range response.Headers {
		headers.Set(s, s2)
	}
	length := int64(0)

	if c := response.Headers["Content-Length"]; len(c) == 1 {
		if c64, err := strconv.ParseInt(c, 10, 64); err == nil {
			length = c64
		}
	} else {
		length = int64(len(response.Body))
	}
	return &http.Response{
		Status:        "",
		StatusCode:    response.Status,
		Header:        headers,
		Body:          io.NopCloser(strings.NewReader(response.Body)),
		ContentLength: length,
		//TransferEncoding: nil,
		//Close:            false,
		//Uncompressed:     false,
		//Trailer:          nil,
		//Request:          nil,
		//TLS:              nil,
		//Proto:            "",
		//ProtoMajor:       0,
		//ProtoMinor:       0,
	}, nil
}
