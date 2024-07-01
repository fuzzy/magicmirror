package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func parseHtml(text string) []string {
	tkn := html.NewTokenizer(strings.NewReader(text))
	retv := []string{}
	var isLi bool
	for {
		tt := tkn.Next()
		switch {
		case tt == html.ErrorToken:
			return retv
		case tt == html.StartTagToken:
			t := tkn.Token()
			isLi = t.Data == "a"
		case tt == html.TextToken:
			t := tkn.Token()
			if isLi {
				retv = append(retv, t.String())
			}
			isLi = false
		}
	}
}

func parseWorker(toParse chan string, toMatch chan string) {
	backMatch := regexp.MustCompile(`^.*\.\./.*$`)
	for {
		select {
		case val := <-toParse:
			_, err := url.Parse(val)
			if err == nil {
				debug(fmt.Sprintf("Parsing: %s", val))
				for _, target := range parseHtml(httpGet(val)) {
					if !backMatch.MatchString(target) {
						toMatch <- fmt.Sprintf("%s%s", val, target)
					}
				}
			}
		}
	}
}
