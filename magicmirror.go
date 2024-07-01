package main

import (
	"flag"
	"fmt"
	"time"
)

var programName = "MagicMirror"
var programVersion = "0.1"
var programPatch = "0"
var programAuthor = "Mike 'Fuzzy' Partin"
var programCopy = "2024"

var noColor = flag.Bool("c", false, "Strip colorized output")
var noOutput = flag.Bool("q", false, "Suppress output")
var showDebug = flag.Bool("d", false, "Show debug messages")
var trimLeading = flag.Int("t", 0, "Trim leading directory components")

func main() {
	var showVersion = flag.Bool("V", false, "Show version")
	var regexPattern = flag.String("r", ".*", "Regex pattern to match")

	flag.Parse()
	urls := flag.Args()

	fmt.Println("")

	if *showVersion {
		info(fmt.Sprintf("%s v%s.%s", programName, programVersion, programPatch))
		info(fmt.Sprintf("Copyright %s by %s", programCopy, programAuthor))
		return
	}

	toParse := make(chan string, 10*(1024*1024))
	toMatch := make(chan string, 10*(1024*1024))
	toFetch := make(chan string, 10*(1024*1024))

	go parseWorker(toParse, toMatch)
	go matchWorker(*regexPattern, toParse, toMatch, toFetch)

	for _, _url := range urls {
		toParse <- _url
	}

	// Wait for all workers to start
	time.Sleep(3 * time.Second)

	// Start one fetch worker
	go fetchWorker(toFetch)

	// Wait for all workers to finish
	for len(toParse) > 0 || len(toMatch) > 0 || len(toFetch) > 0 {
		fmt.Print(fmt.Sprintf("toParse: %-10d || toMatch: %-10d || toFetch: %-10d\r", len(toParse), len(toMatch), len(toFetch)))
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("")

}
