package suit

import (
	"encoding/json"
	"github.com/YReshetko/rest.int.test/util"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type requestSender byte
type assertionStatus bool

const (
	CURL requestSender = iota
	GO
)
const (
	OK     assertionStatus = true
	FAILED assertionStatus = false
)

var requestSenderStr = map[string]requestSender{
	"CURL": CURL,
	"GO":   GO,
}

type Suit struct {
	Description string                 `json:"description,omitempty"`
	Vars        map[string]interface{} `json:"vars,omitempty"`
	Executor    requestSender          `json:"executor"`
	Tests       []Test                 `json:"tests"`
}
type Test struct {
	Lable      string      `json:"label"`
	Command    string      `json:"command"`
	Extracts   []Extract   `json:"extracts"`
	Assertions []Assertion `json:"assertions"`
}
type Extract struct {
	Header   string `json:"header,omitempty"`
	Body     string `json:"body,omitempty"`
	Variable string `json:"var"`
}

type AssertionResult struct {
	Index  int
	Result assertionStatus
	Err    error
}

type TestResult struct {
	Index            int
	Label            string
	AssertionResults []*AssertionResult
	TotalResult      assertionStatus
	ExecutionTime    time.Duration
}

type SuitResult struct {
	Description string
	TestResults []*TestResult
	TotalResult assertionStatus
}

func (status assertionStatus) String() string {
	if status {
		return "OK"
	}
	return "FAILED"
}

type TestRunner interface {
	Run(command string) (head, body map[string]string, executionTime time.Duration)
}

func (a requestSender) MarshalJSON() ([]byte, error) {
	return json.Marshal(a)
}

func (t *requestSender) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*t = requestSenderStr[str]
	return nil
}

func (s requestSender) GetRunner() TestRunner {
	switch s {
	case CURL:
		return curlRunner{}
	case GO:
		panic("GO runner has not been implemented yet")
	default:
		return curlRunner{}
	}
}

func (s Suit) Run() (*SuitResult, error) {
	if s.Tests == nil || len(s.Tests) == 0 {
		return nil, errors.New("Test suit doesn't contain any integration test")
	}

	scope, err := util.Parse(s.Vars)
	if err != nil {
		return nil, err
	}
	scope, err = util.Resolve(scope)
	if err != nil {
		return nil, err
	}
	commandRunner := s.Executor.GetRunner()
	testResults := make([]*TestResult, len(s.Tests))
	suitTotalResult := true
	for i, test := range s.Tests {
		command := filterStringWithTokens(test.Command, scope)
		//log.Println("Run command:", command)
		head, body, duration := commandRunner.Run(command)
		extracts := test.Extracts
		for _, extract := range extracts {
			if head != nil {
				extract.process(scope, head, body)
			}
		}

		assertions := test.Assertions
		assertionResults := make([]*AssertionResult, len(assertions))
		testTotalResult := true
		for i, assertion := range assertions {
			ok, err := assertion.Assert(scope)
			assertionResults[i] = &AssertionResult{
				i + 1, assertionStatus(ok), err,
			}
			testTotalResult = ok && testTotalResult
		}
		testResults[i] = &TestResult{
			Index:            i + 1,
			Label:            test.Lable,
			AssertionResults: assertionResults,
			TotalResult:      assertionStatus(testTotalResult),
			ExecutionTime:    duration,
		}
		suitTotalResult = suitTotalResult && testTotalResult
	}

	return &SuitResult{
		Description: s.Description,
		TestResults: testResults,
		TotalResult: assertionStatus(suitTotalResult),
	}, nil
}

func (e Extract) process(scope, head, body map[string]string) {
	if e.Header != "" {
		value := head[e.Header]
		scope[e.Variable] = value
	} else if e.Body != "" {
		value := body[e.Body]
		scope[e.Variable] = value
		//log.Printf("Extract %s = %s", e.Body, value)
	} else {
		panic("Nothing to extract")
	}
}

func filterStringWithTokens(str string, v map[string]string) string {
	tokens := util.ExtractTokens(str)
	if len(tokens) == 0 {
		return str
	}
	for _, token := range tokens {
		value := util.GetValueByToken(token, v)
		str = strings.Replace(str, token, value, -1)
	}
	return str
}
