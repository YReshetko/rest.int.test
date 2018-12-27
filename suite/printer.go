package suite

import (
	"log"
)

type ResultPrinter interface {
	Print(result *Result)
	PrintAll(results []*Result)
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

func (p *consolePrinter) PrintAll(results []*Result) {
	for _, r := range results {
		p.tab()
		p.Print(r)
		p.stab()
	}
}

func (p *consolePrinter) Print(result *Result) {
	p.p("--- TEST SUITE file -> ", result.FileName, " ---")
	if result.Err != nil {
		p.p("Error -> ", result.Err)
		return
	}
	suiteResult := result.SuiteResult
	p.p(suiteResult.Description)

	testNumber := len(suiteResult.TestResults)
	successTestNumber := 0
	if suiteResult.TotalResult {
		successTestNumber = testNumber
	} else {
		for _, r := range suiteResult.TestResults {
			if r.TotalResult {
				successTestNumber++
			}
		}
	}
	p.pf("Result -> %s, %d/%d\n", suiteResult.TotalResult.String(), successTestNumber, testNumber)
	if p.debug {
		p.tab()
		p.printSuiteDetails(suiteResult.TestResults)
		p.stab()
	}
}

func (p *consolePrinter) printSuiteDetails(testResults []*TestResult) {
	for _, r := range testResults {
		p.pf("[SUITE - %d] %s: %s;\n", r.Index, r.TotalResult.String(), r.Label)
		p.pf("Running time: %fsec;\n", r.ExecutionTime.Seconds())
		if !r.TotalResult {
			for _, a := range r.AssertionResults {
				p.tab()
				p.pf("[ASSERT - %d] %s", a.Index, a.Result.String())
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
	log.Printf(format, values...)
}

func (p *consolePrinter)p(pattern string, values ...interface{})  {
	format := p.printPrefix + pattern
	out := []interface{}{format}
	if len(values) > 0 {
		out = append(out, values...)
	}
	log.Println(out...)
}