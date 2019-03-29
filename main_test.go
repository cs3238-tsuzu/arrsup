// Copyright (c) 2019 tsuzu
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func extractHost(u string) string {
	parsed, err := url.Parse(u)

	if err != nil {
		panic(err)
	}

	return parsed.Host
}

func TestRedirected(t *testing.T) {
	parent := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Println("request for parent", dumpRequestWithoutBody(req))
		rw.Header().Set("Location", `/"hello"/hoge`)
		rw.WriteHeader(http.StatusFound)
	}))
	defer parent.Close()

	handler := generate(parent.URL)

	proxy := httptest.NewServer(handler)
	defer proxy.Close()

	req, _ := http.NewRequest("GET", proxy.URL, nil)

	req.Header.Set("Host", extractHost(parent.URL))
	req.Header.Set("Origin", "http://example.com/")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(dumpResponseWithBody(resp))

	if resp.StatusCode != 200 {
		t.Errorf("status code must be 200, but got %d", resp.StatusCode)
	}

	if ct := resp.Header.Get("Content-Type"); ct != "text/html" {
		t.Errorf("Content-Type must be text/html, but got %s", ct)
	}
}

func TestPassthrough(t *testing.T) {
	parent := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Println("request for parent", dumpRequestWithoutBody(req))
		rw.Header().Set("Location", `/"hello"/hoge`)
		rw.WriteHeader(http.StatusFound)
	}))
	defer parent.Close()
	t.Log("parent: ", parent.URL)

	handler := generate(parent.URL)

	proxy := httptest.NewServer(handler)
	defer proxy.Close()
	t.Log("proxy: ", proxy.URL)

	req, _ := http.NewRequest("GET", proxy.URL, nil)

	req.Header.Set("Host", extractHost(parent.URL))
	req.Header.Set("Origin", parent.URL)

	client := http.Client{}

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(dumpResponseWithBody(resp))

	if resp.StatusCode != http.StatusFound {
		t.Errorf("status code must be 200, but got %d", resp.StatusCode)
	}

	if ct := resp.Header.Get("Location"); ct == "" {
		t.Errorf("Content-Type must be text/html, but got %s", ct)
	}
}
