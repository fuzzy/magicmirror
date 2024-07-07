package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fuzzy/subhuman"
	"github.com/jessevdk/go-flags"
)

var programName = "MagicMirror"
var programVersion = "0.1"
var programPatch = "0"
var programAuthor = "Mike 'Fuzzy' Partin"
var programCopy = "2024"

var matched = 0
var fetched = 0
var totalSize = int64(0)

var opts struct {
	Verbose         []bool   `short:"v" long:"verbose" description:"Show verbose information (twice for debug)"`
	Quiet           bool     `short:"q" long:"quiet" description:"Suppress output"`
	StripColor      bool     `short:"c" long:"strip-color" description:"Strip colorized output"`
	TrimLeading     int      `short:"t" long:"trim-leading" description:"Trim leading directory components"`
	ProcessWorkers  int      `short:"p" long:"process-workers" description:"Number of process workers to start" default:"25"`
	MatchWorkers    int      `short:"m" long:"match-workers" description:"Number of match workers to start" default:"2"`
	DownloadWorkers int      `short:"d" long:"download-workers" description:"Number of download workers to start" default:"2"`
	Regex           []string `short:"r" long:"regex" description:"Regex pattern to match (can be specified multiple times)"`
	URLs            []string `short:"u" long:"url" description:"URLs to fetch (can be specified multiple times)"`
	ShowVersion     bool     `short:"V" long:"version" description:"Show version information"`
}

func main() {
	_, err := flags.Parse(&opts)

	if err != nil {
		os.Exit(1)
	}

	if opts.ShowVersion {
		info(fmt.Sprintf("%s v%s.%s", programName, programVersion, programPatch))
		info(fmt.Sprintf("Copyright %s by %s", programCopy, programAuthor))
		return
	}

	toParse := make(chan string, 25*(1024*1024))
	toMatch := make(chan string, 25*(1024*1024))
	toFetch := make(chan string, 25*(1024*1024))

	for i := 0; i < opts.ProcessWorkers; i++ {
		go parseWorker(toParse, toMatch)
	}

	for i := 0; i < opts.MatchWorkers; i++ {
		go matchWorker(opts.Regex, toParse, toMatch, toFetch)
	}

	for i := 0; i < opts.DownloadWorkers; i++ {
		go fetchWorker(toFetch)
	}

	// record start time
	start := time.Now().Unix()

	// send all URLs to the toParse channel
	for _, _url := range opts.URLs {
		toParse <- _url
	}

	// Wait for all workers to start
	if !opts.Quiet {
		fmt.Println("Waiting for workers to start...")
	}

	time.Sleep(15 * time.Second)

	// Wait for all workers to finish
	for len(toParse) > 0 || len(toMatch) > 0 || len(toFetch) > 0 || fetched < matched {
		if !opts.Quiet {
			_toParse := len(toParse)
			_toMatch := len(toMatch)
			_toFetch := len(toFetch)
			_partOne := fmt.Sprintf("toParse: %-9d || toMatch: %-9d || toFetch: %-9d", _toParse, _toMatch, _toFetch)
			_percent := float64(fetched) / float64(matched) * float64(100)
			_partTwo := fmt.Sprintf("(%d/%d %6s%%", fetched, matched, fmt.Sprintf("%.02f", _percent))
			_totalSize := subhuman.HumanSize(totalSize)
			_speed := subhuman.HumanSize(int64(float64(totalSize) / (float64(time.Now().Unix() - start))))
			_partThree := fmt.Sprintf("|| %s in %s @ %s)", _totalSize, subhuman.HumanTimeColon(time.Now().Unix()-start), _speed)
			fmt.Print(fmt.Sprintf("%s %s %s\r", _partOne, _partTwo, _partThree))
		}
		time.Sleep(10 * time.Millisecond)
	}
	if !opts.Quiet {
		_toParse := len(toParse)
		_toMatch := len(toMatch)
		_toFetch := len(toFetch)
		_partOne := fmt.Sprintf("toParse: %-9d || toMatch: %-9d || toFetch: %-9d", _toParse, _toMatch, _toFetch)
		_percent := float64(fetched) / float64(matched) * float64(100)
		_partTwo := fmt.Sprintf("(%d/%d %6s%%", fetched, matched, fmt.Sprintf("%.02f", _percent))
		_totalSize := subhuman.HumanSize(totalSize)
		_speed := subhuman.HumanSize(int64(float64(totalSize) / (float64(time.Now().Unix() - start))))
		_partThree := fmt.Sprintf("|| %s in %s @ %s)", _totalSize, subhuman.HumanTimeColon(time.Now().Unix()-start), _speed)
		fmt.Print(fmt.Sprintf("%s %s %s\r", _partOne, _partTwo, _partThree))
	}
}
