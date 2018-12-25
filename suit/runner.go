package suit

import "github.com/YReshetko/rest.int.test/util"

type Result struct {
	SuitResult *SuitResult
	FileName   string
	Err        error
}

func Run(iterator util.SuitIterator) []*Result {
	results := []*Result{}
	for iterator.HasNext() {
		suit, fileName := iterator.Next()
		result, err := suit.Run()
		results = append(results, &Result{result, fileName, err})
	}
	return results
}
