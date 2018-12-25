package main

import (
	"flag"
	"github.com/YReshetko/rest.int.test/suit"
	"github.com/YReshetko/rest.int.test/util"
)

var (
	debug     bool
	suitsPath string
	filter    string
	html      bool
)

func init() {
	flag.BoolVar(&debug, "debug", false, "Boolean flag to print test cases with more details")
	flag.StringVar(&suitsPath, "suits", "test", "Path to root folder with tests")
	flag.StringVar(&filter, "filter", "", "Filter to select files to tests used as array via comma and can be used * as matcher")
	flag.BoolVar(&html, "html", false, "Set the flag if you need HTML report (Not implemented)")
	flag.Parse()
}

func main() {
	suitesIterator := util.LoadSuits(suitsPath)
	printer := util.GetPrinter(debug, html)
	suit.Run(suitesIterator, printer)
	//printer.PrintAll(results)
}
