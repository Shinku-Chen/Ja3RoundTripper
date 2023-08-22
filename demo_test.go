package Ja3RoundTripper

import (
	"io"
	"net/http"
	"testing"
)

func TestRoundTripper(t *testing.T) {
	c, _ := (&http.Client{
		Transport: &Ja3RoundTripper{

			Ja3: "771,4865-4866-4867-49196-49195-52393-49200-49199-52392-49162-49161-49172-49171-157-156-53-47-49160-49170-10,0-23-65281-10-11-16-5-13-18-51-45-43-27-21,29-23-24-25,0",
			//HeaderOrder:     nil,
			//	Timeout:         0,
			//DisableRedirect: false,
		},
	}).Get("https://tools.scrapfly.io/api/fp/ja3")
	AA, _ := io.ReadAll(c.Body)
	println(string(AA))
}
