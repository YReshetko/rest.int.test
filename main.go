package main

import (
	"flag"
	"github.com/YReshetko/rest.int.test/suite"
)

var (
	debug     bool
	suitesPath string
	filter    string
	html      bool
)

func init() {
	flag.BoolVar(&debug, "debug", false, "Boolean flag to print test cases with more details")
	flag.StringVar(&suitesPath, "suites", "", "Path to root folder with tests")
	flag.StringVar(&filter, "filter", "", "Filter to select files to tests used as array via comma and can be used * as matcher")
	flag.BoolVar(&html, "html", false, "Set the flag if you need HTML report (Not implemented)")
	flag.Parse()
}

func main() {
	suitesIterator := suite.LoadSuites(suitesPath)
	printer := suite.GetPrinter(debug, html)
	suite.Run(suitesIterator, printer)
	//printer.PrintAll(results)
}
