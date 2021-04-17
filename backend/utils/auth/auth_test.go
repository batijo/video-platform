package auth

import (
	"net/http"
	"testing"
)

type Header struct {
	req *http.Request
	out string
}

var headers = []Header{
	{addHeader("Authorization", "boi"), "Wrong header format"},
	{addHeader("Authorization", "boi fsds"), "Wrong header format"},
	{addHeader("Authorization", "boi sda\n"), "Wrong header format"},
	{addHeader("Authorization", "boi\n  "), "Wrong header format"},
	{addHeader("Authorization", "boi   d "), "Wrong header format"},
	{addHeader("Authorization", "boi    das "), "Wrong header format"},
	{addHeader("Authorization\n", "Bearer sdad.asd.asd."), "Missing 'Authorization' header"},
	{addHeader("Authorization", "Bearer sssdd sad"), "Wrong header format"},
	{addHeader("Authorization", "Bearer sssdd asd asd "), "Wrong header format"},
	{addHeader("Authorization", "Bearer sssdd\n\""), "Wrong header format"},
	{addHeader("Authorization", "Bearer sssdd \""), "Wrong header format"},
	{addHeader("", "Bearer sLsd.sxLzxs.dd"), "Missing 'Authorization' header"},
	{addHeader("Authorization", "Bearer sdad.asd.asd. \n"), "Wrong header format"},
	{addHeader("Authorization", "Bearer sdad.asd.asd.\n"), "Wrong header format"},
	// No error
	{addHeader("Authorization", "Bearer sLsd.sxLzxs.dd"), "sLsd.sxLzxs.dd"},
}

func TestParseAuthHeader(t *testing.T) {
	for _, h := range headers {
		out, err := parseAuthHeader(h.req)
		if err != nil {
			if err.Error() != h.out {
				t.Errorf("DATA: %v EXPECTED: %v, GOT: %v", h.req.Header, h.out, err.Error())
			}
		} else if out != h.out {
			t.Errorf("DATA: %v EXPECTED: %v, GOT: %v", h.req.Header, h.out, out)
		}
	}
}

func addHeader(name, value string) *http.Request {
	req, err := http.NewRequest("", "", nil)
	if err != nil {
		println(err)
	}
	req.Header.Set(name, value)

	return req
}
