package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

func httpGet(uri string) string {
	if os.Getenv("http_proxy") != "" {
		proxy, _ := url.Parse(os.Getenv("http_proxy"))
		http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxy)}
	}
	resp, err := http.Get(uri)
	if err != nil {
		error(fmt.Sprintf("Error: %s", err))
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		error(fmt.Sprintf("Error: %s", err))
		os.Exit(1)
	}
	return string(body)
}

func fetchWorker(toFetch chan string) {
	for {
		select {
		case val := <-toFetch:
			// we should parse the url and set our output file+path here
			uri, _ := url.Parse(val)

			// if *trimLeading is greater than zero, we should trim that many leading directories
			odn := ""
			if *trimLeading == 0 {
				odn = fmt.Sprintf("%s%s", uri.Hostname(), filepath.Dir(uri.Path))
			} else {
				_odn := fmt.Sprintf("%s%s", uri.Hostname(), filepath.Dir(uri.Path))
				odn = strings.Join(strings.Split(_odn, "/")[*trimLeading:], "/")
			}
			ofn := fmt.Sprintf("%s/%s", odn, filepath.Base(uri.Path))

			// if the output directory does not exist, we should create it
			_ = os.MkdirAll(odn, 0755)

			// if the file already exists, we should do a HEAD request and compare the content-length
			dload := true
			stat, err := os.Stat(ofn)
			if err == nil {
				resp, _ := http.Head(val)
				if stat.Size() == resp.ContentLength {
					dload = false
					debug(fmt.Sprintf("Skipping: %s", val))
				}
			}

			if dload {
				ofp, _ := os.OpenFile(ofn, os.O_CREATE|os.O_WRONLY, 0644)
				st := time.Now()
				// actually move the bits
				ofp.WriteString(httpGet(val))
				// close the file
				ofp.Close()
				// and analyze the time it took
				et := time.Since(st)
				sz, _ := os.Stat(ofn)
				sp := int64(0)
				if int64(et.Seconds()) >= 1 {
					sp = sz.Size() / int64(et.Seconds())
				} else {
					sp = sz.Size()
				}
				info(fmt.Sprintf("Fetched: %-45s [%-10s @ %s/s]", ofn, humanize.Bytes(uint64(sz.Size())), humanize.Bytes(uint64(sp))))
			}
		}
	}
}
