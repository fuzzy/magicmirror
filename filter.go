package main

import (
	"fmt"
	"regexp"
)

func matchWorker(patt string, toParse chan string, toMatch chan string, toFetch chan string) {
	backMatch := regexp.MustCompile(`^.*\.\./.*$`)
	dirPattern := regexp.MustCompile(`^.*/$`)
	cusPattern := regexp.MustCompile(patt)
	for {
		select {
		case val := <-toMatch:
			if !backMatch.MatchString(val) {
				if dirPattern.MatchString(val) {
					toParse <- val
				} else if cusPattern.MatchString(val) {
					debug(fmt.Sprintf("Matching: %s", val))
					toFetch <- val
				}
			}
		}
	}
}
