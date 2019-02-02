package main

import (
	"flag"
	"github.com/YReshetko/rest.int.test/suite"
)

var (
	assertion_details bool
	debug             bool
	suitesPath        string
	filter            string
	html              bool
)

func init() {
	flag.BoolVar(&assertion_details, "assertion-details", false, "Boolean flag to print test cases with more details")
	flag.BoolVar(&debug, "debug", false, "Test tool debug mode")
	flag.StringVar(&suitesPath, "suites", "", "Path to root folder with tests")
	flag.StringVar(&filter, "filter", "", "Filter to select files to tests used as array via comma and can be used * as matcher")
	flag.BoolVar(&html, "html", false, "Set the flag if you need HTML report (Not implemented)")
	flag.Parse()
}

func main() {
	suitesIterator := suite.LoadSuites(suitesPath)
	printer := suite.GetPrinter(assertion_details, html)
	suite.Run(suitesIterator, printer, debug)
	//printer.PrintAll(results)
}
