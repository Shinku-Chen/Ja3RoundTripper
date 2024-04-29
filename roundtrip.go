package Ja3RoundTripper

import (
	"errors"
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Ja3RoundTripper struct {
	cycletls.Options
	ProxyFunc func(*http.Request) (*url.URL, error)
	Jar       http.CookieJar
}

func (receiver *Ja3RoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	client := cycletls.Init()

	options := receiver.Options
	if options.Headers == nil {
		options.Headers = make(map[string]string)
	}
	if options.Cookies == nil {
		options.Cookies = make([]cycletls.Cookie, 0)
	}

	if receiver.Jar != nil {
		cookie := receiver.Jar.Cookies(req.URL)
		for _, c := range cookie {
			options.Cookies = append(options.Cookies, cycletls.Cookie{Name: c.Name, Value: c.Value})
		}
	}

	if req.Body != nil {
		if b, e := io.ReadAll(req.Body); e == nil {
			options.Body = string(b)
		} else {
			return nil, e
		}
	}
	if receiver.ProxyFunc != nil {
		if p, e := receiver.ProxyFunc(req); e == nil && p != nil {
			options.Proxy = p.String()
		}
	}

	for s, hs := range req.Header {
		options.Headers[s] = strings.Join(hs, ",")
	}

	response, err := client.Do(req.URL.String(), options, req.Method)
	if err == nil && (response.Status == 0 || (len(response.Headers) == 0 && response.Cookies == nil)) {
		err = errors.New(response.Body)
	}
	//if err != nil {
	//	return nil, err
	//}

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

	{
		delete(headers, "Content-Encoding")
		req.URL, _ = url.Parse(response.FinalUrl)
	}

	{

		responseBody := &http.Response{
			Status:        http.StatusText(response.Status),
			StatusCode:    response.Status,
			Header:        headers,
			Body:          io.NopCloser(strings.NewReader(response.Body)),
			ContentLength: length,
			//TransferEncoding: nil,
			//Close:            false,
			//Uncompressed:     false,
			//Trailer:          nil,
			Request: req,
			//TLS:              nil,
			//Proto:            "",
			//ProtoMajor:       0,
			//ProtoMinor:       0,
		}

		return responseBody, err

	}

}
