package suite

type Result struct {
	SuiteResult *SuiteResult
	FileName   string
	Err        error
}

func Run(iterator SuiteIterator, printer ResultPrinter) []*Result {
	results := []*Result{}
	for iterator.HasNext() {
		suite, fileName := iterator.Next()
		result, err := suite.Run()
		res := &Result{result, fileName, err}
		results = append(results, res)
		printer.Print(res)
	}
	return results
}
