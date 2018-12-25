package util

import (
	"github.com/YReshetko/rest.int.test/suit"
	"log"
)

type ResultPrinter interface {
	Print(result *suit.Result)
	PrintAll(results []*suit.Result)
}
type htmlPrinter struct {
}
type consolePrinter struct {
	debug bool
	printPrefix string
}

func GetPrinter(debug, html bool) ResultPrinter {
	if html {
		panic("Html printer has not implemented yet")
	}
	return &consolePrinter{debug, ""}
}

func (p *consolePrinter) PrintAll(results []*suit.Result) {
	for _, r := range results {
		p.tab()
		p.Print(r)
		p.stab()
	}
}

func (p *consolePrinter) Print(result *suit.Result) {
	p.p("--- TEST SUIT file -> ", result.FileName, " ---")
	if result.Err != nil {
		p.p("Error -> ", result.Err)
		return
	}
	suitResult := result.SuitResult
	p.p(suitResult.TotalResult.String())

	testNumber := len(suitResult.TestResults)
	successTestNumber := 0
	if suitResult.TotalResult {
		successTestNumber = testNumber
	} else {
		for _, r := range suitResult.TestResults {
			if r.TotalResult {
				successTestNumber++
			}
		}
	}
	p.pf("Result -> %s, %b/%b\n", suitResult.TotalResult, successTestNumber, testNumber)
	if p.debug {
		p.tab()
		p.printSuitDetails(suitResult.TestResults)
		p.stab()
	}
}

func (p *consolePrinter) printSuitDetails(testResults []*suit.TestResult) {
	for _, r := range testResults {
		p.pf("[%b] %s: %s;\n", r.Index, r.TotalResult, r.Label)
		p.pf("Running time: %b;\n", r.ExecutionTime.Seconds())
		if !r.TotalResult {
			for _, a := range r.AssertionResults {
				p.tab()
				p.pf("[%b] %s", a.Index, a.Result.String())
				if a.Err != nil {
					p.p("Error: ", a.Err.Error())
				}
				p.stab()
			}
		}
	}
}

func (p *consolePrinter)tab()  {
	p.printPrefix = p.printPrefix + "\t"
}

func (p *consolePrinter)stab()  {
	p.printPrefix = p.printPrefix[1:]
}

func (p *consolePrinter)pf(pattern string, values ...interface{})  {
	format := p.printPrefix + pattern
	log.Printf(format, values)
}

func (p *consolePrinter)p(pattern string, values ...interface{})  {
	format := p.printPrefix + pattern
	log.Println(format, values)
}