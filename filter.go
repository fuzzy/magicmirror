package main

import (
	"fmt"
	"regexp"
)

func matchWorker(patts []string, toParse chan string, toMatch chan string, toFetch chan string) {
	backMatch := regexp.MustCompile(`^.*\.\./.*$`)
	dirPattern := regexp.MustCompile(`^.*/$`)
	for {
		select {
		case val := <-toMatch:
			for _, patt := range patts {
				if !backMatch.MatchString(val) {
					cusPattern := regexp.MustCompile(patt)
					if dirPattern.MatchString(val) {
						toParse <- val
					} else if cusPattern.MatchString(val) {
						debug(fmt.Sprintf("Matching: %s", val))
						toFetch <- val
						matched++
					}
				}
			}
		}
	}
}
