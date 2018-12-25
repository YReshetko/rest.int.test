package suit

import "github.com/YReshetko/rest.int.test/util"

type Result struct {
	SuitResult *SuitResult
	FileName   string
	Err        error
}

func Run(iterator util.SuitIterator, printer util.ResultPrinter) []*Result {
	results := []*Result{}
	for iterator.HasNext() {
		suit, fileName := iterator.Next()
		result, err := suit.Run()
		res := &Result{result, fileName, err}
		results = append(results, res)
		printer.Print(res)
	}
	return results
}
