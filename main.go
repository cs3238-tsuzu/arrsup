package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

const redirect = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="refresh" content="0;URL={{.URL}}"/>
	</head>
</html>
`

func generate(proxy string, passHostHeader bool) http.Handler {
	redirectTemplate := template.Must(template.New("html").Parse(redirect))

	mr := func(resp *http.Response) error {
		req := resp.Request
		origin := req.Header.Get("Origin")

		if req.Method == "GET" && resp.StatusCode/100 == 3 && origin != "" {
			u, err := url.Parse(origin)

			if err != nil {
				return nil
			}
			if u.Host != req.Host {
				fmt.Println(u.Host, req.Host)
				_, err := resp.Location()

				if err != nil {
					return nil
				}

				resp.Status = "OK"
				resp.StatusCode = 200
				if resp.Body != nil {
					resp.Body.Close()
				}
				buf := bytes.NewBuffer(nil)
				redirectTemplate.Execute(buf, map[string]interface{}{"URL": resp.Header.Get("Location")})

				resp.Body = ioutil.NopCloser(buf)
				resp.Header.Del("Location")
				resp.Header.Set("Content-Length", strconv.Itoa(buf.Len()))
				resp.Header.Set("Content-Type", "text/html")

			}
		}

		return nil
	}

	proxyURL, err := url.Parse(proxy)

	if err != nil {
		panic(err)
	}

	director := func(req *http.Request) {
		req.URL.Host = proxyURL.Host
		req.URL.Scheme = proxyURL.Scheme

		if !passHostHeader {
			req.Host = proxyURL.Host
		}
	}

	rp := &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: mr,
		ErrorHandler: func(_ http.ResponseWriter, _ *http.Request, err error) {
			log.Println("reverse proxy error:", err)
		},
	}

	return rp
}

func main() {
	proxy := flag.String("target", "", "target for this proxy")
	listenAddr := flag.String("listen", ":80", "listen address")
	passHostHeader := flag.Bool("pass-host-header", true, "pass host header to target")
	help := flag.Bool("help", false, "show usage")

	flag.Parse()

	if *help {
		flag.Usage()

		return
	}

	panic(http.ListenAndServe(*listenAddr, generate(*proxy, *passHostHeader)))

}
