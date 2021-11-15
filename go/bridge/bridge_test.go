package bridge

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"
)

type HttpHandler struct {
	http.Handler
	t *testing.T
}

func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("X-Foo", "bar")
	w.Header().Add("X-Foo", "baz")
	w.WriteHeader(404)
	w.Write([]byte("test"))
}

func TestServe(t *testing.T) {
	h := &HttpHandler{nil, t}
	req := &Request{
		"test.com",
		"/path?foo=bar",
		"POST",
		map[string]ReqHeaderValues{"Content-Length": ReqHeaderValues{"1"}, "X-Foo": ReqHeaderValues{"bar"}},
		"",
		"a",
	}
	res, err := Serve(h, req)
	if err != nil {
		t.Fail()
	}
	if res.StatusCode != 404 {
		fmt.Printf("status code: %d\n", res.StatusCode)
		fmt.Printf("header: %v\n", res.Headers)
		fmt.Printf("base64 body: %s\n", res.Body)
		t.Fail()
	}

	body, err := base64.StdEncoding.DecodeString(res.Body)
	if string(body) != "test" {

		fmt.Printf("body: %s\n", body)
		t.Fail()
	}
}

func TestParseJsonIntoRequest(t *testing.T) {
	body := `{
		"method": "GET",
		"headers": {
		  "host": "example.com",
		  "x-real-ip": "0.0.0.0",
		  "foo": [
			"bar",
			"baz"
		  ],
		  "x-forwarded-host": "example.com",
		  "accept": "*/*",
		  "x-forwarded-proto": "https",
		  "x-vercel-deployment-url": "example.com",
		  "x-forwarded-for": "0.0.0.0",
		  "user-agent": "curl/7.64.1",
		  "x-vercel-forwarded-for": "0.0.0.0",
		  "x-vercel-id": "dev1::pkmmp-1636755907234-8ee8499e420c"
		},
		"path": "/",
		"host": "example.com"
	  }`
	req, err := ParseJsonIntoRequest(body)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		t.Fail()
	}
	if req.Host != "example.com" {
		t.Errorf("Unexpected host: %s", req.Host)
	}
	if req.Headers["host"][0] != "example.com" {
		t.Errorf("Unexpected header host: %s", req.Headers["host"][0])
	}
	if req.Headers["foo"][0] != "bar" {
		t.Errorf("Unexpected header foo[0]: %s", req.Headers["foo"][0])
	}
	if req.Headers["foo"][1] != "baz" {
		t.Errorf("Unexpected header foo[1]: %s", req.Headers["foo"][1])
	}
}
