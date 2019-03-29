// Copyright (c) 2019 tsuzu
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"net/http"
	"net/http/httputil"
)

func dumpRequestWithoutBody(req *http.Request) string {
	b, err := httputil.DumpRequest(req, false)

	if err != nil {
		return err.Error()
	}

	return string(b)
}

func dumpResponseWithBody(res *http.Response) string {
	b, err := httputil.DumpResponse(res, true)

	if err != nil {
		return err.Error()
	}

	return string(b)
}
