package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
)

var programName = "MagicMirror"
var programVersion = "0.1"
var programPatch = "0"
var programAuthor = "Mike 'Fuzzy' Partin"
var programCopy = "2024"

var matched = 0
var fetched = 0

var opts struct {
	Verbose     []bool   `short:"v" long:"verbose" description:"Show verbose information (twice for debug)"`
	Quiet       bool     `short:"q" long:"quiet" description:"Suppress output"`
	StripColor  bool     `short:"c" long:"strip-color" description:"Strip colorized output"`
	TrimLeading int      `short:"t" long:"trim-leading" description:"Trim leading directory components"`
	Regex       []string `short:"r" long:"regex" description:"Regex pattern to match (can be specified multiple times)"`
	URLs        []string `short:"u" long:"url" description:"URLs to fetch (can be specified multiple times)"`
	ShowVersion bool     `short:"V" long:"version" description:"Show version information"`
}

func main() {
	_, err := flags.Parse(&opts)

	if err != nil {
		os.Exit(1)
	}

	fmt.Println(opts)

	fmt.Println("")

	if opts.ShowVersion {
		info(fmt.Sprintf("%s v%s.%s", programName, programVersion, programPatch))
		info(fmt.Sprintf("Copyright %s by %s", programCopy, programAuthor))
		return
	}

	toParse := make(chan string, 10*(1024*1024))
	toMatch := make(chan string, 10*(1024*1024))
	toFetch := make(chan string, 10*(1024*1024))

	for i := 0; i < 25; i++ {
		go parseWorker(toParse, toMatch)
		go matchWorker(opts.Regex, toParse, toMatch, toFetch)
	}

	for i := 0; i < 2; i++ {
		go fetchWorker(toFetch)
	}

	for _, _url := range opts.URLs {
		toParse <- _url
	}

	// Wait for all workers to start
	time.Sleep(5 * time.Second)

	// Wait for all workers to finish
	for len(toParse) > 0 || len(toMatch) > 0 || len(toFetch) > 0 || fetched < matched {
		_q := fmt.Sprintf("toParse: %-10d || toMatch: %-10d || toFetch: %-10d", len(toParse), len(toMatch), len(toFetch))
		_c := fmt.Sprintf("(%d/%d %6s%%)", fetched, matched, fmt.Sprintf("%.02f", (float64(fetched)/float64(matched))*float64(100)))
		fmt.Print(fmt.Sprintf("%s %s\r", _q, _c))
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("")

}
