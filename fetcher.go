package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
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

	client := &http.Client{}
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")

	resp, err := client.Do(req)
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

func httpFetch(uri string, ofn string) {
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

	tmp := make([]byte, 10485760) // buffer size of 10MB
	ofp, _ := os.OpenFile(ofn, os.O_CREATE|os.O_WRONLY, 0644)
	defer ofp.Close()

	for {
		n, err := resp.Body.Read(tmp)
		if err != nil && err != io.EOF {
			error(fmt.Sprintf("Error: %s", err))
			os.Exit(1)
		}
		if n == 0 {
			break
		}
		ofp.Write(tmp[:n])
		time.Sleep(10 * time.Millisecond)
	}

	return
}

func fetchWorker(toFetch chan string) {
	for {
		select {
		case val := <-toFetch:
			// we should parse the url and set our output file+path here
			uri, _ := url.Parse(val)

			// if *trimLeading is greater than zero, we should trim that many leading directories
			odn := ""
			if opts.TrimLeading == 0 {
				odn = fmt.Sprintf("%s%s", uri.Hostname(), filepath.Dir(uri.Path))
			} else {
				_odn := fmt.Sprintf("%s%s", uri.Hostname(), filepath.Dir(uri.Path))
				odn = strings.Join(strings.Split(_odn, "/")[opts.TrimLeading:], "/")
			}
			ofn := fmt.Sprintf("%s/%s", odn, filepath.Base(uri.Path))

			// if the output directory does not exist, we should create it
			_ = os.MkdirAll(odn, 0755)

			// if the file already exists, we should do a HEAD request and compare the content-length
			dload := true
			if _, err := os.Stat(fmt.Sprintf("%s.lock", ofn)); err != nil {
				// if the lock file does not exist, we should create it
				_, err = os.Create(fmt.Sprintf("%s.lock", ofn))

				dfn := uri.String()
				if len(ofn) > 47 {
					dfn = fmt.Sprintf("...%s", ofn[len(ofn)-47:])
				}
				debug(fmt.Sprintf("Fetching: %-50s", dfn))

				if err != nil {
					error(fmt.Sprintf("Error: %s", err))
					os.Exit(1)
				}
				stat, err := os.Stat(ofn)
				if err == nil {
					resp, _ := http.Head(val)
					if stat.Size() == resp.ContentLength {
						dload = false
						debug(fmt.Sprintf("Skipping: %s", val))
						fetched++
					}
				}

				if dload {
					st := time.Now()
					// actually move the bits
					httpFetch(val, ofn)
					r := rand.Intn(5)
					time.Sleep(time.Duration(r) * time.Millisecond)
					// and analyze the time it took
					et := time.Since(st)
					sz, _ := os.Stat(ofn)
					sp := int64(0)
					if int64(et.Seconds()) >= 1 {
						sp = sz.Size() / int64(et.Seconds())
					} else {
						sp = sz.Size()
					}
					// and truncate the beginning of the filename if it's too long for display
					// leaving room for 3 dots
					dfn = ofn
					if len(ofn) > 47 {
						dfn = fmt.Sprintf("...%s", ofn[len(ofn)-47:])
					}
					// and display the results
					info(fmt.Sprintf("Fetched: %-50s [%-10s @ %10s/s]", dfn, humanize.Bytes(uint64(sz.Size())), humanize.Bytes(uint64(sp))))
					fetched++
				}
				// and remove the lock file
				_ = os.Remove(fmt.Sprintf("%s.lock", ofn))
			} else {
				debug(fmt.Sprintf("Skipping: %s", val))
				fetched++
			}
		}
	}
}
